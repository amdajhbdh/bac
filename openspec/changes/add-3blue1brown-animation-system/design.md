## Context

The BAC Unified platform needs a professional-grade animation system to generate 3Blue1Brown-style educational content. Currently, the animation module (`src/agent/internal/animation/animation.go`) references an empty Noon (Rust) placeholder that was never implemented. The system can only generate basic text-based fallback HTML when animation is attempted.

**Current Architecture:**
- Go-based Agent CLI in `src/agent/`
- Animation module generates Rust code templates (never compiled)
- No actual rendering capability exists
- Database has `animations` table but no real content

**Constraints:**
- Must integrate with Go-based agent (existing architecture)
- Should work on Linux server environment (headless rendering)
- GPU availability may be limited (optimize for CPU)
- Voice narration is optional (fallback if unavailable)
- Must maintain existing PostgreSQL database schema compatibility

## Goals / Non-Goals

**Goals:**
1. Implement complete animation rendering pipeline using Manim Community Edition
2. Generate professional mathematical visualizations (graphs, functions, geometry, calculus)
3. Add physics simulations (mechanics, waves, electromagnetism)
4. Add chemistry visualizations (molecules, reactions)
5. AI-powered code generation from problem solutions
6. Template-based generation for common problem types
7. Multi-quality rendering (preview to production)
8. Caching layer for performance optimization
9. Optional voice narration synthesis
10. REST API and CLI integration

**Non-Goals:**
- Real-time interactive animations (pre-rendered videos only)
- 3D animations beyond basic 3D surfaces (Manim's limitation)
- Mobile app integration (web-only initially)
- User-generated animations (admin/system-generated only)
- Animation editing/preview UI (render-only system)

## Decisions

### 1. Animation Engine: Manim Community Edition over alternatives

**Decision:** Use Manim CE (Community Edition) as primary animation engine.

**Alternatives Considered:**
- **Noon (Rust)**: Placeholder, never implemented, no ecosystem
- **Processing (Java)**: Not Python-based, limited math focus
- **p5.js (JavaScript)**: Web-focused, not professional quality
- **Blender (Python API)**: Overkill for 2D educational content
- **ManimGL**: Legacy version, less maintained than CE
- **ManimVTK**: Good for scientific viz but heavier dependencies

**Rationale:** Manim is the industry standard for educational math content (used by 3Blue1Brown), has active community, excellent documentation, Python-based (fits our stack), and rich plugin ecosystem (physics, chemistry).

### 2. AI Code Generation: Hybrid Template + Ollama

**Decision:** Use three-tier code generation: template-first, Ollama-enhanced, fallback-to-template.

**Alternatives Considered:**
- **Pure Template**: Limited flexibility, can't handle novel problems
- **Pure LLM**: Unreliable output, may generate broken code
- **Fine-tuned Model**: Requires training data, overkill for initial version

**Rationale:** Templates provide reliability for common problems. Ollama adds flexibility for novel cases. Fallback ensures system always produces output.

### 3. Go-Python Integration: Subprocess Bridge

**Decision:** Use Go subprocess to call Python wrapper script.

**Alternatives Considered:**
- **CGO**: Complex, memory management issues
- **gRPC**: Over-engineered for this use case
- **Shared memory**: Platform-specific, risky

**Rationale:** Simple, reliable, easy to debug, good timeout handling, matches existing patterns in codebase.

### 4. Rendering: xvfb for Headless Operation

**Decision:** Use xvfb (virtual framebuffer) for headless rendering on servers.

**Alternatives Considered:**
- **Docker with GPU**: Requires GPU hardware
- **Cloud rendering API**: Costly, external dependency
- **EGL/Mesa**: More complex setup

**Rationale:** Works on any Linux server, no GPU required, well-tested with Manim.

### 5. Voice: Optional Coqui TTS with Fallback

**Decision:** Implement TTS but make it optional - animations work without audio.

**Alternatives Considered:**
- **Required TTS**: Adds complexity, may fail
- **Google Cloud TTS**: External API, cost, privacy
- **ElevenLabs**: Paid API, not open source

**Rationale:** Coqui TTS is open-source, runs locally, good French support. System remains functional without it.

### 6. Caching: Three-Layer Strategy

**Decision:** In-memory (Redis) → Disk → Database.

**Alternatives Considered:**
- **Database only**: Too slow for rendering
- **Disk only**: No distributed caching
- **CDN**: External dependency, cost

**Rationale:** Matches existing BAC architecture (Redis + disk), fast lookups for common problems.

## Risks / Trade-offs

### Risk: Manim Rendering Time → Mitigation: Quality Presets + Async Processing
Manim rendering can take minutes for high-quality output. We implement quality presets (preview/low/medium/high/production), async processing with job queue, and aggressive caching.

### Risk: AI Code Generation Failures → Mitigation: Template Fallback + Validation
Ollama may generate invalid Manim code. We validate generated code before rendering, fall back to templates on failure, and log failures for improvement.

### Risk: Dependency Hell → Mitigation: Virtual Environment + Docker
Manim has complex dependencies (LaTeX, FFmpeg, OpenGL). We use dedicated Python venv, document system deps, and provide Dockerfile.

### Risk: Disk Space Exhaustion → Mitigation: Tiered Storage + Cleanup
Video files are large. We implement disk usage monitoring, automatic cleanup of old renders, and tiered storage (SSD for recent, HDD for archive).

### Risk: LaTeX Rendering Failures → Mitigation: Fallback to No-LaTeX
LaTeX may not be installed or may fail on complex equations. We catch errors and render without LaTeX as fallback.

### Risk: Memory Exhaustion on Large Animations → Mitigation: Resource Limits + Queue
Complex scenes may use too much RAM. We limit concurrent renders, use swap space, and monitor memory usage.

## Migration Plan

**Phase 1: Foundation (Weeks 1-4)**
1. Create Python virtual environment with Manim dependencies
2. Implement Go-Python bridge (subprocess wrapper)
3. Create basic Manim scene rendering
4. Test headless rendering with xvfb
5. Deploy to staging

**Phase 2: Template System (Weeks 5-10)**
1. Implement YAML-based template system
2. Create math templates (function graphs, equations, geometry)
3. Create physics templates (mechanics, waves, EM)
4. Create chemistry templates (molecules, reactions)
5. Test template matching and code generation

**Phase 3: AI Integration (Weeks 11-18)**
1. Implement Ollama code generation prompts
2. Add validation for generated code
3. Implement hybrid generation (template + AI)
4. Add retry logic and error handling

**Phase 4: Video Pipeline (Weeks 19-24)**
1. Implement quality presets
2. Add FFmpeg video processing
3. Implement rendering queue
4. Add caching layer

**Phase 5: Voice (Optional, Weeks 25-28)**
1. Add Coqui TTS integration
2. Implement narration script builder
3. Add audio-video sync
4. Make voice optional in API

**Phase 6: API + CLI (Weeks 29-32)**
1. Implement REST endpoints
2. Add CLI commands
3. Integrate with existing auth
4. Deploy to production

## Open Questions

1. **GPU Rendering**: Should we invest in GPU infrastructure for faster rendering? Current CPU-only approach is slower but simpler.

2. **Concurrent Rendering**: What's the maximum concurrent render jobs we should allow? Limited by RAM/CPU.

3. **Storage Strategy**: Should animation storage use Garage S3 (existing) or local disk? Trade-off: S3 is distributed but slower.

4. **Template Discovery**: How should we identify which template matches a problem? Current: keyword matching. Future: AI classification.

5. **Analytics**: What metrics should we track beyond basic usage? Consider: render time, cache hit rate, template usage, error rates.
