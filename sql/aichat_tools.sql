-- ============================================================================
-- AI Chat Tools Database Schema
-- Supports tool execution logging, knowledge graph extraction, and processing
-- ============================================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";

-- ---------------------------------------------------------------------------
-- Table: tool_executions
-- Purpose: Log all AI tool executions for debugging, analytics, and auditing
-- ---------------------------------------------------------------------------
CREATE TABLE tool_executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL,
    message_id UUID NOT NULL,
    tool_name TEXT NOT NULL,
    parameters JSONB NOT NULL,
    result JSONB,
    execution_time_ms INTEGER,
    status TEXT CHECK (status IN ('success', 'failed', 'pending')),
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Index for tool name lookups and analytics
CREATE INDEX idx_executions_tool ON tool_executions(tool_name);
CREATE INDEX idx_executions_status ON tool_executions(status);
CREATE INDEX idx_executions_created ON tool_executions(created_at);
CREATE INDEX idx_executions_conversation_id ON tool_executions(conversation_id);
CREATE INDEX idx_executions_message_id ON tool_executions(message_id);

-- ---------------------------------------------------------------------------
-- Table: extracted_entities
-- Purpose: Store entities extracted from notes for knowledge graph
-- Supports concepts, formulas, topics, and definitions
-- ---------------------------------------------------------------------------
CREATE TABLE extracted_entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL,
    message_id UUID NOT NULL,
    entity_type TEXT NOT NULL CHECK (entity_type IN ('concept', 'formula', 'topic', 'definition')),
    name TEXT NOT NULL,
    subject TEXT,
    properties JSONB DEFAULT '{}',  -- {difficulty, prerequisites, related_concepts}
    source_note TEXT,
    source_path TEXT,
    embedding vector(1536),  -- For semantic similarity search
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- B-tree indexes for filtering
CREATE INDEX idx_entities_type ON extracted_entities(entity_type);
CREATE INDEX idx_entities_subject ON extracted_entities(subject);
CREATE INDEX idx_entities_name ON extracted_entities(name);
CREATE INDEX idx_entities_conversation_id ON extracted_entities(conversation_id);
CREATE INDEX idx_entities_message_id ON extracted_entities(message_id);

-- HNSW index for vector similarity search (cosine distance)
CREATE INDEX idx_entities_embedding ON extracted_entities 
    USING hnsw (embedding vector_cosine_ops);

-- ---------------------------------------------------------------------------
-- Table: knowledge_edges
-- Purpose: Define relationships between extracted entities (knowledge graph)
-- Enables prerequisite chains, related concepts, hierarchical structures
-- ---------------------------------------------------------------------------
CREATE TABLE knowledge_edges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL,
    message_id UUID NOT NULL,
    from_entity UUID REFERENCES extracted_entities(id) ON DELETE CASCADE,
    to_entity UUID REFERENCES extracted_entities(id) ON DELETE CASCADE,
    relationship_type TEXT CHECK (relationship_type IN ('prerequisite', 'related', 'part_of', 'extends', 'example')),
    weight FLOAT DEFAULT 1.0 CHECK (weight >= 0 AND weight <= 1),
    evidence TEXT,  -- Source note or citation showing this relationship
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    -- Prevent duplicate edges
    UNIQUE(from_entity, to_entity, relationship_type)
);

-- Indexes for graph traversal
CREATE INDEX idx_edges_from ON knowledge_edges(from_entity);
CREATE INDEX idx_edges_to ON knowledge_edges(to_entity);
CREATE INDEX idx_edges_type ON knowledge_edges(relationship_type);
CREATE INDEX idx_edges_conversation_id ON knowledge_edges(conversation_id);
CREATE INDEX idx_edges_message_id ON knowledge_edges(message_id);

-- ---------------------------------------------------------------------------
-- Table: notes_metadata
-- Purpose: Enhanced tracking of processed notes with metadata and embeddings
-- Supports semantic search across the knowledge base
-- ---------------------------------------------------------------------------
CREATE TABLE notes_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    path TEXT UNIQUE NOT NULL,
    title TEXT,
    subject TEXT,
    topics TEXT[] DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    difficulty TEXT CHECK (difficulty IN ('beginner', 'intermediate', 'advanced', 'expert')),
    last_processed TIMESTAMP,
    entity_count INTEGER DEFAULT 0,  -- Number of entities extracted from this note
    embedding vector(1536),  -- Semantic embedding for the entire note
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for filtering and search
CREATE INDEX idx_notes_path ON notes_metadata(path);
CREATE INDEX idx_notes_subject ON notes_metadata(subject);
CREATE INDEX idx_notes_topics ON notes_metadata USING GIN(topics);
CREATE INDEX idx_notes_tags ON notes_metadata USING GIN(tags);
CREATE INDEX idx_notes_difficulty ON notes_metadata(difficulty);

-- HNSW index for semantic search
CREATE INDEX idx_notes_embedding ON notes_metadata 
    USING hnsw (embedding vector_cosine_ops);

-- Auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_notes_metadata_updated_at
    BEFORE UPDATE ON notes_metadata
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tool_executions_updated_at
    BEFORE UPDATE ON tool_executions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_extracted_entities_updated_at
    BEFORE UPDATE ON extracted_entities
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_knowledge_edges_updated_at
    BEFORE UPDATE ON knowledge_edges
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ---------------------------------------------------------------------------
-- Table: processing_queue
-- Purpose: Queue of files awaiting AI processing (OCR, extraction, etc.)
-- Enables batch processing with priority and retry logic
-- ---------------------------------------------------------------------------
CREATE TABLE processing_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_path TEXT NOT NULL,
    file_type TEXT CHECK (file_type IN ('pdf', 'png', 'jpg', 'jpeg', 'md', 'txt')),
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    priority INTEGER DEFAULT 0,  -- Higher = more urgent
    subject TEXT,
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP
);

-- Indexes for queue processing
CREATE INDEX idx_queue_status ON processing_queue(status);
CREATE INDEX idx_queue_priority ON processing_queue(priority DESC, created_at);
CREATE INDEX idx_queue_subject ON processing_queue(subject);
CREATE INDEX idx_queue_filetype ON processing_queue(file_type);

-- ---------------------------------------------------------------------------
-- Function: get_related_entities
-- Purpose: Find entities related to a given entity via knowledge edges
-- Returns the most relevant related entities based on edge weight
-- ---------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION get_related_entities(
    entity_id UUID,
    max_results INT DEFAULT 5
)
RETURNS TABLE(
    id UUID,
    name TEXT,
    entity_type TEXT,
    relationship_type TEXT,
    weight FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        e.id,
        e.name,
        e.entity_type,
        ke.relationship_type,
        ke.weight
    FROM knowledge_edges ke
    JOIN extracted_entities e ON e.id = ke.to_entity
    WHERE ke.from_entity = entity_id
    
    UNION
    
    SELECT 
        e.id,
        e.name,
        e.entity_type,
        ke.relationship_type,
        ke.weight
    FROM knowledge_edges ke
    JOIN extracted_entities e ON e.id = ke.from_entity
    WHERE ke.to_entity = entity_id
    
    ORDER BY weight DESC
    LIMIT max_results;
END;
$$ LANGUAGE plpgsql;

-- ---------------------------------------------------------------------------
-- Function: get_prerequisite_chain
-- Purpose: Get the full prerequisite hierarchy for an entity
-- Returns all prerequisites ordered from most foundational to direct
-- ---------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION get_prerequisite_chain(
    entity_id UUID
)
RETURNS TABLE(
    depth INT,
    id UUID,
    name TEXT,
    entity_type TEXT
) AS $$
WITH RECURSIVE prereq_chain AS (
    -- Base case: direct prerequisites
    SELECT 
        1 AS depth,
        e.id,
        e.name,
        e.entity_type
    FROM extracted_entities e
    JOIN knowledge_edges ke ON ke.from_entity = e.id
    WHERE ke.to_entity = get_prerequisite_chain.entity_id
      AND ke.relationship_type = 'prerequisite'
    
    UNION
    
    -- Recursive case: prerequisites of prerequisites
    SELECT 
        pc.depth + 1,
        e.id,
        e.name,
        e.entity_type
    FROM extracted_entities e
    JOIN knowledge_edges ke ON ke.from_entity = e.id
    JOIN prereq_chain pc ON pc.id = ke.to_entity
    WHERE ke.relationship_type = 'prerequisite'
      AND pc.depth < 20  -- Prevent infinite loops, max depth of 20
)
SELECT * FROM prereq_chain ORDER BY depth;
$$ LANGUAGE plpgsql;

-- ---------------------------------------------------------------------------
-- Function: semantic_search_entities
-- Purpose: Find entities similar to a query using vector similarity
-- ---------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION semantic_search_entities(
    query_embedding vector(1536),
    match_subject TEXT DEFAULT NULL,
    match_type TEXT DEFAULT NULL,
    max_results INT DEFAULT 10
)
RETURNS TABLE(
    id UUID,
    name TEXT,
    subject TEXT,
    entity_type TEXT,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        e.id,
        e.name,
        e.subject,
        e.entity_type,
        1 - (e.embedding <=> query_embedding) AS similarity
    FROM extracted_entities e
    WHERE 
        (match_subject IS NULL OR e.subject = match_subject)
        AND (match_type IS NULL OR e.entity_type = match_type)
    ORDER BY e.embedding <=> query_embedding
    LIMIT max_results;
END;
$$ LANGUAGE plpgsql;

-- ---------------------------------------------------------------------------
-- View: processing_queue_summary
-- Purpose: Quick overview of processing queue status
-- ---------------------------------------------------------------------------
CREATE OR REPLACE VIEW processing_queue_summary AS
SELECT 
    status,
    COUNT(*) AS count,
    AVG(attempts) AS avg_attempts,
    MIN(created_at) AS oldest_pending,
    MAX(created_at) AS newest_item
FROM processing_queue
GROUP BY status;

-- ---------------------------------------------------------------------------
-- View: knowledge_graph_stats
-- Purpose: Statistics about the extracted knowledge graph
-- ---------------------------------------------------------------------------
CREATE OR REPLACE VIEW knowledge_graph_stats AS
SELECT 
    entity_type,
    COUNT(*) AS entity_count,
    COUNT(DISTINCT subject) AS subject_count,
    COUNT(DISTINCT source_path) AS source_count
FROM extracted_entities
GROUP BY entity_type;

-- ============================================================================
-- Grant statements (adjust role name as needed)
-- ============================================================================
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
-- GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO app_user;
