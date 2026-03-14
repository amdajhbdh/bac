//! BAC OCR Service - Main Entry Point

use std::net::SocketAddr;

use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

use bac_ocr_service::pipeline::{OCRPipeline, PipelineConfig};
use bac_ocr_service::server::run_server;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize tracing
    tracing_subscriber::registry()
        .with(tracing_subscriber::fmt::layer())
        .init();
    
    tracing::info!("Starting BAC OCR Service");
    
    // Create OCR pipeline with default config
    let config = PipelineConfig::default();
    let pipeline = OCRPipeline::new(config);
    
    // Get address from environment or use default
    let addr: SocketAddr = std::env::var("OCR_HOST")
        .unwrap_or_else(|_| "127.0.0.1:3000".to_string())
        .parse()?;
    
    tracing::info!("OCR Service initialized");
    tracing::info!("Server listening on {}", addr);
    
    // Run HTTP server
    run_server(addr, pipeline).await;
    
    Ok(())
}
