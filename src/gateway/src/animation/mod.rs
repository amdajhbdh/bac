//! Animation Module - Manim Video Generation
//!
//! Provides animation generation using Manim via Podman

use serde::{Deserialize, Serialize};
use std::collections::VecDeque;
use std::sync::Arc;
use tokio::sync::RwLock;
use uuid::Uuid;

/// Animation job status
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum JobStatus {
    Queued,
    Processing,
    Completed,
    Failed,
}

/// Animation job
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnimationJob {
    pub id: String,
    pub code: String,
    pub status: JobStatus,
    pub created_at: i64,
    pub started_at: Option<i64>,
    pub completed_at: Option<i64>,
    pub output_path: Option<String>,
    pub error: Option<String>,
}

/// Animation request
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnimationRequest {
    pub code: String,
    pub quality: Option<String>,
    pub format: Option<String>,
}

/// Animation response
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnimationResponse {
    pub job_id: String,
    pub status: JobStatus,
    pub output_url: Option<String>,
    pub error: Option<String>,
}

/// Manim configuration
#[derive(Debug, Clone)]
pub struct ManimConfig {
    pub quality: String,
    pub format: String,
    pub preview: bool,
    pub timeout: u64,
}

impl Default for ManimConfig {
    fn default() -> Self {
        Self {
            quality: "medium_quality".to_string(),
            format: "mp4".to_string(),
            preview: false,
            timeout: 300,
        }
    }
}

/// Animation queue
pub struct AnimationQueue {
    jobs: RwLock<VecDeque<AnimationJob>>,
    max_size: usize,
}

impl AnimationQueue {
    pub fn new(max_size: usize) -> Self {
        Self {
            jobs: RwLock::new(VecDeque::new()),
            max_size,
        }
    }

    pub async fn enqueue(&self, code: &str) -> String {
        let id = Uuid::new_v4().to_string();
        let job = AnimationJob {
            id: id.clone(),
            code: code.to_string(),
            status: JobStatus::Queued,
            created_at: chrono::Utc::now().timestamp(),
            started_at: None,
            completed_at: None,
            output_path: None,
            error: None,
        };

        let mut jobs = self.jobs.write().await;
        if jobs.len() >= self.max_size {
            jobs.pop_front();
        }
        jobs.push_back(job);
        
        id
    }

    pub async fn dequeue(&self) -> Option<AnimationJob> {
        let mut jobs = self.jobs.write().await;
        jobs.pop_front()
    }

    pub async fn get_status(&self, id: &str) -> Option<JobStatus> {
        let jobs = self.jobs.read().await;
        jobs.iter().find(|j| j.id == id).map(|j| j.status.clone())
    }

    pub async fn list_jobs(&self) -> Vec<AnimationJob> {
        let jobs = self.jobs.read().await;
        jobs.iter().cloned().collect()
    }
}

/// Manim bridge using Podman
pub struct ManimBridge {
    config: ManimConfig,
    queue: Arc<AnimationQueue>,
    output_dir: String,
}

impl ManimBridge {
    pub fn new(config: ManimConfig, output_dir: String) -> Self {
        Self {
            config,
            queue: Arc::new(AnimationQueue::new(100)),
            output_dir,
        }
    }

    /// Render animation using Manim via Podman
    pub async fn render(&self, code: &str) -> Result<String, String> {
        let job_id = self.queue.enqueue(code).await;
        
        // Create temp file with Manim code
        let temp_dir = std::env::temp_dir();
        let manim_file = temp_dir.join(format!("anim_{}.py", job_id));
        
        let full_code = format!(r#"
from manim import *

{}
"#, code);
        
        std::fs::write(&manim_file, &full_code)
            .map_err(|e| format!("Failed to write code: {}", e))?;

        let manim_path = manim_file.to_string_lossy().to_string();
        
        // Run Manim via Podman
        let output = self.run_manim(&manim_path).await?;
        
        // Cleanup temp file
        let _ = std::fs::remove_file(&manim_file);
        
        Ok(output)
    }

    async fn run_manim(&self, code_path: &str) -> Result<String, String> {
        let work_dir = std::path::Path::new(code_path).parent().unwrap_or(std::path::Path::new("."));
        let work_dir_str = work_dir.display().to_string();
        
        // Create output directory in the work dir to avoid permission issues
        let output_subdir = format!("{}/manim_output", work_dir_str);

        // Ensure output directory exists with proper permissions
        std::fs::create_dir_all(&output_subdir)
            .map_err(|e| format!("Failed to create output dir: {}", e))?;
        
        // Make it writable
        let _ = std::fs::set_permissions(&output_subdir, std::os::unix::fs::PermissionsExt::from_mode(0o777));

        let result = tokio::process::Command::new("podman")
            .args([
                "run", "--rm",
                "-v", &format!("{}:/workspace", work_dir_str),
                "-v", &format!("{}:/output", output_subdir),
                "-w", "/workspace",
                "-e", &format!("MANIM_OUTPUT_TO_RENDER=False"),
                "docker.io/manimcommunity/manim:latest",
                "manim",
                code_path,
                "-ql",
                "--media_dir", "/output",
                "--format", &self.config.format,
            ])
            .output()
            .await
            .map_err(|e| format!("Podman error: {}", e))?;

        if result.status.success() {
            // Find the output file
            let output_dir = std::path::Path::new(&output_subdir);
            if let Ok(entries) = std::fs::read_dir(output_dir) {
                for entry in entries.flatten() {
                    let path = entry.path();
                    if let Some(ext) = path.extension() {
                        if ext == "mp4" || ext == "webm" || ext == "gif" {
                            // Copy to final output location
                            let final_path = format!("{}/animation_{}.{}",
                                self.output_dir,
                                Uuid::new_v4(),
                                ext.display()
                            );
                            std::fs::copy(&path, &final_path)
                                .map_err(|e| format!("Failed to copy output: {}", e))?;
                            return Ok(final_path);
                        }
                    }
                }
            }
            Err("No output file found".to_string())
        } else {
            let stderr = String::from_utf8_lossy(&result.stderr);
            Err(format!("Manim failed: {}", stderr))
        }
    }

    pub fn queue(&self) -> &Arc<AnimationQueue> {
        &self.queue
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_queue_enqueue() {
        let queue = AnimationQueue::new(10);
        let id = queue.enqueue("test code").await;
        assert!(!id.is_empty());
    }

    #[tokio::test]
    async fn test_queue_max_size() {
        let queue = AnimationQueue::new(2);
        queue.enqueue("code1").await;
        queue.enqueue("code2").await;
        queue.enqueue("code3").await;
        
        let jobs = queue.list_jobs().await;
        assert_eq!(jobs.len(), 2);
    }
}
