const fs = require('fs');
const path = require('path');

// French translations for common terms
const translations = {
  "Complex Numbers": "Nombres Complexes",
  "Integrals": "Intégrales",
  "Matrices": "Matrices",
  "Number Theory": "Théorie des Nombres",
  "Dynamics": "Dynamique",
  "Electric Fields": "Champs Électriques",
  "Circuits": "Circuits",
  "Solutions": "Solutions",
  "Organic Chemistry": "Chimie Organique",
  "Kinetics": "Cinétique",
  "Genetics": "Génétique",
  "Reproduction": "Reproduction",
  "Nervous System": "Système Nerveux",
  "Geology": "Géologie",
  "Definition": "Définition",
  "Formula": "Formule",
  "Example": "Exemple",
  "Properties": "Propriétés",
  "Applications": "Applications",
  "Exam Questions": "Questions d'Examen",
  "Possible BAC Questions": "Questions BAC Possibles"
};

function translate(text) {
  return translations[text] || text;
}

const deepDiveTopics = {
  "Complex Numbers": {
    sections: [
      {
        title: "Definition & Representation",
        fr: "Définition et Représentation",
        content: `
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
        `
      },
      {
        title: "Operations",
        fr: "Opérations",
        content: `
## Addition & Subtraction
(a + bi) + (c + di) = (a+c) + (b+d)i

## Multiplication
(a + bi)(c + di) = (ac - bd) + (ad + bc)i

## Division
(a + bi)/(c + di) = [(a + bi)(c - di)]/[(c + di)(c - di)]

## De Moivre's Theorem
[r(cos θ + i sin θ)]^n = r^n(cos nθ + i sin nθ)
        `
      }
    ],
    examQuestions: [
      "Solve z² + 2z + 5 = 0",
      "Find |z| and arg(z) for z = 1 + i√3",
      "Express (1 + i)^8 in rectangular form",
      "Find all cube roots of 8",
      "Prove that |z₁z₂| = |z₁||z₂|",
      "Solve z⁴ = -16",
      "Find the locus of points where |z - 1| = |z + 1|",
      "Express e^(iπ/4) in rectangular form"
    ]
  },
  "Integrals": {
    sections: [
      {
        title: "Fundamental Theorem",
        fr: "Théorème Fondamental",
        content: `
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
        `
      },
      {
        title: "Integration Techniques",
        fr: "Techniques d'Intégration",
        content: `
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
        `
      }
    ],
    examQuestions: [
      "Calculate ∫₀¹ (3x² + 2x) dx",
      "Find ∫ x·eˣ dx using integration by parts",
      "Evaluate ∫ 1/(x² + 1) dx",
      "Calculate ∫₀^(π/2) sin²(x) dx",
      "Find ∫ (2x + 1)/(x² + x + 1) dx",
      "Evaluate ∫ √(1 - x²) dx from -1 to 1",
      "Calculate ∫ x/(x² + 4) dx",
      "Find the area under y = e^(-x) from 0 to ∞"
    ]
  },
  "Matrices": {
    sections: [
      {
        title: "Matrix Operations",
        fr: "Opérations Matricielles",
        content: `
## Basic Operations
- Addition: (A + B)ᵢⱼ = Aᵢⱼ + Bᵢⱼ
- Scalar multiplication: (kA)ᵢⱼ = k·Aᵢⱼ
- Matrix multiplication: (AB)ᵢⱼ = Σₖ Aᵢₖ·Bₖⱼ

## Properties
- (AB)C = A(BC) (associative)
- A(B + C) = AB + AC (distributive)
- AB ≠ BA (not commutative)
- (AB)ᵀ = BᵀAᵀ
        `
      },
      {
        title: "Determinants & Inverses",
        fr: "Déterminants et Inverses",
        content: `
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
        `
      }
    ],
    examQuestions: [
      "Find the determinant of [[1,2,3],[4,5,6],[7,8,9]]",
      "Calculate the inverse of [[2,1],[5,3]]",
      "Solve the system using matrices: 2x + y = 5, x - y = 1",
      "Find eigenvalues of [[3,1],[1,3]]",
      "Prove that det(AB) = det(A)·det(B)",
      "Find the rank of [[1,2,3],[2,4,6]]",
      "Calculate [[1,2],[3,4]]³",
      "Determine if [[1,0],[0,1]] is orthogonal"
    ]
  },
  "Number Theory": {
    sections: [
      {
        title: "Divisibility & Primes",
        fr: "Divisibilité et Nombres Premiers",
        content: `
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
        `
      },
      {
        title: "Modular Arithmetic",
        fr: "Arithmétique Modulaire",
        content: `
## Congruence
a ≡ b (mod n) means n|(a - b)

## Properties
- If a ≡ b (mod n) and c ≡ d (mod n), then:
  - a + c ≡ b + d (mod n)
  - ac ≡ bd (mod n)

## Fermat's Little Theorem
If p is prime and gcd(a,p) = 1: a^(p-1) ≡ 1 (mod p)

## Euler's Theorem
If gcd(a,n) = 1: a^φ(n) ≡ 1 (mod n)
        `
      }
    ],
    examQuestions: [
      "Find gcd(48, 18) using Euclidean algorithm",
      "Prove that if p is prime and p|ab, then p|a or p|b",
      "Solve 3x ≡ 7 (mod 11)",
      "Find the last digit of 7^100",
      "Determine if 17 is prime using trial division",
      "Find lcm(12, 18)",
      "Prove that √2 is irrational",
      "Find all solutions to x² ≡ 1 (mod 8)"
    ]
  },
  "Dynamics": {
    sections: [
      {
        title: "Newton's Laws",
        fr: "Lois de Newton",
        content: `
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
        `
      },
      {
        title: "Motion Analysis",
        fr: "Analyse du Mouvement",
        content: `
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
        `
      }
    ],
    examQuestions: [
      "A 5 kg object accelerates at 2 m/s². Find the net force",
      "Calculate the friction force on a 10 kg block on a 30° incline (μ = 0.2)",
      "Find the velocity after 3 seconds if v₀ = 5 m/s and a = 2 m/s²",
      "A ball is thrown upward at 20 m/s. Find max height",
      "Two objects collide elastically. Find final velocities",
      "Calculate work done by a 50 N force over 10 m",
      "Find the tension in a rope supporting a 20 kg mass",
      "Analyze motion on an inclined plane with friction"
    ]
  },
  "Electric Fields": {
    sections: [
      {
        title: "Electric Field Basics",
        fr: "Bases des Champs Électriques",
        content: `
## Definition
E = F/q (Electric field = Force per unit charge)

## Coulomb's Law
F = k(q₁q₂)/r² where k = 8.99 × 10⁹ N·m²/C²

## Electric Field of Point Charge
E = kq/r²

## Superposition Principle
E_total = E₁ + E₂ + ... (vector sum)
        `
      },
      {
        title: "Potential & Energy",
        fr: "Potentiel et Énergie",
        content: `
## Electric Potential
V = W/q (Potential = Work per unit charge)

## Potential of Point Charge
V = kq/r

## Potential Difference
ΔV = V_B - V_A = -∫ₐᵇ E·dl

## Electric Potential Energy
U = qV = kq₁q₂/r
        `
      }
    ],
    examQuestions: [
      "Calculate the electric field at distance r from a charge q",
      "Find the force between two charges: q₁ = 2 μC, q₂ = 3 μC, r = 0.5 m",
      "Determine the potential at a point due to multiple charges",
      "Calculate the work done moving a charge in an electric field",
      "Find the electric field inside and outside a uniformly charged sphere",
      "Analyze the motion of a charged particle in a uniform field",
      "Calculate the capacitance of a parallel plate capacitor",
      "Find the energy stored in an electric field"
    ]
  },
  "Circuits": {
    sections: [
      {
        title: "Ohm's Law & Resistance",
        fr: "Loi d'Ohm et Résistance",
        content: `
## Ohm's Law
V = IR (Voltage = Current × Resistance)

## Resistance
R = ρL/A where ρ is resistivity, L is length, A is area

## Power
P = VI = I²R = V²/R

## Series & Parallel
- Series: R_total = R₁ + R₂ + ...
- Parallel: 1/R_total = 1/R₁ + 1/R₂ + ...
        `
      },
      {
        title: "Kirchhoff's Laws",
        fr: "Lois de Kirchhoff",
        content: `
## Current Law (KCL)
Sum of currents entering = Sum of currents leaving

## Voltage Law (KVL)
Sum of voltages around a closed loop = 0

## Applications
- Analyze complex circuits
- Find unknown currents and voltages
- Solve multi-loop circuits
        `
      }
    ],
    examQuestions: [
      "Find the current through a 10 Ω resistor with 5 V applied",
      "Calculate total resistance of 3 resistors in series: 2Ω, 3Ω, 5Ω",
      "Find equivalent resistance of 2Ω and 3Ω in parallel",
      "Analyze a circuit with multiple loops using Kirchhoff's laws",
      "Calculate power dissipated in a 100 Ω resistor with 2 A current",
      "Find the current distribution in a parallel circuit",
      "Determine the voltage across each resistor in a series circuit",
      "Analyze a circuit with both series and parallel components"
    ]
  },
  "Solutions": {
    sections: [
      {
        title: "Concentration & Molarity",
        fr: "Concentration et Molarité",
        content: `
## Molarity
M = n/V (moles per liter)

## Molality
m = n/kg_solvent

## Mass Percent
% = (mass_solute / mass_solution) × 100

## Dilution
M₁V₁ = M₂V₂
        `
      },
      {
        title: "Solubility & Saturation",
        fr: "Solubilité et Saturation",
        content: `
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
        `
      }
    ],
    examQuestions: [
      "Calculate molarity of 5 g NaCl in 250 mL solution",
      "Dilute 100 mL of 2 M HCl to 500 mL. Find new concentration",
      "Prepare 1 L of 0.5 M solution from solid solute",
      "Calculate mass percent of solute in solution",
      "Determine molality of a solution",
      "Analyze solubility curves and predict crystallization",
      "Calculate osmotic pressure of a solution",
      "Determine colligative properties of solutions"
    ]
  },
  "Organic Chemistry": {
    sections: [
      {
        title: "Hydrocarbons",
        fr: "Hydrocarbures",
        content: `
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
        `
      },
      {
        title: "Functional Groups",
        fr: "Groupes Fonctionnels",
        content: `
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
        `
      }
    ],
    examQuestions: [
      "Name and draw the structure of pentane isomers",
      "Write the equation for combustion of octane",
      "Predict products of addition reaction with alkene",
      "Distinguish between primary, secondary, tertiary alcohols",
      "Write oxidation reactions of alcohols",
      "Predict esterification products",
      "Analyze substitution vs elimination reactions",
      "Determine molecular formula from combustion data"
    ]
  },
  "Kinetics": {
    sections: [
      {
        title: "Reaction Rates",
        fr: "Vitesses de Réaction",
        content: `
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
        `
      },
      {
        title: "Activation Energy & Catalysts",
        fr: "Énergie d'Activation et Catalyseurs",
        content: `
## Arrhenius Equation
k = Ae^(-Eₐ/RT)

## Temperature Effect
ln(k₂/k₁) = (Eₐ/R)(1/T₁ - 1/T₂)

## Catalysts
- Lower activation energy
- Not consumed in reaction
- Increase reaction rate
        `
      }
    ],
    examQuestions: [
      "Determine reaction order from experimental data",
      "Calculate rate constant from half-life",
      "Predict effect of temperature on reaction rate",
      "Analyze catalyst mechanism",
      "Calculate activation energy from Arrhenius equation",
      "Determine rate law from initial rate data",
      "Predict concentration after time t",
      "Compare reaction rates at different conditions"
    ]
  },
  "Genetics": {
    sections: [
      {
        title: "Mendel's Laws",
        fr: "Lois de Mendel",
        content: `
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
        `
      },
      {
        title: "DNA & Replication",
        fr: "ADN et Réplication",
        content: `
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
        `
      }
    ],
    examQuestions: [
      "Solve monohybrid cross problems",
      "Solve dihybrid cross problems",
      "Predict phenotypic ratios",
      "Analyze pedigree charts",
      "Determine genotypes from phenotypes",
      "Explain DNA replication process",
      "Identify types of mutations",
      "Calculate genetic probabilities"
    ]
  },
  "Reproduction": {
    sections: [
      {
        title: "Human Reproductive System",
        fr: "Système Reproducteur Humain",
        content: `
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
        `
      },
      {
        title: "Fertilization & Development",
        fr: "Fécondation et Développement",
        content: `
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
        `
      }
    ],
    examQuestions: [
      "Describe the menstrual cycle phases",
      "Explain fertilization process",
      "Describe early embryonic development",
      "Identify reproductive structures and functions",
      "Explain hormonal control of reproduction",
      "Describe fetal development stages",
      "Analyze contraceptive methods",
      "Explain pregnancy and childbirth"
    ]
  },
  "Nervous System": {
    sections: [
      {
        title: "Neurons & Synapses",
        fr: "Neurones et Synapses",
        content: `
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
        `
      },
      {
        title: "Brain & Spinal Cord",
        fr: "Cerveau et Moelle Épinière",
        content: `
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
        `
      }
    ],
    examQuestions: [
      "Describe neuron structure and function",
      "Explain action potential generation",
      "Describe synaptic transmission",
      "Identify brain regions and functions",
      "Explain reflex arc mechanism",
      "Describe nervous system organization",
      "Analyze neurotransmitter effects",
      "Explain sensory and motor pathways"
    ]
  },
  "Geology": {
    sections: [
      {
        title: "Plate Tectonics",
        fr: "Tectonique des Plaques",
        content: `
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
        `
      },
      {
        title: "Earthquakes & Volcanoes",
        fr: "Tremblements de Terre et Volcans",
        content: `
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
        `
      }
    ],
    examQuestions: [
      "Describe Earth's internal structure",
      "Explain plate tectonic theory",
      "Analyze plate boundary types",
      "Describe earthquake mechanisms",
      "Explain volcanic activity",
      "Analyze rock cycle processes",
      "Predict geological hazards",
      "Interpret geological maps and data"
    ]
  }
};

