//! Sentence Embeddings using Jina AI API

use anyhow::Result;
use async_trait::async_trait;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone)]
pub struct SentenceEmbeddings {
    model: String,
    api_key: Option<String>,
}

#[derive(Debug, Serialize)]
struct EmbedRequest {
    input: Vec<String>,
    model: String,
}

#[derive(Debug, Deserialize)]
struct EmbedResponse {
    data: Vec<EmbedData>,
}

#[derive(Debug, Deserialize)]
struct EmbedData {
    embedding: Vec<f32>,
}

impl SentenceEmbeddings {
    pub fn new(model: &str) -> Self {
        Self {
            model: model.to_string(),
            api_key: std::env::var("JINA_API_KEY").ok(),
        }
    }

    pub async fn embed(&self, text: &str) -> Result<Vec<f32>> {
        // Try local embedding or API
        if let Some(ref api_key) = self.api_key {
            self.embed_api(text, api_key).await
        } else {
            // Return a simple hash-based embedding for demo
            Ok(self.simple_embed(text))
        }
    }

    async fn embed_api(&self, text: &str, api_key: &str) -> Result<Vec<f32>> {
        let client = reqwest::Client::new();
        
        let response = client.post("https://api.jina.ai/v1/embeddings")
            .header("Authorization", format!("Bearer {}", api_key))
            .json(&EmbedRequest {
                input: vec![text.to_string()],
                model: "jina-embeddings-v2-base-en".to_string(),
            })
            .send()
            .await?
            .json::<EmbedResponse>()
            .await?;

        Ok(response.data.first().map(|d| d.embedding.clone()).unwrap_or_default())
    }

    fn simple_embed(&self, text: &str) -> Vec<f32> {
        // Simple hash-based embedding for when no API key is available
        use std::collections::hash_map::DefaultHasher;
        use std::hash::{Hash, Hasher};
        
        let words: Vec<&str> = text.split_whitespace().collect();
        let dim = 384; // MiniLM-L6-v2 dimension
        
        (0..dim)
            .map(|i| {
                let mut hasher = DefaultHasher::new();
                format!("{}_{}", words.get(i % words.len()).unwrap_or(&""), i).hash(&mut hasher);
                ((hasher.finish() % 1000) as f32 / 1000.0) - 0.5
            })
            .collect()
    }
}
