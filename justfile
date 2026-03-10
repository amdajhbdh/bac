# Build the main agent binary
build:
    cd src/agent && go build -o bin/Agent ./cmd/main.go

# Build CLI tools
build-tools:
    cd src/agent && go build -o bin/kb-cli ./cmd/kb-cli/
    cd src/agent && go build -o bin/semantic-test ./cmd/semantic-test/

# Quick build with all checks
quick-build:
    @echo "🚀 Quick Build Starting..."
    @echo "Checking prerequisites..."
    @command -v go >/dev/null 2>&1 || (echo "❌ Go not installed" && exit 1)
    @command -v podman >/dev/null 2>&1 || (echo "❌ Podman not installed" && exit 1)
    @echo "✅ Prerequisites OK"
    @echo ""
    @echo "Setting up environment..."
    @test -f .env || cp .env.example .env
    @echo "✅ Environment ready"
    @echo ""
    @echo "Starting services..."
    podman-compose up -d postgres redis
    @sleep 5
    @echo "✅ Services started"
    @echo ""
    @echo "Building packages..."
    cd src/agent && go mod tidy
    cd src/agent && go build ./internal/db/ || echo "⚠️  db failed"
    cd src/agent && go build ./internal/memory/ || echo "⚠️  memory failed"
    cd src/agent && go build ./internal/solver/ || echo "⚠️  solver failed"
    cd src/agent && go build ./internal/online/ || echo "⚠️  online failed"
    @echo ""
    @echo "Building main binary..."
    cd src/agent && go build -o bin/Agent ./cmd/main.go
    @echo ""
    @echo "✅ Build complete!"
    @echo "Binary: src/agent/bin/Agent"

# Build and test everything
build-test: quick-build
    @echo ""
    @echo "Testing binary..."
    ./src/agent/bin/Agent -h
    @echo "✅ Binary works"

# Run all tests
test:
    cd src/agent && go test ./...

# Run single test (usage: just test-one TestName ./internal/solver/)
test-one TEST PKG:
    cd src/agent && go test -v -run {{ TEST }} {{ PKG }}

# Run tests with race detector
test-race:
    cd src/agent && go test -race ./...

# Format code
fmt:
    cd src/agent && go fmt ./...

# Run go vet
vet:
    cd src/agent && go vet ./...

# Install dependencies
install:
    cd src/agent && go mod download
    pip install -r src/agent/internal/ocr/requirements.txt || echo "⚠️  Python deps optional"

# Fix common build issues
fix:
    @echo "🔧 Auto-fixing build issues..."
    cd src/agent && go clean -cache -modcache
    cd src/agent && go mod tidy
    cd src/agent && go mod download
    cd src/agent && go mod verify
    @echo "✅ Auto-fix complete"

# Start services with podman
services-up:
    podman-compose up -d postgres redis
    podman compose ps

# Stop services
services-down:
    podman compose down

# View logs
logs SERVICE="":
    #!/usr/bin/env bash
    if [ -z "{{ SERVICE }}" ]; then
        podman-compose logs -f
    else
        podman-compose logs -f {{ SERVICE }}
    fi

# Start development server
dev:
    cd src/agent && go run ./cmd/main.go

# Start REST API
api:
    cd src/api && go run main.go

# Start agent server
server:
    cd src/agent && go run ./cmd/main.go -server -port 8081

# Health check
health:
    @curl -f http://localhost:8081/health || echo "❌ Service not healthy"

# Load sample data
load-data:
    podman-compose exec postgres psql -U bac -d bac < sql/sample_data.sql

# Database shell
db-shell:
    podman-compose exec postgres psql -U bac -d bac

# Redis shell
redis-shell:
    podman-compose exec redis redis-cli

# Clean build artifacts
clean:
    rm -rf src/agent/bin/

# Clean everything including containers
clean-all: clean services-down
    podman volume prune -f

