# BAC Unified - Ecosystem-Leveraged Implementation Plan

## Overview

This plan leverages existing skills, tools, and ecosystem plugins to minimize redundant implementation work.

## Existing Resources

### Skills (`.agents/skills/`)

| Skill | Purpose | Used In |
|-------|---------|---------|
| `postgres-pgvector/` | PostgreSQL + vector similarity | Phases 1, 5 |
| `turso-sqlite/` | Turso SQLite cache | Phase 1 |
| `garage-s3-storage/` | S3-compatible storage | Phases 1, 4 |
| `elasticsearch-hybrid/` | Full-text + vector search | Phases 3, 5 |
| `ollama-local-ai/` | Local LLM inference | Phases 1, 6 |
| `ocr-pipeline/` | OCR with fallbacks | Phase 3 ✅ |
| `rust-noon/` | Animation DSL | Phase 6 |
| `sqlc-generation/` | Type-safe SQL | Phases 1, 2 |
| `jj-workflow/` | Version control | All |
| `go-development/` | Go best practices | Phase 2 |
| `technical-integrations/` | Vendor integration patterns | All |

### Tools

| Tool | Purpose | Used In |
|------|---------|---------|
| **kiro-cli** | Agent CLI orchestration | All |
| **zeroclaw** | Zero-shot automation | Phases 3, 7 |
| **playwright-cli** | Browser automation | Phases 3, 7 |
| **noon** | Animation engine | Phase 6 |
| **ollama** | Local AI | Phases 1, 6 |

### Completed Implementations

| Component | Location | Tests | Status |
|-----------|----------|-------|--------|
| OCR Service | `src/ocr-service/` | 20 | ✅ Done |
| Gateway | `src/gateway/` | 6 | ✅ Done |

---

## Phase Implementation

### Phase 1: Infrastructure (Integration)

**Goal:** Wire existing components using skills

| Task | Using Skill/Tool | New Work |
|------|------------------|----------|
| PostgreSQL + pgvector | `postgres-pgvector` skill | Connect env vars |
| Turso Cache | `turso-sqlite` skill | Connect env vars |
| Garage S3 | `garage-s3-storage` skill | Connect env vars |
| Ollama Integration | `ollama-local-ai` skill | Connect env vars |
| sqlc Generation | `sqlc-generation` skill | Run `sqlc generate` |

**Tests:**
- `sqlc generate` passes
- DB connection test
- S3 upload/download test
- Ollama ping test

---

### Phase 2: API Server (Integration)

**Goal:** Connect Go API to existing infrastructure

| Task | Using Skill/Tool | New Work |
|------|------------------|----------|
| API Endpoints | `go-development` skill | Add routes |
| Auth (JWT) | Existing `services/auth.go` | Wire to DB |
| DB Services | `sqlc-generation` skill | Add queries |
| Rate Limiting | Redis (from Phase 1) | Connect |

**Tests:**
- Endpoint integration tests
- Auth flow tests
- Rate limit tests

---

### Phase 3: Content Aggregation (Partial New)

**Goal:** Complete OCR pipeline, add web scraping

| Task | Using Skill/Tool | New Work |
|------|------------------|----------|
| OCR Pipeline | `ocr-pipeline` skill | ✅ Done |
| Elasticsearch | `elasticsearch-hybrid` skill | Setup index |
| Web Scraping | `playwright-cli` | Setup scrapers |
| YouTube | `zeroclaw` | Create daemon |

**Tests:**
- OCR accuracy tests
- ES search tests
- Scraping validation tests

---

### Phase 4: File Bundling (New)

**Goal:** WRAPX/ARX implementation

| Task | New Work |
|------|----------|
| WRAPX Spec | Create spec doc |
| ARX Integration | Add library |
| Metadata System | Dublin Core + LOM |
| S3 Upload | Connect Garage |

**Tests:**
- Bundle creation tests
- Metadata extraction tests

---

### Phase 5: Question Bank (Integration)

**Goal:** Connect RAG to PostgreSQL

| Task | Using Skill/Tool | New Work |
|------|------------------|----------|
| Vector Search | `postgres-pgvector` skill | Wire to RAG |
| Hybrid Search | `elasticsearch-hybrid` skill | Connect ES |
| Question Import | Existing scripts | Wire to DB |

**Tests:**
- Similarity search tests
- Hybrid query tests

---

### Phase 6: AI Solver + Animation (Integration)

**Goal:** Connect solver and animation

| Task | Using Skill/Tool | New Work |
|------|------------------|----------|
| Ollama Solver | `ollama-local-ai` skill | Connect to Agent |
| Animation | `rust-noon` skill + Manim | Wire to Gateway |
| Multi-Provider | Existing fallback logic | Wire endpoints |

**Tests:**
- Solver integration tests
- Animation generation tests

---

### Phase 7: Prediction Engine (Partial New)

**Goal:** ML pipeline for predictions

| Task | Using Skill/Tool | New Work |
|------|------------------|----------|
| Data Collection | `zeroclaw` | Research patterns |
| ML Training | `ollama-local-ai` | Fine-tune model |
| Prediction API | Existing service | Wire endpoints |

**Tests:**
- Prediction accuracy tests

---

### Phase 8: Gamification (Integration)

**Goal:** Points, badges, leaderboards

| Task | Using Skill/Tool | New Work |
|------|------------------|----------|
| Points System | Existing schema | Add API |
| Achievements | Existing schema | Add API |
| Leaderboards | Existing schema | Add API |

**Tests:**
- Points calculation tests
- Badge award tests

---

### Phase 9: Frontend (New)

**Goal:** Web + PWA

| Task | New Work |
|------|----------|
| React Dashboard | Build UI |
| PWA Features | Add service worker |
| CLI Enhancement | Improve UX |

**Tests:**
- UI component tests
- E2E tests

---

### Phase 10: Deployment (Integration)

**Goal:** Kubernetes + monitoring

| Task | Using Skill/Tool | New Work |
|------|------------------|----------|
| K8s Manifests | Existing k8s/ | Integrate |
| Monitoring | Existing Grafana | Wire metrics |
| CI/CD | Existing configs | Connect |

**Tests:**
- Deployment tests
- Health check tests

---

## Test-Driven Development Workflow

```bash
# 1. Write test first
cd src/gateway && cargo test --test integration

# 2. Run to see failure
cargo test integration

# 3. Implement to pass test
# ... implementation ...

# 4. Refactor
cargo fmt && cargo clippy
```

---

## Key Integrations to Build

### 1. OCR → Gateway
- Gateway calls OCR service HTTP API
- Test: `test_gateway_ocr_request`

### 2. Gateway → Agent
- Agent CLI calls Gateway for OCR/animation
- Test: `test_agent_gateway_call`

### 3. RAG → PostgreSQL
- RAGEngine queries pgvector
- Test: `test_rag_pgvector_search`

### 4. Solver → Animation
- Solver triggers animation generation
- Test: `test_solver_animation_trigger`

---

## Success Criteria

- [ ] All existing skills used where applicable
- [ ] Tests written before implementation
- [ ] Integration tests pass
- [ ] No redundant implementations
- [ ] kiro-cli orchestrates workflows
- [ ] zeroclaw handles automation

---

*Last Updated: 2026-03-06*
*Plan Mode: SPEC + TDD*
