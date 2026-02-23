# Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                    BAC HYPER-SWARM 2026                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌──────────────────────────────────────────────────────────────┐  │
│   │              ORCHESTRATION LAYER (Nushell)                  │  │
│   │                    bac.nu - Main CLI                         │  │
│   └──────────────────────────────────────────────────────────────┘  │
│                                │                                    │
│   ┌────────────────────────────┼────────────────────────────────┐  │
│   │                    AGENT LAYER                               │  │
│   │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐  │  │
│   │  │Ingest  │ │Query   │ │Viz     │ │Exam    │ │NLM     │  │  │
│   │  │Agent   │ │aichat  │ │Agent   │ │Agent   │ │Research│  │  │
│   │  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘  │  │
│   │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐  │  │
│   │  │YouTube │ │Web     │ │OCR     │ │Audio   │ │Flash   │  │  │
│   │  │Daemon  │ │Crawler │ │Tesseract│ │Daemon  │ │Daemon  │  │  │
│   │  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘  │  │
│   └──────────────────────────────────────────────────────────────┘  │
│                                │                                    │
│   ┌────────────────────────────┼────────────────────────────────┐  │
│   │                    INFERENCE LAYER                           │  │
│   │  ┌────────────┐  ┌────────────┐  ┌─────────────────┐     │  │
│   │  │Local       │  │Cloud       │  │  NotebookLM     │     │  │
│   │  │Ollama      │  │kimi/glm/   │  │  (Audio/Quiz)  │     │  │
│   │  │llama3.2:3b │  │minimax    │  │                 │     │  │
│   │  └────────────┘  └────────────┘  └─────────────────┘     │  │
│   └──────────────────────────────────────────────────────────────┘  │
│                                │                                    │
│   ┌────────────────────────────┴────────────────────────────────┐  │
│   │                    DATABASE LAYER                            │  │
│   │              PostgreSQL + pgvector (local)                   │  │
│   └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
│   ┌──────────────────────────────────────────────────────────────┐  │
│   │                    HISTORY MANAGER (jj)                        │  │
│   └──────────────────────────────────────────────────────────────┘  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Components

### Orchestration Layer
- **bac.nu**: Main CLI entry point in Nushell
- Handles all agent coordination

### Agent Layer
- **Ingest Agent**: Rust-based PDF processor
- **Query Agent**: aichat with RAG
- **Viz Agent**: manim-web, mermaid, katex
- **Exam Agent**: Test generation
- **NLM Research**: NotebookLM integration

### Daemon Layer
- **YouTube Daemon**: Search & download (6h interval)
- **Web Crawler**: Scrape education sites (12h interval)
- **NLM Research Daemon**: Auto-discover sources
- **Audio Daemon**: Generate podcasts

### Inference Layer
- **Local**: Ollama (llama3.2:3b)
- **Cloud**: kimi, glm, minimax via ollama
- **NotebookLM**: Audio, quiz, flashcards, mindmaps

### Database Layer
- **PostgreSQL**: Relational data
- **pgvector**: Vector similarity search
- **Tesseract**: OCR for handwritten notes

## Data Flow

1. PDF → Ingest Agent → PostgreSQL + pgvector
2. User Query → Query Agent → RAG → Answer
3. Topic Search → YouTube Daemon → Download → Resources
4. Web Discovery → Crawler → Resources
5. Study Materials → NLM → Audio/Quiz/Flashcards

## Version Control

All changes tracked via jj:
- Config changes
- Script updates
- Database exports
- Documentation

## Sync Mechanism

Since no cloud DB:
1. Export: `pg_dump bac_db > backup.sql`
2. Import: `psql bac_db < backup.sql`
3. Share via USB/cloud storage
4. Track exports in jj
