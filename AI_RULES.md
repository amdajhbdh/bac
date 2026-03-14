# AI Tools Rules - BAC Unified

## Mandatory Rules for All AI Agents

This file contains **RESTRICTIVE RULES** that ALL AI coding agents MUST follow when working in this repository. These rules are **ENFORCED** - do not deviate from them under any circumstances.

---

## 0. VERSION CONTROL - MANDATORY

### REQUIRED TOOL: Jujutsu (jj)
- Use **jj** for ALL version control operations
- NEVER use git directly for commits
- jj provides better history rewrite and working copy management

### ALLOWED COMMANDS:
```bash
jj new                  # Create new change
jj describe             # Edit description
jj commit               # Commit current change
jj log                  # View history
jj diff                 # Show changes
jj push                 # Push to remote
jj pull                 # Pull from remote
jj undo                 # Undo last operation
jj squash               # Squash commits
jj rebase               # Rebase changes
```

### NEVER USE:
- ❌ git commit
- ❌ git push
- ❌ git rebase -i
- ❌ git reset --hard

---

## 0.1. PROJECT PROGRESS - ALWAYS REFER

### MANDATORY: Check Tasks Before Working
- ALWAYS check `openspec/changes/full-bac-unified/tasks.md` before implementing
- Update task status: `- [ ]` → `- [x]` when complete
- NEVER work on already completed tasks
- Report progress: "Working on X/Y tasks complete"

### Progress Check Commands:
```bash
# View current progress
grep -c "\- \[x\]" openspec/changes/full-bac-unified/tasks.md

# List remaining tasks in a phase
grep "### 3.1" openspec/changes/full-bac-unified/tasks.md -A 20
```

---

## 1. DATABASE - STRICT

### ONLY ALLOWED:
- **PostgreSQL** with **pgvector** extension for vector storage
- **Redis** for caching and sessions
- **Turso** (SQLite) for production cache ONLY
- Local SQLite allowed for development

### NEVER USE:
- ❌ MongoDB
- ❌ MySQL without explicit approval
- ❌ SQLite for anything except Turso (prod cache) or local dev
- ❌ Other databases without explicit approval

### REQUIRED TOOLS:
- **sqlc** for type-safe SQL queries (required for all DB operations)
- Generate Go code with `sqlc generate`

### REFERENCE:
- `doc/OPEN_SOURCE_PROJECTS.md` - Database section
- `sql/schema.sql` - Current schema
- Use `TURSO_DB_URL` for production cache (Turso is distributed SQLite)

---

## 2. STORAGE - STRICT

### ONLY ALLOWED:
- **Garage S3** (local development)
- **MinIO** (production alternative)
- Use S3-compatible APIs only

### NEVER USE:
- ❌ AWS S3 directly (unless required)
- ❌ Google Cloud Storage
- ❌ Azure Blob Storage
- ❌ Local filesystem for production blobs

### REFERENCE:
- `config/storage.yaml` - Storage configuration

---

## 3. SEARCH INFRASTRUCTURE - STRICT

### MANDATORY ARCHITECTURE:
```
Hybrid Search = Elasticsearch (full-text) + pgvector (vectors)
```

### ONLY ALLOWED:
- **Elasticsearch** for full-text search
- **pgvector** for vector similarity search
- **PostgreSQL** for metadata

### NEVER USE:
- ❌ Algolia (paid, not self-hosted)
- ❌ MeiliSearch without Elasticsearch fallback
- ❌ Typesense without Elasticsearch
- ❌ Solo pgvector without Elasticsearch for full-text
- ❌ Any other search engine without RFC approval

### REFERENCE:
- `doc/rfc/SEARCH_INFRASTRUCTURE.md`
- `doc/rfc/RESOURCE_SEARCH_SYSTEM.md`

---

## 4. CONTENT PROCESSING - STRICT

### REQUIRED TOOLS (in priority order):

