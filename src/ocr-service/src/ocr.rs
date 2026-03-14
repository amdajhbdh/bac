//! OCR Engine Module
//!
//! Provides OCR functionality using leptess (Tesseract bindings)

use serde::{Deserialize, Serialize};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum OCRError {
    #[error("Image too large: {0} bytes (max 50MB)")]
    ImageTooLarge(usize),

    #[error("Unsupported format: {0}")]
    UnsupportedFormat(String),

    #[error("OCR processing failed: {0}")]
    ProcessingFailed(String),

    #[error("Tesseract error: {0}")]
    TesseractError(String),

    #[error("PDF processing failed: {0}")]
    PDFError(String),
}

pub type OCRErrorResult<T> = Result<T, OCRError>;

/// Supported OCR engines
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum OCREngine {
    /// Tesseract with LSTM neural network
    Lstm,
    /// Legacy Tesseract engine
    Legacy,
}

impl Default for OCREngine {
    fn default() -> Self {
        Self::Lstm
    }
}

impl std::fmt::Display for OCREngine {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            OCREngine::Lstm => write!(f, "tesseract-lstm"),
            OCREngine::Legacy => write!(f, "tesseract-legacy"),
        }
    }
}

/// Result from OCR processing
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OCRResult {
    /// Extracted text
    pub text: String,
    /// Confidence score (0.0 - 1.0)
    pub confidence: f64,
    /// OCR engine used
    pub source: String,
    /// Processing time in milliseconds
    pub processing_time_ms: i64,
}

impl OCRResult {
    pub fn new(text: String, confidence: f64, source: OCREngine, processing_time_ms: i64) -> Self {
        Self {
            text,
            confidence,
            source: source.to_string(),
            processing_time_ms,
        }
    }

    pub fn with_source(
        text: String,
        confidence: f64,
        source: &str,
        processing_time_ms: i64,
    ) -> Self {
        Self {
            text,
            confidence,
            source: source.to_string(),
            processing_time_ms,
        }
    }
}

/// Options for OCR processing
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProcessOptions {
    /// Language codes (e.g., "fra", "eng", "fra+eng")
    pub language: String,
    /// DPI for image processing
    pub dpi: u32,
    /// Confidence threshold for early exit
    pub confidence_threshold: f64,
    /// Enable automatic preprocessing
    pub auto_preprocess: bool,
}

impl Default for ProcessOptions {
    fn default() -> Self {
        Self {
            language: "fra+eng".to_string(),
            dpi: 300,
            confidence_threshold: 0.9,
            auto_preprocess: true,
        }
    }
}

/// Configuration for OCR service
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Config {
    /// Number of worker threads
    pub workers: usize,
    /// Maximum queue size
    pub max_queue_size: usize,
    /// Batch timeout in milliseconds
    pub batch_timeout_ms: u64,
    /// Maximum batch size
    pub batch_size: usize,
    /// Early exit confidence threshold
    pub early_exit_confidence: f64,
    /// Supported languages
    pub languages: Vec<String>,
    /// Data path for Tesseract
    pub tessdata_path: Option<String>,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            workers: 8,
            max_queue_size: 1000,
            batch_timeout_ms: 50,
            batch_size: 10,
            early_exit_confidence: 0.9,
            languages: vec!["fra".to_string(), "eng".to_string(), "ara".to_string()],
            tessdata_path: None,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_ocr_result_creation() {
        let result = OCRResult::new("Bonjour le monde".to_string(), 0.95, OCREngine::Lstm, 500);

        assert_eq!(result.text, "Bonjour le monde");
        assert_eq!(result.confidence, 0.95);
        assert!(result.source.contains("tesseract"));
        assert_eq!(result.processing_time_ms, 500);
    }

    #[test]
    fn test_process_options_defaults() {
        let options = ProcessOptions::default();

        assert_eq!(options.language, "fra+eng");
        assert_eq!(options.dpi, 300);
        assert_eq!(options.confidence_threshold, 0.9);
    }

    #[test]
    fn test_config_defaults() {
        let config = Config::default();

        assert_eq!(config.workers, 8);
        assert_eq!(config.max_queue_size, 1000);
        assert!(config.languages.contains(&"fra".to_string()));
    }
}
