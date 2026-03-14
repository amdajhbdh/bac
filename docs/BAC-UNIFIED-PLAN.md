# BAC Unified - National Exam Preparation System

## Complete Technical Blueprint v2.1

---

## Executive Summary

| Aspect | Specification |
|--------|---------------|
| **Project Name** | BAC Unified - المنظومة الوطنية للباكالوريا |
| **Owner** | Ministry of Education, Islamic Republic of Mauritania |
| **Timeline** | 12 months (4 quarters) |
| **Target Users** | 50,000+ BAC C students annually |
| **Infrastructure** | On-premise Ministry servers + MinIO |
| **AI Approach** | Fine-tuned local models + cloud fallback |
| **Security** | AES-256 encryption + Blockchain audit trail |
| **Offline** | Partial - core features offline, advanced online |
| **Incentives** | Points + Badges + Leaderboards |
| **Animation Engine** | **Noon** (Rust-based, faster than manim) |

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              USER INTERFACES                                     │
├─────────────────────┬─────────────────────┬─────────────────────────────────────┤
│      CLI (Nu)       │   Web PWA (React)   │        Telegram Bot                │
└──────────┬──────────┴──────────┬──────────┴──────────────┬──────────────────────┘
           │                      │                        │
           └──────────────────────┼────────────────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │      API GATEWAY (Go)     │
                    │  • JWT Auth               │
                    │  • Rate Limiting          │
                    │  • Request Caching        │
                    │  • Load Balancing         │
                    │  • SSL/TLS Termination    │
                    └─────────────┬─────────────┘
                                  │
         ┌────────────────────────┼────────────────────────┐
         │                        │                        │
         ▼                        ▼                        ▼
┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
│  SUBMISSION     │   │   SOLVER        │   │  PREDICTION     │
│   MICROSERVICE  │   │   MICROSERVICE  │   │  MICROSERVICE   │
│                 │   │                 │   │                 │
│ • OCR           │   │ • LLM Engine   │   │ • ML Pipeline   │
│ • PDF Parser    │   │ • Step Solver  │   │ • Pattern Miner │
│ • Image Process │   │ • Noon Gen     │   │ • Risk Analyzer │
│ • LaTeX Conv   │   │ • Answer Gen   │   │ • Trend Detector│
└────────┬────────┘   └────────┬────────┘   └────────┬────────┘
         │                     │                     │
         └─────────────────────┼─────────────────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
              ▼                ▼                ▼
       ┌────────────┐   ┌───────────┐   ┌───────────┐
       │ PostgreSQL │   │  MinIO    │   │  Redis    │
       │ +pgvector │   │ (S3 API)  │   │ (Cache)   │
       └────────────┘   └───────────┘   └───────────┘
              │
              ▼
       ┌────────────┐
       │  Blockchain│
       │  (Audit)   │
       └────────────┘
