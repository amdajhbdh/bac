---
tags: [math, calculus, integrals, integration-by-parts]
difficulty: hard
---
> [!summary] 📊 Note Summary
> 
> | Property | Value |
> |----------|-------|
> | **Difficulty** | `easy` #difficulty/easy |
> | **Formulas** | 0 |
> | **Concepts** | 0 |
> | **Related Notes** | 10 |
> | **Word Count** | 842 |
> | **Last Enhanced** | 2026-03-10 |



## 📊 Note Summary

| Property      | Value      |
| ------------- | ---------- |
| Difficulty    | Easy       |
| Formulas      | 0          |
| Concepts      | 0          |
| Related Notes | 10         |
| Word Count    | 771        |
| Last Enhanced | 2026-03-10 |



# Integration by Parts

## The Formula
$$ \integ

Derived from product rule: d(uv) = u dv + v du

## LIATE Rule (Choosing u)
Priority order for selecting u:
1. **L**ogarithmic: ln(x), log(x)
2. **I**nverse trig: arcsin(x), arctan(x)
3. **A**lgebraic: x, x**2, polynomials
4. **T**rigonometric: sin(x), cos(x)
5. **E**xponential: e**x, a**x

Choose u from highest priority, dv gets the rest.

## Basic Examples

### Example 1: Polynomial * Exponential
∫ x·e**x dx

u = x  ->  du = dx
dv = e**x dx  ->  v = e**x

∫ x·e**x dx = x·e**x - ∫ e**x dx
= x·e**x - e**x + C
= e**x(x - 1) + C

### Example 2: Polynomial * Trig
∫ x·sin(x) dx

u = x  ->  du = dx
dv = sin(x) dx  ->  v = -cos(x)

= x·(-cos(x)) - ∫ (-cos(x)) dx
= -x·cos(x) + sin(x) + C

### Example 3: Logarithm
∫ ln(x) dx

u = ln(x)  ->  du = (1/x) dx
dv = dx  ->  v = x

= x·ln(x) - ∫ x·(1/x) dx
= x·ln(x) - ∫ 1 dx
= x·ln(x) - x + C
= x(ln(x) - 1) + C

## Multiple Applications

### Example 4: x**2e**x
∫ x**2·e**x dx

**First application:**
u = x**2  ->  du = 2x dx
dv = e**x dx  ->  v = e**x

= x**2·e**x - ∫ 2x·e**x dx

**Second application on ∫ 2x·e**x dx:**
u = 2x  ->  du = 2 dx
dv = e**x dx  ->  v = e**x

= 2x·e**x - ∫ 2·e**x dx
= 2x·e**x - 2e**x

**Final answer:**
= x**2·e**x - (2x·e**x - 2e**x) + C
= x**2·e**x - 2x·e**x + 2e**x + C
= e**x(x**2 - 2x + 2) + C

### Example 5: x**2sin(x)
∫ x**2·sin(x) dx

**First:**
u = x**2, dv = sin(x) dx
du = 2x dx, v = -cos(x)

= -x**2·cos(x) + ∫ 2x·cos(x) dx

**Second:**
u = 2x, dv = cos(x) dx
du = 2 dx, v = sin(x)

= 2x·sin(x) - ∫ 2·sin(x) dx
= 2x·sin(x) + 2cos(x)

**Final:**
= -x**2·cos(x) + 2x·sin(x) + 2cos(x) + C

## Circular Method

### Example 6: e**x·sin(x)
∫ e**x·sin(x) dx

Let I = ∫ e**x·sin(x) dx

u = sin(x), dv = e**x dx
du = cos(x) dx, v = e**x

I = e**x·sin(x) - ∫ e**x·cos(x) dx

For ∫ e**x·cos(x) dx:
u = cos(x), dv = e**x dx
du = -sin(x) dx, v = e**x

= e**x·cos(x) + ∫ e**x·sin(x) dx
= e**x·cos(x) + I

Substitute back:
I = e**x·sin(x) - (e**x·cos(x) + I)
I = e**x·sin(x) - e**x·cos(x) - I
2I = e**x(sin(x) - cos(x))
I = (e**x/2)(sin(x) - cos(x)) + C

