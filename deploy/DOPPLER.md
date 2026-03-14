# Doppler Configuration for BAC Unified

## Quick Setup

```bash
# Install Doppler CLI
curl -Ls https://doppler.com/install | sh

# Login
doppler login

# Setup project
cd /home/med/Documents/bac
doppler setup

# Set secrets
doppler secrets set DATABASE_URL="postgresql://..."
doppler secrets set JWT_SECRET="$(openssl rand -hex 32)"
```

## Secrets Required

| Secret | Description |
|--------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `JWT_SECRET` | Generate with `openssl rand -hex 32` |
| `OLLAMA_URL` | Ollama endpoint |
| `TURSO_DB_URL` | Turso SQLite (optional) |
| `REDIS_URL` | Redis cache (optional) |
| `TELEGRAM_BOT_TOKEN` | From @BotFather |
| `OPENAI_API_KEY` | For OpenCrust |
| `ANTHROPIC_API_KEY` | For Claude |

## Usage in Code

### Go
```go
import "github.com/dopplerhq/doppler-go"

client, _ := doppler.New()
secrets, _ := client.Get("dev")
dbURL := secrets["DATABASE_URL"]
```

### Render
Add to render.yaml:
```yaml
envVars:
  - key: DATABASE_URL
    fromSecret: DATABASE_URL
```

### Cloudflare Workers
```bash
# Set via wrangler
wrangler secret put DATABASE_URL
```

## Environments

- `dev` - Development
- `staging` - Staging
- `prod` - Production
