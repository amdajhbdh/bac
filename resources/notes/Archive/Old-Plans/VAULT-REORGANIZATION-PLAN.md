# 🎯 VAULT REORGANIZATION PLAN
## Mauritanian BAC Notes - Complete Restructure

---

## 📊 CURRENT STATE ANALYSIS

### Problems Identified:

1. **Scattered Scripts** - 20+ Python/JS scripts in root + 00-Meta/Scripts/
2. **Empty Files** - 7 empty .md files in root
3. **Wrong Excalidraw** - `Induction - Lenz Law.excalidraw.md` (0 bytes) in root
4. **Root Clutter** - 40+ markdown files at root level
5. **Mixed Resources** - 03-Resources has loose files mixed with folders
6. **Duplicate Reports** - Multiple similar report files
7. **Inconsistent Naming** - Some files use spaces, others use hyphens
8. **No Clear Structure** - Hard to navigate

---

## 🎯 TARGET STRUCTURE

```
/home/med/Documents/bac/resources/notes/
├── 00-Meta/                    # Vault metadata & config
│   ├── MOCs/                   # Map of Content files
│   ├── INDEX.md                # Main vault index
│   ├── START-HERE.md          # For new users
│   └── BACKUPS/                # Backup files
│
├── 01-Concepts/                # Learning concepts (organized by subject)
│   ├── Math/
│   │   ├── Algebra/
│   │   ├── Analysis/
│   │   ├── Geometry/
│   │   └── Statistics/
│   ├── Physics/
│   │   ├── Mechanics/
│   │   ├── Electromagnetism/
│   │   ├── Optics/
│   │   └── Thermodynamics/
│   ├── Chemistry/
│   │   ├── Organic/
│   │   ├── Inorganic/
│   │   └── Physical/
│   └── Biology/
│       ├── Genetics/
│       ├── Cell-Biology/
│       └── Ecology/
│
├── 02-Practice/                # Exercises & practice
│   ├── Math/
│   ├── Physics/
│   ├── Chemistry/
│   └── Biology/
│
├── 03-Resources/              # External resources (PDFs, etc.)
│   ├── Math/
│   ├── Physics/
│   ├── Chemistry/
│   └── Biology/
│
├── 04-Exams/                   # Past exam papers
│   ├── BAC-Recent/            # 2023-2025
│   ├── BAC-Archive/           # 2002-2022
│   └── Practice-Tests/         # Mock exams
│
├── 05-Extracted/              # OCR'd content (read-only)
│   ├── Math/
│   ├── Physics/
│   ├── Chemistry/
│   └── Biology/
│
├── 06-Daily/                   # Daily notes & tracking
│   ├── 2026/
│   └── Progress/
│
├── 07-Assets/                  # Images, diagrams
│   ├── Images/
│   ├── Diagrams/
│   └── Animations/
│
├── 08-Templates/              # Note templates
│
├── .opencode/                 # OpenCode config (KEEP)
│
├── Scripts/                   # ◀️ ALL SCRIPTS GO HERE
│   ├── python/                # Python scripts
│   ├── bash/                 # Shell scripts
│   └── js/                   # JavaScript scripts
│
├── Docs/                      # ◀️ DOCUMENTATION
│   ├── QUICK-START.md
│   ├── VAULT-GUIDE.md
│   ├── STYLE-GUIDE.md
│   └── TEMPLATES.md
│
└── Archive/                   # ◀️ OLD FILES GO HERE
    ├── Old-Reports/
    ├── Old-Plans/
    └── Trash/
```

---

## 📋 EXECUTION PHASES

### Phase 1: Cleanup (Quick Wins)
**Time: 10 minutes**

- [ ] Delete empty files: `2026-03-*.md`, `Untitled*.md`
- [ ] Delete broken: `Induction - Lenz Law.excalidraw.md`
- [ ] Remove duplicate content in root

### Phase 2: Scripts Consolidation
**Time: 15 minutes**

- [ ] Create `Scripts/` directory structure:
  ```
  Scripts/
  ├── python/
  │   ├── bac-*.py
  │   ├── content-*.py
  │   ├── ocr-*.py
  │   ├── vault-*.py
  │   └── study-*.py
  ├── bash/
  │   ├── install.sh
  │   ├── ocr-all.sh
  │   ├── ocr-parallel.sh
  │   └── sync-vault.sh
  └── js/
      ├── deep-dive-generator.js
      └── exam-question-generator.js
  ```
- [ ] Move all scripts from root to Scripts/
- [ ] Move scripts from 00-Meta/Scripts/ to Scripts/
- [ ] Update shebangs if needed

### Phase 3: Root Files Organization
**Time: 20 minutes**

**Keep at Root:**
- `README.md` → Move to Docs/ or keep as welcome
- `.opencode/` → KEEP (system)
- `.git/` → KEEP (version control)

**Move to Docs/:**
- `00-READ-ME-FIRST.md` → Docs/
- `QUICK-START.md` → Docs/
- `QUICK-REFERENCE.md` → Docs/
- `QUICK-COMMANDS.md` → Docs/
- `MASTER-INDEX.md` → Docs/

**Move to Archive/Old-Reports/:**
- `*REPORT*.md`
- `*SUMMARY*.md`
- `*COMPLETE*.md`
- `ANALYZER-REPORT.md`
- `AUDIT-REPORT.md`
- `LINK-BUILDER-REPORT.md`
- `SCANNER-REPORT.md`
- `COORDINATOR-FINAL-REPORT.md`
- `STUDY-COACH-REPORT.md`
- `RESOURCE-FINDER-REPORT.md`

