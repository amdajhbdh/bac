use axum::{routing::get, Router};
use serde::{Deserialize, Serialize};

pub fn routes() -> Router {
    Router::new()
        .route("/api/v1/problems", get(list_problems))
        .route("/api/v1/problems/:id", get(get_problem))
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Problem {
    pub id: String,
    pub title: String,
    pub description: String,
    pub subject: Option<String>,
    pub created_at: String,
}

async fn list_problems() -> String {
    "[]".to_string()
}

async fn get_problem() -> String {
    "{}".to_string()
}
