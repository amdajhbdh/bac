//! Gemini API Client
//!
//! Handles all communication with the Gemini API.

use anyhow::{Context, Result};
use reqwest::Client;
use serde::Deserialize;

/// Client for interacting with Gemini API
#[derive(Clone)]
pub struct GeminiClient {
    client: Client,
    api_key: String,
    model: String,
}

impl GeminiClient {
    /// Create a new Gemini client
    pub fn new(api_key: impl Into<String>) -> Self {
        Self {
            client: Client::new(),
            api_key: api_key.into(),
            model: "gemini-2.0-flash".to_string(),
        }
    }

    /// Create a new Gemini client from environment variable
    pub fn from_env() -> Result<Self> {
        let api_key = std::env::var("GEMINI_API_KEY")
            .context("GEMINI_API_KEY environment variable not set")?;
        Ok(Self::new(api_key))
    }

    /// Set the model to use
    pub fn with_model(mut self, model: impl Into<String>) -> Self {
        self.model = model.into();
        self
    }

    /// Generate content using Gemini
    pub async fn generate_content(&self, prompt: &str) -> Result<String> {
        let url = format!(
            "https://generativelanguage.googleapis.com/v1beta/models/{}:generateContent?key={}",
            self.model, self.api_key
        );

        let request_body = serde_json::json!({
            "contents": [{
                "parts": [{
                    "text": prompt
                }]
            }],
            "generationConfig": {
                "temperature": 0.7,
                "maxOutputTokens": 8192
            }
        });

        let response = self.client
            .post(&url)
            .header("Content-Type", "application/json")
            .json(&request_body)
            .send()
            .await
            .context("Failed to send request to Gemini API")?;

        if !response.status().is_success() {
            let status = response.status();
            let body = response.text().await.unwrap_or_default();
            anyhow::bail!("Gemini API error ({}): {}", status, body);
        }

        let response_body: GeminiResponse = response
            .json()
            .await
            .context("Failed to parse Gemini response")?;

        Ok(response_body.extract_text())
    }

    /// Generate embeddings using Gemini
    pub async fn generate_embedding(&self, text: &str) -> Result<Vec<f32>> {
        let url = format!(
            "https://generativelanguage.googleapis.com/v1beta/models/{}:embedContent?key={}",
            self.model, self.api_key
        );

        let request_body = serde_json::json!({
            "content": {
                "parts": [{
                    "text": text
                }]
            }
        });

        let response = self.client
            .post(&url)
            .header("Content-Type", "application/json")
            .json(&request_body)
            .send()
            .await
            .context("Failed to send embedding request to Gemini API")?;

        if !response.status().is_success() {
            let status = response.status();
            let body = response.text().await.unwrap_or_default();
            anyhow::bail!("Gemini API error ({}): {}", status, body);
        }

        let response_body: EmbeddingResponse = response
            .json()
            .await
            .context("Failed to parse embedding response")?;

        Ok(response_body.embedding.values)
    }
}

/// Response structure from Gemini API
#[derive(Debug, Deserialize)]
struct GeminiResponse {
    candidates: Option<Vec<Candidate>>,
}

#[derive(Debug, Deserialize)]
struct Candidate {
    content: Option<Content>,
}

#[derive(Debug, Deserialize)]
struct Content {
    parts: Vec<Part>,
}

#[derive(Debug, Deserialize)]
struct Part {
    text: Option<String>,
}

impl GeminiResponse {
    fn extract_text(&self) -> String {
        self.candidates
            .as_ref()
            .and_then(|c| c.first())
            .and_then(|c| c.content.as_ref())
            .map(|content| {
                content.parts
                    .iter()
                    .filter_map(|p| p.text.clone())
                    .collect::<Vec<_>>()
                    .join("")
            })
            .unwrap_or_default()
    }
}

/// Response structure from embedding API
#[derive(Debug, Deserialize)]
struct EmbeddingResponse {
    embedding: EmbeddingValues,
}

#[derive(Debug, Deserialize)]
struct EmbeddingValues {
    values: Vec<f32>,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_extract_text() {
        let response = GeminiResponse {
            candidates: Some(vec![Candidate {
                content: Some(Content {
                    parts: vec![
                        Part { text: Some("Hello ".to_string()) },
                        Part { text: Some("World".to_string()) },
                    ],
                }),
            }]),
        };
        assert_eq!(response.extract_text(), "Hello World");
    }

    #[test]
    fn test_extract_text_empty() {
        let response = GeminiResponse { candidates: None };
        assert_eq!(response.extract_text(), "");
    }
}
