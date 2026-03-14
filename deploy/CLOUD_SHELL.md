# BAC Unified - Cloud Infrastructure

## Access Options

### 1. Google Cloud Shell (Phone Accessible)
- **URL**: https://shell.cloud.google.com
- **Access**: Any browser (phone/laptop)
- **Features**: Pre-installed gcloud, docker, git, python, node
- **Use**: Development, OCR processing, vault sync

### 2. Render (Deployment)
- **Services**: API, Gateway hosting
- **URL**: https://dashboard.render.com

### 3. Cloudflare Workers
- **Services**: API endpoints, AI agents
- **URL**: https://dash.cloudflare.com

### 4. Fly.io
- **Services**: Edge computing
- **URL**: https://dashboard.fly.io

## Cloud Shell Quick Commands

```bash
# SSH into Cloud Shell from anywhere (phone)
gcloud cloud-shell ssh

# Run OCR in Cloud Shell
cd ~/bac
./bin/ai-ocr extract path/to/pdf

# Sync vault
cd ~/bac/resources/notes
git pull origin main
git add -A && git commit -m "update" && git push

# Deploy services
gcloud run deploy bac-api --source ./src/api
```

## Architecture

```
Phone Browser
     │
     ├──────► Cloud Shell (gcloud shell.cloud.google.com)
     │              └─ OCR, Vault Sync, Development
     │
     ├──────► Telegram/WhatsApp (via OpenCrust)
     │
     └──────► Web Apps (Render/Cloudflare/Fly.io)
```

## Setup Cloud Shell

1. Go to https://shell.cloud.google.com
2. Clone repo:
```bash
git clone https://github.com/amdajhbdh/bac.git ~/bac
cd ~/bac
```
3. Set up environment:
```bash
cp deploy/env.template .env
# Edit .env with your tokens
```
4. Run services:
```bash
docker-compose up -d
```
