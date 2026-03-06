//! BAC OCR Service - High-performance OCR with Rust
//!
//! This crate provides:
//! - Multiple OCR engines (Tesseract, Surya)
//! - Multi-layer fallback pipeline
//! - Batch processing support

pub mod pipeline;

// Re-export pipeline types
pub use pipeline::{OCRPipeline, PipelineConfig, PipelineResult, PipelineError, PipelineError as Error};
