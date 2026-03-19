---
source: codesearch
library: reqwest
package: reqwest
topic: Error handling patterns
fetched: 2026-03-19T00:00:00Z
official_docs: https://docs.rs/reqwest/latest/reqwest/
---

# Error Handling Patterns in reqwest 0.12

## Using reqwest::Error directly

```rust
use reqwest::Client;

#[tokio::main]
async fn main() -> Result<(), reqwest::Error> {
    let client = Client::new();
    
    let response = client
        .get("https://api.example.com/data")
        .send()
        .await?;

    // Check status after send
    if !response.status().is_success() {
        let body = response.text().await?;
        eprintln!("API error: {} - {}", response.status(), body);
        return Err(/* return appropriate error */);
    }

    Ok(())
}
```

## Check specific error conditions

```rust
match result {
    Ok(res) => {
        if res.status().is_client_error() {
            return Err("Server returned a client error".into());
        } else if res.status().is_server_error() {
            return Err("Server returned a server error".into());
        }
        // Process response
    },
    Err(err) => {
        if err.is_connect() {
            println!("Connection error");
        } else if err.is_timeout() {
            println!("Request timed out");
        } else if err.is_decode() {
            println!("Failed to decode response");
        }
        return Err(err.into());
    }
}
```

## Custom Error Types with thiserror

```rust
use thiserror::Error;
use serde::Deserialize;

#[derive(Debug, Error)]
pub enum ApiClientError {
    #[error("HTTP request failed: {0}")]
    Reqwest(#[from] reqwest::Error),

    #[error("Failed to parse JSON response: {0}")]
    Serde(#[from] serde_json::Error),

    #[error("API returned an error: {message} (code: {code})")]
    Api { code: u16, message: String },
}

#[derive(Debug, Deserialize)]
struct ApiErrorResponse {
    message: String,
    code: u16,
}

impl From<reqwest::Error> for ApiClientError {
    fn from(err: reqwest::Error) -> Self {
        ApiClientError::Reqwest(err)
    }
}
```

## Full API Client with Error Handling

```rust
use reqwest::Client;
use serde::{Deserialize, Serialize};
use thiserror::Error;
use std::time::Duration;

#[derive(Debug, Error)]
pub enum TodoClientError {
    #[error("HTTP request failed: {0}")]
    Reqwest(#[from] reqwest::Error),
    
    #[error("Failed to parse JSON: {0}")]
    Serde(#[from] serde_json::Error),
    
    #[error("API error: {message} (code: {code})")]
    Api { code: u16, message: String },
}

pub struct ApiClient {
    client: Client,
    base_url: String,
}

impl ApiClient {
    pub fn new(base_url: &str) -> Result<Self, reqwest::Error> {
        let client = Client::builder()
            .timeout(Duration::from_secs(10))
            .build()?;
        
        Ok(Self {
            client,
            base_url: base_url.to_string(),
        })
    }

    pub async fn create_todo(&self, title: &str) -> Result<Todo, TodoClientError> {
        let url = format!("{}/todos", self.base_url);
        let new_todo = serde_json::json!({ "title": title });

        let response = self.client
            .post(&url)
            .json(&new_todo)
            .send()
            .await?;

        if !response.status().is_success() {
            // Try to parse error response
            if let Ok(error) = response.json::<ApiErrorResponse>().await {
                return Err(TodoClientError::Api {
                    code: response.status().as_u16(),
                    message: error.message,
                });
            }
            return Err(TodoClientError::Api {
                code: response.status().as_u16(),
                message: "Unknown error".to_string(),
            });
        }

        let todo: Todo = response.json().await?;
        Ok(todo)
    }
}
```

## Common reqwest Error Methods

```rust
let result = client.get("https://example.com").send().await;

// Check error type
if result.is_err() {
    let err = result.unwrap_err();
    
    err.is_connect();      // Connection failed
    err.is_timeout();      // Request timed out
    err.is_redirect();    // Too many redirects
    err.is_decode();       // Failed to decode response body
    err.status_code();     // Some(StatusCode) if server responded
}
```

## Best Practices

1. **Always configure timeouts** to prevent hanging requests:
   ```rust
   Client::builder()
       .timeout(Duration::from_secs(10))
       .build()?
   ```

2. **Check status codes** before parsing JSON:
   ```rust
   if response.status().is_success() {
       let data: MyType = response.json().await?;
   }
   ```

3. **Use custom error types** for better error context in production:
   ```rust
   #[derive(Error, Debug)]
   enum ApiError {
       #[error("Network error: {0}")]
       Network(#[from] reqwest::Error),
       #[error("Not found")]
       NotFound,
       #[error("Unauthorized")]
       Unauthorized,
   }
   ```

4. **Handle API error responses** with tagged enums:
   ```rust
   #[derive(Deserialize)]
   #[serde(untagged)]
   enum Response {
       Success(Data),
       Error(ErrorResponse),
   }
   ```
