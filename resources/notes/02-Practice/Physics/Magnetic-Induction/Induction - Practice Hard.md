---
tags: [practice, physics, electromagnetism, induction, hard]
subject: Physics
topic: Magnetic Induction
difficulty: Hard
---

# Induction - Practice Hard

## Advanced Problems

### Problem 1: Complex Moving Loop System

A conducting rod of mass 0.1 kg and length 0.5 m slides without friction on two parallel rails separated by 0.3 m. The system is in a uniform magnetic field of 0.8 T pointing upward. A constant horizontal force of 0.5 N pulls the rod to the right. The rails are connected by a resistor of 5 Ω.

Initially, the rod is at rest at x = 0.

a) Find the terminal velocity of the rod
b) Find the magnetic force at terminal velocity
c) Calculate the power dissipated in the resistor at terminal velocity

**Solution:**

At terminal velocity, net force = 0:
$$F_{ext} = F_{mag}$$
$$F_{ext} = BIl = B \frac{\mathcal{E}}{R} l = B \frac{Bvl}{R} l = \frac{B^2 l^2 v}{R}$$

Solving for v:
$$v = \frac{F_{ext} R}{B^2 l^2} = \frac{0.5 \times 5}{(0.8)^2 \times (0.5)^2} = \frac{2.5}{0.16} = 15.6 \text{ m/s}$$

b) At terminal velocity:
$$\mathcal{E} = Blv = 0.8 \times 0.5 \times 15.6 = 6.24 \text{ V}$$
$$I = \frac{\mathcal{E}}{R} = \frac{6.24}{5} = 1.25 \text{ A}$$
$$F_{mag} = BIl = 0.8 \times 1.25 \times 0.5 = 0.5 \text{ N}$$

c) $P = I^2 R = (1.25)^2 \times 5 = 7.81 \text{ W}$

---

### Problem 2: Two Coils System

Two coils are placed close to each other. Coil 1 has 100 turns, coil 2 has 200 turns. When current in coil 1 changes from 0 to 3 A in 0.2 s, the average charge that flows in coil 2 is 0.015 C. The resistance of coil 2 is 20 Ω.

Find:
a) The mutual inductance M
b) The induced emf in coil 2
c) The magnetic flux through coil 2 per turn from coil 1

**Solution:**

a) $Q = C \times I_{avg}$ is not applicable here. Instead:
$$\mathcal{E}_2 = IR = 0.015 \times 20 = 0.3 \text{ V}$$

From Faraday's law:
$$\mathcal{E}_2 = -M \frac{dI_1}{dt}$$

$$\frac{dI_1}{dt} = \frac{3}{0.2} = 15 \text{ A/s}$$

$$M = \frac{\mathcal{E}_2}{dI_1/dt} = \frac{0.3}{15} = 0.02 \text{ H}$$

b) $\mathcal{E}_2 = 0.3 \text{ V}$ (as calculated above)

c) $\mathcal{E}_2 = -N_2 \frac{d\Phi_{21}}{dt}$

For average values:
$$\Phi_{21} = \frac{\mathcal{E}_2 \times \Delta t}{N_2} = \frac{0.3 \times 0.2}{200} = 3 \times 10^{-4} \text{ Wb}$$

---

### Problem 3: Alternating Generator Analysis

A generator consists of a rectangular coil with 500 turns, dimensions 0.4 m × 0.2 m, rotating at 1200 RPM in a 0.25 T magnetic field. The coil has a resistance of 10 Ω and is connected to a load of 90 Ω.

Find:
a) Maximum emf generated
b) RMS current in the circuit
c) Power delivered to the load
d) The instantaneous emf at t = 0.05 s

**Solution:**

a) $\omega = 2\pi \times \frac{1200}{60} = 40\pi \text{ rad/s}$
$A = 0.4 \times 0.2 = 0.08 \text{ m}^2$
$\mathcal{E}_{max} = NBA\omega = 500 \times 0.25 \times 0.08 \times 40\pi = 628 \text{ V}$

b) Total resistance: $R_{total} = 10 + 90 = 100 \Omega$
$I_{rms} = \frac{\mathcal{E}_{rms}}{R} = \frac{628/\sqrt{2}}{100} = \frac{444}{100} = 4.44 \text{ A}$

c) $P_{load} = I_{rms}^2 R_{load} = (4.44)^2 \times 90 = 1775 \text{ W}$

d) $\mathcal{E}(t) = \mathcal{E}_{max} \sin(\omega t) = 628 \sin(40\pi \times 0.05) = 628 \sin(2\pi) = 0$

---

### Problem 4: RL Circuit Transient Analysis

An inductor of 50 mH, resistor of 25 Ω, and switch are connected in series to a 100 V DC source. At t = 0, the switch is closed.

Find:
a) The initial current
b) The time constant
c) The current at t = 5 ms
d) The energy stored in the inductor at t = 2 ms

**Solution:**

a) At t = 0, inductor acts as open circuit: $I(0) = 0$

b) $\tau = \frac{L}{R} = \frac{50 \times 10^{-3}}{25} = 2 \text{ ms}$

c) $I(t) = I_{max} (1 - e^{-t/\tau}) = \frac{100}{25} (1 - e^{-5 \times 10^{-3}/2 \times 10^{-3}})$
$= 4 (1 - e^{-2.5}) = 4 (1 - 0.0821) = 3.67 \text{ A}$

