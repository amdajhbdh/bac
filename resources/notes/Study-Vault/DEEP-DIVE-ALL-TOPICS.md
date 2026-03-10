# Complex Numbers - Deep Dive

## Nombres Complexes


### Definition & Representation
#### Définition et Représentation


## Definition
A complex number z = a + bi where a, b ∈ ℝ and i² = -1

## Représentation
- Rectangular form: z = a + bi
- Polar form: z = r(cos θ + i sin θ) = r·e^(iθ)
- Exponential form: z = r·e^(iθ)

## Key Properties
- Modulus: |z| = √(a² + b²)
- Argument: arg(z) = θ
- Conjugate: z̄ = a - bi
        

### Operations
#### Opérations


## Addition & Subtraction
(a + bi) + (c + di) = (a+c) + (b+d)i

## Multiplication
(a + bi)(c + di) = (ac - bd) + (ad + bc)i

## Division
(a + bi)/(c + di) = [(a + bi)(c - di)]/[(c + di)(c - di)]

## De Moivre's Theorem
[r(cos θ + i sin θ)]^n = r^n(cos nθ + i sin nθ)
        

## Questions d'Examen
### Questions BAC Possibles

1. Solve z² + 2z + 5 = 0
2. Find |z| and arg(z) for z = 1 + i√3
3. Express (1 + i)^8 in rectangular form
4. Find all cube roots of 8
5. Prove that |z₁z₂| = |z₁||z₂|
6. Solve z⁴ = -16
7. Find the locus of points where |z - 1| = |z + 1|
8. Express e^(iπ/4) in rectangular form

---

# Integrals - Deep Dive

## Intégrales


### Fundamental Theorem
#### Théorème Fondamental


## Definition
∫f(x)dx = F(x) + C where F'(x) = f(x)

## Fundamental Theorem of Calculus
∫ₐᵇ f(x)dx = F(b) - F(a)

## Basic Integrals
- ∫xⁿ dx = x^(n+1)/(n+1) + C (n ≠ -1)
- ∫1/x dx = ln|x| + C
- ∫eˣ dx = eˣ + C
- ∫sin(x) dx = -cos(x) + C
- ∫cos(x) dx = sin(x) + C
        

### Integration Techniques
#### Techniques d'Intégration


## Substitution Method
∫f(g(x))g'(x)dx = ∫f(u)du where u = g(x)

## Integration by Parts
∫u dv = uv - ∫v du
## Partial Fractions
For rational functions: decompose into simpler fractions

## Trigonometric Substitution
- For √(a² - x²): use x = a sin θ
- For √(a² + x²): use x = a tan θ
- For √(x² - a²): use x = a sec θ
        

## Questions d'Examen
### Questions BAC Possibles

1. Calculate ∫₀¹ (3x² + 2x) dx
2. Find ∫ x·eˣ dx using integration by parts
3. Evaluate ∫ 1/(x² + 1) dx
4. Calculate ∫₀^(π/2) sin²(x) dx
5. Find ∫ (2x + 1)/(x² + x + 1) dx
6. Evaluate ∫ √(1 - x²) dx from -1 to 1
7. Calculate ∫ x/(x² + 4) dx
8. Find the area under y = e^(-x) from 0 to ∞

---

# Matrices - Deep Dive

## Matrices


### Matrix Operations
#### Opérations Matricielles


## Basic Operations
- Addition: (A + B)ᵢⱼ = Aᵢⱼ + Bᵢⱼ
- Scalar multiplication: (kA)ᵢⱼ = k·Aᵢⱼ
- Matrix multiplication: (AB)ᵢⱼ = Σₖ Aᵢₖ·Bₖⱼ

## Properties
- (AB)C = A(BC) (associative)
- A(B + C) = AB + AC (distributive)
- AB ≠ BA (not commutative)
- (AB)ᵀ = BᵀAᵀ
        

### Determinants & Inverses
#### Déterminants et Inverses


## 2×2 Matrix
det(A) = ad - bc for A = [[a,b],[c,d]]
A⁻¹ = (1/det(A)) × [[d,-b],[-c,a]]

## 3×3 Matrix (Cofactor Expansion)
det(A) = a₁₁(a₂₂a₃₃ - a₂₃a₃₂) - a₁₂(...) + a₁₃(...)

## Properties
- det(AB) = det(A)·det(B)
- det(Aᵀ) = det(A)
- A is invertible iff det(A) ≠ 0
- AA⁻¹ = I
        

## Questions d'Examen
### Questions BAC Possibles

1. Find the determinant of [[1,2,3],[4,5,6],[7,8,9]]
2. Calculate the inverse of [[2,1],[5,3]]
3. Solve the system using matrices: 2x + y = 5, x - y = 1
4. Find eigenvalues of [[3,1],[1,3]]
5. Prove that det(AB) = det(A)·det(B)
6. Find the rank of [[1,2,3],[2,4,6]]
7. Calculate [[1,2],[3,4]]³
8. Determine if [[1,0],[0,1]] is orthogonal

