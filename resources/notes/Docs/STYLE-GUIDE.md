# 🎨 VAULT STYLE GUIDE
## Mauritanian BAC Notes

---

## 📝 Frontmatter (Required)

All concept files should have this frontmatter:

```yaml
---
tags: [concept, math, algebra]
subject: Mathematics
topic: Algebra
level: Terminale
difficulty: easy|medium|hard
created: 2026-01-15
updated: 2026-03-14
---
```

### Tag Conventions
- `concept` - Learning content
- `practice` - Exercises
- `exam` - Past exam questions
- `formula` - Formula summary
- `summary` - Quick review

---

## 📐 File Naming

| Type | Convention | Example |
|------|------------|---------|
| Concepts | `Topic - Subtopic.md` | `Complex Numbers - Operations.md` |
| Practice | `Topic-Practice-N.md` | `Algebra-Practice-1.md` |
| Exams | `BAC-Year-Serie-Topic.md` | `BAC-2025-D-SN-Ex1-Arithmetique.md` |
| Daily | `YYYY-MM-DD-Title.md` | `2026-03-14-Study-Notes.md` |

---

## 🔢 Math Notation (LaTeX)

### Inline Math
```markdown
The quadratic formula is $x = \frac{-b \pm \sqrt{b^2-4ac}}{2a}$
```

### Display Math
```markdown
$$
\int_0^1 x^2 dx = \left[\frac{x^3}{3}\right]_0^1 = \frac{1}{3}
$$
```

---

## 📢 Callouts

Use these Obsidian callouts:

```markdown
::: tip 💡 Remarque
Important concept to remember!
:::

::: warning ⚠️ Attention
Common mistake to avoid!
:::

::: note 📝 Note
Additional information.
:::

::: example 📋 Exemple
Example problem with solution.
:::

::: definition 📚 Définition
Clear definition of a concept.
:::
```

---

## 🔗 Links

### Internal Links
```markdown
Back to: [[Complex Numbers MOC]] | [[Math MOC]]

See also: [[Algebra - Equations]]
```

### Navigation
```markdown
← Previous: [[Topic A]]
→ Next: [[Topic C]]
```

---

## 📊 Tables

```markdown
| Symbol | Meaning | Example |
|--------|---------|---------|
| $\mathbb{R}$ | Real numbers | $x \in \mathbb{R}$ |
| $\mathbb{C}$ | Complex numbers | $z \in \mathbb{C}$ |
```

---

## 🎯 Headers Structure

```markdown
# Topic Name (H1 - File Title)

## Section 1 (H2)

### Subsection (H3)

#### Detail (H4)
```

---

## 🇫🇷 🇬🇧 Bilingual Notes

When covering concepts in both French and English:

```markdown
## Définition / Definition

**Fr:** La dérivée mesure le taux de variation instantané.
**En:** The derivative measures instantaneous rate of change.
```

---

## ✅ Quick Checklist

- [ ] Frontmatter with tags
- [ ] Consistent file naming
- [ ] Math in LaTeX ($...$)
- [ ] Links to MOCs
- [ ] Callouts for important info
- [ ] French + English where applicable

---

*Last updated: 2026-03-14*
