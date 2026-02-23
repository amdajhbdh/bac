# Setup Guide

## Prerequisites

### Required Tools

| Tool | Install Command | Purpose |
|------|-----------------|---------|
| Nushell | `yay -S nushell` | Shell |
| jj | `yay -S jj-bin` | Version control |
| PostgreSQL | `yay -S postgresql` | Database |
| pgvector | `yay -S postgresql-pgvector` | Vector search |
| Ollama | `yay -S ollama` | Local LLM |
| tesseract | `yay -S tesseract` | OCR |
| aichat | `yay -S aichat` | AI CLI |
| yt-dlp | `yay -S yt-dlp` | YouTube |
| texlive | `yay -S texlive-most` | LaTeX |
| Node.js | `yay -S nodejs` | Runtime |
| npm | Included with nodejs | Package manager |
| playwright | `npm install -g playwright` | Browser automation |

## Installation Steps

### 1. Install System Packages
```bash
yay -S postgresql postgresql-libs postgresql-16
yay -S jj-bin
yay -S aichat
yay -S yt-dlp
yay -S texlive-most
```

### 2. Install Node.js Tools
```bash
npm install -g manim-web mermaid-cli katex playwright
```

### 3. Initialize PostgreSQL
```bash
sudo systemctl enable --now postgresql
sudo -u postgres createuser -s med
sudo -u postgres createdb bac_db
sudo -u postgres psql -d bac_db -c "CREATE EXTENSION IF NOT EXISTS vector;"
sudo -u postgres psql -d bac_db -c "CREATE EXTENSION IF NOT EXISTS pg_trgm;"
```

### 4. Setup jj Repository
```bash
cd ~/bac
jj init
jj commit -m "Initial commit"
```

### 5. Create Database Schema
```bash
psql -d bac_db -f sql/schema.sql
```

### 6. Verify Tools
```bash
nlm login --check
ollama list
jj log
```

## Configuration

### Topics
Edit `config/topics.yaml` to add/remove YouTube search topics.

### Environment
Create `.env` file for database credentials (gitignored).

## Next Steps

1. Create NotebookLM notebook: `nlm notebook create "BAC 2026 Master"`
2. Ingest PDFs to database
3. Start daemons
