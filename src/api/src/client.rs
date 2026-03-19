//! Internal HTTP client for tool service communication
//!
//! Provides a centralized client for proxying requests to gemini-tools,
//! vector-tools, vault-tools, cloud-tools, and graph-tools.

use reqwest::Client;
use std::collections::HashMap;
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::RwLock;
use tracing::{error, info};

/// Default timeout for tool service requests
const DEFAULT_TIMEOUT_SECS: u64 = 30;

/// Configuration for tool service endpoints
#[derive(Debug, Clone)]
pub struct ToolEndpoints {
    /// gemini-tools HTTP address (default: 127.0.0.1:3001)
    pub gemini: String,
    /// vector-tools HTTP address (default: 127.0.0.1:3002)
    pub vector: String,
    /// vault-tools HTTP address (default: 127.0.0.1:3003)
    pub vault: String,
    /// cloud-tools HTTP address (default: 127.0.0.1:3004)
    pub cloud: String,
    /// graph-tools HTTP address (default: 127.0.0.1:3005)
    pub graph: String,
}

impl Default for ToolEndpoints {
    fn default() -> Self {
        Self {
            gemini: "http://127.0.0.1:3001".to_string(),
            vector: "http://127.0.0.1:3002".to_string(),
            vault: "http://127.0.0.1:3003".to_string(),
            cloud: "http://127.0.0.1:3004".to_string(),
            graph: "http://127.0.0.1:3005".to_string(),
        }
    }
}

impl ToolEndpoints {
    /// Create endpoints from environment variables or defaults
    pub fn from_env() -> Self {
        Self {
            gemini: std::env::var("GEMINI_TOOLS_URL")
                .unwrap_or_else(|_| "http://127.0.0.1:3001".to_string()),
            vector: std::env::var("VECTOR_TOOLS_URL")
                .unwrap_or_else(|_| "http://127.0.0.1:3002".to_string()),
            vault: std::env::var("VAULT_TOOLS_URL")
                .unwrap_or_else(|_| "http://127.0.0.1:3003".to_string()),
            cloud: std::env::var("CLOUD_TOOLS_URL")
                .unwrap_or_else(|_| "http://127.0.0.1:3004".to_string()),
            graph: std::env::var("GRAPH_TOOLS_URL")
                .unwrap_or_else(|_| "http://127.0.0.1:3005".to_string()),
        }
    }
}

/// Shared tool client state with connection pooling
#[derive(Clone)]
pub struct ToolClient {
    /// HTTP client with connection pooling
    http: Client,
    /// Tool service endpoints
    endpoints: ToolEndpoints,
    /// Request counter per service (for metrics/debugging)
    counters: Arc<RwLock<HashMap<String, u64>>>,
}

