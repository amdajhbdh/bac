---
tags: [concept, physics, dynamics, forces]
subject: Physics
topic: Dynamics - Forces
difficulty: Medium
---

# Dynamics - Forces

## Introduction

Forces are fundamental interactions that cause objects to accelerate, decelerate, or deform. Understanding forces is essential for analyzing all motion in mechanics.

## Force as a Vector

Force is a vector quantity with both magnitude and direction:

$$\vec{F} = F_x \hat{i} + F_y \hat{j} + F_z \hat{k}$$

The SI unit of force is the Newton (N).

$$1 \text{ N} = 1 \text{ kg} \cdot \text{m/s}^2$$

## Types of Forces

### Contact Forces
- **Normal force** ($N$): Perpendicular to surfaces in contact
- **Friction** ($f$): Parallel to surfaces, opposes motion
- **Tension** ($T$): Pulling force in ropes/strings
- **Applied force** ($F$): External force applied to object

### Field Forces (Non-contact)
- **Gravity** ($mg$): Force due to gravitational field
- **Electric force**: Between charged particles
- **Magnetic force**: Between magnetic poles or moving charges

## Newton's Second Law

The acceleration of an object is directly proportional to the net force and inversely proportional to its mass:

$$\vec{F}_{net} = m\vec{a}$$

Or in component form:
$$F_x = ma_x, \quad F_y = ma_y, \quad F_z = ma_z$$

## Free Body Diagrams

A free body diagram (FBD) shows all forces acting on an object:

1. **Identify the object** to analyze
2. **Draw the object** as a point or shape
3. **Show all forces** as vectors with their directions
4. **Label each force** with its magnitude or symbol
5. **Choose coordinate system** (typically x horizontal, y vertical)

## Common Force Calculations

### Weight (Gravitational Force)

$$W = mg$$

Where:
- $W$ = weight (N)
- $m$ = mass (kg)
- $g$ = gravitational acceleration ($\approx 9.8 \text{ m/s}^2$ on Earth)

### Normal Force

On a horizontal surface: $N = mg$

On an inclined plane: $N = mg\cos\theta$

### Tension in Ropes

For an ideal massless rope:
- Same tension throughout the rope
- Same magnitude on both ends
- Direction pulls away from the object

## Equilibrium

An object is in equilibrium when:
$$\vec{F}_{net} = 0$$

This means:
- **Static equilibrium**: Object at rest ($v = 0$)
- **Dynamic equilibrium**: Object moving at constant velocity ($a = 0$)

### Conditions for Equilibrium

In 2D:
$$\sum F_x = 0, \quad \sum F_y = 0$$

In 3D:
$$\sum F_x = 0, \quad \sum F_y = 0, \quad \sum F_z = 0$$

## Examples

### Example 1: Object on Horizontal Surface

A 5 kg object rests on a horizontal surface. Find the normal force.

**Solution:**
$$N = mg = 5 \times 9.8 = 49 \text{ N}$$

---

### Example 2: Object on Inclined Plane

A 10 kg object rests on a 30° incline. Find the normal force and the component of weight parallel to the incline.

**Solution:**

Normal force: $N = mg\cos\theta = 10 \times 9.8 \times \cos(30°) = 98 \times 0.866 = 84.9 \text{ N}$

Parallel component: $F_{parallel} = mg\sin\theta = 10 \times 9.8 \times \sin(30°) = 98 \times 0.5 = 49 \text{ N}$

---

### Example 3: Two-Hanging Masses

A 3 kg mass hangs from a rope over a frictionless pulley. Find the tension in the rope and the acceleration of the system.

**Solution:**

For the hanging mass: $mg - T = ma$

Since $a = g$ (free fall):
$$T = m(g - a) = 0$$

Actually, the system accelerates: $a = \frac{m_2 - m_1}{m_1 + m_2}g$

With just one mass: $T = 0$ when in free fall...

Wait - a single mass over pulley: The mass accelerates down with $a = g$.

Actually, for a single pulley with one mass:
$$mg - T = ma$$
With $a = g$: $T = 0$

This isn't interesting. Let's do two masses: 3 kg and 5 kg.

$$a = \frac{5-3}{5+3}g = \frac{2}{8} \times 9.8 = 2.45 \text{ m/s}^2$$

For the 5 kg mass (moving down):
$$5g - T = 5a \Rightarrow T = 5(g - a) = 5(9.8 - 2.45) = 36.75 \text{ N}$$

---

### Example 4: Pulling at an Angle

A 20 kg box is pulled across a horizontal floor with a force of 100 N at 30° above horizontal. The coefficient of kinetic friction is 0.3. Find the acceleration.

**Solution:**

Horizontal component: $F_x = 100\cos(30°) = 86.6 \text{ N}$

Vertical component: $F_y = 100\sin(30°) = 50 \text{ N}$

Normal force: $N = mg - F_y = 20 \times 9.8 - 50 = 196 - 50 = 146 \text{ N}$

Friction force: $f_k = \mu_k N = 0.3 \times 146 = 43.8 \text{ N}$

Net horizontal force: $F_{net} = 86.6 - 43.8 = 42.8 \text{ N}$

Acceleration: $a = \frac{F_{net}}{m} = \frac{42.8}{20} = 2.14 \text{ m/s}^2$

---

### Example 5: Atwood Machine

Two masses $m_1 = 2 \text{ kg}$ and $m_2 = 3 \text{ kg}$ are connected by a light rope over a frictionless pulley. Find the acceleration and tension.

**Solution:**

For $m_1$ (assuming $m_2 > m_1$ and moves up):
$$T - m_1g = m_1 a$$

For $m_2$ (moves down):
$$m_2g - T = m_2 a$$

Adding equations:
$$m_2g - m_1g = (m_1 + m_2)a$$
$$a = \frac{m_2 - m_1}{m_1 + m_2}g = \frac{3-2}{2+3} \times 9.8 = \frac{1}{5} \times 9.8 = 1.96 \text{ m/s}^2$$

From $T - m_1g = m_1a$:
$$T = m_1(g + a) = 2(9.8 + 1.96) = 23.5 \text{ N}$$

---

## Force Summary

| Force Type | Formula | Direction |
|------------|---------|-----------|
| Weight | $W = mg$ | Downward |
| Normal | $N$ | Perpendicular to surface |
| Friction | $f = \mu N$ | Opposes motion |
| Tension | $T$ | Along rope, pulls |
| Spring | $F = -kx$ | Opposes displacement |

---

Back to: [[Dynamics MOC]] | [[Physics MOC]]
