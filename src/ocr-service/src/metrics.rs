//! Metrics Module
//!
//! Provides Prometheus metrics for the OCR service

use metrics::{counter, gauge, histogram};
use std::sync::atomic::{AtomicU64, Ordering};

/// OCR service metrics
pub struct OCREmptyMetrics {
    requests_total: AtomicU64,
    errors_total: AtomicU64,
    processing_duration_sum: AtomicU64,
    batch_sizes_sum: AtomicU64,
    batch_count: AtomicU64,
}

impl Default for OCREmptyMetrics {
    fn default() -> Self {
        Self::new()
    }
}

impl OCREmptyMetrics {
    pub fn new() -> Self {
        Self {
            requests_total: AtomicU64::new(0),
            errors_total: AtomicU64::new(0),
            processing_duration_sum: AtomicU64::new(0),
            batch_sizes_sum: AtomicU64::new(0),
            batch_count: AtomicU64::new(0),
        }
    }

    pub fn record_request(&self) {
        self.requests_total.fetch_add(1, Ordering::Relaxed);
        counter!("ocr_requests_total").increment(1);
    }

    pub fn record_error(&self) {
        self.errors_total.fetch_add(1, Ordering::Relaxed);
        counter!("ocr_errors_total").increment(1);
    }

    pub fn record_duration(&self, ms: u64) {
        self.processing_duration_sum
            .fetch_add(ms, Ordering::Relaxed);
        histogram!("ocr_processing_duration_ms").record(ms as f64);
    }

    pub fn record_batch(&self, size: usize) {
        self.batch_sizes_sum
            .fetch_add(size as u64, Ordering::Relaxed);
        self.batch_count.fetch_add(1, Ordering::Relaxed);
        gauge!("ocr_batch_size").set(size as f64);
    }

    pub fn get_requests_total(&self) -> u64 {
        self.requests_total.load(Ordering::Relaxed)
    }

    pub fn get_errors_total(&self) -> u64 {
        self.errors_total.load(Ordering::Relaxed)
    }

    pub fn get_avg_duration_ms(&self) -> f64 {
        let total = self.processing_duration_sum.load(Ordering::Relaxed);
        let count = self.requests_total.load(Ordering::Relaxed);
        if count > 0 {
            total as f64 / count as f64
        } else {
            0.0
        }
    }

    pub fn get_avg_batch_size(&self) -> f64 {
        let total = self.batch_sizes_sum.load(Ordering::Relaxed);
        let count = self.batch_count.load(Ordering::Relaxed);
        if count > 0 {
            total as f64 / count as f64
        } else {
            0.0
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_metrics_counters() {
        let metrics = OCREmptyMetrics::new();

        metrics.record_request();
        metrics.record_request();

        assert_eq!(metrics.get_requests_total(), 2);
    }

    #[test]
    fn test_metrics_errors() {
        let metrics = OCREmptyMetrics::new();

        metrics.record_error();

        assert_eq!(metrics.get_errors_total(), 1);
    }

    #[test]
    fn test_avg_duration() {
        let metrics = OCREmptyMetrics::new();

        metrics.record_duration(100);
        metrics.record_duration(200);

        assert_eq!(metrics.get_avg_duration_ms(), 150.0);
    }
}
