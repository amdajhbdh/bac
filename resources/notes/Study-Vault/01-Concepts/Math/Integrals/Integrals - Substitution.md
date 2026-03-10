---
tags: [math, calculus, integrals, u-substitution]
difficulty: medium
---
> [!summary] 📊 Note Summary
> 
> | Property | Value |
> |----------|-------|
> | **Difficulty** | `hard` #difficulty/hard |
> | **Formulas** | 0 |
> | **Concepts** | 0 |
> | **Related Notes** | 10 |
> | **Word Count** | 681 |
> | **Last Enhanced** | 2026-03-10 |



## 📊 Note Summary

| Property | Value |
|----------|-------|
| Difficulty | Hard |
| Formulas | 0 |
| Concepts | 0 |
| Related Notes | 10 |
| Word Count | 606 |
| Last Enhanced | 2026-03-10 |



# Integration by Substitution (u-Substitution)

## The Method
Used when integrand contains a function and its derivative.

**Strategy**: Let u = g(x), then du = g'(x)dx

## Steps
1. Identify inner function u = g(x)
2. Find du = g'(x)dx
3. Rewrite integral in terms of u
4. Integrate with respect to u
5. Substitute back to x

## Basic Examples

### Example 1: Power of Function
∫ 2x(x**2 + 1)**5 dx

Let u = x**2 + 1
du = 2x dx

∫ u**5 du = u**6/6 + C = (x**2 + 1)**6/6 + C

### Example 2: Exponential
∫ 3x**2e**(x**3) dx

Let u = x**3
du = 3x**2 dx

∫ e**u du = e**u + C = e**(x**3) + C

### Example 3: Trigonometric
∫ sin(x)cos(x) dx

Let u = sin(x)
du = cos(x) dx

∫ u du = u**2/2 + C = sin**2(x)/2 + C

## Definite Integrals with Substitution

### Method 1: Change Limits
∫[a to b] f(g(x))g'(x) dx

Let u = g(x), change limits:
- When x = a, u = g(a)
- When x = b, u = g(b)

∫[g(a) to g(b)] f(u) du

### Example 4
∫[0 to 1] 2x(x**2 + 1)**3 dx

Let u = x**2 + 1, du = 2x dx
When x = 0: u = 1
When x = 1: u = 2

∫[1 to 2] u**3 du = [u**4/4] from 1 to 2
= 16/4 - 1/4 = 15/4

## Advanced Techniques

### Adjusting Constants
∫ 6x(x**2 + 1)**5 dx

Need 2x, have 6x  ->  multiply by 1/3

= 3∫ 2x(x**2 + 1)**5 dx
= 3(x**2 + 1)**6/6 + C
= (x**2 + 1)**6/2 + C

### Example 5: Rational Function
∫ x/(x**2 + 1) dx

Let u = x**2 + 1, du = 2x dx
Need 2x, have x  ->  multiply by 1/2

= (1/2)∫ 2x/(x**2 + 1) dx
= (1/2)∫ 1/u du
= (1/2)ln|u| + C
= (1/2)ln(x**2 + 1) + C

### Example 6: Square Root
∫ x√(x**2 + 4) dx

Let u = x**2 + 4, du = 2x dx

= (1/2)∫ 2x√(x**2 + 4) dx
= (1/2)∫ √u du
= (1/2)(2u**(3/2)/3) + C
= u**(3/2)/3 + C
= (x**2 + 4)**(3/2)/3 + C

## Trigonometric Substitution Preview

### For √(a**2 - x**2): Let x = a sin(θ)
### For √(a**2 + x**2): Let x = a tan(θ)
### For √(x**2 - a**2): Let x = a sec(θ)

## Common Patterns

| Integrand Form | Substitution |
|----------------|--------------|
| f(ax + b) | u = ax + b |
| x·f(x**2) | u = x**2 |
| x**2·f(x**3) | u = x**3 |
| f(x)/x | u = ln(x) |
| e**x·f(e**x) | u = e**x |
| sin(x)·f(cos(x)) | u = cos(x) |
| cos(x)·f(sin(x)) | u = sin(x) |
| sec**2(x)·f(tan(x)) | u = tan(x) |

## Practice Problems

### Problem 1
∫ 4x**3(x**4 + 2)**7 dx

<details>
<summary>Solution</summary>
u = x**4 + 2, du = 4x**3 dx
∫ u**7 du = u**8/8 + C = (x**4 + 2)**8/8 + C
</details>

### Problem 2
∫[0 to π/4] tan(x)sec**2(x) dx

<details>
<summary>Solution</summary>
u = tan(x), du = sec**2(x) dx
When x = 0: u = 0
When x = π/4: u = 1

∫[0 to 1] u du = [u**2/2] = 1/2 - 0 = 1/2
</details>

### Problem 3
∫ (ln(x))/x dx

<details>
<summary>Solution</summary>
u = ln(x), du = (1/x) dx
∫ u du = u**2/2 + C = (ln(x))**2/2 + C
</details>

## Related
- [[Integrals - Fundamentals]]
- [[Integrals - By Parts]]
- [[Integrals - Trig Substitution]]


## 🔗 Related Notes

- [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
- [[Resource Links.md|Resource Links]]
- [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
- [[00-Meta/VAULT-SUMMARY.md|VAULT-SUMMARY]]
- [[00-Meta/DEEP-CONTENT-STATUS.md|DEEP-CONTENT-STATUS]]
- [[00-Meta/MOCs/Math MOC.md|Math MOC]]
- [[00-Meta/VAULT-SUMMARY.md|VAULT-SUMMARY]]
- [[00-Meta/DEEP-CONTENT-STATUS.md|DEEP-CONTENT-STATUS]]
- [[00-Meta/MOCs/Math MOC.md|Math MOC]]
- [[Resource Links.md|Resource Links]]


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
