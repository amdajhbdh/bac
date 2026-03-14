# NLM Caching System Specification

## Overview

A hybrid caching system for NotebookLM (NLM) queries that combines local SQLite for fast lookups with Garage S3 for durable blob storage. This reduces NLM API calls by 60-80% while avoiding rate limits.

## Problem Statement

- **NLM Rate Limits**: ~2-5 seconds delay between operations
- **Session Expiry**: Sessions timeout after ~20 minutes
- **Single Notebook Bottleneck**: All queries hit one notebook
- **No Caching**: Every query hits NLM directly

## Solution Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    NLM Query Flow                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  User Query: "Calcule la dérivée de x³ + 2x"               │
│                      │                                       │
│                      ▼                                       │
│  ┌─────────────────────────────────────────┐                 │
│  │  1. Query Router                        │                 │
│  │     - Extract subject: "math"           │                 │
│  │     - Extract topics: ["derivative"]    │                 │
│  │     - Generate query hash               │                 │
│  └──────────────────┬────────────────────┘                 │
│                     │                                        │
│                     ▼                                        │
│  ┌─────────────────────────────────────────┐                 │
│  │  2. Local Cache Check (SQLite)          │                 │
│  │     - Index: query_hash, subject, tags  │                 │
│  │     - Check TTL (default: 24h)          │                 │
│  └──────────────────┬────────────────────┘                 │
│                     │                                        │
│         ┌───────────┴───────────┐                           │
│         │                       │                           │
│      CACHE HIT              CACHE MISS                      │
│         │                       │                           │
│         ▼                       ▼                           │
│   Return cached           ┌─────────────────────┐            │
│   response              │  3. Rate Limiter    │            │
│                        │     (per notebook)  │            │
│                        └─────────┬───────────┘            │
│                                  │                         │
│                    ┌─────────────┴─────────────┐           │
│                    │                           │           │
│               RATE LIMITED              AVAILABLE            │
│                    │                           │            │
│                    ▼                           ▼            │
│           Use Ollama              ┌───────────────────┐    │
│           fallback                │  4. Query NLM    │    │
│                                   │  (to appropriate │    │
│                                   │   notebook)      │    │
│                                   └────────┬──────────┘    │
│                                            │               │
│                                            ▼               │
│                                   ┌───────────────────┐    │
│                                   │  5. Store Result │    │
│                                   │  - Garage S3     │    │
│                                   │  - SQLite index  │    │
│                                   └───────────────────┘    │
│                                            │               │
│                                            ▼               │
│                                   Return response          │
└─────────────────────────────────────────────────────────────┘
```

## Components

### 1. Query Router (`internal/nlm/router.go`)

**Purpose**: Parse user query to extract subject, topics, and route to appropriate notebook.

**Functions**:
```go
type Route struct {
    NotebookID   string
    Subject     string
    Topics      []string
    QueryHash   string
}

func ExtractRoute(problem string) Route
func SelectNotebook(subject string, topics []string) string
func GenerateQueryHash(problem string) string
```

**Subject Mapping**:
| Subject | Keywords | Notebook |
|---------|----------|----------|
| math | équation, fonction, dérivée, intégrale, limite, polynôme | math-notebook |
| pc | physique, mécanique, électrique, thermodynamique, chimie | pc-notebook |
| svt | biologie, cellule, ADN, génétique, écologie | svt-notebook |
| philosophy | philosophie, morale, éthique, liberté | philosophy-notebook |

### 2. Local Cache Index (`internal/nlm/cache/index.go`)

**Purpose**: SQLite-based fast lookups with TTL tracking.

**Schema**:
```sql
CREATE TABLE cache_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query_hash TEXT NOT NULL,
    query_text TEXT NOT NULL,
    subject TEXT NOT NULL,
    topics TEXT,  -- JSON array
    notebook_id TEXT NOT NULL,
    s3_key TEXT NOT NULL,  -- Garage S3 key
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    access_count INTEGER DEFAULT 1,
    last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(query_hash, subject)
);

