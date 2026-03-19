//! Content Analysis Module
//!
//! Analyzes educational content and extracts structured information.

use crate::client::GeminiClient;
use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};

/// Analysis request
#[derive(Debug, Deserialize)]
pub struct AnalyzeRequest {
    pub content: String,
    pub subject: Option<String>,
}

/// Analysis result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnalysisResult {
    pub summary: String,
    pub difficulty: String,
    pub key_points: Vec<String>,
    pub prerequisites: Vec<String>,
    pub learning_objectives: Vec<String>,
}

/// Analyze educational content
pub async fn analyze(client: &GeminiClient, content: &str, subject: Option<&str>) -> Result<AnalysisResult> {
    let subject_hint = subject
        .map(|s| format!("Subject: {}. ", s))
        .unwrap_or_default();

    let prompt = format!(
        r#"Analyze the following educational content and provide structured analysis in JSON format.
{subject_hint}
Content:
{content}

Provide a JSON response with this structure:
{{
    "summary": "A brief 2-3 sentence summary of the content",
    "difficulty": "beginner|intermediate|advanced",
    "key_points": ["Key point 1", "Key point 2", ...],
    "prerequisites": ["Prerequisite knowledge if any"],
    "learning_objectives": ["What the learner will understand after studying"]
}}

Respond ONLY with valid JSON, no other text."#
    );

    let response = client.generate_content(&prompt).await?;

    // Parse JSON response
    let result: AnalysisResult = serde_json::from_str(&response)
        .or_else(|_| {
            // Try to extract JSON from response if there's extra text
            extract_json_from_response(&response)
        })
        .context("Failed to parse analysis result")?;

    Ok(result)
}

/// Extract JSON from a response that may contain extra text
pub(crate) fn extract_json_from_response(response: &str) -> Result<AnalysisResult> {
    // Find JSON boundaries
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
    fn test_extract_json() {
        let response = r#"Here's the analysis:
{
    "summary": "Test summary",
    "difficulty": "beginner",
    "key_points": ["Point 1"],
    "prerequisites": [],
    "learning_objectives": ["Objective 1"]
}
"#;
        let result = extract_json_from_response(response);
        assert!(result.is_ok());
    }
}
