# 🧬 Biology Exercise 1: Hormone Regulation Animation

**Complete walkthrough of Joussour SN Exercise 1**

---

## 📊 Exercise Content

**Problem:** Two women (M & N) with absent menstruation
- Woman M: FSH=92, LH=60 (very high)
- Woman N: FSH=4, LH=3 (very low)
- Both have zero ovarian hormones

**Questions covered:**
1. Part 1a: Why no menstruation?
2. Part 1b: Could it be pregnancy?
3. Part 2a: Compare FSH/LH levels
4. Part 2b: Propose hypotheses
5. Part 3: Analyze ovarian ultrasound
6. Part 4: Interpret GnRH test
7. Final diagnosis

---

## 🎬 Animation Scenes (7 total)

### Scene 1a: Ovarian Hormones Graph
- Shows zero hormone levels for M & N vs normal cycle
- Duration: ~30 seconds

### Scene 1b: No Menstruation Explanation
- Logic chain: No hormones  ->  No endometrium  ->  No menstruation
- Duration: ~25 seconds

### Scene 2a: FSH/LH Comparison
- Data table + bar chart
- Compares M (high), N (low), Normal (baseline)
- Duration: ~35 seconds

### Scene 2b: Hypotheses
- Woman M: Ovarian failure vs pituitary hyperactivity
- Woman N: Pituitary insufficiency vs hypothalamus problem
- Duration: ~40 seconds

### Scene 3: Ovarian Structure
- Ultrasound comparison showing follicles
- M: No follicles (ovarian failure confirmed)
- N: Follicles present (ovaries intact)
- Duration: ~30 seconds

### Scene 4: GnRH Test
- Injection response graph
- Woman N responds  ->  pituitary functional  ->  hypothalamus problem
- Duration: ~35 seconds

### Final Summary: Complete Diagnosis
- Woman M: Ovarian failure
- Woman N: Hypothalamus insufficiency
- Duration: ~25 seconds

**Total Duration:** ~3 minutes 40 seconds

---

## 🚀 Generate Animations

### Install Manim (first time only)
```bash
pip install manim
```

### Generate All Scenes
```bash
cd ~/Documents/bac/resources/notes/Study-Vault/Animations/Biology
manim -pql Bio_Ex1_HormoneRegulation.py
```

### Generate Single Scene
```bash
# Scene 1a only
manim -pql Bio_Ex1_HormoneRegulation.py Scene1a_OvarianHormones

# Scene 2a only
manim -pql Bio_Ex1_HormoneRegulation.py Scene2a_FSH_LH_Comparison

# Final summary only
manim -pql Bio_Ex1_HormoneRegulation.py FinalSummary
```

### Quality Options
```bash
-ql  # Low (480p) - Fast for testing
-qm  # Medium (720p)
-qh  # High (1080p) - Best for studying
```

---

## 📂 Output Location

Videos saved to:
```
Animations/Biology/media/videos/Bio_Ex1_HormoneRegulation/480p15/
├── Scene1a_OvarianHormones.mp4
├── Scene1b_NoMenstruation.mp4
├── Scene2a_FSH_LH_Comparison.mp4
├── Scene2b_Hypotheses.mp4
├── Scene3_OvarianStructure.mp4
├── Scene4_GnRH_Test.mp4
└── FinalSummary.mp4
```

---

## 🎯 Study Workflow

1. **Watch Scene 1a-1b**  ->  Understand the problem
2. **Watch Scene 2a-2b**  ->  Analyze hormone data
3. **Watch Scene 3**  ->  Interpret ultrasound
4. **Watch Scene 4**  ->  Understand diagnostic test
5. **Watch Final Summary**  ->  Review complete diagnosis
6. **Try solving exercise yourself**
7. **Rewatch specific scenes** as needed

---

## 🔧 Customization

Edit `Bio_Ex1_HormoneRegulation.py` to:
- Change colors
- Adjust timing (run_time parameter)
- Add more details
- Modify text size

---

## ✅ Next Exercises

After this works:
- Bio_Ex2: Synaptic Transmission (8 scenes)
- Bio_Ex3: Immune Response (5 scenes)
- Bio_Ex4: Drosophila Genetics (5 scenes)

Same format, same detail level!