---

# Number Theory - Deep Dive

## Théorie des Nombres


### Divisibility & Primes
#### Divisibilité et Nombres Premiers


## Divisibility
a|b means ∃k ∈ ℤ: b = ka

## Prime Numbers
A prime p > 1 has only divisors 1 and p

## Fundamental Theorem of Arithmetic
Every integer > 1 is uniquely a product of primes

## GCD & LCM
- gcd(a,b): greatest common divisor
- lcm(a,b): least common multiple
- gcd(a,b) × lcm(a,b) = |a × b|
        

### Modular Arithmetic
#### Arithmétique Modulaire


## Congruence
a ≡ b (mod n) means n|(a - b)

## Properties
- If a ≡ b (mod n) and c ≡ d (mod n), then:
  - a + c ≡ b + d (mod n)
  - ac ≡ bd (mod n)

## Fermat's Little Theorem
If p is prime and gcd(a,p) = 1: a^(p-1) ≡ 1 (mod p)

## Euler's Theorem
> [!note] If gcd(a,n) = 1: a^φ(n) ≡ 1 (mod n)
        

## Questions d'Examen
### Questions BAC Possibles

1. Find gcd(48, 18) using Euclidean algorithm
2. Prove that if p is prime and p|ab, then p|a or p|b
3. Solve 3x ≡ 7 (mod 11)
4. Find the last digit of 7^100
5. Determine if 17 is prime using trial division
6. Find lcm(12, 18)
7. Prove that √2 is irrational
8. Find all solutions to x² ≡ 1 (mod 8)

---

# Dynamics - Deep Dive

## Dynamique


### Newton's Laws
#### Lois de Newton


## First Law (Inertia)
An object at rest stays at rest unless acted upon by force

## Second Law
F = ma (Force = mass × acceleration)

## Third Law
For every action, there's an equal and opposite reaction

## Applications
- Weight: W = mg
- Friction: f = μN
- Tension: T in strings/cables
        

### Motion Analysis
#### Analyse du Mouvement


## Kinematics
- v = v₀ + at
- s = v₀t + ½at²
- v² = v₀² + 2as

## Energy
- Kinetic: KE = ½mv²
- Potential: PE = mgh
- Work: W = F·s·cos(θ)

## Momentum
- p = mv
- Impulse: J = FΔt = Δp
        

## Questions d'Examen
### Questions BAC Possibles

1. A 5 kg object accelerates at 2 m/s². Find the net force
2. Calculate the friction force on a 10 kg block on a 30° incline (μ = 0.2)
3. Find the velocity after 3 seconds if v₀ = 5 m/s and a = 2 m/s²
4. A ball is thrown upward at 20 m/s. Find max height
5. Two objects collide elastically. Find final velocities
6. Calculate work done by a 50 N force over 10 m
7. Find the tension in a rope supporting a 20 kg mass
8. Analyze motion on an inclined plane with friction

---

# Electric Fields - Deep Dive

## Champs Électriques


### Electric Field Basics
#### Bases des Champs Électriques


## Definition
E = F/q (Electric field = Force per unit charge)

## Coulomb's Law
F = k(q₁q₂)/r² where k = 8.99 × 10⁹ N·m²/C²

## Electric Field of Point Charge
E = kq/r²

## Superposition Principle
E_total = E₁ + E₂ + ... (vector sum)
        

### Potential & Energy
#### Potentiel et Énergie


## Electric Potential
V = W/q (Potential = Work per unit charge)

## Potential of Point Charge
V = kq/r

## Potential Difference
ΔV = V_B - V_A = -∫ₐᵇ E·dl

## Electric Potential Energy
U = qV = kq₁q₂/r
        

## Questions d'Examen
### Questions BAC Possibles

1. Calculate the electric field at distance r from a charge q
2. Find the force between two charges: q₁ = 2 μC, q₂ = 3 μC, r = 0.5 m
3. Determine the potential at a point due to multiple charges
4. Calculate the work done moving a charge in an electric field
5. Find the electric field inside and outside a uniformly charged sphere
6. Analyze the motion of a charged particle in a uniform field
7. Calculate the capacitance of a parallel plate capacitor
8. Find the energy stored in an electric field

---

# Circuits - Deep Dive

## Circuits


### Ohm's Law & Resistance
#### Loi d'Ohm et Résistance


## Ohm's Law
V = IR (Voltage = Current × Resistance)

## Resistance
R = ρL/A where ρ is resistivity, L is length, A is area

## Power
P = VI = I²R = V²/R

