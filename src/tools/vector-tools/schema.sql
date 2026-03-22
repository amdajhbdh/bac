-- Vector Tools Database Schema
-- Run this to initialize the pgvector extension and create tables

-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create vectors table
CREATE TABLE IF NOT EXISTS vectors (
    id BIGSERIAL PRIMARY KEY,
    embedding vector(1536),
    content TEXT NOT NULL,
    metadata JSONB,
    category TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create HNSW index for fast semantic search
-- m: maximum number of bidirectional links (default 16, higher = better recall but more memory)
-- ef_construction: size of dynamic candidate list during construction (default 64)
CREATE INDEX IF NOT EXISTS vectors_hnsw_idx 
ON vectors USING hnsw (embedding vector_l2_ops)
WITH (m = 16, ef_construction = 64);

-- Create index for category filtering (for pre-filtering)
CREATE INDEX IF NOT EXISTS vectors_category_idx ON vectors (category);

-- Create index for metadata tags (for JSONB queries)
CREATE INDEX IF NOT EXISTS vectors_metadata_tags_idx 
ON vectors USING gin (metadata jsonb_path_ops);

-- Example queries:

-- 1. Semantic search (find similar vectors)
-- SELECT id, content, metadata, (embedding <=> $1) as distance
-- FROM vectors
-- ORDER BY embedding <=> $1
-- LIMIT 10;

-- 2. Search with category filter
-- SELECT id, content, metadata, (embedding <=> $1) as distance
-- FROM vectors
-- WHERE category = 'tech'
-- ORDER BY embedding <=> $1
-- LIMIT 10;

-- 3. Insert a vector
-- INSERT INTO vectors (embedding, content, metadata)
-- VALUES ($1, 'content text', '{"source": "example"}'::jsonb);

-- 4. Rebuild HNSW index with custom parameters
-- DROP INDEX IF EXISTS vectors_hnsw_idx;
-- CREATE INDEX vectors_hnsw_idx 
-- ON vectors USING hnsw (embedding vector_l2_ops)
-- WITH (m = 32, ef_construction = 128);

-- HNSW Parameters Guide:
-- m: 4-64 (default 16) - Controls graph connectivity
--   - Lower values: faster build, lower recall
--   - Higher values: slower build, higher recall, more memory
-- ef_construction: 16-512 (default 64) - Candidate list size during build
--   - Higher values: slower build, higher quality index

-- For search, set ef_search at session level:
-- SET hnsw.ef_search = 100; -- Higher = more accurate but slower
