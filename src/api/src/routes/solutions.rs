use axum::{routing::get, Router};

pub fn routes() -> Router {
    Router::new()
        .route("/api/v1/solutions", get(list_solutions))
        .route("/api/v1/solutions/:id", get(get_solution))
}

async fn list_solutions() -> String {
    "[]".to_string()
}

async fn get_solution() -> String {
    "{}".to_string()
}