## Series & Parallel
- Series: R_total = R₁ + R₂ + ...
- Parallel: 1/R_total = 1/R₁ + 1/R₂ + ...
        

### Kirchhoff's Laws
#### Lois de Kirchhoff


## Current Law (KCL)
Sum of currents entering = Sum of currents leaving

## Voltage Law (KVL)
Sum of voltages around a closed loop = 0

## Applications
- Analyze complex circuits
- Find unknown currents and voltages
- Solve multi-loop circuits
        

## Questions d'Examen
### Questions BAC Possibles

1. Find the current through a 10 Ω resistor with 5 V applied
2. Calculate total resistance of 3 resistors in series: 2Ω, 3Ω, 5Ω
3. Find equivalent resistance of 2Ω and 3Ω in parallel
4. Analyze a circuit with multiple loops using Kirchhoff's laws
5. Calculate power dissipated in a 100 Ω resistor with 2 A current
6. Find the current distribution in a parallel circuit
7. Determine the voltage across each resistor in a series circuit
8. Analyze a circuit with both series and parallel components

---

# Solutions - Deep Dive

## Solutions


### Concentration & Molarity
#### Concentration et Molarité


## Molarity
M = n/V (moles per liter)

## Molality
m = n/kg_solvent

## Mass Percent
% = (mass_solute / mass_solution) × 100

## Dilution
M₁V₁ = M₂V₂
        

### Solubility & Saturation
#### Solubilité et Saturation


## Solubility
Maximum amount of solute that dissolves in solvent

## Saturation
Solution contains maximum dissolved solute

## Supersaturation
Solution contains more solute than normally possible

## Factors Affecting Solubility
- Temperature
- Pressure
- Nature of solute and solvent
        

## Questions d'Examen
### Questions BAC Possibles

1. Calculate molarity of 5 g NaCl in 250 mL solution
2. Dilute 100 mL of 2 M HCl to 500 mL. Find new concentration
3. Prepare 1 L of 0.5 M solution from solid solute
4. Calculate mass percent of solute in solution
5. Determine molality of a solution
6. Analyze solubility curves and predict crystallization
7. Calculate osmotic pressure of a solution
8. Determine colligative properties of solutions

---

# Organic Chemistry - Deep Dive

## Chimie Organique


### Hydrocarbons
#### Hydrocarbures


## Alkanes (CₙH₂ₙ₊₂)
- Saturated hydrocarbons
- Single C-C bonds
- Examples: CH₄, C₂H₆, C₃H₈

## Alkenes (CₙH₂ₙ)
- Unsaturated hydrocarbons
- C=C double bonds
- Examples: C₂H₄, C₃H₆

## Alkynes (CₙH₂ₙ₋₂)
- Triple C≡C bonds
- Examples: C₂H₂, C₃H₄

## Aromatic Compounds
- Benzene ring (C₆H₆)
- Resonance structures
        

### Functional Groups
#### Groupes Fonctionnels


## Alcohols (-OH)
- Primary, secondary, tertiary
- Oxidation to aldehydes/ketones

## Aldehydes & Ketones (C=O)
- Aldehyde: -CHO
- Ketone: -CO-

## Carboxylic Acids (-COOH)
- Acidic functional group
- Esterification reactions

## Esters (-COO-)
- From carboxylic acids + alcohols
- Hydrolysis reactions
        

## Questions d'Examen
### Questions BAC Possibles

1. Name and draw the structure of pentane isomers
2. Write the equation for combustion of octane
3. Predict products of addition reaction with alkene
4. Distinguish between primary, secondary, tertiary alcohols
5. Write oxidation reactions of alcohols
6. Predict esterification products
7. Analyze substitution vs elimination reactions
8. Determine molecular formula from combustion data

---

# Kinetics - Deep Dive

## Cinétique


### Reaction Rates
#### Vitesses de Réaction


## Rate Definition
rate = -Δ[A]/Δt = Δ[B]/Δt

## Rate Law
rate = k[A]ᵐ[B]ⁿ where m, n are orders

## Half-life
t₁/₂ = time for [A] to reduce to half

## First Order
ln[A] = ln[A]₀ - kt
t₁/₂ = 0.693/k

## Second Order
1/[A] = 1/[A]₀ + kt
t₁/₂ = 1/(k[A]₀)
        

### Activation Energy & Catalysts
#### Énergie d'Activation et Catalyseurs


## Arrhenius Equation
k = Ae^(-Eₐ/RT)

## Temperature Effect
ln(k₂/k₁) = (Eₐ/R)(1/T₁ - 1/T₂)

## Catalysts
- Lower activation energy
- Not consumed in reaction
- Increase reaction rate
        

## Questions d'Examen
### Questions BAC Possibles

