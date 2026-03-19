# Axum 0.8.x Full Example Application

## Project Structure

```
src/
├── main.rs           # Entry point
├── routes/
│   ├── mod.rs
│   ├── health.rs
│   └── users.rs
├── handlers/
│   ├── mod.rs
│   └── users.rs
├── models/
│   ├── mod.rs
│   └── user.rs
├── error.rs          # Error types
└── state.rs          # Application state
```

## Cargo.toml

```toml
[package]
name = "api-server"
version = "0.1.0"
edition = "2021"

[dependencies]
# Web framework
axum = { version = "0.7", features = ["macros"] }
tokio = { version = "1", features = ["full"] }

# Middleware
tower = { version = "0.4", features = ["timeout", "limit"] }
tower-http = { version = "0.5", features = [
    "cors",
    "trace",
    "compression-gzip",
    "request-id",
] }

# Serialization
serde = { version = "1", features = ["derive"] }
serde_json = "1"

# Error handling
thiserror = "1"
anyhow = "1"

# Observability
tracing = "0.1"
tracing-subscriber = { version = "0.3", features = ["env-filter", "json"] }
```

## Error Module (src/error.rs)

```rust
use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use serde::Serialize;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("Resource not found: {0}")]
    NotFound(String),
    
    #[error("Validation error: {0}")]
    ValidationError(String),
    
    #[error("Internal server error")]
    InternalError(#[from] anyhow::Error),
}

#[derive(Serialize)]
struct ErrorBody {
    error: String,
    message: String,
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, message) = match &self {
            AppError::NotFound(msg) => (StatusCode::NOT_FOUND, msg.clone()),
            AppError::ValidationError(msg) => (StatusCode::BAD_REQUEST, msg.clone()),
            AppError::InternalError(msg) => {
                tracing::error!("Internal error: {}", msg);
                (StatusCode::INTERNAL_SERVER_ERROR, "Internal server error".to_string())
            }
        };

        let body = Json(ErrorBody {
            error: self.to_string(),
            message,
        });

        (status, body).into_response()
    }
}

pub type Result<T> = std::result::Result<T, AppError>;
```

## Models (src/models/user.rs)

```rust
use serde::{Deserialize, Serialize};

#[derive(Deserialize)]
pub struct CreateUser {
    pub name: String,
    pub email: String,
}

#[derive(Serialize)]
pub struct User {
    pub id: u64,
    pub name: String,
    pub email: String,
}
```

## Handlers (src/handlers/users.rs)

```rust
use crate::error::{AppError, Result};
use crate::models::user::{CreateUser, User};
use axum::{
    extract::{Json, Path, State},
    http::StatusCode,
    response::IntoResponse,
    routing::{get, post},
    Router,
};
use std::sync::Arc;

pub fn router() -> Router<Arc<()>> {
    Router::new()
        .route("/users", get(list_users).post(create_user))
        .route("/users/{id}", get(get_user).put(update_user).delete(delete_user))
}

async fn list_users() -> Json<Vec<User>> {
    Json(vec![
        User { id: 1, name: "Alice".into(), email: "alice@example.com".into() }
    ])
}

async fn create_user(Json(payload): Json<CreateUser>) -> Result<(StatusCode, Json<User>)> {
    if payload.name.is_empty() {
        return Err(AppError::ValidationError("Name cannot be empty".into()));
    }
    
    let user = User {
        id: 1,
        name: payload.name,
        email: payload.email,
    };
    
    Ok((StatusCode::CREATED, Json(user)))
}

async fn get_user(Path(id): Path<u64>) -> Result<Json<User>> {
    if id == 0 {
        return Err(AppError::NotFound(format!("User {} not found", id)));
    }
    
    Ok(Json(User {
        id,
        name: "User".into(),
        email: "user@example.com".into(),
    }))
}

async fn update_user(Path(id): Path<u64>, Json(payload): Json<CreateUser>) -> Result<Json<User>> {
    Ok(Json(User {
        id,
        name: payload.name,
        email: payload.email,
    }))
}

async fn delete_user(Path(id): Path<u64>) -> StatusCode {
    StatusCode::NO_CONTENT
}
```

## Health Routes (src/routes/health.rs)

```rust
use axum::{
    routing::get,
    Router,
};
use serde::Serialize;

#[derive(Serialize)]
struct HealthResponse {
    status: String,
    version: String,
}

pub fn router() -> Router {
    Router::new()
        .route("/health", get(health))
        .route("/ready", get(ready))
}

async fn health() -> &'static str {
    "OK"
}

async fn ready() -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "ready".into(),
        version: env!("CARGO_PKG_VERSION").into(),
    })
}
```

## Main Entry Point (src/main.rs)

```rust
mod error;
mod handlers;
mod models;
mod routes;
mod state;

use axum::Router;
use std::net::SocketAddr;
use tower_http::{
    compression::CompressionLayer,
    cors::{Any, CorsLayer},
    request_id::{MakeRequestUuid, PropagateRequestIdLayer, SetRequestIdLayer},
    trace::TraceLayer,
};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() {
    // Initialize tracing
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "api_server=debug,tower_http=debug".into()),
        )
        .with(tracing_subscriber::fmt::layer().json())
        .init();

    // Build application
    let app = app();

    // Start server
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000")
        .await
        .unwrap();
    
    tracing::info!("Server starting on {}", listener.local_addr().unwrap());
    
    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
}

fn app() -> Router {
    Router::new()
        .merge(routes::health::router())
        .merge(handlers::users::router())
        // Middleware layers (bottom to top)
        .layer(CompressionLayer::new())
        .layer(SetRequestIdLayer::new(MakeRequestUuid))
        .layer(PropagateRequestIdLayer::new())
        .layer(TraceLayer::new_for_http())
        .layer(
            CorsLayer::new()
                .allow_origin(Any)
                .allow_methods(Any)
                .allow_headers(Any),
        )
}

async fn shutdown_signal() {
    let ctrl_c = async {
        tokio::signal::ctrl_c()
            .await
            .expect("Failed to install Ctrl+C handler");
    };

    #[cfg(unix)]
    let terminate = async {
        tokio::signal::unix::signal(tokio::signal::unix::SignalKind::terminate())
            .expect("Failed to install SIGTERM handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c => {},
        _ = terminate => {},
    }

    tracing::info!("Shutdown signal received");
}
```

## Testing (tests/api_test.rs)

```rust
use axum::{
    body::Body,
    http::{Request, StatusCode},
};
use tower::ServiceExt;

#[tokio::test]
async fn health_check() {
    let app = crate::app();
    
    let response = app
        .oneshot(
            Request::builder()
                .method("GET")
                .uri("/health")
                .body(Body::empty())
                .unwrap(),
        )
        .await
        .unwrap();
    
    assert_eq!(response.status(), StatusCode::OK);
}

#[tokio::test]
async fn create_user() {
    let app = crate::app();
    
    let response = app
        .oneshot(
            Request::builder()
                .method("POST")
                .uri("/users")
                .header("Content-Type", "application/json")
                .body(Body::from(r#"{"name":"Test","email":"test@example.com"}"#))
                .unwrap(),
        )
        .await
        .unwrap();
    
    assert_eq!(response.status(), StatusCode::CREATED);
}
```
