---
tags: [concept, physics, electromagnetism, induction, faraday]
subject: Physics
topic: Magnetic Induction
difficulty: Medium
---

# Induction - Faraday's Law

## Introduction

Faraday's law of electromagnetic induction is one of the four Maxwell's equations and describes how a changing magnetic field creates an electric field. This principle is fundamental to the operation of generators, transformers, and many other electrical devices.

## Faraday's Law Statement

The induced electromotive force (emf) in a closed loop equals the negative rate of change of the magnetic flux through the loop:

$$\mathcal{E} = -\frac{d\Phi_B}{dt}$$

Where:
- $\mathcal{E}$ is the induced electromotive force (in volts, V)
- $\Phi_B$ is the magnetic flux (in webers, Wb)
- The negative sign indicates the direction of the induced emf (Lenz's law)

## Magnetic Flux

Magnetic flux is a measure of the amount of magnetic field passing through a surface:

$$\Phi_B = \int_S \vec{B} \cdot d\vec{A} = BA\cos\theta$$

For a uniform magnetic field perpendicular to a flat surface:
- $\Phi_B = BA\cos\theta$ where:
  - $B$ is the magnetic field strength (Tesla, T)
  - $A$ is the area of the surface (m²)
  - $\theta$ is the angle between $\vec{B}$ and the normal to the surface

## Special Cases

### 1. Moving Conductor in Magnetic Field

When a conductor of length $l$ moves with velocity $v$ in a uniform magnetic field $B$:

$$\mathcal{E} = Blv$$

This applies when the motion is perpendicular to both the magnetic field and the conductor.

### 2. Changing Magnetic Field Strength

If the area and orientation remain constant but the magnetic field changes:

$$\mathcal{E} = -A\cos\theta \frac{dB}{dt}$$

### 3. Changing Loop Area

If the magnetic field and orientation remain constant but the area changes:

$$\mathcal{E} = -B\cos\theta \frac{dA}{dt}$$

### 4. Rotating Loop

If a loop rotates in a magnetic field with angular frequency $\omega$:

$$\mathcal{E} = -\frac{d}{dt}(BA\cos(\omega t)) = BA\omega\sin(\omega t)$$

## Examples

### Example 1: Moving Conductor

A conducting rod of length 0.5 m moves with a speed of 4 m/s perpendicular to a uniform magnetic field of 0.8 T. Find the induced emf.

**Solution:**
$$\mathcal{E} = Blv = 0.8 \times 0.5 \times 4 = 1.6 \text{ V}$$

### Example 2: Changing Current in Solenoid

A solenoid with 500 turns per meter and cross-sectional area 2 cm² has current increasing at a rate of 3 A/s. Find the induced emf in a coil of 20 turns looped around the solenoid.

Within the solenoid, the magnetic field is:
$$B = \mu_0 nI$$

The flux through one turn of the external coil is:
$$\Phi_{B,1} = BA = \mu_0 nIA$$

The total flux through the external coil (20 turns) is:
$$\Phi_B = N\Phi_{B,1} = N\mu_0 nIA$$

The induced emf is:
$$\mathcal{E} = -\frac{d\Phi_B}{dt} = -N\mu_0 nA \frac{dI}{dt}$$
$$\mathcal{E} = -20 \times 4\pi \times 10^{-7} \times 500 \times 2 \times 10^{-4} \times 3$$
$$\mathcal{E} = -20 \times 4\pi \times 10^{-7} \times 500 \times 2 \times 10^{-4} \times 3$$
$$\mathcal{E} = -7.54 \times 10^{-5} \text{ V}$$

The magnitude of the induced emf is $7.54 \times 10^{-5}$ V.

### Example 3: Rotating Loop

A rectangular loop of dimensions 10 cm × 20 cm rotates at 60 revolutions per second in a uniform magnetic field of 0.5 T. Find the maximum induced emf.

**Solution:**
Area: $A = 0.1 \times 0.2 = 0.02 \text{ m}^2$
Angular frequency: $\omega = 60 \times 2\pi = 120\pi \text{ rad/s}$

Maximum emf occurs when $\sin(\omega t) = 1$:
$$\mathcal{E}_{max} = BA\omega = 0.5 \times 0.02 \times 120\pi = 3.77 \text{ V}$$

## Applications

### Generators

Electric generators operate on Faraday's law by rotating coils in magnetic fields to produce alternating current.

The emf generated is:
$$\mathcal{E} = NAB\omega\sin(\omega t)$$

Where:
- $N$ = number of turns
- $A$ = area of the coil
- $B$ = magnetic field strength
- $\omega$ = angular frequency

### Induction Cooktops

These use rapidly changing magnetic fields to induce currents in metal cookware, heating it through resistive losses.

## Units and Dimensional Analysis

- Magnetic flux: Weber (Wb) = T·m² = V·s
- Magnetic field: Tesla (T) = Wb/m² = N/(A·m)
- Electromotive force: Volt (V) = W/A = J/C

## Important Considerations

1. **Direction of Induced EMF**: The negative sign in Faraday's law leads to Lenz's law, which states that the induced current opposes the change causing it.

2. **Rate of Change**: Only a *changing* magnetic flux induces an emf. A constant magnetic field produces no induced emf regardless of its strength.

3. **Number of Turns**: For coils with multiple turns, the total emf is proportional to the number of turns.

4. **Path Independence**: The induced emf depends only on the rate of change of flux, not on the shape of the loop or path of integration.

## Relationship to Maxwell's Equations

Faraday's law in differential form is one of Maxwell's equations:

$$\nabla \times \vec{E} = -\frac{\partial \vec{B}}{\partial t}$$

This shows that a time-varying magnetic field creates a circulating electric field.

## Experimental Demonstration

Michael Faraday demonstrated electromagnetic induction with his famous experiment:
1. A primary coil connected to a battery
2. A secondary coil connected to a galvanometer
3. Observing current flow when the primary circuit was switched on/off

This showed that relative motion or changing currents could induce emfs.

## Visualizations

![[Induction - Faraday Law.excalidraw]]

---
Back to: [[Magnetic Induction MOC]] | [[Physics MOC]]