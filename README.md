# BAC Unified - BAC Exam Preparation Platform

## Overview

| Component | Language | Path |
|-----------|----------|------|
| Agent CLI | **Go** | `src/agent/` |
| REST API | Go (Gin) | `src/api/` |
| Animations | Rust (Noon) | `src/noon/` |
| Daemons | Nushell | `daemons/` |
| CLI Orchestration | Nushell | `lib/bac.nu` |
| Cloudflare Worker | TypeScript | `src/cloudflare-worker/` |
| Cloudflare Pages | HTML/JS | `src/cloudflare-pages/` |

## Architecture

```
Nushell CLI (lib/bac.nu)
    └── Go Agent CLI (src/agent/cmd/main.go)
            ├── solver/       Ollama → Cloud APIs → NLM → Online research
            ├── memory/       pgvector similarity search (RAG)
            ├── nlm/          NotebookLM + Turso cache
            ├── ocr/          Tesseract + Surya
            ├── search/       Hybrid: Elasticsearch + pgvector
            ├── animation/    Noon DSL code generation
            ├── bundle/       WRAPX + ARX
            ├── db/           sqlc PostgreSQL queries
            └── cache/        Redis

Database: PostgreSQL + pgvector | Turso (SQLite) | Redis
```

## AI Inference Priority

1. Ollama (llama3.2:3b, local, fast)
2. Cloud APIs (kimi/deepseek/minimax)
3. NLM CLI/MCP (NotebookLM with caching)

## Quick Start

```bash
# Install dependencies
just install

# Build everything
just quick-build

# Start services
just services-up

# Run tests
just test

# Start development
just dev
```

## Commands

Run `just --list` to see all available tasks:

```bash
just build          # Build agent binary
just build-tools    # Build CLI tools
just quick-build    # Build with checks
just test           # Run tests
just test-race      # Run with race detector
just fmt            # Format code
just vet            # Run go vet
just services-up    # Start postgres/redis
just services-down   # Stop services
just dev            # Start dev server
just api            # Start REST API
just server         # Start agent server
just health         # Health check
just load-data      # Load sample data
just db-shell       # PostgreSQL shell
just redis-shell    # Redis shell
just clean          # Clean build artifacts
just deploy         # Deploy to production
just fix            # Auto-fix build issues
```

## Development

```bash
# Go (Agent)
cd src/agent
go build -o bin/Agent ./cmd/main.go
./bin/Agent -h

# Go (REST API)
cd src/api
go run main.go

# Rust (Noon)
cd src/noon
cargo build --release
```

## Testing

```bash
cd src/agent

# All tests
go test ./...

# Single test
go test -v -run TestSolve ./internal/solver/

# With coverage
go test -cover ./...

# Benchmark
go test -bench=BenchmarkMemoryLookup ./internal/memory/
```

## Database

- **Main DB**: PostgreSQL + pgvector (schema: `sql/schema.sql`)
- **NLM cache**: Turso (SQLite) at `src/agent/internal/nlm/cache/`
- **Hot cache**: Redis

Use **sqlc** for all queries - run `sqlc generate` after schema changes.

## Mandatory Rules

- **VCS**: Use **jj (Jujutsu)** - `jj new`, `jj commit`, `jj push`
- **Storage**: Garage S3 or MinIO only
- **OCR**: Tesseract + Surya only (no paid APIs)
- **Auth**: JWT + bcrypt, env vars for secrets
- **Search**: Hybrid = Elasticsearch + pgvector

## Code Style

- Use `log/slog` for logging
- Group imports: stdlib, external, internal
- Wrap errors with context
- Env vars for all config

---

**Built with ❤️ for Mauritanian BAC students**
