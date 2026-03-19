//! BAC API Service
//!
//! REST API for BAC Unified platform
//! 
//! # Port Configuration
//! - API service: :8080
//! - Tool service proxies:
//!   - gemini-tools: :3001
//!   - vector-tools: :3002
//!   - vault-tools: :3003
//!   - cloud-tools: :3004
//!   - graph-tools: :3005

pub mod client;
pub mod models;
pub mod proxy;
pub mod routes;
pub mod state;

use axum::Router;
use std::net::SocketAddr;
use tower_http::cors::CorsLayer;
use tracing_subscriber;

pub async fn run() {
    run_with_client(client::ToolClient::from_env()).await;
}

/// Run the API server with a custom tool client
pub async fn run_with_client(tool_client: client::ToolClient) {
    tracing_subscriber::fmt::init();

    let cors = CorsLayer::permissive();

    let proxy_state = proxy::ProxyState::new(tool_client);

    let app = Router::new()
        .merge(routes::health::routes())
        .merge(routes::problems::routes())
        .merge(routes::solutions::routes())
        .merge(routes::users::routes())
        .merge(routes::auth::routes())
        .merge(proxy::routes(proxy_state))
        .layer(cors);

    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    tracing::info!("Starting BAC API server on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
