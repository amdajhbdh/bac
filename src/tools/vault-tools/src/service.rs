//! HTTP service for vault-tools

use axum::{
    extract::Query,
    routing::{get, post},
    Json, Router,
};
use std::net::SocketAddr;
use tower_http::cors::CorsLayer;

use crate::models::*;
use crate::{link, moc, read, search, write};

/// Start the vault-tools HTTP service
pub async fn run() {
    tracing_subscriber::fmt::init();

    let cors = CorsLayer::permissive();

    let app = Router::new()
        .route("/health", get(health))
        .route("/read", get(read_note))
        .route("/write", post(write_note))
        .route("/link", post(create_link))
        .route("/moc/update", post(update_moc))
        .route("/search", get(search_vault))
        .route("/search/tag", get(search_by_tag))
        .route("/tags", get(list_tags))
        .layer(cors);

    let addr = SocketAddr::from(([0, 0, 0, 0], 3003));
    tracing::info!("Starting vault-tools service on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

// ============================================================================
// HTTP Handlers
// ============================================================================

/// Health check
async fn health() -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "ok".to_string(),
        vault_path: read::get_vault_path().to_string_lossy().to_string(),
    })
}

/// Read note handler
async fn read_note(Query(params): Query<ReadParams>) -> Json<ReadResult> {
    let path = params.path.unwrap_or_default();
    let result = read::read_note(&path).await;
    Json(result)
}

#[derive(serde::Deserialize)]
struct ReadParams {
    path: Option<String>,
}

/// Write note handler
async fn write_note(Json(request): Json<WriteRequest>) -> Json<WriteResult> {
    let result = write::write_note(request).await;
    Json(result)
}

/// Create link handler
async fn create_link(Json(request): Json<LinkRequest>) -> Json<LinkResult> {
    let result = link::create_link(request).await;
    Json(result)
}

/// Update MOC handler
async fn update_moc(Json(request): Json<MocUpdateRequest>) -> Json<MocUpdateResult> {
    let result = moc::update_moc(request).await;
    Json(result)
}

/// Search vault handler
async fn search_vault(Query(params): Query<SearchParams>) -> Json<SearchResults> {
    let query = params.q.unwrap_or_default();
    let result = search::search_vault(&query).await;
    Json(result)
}

#[derive(serde::Deserialize)]
struct SearchParams {
    q: Option<String>,
}

/// Search by tag handler
async fn search_by_tag(Query(params): Query<TagParams>) -> Json<SearchResults> {
    let tag = params.tag.unwrap_or_default();
    let result = search::search_by_tag(&tag).await;
    Json(result)
}

#[derive(serde::Deserialize)]
struct TagParams {
    tag: Option<String>,
}

/// List all tags handler
async fn list_tags() -> Json<Vec<(String, usize)>> {
    let tags = search::get_all_tags().await;
    Json(tags)
}
