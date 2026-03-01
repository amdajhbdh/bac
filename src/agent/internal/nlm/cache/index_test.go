package cache

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIndexInit(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	_, err := NewIndex(dbPath)
	if err != nil {
		t.Fatalf("failed to create index: %v", err)
	}

	if _, err := os.Stat(dbPath); err != nil {
		t.Errorf("database file should exist: %v", err)
	}
}

func TestCacheEntry(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	index, err := NewIndex(dbPath)
	if err != nil {
		t.Fatalf("failed to create index: %v", err)
	}
	defer index.Close()

	entry := &CacheEntry{
		QueryHash:  "testhash123",
		QueryText:  "test problem",
		Subject:    "math",
		Topics:     `["derivative","polynomial"]`,
		NotebookID: "test-notebook-id",
		S3Key:      "responses/2026/03/testhash123-math.json",
		ExpiresAt:  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	}

	if err := index.SetCacheEntry(ctx, entry); err != nil {
		t.Fatalf("failed to set cache entry: %v", err)
	}

	retrieved, err := index.GetCachedResponse(ctx, "testhash123", "math")
	if err != nil {
		t.Fatalf("failed to get cached response: %v", err)
	}

	if retrieved == nil {
		t.Fatal("retrieved entry should not be nil")
	}

	if retrieved.QueryHash != entry.QueryHash {
		t.Errorf("query hash mismatch: got %q, want %q", retrieved.QueryHash, entry.QueryHash)
	}

	if retrieved.QueryText != entry.QueryText {
		t.Errorf("query text mismatch: got %q, want %q", retrieved.QueryText, entry.QueryText)
	}
}

func TestCacheMiss(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	index, err := NewIndex(dbPath)
	if err != nil {
		t.Fatalf("failed to create index: %v", err)
	}
	defer index.Close()

	retrieved, err := index.GetCachedResponse(ctx, "nonexistent", "math")
	if err != nil {
		t.Fatalf("error should be nil for miss: %v", err)
	}

	if retrieved != nil {
		t.Error("retrieved should be nil for cache miss")
	}
}

func TestCleanupExpired(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	index, err := NewIndex(dbPath)
	if err != nil {
		t.Fatalf("failed to create index: %v", err)
	}
	defer index.Close()

	expiredEntry := &CacheEntry{
		QueryHash:  "expiredhash",
		QueryText:  "expired problem",
		Subject:    "math",
		NotebookID: "test-notebook",
		S3Key:      "responses/2026/03/expired.json",
		ExpiresAt:  time.Now().UTC().Add(-1 * time.Hour).Format("2006-01-02 15:04:05"),
	}

	validEntry := &CacheEntry{
		QueryHash:  "validhash",
		QueryText:  "valid problem",
		Subject:    "math",
		NotebookID: "test-notebook",
		S3Key:      "responses/2026/03/valid.json",
		ExpiresAt:  time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	}

	index.SetCacheEntry(ctx, expiredEntry)
	index.SetCacheEntry(ctx, validEntry)

	deleted, err := index.CleanupExpired(ctx)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}

	t.Logf("got %d deleted entries", deleted)

	retrieved, _ := index.GetCachedResponse(ctx, "validhash", "math")
	if retrieved == nil {
		t.Error("valid entry should still exist")
	}
}

func TestIncrementAccessCount(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	index, err := NewIndex(dbPath)
	if err != nil {
		t.Fatalf("failed to create index: %v", err)
	}
	defer index.Close()

	entry := &CacheEntry{
		QueryHash:  "accesshash",
		QueryText:  "test problem",
		Subject:    "math",
		NotebookID: "test-notebook",
		S3Key:      "responses/2026/03/access.json",
		ExpiresAt:  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	}

	index.SetCacheEntry(ctx, entry)

	before, _, _ := index.GetStats(ctx)

	index.IncrementAccessCount(ctx, "accesshash")

	after, _, _ := index.GetStats(ctx)

	if after <= before {
		t.Errorf("access count should have increased")
	}
}

func TestGetStats(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	index, err := NewIndex(dbPath)
	if err != nil {
		t.Fatalf("failed to create index: %v", err)
	}
	defer index.Close()

	hits, total, err := index.GetStats(ctx)
	if err != nil {
		t.Fatalf("get stats failed: %v", err)
	}

	if total != 0 {
		t.Errorf("initial total should be 0, got %d", total)
	}

	entry := &CacheEntry{
		QueryHash:  "statshash",
		QueryText:  "test problem",
		Subject:    "math",
		NotebookID: "test-notebook",
		S3Key:      "responses/2026/03/stats.json",
		ExpiresAt:  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	}
	index.SetCacheEntry(ctx, entry)
	index.IncrementAccessCount(ctx, "statshash")

	hits, total, _ = index.GetStats(ctx)
	if total != 1 {
		t.Errorf("total should be 1, got %d", total)
	}
	if hits < 1 {
		t.Errorf("hits should be at least 1, got %d", hits)
	}
}
