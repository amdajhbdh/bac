# RFC: Search Infrastructure

## 1. Problem Statement

### Problem
Need unified search across all content:
- Questions (existing: pgvector)
- Resources (new: external content)
- Bundles (new: wrapped resources)
- Users (existing: basic)

Current:
- Question search only via pgvector
- No cross-collection search
- Limited ranking

### Requirements
- Sub-100ms latency
- Vector + keyword search
- Faceted filtering
- Ranking optimization
- Multi-language support

---

## 2. Current State

### Existing Search
```
src/agent/internal/
└── memory/
    └── memory.go    # pgvector similarity search
```

### Evidence
- Vector search for questions works
- No full-text search
- No unified index

---

## 3. Goals

1. Hybrid search (vector + keyword)
2. Multi-collection search
3. Faceted filtering
4. Personalized ranking
5. Analytics

---

## 4. Architecture

### Hybrid Search Design

```
┌─────────────────────────────────────────────────────────────────┐
│                       SEARCH INFRASTRUCTURE                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  QUERY                                                            │
│     │                                                              │
│     ▼                                                              │
│  ┌─────────────────────────────┐                               │
│  │      QUERY PARSER            │                               │
│  │  - Parse search query       │                               │
│  │  - Extract filters          │                               │
│  │  - Detect language          │                               │
│  │  - Expand synonyms          │                               │
│  └─────────────┬───────────────┘                               │
│                │                                                  │
│      ┌─────────┴─────────┐                                      │
│      ▼                   ▼                                      │
│  ┌────────────┐   ┌─────────────┐                              │
│  │ FULL-TEXT │   │   VECTOR    │                              │
│  │ (Elastic) │   │ (pgvector)  │                              │
│  └─────┬─────┘   └──────┬──────┘                              │
│        │                 │                                       │
│        └────────┬────────┘                                      │
│                 ▼                                                │
│  ┌─────────────────────────────────────────┐                  │
│  │           RESULT FUSION                   │                  │
│  │  - Reciprocal Rank Fusion               │                  │
│  │  - Score normalization                  │                  │
│  │  - Boost factors                       │                  │
│  └────────────────────┬────────────────────┘                  │
│                       │                                           │
│                       ▼                                           │
│  ┌─────────────────────────────────────────┐                  │
│  │           RERANKING (Optional)            │                  │
│  │  - ML-based reranking                  │                  │
│  │  - Learning to Rank                     │                  │
│  └────────────────────┬────────────────────┘                  │
│                       │                                           │
│                       ▼                                           │
│  ┌─────────────────────────────────────────┐                  │
│  │              RESPONSE                     │                  │
│  │  - Results                              │                  │
│  │  - Facets                               │                  │
│  │  - Suggestions                          │                  │
│  └─────────────────────────────────────────┘                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Collections

| Collection | Type | Index | Search |
|------------|------|-------|--------|
| questions | vector + keyword | pgvector + ES | hybrid |
| resources | vector + keyword | pgvector + ES | hybrid |
| bundles | metadata | PostgreSQL | filtered |
| users | metadata | PostgreSQL | filtered |

---

## 5. Search API

### Unified Search

```go
type SearchRequest struct {
    Query       string        `json:"query"`
    Collections []string     `json:"collections,omitempty"`  // questions, resources, all
    Filters     SearchFilters `json:"filters"`
    Pagination  Pagination    `json:"pagination"`
    Ranking     RankingConfig `json:"ranking"`
}

type SearchFilters struct {
    Subject     string   `json:"subject,omitempty"`
    Chapter    int      `json:"chapter,omitempty"`
    Difficulty int      `json:"difficulty,omitempty"`
    Language   string   `json:"language,omitempty"`
    Source     string   `json:"source,omitempty"`
    MediaType  string   `json:"media_type,omitempty"`
    DateFrom   string   `json:"date_from,omitempty"`
    DateTo     string   `json:"date_to,omitempty"`
}

type SearchResponse struct {
    Results    []SearchHit    `json:"results"`
    Total      int64          `json:"total"`
    Facets     SearchFacets   `json:"facets"`
    Suggest    []string       `json:"suggest,omitempty"`
    QueryTime time.Duration  `json:"query_time_ms"`
}

type SearchHit struct {
    ID          string                 `json:"id"`
    Collection  string                 `json:"collection"`
    Title       string                 `json:"title"`
    Snippet     string                 `json:"snippet"`
    Score       float64                `json:"score"`
    Highlights  []string              `json:"highlights,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
```

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/search | Unified search |
| GET | /api/v1/search/suggest | Autocomplete |
| GET | /api/v1/search/facets | Get facets |
| GET | /api/v1/search/:collection | Collection search |

---

## 6. Query Processing

### Query Pipeline

```
raw query
    │
    ▼
language detection (is it Arabic? French?)
    │
    ▼
tokenization (Arabic: Farasa, French: default)
    │
    ▼
stopwords removal
    │
    ▼
stemming/lemmatization
    │
    ▼
synonym expansion (math → mathematics)
    │
    ▼
query rewrite (add boost for recent, boost for verified)
    │
    ▼
execute search
```

### Multi-Language

| Language | Tokenizer | Stemmer |
|----------|-----------|---------|
| French | standard | french |
| Arabic | Farasa | Farasa |
| English | standard | porter |

---

## 7. Result Fusion

### Reciprocal Rank Fusion

```go
func RRF(results [][]SearchHit, k int) []SearchHit {
    scores := make(map[string]float64)
    
    for _, resultSet := range results {
        for rank, hit := range resultSet {
            scores[hit.ID] += 1.0 / float64(k+rank+1)
        }
    }
    
    // Sort by score
    // Return fused results
}
```

### Score Normalization

```go
func Normalize(scores []float64) []float64 {
    min, max := MinMax(scores)
    normalized := make([]float64, len(scores))
    
    for i, s := range scores {
        if max-min > 0 {
            normalized[i] = (s - min) / (max - min)
        }
    }
    return normalized
}
```

---

## 8. Facets

### Facet Structure

```go
type SearchFacets struct {
    Subject []FacetBucket `json:"subject"`
    Chapter []FacetBucket `json:"chapter"`
    Difficulty []FacetBucket `json:"difficulty"`
    Language []FacetBucket `json:"language"`
    Source   []FacetBucket `json:"source"`
    MediaType []FacetBucket `json:"media_type"`
    Year     []FacetBucket `json:"year"`
}

type FacetBucket struct {
    Value   string `json:"value"`
    Label   string `json:"label"`
    Count   int64  `json:"count"`
}
```

---

## 9. Performance

### Targets

| Metric | Target |
|--------|--------|
| P50 latency | 50ms |
| P95 latency | 100ms |
| P99 latency | 200ms |
| Throughput | 1000 qps |

### Optimization

1. **Caching**
   - Cache frequent queries
   - Cache facets

2. **Indexing**
   - Replica shards
   - Refresh interval

3. **Query**
   - Filter early
   - Limit fields
   - Use filters vs queries

---

## 10. Monitoring

### Metrics
- Query latency (p50, p95, p99)
- Search throughput
- Cache hit rate
- Zero-result queries
- Popular queries

### Alerts
- Latency > 200ms
- Error rate > 1%
- Zero results spike

---

**RFC Status:** Draft
**Created:** 2026-03-01
