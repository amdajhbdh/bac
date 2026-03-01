# NLM Caching System - Tasks

## Phase 1: Core Infrastructure

### 1.1 Query Router
- [x] 1.1.1 Create `internal/nlm/router.go`
- [x] 1.1.2 Implement `ExtractRoute()` function
- [x] 1.1.3 Implement subject keyword matching (math, pc, svt, philosophy)
- [x] 1.1.4 Implement topic extraction from query
- [x] 1.1.5 Implement query hash generation (SHA256)
- [x] 1.1.6 Implement notebook selection based on subject

### 1.2 SQLite Cache Index
- [x] 1.2.1 Create `internal/nlm/cache/index.go`
- [x] 1.2.2 Define cache database schema
- [x] 1.2.3 Implement `Init()` - create tables if not exist
- [x] 1.2.4 Implement `GetCachedResponse()` - lookup by hash
- [x] 1.2.5 Implement `SetCacheEntry()` - store new entry
- [x] 1.2.6 Implement `CleanupExpired()` - remove old entries
- [x] 1.2.7 Implement `IncrementAccessCount()` - track usage

### 1.3 Garage S3 Storage
- [x] 1.3.1 Create `internal/nlm/cache/storage.go`
- [x] 1.3.2 Implement S3 client initialization
- [x] 1.3.3 Implement `StoreResponse()` - upload to Garage
- [x] 1.3.4 Implement `GetResponse()` - download from Garage
- [x] 1.3.5 Implement `DeleteResponse()` - remove from Garage
- [x] 1.3.6 Implement S3 key generation (date-based path)

### 1.4 Rate Limiter
- [x] 1.4.1 Create `internal/nlm/ratelimit.go`
- [x] 1.4.2 Implement `RateLimiter` struct with mutex
- [x] 1.4.3 Implement `Acquire()` - wait if needed
- [x] 1.4.4 Implement `Release()` - allow next request
- [x] 1.4.5 Configure cooldown periods per operation

### 1.5 Configuration
- [x] 1.5.1 Add `config/notebooks.yaml`
- [x] 1.5.2 Define notebook mappings (subject → UUID)
- [x] 1.5.3 Add cache TTL settings
- [x] 1.5.4 Add Garage bucket configuration
- [x] 1.5.5 Add environment variable support

## Phase 2: Integration

### 2.1 Cache System Integration
- [x] 2.1.1 Update `internal/nlm/nlm.go` to use cache
- [x] 2.1.2 Modify `Query()` function to check cache first
- [x] 2.1.3 Modify `Query()` function to store results
- [x] 2.1.4 Add cache hit/miss logging
- [x] 2.1.5 Add cache statistics tracking

### 2.2 Fallback Integration
- [x] 2.2.1 Add Ollama fallback when NLM rate-limited
- [x] 2.2.2 Add online services fallback when both fail
- [x] 2.2.3 Return stale cache if all fallbacks fail
- [x] 2.2.4 Implement graceful degradation

### 2.3 Testing Integration
- [x] 2.3.1 Add unit tests for router
- [x] 2.3.2 Add unit tests for cache index
- [x] 2.3.3 Add integration test for full flow

## Phase 3: Maintenance

### 3.1 Cleanup Daemon
- [x] 3.1.1 Add periodic cleanup task
- [x] 3.1.2 Remove expired SQLite entries
- [x] 3.1.3 Remove expired S3 objects
- [x] 3.1.4 Log cleanup statistics

### 3.2 Monitoring
- [x] 3.2.1 Add cache hit rate metrics
- [x] 3.2.2 Add NLM query count metrics
- [x] 3.2.3 Add rate limit occurrence metrics
- [x] 3.2.4 Log cache statistics on startup

### 3.3 Cache Warmup
- [x] 3.3.1 Add cache warmup on startup
- [x] 3.3.2 Pre-populate common queries
- [x] 3.3.3 Add manual cache clear command
