package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type DBService struct {
	pool      *pgxpool.Pool
	cron      *cron.Cron
	ollamaURL string
}

func NewDBService(connStr string) (*DBService, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	svc := &DBService{
		pool:      pool,
		cron:      cron.New(),
		ollamaURL: os.Getenv("OLLAMA_URL"),
	}
	if svc.ollamaURL == "" {
		svc.ollamaURL = "http://localhost:11434"
	}

	svc.setupCronJobs()

	return svc, nil
}

func (s *DBService) setupCronJobs() {
	s.cron.AddFunc("@hourly", s.cleanupExpiredCache)
	s.cron.AddFunc("@daily", s.generateDailyAnalytics)
	s.cron.Start()
	slog.Info("cron jobs scheduled", "service", "db")
}

func (s *DBService) Close() {
	s.cron.Stop()
	s.pool.Close()
}

func (s *DBService) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *DBService) cleanupExpiredCache() {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, "SELECT cleanup_expired_cache()")
	if err != nil {
		slog.Error("failed to cleanup cache", "error", err)
	}
}

func (s *DBService) generateDailyAnalytics() {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, "REFRESH MATERIALIZED VIEW hourly_active_users")
	if err != nil {
		slog.Error("failed to refresh analytics", "error", err)
	}
}

type ChatSession struct {
	ID           uuid.UUID       `json:"id"`
	UserID       *uuid.UUID      `json:"user_id,omitempty"`
	SessionType  string          `json:"session_type"`
	Provider     string          `json:"provider"`
	Model        string          `json:"model,omitempty"`
	Title        string          `json:"title,omitempty"`
	Context      json.RawMessage `json:"context,omitempty"`
	MessageCount int             `json:"message_count"`
	TokenCount   int             `json:"token_count"`
	IsActive     bool            `json:"is_active"`
	LastMessage  *time.Time      `json:"last_message_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	EndedAt      *time.Time      `json:"ended_at,omitempty"`
}

type ChatMessage struct {
	ID         uuid.UUID       `json:"id"`
	SessionID  uuid.UUID       `json:"session_id"`
	Role       string          `json:"role"`
	Content    string          `json:"content"`
	TokenCount int             `json:"token_count"`
	Model      string          `json:"model,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

func (s *DBService) CreateChatSession(userID *uuid.UUID, sessionType, provider, model string) (*ChatSession, error) {
	ctx := context.Background()
	session := &ChatSession{
		ID:          uuid.New(),
		UserID:      userID,
		SessionType: sessionType,
		Provider:    provider,
		Model:       model,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	_, err := s.pool.Exec(ctx, `
		INSERT INTO ai_chat_sessions (id, user_id, session_type, provider, model, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		session.ID, session.UserID, session.SessionType, session.Provider, session.Model, session.IsActive, session.CreatedAt)

	return session, err
}

func (s *DBService) GetChatSession(id uuid.UUID) (*ChatSession, error) {
	ctx := context.Background()
	var session ChatSession
	var lastMessage pgtype.Timestamptz
	var endedAt pgtype.Timestamptz

	err := s.pool.QueryRow(ctx, `
		SELECT id, user_id, session_type, provider, model, title, context, 
		       message_count, token_count, is_active, last_message_at, created_at, ended_at
		FROM ai_chat_sessions WHERE id = $1`, id).Scan(
		&session.ID, &session.UserID, &session.SessionType, &session.Provider, &session.Model,
		&session.Title, &session.Context, &session.MessageCount, &session.TokenCount,
		&session.IsActive, &lastMessage, &session.CreatedAt, &endedAt)

	if err != nil {
		return nil, err
	}

	if lastMessage.Valid {
		session.LastMessage = &lastMessage.Time
	}
	if endedAt.Valid {
		session.EndedAt = &endedAt.Time
	}

	return &session, nil
}

func (s *DBService) AddChatMessage(sessionID uuid.UUID, role, content, model string, tokenCount int) (*ChatMessage, error) {
	ctx := context.Background()

	msg := &ChatMessage{
		ID:         uuid.New(),
		SessionID:  sessionID,
		Role:       role,
		Content:    content,
		TokenCount: tokenCount,
		Model:      model,
		CreatedAt:  time.Now(),
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO ai_chat_messages (id, session_id, role, content, token_count, model, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		msg.ID, msg.SessionID, msg.Role, msg.Content, msg.TokenCount, msg.Model, msg.CreatedAt)

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE ai_chat_sessions 
		SET message_count = message_count + 1, 
		    token_count = token_count + $1,
		    last_message_at = $2
		WHERE id = $3`,
		tokenCount, msg.CreatedAt, sessionID)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *DBService) GetChatHistory(sessionID uuid.UUID, limit int) ([]ChatMessage, error) {
	ctx := context.Background()
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.pool.Query(ctx, `
		SELECT id, session_id, role, content, token_count, model, metadata, created_at
		FROM ai_chat_messages 
		WHERE session_id = $1
		ORDER BY created_at ASC
		LIMIT $2`, sessionID, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		if err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content,
			&msg.TokenCount, &msg.Model, &msg.Metadata, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (s *DBService) EndChatSession(sessionID uuid.UUID) error {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, `
		UPDATE ai_chat_sessions 
		SET is_active = FALSE, ended_at = $1
		WHERE id = $2`, time.Now(), sessionID)
	return err
}

type RAGDocument struct {
	ID         uuid.UUID       `json:"id"`
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	SourceType string          `json:"source_type,omitempty"`
	SourceURL  string          `json:"source_url,omitempty"`
	SubjectID  *int            `json:"subject_id,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	IndexedAt  time.Time       `json:"indexed_at"`
	CreatedAt  time.Time       `json:"created_at"`
}

type RAGSearchResult struct {
	ID         uuid.UUID       `json:"id"`
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	Similarity float64         `json:"similarity"`
	Metadata   json.RawMessage `json:"metadata"`
}

func (s *DBService) SearchRAG(query string, matchCount int, matchThreshold float64) ([]RAGSearchResult, error) {
	ctx := context.Background()

	embedding, err := s.generateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if matchCount <= 0 {
		matchCount = 5
	}
	if matchThreshold <= 0 {
		matchThreshold = 0.7
	}

	rows, err := s.pool.Query(ctx, `
		SELECT id, title, content, 1 - (embedding <=> $1::vector) as similarity, metadata
		FROM rag_documents
		WHERE embedding IS NOT NULL 
		AND 1 - (embedding <=> $1::vector) > $2
		ORDER BY embedding <=> $1::vector
		LIMIT $3`, embedding, matchThreshold, matchCount)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []RAGSearchResult
	for rows.Next() {
		var r RAGSearchResult
		if err := rows.Scan(&r.ID, &r.Title, &r.Content, &r.Similarity, &r.Metadata); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

func (s *DBService) generateEmbedding(text string) ([]float64, error) {
	type EmbeddingRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}

	type EmbeddingResponse struct {
		Embedding []float64 `json:"embedding"`
	}

	req := EmbeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(s.ollamaURL+"/api/embeddings", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var embedResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, err
	}

	return embedResp.Embedding, nil
}

func (s *DBService) IndexRAGDocument(title, content, sourceType, sourceURL string, subjectID *int, metadata json.RawMessage) (*RAGDocument, error) {
	ctx := context.Background()

	embedding, err := s.generateEmbedding(content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	doc := &RAGDocument{
		ID:         uuid.New(),
		Title:      title,
		Content:    content,
		SourceType: sourceType,
		SourceURL:  sourceURL,
		SubjectID:  subjectID,
		Metadata:   metadata,
		IndexedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO rag_documents (id, title, content, source_type, source_url, subject_id, embedding, metadata, indexed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		doc.ID, doc.Title, doc.Content, doc.SourceType, doc.SourceURL, doc.SubjectID, embedding, doc.Metadata, doc.IndexedAt, doc.CreatedAt)

	return doc, err
}

func (s *DBService) RecordAnalyticsEvent(eventType string, userID, sessionID, questionID *uuid.UUID, subjectID *int, properties json.RawMessage) error {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, `
		INSERT INTO analytics_events (event_type, user_id, session_id, subject_id, question_id, properties, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		eventType, userID, sessionID, subjectID, questionID, properties, time.Now())
	return err
}

func (s *DBService) GetUserActivity(userID uuid.UUID, since time.Time) ([]map[string]interface{}, error) {
	ctx := context.Background()
	rows, err := s.pool.Query(ctx, `
		SELECT activity_type, subject_id, question_id, points_earned, metadata, created_at
		FROM user_activity_timeline
		WHERE user_id = $1 AND created_at > $2
		ORDER BY created_at DESC
		LIMIT 100`, userID, since)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []map[string]interface{}
	for rows.Next() {
		var activityType string
		var subjectID, questionID *int
		var pointsEarned int
		var metadata json.RawMessage
		var createdAt time.Time

		if err := rows.Scan(&activityType, &subjectID, &questionID, &pointsEarned, &metadata, &createdAt); err != nil {
			return nil, err
		}

		activities = append(activities, map[string]interface{}{
			"type":        activityType,
			"subject_id":  subjectID,
			"question_id": questionID,
			"points":      pointsEarned,
			"metadata":    metadata,
			"created_at":  createdAt,
		})
	}

	return activities, nil
}

func (s *DBService) EnqueueBackgroundJob(jobType string, payload json.RawMessage, priority int) (uuid.UUID, error) {
	ctx := context.Background()
	jobID := uuid.New()

	_, err := s.pool.Exec(ctx, `
		INSERT INTO background_jobs (id, job_type, payload, priority, scheduled_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		jobID, jobType, payload, priority, time.Now(), time.Now())

	return jobID, err
}

func (s *DBService) GetPendingJobs(limit int) ([]map[string]interface{}, error) {
	ctx := context.Background()
	if limit <= 0 {
		limit = 10
	}

	rows, err := s.pool.Query(ctx, `
		SELECT id, job_type, payload, priority, attempts, max_attempts, scheduled_at
		FROM background_jobs
		WHERE status = 'pending' AND scheduled_at <= NOW()
		ORDER BY priority DESC, scheduled_at ASC
		LIMIT $1`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var jobType string
		var payload json.RawMessage
		var priority, attempts, maxAttempts int
		var scheduledAt time.Time

		if err := rows.Scan(&id, &jobType, &payload, &priority, &attempts, &maxAttempts, &scheduledAt); err != nil {
			return nil, err
		}

		jobs = append(jobs, map[string]interface{}{
			"id":           id,
			"type":         jobType,
			"payload":      payload,
			"priority":     priority,
			"attempts":     attempts,
			"max_attempts": maxAttempts,
			"scheduled_at": scheduledAt,
		})
	}

	return jobs, nil
}

func (s *DBService) UpdateJobStatus(jobID uuid.UUID, status string, errorMsg string) error {
	ctx := context.Background()

	var completedAt interface{}
	if status == "completed" || status == "failed" {
		completedAt = time.Now()
	}

	_, err := s.pool.Exec(ctx, `
		UPDATE background_jobs 
		SET status = $1, 
		    last_error = $2,
		    completed_at = $3,
		    attempts = attempts + 1
		WHERE id = $4`, status, errorMsg, completedAt, jobID)

	return err
}

func (s *DBService) GetCache(key string) (json.RawMessage, error) {
	ctx := context.Background()
	var value json.RawMessage
	var expiresAt pgtype.Timestamptz

	err := s.pool.QueryRow(ctx, `
		SELECT value, expires_at FROM api_cache WHERE key = $1`, key).Scan(&value, &expiresAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if expiresAt.Valid && expiresAt.Time.Before(time.Now()) {
		s.pool.Exec(ctx, "DELETE FROM api_cache WHERE key = $1", key)
		return nil, nil
	}

	return value, nil
}

func (s *DBService) SetCache(key string, value json.RawMessage, ttl time.Duration) error {
	ctx := context.Background()
	var expiresAt *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		expiresAt = &t
	}

	_, err := s.pool.Exec(ctx, `
		INSERT INTO api_cache (key, value, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (key) DO UPDATE SET value = $2, expires_at = $3`,
		key, value, expiresAt, time.Now())

	return err
}
