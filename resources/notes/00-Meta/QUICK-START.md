---
tags: [meta, guide, quick-start]
---

# 🚀 Quick Start Guide

> [!info] Welcome to Your BAC Study Vault!
> This vault contains **686+ notes**, **11 textbooks**, and **1,571 exercises** to help you master the 7C/7D curriculum.

---

## ⚡ Quick Commands

> [!example] Essential Commands
> 
> ```bash
> # Search vault
> ./vault-cli.py search "matrices"
> 
> # Review due cards
> python3 ../spaced-repetition.py due
> 
> # Check progress
> ./study-cli.py stats
> 
> # Enhance notes
> python3 ../content-enhancer.py enhance-all 01-Concepts
> 
> # Sync to cloud
> jj commit -m "Study session"
> jj git push
> ```

---

## 📁 Vault Structure

> [!note] Folder Organization
> 
> | Folder | Contents | Count |
> |--------|----------|-------|
> | `00-Meta/` | Dashboards, MOCs, guides | 11 |
> | `01-Concepts/` | Enhanced concept notes | 26 |
> | `02-Practice/` | Exercises + textbook problems | 1,571 |
> | `03-Resources/` | 11 PDF textbooks | 11 |
> | `04-Exams/` | BAC papers 2002-2012 | 23 |
> | `05-Extracted/` | Textbook chapters | 550+ |
> | `06-Daily/` | Daily tracking | 2 |
> | `07-Assets/` | PDFs, images | - |
> | `08-Templates/` | Note templates | 3 |

---

## 🎯 Study Workflow

> [!todo] Daily Routine
> 
> ### 1. Morning Review (30 min)
> - [ ] Check spaced repetition queue
> - [ ] Review 10-20 cards
> - [ ] Update progress tracker
> 
> ### 2. Study Session (2-3 hours)
> - [ ] Pick a topic from project dashboard
> - [ ] Read concept notes
> - [ ] Check related notes (Smart Connections)
> - [ ] Practice exercises
> 
> ### 3. Evening Review (30 min)
> - [ ] Summarize what you learned
> - [ ] Add new cards to spaced repetition
> - [ ] Update daily note
> - [ ] Sync to cloud

---

## 📊 Key Dashboards

> [!summary] Important Files
> 
> - **[[00-Meta/Projects/BAC-Prep-2026|BAC Project Dashboard]]** - Main tracking
> - **[[00-Meta/PROGRESS-TABLE|Progress Table]]** - Vault statistics
> - **[[06-Daily/Progress-Dashboard|Daily Progress]]** - Today's stats
> - **[[03-Resources/INDEX|Resource Index]]** - Textbook lookup
> - **[[00-Meta/CONTENT-MAPPING|Content Mapping]]** - Gap analysis

---

## 🔍 Search & Discovery

> [!tip] Finding Content
> 
> ### CLI Search
> ```bash
> ./vault-cli.py search "complex numbers"
> ./vault-cli.py search "رياضيات"  # Arabic search
> ```
> 
> ### Obsidian Search
> - `Ctrl/Cmd + O` - Quick switcher
> - `Ctrl/Cmd + Shift + F` - Global search
> - Smart Connections panel - Related notes
> 
> ### Tags
> - `#difficulty/easy` `#difficulty/medium` `#difficulty/hard`
> - `#status/not-started` `#status/in-progress` `#status/completed`
> - `#math` `#physics` `#chemistry`

---

## 🎓 Study by Subject

