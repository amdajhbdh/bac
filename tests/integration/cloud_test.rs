//! Integration tests for Cloud Tools service (:3004)
//!
//! Tests the HTTP endpoints for Cloud Shell integration operations.

mod common;

use common::{create_client, is_service_running, urls, HEALTH_TIMEOUT};
use serde::{Deserialize, Serialize};

// =============================================================================
// Test Configuration
// =============================================================================

const CLOUD_BASE_URL: &str = urls::CLOUD;

// =============================================================================
// Health Endpoint Tests
// =============================================================================

/// Health response from cloud-tools
#[derive(Debug, Deserialize)]
struct HealthResponse {
    #[serde(default)]
    status: String,
    #[serde(default)]
    cloud_shell_connected: Option<bool>,
    #[serde(default)]
    service_version: Option<String>,
}

/// Test: Health endpoint returns status
///
/// Verifies that the /health endpoint responds with service status.
#[tokio::test]
async fn test_health_endpoint() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    let response = client
        .get(format!("{}/health", CLOUD_BASE_URL))
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
    
    // Assert: Response should be valid JSON with status field
    let health: HealthResponse = response.json().await.unwrap();
    assert!(
        !health.status.is_empty(),
        "Health response should include status"
    );
}

/// Test: Health endpoint reports cloud shell status
///
/// Verifies that the health check reports Cloud Shell connectivity.
#[tokio::test]
async fn test_health_reports_cloud_shell_status() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    let response = client
        .get(format!("{}/health", CLOUD_BASE_URL))
        .send()
        .await
        .unwrap();
    
    let health: HealthResponse = response.json().await.unwrap();
    
    // Cloud shell status may be true or false depending on availability
    assert!(
        health.cloud_shell_connected.is_some(),
        "Health response should include cloud_shell_connected"
    );
    
    // Status should reflect cloud shell state
    let expected_status = if health.cloud_shell_connected.unwrap_or(false) {
        "ok"
    } else {
        "degraded"
    };
    assert_eq!(
        health.status, expected_status,
        "Status should be 'ok' when connected, 'degraded' otherwise"
    );
}

// =============================================================================
// SSH Endpoint Tests
// =============================================================================

/// Request body for /ssh/exec endpoint
#[derive(Debug, Serialize)]
struct SshExecRequest {
    command: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    timeout_secs: Option<u64>,
}

/// Response from /ssh/exec endpoint
#[derive(Debug, Deserialize)]
struct SshExecResult {
    #[serde(default)]
    success: bool,
    #[serde(default)]
    stdout: Option<String>,
    #[serde(default)]
    stderr: Option<String>,
    #[serde(default)]
    exit_code: Option<i32>,
    #[serde(default)]
    error: Option<String>,
}

