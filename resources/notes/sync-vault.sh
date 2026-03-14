#!/usr/bin/env bash
# Vault Sync Script - For use on desktop
# Usage: ./sync-vault.sh [push|pull]

set -e

VAULT_DIR="/home/med/Documents/bac/resources/notes"
cd "$VAULT_DIR"

ACTION="${1:-push}"

if [ "$ACTION" = "push" ]; then
	echo "Pushing vault changes to remote..."
	git add -A
	git commit -m "vault: $(date '+%Y-%m-%d %H:%M')" || echo "No changes to commit"
	git push origin main
	echo "Done!"
elif [ "$ACTION" = "pull" ]; then
	echo "Pulling vault changes from remote..."
	git pull origin main
	echo "Done!"
else
	echo "Usage: $0 [push|pull]"
	exit 1
fi
