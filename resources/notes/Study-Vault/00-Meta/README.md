# Study Vault - Complete Setup Guide

Welcome to your 2-day intensive study vault! This is a fully structured Obsidian/Logseq vault with interconnected notes, practice problems, and resources.

## Quick Start

### For Obsidian
1. Open Obsidian
2. Click "Open folder as vault"
3. Select the `Study-Vault` folder
4. Install community plugins (see below)
5. Start with `000-INDEX.md`

### For Logseq
1. Open Logseq
2. Add graph  ->  Select `Study-Vault` folder
3. Start with `000-INDEX.md`

---

## 📁 Vault Structure

```
Study-Vault/
├── 000-INDEX.md                 # Main hub - START HERE
├── Math MOC.md                  # Mathematics hub
├── Physics MOC.md               # Physics hub
├── Chemistry MOC.md             # Chemistry hub
├── Natural-Sciences MOC.md      # Biology/Geology hub
│
├── Math/
│   ├── Complex-Numbers/
│   │   ├── Complex Numbers MOC.md
│   │   ├── Complex Numbers - Basics.md
│   │   ├── Complex Numbers - Operations.md
│   │   ├── Complex Numbers - Polar Form.md
│   │   ├── Complex Numbers - Practice Easy.md
│   │   ├── Complex Numbers - Practice Medium.md
│   │   └── Complex Numbers - Practice Hard.md
│   ├── Integrals/
│   ├── Number-Theory/
│   └── Matrices/
│
├── Physics/
│   ├── Dynamics/
│   ├── Electromagnetism/
│   ├── Satellites/
│   ├── Energy/
│   └── Circuits/
│
├── Chemistry/
│   ├── Solutions/
│   ├── Organic/
│   └── Kinetics/
│
├── Natural-Sciences/
│   ├── Genetics/
│   ├── Reproduction/
│   ├── Nervous-System/
│   └── Geology/
│
├── Templates/
│   ├── Concept-Note-Template.md
│   ├── Practice-Template.md
│   └── Daily-Note-Template.md
│
├── Assets/
│   ├── Images/
│   ├── Diagrams/
│   └── PDFs/
│
└── Daily-Notes/
```

---

## 🔌 Required Plugins

### Obsidian Community Plugins

#### Essential (Install First)
1. **Excalidraw** - Draw diagrams, math visualizations
2. **Dataview** - Dynamic queries and tables
3. **Calendar** - Daily notes navigation
4. **Tasks** - Task management
5. **Templater** - Advanced templates

#### Highly Recommended
6. **Obsidian Git** - Backup to GitHub
7. **QuickAdd** - Quick note creation
8. **Tag Wrangler** - Tag management
9. **LaTeX Suite** - Fast math typing
10. **Mind Map** - Visual concept maps
11. **Kanban** - Task boards
12. **Advanced Tables** - Better table editing
13. **Spaced Repetition** - Flashcards

### How to Install
1. Settings  ->  Community plugins  ->  Turn off Safe Mode
2. Browse  ->  Search for plugin name  ->  Install  ->  Enable

---

## Using Excalidraw for Diagrams

### Create a Diagram
1. Create new note: `diagram-name.excalidraw`
2. Draw your diagram
3. Embed in notes: `![[diagram-name.excalidraw]]`

### Recommended Diagrams to Create
- **Math**: Complex plane, integration areas, matrix transformations
- **Physics**: Free body diagrams, field lines, circuit diagrams
- **Chemistry**: Molecular structures, reaction mechanisms
- **Biology**: Cell diagrams, genetic crosses, neuron structure

---

## Note-Taking System

### Linking Strategy
- Use `[[Note Name]]` to link between notes
- Use `#tags` for categorization
- Use `**block-id` for block references

### Tags System
- **By Subject**: `#math`, `#physics`, `#chemistry`, `#biology`
- **By Type**: `#concept-note`, `#practice`, `#formula-sheet`, `#diagram`
- **By Difficulty**: `#easy-practice`, `#medium-practice`, `#hard-practice`
- **By Status**: `#todo`, `#in-progress`, `#completed`, `#needs-review`

### MOC (Map of Content) Pattern
Each major topic has a MOC that links to all related notes:
- Main Index  ->  Subject MOC  ->  Topic MOC  ->  Individual Notes

---

## Using Dataview Queries

### Track Progress
```dataview
TABLE status, difficulty
FROM "Math"
WHERE contains(tags, "practice")
SORT difficulty ASC
```

### Find Incomplete Tasks
```dataview
TASK
WHERE !completed
GROUP BY file.folder
```

### List by Tag
```dataview
LIST
FROM #concept-note
WHERE contains(file.path, "Math")
```

---

##  Quick Workflows

### Daily Study Session
1. Open `Daily-Notes/YYYY-MM-DD.md`
2. List today's topics
3. Link to relevant MOCs
4. Track completed practice problems
5. Note areas needing review

### Learning New Concept
1. Start at Subject MOC
2. Navigate to Topic MOC
3. Read concept notes in order
4. Do easy practice problems
5. Progress to medium/hard
6. Create summary diagram

### Review Session
1. Use Dataview to find `#needs-review` notes
2. Redo practice problems
3. Create flashcards with Spaced Repetition plugin
4. Update status tags

