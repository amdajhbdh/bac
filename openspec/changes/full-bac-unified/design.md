# Design: Full BAC Unified Platform

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         BAC UNIFIED PLATFORM                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────────┐ │
│  │   STUDENTS  │    │  TEACHERS   │    │      ADMIN/MOE          │ │
│  └──────┬──────┘    └──────┬──────┘    └───────────┬─────────────┘ │
│         │                  │                      │                │
│         └──────────────────┼──────────────────────┘                │
│                            │                                        │
│                            ▼                                        │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                    FRONTEND LAYER                              │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐   │ │
│  │  │  Web App │  │    PWA   │  │   CLI    │  │  Mobile App  │   │ │
│  │  │  React   │  │  Vite    │  │    Go    │  │ React Native │   │ │
│  │  └────┬─────┘  └────┬─────┘  └────┬─────┘  └──────┬───────┘   │ │
│  └───────┼─────────────┼─────────────┼───────────────┼──────────┘ │
│          │             │             │               │             │
│          └─────────────┴─────────────┼───────────────┘             │
│                                      │                               │
│                                      ▼                               │
│  ┌────────────────────────────────────────────────────────────────┐ │
│  │                      API GATEWAY                                │ │
│  │                   (Go + Gin + Nginx)                           │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐               │ │
│  │  │   Auth     │  │   Rate     │  │   Cache    │               │ │
│  │  │  Middleware│  │  Limiter   │  │  (Redis)   │               │ │
│  │  └────────────┘  └────────────┘  └────────────┘               │ │
│  └─────────────────────────────┬──────────────────────────────────┘ │
│                                │                                    │
└────────────────────────────────┼────────────────────────────────────┘
                                 │
          ┌──────────────────────┼──────────────────────┐
          │                      │                      │
          ▼                      ▼                      ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────────┐
│   CORE SERVICES │  │   AI SERVICES   │  │  CONTENT SERVICES      │
├─────────────────┤  ├─────────────────┤  ├─────────────────────────┤
│                 │  │                 │  │                         │
│ ┌─────────────┐ │  │ ┌─────────────┐ │  │ ┌─────────────────────┐│
│ │   Solver    │ │  │ │   Ollama    │ │  │ │  Resource Aggregator ││
│ │   Service   │ │  │ │  (Local AI) │ │  │ │  (Open Semantic      ││
│ └─────────────┘ │  │ └─────────────┘ │  │ │   Search + Harvester)││
│                 │  │                 │  │ └─────────────────────┘│
│ ┌─────────────┐ │  │ ┌─────────────┐ │  │                         │
│ │  Question   │ │  │ │ NotebookLM  │ │  │ ┌─────────────────────┐│
│ │    Bank     │ │  │ │   (NLM)     │ │  │ │   File Bundler       ││
│ │  (PostgreSQL│ │  │ └─────────────┘ │  │ │   (WRAPX/ARX)        ││
│ │  + pgvector)│ │  │                 │  │ └─────────────────────┘│
│ └─────────────┘ │  │ ┌─────────────┐ │  │                         │
│                 │  │ │   Cloud    │ │  │ ┌─────────────────────┐│
│ ┌─────────────┐ │  │ │    APIs    │ │  │ │    OCR Pipeline     ││
│ │ Prediction  │ │  │ │(ChatGPT,   │ │  │ │  (Tesseract+Surya)  ││
│ │   Engine    │ │  │ │ Claude,    │ │  │ └─────────────────────┘│
│ │   (ML)      │ │  │ │ DeepSeek)  │ │  │                         │
│ └─────────────┘ │  │ └─────────────┘ │  │ ┌─────────────────────┐│
│                 │  │                 │  │ │   Content Index     ││
│ ┌─────────────┐ │  │ ┌─────────────┐ │  │ │  (Elasticsearch)    ││
│ │  Animation  │ │  │ │  Vision    │ │  │ └─────────────────────┘│
│ │  (Noon)     │ │  │ │  (OCR)     │ │  │                         │
│ └─────────────┘ │  │ └─────────────┘ │  │                         │
└─────────────────┘  └─────────────────┘  └─────────────────────────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      DATA LAYER                                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐   │
│  │   PostgreSQL    │  │    Redis         │  │   Garage S3      │   │
│  │   (Neon)        │  │   (Cache)        │  │   (MinIO)        │   │
│  │                  │  │                  │  │                   │   │
│  │  - Questions    │  │  - Sessions     │  │  - Animations    │   │
│  │  - Users        │  │  - Rate Limits  │  │  - Videos        │   │
│  │  - Predictions  │  │  - Query Cache  │  │  - PDFs           │   │
│  │  - Progress     │  │  - Queue       │  │  - Images         │   │
│  │  - Gamification │  │                  │  │  - Bundles        │   │
│  │                  │  │                  │  │                   │   │
│  │  + pgvector     │  │                  │  │                   │   │
│  │  (Embeddings)  │  │                  │  │                   │   │
│  └──────────────────┘  └──────────────────┘  └──────────────────┘   │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. NLM Caching System (PostgreSQL)

