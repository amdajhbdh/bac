# BAC Unified - Agent Instructions

## Project Purpose

**BAC Unified** is a personal learning system for Moroccan BAC exam preparation. It aggregates educational resources, provides AI-powered tutoring via local Ollama models, performs OCR on study materials using Tesseract/Surya, and offers spaced repetition practice—all working offline.

| Component | Language | Path |
|-----------|----------|------|
| Agent CLI | **Go** | `src/agent/` |
| REST API | Go (Gin) | `src/api/` |
| Gateway | Rust (Axum) | `src/gateway/` |
| OCR Service | Rust (Tonic) | `src/ocr-service/` |
| AI-OCR Pipeline | **Go** | `src/ai-ocr/` |
| Web Frontend | React (Bun) | `src/web/` |
| Daemons | Nushell | `daemons/` |
| CLI Orchestration | Nushell | `lib/bac.nu` |
| Cloudflare Worker | TypeScript | `src/cloudflare-worker/` |
| Study Vault | Obsidian | `resources/notes/` |

The `src/api` module imports `src/agent` via a `replace` directive in `go.mod`. Both modules use Go 1.25.

**Database Layer:** PostgreSQL + pgvector (main), Turso (SQLite, NLM cache), Redis (hot cache)

---

## Build / Lint / Test Commands

### Go (src/agent, src/api)

```bash
# Build all Go modules
cd src/agent && go build -o ../../bin/agent .
cd src/api && go build -o ../../bin/api .

# Lint
golangci-lint run ./...

# Run all tests
go test -v ./...

# Run specific test
go test -v -run TestName ./...

# Run tests with coverage
go test -v -cover ./...

# Generate SQLC code (after schema changes)
sqlc generate
```

### Rust (src/gateway, src/ocr-service)

```bash
# Build
cd src/gateway && cargo build --release
cd src/ocr-service && cargo build --release

# Lint
cargo clippy -- -D warnings

# Run tests
cargo test

# Run specific test
cargo test test_name

# Format
cargo fmt
```

### TypeScript / React (src/web, src/cloudflare-worker)

```bash
# Install dependencies
bun install

# Build
bun run build

# Lint
bun run lint

# Type check
bun run typecheck

# Run tests
bun test

# Run specific test
bun test test-name
```

### Python Scripts

```bash
# Lint with ruff
ruff check .
ruff check path/to/file.py

# Format
ruff format .

# Type check (optional)
mypy .
```

### Nushell Scripts

```bash
# Check syntax
nu -c 'source script.nu'

# Run with debug
nu -c 'source script.nu; $nu.scope'
```

---

## Code Style Guidelines

### Go

**Imports** - Group in order, no duplicates:
```go
import (
    "context"
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"

    "bac/src/agent/internal/config"
)
```

**Naming:**
- Types/Interfaces: `PascalCase` (e.g., `UserService`, `Handler`)
- Variables/Functions: `camelCase` (e.g., `getUser`, `processData`)
- Constants: `PascalCase` or `SCREAMING_SNAKE_CASE`
- Packages: `snake_case` (e.g., `internal`, `config`)
- Files: `snake_case.go` (e.g., `user_service.go`)

**Error Handling:**
```go
// GOOD - Wrap with context
if err != nil {
    return fmt.Errorf("fetching user %d: %w", id, err)
}

// GOOD - Handle and return
if err := doSomething(); err != nil {
    return fmt.Errorf("doing X: %w", err)
}
```

**Logging:**
```go
// GOOD - Use slog
log.Info("processing request", "id", requestID)

// BAD - Never use fmt.Print*
fmt.Println("debug info")  // FORBIDDEN
```

**Type Safety:**
- Avoid `interface{}` / `any` where possible
- Use `context.Context` as first parameter
- Return concrete types where possible

