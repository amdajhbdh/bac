---
tags: [concept, physics, dynamics, momentum, impulse]
subject: Physics
topic: Dynamics - Momentum
difficulty: Medium
---

# Dynamics - Momentum

## Introduction

Momentum is a fundamental concept in physics that describes the quantity of motion an object has. It is crucial for analyzing collisions, explosions, and overall motion in mechanics.

## Momentum

Momentum is the product of mass and velocity:

$$\vec{p} = m\vec{v}$$

Where:
- $p$ = momentum (kg·m/s)
- $m$ = mass (kg)
- $v$ = velocity (m/s)

Momentum is a **vector quantity** with the same direction as velocity.

## Impulse

Impulse is the change in momentum caused by a force applied over time:

$$\vec{J} = \vec{F}\Delta t = \Delta\vec{p}$$

Or in integral form:

$$\vec{J} = \int_{t_1}^{t_2} \vec{F}\,dt = \Delta\vec{p}$$

The impulse-momentum theorem states that the impulse on an object equals its change in momentum.

## Conservation of Momentum

For an isolated system (no external forces):

$$\vec{p}_{total, initial} = \vec{p}_{total, final}$$

$$\sum m_i \vec{v}_i = \text{constant}$$

This is one of the most fundamental conservation laws in physics.

## Types of Collisions

### Elastic Collisions

Both momentum and kinetic energy are conserved:

$$m_1v_{1i} + m_2v_{2i} = m_1v_{1f} + m_2v_{2f}$$

$$\frac{1}{2}m_1v_{1i}^2 + \frac{1}{2}m_2v_{2i}^2 = \frac{1}{2}m_1v_{1f}^2 + \frac{1}{2}m_2v_{2f}^2$$

### Inelastic Collisions

Momentum is conserved, kinetic energy is not:

$$m_1v_{1i} + m_2v_{2i} = (m_1 + m_2)v_f$$

**Perfectly inelastic**: Objects stick together after collision.

### Coefficient of Restitution

$$e = \frac{v_{2f} - v_{1f}}{v_{1i} - v_{2i}}$$

- $e = 1$: Perfectly elastic
- $e = 0$: Perfectly inelastic
- $0 < e < 1$: Partially elastic

## Center of Mass

The center of mass of a system is the weighted average position:

$$\vec{r}_{cm} = \frac{\sum m_i \vec{r}_i}{\sum m_i}$$

For continuous mass distribution:

$$\vec{r}_{cm} = \frac{1}{M}\int \vec{r}\,dm$$

The center of mass moves as if all external forces act on a single particle of mass M.

## Motion of Center of Mass

$$M\vec{a}_{cm} = \vec{F}_{external}$$

The center of mass accelerates to Newton's according second law with the total external force.

## Examples

### Example 1: Recoil of a Gun

A 0.05 kg bullet is fired from a 2 kg gun with velocity 500 m/s. Find the recoil velocity of the gun.

**Solution:**

Initial momentum = 0 (system at rest)

Final momentum: $0 = m_b v_b + m_g v_g$

$$v_g = -\frac{m_b v_b}{m_g} = -\frac{0.05 \times 500}{2} = -12.5 \text{ m/s}$$

The negative sign indicates the gun moves opposite to the bullet.

---

### Example 2: Elastic Collision

A 2 kg ball moving at 3 m/s collides elastically with a stationary 2 kg ball. Find their final velocities.

**Solution:**

For equal masses in elastic collision:
$$v_{1f} = v_{2i} = 0$$
$$v_{2f} = v_{1i} = 3 \text{ m/s}$$

The first ball stops, second ball takes all the momentum.

---

### Example 3: Inelastic Collision

A 1000 kg car moving at 20 m/s collides with a 1500 kg car moving at 10 m/s in the same direction. After collision, they stick together. Find their common velocity.

**Solution:**

$$m_1v_1 + m_2v_2 = (m_1 + m_2)v_f$$
$$1000 \times 20 + 1500 \times 10 = 2500 \times v_f$$
$$20000 + 15000 = 2500v_f$$
$$v_f = \frac{35000}{2500} = 14 \text{ m/s}$$

---

### Example 4: Impulse from Force

A 0.5 kg ball hits a wall at 10 m/s and rebounds at 8 m/s. The contact time is 0.01 s. Find the average force on the ball.

**Solution:**

Initial momentum: $p_i = 0.5 \times 10 = 5 \text{ kg·m/s}$
Final momentum: $p_f = 0.5 \times (-8) = -4 \text{ kg·m/s}$

Change in momentum: $\Delta p = -4 - 5 = -9 \text{ kg·m/s}$

$$F_{avg} = \frac{\Delta p}{\Delta t} = \frac{-9}{0.01} = -900 \text{ N}$$

Magnitude: 900 N (opposite to initial direction)

---

### Example 5: Explosion

A 2 kg object at rest explodes into two pieces. One piece (0.5 kg) flies off at 30 m/s. Find the velocity of the other piece.

**Solution:**

Initial momentum: $p_i = 0$

Final momentum: $0 = m_1v_1 + m_2v_2$

$$0 = 0.5 \times 30 + 1.5 \times v_2$$

$$v_2 = -\frac{15}{1.5} = -10 \text{ m/s}$$

The negative sign indicates opposite direction.

---

## Momentum in Multiple Dimensions

In 2D/3D, momentum is conserved separately in each direction:

$$p_{x,initial} = p_{x,final}$$
$$p_{y,initial} = p_{y,final}$$

### 2D Elastic Collision

For two equal masses where one is initially at rest:

- Scattering angle $\theta$
- Both final velocities at $90°$ to each other
- $v_{1f} = v \cos\theta$
- $v_{2f} = v \sin\theta$

## Rocket Propulsion

Rocket thrust comes from momentum conservation:

$$\frac{dm}{dt}v_e = F_{thrust}$$

Where:
- $dm/dt$ = mass ejection rate
- $v_e$ = exhaust velocity

Tsiolkovsky rocket equation:

$$\Delta v = v_e \ln\left(\frac{m_{initial}}{m_{final}}\right)$$

## Important Formulas

| Formula | Description |
|---------|-------------|
| $\vec{p} = m\vec{v}$ | Momentum definition |
| $\vec{J} = \vec{F}\Delta t = \Delta\vec{p}$ | Impulse-momentum |
| $m_1v_1 + m_2v_2 = (m_1+m_2)v_f$ | Perfectly inelastic |
| $e = \frac{v_{2f}-v_{1f}}{v_{1i}-v_{2i}}$ | Coefficient of restitution |
| $\Delta v = v_e \ln(m_i/m_f)$ | Rocket equation |

---

Back to: [[Dynamics MOC]] | [[Physics MOC]]
