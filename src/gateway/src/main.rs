//! BAC Gateway - REST API with OCR and Animation integration

mod animation;
mod rag;

use std::net::SocketAddr;
use std::sync::Arc;

use animation::{AnimationQueue, AnimationRequest, AnimationResponse, JobStatus, ManimBridge, ManimConfig};
use rag::{ChatMode, ChatRequest, ChatResponse, RAGEngine};
use axum::{
    extract::State,
    http::StatusCode,
    response::Json,
    routing::{get, post},
    Router,
};
use reqwest::Client;
use serde::{Deserialize, Serialize};
use tower::ServiceBuilder;
use tower_http::cors::{Any, CorsLayer};
use tower_http::trace::TraceLayer;

#[derive(Clone)]
pub struct AppState {
    pub ocr_client: Arc<OCRClient>,
    pub animation_bridge: Arc<ManimBridge>,
    pub rag_engine: Arc<RAGEngine>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct OCRResponse {
    pub success: bool,
    pub data: Option<OCRData>,
    pub error: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct OCRData {
    pub text: String,
    pub confidence: f64,
    pub engine: String,
    pub processing_time_ms: i64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HealthResponse {
    pub status: String,
    pub gateway_version: String,
    pub ocr_service: String,
    pub animation_service: String,
}

pub struct OCRClient {
    client: Client,
    ocr_url: String,
}

impl OCRClient {
    pub fn new(ocr_url: String) -> Self {
        Self {
            client: Client::new(),
            ocr_url,
        }
    }

    pub async fn process_image(&self, image_data: Vec<u8>, file_name: &str) -> Result<OCRResponse, String> {
        let form = reqwest::multipart::Form::new()
            .part("file", reqwest::multipart::Part::bytes(image_data)
                .file_name(file_name.to_string())
                .mime_str("image/png").map_err(|e| e.to_string())?);

        let response = self.client
            .post(format!("{}/ocr", self.ocr_url))
            .multipart(form)
            .send()
            .await
            .map_err(|e| e.to_string())?;

        let ocr_response: OCRResponse = response.json().await.map_err(|e| e.to_string())?;
        Ok(ocr_response)
    }

    pub async fn process_pdf(&self, pdf_data: Vec<u8>) -> Result<OCRResponse, String> {
        let form = reqwest::multipart::Form::new()
            .part("file", reqwest::multipart::Part::bytes(pdf_data)
                .file_name("document.pdf")
                .mime_str("application/pdf").map_err(|e| e.to_string())?);

        let response = self.client
            .post(format!("{}/pdf", self.ocr_url))
            .multipart(form)
            .send()
            .await
            .map_err(|e| e.to_string())?;

        let ocr_response: OCRResponse = response.json().await.map_err(|e| e.to_string())?;
        Ok(ocr_response)
    }

    pub async fn health_check(&self) -> Result<bool, String> {
        let response = self.client
            .get(format!("{}/health", self.ocr_url))
            .send()
            .await
            .map_err(|e| e.to_string())?;
        Ok(response.status().is_success())
    }
}

async fn health(State(state): State<Arc<AppState>>) -> Json<HealthResponse> {
    let ocr_healthy = state.ocr_client.health_check().await.is_ok();

    Json(HealthResponse {
        status: if ocr_healthy { "healthy".to_string() } else { "degraded".to_string() },
        gateway_version: env!("CARGO_PKG_VERSION").to_string(),
        ocr_service: if ocr_healthy { "connected".to_string() } else { "disconnected".to_string() },
        animation_service: "available".to_string(),
    })
}

async fn process_image(
    State(state): State<Arc<AppState>>,
    mut multipart: axum::extract::Multipart,
) -> Result<Json<OCRResponse>, StatusCode> {
    let field = match multipart.next_field().await {
        Ok(Some(field)) => field,
        Ok(None) => return Err(StatusCode::BAD_REQUEST),
        Err(_) => return Err(StatusCode::BAD_REQUEST),
    };

    let file_name = field.file_name().unwrap_or("image.png").to_string();
    let mime = file_name.split('.').last().unwrap_or("");

    let data = match field.bytes().await {
        Ok(bytes) => bytes.to_vec(),
        Err(_) => return Err(StatusCode::BAD_REQUEST),
    };

    let result = if mime == "pdf" {
        state.ocr_client.process_pdf(data).await
    } else {
        state.ocr_client.process_image(data, &file_name).await
    };

    match result {
        Ok(response) => Ok(Json(response)),
        Err(e) => Ok(Json(OCRResponse {
            success: false,
            data: None,
            error: Some(e),
        })),
    }
}

// Animation endpoints

async fn create_animation(
    State(state): State<Arc<AppState>>,
    Json(req): Json<AnimationRequest>,
) -> Result<Json<AnimationResponse>, StatusCode> {
    let job_id = state.animation_bridge.render(&req.code).await
        .map_err(|e| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(AnimationResponse {
        job_id,
        status: JobStatus::Completed,
        output_url: None,
        error: None,
    }))
}

async fn get_animation_status(
    State(state): State<Arc<AppState>>,
    axum::extract::Path(job_id): axum::extract::Path<String>,
) -> Result<Json<AnimationResponse>, StatusCode> {
    let status = state.animation_bridge.queue().get_status(&job_id).await
        .ok_or(StatusCode::NOT_FOUND)?;

    Ok(Json(AnimationResponse {
        job_id,
        status,
        output_url: None,
        error: None,
    }))
}

async fn list_animations(
    State(state): State<Arc<AppState>>,
) -> Json<Vec<animation::AnimationJob>> {
    Json(state.animation_bridge.queue().list_jobs().await)
}

// Chat endpoints

async fn chat(
    State(state): State<Arc<AppState>>,
    Json(req): Json<ChatRequest>,
) -> Result<Json<ChatResponse>, StatusCode> {
    let mode = req.mode.unwrap_or_else(|| ChatMode::from_query(&req.message));
    
    // Search for relevant context using RAG
    let sources = state.rag_engine.search(&req.message, 3).await;
    
    // Mock response - in production would call AI
    let response = format!("I understand you're asking about: {}. Mode: {:?}. Found {} relevant sources.", 
        req.message, mode, sources.len());
    
    Ok(Json(ChatResponse {
        message: response,
        mode,
        sources: sources.iter().map(|s| rag::Source {
            text: s.text.clone(),
            similarity: s.similarity,
            source_type: "rag".to_string(),
        }).collect(),
        session_id: req.session_id.unwrap_or_else(|| uuid::Uuid::new_v4().to_string()),
    }))
}

async fn chat_history(
    State(_state): State<Arc<AppState>>,
) -> Json<Vec<rag::ChatMessage>> {
    // Mock - would fetch from database
    Json(vec![])
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt::init();

    let ocr_url = std::env::var("OCR_SERVICE_URL")
        .unwrap_or_else(|_| "http://127.0.0.1:3000".to_string());

    let animation_output = std::env::var("ANIMATION_OUTPUT_DIR")
        .unwrap_or_else(|_| "/tmp/animations".to_string());

    std::fs::create_dir_all(&animation_output)?;

    let addr: SocketAddr = std::env::var("GATEWAY_HOST")
        .unwrap_or_else(|_| "127.0.0.1:8080".to_string())
        .parse()?;

    let state = Arc::new(AppState {
        ocr_client: Arc::new(OCRClient::new(ocr_url)),
        animation_bridge: Arc::new(ManimBridge::new(
            ManimConfig::default(),
            animation_output,
        )),
        rag_engine: Arc::new(RAGEngine::new()),
    });

    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods(Any)
        .allow_headers(Any);

    let app = Router::new()
        .route("/health", get(health))
        .route("/ocr", post(process_image))
        // Animation endpoints
        .route("/animation", post(create_animation))
        .route("/animation/:id", get(get_animation_status))
        .route("/animations", get(list_animations))
        // Chat endpoints
        .route("/chat", post(chat))
        .route("/chat/history", get(chat_history))
        .layer(ServiceBuilder::new()
            .layer(TraceLayer::new_for_http())
            .layer(cors))
        .with_state(state);

    tracing::info!("Gateway listening on {}", addr);
    axum::serve(tokio::net::TcpListener::bind(addr).await?, app).await?;

    Ok(())
}

#[cfg(test)]
mod integration_tests {
    use super::*;
    use axum::{
        body::Body,
        http::{Request, StatusCode},
    };
    use tower::ServiceExt;

    fn create_test_state() -> Arc<AppState> {
        let ocr_url = "http://127.0.0.1:9999".to_string();
        let animation_output = "/tmp/test_animations";
        std::fs::create_dir_all(animation_output).ok();

        Arc::new(AppState {
            ocr_client: Arc::new(OCRClient::new(ocr_url)),
            animation_bridge: Arc::new(ManimBridge::new(
                ManimConfig::default(),
                animation_output.to_string(),
            )),
            rag_engine: Arc::new(RAGEngine::new()),
        })
    }

    #[tokio::test]
    async fn test_health_endpoint_returns_200() {
        let state = create_test_state();
        let app = Router::new()
            .route("/health", get(health))
            .with_state(state);

        let response = app
            .oneshot(Request::builder().uri("/health").body(Body::empty()).unwrap())
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_chat_endpoint_returns_response() {
        let state = create_test_state();
        let app = Router::new()
            .route("/chat", post(chat))
            .with_state(state);

        let body = serde_json::to_string(&rag::ChatRequest {
            message: "Hello".to_string(),
            mode: Some(rag::ChatMode::Chat),
            session_id: None,
            context: None,
        }).unwrap();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/chat")
                    .method("POST")
                    .header("Content-Type", "application/json")
                    .body(Body::from(body))
                    .unwrap()
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_animation_endpoint_fails_without_manim() {
        let state = create_test_state();
        let app = Router::new()
            .route("/animation", post(create_animation))
            .with_state(state);

        let body = serde_json::to_string(&AnimationRequest {
            code: "from manim import *\nclass Test(Scene):\n    pass".to_string(),
            quality: Some("low".to_string()),
            format: Some("mp4".to_string()),
        }).unwrap();

        let response = app
            .oneshot(
                Request::builder()
                    .uri("/animation")
                    .method("POST")
                    .header("Content-Type", "application/json")
                    .body(Body::from(body))
                    .unwrap()
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    #[tokio::test]
    async fn test_animation_status_endpoint() {
        let state = create_test_state();
        let app = Router::new()
            .route("/animation/{id}", get(get_animation_status))
            .with_state(state);

        let response = app
            .oneshot(Request::builder().uri("/animation/test-job-123").body(Body::empty()).unwrap())
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::NOT_FOUND);
    }
}
