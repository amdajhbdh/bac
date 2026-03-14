# BAC HYPER-SWARM 2026 - MEGA ENHANCEMENT PLAN 🚀

**Status**: ~80% → 100% + Game-Changing Features  
**Target**: Production-ready national exam system for 50k+ students  
**Timeline**: 12 weeks (3 months)  
**Last Updated**: 2026-03-01

---

## EXECUTIVE SUMMARY

Transform BAC from a prototype into a world-class, AI-powered exam preparation platform by:
1. **Completing all stubs** (Weeks 1-3)
2. **Adding 8 killer features** (Weeks 4-9)
3. **Production hardening** (Weeks 10-12)

**Key Innovations**:
- 🤖 Auto-animation generation using Manimator approach
- 🧠 Spaced repetition with FSRS algorithm
- 📊 Real-time collaborative study rooms
- 🎯 ML-powered exam prediction
- 🌍 Full Arabic/French/English support
- 📱 Progressive Web App (offline-first)
- 🏆 Gamification with badges/leaderboards
- 📈 Advanced analytics dashboard

---

## PHASE 1: COMPLETE THE FOUNDATION (Weeks 1-3)

### Week 1: Critical Stubs & Playwright Automation

**Priority**: HIGH - These block other features

#### 1.1 Complete Playwright Automation
**Files**: `src/agent/internal/online/playwright.go`, `online.go`

**Tasks**:
- [ ] Test DeepSeek automation (manual verification)
- [ ] Test Grok automation (requires X account)
- [ ] Test Claude automation (requires Anthropic account)
- [ ] Test ChatGPT automation (free tier)
- [ ] Add retry logic with exponential backoff
- [ ] Implement session persistence across restarts
- [ ] Add rate limit detection and handling
- [ ] Create fallback chain: DeepSeek → ChatGPT → Claude → Grok

**New Tool**: Use `playwright-python` instead of CLI for better control

```bash
pip install playwright
playwright install chromium
```

**Implementation**:
```go
// src/agent/internal/online/playwright_v2.go
package online

import (
    "context"
    "github.com/playwright-community/playwright-go"
)

type PlaywrightV2 struct {
    browser playwright.Browser
    context playwright.BrowserContext
}

func NewPlaywrightV2() (*PlaywrightV2, error) {
    pw, err := playwright.Run()
    if err != nil {
        return nil, err
    }
    
    browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
        Headless: playwright.Bool(true),
    })
    
    // Load persistent session
    ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{
        StorageState: playwright.String("session.json"),
    })
    
    return &PlaywrightV2{browser: browser, context: ctx}, nil
}
```

#### 1.2 Complete Memory Module
**Files**: `src/agent/internal/memory/memory.go`

**Current Issue**: Uses hash-based fallback instead of proper embeddings

**Tasks**:
- [ ] Integrate proper embedding model (all-MiniLM-L6-v2)
- [ ] Add embedding cache to avoid recomputation
- [ ] Implement hybrid search (vector + keyword)
- [ ] Add relevance scoring with threshold

**New Tool**: `sentence-transformers` via Ollama

```bash
ollama pull nomic-embed-text
```

```go
func generateEmbedding(text string) ([]float32, error) {
    resp, err := http.Post(
        "http://localhost:11434/api/embeddings",
        "application/json",
        strings.NewReader(fmt.Sprintf(`{"model":"nomic-embed-text","prompt":"%s"}`, text)),
    )
    // Parse and return embedding
}
```

#### 1.3 Complete Animation Module
**Files**: `src/agent/internal/animation/animation.go`

**Tasks**:
- [ ] Fix Cargo.toml generation (add all dependencies)
- [ ] Implement video export with ffmpeg
- [ ] Add animation preview (GIF generation)
- [ ] Create animation templates library
- [ ] Add error recovery (fallback to static image)

**New Tool**: `gifski` for high-quality GIF previews

```bash
cargo install gifski
```

---

### Week 2: Database & API Completion

#### 2.1 Complete Database Schema
**Files**: `sql/schema.sql`, `src/agent/internal/db/`

**Tasks**:
- [ ] Add missing tables: `study_sessions`, `flashcards`, `badges`, `leaderboard`
- [ ] Create indexes for performance
- [ ] Add full-text search indexes
- [ ] Implement database migrations system
- [ ] Add seed data for testing

**New Tool**: `golang-migrate` for migrations

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Schema additions**:
```sql
-- Spaced repetition
CREATE TABLE flashcards (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    question_id UUID REFERENCES questions(id),
    front TEXT NOT NULL,
    back TEXT NOT NULL,
    ease_factor FLOAT DEFAULT 2.5,
    interval INT DEFAULT 0,
    repetitions INT DEFAULT 0,
    next_review TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Gamification
CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    description TEXT,
    icon_url TEXT,
    criteria JSONB
);

CREATE TABLE user_badges (
    user_id UUID REFERENCES users(id),
    badge_id INT REFERENCES badges(id),
    earned_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, badge_id)
);

-- Analytics
CREATE TABLE study_sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    questions_attempted INT,
    questions_correct INT,
    topics JSONB,
    metadata JSONB
);
```

