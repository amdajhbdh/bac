#!/bin/bash
# BAC Study - Termux Setup Script
# Run this on Termux (Android)

set -e

echo "=========================================="
echo "  BAC Study - Termux Setup"
echo "=========================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if running in Termux
if [ ! -d "/data/data/com.termux/files" ]; then
	echo -e "${YELLOW}Warning: Not running in Termux. Some features may not work.${NC}"
fi

# Update and install base dependencies
echo -e "${GREEN}[1/6] Installing system dependencies...${NC}"
pkg update -y
pkg upgrade -y
pkg install -y \
	git \
	curl \
	wget \
	python \
	python-pip \
	nodejs \
	npm \
	rust \
	cargo \
	golang \
	ffmpeg \
	tesseract \
	tesseract-data-eng \
	tesseract-data-fra \
	openssh

# Setup Python
echo -e "${GREEN}[2/6] Setting up Python environment...${NC}"

# Fix PEP 668 issue
export PIP_BREAK_SYSTEM_PACKAGES=1

# Create virtual environment
python -m venv .venv
source .venv/bin/activate

# Install core Python packages
pip install --break-system-packages \
	requests \
	typer \
	rich \
	pyyaml \
	fastapi \
	uvicorn

# Install OCR dependencies (CPU-only)
pip install --break-system-packages \
	numpy \
	opencv-python-headless \
	pillow

# Install RAG dependencies
pip install --break-system-packages \
	torch --index-url https://download.pytorch.org/whl/cpu || true

pip install --break-system-packages \
	langchain \
	langchain-community \
	langchain-core \
	langchain-text-splitters \
	langchain-chroma \
	transformers \
	sentence-transformers \
	chromadb \
	rank-bm25 \
	pypdf || true

# Install Node.js packages
echo -e "${GREEN}[3/6] Setting up Node.js...${NC}"
npm install -g wrangler 2>/dev/null || true

# Setup Go
echo -e "${GREEN}[4/6] Setting up Go...${NC}"
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin:$HOME/.local/bin

# Create necessary directories
echo -e "${GREEN}[5/6] Creating directories...${NC}"
mkdir -p cache chroma_db domain_data

# Copy environment template
if [ ! -f .env ]; then
	cp .env.example .env 2>/dev/null || echo "# BAC Study Environment Variables" >.env
	echo -e "${YELLOW}Please edit .env with your API keys${NC}"
fi

# Test installations
echo -e "${GREEN}[6/6] Verifying installations...${NC}"

# Test Python
python --version && echo -e "${GREEN}✓ Python OK${NC}"

# Test Go
go version && echo -e "${GREEN}✓ Go OK${NC}" || echo -e "${YELLOW}⚠ Go not installed${NC}"

# Test Node
node --version && echo -e "${GREEN}✓ Node.js OK${NC}"

# Test Rust
rustc --version && echo -e "${GREEN}✓ Rust OK${NC}" || echo -e "${YELLOW}⚠ Rust not installed${NC}"

# Create convenience aliases
echo ""
echo -e "${GREEN}=========================================="
echo "  Setup Complete!"
echo "==========================================${NC}"

echo ""
echo "Useful commands:"
echo "  source .venv/bin/activate    # Activate Python venv"
echo "  python scripts/rag_pipeline.py  # Run RAG pipeline"
echo "  cd src/api && go run main.go  # Run API server"
echo ""

# Add to ~/.bashrc for persistence
BASHRC="$HOME/.bashrc"
ALIASES="
# BAC Study aliases
alias bac='python scripts/cli/bac-typer.py'
alias bac-activate='source .venv/bin/activate'
alias bac-api='cd src/api && go run main.go'
alias bac-ocr='python scripts/ocr/local-ocr-server.py'
alias bac-rag='python scripts/rag_pipeline.py'
"

if ! grep -q "BAC Study aliases" $BASHRC 2>/dev/null; then
	echo "$ALIASES" >>$BASHRC
	echo -e "${GREEN}Aliases added to ~/.bashrc${NC}"
fi

echo -e "${GREEN}Run 'source ~/.bashrc' to use aliases${NC}"