1. Determine reaction order from experimental data
2. Calculate rate constant from half-life
3. Predict effect of temperature on reaction rate
4. Analyze catalyst mechanism
5. Calculate activation energy from Arrhenius equation
6. Determine rate law from initial rate data
7. Predict concentration after time t
8. Compare reaction rates at different conditions

---

# Genetics - Deep Dive

## Génétique


### Mendel's Laws
#### Lois de Mendel


## First Law (Segregation)
Alleles separate during gamete formation
Each gamete receives one allele

## Second Law (Independent Assortment)
Alleles of different genes assort independently
Applies to genes on different chromosomes

## Punnett Square
Tool for predicting offspring genotypes

## Phenotypic Ratios
- Monohybrid: 3:1 (dominant:recessive)
- Dihybrid: 9:3:3:1
        

### DNA & Replication
#### ADN et Réplication


## DNA Structure
- Double helix
- Base pairs: A-T, G-C
- Sugar-phosphate backbone

## Replication
- Semi-conservative
- DNA polymerase
- Leading and lagging strands

## Mutations
- Point mutations
- Insertions/deletions
- Chromosomal aberrations
        

## Questions d'Examen
### Questions BAC Possibles

1. Solve monohybrid cross problems
2. Solve dihybrid cross problems
3. Predict phenotypic ratios
4. Analyze pedigree charts
5. Determine genotypes from phenotypes
6. Explain DNA replication process
7. Identify types of mutations
8. Calculate genetic probabilities

---

# Reproduction - Deep Dive

## Reproduction


### Human Reproductive System
#### Système Reproducteur Humain


## Male System
- Testes: produce sperm
- Epididymis: stores sperm
- Prostate: produces seminal fluid

## Female System
- Ovaries: produce eggs
- Fallopian tubes: transport egg
- Uterus: site of implantation

## Menstrual Cycle
- Follicular phase
- Ovulation
- Luteal phase
        

### Fertilization & Development
#### Fécondation et Développement


## Fertilization
- Sperm meets egg
- Zygote formation
- Meiosis completion

## Early Development
- Cleavage
- Blastula formation
- Gastrulation

## Embryonic Development
- Organogenesis
- Fetal development
- Birth
        

## Questions d'Examen
### Questions BAC Possibles

1. Describe the menstrual cycle phases
2. Explain fertilization process
3. Describe early embryonic development
4. Identify reproductive structures and functions
5. Explain hormonal control of reproduction
6. Describe fetal development stages
7. Analyze contraceptive methods
8. Explain pregnancy and childbirth

---

# Nervous System - Deep Dive

## Système Nerveux


### Neurons & Synapses
#### Neurones et Synapses


## Neuron Structure
- Soma (cell body)
- Dendrites (receive signals)
- Axon (transmit signals)

## Action Potential
- Resting potential: -70 mV
- Depolarization
- Repolarization
- Hyperpolarization

## Synaptic Transmission
- Neurotransmitter release
- Receptor binding
- Signal propagation
        

### Brain & Spinal Cord
#### Cerveau et Moelle Épinière


## Brain Regions
- Cerebrum: consciousness, thought
- Cerebellum: coordination, balance
- Brainstem: vital functions

## Spinal Cord
- Transmits signals
- Reflex arcs
- Gray and white matter

## Nervous System Organization
- Central: brain and spinal cord
- Peripheral: somatic and autonomic
        

## Questions d'Examen
### Questions BAC Possibles

1. Describe neuron structure and function
2. Explain action potential generation
3. Describe synaptic transmission
4. Identify brain regions and functions
5. Explain reflex arc mechanism
6. Describe nervous system organization
7. Analyze neurotransmitter effects
8. Explain sensory and motor pathways

---

# Geology - Deep Dive

## Géologie


### Plate Tectonics
#### Tectonique des Plaques


## Earth Structure
- Crust: thin outer layer
- Mantle: hot rock
- Core: iron and nickel

## Plate Boundaries
- Divergent: plates separate
- Convergent: plates collide
- Transform: plates slide

## Plate Movement
- Driven by mantle convection
- Seafloor spreading
- Subduction zones
        

### Earthquakes & Volcanoes
#### Tremblements de Terre et Volcans


## Earthquakes
- Caused by plate movement
- Seismic waves
- Magnitude and intensity

## Volcanoes
- Magma eruption
- Types: shield, cinder, composite
- Hazards: lava, ash, gases

## Rock Cycle
- Igneous, sedimentary, metamorphic
- Transformation processes
        

## Questions d'Examen
### Questions BAC Possibles

1. Describe Earth's internal structure
2. Explain plate tectonic theory
3. Analyze plate boundary types
4. Describe earthquake mechanisms
5. Explain volcanic activity
6. Analyze rock cycle processes
7. Predict geological hazards
8. Interpret geological maps and data

---
