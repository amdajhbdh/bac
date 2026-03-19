# Axum Middleware & Layer Patterns

## Basic Middleware Setup with Tower

Axum uses Tower for middleware, providing battle-tested layers:

```rust
use axum::{
    Router,
    routing::get,
};
use tower_http::{
    cors::{CorsLayer, Any},
    trace::TraceLayer,
    compression::CompressionLayer,
};
use std::time::Duration;

fn app() -> Router {
    Router::new()
        .route("/", get(handler))
        .layer(TraceLayer::new_for_http())           // Logging/tracing
        .layer(CompressionLayer::new())               // Gzip compression
        .layer(CorsLayer::permissive())               // CORS headers
        .layer(TimeoutLayer::new(Duration::from_secs(30))); // Request timeout
}

async fn handler() -> &'static str {
    "Hello, World!"
}
```

## Multiple Middleware with ServiceBuilder

```rust
use tower::ServiceBuilder;
use tower_http::{
    cors::CorsLayer,
    trace::TraceLayer,
    compression::CompressionLayer,
    request_id::MakeRequestUuid,
};
use axum::Router;

let app = Router::new()
    .route("/", get(handler))
    .layer(
        ServiceBuilder::new()
            .layer(TraceLayer::new_for_http())
            .layer(CompressionLayer::new())
            .layer(SetRequestIdLayer::new(MakeRequestUuid, request_id::DefaultMakeRequestUri))
            .layer(CorsLayer::very_permissive())
    );
```

## Custom Middleware from Scratch

```rust
use axum::{
    extract::{ConnectInfo, Request},
    middleware::Next,
    response::Response,
};
use std::net::SocketAddr;
use axum::body::Body;

/// Request logging middleware
pub async fn logging_middleware(
    ConnectInfo(addr): ConnectInfo<SocketAddr>,
    request: Request,
    next: Next,
) -> Response {
    tracing::info!("Request from {} to {}", addr, request.uri());
    next.run(request).await
}

/// Request ID middleware
pub async fn request_id_middleware(
    request: Request,
    next: Next,
) -> Response {
    let request_id = uuid::Uuid::new_v4().to_string();
    
    let mut response = next.run(request).await;
    response.headers_mut().insert(
        "X-Request-ID",
        request_id.parse().unwrap(),
    );
    response
}
```

## Applying Middleware to Specific Routes

```rust
use axum::{
    Router,
    routing::get,
    middleware,
};
use tower_http::validate_request::ValidateRequestHeaderLayer;

// Use `route_layer` for middleware that should only run on matched routes
// This prevents 404s from becoming 401s for unauthenticated routes
let app = Router::new()
    .route("/public", get(public_handler))
    .route("/protected", get(protected_handler))
    .route_layer(ValidateRequestHeaderLayer::bearer_auth())
    .layer(middleware::from_fn(logging_middleware));
```

## Middleware with State

```rust
use axum::{
    extract::State,
    middleware::FromFn,
    Router,
    routing::get,
};
use std::sync::Arc;

#[derive(Clone)]
struct AppState {
    db: DatabasePool,
}

async fn auth_middleware(
    State(state): State<Arc<AppState>>,
    request: Request,
    next: Next,
) -> Response {
    // Check authorization header
    if let Some(auth) = request.headers().get("Authorization") {
        if validate_token(auth, &state.db).await {
            return next.run(request).await;
        }
    }
    
    (StatusCode::UNAUTHORIZED, "Unauthorized").into_response()
}

fn app() -> Router {
    let state = Arc::new(AppState { db: todo!() });
    
    Router::new()
        .route("/api/data", get(data_handler))
        .route_layer(FromFn::new(auth_middleware))
        .with_state(state)
}
```

## Required Headers Middleware

```rust
use axum::{
    Router,
    routing::get,
    middleware::FromFn,
    http::{Request, StatusCode},
    response::IntoResponse,
};
use tower_http::validate_request::ValidateRequestHeaderLayer;

let app = Router::new()
    .route("/", get(handler))
    .route_layer(ValidateRequestHeaderLayer::new(
        "Authorization",
        |value: &str| value.starts_with("Bearer "),
    ));

// Or require multiple headers
use tower_http::validate_request::RequireHeaderLayer;

let app = Router::new()
    .route("/", get(handler))
    .layer(RequireHeaderLayer::if_present("Authorization"))
    .layer(RequireHeaderLayer::of("X-Request-ID"));
```

## Middleware Execution Order

Middleware added with `Router::layer()` runs **bottom to top** (added last runs first):

```rust
Router::new()
    .route("/", get(handler))
    .layer(LayerA)  // Runs 3rd (outermost)
    .layer(LayerB)  // Runs 2nd
    .layer(LayerC) // Runs 1st (innermost, closest to handler)
```

For most cases, order doesn't matter unless middleware has dependencies.
