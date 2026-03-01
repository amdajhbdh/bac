package nlm

import (
	"context"
	"os"
	"time"

	"log/slog"

	"github.com/bac-unified/agent/internal/nlm/cache"
	"github.com/bac-unified/agent/internal/solver"
)

type CacheStats struct {
	Hits  int64
	Total int64
	Rate  float64
}

var (
	cacheIndex  *cache.Index
	cacheStore  *cache.Storage
	rateLimiter *RateLimiter
	cacheTTL    = 24 * time.Hour
	cacheStats  CacheStats
)

func InitCache() error {
	slog.Info("initializing NLM cache with Turso (SQLite)")

	dbURL := os.Getenv("NLM_CACHE_DB_URL")
	if dbURL == "" {
		dbURL = os.Getenv("TURSO_DB_URL")
	}

	index, err := cache.NewIndex(dbURL)
	if err != nil {
		slog.Warn("failed to initialize Turso cache index", "error", err)
		return err
	}
	cacheIndex = index

	store, err := cache.NewStorageFromEnv()
	if err != nil {
		slog.Warn("failed to initialize S3 storage", "error", err)
	} else {
		ctx := context.Background()
		if err := store.CreateBucketIfNotExists(ctx); err != nil {
			slog.Warn("failed to create S3 bucket", "error", err)
		}
		cacheStore = store
	}

	rateLimiter = NewRateLimiter()

	if ttlStr := os.Getenv("NLM_CACHE_TTL"); ttlStr != "" {
		if parsed, err := time.ParseDuration(ttlStr); err == nil {
			cacheTTL = parsed
		}
	}

	slog.Info("NLM cache initialized with Turso", "db_url_set", dbURL != "")
	return nil
}

func CloseCache() error {
	if cacheIndex != nil {
		return cacheIndex.Close()
	}
	return nil
}

type CachedQueryResult struct {
	NotebookID string
	Results    string
	Success    bool
	Error      string
	Cached     bool
}

func QueryCached(ctx context.Context, problem string) CachedQueryResult {
	route := ExtractRoute(problem)

	slog.Debug("query routing",
		"subject", route.Subject,
		"topics", route.Topics,
		"hash", route.QueryHash,
	)

	if cacheIndex != nil {
		entry, err := cacheIndex.GetCachedResponse(ctx, route.QueryHash, route.Subject)
		if err == nil && entry != nil && cacheStore != nil {
			s3Data, err := cacheStore.GetResponse(ctx, entry.S3Key)
			if err == nil && s3Data != nil {
				cacheIndex.IncrementAccessCount(ctx, route.QueryHash)
				cacheStats.Hits++
				cacheStats.Total++

				slog.Info("cache hit", "hash", route.QueryHash, "subject", route.Subject)

				return CachedQueryResult{
					NotebookID: entry.NotebookID,
					Results:    s3Data.Response,
					Success:    true,
					Cached:     true,
				}
			}
		}
	}

	slog.Info("cache miss, querying NLM", "hash", route.QueryHash)

	if rateLimiter != nil {
		if err := rateLimiter.Acquire(ctx, route.NotebookID, "query"); err != nil {
			slog.Warn("rate limit acquire failed", "error", err)
		}
	}

	result := Query(ctx, route.NotebookID, problem)

	if rateLimiter != nil {
		rateLimiter.Release(route.NotebookID, "query")
	}

	if result.Success && cacheStore != nil && cacheIndex != nil {
		key, err := cacheStore.StoreResponse(
			ctx,
			route.QueryHash,
			route.Subject,
			problem,
			result.Results,
			route.NotebookID,
			int(cacheTTL.Seconds()),
		)
		if err != nil {
			slog.Warn("failed to store in cache", "error", err)
		} else {
			entry := &cache.CacheEntry{
				QueryHash:  route.QueryHash,
				QueryText:  problem,
				Subject:    route.Subject,
				NotebookID: route.NotebookID,
				S3Key:      key,
				ExpiresAt:  time.Now().Add(cacheTTL).Format("2006-01-02 15:04:05"),
			}
			if err := cacheIndex.SetCacheEntry(ctx, entry); err != nil {
				slog.Warn("failed to index cache entry", "error", err)
			}
		}
	}

	cacheStats.Total++

	return CachedQueryResult{
		NotebookID: result.NotebookID,
		Results:    result.Results,
		Success:    result.Success,
		Error:      result.Error,
		Cached:     false,
	}
}

func QueryCachedWithFallback(ctx context.Context, problem string) CachedQueryResult {
	result := QueryCached(ctx, problem)

	if result.Success {
		return result
	}

	slog.Warn("NLM failed, trying Ollama fallback")

	solverResult, err := solver.Solve(ctx, problem, "")
	if err == nil && solverResult != nil {
		return CachedQueryResult{
			NotebookID: "ollama",
			Results:    solverResult.Solution,
			Success:    true,
			Error:      "",
			Cached:     false,
		}
	}

	if cacheIndex != nil {
		entry, _ := cacheIndex.GetCachedResponse(ctx,
			GenerateQueryHash(problem),
			ExtractRoute(problem).Subject,
		)
		if entry != nil && cacheStore != nil {
			s3Data, _ := cacheStore.GetResponse(ctx, entry.S3Key)
			if s3Data != nil {
				slog.Warn("returning stale cache as last resort")
				return CachedQueryResult{
					NotebookID: entry.NotebookID,
					Results:    s3Data.Response,
					Success:    true,
					Error:      "stale cache",
					Cached:     true,
				}
			}
		}
	}

	return CachedQueryResult{
		NotebookID: "",
		Results:    "",
		Success:    false,
		Error:      "All fallbacks failed",
		Cached:     false,
	}
}

func GetCacheStats() CacheStats {
	if cacheIndex != nil {
		hits, total, _ := cacheIndex.GetStats(context.Background())
		cacheStats.Hits = hits
		cacheStats.Total = total
	}

	if cacheStats.Total > 0 {
		cacheStats.Rate = float64(cacheStats.Hits) / float64(cacheStats.Total)
	}

	return cacheStats
}

func CleanupCache() (int64, error) {
	if cacheIndex == nil {
		return 0, nil
	}

	deleted, err := cacheIndex.CleanupExpired(context.Background())
	if err != nil {
		return 0, err
	}

	slog.Info("cache cleanup completed", "deleted", deleted)
	return deleted, nil
}
