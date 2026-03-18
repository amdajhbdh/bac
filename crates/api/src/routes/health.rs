use axum::{
    routing::get,
    Router,
};

pub fn routes() -> Router {
    Router::new()
        .route("/health", get(health))
        .route("/ready", get(ready))
}

async fn health() -> &'static str {
    "OK"
}

async fn ready() -> &'static str {
    "Ready"
}