d) At t = 2 ms: $I = 4 (1 - e^{-1}) = 4 (1 - 0.368) = 2.53 \text{ A}$
$$E = \frac{1}{2} LI^2 = \frac{1}{2} \times 50 \times 10^{-3} \times (2.53)^2 = 0.16 \text{ J}$$

---

### Problem 5: Transformer with Load

A transformer has 400 primary turns and 40 secondary turns. The primary is connected to 220 V AC and draws 1 A at no-load. When a load is connected, the primary current increases to 5 A.

Find:
a) The secondary voltage under load
b) The power transferred to the load
c) The efficiency of the transformer
d) If the secondary resistance is 2 Ω, find voltage drop

**Solution:**

a) $V_s = V_p \times \frac{N_s}{N_p} = 220 \times \frac{40}{400} = 22 \text{ V}$

b) Assuming ideal transformer at load: $P_{out} = V_p \times I_p \times \eta$
With no-load losses: $P_{no-load} = 220 \times 1 = 220 \text{ W}$

At 5 A primary: $P_{in} = 220 \times 5 = 1100 \text{ W}$
Output = Input - losses = 1100 - 220 = 880 W

c) Efficiency = $\frac{880}{1100} \times 100\% = 80\%$

d) Secondary current: $I_s = I_p \times \frac{N_p}{N_s} = 5 \times \frac{400}{40} = 50 \text{ A}$
Voltage drop: $V_{drop} = I_s \times R_s = 50 \times 2 = 100 \text{ V}$

This is a significant drop! The actual secondary voltage under load would be much lower: 22 - 100 = -78 V (impossible, indicating transformer cannot deliver this current - need thicker wire or fewer turns ratio).

---

### Problem 6: Energy in Changing Magnetic Field

A solenoid with 1000 turns, length 0.2 m, and radius 0.01 m has current increasing at 5 A/s.

Find:
a) The self-inductance
b) The induced emf
c) The rate of energy storage in the magnetic field
d) The energy density inside the solenoid when current is 2 A

**Solution:**

a) $L = \frac{\mu_0 N^2 A}{l} = \frac{4\pi \times 10^{-7} \times 1000^2 \times \pi \times (0.01)^2}{0.2}$
$L = \frac{4\pi \times 10^{-7} \times 10^6 \times \pi \times 10^{-4}}{0.2}$
$L = \frac{4\pi^2 \times 10^{-5}}{0.2} = 6.32 \times 10^{-4} \text{ H}$

b) $\mathcal{E} = -L \frac{dI}{dt} = -6.32 \times 10^{-4} \times 5 = -3.16 \times 10^{-3} \text{ V}$

c) Power (rate of energy storage):
$$P = \frac{dE}{dt} = \frac{1}{2} L \frac{dI^2}{dt} = LI \frac{dI}{dt}$$
At I = 2 A: $P = 6.32 \times 10^{-4} \times 2 \times 5 = 6.32 \times 10^{-3} \text{ W}$

d) At I = 2 A:
$$B = \mu_0 nI = \frac{4\pi \times 10^{-7} \times 1000 \times 2}{0.2} = 4\pi \times 10^{-3} \text{ T}$$
Energy density: $u = \frac{B^2}{2\mu_0} = \frac{(4\pi \times 10^{-3})^2}{2 \times 4\pi \times 10^{-7}} = \frac{16\pi^2 \times 10^{-6}}{8\pi \times 10^{-7}} = 2\pi = 6.28 \text{ J/m}^3$

---

### Problem 7: LC Circuit Oscillation

An inductor of 10 mH is connected in series with a capacitor of 100 μF. Initially, the capacitor is charged to 10 V.

Find:
a) The oscillation frequency
b) The maximum current
c) The current at t = 10 ms
d) The energy oscillating between inductor and capacitor

**Solution:**

a) $f = \frac{1}{2\pi\sqrt{LC}} = \frac{1}{2\pi\sqrt{10 \times 10^{-3} \times 100 \times 10^{-6}}} = \frac{1}{2\pi\sqrt{10^{-6}}} = \frac{1}{2\pi \times 10^{-3}} = 159 \text{ Hz}$

b) Energy in capacitor: $E = \frac{1}{2}CV^2 = \frac{1}{2} \times 100 \times 10^{-6} \times 100 = 5 \times 10^{-3} \text{ J}$

At max current, all energy in inductor:
$$\frac{1}{2}LI_{max}^2 = 5 \times 10^{-3}$$
$$I_{max} = \sqrt{\frac{2E}{L}} = \sqrt{\frac{2 \times 5 \times 10^{-3}}{10 \times 10^{-3}}} = \sqrt{1} = 1 \text{ A}$$

c) $\omega = 2\pi f = 1000 \text{ rad/s}$
$$I(t) = I_{max} \sin(\omega t) = 1 \sin(1000 \times 0.01) = 1 \sin(10) = 1 \sin(10 \text{ rad}) \approx 0.544 \text{ A}$$

d) Total energy: $5 \times 10^{-3}$ J (constant, oscillating between L and C)

---

Back to: [[Magnetic Induction MOC]] | [[Physics MOC]]