#### 2.2 Complete API Services
**Files**: `src/api/internal/services/*.go`

**Tasks**:
- [ ] Wire all handler-to-service connections
- [ ] Add authentication middleware
- [ ] Implement rate limiting
- [ ] Add request validation
- [ ] Create API documentation (Swagger)

**New Tool**: `swaggo` for auto-generated API docs

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g src/api/main.go
```

---

### Week 3: OCR & Document Processing

#### 3.1 Upgrade OCR Pipeline
**Files**: `src/agent/internal/ocr/ocr.go`

**Current**: Basic Tesseract + fallbacks  
**Target**: State-of-the-art multi-modal OCR

**Tasks**:
- [ ] Integrate Surya OCR for layout analysis
- [ ] Add PaddleOCR for Arabic handwriting
- [ ] Implement table extraction
- [ ] Add equation recognition (LaTeX output)
- [ ] Create OCR quality scoring

**New Tools**:
1. **Surya** - Layout + OCR in 90+ languages
```bash
pip install surya-ocr
```

2. **PaddleOCR** - Best for Arabic
```bash
pip install paddleocr
```

3. **Pix2Tex** - Image to LaTeX
```bash
pip install pix2tex
```

**Implementation**:
```go
// src/agent/internal/ocr/surya.go
package ocr

import "os/exec"

func ProcessWithSurya(imagePath string) (*SuryaResult, error) {
    cmd := exec.Command("python", "-c", `
from surya.ocr import run_ocr
from surya.model.detection.model import load_model, load_processor
from PIL import Image

image = Image.open('` + imagePath + `')
model, processor = load_model(), load_processor()
predictions = run_ocr([image], [["en", "ar", "fr"]], model, processor)
print(predictions[0].text_lines)
`)
    output, err := cmd.Output()
    // Parse and return
}
```

#### 3.2 Advanced Document Understanding
**Files**: `src/agent/internal/pdf/extractor.go`

**Tasks**:
- [ ] Integrate LangExtract for structured data
- [ ] Add document classification
- [ ] Extract diagrams and figures
- [ ] Implement smart chunking for RAG

**New Tool**: Google's **LangExtract**

```bash
pip install langextract
```

---

## PHASE 2: KILLER FEATURES (Weeks 4-9)

### Week 4-5: Feature 1 - Auto-Animation Generation (Manimator Approach)

**Goal**: Generate 3Blue1Brown-style animations automatically from solutions

**Architecture**:
```
Problem → Solver → Solution Steps → LLM → Manim/Noon Code → Video
```

**Implementation**:

#### 4.1 Create Animation Generator Service
**File**: `src/agent/internal/animation/manimator.go`

```go
package animation

type Manimator struct {
    llm LLMClient
    noonPath string
}

func (m *Manimator) GenerateFromSolution(solution SolveResult) (*Animation, error) {
    // Step 1: Analyze solution and extract visual elements
    analysis := m.analyzeSolution(solution)
    
    // Step 2: Generate scene description
    sceneDesc := m.generateSceneDescription(analysis)
    
    // Step 3: Convert to Noon code
    noonCode := m.generateNoonCode(sceneDesc)
    
    // Step 4: Compile and render
    video := m.compileAndRender(noonCode)
    
    return video, nil
}
```

**LLM Prompt Template**:
```
You are an expert at creating mathematical animations using the Noon engine (Rust-based, similar to Manim).

Given this solution:
{solution_steps}

Generate Noon code that:
1. Shows the problem statement
2. Animates each solution step
3. Highlights key concepts
4. Uses appropriate colors and transitions

Output valid Rust code using the Noon API.
```

**New Tool**: Fine-tune a small model for code generation

```bash
# Use Qwen3-Coder (best open-source coding model 2026)
ollama pull qwen3-coder:7b
```

#### 4.2 Create Animation Templates Library
**Directory**: `src/agent/internal/animation/templates/`

**Templates**:
- `algebra.rs` - Equation solving
- `calculus.rs` - Derivatives, integrals
- `geometry.rs` - Shapes, proofs
- `physics.rs` - Forces, circuits
- `chemistry.rs` - Molecules, reactions

---

### Week 6: Feature 2 - Spaced Repetition System (FSRS)

**Goal**: Implement scientifically-proven spaced repetition

**Algorithm**: FSRS (Free Spaced Repetition Scheduler) - used by Anki

**Implementation**:

#### 6.1 FSRS Core Engine
**File**: `src/agent/internal/srs/fsrs.go`

```go
package srs