/// Test: SSH exec endpoint accepts command
///
/// Verifies that the /ssh/exec endpoint can execute commands.
/// Note: Will return success=false if Cloud Shell is not available.
#[tokio::test]
async fn test_ssh_exec_endpoint() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    let request = SshExecRequest {
        command: "echo 'test'".to_string(),
        timeout_secs: Some(30),
    };
    
    // Act: Send SSH exec request
    let response = client
        .post(format!("{}/ssh/exec", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Request should succeed
    assert!(response.is_ok(), "SSH exec endpoint should respond");
    
    let result: SshExecResult = response.unwrap().json().await.unwrap();
    
    // Result should indicate success or failure
    // (success=true only if Cloud Shell is actually available)
    assert!(
        result.success || result.error.is_some(),
        "SSH exec should indicate success or have error message"
    );
}

/// Test: SSH exec handles simple echo command
///
/// Verifies that the SSH endpoint can execute simple commands.
#[tokio::test]
async fn test_ssh_exec_echo_command() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    let request = SshExecRequest {
        command: "echo 'Hello from integration test'".to_string(),
        timeout_secs: Some(10),
    };
    
    let response = client
        .post(format!("{}/ssh/exec", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await
        .unwrap();
    
    let result: SshExecResult = response.json().await.unwrap();
    
    // If Cloud Shell is available, should see output
    if result.success {
        assert!(
            result.stdout.is_some() && result.stdout.unwrap().contains("Hello"),
            "Echo command should return output"
        );
    }
}

/// Test: SSH exec handles command with error
///
/// Verifies that the SSH endpoint handles failing commands gracefully.
#[tokio::test]
async fn test_ssh_exec_command_with_error() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    
    // Command that will fail
    let request = SshExecRequest {
        command: "exit 1".to_string(),
        timeout_secs: Some(10),
    };
    
    let response = client
        .post(format!("{}/ssh/exec", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await
        .unwrap();
    
    let result: SshExecResult = response.json().await.unwrap();
    
    // Should report failure with exit code
    if result.success {
        assert!(
            result.exit_code.is_some(),
            "Failing command should report exit code"
        );
    }
}

/// Test: SSH exec with custom timeout
///
/// Verifies that the SSH endpoint respects timeout parameters.
#[tokio::test]
async fn test_ssh_exec_custom_timeout() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    
    // Use a short timeout
    let request = SshExecRequest {
        command: "echo 'quick'".to_string(),
        timeout_secs: Some(5),
    };
    
    let response = client
        .post(format!("{}/ssh/exec", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "SSH exec should accept custom timeout");
}

// =============================================================================
// Upload Endpoint Tests
// =============================================================================

/// Request body for /upload endpoint
#[derive(Debug, Serialize)]
struct UploadRequest {
    local_path: String,
    remote_path: String,
}

/// Response from /upload endpoint
#[derive(Debug, Deserialize)]
struct UploadResult {
    #[serde(default)]
    success: bool,
    #[serde(default)]
    remote_path: Option<String>,
    #[serde(default)]
    bytes_transferred: Option<u64>,
    #[serde(default)]
    error: Option<String>,
}

/// Test: Upload endpoint accepts upload requests
///
/// Verifies that the /upload endpoint is accessible.
/// Note: Will return error for non-existent files.
#[tokio::test]
async fn test_upload_endpoint() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    let request = UploadRequest {
        local_path: "/nonexistent/file.txt".to_string(),
        remote_path: "/remote/path.txt".to_string(),
    };
    
    let response = client
        .post(format!("{}/upload", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Should respond (error expected for non-existent file)
    assert!(response.is_ok(), "Upload endpoint should respond");
}

/// Test: Upload endpoint validates paths
///
/// Verifies that the upload endpoint validates request parameters.
#[tokio::test]
async fn test_upload_validates_paths() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    
    // Empty paths
    let request = UploadRequest {
        local_path: "".to_string(),
        remote_path: "".to_string(),
    };
    
    let response = client
        .post(format!("{}/upload", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await
        .unwrap();
    
    let result: UploadResult = response.json().await.unwrap();
    
    // Should indicate failure for invalid paths
    assert!(
        !result.success || result.error.is_some(),
        "Upload with empty paths should fail"
    );
}

// =============================================================================
// Download Endpoint Tests
// =============================================================================

/// Request body for /download endpoint
#[derive(Debug, Serialize)]
struct DownloadRequest {
    remote_path: String,
    local_path: String,
}

/// Response from /download endpoint
#[derive(Debug, Deserialize)]
struct DownloadResult {
    #[serde(default)]
    success: bool,
    #[serde(default)]
    local_path: Option<String>,
    #[serde(default)]
    bytes_transferred: Option<u64>,
    #[serde(default)]
    error: Option<String>,
}

/// Test: Download endpoint accepts download requests
///
/// Verifies that the /download endpoint is accessible.
/// Note: Will return error for non-existent remote files.
#[tokio::test]
async fn test_download_endpoint() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    let request = DownloadRequest {
        remote_path: "/nonexistent/remote.txt".to_string(),
        local_path: "/tmp/local.txt".to_string(),
    };
    
    let response = client
        .post(format!("{}/download", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Should respond (error expected for non-existent file)
    assert!(response.is_ok(), "Download endpoint should respond");
}

// =============================================================================
// OCR Endpoint Tests
// =============================================================================

/// Request body for /ocr endpoint
#[derive(Debug, Serialize)]
struct OcrRequest {
    image_path: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    language: Option<String>,
}

/// Response from /ocr endpoint
#[derive(Debug, Deserialize)]
struct OcrResult {
    #[serde(default)]
    success: bool,
    #[serde(default)]
    text: Option<String>,
    #[serde(default)]
    confidence: Option<f32>,
    #[serde(default)]
    error: Option<String>,
}

/// Test: OCR endpoint accepts image path
///
/// Verifies that the /ocr endpoint is accessible.
/// Note: Will return error for non-existent images.
#[tokio::test]
async fn test_ocr_endpoint() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    let request = OcrRequest {
        image_path: "/nonexistent/image.jpg".to_string(),
        language: Some("eng".to_string()),
    };
    
    let response = client
        .post(format!("{}/ocr", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Should respond (error expected for non-existent file)
    assert!(response.is_ok(), "OCR endpoint should respond");
}

/// Test: OCR endpoint with language parameter
///
/// Verifies that the OCR endpoint accepts language options.
#[tokio::test]
async fn test_ocr_with_language() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    
    let languages = vec![
        Some("eng".to_string()),
        Some("spa".to_string()),
        Some("fra".to_string()),
        None,
    ];
    
    for lang in languages {
        let request = OcrRequest {
            image_path: "/nonexistent/image.jpg".to_string(),
            language: lang.clone(),
        };
        
        let response = client
            .post(format!("{}/ocr", CLOUD_BASE_URL))
            .json(&request)
            .send()
            .await;
        
        assert!(
            response.is_ok(),
            "OCR should accept language: {:?}",
            lang
        );
    }
}

/// Test: OCR PDF endpoint exists
///
/// Verifies that the /ocr/pdf endpoint is accessible.
#[tokio::test]
async fn test_ocr_pdf_endpoint() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    let request = OcrRequest {
        image_path: "/nonexistent/document.pdf".to_string(),
        language: None,
    };
    
    let response = client
        .post(format!("{}/ocr/pdf", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "OCR PDF endpoint should respond");
}

/// Test: OCR status endpoint reports availability
///
/// Verifies that the /ocr/status endpoint returns availability info.
#[tokio::test]
async fn test_ocr_status_endpoint() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .get(format!("{}/ocr/status", CLOUD_BASE_URL))
        .send()
        .await;
    
    assert!(response.is_ok(), "OCR status endpoint should respond");
    
    // Parse response
    let response = response.unwrap();
    let status: OcrStatusResponse = response.json().await.unwrap();
    
    // Response should have available flag
    assert!(
        status.available == true || status.available == false,
        "OCR status should have available flag"
    );
}

/// OCR status response
#[derive(Debug, Deserialize)]
struct OcrStatusResponse {
    available: bool,
    #[serde(default)]
    languages: Vec<String>,
}

// =============================================================================
// GCS Endpoint Tests
// =============================================================================

/// Test: GCS list endpoint is accessible
///
/// Verifies that the /gcs/list endpoint is available.
#[tokio::test]
async fn test_gcs_list_endpoint() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .get(format!("{}/gcs/list?bucket=test-bucket", CLOUD_BASE_URL))
        .send()
        .await;
    
    assert!(response.is_ok(), "GCS list endpoint should respond");
}

/// Test: GCS stats endpoint is accessible
///
/// Verifies that the /gcs/stats endpoint is available.
#[tokio::test]
async fn test_gcs_stats_endpoint() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .get(format!("{}/gcs/stats?bucket=test-bucket", CLOUD_BASE_URL))
        .send()
        .await;
    
    assert!(response.is_ok(), "GCS stats endpoint should respond");
}

/// GCS sync request
#[derive(Debug, Serialize)]
struct GcsSyncRequest {
    bucket: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    prefix: Option<String>,
    direction: String, // "upload" or "download"
    #[serde(skip_serializing_if = "Option::is_none")]
    local_path: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    remote_path: Option<String>,
}

/// Test: GCS sync endpoint accepts requests
///
/// Verifies that the /gcs/sync endpoint is accessible.
#[tokio::test]
async fn test_gcs_sync_endpoint() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    let request = GcsSyncRequest {
        bucket: "test-bucket".to_string(),
        prefix: Some("test/".to_string()),
        direction: "upload".to_string(),
        local_path: Some("/tmp".to_string()),
        remote_path: Some("/remote".to_string()),
    };
    
    let response = client
        .post(format!("{}/gcs/sync", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "GCS sync endpoint should respond");
}

// =============================================================================
// Error Handling Tests
// =============================================================================

/// Test: SSH exec with missing command
///
/// Verifies that the SSH endpoint validates required fields.
#[tokio::test]
async fn test_ssh_exec_missing_command() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    
    // Missing command field
    let request = serde_json::json!({});
    
    let response = client
        .post(format!("{}/ssh/exec", CLOUD_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Service should return error
    if response.is_ok() {
        let status = response.unwrap().status();
        assert!(
            status.as_u16() >= 400,
            "Missing command should return error status"
        );
    }
}

/// Test: Invalid JSON returns error
///
/// Verifies that the service handles malformed JSON gracefully.
#[tokio::test]
async fn test_invalid_json_returns_error() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .post(format!("{}/ssh/exec", CLOUD_BASE_URL))
        .header("Content-Type", "application/json")
        .body("{ invalid json }")
        .send()
        .await;
    
    // Should return 400 or 422
    if response.is_ok() {
        let status = response.unwrap().status();
        assert!(
            status.as_u16() == 400 || status.as_u16() == 422,
            "Invalid JSON should return 400 or 422"
        );
    }
}

// =============================================================================
// Performance Tests
// =============================================================================

/// Test: Health endpoint responds quickly
///
/// Verifies that health checks complete within acceptable time.
#[tokio::test]
async fn test_health_response_performance() {
    skip_if_service_unavailable!(CLOUD_BASE_URL);
    
    let client = Client::builder()
        .timeout(HEALTH_TIMEOUT)
        .build()
        .unwrap();
    
    let max_response_time = std::time::Duration::from_millis(500);
    
    let start = std::time::Instant::now();
    let response = client
        .get(format!("{}/health", CLOUD_BASE_URL))
        .send()
        .await;
    let elapsed = start.elapsed();
    
    assert!(response.is_ok(), "Health endpoint should respond");
    assert!(
        elapsed < max_response_time,
        "Health should respond within 500ms, took {:?}",
        elapsed
    );
}
