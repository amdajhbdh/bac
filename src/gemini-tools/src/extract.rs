//! Entity Extraction Module
//!
//! Extracts concepts, formulas, and topics from educational content.

use crate::client::GeminiClient;
use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};

/// Extraction request
#[derive(Debug, Deserialize)]
pub struct ExtractRequest {
    pub content: String,
}

/// Extracted entities
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExtractedEntities {
    pub concepts: Vec<String>,
    pub formulas: Vec<String>,
    pub topics: Vec<String>,
    pub definitions: Vec<Definition>,
    pub examples: Vec<String>,
}

/// A definition extracted from content
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Definition {
    pub term: String,
    pub definition: String,
}

/// Extract entities from content
pub async fn extract(client: &GeminiClient, content: &str) -> Result<ExtractedEntities> {
    let prompt = format!(
        r#"Extract structured entities from the following educational content.

Content:
{content}

Provide a JSON response with this structure:
{{
    "concepts": ["Concept 1", "Concept 2", ...],
    "formulas": ["Formula 1", "Formula 2", ...],
    "topics": ["Topic 1", "Topic 2", ...],
    "definitions": [
        {{"term": "Term", "definition": "Definition"}}
    ],
    "examples": ["Example 1", "Example 2", ...]
}}

Respond ONLY with valid JSON, no other text."#
    );

    let response = client.generate_content(&prompt).await?;

    let entities = serde_json::from_str(&response)
        .or_else(|_| extract_json_from_response(&response))
        .context("Failed to parse extracted entities")?;

    Ok(entities)
}

/// Extract JSON from a response that may contain extra text
pub(crate) fn extract_json_from_response(response: &str) -> Result<ExtractedEntities> {
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
    fn test_extract_json_basic() {
        let response = r#"{
    "concepts": ["Newton's Laws", "Force"],
    "formulas": ["F = ma"],
    "topics": ["Physics", "Mechanics"],
    "definitions": [],
    "examples": []
}"#;
        let result: Result<ExtractedEntities, _> = serde_json::from_str(response);
        assert!(result.is_ok());
    }
}
