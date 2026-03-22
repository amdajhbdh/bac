//! SSH client for Cloud Shell operations
//!
//! Handles SSH connection to Cloud Shell instance with connection pooling.

use crate::models::{CloudToolsError, SshExecResult};
use std::process::Stdio;
use std::sync::Arc;
use tokio::process::Command;
use tokio::sync::RwLock;

/// SSH session wrapper with Cloud Shell specific logic
pub struct CloudShellSession {
    connected: bool,
}

/// Connection pool for SSH sessions
pub struct SshConnectionPool {
    sessions: Arc<RwLock<Vec<CloudShellSession>>>,
    max_connections: usize,
}

impl Default for SshConnectionPool {
    fn default() -> Self {
        Self::new(5)
    }
}

impl SshConnectionPool {
    /// Create new connection pool
    pub fn new(max_connections: usize) -> Self {
        Self {
            sessions: Arc::new(RwLock::new(Vec::new())),
            max_connections,
        }
    }

    /// Get available session or create new one
    pub async fn get_session(&self) -> Result<CloudShellSession, CloudToolsError> {
        let mut sessions = self.sessions.write().await;
        
        // Try to reuse an existing session
        if let Some(session) = sessions.pop() {
            if session.connected {
                return Ok(session);
            }
        }

        // Create new session
        let session = self.connect().await?;
        Ok(session)
    }

    /// Return session to pool
    pub async fn return_session(&self, session: CloudShellSession) {
        if self.sessions.read().await.len() < self.max_connections {
            self.sessions.write().await.push(session);
        }
    }

    /// Connect to Cloud Shell via gcloud
    async fn connect(&self) -> Result<CloudShellSession, CloudToolsError> {
        // Check if gcloud is available and authenticated
        let check = Command::new("gcloud")
            .args(["auth", "list", "--filter=status:ACTIVE", "--format=value(account)"])
            .output()
            .await;

        match check {
            Ok(output) if output.status.success() => {
                let account = String::from_utf8_lossy(&output.stdout).trim().to_string();
                if account.is_empty() {
                    return Err(CloudToolsError::CloudShellUnavailable(
                        "No active gcloud authentication".to_string(),
                    ));
                }
                tracing::info!("Authenticated as: {}", account);
            }
            Ok(_) => {
                return Err(CloudToolsError::CloudShellUnavailable(
                    "gcloud not authenticated".to_string(),
                ));
            }
            Err(e) => {
                return Err(CloudToolsError::CloudShellUnavailable(format!(
                    "gcloud CLI not available: {}",
                    e
                )));
            }
        }

        Ok(CloudShellSession { connected: true })
    }
}

/// Execute command on Cloud Shell
pub async fn ssh_exec(command: &str, timeout_secs: Option<u64>) -> SshExecResult {
    let timeout = timeout_secs.unwrap_or(60);

    tracing::debug!("Executing on Cloud Shell: {}", command);

    // Use gcloud cloud-shell ssh with --command flag
    let child = match Command::new("gcloud")
        .args([
            "cloud-shell",
            "ssh",
            "--authorize-session",
            "--quiet",
            "--",
            command,
        ])
        .stdin(Stdio::null())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
    {
        Ok(c) => c,
        Err(e) => {
            return SshExecResult::error(format!("Failed to spawn gcloud command: {}", e));
        }
    };

    let timeout_duration = std::time::Duration::from_secs(timeout);

    // Use a oneshot channel to communicate between spawned task and timeout handler
    let (tx, rx) = tokio::sync::oneshot::channel();

    // Spawn a task that waits for the child and sends the result
    // We use Arc<Mutex<Option<Child>>> to allow killing the process from outside
    use tokio::sync::Mutex;
    use std::sync::Arc;
    let child = Arc::new(Mutex::new(Some(child)));
    let child_for_task = child.clone();
    
    tokio::spawn(async move {
        let mut child_guard = child_for_task.lock().await;
        if let Some(c) = child_guard.take() {
            let output = c.wait_with_output().await;
            let _ = tx.send(output);
        }
    });

    // Wait for either the result or timeout
    let output = tokio::select! {
        result = rx => {
            match result {
                Ok(Ok(output)) => output,
                Ok(Err(e)) => {
                    return SshExecResult::error(format!("Failed to wait for output: {}", e));
                }
                Err(_) => {
                    // Sender dropped, process was killed
                    return SshExecResult::error(format!("Process was killed"));
                }
            }
        }
        _ = tokio::time::sleep(timeout_duration) => {
            // Kill the process on timeout
            let mut child_guard = child.lock().await;
            if let Some(c) = child_guard.as_mut() {
                let _ = c.kill().await;
            }
            return SshExecResult::error(format!("Command timed out after {} seconds", timeout));
        }
    };

    let stdout = String::from_utf8_lossy(&output.stdout).to_string();
    let stderr = String::from_utf8_lossy(&output.stderr).to_string();
    let exit_code = output.status.code().unwrap_or(-1);

    tracing::debug!(
        "Command completed: exit_code={}, stdout_len={}, stderr_len={}",
        exit_code,
        stdout.len(),
        stderr.len()
    );

    SshExecResult::success(stdout, stderr, exit_code)
}

/// Execute command interactively (for SCP/file transfer)
pub async fn ssh_exec_interactive(command: &str) -> Result<(), CloudToolsError> {
    let output = Command::new("bash")
        .args(["-c", command])
        .output()
        .await
        .map_err(|e| CloudToolsError::SshExecution(format!("Failed to execute: {}", e)))?;

    if !output.status.success() {
        let stderr = String::from_utf8_lossy(&output.stderr);
        return Err(CloudToolsError::SshExecution(stderr.to_string()));
    }

    Ok(())
}

/// Get Cloud Shell instance IP
pub async fn get_cloud_shell_ip() -> Result<String, CloudToolsError> {
    let output = Command::new("gcloud")
        .args([
            "cloud-shell",
            "get.SerialPortOutput",
            "--port=1",
            "--zone=us-central1-a",
        ])
        .output()
        .await
        .map_err(|e| CloudToolsError::SshConnection(format!("Failed to get IP: {}", e)))?;

    // Parse IP from serial output
    let output_str = String::from_utf8_lossy(&output.stdout);
    
    // Look for the external IP in the output
    if let Some(ip_line) = output_str.lines().find(|l| l.contains("external IP")) {
        if let Some(ip) = ip_line.split(':').nth(1) {
            return Ok(ip.trim().to_string());
        }
    }

    Err(CloudToolsError::CloudShellUnavailable(
        "Could not determine Cloud Shell IP".to_string(),
    ))
}

/// Check if Cloud Shell is running
pub async fn is_cloud_shell_running() -> bool {
    let output = Command::new("gcloud")
        .args(["cloud-shell", "describe", "--format=value(state)"])
        .output()
        .await;

    match output {
        Ok(out) if out.status.success() => {
            let state = String::from_utf8_lossy(&out.stdout).trim().to_string();
            state == "RUNNING"
        }
        _ => false,
    }
}

/// Get the Cloud Shell scp command prefix
pub fn get_scp_prefix() -> String {
    "gcloud cloud-shell scp".to_string()
}

/// Get the Cloud Shell ssh command prefix
pub fn get_ssh_prefix() -> String {
    "gcloud cloud-shell ssh --quiet --".to_string()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_cloud_shell_status() {
        let running = is_cloud_shell_running().await;
        tracing::info!("Cloud Shell running: {}", running);
        // This test just checks if we can query the status
        // Actual Cloud Shell may or may not be running
    }
}
