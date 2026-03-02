-- name: GetSimilarQuestionsVec :many
-- Get similar questions based on vector embedding similarity
SELECT 
    question_text, 
    solution_text, 
    difficulty
FROM questions 
WHERE question_vector IS NOT NULL 
ORDER BY question_vector <=> $1::vector 
LIMIT 10;

-- name: GetSimilarQuestionsFilteredVec :many
-- Get similar questions with difficulty filter
SELECT 
    q.question_text, 
    q.solution_text, 
    q.difficulty
FROM questions q
WHERE q.question_vector IS NOT NULL 
    AND ($2::int IS NULL OR q.difficulty = $2)
ORDER BY q.question_vector <=> $1::vector 
LIMIT 10;

-- name: InsertQuestion :exec
-- Insert a new question with its solution and embedding
INSERT INTO questions (
    question_text, 
    solution_text, 
    question_vector, 
    topic_tags, 
    difficulty
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: SearchQuestionsByText :many
-- Search questions by text similarity
SELECT 
    question_text, 
    solution_text, 
    difficulty
FROM questions 
WHERE question_text ILIKE $1 
LIMIT $2::int;

-- ============================================================================
-- AI Chat Sessions Queries
-- ============================================================================

-- name: CreateChatSession :one
-- Create a new chat session
INSERT INTO ai_chat_sessions (
    user_id, session_type, provider, model, title, context
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, session_type, provider, model, title, context, message_count, token_count, is_active, last_message_at, created_at, ended_at;

-- name: GetChatSession :one
-- Get a chat session by ID
SELECT id, user_id, session_type, provider, model, title, context, message_count, token_count, is_active, last_message_at, created_at, ended_at 
FROM ai_chat_sessions WHERE id = $1;

-- name: ListChatSessions :many
-- List chat sessions for a user
SELECT id, user_id, session_type, provider, model, title, context, message_count, token_count, is_active, last_message_at, created_at, ended_at
FROM ai_chat_sessions 
WHERE user_id = $1 AND is_active = true
ORDER BY last_message_at DESC 
LIMIT $2::int;

-- name: AddChatMessage :one
-- Add a message to a chat session
INSERT INTO ai_chat_messages (
    session_id, role, content, token_count, model, metadata
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, session_id, role, content, token_count, model, metadata, created_at;

-- name: GetChatMessages :many
-- Get messages for a chat session
SELECT id, session_id, role, content, token_count, model, metadata, created_at
FROM ai_chat_messages 
WHERE session_id = $1 
ORDER BY created_at ASC;

-- ============================================================================
-- RAG Document Queries
-- ============================================================================

-- name: InsertRAGDocument :one
-- Insert a RAG document
INSERT INTO rag_documents (
    title, content, source_type, source_url, subject_id, chapter_id, embedding, metadata
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, title, content, source_type, source_url, subject_id, chapter_id, embedding, metadata, chunk_count, indexed_at, created_at;

-- name: SearchRAGDocumentsVec :many
-- Search RAG documents by vector similarity
SELECT 
    id, title, content, source_type, source_url
FROM rag_documents 
WHERE embedding IS NOT NULL
ORDER BY embedding <=> $1::vector 
LIMIT 10;

-- name: SearchRAGDocumentsFilteredVec :many
-- Search RAG documents with subject/chapter filters
SELECT 
    id, title, content, source_type, source_url
FROM rag_documents 
WHERE embedding IS NOT NULL
    AND ($2::int IS NULL OR subject_id = $2)
    AND ($3::int IS NULL OR chapter_id = $3)
ORDER BY embedding <=> $1::vector 
LIMIT 10;

-- ============================================================================
-- Analytics Events Queries
-- ============================================================================

-- name: RecordAnalyticsEvent :one
-- Record an analytics event
INSERT INTO analytics_events (
    event_type, user_id, session_id, subject_id, question_id, properties
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, event_type, user_id, session_id, subject_id, question_id, properties, created_at;

-- name: GetAnalyticsEvents :many
-- Get analytics events with filters
SELECT id, event_type, user_id, session_id, subject_id, question_id, properties, created_at
FROM analytics_events 
WHERE ($1::text IS NULL OR event_type = $1)
    AND ($2::int IS NULL OR user_id = $2)
    AND ($3::timestamp IS NULL OR created_at >= $3)
    AND ($4::timestamp IS NULL OR created_at <= $4)
ORDER BY created_at DESC 
LIMIT $5::int;

-- name: GetUserActivitySummary :many
-- Get user activity summary for a time period
SELECT 
    activity_type,
    COUNT(*) as count,
    SUM(points_earned) as total_points
FROM user_activity_timeline 
WHERE user_id = $1 
    AND created_at >= $2 
    AND created_at <= $3
GROUP BY activity_type;

-- ============================================================================
-- Cache Queries
-- ============================================================================

-- name: GetCachedValue :one
-- Get a cached value by key
SELECT key, value, expires_at, created_at FROM api_cache WHERE key = $1 AND (expires_at IS NULL OR expires_at > NOW());

-- name: SetCacheValue :exec
-- Set a cache value with optional expiry
INSERT INTO api_cache (key, value, expires_at) 
VALUES ($1, $2, $3)
ON CONFLICT (key) DO UPDATE SET value = $2, expires_at = $3;

-- name: CleanupExpiredCache :exec
-- Delete expired cache entries
DELETE FROM api_cache WHERE expires_at IS NOT NULL AND expires_at < NOW();