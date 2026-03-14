# RFC: File Bundling System (WRAPX/ARX)

## 1. Problem Statement and Constraints

### Problem
BAC Unified needs to package educational resources (PDFs, videos, images, notes) into unified bundles that can:
- Carry comprehensive metadata
- Be searched without extraction
- Work offline
- Handle large files efficiently

Current approach stores raw files in S3 without standardized packaging or rich metadata.

### Constraints
- Must support files up to 10GB
- Must be mountable as read-only filesystem for large archives
- Must preserve original file fidelity
- Must work with existing S3 infrastructure
- Must support Arabic, French, English metadata

---

## 2. Current State Evidence

### Existing Storage
- Garage S3: Raw file storage
- PostgreSQL: Metadata only (limited)
- No unified bundle format

### Evidence Links
- Storage config: `config/storage.yaml`
- Database schema: `sql/schema.sql`

### Gaps
- No standardized bundle format
- No metadata preservation
- No semantic tagging
- No offline bundle access

---

## 3. Goals and Non-Goals

### Goals
1. Implement WRAPX format for metadata-rich bundles
2. Implement ARX format for large mountable archives
3. Create automatic metadata extraction
4. Support bundle versioning
5. Enable offline bundle access

### Non-Goals
1. Real-time bundle streaming
2. Write support (read-only bundles)
3. Encryption at rest (use S3 encryption)
4. Digital rights management

---

## 4. Options Considered

### Option A: WRAPX Only (Selected for small files)

**Format Structure:**
```
my-resource.bacx (or .wrapx)
├── /payload/
│   ├── file1.pdf
│   ├── video.mp4
│   └── image.png
├── /meta/
│   ├── manifest.json      # Full metadata
│   ├── classification.json # Subject, topics
│   ├── tags.json         # User tags
│   └── provenance.json   # Source tracking
├── /embeddings/
│   └── vector.json      # pgvector format
└── /audit/
    └── accesslog.json   # Chain of custody
```

**Pros:**
- Native metadata support
- Self-describing
- JSON-LD compatible
- Policy enforcement possible

**Cons:**
- New format (less tooling)
- Not mountable as FS

**Complexity:** Low | **Risk:** Medium | **Use Case:** Individual resources

---

### Option B: ARX Only (Selected for large archives)

**Format Structure:**
```
large-course.arx
├── /content/
│   ├── /module-1/
│   ├── /module-2/
│   └── /module-3/
├── manifest.json
└── metadata.json
```

**Pros:**
- Mountable as filesystem (FUSE)
- Fast random access
- Proven technology (Rust-based)
- Good compression

**Cons:**
- No built-in metadata
- Requires external index

**Complexity:** Medium | **Risk:** Low | **Use Case:** Large course archives

---

### Option C: Hybrid WRAPX + ARX (SELECTED)

**Decision:** Use both based on file size

- **< 100MB**: WRAPX format
- **>= 100MB**: ARX format

**Rationale:**
- WRAPX: Better for searchable, tagged resources
- ARX: Better for large video courses, offline archives

---

## 5. Chosen Design and Rationale

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    FILE BUNDLING SYSTEM                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  INPUT                                                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │
│  │   PDF    │ │   MP4    │ │   PNG    │ │   TXT    │        │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘        │
│       │            │            │            │                  │
│       └────────────┴────────────┴────────────┘                │
│                            │                                    │
│                            ▼                                    │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              METADATA EXTRACTION                         │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │  │
│  │  │ ExifTool │ │ ffprobe  │ │ pdfinfo  │ │、语言检测 │   │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │  │
│  └─────────────────────────────────────────────────────────┘  │
│                            │                                    │
│                            ▼                                    │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              CLASSIFICATION ENGINE                        │  │
│  │  - Subject: math, pc, svt, philosophy                   │  │
│  │  - Topics: derivative, integral, mechanics...          │  │
│  │  - Language: ar, fr, en                                │  │
│  │  - Difficulty: 1-5                                     │  │
│  └─────────────────────────────────────────────────────────┘  │
│                            │                                    │
│              ┌─────────────┴─────────────┐                    │
│              ▼                           ▼                    │
│  ┌─────────────────────┐    ┌─────────────────────────────┐   │
│  │    WRAPX BUNDLE     │    │       ARX ARCHIVE          │   │
│  │    (< 100MB)        │    │       (>= 100MB)          │   │
│  │                     │    │                             │   │
│  │ ├── payload/        │    │ /content/ (mountable)      │   │
│  │ ├── meta/          │    │ manifest.json              │   │
│  │ ├── embeddings/    │    │ metadata.json              │   │
│  │ └── audit/         │    │ index.json                │   │
│  └─────────┬───────────┘    └─────────────┬───────────────┘   │
│            │                           │                     │
│            └─────────────┬───────────────┘                    │
│                          ▼                                    │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │                    S3 UPLOAD                             │  │
│  │           bac-bundles/bacx/  bac-bundles/arx/           │  │
│  └─────────────────────────────────────────────────────────┘  │
│                          │                                    │
│                          ▼                                    │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              POSTGRES INDEX                             │  │
│  │  - bundle_id, s3_key, metadata, embedding               │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Metadata Schema

