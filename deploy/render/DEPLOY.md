# Deploy to Render

## Option 1: Blueprint (Recommended)

### 1. Commit & Push render.yaml
```bash
git add render.yaml
git commit -m "Add Render deployment"
git push origin first_bm
```

### 2. Open Deploy Link
```
https://dashboard.render.com/blueprint/new?repo=https://github.com/amdajhbdh/bac
```

### 3. Fill Secrets
Set these in the dashboard:
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Generate with `openssl rand -hex 32`
- `OLLAMA_URL` - https://bac-api.amdajhbdh.workers.dev
- `TURSO_DB_URL` - Turso connection (optional)
- `REDIS_URL` - Redis URL (optional)

### 4. Click "Apply"

---

## Option 2: Direct Deploy

Deploy services directly via CLI:
```bash
# Deploy API
render deploy --name bac-api --source ./src/api

# Deploy Agent
render deploy --name bac-agent --source ./src/agent
```

---

## Current Services

| Service | Branch | Status |
|---------|--------|--------|
| bac-api | first_bm | Ready to deploy |
| bac-agent | first_bm | Ready to deploy |

## Environment Variables

Required:
- `PORT` = 8080 (API), 8081 (Agent)
- `DATABASE_URL` = PostgreSQL connection
- `JWT_SECRET` = Generate secure key
- `OLLAMA_URL` = https://bac-api.amdajhbdh.workers.dev
