# Skill: Turso (SQLite) for Production Cache

## Purpose

Turso distributed SQLite for production caching in BAC Unified.

## When to use

- Working with NLM cache (`src/agent/internal/nlm/cache/`)
- Implementing edge caching
- Low-latency database operations

## Why Turso?

- Edge-distributed SQLite
- HTTP-based (low latency)
- Embedded replicas support
- SQLite-compatible

## Connection

```go
import (
    "database/sql"
    _ "github.com/tursodatabase/go-libsql"
)

// Local file
db, _ := sql.Open("libsql", "file:local.db")

// Turso remote
db, _ := sql.Open("libsql", "libsql://my-db.turso.io")

// Turso with auth
db, _ := sql.Open("libsql", "libsql://token@my-db.turso.io")
```

## URL Schemes

| Scheme | Example | Use |
|--------|---------|-----|
| `libsql://` | `libsql://db.turso.io` | Turso cloud |
| `file://` | `file:/path/to/db` | Local file |
| `http://` | `http://localhost:8080` | HTTP replica |
| `https://` | `https://localhost:8080` | HTTPS replica |

## Schema (NLM Cache)

```sql
CREATE TABLE IF NOT EXISTS nlm_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query_hash TEXT NOT NULL,
    query_text TEXT NOT NULL,
    subject TEXT NOT NULL,
    topics TEXT DEFAULT '[]',
    notebook_id TEXT NOT NULL,
    s3_key TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    expires_at TEXT NOT NULL,
    access_count INTEGER DEFAULT 1,
    last_accessed TEXT DEFAULT (datetime('now')),
    UNIQUE(query_hash, subject)
);

CREATE INDEX idx_nlm_cache_hash ON nlm_cache(query_hash);
CREATE INDEX idx_nlm_cache_subject ON nlm_cache(subject);
CREATE INDEX idx_nlm_cache_expires ON nlm_cache(expires_at);
```

## Using sqlc

### Configuration

```yaml
version: "2"
sql:
  - schema: "schema.sql"
    queries: "queries.sql"
    engine: "sqlite"
    database:
      uri: "file:nlm_cache.db"
    gen:
      go:
        package: "cache"
        out: "."
        sql_package: "database/sql"
```

### Query Annotations

```sql
-- name: GetCachedResponse :one
SELECT * FROM nlm_cache 
WHERE query_hash = ?1 AND subject = ?2 AND expires_at > datetime('now');

-- name: SetCacheEntry :exec
INSERT INTO nlm_cache (query_hash, query_text, subject, topics, notebook_id, s3_key, expires_at)
VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7)
ON CONFLICT(query_hash, subject) DO UPDATE SET
    access_count = access_count + 1;
```

## Environment Variables

| Variable | Purpose |
|----------|---------|
| `TURSO_DB_URL` | Turso connection string |
| `NLM_CACHE_DB_URL` | NLM cache DB URL |

## CLI Commands

```bash
# Install Turso CLI
brew install tursodatabase/tap/turso

# Login
turso auth login

# Create database
turso db create bac-cache

# Get connection string
turso db show bac-cache --url

# Open shell
turso db shell bac-cache
```

## Anti-Patterns

- ❌ Using Turso for primary data (use PostgreSQL)
- ❌ Not using sqlc for queries
- ❌ Blocking on network operations
- ❌ Not handling connection errors
