//! Cloud Tools Library
//!
//! Provides Cloud Shell integration via SSH, file transfer, OCR, and GCS sync.
//!
//! # Features
//! - Execute commands on Cloud Shell via SSH
//! - Upload files to Cloud Shell
//! - Download files from Cloud Shell
//! - Run OCR on Cloud Shell
//! - Sync files with GCS buckets

pub mod models;
pub mod service;
pub mod ssh;
pub mod upload;
pub mod download;
pub mod ocr;
pub mod gcs;

pub use models::*;
pub use service::run;

/// Initialize the cloud-tools library.
pub fn init() {
    tracing::info!("cloud-tools initialized");
}
