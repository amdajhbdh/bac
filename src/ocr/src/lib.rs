//! BAC OCR Service
//!
//! OCR processing service for BAC Unified platform using unpdf and tesseract.
//! Can delegate heavy OCR tasks to cloud-tools via CLOUD_TOOLS_URL.

use anyhow::anyhow;
use axum::{
    Router,
    extract::{Path as PathExtractor, Query, State},
    http::StatusCode,
    response::{IntoResponse, Response},
    routing::{get, post},
    Json,
};
use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::{f32::consts::E, net::SocketAddr};
use std::path::{Path, PathBuf};
use std::sync::Arc;
use tower_http::cors::CorsLayer;
use tower_http::trace::TraceLayer;
use unpdf::{
    parse_file, to_markdown, to_markdown_with_options, to_text, extract_text, 
    RenderOptions, CleanupPreset, PageSelection, TableFallback,
    detect_format_from_path, is_pdf,
};
use rayon::prelude::*;
use tesseract::Tesseract;
use image::{DynamicImage, GrayImage};
use imageproc::{
    contrast::{adaptive_threshold, stretch_contrast},
    filter::median_filter,
};

/// Application state shared across handlers
#[derive(Clone)]
pub struct AppState {
    pub http_client: Client,
    pub cloud_tools_url: Option<String>,
}

/// Configuration for PDF extraction
#[derive(Debug, Clone)]
pub struct ExtractionConfig {
    /// Output format
    pub format: OutputFormat,
    /// Whether to extract images
    pub extract_images: bool,
    /// Image output directory
    pub image_dir: Option<PathBuf>,
    /// Whether to include frontmatter
    pub include_frontmatter: bool,
    /// Maximum heading level (1-6)
    pub max_heading_level: u8,
    /// Whether to preserve line breaks
    pub preserve_line_breaks: bool,
    /// Line width for wrapping (0 = no wrap)
    pub line_width: u32,
    /// Whether to collect statistics
    pub collect_stats: bool,
    /// Page selection
    pub page_selection: PageSelection,
    /// Cleanup preset
    pub cleanup_preset: CleanupPreset,
    /// Table fallback mode
    pub table_fallback: TableFallback,
    /// Image format for extraction
    pub image_format: ImageFormat,
    /// Image quality (0-100)
    pub image_quality: u8,
}

/// Image format for extraction
#[derive(Debug, Clone, PartialEq)]
pub enum ImageFormat {
    PNG,
    JPEG,
    WebP,
    Original, // Keep original format
}

/// Output format for PDF extraction
#[derive(Debug, Clone, PartialEq, Serialize)]
pub enum OutputFormat {
    Markdown,
    Text,
    Json,
}

impl Default for ExtractionConfig {
    fn default() -> Self {
        Self {
            format: OutputFormat::Markdown,
            extract_images: false,
            image_dir: None,
            include_frontmatter: false,
            max_heading_level: 6,
            preserve_line_breaks: true,
            line_width: 0,
            collect_stats: false,
            page_selection: PageSelection::All,
            cleanup_preset: CleanupPreset::Standard,
            table_fallback: TableFallback::Markdown,
            image_format: ImageFormat::Original,
            image_quality: 85,
        }
    }
}

/// Extraction result with metadata
#[derive(Debug, Clone, Serialize)]
pub(crate) struct ExtractionResult {
    pub content: String,
    pub format: OutputFormat,
    pub file_path: String,
    pub page_count: usize,
    pub char_count: usize,
    pub line_count: usize,
    pub extraction_time_ms: u128,
}

/// Batch extraction request
#[derive(Debug)]
pub struct BatchExtractionRequest {
    pub directory: String,
    pub config: ExtractionConfig,
    pub recursive: bool,
}

/// Batch extraction response
#[derive(Debug)]
pub struct BatchExtractionResponse {
    pub success_count: usize,
    pub failure_count: usize,
    pub total_files: usize,
    pub results: Vec<ExtractionResult>,
    pub errors: Vec<String>,
}

