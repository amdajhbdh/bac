//! File upload to Cloud Shell
//!
//! Handles uploading files from local machine to Cloud Shell instance.

use crate::models::UploadResult;
use std::path::Path;
use tokio::process::Command;

/// Upload file to Cloud Shell using gcloud scp
pub async fn upload(local_path: &str, remote_path: &str) -> UploadResult {
    tracing::info!("Uploading {} -> Cloud Shell:{}", local_path, remote_path);

    // Validate local file exists
    let path = Path::new(local_path);
    if !path.exists() {
        return UploadResult::error(format!("Local file not found: {}", local_path));
    }

    let metadata = match std::fs::metadata(path) {
        Ok(m) => m,
        Err(e) => return UploadResult::error(format!("Failed to read file metadata: {}", e)),
    };

    let file_size = metadata.len();

    // Use gcloud cloud-shell scp to upload
    // Format: gcloud cloud-shell scp <local> <remote>
    let result = Command::new("gcloud")
        .args([
            "cloud-shell",
            "scp",
            "--quiet",
            local_path,
            &format!(":~{}", remote_path),
        ])
        .output()
        .await;

    match result {
        Ok(output) if output.status.success() => {
            tracing::info!(
                "Upload successful: {} bytes transferred",
                file_size
            );
            UploadResult::success(remote_path.to_string(), file_size)
        }
        Ok(output) => {
            let stderr = String::from_utf8_lossy(&output.stderr);
            UploadResult::error(format!("Upload failed: {}", stderr))
        }
        Err(e) => UploadResult::error(format!("Failed to execute upload: {}", e)),
    }
}

/// Upload file with specific Cloud Shell instance
pub async fn upload_to_instance(
    local_path: &str,
    remote_path: &str,
    instance: &str,
) -> UploadResult {
    tracing::info!(
        "Uploading {} -> {}:{}",
        local_path,
        instance,
        remote_path
    );

    let path = Path::new(local_path);
    if !path.exists() {
        return UploadResult::error(format!("Local file not found: {}", local_path));
    }

    let file_size = match std::fs::metadata(path) {
        Ok(m) => m.len(),
        Err(e) => return UploadResult::error(format!("Failed to read file: {}", e)),
    };

    let result = Command::new("gcloud")
        .args([
            "compute",
            "scp",
            "--zone=us-central1-a",
            local_path,
            &format!("{}:{}", instance, remote_path),
        ])
        .output()
        .await;

    match result {
        Ok(output) if output.status.success() => {
            UploadResult::success(remote_path.to_string(), file_size)
        }
        Ok(output) => {
            let stderr = String::from_utf8_lossy(&output.stderr);
            UploadResult::error(format!("Upload failed: {}", stderr))
        }
        Err(e) => UploadResult::error(format!("Failed to execute upload: {}", e)),
    }
}

/// Upload directory to Cloud Shell
pub async fn upload_directory(local_path: &str, remote_path: &str) -> UploadResult {
    tracing::info!("Uploading directory {} -> Cloud Shell:{}", local_path, remote_path);

    let path = Path::new(local_path);
    if !path.is_dir() {
        return UploadResult::error(format!("Not a directory: {}", local_path));
    }

    let result = Command::new("gcloud")
        .args([
            "cloud-shell",
            "scp",
            "--recurse",
            "--quiet",
            local_path,
            &format!(":~{}", remote_path),
        ])
        .output()
        .await;

    match result {
        Ok(output) if output.status.success() => {
            // Count files transferred
            let count = count_files(path);
            UploadResult::success(remote_path.to_string(), count as u64)
        }
        Ok(output) => {
            let stderr = String::from_utf8_lossy(&output.stderr);
            UploadResult::error(format!("Directory upload failed: {}", stderr))
        }
        Err(e) => UploadResult::error(format!("Failed to execute upload: {}", e)),
    }
}

/// Count files in directory recursively
fn count_files(path: &Path) -> usize {
    let mut count = 0;
    if let Ok(entries) = std::fs::read_dir(path) {
        for entry in entries.flatten() {
            if entry.path().is_dir() {
                count += count_files(&entry.path());
            } else {
                count += 1;
            }
        }
    }
    count
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_count_files() {
        let count = count_files(Path::new("/tmp"));
        tracing::debug!("Files in /tmp: {}", count);
    }
}
