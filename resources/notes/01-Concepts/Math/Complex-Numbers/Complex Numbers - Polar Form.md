---
tags: [concept, math, complex-numbers, polar-form]
subject: Mathematics
topic: Complex Numbers
difficulty: Medium
---

# Complex Numbers - Polar Form

## Definition

The **polar form** of a complex number expresses it in terms of its distance from the origin and its angle with the positive real axis.

## Forms

### Trigonometric Form
$$z = r(\cos \theta + i \sin \theta)$$

### Euler's Form
$$z = re^{i\theta}$$

where:
- r = |z| = modulus (distance from origin)
- θ = arg(z) = argument (angle with positive real axis)

## Conversion from Cartesian to Polar

Given z = a + bi:
1. Calculate modulus: $r = |z| = \sqrt{a^2 + b^2}$
2. Calculate argument: $\theta = \tan^{-1}(b/a)$ (adjust for correct quadrant)

## Conversion from Polar to Cartesian

Given z = r(cos θ + i sin θ):
- Real part: a = r cos θ
- Imaginary part: b = r sin θ

## Advantages of Polar Form

1. **Multiplication** becomes addition of arguments:
   $$z_1 z_2 = r_1 r_2 e^{i(\theta_1 + \theta_2)}$$

2. **Division** becomes subtraction of arguments:
   $$\frac{z_1}{z_2} = \frac{r_1}{r_2} e^{i(\theta_1 - \theta_2)}$$

3. **Powers** become simpler with De Moivre's theorem:
   $$z^n = r^n e^{in\theta}$$

## Examples

### Example 1
Convert z = 1 + √3i to polar form
Solution:
$$r = \sqrt{1^2 + (\sqrt{3})^2} = 2$$
$$\theta = \tan^{-1}(\sqrt{3}/1) = \pi/3$$
$$z = 2(\cos(\pi/3) + i\sin(\pi/3)) = 2e^{i\pi/3}$$

### Example 2
Multiply z₁ = 2e^{iπ/4} and z₂ = 3e^{iπ/6}
Solution:
$$z_1 z_2 = 2 \cdot 3 \cdot e^{i(\pi/4 + \pi/6)} = 6e^{i(5\pi/12)}$$

## Applications

- Alternating current circuits
- Wave interference
- Rotation transformations
- Complex analysis

## Visualizations

![[Complex Numbers - Polar Form.excalidraw]]

---
Back to: [[Complex Numbers MOC]] | [[Math MOC]]