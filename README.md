# BAC Knowledge System

AI-native knowledge management system using aichat, Gemini, and Rust tool crates for semantic search, content generation, and knowledge graph visualization.

## Features

- 📚 **RAG-powered Search** — Semantic similarity search across all notes using pgvector/HNSW
- 🤖 **AI Content Generation** — Generate notes with LaTeX support via Gemini
- 🔗 **Auto-linking** — Connect related concepts using wikilinks and MOCs
- 📊 **Knowledge Graphs** — Visualize note relationships as interactive graphs
- ☁️ **Cloud Integration** — Process via Cloud Shell (Cloudflare/Garage)
 - 🔍 **OCR Pipeline** — 18-provider cascading fallback chain for French/Arabic exam images

## Quick Start

### 1. Install Dependencies

```bash
# Rust toolchain
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# aichat CLI
curl -fsSL https://raw.githubusercontent.com/sigoden/aichat/main/install.sh | sh

# Node.js (for scripts)
npm install
```

### 2. Configure

```bash
# Copy environment template
cp .env.example .env

# Edit with your API keys and connections
# Required: GEMINI_API_KEY, database URLs
```

### 3. Build Tools

```bash
# Build all Rust services
cargo build --release

# Or build individually
cargo build -p gemini-tools --release
cargo build -p vector-tools --release
```

### 4. Start Services

```bash
# Start all services
./scripts/bac-agent-daemon.sh

# Or start individually (see Services section)
cargo run -p gemini-tools --release
cargo run -p vector-tools --release
```

### 5. Use aichat

```bash
# Start with BAC agent configuration
aichat --config config/aichat/config.yaml
```

## Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                         aichat (CLI / Termux)                         │
│  Agent: bac-knowledge  │  RAG: resources/notes/  │  OCR: 18 providers │
└─────────────────────────┬────────────────────────┬───────────────────┘
                          │                        │
          ┌───────────────┼───────────────┬────────┼──────────────────┐
          │               │               │        │                  │
          ▼               ▼               ▼        ▼                  ▼
     ┌───────────┐  ┌───────────┐  ┌──────────┐ ┌──────────┐    ┌─────────┐
     │  gemini   │  │  vector   │  │  vault   │ │  cloud   │    │  Telegram│
     │  -tools   │  │  -tools   │  │  -tools  │ │  -tools  │    │   Bot   │
     │   :3001   │  │   :3002   │  │  :3003   │ │  :3004   │    │   :8081 │
     └───────────┘  └───────────┘  └──────────┘ └──────────┘    └─────────┘
          │               │               │        │                  │
          └───────────────┴───────────────┼────────┼──────────────────┘
                                          │        │
                             ┌─────────────┴────────┴──────────┐
                             │    api-gateway (:8080)           │
                             └─────────────────────────────────┘
                                              │
                          ┌───────────────────┼───────────────────┐
                          │                   │                   │
                          ▼                   ▼                   ▼
                   ┌─────────────┐      ┌────────────┐     ┌────────────┐
                   │ Cloudflare  │      │ Cloud Shell│     │  Local     │
                   │  Worker     │      │  Ollama    │     │  Python    │
                   │ Mistral OCR │      │ LFM2.5-VL  │     │ Paddle/Tess│
                   └─────────────┘      │ Qwen2.5-VL │     └────────────┘
                                        │ LLaVA      │
                                        └────────────┘
```

## OCR Pipeline

The system uses an **18-provider cascading fallback chain** for maximum reliability on French and Arabic exam images.

### Architecture: 3-Tier Fallback

```
Input Image
    │
    ▼
┌─────────────────────────────────────────────────────────┐
│ TIER 1: CLOUD (instant, ~$0.002/page)                    │
│  1. Mistral OCR      ← Primary, best French/Arabic      │
│  2. Workers AI       ← Free, built into Worker          │
│  3. OCR.space        ← Free, no key needed              │
│  4. Yandex Vision    ← Free trial, 1K ops/3mo           │
└────────────────────────────┬────────────────────────────┘
                             │ (fallback if cloud fails)
                             ▼
┌─────────────────────────────────────────────────────────┐
│ TIER 2: CLOUD SHELL (0-cost, needs setup)              │
│  5. LFM2.5-VL  (Liquid AI, 696MB, OCRBench 822)  ⭐   │
│  6. Qwen2.5-VL (Multilingual, Ollama)                  │
│  7. LLaVA       (Stable vision model)                  │
│  8. GLM-OCR     (Best for math equations)              │
└────────────────────────────┬────────────────────────────┘
                             │ (fallback if shell unavailable)
                             ▼
