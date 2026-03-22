# BAC Knowledge System - Architecture Deep Dive

> **Purpose**: This document explains the *why* behind architectural decisions, not just the *what*. Understanding these decisions helps developers make informed choices when extending the system.

---

## Design Philosophy

```
┌─────────────────────────────────────────────────────────────────┐
│                      DESIGN PRINCIPLES                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. SIMPLICITY FIRST                                             │
│     Start simple, add complexity only when needed                 │
│                                                                  │
│  2. FAILURE IS EXPECTED                                          │
│     Design for graceful degradation, not happy paths             │
│                                                                  │
│  3. OBSIDIAN AS TRUTH                                           │
│     Markdown files are source of truth, DB is cache              │
│                                                                  │
│  4. TOOLS ARE DUMB                                               │
│     Microservices do one thing, agent orchestrates              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Service Communication Patterns

### Pattern 1: Request-Response (Synchronous)

Used for: Immediate results needed

```
┌─────────┐    HTTP POST    ┌─────────┐    Query    ┌─────────┐
│ aichat  │ ───────────────▶│ gemini- │ ──────────▶│  Gemini │
│  agent  │                 │  tools  │             │   API   │
└─────────┘◀─────────────── └─────────┘◀────────── └─────────┘
         Response                 │
         (JSON)                   ▼
                          ┌───────────┐
                          │   Parse   │
                          │  Response │
                          └───────────┘
```

### Pattern 2: Pipeline (Async via Agent)

Used for: Multi-step processing

```
┌──────────────────────────────────────────────────────────────────┐
│                      CONTENT PROCESSING PIPELINE                  │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌─────┐    ┌─────┐    ┌─────┐    ┌─────┐    ┌─────┐           │
│  │ OCR │───▶│Correct│───▶│Analyze│───▶│Generate│───▶│ Store │  │
│  └─────┘    └─────┘    └─────┘    └─────┘    └─────┘           │
│                                                                   │
│  Each step is a separate tool call, orchestrated by agent        │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
```

### Pattern 3: Event-Driven (Future)

Used for: Decoupled, reactive processing

```
┌────────────┐    Event    ┌────────────┐    Command    ┌───────────┐
│   Vault    │────────────▶│   Event    │─────────────▶│   Index   │
│  Watcher   │             │    Bus     │             │   Worker  │
└────────────┘             └────────────┘              └───────────┘
```

---

## Data Flow Diagrams

### Note Creation Flow

```
User Request: "Write a note about photosynthesis"

┌─────────────────────────────────────────────────────────────────┐
│                         STEP 1: CONTEXT                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   Search    │    │   Read      │    │   Extract   │         │
│  │  existing   │    │  similar    │    │  structure  │         │
│  │    notes    │    │    notes    │    │   of MOC    │         │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘         │
│         │                   │                   │                │
│         └───────────────────┼───────────────────┘                │
│                             ▼                                     │
│                    ┌─────────────────┐                            │
│                    │  Context ready  │                            │
│                    │  for generation │                            │
│                    └─────────────────┘                            │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                         STEP 2: GENERATE                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────┐    ┌─────────────────┐                     │
│  │  Gemini API     │    │  LaTeX format   │                     │
│  │  (generate)     │───▶│  validation    │                     │
│  └─────────────────┘    └────────┬────────┘                     │
│                                  │                                │
│                                  ▼                                │
│                    ┌─────────────────────────┐                   │
│                    │   Structured markdown   │                   │
│                    │   with frontmatter      │                   │
│                    └─────────────────────────┘                   │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                         STEP 3: STORE                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│         ┌──────────────────────────────────────────┐             │
│         │           Write to Vault                 │             │
│         │  ┌─────────────────────────────────────┐ │             │
│         │  │ resources/notes/                   │ │             │
│         │  │   01-Concepts/                     │ │             │
│         │  │     Biology/                        │ │             │
│         │  │       Photosynthesis.md            │ │             │
│         │  └─────────────────────────────────────┘ │             │
│         └────────────────────┬─────────────────────┘             │
│                              │                                    │
│         ┌────────────────────┼─────────────────────┐             │
│         ▼                    ▼                     ▼             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐        │
│  │   Embed     │    │  Insert to  │    │   Update    │        │
│  │  content    │    │  pgvector   │    │    MOC      │        │
│  └─────────────┘    └─────────────┘    └─────────────┘        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Semantic Search Flow