```json
{
  "@context": "https://schema.org",
  "@type": "LearningResource",
  "name": "Calculus Fundamentals",
  "description": "Introduction to differential calculus",
  "subject": "math",
  "topics": ["derivative", "limit", "function"],
  "difficulty": 2,
  "language": "fr",
  "typicalAgeRange": "15-18",
  "educationalLevel": "BAC C",
  "inLanguage": ["fr", "ar"],
  "learningResourceType": ["course", "lecture notes"],
  "about": [
    {"@type": "Thing", "name": "Mathematics"}
  ],
  "creator": {
    "@type": "Organization",
    "name": "BAC Unified"
  },
  "dateCreated": "2026-01-15",
  "version": "1.0",
  "provenance": {
    "source": "openstax",
    "sourceUrl": "https://openstax.org/books/calculus",
    "license": "CC BY 4.0"
  },
  "statistics": {
    "downloads": 1234,
    "views": 5678,
    "rating": 4.5
  }
}
```

---

## 6. API Design Implications

### Bundle Operations

```go
type BundleService interface {
    // Create new bundle from files
    Create(ctx context.Context, req CreateBundleRequest) (*Bundle, error)
    
    // Read bundle metadata
    Get(ctx context.Context, id string) (*Bundle, error)
    
    // Download bundle
    Download(ctx context.Context, id string) (io.ReadCloser, error)
    
    // Extract/search within bundle
    Search(ctx context.Context, id string, query string) ([]string, error)
    
    // Update metadata
    UpdateMetadata(ctx context.Context, id string, meta *Metadata) error
    
    // Delete bundle
    Delete(ctx context.Context, id string) error
    
    // Convert format
    Convert(ctx context.Context, id string, format Format) (*Bundle, error)
}

type CreateBundleRequest struct {
    Files       []io.Reader
    FileNames   []string
    Metadata    *Metadata
    Format      Format  // WRAPX or ARX
    Compression string  // zstd, lz4, none
}

type Bundle struct {
    ID           string    `json:"id"`
    Format       Format    `json:"format"`
    S3Key        string    `json:"s3_key"`
    SizeBytes    int64     `json:"size_bytes"`
    Metadata     *Metadata `json:"metadata"`
    CreatedAt    time.Time `json:"created_at"`
    Version      string    `json:"version"`
}
```

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/bundles | Create new bundle |
| GET | /api/v1/bundles/:id | Get bundle metadata |
| GET | /api/v1/bundles/:id/download | Download bundle |
| GET | /api/v1/bundles/:id/search | Search within bundle |
| PUT | /api/v1/bundles/:id/metadata | Update metadata |
| DELETE | /api/v1/bundles/:id | Delete bundle |
| POST | /api/v1/bundles/:id/convert | Convert format |

---

## 7. SDK Design

### Go SDK

```go
// BundleClient handles bundle operations
type BundleClient struct {
    endpoint string
    client  *http.Client
}

// CreateBundle creates a new resource bundle
func (c *BundleClient) CreateBundle(ctx context.Context, files []string, meta *Metadata) (*Bundle, error)

// SearchWithin searches text content inside a bundle
func (c *BundleClient) SearchWithin(ctx context.Context, bundleID, query string) ([]SearchResult, error)

// MountBundle mounts an ARX archive as a filesystem
func (c *BundleClient) MountBundle(ctx context.Context, bundleID, mountPoint string) error
```

### CLI

```bash
# Create WRAPX bundle
bac bundle create --files "*.pdf" --metadata meta.json --format wrapx

# Create ARX archive  
bac archive create --dir ./course --compression zstd

# Search within bundle
bac bundle search <bundle-id> "dérivée"

# Extract metadata
bac bundle info <bundle-id>


bac bundle# Convert format convert <bundle-id> --to arx
```

---

## 8. Backward Compatibility

### Phase 1: Coexistence
- Store bundles alongside raw files
- Index bundles in new system
- Keep old metadata in PostgreSQL

### Phase 2: Migration
- Convert existing resources to bundles
- Update search to use bundle metadata
- Deprecate raw file uploads

### Phase 3: Enforce
- Require bundle format for uploads
- Remove raw file support
- Full migration complete

---

## 9. Testing Strategy

### Unit Tests
- [ ] Metadata extraction for each type
- [ ] Format conversion
- [ ] Validation
- [ ] Compression/decompression

### Integration Tests
- [ ] Full bundle creation pipeline
- [ ] S3 upload/download
- [ ] Search within bundles
- [ ] ARX mounting

### Performance Tests
- [ ] Bundle creation time
- [ ] Search latency
- [ ] Mount time for ARX
- [ ] Memory usage

---

## 10. Observability

### Metrics
- Bundles created (by format)
- Bundle size distribution
- Search queries
- Conversion operations
- Storage usage

### Logging
- Bundle creation events
- Access patterns
- Errors and failures

### Alerts
- Failed bundle creation > 5%
- S3 upload failures
- Storage quota warnings

---

## 11. Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|-------------|
| Format obsolescence | Low | High | Use standard formats |
| Large file handling | Medium | Medium | ARX for >100MB |
| Metadata loss | Low | High | Validation checks |
| S3 costs | Medium | Medium | Lifecycle policies |

---

## 12. Implementation

### Tasks
- [ ] Create bundle creation library
- [ ] Implement metadata extraction
- [ ] Build WRAPX encoder/decoder
- [ ] Integrate ARX library
- [ ] Create S3 integration
- [ ] Add search within bundles
- [ ] Build API endpoints
- [ ] Write tests

---

**RFC Status:** Draft
**Created:** 2026-03-01
**Owner:** BAC Unified Team
