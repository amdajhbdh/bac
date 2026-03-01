-- name: GetCachedResponse :one
SELECT 
    id, query_hash, query_text, subject, topics, notebook_id, s3_key,
    created_at, expires_at, access_count, last_accessed
FROM nlm_cache 
WHERE query_hash = ?1 AND subject = ?2 AND expires_at > datetime('now');

-- name: SetCacheEntry :exec
INSERT INTO nlm_cache 
    (query_hash, query_text, subject, topics, notebook_id, s3_key, expires_at, access_count, last_accessed)
VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, 1, datetime('now'))
ON CONFLICT(query_hash, subject) DO UPDATE SET
    access_count = access_count + 1,
    last_accessed = datetime('now');

-- name: CleanupExpired :exec
DELETE FROM nlm_cache WHERE expires_at < datetime('now');

-- name: IncrementAccessCount :exec
UPDATE nlm_cache 
SET access_count = access_count + 1, last_accessed = datetime('now')
WHERE query_hash = ?1;

-- name: GetStats :one
SELECT 
    COALESCE(SUM(access_count), 0) as hits,
    COUNT(*) as total
FROM nlm_cache;

-- name: ClearCache :exec
DELETE FROM nlm_cache;
