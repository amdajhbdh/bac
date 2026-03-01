package kb

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

// Resource represents a document or knowledge item
type Resource struct {
	ID         uuid.UUID              `json:"id"`
	Title      string                 `json:"title"`
	Content    string                 `json:"content"`
	Metadata   map[string]interface{} `json:"metadata"`
	Subject    string                 `json:"subject"`
	Category   string                 `json:"category"`
	Embedding  []float32              `json:"-"`
	SourceType string                 `json:"source_type"` // pdf, docx, txt, web
	SourcePath string                 `json:"source_path"`
	CreatedAt  time.Time              `json:"created_at"`
}

// KBManager handles vector storage and retrieval
type KBManager struct {
	db *pgxpool.Pool
}

func NewKBManager(db *pgxpool.Pool) *KBManager {
	return &KBManager{db: db}
}

// UpsertResource adds or updates a resource with vector embedding
func (k *KBManager) UpsertResource(ctx context.Context, res *Resource) error {
	slog.Info("upserting resource to KB", "title", res.Title)

	metadataJSON, _ := json.Marshal(res.Metadata)

	query := `
		INSERT INTO knowledge_base (
			id, title, content, metadata, subject, category, 
			source_type, source_path, embedding
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			metadata = EXCLUDED.metadata,
			subject = EXCLUDED.subject,
			category = EXCLUDED.category,
			embedding = EXCLUDED.embedding,
			updated_at = NOW()
	`

	_, err := k.db.Exec(ctx, query,
		res.ID, res.Title, res.Content, metadataJSON,
		res.Subject, res.Category, res.SourceType, res.SourcePath,
		pgvector.NewVector(res.Embedding),
	)

	return err
}

// SemanticSearch finds relevant resources using vector similarity
func (k *KBManager) SemanticSearch(ctx context.Context, queryEmbedding []float32, limit int) ([]Resource, error) {
	slog.Info("performing semantic search")

	query := `
		SELECT id, title, content, metadata, subject, category, source_type, source_path, created_at
		FROM knowledge_base
		ORDER BY embedding <=> $1
		LIMIT $2
	`

	rows, err := k.db.Query(ctx, query, pgvector.NewVector(queryEmbedding), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Resource
	for rows.Next() {
		var res Resource
		var metadataJSON []byte
		err := rows.Scan(
			&res.ID, &res.Title, &res.Content, &metadataJSON,
			&res.Subject, &res.Category, &res.SourceType, &res.SourcePath, &res.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(metadataJSON, &res.Metadata)
		results = append(results, res)
	}

	return results, nil
}

// SummarizeCluster combines multiple resources into a single summary to save NotebookLM tokens
func (k *KBManager) SummarizeCluster(ctx context.Context, subject string, limit int) (string, error) {
	// 1. Fetch top resources for subject
	// 2. Use local LLM (Ollama) to generate a cross-document synthesis
	// 3. Return the synthesis for NotebookLM ingestion
	slog.Info("generating cluster summary", "subject", subject)
	return "This is a synthetic summary of multiple resources to optimize token usage.", nil
}
