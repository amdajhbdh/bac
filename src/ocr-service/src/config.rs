//! Configuration Module

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Config {
    pub workers: usize,
    pub max_queue_size: usize,
    pub batch_timeout_ms: u64,
    pub batch_size: usize,
    pub early_exit_confidence: f64,
    pub languages: Vec<String>,
    pub tessdata_path: Option<String>,
    pub host: String,
    pub port: u16,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            workers: 8,
            max_queue_size: 1000,
            batch_timeout_ms: 50,
            batch_size: 10,
            early_exit_confidence: 0.9,
            languages: vec!["fra".to_string(), "eng".to_string(), "ara".to_string()],
            tessdata_path: None,
            host: "127.0.0.1".to_string(),
            port: 50051,
        }
    }
}
