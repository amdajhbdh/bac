---
tags: [physics, circuits, kirchhoff, electricity]
difficulty: medium
---
> [!summary] 📊 Note Summary
> 
> | Property | Value |
> |----------|-------|
> | **Difficulty** | `hard` #difficulty/hard |
> | **Formulas** | 35 |
> | **Concepts** | 0 |
> | **Related Notes** | 10 |
> | **Word Count** | 908 |
> | **Last Enhanced** | 2026-03-10 |



## 📊 Note Summary

| Property | Value |
|----------|-------|
| Difficulty | Hard |
| Formulas | 35 |
| Concepts | 0 |
| Related Notes | 10 |
| Word Count | 833 |
| Last Enhanced | 2026-03-10 |



# Kirchhoff's Laws

## Overview
Two fundamental laws for analyzing complex circuits.

## Kirchhoff's Current Law (KCL)

### Statement
The sum of currents entering a node equals the sum of currents leaving.

$$\sum I_{in} = \sum I_{out}$$

Or: $\sum I = 0$ (with sign convention)

### Physical Basis
Conservation of charge - charge cannot accumulate at a node.

### Example 1: Simple Node
$I_1 = I_2 + I_3$

If $I_1 = 5A$, $I_2 = 2A$, then $I_3 = 3A$

### Example 2: Multiple Currents
$I_1 + I_4 = I_2 + I_3$

## Kirchhoff's Voltage Law (KVL)

### Statement
The sum of voltage drops around any closed loop equals zero.

$$\sum V = 0$$

### Physical Basis
Conservation of energy - potential energy change around closed path is zero.

### Sign Convention
- Voltage rise (- to +): positive
- Voltage drop (+ to -): negative
- Current through resistor: $-IR$ (voltage drop)

### Example 3: Simple Loop
$$V - I \cdot R_1 - I \cdot R_2 = 0$$
$$V = I(R_1 + R_2)$$

## Solving Circuit Problems

### Method 1: Loop Current Method
1. Assign loop currents (clockwise)
2. Apply KVL to each loop
3. Solve system of equations

### Method 2: Node Voltage Method
1. Choose reference node (ground)
2. Assign voltages to other nodes
3. Apply KCL at each node
4. Solve system of equations

## Complex Circuit Examples

### Example 4: Two Loops
Circuit with 12V source, $R_1 = 4\Omega$, $R_2 = 6\Omega$, $R_3 = 3\Omega$

**Loop 1** (left): 
$$12 - 4I_1 - 3(I_1-I_2) = 0$$
$$12 = 7I_1 - 3I_2 \quad ...(1)$$

**Loop 2** (right): 
$$-3(I_2-I_1) - 6I_2 = 0$$
$$0 = -3I_1 + 9I_2 \quad ...(2)$$

From (2): $I_1 = 3I_2$

Substitute in (1): 
$$12 = 7(3I_2) - 3I_2 = 18I_2$$
$$I_2 = \frac{2}{3} A$$
$$I_1 = 2 A$$

Current through $R_3 = I_1 - I_2 = \frac{4}{3} A$

### Example 5: Node Voltage
Two sources (10V and 6V) with $R_1 = 2\Omega$, $R_2 = 4\Omega$

Node A voltage = $V_A$

**KCL at A**:
$$\frac{10 - V_A}{2} + \frac{6 - V_A}{4} = 0$$

Multiply by 4:
$$2(10 - V_A) + (6 - V_A) = 0$$
$$20 - 2V_A + 6 - V_A = 0$$
$$26 = 3V_A$$
$$V_A = \frac{26}{3} V$$

$$I_1 = \frac{10 - 26/3}{2} = \frac{4}{3} A$$
$$I_2 = \frac{26/3 - 6}{4} = \frac{2}{3} A$$

### Example 6: Three Resistors, Two Sources
6V and 12V sources with $R_1 = 3\Omega$, $R_2 = 2\Omega$, $R_3 = 4\Omega$

**Loop 1** (left):
$$6 - 3I_1 - 4(I_1-I_2) = 0$$
$$6 = 7I_1 - 4I_2 \quad ...(1)$$

**Loop 2** (right):
$$-4(I_2-I_1) - 2I_2 - 12 = 0$$
$$-12 = -4I_1 + 6I_2 \quad ...(2)$$

From (1): $I_1 = \frac{6 + 4I_2}{7}$

Substitute in (2):
$$-12 = -4\left(\frac{6 + 4I_2}{7}\right) + 6I_2$$
$$-84 = -24 - 16I_2 + 42I_2$$
$$-60 = 26I_2$$
$$I_2 = -\frac{30}{13} A$$ (negative = opposite direction)

$$I_1 = \frac{6 + 4(-30/13)}{7} = -\frac{6}{91} A$$

## Wheatstone Bridge

