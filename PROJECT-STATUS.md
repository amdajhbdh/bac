# BAC Unified - Status Report

**Last Updated:** 2026-02-26 18:00

## 1. Project Overview

**BAC Unified** is a national exam preparation system for Mauritania's BAC C students. It provides AI-powered exercise solving with step-by-step solutions, question prediction based on historical patterns, and animation visualizations using the Noon (Rust-based) engine. Multi-modal input support includes images, PDFs, URLs, and handwritten notes.

Target users: BAC C students in Mauritania preparing for national exams

## 2. Tech Stack

- **Primary Language:** Go (agent, API)
- **Animation Engine:** Rust (Noon)
- **Legacy CLI:** Nushell (bac.nu)
- **Database:** PostgreSQL + pgvector (Neon)
- **Storage:** MinIO (local S3)
- **AI:** Ollama (local), Cloud APIs (NotebookLM, DeepSeek, Grok, Claude, ChatGPT)
- **Automation:** Playwright

## 3. Current Overall Status

- **Progress:** ~80% complete
- **Phase:** Alpha - Most core functionality working
- **Main modules that are live:**
  - Go Agent CLI with solver, memory, OCR, analyzer modules
  - Go REST API with handlers
  - Noon animation engine (compiled binary)
  - Database schema defined

## 4. What Was Completed

- Go Agent CLI (`src/agent/`) with modular structure
- REST API server (`src/api/`) with handlers for auth, solver, submission, prediction
- Noon animation engine (`src/noon/`) - Rust-based, compiled to release binary
- Database schema (`sql/schema.sql`) - Tables for questions, subjects, chapters, users, predictions
- Nushell CLI (`bac.nu`) - Legacy CLI with command structure
- Multi-service analyzer - Parallel analysis via Ollama, NLM, online services
- OCR pipeline - Image and PDF text extraction
- Memory module - Vector similarity search with pgvector
- Playwright browser automation primitives - Type, click, wait, snapshot functions
- Service-specific selectors - DeepSeek, Grok, Claude, ChatGPT configurations

## 5. Placeholders & Stubs (Incomplete Parts)

### Playwright Automation (Wired, Testing Pending)
- `src/agent/internal/analyzer/analyzer.go:251-319` - Updated to use `online.AutoSolve()` for ChatGPT, Claude, Grok, DeepSeek
- `openspec/changes/complete-playwright-automation/tasks.md:35-37` - Tasks 5.1-5.3 pending (manual testing required)
- `src/agent/cmd/main.go` - Now wires `online.AutoSolve()` when confidence < 0.7 or `--online` flag set

### Online Module
- `src/agent/internal/online/playwright.go` - Browser primitives fully implemented
- `src/agent/internal/online/online.go` - Auto-solve logic implemented and wired in main.go and analyzer.go

### Memory Module (Fixed)
- `src/agent/internal/memory/memory.go:140-200` - Now uses hash-based embedding fallback instead of zeros

### Nushell Daemons (Implemented)
- `daemons/youtube-daemon.nu` - Working: downloads YouTube audio with yt-dlp
- `daemons/webcrawler-daemon.nu` - Working: crawls sites and extracts PDF links
- `daemons/nlm-research-daemon.nu` - Working: starts NLM research (needs notebook ID)
- `daemons/audio-daemon.nu` - Working: checks NLM and lists notebooks

### Database (Connected)
- Schema already applied to Neon (tables: questions, subjects, users, predictions, etc.)
- `src/api/internal/services/db.go` - New DB connection module
- API now connects to Neon via `NEON_DB_URL` env var

### Animation Module (Improved)
- `src/agent/internal/animation/animation.go` - Now includes:
  - Smart noon path discovery
  - `CompileAndExport()` function for one-step generation
  - Fallback on compilation failure
  - Video export via ffmpeg (if available)
  - Proper Cargo.toml generation with serde

### API Services (Wired)
- `src/api/main.go` - All routes now public (no auth required)
- DB connection added to solver, submission, prediction services
- Server can start with `go run src/api/main.go`

## 6. What Will Be Implemented Next

### Next Milestone (immediate next 1-2 weeks)

- [ ] Complete Playwright automation testing (Tasks 5.1-5.3)
  - Test automation with each AI service (DeepSeek, Grok, Claude, ChatGPT)
  - Verify session persistence works correctly
  - Test error handling (rate limits, timeouts)

- [x] Wire up online auto-solve in main agent flow
  - Integrate `online.AutoSolve()` when solver confidence < 0.7
  - Add `--online` flag to force online solve

- [ ] Implement Nushell daemon logic
  - YouTube daemon with actual yt-dlp commands
  - Web crawler with scraping logic
  - NLM research daemon
  - Audio generation daemon

### Future Roadmap

- Phase 2: Complete API service wiring
  - Full handler-to-service connections
  - Authentication flow
  - Database migrations applied

- Phase 3: Animation pipeline
  - AI-to-Noon code generator
  - Video rendering pipeline
  - Integration with solver output

- Phase 4: Database population
  - Seed questions from historical BAC exams
  - User management
  - Predictions storage

- Phase 5: CLI enhancement
  - Full solve command with animation
  - Submit command
  - Predict command
  - Leaderboard
