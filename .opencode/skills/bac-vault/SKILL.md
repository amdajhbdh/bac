---
name: bac-vault
description: Manage Obsidian vault with git sync, search, and statistics. Use when working with resources/notes vault, syncing between devices, or searching notes. Triggers: vault, sync, git, notes, obsidian, search
---

# BAC Vault Management Skill

Manage the study vault with git sync and search.

## Quick Commands

```bash
# Sync vault
cd resources/notes
git pull origin main
# Make changes
git add -A && git commit -m "update" && git push origin main

# Search notes
grep -r "query" --include="*.md" .

# Statistics
find . -name "*.md" | wc -l
find 03-Resources -name "*.pdf" | wc -l
```

## Vault Structure

```
resources/notes/
├── 00-Meta/           # Index, TOC
├── 01-Concepts/       # Subject notes
├── 02-Practice/      # Exercises
├── 03-Resources/      # PDFs to OCR
├── 04-Exams/          # Past papers
├── 05-Extracted/      # OCR output
├── 06-Daily/          # Daily notes
└── 07-Assets/        # Images, media
```

## Mobile Access

1. Clone: `git clone https://github.com/amdajhbdh/notes.git`
2. Edit in Obsidian
3. Push: `git add . && git commit && git push`

## Git Workflow

```bash
# Pull latest
git pull origin main

# Check status
git status

# Add changes
git add -A
git commit -m "description"

# Push
git push origin main
```

## Tips

- Always pull before editing
- Commit frequently
- Use meaningful commit messages
- Keep PDFs in 03-Resources
- OCR'd content goes to 05-Extracted