┌─────────────────────────────────────────────────────────┐
│ TIER 3: LOCAL PYTHON (free, needs pip install)         │
│  9.  PaddleOCR     (Best open-source accuracy)          │
│ 10. EasyOCR        (80+ languages)                      │
│ 11. Tesseract      (Fast, already available)            │
│ 12. Puter.js      (Browser-based, no install)           │
│ ... + 4 more cloud providers for edge cases             │
└─────────────────────────────────────────────────────────┘
                             │
                             ▼
                      Extracted Text
                             │
                             ▼
                      OCR Corrections
                             │
                             ▼
                    Structured Notes (RAG)
```

### Performance Benchmarks

Tested on 9 real BAC exam images (Math, Physics, Chemistry, French):

| Subject | Image | Provider | Chars | Latency |
|---------|-------|----------|-------|---------|
| Math | math_page-1.jpg | Mistral | 2,401 | 3.5s |
| Math | math_page-2.jpg | Mistral | 3,433 | 4.7s |
| Physics | physique_page-1.jpg | Mistral | 3,187 | 2.8s |
| Physics | physique_page-2.jpg | Mistral | 4,202 | 3.4s |
| French | francais_page-1.jpg | Mistral | 2,842 | 2.6s |
| Arabic Math | arabic_math_page-1.jpg | Mistral | 960 | 1.7s |
| Arabic Math | arabic_math_page-2.jpg | Mistral | 536 | 2.5s |
| Chemistry | chimie_page-1.jpg | Mistral | 224 | 1.2s |
| Chemistry | chimie_page-2.jpg | Mistral | 39 | 1.2s |

**Summary:** 9/9 images successful, **17,824 total chars**, ~2.5s average latency. Mistral succeeded on all images, OCR.space was not needed (chain fallback). Cost per page: **$0.002**.

## Services

| Service | Port | Purpose |
|---------|------|---------|
| api-gateway | 8080 | Main API gateway |
| ocr | 3000 | OCR processing (Tesseract) |
| ocr-chain | - | 18-provider OCR chain via scripts |
| ocr-api | 8082 | Local OCR server (FastAPI, optional) |
| gemini-tools | 3001 | Gemini API wrapper |
| vector-tools | 3002 | pgvector semantic search |
| vault-tools | 3003 | Obsidian vault operations |
| cloud-tools | 3004 | Cloud Shell integration |
| graph-tools | 3005 | Knowledge graph generation |

## Usage Examples

### Search notes
```
> search for electric field concepts
> find notes about thermodynamics
```

### Generate note
```
> write a note about thermodynamics
> create a summary of photosynthesis
```

### Process files
```
> process new files in 03-Resources
> extract text from test.jpg
```

### Knowledge graph
```
> show connections for quantum mechanics
> generate graph for neural networks
```

## Project Structure

```
bac/
├── src/
│   ├── api/              # Main API gateway
│   ├── gemini-tools/      # Gemini API wrapper
│   ├── vector-tools/     # pgvector operations
│   ├── vault-tools/      # Obsidian vault ops
│   ├── cloud-tools/      # Cloud Shell integration
│   ├── graph-tools/      # Knowledge graph
│   └── ocr/              # OCR processing
├── config/
│   └── aichat/           # aichat configuration
│       ├── config.yaml    # Main config
│       ├── agents/       # Agent definitions
│       └── rag-template.txt
├── resources/
│   └── notes/            # Obsidian vault
├── scripts/              # Automation scripts
│   ├── ocr_providers.py      # 18-provider OCR chain ⭐
│   ├── ocr_smart_select.py   # Smart provider selection
│   ├── ocr_preprocess.py      # Image preprocessing
│   ├── ocr_config.toml       # Configuration template
│   ├── cloud-shell-ocr.sh    # Cloud Shell Ollama setup
│   ├── cloudflared-tunnel.sh # Cloudflare Tunnel
│   ├── local-ocr-server.py   # FastAPI server
│   ├── test-ocr-chain.py     # Test harness
│   ├── bac-telegram-bot.py   # Telegram bot
│   └── bac-agent-daemon.sh   # Service daemon
├── tests/bac-images/    # Test images for OCR
├── sql/                  # Database schemas
└── docs/                 # Documentation
    ├── API.md
    ├── TROUBLESHOOTING.md
    ├── ocr-api-keys-setup.md  # API key guide ⭐
