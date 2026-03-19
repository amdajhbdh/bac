# Axum 0.8.x Router & Routes Setup

## Basic Router Setup

```rust
use axum::{
    routing::{get, post},
    Router,
};
use std::net::SocketAddr;

#[tokio::main]
async fn main() {
    // Initialize tracing for request logging
    tracing_subscriber::fmt::init();

    // Build the router with routes
    let app = Router::new()
        .route("/", get(root_handler))
        .route("/health", get(health_check))
        .route("/users", get(list_users).post(create_user))
        .route("/users/{id}", get(get_user).put(update_user).delete(delete_user));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

// Handler functions return anything that implements IntoResponse
async fn root_handler() -> &'static str {
    "Hello, World!"
}

async fn health_check() -> &'static str {
    "OK"
}

async fn list_users() -> &'static str {
    "Listing all users"
}

async fn create_user() -> &'static str {
    "Creating user"
}

async fn get_user(Path(user_id): Path<u64>) -> String {
    format!("Fetching user with ID: {}", user_id)
}

async fn update_user(Path(user_id): Path<u64>) -> String {
    format!("Updating user with ID: {}", user_id)
}

async fn delete_user(Path(user_id): Path<u64>) -> String {
    format!("Deleting user with ID: {}", user_id)
}
```

## Path Parameters (Axum 0.8+ Syntax)

**Axum 0.8 introduced new path syntax using curly braces:**
```rust
// Old syntax (Axum 0.7 and earlier)
.route("/users/:id", get(get_user))

// New syntax (Axum 0.8+)
.route("/users/{id}", get(get_user))

// Nested path parameters
.route("/users/{user_id}/posts/{post_id}", get(get_user_post))
```

## Nested Routes & API Versioning

```rust
fn user_routes() -> Router {
    Router::new()
        .route("/", get(list_users).post(create_user))
        .route("/{id}", get(get_user).put(update_user).delete(delete_user))
}

fn post_routes() -> Router {
    Router::new()
        .route("/", get(list_posts).post(create_post))
        .route("/{id}", get(get_post))
}

fn api_v1_routes() -> Router {
    Router::new()
        .nest("/users", user_routes())
        .nest("/posts", post_routes())
}

fn app() -> Router {
    Router::new()
        .nest("/api/v1", api_v1_routes())
        .route("/health", get(health_check))
}
```

## Dependencies (Cargo.toml)

```toml
[dependencies]
axum = "0.7"  # or 0.8 for latest
tokio = { version = "1", features = ["full"] }
serde = { version = "1", features = ["derive"] }
serde_json = "1"
tower = "0.4"
tower-http = { version = "0.5", features = ["cors", "trace"] }
tracing = "0.1"
tracing-subscriber = { version = "0.3", features = ["env-filter"] }
```
