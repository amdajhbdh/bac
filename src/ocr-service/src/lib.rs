//! BAC OCR Service - High-performance OCR with Rust
//!
//! This crate provides:
//! - Multiple OCR engines (Tesseract, Surya)
//! - Multi-layer fallback pipeline
//! - Batch processing support
//! - HTTP server

pub mod config;
pub mod metrics;
pub mod ocr;
pub mod parallel;
pub mod pipeline;
pub mod server;
pub mod service;

// Re-export pipeline types
pub use pipeline::{OCRPipeline, PipelineConfig, PipelineResult, PipelineError};

// Re-export OCR types
pub use ocr::{OCREngine, OCRResult, ProcessOptions, Config as OCRConfig};

// Re-export service types
pub use service::{OCRServiceImpl, OCRServiceMetrics};

// Re-export config types
pub use config::Config;

// Re-export server types
pub use server::{run_server, create_router, AppState, OCRResponse, OCRData, PDFResponse, PDFData};
