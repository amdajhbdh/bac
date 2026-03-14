# Integration Specifications

## Overview

Spec-driven, test-first integration specifications for BAC Unified.

---

## 1. OCR → Gateway Integration

### Contract

**OCR Service Endpoint:**
```
POST /ocr
Content-Type: multipart/form-data

Input: image (file) or pdf (file)
Output: { "text": "...", "confidence": 0.95, "language": "fr" }
```

**Gateway calls OCR:**
```rust
// src/gateway/src/main.rs
async fn ocr_handler(mut ctx: Context) -> Result<impl IntoResponse, StatusCode> {
    let client = reqwest::Client::new();
    let ocr_url = env::var("OCR_SERVICE_URL").unwrap_or_else(|_| "http://localhost:8081".to_string());
    
    let response = client
        .post(format!("{}/ocr", ocr_url))
        .multipart(form)
        .send()
        .await
        .map_err(|_| StatusCode::BAD_GATEWAY)?;
    
    // ...
}
```

### Tests (TDD)

```rust
#[cfg(test)]
mod integration_tests {
    use super::*;
    
    #[tokio::test]
    async fn test_gateway_calls_ocr_service() {
        // Arrange
        let mock_server = mock_server(|m| {
            m.post("/ocr").return_status(200).body(r#"{"text":"test","confidence":0.9}"#);
        }).start();
        
        // Act
        let result = call_ocr(&mock_server.url()).await;
        
        // Assert
        assert!(result.is_ok());
    }
    
    #[tokio::test]
    async fn test_ocr_timeout_fails_gracefully() {
        // Arrange - slow OCR server
        let mock_server = mock_server(|m| {
            m.post("/ocr").delay(10_000).return_status(200);
        }).start();
        
        // Act
        let result = call_ocr_with_timeout(&mock_server.url(), Duration::from_millis(100)).await;
        
        // Assert
        assert!(result.is_err());
    }
}
```

---

## 2. Gateway → Agent Integration

### Contract

**Gateway Endpoints:**
```
POST /chat
POST /animation
GET  /animation/:id
```

**Agent calls Gateway:**
```go
// src/agent/cmd/main.go
func callGateway(endpoint string, payload []byte) (*GatewayResponse, error) {
    gatewayURL := os.Getenv("GATEWAY_URL")
    if gatewayURL == "" {
        gatewayURL = "http://localhost:8080"
    }
    
    resp, err := http.Post(gatewayURL+endpoint, "application/json", bytes.NewBuffer(payload))
    // ...
}
```

### Tests (TDD)

```go
// src/agent/internal/gateway/gateway_test.go
func TestGatewayChat(t *testing.T) {
    // Arrange
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message":"test response"}`))
    }))
    defer server.Close()
    
    // Act
    resp, err := CallGateway(server.URL+"/chat", ChatRequest{Message: "hello"})
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test response", resp.Message)
}
```

---

## 3. RAG → PostgreSQL Integration

### Contract

**RAG searches pgvector:**
```rust
// src/gateway/src/rag/mod.rs
impl RAGEngine {
    pub async fn search(&self, query: &str, limit: usize) -> Vec<SearchResult> {
        // Connect to PostgreSQL using postgres-pgvector skill
        let conn = PgConnection::connect(std::env::var("DATABASE_URL").unwrap()).await;
        
        // Generate embedding
        let embedding = generate_embedding(query).await;
        
        // Query pgvector
        conn.query(
            "SELECT text, 1 - (embedding <=> $1) as similarity FROM documents ORDER BY embedding <=> $1 LIMIT $2",
            embedding,
            limit,
        ).await
    }
}
```

### Tests (TDD)

```rust
#[cfg(test)]
mod pgvector_tests {
    use super::*;
    
    #[tokio::test]
    async fn test_rag_search_returns_similar_documents() {
        // Arrange
        let engine = RAGEngine::new();
        
        // Act
        let results = engine.search("dérivée", 5).await;
        
        // Assert
        assert!(!results.is_empty());
        assert!(results[0].similarity > 0.7);
    }
    
    #[tokio::test]
    async fn test_rag_filtered_search() {
        // Arrange
        let engine = RAGEngine::new();
        
        // Act
        let results = engine.search_filtered(
            "intégrale",
            Some("mathematics"),
            Some(3), // chapter 3
            10,
        ).await;
        
        // Assert
        assert!(results.iter().all(|r| r.metadata["chapter"] == 3));
    }
}
```

---

## 4. Solver → Animation Integration

### Contract

**Solver triggers animation:**
```go
// src/agent/internal/solver/solver.go
func (s *Solver) GenerateAnimation(ctx context.Context, problem string) (string, error) {
    // Call gateway animation endpoint
    payload := AnimationRequest{
        Prompt:    problem,
        Type:      "manim",
        Quality:   "medium",
    }
    
    resp, err := callGateway("/animation", payload)
    if err != nil {
        return "", err
    }
    
    return resp.VideoURL, nil
}
```

### Tests (TDD)

```go
func TestSolverGeneratesAnimation(t *testing.T) {
    // Arrange
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id":"123","status":"processing","video_url":""}`))
    }))
    defer server.Close()
    
    solver := NewSolver(server.URL)
    
    // Act
    result, err := solver.GenerateAnimation(context.Background(), "solve x^2 + 2x + 1 = 0")
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
}
```

---

## Environment Variables

| Variable | Service | Purpose |
|----------|---------|---------|
| `OCR_SERVICE_URL` | Gateway | OCR HTTP endpoint |
| `GATEWAY_URL` | Agent | Gateway HTTP endpoint |
| `DATABASE_URL` | RAG | PostgreSQL connection |
| `OLLAMA_HOST` | Solver | Ollama endpoint |
| `GARAGE_ENDPOINT` | Storage | S3 endpoint |

---

## Test Execution

```bash
# Run integration tests
cd src/gateway && cargo test integration
cd src/agent && go test ./internal/gateway/...

# Run with coverage
cd src/gateway && cargo test --cover
cd src/agent && go test -cover ./...
```

---

*Spec Version: 1.0*
*Created: 2026-03-06*
