---
source: codesearch
library: reqwest
package: reqwest
topic: Setting headers Authorization Content-Type
fetched: 2026-03-19T00:00:00Z
official_docs: https://docs.rs/reqwest/latest/reqwest/
---

# Setting Headers in reqwest 0.12

## Method 1: Per-request headers with `.header()`

```rust
use reqwest::{Client, header};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let client = Client::new();

    let response = client
        .get("https://api.example.com/data")
        .header("Authorization", "Bearer your-token-here")
        .header("Accept", "application/json")
        .header("X-Custom-Header", "custom-value")
        .send()
        .await?;

    Ok(())
}
```

## Method 2: Default headers for all requests

```rust
use reqwest::{Client, header};

let mut headers = header::HeaderMap::new();
headers.insert(
    header::AUTHORIZATION,
    header::HeaderValue::from_static("Bearer your-token-here")
);
headers.insert(
    header::ACCEPT,
    header::HeaderValue::from_static("application/json")
);

let client_with_headers = Client::builder()
    .default_headers(headers)
    .build()?;
```

## Setting Content-Type manually

```rust
use reqwest::header::{CONTENT_TYPE, HeaderValue};

let response = client
    .post("https://api.example.com/data")
    .header(CONTENT_TYPE, HeaderValue::from_static("application/json"))
    .body(json_data.to_string())
    .send()
    .await?;
```

## Basic Authentication

```rust
let res = client
    .get("http://httpbin.org/basic-auth/user/pass")
    .basic_auth("user", Some("pass"))
    .send()
    .await?;
```

## Bearer Token Authentication

```rust
let response = client
    .get("https://api.spotify.com/v1/search")
    .header(AUTHORIZATION, "Bearer [AUTH_TOKEN]")
    .header(CONTENT_TYPE, "application/json")
    .header(ACCEPT, "application/json")
    .send()
    .await?;
```

## Custom User-Agent

```rust
use reqwest::header::{HeaderMap, USER_AGENT, HeaderValue};

let mut headers = HeaderMap::new();
headers.insert(USER_AGENT, HeaderValue::from_static("my-app/1.0"));
headers.insert(CONTENT_TYPE, HeaderValue::from_static("application/json"));

let client = Client::builder()
    .default_headers(headers)
    .build()?;
```
