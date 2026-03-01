# Skill: sqlc Generation

## Purpose

Type-safe SQL queries using sqlc for BAC Unified - NEVER write raw SQL in Go.

## When to use

- Any database operation
- Creating new queries
- Modifying schema

## MANDATORY: Always Use sqlc

❌ **FORBIDDEN**: `db.Query()`, `db.Exec()` with raw SQL
✅ **REQUIRED**: Generated Go code from sqlc

## Setup

### 1. Install sqlc

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### 2. Create sqlc.yaml

```yaml
version: "2"
sql:
  - schema: "schema.sql"
    queries: "queries.sql"
    engine: "postgresql"  # or "sqlite"
    database:
      uri: "${DATABASE_URL}"
    gen:
      go:
        package: "db"
        out: "."
        sql_package: "pgx/v5"  # or "database/sql"
        emit_json_tags: true
        emit_interface: true
```

### 3. Create Schema

```sql
-- schema.sql
CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question TEXT NOT NULL,
    subject TEXT NOT NULL,
    embedding vector(1536),
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 4. Create Queries

```sql
-- queries.sql

-- name: GetQuestion :one
SELECT * FROM questions WHERE id = $1;

-- name: ListQuestions :many
SELECT * FROM questions ORDER BY created_at DESC LIMIT $1;

-- name: CreateQuestion :one
INSERT INTO questions (question, subject, embedding)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM questions WHERE id = $1;
```

### 5. Generate

```bash
sqlc generate
```

## Query Types

| Annotation | Returns |
|------------|---------|
| `:one` | Single row |
| `:many` | Multiple rows |
| `:exec` | No return |
| `:execresult` | Result with rows affected |

## Parameters

### PostgreSQL

```sql
-- name: FindBySubject :many
SELECT * FROM questions WHERE subject = $1;

-- name: Search :many
SELECT * FROM questions WHERE question ILIKE '%' || $1 || '%';
```

### SQLite

```sql
-- name: FindBySubject :many
SELECT * FROM questions WHERE subject = ?1;

-- name: Search :many
SELECT * FROM questions WHERE question LIKE '%' || ?1 || '%';
```

## Using Generated Code

```go
import "github.com/bac-unified/agent/internal/db"

func GetQuestion(ctx context.Context, id string) (*db.Question, error) {
    q := db.New(conn)
    return q.GetQuestion(ctx, id)
}

func ListQuestions(ctx context.Context, limit int32) ([]db.Question, error) {
    q := db.New(conn)
    return q.ListQuestions(ctx, limit)
}

func CreateQuestion(ctx context.Context, arg db.CreateQuestionParams) (*db.Question, error) {
    q := db.New(conn)
    return q.CreateQuestion(ctx, arg)
}
```

## Project Structure

```
src/agent/internal/nlm/cache/
├── schema.sql
├── queries.sql
├── sqlc.yaml
├── db.go           # Generated
├── models.go        # Generated
└── query.sql.go     # Generated
```

## Common Issues

### "failed to resolve schema"

```bash
# Check DATABASE_URL is set
export DATABASE_URL="postgresql://..."
sqlc generate
```

### "query has no parameter"

```sql
-- Wrong
SELECT * FROM questions WHERE id = $1;

-- Right (add parameter)
-- name: GetQuestion :one
SELECT * FROM questions WHERE id = $1;
```

## Anti-Patterns

- ❌ Writing raw SQL in Go
- ❌ Using database/sql directly without sqlc
- ❌ Not running `sqlc generate` after schema changes
- ❌ Hardcoding connection strings