# Deploy to production
deploy:
    @echo "🚀 Deploying..."
    podman build -t bac-agent:latest .
    @echo "✅ Image built"
    @echo "Push with: podman push bac-agent:latest"

# Show all available commands
list:
    @just --list

# ===========================================
# VAULT SYNC TASKS (jj)
# ===========================================

# Sync vault to GitHub using jj
sync-vault:
    #!/usr/bin/env bash
    set -e
    VAULT_PATH="${VAULT_PATH:-/home/med/Documents/bac/resources/notes/Study-Vault}"
    REPO="${GITHUB_REPO:-mednou/study-vault}"
    
    cd "$VAULT_PATH"
    
    # Initialize jj repo if needed
    if [ ! -d ".jj" ]; then
        jj init
        jj config set --user user.email "vault@sync.local"
        jj config set --user user.name "Vault Sync"
    fi
    
    # Check for changes
    jj diff --summary | grep -q . && CHANGES=1 || CHANGES=0
    
    if [ "$CHANGES" = "1" ]; then
        jj describe -r @ -m "Vault sync $(date '+%Y-%m-%d %H:%M')"
        jj new -m "Vault update $(date '+%Y-%m-%d')"
        jj git push --all origin 2>/dev/null || jj git push --all 2>/dev/null || echo "Push skipped"
    else
        echo "No changes to sync"
    fi

# Show vault status
vault-status:
    #!/usr/bin/env bash
    VAULT_PATH="${VAULT_PATH:-/home/med/Documents/bac/resources/notes/Study-Vault}"
    cd "$VAULT_PATH"
    echo "=== Vault Status (jj) ==="
    jj log --limit 3 -r @ 2>/dev/null || jj log --limit 3 2>/dev/null || echo "No jj history"
    echo ""
    echo "=== File Count ==="
    find . -type f -name "*.md" 2>/dev/null | wc -l | xargs echo "Markdown files:"

# ===========================================
# DEPLOYMENT TASKS
# ===========================================

# Deploy to Render
deploy-render:
    @echo "Deploying to Render..."
    @echo "Use: render blueprint create deploy/render/render.yaml"
    @render services list -o json 2>/dev/null | jq '.[] | select(.type=="web") | {name, url}' || echo "Check Render Dashboard"

# Deploy OpenClaw to KiloClaw (hosted)
deploy-kiloclaw:
    @echo "=== KiloClaw Setup ==="
    @echo "1. Go to https://kilo.ai and sign in"
    @echo "2. Click Claw in the sidebar"
    @echo "3. Click Create Instance"
    @echo "4. Select model and configure Telegram"
    @echo "5. Click Create & Provision"
    @echo ""
    @echo "Your KiloClaw will be ready in seconds!"
    @echo "No Render deployment needed - KiloClaw is hosted!"

# Deploy Cloudflare Workers
deploy-cf:
    @echo "Deploying Cloudflare Workers..."
    wrangler deploy --config .kiro/agents/coordinator/wrangler.toml 2>/dev/null || echo "Configure wrangler first"
    wrangler deploy --config .kiro/deploy/research-agent/wrangler.toml 2>/dev/null || echo "Configure wrangler first"

# Deploy vault API to Cloudflare Workers
deploy-vault-worker:
    #!/usr/bin/env bash
    set -e
    echo "Deploying vault-api to Cloudflare Workers..."
    cd .kiro/deploy/vault-worker
    wrangler deploy

# Full sync and deploy
sync-deploy: sync-vault deploy-render deploy-cf deploy-vault-worker
    @echo "Sync and deploy complete!"

# Health check all services
health-all:
    @echo "Checking services..."
    @curl -sf https://vault-api.onrender.com/health && echo "✓ Vault API" || echo "✗ Vault API"
    @curl -sf https://bac-agent.onrender.com/health && echo "✓ BAC Agent" || echo "✗ BAC Agent"

# Default task
default: quick-build