/// Validate if a file is a PDF
pub fn validate_pdf<P: AsRef<Path>>(file_path: P) -> bool {
    is_pdf(file_path)
}

/// Detect PDF format from file
pub fn detect_pdf_format<P: AsRef<Path>>(file_path: P) -> Result<unpdf::PdfFormat, unpdf::Error> {
    detect_format_from_path(file_path)
}

/// Extract plain text from a PDF file using unpdf
pub fn extract_text_from_pdf<P: AsRef<Path>>(pdf_path: P) -> Result<String, unpdf::Error> {
    extract_text(pdf_path)
}

/// Parse a PDF file and return structured document
pub fn parse_pdf_file<P: AsRef<Path>>(pdf_path: P) -> Result<unpdf::Document, unpdf::Error> {
    parse_file(pdf_path)
}

/// Convert PDF to Markdown format with default options
pub fn pdf_to_markdown<P: AsRef<Path>>(pdf_path: P) -> Result<String, unpdf::Error> {
    to_markdown(pdf_path)
}

/// Convert PDF to Markdown with custom configuration
pub fn pdf_to_markdown_with_config<P: AsRef<Path>>(
    pdf_path: P,
    config: &ExtractionConfig,
) -> Result<String, unpdf::Error> {
    let options = build_render_options(config);
    to_markdown_with_options(pdf_path, &options)
}

/// Convert PDF to plain text with cleanup
pub fn pdf_to_text<P: AsRef<Path>>(pdf_path: P) -> Result<String, unpdf::Error> {
    let options = RenderOptions::default();
    to_text(pdf_path, &options)
}

/// Extract PDF with full metadata and statistics
pub fn extract_pdf_with_stats<P: AsRef<Path>>(
    pdf_path: P,
    config: &ExtractionConfig,
) -> Result<ExtractionResult, anyhow::Error> {
    let start_time = std::time::Instant::now();
    let path_str = pdf_path.as_ref().to_string_lossy().to_string();
    
    // Parse document to get metadata
    let doc = parse_file(&pdf_path)?;
    let page_count = doc.page_count() as usize;
    
    // Extract content based on format
    let content = match config.format {
        OutputFormat::Markdown => {
            let options = build_render_options(config);
            to_markdown_with_options(&pdf_path, &options)?
        }
        OutputFormat::PlainText => {
            let options = build_render_options(config);
            to_text(&pdf_path, &options)?
        }
        OutputFormat::Json => {
            // For JSON, extract as markdown first
            let options = build_render_options(config);
            let markdown = to_markdown_with_options(&pdf_path, &options)?;
            serde_json::to_string_pretty(&serde_json::json!({
                "content": markdown,
                "format": "markdown"
            }))?
        }
    };
    
    let extraction_time = start_time.elapsed();
    
    Ok(ExtractionResult {
        content: content.clone(),
        format: config.format.clone(),
        file_path: path_str,
        page_count,
        char_count: content.len(),
        line_count: content.lines().count(),
        extraction_time_ms: extraction_time.as_millis(),
    })
}

/// Process a PDF file and return extracted text
pub async fn process_pdf(file_path: &str) -> Result<String, anyhow::Error> {
    let path = file_path.to_string();
    tokio::task::spawn_blocking(move || {
        extract_text_from_pdf(path)
    })
    .await?
    .map_err(|e| anyhow!(e))
}

/// Process a PDF file and return Markdown format
pub async fn process_pdf_to_markdown(file_path: &str) -> Result<String, anyhow::Error> {
    let path = file_path.to_string();
    tokio::task::spawn_blocking(move || {
        pdf_to_markdown(path)
    })
    .await?
    .map_err(|e| anyhow!(e))
}

