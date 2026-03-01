# Skill: PostgreSQL + pgvector

## Purpose

PostgreSQL with pgvector extension for vector similarity search in BAC Unified.

## When to use

- Working with the main database (`src/agent/internal/db/`)
- Implementing semantic search for questions
- Storing question embeddings

## Database Schema

### Key Tables

```sql
-- Questions with embeddings
CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question TEXT NOT NULL,
    subject TEXT NOT NULL,
    chapter TEXT,
    embedding vector(1536),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Subjects
CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    code TEXT NOT NULL
);

-- Chapters
CREATE TABLE chapters (
    id SERIAL PRIMARY KEY,
    subject_id INTEGER REFERENCES subjects(id),
    name TEXT NOT NULL,
    order_index INTEGER
);
```

## pgvector Operations

### Create Extension

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

### Similarity Search

```sql
-- Find similar questions (cosine distance)
SELECT id, question, embedding <=> $1 AS distance
FROM questions
ORDER BY embedding <=> $1
LIMIT 5;

-- Find similar (L2 distance)
SELECT id, question, embedding <-> $1 AS distance
FROM questions
ORDER BY embedding <-> $1
LIMIT 5;

-- Find similar (inner product)
SELECT id, question, embedding <#> $1 AS distance
FROM questions
ORDER BY embedding <#> $1
LIMIT 5;
```

### Indexes

```sql
-- HNSW index for faster search
CREATE INDEX ON questions USING hnsw (embedding vector_cosine_ops);

-- Or ivfflat for smaller datasets
CREATE INDEX ON questions USING ivfflat (embedding vector_cosine_ops, lists = 100);
```

## Go Integration

### Using sqlc (REQUIRED)

```bash
# Generate from SQL
sqlc generate
```

### Query Example

```go
// From src/agent/internal/db/query.sql.go
func (q *Queries) FindSimilarQuestions(ctx context.Context, arg FindSimilarQuestionsParams) ([]Question, error) {
    rows, err := q.db.QueryContext(ctx, findSimilarQuestions, arg.Embedding, arg.Limit)
    // ...
}
```

## Environment Variables

| Variable | Purpose |
|----------|---------|
| `NEON_DB_URL` | PostgreSQL connection string |
| `PGHOST` | PostgreSQL host |
| `PGPORT` | PostgreSQL port (5432) |
| `PGUSER` | PostgreSQL user |
| `PGPASSWORD` | PostgreSQL password |
| `PGDATABASE` | Database name |

## Connection (Go)

```go
import (
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Connection pool
pool, err := pgxpool.New(ctx, os.Getenv("NEON_DB_URL"))

// Single connection
conn, err := pgx.Connect(ctx, os.Getenv("NEON_DB_URL"))
```

## Neon (Cloud PostgreSQL)

- Managed PostgreSQL with pgvector
- Connection pooling
- Branching for development

```bash
# Install Neon CLI
npm i -g neondb-cli

# Create branch
neon branches create --name dev

# Get connection string
neon projects get-connection-string
```

## Anti-Patterns

- ❌ Using raw SQL - use sqlc
- ❌ Not using connection pooling
- ❌ Missing vector indexes (slow queries)
- ❌ Storing embeddings without normalization
