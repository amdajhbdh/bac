//! Proxy route handlers for tool service integration
//!
//! Proxies API requests to internal tool services:
//! - gemini-tools (:3001) - AI/analyze, extract, generate, embed, correct
//! - vector-tools (:3002) - Semantic search, insert, batch, delete
//! - vault-tools (:3003) - Obsidian vault read/write/link/search
//! - cloud-tools (:3004) - Cloud shell integration
//! - graph-tools (:3005) - Knowledge graph operations

use axum::{
    extract::{Path, State},
    http::StatusCode,
    response::{IntoResponse, Response},
    routing::{get, post},
    Json, Router,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tower_http::cors::CorsLayer;
use tracing::{error, info, warn};

use crate::client::{ToolClient, ToolClientError};

/// Application state for proxy routes
#[derive(Clone)]
pub struct ProxyState {
    pub client: Arc<ToolClient>,
}

impl ProxyState {
    pub fn new(client: ToolClient) -> Self {
        Self {
            client: Arc::new(client),
        }
    }
}

/// Create proxy routes router
pub fn routes(state: ProxyState) -> Router {
    let cors = CorsLayer::permissive();

    Router::new()
        // Gemini tools proxy routes
        .route("/api/gemini/:action", post(gemini_proxy))
        // Vector tools proxy routes
        .route("/api/vector/:action", post(vector_proxy))
        // Vault tools proxy routes
        .route("/api/vault/:action", post(vault_proxy_post))
        .route("/api/vault/:action", get(vault_proxy_get))
        // Cloud tools proxy routes
        .route("/api/cloud/:action", post(cloud_proxy))
        // Graph tools proxy routes
        .route("/api/graph/:action", post(graph_proxy))
        // Health/status endpoint for tool services
        .route("/api/tools/status", get(tools_status))
        .layer(cors)
        .with_state(state)
}

// ============================================================================
// Gemini Tools Proxy Handlers
// ============================================================================

/// Supported gemini actions
#[derive(Debug, Deserialize)]
pub struct GeminiParams {
    pub action: String,
}

/// Proxy to gemini-tools service
///
/// Routes:
/// - POST /api/gemini/analyze - Analyze content with AI
/// - POST /api/gemini/extract - Extract entities from content
/// - POST /api/gemini/generate - Generate study notes
/// - POST /api/gemini/embed - Create embeddings
/// - POST /api/gemini/correct - OCR text correction
async fn gemini_proxy(
    State(state): State<ProxyState>,
    Path(action): Path<String>,
    Json(body): Json<serde_json::Value>,
) -> Result<Response, ProxyError> {
    info!(service = "gemini", action = %action, "Proxying request to gemini-tools");

    let response = state
        .client
        .proxy_gemini(&action, body)
        .await
        .map_err(ProxyError::ToolClient)?;

    proxy_response(response).await
}

// ============================================================================
// Vector Tools Proxy Handlers
// ============================================================================

/// Proxy to vector-tools service
///
/// Routes:
/// - POST /api/vector/search - Semantic search
/// - POST /api/vector/insert - Insert vector
/// - POST /api/vector/batch - Batch insert
/// - DELETE /api/vector/:id - Delete vector
/// - POST /api/vector/rebuild - Rebuild HNSW index
async fn vector_proxy(
    State(state): State<ProxyState>,
    Path(action): Path<String>,
    Json(body): Json<serde_json::Value>,
) -> Result<Response, ProxyError> {
    info!(service = "vector", action = %action, "Proxying request to vector-tools");

    let response = state
        .client
        .proxy_vector(&action, body)
        .await
        .map_err(ProxyError::ToolClient)?;

    proxy_response(response).await
}

// ============================================================================
// Vault Tools Proxy Handlers
// ============================================================================

/// Proxy POST requests to vault-tools
///
/// Routes:
/// - POST /api/vault/write - Write note
/// - POST /api/vault/link - Create link
/// - POST /api/vault/moc/update - Update MOC
async fn vault_proxy_post(
    State(state): State<ProxyState>,
    Path(action): Path<String>,
    Json(body): Json<serde_json::Value>,
) -> Result<Response, ProxyError> {
    info!(service = "vault", action = %action, method = "POST", "Proxying request to vault-tools");

    let response = state
        .client
        .proxy_vault(&action, body)
        .await
        .map_err(ProxyError::ToolClient)?;

    proxy_response(response).await
}

/// Proxy GET requests to vault-tools
///
/// Routes:
/// - GET /api/vault/read - Read note
/// - GET /api/vault/search - Search vault
/// - GET /api/vault/search/tag - Search by tag
/// - GET /api/vault/tags - List all tags
async fn vault_proxy_get(
    State(state): State<ProxyState>,
    Path(action): Path<String>,
) -> Result<Response, ProxyError> {
    info!(service = "vault", action = %action, method = "GET", "Proxying request to vault-tools");

    let response = state
        .client
        .proxy_vault_get(&action)
        .await
        .map_err(ProxyError::ToolClient)?;

    proxy_response(response).await
}

// ============================================================================
// Cloud Tools Proxy Handlers
// ============================================================================

/// Proxy to cloud-tools service
async fn cloud_proxy(
    State(state): State<ProxyState>,
    Path(action): Path<String>,
    Json(body): Json<serde_json::Value>,
) -> Result<Response, ProxyError> {
    info!(service = "cloud", action = %action, "Proxying request to cloud-tools");

    let response = state
        .client
        .proxy_cloud(&action, body)
        .await
        .map_err(ProxyError::ToolClient)?;

    proxy_response(response).await
}

// ============================================================================
// Graph Tools Proxy Handlers
// ============================================================================

/// Proxy to graph-tools service
async fn graph_proxy(
    State(state): State<ProxyState>,
    Path(action): Path<String>,
    Json(body): Json<serde_json::Value>,
) -> Result<Response, ProxyError> {
    info!(service = "graph", action = %action, "Proxying request to graph-tools");

    let response = state
        .client
        .proxy_graph(&action, body)
        .await
        .map_err(ProxyError::ToolClient)?;

    proxy_response(response).await
}

// ============================================================================
// Status Endpoint
// ============================================================================

/// Tool service status response
#[derive(Debug, Serialize)]
pub struct ToolsStatusResponse {
    pub api: ServiceStatus,
    pub gemini: ServiceStatus,
    pub vector: ServiceStatus,
    pub vault: ServiceStatus,
    pub cloud: ServiceStatus,
    pub graph: ServiceStatus,
}

/// Individual service status
#[derive(Debug, Serialize)]
pub struct ServiceStatus {
    pub url: String,
    pub status: String,
}

/// Get status of all tool services
async fn tools_status(State(state): State<ProxyState>) -> Json<ToolsStatusResponse> {
    let endpoints = state.client.endpoints();

    // For now, just return configured endpoints
    // Could be extended to actually ping each service
    Json(ToolsStatusResponse {
        api: ServiceStatus {
            url: "http://0.0.0.0:8080".to_string(),
            status: "running".to_string(),
        },
        gemini: ServiceStatus {
            url: endpoints.gemini.clone(),
            status: "configured".to_string(),
        },
        vector: ServiceStatus {
            url: endpoints.vector.clone(),
            status: "configured".to_string(),
        },
        vault: ServiceStatus {
            url: endpoints.vault.clone(),
            status: "configured".to_string(),
        },
        cloud: ServiceStatus {
            url: endpoints.cloud.clone(),
            status: "configured".to_string(),
        },
        graph: ServiceStatus {
            url: endpoints.graph.clone(),
            status: "configured".to_string(),
        },
    })
}

// ============================================================================
// Helper Functions
// ============================================================================

/// Convert tool service response to API response
async fn proxy_response(response: reqwest::Response) -> Result<Response, ProxyError> {
    let status = response.status();
    let content_type = response.headers().get("content-type").cloned();
    let bytes = response
        .bytes()
        .await
        .map_err(|_| ProxyError::InvalidResponse)?;

    // Clone status for logging
    let status_code = status.as_u16();
    if status.is_success() {
        tracing::debug!(status = status_code, "Tool service responded successfully");
    } else {
        warn!(status = status_code, "Tool service returned non-success status");
    }

    let mut builder = Response::builder().status(status);
    
    // Preserve content-type if present
    if let Some(content_type) = content_type {
        builder = builder.header("content-type", content_type);
    }

    Ok(builder.body(axum::body::Body::from(bytes)).unwrap())
}

// ============================================================================
// Error Types
// ============================================================================

/// Proxy-specific error types
#[derive(Debug)]
pub enum ProxyError {
    ToolClient(ToolClientError),
    InvalidResponse,
    ServiceError(String),
}

impl IntoResponse for ProxyError {
    fn into_response(self) -> Response {
        let (status, message) = match &self {
            ProxyError::ToolClient(e) => {
                error!(error = %e, "Tool client error");
                (
                    StatusCode::BAD_GATEWAY,
                    format!("Tool service unavailable: {}", e),
                )
            }
            ProxyError::InvalidResponse => {
                error!("Invalid response from tool service");
                (
                    StatusCode::BAD_GATEWAY,
                    "Invalid response from tool service".to_string(),
                )
            }
            ProxyError::ServiceError(msg) => {
                error!(error = msg, "Tool service error");
                (StatusCode::BAD_GATEWAY, msg.clone())
            }
        };

        let body = serde_json::json!({
            "error": message,
            "status": "proxy_error"
        });

        (status, Json(body)).into_response()
    }
}

impl From<ToolClientError> for ProxyError {
    fn from(err: ToolClientError) -> Self {
        ProxyError::ToolClient(err)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_proxy_state_creation() {
        let client = ToolClient::with_defaults();
        let state = ProxyState::new(client);
        assert!(state.client.endpoints().gemini.contains("3001"));
    }

    #[tokio::test]
    async fn test_tools_status_response() {
        let client = ToolClient::with_defaults();
        let state = ProxyState::new(client);
        let response = tools_status(State(state)).await;
        
        assert_eq!(response.api.status, "running");
        assert!(response.gemini.url.contains("3001"));
    }
}
