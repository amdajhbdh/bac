#!/bin/bash
# BAC Free Tier - Cloud Shell Setup
# Optimized for Google Cloud Shell free tier (5GB, ephemeral)

set -e

BUCKET_NAME="bac-vault-$(whoami)"
REGION="us-central1"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${GREEN}[BAC]${NC} $*"; }
warn() { echo -e "${YELLOW}[BAC]${NC} $*"; }

# Check gcloud
check_gcloud() {
	if ! command -v gcloud &>/dev/null; then
		log "Installing gcloud..."
		apt-get update && apt-get install -y google-cloud-sdk
	fi
}

# Setup Cloud Storage bucket
setup_bucket() {
	log "Setting up Cloud Storage bucket..."

	# Check if bucket exists
	if gsutil ls gs://${BUCKET_NAME} &>/dev/null; then
		warn "Bucket gs://${BUCKET_NAME} already exists"
	else
		gsutil mb -l ${REGION} gs://${BUCKET_NAME}
		log "Created bucket: gs://${BUCKET_NAME}"
	fi

	# Make bucket public for reading (optional)
	gsutil iam ch allUsers:objectViewer gs://${BUCKET_NAME}
}

# Upload vault to bucket
upload_vault() {
	local vault_dir="${1:-$HOME/bac/resources/notes}"

	log "Uploading vault to Cloud Storage..."
	gsutil -m rsync -r "$vault_dir" gs://${BUCKET_NAME}/vault/
	log "Vault uploaded"
}

# Download vault from bucket
download_vault() {
	local vault_dir="${1:-$HOME/bac/resources/notes}"

	log "Downloading vault from Cloud Storage..."
	mkdir -p "$vault_dir"
	gsutil -m rsync -r gs://${BUCKET_NAME}/vault/ "$vault_dir"
	log "Vault downloaded"
}

# Upload binaries (one-time)
upload_binaries() {
	local bin_dir="${1:-$HOME/bac/bin}"

	if [ ! -d "$bin_dir" ]; then
		warn "Binaries not found at $bin_dir"
		return 1
	fi

	log "Uploading binaries to Cloud Storage..."
	gsutil -m cp -r "$bin_dir"/* gs://${BUCKET_NAME}/binaries/
	log "Binaries uploaded"
}

# Download binaries
download_binaries() {
	local bin_dir="${1:-$HOME/bac/bin}"

	log "Downloading binaries..."
	mkdir -p "$bin_dir"
	gsutil -m cp gs://${BUCKET_NAME}/binaries/* "$bin_dir/" 2>/dev/null || warn "No binaries found"
	log "Binaries ready"
}

# Quick OCR
quick_ocr() {
	local pdf="$1"

	if [ -z "$pdf" ]; then
		echo "Usage: $0 quick-ocr <pdf-file>"
		return 1
	fi

	log "Running OCR on $pdf..."
	cd ~/bac
	USE_CLOUD_MODEL=true CLOUD_MODEL=minimax-m2.5:cloud ./bin/ai-ocr extract "$pdf"
	log "OCR complete"
}

# Sync workflow
sync() {
	log "Syncing with Cloud Storage..."

	# Download latest vault
	download_vault

	# Download latest binaries
	download_binaries

	# Upload changes
	upload_vault

	log "Sync complete"
}

# Status
status() {
	echo "=== BAC Free Tier Status ==="
	echo "Bucket: gs://${BUCKET_NAME}"
	echo ""

	echo "Bucket contents:"
	gsutil ls -la gs://${BUCKET_NAME}/ 2>/dev/null || echo "(empty)"

	echo ""
	echo "Disk usage:"
	df -h ~ | tail -1
}

# Help
help() {
	echo "BAC Free Tier Setup"
	echo ""
	echo "Usage: $0 <command>"
	echo ""
	echo "Commands:"
	echo "  setup          - Initialize bucket and upload vault"
	echo "  sync           - Sync vault with cloud storage"
	echo "  upload-vault   - Upload vault to cloud"
	echo "  download-vault - Download vault from cloud"
	echo "  upload-binaries - Upload binaries (one-time)"
	echo "  quick-ocr <file> - Run OCR on PDF"
	echo "  status         - Show status"
}

# Main
case "${1:-help}" in
setup) check_gcloud && setup_bucket && upload_binaries ;;
sync) sync ;;
upload-vault) upload_vault "$2" ;;
download-vault) download_vault "$2" ;;
upload-binaries) upload_binaries "$2" ;;
download-binaries) download_binaries ;;
quick-ocr) quick_ocr "$2" ;;
status) status ;;
*) help ;;
esac
