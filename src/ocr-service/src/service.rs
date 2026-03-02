//! OCR Service Implementation
//!
//! Provides the main OCR service implementation

use std::sync::Arc;
use tokio::sync::RwLock;

use crate::ocr::{OCRResult, ProcessOptions, Config as OCRConfig};
use crate::parallel::WorkerPool;
use crate::metrics::OCREmptyMetrics;

/// OCR Service implementation
pub struct OCRServiceImpl {
    worker_pool: WorkerPool,
    metrics: Arc<OCREmptyMetrics>,
    config: OCRConfig,
}

impl OCRServiceImpl {
    pub fn new() -> Self {
        let config = OCRConfig::default();
        let worker_pool = WorkerPool::new(config.workers);
        
        Self {
            worker_pool,
            metrics: Arc::new(OCREmptyMetrics::new()),
            config,
        }
    }
    
    pub fn with_config(config: OCRConfig) -> Self {
        let worker_pool = WorkerPool::new(config.workers);
        
        Self {
            worker_pool,
            metrics: Arc::new(OCREmptyMetrics::new()),
            config,
        }
    }
    
    /// Process a single image
    pub async fn process_image(&self, image_data: Vec<u8>, language: String) -> Result<OCRResult, String> {
        // Validate image size
        const MAX_SIZE: usize = 50 * 1024 * 1024; // 50MB
        if image_data.len() > MAX_SIZE {
            return Err(format!("Image too large: {} bytes (max 50MB)", image_data.len()));
        }
        
        self.metrics.record_request();
        
        let options = ProcessOptions {
            language,
            ..Default::default()
        };
        
        // Process in worker pool
        let result = self.worker_pool.process_batch(vec![(image_data, options)]);
        
        match result.into_iter().next() {
            Some(Ok(ocr_result)) => {
                self.metrics.record_duration(ocr_result.processing_time_ms as u64);
                Ok(ocr_result)
            }
            Some(Err(e)) => {
                self.metrics.record_error();
                Err(e.to_string())
            }
            None => Err("No result".to_string()),
        }
    }
    
    /// Process a batch of images
    pub async fn process_batch(&self, images: Vec<(Vec<u8>, String)>) -> Result<Vec<OCRResult>, String> {
        self.metrics.record_batch(images.len());
        
        let items: Vec<_> = images
            .into_iter()
            .map(|(data, lang)| {
                let options = ProcessOptions {
                    language: lang,
                    ..Default::default()
                };
                (data, options)
            })
            .collect();
        
        let results = self.worker_pool.process_batch(items);
        
        let mut ocr_results = Vec::new();
        for result in results {
            match result {
                Ok(ocr) => {
                    self.metrics.record_duration(ocr.processing_time_ms as u64);
                    ocr_results.push(ocr);
                }
                Err(e) => {
                    self.metrics.record_error();
                    return Err(e.to_string());
                }
            }
        }
        
        Ok(ocr_results)
    }
    
    /// Get metrics
    pub fn get_metrics(&self) -> OCRServiceMetrics {
        OCRServiceMetrics {
            requests_total: self.metrics.get_requests_total(),
            errors_total: self.metrics.get_errors_total(),
            avg_duration_ms: self.metrics.get_avg_duration_ms(),
            avg_batch_size: self.metrics.get_avg_batch_size(),
        }
    }
}

impl Default for OCRServiceImpl {
    fn default() -> Self {
        Self::new()
    }
}

#[derive(Debug, Clone)]
pub struct OCRServiceMetrics {
    pub requests_total: u64,
    pub errors_total: u64,
    pub avg_duration_ms: f64,
    pub avg_batch_size: f64,
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[tokio::test]
    async fn test_process_image() {
        let service = OCRServiceImpl::new();
        
        let result = service.process_image(
            vec![0x89, 0x50, 0x4E, 0x47], // PNG magic
            "fra".to_string(),
        ).await;
        
        // Should succeed (returns mock result)
        assert!(result.is_ok());
    }
    
    #[tokio::test]
    async fn test_process_batch() {
        let service = OCRServiceImpl::new();
        
        let images = vec![
            (vec![0x01], "fra".to_string()),
            (vec![0x02], "eng".to_string()),
        ];
        
        let results = service.process_batch(images).await;
        
        assert!(results.is_ok());
        assert_eq!(results.unwrap().len(), 2);
    }
    
    #[test]
    fn test_image_too_large() {
        let service = OCRServiceImpl::new();
        
        // Create image larger than 50MB
        let large_image = vec![0u8; 51 * 1024 * 1024];
        
        let rt = tokio::runtime::Runtime::new().unwrap();
        let result = rt.block_on(service.process_image(large_image, "fra".to_string()));
        
        assert!(result.is_err());
        assert!(result.unwrap_err().contains("too large"));
    }
}