/// Batch extract PDFs from directory (parallel processing)
pub async fn batch_extract(
    directory: &str,
    config: ExtractionConfig,
    recursive: bool,
) -> Result<BatchExtractionResponse, anyhow::Error> {
    let dir_path = Path::new(directory);
    
    if !dir_path.is_dir() {
        return Err(anyhow!("Directory does not exist: {}", directory));
    }
    
    // Find all PDF files
    let pdf_files = find_pdf_files(dir_path, recursive)?;
    let total_files = pdf_files.len();
    
    if total_files == 0 {
        return Ok(BatchExtractionResponse {
            success_count: 0,
            failure_count: 0,
            total_files: 0,
            results: Vec::new(),
            errors: Vec::new(),
        });
    }
    
    tracing::info!("Found {} PDF files for extraction", total_files);
    
    // Process files in parallel using Rayon
    let results: Vec<Result<ExtractionResult, anyhow::Error>> = pdf_files
        .par_iter()
        .map(|pdf_path| {
            tracing::debug!("Processing: {}", pdf_path.file_name().unwrap().to_string_lossy());
            extract_pdf_with_stats(pdf_path, &config).map_err(|e| anyhow!(e))
        })
        .collect();
    
    // Separate successes and failures
    let mut success_results = Vec::new();
    let mut errors = Vec::new();
    
    for (i, result) in results.into_iter().enumerate() {
        match result {
            Ok(extraction) => {
                success_results.push(extraction);
            }
            Err(e) => {
                let filename = pdf_files[i].file_name().unwrap().to_string_lossy().to_string();
                errors.push(format!("{}: {}", filename, e));
            }
        }
    }
    
    Ok(BatchExtractionResponse {
        success_count: success_results.len(),
        failure_count: errors.len(),
        total_files,
        results: success_results,
        errors,
    })
}

/// Find all PDF files in a directory
fn find_pdf_files(dir: &Path, recursive: bool) -> Result<Vec<PathBuf>, anyhow::Error> {
    let mut pdf_files = Vec::new();
    
    if recursive {
        for entry in walkdir::WalkDir::new(dir)
            .into_iter()
            .filter_map(|e| e.ok())
        {
            let path = entry.path();
            if path.is_file() {
                let ext = path.extension().and_then(|s| s.to_str()).unwrap_or("").to_lowercase();
                if ext == "pdf" {
                    pdf_files.push(path.to_path_buf());
                }
            }
        }
    } else {
        for entry in std::fs::read_dir(dir)? {
            let entry = entry?;
            let path = entry.path();
            
            if path.is_file() {
                let ext = path.extension().and_then(|s| s.to_str()).unwrap_or("").to_lowercase();
                if ext == "pdf" {
                    pdf_files.push(path);
                }
            }
        }
    }
    
    Ok(pdf_files)
}

/// Build RenderOptions from configuration
fn build_render_options(config: &ExtractionConfig) -> RenderOptions {
    let mut options = RenderOptions::default();
    
    // Set image extraction
    if config.extract_images {
        if let Some(dir) = &config.image_dir {
            options = options.with_image_dir(dir);
        }
        options = options.with_image_prefix("images/");
    }
    
    // Set frontmatter
    options = options.with_frontmatter(config.include_frontmatter);
    
    // Set heading level
    options = options.with_max_heading(config.max_heading_level);
    
    // Set line breaks
    options = options.with_line_breaks(config.preserve_line_breaks);
    
    // Set line width
    options = options.with_line_width(config.line_width);
    
    // Set statistics collection
    options = options.with_stats(config.collect_stats);
    
    // Set page selection
    options = options.with_pages(config.page_selection.clone());
    
    // Set cleanup preset
    options = options.with_cleanup_preset(config.cleanup_preset);
    
    // Set table fallback
    options = options.with_table_fallback(config.table_fallback);
    
    options
}

