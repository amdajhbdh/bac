package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bac-unified/agent/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

var httpClient = &http.Client{
	Timeout: 60 * time.Second,
}

var (
	pool    *pgxpool.Pool
	queries *db.Queries
)

type SimilarProblem struct {
	Question   string
	Solution   string
	Subject    string
	Chapter    string
	Similarity float64
}

type LookupResult struct {
	SimilarProblems []SimilarProblem
	Context         string
}

func initDB() error {
	if pool != nil {
		return nil
	}

	dbURL := os.Getenv("NEON_DB_URL")
	if dbURL == "" {
		dbURL = "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb?channel_binding=require&sslmode=require"
	}

	p, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("connect to db: %w", err)
	}
	pool = p
	queries = db.New(p)
	return nil
}

type LookupFilters struct {
	Subject string
	Chapter string
}

func Lookup(ctx context.Context, problem string, topK int) LookupResult {
	return LookupWithFilters(ctx, problem, topK, "", "")
}

func LookupWithFilters(ctx context.Context, problem string, topK int, subject, chapter string) LookupResult {
	slog.Info("memory lookup started", "problem", problem, "topK", topK, "subject", subject, "chapter", chapter)

	if err := initDB(); err != nil {
		slog.Warn("database init failed", "error", err)
		return fallbackLookup(ctx, problem, topK)
	}

	embedding, err := generateEmbedding(ctx, problem)
	if err != nil {
		slog.Warn("embedding generation failed, using fallback", "error", err)
		return fallbackLookup(ctx, problem, topK)
	}

	fetchLimit := int64(topK * 3)
	allResults, err := queries.GetSimilarQuestions(ctx, db.GetSimilarQuestionsParams{
		Limit:   fetchLimit,
		Column2: pgvector.NewVector(embedding),
	})
	if err != nil {
		slog.Warn("vector search failed, using fallback", "error", err)
		return fallbackLookup(ctx, problem, topK)
	}

	var results []db.GetSimilarQuestionsRow
	for _, r := range allResults {
		if len(results) >= topK {
			break
		}
		if (subject == "" || matchesSubject(r, subject)) && (chapter == "" || matchesChapter(r, chapter)) {
			results = append(results, r)
		}
	}

	if len(results) == 0 {
		slog.Info("no similar problems found")
		return LookupResult{}
	}

	similar := make([]SimilarProblem, len(results))
	for i, r := range results {
		similar[i] = SimilarProblem{
			Question:   r.QuestionText,
			Solution:   r.SolutionText.String,
			Similarity: r.Similarity,
		}
	}

	var contextParts []string
	for _, r := range similar {
		contextParts = append(contextParts, fmt.Sprintf("Problème similaire: %s\nSolution: %s", r.Question, r.Solution))
	}

	slog.Info("memory lookup complete", "found", len(similar))
	return LookupResult{
		SimilarProblems: similar,
		Context:         strings.Join(contextParts, "\n\n"),
	}
}

func matchesSubject(r db.GetSimilarQuestionsRow, subject string) bool {
	return true
}

func matchesChapter(r db.GetSimilarQuestionsRow, chapter string) bool {
	return true
}

func Store(ctx context.Context, problem, solution, subject, chapter string, concepts []string) error {
	slog.Info("storing problem in memory", "subject", subject, "chapter", chapter)

	if err := initDB(); err != nil {
		return fmt.Errorf("database init: %w", err)
	}

	embedding, err := generateEmbedding(ctx, problem+" "+solution)
	if err != nil {
		return fmt.Errorf("generate embedding: %w", err)
	}

	params := db.InsertQuestionParams{
		Column1: pgtype.Text{String: problem, Valid: true},
		Column2: pgtype.Text{String: solution, Valid: true},
		Column3: pgvector.NewVector(embedding),
		Column4: concepts,
		Column5: pgtype.Int4{Int32: 3, Valid: true},
	}

	err = queries.InsertQuestion(ctx, params)
	if err != nil {
		return fmt.Errorf("insert problem: %w", err)
	}

	slog.Info("problem stored successfully")
	return nil
}

func generateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if len(text) > 8000 {
		text = text[:8000]
	}

	embedding, err := getOllamaEmbedding(ctx, text)
	if err != nil {
		slog.Warn("Ollama embedding failed, using hash-based fallback", "error", err)
		return hashToEmbedding(text), nil
	}

	return embedding, nil
}

func hashToEmbedding(text string) []float32 {
	emb := make([]float32, 768)

	h1 := hashString(text)
	h2 := hashString(text + "salt1")
	h3 := hashString(text + "salt2")

	for i := range emb {
		bitPos := i % 64
		bit1 := (h1 >> uint(bitPos)) & 1
		bit2 := (h2 >> uint(bitPos)) & 1
		bit3 := (h3 >> uint(bitPos)) & 1

		emb[i] = float32(bit1)*0.5 + float32(bit2)*0.3 + float32(bit3)*0.2
		if emb[i] == 0 {
			emb[i] = 0.1
		}
	}

	return emb
}

func hashString(s string) uint64 {
	var h uint64 = 14695981039346656037
	for _, c := range s {
		h ^= uint64(c)
		h *= 1099511628211
	}
	return h
}

func getOllamaEmbedding(ctx context.Context, text string) ([]float32, error) {
	type EmbedRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}

	models := []string{"nomic-embed-text", "mxbai-embed-large", "llama3.2:3b"}

	for _, model := range models {
		reqBody, _ := json.Marshal(EmbedRequest{
			Model:  model,
			Prompt: text,
		})

		resp, err := httpClient.Post("http://127.0.0.1:11434/api/embeddings", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			slog.Debug("model not available", "model", model, "error", err)
			continue
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			slog.Debug("decode failed", "model", model, "error", err)
			continue
		}

		emb, ok := result["embedding"].([]interface{})
		if !ok {
			slog.Debug("invalid embedding format", "model", model)
			continue
		}

		embedding := make([]float32, 768)
		for i := 0; i < 768; i++ {
			if i < len(emb) {
				embedding[i] = float32(emb[i].(float64))
			} else {
				embedding[i] = 0.0
			}
		}

		slog.Info("embedding generated", "model", model, "dim", len(embedding))
		return embedding, nil
	}

	return nil, fmt.Errorf("no embedding model available")
}

func fallbackLookup(ctx context.Context, problem string, topK int) LookupResult {
	if err := initDB(); err != nil {
		return LookupResult{}
	}

	params := db.SearchQuestionsByTextParams{
		QuestionText: "%" + problem + "%",
		Limit:        int64(topK),
	}

	results, err := queries.SearchQuestionsByText(ctx, params)
	if err != nil {
		return LookupResult{}
	}

	similar := make([]SimilarProblem, len(results))
	for i, r := range results {
		similar[i] = SimilarProblem{
			Question: r.QuestionText,
			Solution: r.SolutionText.String,
		}
	}

	if len(similar) == 0 {
		return LookupResult{}
	}

	var contextParts []string
	for _, r := range similar {
		contextParts = append(contextParts, fmt.Sprintf("Problème similaire: %s\nSolution: %s", r.Question, r.Solution))
	}

	return LookupResult{
		SimilarProblems: similar,
		Context:         strings.Join(contextParts, "\n\n"),
	}
}
