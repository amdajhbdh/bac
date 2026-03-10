//! HTTP Server for OCR Service
//!
//! Provides REST API endpoints for OCR processing

use std::net::SocketAddr;
use std::sync::Arc;

use axum::{
    extract::{DefaultBodyLimit, Multipart, State},
    http::StatusCode,
    response::Json,
    routing::{get, post},
    Router,
};
use serde::{Deserialize, Serialize};
use tokio::sync::RwLock;
use tower::ServiceBuilder;
use tower_http::cors::{Any, CorsLayer};
use tower_http::trace::TraceLayer;

use crate::pipeline::{OCRPipeline, PipelineResult, PDFResult};

/// Application state
pub struct AppState {
    pub pipeline: Arc<RwLock<OCRPipeline>>,
}

/// Response for successful OCR
#[derive(Serialize)]
pub struct OCRResponse {
    pub success: bool,
    pub data: Option<OCRData>,
    pub error: Option<String>,
}

#[derive(Serialize)]
pub struct OCRData {
    pub text: String,
    pub confidence: f64,
    pub engine: String,
    pub processing_time_ms: i64,
}

impl From<PipelineResult> for OCRData {
    fn from(result: PipelineResult) -> Self {
        OCRData {
            text: result.text,
            confidence: result.confidence,
            engine: result.engine.to_string(),
            processing_time_ms: result.processing_time_ms,
        }
    }
}

/// Response for PDF OCR
#[derive(Serialize)]
pub struct PDFResponse {
    pub success: bool,
    pub data: Option<PDFData>,
    pub error: Option<String>,
}

#[derive(Serialize)]
pub struct PDFData {
    pub text: String,
    pub total_pages: usize,
    pub average_confidence: f64,
    pub processing_time_ms: i64,
}

impl From<PDFResult> for PDFData {
    fn from(result: PDFResult) -> Self {
        PDFData {
            text: result.combined_text,
            total_pages: result.total_pages,
            average_confidence: result.average_confidence,
            processing_time_ms: result.processing_time_ms,
        }
    }
}

/// Health check response
#[derive(Serialize)]
pub struct HealthResponse {
    pub status: String,
    pub version: String,
}

/// Process image OCR
async fn process_image(
    State(state): State<Arc<AppState>>,
    mut multipart: Multipart,
) -> Result<Json<OCRResponse>, StatusCode> {
    // Get the first file from multipart
    let field = match multipart.next_field().await {
        Ok(Some(field)) => field,
        Ok(None) => return Err(StatusCode::BAD_REQUEST),
        Err(_) => return Err(StatusCode::BAD_REQUEST),
    };

    let file_name = field.file_name().unwrap_or("image").to_string();
    let mime = file_name.split('.').last().unwrap_or("");
    
    let data = match field.bytes().await {
        Ok(bytes) => bytes,
        Err(_) => return Err(StatusCode::BAD_REQUEST),
    };
    
    // Check file size (max 50MB)
    if data.len() > 50 * 1024 * 1024 {
        return Ok(Json(OCRResponse {
            success: false,
            data: None,
            error: Some("File too large (max 50MB)".to_string()),
        }));
    }
    
    let pipeline = state.pipeline.read().await;
    
    // Check if it's a PDF
    if mime == "pdf" {
        match pipeline.process_pdf(&data).await {
            Ok(result) => Ok(Json(OCRResponse {
                success: true,
                data: Some(OCRData {
                    text: result.combined_text,
                    confidence: result.average_confidence,
                    engine: "pdf".to_string(),
                    processing_time_ms: result.processing_time_ms,
                }),
                error: None,
            })),
            Err(e) => Ok(Json(OCRResponse {
                success: false,
                data: None,
                error: Some(e.to_string()),
            })),
        }
    } else {
        match pipeline.process(&data).await {
            Ok(result) => Ok(Json(OCRResponse {
                success: true,
                data: Some(result.into()),
                error: None,
            })),
            Err(e) => Ok(Json(OCRResponse {
                success: false,
                data: None,
                error: Some(e.to_string()),
            })),
        }
    }
}

