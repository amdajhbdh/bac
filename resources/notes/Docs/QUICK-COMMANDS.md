---
tags: [meta, commands, reference]
---

# ⚡ Quick Commands Reference

> [!info] Essential Commands
> Copy-paste these commands for daily use.

---

## 🔍 Search & Discovery

> [!example] Finding Content
> 
> ```bash
> # Search vault (English)
> ./vault-cli.py search "matrices"
> 
> # Search vault (Arabic)
> ./vault-cli.py search "رياضيات"
> 
> # Search with stats
> ./vault-cli.py stats
> 
> # Rebuild search index
> ./vault-cli.py rebuild-index
> ```

---

## 🔄 Spaced Repetition

> [!todo] Review System
> 
> ```bash
> # Check what's due today
> python3 ../spaced-repetition.py due
> 
> # Review cards
> ./study-cli.py review
> 
> # Add note to SR
> python3 ../spaced-repetition.py add "01-Concepts/Math/Matrices/Matrices - Basics.md"
> 
> # View statistics
> python3 ../spaced-repetition.py stats
> ```

---

## 📊 Progress Tracking

> [!stats] Analytics
> 
> ```bash
> # View weekly stats
> ./study-cli.py stats --period week
> 
> # View monthly stats
> ./study-cli.py stats --period month
> 
> # View all-time stats
> ./study-cli.py stats
> ```

---

## 🤖 Agent System

> [!tip] Automated Tasks
> 
> ```bash
> # Scan for gaps
> python3 ../agent-system.py scanner scan_gaps '{}'
> 
> # Find orphaned notes
> python3 ../agent-system.py scanner scan_orphans '{}'
> 
> # Analyze difficulty
> python3 ../agent-system.py analyzer analyze_difficulty '{"path":"01-Concepts/Math"}'
> 
> # Find resources
> python3 ../agent-system.py resource_finder find_videos '{"topic":"matrices"}'
> 
> # Create study plan
> python3 ../agent-system.py study_coach create_plan '{"duration":3,"focus":["Math"]}'
> 
> # Build links
> python3 ../agent-system.py link_builder suggest_links '{"path":"01-Concepts/Math/Matrices/Matrices - Basics.md"}'
> 
> # Generate quiz
> python3 ../agent-system.py quiz_generator generate_quiz '{"topic":"complex numbers","count":10}'
> ```

---

## ✨ Content Enhancement

> [!success] Enhance Notes
> 
> ```bash
> # Enhance single note
> python3 ../content-enhancer.py enhance "01-Concepts/Math/Matrices/Matrices - Basics.md"
> 
> # Enhance all concept notes
> python3 ../content-enhancer.py enhance-all 01-Concepts
> 
> # Enhance textbook notes
> python3 ../content-enhancer.py enhance-all 05-Extracted/Resources
> 
> # Enhance practice problems
> python3 ../content-enhancer.py enhance-all 02-Practice
> 
> # Generate progress table
> python3 ../content-enhancer.py progress
> ```

---

## 🔄 Sync & Backup

> [!example] Cloud Synchronization
> 
> ```bash
> # Commit changes (jj)
> jj commit -m "Study session: matrices"
> 
> # Push to Git
> jj git push
> 
> # Pull from Git
> jj git fetch
> 
> # Send PDFs via Croc
> ./sync-vault.py send
> 
> # Receive PDFs via Croc
> ./sync-vault.py receive <code>
> 
> # Full sync
> ./sync-vault.py full
> ```

---

## 📝 Note Creation

> [!note] Quick Note Templates
> 
> ```bash
> # Create concept note
> cp 08-Templates/Concept-Note-Template.md "01-Concepts/Math/New-Topic.md"
> 
> # Create practice problem
> cp 08-Templates/Practice-Template.md "02-Practice/Math/Problem-1.md"
> 
> # Create daily note
> cp 08-Templates/Study-Session-Template.md "06-Daily/$(date +%Y-%m-%d).md"
> ```

---

## 🔧 Maintenance

> [!warning] System Maintenance
> 
> ```bash
> # Check vault status
> jj status
> 
> # View recent changes
> jj log
> 
> # Check database size
> du -h .vault-index.db .agents.db .spaced-repetition.db .progress.db
> 
> # Count notes
> find . -name "*.md" -type f | wc -l
> 
> # Find large files
> find . -type f -size +10M
> ```

---

## 📊 Statistics

> [!stats] Quick Stats
> 
> ```bash
> # Count notes by folder
> find 01-Concepts -name "*.md" | wc -l
> find 05-Extracted -name "*.md" | wc -l
> 
> # Count exercises
> find 02-Practice -name "*.md" | wc -l
> 
> # Count BAC papers
> find 04-Exams -name "*.md" | wc -l
> 
> # Database stats
> sqlite3 .vault-index.db "SELECT COUNT(*) FROM notes"
> sqlite3 .spaced-repetition.db "SELECT COUNT(*) FROM cards"
> ```

---

## 🎯 Study Workflow

> [!success] Daily Routine
> 
> ```bash
> # Morning (30 min)
> python3 ../spaced-repetition.py due
> ./study-cli.py review
> 
> # Study Session (2-3 hours)
> ./vault-cli.py search "today's topic"
> # Read notes in Obsidian
> # Practice exercises
> 
> # Evening (30 min)
> ./study-cli.py stats --period week
> python3 ../content-enhancer.py progress
> jj commit -m "Study session: $(date +%Y-%m-%d)"
> jj git push
> ```

---

## 🚀 One-Liners

> [!tip] Power User Commands
> 
> ```bash
> # Find all hard topics
> grep -r "difficulty: hard" 01-Concepts
> 
> # Find incomplete notes
> grep -r "status: not-started" 01-Concepts
> 
> # Count formulas
> grep -r '\$\$' 01-Concepts | wc -l
> 
> # Find notes without enhancements
> find 01-Concepts -name "*.md" -exec grep -L "summary" {} \;
> 
> # List all tags
> grep -rh "^tags:" 01-Concepts | sort | uniq
> ```

---

## 📱 Obsidian Shortcuts

> [!note] Keyboard Shortcuts
> 
> | Action | Shortcut |
> |--------|----------|
> | Quick switcher | `Ctrl/Cmd + O` |
> | Global search | `Ctrl/Cmd + Shift + F` |
> | Command palette | `Ctrl/Cmd + P` |
> | Toggle edit/preview | `Ctrl/Cmd + E` |
> | Open graph view | `Ctrl/Cmd + G` |
> | Create new note | `Ctrl/Cmd + N` |
> | Open settings | `Ctrl/Cmd + ,` |

---

## 🔗 Quick Navigation

> [!related] Important Files
> 
> ```bash
> # Open in Obsidian
> obsidian://open?vault=Study-Vault&file=00-Meta/Projects/BAC-Prep-2026
> 
> # Open in terminal
> cd /home/med/Documents/bac/resources/notes/Study-Vault
> 
> # Open specific note
> vim "01-Concepts/Math/Matrices/Matrices - Basics.md"
> ```

---

**Bookmark this page for quick reference! ⚡**

