# API Documentation

## Main CLI (bac.nu)

### Database Commands
```bash
bac db init           # Initialize database
bac db export <file>  # Export database
bac db import <file>  # Import database
```

### Ingest Commands
```bash
bac ingest <pdf>              # Ingest single PDF
bac ingest-all                # Ingest all PDFs in db/pdfs/
bac ocr <image>               # OCR handwritten image
```

### Query Commands
```bash
bac ask <question>            # Query knowledge base
bac ask --nlm <question>     # Query via NotebookLM
bac ask --ollama <question>  # Query via Ollama
```

### NLM Commands
```bash
bac nlm create-notebook <name>   # Create notebook
bac nlm add-sources <notebook>   # Add PDFs to notebook
bac nlm audio <notebook>         # Generate audio
bac nlm quiz <notebook>          # Generate quiz
bac nlm flashcards <notebook>    # Generate flashcards
bac nlm mindmap <notebook>      # Generate mindmap
bac nlm research <query>         # Start research
```

### Daemon Commands
```bash
bac daemon youtube start    # Start YouTube daemon
bac daemon web start        # Start web crawler
bac daemon nlm start       # Start NLM research
bac daemon status           # Check daemon status
```

### Viz Commands
```bash
bac viz animate <topic>     # Generate animation
bac viz formula <latex>     # Render formula
bac viz plot <function>     # Plot function
```

### Config Commands
```bash
bac config list-topics      # List search topics
bac config add-topic <topic> # Add topic
bac config remove-topic <topic> # Remove topic
```

## Individual Tools

### nlm (NotebookLM)
```bash
nlm notebook list
nlm notebook create "name"
nlm notebook query <id> <question>
nlm source add <notebook-id> <file>
nlm audio create <notebook-id>
nlm quiz create <notebook-id>
nlm flashcards create <notebook-id>
nlm research start <query>
```

### ollama
```bash
ollama list
ollama run llama3.2:3b <prompt>
ollama run kimi-k2.5:cloud <prompt>
```

### aichat
```bash
aichat --rag db/pdfs "question"
aichat -m llama3.2:3b
```

### jj (Version Control)
```bash
jj log
jj commit -m "message"
jj bookmark add <name>
jj git export
jj git import
```

### tesseract (OCR)
```bash
tesseract input.png output -l ara+fra
```

### yt-dlp
```bash
yt-dlp --search "bac mathematics" --max-downloads 5
```
