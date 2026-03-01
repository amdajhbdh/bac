# Integration Inventory - BAC Unified

## Overview

This document captures the complete inventory of external integrations for the BAC Unified platform.

## Current Integrations

### 1. AI Services

| Integration | Surface | Status | Reusability |
|------------|---------|--------|-------------|
| Ollama | CLI/API | Implemented | HIGH |
| NotebookLM (nlm CLI) | CLI | Implemented | HIGH |
| ChatGPT | Playwright | Implemented | MEDIUM |
| Claude | Playwright | Implemented | MEDIUM |
| DeepSeek | Playwright | Implemented | MEDIUM |
| Grok | Playwright | Implemented | MEDIUM |

### 2. Database & Storage

| Integration | Surface | Status | Reusability |
|------------|---------|--------|-------------|
| PostgreSQL | Go library (jackc/pgx) | Implemented | HIGH |
| pgvector | PostgreSQL extension | Implemented | HIGH |
| Neon | PostgreSQL cloud | Configured | HIGH |
| Garage S3 | S3 API | Implemented | HIGH |
| MinIO | S3 API | Planned | HIGH |
| Redis | Go library | Planned | HIGH |

### 3. Content Sources

| Integration | Surface | Status | Reusability |
|------------|---------|--------|-------------|
| YouTube (yt-dlp) | CLI | Implemented | HIGH |
| Web Scraping (Playwright) | Go library | Implemented | HIGH |
| OER Commons | API | Planned | HIGH |
| OpenStax | API | Planned | MEDIUM |
| Khan Academy | Scraping | Planned | MEDIUM |
| MIT OCW | API | Planned | MEDIUM |

### 4. Content Processing

| Integration | Surface | Status | Reusability |
|------------|---------|--------|-------------|
| Tesseract | CLI | Implemented | HIGH |
| Surya (ML OCR) | Python | Planned | HIGH |
| pdfcpu | Go library | Planned | HIGH |
| Poppler | CLI | Planned | HIGH |
| Apache Tika | Java | Planned | HIGH |

### 5. Search

| Integration | Surface | Status | Reusability |
|------------|---------|--------|-------------|
| pgvector | PostgreSQL extension | Implemented | HIGH |
| Elasticsearch | REST API | Planned | HIGH |
| OpenSearch | REST API | Planned | HIGH |
| MeiliSearch | REST API | Planned | MEDIUM |

### 6. File Bundling

| Integration | Surface | Status | Reusability |
|------------|---------|--------|-------------|
| WRAPX | Custom format | RFC | HIGH |
| ARX | Rust library | RFC | HIGH |
| Research Object Bundle | Spec | RFC | HIGH |
| BagIt | Format | RFC | MEDIUM |

### 7. Monitoring

| Integration | Surface | Status | Reusability |
|------------|---------|--------|-------------|
| Prometheus | Go library | Planned | HIGH |
| Grafana | REST API | Planned | HIGH |
| Sentry | Go library | Planned | MEDIUM |

### 8. Authentication

| Integration | Surface | Status | Reusability |
|------------|---------|--------|-------------|
| JWT | Go library | Planned | HIGH |
| bcrypt | Go stdlib | Implemented | HIGH |
| TOTP | Go library | Planned | HIGH |
| OAuth2 | Various | Planned | MEDIUM |

---

## Pattern Classification

### Adapter Pattern
- Used for: AI services (Ollama, NLM, Cloud APIs)
- Characteristics: Interface abstraction, retry logic, fallback chain

### CLI Wrapper Pattern
- Used for: yt-dlp, Tesseract, ffmpeg
- Characteristics: Process spawning, output parsing, timeout handling

### REST API Client Pattern
- Used for: Elasticsearch, Cloud APIs
- Characteristics: HTTP client, pagination, error translation

### Database Driver Pattern
- Used for: PostgreSQL, Redis
- Characteristics: Connection pooling, query building, migrations

### File Format Pattern
- Used for: WRAPX, ARX, BagIt
- Characteristics: Stream processing, metadata extraction

---

## Integration Evidence

### Go Agent Integrations

```
src/agent/internal/
├── solver/          # Ollama integration
├── nlm/            # NotebookLM CLI integration  
├── online/         # Playwright cloud APIs
├── memory/         # pgvector integration
├── ocr/            # Tesseract integration
└── cache/          # PostgreSQL cache (NEW)
```

### Evidence Links

- Ollama: `internal/solver/solver.go` - LLM inference
- NLM: `internal/nlm/nlm.go` - NotebookLM CLI wrapper
- Online: `internal/online/cloudapi.go` - Cloud API clients
- Memory: `internal/memory/memory.go` - Vector search
- OCR: `internal/ocr/ocr.go` - Tesseract CLI wrapper

---

## Reusable Modules

### 1. Rate Limiter (`internal/nlm/ratelimit.go`)
- Per-resource rate limiting
- Cooldown management
- Can be extracted for general use

### 2. Cache System (`internal/nlm/cache/`)
- PostgreSQL-backed caching
- TTL management
- S3 integration for blob storage
- Can be reused for any cache

### 3. Router (`internal/nlm/router.go`)
- Subject/topic extraction
- Query hashing
- Notebook selection

---

## Anti-Patterns to Fix

1. **Playwright in-process**: Current implementation spawns browser per request
   - Should: Use persistent browser sessions

2. **No connection pooling**: Each service creates new connections
   - Should: Implement connection pool

3. **Hardcoded credentials**: Some integrations have hardcoded values
   - Should: All via environment variables

4. **No circuit breaker**: Failures cascade
   - Should: Implement circuit breaker pattern

---

## Coverage Map

| Category | Unit Tests | Integration Tests | E2E Tests |
|----------|------------|-------------------|------------|
| AI Services | Partial | None | None |
| Database | None | None | Manual |
| Storage | None | None | Manual |
| Search | None | None | Manual |
| Content Processing | None | None | Manual |

---

## Next Steps

1. Add unit tests for existing integrations
2. Implement connection pooling
3. Add circuit breakers
4. Create integration test suite
5. Document error handling patterns