| Priority | Tool | Purpose | URL |
|----------|------|---------|-----|
| 1 | Tesseract | OCR | github.com/tesseract-ocr |
| 2 | Surya | ML OCR | github.com/VikParuchuri/surya |
| 3 | pdfcpu | PDF processing | pdfcpu.io |
| 4 | Poppler | PDF utilities | poppler.freedesktop.org |
| 5 | Apache Tika | Content extraction | tika.apache.org |
| 6 | yt-dlp | YouTube extraction | github.com/yt-dlp |
| 7 | Playwright | Web scraping | playwright.dev |

### NEVER USE:
- ❌ Google Vision API (requires billing)
- ❌ AWS Textract (requires AWS account)
- ❌ Azure Computer Vision (requires Azure)
- ❌ Paid OCR services
- ❌ Closed-source tools

### REFERENCE:
- `doc/rfc/CONTENT_PROCESSING_PIPELINE.md`

---

## 5. FILE BUNDLING - STRICT

### MANDATORY FORMATS:

| Size | Format | Purpose |
|------|--------|---------|
| < 100MB | **WRAPX** | Metadata-rich bundles |
| >= 100MB | **ARX** | Mountable archives |

### REQUIRED STRUCTURE:
```
WRAPX Bundle (.bacx):
├── /payload/           # Original files
├── /meta/
│   ├── manifest.json   # Full metadata
│   ├── classification.json
│   └── tags.json
├── /embeddings/       # pgvector format
└── /audit/           # Chain of custody
```

### NEVER USE:
- ❌ tar/zip without metadata wrapper
- ❌ Proprietary formats
- ❌ Non-standard containers

### REFERENCE:
- `doc/rfc/FILE_BUNDLING_SYSTEM.md`

---

## 6. AUTHENTICATION - STRICT

### REQUIRED:
- **JWT** for token authentication
- **bcrypt** for password hashing
- **TOTP** (Google Authenticator) for MFA

### NEVER USE:
- ❌ Basic auth (except health endpoints)
- ❌ API keys in URLs
- ❌ Hardcoded credentials
- ❌ Session tokens in URL parameters

### ENVIRONMENT VARIABLES:
All secrets MUST use:
- `NEON_DB_URL`
- `GARAGE_ENDPOINT`
- `GARAGE_ACCESS_KEY`
- `GARAGE_SECRET_KEY`
- `REDIS_URL`
- `JWT_SECRET`
- `OLLAMA_HOST`

### REFERENCE:
- `AGENTS.md` - Environment Variables section

---

## 7. AI/LLM - STRICT

### PRIORITY ORDER (Offline First):

1. **Ollama** (local) - ALWAYS try first
2. **NotebookLM** (nlm CLI) - Second choice
3. Cloud APIs - LAST RESORT only

### NEVER USE:
- ❌ Cloud APIs as primary (only fallback)
- ❌ Paid AI APIs without approval
- ❌ API keys in code (use env vars)

### REFERENCE:
- `docs/AGENTS.md` - NLM Maximization rules

---

## 8. LOGGING - STRICT

### MANDATORY:
- Use **`log/slog`** (Go built-in) for ALL logging
- Structured JSON logging for production

### NEVER USE:
- ❌ `fmt.Print*` for logging
- ❌ `print()` statements
- ❌ Any other logging library unless required

### REFERENCE:
- `.cursorrules` - "never write useless boring print statements"

---

## 9. CODING STANDARDS

### MANDATORY:

```go
// 1. Error handling - ALWAYS wrap with context
if err != nil {
    return nil, fmt.Errorf("doing something: %w", err)
}

// 2. Imports - Group stdlib, external, internal
import (
    // Standard library
    "context"
    "log/slog"
    
    // External packages
    "github.com/gin-gonic/gin"
    
    // Internal packages
    "github.com/bac-unified/agent/internal/solver"
)

// 3. Configuration - NEVER hardcode
dbURL := os.Getenv("NEON_DB_URL")
if dbURL == "" {
    dbURL = "postgres://localhost/bac"
}
```