**Testing:**
```go
func TestUserService(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "test", false},
        {"empty input", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := process(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("process() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Rust

**Naming:**
- Variables/Functions: `snake_case`
- Types: `PascalCase`
- Constants: `SCREAMING_SNAKE_CASE`

**Error Handling:**
```rust
// Use Result and anyhow/thiserror
fn process() -> Result<Output, MyError> {
    let data = something().context("failed to do X")?;
    Ok(data)
}
```

**Logging:**
```rust
use tracing::{info, error};

info!("processing request {}", id);
error!("failed: {}", err);
```

**Testing:**
```rust
#[cfg(test)]
mod tests {
    #[test]
    fn test_something() {
        assert_eq!(2 + 2, 4);
    }
}
```

### TypeScript / JavaScript

**Naming:**
- Variables/Functions: `camelCase`
- Components/Classes: `PascalCase`
- Files: `kebab-case.ts` or `camelCase.ts`

**Imports:**
```typescript
// Group: external → internal
import React from 'react';
import { useState, useEffect } from 'react';

import { api } from '@/lib/api';
import { User } from '@/types';
```

**Type Safety:**
- Always use TypeScript types
- Avoid `any`, use `unknown` if needed

### Python

**Naming:**
- Functions/Variables: `snake_case`
- Classes: `PascalCase`
- Constants: `SCREAMING_SNAKE_CASE`

**Imports:**
```python
# Group: stdlib → external → internal
import os
import sys
from pathlib import Path

import requests
from dotenv import load_dotenv

from . import local_module
```

**Type Hints:**
```python
def process_user(user_id: int) -> User | None:
    ...
```

---

## Mandatory Rules

| Rule | Requirement |
|------|-------------|
| **VCS** | Use **jj** (Jujutsu). Never git directly |
| **Database** | PostgreSQL + pgvector for main, Turso for cache |
| **DB Tools** | Use **sqlc** for all queries |
| **Search** | Hybrid = Elasticsearch + pgvector |
| **Storage** | Garage S3 or MinIO only (S3-compatible) |
| **OCR** | Tesseract + Surya only. Never paid APIs |
| **Auth** | JWT + bcrypt. All secrets via env vars |
| **Logging** | Use log/slog (Go), tracing (Rust). Never fmt.Print* |
| **Bundling** | WRAPX (<100MB), ARX (≥100MB) with metadata |
| **Monitoring** | Prometheus + Grafana |

---

## Before Working

1. Check `specs/` — update task status on completion
2. Read `doc/rfc/*.md` for architecture decisions
3. Read `doc/OPEN_SOURCE_PROJECTS.md` for allowed tool list
4. Check `.cursorrules` for additional strict rules

---

## Vault (resources/notes/)

The study vault is a standalone git repository:
- **Remote**: https://github.com/amdajhbdh/notes.git
- **Sync**: Git-based (not Obsidian Sync)
- **Phone**: Clone from GitHub, use git pull/push workflow

### Quick Vault Commands
```bash
cd resources/notes

# Sync from remote
./sync-vault.sh pull

# Push changes
./sync-vault.sh push

# Run AI-OCR on a PDF
../../bin/ai-ocr extract path/to/file.pdf

# Batch OCR
../../bin/ai-ocr batch 03-Resources/
```

---

## Multi-Cloud Infrastructure

| Provider | Service | Purpose |
|----------|---------|---------|
| **CloudShell** | gcloud CLI | Development, OCR processing |
| **Cloudflare** | Workers + Pages | API endpoints, agents |
| **Render** | Go/Rust hosting | API service hosting |
| **Fly.io** | Edge computing | Low-latency requests |
| **Doppler** | Secrets management | Environment variables |

---

## AI Inference Priority

| Priority | Model | Use Case |
|----------|-------|----------|
| 1 | Ollama (llama3.2:3b) | Local, fast inference |
| 2 | Cloud APIs | kimi, deepseek, minimax |
| 3 | NLM CLI/MCP | NotebookLM with caching |

---

RULES:
    ALWAYS use context-mode mcp
    NEVER use fmt.Print* - use logging instead
    NEVER hardcode secrets - use environment variables
*Last Updated: 2026-03-14*
