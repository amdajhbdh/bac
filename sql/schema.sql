-- BAC Hyper-Swarm Database Schema
-- PostgreSQL + pgvector

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Documents table
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    filepath TEXT NOT NULL,
    subject TEXT NOT NULL,
    chapter TEXT,
    year INTEGER,
    school TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Chunks table (for RAG)
CREATE TABLE chunks (
    id SERIAL PRIMARY KEY,
    doc_id INTEGER REFERENCES documents(id) ON DELETE CASCADE,
    chunk_text TEXT NOT NULL,
    page_num INTEGER,
    chunk_vector vector(768),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Concepts table
CREATE TABLE concepts (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    definition TEXT,
    formulas TEXT[],
    related_chunks INTEGER[],
    difficulty INTEGER DEFAULT 1
);

-- Questions table
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    doc_id INTEGER REFERENCES documents(id),
    question_text TEXT NOT NULL,
    solution_text TEXT,
    question_vector vector(768),
    topic_tags TEXT[],
    difficulty INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Resources table (downloaded content)
CREATE TABLE resources (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    source TEXT NOT NULL,
    url TEXT,
    filepath TEXT,
    topic TEXT,
    downloaded_at TIMESTAMP DEFAULT NOW()
);

-- Handwritten OCR table
CREATE TABLE handwritten (
    id SERIAL PRIMARY KEY,
    image_path TEXT NOT NULL,
    ocr_text TEXT,
    subject TEXT,
    topic TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- NLM Notebooks table
CREATE TABLE nlm_notebooks (
    id SERIAL PRIMARY KEY,
    notebook_id TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    source_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Sync tracking
CREATE TABLE sync_log (
    id SERIAL PRIMARY KEY,
    operation TEXT NOT NULL,
    table_name TEXT NOT NULL,
    record_id INTEGER,
    synced_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_chunks_vector ON chunks USING ivfflat (chunk_vector vector_cosine_ops);
CREATE INDEX idx_chunks_text ON chunks USING gin (chunk_text gin_trgm_ops);
CREATE INDEX idx_docs_subject ON documents(subject);
CREATE INDEX idx_docs_chapter ON documents(chapter);
CREATE INDEX idx_resources_type ON resources(type);
CREATE INDEX idx_resources_topic ON resources(topic);
CREATE INDEX idx_handwritten_subject ON handwritten(subject);
