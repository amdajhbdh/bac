# Async Storage Pipeline Specification

## Overview

Non-blocking storage pipeline with adaptive batching, debouncing, and parallel embedding generation. Reduces storage latency from 2-3s to <100ms.

## Problem Statement

- **Current**: Synchronous DB writes block entire pipeline
- **Bottleneck**: Each store operation waits for DB commit
- **Goal**: <100ms storage latency with batching

## Solution Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  Async Storage Pipeline                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌─────────────┐   │
│  │   Request    │───▶│  Debouncer   │───▶│   Batch     │   │
│  │   Channel    │    │  (100ms)     │    │   Collector │   │
│  └──────────────┘    └──────────────┘    └──────┬──────┘   │
│                                                  │            │
│                                                  ▼            │
│  ┌──────────────┐    ┌──────────────┐    ┌─────────────┐   │
│  │   Embedding │    │   Parallel   │    │   Batch     │   │
│  │   Generator │───▶│   DB Insert  │───▶│   Response  │   │
│  │   (Rayon)   │    │              │    │   Writer   │   │
│  └──────────────┘    └──────────────┘    └─────────────┘   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## ADDED Requirements

### Requirement: Non-blocking storage API
The storage pipeline SHALL accept storage requests without waiting for completion.

#### Scenario: Store request queued immediately
- **WHEN** client calls StoreAsync()
- **THEN** request is queued in channel
- **AND** returns immediately (< 1ms)

#### Scenario: Client notified on completion
- **WHEN** storage completes (success or failure)
- **THEN** callback/channel is notified
- **AND** client can handle errors

### Requirement: Adaptive batching
The pipeline SHALL batch requests based on queue depth and time.

#### Scenario: Batch collected until size limit
- **WHEN** batch reaches 50 items
- **THEN** batch is flushed to database
- **AND** processing starts for next batch

#### Scenario: Batch flushed on timeout
- **WHEN** batch has items but timeout (500ms) expires
- **THEN** partial batch is flushed
- **AND** timer resets for next batch

#### Scenario: High load increases batch size
- **WHEN** queue depth > 500
- **THEN** batch size increases to 100
- **AND** processing parallelism increases

### Requirement: Parallel embedding generation
The pipeline SHALL generate embeddings in parallel using Rayon.

#### Scenario: Batch embeddings generated concurrently
- **WHEN** batch of 50 items processed
- **THEN** embeddings generated in parallel (8 workers)
- **AND** total time ≈ single embedding time

#### Scenario: Embedding fails for one item
- **WHEN** embedding generation fails for one item
- **THEN** other items continue processing
- **AND** failed item is logged for retry

### Requirement: Write-behind caching
The pipeline SHALL write to Redis immediately for fast read-after-write.

#### Scenario: Fast read-after-write
- **WHEN** item stored
- **THEN** Redis updated immediately
- **AND** PostgreSQL updated asynchronously

#### Scenario: Redis unavailable
- **WHEN** Redis write fails
- **THEN** PostgreSQL write still proceeds
- **AND** error is logged but not fatal

### Requirement: Overflow handling
The pipeline SHALL handle overflow gracefully when channel is full.

#### Scenario: Channel full triggers direct write
- **WHEN** storage channel is full (1000 items)
- **THEN** request is written directly (bypass batch)
- **AND** channel capacity is expanded

### Requirement: Metrics and monitoring
The pipeline SHALL expose metrics for monitoring.

#### Scenario: Storage metrics accessible
- **WHEN** client queries /metrics
- **THEN** returns:
  - storage_queue_depth (gauge)
  - storage_batch_size (histogram)
  - storage_latency_ms (histogram)
  - storage_errors_total (counter)
  - storage_throughput (counter)

## Interface

### Go API

```go
type StoragePipeline struct {
    batchChan chan StoreRequest
    db        *db.Queries
    redis     *redis.Client
}

type StoreRequest struct {
    Problem  string
    Solution string
    Subject  string
    Chapter  string
    Concepts []string
    Done     chan error
}

func NewStoragePipeline(db *db.Queries, redis *redis.Client) *StoragePipeline

func (sp *StoragePipeline) StoreAsync(ctx context.Context, req StoreRequest) error

func (sp *StoragePipeline) StoreBatch(ctx context.Context, items []StoreRequest) error

func (sp *StoragePipeline) Close() error

func (sp *StoragePipeline) GetMetrics() StorageMetrics
```

### Configuration

```yaml
storage:
  channel_capacity: 1000
  batch_size: 50
  batch_timeout_ms: 500
  debounce_ms: 100
  embedding_workers: 8
  overflow_bypass: true
  redis_write_through: true
```

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Store latency (async) | < 100ms | p99 |
| Batch insert | 100 items/s | sustained |
| Embedding generation | 100/s | parallel |
| Memory usage | < 500MB | peak |
| Error rate | < 0.1% | failed writes |

## Data Flow

```
1. StoreAsync(request) → channel
2. Debouncer accumulates batch (100ms window)
3. Batch collector gathers up to 50 items
4. Parallel embedding generation (Rayon)
5. Batch insert to PostgreSQL
6. Write-through to Redis
7. Notify completion via callback
```

## Acceptance Criteria

1. **Non-blocking**: StoreAsync returns in < 1ms
2. **Batched**: 50 items inserted in single transaction
3. **Parallel**: Embeddings generated concurrently
4. **Resilient**: Graceful degradation on errors
5. **Monitored**: Metrics exposed for Prometheus
