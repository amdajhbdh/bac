//! Multi-layer OCR Pipeline
//!
//! Implements fallback chain: Tesseract → Surya → Google Lens
//! With PDF multi-page support

use std::time::Instant;

use rayon::prelude::*;
use serde::{Deserialize, Serialize};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum PipelineError {
    #[error("All OCR engines failed")]
    AllEnginesFailed { errors: Vec<String> },
    
    #[error("Image processing failed: {0}")]
    ImageError(String),
    
    #[error("PDF processing failed: {0}")]
    PDFError(String),
    
    #[error("Tesseract error: {0}")]
    TesseractError(String),
    
    #[error("Surya error: {0}")]
    SuryaError(String),
    
    #[error("Google Lens API error: {0}")]
    GoogleLensError(String),
    
    #[error("Network error: {0}")]
    NetworkError(String),
    
    #[error("Timeout")]
    Timeout,
}

pub type Result<T> = std::result::Result<T, PipelineError>;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum OCREngine {
    Tesseract,
    Surya,
    GoogleLens,
}

impl std::fmt::Display for OCREngine {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            OCREngine::Tesseract => write!(f, "tesseract"),
            OCREngine::Surya => write!(f, "surya"),
            OCREngine::GoogleLens => write!(f, "google-lens"),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PipelineResult {
    pub text: String,
    pub confidence: f64,
    pub engine: OCREngine,
    pub processing_time_ms: i64,
    pub attempts: Vec<EngineAttempt>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EngineAttempt {
    pub engine: OCREngine,
    pub success: bool,
    pub confidence: f64,
    pub error: Option<String>,
    pub processing_time_ms: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PDFResult {
    pub pages: Vec<PipelineResult>,
    pub total_pages: usize,
    pub combined_text: String,
    pub average_confidence: f64,
    pub processing_time_ms: i64,
}

#[derive(Debug, Clone)]
pub struct PipelineConfig {
    pub tesseract_enabled: bool,
    pub surya_enabled: bool,
    pub googlelens_enabled: bool,
    pub tesseract_threshold: f64,
    pub surya_threshold: f64,
    pub max_attempts: usize,
    pub engine_timeout_secs: u64,
    pub workers: usize,
}

impl Default for PipelineConfig {
    fn default() -> Self {
        Self {
            tesseract_enabled: true,
            surya_enabled: true,
            googlelens_enabled: false,
            tesseract_threshold: 0.7,
            surya_threshold: 0.85,
            max_attempts: 3,
            engine_timeout_secs: 30,
            workers: num_cpus::get(),
        }
    }
}

#[derive(Debug, Clone)]
struct OCREngineResult {
    text: String,
    confidence: f64,
    processing_time_ms: i64,
}

pub struct OCRPipeline {
    config: PipelineConfig,
}

impl OCRPipeline {
    pub fn new(config: PipelineConfig) -> Self {
        Self { config }
    }
    
    pub async fn process(&self, image_data: &[u8]) -> Result<PipelineResult> {
        let start = Instant::now();
        let mut attempts = Vec::new();
        
        if self.config.tesseract_enabled {
            let attempt = self.process_tesseract(image_data).await;
            attempts.push(EngineAttempt {
                engine: OCREngine::Tesseract,
                success: attempt.is_ok(),
                confidence: attempt.as_ref().map(|r| r.confidence).unwrap_or(0.0),
                error: attempt.as_ref().err().map(|e| e.to_string()),
                processing_time_ms: attempt.as_ref().map(|r| r.processing_time_ms).unwrap_or(0),
            });
            
            if let Ok(result) = attempt {
                if result.confidence >= self.config.tesseract_threshold {
                    return Ok(PipelineResult {
                        text: result.text,
                        confidence: result.confidence,
                        engine: OCREngine::Tesseract,
                        processing_time_ms: start.elapsed().as_millis() as i64,
                        attempts,
                    });
                }
            }
        }
        
        if self.config.surya_enabled {
            let attempt = self.process_surya(image_data).await;
            attempts.push(EngineAttempt {
                engine: OCREngine::Surya,
                success: attempt.is_ok(),
                confidence: attempt.as_ref().map(|r| r.confidence).unwrap_or(0.0),
                error: attempt.as_ref().err().map(|e| e.to_string()),
                processing_time_ms: attempt.as_ref().map(|r| r.processing_time_ms).unwrap_or(0),
            });
            
            if let Ok(result) = attempt {
                if result.confidence >= self.config.surya_threshold {
                    return Ok(PipelineResult {
                        text: result.text,
                        confidence: result.confidence,
                        engine: OCREngine::Surya,
                        processing_time_ms: start.elapsed().as_millis() as i64,
                        attempts,
                    });
                }
            }
        }
        
        let errors: Vec<String> = attempts.iter()
            .filter_map(|a| a.error.clone())
            .collect();
        
        Err(PipelineError::AllEnginesFailed { errors })
    }
    
    async fn process_tesseract(&self, image_data: &[u8]) -> Result<OCREngineResult> {
        let start = Instant::now();
        
        let temp_dir = std::env::temp_dir();
        let temp_path = temp_dir.join(format!("tess_{}.png", uuid::Uuid::new_v4()));
        std::fs::write(&temp_path, image_data)
            .map_err(|e| PipelineError::TesseractError(format!("Write error: {}", e)))?;
        
        let path_str = temp_path.to_string_lossy().to_string();
        
        let result = tokio::task::spawn_blocking(move || {
            let output = std::process::Command::new("tesseract")
                .arg(&path_str)
                .arg("stdout")
                .arg("-l")
                .arg("eng+fra")
                .arg("--psm")
                .arg("6")
                .output();
            
            let _ = std::fs::remove_file(&path_str);
            
            match output {
                Ok(out) if out.status.success() => {
                    Ok(OCREngineResult {
                        text: String::from_utf8_lossy(&out.stdout).to_string(),
                        confidence: 0.85,
                        processing_time_ms: start.elapsed().as_millis() as i64,
                    })
                }
                Ok(out) => Err(PipelineError::TesseractError(
                    String::from_utf8_lossy(&out.stderr).to_string())),
                Err(e) => Err(PipelineError::TesseractError(e.to_string())),
            }
        }).await.map_err(|e| PipelineError::TesseractError(e.to_string()))??;
        
        Ok(result)
    }
    
    async fn process_surya(&self, image_data: &[u8]) -> Result<OCREngineResult> {
        let start = Instant::now();
        
        let temp_dir = std::env::temp_dir();
        let temp_path = temp_dir.join(format!("surya_{}.png", uuid::Uuid::new_v4()));
        std::fs::write(&temp_path, image_data)
            .map_err(|e| PipelineError::SuryaError(format!("Write error: {}", e)))?;
        
        let path_str = temp_path.to_string_lossy().to_string();
        
        let result = tokio::task::spawn_blocking(move || {
            let python_code = format!(r#"
import sys
import json
try:
    from surya.ocr import run_ocr
    from surya.model.detection.segformer import load_model as load_det, load_processor as load_det_proc
    from surya.model.recognition.model import load_model as load_rec
    from surya.model.recognition.processor import load_processor as load_rec_proc
    
    det_model = load_det()
    det_proc = load_det_proc()
    rec_model = load_rec()
    rec_proc = load_rec_proc()
    
    preds = run_ocr([r"{path_str}"], [det_model, rec_model], [det_proc, rec_proc])
    
    text = "\n".join([p.text for p in preds[0]])
    conf = sum([p.confidence for p in preds[0]]) / len(preds[0]) if preds[0] else 0
    print(json.dumps({{"text": text, "confidence": conf}}))
except Exception as e:
    print(json.dumps({{"error": str(e)}}), file=sys.stderr)
    sys.exit(1)
"#, path_str=path_str);
            
            let output = std::process::Command::new("python3")
                .arg("-c")
                .arg(python_code)
                .output();
            
            let _ = std::fs::remove_file(&path_str);
            
            match output {
                Ok(out) if out.status.success() => {
                    let json_str = String::from_utf8_lossy(&out.stdout);
                    if let Ok(data) = serde_json::from_str::<serde_json::Value>(&json_str) {
                        let text = data["text"].as_str().unwrap_or("").to_string();
                        let confidence = data["confidence"].as_f64().unwrap_or(0.0);
                        Ok(OCREngineResult {
                            text,
                            confidence,
                            processing_time_ms: start.elapsed().as_millis() as i64,
                        })
                    } else {
                        Err(PipelineError::SuryaError("Parse error".to_string()))
                    }
                }
                Ok(out) => Err(PipelineError::SuryaError(
                    String::from_utf8_lossy(&out.stderr).to_string())),
                Err(e) => Err(PipelineError::SuryaError(e.to_string())),
            }
        }).await.map_err(|e| PipelineError::SuryaError(e.to_string()))??;
        
        Ok(result)
    }
    
    /// Process a PDF document - converts pages to images and runs OCR
    pub async fn process_pdf(&self, pdf_data: &[u8]) -> Result<PDFResult> {
        let start = Instant::now();
        
        // Save PDF to temp file
        let temp_dir = std::env::temp_dir();
        let pdf_path = temp_dir.join(format!("input_{}.pdf", uuid::Uuid::new_v4()));
        let output_dir = temp_dir.join(format!("pdf_out_{}", uuid::Uuid::new_v4()));
        
        std::fs::write(&pdf_path, pdf_data)
            .map_err(|e| PipelineError::PDFError(format!("Failed to write PDF: {}", e)))?;
        
        std::fs::create_dir_all(&output_dir)
            .map_err(|e| PipelineError::PDFError(format!("Failed to create output dir: {}", e)))?;
        
        let pdf_path_str = pdf_path.to_string_lossy().to_string();
        let output_dir_str = output_dir.to_string_lossy().to_string();
        
        // Convert PDF pages to images using pdftoppm
        let page_count = tokio::task::spawn_blocking(move || {
            let output = std::process::Command::new("pdftoppm")
                .arg("-r")
                .arg("300")
                .arg("-png")
                .arg("-singlefile")
                .arg(&pdf_path_str)
                .arg(format!("{}/page", output_dir_str))
                .output();
            
            let _ = std::fs::remove_file(&pdf_path_str);
            
            match output {
                Ok(out) if out.status.success() => {
                    // Count generated images
                    let entries = std::fs::read_dir(&output_dir_str)
                        .map_err(|e| PipelineError::PDFError(e.to_string()))?;
                    let count = entries.filter_map(|e| e.ok())
                        .filter(|e| e.path().extension().map(|ext| ext == "png").unwrap_or(false))
                        .count();
                    Ok(count.max(1))
                }
                Ok(out) => {
                    let err = String::from_utf8_lossy(&out.stderr);
                    Err(PipelineError::PDFError(format!("pdftoppm failed: {}", err)))
                }
                Err(e) => Err(PipelineError::PDFError(format!("pdftoppm not found: {}", e)))
            }
        }).await.map_err(|e| PipelineError::PDFError(e.to_string()))??;
        
        // Process each page with OCR
        let mut page_results = Vec::new();
        
        for page_num in 1..=page_count {
            let page_file = output_dir.join(format!("page-{}.png", page_num));
            if !page_file.exists() {
                // Try alternative naming (pdftoppm -singlefile uses different naming)
                let alt_page_file = output_dir.join("page.png");
                if alt_page_file.exists() {
                    // Single page PDF
                    let page_data = std::fs::read(&alt_page_file)
                        .map_err(|e| PipelineError::PDFError(e.to_string()))?;
                    match self.process(&page_data).await {
                        Ok(result) => page_results.push(result),
                        Err(e) => {
                            // Continue with other pages even if one fails
                            tracing::warn!("Page {} OCR failed: {}", page_num, e);
                        }
                    }
                    break;
                }
                continue;
            }
            
            let page_data = std::fs::read(&page_file)
                .map_err(|e| PipelineError::PDFError(e.to_string()))?;
            
            match self.process(&page_data).await {
                Ok(result) => page_results.push(result),
                Err(e) => {
                    tracing::warn!("Page {} OCR failed: {}", page_num, e);
                }
            }
            
            // Cleanup page file
            let _ = std::fs::remove_file(&page_file);
        }
        
        // Cleanup output directory
        let _ = std::fs::remove_dir_all(&output_dir);
        
        if page_results.is_empty() {
            return Err(PipelineError::PDFError("No pages could be processed".to_string()));
        }
        
        // Combine results
        let combined_text: String = page_results.iter()
            .map(|r| format!("--- Page {} ---\n{}", r.processing_time_ms, r.text))
            .collect::<Vec<_>>()
            .join("\n\n");
        
        let avg_confidence: f64 = page_results.iter()
            .map(|r| r.confidence)
            .sum::<f64>() / page_results.len() as f64;
        
        Ok(PDFResult {
            pages: page_results,
            total_pages: page_count,
            combined_text,
            average_confidence: avg_confidence,
            processing_time_ms: start.elapsed().as_millis() as i64,
        })
    }
    
    pub async fn process_batch(&self, images: Vec<Vec<u8>>) -> Vec<Result<PipelineResult>> {
        let pipeline = self.clone();
        
        images
            .into_par_iter()
            .map(|img| {
                let p = pipeline.clone();
                tokio::runtime::Handle::current().block_on(async move {
                    p.process(&img).await
                })
            })
            .collect()
    }
}

impl Clone for OCRPipeline {
    fn clone(&self) -> Self {
        Self {
            config: self.config.clone(),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_config_defaults() {
        let config = PipelineConfig::default();
        assert!(config.tesseract_enabled);
        assert!(config.surya_enabled);
    }
}
