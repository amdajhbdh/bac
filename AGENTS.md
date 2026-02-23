# BAC Hyper-Swarm Agents Rules

## Core Principles

1. **Local First** - Always prioritize local processing (Ollama) before cloud
2. **Privacy** - No data leaves user's control
3. **Offline Capable** - Core functionality works without internet
4. **NLM Maximization** - Exploit ALL NotebookLM features

## Agent Behaviors

### YouTube Daemon

- Search open topics from `config/topics.yaml`
- Download educational content via yt-dlp
- Extract audio for offline study
- Store in `resources/youtube/`

### Web Crawler Daemon

- Scrape Mauritanian education sites: minedu.mr, bacmr.net
- Download PDFs to `resources/web/`
- Run every 12 hours via systemd timer

### OCR Agent (Handwritten)

- Use Tesseract (free, Arabic + French support)
- Never use paid APIs (Google Vision requires billing)
- Store OCR'd text in database
- Location: `resources/handwritten/`

### Database Agent

- Local PostgreSQL + pgvector only
- No cloud services requiring credit card
- Export/import via jj for sync between devices

## NotebookLM (nlm) - FULL FEATURE EXPLOITATION

All nlm commands MUST be used:

### Notebook Management

```bash
nlm notebook create "BAC 2026 Master"     # Create master notebook
nlm notebook list                         # List all notebooks
nlm notebook describe <id>              # Get AI summary
nlm notebook query <id> <question>       # Chat with sources
nlm notebook rename <id> <new-name>     # Rename
nlm notebook delete <id>                 # Delete
```

### Source Management

```bash
nlm source add <notebook-id> <file>     # Add PDF/image/audio
nlm source list <notebook-id>           # List sources
nlm source describe <notebook-id> <id> # AI summary
nlm source content <notebook-id> <id>   # Raw content
nlm source delete <notebook-id> <id>    # Remove source
nlm source sync                         # Sync Drive sources
```

### Content Generation (USE ALL!)

```bash
# Audio Overview (Podcast)
nlm audio create <notebook-id>                    # Generate podcast
nlm audio create <notebook-id> --topic "dérivées" # Topic-specific
nlm download audio <notebook-id>                  # Download MP3

# Quiz Generation
nlm quiz create <notebook-id>                     # Auto-generate quiz
nlm quiz create <notebook-id> --topic "intégrales"

# Flashcards
nlm flashcards create <notebook-id>
nlm flashcards create <notebook-id> --topic "chemistry"

# Mind Maps
nlm mindmap create <notebook-id>
nlm mindmap create <notebook-id> --topic "biology"

# Video Overview
nlm video create <notebook-id>

# Infographics
nlm infographic create <notebook-id>

# Reports
nlm report create <notebook-id>

# Data Tables
nlm data-table create <notebook-id>

# Slide Decks
nlm slides create <notebook-id>
```

### Research & Discovery

```bash
nlm research start <query>              # Start research
nlm research start <query> --mode deep  # Deep research
nlm research status                      # Check progress
nlm research import                      # Import results
```

### Export & Sharing

```bash
nlm export gdocs <notebook-id>         # Export to Google Docs
nlm export sheets <notebook-id>         # Export to Sheets
nlm share <notebook-id>                 # Share notebook
nlm download <notebook-id>              # Download all artifacts
```

### Chat Configuration

```bash
nlm chat configure                      # Configure AI behavior
nlm chat start <notebook-id>            # Interactive chat
```

## Topic Search

Open topics stored in `config/topics.yaml`:

- User can add: `bac config add-topic "new topic"`
- Topics auto-discover new content

Default topics:

- BAC mathématiques Terminale
- BAC physics derivatives integrals
- Terminale C chimie organique
- SVT Terminale biology
- BAC français dissertation

## Database Rules

1. **Primary**: Local PostgreSQL + pgvector
2. **No Neon/Cloud DB** - requires credit card
3. **Sync**: Export to SQL file, import on other devices
4. **Backup**: Daily via jj commit

## Visualization

Use only free tools:

- manim-web (npm) - Math animations
- mermaid-cli (npm) - Diagrams
- katex (npm) - LaTeX rendering
- Tesseract - OCR (Arabic + French)

NO paid services (Google Cloud Vision requires billing)

## History Management (jj)

All changes tracked via jj:

```bash
jj init                                  # Initialize
jj commit -m "Add new PDF"               # Save changes
jj log                                   # View history
jj bookmark add checkpoint-v1            # Create checkpoint
jj git export                            # Export to git
jj git import                            # Import from git
```

## Progress Tracking

All progress in `progress.md`:

- Daily log entries
- Checkbox completion
- Notes section

## Error Handling

1. NLM unavailable → Use local Ollama
2. Internet down → Work offline, sync later
3. Database error → Check PostgreSQL status
4. OCR fails → Try different image preprocessing

## Implementation Priority

1. Create AGENTS.md ✓
2. Create progress.md
3. Create docs/
4. Initialize jj repo
5. Install tools
6. Setup PostgreSQL
7. Create NLM notebook
8. Ingest PDFs
9. Start daemons

---

*Last Updated: 2026-02-23*
*Version: 1.0*