```

## OCR Scripts

### Quick OCR (Mistral, works now)

```bash
# Single image
MISTRAL_API_KEY="your-key" python3 scripts/ocr_providers.py image.jpg

# Chain with fallback
MISTRAL_API_KEY="your-key" python3 scripts/ocr_providers.py image.jpg -p mistral -p ocrspace

# Parallel mode (fastest)
MISTRAL_API_KEY="your-key" python3 scripts/ocr_providers.py image.jpg --parallel

# JSON output
MISTRAL_API_KEY="your-key" python3 scripts/ocr_providers.py image.jpg -j

# List all providers
python3 scripts/ocr_providers.py --list
```

### Smart Selection (auto-detects best provider)

```bash
# Analyze image and select optimal provider
python3 scripts/ocr_smart_select.py image.jpg

# Analyze only (no OCR)
python3 scripts/ocr_smart_select.py image.jpg --analyze-only

# Specific modes
python3 scripts/ocr_smart_select.py image.jpg --mode math
python3 scripts/ocr_smart_select.py image.jpg --mode speed
python3 scripts/ocr_smart_select.py image.jpg --mode french_arabic
python3 scripts/ocr_smart_select.py image.jpg --mode cost
```

### Image Preprocessing

```bash
# Enhance contrast
python3 scripts/ocr_preprocess.py input.jpg output.jpg --enhance

# Binary (B&W, best for tesseract)
python3 scripts/ocr_preprocess.py input.jpg output.jpg --binary --deskew

# Provider-specific optimization
python3 scripts/ocr_preprocess.py input.jpg output.jpg --smart tesseract
python3 scripts/ocr_preprocess.py input.jpg output.jpg --smart mistral
```

### Telegram Bot

```bash
# Start bot (receives photos and extracts text)
python3 scripts/bac-telegram-bot.py

# Bot commands:
# /ocr-providers   - List all OCR providers
# /ocr-test        - Test OCR chain
# /ocr <provider>  - Use specific provider
```

### Test Harness

```bash
# Test specific provider
python3 scripts/test-ocr-chain.py --provider mistral --chain

# Test all providers
python3 scripts/test-ocr-chain.py --all

# Parallel mode
python3 scripts/test-ocr-chain.py --parallel --all

# JSON output
python3 scripts/test-ocr-chain.py --all --json
```

### Cloud Shell (Tier 2 providers)

```bash
# In Cloud Shell (ssh):
bash scripts/cloud-shell-ocr.sh

# Start Cloudflare Tunnel (for remote access):
bash scripts/cloudflared-tunnel.sh

# Edit ocr_config.toml with tunnel URL
```

## Termux / Mobile CLI

All OCR scripts run on **Termux** (Android) with no modifications:

```bash
# Install Termux from F-Droid (recommended over Play Store)
# Then:

pkg update && pkg upgrade
pkg in python tesseract git

# Clone repo
git clone https://github.com/your-repo/bac.git
cd bac

# Copy env
cp .env.example .env
nano .env  # add MISTRAL_API_KEY

# Run OCR
MISTRAL_API_KEY="your-key" python3 scripts/ocr_providers.py image.jpg

# Install languages for Tesseract
pkg in tesseract-lang  # includes French + Arabic

# For PaddleOCR / EasyOCR:
pip install paddlepaddle paddleocr easyocr
```

## Configuration

### Environment Variables

```bash
# Core
GATEWAY_HOST=127.0.0.1:8080

# Database (Neon PostgreSQL)
NEON_DB_URL=postgresql://user:pass@host/db?sslmode=require

# AI (Gemini)
GEMINI_API_KEY=your_api_key

# Storage (Garage S3)
GARAGE_ENDPOINT=http://localhost:3900
```

### aichat Tools

| Tool | Service | Description |
|------|---------|-------------|
| `gemini_*` | gemini-tools | Generate, extract, embed content |
| `vector_*` | vector-tools | Semantic search operations |
| `vault_*` | vault-tools | Read/write vault notes |
| `cloud_*` | cloud-tools | SSH, upload, download |
| `graph_*` | graph-tools | Generate/extract graphs |

## Development

```bash
# Run all tests
cargo test

# Build release
cargo build --release

# Run specific service
cargo run -p gemini-tools --release

# Watch mode (dev)
cargo watch -x build
```

## License

MIT OR Apache-2.0
