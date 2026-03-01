# RFC: Resource Search System

## 1. Problem Statement and Constraints

### Problem
BAC Unified needs to aggregate educational content from multiple sources (YouTube, web OER repositories,) and provide unified search across all content. Current system lacks:
- Unified content indexing
- Cross-source search capability
- Automatic content extraction and classification
- Metadata standardization

### Constraints
- Must work offline for core functionality
- Must use open source technologies
- Must minimize API costs (avoid paid APIs unless necessary)
- Must handle Arabic, French, and English content
- Must respect source terms of service

### Requirements
- [ ] Aggregate from 10+ content sources
- [ ] Full-text search with <100ms latency
- [ ] Vector similarity search for semantic queries
- [ ] Automatic content classification
- [ ] Support for video, PDF, images, text
- [ ] Faceted search (subject, chapter, difficulty)

---

## 2. Current State Evidence

### Existing Components

```
src/agent/internal/
├── online/           # Playwright scraping, yt-dlp
├── ocr/             # Tesseract integration
├── memory/          # pgvector for questions
└── nlm/             # NotebookLM integration
```

### Evidence Links
- YouTube: `src/agent/internal/online/playwright.go` - existing yt-dlp
- Web scraping: `src/agent/internal/online/` - Playwright setup
- OCR: `src/agent/internal/ocr/ocr.go` - Tesseract wrapper

### Gaps
- No centralized content index
- No unified metadata schema
- No cross-source search
- No content classification

---

## 3. Goals and Non-Goals

### Goals
1. Create unified resource index
2. Implement hybrid search (full-text + vector)
3. Automate content extraction pipeline
4. Standardize metadata across sources
5. Provide faceted search capabilities

### Non-Goals
1. Real-time streaming of content
2. Social features (sharing, comments)
3. User-generated content hosting
4. Payment/procurement integration
5. Real-time collaboration

---

## 4. Options Considered

### Option A: Elasticsearch-Centric (Selected)

**Architecture:**
```
Content Sources → Extractors → Elasticsearch ← pgvector
                              ↓
                          PostgreSQL (metadata)
```

**Pros:**
- Proven full-text search
- Scalable horizontally
- Rich faceting support
- Active community

**Cons:**
- Resource intensive
- Complex setup
- Memory hungry

**Complexity:** Medium | **Risk:** Low | **Cost:** Low (self-hosted)

---

### Option B: MeiliSearch-Centric

**Architecture:**
```
Content Sources → Extractors → MeiliSearch
                              ↓
                          PostgreSQL (metadata)
```

**Pros:**
- Faster implementation
- Simpler setup
- Better developer experience

**Cons:**
- Less mature than Elasticsearch
- Limited vector support (newer)
- Smaller community

**Complexity:** Low | **Risk:** Medium | **Cost:** Low

---

### Option C: Open Semantic Search

**Architecture:**
```
Open Semantic Search (includes ETL + search + UI)
```

**Pros:**
- Complete solution
- Includes ETL pipeline
- Built-in connectors

**Cons:**
- Tightly coupled
- Python-based (different stack)
- Harder to customize

**Complexity:** Low | **Risk:** Low | **Cost:** Low

---

### Option D: Hybrid (PostgreSQL + Elasticsearch)

**Architecture:**
```
Full-text Search → Elasticsearch
Vector Search → pgvector (PostgreSQL)
Metadata → PostgreSQL
```

**Pros:**
- Best of both worlds
- Distributed can scale
- pgvector excellent

**Cons:**
- Two systems to maintain
- Synchronization complexity

**Complexity:** High | **Risk:** Low | **Cost:** Medium

---

## 5. Chosen Design and Rationale

### Selected: Option D (Hybrid PostgreSQL + Elasticsearch)

**Rationale:**
1. **pgvector is already in use** - Leverages existing infrastructure
2. **Elasticsearch for full-text** - Better for complex queries
3. **PostgreSQL for metadata** - Already the primary DB
4. **Future-proof** - Both technologies are well-maintained

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    RESOURCE SEARCH SYSTEM                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐   │
│  │   YouTube    │    │    OER API   │    │  Web Pages  │   │
│  │   (yt-dlp)   │    │(REST/Scraping)│    │(Playwright) │   │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘   │
│         │                   │                   │             │
│         └───────────────────┼───────────────────┘             │
│                             │                                 │
│                             ▼                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              EXTRACTION PIPELINE                         │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │   │
│  │  │   OCR    │ │   PDF    │ │ Metadata │ │  Text  │ │   │
│  │  │Tesseract │ │  pdfcpu  │ │ExifTool │ │Extract │ │   │
│  │  │  Surya   │ │ Poppler  │ │  Tika   │ │        │ │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └────────┘ │   │
│  └─────────────────────────────────────────────────────────┘   │
│                             │                                 │
│                             ▼                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              CLASSIFICATION ENGINE                       │   │
│  │  - Subject classification (math, pc, svt, philosophy)   │   │
│  │  - Topic extraction (derivative, integral, etc.)        │   │
│  │  - Language detection (AR, FR, EN)                     │   │
│  │  - Difficulty estimation                                │   │
│  └─────────────────────────────────────────────────────────┘   │
│                             │                                 │
│              ┌──────────────┼──────────────┐                 │
│              ▼              ▼              ▼                 │
│  ┌─────────────────┐ ┌─────────────┐ ┌─────────────────┐   │
│  │  Elasticsearch   │ │ PostgreSQL   │ │   Garage S3     │   │
│  │  (Full-text)    │ │ + pgvector   │ │   (Blobs)       │   │
│  │                 │ │ (Vectors)    │ │                 │   │
│  └────────┬────────┘ └──────┬───────┘ └────────┬────────┘   │
│           │                 │                  │             │
│           └─────────────────┼──────────────────┘             │
│                             │                                 │
│                             ▼                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    SEARCH API                            │   │
│  │  - Hybrid search (text + vector)                        │   │
│  │  - Faceted filtering                                    │   │
│  │  - Ranking and relevance                                │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 6. API Design Implications

