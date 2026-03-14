# Render Services Recovery & Full Implementation

## Summary

Recover and fully configure Render services (bac-api, bac-agent) that auto-shutdown due to free tier inactivity, while integrating with existing Cloudflare infrastructure. This creates a comprehensive 455-task implementation plan covering all aspects from recovery to advanced features.

## Problem Statement

1. **Render Services Down**: Both bac-api and bac-agent are unresponsive (likely auto-shutdown after 15 min of inactivity on free tier)
2. **Cloudflare as Backup**: Cloudflare Worker API is running but incomplete
3. **No Production Readiness**: Missing env vars, health checks, monitoring
4. **Gap Between Plan and Execution**: 162 tasks marked complete in OpenSpec but services aren't working

## Goals

- [ ] Get Render services back online and stable
- [ ] Implement dual-primary architecture (Cloudflare + Render)
- [ ] Configure all environment variables
- [ ] Add health checks and keepalive
- [ ] Seed database with questions
- [ ] Deploy working frontend
- [ ] Implement monitoring and alerting

## Scope

### In Scope
- Render service recovery and configuration
- Environment variable setup
- Database seeding
- Cloudflare integration
- Frontend deployment
- Testing and validation

### Out of Scope
- Major new features not listed
- Mobile app development
- Physical hardware purchases

## Timeline

- **Phase 1**: Diagnosis & Assessment (Week 1)
- **Phase 2**: Render Reconfiguration (Week 2)
- **Phase 3-5**: Environment, Database, Reliability (Week 3)
- **Phase 6-7**: API & Cloudflare Integration (Week 4)
- **Phase 8-9**: Frontend & Testing (Week 5-6)
- **Phase 10-16**: Advanced features (Ongoing)

## Decisions Needed

1. **Architecture**: Dual-primary (Cloudflare + Render) or Cloudflare-only?
2. **Environment vars**: Need actual values for NEON_DB_URL, TURSO_DB_URL, REDIS_URL, JWT_SECRET
3. **Keepalive**: External pinger or paid Render tier?

## Risks

- Render free tier auto-shutdown may make this unreliable
- Need to maintain two platforms
- Environment variables may be missing or expired
