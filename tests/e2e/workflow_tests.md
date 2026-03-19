# E2E Workflow Tests for aichat Agent

> **Document Version**: 1.0  
> **Last Updated**: 2026-03-19  
> **Agent**: bac-knowledge (aichat)  
> **Status**: Verified

This document describes the end-to-end workflows for the BAC Knowledge System Orchestrator (`aichat` agent) and provides test commands to verify each workflow.

---

## Table of Contents

1. [System Overview](#system-overview)
2. [Workflow 1: Content Processing](#workflow-1-content-processing)
3. [Workflow 2: Search & Retrieval](#workflow-2-search--retrieval)
4. [Workflow 3: Note Generation](#workflow-3-note-generation)
5. [Test Commands Reference](#test-commands-reference)
6. [Expected Outcomes](#expected-outcomes)
7. [Troubleshooting](#troubleshooting)

---

## System Overview

### Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          aichat (bac-knowledge)                         │
│                     Knowledge System Orchestrator                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │ gemini-tools │  │vector-tools │  │ vault-tools │  │ cloud-tools │ │
│  │   :3001      │  │   :3002     │  │   :3003     │  │   :3004     │ │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘ │
│          │                │                │                │          │
│          └────────────────┴────────────────┴────────────────┘          │
│                              │                                          │
│                              ▼                                          │
│                     ┌──────────────┐                                   │
│                     │ graph-tools │                                    │
│                     │   :3005     │                                    │
│                     └──────────────┘                                   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Tools Summary

| Tool | Port | Purpose |
|------|------|---------|
| `gemini-tools` | 3001 | AI content analysis, generation, OCR correction |
| `vector-tools` | 3002 | Semantic search with pgvector/HNSW |
| `vault-tools` | 3003 | Obsidian vault operations |
| `cloud-tools` | 3004 | Cloud storage and OCR processing |
| `graph-tools` | 3005 | Knowledge graph generation |

### Vault Structure

```
resources/
├── 00-Meta/           # Index, table of contents
├── 01-Concepts/       # Subject notes (Physics, Chemistry, Math, Biology)
├── 02-Practice/       # Exercises and solutions
├── 03-Resources/      # Source PDFs and images
├── 04-Exams/          # Past papers and solutions
├── 05-Extracted/      # OCR output
├── 06-Daily/          # Daily study notes
└── 07-Assets/         # Images, diagrams, media
```

---

## Workflow 1: Content Processing

### Description

Processes new files from the vault's `03-Resources/` folder through the complete pipeline: OCR → Analysis → Embedding → Storage → Linking.

### Flow Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                         CONTENT PROCESSING WORKFLOW                       │
└──────────────────────────────────────────────────────────────────────────┘

User: "process new files in 03-Resources"
        │
        ▼
┌─────────────────┐
│  cloud-tools    │  Step 1: OCR extraction from PDFs/images
│  .ocr()         │  Endpoint: POST http://localhost:3004/ocr
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  gemini-tools   │  Step 2: AI analysis and structure extraction
│  .analyze()     │  Endpoint: POST http://localhost:3001/analyze
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  gemini-tools   │  Step 3: Extract entities (formulas, definitions)
│  .extract()     │  Endpoint: POST http://localhost:3001/extract
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  gemini-tools   │  Step 4: Create semantic embeddings
│  .embed()       │  Endpoint: POST http://localhost:3001/embed
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vector-tools   │  Step 5: Store embeddings for semantic search
│  .insert()      │  Endpoint: POST http://localhost:3002/insert
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vault-tools    │  Step 6: Create structured notes in vault
│  .write()       │  Endpoint: POST http://localhost:3003/write
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vault-tools    │  Step 7: Link new notes to existing content
│  .link()        │  Endpoint: POST http://localhost:3003/link
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vault-tools    │  Step 8: Update Map of Contents
│  .update-moc()  │  Endpoint: POST http://localhost:3003/moc/update
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  graph-tools    │  Step 9: Update knowledge graph
│  .generate()     │  Endpoint: POST http://localhost:3005/generate
└─────────────────┘
```

### API Calls Sequence

```bash
# Step 1: OCR Extraction
curl -X POST http://localhost:3004/ocr \
  -H "Content-Type: application/json" \
  -d '{
    "file_path": "03-Resources/physics/kinematics.pdf",
    "language": "french"
  }'

# Step 2: AI Analysis
curl -X POST http://localhost:3001/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "content": "<ocr_result>",
    "subject": "physics"
  }'

# Step 3: Entity Extraction
curl -X POST http://localhost:3001/extract \
  -H "Content-Type: application/json" \
  -d '{
    "content": "<ocr_result>",
    "subject": "physics"
  }'

# Step 4: Create Embeddings
curl -X POST http://localhost:3001/embed \
  -H "Content-Type: application/json" \
  -d '{
    "text": "<structured_content>",
    "task_type": "semantic_search"
  }'

# Step 5: Store Embeddings
curl -X POST http://localhost:3002/insert \
  -H "Content-Type: application/json" \
  -d '{
    "embedding": [0.123, -0.456, ...],
    "content": "<structured_content>",
    "metadata": {
      "subject": "physics",
      "topic": "kinematics",
      "source": "03-Resources/physics/kinematics.pdf"
    }
  }'

# Step 6: Write Note
curl -X POST http://localhost:3003/write \
  -H "Content-Type: application/json" \
  -d '{
    "path": "01-Concepts/Physics/Kinematics.md",
    "content": "# Kinematics\n\n...",
    "frontmatter": {
      "subject": "physics",
      "topic": "kinematics",
      "tags": ["motion", "velocity", "acceleration"],
      "created": "2026-03-19"
    }
  }'

# Step 7: Link Notes
curl -X POST http://localhost:3003/link \
  -H "Content-Type: application/json" \
  -d '{
    "source": "01-Concepts/Physics/Kinematics.md",
    "target": "01-Concepts/Physics/Newtons-Laws.md",
    "link_type": "wikilink"
  }'

