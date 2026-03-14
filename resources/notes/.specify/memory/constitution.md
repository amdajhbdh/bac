# BAC Study Vault Constitution

<!-- Sync Impact Report (2026-03-11):
Version: N/A → 1.0.0 (initial)
Added: 5 Core Principles, 2 Additional Sections, Governance
Templates verified: plan-template.md ✅ | spec-template.md ✅ | tasks-template.md ✅
-->

## Core Principles

### I. Mathematical Rigor
All mathematical and scientific content MUST use LaTeX notation. Inline math uses `$...$`, display math uses `$$...$$`. Forbidden: ASCII subscripts (H₂O), superscripts (x²), Greek letters (α, β, γ), special symbols (≤, ≥, ∫). Rationale: Ensures consistent rendering across all platforms and maintains academic quality.

### II. Content Organization
Files follow structure: `Topic-Subtopic.md`, include level `(7ème/Terminale)`, add source `*Source: ...*`. Section headers use markdown hierarchy. Exercise format: Énoncé → Méthode de résolution → Résultat. Rationale: Enables efficient navigation and ensures content is appropriately classified by academic level.

### III. Quality Standards
All content must include units, proper notation ($[H_3O^+]$ not $[H3O+]$), equilibrium arrows ($\rightleftharpoons$), and correct chemical formulas. Language: French for Mauritanian BAC, English for general concepts, bilingual when helpful. Rationale: Maintains academic integrity and supports diverse learners.

### IV. Automated Validation
All content passes: ASCII math scan, LaTeX closure check, level verification, source attribution check. Python code passes lint (`ruff check .`) and formatting (`ruff format .`). Rationale: Catches errors before publication and ensures consistency across 900+ files.

### V. Agent Responsibilities
Content Agents: Convert ASCII math to LaTeX, maintain source attribution. QA Agents: Scan for non-LaTeX notation, verify equations. Code Agents: Follow Python style (type hints, docstrings), run linting. Update Agents: Apply formatting rules, maintain consistency. Rationale: Clear role definitions prevent work duplication and ensure accountability.

---

## Quality Standards

### Technology Stack
- **Python**: 3.11+, use ruff for lint/format, pytest for tests
- **Storage**: SQLite (vault-index.db, spaced-repetition.db, progress.db, agents.db)
- **Search**: FTS5 for full-text search with Arabic + English support
- **Sync**: jj/Git for version control, Croc for P2P file transfer

### Content Requirements
- **Math Notation**: LaTeX mandatory (see Core Principles)
- **File Naming**: kebab-case with level indicator
- **Source Attribution**: Required for all extracted content
- **Level Classification**: Basic (6ème), Intermediate (7ème C/D), Advanced (7ème/Terminale)

---

## Development Workflow

### Python Code Standards
- **Imports**: Group stdlib → third-party → local; avoid wildcards; sort alphabetically
- **Formatting**: 100 chars max, 4 spaces, trailing commas
- **Types**: Use type hints, prefer `Optional[X]`, concrete types for public APIs
- **Naming**: Classes `PascalCase`, funcs/vars `snake_case`, constants `SCREAMING_SNAKE_CASE`
- **Error Handling**: Specific exceptions, meaningful messages, never silent catch
- **Testing**: Run single test with `pytest tests/test_file.py::test_function_name`

### Content Creation Standards
- Convert all ASCII/Unicode math to LaTeX during extraction
- Maintain source attribution and academic level
- Organize into appropriate subject folders
- Include callouts: summary, related, progress, stats

### Review Process
- All PRs must verify LaTeX compliance
- Code changes require `ruff check . --fix` before commit
- Content extraction requires source verification
- Complexity must be justified in plan documents

---

## Governance

### Amendment Procedure
1. Proposal documents specific changes with rationale
2. Review against existing principles for consistency
3. Version bump per semantic rules:
   - MAJOR: Backward incompatible governance/principle changes
   - MINOR: New principle or materially expanded guidance
   - PATCH: Clarifications, wording, typo fixes
4. Update LAST_AMENDED_DATE to current date
5. Sync Impact Report added as HTML comment

### Compliance
- All templates must align with constitution principles
- Plan template Constitution Check gate enforced
- Spec template user story independence required
- Tasks template phase dependencies honored

### References
- Runtime guidance: `AGENTS.md` for agent-specific rules
- Quick reference: `QUICK-COMMANDS.md`
- Templates: `.specify/templates/`

---

**Version**: 1.0.0 | **Ratified**: 2026-03-11 | **Last Amended**: 2026-03-11