### Endpoints

```go
// ResourceService handles content operations
type ResourceService interface {
    // Search across all content
    Search(ctx context.Context, req SearchRequest) (*SearchResponse, error)
    
    // Index new content
    Index(ctx context.Context, req IndexRequest) (*IndexResponse, error)
    
    // Get resource metadata
    Get(ctx context.Context, id string) (*Resource, error)
    
    // Find similar resources
    FindSimilar(ctx context.Context, id string, limit int) ([]*Resource, error)
    
    // Delete resource
    Delete(ctx context.Context, id string) error
}

// SearchRequest represents a search query
type SearchRequest struct {
    Query         string   `json:"query"`
    Subject       string   `json:"subject,omitempty"`
    Chapter      int      `json:"chapter,omitempty"`
    Difficulty   int      `json:"difficulty,omitempty"`
    Source       string   `json:"source,omitempty"`
    MediaType    string   `json:"media_type,omitempty"`  // video, pdf, text, image
    Language     string   `json:"language,omitempty"`
    Limit        int      `json:"limit,omitempty"`
    Offset       int      `json:"offset,omitempty"`
}

// SearchResponse contains search results
type SearchResponse struct {
    Results    []SearchResult `json:"results"`
    Total      int           `json:"total"`
    Facets     Facets        `json:"facets"`
    QueryTime  time.Duration `json:"query_time"`
}

// SearchResult represents a single hit
type SearchResult struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Subject     string    `json:"subject"`
    Chapter     int       `json:"chapter"`
    MediaType   string    `json:"media_type"`
    Source      string    `json:"source"`
    URL         string    `json:"url,omitempty"`
    Score       float64   `json:"score"`
    Highlights  []string `json:"highlights,omitempty"`
}
```

### API Contracts

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/resources/search | Search resources |
| POST | /api/v1/resources/index | Index new resource |
| GET | /api/v1/resources/:id | Get resource |
| GET | /api/v1/resources/:id/download | Download content |
| GET | /api/v1/resources/similar/:id | Find similar |
| DELETE | /api/v1/resources/:id | Delete resource |
| GET | /api/v1/resources/facets | Get available facets |

---

## 7. SDK Design Implications

### Go Client

```go
// Client provides resource search functionality
type Client struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

// Search performs a hybrid search
func (c *Client) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)

// Index adds a new resource
func (c *Client) Index(ctx context.Context, req *IndexRequest) (*IndexResponse, error)

// Stream allows searching with filters
func (c *Client) Stream(ctx context.Context, req *SearchRequest) (*SearchIterator, error)
```

### Python Client

```python
class ResourceClient:
    def search(self, query: str, **filters) -> SearchResponse:
        """Full-text + vector hybrid search"""
        
    def index(self, resource: Resource) -> IndexResponse:
        """Add resource to index"""
        
    def similar(self, resource_id: str, limit: int = 10) -> List[Resource]:
        """Find similar resources"""
```

---

## 8. Backward Compatibility and Migration

### Phase 1: Dual Operation
- Run Elasticsearch alongside existing search
- Mirror writes to both systems
- Compare results during validation

### Phase 2: Read Replica
- Direct read queries to Elasticsearch
- Continue writing to both systems
- Monitor for discrepancies

### Phase 3: Full Migration
- Switch all reads to Elasticsearch
- Deprecate old search
- Remove old implementation

### Migration Timeline
- Week 1-2: Infrastructure setup
- Week 3-4: Pipeline development
- Week 5-6: Dual operation testing
- Week 7-8: Gradual rollout

---

## 9. Testing Strategy

### Unit Tests
- [ ] Extractor tests for each content type
- [ ] Classification model tests
- [ ] API endpoint tests
- [ ] Search ranking tests

### Integration Tests
- [ ] End-to-end indexing pipeline
- [ ] Search across multiple sources
- [ ] Faceted filtering
- [ ] Similarity search

### Performance Tests
- [ ] Indexing throughput (docs/second)
- [ ] Search latency (p50, p95, p99)
- [ ] Vector search performance
- [ ] Concurrency handling

---

## 10. Rollout and Observability Plan

### Metrics to Track
- Indexing rate (resources/minute)
- Search latency (p50, p95, p99)
- Cache hit rate
- Error rate by source
- Queue depth

### Dashboards
- Indexing pipeline health
- Search performance
- Source coverage
- Error rates

### Alerting
- Queue backup > 1000 items
- Search latency > 500ms
- Error rate > 5%
- Source unavailable > 5 minutes

---

## 11. Risks and Open Questions

### Risks
| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|-------------|
| Elasticsearch resource usage | High | Medium | Right-size instance, monitor |
| Classification accuracy | Medium | High | Human review queue |
| Source API rate limits | Medium | Medium | Implement backoff |
| Data sync inconsistencies | Low | High | Dual-write validation |

### Open Questions
1. Should we include user-uploaded content in search?
2. How to handle copyrighted content?
3. What retention policy for cached content?
4. How to prioritize sources in ranking?

---

## 12. Implementation Tasks

- [ ] Setup Elasticsearch cluster
- [ ] Create content extractors
- [ ] Build classification pipeline
- [ ] Implement search API
- [ ] Add vector similarity
- [ ] Create dashboard
- [ ] Write tests
- [ ] Document API

---

**RFC Status:** Draft
**Created:** 2026-03-01
**Owner:** BAC Unified Team