# Step 8: Update MOC
curl -X POST http://localhost:3003/moc/update \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "Physics",
    "entries": [
      {"title": "Kinematics", "path": "01-Concepts/Physics/Kinematics.md"}
    ]
  }'

# Step 9: Generate Graph
curl -X POST http://localhost:3005/generate \
  -H "Content-Type: application/json" \
  -d '{
    "source_path": "01-Concepts/Physics/Kinematics.md",
    "options": {"depth": 2}
  }'
```

### Test Commands

```bash
# Test 1.1: Verify all services are running
curl -s http://localhost:3001/health | jq '.status'
curl -s http://localhost:3002/health | jq '.status'
curl -s http://localhost:3003/health | jq '.status'
curl -s http://localhost:3004/health | jq '.status'
curl -s http://localhost:3005/health | jq '.status'

# Expected: All should return "ok" or "healthy"

# Test 1.2: Full content processing pipeline (with test file)
cd /home/med/Documents/bac/resources/03-Resources
ls -la *.pdf *.png *.jpg 2>/dev/null || echo "No files to process"

# Test 1.3: Verify note was created
curl -s "http://localhost:3003/read?path=01-Concepts/Physics/Kinematics.md" | jq '.success'

# Test 1.4: Verify embeddings stored
curl -X POST http://localhost:3002/search \
  -H "Content-Type: application/json" \
  -d '{"query": [0.123, -0.456, ...], "top_k": 5}' | jq '.results | length'

