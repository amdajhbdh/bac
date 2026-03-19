//! Gemini Tools Library
//!
//! Provides Gemini API wrapper functionality with HTTP service and CLI.

pub mod analyze;
pub mod client;
pub mod correct;
pub mod embed;
pub mod extract;
pub mod generate;
pub mod service;

pub use client::GeminiClient;
pub use service::{run_server, AppState, create_router};

/// Initialize the gemini-tools library.
pub fn init() {
    tracing::info!("gemini-tools initialized");
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_init() {
        init();
        // Just verify it doesn't panic
    }
}
