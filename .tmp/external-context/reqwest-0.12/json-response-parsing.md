---
source: codesearch
library: reqwest
package: reqwest
topic: Response JSON parsing
fetched: 2026-03-19T00:00:00Z
official_docs: https://docs.rs/reqwest/latest/reqwest/
---

# Response JSON Parsing in reqwest 0.12

## Basic JSON Response Parsing

```rust
use serde::Deserialize;

#[derive(Debug, Deserialize)]
struct User {
    id: u64,
    name: String,
    email: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let client = Client::new();
    
    let response = client
        .get("https://api.example.com/users/1")
        .send()
        .await?;

    // Only parse if status is OK
    if response.status().is_success() {
        let user: User = response.json().await?;
        println!("User: {:?}", user);
    }

    Ok(())
}
```

## Quick one-liner for simple GET requests

```rust
let user: User = reqwest::get("https://api.example.com/users/1")
    .await?
    .json()
    .await?;
```

## Get raw text response

```rust
let body = response.text().await?;
println!("Body: {}", body);
```

## Handling Union/Tagged Responses

For APIs that return different JSON structures for success/error:

```rust
use serde::Deserialize;

#[derive(Serialize, Deserialize)]
#[serde(untagged)]
enum ApiResponse {
    Success {
        id: u64,
        name: String,
    },
    Error {
        errors: Vec<String>,
    },
}

// Usage
match response.json::<ApiResponse>().await {
    Ok(ApiResponse::Success(data)) => println!("Success: {:?}", data),
    Ok(ApiResponse::Error(err)) => println!("API Error: {:?}", err.errors),
    Err(e) => println!("Parse error: {}", e),
}
```

## Check status before parsing

```rust
use reqwest::StatusCode;

match response.status() {
    StatusCode::OK => {
        let user: User = response.json().await?;
        println!("User: {:?}", user);
    },
    StatusCode::UNAUTHORIZED => {
        println!("Need to grab a new token");
    },
    StatusCode::NOT_FOUND => {
        println!("Resource not found");
    },
    _ => {
        println!("Unexpected status: {}", response.status());
    },
}
```

## Parse to serde_json::Value (dynamic)

```rust
use serde_json::Value;

let json: Value = response.json().await?;
println!("{:#}", json);
```
