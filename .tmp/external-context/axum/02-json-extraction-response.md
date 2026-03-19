# Axum JSON Extraction & Response

## Extracting JSON Request Body

Use the `Json` extractor to parse JSON from request bodies:

```rust
use axum::{
    extract::Json,
    http::StatusCode,
    response::IntoResponse,
    routing::post,
    Router,
};
use serde::{Deserialize, Serialize};

// Request payload
#[derive(Deserialize)]
struct CreateUser {
    username: String,
    email: String,
}

// Response payload
#[derive(Serialize)]
struct User {
    id: u64,
    username: String,
    email: String,
}

async fn create_user(
    Json(payload): Json<CreateUser>
) -> impl IntoResponse {
    // In a real app, you would save to a database here
    let user = User {
        id: 1,
        username: payload.username,
        email: payload.email,
    };

    // Return 201 Created with the user data
    (StatusCode::CREATED, Json(user))
}
```

## Returning JSON Responses

```rust
use axum::{
    extract::Json,
    response::IntoResponse,
    routing::get,
    Router,
};
use serde::Serialize;

#[derive(Serialize)]
struct Message {
    message: String,
    code: u32,
}

async fn get_message() -> Json<Message> {
    Json(Message {
        message: "Hello, World!".to_string(),
        code: 200,
    })
}

// Multiple ways to return JSON
async fn different_json_responses() -> impl IntoResponse {
    // 1. Just Json (returns 200 OK)
    Json(serde_json::json!({"key": "value"}))

    // 2. Status code + Json (returns custom status)
    (StatusCode::CREATED, Json(serde_json::json!({"id": 1})))

    // 3. With headers
    (
        StatusCode::OK,
        [("X-Custom-Header", "value")],
        Json(serde_json::json!({"data": "test"})),
    )
}
```

## Query Parameters Extraction

```rust
use axum::{
    extract::{Path, Query},
    routing::get,
    Router,
};
use serde::Deserialize;

#[derive(Deserialize)]
struct PaginationParams {
    page: Option<u64>,
    limit: Option<u64>,
}

async fn list_items(Query(params): Query<PaginationParams>) -> String {
    let page = params.page.unwrap_or(1);
    let limit = params.limit.unwrap_or(20);
    format!("Listing items - page: {}, limit: {}", page, limit)
}

// Multiple extractors can be combined
async fn get_user_posts(
    Path(user_id): Path<u64>,
    Query(params): Query<PaginationParams>,
) -> String {
    format!(
        "User {} posts - page: {}, limit: {}",
        user_id,
        params.page.unwrap_or(1),
        params.limit.unwrap_or(10)
    )
}
```

## Custom Request Validation Extractor

```rust
use axum::{
    async_trait,
    extract::{FromRequest, Request},
    http::StatusCode,
    Json,
};
use serde::de::DeserializeOwned;
use validator::Validate;

/// JSON extractor that automatically validates the payload
pub struct ValidatedJson<T>(pub T);

#[async_trait]
impl<S, B, T> FromRequest<S, B> for ValidatedJson<T>
where
    B: http_body::Body + Send,
    B::Data: Send,
    B::Error: Into<axum::Error>,
    S: Send + Sync,
    T: DeserializeOwned + Validate,
{
    type Rejection = (StatusCode, String);

    async fn from_request(req: Request<B>, state: &S) -> Result<Self, Self::Rejection> {
        let Json(value) = Json::<T>::from_request(req, state)
            .await
            .map_err(|e| (StatusCode::BAD_REQUEST, e.to_string()))?;

        value.validate().map_err(|e| {
            (StatusCode::UNPROCESSABLE_ENTITY, e.to_string())
        })?;

        Ok(ValidatedJson(value))
    }
}

// Usage
async fn create_post(
    ValidatedJson(payload): ValidatedJson<CreatePost>,
) -> Json<Post> {
    // payload is already validated
    todo!()
}
```
