# Implementation Tasks - 3Blue1Brown Animation System

**Approach**: Test-Driven Development (TDD)
**Prerequisite**: Run tests before implementation to verify requirements

---

## 1. Foundation & Environment Setup

- [x] 1.1 Create Python virtual environment `bac-animations` with requirements.txt
  - **Test**: `python -c "import manim; print(manim.__version__)"` returns version
  - **Test**: `python -c "import manim_physics; print('physics OK')"` imports successfully
  - **Test**: `python -c "import manim_chemistry; print('chemistry OK')"` imports successfully
  - **Test**: `ffmpeg -version` returns version
  - **Test**: `xvfb-run --auto-servernum manim --version` renders without display error

- [x] 1.2 Create bin/manim_renderer.py - Python wrapper for Manim rendering
  - **Test**: `python bin/manim_renderer.py --test` returns success
  - **Test**: Renders sample scene to output directory

- [x] 1.3 Add system dependencies documentation (Dockerfile/apt)
  - **Test**: All dependencies install without error on fresh Ubuntu 22.04

- [x] 1.4 Create Go-Python bridge in src/agent/internal/animation/bridge.go
  - **Test**: Go calls Python renderer and receives valid JSON response
  - **Test**: Timeout handling works (process killed after limit)
  - **Test**: Error propagation from Python to Go

---

## 2. Animation Engine (Core Rendering)

- [ ] 2.1 Implement renderer/renderer.go with quality presets
  - **Test**: All 5 quality presets produce correct resolution/fps
  - **Test**: Custom width/height/fps parameters work

- [ ] 2.2 Implement headless rendering with xvfb
  - **Test**: Renders successfully on server without display
  - **Test**: xvfb-run invoked correctly

- [ ] 2.3 Implement error handling for rendering failures
  - **Test**: Invalid scene code returns parse error
  - **Test**: LaTeX errors handled gracefully
  - **Test**: Out of memory errors handled with useful message

- [ ] 2.4 Create unit tests for renderer
  - **Test**: `go test -v ./internal/animation/renderer/...` passes

---

## 3. Template System

- [ ] 3.1 Create template manager in templates/manager.go
  - **Test**: Loads all YAML templates from templates/ directory
  - **Test**: Templates parsed into struct correctly

- [ ] 3.2 Implement template matching (keywords + subject)
  - **Test**: "graph f(x)" matches function_graph template
  - **Test**: "solve equation" matches equation_solving template
  - **Test**: Physics problems match physics templates
  - **Test**: Unknown problems fallback to generic template

- [ ] 3.3 Create math templates (5 core templates)
  - [ ] 3.3.1 function_graph.yaml - Graph plotting
    - **Test**: Template renders valid Manim code
    - **Test**: Generated code produces video
  - [ ] 3.3.2 equation_solving.yaml - Step-by-step solving
  - [ ] 3.3.3 geometry_proof.yaml - Geometric proofs
  - [ ] 3.3.4 calculus_derivative.yaml - Derivatives
  - [ ] 3.3.5 calculus_integral.yaml - Integrals

- [ ] 3.4 Create physics templates (3 core templates)
  - [ ] 3.4.1 mechanics_projectile.yaml - Projectile motion
  - [ ] 3.4.2 waves_simple.yaml - Wave animations
  - [ ] 3.4.3 electromagnetism_field.yaml - Electric fields

- [ ] 3.5 Create chemistry templates (2 core templates)
  - [ ] 3.5.1 molecule_2d.yaml - 2D molecular structures
  - [ ] 3.5.2 reaction_simple.yaml - Chemical reactions

- [ ] 3.6 Create unit tests for template system
  - **Test**: `go test -v ./internal/animation/templates/...` passes

---

## 4. AI Code Generator

- [ ] 4.1 Implement codegen/generator.go - main generator interface
  - **Test**: Generate returns valid Manim code
  - **Test**: Three generation levels work (template, ollama, hybrid)

- [ ] 4.2 Implement prompt builder for Ollama
  - **Test**: Math problems generate appropriate prompts
  - **Test**: Physics problems generate physics-specific prompts
  - **Test**: Chemistry problems generate chemistry prompts

- [ ] 4.3 Implement code validator
  - **Test**: Valid code passes validation
  - **Test**: Invalid code detected (missing imports, no Scene class)
  - **Test**: Template fallback triggered on validation failure

- [ ] 4.4 Implement variable extraction from problems
  - **Test**: Extracts function expressions from "f(x) = x^2"
  - **Test**: Extracts numeric values
  - **Test**: Extracts geometric parameters

- [ ] 4.5 Implement Ollama integration
  - **Test**: Ollama generates Manim code
  - **Test**: Timeout handling works
  - **Test**: Fallback to template on failure

- [ ] 4.6 Create integration tests for code generator
  - **Test**: Full pipeline: problem → code → render → video
  - **Test**: `go test -v ./internal/animation/codegen/...` passes

---

## 5. Rendering Pipeline

- [ ] 5.1 Implement pipeline/orchestrator.go
  - **Test**: Full pipeline executes in correct order
  - **Test**: Code gen → Render → Video → Cache

- [ ] 5.2 Implement queue system for async rendering
  - [ ] 5.2.1 Job queue with priority
    - **Test**: Higher priority jobs processed first
    - **Test**: Jobs wait when concurrency limit reached
  - [ ] 5.2.2 Job status tracking
    - **Test**: Status updates correctly (pending → processing → completed/failed)
    - **Test**: Progress percentage reported

