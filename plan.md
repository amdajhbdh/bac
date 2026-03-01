<!--toc:start-->

- [1. North‑Star Vision: BAC Hyper‑Tutor 2026+](#1-northstar-vision-bac-hypertutor-2026)
- [2. System Pillars](#2-system-pillars)
- [3. Multi‑Phase Roadmap](#3-multiphase-roadmap)
  - [Phase 1 – Solidify the Core CLI + DB Brain (0–3 months)](#phase-1-solidify-the-core-cli-db-brain-03-months)
  - [Phase 2 – True Hyper‑Tutor Pipeline (3–9 months)](#phase-2-true-hypertutor-pipeline-39-months)
  - [Phase 3 – Automation Swarm & Continuous Ingestion (6–18 months)](#phase-3-automation-swarm-continuous-ingestion-618-months)
  - [Phase 4 – Noon‑First Visual Tutor (9–24 months)](#phase-4-noonfirst-visual-tutor-924-months)
  - [Phase 5 – Prediction Engine & National‑Scale Features (12–36 months)](#phase-5-prediction-engine-nationalscale-features-1236-months)
  - [Phase 6 – Community, Ecosystem & Governance (18–48 months)](#phase-6-community-ecosystem-governance-1848-months)
- [4. How AGENTS & Daemons Fit In](#4-how-agents-daemons-fit-in)
- [5. Deep Dives You’ll “Dive In” To](#5-deep-dives-youll-dive-in-to)
<!--toc:end-->

Here’s a long‑term, “big arc” plan focused on _upgrading_ BAC Unified into the hyper‑tutor you want, grounded in your current repo (Nushell CLI, Python solver, Noon/Rust, Postgres, MinIO).

I’ll structure it like this:

1. North‑star vision (what “done” feels like)
2. System pillars
3. Multi‑year roadmap (phases)
4. How AGENTS/daemons fit in
5. “Deep dives” you’ll do as you go

You can keep this as `docs/MASTER_PLAN.md`.

---

## 1. North‑Star Vision: BAC Hyper‑Tutor 2026+

In 1–2 years, BAC Unified should feel like:

- A **single command** (`bac`) that orchestrates:
  - Solving, explaining, visualizing, and tracking every BAC C exercise,
  - Across text, images, PDFs, and web resources.
- A **national‑level question bank**:
  - All past exams, regional mocks, teacher sheets, textbook questions,
  - Cleanly structured, searchable by concept, difficulty, and year.
- A **hyper‑tutor** that:
  - Understands the student’s _history_ (weak/strong concepts),
  - Proposes adaptive study plans,
  - Generates step‑by‑step solutions, plus **Noon animations**,
  - Works fully offline with local models, and “upgrades” itself online when allowed.
- A **safe automation swarm**:
  - Crawlers, OCR/ingestion daemons, research agents, prediction jobs,
  - All sandboxed and observable, with human‑in‑the‑loop for risky actions.
- An **audited, explainable system**:
  - Every solution and prediction traceable,
  - Blockchain/audit trail integrated but doesn’t slow development.

---

## 2. System Pillars

These are the “axes” you’ll continuously improve:

1. **Knowledge Base**
   - Postgres + pgvector as the _source of truth_.
   - Tables: problems, concepts, solution_steps, media (Noon animations), tool_runs, students, predictions, audits.

2. **Orchestration & Interface**
   - `bac.nu` as the _only_ student‑facing entrypoint.
   - Internal modules (`lib/bac.nu`, Python, Rust) as “tools” the orchestrator uses.

3. **Solving & Teaching**
   - Pipelines:
     - Parse → Retrieve → Solve → Teach → Visualize → Log.
   - Subject‑specific solvers for Math/Physics/Chemistry/French.

4. **Automation & Agents**
   - Daemons + agents (crawlers, OCR, YouTube, research) that:
     - Ingest new problems,
     - Enrich concepts,
     - Generate hints and variants.

5. **Visualization & Experience**
   - Noon/Rust as first‑class part of solving.
   - `--animate` is not a gimmick but a normal expectation.

6. **Safety & Observability**
   - Clear safety envelopes for agents.
   - Logs, metrics, and simple dashboards.

---

## 3. Multi‑Phase Roadmap

### Phase 1 – Solidify the Core CLI + DB Brain (0–3 months)

**Goal:** Make `bac` a rock‑solid orchestrator that writes everything to the DB and can grow safely.

**Key upgrades:**

1. **Unify the CLI contract**
   - Finalize public commands:
     - `bac solve`, `bac submit`, `bac practice`, `bac predict`, `bac stats`, `bac audit`, `bac storage`.
   - Introduce a “student‑first” alias layer:
     - `bac ask`, `bac study`, `bac exam-sim` that internally map to existing commands.

2. **Database as central brain**
   - Expand schema beyond the current `save-to-database` usage:
     - `problems`: text, subject, parsed metadata.
     - `solutions`: linked to problems, solver name, quality flags.
     - `concepts`: canonical list, descriptions, embeddings.
     - `student_profiles`: per‑student baseline info.
     - `student_concepts`: mastery per concept.
     - `tool_runs`: any Python/agent/Noon invocation.
   - Make all major CLI actions go through the DB:
     - `solve` saves problem + solution,
     - `practice` pulls questions from DB,
     - `predict` uses DB history.

3. **Structured parsing pipeline**
   - Add a “ProblemParser” stage before `python-solve`:
     - Extract subject, chapter, year (if exam), estimated difficulty, concept tags.
   - Initially rule‑based → then LLM‑based.

4. **Core logging & basic safety**
   - `docs/SAFETY.md` (as discussed previously).
   - `docs/TOOLS.md` listing:
     - All tools (Python solver, Noon, daemons, future agents),
     - Safety constraints,
     - A short usage log.
   - In `bac.nu`, add:
     - `--debug` or `BAC_DEBUG` env to trace each stage (parse/retrieve/solve/animate).

**Outcome:** A robust, transparent CLI + DB spine that other pieces plug into.

---

### Phase 2 – True Hyper‑Tutor Pipeline (3–9 months)

**Goal:** Turn `bac solve` into a real hyper‑tutor: structured reasoning, RAG, teaching, and early personalization.

1. **Deep Problem Representation**
   - After parsing, store a **normalized problem representation**:
     - “Given f(x) = …, compute f’(2)” as structured math, not just text.
   - Start building a library of **problem templates**:
     - Parameterized exercises for automatic variant generation.

2. **Concept Universe**
   - Build out a `concepts` universe for BAC C:
     - For each concept: name, short description, examples, common mistakes.
   - Link each `problem` to multiple `concepts` (many‑to‑many).

3. **Retrieval‑Augmented Solving**
   - Before calling `python-solve`:
     - Retrieve similar problems from `problems`+`problem_embeddings`,
     - Retrieve concept notes.
   - Pass this context into the solver:
     - Either via local LLM prompt,
     - Or as structured hints the Python solver understands.

4. **Teacher Layer**
   - After getting a solution:
     - Add a “TeacherAgent” pass that:
       - Breaks solution into atomic steps,
       - Adds hints after each step,
       - Generates a short recap in simple French/Arabic.
   - Integrate into CLI:
     - `bac solve "..." --mode teacher`
     - Or default behavior if the user is a student.

5. **Early personalization**
   - When a student interacts:
     - Link actions to a `student_id`.
   - Track:
     - Which concepts appear in solved problems,
     - Performance on quizzes/practice,
     - Generate a **weakness map**.
   - Enable:
     - `bac study` to propose a 20–30 minute adaptive session.

**Outcome:** For any exercise, you get: parse → context → stepwise solution → explanation suited to the student.

---

### Phase 3 – Automation Swarm & Continuous Ingestion (6–18 months)

**Goal:** The system continuously discovers, ingests, and curates new content with agents/daemons.

1. **Ingestion workflows**
   - Extend existing daemons (`daemons/`) into end‑to‑end flows:
     - **PDF pipeline**:
       - Watch `db/pdfs/` → OCR (if needed) → split into problems → parse → store.
     - **YouTube pipeline**:
       - Watch configured channels → download edu videos → transcribe → detect problems within transcripts.
     - **Web crawler pipeline**:
       - Crawl education sites, past exam repositories, GitHub repos.

2. **Agent patterns**
   - For each pipeline, define a clear agent role:
     - _Crawler Agent_: finds raw content (URLs, PDFs).
     - _Extractor Agent_: turns raw content into structured problems.
     - _Annotator Agent_: tags concepts, difficulty, adds metadata.
   - All agents:
     - Run under restricted OS user,
     - Only touch `resources/`, `db/pdfs/`, and scratch dirs,
     - Log every run to `tool_runs`.

3. **Research agents**
   - For missing explanations or unclear topics:
     - **ResearchAgent** that:
       - Uses online LLMs via Playwright/MCP,
       - Searches NotebookLM notebooks if you continue to use them,
       - Produces a “summary for DB” that is:
         - Short,
         - Attributed,
         - Stored into `concepts` or `notes`.

4. **Scheduling & monitoring**
   - Use `systemd`, `cron`, or a workflow tool (n8n) to:
     - Run ingestion jobs daily,
     - Run prediction/analysis weekly.
   - Add basic observability:
     - Log volumes per day,
     - Number of new problems,
     - Failures per pipeline.

**Outcome:** The knowledge base grows continuously with minimal manual work, still under your control and auditable.

---

### Phase 4 – Noon‑First Visual Tutor (9–24 months)

**Goal:** Visual explanations become a core part of the experience, not an add‑on.

1. **Canonical solution representation**
   - Extend solver output format to include:
     - Entities (functions, intervals, points, lines, vectors),
     - Operations (differentiate, integrate, solve equation, transform),
     - Geometry objects (triangles, circles, transformations).
   - Use this to generate both:
     - Human‑readable steps,
     - Machine‑readable instructions for Noon.

2. **Noon integration**
   - Build a dedicated **Noon submodule** (Rust project) that:
     - Receives a JSON description of solution,
     - Generates a `.mp4` or `.webm` animation,
     - Stores a reference in `media` table (linked to problem/solution).
   - Add CLI usage:
     - `bac solve "..." -a` → triggers Noon,
     - `bac animate <problem-id>` → re‑renders or replays.

3. **Template scenes**
   - Create a library of Noon templates by concept:
     - Derivative as slope,
     - Area under curve,
     - Trig unit circle,
     - Motion in 1D/2D for physics.
   - The solver chooses template IDs and fills parameters.

4. **Interactive flows (later)**
   - Plan for future:
     - `bac study` sessions that include:
       - Watch animation,
       - Answer a question about it,
       - Get “What if…?” interactive follow‑ups.

**Outcome:** Visual intuition is baked into the system; many problems automatically get high‑quality animations.

---

### Phase 5 – Prediction Engine & National‑Scale Features (12–36 months)

**Goal:** Use the accumulated data to predict and simulate exams.

1. **Prediction engine**
   - Use your DB of past exams:
     - For each year, subject, topic:
       - Frequency counts,
       - Difficulty distributions.
   - Simple early models:
     - Weighted frequency + recency to predict high‑probability topics.
   - Later models:
     - Sequence models over years,
     - Multi‑factor models using curriculum changes, teacher feedback.

2. **Exam simulation**
   - `bac exam-sim`:
     - Generates full mock exams with:
       - Realistic question mix,
       - Appropriate difficulty curve,
       - Coverage across major concepts.
   - After completion:
     - Score student,
     - Update concept mastery,
     - Give targeted recommendations.

3. **Nation‑scale readiness**
   - Support multiple schools/classes:
     - `student_profiles` with school/class metadata.
   - Dashboard ideas (for later when you have users):
     - Aggregate stats by school, anonymized.

**Outcome:** BAC Unified shifts from “solver” to **co‑pilot for exam strategy**.

---

### Phase 6 – Community, Ecosystem & Governance (18–48 months)

**Goal:** Make BAC Unified sustainable, extensible, and community‑driven without compromising safety.

1. **Plugin/skill system**
   - Inspired by OpenClaw “skills”:
     - Define a formal contract for new tools:
       - Input/Output schema,
       - Permissions,
       - Logging requirements.
   - Allow advanced users to add:
     - Custom solvers,
     - Additional visualization tools,
     - New ingestion workflows.

2. **Governance & content QA**
   - Add moderation/verification fields:
     - `solutions` and `concepts` have:
       - `source` (human/agent),
       - `verified_by` (teacher id),
       - `status` (draft/verified/flagged).
   - Allow teachers to:
     - Approve or correct solutions,
     - Contribute problems and notes.

3. **Documentation & training**
   - Expand `AGENTS.md` into a full **Agent Playbook**:
     - How to write new agents safely,
     - Example flows,
     - Safety pitfalls.
   - Keep `README.md` student‑simple; move dev complexity into `docs/`.

**Outcome:** A mature project that others can extend safely.

---

## 4. How AGENTS & Daemons Fit In

Given your current `AGENTS.md` and daemons directory, the long‑term pattern is:

- **bac.nu = conductor**, AGENTS/daemons = instruments.
- Each daemon/agent:
  - Has:
    - Single responsibility (OCR pipeline, YouTube ingestion, web crawl, research),
    - Clear inputs/outputs,
    - Config in `config/`,
    - Logs to `tool_runs`.
  - Never “owns” the student experience; only `bac` does.

You evolve `AGENTS.md` into:

- Role descriptions (Crawler, Extractor, Annotator, ResearchAgent, Visualizer).
- For each:
  - Permissions,
  - Triggers (time‑based, file‑based, manual),
  - Failure handling.

---

## 5. Deep Dives You’ll “Dive In” To

Over the next 1–2 years, you’ll repeatedly go deep on:

1. **Nushell Orchestration Patterns**
   - Clean separation between:
     - CLI parsing,
     - Orchestration logic,
     - Tool invocations (Python, Rust, DB).
   - Strong `--debug` tracing and error handling.

2. **Python Solver Evolution**
   - From single‑model “llama3.2:3b” to:
     - Ensemble of local + online models,
     - Self‑consistency checks,
     - Tool‑using solvers (symbolic algebra libraries, etc.).

3. **Rust/Noon**
   - Building a general scene description language for math/physics.
   - Performance and rendering pipelines.

4. **Data Modeling & Analytics**
   - Evolving the schema as:
     - You discover more problem types,
     - You add prediction and mastery tracking.

5. **Safety & Trust**
   - Continuously refining:
     - Which agents are allowed to do what,
     - How to observe them,
     - How to recover from mistakes.

---

If you want, next step I can:

- Turn this into an actual `docs/MASTER_PLAN.md` file with more concrete sub‑tasks under each phase, or
- Zoom in on **one phase** (e.g., “Hyper‑Tutor Pipeline in 3–9 months”) and break it into implementable milestones compatible with your current `bac.nu` + Python solver structure.