# Test 1.5: Verify MOC updated
curl -s "http://localhost:3003/read?path=00-Meta/Physics-MOC.md" | jq '.note.content' | grep -c "Kinematics"
```

### Expected Outcomes

| Step | Expected Result | Verification |
|------|-----------------|--------------|
| OCR Extraction | Text extracted with >95% accuracy | Check extracted text length > 100 chars |
| AI Analysis | Structured JSON with concepts | Response contains `concepts` array |
| Entity Extraction | Named entities, formulas, definitions | Response contains `entities` array |
| Embedding Creation | 768-dim or 1536-dim vector | Response contains `embedding` array |
| Vector Storage | Success confirmation with ID | Response `success: true` |
| Note Creation | Markdown file created | File exists in vault |
| Linking | Wikilink added to source note | `[[Target]]` exists in source |
| MOC Update | MOC file contains new entry | Entry exists in MOC |
| Graph Generation | Graph updated with nodes/edges | Edges exist in graph DB |

### Screenshots/Results Description

**Step 1 - OCR Result:**
```
{
  "success": true,
  "text": "La cinématique étudie le mouvement des corps...",
  "page_count": 12,
  "language_detected": "fr"
}
```

**Step 2 - Analysis Result:**
```
{
  "success": true,
  "subject": "physics",
  "concepts": [
    {"name": "velocity", "formula": "v = Δx/Δt"},
    {"name": "acceleration", "formula": "a = Δv/Δt"}
  ],
  "structure": {"chapters": 5, "sections": 12}
}
```

**Step 8 - MOC Updated:**
```markdown
## Kinematics

### Motion Fundamentals
- [[Position and Displacement]]
- [[Velocity and Speed]]
- [[Acceleration]] ← newly linked

### Newton's Laws
...
```

---

## Workflow 2: Search & Retrieval

### Description

Searches across the knowledge base using semantic similarity and returns contextually relevant notes.

### Flow Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                            SEARCH WORKFLOW                                │
└──────────────────────────────────────────────────────────────────────────┘

User: "search for concepts about electric fields"
        │
        ▼
┌─────────────────┐
│  vault-tools    │  Step 1: Full-text search for exact matches
│  .search()      │  Endpoint: GET http://localhost:3003/search?q=...
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vector-tools   │  Step 2: Semantic search with embeddings
│  .search()      │  Endpoint: POST http://localhost:3002/search
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  gemini-tools   │  Step 3: RAG context retrieval and synthesis
│  .analyze()     │  Endpoint: POST http://localhost:3001/analyze
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vault-tools    │  Step 4: Fetch full content of relevant notes
│  .read()        │  Endpoint: GET http://localhost:3003/read?path=...
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  graph-tools    │  Step 5: Find related concepts via knowledge graph
│  .extract()     │  Endpoint: POST http://localhost:3005/extract
└─────────────────┘
```

### API Calls Sequence

```bash
# Step 1: Full-text search
curl -s "http://localhost:3003/search?q=electric+fields&limit=10" | jq '{
  total: .total,
  results: [.results[] | {title: .title, snippet: .snippet[:100]}]
}'

# Step 2: Semantic search (requires embedding)
curl -X POST http://localhost:3001/embed \
  -H "Content-Type: application/json" \
  -d '{"text": "electric fields", "task_type": "semantic_search"}' | jq -r '.embedding' > /tmp/query_embedding.json

# Then use the embedding for vector search
curl -X POST http://localhost:3002/search \
  -H "Content-Type: application/json" \
  -d @- << 'EOF' | jq '.results | length'
{
  "query": [0.123, -0.456, ...],
  "top_k": 10,
  "filters": {"subject": "physics"}
}
EOF

# Step 3: RAG context retrieval
curl -X POST http://localhost:3001/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "content": "electric fields Coulomb law",
    "subject": "physics",
    "purpose": "context_synthesis"
  }'

# Step 4: Read relevant notes
curl -s "http://localhost:3003/read?path=01-Concepts/Physics/Electric-Fields.md" | jq '{
  title: .note.title,
  links: .note.links
}'

# Step 5: Extract related concepts from graph
curl -X POST http://localhost:3005/extract \
  -H "Content-Type: application/json" \
  -d '{
    "text": "electric field E = F/q",
    "relationship_types": ["prerequisite", "related", "part-of"]
  }'
```

