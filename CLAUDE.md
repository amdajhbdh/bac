# BAC Unified - Claude/AI Rules

This file contains restrictive rules for all AI coding agents.

## Quick Reference

| Category | Allowed | Forbidden |
|----------|---------|------------|
| VCS | jj (Jujutsu) | git directly |
| Database | PostgreSQL + pgvector, Turso (SQLite) | MongoDB, MySQL |
| DB Tools | sqlc | raw sql, ORM |
| Search | Elasticsearch + pgvector | Algolia, Typesense alone |
| Storage | Garage S3, MinIO | AWS S3 direct, GCS |
| OCR | Tesseract, Surya | Google Vision, AWS Textract |
| Bundles | WRAPX, ARX | tar/zip without metadata |
| Logging | log/slog | fmt.Print*, print() |
| Auth | JWT + bcrypt | Basic auth, hardcoded secrets |

## Critical Rules

### 0. VERSION CONTROL
- Use **jj** (Jujutsu) - NEVER use git directly
- `jj new`, `jj commit`, `jj push`, `jj log`

### 0.1. PROJECT PROGRESS
- ALWAYS check `openspec/changes/full-bac-unified/tasks.md` before working
- Update completed tasks: `- [ ]` → `- [x]`
- Report progress in responses: "X/Y tasks complete"

### 1. DATABASE
- PostgreSQL + pgvector is MANDATORY for main database
- Turso (SQLite) is REQUIRED for production cache
- Use **sqlc** for all database queries (generate with `sqlc generate`)

### 2. SEARCH
```
Hybrid Search = Elasticsearch (full-text) + pgvector (vectors)
```
- Never use one without the other for production

### 3. LOGGING
- ALWAYS use `log/slog`
- NEVER use `fmt.Print*` or `print()`

### 4. ERRORS
```go
if err != nil {
    return nil, fmt.Errorf("doing something: %w", err)
}
```

### 5. CONFIG
```go
// GOOD
dbURL := os.Getenv("NEON_DB_URL")

// BAD  
dbURL := "postgres://localhost/bac"
```

### 6. IMPORTS
```go
import (
    // Standard library
    "context"
    "log/slog"
    
    // External packages
    "github.com/gin-gonic/gin"
    
    // Internal packages
    "github.com/bac-unified/agent/internal/solver"
)
```

### 7. TESTING
```bash
# Single test
go test -v -run TestName ./...

# With race detector
go test -race ./...
```

## Required Reading

Before coding, read these files:
1. `AGENTS.md` - Canonical rules
2. `AI_RULES.md` - Comprehensive restrictions
3. `doc/rfc/*.md` - Architecture decisions
4. `doc/OPEN_SOURCE_PROJECTS.md` - Allowed tools

## Integration Decisions (Locked)

| Category | Decision |
|----------|----------|
| Database | PostgreSQL + pgvector, Turso |
| Search | Elasticsearch + pgvector |
| Storage | Garage S3 |
| Cache | Turso (SQLite) + Redis |
| OCR | Tesseract + Surya |
| Bundling | WRAPX + ARX |

## Violations

Violations will be corrected without warning.

---

*Last Updated: 2026-03-01*