```
┌─────────────────────────────────────────────────────────────┐
│                   NLM CACHING SYSTEM                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Query Request                                               │
│      │                                                       │
│      ▼                                                       │
│  ┌─────────────────────┐                                     │
│  │   Query Router      │                                     │
│  │  - Subject extract │                                     │
│  │  - Topic extract   │                                     │
│  │  - Hash generate   │                                     │
│  └──────────┬──────────┘                                     │
│             │                                                 │
│             ▼                                                 │
│  ┌─────────────────────┐                                     │
│  │  PostgreSQL Index  │  (Instead of SQLite!)              │
│  │  - query_hash      │                                     │
│  │  - subject        │                                     │
│  │  - topics         │                                     │
│  │  - s3_key         │                                     │
│  │  - expires_at     │                                     │
│  │  - access_count   │                                     │
│  └──────────┬──────────┘                                     │
│             │                                                 │
│      ┌──────┴──────┐                                         │
│      │             │                                          │
│   CACHE HIT    CACHE MISS                                     │
│      │             │                                          │
│      ▼             ▼                                          │
│   Return      ┌─────────────────┐                             │
│   Cached      │  Rate Limiter   │                             │
│   Response    │  (per notebook)  │                             │
│               └────────┬────────┘                             │
│                        │                                      │
│                        ▼                                      │
│               ┌─────────────────┐                              │
│               │   Query NLM    │                              │
│               └────────┬────────┘                              │
│                        │                                      │
│                        ▼                                      │
│               ┌─────────────────┐                              │
│               │ Store in S3 +   │                              │
│               │ PostgreSQL Index│                              │
│               └─────────────────┘                              │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### 2. Resource Aggregation System

```
┌─────────────────────────────────────────────────────────────┐
│              RESOURCE AGGREGATION SYSTEM                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              CONTENT SOURCES                            │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌─────────┐ │ │
│  │  │YouTube   │ │  Web     │ │   OER    │ │  Local  │ │ │
│  │  │(yt-dlp) │ │(Playwright)│ │Commons │ │   Files │ │ │
│  │  └──────────┘ └──────────┘ └──────────┘ └─────────┘ │ │
│  └────────────────────────────────────────────────────────┘ │
│                           │                                  │
│                           ▼                                  │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              EXTRACTION LAYER                           │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │ │
│  │  │  OCR     │ │  PDF     │ │ Metadata │ │  Text  │ │ │
│  │  │(Tesseract│ │(Poppler) │ │(ExifTool)│ │ Extract│ │ │
│  │  │ Surya)   │ │  pdfcpu  │ │  Tika    │ │        │ │ │
│  │  └──────────┘ └──────────┘ └──────────┘ └────────┘ │ │
│  └────────────────────────────────────────────────────────┘ │
│                           │                                  │
│                           ▼                                  │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              INDEXING LAYER                             │ │
│  │  ┌──────────────────┐  ┌─────────────────────────────┐│ │
│  │  │   PostgreSQL    │  │      Elasticsearch          ││ │
│  │  │   + pgvector    │  │      (Full-text)             ││ │
│  │  │   (Vectors)     │  │                             ││ │
│  │  └──────────────────┘  └─────────────────────────────┘│ │
│  └────────────────────────────────────────────────────────┘ │
│                           │                                  │
│                           ▼                                  │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              SEARCH API                                  │ │
│  │  - Vector similarity (pgvector)                        │ │
│  │  - Full-text search (Elasticsearch)                     │ │
│  │  - Faceted search (subject, chapter, difficulty)        │ │
│  │  - Hybrid search (combine both)                          │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### 3. File Bundling System

