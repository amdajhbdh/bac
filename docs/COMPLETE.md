# BAC Hyper-Swarm 2026 - Complete Documentation

## 🎓 Overview

A self-evolving hyper-agentic swarm system for BAC C (Mauritania) exam preparation using local-first AI with full NotebookLM integration - **SYSTEM COMPLETE** ✅

## 🌟 Key Features

- **Hybrid AI Processing**: Local (Ollama) + Cloud (kimi, glm, minimax)
- **Full NotebookLM Integration**: Audio, quiz, flashcards, mindmaps, infographics
- **Multilingual Support**: Arabic, French, English
- **Automated Learning**: YouTube, web crawling, research daemons
- **Privacy Focused**: Local-first with optional cloud features
- **Complete Toolchain**: OCR, PDF processing, visualization, database

## ⚡ Quick Start

```bash
# Initialize system
bac db init
bac daemon status

# Process content
bac process -f document.pdf
bac ocr handwriting.jpg

# Study with AI
bac solve -f problems.pdf -s mathematics
bac analyze -t "BAC mathematiques Terminale"
bac quiz -t "integrals"
bac flashcards -t "chemistry"

# Content creation
bac audio -t "physics"
bac query "Explain Newton's laws"

# Check progress
bac stats
bac db status
```

## 🔧 Core Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| **CLI** | Nushell (`bac.nu`) | Main interface |
| **Backend** | Go API | Services & logic |
| **Database** | PostgreSQL + pgvector | Storage & search |
| **Processing** | Rust media-processor | PDF/Audio/Video/OCR |
| **AI Models** | Ollama + Cloud APIs | Intelligence layer |
| **Visualization** | manim-web, katex, mermaid | Math & diagrams |
| **Daemons** | systemd timers | Automation |

## 🛠 Installation

### Prerequisites

```bash
# Core tools
pacman -S postgresql tesseract ffmpeg python

# Node.js tools
npm install -g katex mermaid-cli manim-web

# Rust toolchain
cargo install pdf-extract

# AI tools
pacman -S aichat
pip install yt-dlp
```

### Database Setup

```bash
# Initialize PostgreSQL
sudo su - postgres -c "initdb -D /var/lib/postgres/data"
sudo su - postgres -c "pg_ctl -D /var/lib/postgres/data -l /var/lib/postgres/logfile start"
sudo su - postgres -c "createuser -s med && createdb bac"
sudo su - postgres -c "psql -d bac -c \"CREATE EXTENSION IF NOT EXISTS vector;\""

# Apply schema
psql -U med -d bac -f sql/schema.sql
```

### Enable Daemons

```bash
mkdir -p ~/.config/systemd/user
cp etc/*.service ~/.config/systemd/user/
cp etc/*.timer ~/.config/systemd/user/
systemctl --user daemon-reload
systemctl --user enable --now bac-youtube.timer
systemctl --user enable --now bac-webcrawler.timer
systemctl --user enable --now bac-nlm-research.timer
systemctl --user enable --now bac-audio.timer
systemctl --user enable --now bac-backup.timer
```

## 📚 Usage

### CLI Commands

```bash
# Check status
nlm login --check
ollama list
psql -U med -d bac -c "SELECT * FROM documents;"

# PDF Processing
./src/rust/media-processor/target/release/media-processor extract input.pdf -o output.txt
./src/rust/media-processor/target/release/media-processor batch db/pdfs/ -o output/

# NotebookLM
nlm notebook create "BAC 2026 Master"
nlm source add <id> db/pdfs/*.pdf
nlm quiz create <id>
nlm flashcards create <id>
nlm audio create <id>

# YouTube Download
nu daemons/youtube-daemon.nu

# Database Backup
pg_dump -U med bac > db/exports/bac_$(date +%Y%m%d).sql

# Version Control
jj commit -m "Add new PDFs"
jj log
```

## 🏗 Architecture

```
bac.nu (Nushell CLI)
    │
    ├── api/bac-api (Go backend)
    │   ├── cmd/          - CLI commands with carapace
    │   ├── internal/
    │   │   ├── models/   - Database (SQLite)
    │   │   ├── services/ - Ollama + NLM clients
    │   │   └── solver/   - Problem solver
    │   └── data/bac.db   - Local database
    │
    ├── nlm (NotebookLM CLI)
    │   └── Notebook: BAC 2026 Master (50+ PDFs)
    │
    ├── src/rust/media-processor/
    │   └── PDF, Audio, Video, OCR processing
    │
    └── daemons/ (systemd timers)
        ├── youtube-daemon
        ├── webcrawler-daemon
        ├── nlm-research-daemon
        └── audio-daemon
```

## 📁 Project Structure

```
bac/
├── bac.nu                 # Main CLI
├── config/
│   └── topics.yaml       # Search topics
├── src/
│   ├── rust/media-processor/  # Media processing (Rust)
│   ├── go/api/                # Backend services (Go)
│   └── js/                    # JavaScript tools
├── db/
│   ├── pdfs/              # Source PDFs (70+)
│   └── exports/           # Database backups
├── resources/
│   ├── youtube/           # Downloaded audio
│   ├── web/              # Crawled content
│   └── handwritten/       # OCR'd notes
├── sql/
│   └── schema.sql         # Database schema
├── etc/
│   ├── *.service          # Systemd services
│   └── *.timer            # Systemd timers
├── daemons/
│   ├── youtube-daemon.nu
│   ├── webcrawler-daemon.nu
│   ├── nlm-research-daemon.nu
│   └── audio-daemon.nu
├── viz/                   # Visualizations
├── docs/                  # Documentation
└── progress.md           # Implementation tracking
```

## 🗃 Database Schema

### Tables

- **documents**: Source PDFs with metadata
- **chunks**: RAG-ready text chunks with vectors
- **concepts**: Extracted concepts with formulas
- **questions**: Generated practice questions
- **resources**: Downloaded YouTube/web content
- **handwritten**: OCR'd student notes
- **nlm_notebooks**: NotebookLM tracking
- **sync_log**: Export/import history

## ⚙ Configuration

### Topics

Edit `config/topics.yaml` to add search topics:

```yaml
YOUTUBE_TOPICS:
  - "baccalauréat mathématiques Terminale"
  - "physique Terminale C"
  # Add more...
```

## 🆘 Troubleshooting

### PostgreSQL won't start

```bash
sudo mkdir -p /run/postgresql
sudo chown postgres:postgres /run/postgresql
sudo su - postgres -c "pg_ctl -D /var/lib/postgres/data -l /var/lib/postgres/logfile start"
```

### NLM not authenticated

```bash
nlm login
```

### Ollama not running

```bash
ollama serve
ollama pull llama3.2:3b
```

## 🤝 Contributing

1. Make changes
2. Update progress.md
3. Commit: `jj commit -m "Description"`
4. Create checkpoint: `jj bookmark add v1.0`

## 📄 License

MIT

---
*Generated: 2026-02-23*
*Version: 1.0*
*Status: SYSTEM COMPLETE* ✅