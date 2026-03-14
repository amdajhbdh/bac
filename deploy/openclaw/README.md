# KiloClaw Integration for BAC Study Vault

## Quick Setup

1. Go to [kilo.ai](https://kilo.ai) and sign in
2. Navigate to **Claw** in the sidebar
3. Click **Create Instance**
4. Select your model (default: claude-sonnet)
5. Configure Telegram (optional): Add your bot token
6. Click **Create & Provision**

Your KiloClaw instance will be ready in seconds!

## Configuration

After creating your KiloClaw instance, configure it:

```
KILO_CLAW_URL=https://your-instance.kiloclaw.ai
KILO_CLAW_API_KEY=your-api-key
```

## Telegram Commands

Once connected to Telegram:

| Command | Description |
|---------|-------------|
| `/scan` | Trigger content scanner |
| `/report` | Get daily report |
| `/query <text>` | Search vault |
| `/status` | Agent status |
| `/recommend` | Study recommendation |
| `/quiz` | Generate practice quiz |
| `/links` | Find related notes |
| `/help` | Show all commands |

## Environment Variables

Set these for the KiloClaw to access your vault:

```bash
# KiloClaw (hosted - no setup needed)
KILO_CLAW_URL=https://your-instance.kiloclaw.ai

# Vault API (deployed on Render)
VAULT_API_URL=https://vault-api.onrender.com

# Cloudflare Workers
COORDINATOR_URL=https://vault-coordinator.your-subdomain.workers.dev
RESEARCH_URL=https://vault-research-agent.your-subdomain.workers.dev
```

## Architecture

```
Telegram → KiloClaw (kilo.ai) 
                    ↓
           Vault API (Render + R2)
                    ↓
           Cloudflare Workers
                    ↓
           Study Vault (GitHub sync)
```

## Deploy Other Services

```bash
# Deploy to Render
just deploy-render

# Deploy Cloudflare Workers
just deploy-cf

# Full deployment
just full-deploy
```
