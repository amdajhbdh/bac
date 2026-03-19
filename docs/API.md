# BAC API Documentation

REST API endpoints for all BAC services.

## Gateway (Port 8080)

Main API gateway routing requests to all services.

### Health Check

```bash
GET /health
```

Response:
```json
{
  "status": "healthy",
  "services": {
    "gemini": "up",
    "vector": "up",
    "vault": "up"
  }
}
```

---

## gemini-tools (Port 3001)

AI content generation and analysis.

### Health

```bash
GET /health
```

### Generate Content

```bash
POST /generate
Content-Type: application/json

{
  "prompt": "Write a summary of thermodynamics",
  "subject": "physics",
  "format": "markdown"
}
```

Response:
```json
{
  "content": "# Thermodynamics\n\nThermodynamics is...",
  "tokens_used": 512
}
```

### Extract Content

```bash
POST /extract
Content-Type: application/json

{
  "text": "Raw text to analyze...",
  "action": "summarize"
}
```

### Embed Text

```bash
POST /embed
Content-Type: application/json

{
  "text": "Text to embed",
  "model": "gemini-embedding-2-preview"
}
```

Response:
```json
{
  "embedding": [0.1, 0.2, 0.3, ...],
  "dimensions": 768
}
```

---

## vector-tools (Port 3002)

Semantic search using pgvector with HNSW index.

### Health

```bash
GET /health
```

### Semantic Search

```bash
POST /search
Content-Type: application/json

{
  "query": [0.1, 0.2, 0.3, ...],
  "top_k": 10,
  "filters": {
    "category": "physics"
  }
}
```

Response:
```json
{
  "results": [
    {
      "id": 1,
      "content": "Document text...",
      "metadata": {"source": "notes"},
      "similarity": 0.95
    }
  ],
  "query_time_ms": 5
}
```

### Insert Vector

```bash
POST /insert
Content-Type: application/json

{
  "embedding": [0.1, 0.2, 0.3, ...],
  "content": "Document content",
  "metadata": {
    "category": "physics",
    "tags": ["thermodynamics"]
  }
}
```

### Batch Insert

```bash
POST /batch
Content-Type: application/json

{
  "records": [
    {"embedding": [0.1, 0.2], "content": "Doc 1"},
    {"embedding": [0.3, 0.4], "content": "Doc 2"}
  ]
}
```

### Delete Vector

```bash
DELETE /delete/{id}
```

### Rebuild Index

```bash
POST /rebuild
Content-Type: application/json

{
  "m": 16,
  "ef_construction": 64
}
```

---

## vault-tools (Port 3003)

Obsidian vault operations.

### Health

```bash
GET /health
```

### Read Note

```bash
POST /read
Content-Type: application/json

{
  "path": "Physics/Thermodynamics.md"
}
```

### Write Note

```bash
POST /write
Content-Type: application/json

{
  "path": "Physics/Thermodynamics.md",
  "content": "# Thermodynamics\n\n...",
  "frontmatter": {
    "tags": ["physics"],
    "created": "2024-01-01"
  }
}
```

### Search Vault

```bash
POST /search
Content-Type: application/json

{
  "query": "thermodynamics",
  "limit": 10
}
```

### Get MOC (Map of Content)

```bash
POST /moc
Content-Type: application/json

{
  "folder": "Physics"
}
```

### Create Link

```bash
POST /link
Content-Type: application/json

{
  "source": "Physics/Thermo.md",
  "target": "Physics/Entropy.md",
  "type": "bidirectional"
}
```

---

## cloud-tools (Port 3004)

Cloud Shell integration.

### Health

```bash
GET /health
```

### SSH Command

```bash
POST /ssh
Content-Type: application/json

{
  "command": "ls -la",
  "timeout": 30
}
```

### Upload File

```bash
POST /upload
Content-Type: multipart/form-data

file: <binary>
destination: /remote/path
```

### Download File

```bash
POST /download
Content-Type: application/json

{
  "remote_path": "/remote/file.txt"
}
```

### OCR via Cloud

```bash
POST /ocr
Content-Type: application/json

{
  "image_url": "https://...",
  "language": "ara+eng"
}
```

---

## graph-tools (Port 3005)

Knowledge graph operations.

### Health

```bash
GET /health
```

### Generate Graph

```bash
POST /generate
Content-Type: application/json

{
  "note_path": "Physics/Thermodynamics.md",
  "depth": 2,
  "format": "json"
}
```

Response:
```json
{
  "nodes": [
    {"id": "thermodynamics", "label": "Thermodynamics"},
    {"id": "entropy", "label": "Entropy"}
  ],
  "edges": [
    {"source": "thermodynamics", "target": "entropy"}
  ]
}
```

### Export Graph

```bash
POST /export
Content-Type: application/json

{
  "format": "svg",
  "layout": "force-directed"
}
```

---

## ocr (Port 3000)

OCR processing service.

### Process Image

```bash
POST /process
Content-Type: multipart/form-data

file: <image>
language: ara+eng+fra
```

Response:
```json
{
  "text": "Extracted text...",
  "confidence": 0.95,
  "language": "eng"
}
```

### Process PDF

```bash
POST /pdf
Content-Type: application/json

{
  "url": "https://example.com/doc.pdf",
  "pages": [1, 2, 3]
}
```

---

## Error Responses

All endpoints return errors in this format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid embedding dimensions",
    "details": {
      "expected": 768,
      "received": 512
    }
  }
}
```

### Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 400 | Bad Request |
| 404 | Not Found |
| 429 | Rate Limited |
| 500 | Internal Error |
