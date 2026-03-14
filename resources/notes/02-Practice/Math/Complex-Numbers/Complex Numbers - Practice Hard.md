---
tags: [practice, math, complex-numbers, hard]
difficulty: Hard
total_problems: 5
---

# Complex Numbers - Practice Hard

← Back to [[Complex Numbers MOC]]

## Instructions
- Time limit: 45 minutes
- Show all work
- Check answers at bottom

---

## Problems

### Problem 1
Find all values of z such that z³ = -8.

### Problem 2
Prove that for any complex number z, |zⁿ| = |z|ⁿ for any positive integer n.

### Problem 3
Let z₁, z₂ be complex numbers such that |z₁| = |z₂| = 1 and z₁ + z₂ = 1. Find z₁ and z₂.

### Problem 4
Find the sum of all solutions to the equation (z + 1)⁴ = 16z⁴.

### Problem 5
Let z = cos(π/7) + i sin(π/7). Find the value of:
$$(1 + z)(1 + z^2)(1 + z^3)(1 + z^4)(1 + z^5)(1 + z^6)$$

---

## Answers

<details>
<summary>Click to reveal answers</summary>

### Problem 1 Solution
We need to find cube roots of -8.
-8 = 8e^{iπ} in polar form
Cube roots are of the form:
$$z_k = \sqrt[3]{8} e^{i(\pi + 2kπ)/3} = 2e^{i(π + 2kπ)/3}$$ for k = 0,1,2

For k = 0: $z_0 = 2e^{iπ/3} = 2(\cos(π/3) + i\sin(π/3)) = 1 + i\sqrt{3}$
For k = 1: $z_1 = 2e^{iπ} = -2$
For k = 2: $z_2 = 2e^{i5π/3} = 2(\cos(5π/3) + i\sin(5π/3)) = 1 - i\sqrt{3}$

Answer: $z = 1 + i\sqrt{3}, -2, 1 - i\sqrt{3}$

### Problem 2 Solution
Let z = re^{iθ} where r = |z|.
Then zⁿ = rⁿe^{inθ}, so |zⁿ| = |rⁿe^{inθ}| = rⁿ = |z|ⁿ.
Therefore |zⁿ| = |z|ⁿ.

### Problem 3 Solution
Let z₁ = e^{iα} and z₂ = e^{iβ} (since |z₁| = |z₂| = 1).
From z₁ + z₂ = 1:
e^{iα} + e^{iβ} = 1
(cos α + cos β) + i(sin α + sin β) = 1 + 0i

This gives us:
cos α + cos β = 1
sin α + sin β = 0

From the second equation: sin β = -sin α, which means β = -α or β = π + α.
Substituting into the first equation:
- Case 1: β = -α → cos α + cos(-α) = 2cos α = 1 → cos α = 1/2 → α = ±π/3
- Case 2: β = π + α → cos α + cos(π + α) = cos α - cos α = 0 ≠ 1 (contradiction)

Therefore α = π/3, β = -π/3 or α = -π/3, β = π/3.
So z₁ = 1/2 + i√3/2, z₂ = 1/2 - i√3/2 (or vice versa).

### Problem 4 Solution
Rewriting: (z + 1)⁴ = 16z⁴
Taking fourth root: z + 1 = 2ze^{ikπ/2} for k = 0,1,2,3

For k = 0: z + 1 = 2z → z = 1
For k = 1: z + 1 = 2ze^{iπ/2} = 2zi → z + 1 = 2zi → z(1 - 2i) = -1 → z = -1/(1-2i) = -(1+2i)/5
For k = 2: z + 1 = 2ze^{iπ} = -2z → z + 1 = -2z → 3z = -1 → z = -1/3
For k = 3: z + 1 = 2ze^{i3π/2} = -2zi → z + 1 = -2zi → z(1 + 2i) = -1 → z = -1/(1+2i) = -(1-2i)/5

Sum of solutions: 1 + (-1-2i)/5 + (-1/3) + (-1+2i)/5 = 1 - 1/3 - 2/5 = 2/3 - 2/5 = (10-6)/15 = 4/15

### Problem 5 Solution
Note that z⁷ = 1 since z = e^{iπ/7}.
We want to find $(1+z)(1+z²)...(1+z⁶)$. 

Consider the polynomial f(w) = w⁷ - 1 = (w-1)(w-z)(w-z²)...(w-z⁶)
Also, f(w) = (w-1)(1 + w + w² + ... + w⁶)

Therefore: (w-z)(w-z²)...(w-z⁶) = 1 + w + w² + ... + w⁶

Setting w = -1: (-1-z)(-1-z²)...(-1-z⁶) = 1 - 1 + 1 - 1 + 1 - 1 + 1 = 1

Since (-1-zᵏ) = -(1+zᵏ), we have:
(-1)⁶(1+z)(1+z²)...(1+z⁶) = 1
Therefore: (1+z)(1+z²)...(1+z⁶) = 1

Answer: 1

</details>

---

## Score Yourself
- 100%: Excellent! Master of complex numbers
- 80%: Very good! Minor review needed
- 60%: Good! Review key concepts
- <60%: Review polar forms, De Moivre's theorem and roots

---

## Related Practice
- [[Complex Numbers - Practice Easy]]
- [[Complex Numbers - Practice Medium]]

Back to: [[Complex Numbers MOC]] | [[Math MOC]]