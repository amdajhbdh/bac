# BAC Knowledge System - Limitations & Known Issues

> **Honest Assessment**: This document documents the current limitations, trade-offs, and technical debt in the BAC Knowledge System. Understanding these constraints helps set realistic expectations and guides future improvements.

---

## Current Limitations

### Session & State Management

| Limitation | Impact | Workaround |
|------------|--------|------------|
| **No persistent sessions** | Context lost on restart | Manual context restoration |
| **Session compression** | May lose subtle context | Keep sessions focused |
| **No session history** | Can't revisit past conversations | External note-taking |

### Database & Search

| Limitation | Impact | Workaround |
|------------|--------|------------|
| **HNSW rebuild is slow** | Index changes take time | Rebuild during off-peak |
| **Chunk size trade-offs** | 2048 tokens may split concepts | Manual re-chunking |
| **No hybrid search** | Keyword search not available | Semantic search only |
| **Neon connection latency** | Remote DB adds delay | Accept ~50-200ms latency |

### Cloud Integration

| Limitation | Impact | Workaround |
|------------|--------|------------|
| **Cloud Shell 5GB limit** | Limited heavy processing | Clean old files regularly |
| **No real-time sync** | Manual git pull/push | Use sync scripts |
| **SSH connection drops** | Commands may fail mid-run | Add retry logic |
| **GCS quota limits** | Storage costs apply | Monitor usage |

### OCR Pipeline

| Limitation | Impact | Workaround |
|------------|--------|------------|
| **Tesseract quality** | Handwriting/scans struggle | Use clean printed sources |
| **No layout detection** | Multi-column issues | Pre-process images |
| **Language detection** | May misdetect mixed content | Specify language explicitly |

### External Dependencies

| Limitation | Impact | Workaround |
|------------|--------|------------|
| **No offline mode** | Requires API keys | Work offline, sync later |
| **Gemini rate limits** | Burst requests blocked | Implement backoff |
| **No retry logic** | Failed requests need manual retry | Build custom retry wrapper |

---

## Known Issues

### Session Management

| Issue | Severity | Description |
|-------|----------|-------------|
| Context loss on restart | Medium | Session state not persisted |
| Long context degrades | Low | Very long conversations may lose early context |
| No session export/import | Low | Can't save session state |

### File Processing

| Issue | Severity | Description |
|-------|----------|-------------|
| Large file processing limited | Medium | Files >50MB may timeout |
| No PDF text extraction fallback | Medium | Corrupt PDFs fail entirely |
| Image rotation not handled | Low | Rotated images need manual correction |

### Git Integration

| Issue | Severity | Description |
|-------|----------|-------------|
| Manual conflict resolution | Medium | Concurrent edits need manual merge |
| No branch management | Low | Single branch workflow |
| Large files not tracked | Medium | Git LFS not configured |

### Tool Crates

| Issue | Severity | Description |
|-------|----------|-------------|
| Incomplete error messages | Low | Errors often generic |
| No health check aggregation | Low | Can't easily see overall system health |
| Timeout handling inconsistent | Low | Some tools hang, others timeout quickly |

---

## Technical Debt

### Infrastructure

| Debt Item | Priority | Description |
|-----------|----------|-------------|
| No authentication on services | High | All endpoints open |
| No retry logic | High | Failed requests need manual retry |
| Minimal logging | Medium | Hard to debug issues |
| No metrics/monitoring | Medium | No visibility into performance |

### Code Quality

| Debt Item | Priority | Description |
|-----------|----------|-------------|
| Tool crates are stubs | High | Basic implementation, missing features |
| No unit test coverage | High | Quality not validated |
| Error handling incomplete | Medium | Many error cases unhandled |
| No integration tests | Medium | Service interactions not tested |

### Data Management

| Debt Item | Priority | Description |
|-----------|----------|-------------|
| No version history for notes | Medium | Can't see changes over time |
| No backup automation | Medium | Manual backup process |
| No data migration scripts | Low | Schema changes need manual updates |

