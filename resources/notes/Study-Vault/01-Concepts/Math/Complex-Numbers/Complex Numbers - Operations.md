---
tags: [concept-note, math, complex-numbers, operations]
aliases: [Complex Arithmetic]
status: active
---
> [!summary] 📊 Note Summary
> 
> | Property | Value |
> |----------|-------|
> | **Difficulty** | `hard` #difficulty/hard |
> | **Formulas** | 15 |
> | **Concepts** | 3 |
> | **Related Notes** | 10 |
> | **Word Count** | 606 |
> | **Last Enhanced** | 2026-03-10 |



## 📊 Note Summary

| Property | Value |
|----------|-------|
| Difficulty | Hard |
| Formulas | 15 |
| Concepts | 3 |
| Related Notes | 10 |
| Word Count | 535 |
| Last Enhanced | 2026-03-10 |



# Complex Numbers - Operations

← [[Complex Numbers - Basics]] | [[Complex Numbers MOC]]

## ➕ Addition and Subtraction

Add/subtract real and imaginary parts separately:

$$(a + bi) + (c + di) = (a + c) + (b + d)i$$
$$(a + bi) - (c + di) = (a - c) + (b - d)i$$

### Example
$$(3 + 4i) + (2 - 5i) = (3+2) + (4-5)i = 5 - i$$

## ✖Multiplication

Use FOIL method and remember $i**2 = -1$:

$$(a + bi)(c + di) = ac + adi + bci + bdi**2$$
$$= ac + (ad + bc)i + bd(-1)$$
$$= (ac - bd) + (ad + bc)i$$

### Example
$$(2 + 3i)(1 + 4i) = 2 + 8i + 3i + 12i**2$$
$$= 2 + 11i + 12(-1) = -10 + 11i$$

## ➗ Division

Multiply numerator and denominator by the **conjugate** of denominator:

$$\frac{a + bi}{c + di} = \frac{(a + bi)(c - di)}{(c + di)(c - di)}$$

The denominator becomes real: $(c + di)(c - di) = c**2 + d**2$

### Example
$$\frac{3 + 4i}{1 - 2i} = \frac{(3 + 4i)(1 + 2i)}{(1 - 2i)(1 + 2i)}$$
$$= \frac{3 + 6i + 4i + 8i**2}{1 + 4} = \frac{3 + 10i - 8}{5}$$
$$= \frac{-5 + 10i}{5} = -1 + 2i$$

##Complex Conjugate

The conjugate of $z = a + bi$ is:
$$\bar{z} = a - bi$$

### Properties
- $z + \bar{z} = 2a$ (real)
- $z - \bar{z} = 2bi$ (imaginary)
- $z \cdot \bar{z} = a**2 + b**2$ (real, always positive)
- $\overline{z_1 + z_2} = \bar{z_1} + \bar{z_2}$
- $\overline{z_1 \cdot z_2} = \bar{z_1} \cdot \bar{z_2}$

##Modulus (Absolute Value)

$$|z| = |a + bi| = \sqrt{a**2 + b**2}$$

This is the distance from origin to point $(a, b)$ in complex plane.

### Properties
- $|z| \geq 0$, and $|z| = 0 \iff z = 0$
- $|z_1 \cdot z_2| = |z_1| \cdot |z_2|$
- $|z_1 + z_2| \leq |z_1| + |z_2|$ (triangle inequality)
- $|z| = |\bar{z}|$
- $z \cdot \bar{z} = |z|**2$

### Example
For $z = 3 + 4i$:
$$|z| = \sqrt{3**2 + 4**2} = \sqrt{9 + 16} = \sqrt{25} = 5$$

## Practice Problems

### Easy
1. $(5 + 2i) + (3 - 4i) = ?$
2. $(1 + i)(1 - i) = ?$
3. Find $|2 + 2i|$

### Medium
4. $\frac{2 + i}{3 - i} = ?$
5. If $z = 1 + 2i$, find $z**2$
6. Prove: $|z_1 \cdot z_2| = |z_1| \cdot |z_2|$

### Hard
7. Solve: $z**2 = 3 + 4i$
8. Find all $z$ where $|z - 1| = |z + 1|$

## Related Notes
- [[Complex Numbers - Conjugates]] - Deep dive into conjugates
- [[Complex Numbers - Modulus and Argument]] - Polar representation
- [[Complex Numbers - Polar Form]] - Alternative way to multiply/divide
- [[Complex Numbers - Practice Easy]] - More practice problems

##Video Resources
- [Khan Academy - Complex Number Operations](https://www.khanacademy.org/math/algebra2/x2ec2f6f830c9fb89:complex)
- [Organic Chemistry Tutor - Complex Numbers](https://www.youtube.com/watch?v=T647CGsuOVU)

---

← [[Complex Numbers - Basics]] | Next: [[Complex Numbers - Polar Form]]  -> 

Back to: [[Complex Numbers MOC]] | [[Math MOC]]


## 🔗 Related Notes

- [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
- [[Resource Links.md|Resource Links]]
- [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
- [[Animations/ALL-EXERCISES-COVERED.md|ALL-EXERCISES-COVERED]]
- [[00-Meta/QUICK-START.md|QUICK-START]]
- [[00-Meta/VAULT-SUMMARY.md|VAULT-SUMMARY]]
- [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
- [[Resource Links.md|Resource Links]]
- [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
- [[QUICK-REFERENCE.md|QUICK-REFERENCE]]


> [!related] 🔗 Related Notes
> 
> - [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
> - [[Resource Links.md|Resource Links]]
> - [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
> - [[Animations/ALL-EXERCISES-COVERED.md|ALL-EXERCISES-COVERED]]
> - [[00-Meta/QUICK-START.md|QUICK-START]]
> - [[00-Meta/VAULT-SUMMARY.md|VAULT-SUMMARY]]
> - [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
> - [[Resource Links.md|Resource Links]]
> - [[ANIMATION-SYSTEM-COMPLETE.md|ANIMATION-SYSTEM-COMPLETE]]
> - [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
