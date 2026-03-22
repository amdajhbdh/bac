//! Data models for vector-tools

use serde::{Deserialize, Serialize};

/// Search query request
#[derive(Debug, Deserialize)]
pub struct SearchRequest {
    pub query: Vec<f32>,
    pub top_k: Option<usize>,
    pub filters: Option<SearchFilters>,
}

impl Default for SearchRequest {
    fn default() -> Self {
        Self {
            query: Vec::new(),
            top_k: Some(10),
            filters: None,
        }
    }
}

/// Optional filters for search
#[derive(Debug, Deserialize, Default)]
pub struct SearchFilters {
    pub category: Option<String>,
    pub tags: Option<Vec<String>>,
    pub metadata: Option<serde_json::Value>,
}

/// Search result item
#[derive(Debug, Serialize)]
pub struct SearchResult {
    pub id: i64,
    pub content: String,
    pub metadata: Option<serde_json::Value>,
    pub similarity: f64,
}

/// Search results response
#[derive(Debug, Serialize)]
pub struct SearchResults {
    pub results: Vec<SearchResult>,
    pub query_time_ms: Option<u64>,
}

/// Insert vector request
#[derive(Debug, Deserialize)]
pub struct InsertRequest {
    pub id: Option<i64>,
    pub embedding: Vec<f32>,
    pub content: String,
    pub metadata: Option<serde_json::Value>,
}

/// Insert result response
#[derive(Debug, Serialize)]
pub struct InsertResult {
    pub id: i64,
    pub success: bool,
}

/// Batch insert request
#[derive(Debug, Deserialize)]
pub struct BatchInsertRequest {
    pub records: Vec<BatchRecord>,
}

impl BatchInsertRequest {
    pub fn validate(&self) -> Result<(), String> {
        if self.records.is_empty() {
            return Err("Records cannot be empty".to_string());
        }
        if self.records.len() > 1000 {
            return Err("Batch size exceeds maximum of 1000".to_string());
        }
        for (i, record) in self.records.iter().enumerate() {
            if record.embedding.is_empty() {
                return Err(format!("Record {} has empty embedding", i));
            }
            if record.content.is_empty() {
                return Err(format!("Record {} has empty content", i));
            }
        }
        Ok(())
    }
}

/// Single record for batch insert
#[derive(Debug, Deserialize)]
pub struct BatchRecord {
    pub id: Option<i64>,
    pub embedding: Vec<f32>,
    pub content: String,
    pub metadata: Option<serde_json::Value>,
}

/// Batch insert result
#[derive(Debug, Serialize)]
pub struct BatchResult {
    pub inserted: usize,
    pub failed: usize,
    pub ids: Vec<i64>,
    pub errors: Vec<String>,
}

/// Delete request (ID from path)
#[derive(Debug, Deserialize)]
pub struct DeleteRequest {
    pub id: i64,
}

/// Delete result
#[derive(Debug, Serialize)]
pub struct DeleteResult {
    pub id: i64,
    pub success: bool,
}

/// Rebuild index request
#[derive(Debug, Deserialize)]
pub struct RebuildRequest {
    pub m: Option<i32>,
    pub ef_construction: Option<i32>,
}

impl Default for RebuildRequest {
    fn default() -> Self {
        Self {
            m: Some(16),
            ef_construction: Some(64),
        }
    }
}

/// Rebuild result
#[derive(Debug, Serialize)]
pub struct RebuildResult {
    pub success: bool,
    pub message: String,
    pub duration_ms: u64,
}

/// Health check response
#[derive(Debug, Serialize)]
pub struct HealthResponse {
    pub status: String,
    pub version: String,
    pub database: String,
}