type Card struct {
    ID           string
    Difficulty   float64  // 1-10
    Stability    float64  // Days
    Retrievability float64 // 0-1
    LastReview   time.Time
    NextReview   time.Time
}

type FSRS struct {
    w []float64 // 17 parameters
}

func NewFSRS() *FSRS {
    // Default FSRS-5 parameters
    return &FSRS{
        w: []float64{0.4072, 1.1829, 3.1262, 15.4722, 7.2102, 
                     0.5316, 1.0651, 0.0234, 1.616, 0.1544,
                     1.0824, 1.9813, 0.0953, 0.2975, 2.2042,
                     0.2407, 2.9466},
    }
}

func (f *FSRS) Schedule(card *Card, rating int) time.Time {
    // FSRS algorithm implementation
    // Returns next review date
}
```

**Reference**: https://github.com/open-spaced-repetition/fsrs-rs

#### 6.2 Flashcard Service
**File**: `src/api/internal/services/flashcard.go`

**Features**:
- Auto-generate flashcards from questions
- Track review history
- Adaptive scheduling
- Export to Anki format

---

### Week 7: Feature 3 - Real-time Collaborative Study Rooms

**Goal**: Students can study together in real-time

**Tech Stack**:
- WebSockets (Go)
- Redis for pub/sub
- Presence tracking

**Implementation**:

#### 7.1 WebSocket Server
**File**: `src/api/internal/realtime/rooms.go`

```go
package realtime

import (
    "github.com/gorilla/websocket"
    "github.com/redis/go-redis/v9"
)

type StudyRoom struct {
    ID          string
    Topic       string
    Participants map[string]*Participant
    Messages    []Message
    redis       *redis.Client
}

func (r *StudyRoom) Join(userID string, conn *websocket.Conn) {
    // Add participant
    // Broadcast join event
    // Sync current state
}

func (r *StudyRoom) SendMessage(msg Message) {
    // Broadcast to all participants
    r.redis.Publish(ctx, "room:"+r.ID, msg)
}
```

#### 7.2 Features
- [ ] Text chat
- [ ] Shared whiteboard
- [ ] Question sharing
- [ ] Voice chat (WebRTC)
- [ ] Screen sharing

**New Tools**:
- `gorilla/websocket` - WebSocket library
- `pion/webrtc` - WebRTC for Go

---

### Week 8: Feature 4 - ML-Powered Exam Prediction

**Goal**: Predict likely exam questions using ML

**Approach**:
1. Pattern analysis of past 20 years
2. Topic frequency analysis
3. Difficulty progression
4. LLM-based question generation

**Implementation**:

#### 8.1 Pattern Analyzer
**File**: `src/agent/internal/prediction/analyzer.go`

```go
package prediction

type PatternAnalyzer struct {
    db *pgxpool.Pool
}

func (p *PatternAnalyzer) AnalyzeHistoricalPatterns(subject string, years int) *Patterns {
    // Query past exams
    // Extract topics, difficulty, question types
    // Find recurring patterns
    // Calculate probabilities
}

func (p *PatternAnalyzer) PredictTopics(subject string, year int) []TopicPrediction {
    patterns := p.AnalyzeHistoricalPatterns(subject, 20)
    
    // ML model: Random Forest or XGBoost
    predictions := p.model.Predict(patterns)
    
    return predictions
}
```

**New Tool**: `goml` - ML library for Go

```bash
go get github.com/cdipaolo/goml
```

#### 8.2 Question Generator
**File**: `src/agent/internal/prediction/generator.go`

Use LLM to generate practice questions based on predictions:

```go
func (g *Generator) GenerateQuestions(topic TopicPrediction) []Question {
    prompt := fmt.Sprintf(`
Generate 5 BAC-level exam questions about: %s
Difficulty: %d/10
Style: Similar to Mauritanian BAC exams
Include: Multiple choice and essay questions
`, topic.Name, topic.Difficulty)
    
    response := g.llm.Generate(prompt)
    return parseQuestions(response)
}
```

---

### Week 9: Feature 5 - Progressive Web App (PWA)

**Goal**: Mobile-friendly, offline-capable web interface

**Tech Stack**:
- React + TypeScript
- Vite for build
- Workbox for service worker
- IndexedDB for offline storage

**Implementation**:

#### 9.1 Create Web Frontend
**Directory**: `src/web-v2/`

```bash
cd src
npm create vite@latest web-v2 -- --template react-ts
cd web-v2
npm install workbox-webpack-plugin localforage
```

**Features**:
- [ ] Responsive design (mobile-first)
- [ ] Offline mode with sync
- [ ] Push notifications
- [ ] Install as app
- [ ] Dark mode

#### 9.2 Service Worker
**File**: `src/web-v2/src/sw.ts`

```typescript
import { precacheAndRoute } from 'workbox-precaching';
import { registerRoute } from 'workbox-routing';
import { CacheFirst, NetworkFirst } from 'workbox-strategies';