### Documentation

| Debt Item | Priority | Description |
|-----------|----------|-------------|
| API docs incomplete | Low | Some endpoints undocumented |
| No architecture diagrams | Low | System design not visualized |
| Missing changelog | Low | Can't track system changes |

---

## Trade-offs

### RAG Chunk Size: 2048 tokens

**Decision**: Fixed 2048 token chunks for embeddings.

| Pros | Cons |
|------|------|
| Consistent embedding quality | May split related concepts |
| Predictable memory usage | Some context fragmentation |
| Simple implementation | Not optimal for all content types |

**Recommendation**: For critical concepts, manually ensure related content stays together.

### Semantic Search Only

**Decision**: Pure vector similarity search, no BM25/keyword search.

| Pros | Cons |
|------|------|
| Meaning-based matching | Can't find exact phrases |
| Handles synonyms | May miss specific terminology |
| Concept-level retrieval | Less precise than keyword search |

**Recommendation**: Use descriptive queries, not exact phrases.

### Single User Design

**Decision**: No multi-user support, single vault per instance.

| Pros | Cons |
|------|------|
| Simple architecture | Can't share with family/class |
| No auth overhead | Security relies on file permissions |
| No sync conflicts | One source of truth |

**Recommendation**: Use separate vaults for separate users.

### Obsidian as Storage

**Decision**: Markdown files with YAML frontmatter.

| Pros | Cons |
|------|------|
| Human-readable | No database transactions |
| Easy to edit manually | Limited query capability |
| Portable | No built-in versioning |
| Standard format | Large vaults can be slow |

**Recommendation**: Keep vaults organized with MOCs, archive old notes.

---

## Environment Limitations

### Termux (Mobile)

| Limitation | Impact |
|------------|--------|
| Limited RAM | May struggle with large notes |
| No systemd | Manual service management |
| Battery constraints | Long operations drain battery |
| Network dependency | Cloud operations need connectivity |

### Desktop

| Limitation | Impact |
|------------|--------|
| Cross-platform differences | Some paths/commands differ |
| Resource variability | High-end vs low-end machines |
| Display differences | UI may vary across screens |

---

## Third-Party Dependencies

### aichat

- Version compatibility may break
- Plugin ecosystem limited
- Configuration changes on updates

### Neon (PostgreSQL)

- Connection limits on free tier
- Cold start latency
- Cost at scale

### Gemini API

- Rate limits (15-60 requests/min)
- Token limits per request
- Quota costs
- Regional availability

### Garage (S3-compatible)

- Less mature than AWS S3
- Limited documentation
- Performance varies

---

## Risk Assessment

### High Risk

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| API key exposure | Low | Critical | Never commit .env |
| Data loss | Medium | High | Regular git commits |
| Service downtime | Medium | Medium | Monitor health endpoints |

### Medium Risk

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Rate limit exhaustion | Medium | Low | Implement backoff |
| Large vault performance | Low | Medium | Archive old notes |
| Corrupt vault | Low | High | Git versioning |

### Low Risk

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Tool crashes | Low | Low | Auto-restart daemon |
| Memory exhaustion | Low | Low | Monitor resource usage |

---

## Future Improvements Priority

Based on impact vs effort:

| Priority | Improvement | Impact | Effort |
|----------|-------------|--------|--------|
| 1 | Add authentication | High | Medium |
| 2 | Implement retry logic | High | Low |
| 3 | Add health monitoring | Medium | Low |
| 4 | Session persistence | High | Medium |
| 5 | Test coverage | Medium | High |

---

## Getting Help

For issues not covered here:

1. Check [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)
2. Review service logs with `RUST_LOG=debug`
3. Check health endpoints
4. Open an issue with reproduction steps

---

*This document will be updated as limitations are addressed.*
