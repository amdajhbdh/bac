//! Embedding Module
//!
//! Generates embeddings via Gemini API.

use crate::client::GeminiClient;
use anyhow::Result;
use serde::{Deserialize, Serialize};

/// Embedding request
#[derive(Debug, Deserialize)]
pub struct EmbedRequest {
    pub text: String,
    pub task_type: Option<String>,
}

/// Embedding result
#[derive(Debug, Serialize, Deserialize)]
pub struct Embedding {
    pub values: Vec<f32>,
    pub model: String,
    pub task_type: Option<String>,
}

/// Generate embeddings
pub async fn embed(
    client: &GeminiClient,
    text: &str,
    task_type: Option<&str>,
) -> Result<Embedding> {
    let values = client.generate_embedding(text).await?;

    Ok(Embedding {
        values,
        model: "gemini-embedding".to_string(),
        task_type: task_type.map(|s| s.to_string()),
    })
}

#[cfg(test)]
mod tests {
    // Embedding tests would require mocking the client
    // which is more involved - tested via integration tests
}
