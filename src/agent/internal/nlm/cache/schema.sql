-- NLM Cache Schema for Turso (SQLite)

CREATE TABLE IF NOT EXISTS nlm_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query_hash TEXT NOT NULL,
    query_text TEXT NOT NULL,
    subject TEXT NOT NULL,
    topics TEXT DEFAULT '[]',
    notebook_id TEXT NOT NULL,
    s3_key TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    expires_at TEXT NOT NULL,
    access_count INTEGER DEFAULT 1,
    last_accessed TEXT DEFAULT (datetime('now')),
    UNIQUE(query_hash, subject)
);

CREATE INDEX IF NOT EXISTS idx_nlm_cache_hash ON nlm_cache(query_hash);
CREATE INDEX IF NOT EXISTS idx_nlm_cache_subject ON nlm_cache(subject);
CREATE INDEX IF NOT EXISTS idx_nlm_cache_expires ON nlm_cache(expires_at);
