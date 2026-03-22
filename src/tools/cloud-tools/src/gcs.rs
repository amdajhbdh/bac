//! GCS (Google Cloud Storage) sync operations
//!
//! Syncs files between Cloud Shell and GCS buckets using gsutil.

use crate::models::{CloudToolsError, GcsSyncResult, SyncDirection};
use crate::ssh::ssh_exec;
use std::path::Path;

/// Sync files with GCS bucket
pub async fn sync_gcs(
    bucket: &str,
    direction: SyncDirection,
    local_path: Option<&str>,
    remote_path: Option<&str>,
    prefix: Option<&str>,
) -> GcsSyncResult {
    tracing::info!(
        "GCS sync: {} -> bucket {} ({})",
        direction_to_string(&direction),
        bucket,
        prefix.unwrap_or("")
    );

    // Validate inputs
    let bucket_url = format!("gs://{}", bucket);

    let result = match direction {
        SyncDirection::Upload => {
            let src = local_path.ok_or_else(|| GcsSyncResult::error("local_path required for upload".to_string()));
            let dst_path = remote_path.or(prefix).unwrap_or("");
            upload_to_gcs(src.unwrap(), &bucket_url, dst_path).await
        }
        SyncDirection::Download => {
            let dst = local_path.ok_or_else(|| GcsSyncResult::error("local_path required for download".to_string()));
            let src_path = remote_path.or(prefix).unwrap_or("");
            download_from_gcs(&bucket_url, src_path, dst.unwrap()).await
        }
    };

    result.unwrap_or_else(|e| GcsSyncResult::error(e.to_string()))
}

/// Upload files to GCS bucket
async fn upload_to_gcs(
    local_path: &str,
    bucket_url: &str,
    remote_path: &str,
) -> Result<GcsSyncResult, CloudToolsError> {
    let path = Path::new(local_path);
    if !path.exists() {
        return Err(CloudToolsError::FileTransfer(format!(
            "Local path not found: {}",
            local_path
        )));
    }

    let dest = if remote_path.is_empty() {
        bucket_url.to_string()
    } else {
        format!("{}/{}", bucket_url, remote_path.trim_start_matches('/'))
    };

    // Build gsutil command
    let gsutil_cmd = if path.is_dir() {
        format!("gsutil -m cp -r \"{}\" \"{}\"", local_path, dest)
    } else {
        format!("gsutil cp \"{}\" \"{}\"", local_path, dest)
    };

    let result = ssh_exec(&gsutil_cmd, Some(300)).await;

    if result.success {
        let files = count_files_in_path(path);
        Ok(GcsSyncResult::success(files, 0)) // Bytes tracking would need stat
    } else {
        Err(CloudToolsError::GcsOperation(
            result.stderr.unwrap_or_else(|| "Upload failed".to_string()),
        ))
    }
}

/// Download files from GCS bucket
async fn download_from_gcs(
    bucket_url: &str,
    remote_path: &str,
    local_path: &str,
) -> Result<GcsSyncResult, CloudToolsError> {
    let source = if remote_path.is_empty() {
        bucket_url.to_string()
    } else {
        format!("{}/{}", bucket_url, remote_path.trim_start_matches('/'))
    };

    let gsutil_cmd = format!("gsutil -m cp -r \"{}\" \"{}\"", source, local_path);

    let result = ssh_exec(&gsutil_cmd, Some(300)).await;

    if result.success {
        let files = count_files_in_path(Path::new(local_path));
        Ok(GcsSyncResult::success(files, 0))
    } else {
        Err(CloudToolsError::GcsOperation(
            result.stderr.unwrap_or_else(|| "Download failed".to_string()),
        ))
    }
}

/// List contents of a GCS bucket
pub async fn list_bucket(bucket: &str, prefix: Option<&str>) -> Result<Vec<String>, CloudToolsError> {
    let bucket_url = match prefix {
        Some(p) => format!("gs://{}/{}", bucket, p.trim_start_matches('/')),
        None => format!("gs://{}", bucket),
    };

    let result = ssh_exec(&format!("gsutil ls \"{}\"", bucket_url), Some(60)).await;

    if result.success {
        let files: Vec<String> = result
            .stdout
            .unwrap_or_default()
            .lines()
            .filter(|l| !l.is_empty())
            .map(|s| s.to_string())
            .collect();
        Ok(files)
    } else {
        Err(CloudToolsError::GcsOperation(
            result.stderr.unwrap_or_else(|| "List failed".to_string()),
        ))
    }
}

/// Check if a bucket exists and is accessible
pub async fn bucket_exists(bucket: &str) -> bool {
    let result = ssh_exec(&format!("gsutil ls -b gs://{}", bucket), Some(30)).await;
    result.success
}

/// Get bucket size and file count
pub async fn bucket_stats(bucket: &str) -> Result<(u64, usize), CloudToolsError> {
    let result = ssh_exec(
        &format!("gsutil du -s gs://{}", bucket),
        Some(60),
    )
    .await;

    if result.success {
        let output = result.stdout.unwrap_or_default();
        // Parse "12345 gs://bucket" format
        if let Some(size_str) = output.split_whitespace().next() {
            let size: u64 = size_str.parse().unwrap_or(0);
            // Count files
            let count_result = ssh_exec(
                &format!("gsutil ls gs://{}/** 2>/dev/null | wc -l", bucket),
                Some(60),
            )
            .await;
            let count: usize = count_result
                .stdout
                .unwrap_or_default()
                .trim()
                .parse()
                .unwrap_or(0);
            return Ok((size, count));
        }
    }

    Err(CloudToolsError::GcsOperation(
        result.stderr.unwrap_or_else(|| "Stats failed".to_string()),
    ))
}

/// Count files in path
fn count_files_in_path(path: &Path) -> usize {
    if path.is_file() {
        return 1;
    }

    let mut count = 0;
    if let Ok(entries) = std::fs::read_dir(path) {
        for entry in entries.flatten() {
            if entry.path().is_dir() {
                count += count_files_in_path(&entry.path());
            } else {
                count += 1;
            }
        }
    }
    count
}

/// Convert direction to string
fn direction_to_string(direction: &SyncDirection) -> &'static str {
    match direction {
        SyncDirection::Upload => "upload",
        SyncDirection::Download => "download",
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_count_files() {
        let count = count_files_in_path(Path::new("/tmp"));
        tracing::debug!("Files in /tmp: {}", count);
    }

    #[tokio::test]
    async fn test_bucket_exists() {
        let exists = bucket_exists("nonexistent-bucket-12345xyz").await;
        tracing::info!("Bucket exists: {}", exists);
        // Should be false for non-existent bucket
    }
}
