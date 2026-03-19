use axum::{routing::post, Router};

pub fn routes() -> Router {
    Router::new()
        .route("/api/v1/auth/login", post(login))
        .route("/api/v1/auth/register", post(register))
        .route("/api/v1/auth/logout", post(logout))
}

async fn login() -> String {
    r#"{"token": "dummy-token"}"#.to_string()
}

async fn register() -> String {
    r#"{"token": "dummy-token"}"#.to_string()
}

async fn logout() -> String {
    r#"{"success": true}"#.to_string()
}