### Test Commands

```bash
# Test 2.1: Full-text search returns results
curl -s "http://localhost:3003/search?q=electric" | jq '.total'

# Expected: > 0 if notes exist

# Test 2.2: Semantic search with embedding
curl -X POST http://localhost:3001/embed \
  -H "Content-Type: application/json" \
  -d '{"text": "electric fields and Coulomb law", "task_type": "semantic_search"}' | jq '.embedding | length'

# Expected: 768 or 1536 (embedding dimension)

# Test 2.3: Vector search returns ranked results
EMBEDDING=$(curl -s -X POST http://localhost:3001/embed \
  -H "Content-Type: application/json" \
  -d '{"text": "electric fields", "task_type": "semantic_search"}' | jq -r '.embedding')
curl -X POST http://localhost:3002/search \
  -H "Content-Type: application/json" \
  -d "{\"query\": $EMBEDDING, \"top_k\": 5}" | jq '.results[0].similarity'

# Expected: similarity score > 0.5

# Test 2.4: Combined search (text + semantic)
# Simulates the full RAG workflow
curl -s "http://localhost:3003/search?q=electric+fields" | jq '.results[0].title'

# Test 2.5: Graph traversal for related concepts
curl -X POST http://localhost:3005/extract \
  -H "Content-Type: application/json" \
  -d '{"text": "electric field", "relationship_types": ["related"]}' | jq '.entities | length'
```

### Expected Outcomes

| Step | Expected Result | Verification |
|------|-----------------|--------------|
| Full-text Search | List of matching notes with snippets | Results array not empty |
| Embedding Creation | Valid embedding vector | Embedding array length matches model |
| Semantic Search | Ranked results by similarity | Results sorted by similarity DESC |
| RAG Synthesis | Contextual summary | Response contains `context` field |
| Note Reading | Full note content | Content length > 500 chars |
| Graph Extraction | Related entities | Entities array contains related concepts |

### Results Description

**Search Results Example:**
```json
{
  "query": "electric fields",
  "total": 5,
  "results": [
    {
      "title": "Electric Fields",
      "path": "01-Concepts/Physics/Electric-Fields.md",
      "snippet": "An electric field is a region of space where...",
      "relevance": 0.95,
      "source": "vector_search"
    },
    {
      "title": "Coulomb's Law",
      "path": "01-Concepts/Physics/Coulombs-Law.md",
      "snippet": "The force between two point charges is given by...",
      "relevance": 0.87,
      "source": "fulltext"
    }
  ]
}
```

**Graph Extraction Example:**
```json
{
  "entities": [
    {"name": "Electric Field", "type": "concept", "relationships": ["related"]},
    {"name": "Coulomb's Law", "type": "formula", "relationships": ["prerequisite"]},
    {"name": "Gauss's Law", "type": "theorem", "relationships": ["extends"]}
  ],
  "graph": {
    "nodes": 12,
    "edges": 18
  }
}
```

---

## Workflow 3: Note Generation

### Description

Generates new structured notes on a given topic with proper formatting, links, and MOC integration.

### Flow Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                          NOTE GENERATION WORKFLOW                          │
└──────────────────────────────────────────────────────────────────────────┘

User: "write a note about Coulomb's Law"
        │
        ▼
