# Memory Lookup Specification (MODIFIED)

## Overview

Memory lookup with multi-tier caching strategy: L1 (USearch hot), L2 (Redis warm), L3 (pgvector cold). Provides sub-millisecond similarity search for hot queries.

## MODIFIED Requirements

### Requirement: Multi-tier cache fallback
The memory lookup SHALL check caches in order: USearch → Redis → pgvector.

#### Scenario: Hot query hits USearch cache
- **WHEN** query matches USearch index
- **THEN** returns results in < 1ms
- **AND** increments hit counter for USearch

#### Scenario: Warm query hits Redis cache
- **WHEN** query not in USearch but in Redis
- **THEN** returns results in < 5ms
- **AND** populates USearch for next time

#### Scenario: Cold query hits pgvector
- **WHEN** query not in Redis
- **THEN** queries pgvector in < 10ms
- **AND** populates Redis and USearch

### Requirement: Cache population on miss
The lookup SHALL populate warmer caches on cache miss.

#### Scenario: Cold query populates caches
- **WHEN** pgvector returns results
- **THEN** results are stored in Redis
- **AND** background job populates USearch

#### Scenario: Background warming
- **WHEN** service idle and cache warm
- **THEN** background job loads more popular queries
- **AND** cache hit rate improves

### Requirement: Parallel embedding generation
The lookup SHALL generate embeddings in parallel with other operations.

#### Scenario: Embedding generated async
- **WHEN** query requires embedding
- **THEN** embedding generated in parallel with cache check
- **AND** doesn't block response if cache hit

### Requirement: Filter support
The lookup SHALL support filtering by subject and chapter.

#### Scenario: Filter by subject
- **WHEN** lookup called with subject filter
- **THEN** results filtered to matching subject
- **AND** subject metadata included in results

#### Scenario: Filter by chapter
- **WHEN** lookup called with chapter filter
- **THEN** results filtered to matching chapter
- **AND** chapter metadata included in results

### Requirement: Metrics exposed
The lookup SHALL expose cache metrics.

#### Scenario: Cache metrics accessible
- **WHEN** client queries /metrics
- **THEN** returns:
  - memory_lookup_total (counter)
  - memory_usearch_hits (counter)
  - memory_redis_hits (counter)
  - memory_pgvector_hits (counter)
  - memory_cache_hit_rate (gauge)
  - memory_query_latency_ms (histogram)

## ADDED Requirements

### Requirement: Cache invalidation
The lookup SHALL support invalidating cache entries.

#### Scenario: Invalidate single entry
- **WHEN** invalidation request for ID
- **THEN** entry removed from USearch, Redis, and pgvector
- **AND** returns success

#### Scenario: Invalidate by subject
- **WHEN** invalidation request for subject
- **THEN** all entries for subject removed from caches
- **AND** returns count of invalidated entries

## Interface

```go
type MemoryLookup struct {
    usearch *usearch.Cache
    redis   *redis.Client
    pg      *db.Queries
}

type LookupFilters struct {
    Subject string
    Chapter string
    TopK    int
}

type LookupResult struct {
    Problems  []SimilarProblem
    Context   string
    CacheTier string  // "usearch", "redis", "pgvector"
    LatencyMs int64
}

func Lookup(ctx context.Context, problem string, filters LookupFilters) LookupResult

func LookupWithFilters(ctx context.Context, problem string, topK int, subject, chapter string) LookupResult

func Invalidate(ctx context.Context, id uint64) error

func InvalidateBySubject(ctx context.Context, subject string) (int64, error)
```

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| USearch latency | < 1ms | p99 |
| Redis latency | < 5ms | p99 |
| pgvector latency | < 10ms | p99 |
| Overall hit rate | > 90% | sustained |
| Embedding gen | 100/s | parallel |
