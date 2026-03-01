# BAC Unified - Agent Instructions

## Project Overview

| Component | Language | Path |
|-----------|----------|------|
| Agent CLI | **Go** | `src/agent/` |
| REST API | Go (Gin) | `src/api/` |
| Animations | Rust (Noon) | `src/noon/` |

## Build Commands

### Go (Agent & API)

```bash
cd src/agent

# Build binary
go build -o bin/Agent ./cmd/main.go

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o bin/Agent-linux ./cmd/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/Agent-darwin ./cmd/main.go

# Run with flags
./bin/Agent -animate "problem"     # Generate animation
./bin/Agent -image photo.jpg       # OCR from image
./bin/Agent -pdf exam.pdf          # OCR from PDF
./bin/Agent -analyze resource      # Full analysis
./bin/Agent -online                # Force online auto-solve
```

### Rust (Noon)

```bash
cd src/noon
cargo build --release
cargo test
cargo check
cargo fmt
cargo clippy -- -D warnings
```

## Testing

```bash
cd src/agent

# All tests
go test ./...

# Single test (most common)
go test -v -run TestSolve ./internal/solver/

# With coverage
go test -cover ./...

# With race detector
go test -race ./...

# Benchmark
go test -bench=BenchmarkMemoryLookup ./internal/memory/
```

### Test Naming

- Test files: `*_test.go`
- Unit tests: `TestFunctionName`
- Benchmarks: `BenchmarkFunctionName`

## Code Style

### Logging: Use `log/slog`

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.Info("starting server", "port", 8080)
slog.Error("failed to connect", "err", err)

// NEVER use: fmt.Println, print()
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

## Database

- Use **sqlc** for all SQL queries: `cd src/agent && sqlc generate`
- PostgreSQL + pgvector for main DB
- Turso (SQLite) for production cache only
- Schema: `sql/schema.sql`

## Cursor Rules (from .cursorrules)

- **VCS**: Use jj (Jujutsu). NEVER use git directly. `jj new`, `jj commit`, `jj push`
- **DB Tools**: sqlc for all queries. Run `sqlc generate` after schema changes
- **Search**: Hybrid = Elasticsearch (full-text) + pgvector (vectors)
- **OCR**: Tesseract + Surya. NEVER use paid APIs
- **Auth**: JWT + bcrypt. Environment variables for ALL secrets
- **AI Priority**: Ollama → NotebookLM → Cloud (fallback)
- **Monitoring**: Prometheus + Grafana for latency, errors, cache hit rate

## Code Review Checklist

- [ ] Tests pass (`go test ./...`)
- [ ] Code formatted (`go fmt`)
- [ ] No vet warnings (`go vet`)
- [ ] Race clean (`go test -race`)
- [ ] Logging uses `slog`
- [ ] Errors wrapped
- [ ] Env vars for config
- [ ] Check openspec/changes/full-bac-unified/tasks.md before working

---

*Last Updated: 2026-03-01*
