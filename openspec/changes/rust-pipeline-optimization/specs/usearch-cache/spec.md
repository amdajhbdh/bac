# USearch Vector Cache Specification

## Overview

In-memory vector similarity search using USearch for sub-millisecond query performance. Acts as a hot cache layer in front of pgvector.

## Problem Statement

- **Current**: pgvector queries take 5-10ms
- **Bottleneck**: Disk I/O, connection overhead
- **Goal**: < 1ms for hot queries

## Solution Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  Vector Search Architecture                    │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                  Query Request                        │   │
│  └─────────────────────┬───────────────────────────────┘   │
│                        │                                      │
│                        ▼                                      │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              USearch Hot Cache                       │   │
│  │              (In-Memory Index)                       │   │
│  │                                                    │   │
│  │  Index: top 100k most accessed questions          │   │
│  │  Dimension: 1536                                   │   │
│  │  Metric: Cosine                                    │   │
│  │  Quantization: fp16                                │   │
│  │                                                    │   │
│  │  Latency: < 1ms                                   │   │
│  └─────────────────────┬───────────────────────────────┘   │
│                        │                                      │
│           ┌────────────┴────────────┐                        │
│           │                         │                         │
│      CACHE HIT                CACHE MISS                      │
│           │                         │                         │
│           ▼                         ▼                         │
│  ┌─────────────────┐    ┌─────────────────────┐           │
│  │ Return results  │    │ Query pgvector      │           │
│  │ from memory     │    │ (5-10ms)           │           │
│  │ (< 1ms)        │    │                     │           │
│  └─────────────────┘    │ Populate cache     │           │
│                         │ for next time      │           │
│                         └─────────────────────┘           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## ADDED Requirements

### Requirement: Sub-millisecond similarity search
The USearch cache SHALL return similarity search results in < 1ms for cached vectors.

#### Scenario: Hot query returns instantly
- **WHEN** query vector matches cached index
- **THEN** results returned in < 1ms
- **AND** top-k results include similarity scores

#### Scenario: Cache miss falls back to pgvector
- **WHEN** query vector not in USearch cache
- **THEN** pgvector is queried
- **AND** results are cached for next time

### Requirement: LRU eviction for memory management
The cache SHALL evict least-recently-used entries when memory limit reached.

#### Scenario: Cache reaches memory limit
- **WHEN** cache size exceeds max_memory (e.g., 1GB)
- **THEN** least recently used entries are evicted
- **AND** new entries are added

#### Scenario: Entry accessed updates LRU position
- **WHEN** entry is returned from search
- **THEN** entry's last_accessed timestamp is updated
- **AND** entry moves to front of LRU queue

### Requirement: FP16 quantization for memory efficiency
The cache SHALL use fp16 quantization to reduce memory usage by 50%.

#### Scenario: FP16 quantization reduces memory
- **WHEN** vectors stored in fp16 format
- **THEN** memory usage is ~50% of fp32
- **AND** search accuracy remains > 99%

### Requirement: Automatic cache warming
The cache SHALL be warmed with popular queries on startup.

#### Scenario: Cache warmed on startup
- **WHEN** service starts
- **THEN** most accessed queries are loaded into USearch
- **AND** cache hit rate is immediately high

#### Scenario: Background warming job
- **WHEN** cache hit rate drops below 80%
- **THEN** background job loads more popular queries
- **AND** cache hit rate recovers

### Requirement: Index persistence to disk
The cache SHALL persist index to disk for fast restart.

#### Scenario: Service restarts with disk index
- **WHEN** service restarts
- **THEN** index is loaded from disk
- **AND** cache is immediately warm

### Requirement: Metrics exposed for monitoring
The cache SHALL expose Prometheus metrics.

#### Scenario: Cache metrics accessible
- **WHEN** client queries /metrics
- **THEN** returns:
  - usearch_cache_size (gauge)
  - usearch_cache_hits_total (counter)
  - usearch_cache_misses_total (counter)
  - usearch_hit_rate (gauge)
  - usearch_query_latency_ms (histogram)

## Interface

### Go API

```go
type USearchCache struct {
    index *usearch.Index
    pg    *db.Queries
    redis *redis.Client
}

func NewUSearchCache(cfg Config) (*USearchCache, error)

func (c *USearchCache) Search(ctx context.Context, query []float32, k int) ([]SearchResult, error)

func (c *USearchCache) Add(ctx context.Context, id uint64, vector []float32, problem string) error

func (c *USearchCache) Warm(ctx context.Context) error

func (c *USearchCache) GetStats() CacheStats
```

### Configuration

```yaml
usearch:
  dimension: 1536
  metric: cosine
  quantization: fp16
  connectivity: 16
  expansion_add: 128
  expansion_search: 64
  max_memory_mb: 2048
  warm_top_k: 100000
  warm_queries: 1000
```

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Query latency (hot) | < 1ms | p99 |
| Query latency (cold) | < 10ms | p99 |
| Cache hit rate | > 90% | sustained |
| Memory usage | < 2GB | peak |
| Index load time | < 5s | from disk |
| Batch insert | 10k/s | sustained |

## Data Flow

```
1. Generate embedding for query
2. Check USearch (hot) → if hit, return
3. Check Redis (warm) → if hit, populate USearch, return
4. Query pgvector (cold) → populate Redis + USearch
5. Return results
```

## Acceptance Criteria

1. **Performance**: Hot queries < 1ms latency
2. **Hit Rate**: > 90% for popular queries
3. **Memory**: < 2GB for 100k vectors
4. **Persistence**: Fast restart from disk
5. **Monitoring**: Metrics exposed for Grafana
