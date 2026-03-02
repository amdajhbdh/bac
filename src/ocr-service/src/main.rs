//! BAC OCR Service - Main Entry Point

use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize tracing
    tracing_subscriber::registry()
        .with(tracing_subscriber::fmt::layer())
        .init();
    
    tracing::info!("Starting BAC OCR Service");
    
    // Create OCR service
    let service = bac_ocr_service::service::OCRServiceImpl::new();
    
    // Get metrics
    let metrics = service.get_metrics();
    
    tracing::info!("OCR Service initialized");
    tracing::info!("Metrics: {:?}", metrics);
    
    Ok(())
}
