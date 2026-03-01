package cache

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/tursodatabase/go-libsql"
)

type Index struct {
	db  *Queries
	raw *sql.DB
}

func NewIndex(dbURL string) (*Index, error) {
	if dbURL == "" {
		dbURL = os.Getenv("NLM_CACHE_DB_URL")
		if dbURL == "" {
			dbURL = os.Getenv("TURSO_DB_URL")
		}
	}

	if dbURL == "" {
		return nil, fmt.Errorf("database URL not configured")
	}

	if !strings.HasPrefix(dbURL, "libsql://") && !strings.HasPrefix(dbURL, "file://") &&
		!strings.HasPrefix(dbURL, "http://") && !strings.HasPrefix(dbURL, "https://") {
		dbURL = "file:" + dbURL
	}

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("init schema: %w", err)
	}

	return &Index{db: New(db), raw: db}, nil
}

func initSchema(db *sql.DB) error {
	schema := `
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
	`
	_, err := db.Exec(schema)
	return err
}

type CacheEntry struct {
	ID           int64
	QueryHash    string
	QueryText    string
	Subject      string
	Topics       string
	NotebookID   string
	S3Key        string
	CreatedAt    string
	ExpiresAt    string
	AccessCount  int64
	LastAccessed string
}

func (i *Index) GetCachedResponse(ctx context.Context, queryHash, subject string) (*CacheEntry, error) {
	result, err := i.db.GetCachedResponse(ctx, GetCachedResponseParams{
		QueryHash: queryHash,
		Subject:   subject,
	})
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query cache: %w", err)
	}

	return &CacheEntry{
		ID:           result.ID,
		QueryHash:    result.QueryHash,
		QueryText:    result.QueryText,
		Subject:      result.Subject,
		Topics:       result.Topics.String,
		NotebookID:   result.NotebookID,
		S3Key:        result.S3Key,
		CreatedAt:    result.CreatedAt.String,
		ExpiresAt:    result.ExpiresAt,
		AccessCount:  result.AccessCount.Int64,
		LastAccessed: result.LastAccessed.String,
	}, nil
}

func (i *Index) SetCacheEntry(ctx context.Context, entry *CacheEntry) error {
	return i.db.SetCacheEntry(ctx, SetCacheEntryParams{
		QueryHash:  entry.QueryHash,
		QueryText:  entry.QueryText,
		Subject:    entry.Subject,
		Topics:     nullString(entry.Topics),
		NotebookID: entry.NotebookID,
		S3Key:      entry.S3Key,
		ExpiresAt:  entry.ExpiresAt,
	})
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func (i *Index) CleanupExpired(ctx context.Context) (int64, error) {
	err := i.db.CleanupExpired(ctx)
	return 0, err
}

func (i *Index) IncrementAccessCount(ctx context.Context, queryHash string) error {
	return i.db.IncrementAccessCount(ctx, queryHash)
}

func (i *Index) GetStats(ctx context.Context) (hits int64, total int64, err error) {
	result, err := i.db.GetStats(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("get stats: %w", err)
	}
	return toInt64(result.Hits), result.Total, nil
}

func toInt64(v interface{}) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	default:
		return 0
	}
}

func (i *Index) Close() error {
	return i.raw.Close()
}

func (i *Index) Clear(ctx context.Context) error {
	return i.db.ClearCache(ctx)
}