### REFERENCE:
- `AGENTS.md` - Code Style Guidelines section

---

## 10. TESTING - STRICT

### MANDATORY:
- Tests in same package with `_test.go` suffix
- Table-driven test pattern
- Single test: `go test -v -run TestName ./...`
- Run with race detector: `go test -race ./...`

### NEVER USE:
- ❌ Skip tests to "save time"
- ❌ Test files in separate packages without reason

### REFERENCE:
- `AGENTS.md` - Testing section

---

## 11. API DESIGN - STRICT

### REQUIRED PATTERNS:

```go
// Response wrapper
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// Error handling
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

### NEVER USE:
- ❌ Different response formats per endpoint
- ❌ Error messages in multiple languages
- ❌ Expose internal errors to clients

---

## 12. DEVOPS - STRICT

### ONLY ALLOWED:

| Category | Tool | Purpose |
|----------|------|---------|
| Container | Docker | Containers |
| Orchestration | Kubernetes/K3s | Scale |
| CI/CD | GitHub Actions | Automation |
| Monitoring | Prometheus + Grafana | Observability |
| Secrets | HashiCorp Vault | Secrets (or env vars) |

### NEVER USE:
- ❌ Heroku (paid)
- ❌ Vercel (not self-hostable)
- ❌ Proprietary CI/CD
- ❌ Hardcoded secrets

---

## 13. MONITORING - REQUIRED

### MANDATORY METRICS:
- Request latency (p50, p95, p99)
- Error rate
- Cache hit rate
- Queue depth

### REQUIRED TOOLS:
- **Prometheus** for metrics
- **Grafana** for dashboards
- **slog** for structured logging

### REFERENCE:
- `doc/rfc/SEARCH_INFRASTRUCTURE.md` - Monitoring section

---

## 14. CACHE STRATEGY - MANDATORY

### REQUIRED CACHE LAYERS:

```
Request → Redis (hot cache)
            ↓ (miss)
         PostgreSQL (NLM cache)
            ↓ (miss)
         NLM/AI Service
            ↓
         Store result in cache
```

### NEVER USE:
- ❌ No caching for expensive operations
- ❌ Cache without TTL
- ❌ Cache without invalidation strategy

---

## 15. INTEGRATION DECISIONS (LOCKED)

These decisions are **FINAL** - do not revisit:

| Category | Decision | Rationale |
|----------|----------|-----------|
| Database | PostgreSQL + pgvector | Already implemented |
| Search | Elasticsearch + pgvector | RFC approved |
| Storage | Garage S3 | Self-hosted |
| Cache | PostgreSQL (not SQLite) | Performance |
| Bundling | WRAPX + ARX | RFC approved |
| OCR | Tesseract + Surya | Free, Arabic support |
| Auth | JWT + bcrypt | Industry standard |

---

## 16. REFERENCE CHECKLIST

Before writing any code, verify:

- [ ] Using PostgreSQL (not SQLite) for data?
- [ ] Using Garage S3 (not direct cloud)?
- [ ] Using Elasticsearch + pgvector for search?
- [ ] Using Tesseract/Surya for OCR?
- [ ] Using WRAPX/ARX for bundles?
- [ ] Using `log/slog` for logging?
- [ ] Using environment variables for config?
- [ ] Wrapping errors with context?
- [ ] Grouping imports correctly?
- [ ] Following table-driven test pattern?
- [ ] Using JWT + bcrypt for auth?

---

## VIOLATIONS

Any AI agent found violating these rules will be **CORRECTED** without warning. Repeated violations may result in repository access restrictions.

---

## EMERGENCY EXCEPTIONS

If you believe a rule must be broken, you MUST:
1. Document the exception in a comment
2. Get explicit approval from repository owner
3. Create an RFC explaining why

---

*Last Updated: 2026-03-01*
*Enforcement: IMMEDIATE*
*Contact: Repository Maintainers*
