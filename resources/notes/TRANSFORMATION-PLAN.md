# BAC Study Vault - Transformation Plan

## Vision
Transform `resources/notes/` into a professional textbook-style vault accessible from everywhere (phone + desktop) via git sync.

## Current State
- Vault already has git remote: `https://github.com/amdajhbdh/notes.git`
- Contains ~80+ PDFs (resources from teachers, exam papers)
- Mixed structure with scripts, databases, and notes

## Target Book Structure

```
Vault (Standalone Git Repo)
├── 00-FrontMatter/
│   ├── COVER.md              # Book cover with title, author
│   ├── PREFACE.md            # Why this book, how to use
│   ├── TABLE_OF_CONTENTS.md  # Auto-generated TOC
│   └── ACKNOWLEDGMENTS.md
│
├── 01-Foundations/           # Part I: Core Concepts
│   ├── 01-Mathematics/
│   │   ├── 01-Algebra/
│   │   ├── 02-Analysis/
│   │   └── 03-Geometry/
│   ├── 02-Physics/
│   │   ├── 01-Mechanics/
│   │   ├── 02-Electricity/
│   │   └── 03-Optics/
│   ├── 03-Chemistry/
│   │   ├── 01-Organic/
│   │   └── 02-Inorganic/
│   └── 04-Biology/
│       ├── 01-CellBiology/
│       ├── 02-Genetics/
│       └── 03-Ecology/
│
├── 02-Practice/             # Part II: Exercises & Problems
│   ├── 01-MathExercises/
│   ├── 02-PhysicsExercises/
│   ├── 03-ChemistryExercises/
│   └── 04-BiologyExercises/
│
├── 03-Resources/            # Part III: Teacher Resources
│   ├── 01-TeacherNotes/     # Handwritten teacher lessons (OCR'd)
│   ├── 02-Series/           # Exercise series
│   └── 03-BookExcerpts/     # Scanned book pages
│
├── 04-Exams/               # Part IV: Past Exams
│   ├── 01-BAC-Papers/      # 2000-2025
│   │   ├── 2025/
│   │   ├── 2024/
│   │   └── ...
│   └── 02-MockExams/
│
├── 05-Reference/           # Part V: Quick Reference
│   ├── FORMULAS.md         # All formulas
│   ├── DEFINITIONS.md      # Key definitions
│   └── GLOSSARY.md         # Subject glossary
│
├── 06-Appendices/
│   ├── SOLUTIONS/          # Full solutions
│   └── ANNEXES/
│
└── .obsidian/              # Obsidian config (synced)
```

## Phase 1: Vault Cleanup & Restructure

### Task 1.1: Remove Clutter
- [ ] Remove Python scripts (*.py) - move to separate repo
- [ ] Remove databases (*.db files)
- [ ] Remove large generated files (*.json > 1MB)
- [ ] Remove build artifacts (target/, Cargo.lock)
- [ ] Keep only: markdown, images, PDFs, Obsidian config

### Task 1.2: Create Book Structure
- [ ] Create 00-FrontMatter folder with proper front matter
- [ ] Create TOC file with links to all sections
- [ ] Move existing content to proper locations

## Phase 2: OCR Pipeline

### Task 2.1: Install & Configure OCR
- [ ] Install Tesseract: `sudo apt install tesseract-ocr tesseract-ara`
- [ ] Install Surya (if needed)
- [ ] Test OCR on sample PDF

### Task 2.2: OCR Resources
- [ ] OCR all PDFs in 03-Resources/
- [ ] OCR all PDFs in 07-Assets/PDFs/
- [ ] OCR mauritanian-bac PDFs
- [ ] Save as markdown in 03-Resources/OCR/

### Task 2.3: Format Extracted Content
- [ ] Add proper front matter to each extracted note
- [ ] Add wikilinks between related concepts
- [ ] Add tags for searchability

## Phase 3: Phone Sync

### Task 3.1: Git-Based Sync
- [ ] Ensure git credentials configured
- [ ] Create simple sync script
- [ ] Document sync process for mobile

### Task 3.2: Mobile Setup Instructions
- [ ] Install Obsidian on mobile
- [ ] Clone vault from GitHub
- [ ] Configure sync (git pull/push workflow)

## Phase 4: Multi-Cloud Infrastructure (Future)

### Task 4.1: CloudShell
- [ ] Development environment
- [ ] OCR processing
- [ ] Git operations

### Task 4.2: Cloudflare
- [ ] Workers for API
- [ ] Pages for static content
- [ ] R2 for storage

### Task 4.3: Render
- [ ] Go API hosting
- [ ] Agent CLI hosting

### Task 4.4: Fly.io
- [ ] Edge services
- [ ] Low-latency requests

### Task 4.5: Doppler
- [ ] Secrets management
- [ ] Environment variables

## Quick Commands

```bash
# Sync vault from phone
cd /path/to/vault && git pull

# After editing on phone
git add . && git commit -m "update" && git push

# OCR a PDF
tesseract -l ara+eng input.pdf output

# Batch OCR
for f in *.pdf; do tesseract -l ara+eng "$f" "${f%.pdf}"; done
```

## Status: IN PROGRESS
Last Updated: 2026-03-13