- [ ] 5.3 Implement video processing with FFmpeg
  - [ ] 5.3.1 MP4 output
    - **Test**: Produces valid MP4 file
    - **Test**: Correct resolution and fps
  - [ ] 5.3.2 GIF output
    - **Test**: Produces optimized GIF
  - [ ] 5.3.3 Audio merging
    - **Test**: Audio syncs with video

- [ ] 5.4 Create pipeline tests
  - **Test**: `go test -v ./internal/animation/pipeline/...` passes

---

## 6. Caching System

- [ ] 6.1 Implement cache/key.go - cache key generation
  - **Test**: SHA256(problem+solution+quality+style) generates consistent keys
  - **Test**: Same inputs produce same key

- [ ] 6.2 Implement cache/memory.go - Redis cache layer
  - [ ] 6.2.1 Store and retrieve
    - **Test**: Cache stores with 24h TTL
    - **Test**: Cache returns hit within TTL
    - **Test**: Cache returns miss after TTL
  - [ ] 6.2.2 Cache lookup
    - **Test**: Checked on every render request

- [ ] 6.3 Implement cache/disk.go - disk cache layer
  - [ ] 6.3.1 File storage
    - **Test**: Videos stored in /data/animations/
    - **Test**: Organized by hash prefix
  - [ ] 6.3.2 Cleanup
    - **Test**: Deletes oldest when > 80% disk usage

- [ ] 6.4 Create cache tests
  - **Test**: Three-layer cache lookup works correctly
  - **Test**: `go test -v ./internal/animation/cache/...` passes

---

## 7. Voice Narration (Optional)

- [ ] 7.1 Implement voice/script_builder.go
  - [ ] 7.1.1 French script generation
    - **Test**: Generates French narration text
  - [ ] 7.1.2 English script generation
    - **Test**: Generates English narration text
  - [ ] 7.1.3 Timing markers
    - **Test**: Timing markers align with animation steps

- [ ] 7.2 Implement voice/synthesizer.go - Coqui TTS
  - **Test**: Synthesizes speech from script
  - **Test**: Handles TTS failures gracefully

- [ ] 7.3 Implement voice/mixer.go - audio-video sync
  - **Test**: Audio merged with video correctly

- [ ] 7.4 Create voice tests
  - **Test**: Voice feature can be disabled and system still works

---

## 8. REST API

- [ ] 8.1 Implement POST /animations/generate endpoint
  - **Test**: Returns job ID for async requests
  - **Test**: Waits and returns URL for sync requests

- [ ] 8.2 Implement GET /animations/:id endpoints
  - [ ] 8.2.1 Get animation metadata
    - **Test**: Returns video URL, duration, status
  - [ ] 8.2.2 Stream video
    - **Test**: Streams with correct content-type
  - [ ] 8.2.3 List animations
    - **Test**: Returns paginated list

- [ ] 8.3 Implement GET /templates endpoints
  - **Test**: List all templates
  - **Test**: Get single template details

- [ ] 8.4 Implement GET /queue/status endpoint
  - **Test**: Returns queue depth, processing count

- [ ] 8.5 Implement error handling
  - [ ] 8.5.1 400 for invalid requests
  - [ ] 8.5.2 404 for not found
  - [ ] 8.5.3 500 for internal errors

- [ ] 8.6 Create API tests
  - **Test**: All endpoints return correct status codes
  - **Test**: `go test -v ./api/...` passes

---

## 9. CLI Integration

- [ ] 9.1 Implement bac animate command
  - **Test**: `bac animate "problem text"` generates animation
  - **Test**: Options --subject, --quality work

- [ ] 9.2 Implement bac animations list command
  - **Test**: Displays table of animations

- [ ] 9.3 Implement bac templates command
  - **Test**: `bac templates list` shows templates
  - **Test**: `bac templates show <id>` shows details

- [ ] 9.4 Implement bac animate --audio
  - **Test**: Generates with voice narration

- [ ] 9.5 Create CLI tests
  - **Test**: All commands work from nushell

---

## 10. Database Schema

- [ ] 10.1 Add columns to animations table
  - **Test**: migrations apply without error

- [ ] 10.2 Create animation_templates table
  - **Test**: CRUD operations work

- [ ] 10.3 Create rendering_jobs table
  - **Test**: Job tracking works

- [ ] 10.4 Create narration_voices table
  - **Test**: Voice management works

- [ ] 10.5 Create indexes
  - **Test**: Queries perform well

---

## 11. Monitoring & Observability

- [ ] 11.1 Add Prometheus metrics
  - [ ] 11.1.1 Request counters
    - **Test**: /metrics shows animation_requests_total
  - [ ] 11.1.2 Duration histograms
    - **Test**: /metrics shows render duration
  - [ ] 11.1.3 Cache hit rate
    - **Test**: /metrics shows cache hit/miss

- [ ] 11.2 Add structured logging
  - **Test**: All operations log with appropriate level
  - **Test**: Log includes correlation IDs

---

## 12. End-to-End Testing

- [ ] 12.1 Integration test: Full pipeline
  - **Test**: Problem → API → Render → Video → Response

- [ ] 12.2 Integration test: With caching
  - **Test**: Second identical request returns cached

- [ ] 12.3 Integration test: With voice
  - **Test**: Animation includes audio track

- [ ] 12.4 Load test
  - **Test**: 10 concurrent renders work
  - **Test**: Queue handles 100 pending jobs

- [ ] 12.5 Full test suite passes
  - **Test**: `go test ./...` passes
  - **Test**: `python -m pytest tests/` passes

---

## Running Tests

```bash
# Go tests
cd src/agent && go test -v ./internal/animation/...

# Python tests
cd bin && python -m pytest tests/ -v

# Integration tests
./scripts/run_e2e_tests.sh

# All tests before commit
go test ./... && python -m pytest tests/
```
