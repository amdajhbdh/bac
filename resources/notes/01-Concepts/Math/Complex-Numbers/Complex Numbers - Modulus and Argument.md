---
tags: [concept, math, complex-numbers, modulus, argument]
subject: Mathematics
topic: Complex Numbers
difficulty: Medium
---

# Complex Numbers - Modulus and Argument

## Modulus (Absolute Value)

The **modulus** of a complex number z = a + bi is defined as:
$$|z| = \sqrt{a^2 + b^2}$$

It represents the distance from the origin to the point (a,b) in the complex plane.

## Argument

The **argument** of a complex number z = a + bi is the angle θ that the vector makes with the positive real axis:
$$\theta = \arg(z) = \tan^{-1}\left(\frac{b}{a}\right)$$

The principal argument is usually taken in the interval (-π, π].

## Polar Form

Using modulus and argument, any complex number can be written in polar form:
$$z = r(\cos \theta + i \sin \theta) = r e^{i\theta}$$

where r = |z| and θ = arg(z).

## Properties

1. $|z_1 z_2| = |z_1| |z_2|$
2. $|z_1/z_2| = |z_1|/|z_2|$ (when z₂ ≠ 0)
3. $|\overline{z}| = |z|$
4. $|z|^2 = z \overline{z}$

## Examples

### Example 1
Find modulus and argument of z = 1 + i
Solution:
$$|z| = \sqrt{1^2 + 1^2} = \sqrt{2}$$
$$\arg(z) = \tan^{-1}(1/1) = \pi/4$$

### Example 2
Convert z = 2 + 2i to polar form
Solution:
$$|z| = \sqrt{4 + 4} = 2\sqrt{2}$$
$$\arg(z) = \tan^{-1}(1) = \pi/4$$
$$z = 2\sqrt{2}(\cos(\pi/4) + i\sin(\pi/4))$$

## Applications

- Electrical engineering (phasors)
- Signal processing
- Quantum mechanics
- Geometry problems

## Visualizations

![[Complex Numbers - Modulus and Argument.excalidraw]]

---
Back to: [[Complex Numbers MOC]] | [[Math MOC]]