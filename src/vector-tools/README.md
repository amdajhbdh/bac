# Vector Tools

pgvector operations service with HNSW index support for semantic search.

## Features

- **Semantic Search**: Fast (<10ms) similarity search using HNSW index
- **Vector Insert**: Store embeddings with metadata
- **Batch Operations**: Bulk insert for efficient data loading
- **Metadata Filtering**: Filter by category, tags, and custom metadata
- **Index Management**: Rebuild HNSW index with configurable parameters

## Service Port

Runs on port `3002`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/search` | Semantic search |
| POST | `/insert` | Insert single vector |
| POST | `/batch` | Batch insert |
| DELETE | `/delete/:id` | Delete by ID |
| POST | `/rebuild` | Rebuild HNSW index |

## Request/Response Examples

### Search

```bash
curl -X POST http://localhost:3002/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": [0.1, 0.2, 0.3, ...],
    "top_k": 10,
    "filters": {
      "category": "tech"
    }
  }'
```

Response:
```json
{
  "results": [
    {
      "id": 1,
      "content": "Document text...",
      "metadata": {"source": "example"},
      "similarity": 0.95
    }
  ],
  "query_time_ms": 5
}
```

### Insert

```bash
curl -X POST http://localhost:3002/insert \
  -H "Content-Type: application/json" \
  -d '{
    "embedding": [0.1, 0.2, 0.3, ...],
    "content": "Document content",
    "metadata": {"category": "tech", "tags": ["rust", "pgvector"]}
  }'
```

### Batch Insert

```bash
curl -X POST http://localhost:3002/batch \
  -H "Content-Type: application/json" \
  -d '{
    "records": [
      {"embedding": [0.1, 0.2], "content": "Doc 1"},
      {"embedding": [0.3, 0.4], "content": "Doc 2"}
    ]
  }'
```

### Delete

```bash
curl -X DELETE http://localhost:3002/delete/123
```

### Rebuild Index

```bash
curl -X POST http://localhost:3002/rebuild \
  -H "Content-Type: application/json" \
  -d '{"m": 16, "ef_construction": 64}'
```

## Configuration

Set via environment variable:
```bash
DATABASE_URL=postgresql://user:pass@host:5432/dbname?sslmode=require
```

## Running

```bash
# Build
cargo build -p vector-tools

# Run
cargo run -p vector-tools

# Or with custom port
PORT=3002 cargo run -p vector-tools
```

## Architecture

```
src/vector-tools/
├── Cargo.toml
├── schema.sql           # Database schema
├── README.md
└── src/
    ├── lib.rs           # Library exports
    ├── main.rs          # Binary entrypoint
    ├── client.rs        # PostgreSQL/pgvector client
    ├── service.rs       # HTTP service (Axum)
    ├── search.rs        # Semantic search
    ├── insert.rs        # Single insert
    ├── batch.rs         # Batch operations
    ├── delete.rs        # Delete + rebuild
    ├── models.rs        # Request/response types
    └── error.rs         # Error handling
```

## HNSW Index Parameters

- **m**: Maximum connections per node (4-64, default: 16)
- **ef_construction**: Build-time candidate list size (16-512, default: 64)

Higher values = better recall but more memory and slower builds.
