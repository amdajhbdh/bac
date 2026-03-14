---
tags: [practice, physics, electromagnetism, induction, medium]
subject: Physics
topic: Magnetic Induction
difficulty: Medium
---

# Induction - Practice Medium

## Problems

### Problem 1: Solenoid with Changing Current

A solenoid has 500 turns per meter and a cross-sectional area of 4 cm². A coil of 100 turns is wrapped around it. If the current in the solenoid increases uniformly from 0 to 2 A in 0.5 s, find:

a) The induced emf in the coil
b) The direction of the induced current

**Solution:**

a) Magnetic field in solenoid: $B = \mu_0 nI$
$$\frac{dB}{dt} = \mu_0 n \frac{dI}{dt} = 4\pi \times 10^{-7} \times 500 \times \frac{2}{0.5} = 8\pi \times 10^{-4} \text{ T/s}$$

Flux through one turn of coil:
$$\Phi = BA = \mu_0 n I A$$

Rate of change:
$$\frac{d\Phi}{dt} = \mu_0 n A \frac{dI}{dt}$$

Induced emf:
$$\mathcal{E} = -N \frac{d\Phi}{dt} = -100 \times 4\pi \times 10^{-7} \times 500 \times 4 \times 10^{-4} \times \frac{2}{0.5}$$
$$\mathcal{E} = -100 \times 4\pi \times 10^{-7} \times 500 \times 4 \times 10^{-4} \times 4$$
$$\mathcal{E} = -0.503 \text{ V}$$

Magnitude: 0.503 V

b) The current in the solenoid is increasing, creating an increasing magnetic field pointing to the right (using right-hand rule). According to Lenz's law, the induced current in the outer coil will create a magnetic field pointing LEFT to oppose this increase. Using the right-hand rule, this corresponds to a counterclockwise current as viewed from the left.

---

### Problem 2: Moving Rod

A metal rod of length 0.5 m slides on parallel conducting rails in a uniform magnetic field of 0.4 T. The rails are connected by a resistor of 10 Ω. The rod moves with constant velocity 2 m/s.

a) Find the induced emf
b) Find the current in the circuit
c) What force is required to maintain constant velocity?

**Solution:**

a) $\mathcal{E} = Blv = 0.4 \times 0.5 \times 2 = 0.4 \text{ V}$

b) $I = \frac{\mathcal{E}}{R} = \frac{0.4}{10} = 0.04 \text{ A}$

c) The magnetic force on the rod opposes motion:
$$F_{mag} = BIl = 0.4 \times 0.04 \times 0.5 = 8 \times 10^{-3} \text{ N}$$

To maintain constant velocity, an external force of 8 mN must be applied in the direction of motion.

---

### Problem 3: Rectangular Loop Pulled from Field

A rectangular loop with dimensions 0.2 m × 0.3 m and resistance 2 Ω is initially fully inside a uniform magnetic field of 0.5 T. The loop is pulled out of the field in 0.2 s at constant velocity, maintaining its orientation.

a) Find the average induced emf
b) Find the induced current
c) Calculate the energy dissipated as heat

**Solution:**

a) Initial flux: $\Phi_i = BA = 0.5 \times 0.2 \times 0.3 = 0.03 \text{ Wb}$
Final flux: $\Phi_f = 0$ (outside field)

$$\mathcal{E}_{avg} = -\frac{\Delta\Phi}{\Delta t} = -\frac{0 - 0.03}{0.2} = 0.15 \text{ V}$$

b) $I = \frac{\mathcal{E}}{R} = \frac{0.15}{2} = 0.075 \text{ A}$

c) Energy dissipated: $E = I^2 R t = (0.075)^2 \times 2 \times 0.2 = 2.25 \times 10^{-3} \text{ J}$

---

### Problem 4: Rotating Circular Coil

A circular coil with 50 turns and radius 10 cm rotates at 1800 RPM in a uniform magnetic field of 0.3 T. The axis of rotation is perpendicular to the field.

Find:
a) The angular frequency
b) The maximum induced emf

**Solution:**

a) $\omega = 2\pi f = 2\pi \times \frac{1800}{60} = 60\pi \text{ rad/s}$

b) $A = \pi r^2 = \pi \times (0.1)^2 = 0.0314 \text{ m}^2$
$$\mathcal{E}_{max} = NBA\omega = 50 \times 0.3 \times 0.0314 \times 60\pi = 8.88 \text{ V}$$

---

### Problem 5: Transformer Calculations

A step-up transformer has 200 turns in the primary and 2000 turns in the secondary. The primary is connected to 230 V AC and draws 2 A.

Find:
a) Secondary voltage
b) Secondary current (assuming ideal transformer)
c) Power output

**Solution:**

a) $V_s = V_p \times \frac{N_s}{N_p} = 230 \times \frac{2000}{200} = 2300 \text{ V}$

b) For ideal transformer: $P_p = P_s$
$$I_s = I_p \times \frac{V_p}{V_s} = 2 \times \frac{230}{2300} = 0.2 \text{ A}$$

c) $P = V_s I_s = 2300 \times 0.2 = 460 \text{ W}$

---

### Problem 6: Self-Induction

A solenoid 0.5 m long has 500 turns and cross-sectional area 3 cm². Calculate:

a) Its inductance
b) The emf induced when current changes at 10 A/s

**Solution:**

a) $L = \frac{\mu_0 N^2 A}{l} = \frac{4\pi \times 10^{-7} \times 500^2 \times 3 \times 10^{-4}}{0.5}$
$$L = 1.88 \times 10^{-4} \text{ H}$$

b) $\mathcal{E} = -L \frac{dI}{dt} = -1.88 \times 10^{-4} \times 10 = -1.88 \times 10^{-3} \text{ V}$

Magnitude: 1.88 mV

---

### Problem 7: RL Circuit Time Constant

An inductor of 10 mH is connected in series with a 100 Ω resistor. Find:

a) The time constant
b) The time for current to reach 63% of maximum
c) The current after one time constant if V = 10 V

**Solution:**

a) $\tau = \frac{L}{R} = \frac{10 \times 10^{-3}}{100} = 10^{-4} \text{ s} = 0.1 \text{ ms}$

b) The time constant is 0.1 ms (by definition)

c) $I_{max} = \frac{V}{R} = \frac{10}{100} = 0.1 \text{ A}$
After one time constant: $I = 0.632 \times 0.1 = 0.0632 \text{ A}$

---

Back to: [[Magnetic Induction MOC]] | [[Physics MOC]]
