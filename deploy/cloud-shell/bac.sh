#!/bin/bash
# BAC Unified - Cloud Shell Management Script
# Run from Google Cloud Shell

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="${REPO_DIR:-$HOME/bac}"
VAULT_DIR="$REPO_DIR/resources/notes"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${GREEN}[BAC]${NC} $*"; }
warn() { echo -e "${YELLOW}[BAC]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; }

# Ensure repo exists
init() {
	if [ ! -d "$REPO_DIR" ]; then
		log "Cloning repository..."
		git clone https://github.com/amdajhbdh/bac.git "$REPO_DIR"
	fi
	cd "$REPO_DIR"
	log "Repository ready at $REPO_DIR"
}

# Sync vault
vault-sync() {
	cd "$VAULT_DIR"
	log "Syncing vault..."
	git pull origin main
	git add -A
	git commit -m "cloud-shell: $(date '+%Y-%m-%d %H:%M')" || warn "No changes to commit"
	git push origin main
	log "Vault synced"
}

# Run OCR
ocr() {
	local pdf="$1"
	if [ -z "$pdf" ]; then
		error "Usage: bac ocr <pdf-path>"
		return 1
	fi
	log "Running OCR on $pdf..."
	cd "$REPO_DIR"
	OLLAMA_URL=http://localhost:11434 USE_CLOUD_MODEL=true CLOUD_MODEL=minimax-m2.5:cloud ./bin/ai-ocr extract "$pdf"
	log "OCR complete"
}

# Start services
start() {
	cd "$REPO_DIR/deploy"
	log "Starting services..."
	docker-compose up -d
	log "Services started"
}

# Stop services
stop() {
	cd "$REPO_DIR/deploy"
	log "Stopping services..."
	docker-compose down
	log "Services stopped"
}

# Status
status() {
	echo "=== BAC Status ==="
	echo "Repo: $REPO_DIR"
	echo "Vault: $VAULT_DIR"

	echo ""
	echo "=== Vault Stats ==="
	cd "$VAULT_DIR"
	echo "Notes: $(find . -name '*.md' | wc -l)"
	echo "PDFs: $(find 03-Resources -name '*.pdf' 2>/dev/null | wc -l)"

	echo ""
	echo "=== Git Status ==="
	git status --short

	echo ""
	echo "=== Docker ==="
	cd "$REPO_DIR/deploy"
	docker-compose ps 2>/dev/null || echo "Docker not running"
}

# Help
help() {
	echo "BAC Unified - Cloud Shell Management"
	echo ""
	echo "Usage: $0 <command>"
	echo ""
	echo "Commands:"
	echo "  init        - Initialize repository"
	echo "  vault-sync  - Sync vault with GitHub"
	echo "  ocr <file> - Run OCR on PDF"
	echo "  start      - Start all services"
	echo "  stop       - Stop all services"
	echo "  status     - Show status"
	echo "  help       - Show this help"
}

# Main
case "${1:-help}" in
init) init ;;
vault-sync) vault-sync ;;
ocr) ocr "$2" ;;
start) start ;;
stop) stop ;;
status) status ;;
*) help ;;
esac
