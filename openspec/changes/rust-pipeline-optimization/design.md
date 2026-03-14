# Design: Rust-Powered Pipeline Optimization

## Context

### Current Architecture
The existing pipeline processes documents sequentially:
```
Input → OCR (Go + exec) → Solver (Go) → Storage (pgvector)
```

**Bottlenecks identified:**
1. **OCR**: 3-5s per image - serial execution, no batching
2. **Solver**: Blocks waiting for OCR, then waits for embedding generation
3. **Storage**: Synchronous writes block entire pipeline
4. **Search**: pgvector 5-10ms latency compounds at scale

### Constraints
- Must maintain backward compatibility with existing Go code
- Existing Tesseract traineddata files must work
- French/Arabic OCR support required
- Deployment: Linux servers (no macOS-specific features)

### Stakeholders
- Backend team (Go maintainers)
- ML/AI team (model training)
- DevOps (deployment, monitoring)

## Goals / Non-Goals

**Goals:**
- Reduce end-to-end latency by 80% (15-30s → 3-6s)
- Increase throughput by 50x (10 → 500 docs/min)
- Add vector search latency to <1ms for hot queries
- Maintain 99.9% availability during migration

**Non-Goals:**
- Replace entire Go codebase with Rust
- Add real-time video processing
- Implement distributed processing across multiple machines (Phase 1)
- Migrate existing stored questions to new format

## Decisions

### D1: Rust OCR Service (leptess over tesseract-rs)

**Decision**: Use leptess for Rust Tesseract bindings

**Alternatives considered:**
- `ocrs`: Pure Rust but no French support yet
- `tesseract-rs`: Lower-level, less safe
- Python + subprocess: Current approach, slow

**Rationale**: leptess provides thread-safe access to Tesseract with LSTM engine, maintains compatibility with existing French/Arabic traineddata, and offers both high-level and low-level APIs.

### D2: Fluvio for Streaming (over Kafka/RabbitMQ)

**Decision**: Use Fluvio for event streaming

**Alternatives considered:**
- Kafka: Mature but Java-heavy, complex运维
- RabbitMQ: Good for queues but not streaming
- In-memory channels: Simple but no persistence

**Rationale**: Fluvio is native Rust with WASM transforms, automatic backpressure, and simpler deployment. Good fit for Go/Rust hybrid architecture.

### D3: USearch for Vector Cache (over Qdrant)

**Decision**: Use USearch as embedded hot cache

**Alternatives considered:**
- Qdrant: Full vector DB, requires separate service
- pgvector: Current solution, sufficient for cold data
- Pinecone: Cloud-only, not suitable

**Rationale**: USearch is embedded (no separate service), 10x faster than FAISS, supports disk-backed indexes for larger datasets. Will use pgvector for cold data, USearch for hot cache.

### D4: Go + Rust Communication (gRPC over HTTP)

**Decision**: Use gRPC for Go ↔ Rust communication

**Alternatives considered:**
- HTTP REST: Simpler but more overhead
- Unix sockets + JSON: Lower latency but less standard
- Shared memory: Fastest but complex

**Rationale**: gRPC provides type-safe contracts, bidirectional streaming for progress updates, and good Go/Rust support. HTTP/2 under the hood.

### D5: Async Storage with Debouncing

**Decision**: Implement debounced batch writes for storage

**Alternatives considered:**
- Direct writes: Current approach, blocks
- Write-through cache: More complex
- Transactional outbox: Requires additional table

**Rationale**: Simple debouncing (100ms) with batch writes (up to 50 items) reduces DB load while maintaining near-real-time persistence.

## Risks / Trade-offs

### R1: Rust/Go Integration Complexity
**Risk**: FFI overhead, debugging across language boundary
**Mitigation**: 
- Use gRPC for clean separation
- Add logging at both ends
- Create integration tests covering the full pipeline

### R2: USearch Memory Usage
**Risk**: Large embedding index could consume excessive RAM
**Mitigation**:
- Use fp16 quantization (50% memory reduction)
- Implement LRU eviction for oldest entries
- Disk-backed indexes for cold data

### R3: Fluvio Learning Curve
**Risk**: Team unfamiliar with Fluvio
**Mitigation**:
- Start with simple topic consumer
- Add monitoring dashboards early
- Document common patterns

### R4: OCR Accuracy Regression
**Risk**: leptess produces different results than tesseract CLI
**Mitigation**:
- Run A/B comparison during implementation
- Fall back to CLI if confidence drops
- Keep Python tesseract as backup option

### R5: Deployment Complexity
**Risk**: New services = new failure modes
**Mitigation**:
- Feature flags for gradual rollout
- Circuit breaker pattern for Rust service calls
- Comprehensive health checks

## Migration Plan

### Phase 1: Shadow Mode (Week 1-2)
1. Deploy Rust OCR service alongside existing
2. Run both in parallel, compare results
3. Log metrics for both paths
4. No traffic redirection yet

### Phase 2: Canary Release (Week 3)
1. Route 10% traffic to Rust pipeline
2. Monitor error rates, latency
3. Gradually increase if metrics OK
4. Roll back if issues detected

### Phase 3: Full Migration (Week 4)
1. Migrate 100% to optimized pipeline
2. Keep old code for 2 weeks
3. Monitor closely
4. Remove old code after stabilization

### Rollback Plan
1. Revert to old pipeline via feature flag
2. Rust service continues running in shadow
3. Debug issues with live comparison
4. Re-enable when fixed

## Open Questions

### Q1: USearch Index Size
**Question**: How many questions fit in memory for <1ms lookup?
**Answer needed**: Depends on average embedding storage. 1M questions × 3KB = 3GB. Need to benchmark actual usage.

### Q2: Fluvio Deployment
**Question**: Self-hosted Fluvio or Qdrant Cloud?
**Recommendation**: Start self-hosted on existing infrastructure, migrate to cloud if needed.

### Q3: OCR Quality Thresholds
**Question**: What confidence threshold triggers early exit?
**Recommendation**: Start at 0.9, adjust based on production data.

### Q4: Cache Warming Strategy
**Question**: Pre-populate cache on startup or lazily?
**Recommendation**: Lazy with background warming of popular queries.
