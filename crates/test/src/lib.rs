//! BAC Test Infrastructure
//!
//! Provides test utilities for integration testing.
//!
//! ## Usage
//!
//! ```rust
//! use bac_test::TestConfig;
//! ```

pub struct TestConfig {
    pub database_url: String,
    pub redis_url: String,
}

impl TestConfig {
    pub fn new() -> Self {
        Self {
            database_url: std::env::var("DATABASE_URL")
                .unwrap_or_else(|_| "postgres://bac:bac_test@localhost:5432/bac_test".to_string()),
            redis_url: std::env::var("REDIS_URL")
                .unwrap_or_else(|_| "redis://localhost:6379".to_string()),
        }
    }
}

impl Default for TestConfig {
    fn default() -> Self {
        Self::new()
    }
}
