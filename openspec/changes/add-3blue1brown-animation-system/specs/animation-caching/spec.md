## ADDED Requirements

### Requirement: Cache system stores rendered animations
The caching system SHALL store rendered animations to avoid re-rendering identical requests.

#### Scenario: Store in memory cache
- **WHEN** animation renders successfully
- **THEN** the result SHALL be stored in Redis with 24-hour TTL

#### Scenario: Store on disk cache
- **WHEN** animation renders successfully
- **THEN** the video file SHALL be stored in /data/animations/ directory

#### Scenario: Check memory cache first
- **WHEN** new render request arrives
- **THEN** the system SHALL check Redis cache using hash of problem+solution

### Requirement: Cache key generation
The caching system SHALL generate consistent cache keys from request parameters.

#### Scenario: Generate cache key
- **WHEN** render request is received
- **THEN** cache key SHALL be SHA256(problem + solution + quality + style)

#### Scenario: Same problem returns cached result
- **WHEN** identical request is made within TTL
- **THEN** cached video path SHALL be returned immediately

### Requirement: Cache invalidation
The cache system SHALL support manual and automatic invalidation.

#### Scenario: Manual invalidation
- **WHEN** DELETE /animations/:id/cache is called
- **THEN** the cache entry SHALL be removed from all layers

#### Scenario: TTL expiration
- **WHEN** cache entry exceeds TTL
- **THEN** the entry SHALL be automatically removed

### Requirement: Disk cache management
The disk cache SHALL implement automatic cleanup to prevent disk exhaustion.

#### Scenario: Check disk space before caching
- **WHEN** storing to disk cache
- **THEN** system SHALL verify available space exceeds minimum threshold

#### Scenario: Cleanup old cached files
- **WHEN** disk usage exceeds 80%
- **THEN** system SHALL delete oldest cached files until usage below 60%