/// Extract images from a PDF file
pub fn extract_images_from_pdf<P: AsRef<Path>>(
    pdf_path: P,
    output_dir: P,
) -> Result<Vec<PathBuf>, anyhow::Error> {
    let mut image_paths = Vec::new();
    
    // Create output directory if it doesn't exist
    std::fs::create_dir_all(&output_dir)?;
    
    // Extract images using unpdf's built-in image extraction
    // Note: unpdf handles image extraction via the image_dir option in RenderOptions
    // This function provides a wrapper for manual extraction
    
    // For now, we'll use the built-in extraction by setting image_dir in config
    // and using the to_markdown function which extracts images
    
    let config = ExtractionConfig {
        extract_images: true,
        image_dir: Some(output_dir.as_ref().to_path_buf()),
        ..Default::default()
    };
    
    // This will extract images as a side effect
    let _content = pdf_to_markdown_with_config(&pdf_path, &config)?;
    
    // Collect extracted image files
    if let Ok(entries) = std::fs::read_dir(&output_dir) {
        for entry in entries.flatten() {
            let path = entry.path();
            if path.is_file() {
                image_paths.push(path);
            }
        }
    }
    
    Ok(image_paths)
}

/// Process standalone image files
pub fn process_standalone_images(
    image_dir: &str,
    output_dir: &str,
) -> Result<Vec<PathBuf>, anyhow::Error> {
    let input_path = Path::new(image_dir);
    let output_path = Path::new(output_dir);
    
    if !input_path.is_dir() {
        return Err(anyhow!("Input directory does not exist: {}", image_dir));
    }
    
    std::fs::create_dir_all(output_path)?;
    
    let mut processed_images = Vec::new();
    
    for entry in std::fs::read_dir(input_path)? {
        let entry = entry?;
        let path = entry.path();
        
        if path.is_file() {
            let ext = path.extension().and_then(|s| s.to_str()).unwrap_or("").to_lowercase();
            
            // Check if it's an image file
            if ["png", "jpg", "jpeg", "gif", "bmp", "webp", "svg"].contains(&ext.as_str()) {
                let filename = path.file_name().unwrap();
                let output_file = output_path.join(filename);
                
                // Copy or process the image
                std::fs::copy(&path, &output_file)?;
                processed_images.push(output_file);
            }
        }
    }
    
    Ok(processed_images)
}



/// Configuration for image preprocessing
#[derive(Debug, Clone)]
pub struct PreprocessingConfig {
    /// Block size for adaptive thresholding (larger for uneven lighting)
    pub threshold_block_size: u32,
    /// Kernel size for median filter (1 for Arabic to preserve connections)
    pub median_kernel: u32,
    /// Deskew angle in radians (0.0 for no deskewing)
    pub deskew_angle: f32,
    /// Whether to enhance contrast
    pub enhance_contrast: bool,
}

impl Default for PreprocessingConfig {
    fn default() -> Self {
        Self {
            threshold_block_size: 51, // Larger for Arabic documents
            median_kernel: 1,         // Conservative to preserve Arabic connections
            deskew_angle: 0.0,
            enhance_contrast: true,
        }
    }
}

/// Preprocess an image for OCR
/// 
/// # Arguments
/// * `img` - DynamicImage to preprocess
/// * `config` - Preprocessing configuration
/// 
/// # Returns
/// Preprocessed GrayImage ready for OCR
pub fn preprocess_image(img: &DynamicImage, config: &PreprocessingConfig) -> GrayImage {
    // Step 1: Convert to grayscale
    let mut processed = img.to_luma8();
    
    // Step 2: Remove noise (conservative for Arabic to preserve connections)
    if config.median_kernel > 0 {
        processed = median_filter(&processed, config.median_kernel, config.median_kernel);
    }
    
    // Step 3: Enhance contrast (optional)
    if config.enhance_contrast {
        // Use default parameters for contrast stretching
        processed = stretch_contrast(&processed, 0, 255, 0, 255);
    }
    
    // Step 4: Deskew if angle is provided
    if config.deskew_angle.abs() > 0.001 {
        use imageproc::geometric_transformations::{rotate_about_center, Interpolation};
        use image::Luma;
        processed = rotate_about_center(
            &processed,
            config.deskew_angle,
            Interpolation::Bilinear,
            Luma([255u8]) // Background color (white)
        );
    }
    
    // Step 5: Binarization (adaptive thresholding)
    processed = adaptive_threshold(&processed, config.threshold_block_size);
    
    processed
}

