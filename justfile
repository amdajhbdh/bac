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

# Deploy script
deploy:
    ./scripts/deploy.sh

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
