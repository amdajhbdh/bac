//! Basic tests for cloud-tools crate

use cloud_tools::models::*;

#[test]
fn test_ssh_exec_result_success() {
    let result = SshExecResult::success("Hello, World!".to_string(), String::new(), 0);

    assert!(result.success);
    assert_eq!(result.stdout, Some("Hello, World!".to_string()));
    assert_eq!(result.stderr, Some(String::new()));
    assert_eq!(result.exit_code, Some(0));
    assert!(result.error.is_none());
}

#[test]
fn test_ssh_exec_result_error() {
    let result = SshExecResult::error("Command failed".to_string());

    assert!(!result.success);
    assert!(result.stdout.is_none());
    assert!(result.stderr.is_none());
    assert!(result.exit_code.is_none());
    assert_eq!(result.error, Some("Command failed".to_string()));
}

#[test]
fn test_upload_result_success() {
    let result = UploadResult::success("/remote/path/file.txt".to_string(), 1024);

    assert!(result.success);
    assert_eq!(
        result.remote_path,
        Some("/remote/path/file.txt".to_string())
    );
    assert_eq!(result.bytes_transferred, Some(1024));
    assert!(result.error.is_none());
}

#[test]
fn test_upload_result_error() {
    let result = UploadResult::error("File not found".to_string());

    assert!(!result.success);
    assert!(result.remote_path.is_none());
    assert!(result.bytes_transferred.is_none());
    assert_eq!(result.error, Some("File not found".to_string()));
}

#[test]
fn test_download_result_success() {
    let result = DownloadResult::success("/local/path/file.txt".to_string(), 2048);

    assert!(result.success);
    assert_eq!(result.local_path, Some("/local/path/file.txt".to_string()));
    assert_eq!(result.bytes_transferred, Some(2048));
    assert!(result.error.is_none());
}

#[test]
fn test_ocr_result_success() {
    let result = OcrResult::success("Extracted text".to_string(), 0.95);

    assert!(result.success);
    assert_eq!(result.text, Some("Extracted text".to_string()));
    assert!((result.confidence.unwrap() - 0.95).abs() < 0.001);
    assert!(result.error.is_none());
}

#[test]
fn test_ocr_result_error() {
    let result = OcrResult::error("OCR failed".to_string());

    assert!(!result.success);
    assert!(result.text.is_none());
    assert!(result.confidence.is_none());
    assert_eq!(result.error, Some("OCR failed".to_string()));
}

#[test]
fn test_gcs_sync_result_success() {
    let result = GcsSyncResult::success(10, 10240);

    assert!(result.success);
    assert_eq!(result.files_synced, Some(10));
    assert_eq!(result.bytes_transferred, Some(10240));
    assert!(result.error.is_none());
}

#[test]
fn test_gcs_sync_direction() {
    let upload = SyncDirection::Upload;
    let download = SyncDirection::Download;

    assert_eq!(upload, SyncDirection::Upload);
    assert_eq!(download, SyncDirection::Download);
}

#[test]
fn test_health_response_healthy() {
    let response = HealthResponse::healthy(true);

    assert_eq!(response.status, "ok");
    assert!(response.cloud_shell_connected);
    assert_eq!(response.service_version, env!("CARGO_PKG_VERSION"));
}

#[test]
fn test_health_response_degraded() {
    let response = HealthResponse::healthy(false);

    assert_eq!(response.status, "degraded");
    assert!(!response.cloud_shell_connected);
}

#[test]
fn test_cloud_tools_error_display() {
    use cloud_tools::models::CloudToolsError;

    let err = CloudToolsError::SshConnection("Connection refused".to_string());
    assert_eq!(err.to_string(), "SSH connection failed: Connection refused");

    let err = CloudToolsError::GcsOperation("Access denied".to_string());
    assert_eq!(err.to_string(), "GCS operation failed: Access denied");

    let err = CloudToolsError::Timeout("30 seconds".to_string());
    assert_eq!(err.to_string(), "Timeout: 30 seconds");
}

#[test]
fn test_ssh_exec_request_deserialize() {
    let json = r#"{"command": "ls -la", "timeout_secs": 60}"#;
    let request: SshExecRequest = serde_json::from_str(json).unwrap();

    assert_eq!(request.command, "ls -la");
    assert_eq!(request.timeout_secs, Some(60));
}

#[test]
fn test_upload_request_deserialize() {
    let json = r#"{"local_path": "/tmp/file.txt", "remote_path": "/home/user/file.txt"}"#;
    let request: UploadRequest = serde_json::from_str(json).unwrap();

    assert_eq!(request.local_path, "/tmp/file.txt");
    assert_eq!(request.remote_path, "/home/user/file.txt");
}

#[test]
fn test_ocr_request_deserialize() {
    let json = r#"{"image_path": "/tmp/image.png", "language": "eng"}"#;
    let request: OcrRequest = serde_json::from_str(json).unwrap();

    assert_eq!(request.image_path, "/tmp/image.png");
    assert_eq!(request.language, Some("eng".to_string()));
}

#[test]
fn test_gcs_sync_request_upload_deserialize() {
    let json = r#"{
        "bucket": "my-bucket",
        "direction": "upload",
        "local_path": "/tmp/data",
        "remote_path": "/data"
    }"#;
    let request: GcsSyncRequest = serde_json::from_str(json).unwrap();

    assert_eq!(request.bucket, "my-bucket");
    assert_eq!(request.direction, SyncDirection::Upload);
    assert_eq!(request.local_path, Some("/tmp/data".to_string()));
    assert_eq!(request.remote_path, Some("/data".to_string()));
}

#[test]
fn test_gcs_sync_request_download_deserialize() {
    let json = r#"{
        "bucket": "my-bucket",
        "direction": "download",
        "prefix": "exports/",
        "local_path": "/tmp/downloads"
    }"#;
    let request: GcsSyncRequest = serde_json::from_str(json).unwrap();

    assert_eq!(request.bucket, "my-bucket");
    assert_eq!(request.direction, SyncDirection::Download);
    assert_eq!(request.prefix, Some("exports/".to_string()));
    assert_eq!(request.local_path, Some("/tmp/downloads".to_string()));
}
