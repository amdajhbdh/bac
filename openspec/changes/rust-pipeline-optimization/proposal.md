# Proposal: Rust-Powered OCR Pipeline Optimization

## Why

The current OCR → Solve → Store pipeline runs sequentially in Go with significant performance bottlenecks:
- OCR processing takes 3-5 seconds per image due to serial execution
- Solver waits for OCR to complete before starting
- Storage operations block the entire pipeline
- Vector search (pgvector) has 5-10ms latency that compounds at scale

By introducing Rust-powered components (leptess OCR, Fluvio streaming, USearch vector cache), we can achieve:
- **5x faster** single document processing (15-30s → 3-6s)
- **50x higher** throughput (10 docs/min → 500 docs/min)
- **95% reduction** in vector search latency (5-10ms → 0.1-0.5ms)

## What Changes

### New Components
1. **Rust OCR Service** (`src/ocr-service/`) - High-performance OCR using leptess + Rayon parallelism
2. **Fluvio Streaming Pipeline** - Event-driven processing with automatic backpressure
3. **USearch Vector Cache** - In-memory hot cache for instant similarity search
4. **Async Storage Pipeline** - Non-blocking batch writes with debouncing

### Modified Components
1. **Solver Service** - Hybrid cold/warm query strategy with USearch cache
2. **Memory Module** - Multi-level caching (L1: USearch, L2: Redis, L3: pgvector)
3. **Main Pipeline** - Orchestrates all stages with progress streaming

### Breaking Changes
- None - new components are additive
- Existing Go code continues to work as fallback

## Capabilities

### New Capabilities
- `rust-ocr-service`: Native Rust OCR with 8-way parallelism and batch processing
- `fluvio-streaming`: Real-time event streaming with WASM transforms
- `usearch-cache`: Sub-millisecond vector similarity search
- `async-storage`: Non-blocking storage with adaptive batching

### Modified Capabilities
- `pipeline-orchestration`: Enhanced to coordinate Rust + Go hybrid execution
- `memory-lookup`: Now uses 3-tier caching strategy

## Impact

### Code Changes
- New: `src/ocr-service/` (Rust)
- New: `src/fluvio/` (Go client)
- New: `src/cache/` (USearch bindings)
- Modified: `src/agent/internal/solver/`
- Modified: `src/agent/internal/memory/`

### Dependencies Added
- Rust: leptess, pdf-rs, rayon, tokio, flume, usearch
- Go: fluvio-go, usearch (via CGO)

### Systems Affected
- OCR processing (moves from Go+exec to Rust service)
- Vector search (adds USearch hot cache layer)
- Storage (becomes async with batching)
