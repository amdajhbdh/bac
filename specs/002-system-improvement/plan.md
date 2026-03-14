# Implementation Plan: System Improvement

**Branch**: `002-system-improvement` | **Date**: 2026-03-11 | **Spec**: specs/002-system-improvement/spec.md

## Summary

Build core backend infrastructure for the BAC Unified learning system. This includes Go backend with PostgreSQL, OCR processing for study materials, hybrid search, and AI tutoring capabilities. The system prioritizes offline-first operation with local AI inference.

## Technical Context

**Language/Version**: Go 1.25, Rust (gateway/ocr), React 18 (web), TypeScript (workers)  
**Primary Dependencies**: Gin (Go web), pgx (DB), Tesseract/Surya (OCR), Elasticsearch + pgvector (search), Ollama (AI)  
**Storage**: PostgreSQL + pgvector (main), Redis (cache), Turso (NLM cache)  
**Testing**: Go test, Vitest, Cargo test  
**Target Platform**: Linux server, self-hosted  
**Project Type**: web-service/cli/desktop-app  
**Performance Goals**: Search <2s, OCR 10-page PDF <60s, AI response streaming  
**Constraints**: Offline-first, local AI priority, no paid APIs  
**Scale**: Single user to classroom use

## Constitution Check

Per the Constitution (Learning-First, Offline-First, Local AI Priority, OER, Privacy):
- [x] All features serve learning outcomes
- [x] Core features work offline
- [x] Local AI (Ollama) prioritized over cloud
- [x] Uses open educational resources
- [x] No user tracking/telemetry

## Project Structure

### Source Code

```
src/
├── agent/              # Go CLI/API (70+ files)
│   ├── cmd/           # Entry points
│   ├── internal/      # Core packages
│   │   ├── db/        # PostgreSQL + pgvector
│   │   ├── ocr/       # Tesseract/Surya
│   │   ├── search/    # Elasticsearch hybrid
│   │   ├── ai/        # Ollama + NLM
│   │   └── ...
│   └── pkg/
├── api/               # REST API (Gin)
│   ├── internal/
│   │   ├── handlers/  # HTTP endpoints
│   │   ├── services/  # Business logic
│   │   └── models/    # Data models
│   └── docs/          # Swagger
├── gateway/            # Rust Axum
├── web/               # React/Bun
└── cloudflare-worker/  # TypeScript Workers
```

### Database Schema

- `sql/schema.sql` - PostgreSQL with pgvector
- Tables: materials, questions, answers, embeddings

## Complexity Tracking

| Area | Decision | Rationale |
|------|----------|-----------|
| DB | PostgreSQL + pgvector | Vector search for semantic matching |
| Search | Elasticsearch + pgvector RRF | Hybrid full-text + semantic |
| AI | Ollama → NLM → Cloud | Offline-first priority |
| OCR | Tesseract + Surya | Free, local-only |

All decisions align with constitution principles.