## Definite Integrals

### Example 7
∫[1 to e] ln(x) dx

u = ln(x), dv = dx
du = (1/x) dx, v = x

= [x·ln(x) - x] from 1 to e
= (e·ln(e) - e) - (1·ln(1) - 1)
= (e·1 - e) - (0 - 1)
= 0 + 1 = 1

## Tabular Method (for polynomials)

For ∫ P(x)·f(x) dx where P is polynomial:

| Derivative | Integral | Sign |
|------------|----------|------|
| P(x) | f(x) | + |
| P'(x) | ∫f(x) | - |
| P''(x) | ∫∫f(x) | + |
| ... | ... | ... |

### Example 8: x**3e**x
| D | I | Sign |
|---|---|------|
| x**3 | e**x | + |
| 3x**2 | e**x | - |
| 6x | e**x | + |
| 6 | e**x | - |
| 0 | e**x | + |

Result: x**3e**x - 3x**2e**x + 6xe**x - 6e**x + C
= e**x(x**3 - 3x**2 + 6x - 6) + C

## Special Cases

### Reduction Formulas
∫ x**n·e**x dx = x**n·e**x - n∫ x**n**{-1}·e**x dx

∫ x**n·sin(x) dx = -x**n·cos(x) + n∫ x**n**{-1}·cos(x) dx

## Practice Problems

### Problem 1
∫ x·cos(x) dx

<details>
<summary>Solution</summary>
u = x, dv = cos(x) dx
du = dx, v = sin(x)
= x·sin(x) - ∫ sin(x) dx
= x·sin(x) + cos(x) + C
</details>

### Problem 2
∫ x**2·ln(x) dx

<details>
<summary>Solution</summary>
u = ln(x), dv = x**2 dx
du = (1/x) dx, v = x**3/3
= (x**3/3)·ln(x) - ∫ (x**3/3)·(1/x) dx
= (x**3/3)·ln(x) - (1/3)∫ x**2 dx
= (x**3/3)·ln(x) - x**3/9 + C
</details>

### Problem 3
∫ e**x·cos(x) dx

<details>
<summary>Solution</summary>
Let I = ∫ e**x·cos(x) dx
Apply twice, get: 2I = e**x(cos(x) + sin(x))
I = (e**x/2)(cos(x) + sin(x)) + C
</details>

## Related
- [[Integrals - Substitution]]
- [[Integrals - Partial Fractions]]
- [[Integrals - Applications]]


## 🔗 Related Notes

- [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
- [[Resource Links.md|Resource Links]]
- [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
- [[Animations/ANIMATION_SPEC.md|ANIMATION_SPEC]]
- [[Animations/Biology/README.md|README]]
- [[00-Meta/VAULT-SUMMARY.md|VAULT-SUMMARY]]
- [[Resource Links.md|Resource Links]]
- [[ANIMATION-SYSTEM-COMPLETE.md|ANIMATION-SYSTEM-COMPLETE]]
- [[Animations/Biology/README.md|README]]
- [[QUICK-REFERENCE.md|QUICK-REFERENCE]]


> [!related] 🔗 Related Notes
> 
> - [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
> - [[Resource Links.md|Resource Links]]
> - [[ANIMATION-SYSTEM-COMPLETE.md|ANIMATION-SYSTEM-COMPLETE]]
> - [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
> - [[ANIMATION-SYSTEM-COMPLETE.md|ANIMATION-SYSTEM-COMPLETE]]
> - [[Animations/ALL-EXERCISES-COVERED.md|ALL-EXERCISES-COVERED]]
> - [[00-Meta/DEEP-CONTENT-STATUS.md|DEEP-CONTENT-STATUS]]
> - [[00-Meta/MOCs/Chemistry MOC.md|Chemistry MOC]]
> - [[01-Concepts/Math/Complex-Numbers/Complex Numbers - Operations.md|Complex Numbers - Operations]]
> - [[Animations/ANIMATION_SPEC.md|ANIMATION_SPEC]]
