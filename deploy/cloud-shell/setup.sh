#!/bin/bash
# BAC Cloud Shell Setup
# Run this in Cloud Shell or local gcloud environment

set -e

echo "=== BAC Cloud Shell Setup ==="

# Clone repo if not exists
if [ ! -d "$HOME/bac" ]; then
	echo "Cloning repository..."
	git clone https://github.com/amdajhbdh/bac.git "$HOME/bac"
fi

cd "$HOME/bac"

# Make scripts executable
chmod +x deploy/cloud-shell/*.sh

# Check environment
echo ""
echo "=== Environment ==="
echo "Project: $(gcloud config get-value project 2>/dev/null)"
echo "Zone: $(gcloud config get-value compute/zone 2>/dev/null)"

# Setup Cloud Storage bucket
BUCKET_NAME="bac-vault-$(whoami)"
echo ""
echo "=== Setting up Cloud Storage ==="

if gsutil ls gs://${BUCKET_NAME} &>/dev/null; then
	echo "Bucket gs://${BUCKET_NAME} exists"
else
	gsutil mb -l us-central1 gs://${BUCKET_NAME}
	echo "Created bucket: gs://${BUCKET_NAME}"
fi

# Sync vault to bucket
echo ""
echo "=== Syncing vault to Cloud Storage ==="
gsutil -m rsync -r resources/notes/ gs://${BUCKET_NAME}/vault/

# Upload binaries
echo ""
echo "=== Uploading binaries ==="
mkdir -p bin
gsutil -m cp bin/* gs://${BUCKET_NAME}/binaries/ 2>/dev/null || echo "No binaries to upload"

echo ""
echo "=== Setup Complete ==="
echo ""
echo "Access from phone:"
echo "1. shell.cloud.google.com"
echo "2. Run: cd ~/bac && ./deploy/cloud-shell/bac-free-tier.sh sync"
echo ""
echo "Bucket: gs://${BUCKET_NAME}"
