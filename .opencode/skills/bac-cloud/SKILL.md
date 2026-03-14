---
name: bac-cloud
description: Google Cloud Shell management for BAC project. Use when working with cloud shell, deploying to render/cloudflare, or managing cloud resources. Triggers: gcloud, cloud-shell, render, deploy, cloudflare, fly
---

# BAC Cloud Infrastructure Skill

Manage cloud services and deployments.

## Cloud Shell (Phone Accessible)

**URL**: https://shell.cloud.google.com

### Quick Commands

```bash
# Clone repo
git clone https://github.com/amdajhbdh/bac.git ~/bac
cd ~/bac

# Setup
./deploy/cloud-shell/bac-free-tier.sh setup

# Sync vault
./deploy/cloud-shell/bac-free-tier.sh sync

# Run OCR
./deploy/cloud-shell/bac-free-tier.sh quick-ocr file.pdf
```

## Services

| Service | Purpose | Free Tier |
|---------|---------|-----------|
| Cloud Shell | Development | 5GB, 90min |
| Cloud Storage | Vault storage | 5GB |
| Cloud Run | API hosting | 180k vCPU-s/day |
| Cloudflare | Workers/Workers AI | Generous free |

## Deploy Commands

```bash
# Render
gcloud run deploy bac-api --source ./src/api

# Cloudflare
wrangler deploy

# Fly.io
fly deploy
```

## Environment Variables

Set in `.env`:
- `DATABASE_URL` - PostgreSQL
- `JWT_SECRET` - Generate with `openssl rand -hex 32`
- `TELEGRAM_BOT_TOKEN` - From @BotFather
- `OPENAI_API_KEY` - For OpenCrust

## Tips

- Use Cloud Storage for vault (persistent across sessions)
- Upload binaries once to GCS
- Cloud Run for always-on API
- Cloudflare Workers for low-latency AI
