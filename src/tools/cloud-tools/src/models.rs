//! Data models for cloud-tools

use serde::{Deserialize, Serialize};

// ============================================================================
// SSH Models
// ============================================================================

/// SSH command execution request
#[derive(Debug, Deserialize)]
pub struct SshExecRequest {
    pub command: String,
    pub timeout_secs: Option<u64>,
}

/// SSH command execution result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SshExecResult {
    pub success: bool,
    pub stdout: Option<String>,
    pub stderr: Option<String>,
    pub exit_code: Option<i32>,
    pub error: Option<String>,
}

impl SshExecResult {
    pub fn success(stdout: String, stderr: String, exit_code: i32) -> Self {
        Self {
            success: true,
            stdout: Some(stdout),
            stderr: Some(stderr),
            exit_code: Some(exit_code),
            error: None,
        }
    }

    pub fn error(msg: String) -> Self {
        Self {
            success: false,
            stdout: None,
            stderr: None,
            exit_code: None,
            error: Some(msg),
        }
    }
}

// ============================================================================
// Upload Models
// ============================================================================

/// File upload request
#[derive(Debug, Deserialize)]
pub struct UploadRequest {
    pub local_path: String,
    pub remote_path: String,
}

/// File upload result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UploadResult {
    pub success: bool,
    pub remote_path: Option<String>,
    pub bytes_transferred: Option<u64>,
    pub error: Option<String>,
}

impl UploadResult {
    pub fn success(remote_path: String, bytes_transferred: u64) -> Self {
        Self {
            success: true,
            remote_path: Some(remote_path),
            bytes_transferred: Some(bytes_transferred),
            error: None,
        }
    }

    pub fn error(msg: String) -> Self {
        Self {
            success: false,
            remote_path: None,
            bytes_transferred: None,
            error: Some(msg),
        }
    }
}

// ============================================================================
// Download Models
// ============================================================================

/// File download request
#[derive(Debug, Deserialize)]
pub struct DownloadRequest {
    pub remote_path: String,
    pub local_path: String,
}

/// File download result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DownloadResult {
    pub success: bool,
    pub local_path: Option<String>,
    pub bytes_transferred: Option<u64>,
    pub error: Option<String>,
}

impl DownloadResult {
    pub fn success(local_path: String, bytes_transferred: u64) -> Self {
        Self {
            success: true,
            local_path: Some(local_path),
            bytes_transferred: Some(bytes_transferred),
            error: None,
        }
    }

    pub fn error(msg: String) -> Self {
        Self {
            success: false,
            local_path: None,
            bytes_transferred: None,
            error: Some(msg),
        }
    }
}

// ============================================================================
// OCR Models
// ============================================================================

/// OCR request
#[derive(Debug, Deserialize)]
pub struct OcrRequest {
    pub image_path: String,
    pub language: Option<String>,
}

/// OCR result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OcrResult {
    pub success: bool,
    pub text: Option<String>,
    pub confidence: Option<f32>,
    pub error: Option<String>,
}

impl OcrResult {
    pub fn success(text: String, confidence: f32) -> Self {
        Self {
            success: true,
            text: Some(text),
            confidence: Some(confidence),
            error: None,
        }
    }

    pub fn error(msg: String) -> Self {
        Self {
            success: false,
            text: None,
            confidence: None,
            error: Some(msg),
        }
    }
}

// ============================================================================
// GCS Sync Models
// ============================================================================

/// GCS sync direction
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum SyncDirection {
    Upload,
    Download,
}

/// GCS sync request
#[derive(Debug, Deserialize)]
pub struct GcsSyncRequest {
    pub bucket: String,
    pub prefix: Option<String>,
    pub direction: SyncDirection,
    pub local_path: Option<String>,
    pub remote_path: Option<String>,
}

/// GCS sync result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GcsSyncResult {
    pub success: bool,
    pub files_synced: Option<usize>,
    pub bytes_transferred: Option<u64>,
    pub error: Option<String>,
}

impl GcsSyncResult {
    pub fn success(files_synced: usize, bytes_transferred: u64) -> Self {
        Self {
            success: true,
            files_synced: Some(files_synced),
            bytes_transferred: Some(bytes_transferred),
            error: None,
        }
    }

    pub fn error(msg: String) -> Self {
        Self {
            success: false,
            files_synced: None,
            bytes_transferred: None,
            error: Some(msg),
        }
    }
}

// ============================================================================
// Health Models
// ============================================================================

/// Health check response
#[derive(Debug, Serialize)]
pub struct HealthResponse {
    pub status: String,
    pub cloud_shell_connected: bool,
    pub service_version: String,
}

impl HealthResponse {
    pub fn healthy(connected: bool) -> Self {
        Self {
            status: if connected { "ok" } else { "degraded" }.to_string(),
            cloud_shell_connected: connected,
            service_version: env!("CARGO_PKG_VERSION").to_string(),
        }
    }
}

// ============================================================================
// Error Types
// ============================================================================

/// Cloud tools error types
#[derive(Debug, thiserror::Error)]
pub enum CloudToolsError {
    #[error("SSH connection failed: {0}")]
    SshConnection(String),

    #[error("SSH command execution failed: {0}")]
    SshExecution(String),

    #[error("File transfer failed: {0}")]
    FileTransfer(String),

    #[error("Cloud Shell not available: {0}")]
    CloudShellUnavailable(String),

    #[error("GCS operation failed: {0}")]
    GcsOperation(String),

    #[error("Timeout: {0}")]
    Timeout(String),

    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),
}

impl serde::Serialize for CloudToolsError {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        serializer.serialize_str(&self.to_string())
    }
}
