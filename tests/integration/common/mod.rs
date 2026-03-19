//! Common test utilities for integration tests
//!
//! Shared helpers for HTTP testing, service setup, and assertion utilities.

use std::time::Duration;

/// Default timeout for HTTP requests
pub const DEFAULT_TIMEOUT: Duration = Duration::from_secs(30);

/// Short timeout for health checks
pub const HEALTH_TIMEOUT: Duration = Duration::from_secs(5);

/// Service base URLs
pub mod urls {
    /// Gemini Tools service URL
    pub const GEMINI: &str = "http://localhost:3001";
    
    /// Vector Tools service URL
    pub const VECTOR: &str = "http://localhost:3002";
    
    /// Vault Tools service URL
    pub const VAULT: &str = "http://localhost:3003";
    
    /// Cloud Tools service URL
    pub const CLOUD: &str = "http://localhost:3004";
    
    /// Graph Tools service URL
    pub const GRAPH: &str = "http://localhost:3005";
}

/// Create a new HTTP client with default configuration
pub fn create_client() -> reqwest::Client {
    reqwest::Client::builder()
        .timeout(DEFAULT_TIMEOUT)
        .build()
        .expect("Failed to create HTTP client")
}

/// Create a client with a custom timeout
pub fn create_client_with_timeout(timeout: Duration) -> reqwest::Client {
    reqwest::Client::builder()
        .timeout(timeout)
        .build()
        .expect("Failed to create HTTP client")
}

/// Wait for a service to become available
/// 
/// # Arguments
/// * `url` - The health endpoint URL to check
/// * `max_attempts` - Maximum number of attempts
/// * `interval` - Time between attempts
/// 
/// # Returns
/// `true` if the service became available, `false` otherwise
pub async fn wait_for_service(url: &str, max_attempts: u32, interval: Duration) -> bool {
    let client = create_client_with_timeout(HEALTH_TIMEOUT);
    
    for attempt in 1..=max_attempts {
        match client.get(url).send().await {
            Ok(response) if response.status().is_success() => {
                return true;
            }
            _ => {
                if attempt < max_attempts {
                    tokio::time::sleep(interval).await;
                }
            }
        }
    }
    
    false
}

/// Check if a service is running by attempting a health check
pub async fn is_service_running(url: &str) -> bool {
    wait_for_service(url, 1, Duration::ZERO).await
}

/// Test result type for integration tests
pub type TestResult<T = ()> = Result<T, TestError>;

/// Errors that can occur during integration testing
#[derive(Debug)]
pub enum TestError {
    /// Service is not available
    ServiceUnavailable(String),
    
    /// HTTP request failed
    RequestFailed(reqwest::Error),
    
    /// Unexpected response status
    UnexpectedStatus(String, reqwest::StatusCode),
    
    /// JSON serialization/deserialization failed
    SerializationError(String),
    
    /// Assertion failed
    AssertionFailed(String),
    
    /// Test skipped
    Skipped(String),
}

impl std::fmt::Display for TestError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            TestError::ServiceUnavailable(url) => {
                write!(f, "Service unavailable: {}", url)
            }
            TestError::RequestFailed(e) => {
                write!(f, "Request failed: {}", e)
            }
            TestError::UnexpectedStatus(msg, status) => {
                write!(f, "{} - Status: {}", msg, status)
            }
            TestError::SerializationError(msg) => {
                write!(f, "Serialization error: {}", msg)
            }
            TestError::AssertionFailed(msg) => {
                write!(f, "Assertion failed: {}", msg)
            }
            TestError::Skipped(msg) => {
                write!(f, "Test skipped: {}", msg)
            }
        }
    }
}

impl From<reqwest::Error> for TestError {
    fn from(err: reqwest::Error) -> Self {
        TestError::RequestFailed(err)
    }
}

impl From<serde_json::Error> for TestError {
    fn from(err: serde_json::Error) -> Self {
        TestError::SerializationError(err.to_string())
    }
}

/// Assert that a response is successful (2xx status)
pub fn assert_success_response(response: &reqwest::Response) -> TestResult<()> {
    if !response.status().is_success() {
        return Err(TestError::UnexpectedStatus(
            format!("Expected success status, got {:?}", response.status()),
            response.status(),
        ));
    }
    Ok(())
}

/// Assert that a response has a specific status code
pub fn assert_status(response: &reqwest::Response, expected: reqwest::StatusCode) -> TestResult<()> {
    let actual = response.status();
    if actual != expected {
        return Err(TestError::UnexpectedStatus(
            format!("Expected status {}, got {}", expected, actual),
            actual,
        ));
    }
    Ok(())
}

/// Helper macro to skip a test if the service is not available
#[macro_export]
macro_rules! skip_if_service_unavailable {
    ($url:expr) => {
        if !$crate::common::is_service_running($url).await {
            eprintln!(
                "Skipping test: service at {} is not available",
                $url
            );
            return;
        }
    };
}

/// Helper macro to wrap test assertions
#[macro_export]
macro_rules! assert_test {
    ($condition:expr, $msg:expr) => {
        if !$condition {
            return Err($crate::common::TestError::AssertionFailed($msg.to_string()));
        }
    };
}