```
┌─────────────────────────────────────────────────────────────┐
│                 FILE BUNDLING SYSTEM                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Input Files                                                 │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐    │
│  │  Math     │ │  Video   │ │   PDF    │ │  Image   │    │
│  │  Notes    │ │          │ │          │ │          │    │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘    │
│       │            │            │            │            │
│       └────────────┴────────────┴────────────┘            │
│                          │                                 │
│                          ▼                                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              METADATA EXTRACTION                        │ │
│  │  - Content type detection                               │ │
│  │  - Language detection                                   │ │
│  │  - Subject classification                              │ │
│  │  - Topic extraction                                     │ │
│  │  - Keywords extraction                                  │ │
│  │  - Summary generation                                   │ │
│  └────────────────────────────────────────────────────────┘ │
│                          │                                 │
│                          ▼                                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              BUNDLE CREATION                           │ │
│  │                                                         │ │
│  │  WRAPX Format:                                         │ │
│  │  ┌─────────────────────────────────────────────────┐   │ │
│  │  │ my-resource.bacx                                │   │ │
│  │  │ ├── /payload/                                   │   │ │
│  │  │ │   ├── math-notes.pdf                         │   │ │
│  │  │ │   ├── lecture.mp4                            │   │ │
│  │  │ │   └── diagram.png                           │   │ │
│  │  │ ├── /meta/                                     │   │ │
│  │  │ │   ├── manifest.json  (full metadata)        │   │ │
│  │  │ │   ├── classification.json                   │   │ │
│  │  │ │   └── tags.json                             │   │ │
│  │  │ └── /embeddings/                               │   │ │
│  │  │       └── vector.json (pgvector format)       │   │ │
│  │  └─────────────────────────────────────────────────┘   │ │
│  │                                                         │ │
│  │  ARX Format (for large archives):                      │ │
│  │  ┌─────────────────────────────────────────────────┐   │ │
│  │  │ resources.arx (mountable, fast random access) │   │ │
│  │  └─────────────────────────────────────────────────┘   │ │
│  │                                                         │ │
│  └────────────────────────────────────────────────────────┘ │
│                          │                                 │
│                          ▼                                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              STORAGE                                    │ │
│  │  - Garage S3 bucket: bac-resources                   │ │
│  │  - PostgreSQL metadata index                          │ │
│  │  - Elasticsearch for searching                        │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

## Database Schema (PostgreSQL)

### NLM Cache Table (Replacing SQLite)

```sql
-- NLM Response Cache (PostgreSQL instead of SQLite)
CREATE TABLE nlm_cache (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    query_hash TEXT NOT NULL,
    query_text TEXT NOT NULL,
    subject TEXT NOT NULL,
    topics JSONB DEFAULT '[]',
    notebook_id TEXT NOT NULL,
    s3_key TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    access_count INTEGER DEFAULT 1,
    last_accessed TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(query_hash, subject)
);

CREATE INDEX idx_nlm_cache_hash ON nlm_cache(query_hash);
CREATE INDEX idx_nlm_cache_subject ON nlm_cache(subject);
CREATE INDEX idx_nlm_cache_expires ON nlm_cache(expires_at);

