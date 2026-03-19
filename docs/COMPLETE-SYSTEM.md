# BAC Knowledge System - Complete Documentation

AI-native knowledge management system using aichat, Gemini, and Rust tool crates for semantic search, content generation, and knowledge graph visualization.

## Executive Summary

The BAC Knowledge System transforms how students prepare for exams by:

- **Processing** study materials automatically (OCR → AI correction → structured notes)
- **Searching** across all content semantically (find concepts by meaning, not keywords)
- **Connecting** related ideas through automatic wikilinks and knowledge graphs
- **Generating** quality study notes with LaTeX math support

**Target Users**: Students preparing for BAC (Baccalauréat) exams in Physics, Chemistry, Mathematics, and Biology.

**Tech Stack**: aichat CLI + Rust microservices + PostgreSQL/pgvector + Obsidian vault

---

## System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                          USER INTERFACE                              │
│            (Phone: Termux + aichat) or (Desktop)                     │
└─────────────────────────────────────┬───────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      aichat Agent Layer                              │
│  ┌──────────────────┐  ┌────────────────┐  ┌──────────────────────┐  │
│  │  bac-knowledge   │  │  RAG Engine    │  │   Tool Router       │  │
│  │  (Orchestrator)  │  │  (semantic)    │  │   (5 microservices) │  │
│  └──────────────────┘  └────────────────┘  └──────────────────────┘  │
└─────────────────────────────────────┬───────────────────────────────┘
                                      │
          ┌───────────────────────────┼───────────────────────────┐
          │                           │                           │
          ▼                           ▼                           ▼
┌──────────────────┐       ┌──────────────────┐       ┌──────────────────┐
│   Gemini API     │       │   pgvector DB    │       │  Obsidian Vault  │
│   (AI Model)     │       │   (Neon Cloud)   │       │  (File Storage)  │
└──────────────────┘       └──────────────────┘       └──────────────────┘
```

### Data Flow

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│  Input  │───▶│   OCR   │───▶│ Correct │───▶│ Analyze │───▶│ Generate│
│ (PDF)   │    │         │    │  (AI)   │    │         │    │  (Note) │
└─────────┘    └─────────┘    └─────────┘    └─────────┘    └────┬────┘
                                                                  │
                    ┌────────────────────────────────────────────┘
                    │
                    ▼
        ┌───────────────────┐       ┌───────────────────┐
        │  Obsidian Vault   │       │   pgvector DB     │
        │  (file storage)   │       │  (embeddings)     │
        └───────────────────┘       └───────────────────┘
```

---

## Architecture Diagram

```
Phone (Termux + aichat)
         │
         ▼
┌─────────────────────────────────────────────────────────────────┐
│                         aichat Agent Layer                       │
│  Agent: bac-knowledge  │  RAG: resources/notes/                 │
└─────────────────────────┬───────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┬─────────────────┐
        │                 │                 │                 │
        ▼                 ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│ gemini-tools   │ │ vector-tools   │ │ vault-tools   │ │ cloud-tools   │
│    :3001      │ │    :3002       │ │    :3003      │ │    :3004      │
│               │ │                │ │               │ │               │
│ • analyze     │ │ • search       │ │ • read        │ │ • SSH exec    │
│ • embed       │ │ • insert       │ │ • write       │ │ • upload      │
│ • generate    │ │ • batch        │ │ • link        │ │ • OCR         │
│ • correct     │ │ • delete       │ │ • moc         │ │ • GCS sync    │
│ • extract     │ │ • rebuild      │ │ • search      │ │               │
└───────────────┘ └───────────────┘ └───────────────┘ └───────────────┘
        │                 │                 │                 │
        └─────────────────┴─────────────────┴─────────────────┘
                              │
              ┌───────────────┴───────────────┐
              ▼                               ▼
┌─────────────────────────┐     ┌─────────────────────────┐
│      Neon pgvector      │     │      Obsidian Vault     │
│   (semantic database)   │     │     (file storage)      │
│                         │     │                         │
│ • HNSW index (m=16)     │     │ • Markdown files        │
│ • 768-dim embeddings    │     │ • Wikilinks             │
│ • Similarity search     │     │ • YAML frontmatter      │
└─────────────────────────┘     └─────────────────────────┘
```

---

## Service Descriptions

### gemini-tools (Port 3001)

AI content analysis and generation powered by Google's Gemini.

| Function | Purpose | Parameters |
|----------|---------|------------|
| `analyze` | Analyze content structure | `content`, `subject` |
| `extract` | Extract entities/facts | `content` |
| `generate` | Create study notes | `topic`, `subject`, `format` |
| `embed` | Create semantic embeddings | `text`, `task_type` |
| `correct` | Fix OCR text errors | `ocr_text`, `language` |

**Example Request:**
```bash
curl -X POST http://localhost:3001/generate \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Write about Newton's laws", "subject": "physics"}'
```

### vector-tools (Port 3002)

Semantic search operations using pgvector HNSW indexes.

