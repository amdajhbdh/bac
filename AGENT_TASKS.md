# Agent Task Board - BAC Unified

Multi-instance coordination for AI coding agents.

---

## Instance Mapping

| Instance | Tool | Primary Focus | Status |
|----------|------|---------------|--------|
| **opencode** | opencode | Gateway/REST API | Active |
| **kiro** | kiro-cli | Agent CLI (Go) | Standby |
| **kilocode** | kilocode | OCR Service (Rust) | Active |

---

## Current Sprint

### Phase 1: OCR Integration

**Started**: 2026-03-06
**Goal**: Complete full OCR system with Surya + PDF support

---

## Task Board

### opencode (Gateway)

| ID | Task | Status | Notes |
|----|------|--------|-------|
| G-001 | Create AGENT_TASKS.md | ✅ Done | This file |
| G-002 | Create Rust gateway structure | ✅ Done | src/gateway/ created |
| G-003 | Implement RAG chat modes | ✅ Done | Query/Chat/Agent/Auto |
| G-004 | Implement animation queue | ✅ Done | queue.rs + bridge.rs |
| G-005 | Fix Manim podman permissions | ✅ Done | --user flag + dirs |
| G-006 | Connect OCR client to gateway | ✅ Done | HTTP client added |
| G-007 | Create HTTP endpoints | ✅ Done | /ocr, /health |
| G-008 | Add gateway integration tests | ✅ Done | 4 tests added |
| G-009 | Fix ChatMode test failures | ✅ Done | Word boundary matching |
| G-010 | Create ecosystem plan docs | ✅ Done | .opencode/plans/ |
| G-011 | Add Go gateway tests | ✅ Done | 4 tests in agent/ |
| G-012 | Phase 5: RAG→PostgreSQL | ✅ Done | sqlx + pgvector queries |
| G-013 | Add embedding generation (Ollama) | ✅ Done | generate_embedding() |
| G-014 | Add hybrid/range search queries | ✅ Done | hybrid_search, range_search |

### kiro (Agent CLI - Go)

| ID | Task | Status | Notes |
|----|------|--------|-------|
| A-001 | Connect to OCR service | ✅ Done | HTTP client ready |
| A-002 | Add Gateway client tests | ✅ Done | 4 integration tests |
| A-003 | Enhance solver with animations | ✅ Done | SolveWithAnimation + tests |
| A-004 | Update daemons for new services | ⏳ Pending | Nushell scripts |

### kilocode (OCR Service - Rust)

| ID | Task | Status | Notes |
|----|------|--------|-------|
| O-001 | Create pipeline.rs | ✅ Done | Tesseract fallback |
| O-002 | Add Surya OCR integration | ✅ Done | Python subprocess |
| O-003 | Add PDF processing | ✅ Done | pdftoppm + multi-page |
| O-004 | Create HTTP server | ✅ Done | Axum endpoints |
| O-005 | Add French/Arabic models | ✅ Done | Language support added |

---

## Coordination Rules

### Before Starting Work

1. ✅ Check this file
2. ✅ Verify task is unclaimed or assigned to you
3. ✅ Update status when starting

### While Working

1. ❌ Never modify another agent's active files
2. ❌ Don't commit if tests fail
3. ✅ Document any blocking issues

### After Session

1. ✅ Update task status in this file
2. ✅ Note any incomplete work
3. ✅ Commit with jj if meaningful progress

---

## File Ownership

| Directory | Owner |
|-----------|-------|
| `src/gateway/` | opencode |
| `src/agent/` | kiro |
| `src/ocr-service/` | kilocode |
| `daemons/` | kiro |
| `lib/` | kiro |
| `src/noon/` | Shared (animation) |

---

## Communication

- Update this file after each session
- Use TODO comments in code for follow-up
- Blockers → Add to Issues section below

---

## Issues / Blockers

- [x] OCR service needs HTTP endpoint before gateway can connect (Done)
- [x] Surya requires Python environment setup (Workaround: fallback to Tesseract)
- [x] Need test images for OCR validation (Use db/pdfs/)
- [ ] LSP errors in src/agent/internal/cache/ - need fixing
- [ ] ocr-service reqwest blocking - need fixing

---

## Quick Commands

```bash
# Check your tasks
grep -E "^## " AGENT_TASKS.md
grep "opencode" AGENT_TASKS.md | grep "⏳"

# Update status (use these marks)
✅ Done
⏳ Pending/In Progress
❌ Blocked
```

---

*Last Updated: 2026-03-06*
*Maintained by: All agents*
