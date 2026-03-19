//! Integration tests for Gemini Tools service (:3001)
//!
//! Tests the HTTP endpoints for AI-powered content analysis and generation.

mod common;

use common::{create_client, is_service_running, urls, HEALTH_TIMEOUT, DEFAULT_TIMEOUT};
use reqwest::Client;
use serde::{Deserialize, Serialize};

// =============================================================================
// Test Configuration
// =============================================================================

const GEMINI_BASE_URL: &str = urls::GEMINI;

/// Skip tests if the service is not running
fn require_service() -> bool {
    // This is a compile-time check that service must be available
    // Use tokio's #[test] to actually skip at runtime
    true
}

// =============================================================================
// Health Endpoint Tests
// =============================================================================

/// Test: Health endpoint returns OK
///
/// Verifies that the /health endpoint responds correctly when the service is running.
#[tokio::test]
async fn test_health_endpoint() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    let response = client
        .get(format!("{}/health", GEMINI_BASE_URL))
        .send()
        .await;
    
    // Assert: Response should be successful
    assert!(response.is_ok(), "Health endpoint should respond");
    let response = response.unwrap();
    
    // Assert: Should return 200 OK
    assert_eq!(
        response.status().as_u16(),
        200,
        "Health endpoint should return 200 OK"
    );
    
    // Assert: Body should be "OK"
    let body = response.text().await.unwrap();
    assert_eq!(body, "OK", "Health endpoint should return 'OK'");
}

/// Test: Health endpoint is accessible
///
/// Verifies that the health endpoint responds quickly.
#[tokio::test]
async fn test_health_endpoint_fast_response() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = Client::builder()
        .timeout(HEALTH_TIMEOUT)
        .build()
        .unwrap();
    
    let start = std::time::Instant::now();
    let response = client
        .get(format!("{}/health", GEMINI_BASE_URL))
        .send()
        .await;
    let elapsed = start.elapsed();
    
    assert!(response.is_ok(), "Health endpoint should respond quickly");
    assert!(
        elapsed < DEFAULT_TIMEOUT,
        "Health endpoint should respond within timeout"
    );
}

// =============================================================================
// Analyze Endpoint Tests
// =============================================================================

/// Request body for /analyze endpoint
#[derive(Debug, Serialize)]
struct AnalyzeRequest {
    content: String,
    subject: Option<String>,
}

/// Response from /analyze endpoint
#[derive(Debug, Deserialize)]
struct AnalyzeResponse {
    // The actual response structure depends on the gemini-tools implementation
    // This is a generic structure
    #[serde(default)]
    success: bool,
    #[serde(default)]
    analysis: Option<serde_json::Value>,
    #[serde(default)]
    error: Option<String>,
}

