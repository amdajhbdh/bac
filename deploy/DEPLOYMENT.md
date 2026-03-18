# BAC Unified Deployment Guide

## Deploy to Render.com

1. **Create Render Account**
   - Go to [render.com](https://render.com)
   - Connect GitHub account
   - Enable payment method (required for non-free services)

2. **Configure Render.yaml**
   ```bash
   cd /home/med/Documents/bac
   git add render.yml
   git commit -m "Add Render deployment config"
   ```

3. **Set Environment Variables**
   ```bash
   # Copy secret to Render dashboard
   cp deploy/render/.env.example .env
   # Fill in your values and push to repository
   ```

4. **Deploy via Render**
   - Push to GitHub
   - Render auto-deploys on `master` branch
   - Error "branch master could not be found" → create `master` branch or switch `branch` in render.yml to your current branch (e.g. `main`)
   - Error "need_payment_info" → add billing to your Render account

## Deploy to Cloudflare Pages

1. **Cloudflare Setup**
   ```bash
   # Install Wrangler
   npm install -g wrangler
   ```

2. **Configure Deployment**
   ```bash
   cd deploy/cloudflare
   wrangler pages deploy --project-name bac-unified
   ```

## Manual Commands

```bash
# Test deployment locally
docker-compose up

# Deploy specific service
render services create bac-api-gateway ./render.yml

# View logs
render logs bac-api-gateway

# Redeploy
render services update bac-api-gateway ./render.yml
```

## Service URLs
- API Gateway: https://bac-api-gateway.onrender.com
- Messaging Service: https://bac-messaging-service.onrender.com
- AI Agent: https://bac-ai-agent.onrender.com
- Cloudflare: https://bac-unified.pages.dev