| Function | Purpose | Parameters |
|----------|---------|------------|
| `search` | Semantic similarity search | `query`, `top_k`, `filters` |
| `insert` | Store embedded content | `text`, `embedding`, `metadata` |
| `batch` | Bulk insert records | `records[]` |
| `delete` | Remove content by ID | `id` |
| `upsert` | Update existing content | `id`, `text`, `embedding` |
| `rebuild` | Rebuild HNSW index | `m`, `ef_construction` |

**Example Request:**
```bash
curl -X POST http://localhost:3002/search \
  -H "Content-Type: application/json" \
  -d '{"query": [0.1, 0.2, ...], "top_k": 5}'
```

### vault-tools (Port 3003)

Obsidian vault operations for note management.

| Function | Purpose | Parameters |
|----------|---------|------------|
| `read` | Read vault note | `path` |
| `write` | Create/update note | `path`, `content`, `frontmatter` |
| `link` | Create wikilink | `source`, `target`, `alias` |
| `moc` | Update table of contents | `folder` |
| `search` | Full-text search vault | `query`, `path` |

**Example Request:**
```bash
curl -X POST http://localhost:3003/write \
  -H "Content-Type: application/json" \
  -d '{"path": "Physics/Thermodynamics.md", "content": "# Thermo..."}'
```

### cloud-tools (Port 3004)

Cloud Shell integration for heavy processing.

| Function | Purpose | Parameters |
|----------|---------|------------|
| `ssh` | Execute remote command | `host`, `command`, `timeout` |
| `upload` | Upload file to cloud | `source`, `destination` |
| `download` | Download file from cloud | `url`, `destination` |
| `ocr` | OCR via Cloud Shell | `file_path`, `language` |
| `gcs_sync` | Sync with Google Cloud Storage | `bucket`, `direction` |

**Example Request:**
```bash
curl -X POST http://localhost:3004/ssh \
  -H "Content-Type: application/json" \
  -d '{"command": "ls -la", "timeout": 30}'
```

### graph-tools (Port 3005)

Knowledge graph generation and visualization.

| Function | Purpose | Parameters |
|----------|---------|------------|
| `generate` | Generate graph from notes | `source_path`, `depth`, `format` |
| `extract` | Extract entities/relationships | `text`, `relationship_types` |
| `export` | Export to format | `format`, `output_path` |

**Example Request:**
```bash
curl -X POST http://localhost:3005/generate \
  -H "Content-Type: application/json" \
  -d '{"note_path": "Physics/Thermo.md", "depth": 2}'
```

---

## Configuration Guide

### Environment Variables

```bash
# Core Services
GATEWAY_HOST=127.0.0.1:8080
GEMINI_TOOLS_URL=http://localhost:3001
VECTOR_TOOLS_URL=http://localhost:3002
VAULT_TOOLS_URL=http://localhost:3003
CLOUD_TOOLS_URL=http://localhost:3004
GRAPH_TOOLS_URL=http://localhost:3005

# Database (Neon PostgreSQL with pgvector)
NEON_DB_URL=postgresql://user:pass@host/db?sslmode=require

# AI Provider
GEMINI_API_KEY=your_api_key_here

# Storage (Garage S3 or GCS)
GARAGE_ENDPOINT=http://localhost:3900
GCS_BUCKET=your-bucket-name

# Vault Path
VAULT_PATH=/path/to/resources/notes

# Cloud Shell
CLOUD_SHELL_HOST=your-instance
CLOUD_SHELL_KEY=/path/to/ssh-key
```

### aichat Configuration

Located in `config/aichat/`:

```
config/aichat/
├── config.yaml              # Main aichat configuration
├── agents/
│   └── bac-knowledge/
│       └── instructions.md  # Agent prompt (324 lines)
├── roles/
│   ├── physicist.md         # Physics expert persona
│   ├── chemist.md           # Chemistry expert persona
│   ├── biologist.md         # Biology expert persona
│   └── math_gen.md          # Math generation persona
├── macros/
│   ├── process.yaml         # Content processing workflow
│   └── sync-search.yaml     # Sync + search workflow
└── rag-template.txt         # RAG response template
```

### Agent Instructions

The `bac-knowledge` agent instructions define:
- **Identity**: Orchestrator for BAC exam prep
- **Capabilities**: Content processing, note generation, semantic search
- **Tool references**: All available tools with parameters
- **Formatting rules**: LaTeX syntax, wikilink conventions
- **Workflows**: Example processing pipelines

### RAG Template

Defines how retrieved context is formatted in responses:

```markdown
## Relevant Information

{{context}}

## Answer

{{answer}}
```

---

## Installation Steps

### 1. Install Dependencies

```bash
# Rust toolchain
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source ~/.cargo/env

# aichat CLI
curl -fsSL https://raw.githubusercontent.com/sigoden/aichat/main/install.sh | sh

# Node.js (for scripts)
npm install
```

### 2. Configure Environment

