//! HTTP service for cloud-tools
//!
//! Provides REST API for Cloud Shell operations.

use axum::{
    extract::Query,
    routing::{get, post},
    Json, Router,
};
use std::net::SocketAddr;
use tower_http::cors::CorsLayer;

use crate::download;
use crate::gcs;
use crate::models::*;
use crate::ocr;
use crate::ssh;
use crate::upload;

// ============================================================================
// Service Entry Point
// ============================================================================

/// Start the cloud-tools HTTP service
pub async fn run() {
    tracing_subscriber::fmt::init();

    let cors = CorsLayer::permissive();

    let app = Router::new()
        // Health
        .route("/health", get(health))
        // SSH operations
        .route("/ssh/exec", post(ssh_exec_handler))
        // File transfer
        .route("/upload", post(upload_handler))
        .route("/download", post(download_handler))
        // OCR
        .route("/ocr", post(ocr_handler))
        .route("/ocr/pdf", post(pdf_ocr_handler))
        .route("/ocr/status", get(ocr_status))
        // GCS operations
        .route("/gcs/sync", post(gcs_sync_handler))
        .route("/gcs/list", get(gcs_list_handler))
        .route("/gcs/stats", get(gcs_stats_handler))
        .layer(cors);

    let addr = SocketAddr::from(([0, 0, 0, 0], 3004));
    tracing::info!("Starting cloud-tools service on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

// ============================================================================
// Health Endpoint
// ============================================================================

/// Health check
async fn health() -> Json<HealthResponse> {
    let connected = ssh::is_cloud_shell_running().await;
    Json(HealthResponse::healthy(connected))
}

// ============================================================================
// SSH Endpoints
// ============================================================================

/// Execute SSH command on Cloud Shell
async fn ssh_exec_handler(Json(request): Json<SshExecRequest>) -> Json<SshExecResult> {
    let result = ssh::ssh_exec(&request.command, request.timeout_secs).await;
    Json(result)
}

// ============================================================================
// File Transfer Endpoints
// ============================================================================

/// Upload file to Cloud Shell
async fn upload_handler(Json(request): Json<UploadRequest>) -> Json<UploadResult> {
    let result = upload::upload(&request.local_path, &request.remote_path).await;
    Json(result)
}

/// Download file from Cloud Shell
async fn download_handler(Json(request): Json<DownloadRequest>) -> Json<DownloadResult> {
    let result = download::download(&request.remote_path, &request.local_path).await;
    Json(result)
}

// ============================================================================
// OCR Endpoints
// ============================================================================

/// Run OCR on image file
async fn ocr_handler(Json(request): Json<OcrRequest>) -> Json<OcrResult> {
    let result = ocr::run_ocr(&request.image_path, request.language.as_deref()).await;
    Json(result)
}

/// Run OCR on PDF file
async fn pdf_ocr_handler(Json(request): Json<OcrRequest>) -> Json<OcrResult> {
    let result = ocr::run_pdf_ocr(&request.image_path, request.language.as_deref()).await;
    Json(result)
}

/// Check OCR availability
async fn ocr_status() -> Json<OcrStatusResponse> {
    let available = ocr::check_ocr_availability().await;
    let languages = if available {
        ocr::get_ocr_languages().await
    } else {
        vec![]
    };

    Json(OcrStatusResponse {
        available,
        languages,
    })
}

/// OCR status response
#[derive(serde::Serialize)]
struct OcrStatusResponse {
    available: bool,
    languages: Vec<String>,
}

// ============================================================================
// GCS Endpoints
// ============================================================================

/// GCS list parameters
#[derive(serde::Deserialize)]
struct GcsListParams {
    bucket: String,
    prefix: Option<String>,
}

/// GCS stats parameters
#[derive(serde::Deserialize)]
struct GcsStatsParams {
    bucket: String,
}

/// Sync files with GCS bucket
async fn gcs_sync_handler(Json(request): Json<GcsSyncRequest>) -> Json<GcsSyncResult> {
    let result = gcs::sync_gcs(
        &request.bucket,
        request.direction,
        request.local_path.as_deref(),
        request.remote_path.as_deref(),
        request.prefix.as_deref(),
    )
    .await;
    Json(result)
}

/// List GCS bucket contents
async fn gcs_list_handler(Query(params): Query<GcsListParams>) -> Json<GcsListResponse> {
    match gcs::list_bucket(&params.bucket, params.prefix.as_deref()).await {
        Ok(files) => Json(GcsListResponse {
            success: true,
            bucket: params.bucket,
            files,
            error: None,
        }),
        Err(e) => Json(GcsListResponse {
            success: false,
            bucket: params.bucket,
            files: vec![],
            error: Some(e.to_string()),
        }),
    }
}

/// GCS list response
#[derive(serde::Serialize)]
struct GcsListResponse {
    success: bool,
    bucket: String,
    files: Vec<String>,
    error: Option<String>,
}

/// Get GCS bucket statistics
async fn gcs_stats_handler(Query(params): Query<GcsStatsParams>) -> Json<GcsStatsResponse> {
    match gcs::bucket_stats(&params.bucket).await {
        Ok((size, count)) => Json(GcsStatsResponse {
            success: true,
            bucket: params.bucket,
            size_bytes: size,
            file_count: count,
            error: None,
        }),
        Err(e) => Json(GcsStatsResponse {
            success: false,
            bucket: params.bucket,
            size_bytes: 0,
            file_count: 0,
            error: Some(e.to_string()),
        }),
    }
}

/// GCS stats response
#[derive(serde::Serialize)]
struct GcsStatsResponse {
    success: bool,
    bucket: String,
    size_bytes: u64,
    file_count: usize,
    error: Option<String>,
}