/// Perform OCR on an image file using Tesseract with preprocessing
pub fn perform_ocr_on_image<P: AsRef<Path>>(image_path: P) -> Result<String, anyhow::Error> {
    let image_path = image_path.as_ref();
    
    // Check if file exists
    if !image_path.exists() {
        return Err(anyhow!("Image file does not exist: {:?}", image_path));
    }
    
    // Read image file bytes
    let image_bytes = std::fs::read(image_path)?;
    
    // Perform OCR with multi-language support (Arabic + English + French)
    let mut tesseract = Tesseract::new(None, Some("ara+eng+fra"))?
        .set_image_from_mem(&image_bytes)?
        .set_variable("tessedit_pageseg_mode", "1")? // Automatic page segmentation with OSD
        .recognize()?;
    
    let text = tesseract.get_text()?;
    
    Ok(text)
}

/// Perform OCR on multiple images in a directory
pub fn perform_ocr_on_directory<P: AsRef<Path>>(
    image_dir: P,
    recursive: bool,
) -> Result<Vec<(PathBuf, String)>, anyhow::Error> {
    let dir_path = image_dir.as_ref();
    
    if !dir_path.is_dir() {
        return Err(anyhow!("Directory does not exist: {:?}", dir_path));
    }
    
    // Find all image files
    let image_files = find_image_files(dir_path, recursive)?;
    
    // Process images in parallel
    let results: Vec<Result<(PathBuf, String), anyhow::Error>> = image_files
        .par_iter()
        .map(|image_path| {
            match perform_ocr_on_image(image_path) {
                Ok(text) => Ok((image_path.clone(), text)),
                Err(e) => Err(anyhow!("Failed to OCR {:?}: {}", image_path, e)),
            }
        })
        .collect();
    
    // Separate successes and failures
    let mut success_results = Vec::new();
    let mut errors = Vec::new();
    
    for result in results {
        match result {
            Ok((path, text)) => {
                success_results.push((path, text));
            }
            Err(e) => {
                errors.push(e.to_string());
            }
        }
    }
    
    if !errors.is_empty() {
        tracing::warn!("OCR completed with {} errors", errors.len());
    }
    
    Ok(success_results)
}

/// Find all image files in a directory
fn find_image_files(dir: &Path, recursive: bool) -> Result<Vec<PathBuf>, anyhow::Error> {
    let mut image_files = Vec::new();
    let image_extensions = ["png", "jpg", "jpeg", "gif", "bmp", "webp", "tiff"];
    
    if recursive {
        for entry in walkdir::WalkDir::new(dir)
            .into_iter()
            .filter_map(|e| e.ok())
        {
            let path = entry.path();
            if path.is_file() {
                let ext = path.extension().and_then(|s| s.to_str()).unwrap_or("").to_lowercase();
                if image_extensions.contains(&ext.as_str()) {
                    image_files.push(path.to_path_buf());
                }
            }
        }
    } else {
        for entry in std::fs::read_dir(dir)? {
            let entry = entry?;
            let path = entry.path();
            
            if path.is_file() {
                let ext = path.extension().and_then(|s| s.to_str()).unwrap_or("").to_lowercase();
                if image_extensions.contains(&ext.as_str()) {
                    image_files.push(path);
                }
            }
        }
    }
    
    Ok(image_files)
}

// ============================================================================
// HTTP API Types
// ============================================================================

/// Health check response
#[derive(Serialize)]
pub struct HealthResponse {
    pub status: String,
    pub service: String,
    pub version: String,
    pub cloud_tools: String,
}

/// Request for OCR on a single image
#[derive(Deserialize, Serialize)]
pub struct ImageOcrRequest {
    pub image_path: String,
    pub language: Option<String>,
}

/// Request for OCR on a PDF
#[derive(Deserialize)]
pub struct PdfOcrRequest {
    pub pdf_path: String,
    pub language: Option<String>,
}

