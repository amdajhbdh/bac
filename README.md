# BAC Knowledge System

AI-native knowledge management system using aichat, Gemini, and Rust tool crates for semantic search, content generation, and knowledge graph visualization.

## Features

- 📚 **RAG-powered Search** — Semantic similarity search across all notes using pgvector/HNSW
- 🤖 **AI Content Generation** — Generate notes with LaTeX support via Gemini
- 🔗 **Auto-linking** — Connect related concepts using wikilinks and MOCs
- 📊 **Knowledge Graphs** — Visualize note relationships as interactive graphs
- ☁️ **Cloud Integration** — Process via Cloud Shell (Cloudflare/Garage)
- 🔍 **OCR Pipeline** — Extract text from images and PDFs

## Quick Start

### 1. Install Dependencies

```bash
# Rust toolchain
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# aichat CLI
curl -fsSL https://raw.githubusercontent.com/sigoden/aichat/main/install.sh | sh

# Node.js (for scripts)
npm install
```

### 2. Configure

```bash
# Copy environment template
cp .env.example .env

# Edit with your API keys and connections
# Required: GEMINI_API_KEY, database URLs
```

### 3. Build Tools

```bash
# Build all Rust services
cargo build --release

# Or build individually
cargo build -p gemini-tools --release
cargo build -p vector-tools --release
```

### 4. Start Services

```bash
# Start all services
./scripts/bac-agent-daemon.sh

# Or start individually (see Services section)
cargo run -p gemini-tools --release
cargo run -p vector-tools --release
```

### 5. Use aichat

```bash
# Start with BAC agent configuration
aichat --config config/aichat/config.yaml
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    aichat (CLI)                          │
│  Agent: bac-knowledge  │  RAG: resources/notes/          │
└─────────────────────────┬───────────────────────────────┘
                          │
          ┌───────────────┼───────────────┬──────────────┐
          │               │               │              │
          ▼               ▼               ▼              ▼
    ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐
    │  gemini   │  │  vector   │  │   vault   │  │   cloud   │
    │  -tools   │  │  -tools   │  │  -tools   │  │  -tools   │
    │   :3001   │  │   :3002   │  │   :3003   │  │   :3004   │
    └───────────┘  └───────────┘  └───────────┘  └───────────┘
          │               │               │              │
          └───────────────┴───────────────┼──────────────┘
                                          │
                            ┌─────────────┴─────────────┐
                            │       api-gateway (:8080)   │
                            └─────────────────────────────┘
```

## Services

| Service | Port | Purpose |
|---------|------|---------|
| api-gateway | 8080 | Main API gateway |
| ocr | 3000 | OCR processing (Tesseract) |
| gemini-tools | 3001 | Gemini API wrapper |
| vector-tools | 3002 | pgvector semantic search |
| vault-tools | 3003 | Obsidian vault operations |
| cloud-tools | 3004 | Cloud Shell integration |
| graph-tools | 3005 | Knowledge graph generation |

## Usage Examples

### Search notes
```
> search for electric field concepts
> find notes about thermodynamics
```

### Generate note
```
> write a note about thermodynamics
> create a summary of photosynthesis
```

### Process files
```
> process new files in 03-Resources
> extract text from test.jpg
```

### Knowledge graph
```
> show connections for quantum mechanics
> generate graph for neural networks
```

## Project Structure

```
bac/
├── src/
│   ├── api/              # Main API gateway
│   ├── gemini-tools/      # Gemini API wrapper
│   ├── vector-tools/     # pgvector operations
│   ├── vault-tools/      # Obsidian vault ops
│   ├── cloud-tools/      # Cloud Shell integration
│   ├── graph-tools/      # Knowledge graph
│   └── ocr/              # OCR processing
├── config/
│   └── aichat/           # aichat configuration
│       ├── config.yaml    # Main config
│       ├── agents/       # Agent definitions
│       └── rag-template.txt
├── resources/
│   └── notes/            # Obsidian vault
├── scripts/              # Automation scripts
├── sql/                  # Database schemas
└── docs/                 # Documentation
    ├── API.md
    └── TROUBLESHOOTING.md
```

## Configuration

### Environment Variables

```bash
# Core
GATEWAY_HOST=127.0.0.1:8080

# Database (Neon PostgreSQL)
NEON_DB_URL=postgresql://user:pass@host/db?sslmode=require

# AI (Gemini)
GEMINI_API_KEY=your_api_key

# Storage (Garage S3)
GARAGE_ENDPOINT=http://localhost:3900
```

### aichat Tools

| Tool | Service | Description |
|------|---------|-------------|
| `gemini_*` | gemini-tools | Generate, extract, embed content |
| `vector_*` | vector-tools | Semantic search operations |
| `vault_*` | vault-tools | Read/write vault notes |
| `cloud_*` | cloud-tools | SSH, upload, download |
| `graph_*` | graph-tools | Generate/extract graphs |

## Development

```bash
# Run all tests
cargo test

# Build release
cargo build --release

# Run specific service
cargo run -p gemini-tools --release

# Watch mode (dev)
cargo watch -x build
```

## License

MIT OR Apache-2.0