// Generate comprehensive notes
function generateDeepDiveNotes() {
  const output = [];
  
  for (const [topic, data] of Object.entries(deepDiveTopics)) {
    output.push(`# ${topic} - Deep Dive\n`);
    output.push(`## ${translate(topic)}\n`);
    
    // Sections
    data.sections.forEach(section => {
      output.push(`\n### ${section.title}`);
      output.push(`#### ${section.fr}\n`);
      output.push(section.content);
    });
    
    // Exam Questions
    output.push(`\n## ${translate("Exam Questions")}`);
    output.push(`### ${translate("Possible BAC Questions")}\n`);
    data.examQuestions.forEach((q, i) => {
      output.push(`${i + 1}. ${q}`);
    });
    
    output.push('\n---\n');
  }
  
  return output.join('\n');
}

// Generate HTML version with translations
function generateHTMLVersion() {
  let html = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>BAC Deep Dive - All Topics</title>
  <style>
    body { font-family: Arial, sans-serif; max-width: 1200px; margin: 0 auto; padding: 20px; background: #f5f5f5; }
    .topic { background: white; margin: 20px 0; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
    .topic h2 { color: #2c3e50; border-bottom: 3px solid #3498db; padding-bottom: 10px; }
    .section { margin: 15px 0; padding: 15px; background: #ecf0f1; border-left: 4px solid #3498db; }
    .section h3 { margin-top: 0; color: #2c3e50; }
    .fr { color: #e74c3c; font-style: italic; font-size: 0.9em; }
    .questions { background: #fff3cd; padding: 15px; border-radius: 5px; margin-top: 15px; }
    .questions h4 { color: #856404; }
    .questions ol { margin: 10px 0; }
    .questions li { margin: 5px 0; }
    code { background: #f4f4f4; padding: 2px 6px; border-radius: 3px; }
  </style>
</head>
<body>
  <h1>BAC Examination - Complete Deep Dive</h1>
  <p>All topics with French translations and potential exam questions</p>
`;

  for (const [topic, data] of Object.entries(deepDiveTopics)) {
    html += `<div class="topic">
      <h2>${topic} <span class="fr">${translate(topic)}</span></h2>`;
    
    data.sections.forEach(section => {
      html += `<div class="section">
        <h3>${section.title} <span class="fr">${section.fr}</span></h3>
        <pre>${section.content}</pre>
      </div>`;
    });
    
    html += `<div class="questions">
      <h4>Possible BAC Questions</h4>
      <ol>`;
    
    data.examQuestions.forEach(q => {
      html += `<li>${q}</li>`;
    });
    
    html += `</ol></div></div>`;
  }
  
  html += `</body></html>`;
  return html;
}

// Write files
const mdContent = generateDeepDiveNotes();
const htmlContent = generateHTMLVersion();

fs.writeFileSync('/home/med/Documents/bac/resources/notes/Study-Vault/DEEP-DIVE-ALL-TOPICS.md', mdContent);
fs.writeFileSync('/home/med/Documents/bac/resources/notes/Study-Vault/DEEP-DIVE-ALL-TOPICS.html', htmlContent);

console.log('✓ Deep dive notes generated');
console.log('✓ Markdown: DEEP-DIVE-ALL-TOPICS.md');
console.log('✓ HTML: DEEP-DIVE-ALL-TOPICS.html');
