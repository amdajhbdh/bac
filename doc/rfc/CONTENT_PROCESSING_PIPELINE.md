# RFC: Content Processing Pipeline

## 1. Problem Statement

### Problem
BAC Unified needs to automatically extract, process, and classify content from multiple sources:
- PDFs (exam papers, textbooks)
- Images (diagrams, handwritten notes)
- Videos (YouTube, educational content)
- Web pages (OER repositories)

Current state:
- Tesseract for OCR (existing)
- No unified pipeline
- Manual classification

### Requirements
- Automatic text extraction
- OCR for images and PDFs
- Language detection (AR, FR, EN)
- Subject classification
- Quality assessment

---

## 2. Current State Evidence

### Existing Components
```
src/agent/internal/
├── ocr/              # Tesseract wrapper
├── online/           # YouTube/web scraping
└── analyzer/         # Multi-service analysis
```

### Evidence Links
- OCR: `internal/ocr/ocr.go` - Tesseract integration
- Video: `internal/online/` - yt-dlp

### Gaps
- No PDF processing pipeline
- No metadata extraction
- No language detection
- No automatic classification

---

## 3. Goals and Non-Goals

### Goals
1. Unified content extraction pipeline
2. Multi-format OCR (image, PDF, handwritten)
3. Automatic classification
4. Language detection
5. Quality scoring

### Non-Goals
1. Real-time streaming processing
2. Video frame extraction
3. Audio transcription
4. Translation

---

## 4. Options Considered

### Option A: Chain of CLI Tools

**Architecture:**
```
Input → Tesseract → pdfcpu → ffprobe → ExifTool → Output
```

**Pros:**
- Simple to implement
- Easy to debug
- Flexible

**Cons:**
- No parallel processing
- Error handling complex
- Slow for large batches

**Complexity:** Low | **Risk:** Low

---

### Option B: Processing Queue (Selected)

**Architecture:**
```
Input → Queue → Workers → Results → Index
         ↑
     Redis/PostgreSQL
```

**Pros:**
- Scalable
- Parallel processing
- Error resilience
- Monitoring

**Cons:**
- More complex setup
- Infrastructure needed

**Complexity:** Medium | **Risk:** Low

---

### Option C: Distributed Processing

**Architecture:**
```
Input → API → Kubernetes Jobs → Results
```

**Pros:**
- Highly scalable
- Auto-scaling
- Managed

**Cons:**
- Complex infrastructure
- Cost management

**Complexity:** High | **Risk:** Medium

---

## 5. Chosen Design

### Pipeline Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    CONTENT PROCESSING PIPELINE                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────┐                                                   │
│  │  INPUT    │                                                   │
│  │  Queue    │                                                   │
│  └────┬─────┘                                                   │
│       │                                                          │
│       ▼                                                          │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  STAGE 1: RECEIPT                         │   │
│  │  - Validate file format                                  │   │
│  │  - Calculate checksum                                    │   │
│  │  - Determine processing type                             │   │
│  │  - Create processing job                                  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│                              ▼                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  STAGE 2: EXTRACTION                      │   │
│  │                                                           │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │   │
│  │  │   PDF     │ │   OCR    │ │  Video   │ │  Web   │ │   │
│  │  │ (pdfcpu)  │ │(Tesseract│ │ (yt-dlp) │ │Scraper │ │   │
│  │  │  Poppler  │ │  Surya)  │ │          │ │        │ │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └────────┘ │   │
│  │                                                           │   │
│  │  Output: Extracted text, images, metadata               │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│                              ▼                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  STAGE 3: ANALYSIS                        │   │
│  │                                                           │   │
│  │  ┌──────────────┐ ┌──────────────┐ ┌─────────────┐  │   │
│  │  │    Language   │ │  Subject     │ │   Topic    │  │   │
│  │  │   Detection   │ │Classification │ │ Extraction  │  │   │
│  │  │  (fastText)   │ │   (ML)       │ │   (NLP)     │  │   │
│  │  └──────────────┘ └──────────────┘ └─────────────┘  │   │
│  │                                                           │   │
│  │  Output: Language, subject, topics, entities           │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│                              ▼                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  STAGE 4: ENRICHMENT                     │   │
│  │                                                           │   │
│  │  - Generate embeddings (pgvector)                        │   │
│  │  - Calculate quality score                                │   │
│  │  - Extract keywords                                      │   │
│  │  - Create summary                                        │   │
│  │                                                           │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                  │
│                              ▼                                  │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  STAGE 5: STORAGE                         │   │
│  │                                                           │   │
│  │  - Save to S3                                            │   │
│  │  - Index in PostgreSQL                                    │   │
│  │  - Add to search index                                   │   │
│  │                                                           │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Processing Types

