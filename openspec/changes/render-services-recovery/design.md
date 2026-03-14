# Render Services Recovery - Design

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CURRENT STATE                                 │
├─────────────────────────────────────────────────────────────────────┤
│  Cloudflare Worker     ✅  https://bac-api.amdajhbdh.workers.dev   │
│  ├── /solve (AI)       ✅  Working                                │
│  ├── /questions        ✅  6 questions                            │
│  └── ...               ✅  All endpoints implemented              │
├─────────────────────────────────────────────────────────────────────┤
│  Render bac-api       ❌  Not responding (free tier shutdown)      │
│  Render bac-agent    ❌  Not responding (free tier shutdown)      │
└─────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────┐
│                        TARGET STATE                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐    │
│  │   User       │──────▶│ Cloudflare   │──────▶│   Neon DB    │    │
│  │   (Web)      │       │   Worker     │       │ +TimescaleDB │    │
│  └──────────────┘       └──────┬───────┘      └──────────────┘    │
│                                │                                     │
│                      ┌─────────┴─────────┐                        │
│                      │     Failover      │                        │
│                      │    (if primary    │                        │
│                      │     fails)        │                        │
│                      └─────────┬─────────┘                        │
│                                │                                     │
│                                ▼                                     │
│                       ┌──────────────┐                             │
│                       │   Render     │                             │
│                       │   (Backup)   │                             │
│                       └──────────────┘                             │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Components

### 1. Cloudflare Worker (Primary)
- **URL**: https://bac-api.amdajhbdh.workers.dev
- **Database**: D1 (SQLite)
- **AI**: Workers AI (Llama 3.2)
- **Storage**: R2 (when enabled)
- **Search**: Vectorize

### 2. Render Services (Backup)
- **bac-api**: Go API server
- **bac-agent**: Go CLI server
- **Database**: Neon PostgreSQL
- **Cache**: Turso + Redis

### 3. Frontend
- **Cloudflare Pages**: Primary
- **src/cloudflare-pages/index.html**: Ready to deploy

## Environment Variables

### Required for Render
```bash
NEON_DB_URL=postgresql://...  # From Neon console
TURSO_DB_URL=libsql://...     # From Turso dashboard  
REDIS_URL=redis://...         # From Render or external
JWT_SECRET=...                # Generated secure string
OLLAMA_URL=...                # Cloudflare Worker URL
```

## Health Check Strategy

### Render Health Checks
- **Path**: /health
- **Interval**: 1 minute
- **Timeout**: 10 seconds
- **Start period**: 30 seconds

### Keepalive Solutions
1. **External pinger**: curl every 10 minutes
2. **Paid tier**: Always-on (~$7/service/month)
3. **Cloudflare primary**: Disable Render auto-sleep

## Data Sync Strategy

### Questions
- **Primary**: Cloudflare D1
- **Backup**: Neon PostgreSQL
- **Sync**: Manual or scheduled export

### Users
- **Primary**: Neon PostgreSQL
- **Secondary**: Cloudflare D1
- **Sync**: Bidirectional on write

## API Endpoints

### Cloudflare Worker (Current)
| Endpoint | Method | Status |
|---------|--------|--------|
| /health | GET | ✅ |
| /solve | POST | ✅ |
| /questions | GET/POST | ✅ |
| /search | POST | ✅ |
| /auth/register | POST | ✅ |
| /auth/login | POST | ✅ |
| /user/profile | POST | ✅ |
| /user/points | POST | ✅ |
| /predictions | GET | ✅ |
| /predictions/generate | POST | ✅ |
| /leaderboard | GET | ✅ |
| /practice/answer | POST | ✅ |
| /practice/questions | GET | ✅ |

### To Add (Cloudflare)
| Endpoint | Method | Status |
|---------|--------|--------|
| /upload | POST | ❌ |
| /ocr | POST | ❌ |
| /animate | POST | ❌ |
| /analyze | POST | ❌ |

### To Add (Render)
| Endpoint | Method | Status |
|---------|--------|--------|
| /api/v1/* | * | ❌ |

## Testing Strategy

### Smoke Tests
- [ ] GET /health returns 200
- [ ] POST /solve returns solution
- [ ] GET /questions returns data

### Integration Tests
- [ ] Full solve flow
- [ ] User registration → login → solve → points
- [ ] Leaderboard updates

### Load Tests
- [ ] 10 concurrent users
- [ ] 50 concurrent users
- [ ] 100 concurrent users

## Monitoring

### Metrics to Track
- Request count
- Response latency
- Error rate
- AI call latency
- Database query time
- Cache hit rate

### Alert Conditions
- Error rate > 5%
- Latency > 5 seconds
- Service down

## Rollback Plan

If deployment fails:
1. Revert to previous commit
2. Use Cloudflare Worker as sole provider
3. Disable Render services temporarily

## Dependencies

- [ ] GitHub repository accessible
- [ ] Neon DB active
- [ ] Turso DB active
- [ ] Render account accessible
- [ ] Cloudflare account accessible
