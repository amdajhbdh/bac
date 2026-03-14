# Google Cloud Shell (Free Tier) Setup

## Free Tier Limitations
- **5GB** persistent disk (home directory)
- **Sessions** expire after 90 min inactivity
- **CPU**: Shared, not always-on
- **No** custom domains

## Optimized Setup

### 1. Store Vault on Cloud Storage (instead of disk)
```bash
# Create bucket (free tier: 5GB free)
gsutil mb -l us-central1 gs://bac-vault-$(whoami)

# Sync vault to bucket
gsutil -m rsync -r gs://bac-vault-$(whoami)/ ~/bac/resources/notes/
```

### 2. Use Cloud Scheduler for OCR Jobs
```bash
# Create Cloud Storage trigger
gcloud functions deploy ocr-trigger \
  --runtime python311 \
  --trigger-bucket gs://bac-vault-$(whoami) \
  --entry-point process_pdf
```

### 3. Keep Binaries in Cloud Storage
```bash
# Upload once, mount anywhere
gsutil cp bin/* gs://bac-unified-binaries/
```

## Architecture (Free Tier)
```
Phone ──► Cloud Shell (ephemeral)
            │
            ├─► Cloud Storage (vault + binaries)
            │
            └─► Cloud Functions (OCR trigger)
```

## Free Tier Commands
```bash
# In Cloud Shell
./bac.sh mount-gcs     # Mount bucket as vault
./bac.sh upload-binaries  # Upload once
./bac.sh schedule-ocr   # Schedule daily OCR
```

## Alternative: Use Cloud Run (Always-on)
```bash
# Deploy API to Cloud Run (free tier: 180k vCPU-seconds/day)
gcloud run deploy bac-api \
  --source ./src/api \
  --region us-central1 \
  --allow-unauthenticated
```

## Quick Setup Script
```bash
# Run in Cloud Shell
cd ~
git clone https://github.com/amdajhbdh/bac.git
cd bac/deploy/cloud-shell
chmod +x bac-free-tier.sh
./bac-free-tier.sh setup
```
