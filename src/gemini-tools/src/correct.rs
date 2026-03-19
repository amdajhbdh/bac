//! OCR Correction Module
//!
//! AI-corrects OCR errors in text.

use crate::client::GeminiClient;
use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};

/// Correction request
#[derive(Debug, Deserialize)]
pub struct CorrectRequest {
    pub ocr_text: String,
    pub language: Option<String>,
}

/// Corrected text result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CorrectedText {
    pub corrected: String,
    pub corrections: Vec<Correction>,
}

/// A single correction made
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Correction {
    pub original: String,
    pub corrected: String,
    pub reason: String,
}

/// AI-correct OCR text
pub async fn correct(
    client: &GeminiClient,
    ocr_text: &str,
    language: Option<&str>,
) -> Result<CorrectedText> {
    let language_hint = language
        .map(|l| format!("The text is in {}. ", l))
        .unwrap_or_else(|| "Detect the language automatically. ".to_string());

    let prompt = format!(
        r#"Correct OCR errors in the following text.
{language_hint}Common OCR errors include:
- Misrecognized characters (0/O, 1/l/I, 5/S, etc.)
- Missing or extra spaces
- Broken words
- Missing punctuation
- Case errors

Text:
{ocr_text}

Provide a JSON response with this structure:
{{
    "corrected": "The fully corrected text",
    "corrections": [
        {{"original": "misread", "corrected": "misread", "reason": "Character recognition error"}}
    ]
}}

Respond ONLY with valid JSON, no other text."#
    );

    let response = client.generate_content(&prompt).await?;

    let corrected = serde_json::from_str(&response)
        .or_else(|_| extract_json_from_response(&response))
        .context("Failed to parse corrected text")?;

    Ok(corrected)
}

/// Extract JSON from response
pub(crate) fn extract_json_from_response(response: &str) -> Result<CorrectedText> {
    let start = match response.find('{') {
        Some(s) => s,
        None => anyhow::bail!("No JSON found in response"),
    };
    let end = match response.rfind('}') {
        Some(e) => e,
        None => anyhow::bail!("No JSON found in response"),
    };
    let json_str = &response[start..=end];
    serde_json::from_str(json_str).context("Failed to parse extracted JSON")
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_corrected() {
        let json = r#"{
    "corrected": "Hello World",
    "corrections": [
        {"original": "He110", "corrected": "Hello", "reason": "Number recognition error"}
    ]
}"#;
        let result: Result<CorrectedText, _> = serde_json::from_str(json);
        assert!(result.is_ok());
        let corrected = result.unwrap();
        assert_eq!(corrected.corrected, "Hello World");
        assert_eq!(corrected.corrections.len(), 1);
    }
}
