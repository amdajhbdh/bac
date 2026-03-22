//! File download from Cloud Shell
//!
//! Handles downloading files from Cloud Shell instance to local machine.

use crate::models::DownloadResult;
use std::path::Path;
use tokio::fs;
use tokio::process::Command;

/// Download file from Cloud Shell using gcloud scp
pub async fn download(remote_path: &str, local_path: &str) -> DownloadResult {
    tracing::info!("Downloading Cloud Shell:{} -> {}", remote_path, local_path);

    // Ensure parent directory exists
    if let Some(parent) = Path::new(local_path).parent() {
        if !parent.exists() {
            if let Err(e) = fs::create_dir_all(parent).await {
                return DownloadResult::error(format!("Failed to create directory: {}", e));
            }
        }
    }

    // Use gcloud cloud-shell scp to download
    // Format: gcloud cloud-shell scp <remote> <local>
    let result = Command::new("gcloud")
        .args([
            "cloud-shell",
            "scp",
            "--quiet",
            &format!(":~{}", remote_path),
            local_path,
        ])
        .output()
        .await;

    match result {
        Ok(output) if output.status.success() => {
            // Get file size
            let file_size = match std::fs::metadata(local_path) {
                Ok(m) => m.len(),
                Err(_) => 0,
            };

            tracing::info!("Download successful: {} bytes", file_size);
            DownloadResult::success(local_path.to_string(), file_size)
        }
        Ok(output) => {
            let stderr = String::from_utf8_lossy(&output.stderr);
            DownloadResult::error(format!("Download failed: {}", stderr))
        }
        Err(e) => DownloadResult::error(format!("Failed to execute download: {}", e)),
    }
}

/// Download file from specific Cloud Shell instance
pub async fn download_from_instance(
    remote_path: &str,
    local_path: &str,
    instance: &str,
) -> DownloadResult {
    tracing::info!(
        "Downloading {}:{} -> {}",
        instance,
        remote_path,
        local_path
    );

    // Ensure parent directory exists
    if let Some(parent) = Path::new(local_path).parent() {
        if !parent.exists() {
            if let Err(e) = fs::create_dir_all(parent).await {
                return DownloadResult::error(format!("Failed to create directory: {}", e));
            }
        }
    }

    let result = Command::new("gcloud")
        .args([
            "compute",
            "scp",
            "--zone=us-central1-a",
            &format!("{}:{}", instance, remote_path),
            local_path,
        ])
        .output()
        .await;

    match result {
        Ok(output) if output.status.success() => {
            let file_size = match std::fs::metadata(local_path) {
                Ok(m) => m.len(),
                Err(_) => 0,
            };
            DownloadResult::success(local_path.to_string(), file_size)
        }
        Ok(output) => {
            let stderr = String::from_utf8_lossy(&output.stderr);
            DownloadResult::error(format!("Download failed: {}", stderr))
        }
        Err(e) => DownloadResult::error(format!("Failed to execute download: {}", e)),
    }
}

/// Download directory from Cloud Shell
pub async fn download_directory(remote_path: &str, local_path: &str) -> DownloadResult {
    tracing::info!(
        "Downloading directory Cloud Shell:{} -> {}",
        remote_path,
        local_path
    );

    // Ensure parent directory exists
    if let Some(parent) = Path::new(local_path).parent() {
        if !parent.exists() {
            if let Err(e) = fs::create_dir_all(parent).await {
                return DownloadResult::error(format!("Failed to create directory: {}", e));
            }
        }
    }

    let result = Command::new("gcloud")
        .args([
            "cloud-shell",
            "scp",
            "--recurse",
            "--quiet",
            &format!(":~{}", remote_path),
            local_path,
        ])
        .output()
        .await;

    match result {
        Ok(output) if output.status.success() => {
            // Count downloaded files
            let count = count_files(Path::new(local_path));
            tracing::info!("Directory download successful: {} files", count);
            DownloadResult::success(local_path.to_string(), count as u64)
        }
        Ok(output) => {
            let stderr = String::from_utf8_lossy(&output.stderr);
            DownloadResult::error(format!("Directory download failed: {}", stderr))
        }
        Err(e) => DownloadResult::error(format!("Failed to execute download: {}", e)),
    }
}

/// Count files in directory
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