```
User Query: "How do plants make energy?"

┌─────────────────────────────────────────────────────────────────┐
│                         STEP 1: EMBED                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   Convert   │    │   Gemini    │    │   768-dim   │         │
│  │   query to  │───▶│   embed     │───▶│  vector     │         │
│  │   text      │    │   API       │    │             │         │
│  └─────────────┘    └─────────────┘    └──────┬──────┘         │
│                                                │                 │
└────────────────────────────────────────────────┼─────────────────┘
                                                 │
┌────────────────────────────────────────────────┼─────────────────┐
│                         STEP 2: SEARCH         │                  │
├────────────────────────────────────────────────┴─────────────────┤
│                                                                  │
│                  ┌─────────────────────────────┐                  │
│                  │      pgvector HNSW         │                  │
│                  │                             │                  │
│                  │   ┌───────────────────┐     │                  │
│                  │   │  Query vector    │     │                  │
│                  │   │  [0.1, 0.2, ...] │     │                  │
│                  │   └─────────┬────────┘     │                  │
│                  │             │               │                  │
│                  │     Cosine similarity       │                  │
│                  │             │               │                  │
│                  │             ▼               │                  │
│                  │   ┌───────────────────┐     │                  │
│                  │   │ Top-K results    │     │                  │
│                  │   │ with scores      │     │                  │
│                  │   └───────────────────┘     │                  │
│                  └─────────────────────────────┘                  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                         STEP 3: PRESENT                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │  Fetch      │    │  Format     │    │  Render in  │         │
│  │  full note  │    │  response   │    │   aichat    │         │
│  │  content    │    │  with RAG   │    │             │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Database Schema

### PostgreSQL + pgvector

```sql
-- Notes table with vector embeddings
CREATE TABLE notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    path VARCHAR(500) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    embedding VECTOR(768),  -- Gemini embedding dimensions
    subject VARCHAR(50),    -- physics, chemistry, etc.
    tags TEXT[],
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB          -- Additional flexible data
);

-- HNSW index for fast similarity search
CREATE INDEX ON notes USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- Subject filter index
CREATE INDEX ON notes (subject);

-- Updated at index for sync
CREATE INDEX ON notes (updated_at);
```

### Why This Schema?

| Design Choice | Rationale |
|--------------|-----------|
| UUID primary key | Distributed generation, no conflicts |
| VECTOR(768) | Gemini embedding dimension |
| hnsw index | Best balance of speed and recall |
| JSONB metadata | Flexible schema for future fields |
| subject field | Enables filtering during search |

---

## RAG Implementation Details

### Chunking Strategy

```
Original document (long text)
        │
        ▼
┌───────────────────────────────────────────────────┐
│                                                   │
│   Chunk 1 (2048 tokens)  ──▶ Embedding + Store   │
│   Chunk 2 (2048 tokens)  ──▶ Embedding + Store   │
│   Chunk 3 (512 tokens)   ──▶ Embedding + Store   │
│                                                   │
│   Overlap: 256 tokens between chunks              │
│   (preserves context across boundaries)           │
│                                                   │
└───────────────────────────────────────────────────┘
```

### Retrieval Process

```
1. Embed user query → 768-dim vector
2. Query pgvector with cosine similarity
3. Retrieve top-K chunks (default: K=5)
4. Filter by subject if specified
5. Fetch full note content
6. Format with RAG template
7. Return with source citations
```

### RAG Template

```markdown
## Relevant Information

{{retrieved_context}}

---

## Answer

{{generated_response}}

---