# ===========================================
# EXTRACTION TASKS
# ===========================================

# Build extraction system
build-extraction:
    @echo "Building extraction system..."
    cd src/agent && go build -o bin/Agent ./cmd/main.go
    @echo "✅ Build complete!"

# Extract from PDF using extraction system
extract FILE:
    source venv_ocr/bin/activate && ./src/agent/bin/Agent -extract {{ FILE }}

# Fast extraction mode
extract-fast FILE:
    ./src/agent/bin/Agent -extract -fast {{ FILE }}

# Extract to Cloudflare
extract-cloudflare:
    ./scripts/extract-cloudflare.sh

# Extract batch of files
extract-batch:
    ./scripts/extract-batch.sh

# ===========================================
# PDF PROCESSING TASKS
# ===========================================

# Convert PDF to Markdown
pdf2md PDF:
    @PDF="{{ PDF }}"; BASE=$$(basename "$$PDF" .pdf); DIR=$$(dirname "$$PDF"); \
    GRAPHS_DIR="$$DIR/graphs/$$BASE"; \
    mkdir -p "$$GRAPHS_DIR"; \
    python3 scripts/extract_graphs.py "$$PDF" "$$GRAPHS_DIR" 2>&1 | grep "✅" || true; \
    python3 scripts/json2md.py "$$PDF.extraction.json" "$$PDF.md"

# Fast PDF to Markdown
pdf2md-fast PDF:
    ./scripts/pdf2md-fast.sh {{ PDF }}

# Simple PDF to Markdown
pdf2md-simple PDF:
    python3 scripts/pdf_ocr.py {{ PDF }}

# ===========================================
# BATCH PROCESSING TASKS
# ===========================================

# Batch extract all PDFs
batch-extract:
    ./scripts/batch-extract.sh

# Batch process to Cloudflare
batch-process-cloud:
    ./scripts/batch-process-cloud.sh

# Batch upload to Cloudflare
batch-upload-cloudflare:
    ./scripts/batch-upload-cloudflare.sh

# ===========================================
# BUILD TASKS
# ===========================================

# Quick build script
quick-build-script:
    ./scripts/quick-build.sh

# Build and test script
build-and-test:
    ./scripts/build-and-test.sh

# Build extractor
build-extractor:
    ./scripts/build-extractor.sh

# Fix build issues
fix-build:
    ./scripts/fix-build.sh

# ===========================================
# DEPLOY TASKS
# ===========================================

# Deploy to production
deploy-prod:
    ./scripts/deploy-production.sh

# ===========================================
# TEST TASKS
# ===========================================

# Test extraction
test-extraction:
    ./scripts/test-extraction.sh

# Test all
test-all:
    ./scripts/test-all.sh

# Test semantic search
test-semantic:
    ./scripts/test-semantic-search.sh

# ===========================================
# INIT TASKS
# ===========================================

# Initialize platform
init-platform:
    ./scripts/init-platform.sh

# Setup extraction
setup-extraction:
    ./scripts/setup-extraction.sh

# ===========================================
# UTILITY TASKS
# ===========================================

# Start MinIO
start-minio:
    ./scripts/start-minio.sh

# Analyze graphs
analyze-graphs:
    python3 scripts/analyze_graphs.py

# Process knowledge base
process-kb:
    python3 scripts/process_kb.py

# Benchmark KB
benchmark-kb:
    python3 scripts/benchmark_kb.py

# Scan PDF
scan-pdf PDF:
    python3 scripts/scan_pdf.py {{ PDF }}

# Open Zen browser
zen-open:
    ./scripts/zen-open.sh

# ===========================================
# MISE TASKS
# ===========================================

# Install mise tools
mise-install:
    @echo "Installing mise tools..."
    @if command -v mise &> /dev/null; then \
        mise install; \
    else \
        echo "Install mise: https://mise.jdx.dev/getting-started.html"; \
    fi