CREATE INDEX idx_query_hash ON cache_entries(query_hash);
CREATE INDEX idx_subject ON cache_entries(subject);
CREATE INDEX idx_expires ON cache_entries(expires_at);
```

**Functions**:
```go
type CacheEntry struct {
    QueryHash   string
    QueryText   string
    Subject     string
    Topics      []string
    S3Key       string
    ExpiresAt   time.Time
}

func GetCachedResponse(queryHash, subject string) (*CacheEntry, error)
func SetCacheEntry(entry *CacheEntry) error
func CleanupExpired() (int64, error)
func IncrementAccessCount(queryHash string) error
```

### 3. Garage S3 Storage (`internal/nlm/cache/storage.go`)

**Purpose**: Durable blob storage for cached responses.

**S3 Structure**:
```
nlm-cache-bucket/
├── responses/
│   ├── 2026/
│   │   ├── 02/
│   │   │   ├── {query_hash}-math-{timestamp}.json
│   │   │   ├── {query_hash}-pc-{timestamp}.json
│   │   │   └── {query_hash}-svt-{timestamp}.json
│   │   └── 03/
│   └── ...
└── metadata/
    └── index.json  -- backup index
```

**S3 Object**:
```json
{
  "query": "Calcule la dérivée de x³ + 2x",
  "query_hash": "a1b2c3d4e5f6",
  "subject": "math",
  "topics": ["derivative", "polynomial"],
  "notebook_id": "uuid",
  "response": "La dérivée de f(x) = x³ + 2x est...",
  "created_at": "2026-02-28T10:00:00Z",
  "ttl_seconds": 86400
}
```

**Functions**:
```go
func StoreResponse(ctx context.Context, key string, data []byte) error
func GetResponse(ctx context.Context, key string) ([]byte, error)
func DeleteResponse(ctx context.Context, key string) error
func ListExpiredKeys(ctx context.Context, before time.Time) ([]string, error)
```

### 4. Rate Limiter (`internal/nlm/ratelimit.go`)

**Purpose**: Per-notebook rate limiting to avoid hitting NLM limits.

**Implementation**:
```go
type RateLimiter struct {
    mu sync.Mutex
    buckets map[string]*Bucket
}

type Bucket struct {
    lastRequest time.Time
    cooldown    time.Duration  // 5 seconds default
}

func (r *RateLimiter) Acquire(ctx context.Context, notebookID string) error
func (r *RateLimiter) Release(notebookID string)
```

**Cooldown Periods**:
| Operation | Cooldown |
|-----------|----------|
| Query | 2 seconds |
| Source add | 2 seconds |
| Audio generate | 5 seconds |
| Research | 5 seconds |

### 5. Notebook Registry (`config/notebooks.yaml`)

**Configuration**:
```yaml
notebooks:
  math:
    id: "uuid-for-math-notebook"
    name: "BAC Math"
    topics: ["algebra", "calculus", "geometry", "trigonometry"]
    
  pc:
    id: "uuid-for-pc-notebook"
    name: "BAC Physique-Chimie"
    topics: ["mechanics", "optics", "thermodynamics", "chemistry"]
    
  svt:
    id: "uuid-for-svt-notebook"
    name: "BAC SVT"
    topics: ["cell", "genetics", "ecology", "evolution"]
    
  philosophy:
    id: "uuid-for-philosophy-notebook"
    name: "BAC Philosophie"
    topics: ["ethics", "metaphysics", "epistemology"]

cache:
  default_ttl_seconds: 86400  # 24 hours
  max_entries: 10000
  garage_bucket: "nlm-cache"
```

## Data Flow

### Query with Cache Hit
```
1. Router.extractRoute("dérivée de x²")
   → {subject: "math", topics: ["derivative"], hash: "abc123"}

2. cache.GetCachedResponse("abc123", "math")
   → CacheEntry{s3_key: "responses/2026/02/abc123-math.json"}

3. storage.GetResponse("responses/2026/02/abc123-math.json")
   → {response: "La dérivée est 2x...", ttl: 86400}

4. cache.IncrementAccessCount("abc123")
   → Returns response to user
