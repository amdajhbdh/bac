# Open Source Projects Reference - BAC Unified

## Overview

This document catalogs all open source projects used or considered for BAC Unified platform, organized by function.

---

## Core Technologies

### Programming Languages

| Project | Version | Purpose | URL |
|---------|---------|---------|-----|
| Go | 1.22+ | Primary language for API and CLI | golang.org |
| Rust | stable | Animation engine (Noon) | rust-lang.org |
| Python | 3.12+ | ML, scripts | python.org |
| TypeScript | 5.0+ | Frontend | typescriptlang.org |
| Nushell | latest | Legacy CLI | nushell.sh |

### Web Frameworks

| Project | Purpose | URL |
|---------|---------|-----|
| Gin | Go HTTP router | github.com/gin-gonic/gin |
| Echo | Alternative Go framework | github.com/labstack/echo |
| React | UI framework | react.dev |
| Next.js | Full-stack React | nextjs.org |
| Vite | Build tool | vitejs.dev |

### Databases

| Project | Purpose | URL |
|---------|---------|-----|
| PostgreSQL | Primary database | postgresql.org |
| pgvector | Vector similarity | github.com/pgvector/pgvector-go |
| Redis | Cache, sessions | redis.io |
| SQLite | Local cache | sqlite.org |
| Neon | Cloud PostgreSQL | neon.tech |

### Storage

| Project | Purpose | URL |
|---------|---------|-----|
| Garage S3 | Local S3-compatible | garagehq.org |
| MinIO | Production S3 | min.io |

---

## AI & ML

### Local AI

| Project | Purpose | URL |
|---------|---------|-----|
| Ollama | Local LLM inference | ollama.ai |
| llama.cpp | GGUF model inference | github.com/ggerganov/llama.cpp |
| GPT4All | Local models | github.com/nomic-ai/gpt4all |
| LM Studio | Model management | lmstudio.ai |

### Cloud AI Services

| Service | Purpose | API Docs |
|---------|---------|----------|
| NotebookLM | Research, audio | Via nlm CLI |
| DeepSeek | Reasoning | deepseek.com |
| Claude | Anthropic | claude.ai |
| ChatGPT | OpenAI | openai.com |

### ML Pipeline

| Project | Purpose | URL |
|---------|---------|-----|
| scikit-learn | ML models | scikit-learn.org |
| PyTorch | Deep learning | pytorch.org |
| Jupyter | Notebooks | jupyter.org |
| MLflow | ML lifecycle | mlflow.org |

---

## Content Aggregation

### OER Platforms

| Project | Purpose | URL |
|---------|---------|-----|
| Open Semantic Search | Enterprise search | opensemanticsearch.org |
| Harvester | OER aggregation | github.com/surfedushare/harvester |
| VuFind | Library portal | vufind.org |
| Learnweb | Collaborative learning | github.com/l3s-learnweb/learnweb |
| Invenio | Digital library | github.com/inveniosoftware/invenio |

### OER Repositories

| Repository | Content | URL |
|-----------|---------|-----|
| OER Commons | K-20 OER | oercommons.org |
| OpenStax | Textbooks | openstax.org |
| MIT OCW | Courses | ocw.mit.edu |
| Khan Academy | K-12 | khanacademy.org |
| PhET | Simulations | phet.colorado.edu |
| BC Campus | Textbooks | bccampus.ca |
| LibreTexts | STEM | libretexts.org |

### Academic Search

| Engine | Purpose | URL |
|--------|---------|-----|
| BASE | Academic search | base-search.net |
| Google Scholar | Papers | scholar.google.com |
| CORE | Open access | core.ac.uk |
| Semantic Scholar | AI research | semanticscholar.org |

---

## Web Scraping

| Project | Purpose | URL |
|---------|---------|-----|
| Playwright | Browser automation | playwright.dev |
| Puppeteer | Chrome automation | github.com/puppeteer/puppeteer |
| Scrapy | Web scraping | scrapy.org |
| Colly | Go scraping | github.com/gocolly/colly |
| Crawl4AI | AI web scraping | github.com/codebybug/crawl4ai |
| yt-dlp | YouTube extraction | github.com/yt-dlp/yt-dlp |
| gallery-dl | Media downloading | github.com/mikf/gallery-dl |

---

## Document Processing

### OCR

| Project | Purpose | URL |
|---------|---------|-----|
| Tesseract | OCR engine | github.com/tesseract-ocr/tesseract |
| Surya | ML OCR | github.com/VikParuchuri/surya |
| DocTR | Document OCR | github.com/mindee/doctr |

### PDF

