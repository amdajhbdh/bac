package nlm

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type CleanupConfig struct {
	Interval     time.Duration
	StartupDelay time.Duration
}

var defaultCleanupConfig = CleanupConfig{
	Interval:     24 * time.Hour,
	StartupDelay: 5 * time.Minute,
}

func StartCleanupDaemon(ctx context.Context, config CleanupConfig) {
	if config.Interval == 0 {
		config = defaultCleanupConfig
	}

	timer := time.NewTimer(config.StartupDelay)
	defer timer.Stop()

	slog.Info("cleanup daemon starting", "interval", config.Interval, "delay", config.StartupDelay)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			slog.Info("cleanup daemon stopping")
			return
		case <-sigChan:
			slog.Info("cleanup daemon received shutdown signal")
			return
		case <-timer.C:
			slog.Info("running scheduled cache cleanup")
			deleted, err := CleanupCache()
			if err != nil {
				slog.Error("cleanup failed", "error", err)
			} else {
				slog.Info("cleanup completed", "deleted", deleted)
			}
			timer.Reset(config.Interval)
		}
	}
}

func StartCleanupDaemonFromEnv(ctx context.Context) {
	intervalStr := os.Getenv("NLM_CLEANUP_INTERVAL")
	interval := defaultCleanupConfig.Interval

	if intervalStr != "" {
		if parsed, err := time.ParseDuration(intervalStr); err == nil {
			interval = parsed
		}
	}

	StartCleanupDaemon(ctx, CleanupConfig{
		Interval:     interval,
		StartupDelay: defaultCleanupConfig.StartupDelay,
	})
}

type CacheMetrics struct {
	Hits           int64   `json:"hits"`
	Total          int64   `json:"total"`
	HitRate        float64 `json:"hit_rate"`
	ExpiredCleaned int64   `json:"expired_cleaned"`
	StorageSize    int64   `json:"storage_size_bytes"`
}

func GetCacheMetrics(ctx context.Context) CacheMetrics {
	stats := GetCacheStats()

	metrics := CacheMetrics{
		Hits:    stats.Hits,
		Total:   stats.Total,
		HitRate: stats.Rate,
	}

	return metrics
}

func LogCacheStats() {
	stats := GetCacheStats()
	slog.Info("cache statistics",
		"hits", stats.Hits,
		"total", stats.Total,
		"hit_rate", stats.Rate,
	)
}

func WarmupCache(ctx context.Context, problems []string) {
	slog.Info("starting cache warmup", "count", len(problems))

	for i, problem := range problems {
		select {
		case <-ctx.Done():
			slog.Info("cache warmup cancelled")
			return
		default:
		}

		slog.Debug("warming cache", "problem", problem, "progress", i+1, "total", len(problems))
		QueryCached(ctx, problem)

		if i < len(problems)-1 {
			time.Sleep(2 * time.Second)
		}
	}

	slog.Info("cache warmup completed")
}

var commonProblems = []string{
	"Résous l'équation x² + x - 2 = 0",
	"Calcule la dérivée de f(x) = x³ + 2x",
	"Quelle est la force exercée sur une masse de 5kg?",
	"Explique la photosynthèse",
	"Qu'est-ce que la liberté en philosophie?",
	"Résous: 2x + 3 = 7",
	"Calcule l'intégrale de x² dx",
	"Quelle est la vitesse après 10s?",
	"Décris la mitose",
	"Explique le concept de justice",
}

func WarmupDefaultCache(ctx context.Context) {
	WarmupCache(ctx, commonProblems)
}

func ClearCache(ctx context.Context) error {
	if cacheIndex == nil {
		return nil
	}

	err := cacheIndex.Clear(ctx)
	if err != nil {
		return err
	}

	slog.Info("cache cleared")
	return nil
}
