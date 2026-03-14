#!/bin/bash
# Vault Sync Script - Sync local vault to GitHub for cloud access
# Usage: ./sync-vault.sh

set -e

VAULT_PATH="${VAULT_PATH:-/home/med/Documents/bac/resources/notes/Study-Vault}"
GITHUB_REPO="${GITHUB_REPO:-mednou/study-vault}"
BRANCH="${BRANCH:-main}"

echo "=== Vault Sync to GitHub ==="
echo "Vault: $VAULT_PATH"
echo "Repo: $GITHUB_REPO"
echo ""

cd "$VAULT_PATH"

# Check if git is initialized
if [ ! -d ".git" ]; then
	echo "Initializing git repo..."
	git init
	git remote add origin "https://github.com/$GITHUB_REPO.git"
fi

# Add all changes
git add -A

# Check if there are changes
if git diff --staged --quiet; then
	echo "No changes to sync"
	exit 0
fi

# Commit changes
TIMESTAMP=$(date "+%Y-%m-%d %H:%M")
git commit -m "Vault sync: $TIMESTAMP"

# Push to GitHub
echo "Pushing to GitHub..."
git push -u origin "$BRANCH" 2>/dev/null || echo "Push failed - token may not have write access"

echo "=== Sync Complete ==="