impl ToolClient {
    /// Create a new tool client with connection pooling
    pub fn new(endpoints: ToolEndpoints) -> Self {
        let http = Client::builder()
            .pool_max_idle_per_host(5)
            .pool_idle_timeout(Duration::from_secs(60))
            .timeout(Duration::from_secs(DEFAULT_TIMEOUT_SECS))
            .build()
            .expect("Failed to create HTTP client");

        info!(
            "ToolClient initialized with endpoints: gemini={}, vector={}, vault={}, cloud={}, graph={}",
            endpoints.gemini, endpoints.vector, endpoints.vault, endpoints.cloud, endpoints.graph
        );

        Self {
            http,
            endpoints,
            counters: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    /// Create with default endpoints
    pub fn with_defaults() -> Self {
        Self::new(ToolEndpoints::default())
    }

    /// Create from environment variables
    pub fn from_env() -> Self {
        Self::new(ToolEndpoints::from_env())
    }

    /// Increment and get request counter for a service
    async fn increment_counter(&self, service: &str) -> u64 {
        let mut counters = self.counters.write().await;
        let count = counters.entry(service.to_string()).or_insert(0);
        *count += 1;
        *count
    }

    /// Get current counter for a service
    pub async fn get_counter(&self, service: &str) -> u64 {
        let counters = self.counters.read().await;
        counters.get(service).copied().unwrap_or(0)
    }

    /// Proxy a request to gemini-tools
    pub async fn proxy_gemini(
        &self,
        action: &str,
        body: serde_json::Value,
    ) -> Result<reqwest::Response, ToolClientError> {
        let counter = self.increment_counter("gemini").await;
        let url = format!("{}/{}", self.endpoints.gemini, action);
        
        tracing::debug!(target: "tool_proxy", service = "gemini", action = action, request_num = counter, url = %url);
        
        self.http
            .post(&url)
            .json(&body)
            .send()
            .await
            .map_err(|e| {
                error!(service = "gemini", action = action, error = %e, "Failed to proxy to gemini-tools");
                ToolClientError::ConnectionFailed(e.to_string())
            })
    }

    /// Proxy a request to vector-tools
    pub async fn proxy_vector(
        &self,
        action: &str,
        body: serde_json::Value,
    ) -> Result<reqwest::Response, ToolClientError> {
        let counter = self.increment_counter("vector").await;
        let url = format!("{}/{}", self.endpoints.vector, action);
        
        tracing::debug!(target: "tool_proxy", service = "vector", action = action, request_num = counter, url = %url);
        
        self.http
            .post(&url)
            .json(&body)
            .send()
            .await
            .map_err(|e| {
                error!(service = "vector", action = action, error = %e, "Failed to proxy to vector-tools");
                ToolClientError::ConnectionFailed(e.to_string())
            })
    }

    /// Proxy a request to vault-tools
    pub async fn proxy_vault(
        &self,
        action: &str,
        body: serde_json::Value,
    ) -> Result<reqwest::Response, ToolClientError> {
        let counter = self.increment_counter("vault").await;
        let url = format!("{}/{}", self.endpoints.vault, action);
        
        tracing::debug!(target: "tool_proxy", service = "vault", action = action, request_num = counter, url = %url);
        
        self.http
            .post(&url)
            .json(&body)
            .send()
            .await
            .map_err(|e| {
                error!(service = "vault", action = action, error = %e, "Failed to proxy to vault-tools");
                ToolClientError::ConnectionFailed(e.to_string())
            })
    }

    /// Proxy a GET request to vault-tools (for read/search operations)
    pub async fn proxy_vault_get(
        &self,
        action: &str,
    ) -> Result<reqwest::Response, ToolClientError> {
        let counter = self.increment_counter("vault").await;
        let url = format!("{}/{}", self.endpoints.vault, action);
        
        tracing::debug!(target: "tool_proxy", service = "vault", method = "GET", action = action, request_num = counter, url = %url);
        
        self.http
            .get(&url)
            .send()
            .await
            .map_err(|e| {
                error!(service = "vault", action = action, error = %e, "Failed to proxy to vault-tools");
                ToolClientError::ConnectionFailed(e.to_string())
            })
    }

    /// Proxy a request to cloud-tools
    pub async fn proxy_cloud(
        &self,
        action: &str,
        body: serde_json::Value,
    ) -> Result<reqwest::Response, ToolClientError> {
        let counter = self.increment_counter("cloud").await;
        let url = format!("{}/{}", self.endpoints.cloud, action);
        
        tracing::debug!(target: "tool_proxy", service = "cloud", action = action, request_num = counter, url = %url);
        
        self.http
            .post(&url)
            .json(&body)
            .send()
            .await
            .map_err(|e| {
                error!(service = "cloud", action = action, error = %e, "Failed to proxy to cloud-tools");
                ToolClientError::ConnectionFailed(e.to_string())
            })
    }

    /// Proxy a request to graph-tools
    pub async fn proxy_graph(
        &self,
        action: &str,
        body: serde_json::Value,
    ) -> Result<reqwest::Response, ToolClientError> {
        let counter = self.increment_counter("graph").await;
        let url = format!("{}/{}", self.endpoints.graph, action);
        
        tracing::debug!(target = "tool_proxy", service = "graph", action = action, request_num = counter, url = %url);
        
        self.http
            .post(&url)
            .json(&body)
            .send()
            .await
            .map_err(|e| {
                error!(service = "graph", action = action, error = %e, "Failed to proxy to graph-tools");
                ToolClientError::ConnectionFailed(e.to_string())
            })
    }

    /// Get endpoints configuration
    pub fn endpoints(&self) -> &ToolEndpoints {
        &self.endpoints
    }
}

/// Error types for tool client operations
#[derive(Debug, thiserror::Error)]
pub enum ToolClientError {
    #[error("Failed to connect to tool service: {0}")]
    ConnectionFailed(String),
    
    #[error("Tool service returned error: {0}")]
    ServiceError(String),
    
    #[error("Invalid response from tool service")]
    InvalidResponse,
    
    #[error("Request timeout")]
    Timeout,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_tool_endpoints_default() {
        let endpoints = ToolEndpoints::default();
        assert_eq!(endpoints.gemini, "http://127.0.0.1:3001");
        assert_eq!(endpoints.vector, "http://127.0.0.1:3002");
        assert_eq!(endpoints.vault, "http://127.0.0.1:3003");
    }

    #[test]
    fn test_tool_endpoints_custom() {
        let endpoints = ToolEndpoints {
            gemini: "http://localhost:9001".to_string(),
            vector: "http://localhost:9002".to_string(),
            vault: "http://localhost:9003".to_string(),
            cloud: "http://localhost:9004".to_string(),
            graph: "http://localhost:9005".to_string(),
        };
        assert_eq!(endpoints.gemini, "http://localhost:9001");
    }

    #[tokio::test]
    async fn test_tool_client_counter() {
        let client = ToolClient::with_defaults();
        assert_eq!(client.get_counter("gemini").await, 0);
        
        // Increment counters
        let _ = client.proxy_gemini("test", serde_json::json!({})).await;
        let _ = client.proxy_gemini("test2", serde_json::json!({})).await;
        
        assert!(client.get_counter("gemini").await >= 2);
    }
}
