---
tags: [concept, math, complex-numbers, de-moivre]
subject: Mathematics
topic: Complex Numbers
difficulty: Hard
---

# Complex Numbers - De Moivre's Theorem

## Statement

**De Moivre's Theorem** states that for any complex number in polar form and any integer n:
$$(r(\cos \theta + i \sin \theta))^n = r^n(\cos(n\theta) + i \sin(n\theta))$$

In exponential form:
$$(re^{i\theta})^n = r^n e^{in\theta}$$

## Proof Outline

The theorem can be proven by mathematical induction:
1. Base case: n = 1 (trivially true)
2. Assume true for n = k
3. Prove true for n = k + 1 using trigonometric identities

## Applications

### Finding Powers of Complex Numbers

To compute zⁿ where z = a + bi:
1. Convert z to polar form
2. Apply De Moivre's theorem
3. Convert back to Cartesian form if needed

### Finding Roots of Complex Numbers

To find the nth roots of z = r(cos θ + i sin θ):
$$z_k = \sqrt[n]{r}\left(\cos\left(\frac{\theta + 2k\pi}{n}\right) + i \sin\left(\frac{\theta + 2k\pi}{n}\right)\right)$$

for k = 0, 1, 2, ..., n-1.

## Examples

### Example 1: Computing a Power
Find (1 + i)⁴

Solution:
1. Convert to polar: $1 + i = \sqrt{2}e^{i\pi/4}$
2. Apply De Moivre: $(\sqrt{2}e^{i\pi/4})^4 = (\sqrt{2})^4 e^{i\pi} = 4e^{i\pi} = -4$

### Example 2: Finding Cube Roots
Find all cube roots of 8.

Solution:
1. Express 8 as complex number: 8 + 0i
2. Polar form: r = 8, θ = 0
3. Cube roots formula:
   $$z_k = \sqrt[3]{8}\left(\cos\left(\frac{0 + 2k\pi}{3}\right) + i \sin\left(\frac{0 + 2k\pi}{3}\right)\right) = 2\left(\cos\left(\frac{2k\pi}{3}\right) + i \sin\left(\frac{2k\pi}{3}\right)\right)$$

   For k = 0: z₀ = 2(cos(0) + i sin(0)) = 2
   For k = 1: z₁ = 2(cos(2π/3) + i sin(2π/3)) = 2(-1/2 + i√3/2) = -1 + i√3
   For k = 2: z₂ = 2(cos(4π/3) + i sin(4π/3)) = 2(-1/2 - i√3/2) = -1 - i√3

## Connection to Euler's Formula

De Moivre's theorem is closely related to Euler's formula:
$$e^{i\theta} = \cos \theta + i \sin \theta$$

This makes the exponential form particularly useful:
$$(e^{i\theta})^n = e^{in\theta}$$

## Important Identities Derived from De Moivre's Theorem

1. **Multiple Angle Formulas**:
   $$\cos(n\theta) + i\sin(n\theta) = (\cos \theta + i \sin \theta)^n$$

2. **Chebyshev Polynomials**:
   Expanding the right side gives polynomials in cos θ.

## Advanced Applications

- Solving polynomial equations
- Series expansions
- Fourier analysis
- Quantum mechanics

## Visualizations

![[Complex Numbers - De Moivre Theorem.excalidraw]]

---
Back to: [[Complex Numbers MOC]] | [[Math MOC]]