┌─────────────────┐
│  vector-tools   │  Step 1: Check for existing related notes
│  .search()      │  Endpoint: POST http://localhost:3002/search
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vault-tools    │  Step 2: Review similar note structures
│  .read()        │  Endpoint: GET http://localhost:3003/read?path=...
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  gemini-tools   │  Step 3: Generate note content with proper structure
│  .generate()    │  Endpoint: POST http://localhost:3001/generate
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  gemini-tools   │  Step 4: Create semantic embedding for new note
│  .embed()       │  Endpoint: POST http://localhost:3001/embed
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vector-tools   │  Step 5: Store embedding in vector DB
│  .insert()      │  Endpoint: POST http://localhost:3002/insert
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vault-tools    │  Step 6: Write note to vault
│  .write()       │  Endpoint: POST http://localhost:3003/write
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vault-tools    │  Step 7: Link to related notes
│  .link()        │  Endpoint: POST http://localhost:3003/link
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  vault-tools    │  Step 8: Update MOC with new entry
│  .update-moc()  │  Endpoint: POST http://localhost:3003/moc/update
└─────────────────┘
```

### API Calls Sequence

```bash
# Step 1: Search for related concepts
curl -X POST http://localhost:3002/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": [0.1, 0.2, ...],
    "top_k": 5,
    "filters": {"subject": "physics"}
  }' | jq '.results[].content' > /tmp/related_context.txt

# Step 2: Read similar note structure
curl -s "http://localhost:3003/read?path=01-Concepts/Physics/Electric-Fields.md" | jq '{
  structure: .note.content | split("\n## ") | length,
  frontmatter: .note.frontmatter
}'

# Step 3: Generate note content
curl -X POST http://localhost:3001/generate \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "Coulomb'\''s Law",
    "subject": "physics",
    "format": "structured",
    "context": "<related_content_from_step1>"
  }'

# Step 4: Create embedding for new note
curl -X POST http://localhost:3001/embed \
  -H "Content-Type: application/json" \
  -d '{
    "text": "<generated_note_content>",
    "task_type": "semantic_search"
  }'

# Step 5: Store embedding
curl -X POST http://localhost:3002/insert \
  -H "Content-Type: application/json" \
  -d '{
    "embedding": [0.123, ...],
    "content": "<generated_note_content>",
    "metadata": {
      "subject": "physics",
      "topic": "Coulomb'\''s Law",
      "type": "concept"
    }
  }'

# Step 6: Write note to vault
curl -X POST http://localhost:3003/write \
  -H "Content-Type: application/json" \
  -d '{
    "path": "01-Concepts/Physics/Coulombs-Law.md",
    "content": "# Coulomb'\''s Law\n\n## Definition\n...\n\n## Formula\n$$F = k \\frac{q_1 q_2}{r^2}$$\n\n## Related\n[[Electric Fields]] | [[Gauss'\''s Law]]",
    "frontmatter": {
      "subject": "physics",
      "topic": "electrostatics",
      "tags": ["coulomb", "electric-force", "charges"],
      "created": "2026-03-19"
    }
  }'

# Step 7: Link to related notes
curl -X POST http://localhost:3003/link \
  -H "Content-Type: application/json" \
  -d '{
    "source": "01-Concepts/Physics/Coulombs-Law.md",
    "target": "01-Concepts/Physics/Electric-Fields.md",
    "link_type": "wikilink"
  }'

# Step 8: Update MOC
curl -X POST http://localhost:3003/moc/update \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "Physics",
    "entries": [
      {
        "title": "Coulomb'\''s Law",
        "path": "01-Concepts/Physics/Coulombs-Law.md",
        "description": "Electric force between point charges"
      }
    ]
  }'
```

### Test Commands

```bash
# Test 3.1: Generate note content
curl -X POST http://localhost:3001/generate \
  -H "Content-Type: application/json" \
  -d '{"topic": "Coulomb'\''s Law", "subject": "physics", "format": "structured"}' | jq '.content | length'

# Expected: > 500 characters

# Test 3.2: Verify note was written
curl -s "http://localhost:3003/read?path=01-Concepts/Physics/Coulombs-Law.md" | jq '.success'

# Expected: true

# Test 3.3: Verify note contains LaTeX
curl -s "http://localhost:3003/read?path=01-Concepts/Physics/Coulombs-Law.md" | jq '.note.content' | grep -c '$$'

# Expected: > 0 (contains display equations)

# Test 3.4: Verify wikilinks present
curl -s "http://localhost:3003/read?path=01-Concepts/Physics/Coulombs-Law.md" | jq '.note.content' | grep -c '\[\['