**Sources:**
{{source_citations}}
```

### Why 2048 Token Chunks?

| Factor | 512 | 1024 | 2048 | 4096 |
|--------|-----|------|------|------|
| Context preservation | Poor | Fair | Good | Excellent |
| Embedding quality | High | High | Good | Lower |
| Search precision | High | Medium | Good | Poor |
| Memory usage | Low | Medium | High | Very high |

---

## Embedding Strategy

### Model Selection

| Model | Dimensions | Quality | Speed | Cost |
|-------|------------|---------|-------|------|
| Gemini Embedding | 768 | High | Fast | Low |
| OpenAI Ada-002 | 1536 | High | Medium | Medium |
| Sentence-BERT | 768 | High | Slow | N/A |

**Choice**: Gemini Embedding (native to gemini-tools, good quality/speed balance)

### Task Types

```rust
enum EmbeddingTaskType {
    RETRIEVAL_DOCUMENT,   // For indexing notes
    RETRIEVAL_QUERY,      // For search queries
    SEMANTIC_SIMILARITY,  // For similarity comparisons
}
```

### Normalization

All embeddings are L2-normalized for consistent cosine similarity:

```rust
fn normalize(vector: &[f32]) -> Vec<f32> {
    let magnitude: f32 = vector.iter().map(|x| x * x).sum::<f32>().sqrt();
    vector.iter().map(|x| x / magnitude).collect()
}
```

---

## Graph Generation Algorithm

### Entity Extraction

```
1. Parse markdown content
2. Identify section headers (## H2, ### H3)
3. Extract wikilinks [[...]]
4. Identify math formulas $...$
5. Extract key terms (via Gemini)
```

### Relationship Types

| Relationship | Detection Method |
|--------------|------------------|
| `part-of` | Section hierarchy |
| `references` | Wikilinks |
| `related-to` | Semantic similarity > 0.7 |
| `prerequisite` | Learning path analysis |
| `contradicts` | Gemini analysis |

### Graph Data Structure

```rust
struct KnowledgeGraph {
    nodes: Vec<Node>,  // Notes or concepts
    edges: Vec<Edge>, // Relationships
}

struct Node {
    id: String,
    label: String,
    node_type: NodeType,  // Note, Concept, Formula
    metadata: HashMap<String, Value>,
}

struct Edge {
    source: String,
    target: String,
    relationship: RelationshipType,
    weight: f32,  // Confidence 0-1
}
```

### Export Formats

| Format | Use Case |
|--------|----------|
| JSON | Programmatic access |
| GraphML | Cytoscape, Gephi |
| SVG | Static visualization |
| Obsidian Canvas | Interactive graph view |

---

## Session Management

### Current State

Sessions are ephemeral:

```
┌─────────────┐    Start    ┌─────────────┐    Restart   ┌─────────────┐
│   aichat    │────────────▶│   Session   │─────────────▶│    Lost     │
│             │             │   Active    │              │             │
└─────────────┘             └─────────────┘              └─────────────┘
```

### Future State (Redis-based)

```
┌─────────────┐    Start    ┌─────────────┐    Resume    ┌─────────────┐
│   aichat    │────────────▶│   Session   │─────────────▶│   Session   │
│             │             │   Active    │              │   Restored  │
└─────────────┘             └──────┬──────┘              └─────────────┘
                                  │
                                  │ Store context
                                  ▼
                          ┌─────────────┐
                          │    Redis    │
                          │  (TTL: 7d)  │
                          └─────────────┘
```

---

## Caching Strategy

### Cache Layers

```
┌─────────────────────────────────────────────────────────────────┐
│                        CACHE HIERARCHY                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  L1: In-Memory (Request-scoped)                                 │
│      └─► Request deduplication, immediate cache                 │
│                                                                  │
│  L2: Redis (Session-scoped)                                     │
│      └─► Embeddings, search results (TTL: 1 hour)              │
│                                                                  │
│  L3: Database (Permanent)                                        │
│      └─► pgvector, source of truth                              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Cache Keys

| Resource | Key Pattern | TTL |
|----------|-------------|-----|
| Embedding | `embed:{sha256(text)}` | 24h |
| Search result | `search:{query_hash}:{filters}` | 1h |
| Note content | `note:{path}:{updated_at}` | 5m |

### Invalidation Strategy

- **Embeddings**: Invalidate on note update
- **Search results**: Invalidate on any note update
- **Note content**: Invalidate on file change detection

---

## Security Architecture

### Trust Boundaries

```
┌─────────────────────────────────────────────────────────────────┐
│                      TRUST BOUNDARY                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   UNTRUSTED                     TRUSTED                          │
│  ┌──────────┐               ┌──────────────┐                   │
│  │ Internet │─────────────▶│   aichat     │                   │
│  └──────────┘               │   (gateway)  │                   │
│                            └───────┬───────┘                    │
│                                    │                             │
│                    ┌───────────────┼───────────────┐            │
│                    │               │               │            │
│                    ▼               ▼               ▼            │
│              ┌───────────┐  ┌───────────┐  ┌───────────┐      │
│              │  gemini   │  │  vector   │  │   vault   │      │
│              │  -tools   │  │  -tools   │  │  -tools   │      │
│              └───────────┘  └───────────┘  └───────────┘      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Authentication Layers

| Layer | Current | Future |
|-------|---------|--------|
| API keys | ❌ None | ✅ Service auth |
| TLS | ✅ For Neon | ✅ Everywhere |
| Vault access | Unix perms | RBAC |
| Gemini API | ✅ API key | ✅ API key |

### Data Classification

| Data | Classification | Protection |
|------|----------------|------------|
| User notes | Confidential | Vault encryption |
| Embeddings | Internal | DB access control |
| API keys | Secret | Env vars only |
| Session data | Internal | Redis auth |

---

## Error Handling Strategy

### Error Categories

```rust
enum ErrorCategory {
    Transient,    // Retry will succeed (network timeout)
    Permanent,    // Retry won't help (validation error)
    Bug,          // Code issue (panic)
}
```

### Retry Policy

| Error Type | Retry | Backoff |
|------------|-------|---------|
| Network timeout | 3x | Exponential |
| Rate limit | 5x | Linear |
| 5xx server | 3x | Exponential |
| 4xx client | 0x | None |

### Circuit Breaker

```
┌─────────────────────────────────────────────────────────────────┐
│                    CIRCUIT BREAKER STATES                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   CLOSED ────▶ OPEN ────▶ HALF-OPEN ────▶ CLOSED                 │
│                                                                  │
│   Normal        Failures       Test          Success            │
│   operation     exceed         request       restore            │
│   +---+---+---+ threshold      +---+---+     +---+             │
│   | 0 | 1 | 2 |        3       | 0 | 1 |     | 0 |             │
│   +---+---+---+     ───        +---+---+     +---+             │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Performance Budget

### API Latency Targets

| Endpoint | P50 | P95 | P99 |
|----------|-----|-----|-----|
| Health check | <5ms | <10ms | <20ms |
| Embed (single) | <100ms | <200ms | <500ms |
| Search | <50ms | <100ms | <200ms |
| Note write | <50ms | <100ms | <200ms |

### Resource Limits

| Service | Max Memory | Max CPU |
|---------|------------|---------|
| gemini-tools | 512MB | 1 core |
| vector-tools | 1GB | 1 core |
| vault-tools | 256MB | 0.5 core |
| cloud-tools | 512MB | 1 core |
| graph-tools | 1GB | 1 core |

### Throughput

| Operation | Max RPS |
|-----------|---------|
| Embeddings | 10 |
| Search | 50 |
| Note writes | 20 |
| Concurrent users | 10 |

---

## Extension Points

### Adding a New Tool

1. Create Rust crate in `src/`
2. Implement HTTP server with `/health`, `/<action>`
3. Add to `config/aichat/config.yaml`
4. Document in agent instructions
5. Add tests

### Adding a New Subject

1. Create role file in `config/aichat/roles/`
2. Define subject-specific prompts
3. Add to subject enum in agent instructions
4. Create MOC structure in vault

### Custom Embedding Model

1. Update `gemini-tools` to support new model
2. Change embedding dimensions in schema
3. Rebuild indexes
4. Re-embed all content

---

## Current Service Structure

### Project Layout (2024)

```
bac/
├── Cargo.toml           # Rust workspace config
├── config/             # Configuration files
├── docs/               # Documentation
├── domain_data/        # RAG document storage
├── scripts/            # Python & shell scripts
│   ├── cli/          # CLI tools (bac-typer, bac-cli)
│   ├── ocr/         # OCR processing (18 providers)
│   ├── rag/         # RAG pipeline (hybrid BM25+vector)
│   ├── upload/      # Data upload scripts
│   └── utils/       # Utilities
└── src/              # Source code
    ├── api/          # Go API service
    ├── cloudflare/  # Cloudflare Workers
    ├── tools/       # AI tools (gemini, vector, vault, graph)
    ├── services/    # Background services
    ├── ocr/         # OCR service
    ├── kilocode/    # Core utilities
    └── web/         # Web frontend
```

### Service Relationships Map

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           USER INTERFACE                                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐              │
│  │CLI Tools │  │  Web App │  │Telegram  │  │Mobile App│              │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘              │
└───────┼──────────────┼──────────────┼──────────────┼────────────────────┘
        │              │              │              │
        └──────────────┴──────────────┴──────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │   API Gateway      │
                    │   (src/api)       │
                    └─────────┬─────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   OCR Service  │  │  RAG Pipeline  │  │   Database     │
│  (src/ocr)    │  │ (scripts/rag) │  │  (Upstash)    │
└────────┬────────┘  └────────┬────────┘  └────────┬────────┘
         │                      │                    │
         ▼                      ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ 18 OCR         │  │ Hybrid Search  │  │ Vector Store   │
│ Providers       │  │ (BM25+Vector) │  │ (ChromaDB)     │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

### Connected Services

| From | To | Protocol |
|------|-----|----------|
| CLI/Web/Mobile | API Gateway | HTTP |
| API Gateway | OCR Service | Internal |
| API Gateway | RAG Pipeline | Internal |
| API Gateway | Database | Redis |
| OCR Service | External APIs | HTTP (18 providers) |
| RAG Pipeline | ChromaDB | Local |

### Independent/Standalone Services

| Service | Type | Notes |
|---------|------|-------|
| `src/tools/*` | Utilities | Used by main services |
| `src/cloudflare/*` | Deployment | Public API |
| `conference_worker.ts` | Standalone | Independent |
| `kilocode/` | Utilities | Core code |

---

*Architecture decisions are documented in ADRs (Architecture Decision Records) in `/docs/adr/`.*
