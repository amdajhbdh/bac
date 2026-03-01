package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bac-unified/agent/internal/kb"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Flags
	action := flag.String("action", "search", "Action: search, index, summarize")
	query := flag.String("query", "", "Search query")
	flag.String("file", "", "File to index") // TODO: implement file indexing
	subject := flag.String("subject", "Math", "Subject filter")

	flag.Parse()

	// DB Connection (Neon)
	dbURL := "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb"

	ctx := context.Background()
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("failed to connect to DB", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	manager := kb.NewKBManager(db)

	switch *action {
	case "search":
		if *query == "" {
			fmt.Println("Error: -query is required for search")
			return
		}

		// In real usage, we'd generate an embedding for the query first
		// Using a dummy embedding for demo
		dummyEmbed := make([]float32, 1536)

		start := time.Now()
		results, err := manager.SemanticSearch(ctx, dummyEmbed, 5)
		if err != nil {
			slog.Error("search failed", "err", err)
			return
		}

		fmt.Printf("Found %d results in %v\n", len(results), time.Since(start))
		for _, r := range results {
			fmt.Printf("- [%s] %s (Source: %s)\n", r.Subject, r.Title, r.SourceType)
		}

	case "summarize":
		summary, err := manager.SummarizeCluster(ctx, *subject, 5)
		if err != nil {
			slog.Error("summarization failed", "err", err)
			return
		}
		fmt.Println("--- Cluster Summary ---")
		fmt.Println(summary)

	default:
		fmt.Printf("Unknown action: %s\n", *action)
	}
}
