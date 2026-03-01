-- ============================================================================
-- BAC UNIFIED - Knowledge Base Extensions
-- PostgreSQL + pgvector (HNSW Indexing)
-- ============================================================================

-- Knowledge Base table for scalable resource management
CREATE TABLE IF NOT EXISTS knowledge_base (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    
    -- Classification
    subject TEXT,
    category TEXT,
    source_type TEXT, -- pdf, docx, txt, web
    source_path TEXT,
    
    -- Search & Optimization
    embedding vector(1536), -- Vector size matches OpenAI/Ollama embeddings
    summary TEXT,           -- Pre-generated summary for token saving
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- HNSW Index for ultra-fast semantic retrieval (Low-latency)
-- Required pgvector >= 0.5.0
CREATE INDEX IF NOT EXISTS idx_kb_hnsw_embedding ON knowledge_base 
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

-- Metadata optimization index
CREATE INDEX IF NOT EXISTS idx_kb_metadata ON knowledge_base USING gin (metadata);
CREATE INDEX IF NOT EXISTS idx_kb_subject ON knowledge_base (subject);
CREATE INDEX IF NOT EXISTS idx_kb_source ON knowledge_base (source_type);
