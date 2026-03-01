# Proposal: Full BAC Unified Platform

## Summary

Implement a complete national AI-powered exam preparation platform for Mauritania's BAC C students, serving 100,000+ students across 2,000+ schools.

## Problem

- No unified digital platform for BAC exam preparation
- Students lack access to AI-powered tutoring
- Historical questions not digitized or searchable
- No question prediction system for exams
- Limited access to educational resources

## Solution

**BAC Unified** - A comprehensive platform providing:
- Intelligent Question Bank with vector similarity search
- AI Tutor with step-by-step solutions and animations
- Resource Aggregation from multiple educational sources
- ML-based Question Predictions
- Offline-first capability with local AI

## Scope

This change implements the complete 2-year roadmap:

### Year 1: Foundation
- Core infrastructure (Go API, PostgreSQL, Garage S3)
- NLM Caching System with PostgreSQL (instead of SQLite)
- User authentication and management
- Question Bank with OCR pipeline
- AI Solver with Noon animations

### Year 2: Scale
- Prediction engine with ML
- Gamification system
- Mobile PWA
- National deployment

## Key Technologies

- **Language**: Go, Rust, TypeScript, Python
- **Database**: PostgreSQL + pgvector (Neon)
- **Storage**: Garage S3 + MinIO
- **AI**: Ollama (local), NotebookLM, Cloud APIs
- **Search**: Elasticsearch, Open Semantic Search
- **Scraping**: Playwright, yt-dlp
- **OCR**: Tesseract, Surya
- **File Bundling**: WRAPX, ARX, Research Object Bundle
- **Auth**: JWT, bcrypt, Redis sessions
- **Monitoring**: Prometheus, Grafana

## Non-Goals

- Hardware token MFA (TOTP only)
- Mobile native apps (PWA first)
- Real-time collaboration
- Blockchain integration (audit only)

## Timeline

- Year 1 (2026): Foundation - 4 quarters
- Year 2 (2027): Scale - 4 quarters

## Success Metrics

| Metric | Year 1 | Year 2 |
|--------|--------|--------|
| Active Users | 10,000 | 100,000 |
| Schools | 100 | 2,000 |
| Questions | 10,000 | 100,000 |
| Cache Hit Rate | 60% | 80% |
