//! BAC API Service
//!
//! REST API for BAC Unified platform

pub mod models;
pub mod routes;
pub mod state;

use axum::Router;
use std::net::SocketAddr;
use tower_http::cors::CorsLayer;
use tracing_subscriber;

pub async fn run() {
    tracing_subscriber::fmt::init();

    let cors = CorsLayer::permissive();

    let app = Router::new()
        .merge(routes::health::routes())
        .merge(routes::problems::routes())
        .merge(routes::solutions::routes())
        .merge(routes::users::routes())
        .merge(routes::auth::routes())
        .layer(cors);

    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    tracing::info!("Starting BAC API server on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