**Move to Archive/Old-Plans/:**
- `*PLAN*.md`
- `2Days Plan.md`
- `DAILY-STUDY-SCHEDULE.md`
- `TRANSFORMATION-PLAN.md`
- `RESEARCH-ANALYSIS-PLAN.md`

**Move to Archive/Other/:**
- `LICENSE`
- `COMPLETE-SYSTEM-READY.txt`
- `COMPLETION-SUMMARY.txt`
- `FINAL-SUMMARY.txt`
- `*.txt` files

### Phase 4: Content Standardization
**Time: 30 minutes**

- [ ] Add consistent frontmatter to all concept files:
  ```yaml
  ---
  tags: [concept, math, algebra]
  subject: Mathematics
  topic: Algebra
  level: 7ème/Terminale
  difficulty: easy|medium|hard
  ---
  ```
- [ ] Add consistent section headers
- [ ] Ensure LaTeX formatting ($...$) for math
- [ ] Fix any remaining broken links
- [ ] Add French/English bilingual tags where applicable

### Phase 5: Styling Enhancement
**Time: 20 minutes**

- [ ] Create `Docs/STYLE-GUIDE.md` with:
  - Callout formats (::: tip, ::: warning)
  - Math notation rules
  - File naming conventions
  - Link standards
- [ ] Update `04-Exams/BAC-exam-styles.css` with:
  - Better typography
  - Callout styling
  - Code block styling
  - Table styling

### Phase 6: Index & Navigation
**Time: 15 minutes**

- [ ] Create comprehensive `00-Meta/INDEX.md`:
  ```markdown
  # Vault Index
  
  ## 📚 By Subject
  - [[01-Concepts/Math/INDEX]]
  - [[01-Concepts/Physics/INDEX]]
  - [[01-Concepts/Chemistry/INDEX]]
  - [[01-Concepts/Biology/INDEX]]
  
  ## 📝 By Type
  - Concepts: [[01-Concepts/INDEX]]
  - Practice: [[02-Practice/INDEX]]
  - Resources: [[03-Resources/INDEX]]
  - Exams: [[04-Exams/INDEX]]
  
  ## 🚀 Quick Links
  - [[00-Meta/START-HERE]]
  - [[Docs/QUICK-START]]
  - [[Docs/VAULT-GUIDE]]
  ```
- [ ] Create INDEX.md files in each main folder
- [ ] Add breadcrumb navigation to concept files

### Phase 7: Scripts Automation
**Time: 15 minutes**

- [ ] Create `Scripts/bash/vault-organize.sh`:
  - Auto-organize new files
  - Check for empty files
  - Report statistics
- [ ] Create `Scripts/bash/vault-check.sh`:
  - Check for broken links
  - Check for missing frontmatter
  - Report vault health

---

## 📝 FILE NAMING CONVENTIONS

| Type | Convention | Example |
|------|------------|---------|
| Concepts | `Topic - Subtopic.md` | `Complex Numbers - Operations.md` |
| Practice | `Topic-Practice-N.md` | `Algebra-Practice-1.md` |
| Exams | `BAC-Year-Serie-Topic.md` | `BAC-2025-D-SN-Ex1-Arithmetique.md` |
| Daily | `YYYY-MM-DD-Title.md` | `2026-03-14-Study-Notes.md` |
| Templates | `Type-Template.md` | `Concept-Note-Template.md` |

---

## 🎨 CONTENT STYLE

### Frontmatter (Required)
```yaml
---
tags: [concept, math, complex-numbers]
subject: Mathematics
topic: Complex Numbers
level: Terminale
difficulty: medium
created: 2026-01-15
updated: 2026-03-14
---
```

### Section Headers
```markdown
# Main Topic - Subtopic (Level)

## Concepts Clés / Key Concepts

### Définition / Definition

### Exemples / Examples

## Exercice / Exercise
```

### Math Notation (LaTeX)
```markdown
Inline: $x = \frac{-b \pm \sqrt{b^2-4ac}}{2a}$

Display:
$$
\int_0^1 x^2 dx = \left[\frac{x^3}{3}\right]_0^1 = \frac{1}{3}
$$
```

### Callouts
```markdown
::: tip 💡 Remarque
Important concept here!
:::

::: warning ⚠️ Attention
Common mistake to avoid!
:::

::: note 📝 Note
Additional information.
:::
```

---

## ✅ VALIDATION CHECKLIST

After reorganization, verify:

- [ ] All scripts in `Scripts/` folder
- [ ] No empty files in root
- [ ] No broken links (run `obsidian unresolved`)
- [ ] All concept files have frontmatter
- [ ] Consistent naming throughout
- [ ] INDEX files in each folder
- [ ] README at root explains structure
- [ ] Vault loads correctly in Obsidian
- [ ] All daily notes in 06-Daily/
- [ ] All resources in 03-Resources/

---

## ⏱️ ESTIMATED TIME

| Phase | Time |
|-------|------|
| Phase 1: Cleanup | 10 min |
| Phase 2: Scripts | 15 min |
| Phase 3: Root Files | 20 min |
| Phase 4: Standardization | 30 min |
| Phase 5: Styling | 20 min |
| Phase 6: Navigation | 15 min |
| Phase 7: Automation | 15 min |
| **Total** | **~2 hours** |

---

## 🚀 QUICK START

Want me to execute this plan? I can:

1. **Start Now** - Begin with Phase 1 (quick cleanup)
2. **Do Everything** - Execute all 7 phases
3. **Just Plan** - Save this plan and do it manually
4. **Ask Questions** - Clarify any part of the plan

What would you like me to do?