```bash
# Copy environment template
cp .env.example .env

# Edit with your credentials
nano .env
```

### 3. Set Up Database

```bash
# Create Neon database
# Go to https://neon.tech and create a project

# Run schema
psql $NEON_DB_URL -f sql/schema.sql

# Verify pgvector
psql $NEON_DB_URL -c "CREATE EXTENSION IF NOT EXISTS vector;"
psql $NEON_DB_URL -c "SELECT extversion FROM pg_extension WHERE extname = 'vector';"
```

### 4. Build Services

```bash
# Build all Rust services
cargo build --release

# Or build individually
cargo build -p gemini-tools --release
cargo build -p vector-tools --release
cargo build -p vault-tools --release
cargo build -p cloud-tools --release
cargo build -p graph-tools --release
```

### 5. Start Services

```bash
# Start all services via daemon script
./scripts/bac-agent-daemon.sh

# Or start individually
cargo run -p gemini-tools --release &
cargo run -p vector-tools --release &
cargo run -p vault-tools --release &
cargo run -p cloud-tools --release &
cargo run -p graph-tools --release &
```

### 6. Verify Services

```bash
# Check all services
curl http://localhost:3001/health  # gemini-tools
curl http://localhost:3002/health  # vector-tools
curl http://localhost:3003/health  # vault-tools
curl http://localhost:3004/health  # cloud-tools
curl http://localhost:3005/health  # graph-tools
```

### 7. Initialize Vault

```bash
# Create vault structure
mkdir -p resources/notes/{00-Meta,01-Concepts,02-Practice,03-Resources,04-Exams}

# Set permissions
chmod -R 755 resources/notes
```

### 8. Start Using

```bash
# Start aichat with BAC configuration
aichat --config config/aichat/config.yaml
```

---

## Usage Examples

### Search for Notes

```bash
# Natural language search
> search for electric field concepts
> find notes about thermodynamics
> what does my notes say about photosynthesis
```

### Generate Notes

```bash
# Create study note
> write a note about thermodynamics
> create a summary of photosynthesis
> generate practice problems for kinematics
```

### Process Files

```bash
# OCR and extract
> process new files in 03-Resources
> extract text from test.jpg
> analyze this PDF for key concepts
```

### Knowledge Graphs

```bash
# Generate graph
> show connections for quantum mechanics
> generate graph for neural networks
> visualize the chemistry MOC
```

### Auto-linking

```bash
# Link related notes
> link all thermodynamics notes together
> create bidirectional links for Newton's laws
```

---

## Security Considerations

### API Keys

| Key | Risk Level | Protection |
|-----|------------|------------|
| GEMINI_API_KEY | High | Never commit to git; use .env |
| NEON_DB_URL | High | Contains credentials; SSL required |
| Cloud SSH keys | High | 0600 permissions; gitignored |

### Service Authentication

**Current State**: No authentication on microservice endpoints.

**Recommendations**:
- Bind services to localhost only
- Use firewall rules for production
- Add API key authentication for cloud deployments

### Data Privacy

- All notes stored locally in Obsidian vault
- Embeddings stored in Neon (cloud database)
- Consider E2E encryption for sensitive materials

### Git Security

```bash
# Ensure sensitive files are ignored
# .gitignore should contain:
.env
*.key
*.pem
~/.aichat/
```

---

## Performance Notes

### pgvector HNSW Index

Default configuration for good quality/speed balance:

```sql
CREATE INDEX ON notes USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);
```

| Parameter | Default | Speed | Quality |
|-----------|---------|-------|---------|
| `m` | 16 | Faster = lower | Higher = better |
| `ef_construction` | 64 | Faster = lower | Higher = better |
| `ef_search` | 40 | Faster = lower | Higher = better |

### Embedding Dimensions

- **Gemini embeddings**: 768 dimensions
- **Chunk size**: 2048 tokens (configurable)
- **Batch size**: 100 records for bulk inserts

### Memory Usage

| Service | Typical RAM |
|---------|-------------|
| gemini-tools | ~100MB |
| vector-tools | ~200MB |
| vault-tools | ~50MB |
| cloud-tools | ~100MB |
| graph-tools | ~150MB |

### Latency Expectations

| Operation | Typical Latency |
|-----------|-----------------|
| Semantic search | 50-200ms |
| Note generation | 2-5s |
| OCR processing | 5-30s |
| Graph generation | 1-10s |

---

## Monitoring

### Health Checks

```bash
# All services
curl http://localhost:8080/health

# Individual
curl http://localhost:3001/health
curl http://localhost:3002/health
```

### Debug Logging

```bash
# Enable debug output
RUST_LOG=debug cargo run -p gemini-tools

# Full trace
RUST_LOG=trace cargo run -p vector-tools
```

### Metrics

Services expose Prometheus metrics at `/metrics`:
- Request counts
- Latency histograms
- Error rates

---

## Troubleshooting

See [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) for common issues and solutions.

---

## License

MIT OR Apache-2.0