| Type | Tools | Output |
|-------|-------|--------|
| PDF Text | pdfcpu, Poppler | Plain text |
| PDF OCR | Tesseract, Surya | Plain text |
| Image OCR | Tesseract | Plain text |
| Video Audio | yt-dlp | Audio file |
| Video Metadata | ffprobe | JSON metadata |
| Web Content | Playwright | HTML, text |

---

## 6. API Design

### Job Management

```go
type ProcessingJob struct {
    ID          string         `json:"id"`
    InputURL    string        `json:"input_url"`
    InputType   string        `json:"input_type"`  // pdf, image, video, url
    Status      JobStatus     `json:"status"`
    Progress    int           `json:"progress"`    // 0-100
    Stage       string        `json:"stage"`       // receipt, extraction, analysis, enrichment, storage
    Result      *JobResult    `json:"result,omitempty"`
    Error       string        `json:"error,omitempty"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
}

type JobResult struct {
    Text          string            `json:"text,omitempty"`
    Language      string            `json:"language,omitempty"`
    Subject       string            `json:"subject,omitempty"`
    Topics        []string          `json:"topics,omitempty"`
    Embedding     []float32         `json:"embedding,omitempty"`
    QualityScore float32           `json:"quality_score,omitempty"`
    Metadata      map[string]string `json:"metadata,omitempty"`
}
```

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/process | Submit processing job |
| GET | /api/v1/process/:id | Get job status |
| GET | /api/v1/process/:id/result | Get result |
| DELETE | /api/v1/process/:id | Cancel job |
| GET | /api/v1/process | List jobs |

---

## 7. Worker Architecture

```go
type Worker struct {
    ID         string
    Queue      string
    Processes  []Processor
    Metrics    *Metrics
}

type Processor interface {
    CanProcess(job *ProcessingJob) bool
    Process(ctx context.Context, job *ProcessingJob) (*JobResult, error)
}

// Processors
type PDFTextProcessor struct{}
type PDFOCRProcessor struct{}
type ImageOCRProcessor struct{}
type VideoProcessor struct{}
type WebProcessor struct{}
```

### Processing Flow

```
1. Worker pulls job from queue
2. Validate job parameters
3. Download input file (if URL)
4. Extract text based on type
5. Run analysis (language, subject, topics)
6. Generate embeddings
7. Calculate quality score
8. Upload to S3
9. Index in PostgreSQL
10. Update job status
11. Acknowledge completion
```

---

## 8. Error Handling

### Retry Strategy

| Stage | Retries | Backoff |
|-------|----------|---------|
| Download | 3 | Exponential |
| OCR | 2 | Linear |
| Analysis | 3 | Exponential |
| Storage | 5 | Linear |

### Error Types

```go
type ProcessingError struct {
    Code    string  // DOWNLOAD_FAILED, OCR_FAILED, etc.
    Message string
    Stage  string
    Retry  bool
}
```

---

## 9. Monitoring

### Metrics

| Metric | Type | Description |
|--------|------|-------------|
| jobs_submitted | Counter | Total jobs |
| jobs_completed | Counter | Successful jobs |
| jobs_failed | Counter | Failed jobs |
| processing_duration | Histogram | Time per job |
| queue_depth | Gauge | Pending jobs |
| stage_duration | Histogram | Time per stage |

### Dashboards
- Job throughput
- Queue depth by type
- Error rate by stage
- Processing latency

---

## 10. Implementation

### Tasks
- [ ] Create job queue (Redis)
- [ ] Implement worker pool
- [ ] Build PDF processor
- [ ] Build OCR processor
- [ ] Build video processor
- [ ] Add language detection
- [ ] Add classification
- [ ] Build API
- [ ] Add monitoring

---

**RFC Status:** Draft
**Created:** 2026-03-01
