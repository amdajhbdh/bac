# Axum Error Handling Patterns

## Centralized Error Handling with ThisError

```rust
use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use serde::Serialize;
use thiserror::Error;

/// Application error types
#[derive(Error, Debug)]
pub enum AppError {
    #[error("Resource not found: {0}")]
    NotFound(String),

    #[error("Validation error: {0}")]
    ValidationError(String),

    #[error("Unauthorized: {0}")]
    Unauthorized(String),

    #[error("Internal server error")]
    InternalError(#[from] anyhow::Error),
}

/// Error response body
#[derive(Serialize)]
struct ErrorResponse {
    error: String,
    message: String,
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, error_message) = match &self {
            AppError::NotFound(msg) => (StatusCode::NOT_FOUND, msg.clone()),
            AppError::ValidationError(msg) => (StatusCode::BAD_REQUEST, msg.clone()),
            AppError::Unauthorized(msg) => (StatusCode::UNAUTHORIZED, msg.clone()),
            AppError::InternalError(msg) => {
                tracing::error!("Internal error: {}", msg);
                (StatusCode::INTERNAL_SERVER_ERROR, "Internal server error".to_string())
            }
        };

        let body = Json(ErrorResponse {
            error: self.to_string(),
            message: error_message,
        });

        (status, body).into_response()
    }
}

/// Result type alias for convenience
pub type Result<T> = std::result::Result<T, AppError>;
```

## Using Error Handling in Handlers

```rust
use crate::error::{AppError, Result};
use axum::{extract::Path, Json, Router, routing::get};
use serde::{Deserialize, Serialize};

#[derive(Serialize)]
struct User {
    id: u64,
    name: String,
}

// Handler returning Result with automatic error conversion
async fn get_user(Path(id): Path<u64>) -> Result<Json<User>> {
    let user = find_user_by_id(id)
        .await
        .map_err(|_| AppError::NotFound(format!("User {} not found", id)))?;

    Ok(Json(user))
}

async fn find_user_by_id(id: u64) -> anyhow::result::Result<User> {
    // Simulated database lookup
    if id == 0 {
        anyhow::bail!("User not found");
    }
    Ok(User { id, name: "John".to_string() })
}
```

## Handling Rejections with Custom Fallback

```rust
use axum::{
    body::Body,
    extract::{Request, State},
    http::StatusCode,
    response::Response,
    routing::get,
    Router,
};
use serde::Serialize;
use std::convert::Infallible;

#[derive(Serialize)]
struct RejectionBody {
    error: String,
    status: u16,
}

async fn handle_rejection(
    State(_state): State<()>,
    request: Request,
) -> Response {
    let status = request
        .extensions()
        .get::<axum::extract::rejection::PathRejection>()
        .map(|p| p.status())
        .or_else(|| {
            request
                .extensions()
                .get::<axum::extract::rejection::JsonRejection>()
                .map(|p| p.status())
        })
        .unwrap_or(StatusCode::INTERNAL_SERVER_ERROR);

    let body = RejectionBody {
        error: "Request rejected".to_string(),
        status: status.as_u16(),
    };

    (status, Json(body)).into_response()
}

// Apply to router
let app = Router::new()
    .route("/", get(handler))
    .route("/fallback", get(fallback_handler))
    .fallback(handle_rejection);
```

## Using HandleErrorLayer for Middleware

```rust
use axum::{
    error_handling::HandleErrorLayer,
    BoxError,
    Router,
    routing::get,
};
use tower::ServiceBuilder;
use std::time::Duration;

let app = Router::new()
    .route("/", get(root))
    .layer(
        ServiceBuilder::new()
            .layer(HandleErrorLayer::new(|error: BoxError| async move {
                if error.is::<tower::timeout::Elapsed>() {
                    (
                        StatusCode::REQUEST_TIMEOUT,
                        Json(serde_json::json!({"error": "Request timeout"})),
                    )
                } else {
                    (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        Json(serde_json::json!({"error": "Internal server error"})),
                    )
                }
            }))
            .timeout(Duration::from_secs(30))
    );
```