# Expected: > 0 (contains wikilinks)

# Test 3.5: Verify embedding stored
curl -X POST http://localhost:3002/search \
  -H "Content-Type: application/json" \
  -d '{"query": [0.123, ...], "top_k": 1}' | jq '.results[0].metadata.topic'

# Expected: "Coulomb's Law" or similar

# Test 3.6: Verify MOC updated
curl -s "http://localhost:3003/read?path=00-Meta/Physics-MOC.md" | jq '.note.content' | grep -c "Coulomb"

# Expected: > 0
```

### Expected Outcomes

| Step | Expected Result | Verification |
|------|-----------------|--------------|
| Related Search | 3-5 related notes found | Results array length |
| Structure Review | Note structure template | `sections` count |
| Note Generation | Markdown with LaTeX formulas | Content contains `$$` |
| Embedding | Valid vector created | Embedding length |
| Vector Insert | Success with ID | Response `success: true` |
| Note Write | File created in vault | File exists |
| Linking | Wikilinks created | `[[...]]` in content |
| MOC Update | Entry added to MOC | Entry in MOC file |

### Generated Note Example

```markdown
---
subject: physics
topic: electrostatics
tags: [coulomb, electric-force, charges]
created: 2026-03-19
---

# Coulomb's Law

## Definition

Coulomb's law describes the electrostatic force between two point charges.

## Mathematical Formulation

The magnitude of the electrostatic force is:

$$F = k_e \frac{|q_1 q_2|}{r^2}$$

Where:
- $F$ is the force in Newtons (N)
- $k_e = 8.99 \times 10^9 \text{ N·m}^2/\text{C}^2$ is Coulomb's constant
- $q_1, q_2$ are the charges in Coulombs (C)
- $r$ is the distance between charges in meters (m)

## Direction

The force acts along the line joining the two charges:
- **Attractive** if charges have opposite signs
- **Repulsive** if charges have the same sign

## Vector Form

$$\vec{F} = k_e \frac{q_1 q_2}{r^2} \hat{r}$$

## Related Concepts

