//! BAC OCR Service - High-performance OCR with Rust
//!
//! This crate provides:
//! - Parallel OCR processing with Rayon
//! - Multiple OCR engines (Tesseract LSTM, legacy)
//! - Batch processing support
//! - gRPC interface

pub mod ocr;
pub mod parallel;
pub mod metrics;
pub mod config;

#[cfg(feature = "grpc")]
pub mod service;

pub use ocr::{OCREngine, OCRResult, ProcessOptions};
pub use parallel::WorkerPool;
pub use metrics::OCREmptyMetrics;
pub use config::Config;
