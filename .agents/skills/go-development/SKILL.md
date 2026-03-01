# Skill: Go Development

## Purpose

Go development patterns, testing, and project structure for BAC Unified.

## When to use

- Working on any Go code in `src/agent/` or `src/api/`
- Creating new packages or handlers
- Writing tests or benchmarks

## Project Structure

```
src/agent/
├── cmd/main.go           # Entry point
└── internal/
    ├── solver/            # Math problem solver
    ├── memory/            # Vector DB (pgvector)
    ├── ocr/               # OCR (Tesseract)
    ├── analyzer/          # Multi-service analysis
    ├── animation/         # Noon animations
    ├── nlm/               # NotebookLM integration
    ├── online/            # Playwright research
    ├── db/                # Database (sqlc)
    └── kb/                # Knowledge base

src/api/
├── cmd/main.go
└── internal/
    ├── handlers/          # HTTP handlers
    ├── services/          # Business logic
    └── models/            # Data models
```

## Mandatory Patterns

### Imports (grouped)

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

### Error Handling

```go
// GOOD
if err != nil {
    return nil, fmt.Errorf("solve problem: %w", err)
}

// Check sentinel errors
if errors.Is(err, context.DeadlineExceeded) {
    return nil, fmt.Errorf("timeout: %w", err)
}
```

### Logging

```go
// GOOD - structured logging
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.Info("starting server", "port", 8080)
slog.Error("failed to connect", "err", err)

// NEVER use
fmt.Println("Starting server")
print("some value")
```

### Configuration

```go
// GOOD
dbURL := os.Getenv("NEON_DB_URL")
if dbURL == "" {
    dbURL = "postgres://localhost/bac"
}

// BAD
dbURL := "postgres://localhost/bac"
```

### Naming

```go
// Exported: PascalCase
func Solve(ctx context.Context, problem string) (*SolveResult, error)

// Unexported: camelCase
func memoryLookup(ctx context.Context, text string) []Similar

// Constants: PascalCase exported, camelCase unexported
const DefaultTimeout = 30 * time.Second
const maxRetries = 3
```

## Commands

```bash
# Build
cd src/agent && go build -o bin/Agent ./cmd/main.go
cd src/api && go build -o bin/api ./cmd/main.go

# Test
go test -v ./...
go test -v -run TestName ./internal/solver/
go test -cover ./...
go test -race ./...

# Lint
go fmt ./...
go vet ./...
golangci-lint run

# Generate sqlc
sqlc generate
```

## Table-Driven Tests

```go
func TestDetermineSubject(t *testing.T) {
    tests := []struct {
        name    string
        problem string
        want    string
    }{
        {"quadratic", "Résous: x² + x - 2 = 0", "math"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := determineSubject(tt.problem)
            if got != tt.want {
                t.Errorf("got %q, want %q", got, tt.want)
            }
        })
    }
}
```

## Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `NEON_DB_URL` | PostgreSQL + pgvector | - |
| `TURSO_DB_URL` | SQLite cache | - |
| `OLLAMA_HOST` | Local AI | localhost:11434 |
| `GARAGE_ENDPOINT` | S3 storage | localhost:3900 |
| `PORT` | API port | 8080 |

## Dependencies

- **Database**: PostgreSQL + pgvector, Turso (SQLite)
- **AI**: Ollama, NotebookLM
- **Storage**: Garage S3
- **Tools**: sqlc, golangci-lint

## Anti-Patterns

- ❌ Using `fmt.Print*` instead of `log/slog`
- ❌ Hardcoding configuration values
- ❌ Using raw SQL - use sqlc
- ❌ Global mutable state
- ❌ Ignoring errors
- ❌ Not using context for cancellation