// Precache static assets
precacheAndRoute(self.__WB_MANIFEST);

// Cache API responses
registerRoute(
  /\/api\//,
  new NetworkFirst({
    cacheName: 'api-cache',
    networkTimeoutSeconds: 3,
  })
);

// Cache images
registerRoute(
  /\.(?:png|jpg|jpeg|svg|gif)$/,
  new CacheFirst({
    cacheName: 'image-cache',
  })
);
```

---

## PHASE 3: PRODUCTION HARDENING (Weeks 10-12)

### Week 10: Performance & Scalability

#### 10.1 Database Optimization
- [ ] Add connection pooling
- [ ] Implement query caching
- [ ] Create materialized views
- [ ] Add read replicas
- [ ] Optimize indexes

#### 10.2 API Optimization
- [ ] Add response caching (Redis)
- [ ] Implement rate limiting
- [ ] Add request batching
- [ ] Enable HTTP/2
- [ ] Add CDN for static assets

**New Tool**: `Cloudflare Workers` for edge caching

---

### Week 11: Monitoring & Observability

#### 11.1 Metrics
**Tool**: Prometheus + Grafana

```bash
docker-compose up -d prometheus grafana
```

**Metrics to track**:
- Request latency (p50, p95, p99)
- Error rates
- Database query time
- Cache hit rate
- Active users
- Question solve rate

#### 11.2 Logging
**Tool**: Loki + Promtail

**Structured logging**:
```go
slog.Info("question solved",
    "user_id", userID,
    "question_id", questionID,
    "duration_ms", duration,
    "model", model,
    "confidence", confidence,
)
```

#### 11.3 Tracing
**Tool**: Jaeger

Track request flow across services

---

### Week 12: Security & Deployment

#### 12.1 Security Hardening
- [ ] Add HTTPS (Let's Encrypt)
- [ ] Implement CORS properly
- [ ] Add CSRF protection
- [ ] Sanitize all inputs
- [ ] Add SQL injection prevention
- [ ] Implement rate limiting
- [ ] Add DDoS protection

#### 12.2 Deployment
**Tool**: Docker + Kubernetes

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: bac-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: bac-api
  template:
    metadata:
      labels:
        app: bac-api
    spec:
      containers:
      - name: api
        image: bac-unified/api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
```

---

## BONUS FEATURES (If Time Permits)

### B1. Voice Assistant (TTS/STT)
**Tools**:
- Chatterbox TTS (multilingual, open-source)
- Whisper for STT

### B2. Knowledge Graph Visualization
**Tool**: Neo4j + D3.js

### B3. Peer Matching Algorithm
Match students studying same topics

### B4. Blockchain Certificates
Issue verifiable completion certificates

---

## OPEN-SOURCE TOOLS SUMMARY

| Category | Tool | Purpose |
|----------|------|---------|
| **OCR** | Surya | Layout + OCR (90+ langs) |
| | PaddleOCR | Arabic handwriting |
| | Pix2Tex | Image → LaTeX |
| **AI/LLM** | Qwen3-Coder | Code generation |
| | Nomic-Embed | Embeddings |
| | Ollama | Local LLM runtime |
| **Animation** | Noon | Rust animation engine |
| | Gifski | GIF generation |
| **Database** | PostgreSQL | Main database |
| | pgvector | Vector search |
| | Redis | Caching + pub/sub |
| **Web** | React + Vite | Frontend |
| | Workbox | Service worker |
| **Monitoring** | Prometheus | Metrics |
| | Grafana | Dashboards |
| | Loki | Logging |
| | Jaeger | Tracing |
| **Deployment** | Docker | Containers |
| | Kubernetes | Orchestration |

---

## SUCCESS METRICS

**Technical**:
- [ ] 100% test coverage for critical paths
- [ ] <200ms API response time (p95)
- [ ] 99.9% uptime
- [ ] <5% error rate
- [ ] 80%+ cache hit rate

**User**:
- [ ] 50k+ registered students
- [ ] 1M+ questions solved
- [ ] 100k+ animations generated
- [ ] 80%+ user satisfaction
- [ ] 70%+ exam pass rate improvement

---

## NEXT STEPS

1. **Review this plan** - Adjust priorities
2. **Set up project board** - Track progress
3. **Start Week 1** - Complete Playwright automation
4. **Daily standups** - 15min progress check
5. **Weekly demos** - Show working features

**Ready to start building?** 🚀
