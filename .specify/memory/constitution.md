# BAC Unified Constitution

## Core Principles

### I. Learning-First
Every feature must directly improve learning outcomes for Moroccan BAC students. Features that do not contribute to effective study, practice, or understanding are not implemented. The system prioritizes educational value over novelty.

### II. Offline-First
The system MUST function fully offline for core study features. Students in areas with limited connectivity must be able to practice problems, review notes, and use spaced repetition. Cloud services are optional enhancements, not requirements.

### III. Local AI Priority
AI inference follows strict priority: (1) Ollama local models, (2) NotebookLM with caching, (3) Cloud APIs as fallback. Paid AI services are never required. The system maximizes local processing to reduce costs and ensure availability.

### IV. Open Educational Resources
All bundled study materials must be freely available or properly licensed. No proprietary paid content. Content aggregation sources include Khan Academy, OpenStax, MIT OCW, and other OER providers. Attribution must be maintained.

### V. Privacy-Respecting
Student study data, progress, and learning patterns remain local unless explicitly shared by the student. No telemetry that could identify individual学生的学习 habits. All data stored on user's device or self-hosted infrastructure.

### VI. Vault-First (Obsidian as Primary Interface)
The study vault is the central interface for all learning content. The vault MUST be:
- A standalone git repository accessible from phone and desktop
- Structured like professional textbooks (Harvard-style)
- Synced via git (not Obsidian Sync paid service)
- Self-contained with all content in markdown format

### VII. AI-Managed OCR Pipeline
All OCR operations MUST be managed by AI to ensure quality and structure. Raw Tesseract output is not acceptable - it must be processed through:
- **Ollama (llama3.2:3b)** for error correction and formatting
- Structure extraction: convert raw OCR to proper markdown with headings, lists, formulas
- Wikilink generation: connect related concepts automatically
- The pipeline handles Arabic/French character correction, mathematical notation, and layout understanding

### VIII. OCR Pipeline for All Resources
All PDF resources (teacher handwritten lessons, scanned books, exercise sheets) MUST be OCR'd and converted to searchable markdown. The pipeline:
- Uses Tesseract with Arabic + French + English support
- Preserves original PDF sources alongside extracted text
- Adds proper front matter and wikilinks for cross-referencing

### IX. Multi-Cloud Infrastructure
The system is distributed across multiple cloud providers for resilience and performance:
- **CloudShell**: Development environment, OCR processing
- **Cloudflare Workers**: API endpoints, agent coordination
- **Render**: Go/Rust service hosting
- **Fly.io**: Edge computing, low-latency requests
- **Doppler**: Secrets management and environment configuration

## Technology Constraints

**Stack Requirements**: Go 1.25+ for CLI/API, Rust for performance-critical services, React for web UI, Nushell for scripting. PostgreSQL + pgvector for main database, Redis for caching.

**Vault Stack**: Obsidian for note-taking, Tesseract for OCR, Git for sync.

**Prohibited**: Paid OCR APIs (must use Tesseract/Surya), proprietary storage (must use S3-compatible), cloud-only AI (must have local fallback), Obsidian Sync (use git-based sync instead).

**Quality Gates**: All Go code passes golangci-lint, all Rust code passes clippy, all TypeScript passes ESLint + typecheck.

## Development Workflow

**Phased Implementation**: Features proceed through Setup → Foundational → User Stories → Polish. Foundational phase must complete before any user story work begins.

**Test Discipline**: User stories include independent test scenarios. Tests are written before implementation to define acceptance criteria. Each story must be independently testable and deployable.

**Constitution Compliance**: All PRs must verify compliance with these principles. Complexity must be justified against learning value. Violations require explicit approval with documented rationale.

## Governance

This constitution supersedes all other development practices. Amendments require: (1) documented proposal explaining the change, (2) rationale tied to project purpose, (3) migration plan if applicable, (4) update to AGENTS.md with new commands if relevant.

The constitution is reviewed quarterly to ensure alignment with the personal learning system mission.

**Version**: 1.2.0 | **Ratified**: 2026-03-11 | **Last Amended**: 2026-03-14