> [!note] Subject Navigation
> 
> ### Mathematics (490 chapters)
> - [[01-Concepts/Math/Complex-Numbers/Complex Numbers - Basics|Complex Numbers]]
> - [[01-Concepts/Math/Matrices/Matrices - Basics|Matrices]]
> - [[01-Concepts/Math/Integrals/Integrals - Fundamentals|Integrals]]
> - Browse: `05-Extracted/Resources/Math/`
> 
> ### Physics (275 chapters)
> - [[01-Concepts/Physics/Dynamics/Dynamics - Newton Laws|Newton's Laws]]
> - [[01-Concepts/Physics/Electromagnetism/EM - Electric Fields|Electric Fields]]
> - [[01-Concepts/Physics/Circuits/Circuits - Ohm Law|Ohm's Law]]
> - Browse: `05-Extracted/Resources/Physics/`
> 
> ### Chemistry (176 chapters)
> - [[01-Concepts/Chemistry/Solutions/Solutions - Concentration|Solutions]]
> - [[01-Concepts/Chemistry/Organic/Organic - Hydrocarbons|Organic Chemistry]]
> - [[01-Concepts/Chemistry/Kinetics/Kinetics - Reaction Rates|Kinetics]]
> - Browse: `05-Extracted/Resources/Chemistry/`

---

## 🔄 Spaced Repetition

> [!success] Review System
> 
> ### Check What's Due
> ```bash
> python3 ../spaced-repetition.py due
> ```
> 
> ### Review Cards
> ```bash
> ./study-cli.py review
> ```
> 
> ### Add New Cards
> 1. Open any concept note
> 2. Add to frontmatter:
> ```yaml
> ---
> sr-due: 2026-03-10
> sr-interval: 1
> sr-ease: 250
> ---
> ```

---

## 📈 Track Progress

> [!stats] Progress Tracking
> 
> ### View Statistics
> ```bash
> ./study-cli.py stats --period week
> ```
> 
> ### Update Progress
> - Mark notes as completed in frontmatter
> - Update `progress: 0` → `progress: 100`
> - Change `status: not-started` → `status: completed`
> 
> ### View Dashboard
> - Open: `00-Meta/Projects/BAC-Prep-2026.md`
> - Check weak areas
> - Review mastered topics

---

## 🔗 Smart Connections

> [!related] Discover Related Content
> 
> ### Enable Plugin
> 1. Settings → Community plugins
> 2. Enable "Smart Connections"
> 3. View panel on right sidebar
> 
> ### Usage
> - Open any note
> - Check Smart Connections panel
> - See automatically related notes
> - Click to navigate

---

## 📱 Sync to Cloud

> [!example] Synchronization
> 
> ### Sync Notes (jj/Git)
> ```bash
> jj commit -m "Study session: matrices"
> jj git push
> ```
> 
> ### Sync PDFs (Croc)
> ```bash
> ./sync-vault.py send
> # On Cloud Shell:
> croc <code>
> ```
> 
> ### Full Sync
> ```bash
> ./sync-vault.py full
> ```

---

## 🎨 Customization

> [!tip] Personalize Your Vault
> 
> ### Enable CSS Styling
> 1. Settings → Appearance → CSS snippets
> 2. Enable "vault-styling"
> 3. Reload Obsidian
> 
> ### Customize Templates
> - Edit files in `08-Templates/`
> - Use Templater plugin
> 
> ### Add Plugins
> - Settings → Community plugins
> - Browse and install
> - Recommended: Dataview, Projects, Smart Connections

---

## 🆘 Troubleshooting

> [!warning] Common Issues
> 
> ### Search Not Working
> ```bash
> # Rebuild index
> ./vault-cli.py rebuild-index
> ```
> 
> ### Sync Failed
> ```bash
> # Check jj status
> jj status
> jj log
> ```
> 
> ### Missing Notes
> - Check `.gitignore` (PDFs excluded)
> - Use Croc for large files
> 
> ### Plugin Errors
> - Disable and re-enable plugin
> - Check plugin settings
> - Restart Obsidian

---

## 📚 Documentation

> [!note] More Guides
> 
> - [[00-Meta/VAULT-SUMMARY|Vault Summary]] - Complete overview
> - [[00-Meta/SYNC-GUIDE|Sync Guide]] - Detailed sync instructions
> - [[00-Meta/MCP-SERVER-GUIDE|MCP Guide]] - Server usage
> - [[00-Meta/TASK-12-COMPLETE|Enhancement Guide]] - Content enhancement
> - [[00-Meta/TASK-12-STYLING-COMPLETE|Styling Guide]] - Callout usage

---

## 🎯 Next Steps

> [!success] Get Started!
> 
> 1. ✅ Enable CSS snippet for styling
> 2. ✅ Open project dashboard: `00-Meta/Projects/BAC-Prep-2026.md`
> 3. ✅ Pick your first topic to study
> 4. ✅ Review related notes with Smart Connections
> 5. ✅ Practice exercises from `02-Practice/`
> 6. ✅ Add cards to spaced repetition
> 7. ✅ Track your progress daily
> 8. ✅ Sync to cloud regularly

---

**Good luck with your BAC preparation! 🎓✨**

