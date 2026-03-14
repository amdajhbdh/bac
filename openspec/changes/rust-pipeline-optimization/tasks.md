# Rust Pipeline Optimization - Tasks

## 1. Rust OCR Service (TDD)

### Week 1-2: Setup and Core

- [x] 1.1 Create `src/ocr-service/` Rust project with Cargo.toml
- [x] 1.2 Add dependencies: leptess, pdf-rs, rayon, tokio, flume, prost
- [x] 1.3 Create `Cargo.toml` with leptess bindings
- [x] 1.4 Write tests: `tests/ocr_service_test.rs` - verify OCR result structure

### Week 1-2: OCR Implementation

- [x] 1.5 Implement `src/ocr.rs` - LepTess wrapper with French support
- [x] 1.6 Write TDD test: verify French text extraction accuracy
- [x] 1.7 Implement `src/parallel.rs` - Rayon worker pool
- [x] 1.8 Write TDD test: verify parallel processing speedup (10 images < 5s)

### Week 1-2: gRPC Interface

- [x] 1.9 Define `proto/ocr.proto` - gRPC service definition
- [ ] 1.10 Implement `src/server.rs` - gRPC server
- [ ] 1.11 Write TDD test: verify gRPC request/response flow
- [ ] 1.12 Add Prometheus metrics endpoint

## 2. Streaming Pipeline (TDD)

### Week 3: Core Implementation

- [x] 2.1 Create streaming package with in-memory channels (replaces Fluvio)
- [x] 2.2 Implement StreamProcessor with topics and partitions
- [x] 2.3 Write TDD test: verify topic/message flow
- [x] 2.4 Define topics in code (OCR, Solver, Storage)
- [x] 2.5 Implement publish/subscribe with backpressure
- [x] 2.6 Write TDD test for worker and batch processing

### Week 3: Pipeline Integration

- [ ] 2.7 Integrate with existing OCR service
- [ ] 2.8 Integrate with solver service
- [ ] 2.9 Add context enrichment between stages

## 3. USearch Cache (TDD)

### Week 4: Setup

- [ ] 3.1 Create `src/cache/` Go package with CGO USearch bindings
- [ ] 3.2 Add USearch C library to project
- [ ] 3.3 Write TDD test: verify USearch index initialization

### Week 4: Core Implementation

- [ ] 3.4 Implement `cache/usearch.go` - hot cache with LRU eviction
- [ ] 3.5 Write TDD test: verify LRU eviction when memory limit reached
- [ ] 3.6 Implement FP16 quantization
- [ ] 3.7 Write TDD test: verify memory usage < 50% of fp32

### Week 4: Integration

- [ ] 3.8 Implement fallback to pgvector on cache miss
- [ ] 3.9 Write TDD test: verify cache miss triggers pgvector query
- [ ] 3.10 Implement cache warming on startup
- [ ] 3.11 Write TDD test: verify cache populated after warming

## 4. Async Storage Pipeline (TDD)

### Week 5: Setup

- [ ] 4.1 Create `src/storage/` Go package
- [ ] 4.2 Write TDD test: verify channel creation

### Week 5: Core Implementation

- [ ] 4.3 Implement `storage/pipeline.go` - debounced batch collector
- [ ] 4.4 Write TDD test: verify batch flush on size limit (50 items)
- [ ] 4.5 Write TDD test: verify batch flush on timeout (500ms)
- [ ] 4.6 Implement parallel embedding generation with Rayon
- [ ] 4.7 Write TDD test: verify parallel embedding speedup

### Week 5: Integration

- [ ] 4.8 Implement write-through to Redis
- [ ] 4.9 Write TDD test: verify Redis updated immediately
- [ ] 4.10 Implement overflow handling (direct write when full)
- [ ] 4.11 Write TDD test: verify overflow bypass works

## 5. Memory Lookup Enhancement (TDD)

### Week 6: Multi-tier Cache

- [ ] 5.1 Update `src/agent/internal/memory/memory.go` - 3-tier fallback
- [ ] 5.2 Write TDD test: verify USearch → Redis → pgvector fallback
- [ ] 5.3 Implement cache population on miss
- [ ] 5.4 Write TDD test: verify caches populated after pgvector query

### Week 6: Integration

- [ ] 5.5 Implement filter support (subject, chapter)
- [ ] 5.6 Write TDD test: verify subject filter works
- [ ] 5.7 Add metrics for cache hit rates
- [ ] 5.8 Write TDD test: verify metrics exposed correctly

## 6. Pipeline Orchestration (TDD)

### Week 7: Core Implementation

- [ ] 6.1 Create `src/pipeline/` Go package
- [ ] 6.2 Implement `pipeline/orchestrator.go` - stage coordination
- [ ] 6.3 Write TDD test: verify full pipeline execution
- [ ] 6.4 Implement concurrent stage execution (OCR + memory in parallel)
- [ ] 6.5 Write TDD test: verify parallel execution time savings

### Week 7: Resilience

- [ ] 6.6 Implement circuit breaker for Rust OCR calls
- [ ] 6.7 Write TDD test: verify circuit opens after failures
- [ ] 6.8 Implement graceful shutdown with drain
- [ ] 6.9 Write TDD test: verify in-flight requests complete on shutdown

### Week 7: Observability

- [ ] 6.10 Add streaming progress updates via channels
- [ ] 6.11 Write TDD test: verify progress events emitted
- [ ] 6.12 Implement health check endpoint
- [ ] 6.13 Write TDD test: verify health check returns correct status

## 7. Integration Testing

### Week 8: E2E Tests

- [ ] 7.1 Run shadow mode with old and new pipelines
- [ ] 7.2 Compare OCR accuracy between implementations
- [ ] 7.3 Measure end-to-end latency improvement
- [ ] 7.4 Run load test with 1000 documents

### Week 8: Performance Tuning

- [ ] 7.5 Optimize worker pool sizes based on benchmarks
- [ ] 7.6 Tune batch sizes and timeouts
- [ ] 7.7 Verify 80% latency reduction target met
- [ ] 7.8 Verify 50x throughput increase target met

## 8. Deployment

### Week 9: Production

- [ ] 8.1 Create Docker Compose for Rust OCR service
- [ ] 8.2 Configure Fluvio in production
- [ ] 8.3 Setup monitoring dashboards (Grafana)
- [ ] 8.4 Deploy to staging environment

### Week 9: Migration

- [ ] 8.5 Implement canary release (10% traffic)
- [ ] 8.6 Monitor error rates and latency
- [ ] 8.7 Gradual rollout to 100%
- [ ] 8.8 Remove legacy code after 2 weeks

## TDD Testing Strategy

Each component follows:
1. **Write failing test** - Define expected behavior
2. **Implement minimal code** - Make test pass
3. **Refactor** - Improve code while keeping tests green

### Test Coverage Targets

| Component | Target Coverage |
|-----------|----------------|
| Rust OCR | > 80% |
| USearch Cache | > 85% |
| Async Storage | > 80% |
| Pipeline | > 75% |

### Performance Test Targets

| Metric | Target | Deadline |
|--------|--------|----------|
| OCR latency | < 1s | Week 2 |
| End-to-end | < 10s | Week 7 |
| Throughput | 100/min | Week 8 |
| Cache hit rate | > 90% | Week 6 |

---

## Dependencies

```
Week 1-2: Rust OCR
Week 3: Fluvio  
Week 4: USearch
Week 5: Async Storage
Week 6: Memory Enhancement
Week 7: Pipeline
Week 8: Testing
Week 9: Deployment
```