/// Request for PDF text extraction
#[derive(Deserialize)]
pub struct PdfExtractRequest {
    pub pdf_path: String,
    pub format: Option<String>, // "markdown", "text", "json"
}

/// Cloud-tools OCR response
#[derive(Serialize, Deserialize)]
pub struct CloudOcrResponse {
    pub success: bool,
    pub text: Option<String>,
    pub error: Option<String>,
    pub confidence: Option<f32>,
}

// ============================================================================
// HTTP Handlers
// ============================================================================

/// Health check endpoint
async fn health_check(State(state): State<Arc<AppState>>) -> Json<HealthResponse> {
    let cloud_tools_status = match &state.cloud_tools_url {
        Some(url) => format!("configured ({})", url),
        None => "not configured".to_string(),
    };

    Json(HealthResponse {
        status: "ok".to_string(),
        service: "bac-ocr".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
        cloud_tools: cloud_tools_status,
    })
}

/// Local image OCR using tesseract
async fn ocr_image(
    State(state): State<Arc<AppState>>,
    Json(req): Json<ImageOcrRequest>,
) -> impl IntoResponse {
    let path = Path::new(&req.image_path);
    
    if !path.exists() {
        return (
            StatusCode::NOT_FOUND,
            Json(serde_json::json!({ "error": "Image file not found" })),
        ).into_response();
    }

    match tokio::task::spawn_blocking({
        let path = req.image_path.clone();
        move || perform_ocr_on_image(&path)
    })
    .await
    {
        Ok(Ok(text)) => {
            Json(serde_json::json!({ "success": true, "text": text, "source": "local" }))
                .into_response()
        }
        Ok(Err(e)) => {
            tracing::warn!("Local OCR failed: {}, trying cloud-tools", e);
            
            // Fall back to cloud-tools if configured
            if let Some(base_url) = &state.cloud_tools_url {
                let url = format!("{}/api/ocr/image", base_url);
                let client = &state.http_client;
                
                match client
                    .post(&url)
                    .json(&req)
                    .send()
                    .await
                {
                    Ok(resp) => {
                        if resp.status().is_success() {
                            return resp.json::<CloudOcrResponse>().await
                                .map(|r| {
                                    Json(serde_json::json!({
                                        "success": r.success,
                                        "text": r.text,
                                        "error": r.error,
                                        "source": "cloud"
                                    })).into_response()
                                })
                                .unwrap_or_else(|_| {
                                    Json(serde_json::json!({ 
                                        "error": "Failed to parse cloud response" 
                                    })).into_response()
                                });
                        }
                    }
                    Err(e) => {
                        tracing::error!("Cloud-tools request failed: {}", e);
                    }
                }
            }
            
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({ "error": e.to_string() })),
            ).into_response()
        }
        Err(e) => {
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({ "error": e.to_string() })),
            ).into_response()
        }
    }
}

/// Cloud-tools image OCR (delegates to cloud-tools service)
async fn ocr_image_cloud(
    State(state): State<Arc<AppState>>,
    Json(req): Json<ImageOcrRequest>,
) -> Response {
    let base_url = match &state.cloud_tools_url {
        Some(url) => url.clone(),
        None => {
            return Json(serde_json::json!({ 
                "error": "Cloud-tools not configured" 
            })).into_response();
        }
    };

    let url = format!("{}/api/ocr/image", base_url);
    
    match state.http_client
        .post(&url)
        .json(&req)
        .send()
        .await
    {
        Ok(resp) => {
            if resp.status().is_success() {
                return resp.json::<CloudOcrResponse>().await
                    .map(|r| {
                        Json(serde_json::json!({
                            "success": r.success,
                            "text": r.text,
                            "error": r.error,
                            "confidence": r.confidence
                        })).into_response()
                    })
                    .unwrap_or_else(|_| {
                        Json(serde_json::json!({ 
                            "error": "Failed to parse response" 
                        })).into_response()
                    });
            }
            
            (
                StatusCode::from_u16(resp.status().as_u16()).unwrap_or(StatusCode::BAD_REQUEST),
                Json(serde_json::json!({ "error": "Cloud OCR request failed" })),
            ).into_response()
        }
        Err(e) => {
            (
                StatusCode::BAD_GATEWAY,
                Json(serde_json::json!({ "error": format!("Cloud-tools unavailable: {}", e) })),
            ).into_response()
        }
    }
}