```

### Query with Cache Miss
```
1. Router.extractRoute("dérivée de x³")
   → {subject: "math", topics: ["derivative"], hash: "def456"}

2. cache.GetCachedResponse("def456", "math")
   → nil (cache miss)

3. rateLimiter.Acquire("math-notebook-uuid")
   → Wait if needed, then proceed

4. nlm.Query(notebookID, "dérivée de x³")
   → {response: "La dérivée est 3x² + ..."}

5. storage.StoreResponse("responses/2026/02/def456-math.json", data)
   → S3 upload successful

6. cache.SetCacheEntry(entry)
   → SQLite index updated

7. rateLimiter.Release("math-notebook-uuid")
   → Returns response to user
```

## Fallback Strategy

When NLM is unavailable or rate-limited:

```
┌─────────────────────────────────────────┐
│           NLM Unavailable               │
├─────────────────────────────────────────┤
│                                          │
│  1. Check cache (even if expired)        │
│     └── Return stale cache if exists     │
│                                          │
│  2. Try Ollama (local)                  │
│     └── Use solver.Solve()               │
│                                          │
│  3. Try online services (Playwright)     │
│     └── online.AutoSolve()               │
│                                          │
│  4. Return fallback response             │
│     └── "Veuillez réessayer plus tard"  │
└─────────────────────────────────────────┘
```

## Implementation Tasks

### Phase 1: Core Infrastructure
- [ ] 1.1 Create `internal/nlm/router.go` - Query routing
- [ ] 1.2 Create `internal/nlm/cache/index.go` - SQLite cache
- [ ] 1.3 Create `internal/nlm/cache/storage.go` - Garage S3 storage
- [ ] 1.4 Create `internal/nlm/ratelimit.go` - Rate limiting
- [ ] 1.5 Add `config/notebooks.yaml` - Notebook configuration

### Phase 2: Integration
- [ ] 2.1 Update `internal/nlm/nlm.go` to use cache
- [ ] 2.2 Add Garage S3 credentials to config
- [ ] 2.3 Implement fallback to Ollama when NLM limited

### Phase 3: Maintenance
- [ ] 3.1 Add cleanup daemon for expired entries
- [ ] 3.2 Add cache statistics/metrics
- [ ] 3.3 Add cache warmup on startup

## Environment Variables

```bash
# Garage S3
GARAGE_ENDPOINT="http://localhost:3900"
GARAGE_ACCESS_KEY="GKd001902bdf51fcdbeaf4d4a5"
GARAGE_SECRET_KEY="e92117c54a666734d6809f75251c09862bffe6ed610754f2e6e31732267e80b5"
GARAGE_BUCKET="nlm-cache"

# Cache settings
NLM_CACHE_TTL=86400
NLM_CACHE_MAX_ENTRIES=10000

# Rate limiting
NLM_COOLDOWN_QUERY=2
NLM_COOLDOWN_GENERATE=5
```

## Metrics

| Metric | Description |
|--------|-------------|
| cache_hits_total | Total cache hits |
| cache_misses_total | Total cache misses |
| nlm_queries_total | NLM queries made |
| nlm_rate_limited_total | Rate limit occurrences |
| cache_hit_rate | hits / (hits + misses) |
| avg_query_latency_ms | Average query time |

## Benefits

| Metric | Before | After |
|--------|--------|-------|
| NLM API calls | Every query | ~20-40% of queries |
| Cache hit rate | 0% | 60-80% |
| Rate limit errors | Frequent | Minimal |
| Offline support | None | Cache fallback |
| Response latency | 2-10s | <100ms (cache) |

## Files to Create

```
src/agent/internal/nlm/
├── router.go           # Query routing and subject extraction
├── cache/
│   ├── index.go       # SQLite cache index
│   └── storage.go     # Garage S3 storage
├── ratelimit.go       # Per-notebook rate limiting
└── nlm.go             # Updated to use cache system
```

## Testing Strategy

1. **Unit tests** for router (subject extraction)
2. **Integration tests** for cache (SQLite + S3)
3. **Load tests** for rate limiter
4. **E2E tests** full query flow with mock NLM
