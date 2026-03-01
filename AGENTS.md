# BAC Unified - Agent Instructions

<!--toc:start-->
- [BAC Unified - Agent Instructions](#bac-unified---agent-instructions)
  - [Project Vision](#project-vision)
  - [Project Overview](#project-overview)
  - [Quick Start](#quick-start)
  - [Build Commands](#build-commands)
  - [Testing](#testing)
  - [Code Style](#code-style)
  - [Agent Components](#agent-components)
  - [Database](#database)
  - [External Services](#external-services)
  - [Environment Variables](#environment-variables)
  - [Common Tasks](#common-tasks)
  - [Documentation](#documentation)
  - [Development Workflow](#development-workflow)
  - [OpenSpec Workflow](#openspec-workflow)
<!--toc:end-->

## Project Vision

**BAC Unified** is a national exam preparation system for Mauritania's BAC C students:
- Multi-modal input (images, PDFs, URLs)
- AI-powered exercise solver with step-by-step solutions
- Question prediction based on historical patterns
- Animation visualizations using **noon** (Rust-based)
- Full audit trail
- **Offline-first** - Local processing before cloud

## Project Overview

| Component | Language | Path |
|-----------|----------|------|
| Agent CLI | **Go** | `src/agent/` |
| REST API | Go (Gin) | `src/api/` |
| Animations | Rust (Noon) | `src/noon/` |
| Legacy CLI | Nushell | `bac.nu`, `lib/` |

- **Database**: PostgreSQL + pgvector (Neon), Turso (SQLite for production cache)
- **Storage**: Garage S3 (local), MinIO (production)
- **AI**: Ollama (local), NotebookLM, Cloud APIs (fallback)

## Quick Start

```bash
# Clone and setup
git clone https://github.com/bac-unified/bac.git
cd bac

# Agent - build and run
cd src/agent
go build -o bin/Agent ./cmd/main.go
./bin/Agent "solve x² + x - 2 = 0"

# API - build and run
cd src/api
go build -o bin/api ./cmd/main.go
./bin/api

# Run tests
cd src/agent && go test -v ./...
```

## Build Commands

### Agent CLI

```bash
cd src/agent

# Build binary
go build -o bin/Agent ./cmd/main.go

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o bin/Agent-linux ./cmd/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/Agent-darwin ./cmd/main.go
GOOS=windows GOARCH=amd64 go build -o bin/Agent.exe ./cmd/main.go

# Run with flags
./bin/Agent -animate "problem"     # Generate animation
./bin/Agent -image photo.jpg      # OCR from image
./bin/Agent -pdf exam.pdf         # OCR from PDF
./bin/Agent -analyze resource     # Full analysis
./bin/Agent -online               # Force online auto-solve
```

### REST API

```bash
cd src/api

# Build
go build -o bin/api ./cmd/main.go

# Run
./bin/api
PORT=8080 ./bin/api

# Hot reload
air -c .air.toml
```

### Rust (Noon)

```bash
cd src/noon

# Build
cargo build --release

# Run tests
cargo test

# Check
cargo check

# Format
cargo fmt

# Lint
cargo clippy -- -D warnings
```

### Database

```bash
# Apply schema to Neon
psql $NEON_DB_URL -f sql/schema.sql

# Connect interactively
psql $NEON_DB_URL
```

## Testing

### Run Tests

```bash
cd src/agent

# All tests
go test ./...

# Single test
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

### Table-Driven Tests

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

## Code Style

### Logging: Use `log/slog`

```go
import "log/slog"

// GOOD - structured logging
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.Info("starting server", "port", 8080)
slog.Error("failed to connect", "err", err)

// NEVER use these
fmt.Println("Starting server")
print("some value")
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
// GOOD
if err != nil {
    return nil, fmt.Errorf("solve with ollama: %w", err)
}

// Check sentinel errors
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

### Config: Environment variables, no hardcoded values

```go
// GOOD
dbURL := os.Getenv("NEON_DB_URL")
if dbURL == "" {
    dbURL = "postgres://localhost/bac"
}

// BAD - never hardcode
dbURL := "postgres://localhost/bac"
```

### Cursor Rules

> **never write useless boring print statements**

Always use structured logging with `log/slog`.

## Agent Components

| Package | Purpose | Path |
|---------|---------|------|
| `solver` | Solve math problems via Ollama | `internal/solver/` |
| `memory` | Vector DB lookup/storage (pgvector) | `internal/memory/` |
| `ocr` | OCR (Tesseract, Ollama Vision) | `internal/ocr/` |
| `analyzer` | Multi-service parallel analysis | `internal/analyzer/` |
| `animation` | Noon animation generation | `internal/animation/` |
| `nlm` | NotebookLM integration | `internal/nlm/` |
| `online` | Playwright research (ChatGPT, Claude) | `internal/online/` |
| `db` | Database connection and queries | `internal/db/` |
| `kb` | Knowledge base | `internal/kb/` |

## Database

### Schema Location

- Main schema: `sql/schema.sql`
- Queries: `sql/query.sql`
- Knowledge base: `sql/knowledge_base.sql`

### sqlc - Type-Safe SQL

All database queries MUST use **sqlc** for type-safe Go code:

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate Go code from SQL
cd src/agent && sqlc generate
```

- Write SQL in `sqlc.yaml` config
- Generate Go with `sqlc generate`
- NEVER write raw database code

### Key Tables

- `questions` - Question bank with embeddings
- `subjects` - Subject categories
- `chapters` - Chapter organization
- `users` - User data
- `predictions` - Question predictions
- `animations` - Noon animations

### pgvector Search

```sql
-- Find similar questions
SELECT id, question, embedding <=> $1 AS distance
FROM questions
ORDER BY embedding <=> $1
LIMIT 5;
```

## External Services

| Service | Endpoint | Purpose | Env Variable |
|---------|----------|---------|--------------|
| Ollama | localhost:11434 | Local AI | `OLLAMA_HOST` |
| Neon DB | Cloud (env) | PostgreSQL + pgvector | `NEON_DB_URL` |
| Turso | libsql://*.turso.io | SQLite (prod cache) | `TURSO_DB_URL` |
| Garage S3 | localhost:3900 | Storage | `GARAGE_ENDPOINT` |
| NLM | Cloud API | NotebookLM | Via `nlm` CLI |

## Environment Variables

```bash
# Database
NEON_DB_URL="postgresql://user:pass@host/db"
TURSO_DB_URL="libsql://token@host/db"

# Storage
GARAGE_ENDPOINT="http://localhost:3900"
GARAGE_ACCESS_KEY="..."
GARAGE_SECRET_KEY="..."
GARAGE_BUCKET="bac-resources"

# AI
OLLAMA_HOST="http://localhost:11434"

# Cache (Turso SQLite)
TURSO_DB_URL="libsql://token@turso.io/db"
NLM_CACHE_DB_URL="${TURSO_DB_URL}"
NLM_CACHE_TTL=86400

# Redis
REDIS_URL="redis://localhost:6379"

# API
PORT=8080
JWT_SECRET="..."
```

## Common Tasks

### New Agent Command

1. Create `src/agent/internal/<module>/`
2. Implement function in new file
3. Add to `cmd/main.go`
4. Add tests in `<module>_test.go`
5. Run: `go test -v ./internal/<module>/`

### New API Endpoint

1. Handler: `src/api/internal/handlers/`
2. Service: `src/api/internal/services/`
3. Register: `src/api/main.go`

### Database Migration

1. Add SQL to `sql/schema.sql`
2. Test locally
3. Apply: `psql $NEON_URL -f sql/schema.sql`

## Documentation

### Files

| File | Purpose |
|------|---------|
| `AGENTS.md` | Agent instructions (this file) |
| `docs/AGENTS.md` | Agent daemon rules |
| `README.md` | Project overview |
| `PROBLEM.md` | Sample problem |
| `PROJECT-STATUS.md` | Status report |
| `doc/PLAN/STRATEGIC_PLAN.md` | Long-term plan |

### Alias Policy

- **`AGENTS.md` (root)** is canonical for agent instructions
- **`docs/AGENTS.md`** contains agent-specific daemon rules (scope-specific)
- Create symlink: `docs/AGENTS.md` → `AGENTS.md` or keep separate with clear scope

## Development Workflow

### Setup

```bash
# Install dependencies
cd src/agent && go mod download
cd src/noon && cargo fetch

# Set environment
cp .env.example .env
# Edit .env with credentials
```

### Typical Cycle

```bash
# 1. Create new change
jj new -m "Add new feature"

# 2. Make changes
# - Edit code
# - Run tests: go test -v ./...

# 3. Format and lint
go fmt ./...
go vet ./...
golangci-lint run

# 4. Commit
jj describe -m "Add new feature"
jj commit

# 5. Test suite
go test -race -cover ./...

# 6. Push to remote
jj push
```

### Code Review Checklist

- [ ] Tests pass (`go test ./...`)
- [ ] Code formatted (`go fmt`)
- [ ] No vet warnings (`go vet`)
- [ ] Race clean (`go test -race`)
- [ ] Logging uses `slog`
- [ ] Errors wrapped
- [ ] Env vars for config

## OpenSpec Workflow

### Creating a Change

```bash
# Create new change
mkdir -p openspec/changes/<change-name>

# Add artifacts
openspec/changes/<change-name>/
├── .openspec.yaml    # metadata
├── proposal.md      # what & why
├── design.md       # how
└── tasks.md        # implementation
```

### Running Tasks

```bash
# List changes
openspec list

# Check status
openspec status --change <name>

# Apply (implement tasks)
openspec apply <name>

# Archive when complete
openspec archive <name>
```

### Skills

Use skill tool for complex tasks:

```bash
# Load openspec-apply-change skill
skill openspec-apply-change

# Load openspec-propose skill  
skill openspec-propose
```

---

*Last Updated: 2026-03-01*
*Canonical: AGENTS.md (root)*
