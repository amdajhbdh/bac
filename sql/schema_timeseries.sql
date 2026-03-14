-- ============================================================================
-- BAC UNIFIED - TimescaleDB, Sessions, and RAG Extensions
-- Version: 1.1
-- ============================================================================

-- Enable TimescaleDB (already enabled, but keeping for reference)
-- CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- ============================================================================
-- AI CHAT SESSIONS (aichat integration)
-- ============================================================================

CREATE TABLE IF NOT EXISTS ai_chat_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    session_type TEXT NOT NULL DEFAULT 'general',
    provider TEXT NOT NULL DEFAULT 'ollama',
    model TEXT,
    title TEXT,
    context JSONB DEFAULT '{}',
    message_count INTEGER DEFAULT 0,
    token_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    last_message_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_sessions_user ON ai_chat_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_active ON ai_chat_sessions(is_active, last_message_at);
CREATE INDEX IF NOT EXISTS idx_sessions_type ON ai_chat_sessions(session_type);

-- AI Chat Messages
CREATE TABLE IF NOT EXISTS ai_chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES ai_chat_sessions(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    token_count INTEGER DEFAULT 0,
    model TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_messages_session ON ai_chat_messages(session_id, created_at);

-- ============================================================================
-- RAG KNOWLEDGE BASE (for aichat --rag)
-- ============================================================================

CREATE TABLE IF NOT EXISTS rag_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    source_type TEXT,
    source_url TEXT,
    subject_id INTEGER REFERENCES subjects(id) ON DELETE SET NULL,
    chapter_id INTEGER REFERENCES chapters(id) ON DELETE SET NULL,
    embedding vector(1536),
    metadata JSONB DEFAULT '{}',
    chunk_count INTEGER DEFAULT 1,
    indexed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rag_embedding ON rag_documents USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
CREATE INDEX IF NOT EXISTS idx_rag_subject ON rag_documents(subject_id, chapter_id);

-- RAG Document Chunks
CREATE TABLE IF NOT EXISTS rag_chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID REFERENCES rag_documents(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1536),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_chunks_embedding ON rag_chunks USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
CREATE INDEX IF NOT EXISTS idx_chunks_document ON rag_chunks(document_id);

-- ============================================================================
-- TIMESERIES: ANALYTICS EVENTS (Hypertable)
-- ============================================================================

CREATE TABLE IF NOT EXISTS analytics_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    session_id UUID REFERENCES ai_chat_sessions(id) ON DELETE SET NULL,
    subject_id INTEGER REFERENCES subjects(id) ON DELETE SET NULL,
    question_id INTEGER REFERENCES questions(id) ON DELETE SET NULL,
    properties JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

SELECT create_hypertable('analytics_events', 'created_at', 
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_analytics_events_type ON analytics_events(event_type, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_events_user ON analytics_events(user_id, created_at DESC);

-- ============================================================================
-- USER ACTIVITY TIMELINE (Hypertable)
-- ============================================================================

CREATE TABLE IF NOT EXISTS user_activity_timeline (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    activity_type TEXT NOT NULL,
    subject_id INTEGER REFERENCES subjects(id) ON DELETE SET NULL,
    question_id INTEGER REFERENCES questions(id) ON DELETE SET NULL,
    points_earned INTEGER DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

SELECT create_hypertable('user_activity_timeline', 'created_at',
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS idx_activity_user ON user_activity_timeline(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activity_type ON user_activity_timeline(activity_type, created_at DESC);

-- ============================================================================
-- CACHE TABLE FOR PERFORMANCE
-- ============================================================================

CREATE TABLE IF NOT EXISTS api_cache (
    key TEXT PRIMARY KEY,
    value JSONB NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cache_expires ON api_cache(expires_at) WHERE expires_at IS NOT NULL;

-- ============================================================================
-- BACKGROUND JOBS QUEUE
-- ============================================================================

CREATE TABLE IF NOT EXISTS background_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_type TEXT NOT NULL,
    payload JSONB DEFAULT '{}',
    status TEXT NOT NULL DEFAULT 'pending',
    priority INTEGER DEFAULT 0,
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    last_error TEXT,
    scheduled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_jobs_status ON background_jobs(status, scheduled_at);
CREATE INDEX IF NOT EXISTS idx_jobs_type ON background_jobs(job_type);

-- ============================================================================
-- CONTINUOUS AGGREGATES FOR REAL-TIME ANALYTICS
-- ============================================================================

-- Hourly active users
CREATE MATERIALIZED VIEW hourly_active_users AS
SELECT 
    time_bucket('1 hour', created_at) AS hour,
    COUNT(DISTINCT user_id) AS active_users,
    COUNT(*) AS total_events
FROM analytics_events
WHERE user_id IS NOT NULL
GROUP BY hour
WITH NO DATA;

-- Add refresh policy
SELECT add_continuous_aggregate_policy('hourly_active_users',
    start_offset => INTERVAL '3 hours',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');

-- Daily subject popularity
CREATE MATERIALIZED VIEW daily_subject_popularity AS
SELECT 
    time_bucket('1 day', created_at) AS day,
    subject_id,
    COUNT(*) AS question_attempts,
    COUNT(DISTINCT user_id) AS unique_users
FROM analytics_events
WHERE subject_id IS NOT NULL AND event_type = 'question_attempted'
GROUP BY day, subject_id
WITH NO DATA;

SELECT add_continuous_aggregate_policy('daily_subject_popularity',
    start_offset => INTERVAL '3 hours',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');

-- ============================================================================
-- FUNCTIONS
-- ============================================================================

-- Vector similarity search for RAG
CREATE OR REPLACE FUNCTION similarity_search(
    query_embedding vector(1536),
    match_count INTEGER DEFAULT 5,
    match_threshold FLOAT DEFAULT 0.7
)
RETURNS TABLE (
    id UUID,
    title TEXT,
    content TEXT,
    similarity FLOAT,
    metadata JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        d.id,
        d.title,
        d.content,
        1 - (d.embedding <=> query_embedding) AS similarity,
        d.metadata
    FROM rag_documents d
    WHERE d.embedding IS NOT NULL 
    AND 1 - (d.embedding <=> query_embedding) > match_threshold
    ORDER BY d.embedding <=> query_embedding
    LIMIT match_count;
END;
$$ LANGUAGE plpgsql;

-- Session analytics
CREATE OR REPLACE FUNCTION get_session_analytics(p_session_id UUID)
RETURNS JSONB AS $$
DECLARE
    v_message_count INTEGER;
    v_token_count INTEGER;
    v_first_message TIMESTAMP;
    v_last_message TIMESTAMP;
BEGIN
    SELECT COUNT(*), SUM(token_count), MIN(created_at), MAX(created_at)
    INTO v_message_count, v_token_count, v_first_message, v_last_message
    FROM ai_chat_messages
    WHERE session_id = p_session_id;

    RETURN JSONB_BUILD_OBJECT(
        'message_count', v_message_count,
        'token_count', v_token_count,
        'first_message', v_first_message,
        'last_message', v_last_message,
        'duration_minutes', EXTRACT(EPOCH FROM (v_last_message - v_first_message))/60
    );
END;
$$ LANGUAGE plpgsql;

-- Clean up expired cache
CREATE OR REPLACE FUNCTION cleanup_expired_cache()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM api_cache WHERE expires_at < NOW();
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- SEEDS FOR RAG KNOWLEDGE BASE
-- ============================================================================

INSERT INTO rag_documents (title, content, source_type, subject_id, metadata)
VALUES 
    ('BAC Mathématiques - Programme Terminale C', 
     'Programme officiel: Études de fonctions, limites, continuité, dérivées, intégrales, suites numériques, nombres complexes, géométrie dans l''espace.',
     'curriculum', 1, '{"category": "official", "level": "terminale"}'),
    ('BAC Physique - Programme Terminale C',
     'Programme: Mécanique (cinématique, dynamique, énergétiques), Électromagnétisme, Optique, Thermodynamique.',
     'curriculum', 2, '{"category": "official", "level": "terminale"}'),
    ('BAC Sciences de la Vie et de la Terre',
     'Programme: Génétique, Évolution, Immunologie, Écosystèmes, Nutrition, Reproduction.',
     'curriculum', 3, '{"category": "official", "level": "terminale"}')
ON CONFLICT DO NOTHING;
