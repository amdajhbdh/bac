//! HTTP service for vector-tools

use std::net::SocketAddr;

use axum::{
    extract::{Path, State},
    routing::{delete as delete_route, get, post},
    Json, Router,
};
use tower_http::cors::CorsLayer;
use tracing_subscriber;

use crate::batch::batch_insert;
use crate::client::PgVectorClient;
use crate::delete::{delete_vector, rebuild_index};
use crate::error::VectorError;
use crate::insert::insert;
use crate::models::{
    BatchInsertRequest, DeleteResult, HealthResponse, InsertRequest, InsertResult,
    RebuildRequest, RebuildResult, SearchRequest, SearchResults,
};
use crate::search;

/// Application state shared across handlers
#[derive(Clone)]
pub struct AppState {
    pub client: PgVectorClient,
}

/// Start the vector-tools HTTP service
pub async fn run(port: u16) -> Result<(), VectorError> {
    tracing_subscriber::fmt::init();

    // Load database URL from environment or config
    let db_url = std::env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb?sslmode=require".to_string());

    let client = PgVectorClient::new(&db_url).await?;
    
    // Initialize schema (create tables if not exist)
    client.init_schema().await?;

    let state = AppState { client };

    let cors = CorsLayer::permissive();

    let app = Router::new()
        .route("/health", get(health))
        .route("/search", post(search_handler))
        .route("/insert", post(insert_handler))
        .route("/batch", post(batch_handler))
        .route("/delete/:id", delete_route(delete_handler))
        .route("/rebuild", post(rebuild_handler))
        .layer(cors)
        .with_state(state);

    let addr = SocketAddr::from(([0, 0, 0, 0], port));
    tracing::info!("Starting vector-tools service on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}

// Handler: Health check
async fn health(State(state): State<AppState>) -> Result<Json<HealthResponse>, VectorError> {
    let db_version = state.client.health_check().await?;
    Ok(Json(HealthResponse {
        status: "healthy".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
        database: db_version,
    }))
}

// Handler: Semantic search
async fn search_handler(
    State(state): State<AppState>,
    Json(request): Json<SearchRequest>,
) -> Result<Json<SearchResults>, VectorError> {
    let results = search::search(&state.client, &request).await?;
    Ok(Json(results))
}

// Handler: Insert vector
async fn insert_handler(
    State(state): State<AppState>,
    Json(request): Json<InsertRequest>,
) -> Result<Json<InsertResult>, VectorError> {
    let result = insert(&state.client, &request).await?;
    Ok(Json(result))
}

// Handler: Batch insert
async fn batch_handler(
    State(state): State<AppState>,
    Json(request): Json<BatchInsertRequest>,
) -> Result<Json<crate::models::BatchResult>, VectorError> {
    let result = batch_insert(&state.client, &request).await?;
    Ok(Json(result))
}

// Handler: Delete vector
async fn delete_handler(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> Result<Json<DeleteResult>, VectorError> {
    let result = delete_vector(&state.client, id).await?;
    Ok(Json(result))
}

// Handler: Rebuild HNSW index
async fn rebuild_handler(
    State(state): State<AppState>,
    Json(request): Json<RebuildRequest>,
) -> Result<Json<RebuildResult>, VectorError> {
    let (success, message, duration_ms) =
        rebuild_index(&state.client, request.m, request.ef_construction).await?;
    Ok(Json(RebuildResult {
        success,
        message,
        duration_ms,
    }))
}
