//! Parallel Processing Module
//!
//! Provides worker pool for parallel OCR processing using Rayon

use rayon::prelude::*;
use std::sync::Arc;
use std::time::Instant;
use tokio::sync::mpsc;

use crate::ocr::{OCRError, OCRResult, ProcessOptions};

/// Worker pool for parallel OCR processing
pub struct WorkerPool {
    workers: usize,
    queue: Arc<flume::Sender<WorkItem>>,
}

struct WorkItem {
    data: Vec<u8>,
    options: ProcessOptions,
    result_tx: mpsc::Sender<Result<OCRResult, OCRError>>,
}

impl WorkerPool {
    /// Create a new worker pool
    pub fn new(workers: usize) -> Self {
        let (tx, rx) = flume::unbounded();

        // Spawn worker threads
        for _ in 0..workers {
            let rx = rx.clone();
            std::thread::spawn(move || {
                Self::worker_loop(rx);
            });
        }

        Self { workers, queue: tx }
    }

    fn worker_loop(rx: flume::Receiver<WorkItem>) {
        while let Ok(work) = rx.recv() {
            // Process work item
            let result = Self::process_single(&work.data, &work.options);
            let _ = work.result_tx.blocking_send(result);
        }
    }

    fn process_single(data: &[u8], options: &ProcessOptions) -> Result<OCRResult, OCRError> {
        let start = Instant::now();

        // Placeholder: actual OCR processing would happen here
        // For now, return a mock result
        let result = OCRResult::new(
            String::from("Extracted text"),
            0.85,
            crate::ocr::OCREngine::Lstm,
            start.elapsed().as_millis() as i64,
        );

        Ok(result)
    }

    /// Process multiple items in parallel
    pub fn process_batch(
        &self,
        items: Vec<(Vec<u8>, ProcessOptions)>,
    ) -> Vec<Result<OCRResult, OCRError>> {
        items
            .par_iter()
            .map(|(data, options)| Self::process_single(data, options))
            .collect()
    }

    /// Get number of workers
    pub fn workers(&self) -> usize {
        self.workers
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_worker_pool_creation() {
        let pool = WorkerPool::new(4);
        assert_eq!(pool.workers(), 4);
    }

    #[test]
    fn test_batch_processing() {
        let pool = WorkerPool::new(4);

        let items: Vec<_> = (0..10)
            .map(|i| (vec![i as u8], ProcessOptions::default()))
            .collect();

        let results = pool.process_batch(items);

        assert_eq!(results.len(), 10);
    }
}
