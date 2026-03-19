---
source: codesearch
library: reqwest
package: reqwest
topic: POST requests with JSON body
fetched: 2026-03-19T00:00:00Z
official_docs: https://docs.rs/reqwest/latest/reqwest/
---

# POST Requests with JSON Body in reqwest 0.12

## Dependencies

```toml
[dependencies]
reqwest = { version = "0.12", features = ["json"] }
tokio = { version = "1", features = ["full"] }
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
```

## Basic POST with JSON

The `.json()` method automatically serializes your struct to JSON and sets `Content-Type: application/json`:

```rust
use reqwest::Client;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize)]
struct NewUser {
    name: String,
    email: String,
    role: String,
}

#[derive(Debug, Deserialize)]
struct ApiResponse {
    id: u64,
    name: String,
    email: String,
    role: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let client = Client::new();
    
    let new_user = NewUser {
        name: "Alice".to_string(),
        email: "alice@example.com".to_string(),
        role: "admin".to_string(),
    };

    // Send POST request with JSON body
    let response = client
        .post("https://httpbin.org/post")
        .json(&new_user)  // Automatically serializes and sets Content-Type
        .send()
        .await?;

    println!("Status: {}", response.status());
    
    // Parse JSON response
    let body: ApiResponse = response.json().await?;
    println!("Created user: {:?}", body);
    
    Ok(())
}
```

## POST with raw JSON body

```rust
use serde_json::json;

let mut map = std::collections::HashMap::new();
map.insert("lang", "rust");
map.insert("body", "json");

let res = client
    .post("http://httpbin.org/post")
    .json(&map)
    .send()
    .await?;
```

## POST with form data

```rust
// This will POST a body of `foo=bar&baz=quux`
let params = [("foo", "bar"), ("baz", "quux")];
let client = Client::new();
let res = client
    .post("http://httpbin.org/post")
    .form(&params)
    .send()?;
```
