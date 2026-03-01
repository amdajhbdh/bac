# NLM Caching System - Tasks

## Phase 1: Core Infrastructure

### 1.1 Query Router
- [ ] 1.1.1 Create `internal/nlm/router.go`
- [ ] 1.1.2 Implement `ExtractRoute()` function
- [ ] 1.1.3 Implement subject keyword matching (math, pc, svt, philosophy)
- [ ] 1.1.4 Implement topic extraction from query
- [ ] 1.1.5 Implement query hash generation (SHA256)
- [ ] 1.1.6 Implement notebook selection based on subject

### 1.2 SQLite Cache Index
- [ ] 1.2.1 Create `internal/nlm/cache/index.go`
- [ ] 1.2.2 Define cache database schema
- [ ] 1.2.3 Implement `Init()` - create tables if not exist
- [ ] 1.2.4 Implement `GetCachedResponse()` - lookup by hash
- [ ] 1.2.5 Implement `SetCacheEntry()` - store new entry
- [ ] 1.2.6 Implement `CleanupExpired()` - remove old entries
- [ ] 1.2.7 Implement `IncrementAccessCount()` - track usage

### 1.3 Garage S3 Storage
- [ ] 1.3.1 Create `internal/nlm/cache/storage.go`
- [ ] 1.3.2 Implement S3 client initialization
- [ ] 1.3.3 Implement `StoreResponse()` - upload to Garage
- [ ] 1.3.4 Implement `GetResponse()` - download from Garage
- [ ] 1.3.5 Implement `DeleteResponse()` - remove from Garage
- [ ] 1.3.6 Implement S3 key generation (date-based path)

### 1.4 Rate Limiter
- [ ] 1.4.1 Create `internal/nlm/ratelimit.go`
- [ ] 1.4.2 Implement `RateLimiter` struct with mutex
- [ ] 1.4.3 Implement `Acquire()` - wait if needed
- [ ] 1.4.4 Implement `Release()` - allow next request
- [ ] 1.4.5 Configure cooldown periods per operation

### 1.5 Configuration
- [ ] 1.5.1 Add `config/notebooks.yaml`
- [ ] 1.5.2 Define notebook mappings (subject → UUID)
- [ ] 1.5.3 Add cache TTL settings
- [ ] 1.5.4 Add Garage bucket configuration
- [ ] 1.5.5 Add environment variable support

## Phase 2: Integration

### 2.1 Cache System Integration
- [ ] 2.1.1 Update `internal/nlm/nlm.go` to use cache
- [ ] 2.1.2 Modify `Query()` function to check cache first
- [ ] 2.1.3 Modify `Query()` function to store results
- [ ] 2.1.4 Add cache hit/miss logging
- [ ] 2.1.5 Add cache statistics tracking

### 2.2 Fallback Integration
- [ ] 2.2.1 Add Ollama fallback when NLM rate-limited
- [ ] 2.2.2 Add online services fallback when both fail
- [ ] 2.2.3 Return stale cache if all fallbacks fail
- [ ] 2.2.4 Implement graceful degradation

### 2.3 Testing Integration
- [ ] 2.3.1 Add unit tests for router
- [ ] 2.3.2 Add unit tests for cache index
- [ ] 2.3.3 Add integration test for full flow
- [ ] 2.3.4 Test rate limiter behavior

## Phase 3: Maintenance

### 3.1 Cleanup Daemon
- [ ] 3.1.1 Add periodic cleanup task
- [ ] 3.1.2 Remove expired SQLite entries
- [ ] 3.1.3 Remove expired S3 objects
- [ ] 3.1.4 Log cleanup statistics

### 3.2 Monitoring
- [ ] 3.2.1 Add cache hit rate metrics
- [ ] 3.2.2 Add NLM query count metrics
- [ ] 3.2.3 Add rate limit occurrence metrics
- [ ] 3.2.4 Log cache statistics on startup

### 3.3 Cache Warmup
- [ ] 3.3.1 Add cache warmup on startup
- [ ] 3.3.2 Pre-populate common queries
- [ ] 3.3.3 Add manual cache clear command
