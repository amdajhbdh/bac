---
project: BAC Preparation 2026
type: project
status: active
start-date: 2026-03-01
target-date: 2026-06-15
---

# 🎓 BAC Preparation 2026

> [!info] Project Overview
> **Goal:** Master all 7C/7D curriculum topics for Mauritanian Baccalaureate
> 
> **Timeline:** 3.5 months (March - June 2026)
> 
> **Resources:** 11 textbooks, 1,571 exercises, 686+ notes

---

## 📊 Progress Dashboard

```dataview
TABLE 
  difficulty as "Difficulty",
  status as "Status",
  progress as "Progress %",
  next-review as "Next Review"
FROM "01-Concepts"
WHERE project = "BAC Preparation 2026"
SORT difficulty ASC, status ASC
```

---

> [!stats] 📈 Statistics
> 
> | Metric | Value |
> |--------|-------|
> | Total Topics | 686+ |
> | Completed | 0 |
> | In Progress | 100 |
> | Not Started | 586 |
> | Exercises Done | 0/1571 |
> | Study Streak | 0 days |

---

> [!todo] 🎯 Current Focus
> 
> ### This Week
> - [ ] Review matrices (Math)
> - [ ] Practice mechanics problems (Physics)
> - [ ] Study organic chemistry (Chemistry)
> 
> ### This Month
> - [ ] Complete all Math chapters
> - [ ] Master Physics fundamentals
> - [ ] Finish Chemistry basics

---

## 📚 By Subject

### Mathematics (490 chapters)
```dataview
TABLE status, difficulty, progress
FROM "01-Concepts/Math" OR "05-Extracted/Resources/Math"
WHERE project = "BAC Preparation 2026"
LIMIT 10
```

### Physics (275 chapters)
```dataview
TABLE status, difficulty, progress
FROM "01-Concepts/Physics" OR "05-Extracted/Resources/Physics"
WHERE project = "BAC Preparation 2026"
LIMIT 10
```

### Chemistry (176 chapters)
```dataview
TABLE status, difficulty, progress
FROM "01-Concepts/Chemistry" OR "05-Extracted/Resources/Chemistry"
WHERE project = "BAC Preparation 2026"
LIMIT 10
```

---

> [!warning] 🔥 Weak Areas
> 
> ```dataview
> TABLE difficulty, progress, last-reviewed
> FROM "01-Concepts"
> WHERE progress < 50 AND difficulty = "hard"
> SORT progress ASC
> LIMIT 10
> ```

---

> [!success] ⭐ Mastered Topics
> 
> ```dataview
> TABLE difficulty, progress, mastery-date
> FROM "01-Concepts"
> WHERE progress >= 90
> SORT mastery-date DESC
> LIMIT 10
> ```

---

> [!note] 📅 Study Schedule
> 
> | Day | Focus | Hours | Status |
> |-----|-------|-------|--------|
> | Monday | Math | 3 | ⏳ |
> | Tuesday | Physics | 3 | ⏳ |
> | Wednesday | Chemistry | 3 | ⏳ |
> | Thursday | Math | 3 | ⏳ |
> | Friday | Physics | 3 | ⏳ |
> | Saturday | Review | 4 | ⏳ |
> | Sunday | Practice | 4 | ⏳ |
> 
> **Total:** 23 hours/week

---

> [!example] 🎓 Exam Preparation
> 
> ### BAC Papers (2002-2012)
> - [ ] Math papers (10 years)
> - [ ] Physics papers (10 years)
> - [ ] Chemistry papers (10 years)
> 
> ### Mock Exams
> - [ ] Full mock exam 1
> - [ ] Full mock exam 2
> - [ ] Full mock exam 3

---

> [!tip] 📝 Study Tips
> - Focus on weak areas first
> - Practice daily with spaced repetition
> - Review past BAC papers weekly
> - Track progress in dashboard
> - Use Smart Connections to discover related topics

