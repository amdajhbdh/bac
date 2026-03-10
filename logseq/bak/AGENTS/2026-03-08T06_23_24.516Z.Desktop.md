# BAC Unified - Agent Instructions

## Project Overview

| Component | Language | Path |
|-----------|----------|------|
| Agent CLI | **Go** | `src/agent/` |
| REST API | Go (Gin) | `src/api/` |
| Animations | Rust (Noon) | `src/noon/` |
| Daemons | Nushell | `daemons/` |
| CLI Orchestration | Nushell | `lib/bac.nu` |
| Cloudflare Worker | TypeScript | `src/cloudflare-worker/` |
| Cloudflare Pages | HTML/JS | `src/cloudflare-pages/` |
 
 The `src/api` module imports `src/agent` via a `replace` directive in `go.mod`. Both modules are `github.com/bac-unified/agent` and `github.com/bac-unified/api` respectively, using Go 1.25.
 
 ## Architecture
 
 ```
 Nushell CLI (lib/bac.nu)
     └── Go Agent CLI (src/agent/cmd/main.go)
             ├── solver/       Ollama (local) → Cloud APIs → NLM → Online research
             ├── memory/       pgvector similarity search (RAG)
             ├── nlm/          NotebookLM integration + Turso (SQLite) cache
             ├── ocr/          Tesseract + Surya
             ├── search/       Hybrid: Elasticsearch (full-text) + pgvector (vector)
             ├── animation/    Noon DSL code generation
             ├── bundle/       WRAPX (<100MB) + ARX (≥100MB)
             ├── aggregator/   Khan Academy, OpenStax, MIT OCW, OER
             ├── db/           sqlc-generated PostgreSQL queries
             └── cache/        Redis
 
 Background Daemons (daemons/*.nu)
     youtube-daemon.nu | webcrawler-daemon.nu | audio-daemon.nu
     nlm-research-daemon.nu | kb-watcher.nu
 
 Database Layer
     PostgreSQL + pgvector  → main DB (questions, users, gamification, audit)
     Turso (SQLite)         → NLM response cache (src/agent/internal/nlm/cache/)
     Redis                  → hot cache
 ```
 
 ### AI Inference Priority
 ```
 Ollama (llama3.2:3b, local, fast)
     → Cloud APIs (kimi/deepseek/minimax/ministral, medium)
     → NLM CLI/MCP (NotebookLM with caching)
 ```
 
 ### Hybrid Search (Reciprocal Rank Fusion)
 1. Parse query, extract filters (subject, chapter, difficulty, language, source)
 2. Run in parallel: Elasticsearch full-text (French analyzer) + pgvector cosine similarity (1536-dim)
 3. Fuse results with RRF normalization
 4. Return with facets + suggestions
 
 ## Before Working
 
 1. Check `openspec/changes/full-bac-unified/tasks.md` for current task status — update `- [ ]` → `- [x]` on completion
 2. Read `doc/rfc/*.md` for architecture decisions
 3. Read `doc/OPEN_SOURCE_PROJECTS.md` for the allowed tool list
 
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

 ### Go (REST API)
 
 ```bash
 cd src/api
 go run main.go
 # Swagger docs available at /swagger/*any
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

Test files: `*_test.go`
Unit tests: `TestFunctionName`
Benchmarks: `BenchmarkFunctionName`

 ## Database
 
 - **Main DB**: PostgreSQL + pgvector. Schema: `sql/schema.sql`
 - **NLM cache**: Turso (SQLite) at `src/agent/internal/nlm/cache/` — has its own `schema.sql` and sqlc queries
 - **Hot cache**: Redis
 - Use **sqlc** for all SQL queries — run `sqlc generate` after any schema change
 - Never use raw SQL strings or an ORM
 
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
 
 ### Key Schema Tables
 - `questions` — core table: `embedding vector(1536)`, `content_hash` (MD5 dedup), `solution_steps JSONB`, `ai_concepts JSONB`
 - `users` — gamification fields: `points`, `level`, `streak_days`, `accuracy_rate`
 - `predictions` — ML exam topic predictions with `confidence_score`
 - `animations` — Noon DSL code + rendered video paths
 - `audit_log` — blockchain-style chain (`transaction_hash`, `previous_hash`)
 - `storage_objects` — Garage S3 object metadata
 
 ## Mandatory Rules
 
 ### VCS
 - Use **jj (Jujutsu)**. NEVER use `git` directly.
 - `jj new`, `jj commit`, `jj push`, `jj log`
 
 ### Storage
 - **Garage S3** or MinIO only (S3-compatible API). Never AWS S3 direct or GCS.
 
 ### Bundling
 - **WRAPX** for bundles <100MB, **ARX** for ≥100MB. Always include Dublin Core + LOM metadata.
 
 ### OCR
 - Tesseract + Surya only. NEVER use paid APIs (Google Vision, AWS Textract).
 
 ### Auth
 - JWT + bcrypt. All secrets via environment variables.
 
 ### Search
 - Hybrid = Elasticsearch (full-text) + pgvector (vectors). Never one alone in production.
 
 ### Monitoring
 - Prometheus + Grafana. Track: latency, errors, cache hit rate.
 
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

## Code Review Checklist

- [ ] Check `openspec/changes/full-bac-unified/tasks.md` before and after
- [ ] Tests pass (`go test ./...`)
- [ ] Code formatted (`go fmt`)
- [ ] No vet warnings (`go vet`)
- [ ] Race clean (`go test -race`)
- [ ] Logging uses `slog`
- [ ] Errors wrapped
- [ ] Env vars for config
- [ ] `sqlc generate` run if schema changed
- [ ] VCS: `jj commit` not `git commit`
- [ ] Check openspec/changes/full-bac-unified/tasks.md before working

---

*Last Updated: 2026-03-01*
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
 
## Code Review Checklist

[ ] Check `openspec/changes/full-bac-unified/tasks.md` before and after
[ ] Tests pass (`go test ./...`)
[ ] Code formatted (`go fmt`)
[ ] No vet warnings (`go vet`)
[ ] Race clean (`go test -race`)
[ ] Logging uses `slog`
[ ] Errors wrapped
[ ] Env vars for config
[ ] `sqlc generate` run if schema changed
[ ] VCS: `jj commit` not `git commit`
[ ] Check openspec/changes/full-bac-unified/tasks.md before working

---

*Last Updated: 2026-03-06*
