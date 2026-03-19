use axum::{routing::get, Router};

pub fn routes() -> Router {
    Router::new()
        .route("/api/v1/users", get(list_users))
        .route("/api/v1/users/:id", get(get_user))
}

async fn list_users() -> String {
    "[]".to_string()
}

async fn get_user() -> String {
    "{}".to_string()
}
