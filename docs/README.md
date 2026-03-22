# BAC Study - Project Documentation

## Overview

BAC Study is a comprehensive study aid system for Mauritanian BAC (Baccalaureate) exams, featuring:
- OCR-powered document extraction
- RAG-based question answering
- Multi-provider AI services
- Cloudflare Workers deployment

## Project Structure

```
bac/
├── Cargo.toml           # Rust workspace config
├── config/             # Configuration files
├── docs/               # Documentation
├── domain_data/        # RAG document storage
├── scripts/            # Python & shell scripts
│   ├── cli/          # CLI tools
│   ├── ocr/          # OCR processing
│   ├── rag/          # RAG pipeline
│   ├── upload/       # Data upload
│   └── utils/        # Utilities
└── src/               # Source code
    ├── api/          # Go API service
    ├── cloudflare/   # Cloudflare Workers
    ├── tools/        # AI tools (gemini, vector, etc.)
    ├── services/     # Backend services
    ├── ocr/          # OCR service
    ├── kilocode/     # Core code
    └── web/          # Web frontend
```

## Quick Start (Termux)

### 1. Install Dependencies

```bash
# Update and install base packages
pkg update && pkg upgrade -y
pkg install -y git curl wget python nodejs rust

# Clone repository
git clone https://github.com/amdajhbdh/bac.git
cd bac

# Run setup
bash scripts/setup-termux.sh
```

### 2. Python Environment

```bash
# Create virtual environment
python -m venv .venv
source .venv/bin/activate

# Install Python dependencies
pip install --break-system-packages -r scripts/requirements.txt
```

### 3. Run Services

```bash
# Start RAG pipeline
python scripts/rag_pipeline.py

# Start OCR server
python scripts/ocr/local-ocr-server.py

# Start API
cd src/api && go run main.go
```

## Available Scripts

### OCR Scripts (`scripts/ocr/`)
- `ocr_providers.py` - 18-provider OCR chain
- `test-ocr-chain.py` - Test OCR providers
- `local-ocr-server.py` - FastAPI OCR server

### RAG Scripts (`scripts/rag/`)
- `rag_pipeline.py` - Hybrid BM25 + vector retrieval
- `eval_rag.py` - Evaluate RAG performance

### CLI Scripts (`scripts/cli/`)
- `bac-cli.py` - Main CLI interface
- `bac-typer.py` - Type assistant

### Upload Scripts (`scripts/upload/`)
- `upload-to-upstash.py` - Upload to Upstash
- `upload-batch-worker.py` - Batch upload via worker

## Environment Variables

Create `.env` file:

```env
# API Keys
ANTHROPIC_API_KEY=your_key
YANDEX_FOLDER_ID=your_folder
YANDEX_IAM_TOKEN=your_token

# OCR Providers
ABBYY_APP_ID=your_id
ABBYY_PASSWORD=your_pass
CLOUDMERSIVE_API_KEY=your_key

# Database
UPSTASH_REST_URL=your_url
UPSTASH_REST_TOKEN=your_token
```

## Cloudflare Deployment

```bash
# Deploy OCR worker
cd src/cloudflare/worker
wrangler deploy

# Set secrets
wrangler secret put ANTHROPIC_API_KEY
```

## API Endpoints

| Service | Endpoint | Description |
|---------|----------|-------------|
| API | `localhost:8080` | Go API server |
| OCR | `localhost:8000` | OCR FastAPI |
| RAG | `localhost:5000` | RAG pipeline |

## Tech Stack

- **Backend**: Go, Rust, Python
- **AI/ML**: LangChain, Ollama, Claude, Gemini
- **OCR**: PaddleOCR, EasyOCR, 18+ providers
- **Storage**: Upstash (Redis), ChromaDB
- **Deployment**: Cloudflare Workers, Fly.io

## License

MIT