### Balanced Bridge
When balanced (no current through middle resistor):
$$\frac{R_1}{R_2} = \frac{R_3}{R_4}$$

### Example 7
$R_1 = 10\Omega$, $R_2 = 20\Omega$, $R_3 = 15\Omega$

Find $R_4$ for balance:
$$\frac{10}{20} = \frac{15}{R_4}$$
$$R_4 = 30\Omega$$

## Superposition Theorem

For circuits with multiple sources:
1. Consider each source separately (replace others with shorts/opens)
2. Calculate currents/voltages for each
3. Sum the results

### Example 8
Two voltage sources in circuit:
- With $V_1$ only ($V_2 = 0$): $I_1' = 2A$
- With $V_2$ only ($V_1 = 0$): $I_1'' = -1A$
- Total: $I_1 = 2 + (-1) = 1A$

## Thevenin's Theorem

Any linear circuit can be replaced by:
- $V_{th}$ (Thevenin voltage): open-circuit voltage
- $R_{th}$ (Thevenin resistance): resistance with sources off

### Finding Thevenin Equivalent
1. Remove load resistor
2. Find $V_{oc}$ (open circuit voltage) = $V_{th}$
3. Turn off sources (V -> short, I -> open)
4. Find resistance looking into terminals = $R_{th}$

## Norton's Theorem

Equivalent to Thevenin:
- $I_N$ (Norton current): short-circuit current
- $R_N$ (Norton resistance) = $R_{th}$

Conversion:
$$I_N = \frac{V_{th}}{R_{th}}$$
$$V_{th} = I_N \cdot R_N$$

## Practice Problems

### Problem 1
10V source with $2\Omega$, $3\Omega$, and $4\Omega$ resistors.
Find current through $4\Omega$ resistor.

<details>
<summary>Solution</summary>

Loop 1: $10 - 2I_1 - 4(I_1-I_2) = 0 \rightarrow 10 = 6I_1 - 4I_2$

Loop 2: $-4(I_2-I_1) - 3I_2 = 0 \rightarrow 0 = -4I_1 + 7I_2$

Solving: $I_1 = 2A$, $I_2 = \frac{8}{7} A$

$I_{4\Omega} = I_1 - I_2 = \frac{6}{7} A$
</details>

### Problem 2
Node with three branches: $I_1 = 5A$ (in), $I_2 = 2A$ (out), $I_3 = ?$

<details>
<summary>Solution</summary>

KCL: $I_1 = I_2 + I_3$

$5 = 2 + I_3$

$I_3 = 3A$ (out)
</details>

### Problem 3
Wheatstone bridge: $R_1 = 100\Omega$, $R_2 = 200\Omega$, $R_4 = 150\Omega$

Find $R_3$ for balance.

<details>
<summary>Solution</summary>

$$\frac{R_1}{R_2} = \frac{R_3}{R_4}$$
$$\frac{100}{200} = \frac{R_3}{150}$$
$$R_3 = 75\Omega$$
</details>

## Related
- [[Circuits - Ohm Law]]
- [[Circuits - RC Circuits]]
- [[Circuits - Power]]


## 🔗 Related Notes

- [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
- [[00-Meta/DEEP-CONTENT-STATUS.md|DEEP-CONTENT-STATUS]]
- [[00-Meta/MOCs/Physics MOC.md|Physics MOC]]
- [[Resource Links.md|Resource Links]]
- [[VAULT-COMPLETION-REPORT.md|VAULT-COMPLETION-REPORT]]
- [[Animations/ANIMATION_SPEC.md|ANIMATION_SPEC]]
- [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
- [[Animations/ALL-EXERCISES-COVERED.md|ALL-EXERCISES-COVERED]]
- [[00-Meta/MOCs/Natural Sciences MOC.md|Natural Sciences MOC]]
- [[Animations/ANIMATION_SPEC.md|ANIMATION_SPEC]]


> [!related] 🔗 Related Notes
> 
> - [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
> - [[Resource Links.md|Resource Links]]
> - [[ANIMATION-SYSTEM-COMPLETE.md|ANIMATION-SYSTEM-COMPLETE]]
> - [[QUICK-REFERENCE.md|QUICK-REFERENCE]]
> - [[ANIMATION-SYSTEM-COMPLETE.md|ANIMATION-SYSTEM-COMPLETE]]
> - [[Animations/ALL-EXERCISES-COVERED.md|ALL-EXERCISES-COVERED]]
> - [[00-Meta/DEEP-CONTENT-STATUS.md|DEEP-CONTENT-STATUS]]
> - [[00-Meta/MOCs/Chemistry MOC.md|Chemistry MOC]]
> - [[01-Concepts/Math/Complex-Numbers/Complex Numbers - Operations.md|Complex Numbers - Operations]]
> - [[Animations/ANIMATION_SPEC.md|ANIMATION_SPEC]]