---

## Study Strategy

### Day 1 (8 hours)
**Morning (3h)**: Math
- Complex Numbers: 1.5h
- Integrals: 1.5h

**Afternoon (3h)**: Physics
- Dynamics: 1.5h
- Electromagnetism: 1.5h

**Evening (2h)**: Chemistry + Biology
- Solutions Chemistry: 1h
- Genetics: 1h

### Day 2 (8 hours)
**Morning (3h)**: Math
- Number Theory: 1.5h
- Matrices: 1.5h

**Afternoon (3h)**: Physics
- Satellites: 1h
- Energy: 1h
- Circuits: 1h

**Evening (2h)**: Chemistry + Biology
- Organic Chemistry: 1h
- Nervous System: 1h

---

## Pro Tips

### Math
- Use Excalidraw for complex plane diagrams
- Practice problems in order: easy  ->  medium  ->  hard
- Create formula sheets for quick reference
- Use Desmos/GeoGebra for visualization

### Physics
- Draw diagrams for EVERY problem
- Use PhET simulations
- Create free body diagrams in Excalidraw
- Link physics concepts to math (calculus, vectors)

### Chemistry
- Draw molecular structures in Excalidraw
- Create reaction mechanism flowcharts
- Use ChemDraw or similar for complex molecules
- Link to real-world applications

### Biology
- Create visual diagrams (cells, organs, systems)
- Use Anki/Spaced Repetition for terminology
- Draw genetic crosses and Punnett squares
- Link concepts across topics

---

## External Tools Integration

### Math Tools
- [Desmos](https://www.desmos.com/) - Graphing
- [GeoGebra](https://www.geogebra.org/) - Geometry
- [Wolfram Alpha](https://www.wolframalpha.com/) - Computation
- [Symbolab](https://www.symbolab.com/) - Step-by-step solutions

### Physics Tools
- [PhET Simulations](https://phet.colorado.edu/)
- [HyperPhysics](http://hyperphysics.phy-astr.gsu.edu/)
- [Physics Classroom](https://www.physicsclassroom.com/)

### Chemistry Tools
- [ChemDraw](https://www.perkinelmer.com/category/chemdraw)
- [PubChem](https://pubchem.ncbi.nlm.nih.gov/)
- [Periodic Table](https://ptable.com/)

### Biology Tools
- [Khan Academy Biology](https://www.khanacademy.org/science/biology)
- [NCBI](https://www.ncbi.nlm.nih.gov/)
- [BioDigital Human](https://www.biodigital.com/)

---

##  Learning Techniques

### Active Recall
- Close notes and explain concepts aloud
- Teach concepts to someone else
- Create questions before looking at practice problems

### Spaced Repetition
- Review notes after 1 day, 3 days, 7 days
- Use Obsidian Spaced Repetition plugin
- Create flashcards for key formulas

### Feynman Technique
1. Choose a concept
2. Explain it in simple terms
3. Identify gaps in understanding
4. Review and simplify further

### Interleaving
- Mix different topics in one session
- Don't study one topic for too long
- Switch between subjects every 1-2 hours

---

##  Mobile Access

### Obsidian Mobile
- Install Obsidian app on phone/tablet
- Sync vault via iCloud/Dropbox/Obsidian Sync
- Review notes on the go

### Logseq Mobile
- Install Logseq app
- Sync via iCloud/Dropbox
- Quick capture ideas

---

##  Backup Strategy

### Option 1: Obsidian Git Plugin
1. Create GitHub repository
2. Install Obsidian Git plugin
3. Configure auto-backup every 30 minutes

### Option 2: Cloud Sync
- Use Dropbox, Google Drive, or iCloud
- Place vault in synced folder
- Automatic backup

### Option 3: Manual Backup
- Copy vault folder weekly
- Use external drive or cloud storage

---

## 🆘 Troubleshooting

### Plugins Not Working
- Check if Safe Mode is off
- Restart Obsidian
- Reinstall plugin

### Links Not Working
- Use `[[Note Name]]` format
- Check file exists
- Ensure correct capitalization

### Math Not Rendering
- Use `$inline$` for inline math
- Use `$$block$$` for block math
- Install LaTeX Suite plugin

### Excalidraw Issues
- Update to latest version
- Check file extension is `.excalidraw`
- Restart Obsidian

---

## Track Your Progress

Create a progress note:

```markdown
# Progress Tracker

## Math
- [x] Complex Numbers - Easy (20/20)
- [x] Complex Numbers - Medium (15/20)
- [ ] Complex Numbers - Hard (0/20)
- [ ] Integrals - Easy (0/20)

## Physics
- [ ] Dynamics - Easy (0/20)

## Chemistry
- [ ] Solutions - Easy (0/20)

## Biology
- [ ] Genetics - Easy (0/20)
```

---

## 🎉 You're Ready!

Start with `000-INDEX.md` and follow the study plan. Good luck! 

---

**Questions?** Check the templates in `Templates/` folder for examples.

**Need more resources?** Each MOC has curated links to videos, websites, and practice problems.

**Stuck on a concept?** Use the linking system to find related notes and build understanding.
