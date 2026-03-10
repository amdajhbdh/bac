# BAC Unified - Agent Instructions

## Project Overview

| Component | Language | Path |
|-----------|----------|------|
| Agent CLI | **Go** | `src/agent/` |
| REST API | Go (Gin) | `src/api/` |
| Gateway | Rust (Axum) | `src/gateway/` |
| OCR Service | Rust (Tonic) | `src/ocr-service/` |
| Animations | Rust (Noon) | `src/noon/` |
| Web Frontend | React (Bun) | `src/web/` |
| Daemons | Nushell | `daemons/` |
| CLI Orchestration | Nushell | `lib/bac.nu` |
| Cloudflare Worker | TypeScript | `src/cloudflare-worker/` |
| Cloudflare Pages | HTML/JS | `src/cloudflare-pages/` |
| Study Vault | Obsidian | `resources/vault/` |

The `src/api` module imports `src/agent` via a `replace` directive in `go.mod`. Both modules use Go 1.25.

## Architecture

```
Nushell CLI (lib/bac.nu)
    └── Go Agent CLI (src/agent/cmd/main.go)
            ├── solver/       Ollama (local) → Cloud APIs → NLM → Online research
            ├── memory/       pgvector similarity search (RAG)
            ├── nlm/          NotebookLM integration + Turso (SQLite) cache
            ├── ocr/          Tesseract + Surya
            ├── search/       Hybrid: Elasticsearch (full-text) + pgvector (vector)
            ├── animation/    Noon DSL code generation
            ├── bundle/       WRAPX (<100MB) + ARX (≥100MB)
            ├── aggregator/   Khan Academy, OpenStax, MIT OCW, OER
            ├── db/           sqlc-generated PostgreSQL queries
            └── cache/        Redis
```

**Database Layer:** PostgreSQL + pgvector (main), Turso (SQLite, NLM cache), Redis (hot cache)

## AI Inference Priority

| Priority | Model | Use Case |
|----------|-------|----------|
| 1 | Ollama (llama3.2:3b) | Local, fast inference |
| 2 | Cloud APIs | kimi, deepseek, minimax |
| 3 | NLM CLI/MCP | NotebookLM with caching |

## Hybrid Search (RRF)

1. Parse query, extract filters (subject, chapter, difficulty, language, source)
2. Run: Elasticsearch full-text + pgvector cosine similarity (1536-dim)
3. Fuse with RRF normalization
4. Return with facets + suggestions

## Agent Ecosystem

The system uses a multi-agent architecture for vault enhancement:

- **Coordinator Agent**: Main orchestrator on Cloudflare Workers (Durable Objects)
- **Content Scanner Agent**: PDF/image OCR using Workers AI vision models
- **Resource Finder Agent**: Web search for educational content
- **Study Coach Agent**: Progress tracking and study recommendations
- **Link Builder Agent**: Knowledge graph and cross-reference generation

### MCP Servers
- **Vault MCP**: Local file operations, search, link management
- **Image Processor**: OCR, diagram recognition, equation extraction
- **Web Search**: YouTube, educational sites

### Integrations
- **Telegram**: Ironclaw/OpenClaw gateway for commands
- **Cloudflare Workers**: AI processing, stateful agents
- **Kiro CLI**: Orchestration, hooks, agent management

---

## Testing

- Go: `go test -v -run TestName ./...`
- Rust: `cargo test test_name`
- React: `npm run test`

---

## Code Style Guidelines

### Imports (Go)
Group in order: standard library → external packages → internal packages

### Naming Conventions
- Types/Interfaces: PascalCase
- Variables/Functions: camelCase
- Constants: PascalCase or SCREAMING_SNAKE_CASE
- Packages: snake_case
- Files: snake_case.go

### Error Handling
Wrap errors with context using `fmt.Errorf` with `%w`

### Logging
Use `log/slog` exclusively. Never `fmt.Print*` or `print()`

### Configuration
All config via environment variables. Never hardcode

### Type Safety
- Avoid `interface{}` / `any`
- Use `context.Context` as first parameter
- Return concrete types where possible

### Testing
Use table-driven tests with subtests

---

## Mandatory Rules

| Rule | Requirement |
|------|-------------|
| **VCS** | Use **jj** (Jujutsu). Never git directly |
| **Storage** | Garage S3 or MinIO only (S3-compatible) |
| **Bundling** | WRAPX (<100MB), ARX (≥100MB) with metadata |
| **OCR** | Tesseract + Surya only. Never paid APIs |
| **Auth** | JWT + bcrypt. All secrets via env vars |
| **Search** | Hybrid = Elasticsearch + pgvector |
| **Monitoring** | Prometheus + Grafana |
| **AI Priority** | Ollama → NotebookLM → Cloud fallback |

---

## Database

- **Main DB**: PostgreSQL + pgvector (`sql/schema.sql`)
- **NLM cache**: Turso (SQLite) at `src/agent/internal/nlm/cache/`
- **Hot cache**: Redis
- **Rule**: Use **sqlc** for all queries. Never raw SQL or ORM

---

## Study Vault

Located at `resources/vault/`

### Math
- Complex Numbers, Integrals, Number Theory, Matrices

### Physics
- Dynamics, Electromagnetism, Lorentz Force, Magnetic Induction, Satellites, Energy, AC/DC Circuits

### Chemistry
- Solutions, Organic Chemistry, Kinetics

### Natural Sciences
- Genetics, Human Reproduction, Nervous System, Geology (Mauritania)

### Features
- 1200+ practice problems (Easy, Medium, Hard)
- 100+ video links (YouTube playlists)
- 50+ interactive tools (Desmos, PhET, Molview, BioDigital)
- Excalidraw diagrams, Dataview progress tracking

### Key Files
- `000-INDEX.md` - Main navigation
- `QUICK-START.md` - Study schedule
- `Resource Links.md` - External resources

---

## Before Working

1. Check `openspec/changes/full-bac-unified/tasks.md` — update `- [ ]` → `- [x]` on completion
2. Read `doc/rfc/*.md` for architecture decisions
3. Read `doc/OPEN_SOURCE_PROJECTS.md` for allowed tool list

---

*Last Updated: 2026-03-08*
