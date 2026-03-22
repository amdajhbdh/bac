//! Note Generation Module
//!
//! Generates markdown notes with LaTeX support.

use crate::client::GeminiClient;
use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};

/// Generation request
#[derive(Debug, Deserialize)]
pub struct GenerateRequest {
    pub topic: String,
    pub subject: Option<String>,
    pub format: Option<String>,
}

/// Generated note
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GeneratedNote {
    pub title: String,
    pub content: String,
    pub format: String,
}

/// Generate markdown notes
pub async fn generate(
    client: &GeminiClient,
    topic: &str,
    subject: Option<&str>,
    format: Option<&str>,
) -> Result<GeneratedNote> {
    let format = format.unwrap_or("markdown");
    let subject_hint = subject
        .map(|s| format!("Subject: {}. ", s))
        .unwrap_or_default();

    let format_instructions = match format {
        "latex" => "Use LaTeX for all formulas and equations. Wrap math in $...$ for inline and $$...$$ for block equations.",
        "markdown" | _ => "Use Markdown with LaTeX support for formulas. Use $...$ for inline math and $$...$$ for block equations.",
    };

    let prompt = format!(
        r#"Generate comprehensive study notes on the following topic.
{subject_hint}

Topic: {topic}

Requirements:
- {format_instructions}
- Include clear sections with headers
- Include examples where appropriate
- Make it educational and well-structured
- Use proper heading hierarchy (## for sections, ### for subsections)

Return a JSON response with this structure:
{{
    "title": "Note Title",
    "content": "The complete markdown note content",
    "format": "{format}"
}}

Respond ONLY with valid JSON, no other text."#
    );

    let response = client.generate_content(&prompt).await?;

    // Try to parse as JSON first
    let note: Result<GeneratedNote, _> = serde_json::from_str(&response)
        .or_else(|_| {
            // If parsing fails, wrap the response as content
            Ok::<GeneratedNote, serde_json::Error>(GeneratedNote {
                title: topic.to_string(),
                content: extract_json_content(&response).unwrap_or(response),
                format: format.to_string(),
            })
        });

    note.context("Failed to parse generated note")
}

/// Extract content from JSON response or return original
fn extract_json_content(response: &str) -> Option<String> {
    if let Ok(note) = serde_json::from_str::<GeneratedNote>(response) {
        return Some(note.content);
    }

    // Try to extract JSON and get content field
    if let Some(start) = response.find('{') {
        if let Some(end) = response.rfind('}') {
            let json_str = &response[start..=end];
            if let Ok(value) = serde_json::from_str::<serde_json::Value>(json_str) {
                return value.get("content")
                    .and_then(|v| v.as_str())
                    .map(|s| s.to_string());
            }
        }
    }
    None
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_extract_content() {
        let response = r#"{
    "title": "Test Title",
    "content": "This is the content",
    "format": "markdown"
}"#;
        let content = extract_json_content(response);
        assert_eq!(content, Some("This is the content".to_string()));
    }
}