/// Process PDF OCR
async fn process_pdf(
    State(state): State<Arc<AppState>>,
    mut multipart: Multipart,
) -> Result<Json<PDFResponse>, StatusCode> {
    let field = match multipart.next_field().await {
        Ok(Some(field)) => field,
        Ok(None) => return Err(StatusCode::BAD_REQUEST),
        Err(_) => return Err(StatusCode::BAD_REQUEST),
    };

    let data = match field.bytes().await {
        Ok(data) => data,
        Err(_) => return Err(StatusCode::BAD_REQUEST),
    };

    if data.len() > 50 * 1024 * 1024 {
        return Ok(Json(PDFResponse {
            success: false,
            data: None,
            error: Some("File too large (max 50MB)".to_string()),
        }));
    }

    let pipeline = state.pipeline.read().await;
    
    match pipeline.process_pdf(&data).await {
        Ok(result) => Ok(Json(PDFResponse {
            success: true,
            data: Some(result.into()),
            error: None,
        })),
        Err(e) => Ok(Json(PDFResponse {
            success: false,
            data: None,
            error: Some(e.to_string()),
        })),
    }
}

/// Process image from URL
async fn process_image_url(
    State(state): State<Arc<AppState>>,
    Json(payload): Json<ImageUrlRequest>,
) -> Result<Json<OCRResponse>, StatusCode> {
    // Fetch image from URL
    let client = reqwest::Client::new();
    let response = match client.get(&payload.url).send().await {
        Ok(resp) => resp,
        Err(_) => return Err(StatusCode::BAD_REQUEST),
    };

    let data = match response.bytes().await {
        Ok(bytes) => bytes,
        Err(_) => return Err(StatusCode::BAD_REQUEST),
    };

    let pipeline = state.pipeline.read().await;
    
    match pipeline.process(&data).await {
        Ok(result) => Ok(Json(OCRResponse {
            success: true,
            data: Some(result.into()),
            error: None,
        })),
        Err(e) => Ok(Json(OCRResponse {
            success: false,
            data: None,
            error: Some(e.to_string()),
        })),
    }
}

#[derive(Deserialize)]
pub struct ImageUrlRequest {
    url: String,
}

/// Batch process images
async fn process_batch(
    State(state): State<Arc<AppState>>,
    Json(payload): Json<BatchRequest>,
) -> Result<Json<BatchResponse>, StatusCode> {
    let image_count = payload.images.len();
    let pipeline = state.pipeline.read().await;
    
    let results = pipeline.process_batch(payload.images).await;
    
    let mut successful = 0;
    let mut failed = 0;
    let responses: Vec<OCRResponse> = results
        .into_iter()
        .map(|r| match r {
            Ok(result) => {
                successful += 1;
                OCRResponse {
                    success: true,
                    data: Some(result.into()),
                    error: None,
                }
            }
            Err(e) => {
                failed += 1;
                OCRResponse {
                    success: false,
                    data: None,
                    error: Some(e.to_string()),
                }
            }
        })
        .collect();

    Ok(Json(BatchResponse {
        total: image_count,
        successful,
        failed,
        results: responses,
    }))
}

#[derive(Deserialize)]
pub struct BatchRequest {
    images: Vec<Vec<u8>>,
}

#[derive(Serialize)]
pub struct BatchResponse {
    total: usize,
    successful: usize,
    failed: usize,
    results: Vec<OCRResponse>,
}

/// Health check
async fn health() -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "healthy".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
    })
}

/// Create the router
pub fn create_router(pipeline: OCRPipeline) -> Router {
    let state = Arc::new(AppState {
        pipeline: Arc::new(RwLock::new(pipeline)),
    });

    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods(Any)
        .allow_headers(Any);

    Router::new()
        .route("/health", get(health))
        .route("/ocr", post(process_image))
        .route("/ocr/url", post(process_image_url))
        .route("/ocr/batch", post(process_batch))
        .route("/pdf", post(process_pdf))
        .layer(ServiceBuilder::new()
            .layer(TraceLayer::new_for_http())
            .layer(cors)
            .layer(DefaultBodyLimit::max(50 * 1024 * 1024)))
        .with_state(state)
}

/// Run the HTTP server
pub async fn run_server(addr: SocketAddr, pipeline: OCRPipeline) {
    let router = create_router(pipeline);
    
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    
    tracing::info!("OCR HTTP server listening on {}", addr);
    
    axum::serve(listener, router).await.unwrap();
}