```

---

## Media-to-Noon Animation Pipeline

The core innovation: **AI analyzes submitted media and generates Noon animations automatically**

### 2.1 Input Processing Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MEDIA-TO-ANIMATION PIPELINE                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  USER SUBMISSION                                                            │
│  ──────────────                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐                    │
│  │  Image   │  │   PDF    │  │   URL    │  │  Text    │                    │
│  │(handwr.) │  │(exerc.)  │  │(video)   │  │(problem) │                    │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘                    │
│       │             │             │             │                             │
│       └─────────────┼─────────────┼─────────────┘                             │
│                     ▼                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    OCR LAYER (Multi-Strategy)                       │   │
│  │                                                                      │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │   │
│  │  │  Tesseract   │  │   EasyOCR    │  │  PaddleOCR  │            │   │
│  │  │  (Primary)   │  │  (Fallback)  │  │ (Advanced)  │            │   │
│  │  │              │  │              │  │              │            │   │
│  │  │ + Arabic    │  │ + Math       │  │ + Tables    │            │   │
│  │  │ + French    │  │   detection  │  │              │            │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘            │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                     │                                                       │
│                     ▼                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    AI ANALYSIS LAYER                                  │   │
│  │                                                                      │   │
│  │  ┌──────────────────────────────────────────────────────────────┐  │   │
│  │  │                    QUESTION UNDERSTANDING                       │  │   │
│  │  │                                                               │  │   │
│  │  │  Input: "Solve: ∫(x² + 2x)dx from 0 to 2"                     │  │   │
│  │  │                                                               │  │   │
│  │  │  Output:                                                      │  │   │
│  │  │  {                                                            │  │   │
│  │  │    "subject": "mathematics",                                 │  │   │
│  │  │    "topic": "integrals",                                     │  │   │
│  │  │    "type": "calculation",                                    │  │   │
│  │  │    "difficulty": 3,                                          │  │   │
│  │  │    "concepts": ["antiderivative", "definite integral"],     │  │   │
│  │  │    "visual_elements": ["graph", "area_shading"],            │  │   │
│  │  │    "solution_steps": [                                        │  │   │
│  │  │      "expand: x² + 2x",                                      │  │   │
│  │  │      "antiderivative: x³/3 + x²",                            │  │   │
│  │  │      "apply bounds: F(2) - F(0)",                            │  │   │
│  │  │      "calculate: 8/3 + 4 - 0"                                 │  │   │
│  │  │    ]                                                          │  │   │
│  │  │  }                                                            │  │   │
│  │  └──────────────────────────────────────────────────────────────┘  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                     │                                                       │
│                     ▼                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    SOLVER ENGINE                                     │   │
│  │                                                                      │   │
│  │  ┌───────────────────────┐  ┌───────────────────────┐              │   │
│  │  │    LOCAL OLLAMA       │  │   CLOUD FALLBACK     │              │   │
│  │  │  (Fine-tuned BAC)     │  │  (kimi/glm/minimax)  │              │   │
│  │  └───────────────────────┘  └───────────────────────┘              │   │
│  │                                                                      │   │
│  │  Output: Step-by-step solution + all possible answers              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                     │                                                       │
│                     ▼                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    NOON ANIMATION GENERATOR                          │   │
│  │                                                                      │   │
│  │  ┌──────────────────────────────────────────────────────────────┐  │   │
│  │  │              AI-GENERATED NOON CODE                            │  │   │
│  │  │                                                               │  │   │
│  │  │  // Auto-generated from solution steps                        │  │   │
│  │  │  fn scene(win_rect: Rect) -> Scene {                         │  │   │
│  │  │      let mut scene = Scene::new(win_rect);                  │  │   │
│  │  │                                                               │  │   │
│  │  │      // Step 1: Draw axes                                     │  │   │
│  │  │      let axes = scene.axes().with_range(-1, 3).make();      │  │   │
│  │  │      scene.play(axes.show_creation()).run_time(1.0);        │  │   │
│  │  │                                                               │  │   │
│  │  │      // Step 2: Plot function                                 │  │   │
│  │  │      let graph = axes.plot(|x| x*x + 2.0*x).make();         │  │   │
│  │  │      scene.play(graph.show_creation()).run_time(1.5);       │  │   │
│  │  │                                                               │  │   │
│  │  │      // Step 3: Shade area under curve                       │  │   │
│  │  │      let area = axes.fill_between(graph, 0.0, 2.0);          │  │   │
│  │  │      scene.play(area.fade_in()).run_time(1.0);               │  │   │
│  │  │                                                               │  │   │
│  │  │      // Step 4: Show calculation                              │  │   │
│  │  │      let text = scene.text("∫ = 20/3").make();              │  │   │
│  │  │      scene.play(text.show_creation()).run_time(1.0);        │  │   │
│  │  │  }                                                            │  │   │
│  │  └──────────────────────────────────────────────────────────────┘  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                     │                                                       │
│                     ▼                                                       │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    OUTPUT                                             │   │
│  │                                                                      │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐             │   │
│  │  │   MP4 Video  │  │  Thumbnail   │  │   Metadata   │             │   │
│  │  │  (Animation) │  │   (Preview)  │  │   (JSON)     │             │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘             │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Noon Animation Types

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    NOON ANIMATION TYPES BY SUBJECT                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    MATHEMATICS                                        │   │
│  │                                                                      │   │
│  │  Graph Animations              Geometry Animations                   │   │
│  │  ─────────────────            ──────────────────                    │   │
│  │  • Function plotting          • Triangle constructions              │   │
│  │  • Curve transformations     • Circle theorems                     │   │
│  │  • Area under curves         • Geometric proofs                    │   │
│  │  • 3D surfaces               • Transformations                     │   │
│  │                                                                      │   │
│  │  Calculus                    Algebra                               │   │
│  │  ────────                    ────────                               │   │
│  • Derivative visualization    • Matrix operations                     │   │
│  • Integral approximation      • Vector operations                     │   │
│  • Limit processes             • Complex numbers                       │   │
│  • Series convergence          • Polynomial roots                      │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    PHYSICS                                            │   │
│  │                                                                      │   │
│  │  Mechanics                   Electricity                             │   │
│  │  ──────────                 ──────────                               │   │
│  │  • Projectile motion        • Circuit diagrams                      │   │
│  │  • Pendulum oscillation    • Resistor networks                     │   │
│  │  • Spring oscillations     • Capacitor charge/discharge             │   │
│  • Collisions                 • Electromagnetic induction              │   │
│  • Circular motion            • Wave propagation                      │   │
│  │                                                                      │   │
│  │  Optics                     Thermodynamics                          │   │
│  │  ──────                    ────────────                             │   │
│  │  • Ray diagrams            • Heat transfer                         │   │
│  • Lens combinations          • Gas processes                          │   │
│  • Interference patterns      • Entropy visualization                  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    CHEMISTRY                                          │   │
│  │                                                                      │   │
│  │  Organic Chemistry          Reactions                                │   │
│  │  ───────────────           ─────────                                │   │
│  │  • Molecular structures    • Reaction mechanisms                    │   │
│  │  • Isomer transformations  • Bond breaking/forming                  │   │
│  │  • Reaction pathways       • Equilibrium shifts                     │   │
│  │                                                                      │   │
│  │  Structure                  Kinetics                                 │   │
│  │  ────────                  ────────                                 │   │
│  │  • Electron orbitals       • Collision theory                       │   │
│  │  • Crystal structures      • Rate curves                           │   │
│  │  • Chemical bonds           • Activation energy                      │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Implementation Phases (12 Months)

### Phase 1: Foundation (Months 1-2)
- [ ] Infrastructure setup (MinIO, PostgreSQL)
- [ ] Database schema design
- [ ] Authentication system
- [ ] Core API development
- [ ] Clone and setup noon repository

### Phase 2: Input System (Months 2-3)
- [ ] Multi-OCR pipeline
- [ ] PDF parser with LaTeX extraction
- [ ] Image preprocessing
- [ ] AI question analyzer

### Phase 3: Solver Engine (Months 3-5)
- [ ] LLM integration
- [ ] Step-by-step solution generator
- [ ] Multi-answer generator
- [ ] Solution validation

### Phase 4: Noon Animation (Months 5-6)
- [ ] Install noon from https://github.com/yongkyuns/noon
- [ ] AI-to-Noon code generator
- [ ] Math animations
- [ ] Physics/chemistry animations

### Phase 5: Prediction System (Months 6-7)
- [ ] Historical data import
- [ ] ML pattern analysis
- [ ] Prediction API

### Phase 6: User Features (Months 7-9)
- [ ] Progress tracking
- [ ] Points/badges system
- [ ] Leaderboards
- [ ] Gamification

### Phase 7: Testing & Launch (Months 9-12)
- [ ] Beta testing
- [ ] Ministry approval
- [ ] Public launch

---

## Technology Stack

| Component | Technology |
|-----------|------------|
| **API Server** | Go (Gin) |
| **Database** | PostgreSQL + pgvector |
| **Storage** | MinIO (local S3) |
| **AI/LLM** | Ollama (local) + Cloud fallback |
| **OCR** | Tesseract, EasyOCR, PaddleOCR |
| **Animation** | **Noon** (Rust) |
| **CLI** | Nushell |
| **Auth** | JWT |
| **Blockchain** | Hyperledger Fabric |
| **Deployment** | Docker + Kubernetes |

---

## CLI Commands (Year 1)

```bash
# Submit content
bac submit image.jpg --subject physics
bac submit pdf exercises.pdf --chapter "optics"
bac submit url https://bacmr.net/exam2024

# Solve with animation
bac solve "∫(x² + 2x)dx from 0 to 2"
bac solve image.jpg --animate

# Practice
bac practice --subject math --difficulty hard
bac practice --chapter integrals

# Predictions (PUBLIC)
bac predict --subject math
bac predict --year 2027

# Progress & Rewards
bac stats
bac leaderboard
bac badges
```

---

## Documentation Standards

All AI actions must be documented:

1. **API Logs**: Every request logged with timestamp, user, action
2. **Solution Audit**: All solver outputs stored with version
3. **Animation Metadata**: JSON with generation parameters
4. **User Actions**: Progress, submissions, verifications
5. **Blockchain Trail**: Immutable audit log for all changes

---

## Success Metrics

| Metric | Year 1 Target |
|--------|---------------|
| Total Students | 50,000+ |
| Active Users | 30,000+ |
| Questions Bank | 100,000+ |
| Questions Solved | 500,000+ |
| Animations Generated | 100,000+ |
| Predictions Published | 200+ |
| Average Accuracy | 75%+ |
| System Uptime | 99.5% |

---

*Generated: 2026-02-23*
*Version: 2.1*
*Status: PLANNING*