| Project | Purpose | URL |
|---------|---------|-----|
| pdfcpu | PDF processing | pdfcpu.io |
| Poppler | PDF utilities | poppler.freedesktop.org |
| MuPDF | PDF rendering | mupdf.com |
| pdfminer | Python PDF | github.com/euske/pdfminer |
| Camelot | Table extraction | github.com/camelot-dev/camelot |

### Content Extraction

| Project | Purpose | URL |
|---------|---------|-----|
| Apache Tika | Content extraction | tika.apache.org |
| ExifTool | Metadata | exiftool.org |
| ffprobe | Media metadata | ffmpeg.org |

---

## Search

| Project | Purpose | URL |
|---------|---------|-----|
| Elasticsearch | Full-text search | elastic.co |
| OpenSearch | Search engine | opensearch.org |
| MeiliSearch | Fast search | meilisearch.com |
| Typesense | Typing search | typesense.org |
| Apache Solr | Enterprise search | solr.apache.org |

---

## File Bundling

### Container Formats

| Project | Purpose | URL |
|---------|---------|-----|
| WRAPX | Metadata containers | github.com/beroccaboy/wrapx-spec |
| ARX | Mountable archives | github.com/jubako/arx |
| Research Object | Scientific bundles | researchobject.org |
| BagIt | Preservation | rfc-editor.org/rfc/rfc8493.html |

### Metadata Standards

| Standard | Purpose | URL |
|----------|---------|-----|
| Dublin Core | General metadata | dublincore.org |
| LOM | Learning objects | ieee.org |
| JSON-LD | Linked data | json-ld.org |
| RDF | Semantic web | w3.org/RDF |
| SKOS | Taxonomies | w3.org/2004/02/skos |

---

## Media Processing

| Project | Purpose | URL |
|---------|---------|-----|
| FFmpeg | Video/audio | ffmpeg.org |
| HandBrake | Encoding | handbrake.fr |
| ImageMagick | Image processing | imagemagick.org |
| Libvips | Fast images | libvips.org |

---

## Authentication & Security

| Project | Purpose | URL |
|---------|---------|-----|
| JWT | Token auth | github.com/golang-jwt/jwt |
| Casbin | Authorization | casbin.org |
| bcrypt | Password hashing | golang.org/x/crypto/bcrypt |
| argon2 | Strong hashing | github.com/p-h-c/argon2 |
| otp | TOTP MFA | github.com/pquerna/otp |

---

## API Development

| Project | Purpose | URL |
|---------|---------|-----|
| Swagger/OpenAPI | API docs | swagger.io |
| swaggo | Go swagger | github.com/swaggo/swag |
| gRPC | RPC | grpc.io |
| buf | Protobuf | buf.build |

---

## Monitoring

### Metrics

| Project | Purpose | URL |
|---------|---------|-----|
| Prometheus | Metrics | prometheus.io |
| Grafana | Visualization | grafana.com |
| OpenTelemetry | Tracing | opentelemetry.io |

### Logging

| Project | Purpose | URL |
|---------|---------|-----|
| slog | Go logging | (built-in) |
| Zap | Go logging | uber-go.github.io/zap |
| Loki | Log aggregation | grafana.com/loki |
| ELK Stack | Log aggregation | elastic.co |

### Error Tracking

| Project | Purpose | URL |
|---------|---------|-----|
| Sentry | Error tracking | sentry.io |
| Rollbar | Monitoring | rollbar.com |

---

## Frontend

### UI

| Project | Purpose | URL |
|---------|---------|-----|
| Tailwind CSS | Styling | tailwindcss.com |
| Radix UI | Components | radix-ui.com |
| Shadcn/UI | UI kit | github.com/shadcn-ui/ui |
| Framer Motion | Animations | framer.com/motion |
| Recharts | Charts | recharts.org |
| React Query | Data fetching | tanstack.com/query |
| Zustand | State | zustand-demo.justforuse |

### Mobile

| Project | Purpose | URL |
|---------|---------|-----|
| PWA | Web mobile | web.dev/progressive-web-apps |
| React Native | Native apps | reactnative.dev |
| Capacitor | Cross-platform | capacitorjs.com |
| Expo | RN dev | expo.dev |

### Offline Storage

| Project | Purpose | URL |
|---------|---------|-----|
| WatermelonDB | Offline DB | watermelonjs.com |
| PouchDB | Offline DB | pouchdb.com |
| RxDB | Real-time DB | rxdb.info |

---

## DevOps

### Container

