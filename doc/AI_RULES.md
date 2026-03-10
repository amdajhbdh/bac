# BAC Unified Platform - AI & Rules Documentation

## Overview

The BAC Unified platform is a comprehensive educational AI system for baccalaureate exam preparation, featuring:
- OCR for document processing (Tesseract, OCR.space, Google Vision)
- Animation generation (Manim via Podman)
- Chat modes (Query, Chat, Agent, Automatic)
- Hybrid search (Elasticsearch + pgvector)
- Multi-provider AI fallback (Ollama → NLM → Cloud APIs)

## AI Inference Priority

```
Ollama (llama3.2:3b, local, fast)
    → Cloud APIs (kimi/deepseek/minimax/ministral, medium)
    → NLM CLI/MCP (NotebookLM with caching)
```

## Chat Modes

### Mode Selection

| Mode | Trigger Keywords | Use Case |
|------|-----------------|----------|
| **Query** | who, what, where, when, how, define | Quick factual answers |
| **Chat** | hello, hi, thanks, please | Conversational interaction |
| **Agent** | solve, calculate, find, search, look up | Complex tasks requiring tools |
| **Automatic** | (none or complex) | Auto-detect based on query analysis |

### Auto-Select Logic

```rust
fn select_mode(query: &str) -> ChatMode {
    let query_lower = query.to_lowercase();
    
    // Query mode keywords
    let query_keywords = ["who", "what", "where", "when", "how", "define", "explain"];
    if query_keywords.iter().any(|k| query_lower.contains(k)) {
        return ChatMode::Query;
    }
    
    // Chat mode keywords
    let chat_keywords = ["hello", "hi", "thanks", "please", "goodbye"];
    if chat_keywords.iter().any(|k| query_lower.contains(k)) {
        return ChatMode::Chat;
    }
    
    // Agent mode keywords
    let agent_keywords = ["solve", "calculate", "find", "search", "look up", "compare"];
    if agent_keywords.iter().any(|k| query_lower.contains(k)) {
        return ChatMode::Agent;
    }
    
    ChatMode::Automatic
}
```

## OCR Pipeline

### Engine Priority

```
1. Tesseract (local, fast)
   → Surya (if enabled, needs Python)
   → OCR.space (free cloud, auto-language)
   → Google Vision (requires API key)
```

### Preprocessing Pipeline

1. Convert to grayscale
2. Denoise (median filter × 2)
3. Contrast enhancement (2.0x)
4. Brightness adjustment (1.1x)
5. Sharpen filter
6. Binarize (threshold)

### Supported Languages

- French (fra)
- Arabic (ara)  
- English (eng)

## Animation System

### Manim via Podman

```rust
// Animation request
{
    "code": "class Test(Scene):\n    def construct(self):\n        self.play(Write(Text('Hello')))"
}

// Response
{
    "job_id": "uuid",
    "status": "Completed",
    "output_url": "/path/to/video.mp4"
}
```

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/animation` | POST | Create animation job |
| `/animation/:id` | GET | Get job status |
| `/animations` | GET | List all jobs |

### Quality Options

- `l` (low) - 854×480
- `m` (medium) - 1280×720
- `h` (high) - 1920×1080
- `p` (production) - 2560×1440

## Search System

### Hybrid Search (Reciprocal Rank Fusion)

1. Parse query, extract filters (subject, chapter, difficulty, language, source)
2. Run in parallel: 
   - Elasticsearch full-text (French analyzer)
   - pgvector cosine similarity (1536-dim)
3. Fuse results with RRF normalization
4. Return with facets + suggestions

## Code Style Rules

### Logging: Use `log/slog`

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.Info("starting server", "port", 8080)
slog.Error("failed to connect", "err", err)
```

### Imports: Group stdlib, external, internal

```go
import (
    // Standard library
    "context"
    "encoding/json"
    "log/slog"
    "os"
    "time"

    // External packages
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    // Internal packages
    "github.com/bac-unified/agent/internal/memory"
    "github.com/bac-unified/agent/internal/solver"
)
```

### Error Handling: Wrap with context

```go
if err != nil {
    return nil, fmt.Errorf("solve with ollama: %w", err)
}

if errors.Is(err, context.DeadlineExceeded) {
    return nil, fmt.Errorf("timeout: %w", err)
}
```

### Naming Conventions

```go
// Exported: PascalCase
func Solve(ctx context.Context, problem string) (*SolveResult, error)

// Unexported: camelCase
func memoryLookup(ctx context.Context, text string) []Similar

// Constants: PascalCase exported, camelCase unexported
const DefaultTimeout = 30 * time.Second
const maxRetries = 3
```

### Config: Environment variables only

```go
dbURL := os.Getenv("NEON_DB_URL")
if dbURL == "" {
    dbURL = "postgres://localhost/bac"
}
// NEVER hardcode credentials or URLs
```

## Database Rules

- **Main DB**: PostgreSQL + pgvector
- **NLM cache**: Turso (SQLite)
- **Hot cache**: Redis
- Use **sqlc** for all SQL queries
- Never use raw SQL strings or ORM

## Storage Rules

- **Garage S3** or MinIO only (S3-compatible API)
- Never AWS S3 direct or GCS

## VCS Rules

- Use **jj (Jujutsu)**. NEVER use `git` directly.
- `jj new`, `jj commit`, `jj push`, `jj log`

## Monitoring

- Prometheus + Grafana
- Track: latency, errors, cache hit rate

## Deployment

### Container-Ready

- Multi-stage Dockerfile
- Health checks
- Configuration via environment variables
- Graceful shutdown

### Security

- JWT + bcrypt for auth
- All secrets via environment variables
- Rate limiting (100 requests/minute)
- Input validation

## Key Dependencies

| Component | Language | Version |
|-----------|----------|---------|
| Agent CLI | Go | 1.25 |
| Gateway | Rust | 2021 |
| Animation | Rust (Noon) | latest |
| OCR | Rust | - |
| HTTP | Axum | 0.7+ |
| DB | PostgreSQL + pgvector | 15+ |

## Testing Requirements

### Test Naming

- Test files: `*_test.go` / `*_test.rs`
- Unit tests: `TestFunctionName`
- Benchmarks: `BenchmarkFunctionName`

### Test Commands

```bash
# Go tests
cd src/agent && go test ./...

# Rust tests  
cd src/gateway && cargo test
cd src/ocr-service && cargo test
```

## Architecture

```
Nushell CLI (lib/bac.nu)
    └── Go Agent CLI (src/agent/cmd/main.go)
            ├── solver/       Ollama → Cloud APIs → NLM → Online research
            ├── memory/       pgvector similarity search (RAG)
            ├── nlm/         NotebookLM + Turso cache
            ├── ocr/         Tesseract + Surya + Google Vision
            ├── search/      Elasticsearch + pgvector
            ├── animation/   Noon DSL code generation
            └── bundle/      WRAPX + ARX
```
