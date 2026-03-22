//! Error types for vector-tools

use thiserror::Error;

/// Vector-tools error types
#[derive(Debug, Error)]
pub enum VectorError {
    #[error("Database error: {0}")]
    Database(#[from] tokio_postgres::Error),

    #[error("Pool error: {0}")]
    Pool(#[from] deadpool_postgres::PoolError),

    #[error("Serialization error: {0}")]
    Serialization(#[from] serde_json::Error),

    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),

    #[error("Pool creation error: {0}")]
    PoolCreate(#[from] deadpool_postgres::CreatePoolError),

    #[error("Validation error: {0}")]
    Validation(String),

    #[error("Not found: {0}")]
    NotFound(String),

    #[error("Index error: {0}")]
    Index(String),
}

impl axum::response::IntoResponse for VectorError {
    fn into_response(self) -> axum::response::Response {
        let (_status, message) = match &self {
            VectorError::NotFound(_) => (axum::http::StatusCode::NOT_FOUND, self.to_string()),
            VectorError::Validation(_) => (axum::http::StatusCode::BAD_REQUEST, self.to_string()),
            _ => (
                axum::http::StatusCode::INTERNAL_SERVER_ERROR,
                self.to_string(),
            ),
        };

        axum::Json(serde_json::json!({
            "error": message,
            "type": std::any::type_name::<Self>()
        }))
        .into_response()
    }
}