| Project | Purpose | URL |
|---------|---------|-----|
| Docker | Containers | docker.com |
| Podman | Container engine | podman.io |
| Kubernetes | Orchestration | kubernetes.io |
| K3s | Light K8s | k3s.io |
| Helm | K8s packages | helm.sh |

### CI/CD

| Project | Purpose | URL |
|---------|---------|-----|
| GitHub Actions | CI/CD | github.com/features/actions |
| ArgoCD | GitOps | argoproj.io/argo-cd |
| Jenkins | Automation | jenkins.io |

### Infrastructure

| Project | Purpose | URL |
|---------|---------|-----|
| Cloudflare | CDN, Security | cloudflare.com |
| Traefik | Reverse proxy | traefik.io |
| Nginx | Web server | nginx.com |
| Caddy | Auto HTTPS | caddyserver.com |

---

## Testing

| Project | Purpose | URL |
|---------|---------|-----|
| testify | Go testing | github.com/stretchr/testify |
| Ginkgo | BDD Go | onsi.github.io/ginkgo |
| Playwright | E2E testing | playwright.dev |
| Cypress | JS E2E | cypress.io |
| K6 | Load testing | k6.io |
| Locust | Python load | locust.io |

### Code Quality

| Project | Purpose | URL |
|---------|---------|-----|
| golangci-lint | Go linter | github.com/golangci/golangci-lint |
| gofumpt | Go formatting | github.com/mvdan/gofumpt |
| staticcheck | Go analysis | staticcheck.io |
| SonarQube | Code quality | sonarsource.com |

---

## Content Delivery

### CDN

| Project | Purpose | URL |
|---------|---------|-----|
| Cloudflare | CDN | cloudflare.com |
| jsDelivr | NPM CDN | jsdelivr.com |
| BunnyCDN | Budget CDN | bunnycdn.com |

### Video

| Project | Purpose | URL |
|---------|---------|-----|
| Video.js | Player | videojs.com |
| Plyr | Audio/video | plyr.io |
| Howler.js | Audio | howlerjs.com |

---

## Compression

| Project | Purpose | URL |
|---------|---------|-----|
| zstd | Fast compression | facebook.github.io/zstd |
| LZ4 | Ultra fast | lz4.github.io |
| xz | High compression | tukaani.org/xz |
| pigz | Parallel gzip | zlib.net/pigz |

---

## Knowledge Graph

| Project | Purpose | URL |
|---------|---------|-----|
| Apache Jena | RDF | jena.apache.org |
| Neo4j | Graph DB | neo4j.com |
| rdflib | RDF Python | rdflib.dev |
| GraphDB | RDF store | ontotext.com |

---

## Usage by Phase

### Phase 1: Infrastructure (Q1 2026)

| Project | Use | Status |
|---------|-----|--------|
| PostgreSQL + pgvector | Database | Existing |
| Garage S3 | Storage | Existing |
| Redis | Cache | To add |
| Ollama | Local AI | Existing |
| JWT, bcrypt | Auth | To add |

### Phase 2: Content (Q2-Q3 2026)

| Project | Use | Status |
|---------|-----|--------|
| Playwright | Scraping | Existing |
| yt-dlp | YouTube | Existing |
| Tesseract | OCR | Existing |
| Elasticsearch | Search | To add |
| Apache Tika | Extraction | To add |
| pdfcpu | PDF processing | To add |

### Phase 3: Bundling (Q3 2026)

| Project | Use | Status |
|---------|-----|--------|
| WRAPX | Metadata bundles | RFC |
| ARX | Large archives | RFC |
| Research Object | Scientific | RFC |

### Phase 4: Search (Q4 2026)

| Project | Use | Status |
|---------|-----|--------|
| Elasticsearch | Full-text | RFC |
| pgvector | Vectors | Existing |
| MeiliSearch | Fast search | Consider |

---

## Integration Decision Registry

| Decision | Project | Rationale |
|----------|---------|-----------|
| Database | PostgreSQL + pgvector | Already in use, excellent vector support |
| Storage | Garage S3 | Self-hosted, S3-compatible |
| Local AI | Ollama | Offline capability |
| OCR | Tesseract + Surya | Free, Arabic support |
| Search | Elasticsearch + pgvector | Hybrid search requirement |
| Bundling | WRAPX + ARX | Metadata + mountable |
| Cache | Redis | Sessions, rate limiting |
| Auth | JWT + bcrypt | Industry standard |

---

## References

- All URLs verified as of 2026-03-01
- See `doc/rfc/` for detailed RFCs
- See `doc/PLAN/STRATEGIC_PLAN.md` for complete roadmap

---

*Last Updated: 2026-03-01*
*Part of BAC Unified Documentation*
