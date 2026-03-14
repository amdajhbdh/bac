# Skill: Elasticsearch + pgvector Hybrid Search

## Purpose

Hybrid search combining Elasticsearch (full-text) with pgvector (semantic) for BAC Unified.

## When to use

- Implementing question search
- Building resource lookup
- Creating hybrid query pipeline

## Architecture

```
Content Sources → Extractors → Elasticsearch ← pgvector
                                             ↓
                                    PostgreSQL (metadata)
                                             ↓
                                        Garage S3
```

## Search Flow

1. **Full-text**: Elasticsearch for keyword matching, scoring
2. **Vector**: pgvector for semantic similarity
3. **Merge**: RRF (Reciprocal Rank Fusion) to combine results

## Elasticsearch Index

```json
{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "question": { "type": "text", "analyzer": "french" },
      "subject": { "type": "keyword" },
      "chapter": { "type": "keyword" },
      "answer": { "type": "text" },
      "embedding": { "type": "dense_vector", "dims": 1536 }
    }
  }
}
```

## Hybrid Query (RRF)

```python
def hybrid_search(query, k=10, alpha=0.5):
    # Full-text search
    es_results = elasticsearch.search(
        index="questions",
        body={"query": {"match": {"question": query}}}
    )
    
    # Vector search  
    vector = embed_model.encode(query)
    pg_results = postgres.query("""
        SELECT id, question, embedding <=> $1 AS distance
        FROM questions
        ORDER BY embedding <=> $1
        LIMIT $2
    """, vector, k)
    
    # RRF fusion
    return rrf_fusion(es_results, pg_results, k, alpha)
```

## Go Client

```go
import (
    "github.com/elastic/go-elasticsearch/v8"
)

es, _ := elasticsearch.NewClient(elasticsearch.Config{
    Addresses: []string{"http://localhost:9200"},
})

res, err := es.Search(
    es.Search.WithIndex("questions"),
    es.Search.WithBody(strings.NewReader(`{
        "query": {
            "match": {"question": "dérivée"}
        }
    }`)),
)
```

## Environment

| Variable | Purpose |
|----------|---------|
| `ELASTICSEARCH_URL` | ES endpoint |
| `ELASTICSEARCH_INDEX` | Index name |

## RFC Reference

See `doc/rfc/SEARCH_INFRASTRUCTURE.md` for full architecture.

## Anti-Patterns

- ❌ Using only Elasticsearch (no vectors)
- ❌ Using only pgvector (no full-text)
- ❌ Not handling embedding generation
- ❌ Not syncing ES with PostgreSQL
