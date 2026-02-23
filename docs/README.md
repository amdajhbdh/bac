# BAC Hyper-Swarm 2026

A self-evolving hyper-agentic swarm system for BAC C exam preparation.

## Features

- **Local LLM**: Ollama with llama3.2:3b
- **Cloud LLM**: kimi, glm, minimax via ollama
- **NotebookLM**: Full feature exploitation (audio, quiz, flashcards, mindmaps)
- **PostgreSQL + pgvector**: Local vector database
- **Tesseract OCR**: Free Arabic + French support
- **YouTube Daemon**: Auto-download educational content
- **Web Crawler**: Scrape Mauritanian education sites
- **Visualization**: manim-web, mermaid, katex
- **Version Control**: jj (Jujutsu)

## Quick Start

```bash
# Check status
nlm login --check
ollama list
jj log

# Create notebook
nlm notebook create "BAC 2026 Master"

# Add PDFs
nlm source add <notebook-id> db/pdfs/*.pdf
```

## Documentation

See `docs/` for detailed documentation.

## Commands

See `docs/API.md` for all CLI commands.