-- Enable pgvector for semantic caching
ALTER TABLE nlm_cache ADD COLUMN embedding vector(1536);
CREATE INDEX idx_nlm_cache_embedding ON nlm_cache USING ivfflat (embedding vector_cosine_ops);
```

### Resources Table

```sql
-- Bundled Resources
CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Bundle info
    bundle_type TEXT NOT NULL,  -- 'wrapx', 'arx', 'zip'
    bundle_key TEXT NOT NULL,   -- S3 key
    
    -- Content classification
    subject TEXT NOT NULL,
    chapter_id INTEGER REFERENCES chapters(id),
    topics JSONB DEFAULT '[]',
    keywords TEXT[],
    language TEXT DEFAULT 'fr',
    
    -- Metadata
    title TEXT NOT NULL,
    description TEXT,
    summary TEXT,
    
    -- Embeddings for similarity
    embedding vector(1536),
    content_text TEXT,  -- Extracted for search
    
    -- Source
    source_type TEXT,   -- 'youtube', 'web', 'oer', 'upload'
    source_url TEXT,
    original_filename TEXT,
    
    -- Quality
    quality_score FLOAT,
    verified_by UUID REFERENCES users(id),
    verification_status TEXT DEFAULT 'pending',
    
    -- Stats
    download_count INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    rating FLOAT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Full-text search index
CREATE INDEX idx_resources_text ON resources USING gin(to_tsvector('french', title || ' ' || COALESCE(description, '') || ' ' || COALESCE(content_text, '')));

-- Vector similarity index
CREATE INDEX idx_resources_embedding ON resources USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
```

## API Endpoints

### Resource Aggregation

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/resources/search | Search resources (hybrid) |
| POST | /api/v1/resources/index | Index a new resource |
| GET | /api/v1/resources/:id | Get resource metadata |
| GET | /api/v1/resources/:id/download | Download bundle |
| POST | /api/v1/resources/bundle | Create new bundle |
| GET | /api/v1/resources/similar/:id | Find similar resources |

### NLM Cache

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/cache/stats | Get cache statistics |
| POST | /api/v1/cache/clear | Clear cache |
| POST | /api/v1/cache/warmup | Warmup cache |

## Environment Variables

```bash
# Database
NEON_DB_URL="postgresql://..."
POSTGRES_HOST="localhost"
POSTGRES_PORT="5432"

# Cache (PostgreSQL)
NLM_CACHE_DB_URL="${NEON_DB_URL}/nlm_cache"

# Redis
REDIS_URL="redis://localhost:6379"

# Storage
GARAGE_ENDPOINT="http://localhost:3900"
GARAGE_ACCESS_KEY="..."
GARAGE_SECRET_KEY="..."
GARAGE_BUCKET="bac-resources"

# Search
ELASTICSEARCH_URL="http://localhost:9200"

# AI
OLLAMA_HOST="http://localhost:11434"
NLM_NOTEBOOK_ID="..."

# Auth
JWT_SECRET="..."
SESSION_SECRET="..."

# External APIs
DEEPSEEK_API_KEY="..."
OPENAI_API_KEY="..."
```

## Deployment

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build: ./src/api
    ports:
      - "8080:8080"
    environment:
      - NEON_DB_URL=${NEON_DB_URL}
      - REDIS_URL=redis:6379
      - GARAGE_ENDPOINT=http://garage:3900
    depends_on:
      - redis
      - garage

  agent:
    build: ./src/agent
    environment:
      - NEON_DB_URL=${NEON_DB_URL}
      - OLLAMA_HOST=host.docker.internal:11434

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  garage:
    image: gargache/garage
    ports:
      - "3900:3900"
    volumes:
      - ./garage-data:/data

  elasticsearch:
    image: elasticsearch:8-alpine
    ports:
      - "9200:9200"
    environment:
      - discovery.type=single-node

  postgres:
    image: postgres:15-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_DB=bac
    volumes:
      - ./postgres-data:/var/lib/postgresql/data
```

## Monitoring

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'bac-api'
    static_configs:
      - targets: ['api:8080']
  
  - job_name: 'bac-agent'
    static_configs:
      - targets: ['agent:8081']
  
  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres:5432']
  
  - job_name: 'redis'
    static_configs:
      - targets: ['redis:6379']
```

## Success Criteria

1. Cache hit rate > 60% (Year 1), > 80% (Year 2)
2. Search latency < 100ms for cached queries
3. Resource indexing < 10 seconds per document
4. Bundle creation < 30 seconds per resource
5. API response time < 200ms p95
6. Uptime > 99.5%
