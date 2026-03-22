//! HTTP Service for Graph Tools
//!
//! Provides REST API for knowledge graph operations.

use axum::{
    extract::Query,
    routing::{get, post},
    Json, Router,
};
use std::net::SocketAddr;
use tower_http::cors::CorsLayer;

use crate::canvas;
use crate::export;
use crate::extract;
use crate::models::{CanvasFile, Entities, HealthResponse, VisualizationUrl};
use crate::relationships;
use crate::visualize;

// ============================================================================
// Service Entry Point
// ============================================================================

/// Start the graph-tools HTTP service on port 3005
pub async fn run() {
    tracing_subscriber::fmt::init();

    let cors = CorsLayer::permissive();

    let app = Router::new()
        .route("/health", get(health))
        .route("/canvas", post(generate_canvas))
        .route("/extract", post(extract_entities))
        .route("/relationships", post(build_rel_graph))
        .route("/export", get(export_graph))
        .route("/visualize", post(create_visualization))
        .layer(cors);

    let addr = SocketAddr::from(([0, 0, 0, 0], 3005));
    tracing::info!("Starting graph-tools service on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr)
        .await
        .expect("Failed to bind to port 3005");
    axum::serve(listener, app)
        .await
        .expect("Server error");
}

// ============================================================================
// HTTP Handlers
// ============================================================================

/// Health check endpoint
async fn health() -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "ok".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
        service: "graph-tools".to_string(),
    })
}

/// Request to generate canvas
#[derive(serde::Deserialize)]
pub struct CanvasRequest {
    pub subject: String,
}

/// Generate canvas file
async fn generate_canvas(Json(request): Json<CanvasRequest>) -> Json<CanvasFile> {
    let canvas = canvas::generate_canvas(&request.subject);
    Json(canvas)
}

/// Request to extract entities
#[derive(serde::Deserialize)]
pub struct ExtractRequest {
    pub content: Option<String>,
    pub source: Option<String>,
}

/// Extract entities from notes
async fn extract_entities(Json(request): Json<ExtractRequest>) -> Json<Entities> {
    let entities = if let (Some(content), Some(source)) = (&request.content, &request.source) {
        extract::extract_from_content(content, source)
    } else {
        extract::extract_entities_from_notes()
    };
    Json(entities)
}

/// Build relationships graph
async fn build_rel_graph(Json(_request): Json<CanvasRequest>) -> Json<CanvasFile> {
    let entities = extract::extract_entities_from_notes();
    let canvas = relationships::build_from_entities(&entities);
    Json(canvas)
}

/// Export graph as JSON
async fn export_graph(Query(params): Query<ExportParams>) -> String {
    let subject = params.subject.as_deref().unwrap_or("graph");
    let canvas = relationships::build_relationships();
    export::export_json(&canvas, subject)
}

#[derive(serde::Deserialize)]
struct ExportParams {
    subject: Option<String>,
}

/// Request for visualization
#[derive(serde::Deserialize)]
pub struct VisualizeRequest {
    pub subject: String,
}

/// Generate visualization URL
async fn create_visualization(Json(request): Json<VisualizeRequest>) -> Json<VisualizationUrl> {
    let viz = visualize::generate_visualization(&request.subject);
    Json(viz)
}

// ============================================================================
// Module Exports for Testing
// ============================================================================

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_health() {
        let response = health().await;
        assert_eq!(response.status, "ok");
        assert_eq!(response.service, "graph-tools");
    }

    #[tokio::test]
    async fn test_generate_canvas_endpoint() {
        let request = CanvasRequest {
            subject: "Test Subject".to_string(),
        };
        let response = generate_canvas(Json(request)).await;
        assert!(!response.nodes.is_empty());
        assert_eq!(response.nodes[0].text, "Test Subject");
    }

    #[tokio::test]
    async fn test_extract_entities_endpoint() {
        let request = ExtractRequest {
            content: Some("## Electric Field\nForce:: A push".to_string()),
            source: Some("physics.md".to_string()),
        };
        let response = extract_entities(Json(request)).await;
        assert!(!response.concepts.is_empty() || !response.definitions.is_empty());
    }

    #[tokio::test]
    async fn test_visualize_endpoint() {
        let request = VisualizeRequest {
            subject: "Physics".to_string(),
        };
        let response = create_visualization(Json(request)).await;
        assert!(response.url.contains("mermaid"));
        assert!(response.embed_url.contains("mermaid"));
    }
}
