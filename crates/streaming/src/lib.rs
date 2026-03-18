//! BAC Streaming Service
//!
//! Streaming service for BAC Unified platform

use axum::Router;
use std::net::SocketAddr;
use tower_http::cors::CorsLayer;
use tracing_subscriber;

pub async fn run() {
    tracing_subscriber::fmt::init();

    let cors = CorsLayer::permissive();

    let app = Router::new()
        .layer(cors);

    let addr = SocketAddr::from(([0, 0, 0, 0], 8083));
    tracing::info!("Starting BAC Streaming server on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
