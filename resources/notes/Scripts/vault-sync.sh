#!/bin/bash
# Vault Sync Script for ZeroClaw integration
# Manages git sync for mobile Obsidian access

set -e

VAULT_DIR="${VAULT_DIR:-/home/med/Documents/bac/resources/notes}"
TELEGRAM_BOT_TOKEN="${TELEGRAM_BOT_TOKEN}"
NOTIFICATION_CHAT_ID="${NOTIFICATION_CHAT_ID}"

cd "$VAULT_DIR"

log() {
	echo "[$(date '+%Y-%m-%d %H:%M')] $*"
}

notify() {
	if [ -n "$TELEGRAM_BOT_TOKEN" ] && [ -n "$NOTIFICATION_CHAT_ID" ]; then
		curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage" \
			-d "chat_id=$NOTIFICATION_CHAT_ID" \
			-d "text=$1" >/dev/null 2>&1
	fi
}

# Check for changes
has_changes() {
	git diff --quiet HEAD 2>/dev/null && return 1 || return 0
}

# Pull from remote
sync_pull() {
	log "Pulling vault from remote..."
	git pull origin main
	notify "Vault updated: pulled from remote"
}

# Push to remote
sync_push() {
	if has_changes; then
		log "Pushing vault to remote..."
		git add -A
		git commit -m "vault: $(date '+%Y-%m-%d %H:%M')"
		git push origin main
		notify "Vault updated: pushed changes"
	else
		log "No changes to push"
	fi
}

# Run OCR on new PDFs
run_ocr() {
	log "Checking for new PDFs..."
	cd "$VAULT_DIR"

	for pdf in $(find 03-Resources -name "*.pdf" -newer 05-Extracted 2>/dev/null); do
		log "Processing: $pdf"
		python3 /home/med/Documents/bac/src/ai-ocr/main.go extract "$pdf"
	done
}

# Full sync
sync() {
	sync_pull
	run_ocr
	sync_push
}

# Command handler for ZeroClaw
case "${1:-sync}" in
pull)
	sync_pull
	;;
push)
	sync_push
	;;
ocr)
	run_ocr
	;;
sync)
	sync
	;;
status)
	git status --short
	echo "---"
	echo "Last commit: $(git log -1 --format='%h %s' HEAD)"
	;;
*)
	echo "Usage: $0 {pull|push|ocr|sync|status}"
	exit 1
	;;
esac
