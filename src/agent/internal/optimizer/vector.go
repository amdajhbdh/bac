package optimizer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VectorOptimizer struct {
	pool *pgxpool.Pool
}

func NewVectorOptimizer(pool *pgxpool.Pool) *VectorOptimizer {
	return &VectorOptimizer{pool: pool}
}

func (o *VectorOptimizer) EnsureIndexes(ctx context.Context) error {
	slog.Info("creating vector search indexes")

	indexes := []string{
		// HNSW index for vector similarity search - much faster than IVFFlat
		`CREATE INDEX IF NOT EXISTS idx_questions_vector_hnsw 
		 ON questions USING hnsw (question_vector vector_cosine_ops)
		 WITH (m = 16, ef_construction = 64)`,

		// Composite index for filtered vector searches
		`CREATE INDEX IF NOT EXISTS idx_questions_subject_vector 
		 ON questions USING hnsw (question_vector vector_cosine_ops)
		 WHERE subject IS NOT NULL`,

		// Index for chapter filtered searches
		`CREATE INDEX IF NOT EXISTS idx_questions_chapter_vector 
		 ON questions USING hnsw (question_vector vector_cosine_ops)
		 WHERE chapter IS NOT NULL`,

		// B-tree indexes for standard filtering
		`CREATE INDEX IF NOT EXISTS idx_questions_subject ON questions (subject)`,
		`CREATE INDEX IF NOT EXISTS idx_questions_chapter ON questions (chapter)`,
		`CREATE INDEX IF NOT EXISTS idx_questions_difficulty ON questions (difficulty)`,
		`CREATE INDEX IF NOT EXISTS idx_questions_created_at ON questions (created_at DESC)`,
	}

	for _, idx := range indexes {
		if _, err := o.pool.Exec(ctx, idx); err != nil {
			slog.Warn("index creation failed", "index", idx[:50], "error", err)
		}
	}

	slog.Info("vector indexes created")
	return nil
}

func (o *VectorOptimizer) AnalyzeTables(ctx context.Context) error {
	tables := []string{"questions", "users", "predictions", "resources"}

	for _, table := range tables {
		_, err := o.pool.Exec(ctx, "ANALYZE "+table)
		if err != nil {
			slog.Warn("analyze failed", "table", table, "error", err)
		}
	}

	return nil
}

func (o *VectorOptimizer) OptimizeConnectionPool(ctx context.Context) error {
	// Set optimal pgvector settings
	settings := []string{
		"SET work_mem = '256MB'",
		"SET maintenance_work_mem = '512MB'",
		"SET max_parallel_workers_per_gather = 4",
	}

	for _, setting := range settings {
		if _, err := o.pool.Exec(ctx, setting); err != nil {
			slog.Warn("setting optimization failed", "setting", setting, "error", err)
		}
	}

	return nil
}

func ConnectAndOptimize(dbURL string) (*pgxpool.Pool, error) {
	if dbURL == "" {
		dbURL = os.Getenv("NEON_DB_URL")
	}
	if dbURL == "" {
		dbURL = "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb?sslmode=require"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	optimizer := NewVectorOptimizer(pool)
	if err := optimizer.EnsureIndexes(ctx); err != nil {
		slog.Warn("index optimization failed", "error", err)
	}

	if err := optimizer.AnalyzeTables(ctx); err != nil {
		slog.Warn("table analysis failed", "error", err)
	}

	return pool, nil
}