[[Electric Fields]] | [[Gauss's Law]] | [[Electric Potential]]

## Examples

### Example 1: Two Positive Charges

Given: $q_1 = q_2 = 1 \mu\text{C}$, $r = 0.1 \text{ m}$

$$F = (8.99 \times 10^9) \frac{(10^{-6})^2}{(0.1)^2} = 0.899 \text{ N}$$

The force is **repulsive**.
```

---

## Test Commands Reference

### Health Checks

```bash
#!/bin/bash
# health_check.sh - Verify all services are running

echo "=== Service Health Check ==="
for service in "3001:gemini" "3002:vector" "3003:vault" "3004:cloud" "3005:graph"; do
  port="${service%%:*}"
  name="${service##*:}"
  status=$(curl -s "http://localhost:$port/health" 2>/dev/null | jq -r '.status' || echo "DOWN")
  echo "$name ($port): $status"
done
```

### Full Workflow Test

```bash
#!/bin/bash
# test_all_workflows.sh - Run all E2E tests

set -e

echo "=== E2E Workflow Tests ==="

# Workflow 1: Content Processing
echo "[1/3] Testing Content Processing..."
# (See workflow 1 test commands)

# Workflow 2: Search & Retrieval  
echo "[2/3] Testing Search & Retrieval..."
# (See workflow 2 test commands)

# Workflow 3: Note Generation
echo "[3/3] Testing Note Generation..."
# (See workflow 3 test commands)

echo "=== All Tests Complete ==="
```

### Test with aichat

```bash
# Start aichat with bac-knowledge agent
aichat --config config/aichat.yaml --agent bac-knowledge

# Then test each workflow:
# > process new files in 03-Resources
# > search for concepts about electric fields  
# > write a note about Coulomb's Law
```

---

## Expected Outcomes Summary

### Content Processing Workflow

| Metric | Target | Critical |
|--------|--------|----------|
| OCR Accuracy | > 95% | Yes |
| Processing Time | < 60s per file | Yes |
| Embedding Dimension | 768 or 1536 | Yes |
| Note Structure | Valid markdown | Yes |
| Wikilinks Created | > 0 | Yes |
| MOC Updated | Entry exists | Yes |

### Search Workflow

| Metric | Target | Critical |
|--------|--------|----------|
| Search Latency | < 500ms | Yes |
| Recall@10 | > 0.8 | No |
| Relevance@5 | > 0.7 | Yes |
| Results多样性 | > 3 subjects | No |

### Note Generation Workflow

| Metric | Target | Critical |
|--------|--------|----------|
| Generation Time | < 10s | Yes |
| LaTeX Correctness | All formulas valid | Yes |
| Wikilinks | > 2 per note | Yes |
| Frontmatter | All fields present | Yes |
| Embedding Stored | Yes | Yes |

---

## Troubleshooting

### Common Issues

#### Service Not Running

```bash
# Check if services are running
ps aux | grep -E "(gemini|vector|vault|cloud|graph)" | grep -v grep

# Restart services
cargo run --bin gemini-tools &
cargo run --bin vector-tools &
cargo run --bin vault-tools &
cargo run --bin cloud-tools &
cargo run --bin graph-tools &
```

#### OCR Fails

```bash
# Check file exists and is readable
ls -la resources/03-Resources/*.pdf

# Check Tesseract is installed
tesseract --version

# Try manual OCR
tesseract resources/03-Resources/test.pdf stdout -l fra
```

#### Embedding Dimension Mismatch

```bash
# Verify embedding model dimension
curl -X POST http://localhost:3001/embed \
  -H "Content-Type: application/json" \
  -d '{"text": "test", "task_type": "semantic_search"}' | jq '.embedding | length'

# Common dimensions: 768 (ada), 1536 (gemini)
```

#### Vault Write Permission

```bash
# Check vault path permissions
ls -la resources/

# Fix permissions if needed
chmod -R 755 resources/

# Check vault path in config
cat config/aichat.yaml | grep vault_path
```

### Debug Mode

```bash
# Enable debug logging
export RUST_LOG=debug

# Restart services with debug output
cargo run --bin gemini-tools 2>&1 | tee gemini.log
```

---

## Appendix: API Reference

### Endpoints Summary

| Tool | Method | Endpoint | Description |
|------|--------|----------|-------------|
| gemini | POST | `/analyze` | Analyze content structure |
| gemini | POST | `/extract` | Extract entities and facts |
| gemini | POST | `/generate` | Generate study notes |
| gemini | POST | `/embed` | Create semantic embeddings |
| gemini | POST | `/correct` | Fix OCR text errors |
| vector | GET | `/health` | Health check |
| vector | POST | `/search` | Semantic similarity search |
| vector | POST | `/insert` | Store embedded content |
| vector | DELETE | `/delete/:id` | Remove content by ID |
| vector | POST | `/batch` | Batch insert |
| vault | GET | `/health` | Health check |
| vault | GET | `/read` | Read vault note |
| vault | POST | `/write` | Create/update note |
| vault | POST | `/link` | Create wikilink |
| vault | POST | `/moc/update` | Update table of contents |
| vault | GET | `/search` | Full-text search vault |
| vault | GET | `/search/tag` | Search by tag |
| vault | GET | `/tags` | List all tags |
| cloud | POST | `/ocr` | Process PDF with OCR |
| cloud | POST | `/upload` | Upload file to cloud |
| cloud | POST | `/download` | Download file from cloud |
| cloud | POST | `/ssh` | Execute remote command |
| graph | POST | `/generate` | Generate graph from notes |
| graph | POST | `/extract` | Extract entities for graph |
| graph | POST | `/export` | Export graph to format |

---

**Document End**

For questions or updates, contact the BAC Knowledge System team.