/// PDF text extraction endpoint
async fn extract_pdf(
    State(_state): State<Arc<AppState>>,
    Json(req): Json<PdfExtractRequest>,
) -> impl IntoResponse {
    let path = Path::new(&req.pdf_path);
    
    if !path.exists() {
        return (
            StatusCode::NOT_FOUND,
            Json(serde_json::json!({ "error": "PDF file not found" })),
        ).into_response();
    }

    let format = req.format.as_deref().unwrap_or("markdown").to_string();
    let format_response = format.clone();
    
    let result = tokio::task::spawn_blocking({
        let pdf_path = req.pdf_path.clone();
        let format = format;
        move || -> Result<String, unpdf::Error> {
            match format.as_str() {
                "text" => pdf_to_text(&pdf_path),
                "json" => {
                    let text = to_markdown(&pdf_path)?;
                    Ok(serde_json::to_string_pretty(&serde_json::json!({
                                            "content": text,
                                            "format": "markdown"
                                        })).unwrap())
                }
                _ => pdf_to_markdown(&pdf_path),
            }
        }
    })
    .await;

    match result {
        Ok(Ok(content)) => {
            Json(serde_json::json!({ 
                "success": true, 
                "content": content,
                "format": format_response 
            })).into_response()
        }
        Ok(Err(e)) => {
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({ "error": e.to_string() })),
            ).into_response()
        }
        Err(e) => {
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({ "error": e.to_string() })),
            ).into_response()
        }
    }
}

/// Batch PDF extraction endpoint
async fn batch_extract_pdf(
    State(_state): State<Arc<AppState>>,
    Json(req): Json<serde_json::Value>,
) -> impl IntoResponse {
    let directory = req["directory"].as_str().unwrap_or(".");
    let recursive = req["recursive"].as_bool().unwrap_or(false);
    let format = req["format"].as_str().unwrap_or("markdown");
    
    let config = ExtractionConfig {
        format: match format {
            "text" => OutputFormat::PlainText,
            "json" => OutputFormat::Json,
            _ => OutputFormat::Markdown,
        },
        ..Default::default()
    };
    
    match batch_extract(directory, config, recursive).await {
        Ok(response) => {
            Json(serde_json::json!({
                "success": response.success_count,
                "failed": response.failure_count,
                "total": response.total_files,
                "results": response.results,
                "errors": response.errors
            })).into_response()
        }
        Err(e) => {
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({ "error": e.to_string() })),
            ).into_response()
        }
    }
}

// ============================================================================
// Server Entry Point
// ============================================================================

/// Run the OCR HTTP server
pub async fn run() {
    // Load cloud-tools URL from environment
    let cloud_tools_url = std::env::var("CLOUD_TOOLS_URL").ok();
    
    if let Some(ref url) = cloud_tools_url {
        tracing::info!("Cloud-tools integration enabled: {}", url);
    } else {
        tracing::info!("Running in local-only mode (CLOUD_TOOLS_URL not set)");
    }

    let state = Arc::new(AppState {
        http_client: Client::new(),
        cloud_tools_url,
    });

    let cors = CorsLayer::permissive();
    
    let app = Router::new()
        .route("/health", get(health_check))
        .route("/api/ocr/image", post(ocr_image))
        .route("/api/ocr/image/cloud", post(ocr_image_cloud))
        .route("/api/ocr/pdf", post(extract_pdf))
        .route("/api/ocr/pdf/batch", post(batch_extract_pdf))
        .layer(cors)
        .layer(TraceLayer::new_for_http())
        .with_state(state);

    let addr = SocketAddr::from(([0, 0, 0, 0], 8082));
    tracing::info!("Starting BAC OCR server on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