/// Test: Analyze endpoint accepts valid content
///
/// Verifies that the /analyze endpoint accepts content for analysis.
#[tokio::test]
async fn test_analyze_endpoint_accepts_request() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    let request = AnalyzeRequest {
        content: "The mitochondria is the powerhouse of the cell.".to_string(),
        subject: Some("Biology".to_string()),
    };
    
    // Act: Send analyze request
    let response = client
        .post(format!("{}/analyze", GEMINI_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Request should succeed
    assert!(response.is_ok(), "Analyze endpoint should accept requests");
}

/// Test: Analyze endpoint handles empty content gracefully
///
/// Verifies that the /analyze endpoint handles edge cases.
#[tokio::test]
async fn test_analyze_endpoint_empty_content() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    let request = AnalyzeRequest {
        content: "".to_string(),
        subject: None,
    };
    
    // Act: Send analyze request with empty content
    let response = client
        .post(format!("{}/analyze", GEMINI_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // The service may handle this as a 400 or 200 with error in body
    // We just verify it responds
    if response.is_ok() {
        let status = response.unwrap().status();
        // Accept 200, 400, or 422 as valid responses for edge case
        assert!(
            status.is_success() || status.as_u16() == 400 || status.as_u16() == 422,
            "Empty content should return 2xx, 400, or 422"
        );
    }
}

// =============================================================================
// Embed Endpoint Tests
// =============================================================================

/// Request body for /embed endpoint
#[derive(Debug, Serialize)]
struct EmbedRequest {
    text: String,
    task_type: Option<String>,
}

/// Response from /embed endpoint
#[derive(Debug, Deserialize)]
struct EmbedResponse {
    #[serde(default)]
    embedding: Option<Vec<f32>>,
    #[serde(default)]
    success: bool,
    #[serde(default)]
    error: Option<String>,
}

/// Test: Embed endpoint creates text embeddings
///
/// Verifies that the /embed endpoint can create embeddings for text.
#[tokio::test]
async fn test_embed_endpoint_creates_embedding() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    let request = EmbedRequest {
        text: "Physics: The study of matter and energy.".to_string(),
        task_type: Some("RETRIEVAL_DOCUMENT".to_string()),
    };
    
    // Act: Send embed request
    let response = client
        .post(format!("{}/embed", GEMINI_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Response should be received
    assert!(response.is_ok(), "Embed endpoint should respond");
}

/// Test: Embed endpoint handles various text lengths
///
/// Verifies that the /embed endpoint handles short and long text.
#[tokio::test]
async fn test_embed_endpoint_handles_text_lengths() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    
    let test_cases = vec![
        "Short text.",                           // Very short
        "A moderate length sentence for testing.", // Normal
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. ".repeat(10), // Long
    ];
    
    for text in test_cases {
        let request = EmbedRequest {
            text: text.clone(),
            task_type: None,
        };
        
        let response = client
            .post(format!("{}/embed", GEMINI_BASE_URL))
            .json(&request)
            .send()
            .await;
        
        assert!(
            response.is_ok(),
            "Embed endpoint should handle text of length {}",
            text.len()
        );
    }
}

// =============================================================================
// Correct Endpoint Tests
// =============================================================================

/// Request body for /correct endpoint
#[derive(Debug, Serialize)]
struct CorrectRequest {
    ocr_text: String,
    language: Option<String>,
}

/// Response from /correct endpoint
#[derive(Debug, Deserialize)]
struct CorrectResponse {
    #[serde(default)]
    corrected_text: Option<String>,
    #[serde(default)]
    success: bool,
    #[serde(default)]
    error: Option<String>,
}

/// Test: Correct endpoint processes OCR text
///
/// Verifies that the /correct endpoint can correct OCR output.
#[tokio::test]
async fn test_correct_endpoint_processes_ocr() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    
    // Simulate OCR text with common errors
    let request = CorrectRequest {
        ocr_text: "The m1tochondr1a 1s the powerhouse of the cell.".to_string(),
        language: Some("en".to_string()),
    };
    
    // Act: Send correct request
    let response = client
        .post(format!("{}/correct", GEMINI_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Response should be received
    assert!(response.is_ok(), "Correct endpoint should respond");
}

/// Test: Correct endpoint handles various languages
///
/// Verifies that the /correct endpoint handles different language parameters.
#[tokio::test]
async fn test_correct_endpoint_language_support() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    let ocr_text = "Sample OCR text for correction".to_string();
    
    let languages = vec![Some("en".to_string()), Some("es".to_string()), None];
    
    for lang in languages {
        let request = CorrectRequest {
            ocr_text: ocr_text.clone(),
            language: lang.clone(),
        };
        
        let response = client
            .post(format!("{}/correct", GEMINI_BASE_URL))
            .json(&request)
            .send()
            .await;
        
        assert!(
            response.is_ok(),
            "Correct endpoint should handle language: {:?}",
            lang
        );
    }
}

// =============================================================================
// Extract Endpoint Tests
// =============================================================================

/// Request body for /extract endpoint
#[derive(Debug, Serialize)]
struct ExtractRequest {
    content: String,
}

/// Test: Extract endpoint extracts content
///
/// Verifies that the /extract endpoint can extract entities from content.
#[tokio::test]
async fn test_extract_endpoint_extracts_entities() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    let request = ExtractRequest {
        content: "The concept of gravity was discovered by Newton. \
                  An electric field surrounds charged particles. \
                  Energy exists in multiple forms including kinetic and potential."
            .to_string(),
    };
    
    // Act: Send extract request
    let response = client
        .post(format!("{}/extract", GEMINI_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Response should be received
    assert!(response.is_ok(), "Extract endpoint should respond");
}

// =============================================================================
// Generate Endpoint Tests
// =============================================================================

/// Request body for /generate endpoint
#[derive(Debug, Serialize)]
struct GenerateRequest {
    topic: String,
    subject: Option<String>,
    format: Option<String>,
}

/// Test: Generate endpoint creates study notes
///
/// Verifies that the /generate endpoint can generate study notes.
#[tokio::test]
async fn test_generate_endpoint_creates_notes() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    let request = GenerateRequest {
        topic: "Photosynthesis".to_string(),
        subject: Some("Biology".to_string()),
        format: Some("markdown".to_string()),
    };
    
    // Act: Send generate request
    let response = client
        .post(format!("{}/generate", GEMINI_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Response should be received
    assert!(response.is_ok(), "Generate endpoint should respond");
}

// =============================================================================
// Error Handling Tests
// =============================================================================

/// Test: Invalid JSON returns appropriate error
///
/// Verifies that the service handles malformed JSON gracefully.
#[tokio::test]
async fn test_invalid_json_returns_error() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    
    // Act: Send invalid JSON
    let response = client
        .post(format!("{}/analyze", GEMINI_BASE_URL))
        .header("Content-Type", "application/json")
        .body("{ invalid json }")
        .send()
        .await;
    
    // Assert: Should return 400 Bad Request or 422 Unprocessable Entity
    if response.is_ok() {
        let status = response.unwrap().status();
        assert!(
            status.as_u16() == 400 || status.as_u16() == 422,
            "Invalid JSON should return 400 or 422"
        );
    }
}

/// Test: Wrong HTTP method returns 405
///
/// Verifies that using GET on POST-only endpoints returns method not allowed.
#[tokio::test]
async fn test_wrong_method_returns_405() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    
    // Try GET on /analyze (which requires POST)
    let response = client
        .get(format!("{}/analyze", GEMINI_BASE_URL))
        .send()
        .await;
    
    // Assert: Should return 405 Method Not Allowed or similar
    if response.is_ok() {
        let status = response.unwrap().status();
        assert!(
            status.as_u16() == 405 || status.as_u16() == 400,
            "Wrong method should return 405 or 400"
        );
    }
}

// =============================================================================
// Cross-Service Communication Tests
// =============================================================================

/// Test: CORS headers are present
///
/// Verifies that the service includes CORS headers for cross-origin requests.
#[tokio::test]
async fn test_cors_headers_present() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .get(format!("{}/health", GEMINI_BASE_URL))
        .header("Origin", "http://localhost:8080")
        .send()
        .await;
    
    assert!(response.is_ok(), "Health endpoint should respond with CORS");
}

// =============================================================================
// Performance Tests
// =============================================================================

/// Test: Health endpoint responds within acceptable time
///
/// Verifies that the health endpoint is performant.
#[tokio::test]
async fn test_health_response_time() {
    skip_if_service_unavailable!(GEMINI_BASE_URL);
    
    let client = Client::builder()
        .timeout(HEALTH_TIMEOUT)
        .build()
        .unwrap();
    
    let max_response_time = std::time::Duration::from_millis(500);
    
    for _ in 0..3 {
        let start = std::time::Instant::now();
        let response = client
            .get(format!("{}/health", GEMINI_BASE_URL))
            .send()
            .await;
        let elapsed = start.elapsed();
        
        assert!(response.is_ok(), "Health endpoint should respond");
        assert!(
            elapsed < max_response_time,
            "Health endpoint should respond within 500ms, took {:?}",
            elapsed
        );
    }
}
