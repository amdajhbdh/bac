//! HTTP Service Module
//!
//! Axum-based HTTP server for gemini-tools API.

use crate::client::GeminiClient;
use crate::{analyze, correct, embed, extract, generate};
use axum::{
    extract::Json,
    http::StatusCode,
    response::IntoResponse,
    routing::{get, post},
    Router,
};
use std::net::SocketAddr;
use std::sync::Arc;
use tower_http::cors::CorsLayer;
use tracing::info;

/// Application state shared across handlers
#[derive(Clone)]
pub struct AppState {
    pub client: Arc<GeminiClient>,
}

/// Create and configure the router
pub fn create_router(state: AppState) -> Router {
    let cors = CorsLayer::permissive();

    Router::new()
        .route("/health", get(health))
        .route("/analyze", post(handle_analyze))
        .route("/extract", post(handle_extract))
        .route("/generate", post(handle_generate))
        .route("/embed", post(handle_embed))
        .route("/correct", post(handle_correct))
        .layer(cors)
        .with_state(state)
}

/// Start the HTTP server
pub async fn run_server(port: u16, client: GeminiClient) -> anyhow::Result<()> {
    let state = AppState {
        client: Arc::new(client),
    };

    let app = create_router(state);

    let addr = SocketAddr::from(([0, 0, 0, 0], port));
    info!("Starting gemini-tools HTTP server on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}

// Health check handler
async fn health() -> &'static str {
    "OK"
}

// Analyze handler
async fn handle_analyze(
    state: axum::extract::State<AppState>,
    Json(req): Json<analyze::AnalyzeRequest>,
) -> Result<Json<analyze::AnalysisResult>, AppError> {
    let result = analyze::analyze(&state.client, &req.content, req.subject.as_deref())
        .await
        .map_err(AppError::Anyhow)?;
    Ok(Json(result))
}

// Extract handler
async fn handle_extract(
    state: axum::extract::State<AppState>,
    Json(req): Json<extract::ExtractRequest>,
) -> Result<Json<extract::ExtractedEntities>, AppError> {
    let result = extract::extract(&state.client, &req.content)
        .await
        .map_err(AppError::Anyhow)?;
    Ok(Json(result))
}

// Generate handler
async fn handle_generate(
    state: axum::extract::State<AppState>,
    Json(req): Json<generate::GenerateRequest>,
) -> Result<Json<generate::GeneratedNote>, AppError> {
    let result = generate::generate(
        &state.client,
        &req.topic,
        req.subject.as_deref(),
        req.format.as_deref(),
    )
    .await
    .map_err(AppError::Anyhow)?;
    Ok(Json(result))
}

// Embed handler
async fn handle_embed(
    state: axum::extract::State<AppState>,
    Json(req): Json<embed::EmbedRequest>,
) -> Result<Json<embed::Embedding>, AppError> {
    let result = embed::embed(&state.client, &req.text, req.task_type.as_deref())
        .await
        .map_err(AppError::Anyhow)?;
    Ok(Json(result))
}

// Correct handler
async fn handle_correct(
    state: axum::extract::State<AppState>,
    Json(req): Json<correct::CorrectRequest>,
) -> Result<Json<correct::CorrectedText>, AppError> {
    let result = correct::correct(&state.client, &req.ocr_text, req.language.as_deref())
        .await
        .map_err(AppError::Anyhow)?;
    Ok(Json(result))
}

/// Application error type
#[derive(Debug)]
pub enum AppError {
    Anyhow(anyhow::Error),
}

impl std::fmt::Display for AppError {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        match self {
            AppError::Anyhow(e) => write!(f, "{}", e),
        }
    }
}

impl IntoResponse for AppError {
    fn into_response(self) -> axum::response::Response {
        let (status, message) = match self {
            AppError::Anyhow(ref e) => {
                tracing::error!("Application error: {}", e);
                (StatusCode::INTERNAL_SERVER_ERROR, e.to_string())
            }
        };

        let body = serde_json::json!({
            "error": message
        });

        (status, Json(body)).into_response()
    }
}

impl From<anyhow::Error> for AppError {
    fn from(err: anyhow::Error) -> Self {
        AppError::Anyhow(err)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_app_error_display() {
        let err = AppError::Anyhow(anyhow::anyhow!("test error"));
        let display = format!("{}", err);
        assert!(display.contains("test error"));
    }
}
