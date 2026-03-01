# TOOLS.md - BAC Unified Tool Registry

## Overview
Central registry of all tools, agents, and services used by BAC Unified.
This file serves as both documentation and usage log.

---

## Core Local Tools

### ollama
- **Type**: Local LLM
- **Purpose**: Primary AI solver for math/physics problems
- **Model**: llama3.2:3b
- **Endpoint**: http://localhost:11434
- **Usage**:
  ```bash
  curl -X POST http://localhost:11434/api/generate -d '{"model":"llama3.2:3b","prompt":"..."}'
  ```

### postgres (Neon)
- **Type**: Database + Vector
- **Purpose**: Store problems, solutions, concepts, embeddings
- **Connection**: postgresql://neondb_owner:***@ep-***.aws.neon.tech/neondb
- **Extensions**: pgvector, uuid-ossp, pg_trgm

### playwright-cli
- **Type**: Browser Automation
- **Purpose**: Automate web interactions, query online AI models
- **Commands**:
  ```bash
  playwright-cli open <url>
  playwright-cli type "prompt"
  playwright-cli snapshot
  ```

### noon (Rust)
- **Type**: Animation Engine
- **Purpose**: Generate math/physics visualizations
- **Location**: src/noon/
- **Usage**: `cargo run --example <name>`

---

## Online AI Models (via Playwright)

### deepseek
- **URL**: https://chat.deepseek.com
- **Access**: Browser automation (playwright-cli)
- **Status**: Free tier available
- **Usage**: Query via web interface

### grok
- **URL**: https://grok.com
- **Access**: Browser automation
- **Status**: May require login

### claude
- **URL**: https://claude.ai
- **Access**: Browser automation
- **Status**: May require login

### chatgpt
- **URL**: https://chatgpt.com
- **Access**: Browser automation
- **Status**: Free tier available

---

## Agent System

### orchestrator
- **Type**: Main Agent
- **File**: lib/agent/orchestrator.nu
- **Purpose**: Coordinate all sub-agents, manage flow

### Input Agents
- **classifier**: Detect subject, chapter, concepts
- **ocr**: Extract text from images
- **pdf_parser**: Extract problems from PDFs

### Solver Agents
- **local_solver**: Use Ollama
- **online_solver**: Use Playwright + online AIs

### Learning Agents
- **memory**: Vector DB operations
- **teacher**: Generate explanations, quizzes

### Visualizer
- **noon_agent**: Generate animations

---

## Database Schema

### problems
```sql
CREATE TABLE problems (
    id UUID PRIMARY KEY,
    raw_text TEXT,
    subject VARCHAR(50),
    chapter VARCHAR(100),
    concepts TEXT[],
    difficulty INT,
    embedding vector(1536),
    created_at TIMESTAMP
);
```

### solutions
```sql
CREATE TABLE solutions (
    id UUID PRIMARY KEY,
    problem_id UUID REFERENCES problems(id),
    solver_name TEXT,
    solution_text TEXT,
    steps JSONB,
    quality_score FLOAT
);
```

### concepts
```sql
CREATE TABLE concepts (
    id SERIAL PRIMARY KEY,
    subject VARCHAR(50),
    chapter VARCHAR(100),
    name VARCHAR(200),
    description TEXT,
    embedding vector(384)
);
```

### tool_runs
```sql
CREATE TABLE tool_runs (
    id SERIAL PRIMARY KEY,
    agent_name TEXT,
    action TEXT,
    input JSONB,
    output JSONB,
    duration_ms INT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## Usage Log

Format: `- DATE: tool=NAME action=ACTION notes="..."`

### 2026-02-25
- tool=orchestrator action=init notes="Created agent pipeline"
- tool=local_solver action=solve notes="Solved x²+x-2=0 with llama3.2:3b"
- tool=noon_agent action=generate notes="Created animation for quadratic equation"
- tool=agent action=solve notes="Tested bac agent command with x+5=10"

---

## Safety Constraints

1. **No external API keys** unless explicitly configured
2. **All agent actions** logged to tool_runs
3. **Online AIs** used only as fallback when local fails
4. **No data leaves** user's machine without explicit consent

---

## Adding New Tools

To add a new tool:

1. Add entry to appropriate section above
2. Document in AGENTS.md
3. Log initial usage in Usage Log
4. Update schema if needed

---

*Last Updated: 2026-02-25*