# Install Go via mise
mise-install-go:
    mise install go@1.25

# Install Node via mise
mise-install-node:
    mise install node@22

# Install Python via mise
mise-install-python:
    mise install python@3.11

# Install Rust via mise
mise-install-rust:
    mise install rust@latest

# Install all tools
mise-install-all: mise-install-go mise-install-node mise-install-python mise-install-rust

# Show mise status
mise-status:
    @if command -v mise &> /dev/null; then \
        mise status; \
    else \
        echo "mise not installed"; \
    fi

# ===========================================
# OCR SERVICE TASKS
# ===========================================

# Build OCR service
ocr-build:
    cd src/ocr-service && cargo build

# Run OCR service tests
ocr-test:
    cd src/ocr-service && cargo test

# Start OCR service
ocr-start:
    cd src/ocr-service && cargo run --bin ocr-service

# Generate OCR test data
ocr-gen-test-data:
    python3 openspec/changes/bac-ocr-service-requirements/generate_test_data.py

# Run OCR integration tests
ocr-test-integration:
    @echo "Testing OCR service endpoints..."
    @curl -s http://127.0.0.1:3000/health | jq -e '.status == "healthy"' && echo "✓ Service healthy"
    @curl -s -X POST http://127.0.0.1:3000/ocr -F "file=@openspec/changes/bac-ocr-service-requirements/test_data/basic_english.png" | jq -e '.success == true' && echo "✓ English OCR"
    @curl -s -X POST http://127.0.0.1:3000/ocr -F "file=@openspec/changes/bac-ocr-service-requirements/test_data/basic_french.png" | jq -e '.success == true' && echo "✓ French OCR"
    @curl -s -X POST http://127.0.0.1:3000/pdf -F "file=@openspec/changes/bac-ocr-service-requirements/test_data/multipage_test.pdf" | jq -e '.success == true' && echo "✓ PDF OCR"
    @echo "All OCR tests passed!"

# Test OCR with specific image
ocr-test-image FILE:
    curl -s -X POST http://127.0.0.1:3000/ocr -F "file={{ FILE }}" | jq .

# Test OCR with PDF
ocr-test-pdf FILE:
    curl -s -X POST http://127.0.0.1:3000/pdf -F "file={{ FILE }}" | jq .

# Test OCR with Arabic PDF
ocr-test-arabic FILE:
    curl -s -X POST http://127.0.0.1:3000/pdf -F "file={{ FILE }}" | jq '.data.text' | head -c 200

# Process all PDFs in db/pdfs/
ocr-process-pdfs:
    @echo "Processing all PDFs in db/pdfs/..."
    @mkdir -p db/pdfs/processed
    @for f in db/pdfs/*.pdf; do \
        filename=$$(basename "$$f"); \
        echo "Processing: $$filename"; \
        curl -s -X POST http://127.0.0.1:3000/pdf -F "file=@$$f" > "db/pdfs/processed/$$filename.json" || echo "Failed: $$filename"; \
    done
    @echo "Done! Output in db/pdfs/processed/"

# ===========================================
# ANIMATION TASKS
# ===========================================

# Build gateway with animation
gateway-build:
    cd src/gateway && cargo build

# Run gateway tests
gateway-test:
    cd src/gateway && cargo test

# Start gateway
gateway-start:
    cd src/gateway && cargo run

# Check Manim availability
manim-check:
    @podman images | grep manim && echo "Manim available!" || echo "Manim not found"

# Test Manim with sample
manim-test:
    @echo "Testing Manim..."
    @mkdir -p /tmp/manim_test
    @echo 'from manim import *\nclass Test(Scene):\n    def construct(self):\n        self.play(Write(Text("Hello BAC!")))' > /tmp/manim_test/test.py
    @podman run --rm -v /tmp/manim_test:/workspace docker.io/manimcommunity/manim:latest manim test.py -ql || echo "Manim test failed"
