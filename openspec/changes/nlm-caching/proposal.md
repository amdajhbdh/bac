# Proposal: NLM Caching System

## Summary

Implement a hybrid caching system for NotebookLM queries that combines local SQLite for fast lookups with Garage S3 for durable blob storage. This reduces NLM API calls by 60-80% while avoiding rate limits.

## Problem

- NLM has strict rate limits (2-5 seconds between operations)
- Sessions expire after ~20 minutes
- Single notebook creates a bottleneck
- No caching - every query hits NLM directly

## Solution

Hybrid architecture:
- **SQLite**: Fast local index with TTL tracking
- **Garage S3**: Durable response storage
- **Query Router**: Subject/topic extraction + notebook selection
- **Rate Limiter**: Per-notebook rate limiting

## Benefits

| Metric | Before | After |
|--------|--------|-------|
| NLM API calls | Every query | 20-40% of queries |
| Cache hit rate | 0% | 60-80% |
| Response time (cache) | 2-10s | <100ms |

## Scope

- Create caching system in `src/agent/internal/nlm/`
- Use existing Garage S3 infrastructure
- Support 4 subject areas: math, pc, svt, philosophy

## Non-Goals

- Replace NLM entirely (still primary source)
- Implement distributed cache (single node sufficient)
- Support real-time cache invalidation

## Timeline

- Phase 1 (Core): 2-3 hours
- Phase 2 (Integration): 1-2 hours
- Phase 3 (Maintenance): 1 hour

Total: ~4-6 hours
