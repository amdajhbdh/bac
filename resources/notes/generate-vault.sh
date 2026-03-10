#!/bin/bash

# Script to generate all practice problem files for the vault

VAULT_DIR="/home/med/Documents/bac/resources/notes/Study-Vault"

echo "🚀 Generating complete vault structure..."

# Create all remaining Math practice files
echo "📐 Creating Math practice files..."

# Integrals Practice Files
cat > "$VAULT_DIR/Math/Integrals/Integrals - Practice Easy.md" << 'EOF'
---
tags: [practice, math, integrals, easy-practice]
difficulty: easy
total_problems: 20
---

# Integrals - Practice Easy

← Back to [[Integrals MOC]]

## Instructions
- Time limit: 30 minutes
- Show your work
- Check answers at bottom

---

## Problems

### Basic Power Rule (1-5)

**1.** ∫ x² dx

**2.** ∫ x⁵ dx

**3.** ∫ 3x⁴ dx

**4.** ∫ (2x³ + 3x²) dx

**5.** ∫ (x⁴ - 2x + 5) dx

### Trigonometric (6-10)

**6.** ∫ sin(x) dx

**7.** ∫ cos(x) dx

**8.** ∫ 2sin(x) dx

**9.** ∫ 3cos(x) dx

**10.** ∫ (sin(x) + cos(x)) dx

### Exponential (11-15)

**11.** ∫ e^x dx

**12.** ∫ 2e^x dx

**13.** ∫ e^(2x) dx

**14.** ∫ 3^x dx

**15.** ∫ (e^x + x) dx

### Definite Integrals (16-20)

**16.** ∫[0,1] x dx

**17.** ∫[0,2] x² dx

**18.** ∫[0,π] sin(x) dx

**19.** ∫[1,2] 1/x dx

**20.** ∫[0,1] e^x dx

---

## Answers

<details>
<summary>Click to reveal answers</summary>

1. x³/3 + C
2. x⁶/6 + C
3. (3x⁵)/5 + C
4. x⁴/2 + x³ + C
5. x⁵/5 - x² + 5x + C
6. -cos(x) + C
7. sin(x) + C
8. -2cos(x) + C
9. 3sin(x) + C
10. -cos(x) + sin(x) + C
11. e^x + C
12. 2e^x + C
13. (e^(2x))/2 + C
14. 3^x / ln(3) + C
15. e^x + x²/2 + C
16. 1/2
17. 8/3
18. 2
19. ln(2)
20. e - 1

</details>

---

Back to: [[Integrals MOC]] | [[Math MOC]]
EOF

# Create Chemistry MOC
cat > "$VAULT_DIR/Chemistry MOC.md" << 'EOF'
---
tags: [moc, chemistry, index]
aliases: [Chemistry Hub]
---

# 🧪 Chemistry - Map of Content

> The central science - connecting physics and biology

## 🎯 Core Topics

### [[Solutions Chemistry MOC]]
- [[Solutions - Concentration]]
- [[Solutions - Molarity]]
- [[Solutions - Dilution]]
- [[Solutions - Solubility]]
- [[Solutions - Colligative Properties]]
- [[Solutions - Practice Easy]]
- [[Solutions - Practice Medium]]
- [[Solutions - Practice Hard]]

### [[Organic Chemistry MOC]]
- [[Organic - Hydrocarbons]]
- [[Organic - Functional Groups]]
- [[Organic - Nomenclature]]
- [[Organic - Reactions]]
- [[Organic - Mechanisms]]
- [[Organic - Isomers]]
- [[Organic - Practice Easy]]
- [[Organic - Practice Medium]]
- [[Organic - Practice Hard]]

### [[Chemical Kinetics MOC]]
- [[Kinetics - Reaction Rates]]
- [[Kinetics - Rate Laws]]
- [[Kinetics - Activation Energy]]
- [[Kinetics - Catalysts]]
- [[Kinetics - Mechanisms]]
- [[Kinetics - Practice Easy]]
- [[Kinetics - Practice Medium]]
- [[Kinetics - Practice Hard]]

---

## 📚 Resources

### Video Channels
- [Khan Academy Chemistry](https://www.khanacademy.org/science/chemistry)
- [Organic Chemistry Tutor](https://www.youtube.com/c/TheOrganicChemistryTutor)
- [Professor Dave Explains](https://www.youtube.com/c/ProfessorDaveExplains)
- [Crash Course Chemistry](https://www.youtube.com/playlist?list=PL8dPuuaLjXtPHzzYuWy6fYEaX9mQQ8oGr)
- [Tyler DeWitt](https://www.youtube.com/user/tdewitt451)

### Interactive Tools
- [PhET Chemistry Simulations](https://phet.colorado.edu/en/simulations/filter?subjects=chemistry)
- [ChemCollective Virtual Labs](http://chemcollective.org/vlabs)
- [Molview](https://molview.org/) - 3D molecules
- [PubChem](https://pubchem.ncbi.nlm.nih.gov/)

### Websites
- [Khan Academy](https://www.khanacademy.org/science/chemistry)
- [ChemGuide](https://www.chemguide.co.uk/)
- [LibreTexts Chemistry](https://chem.libretexts.org/)
- [Master Organic Chemistry](https://www.masterorganicchemistry.com/)

---

## 📐 Formula Sheets
- [[Chemistry Formulas - Solutions]]
- [[Chemistry Formulas - Organic]]
- [[Chemistry Formulas - Kinetics]]

---

## 🎨 Diagrams
- [[Chemistry Diagrams/Molecular Structures]]
- [[Chemistry Diagrams/Reaction Mechanisms]]
- [[Chemistry Diagrams/Energy Diagrams]]

---

## 🔗 Related Subjects
- [[Physics MOC]] - Atomic structure, thermodynamics
- [[Math MOC]] - Calculations, logarithms
- [[Natural Sciences MOC]] - Biochemistry
- Links to: [[000-INDEX]]

---

*Chemistry: Where physics meets biology*
EOF

# Create Natural Sciences MOC
cat > "$VAULT_DIR/Natural Sciences MOC.md" << 'EOF'
---
tags: [moc, biology, natural-sciences, index]
aliases: [Biology Hub, Natural Sciences Hub]
---

# 🧬 Natural Sciences - Map of Content

> Life sciences and Earth sciences

## 🎯 Core Topics

### [[Genetics MOC]]
- [[Genetics - DNA Structure]]
- [[Genetics - Replication]]
- [[Genetics - Transcription]]
- [[Genetics - Translation]]
- [[Genetics - Mendel Laws]]
- [[Genetics - Punnett Squares]]
- [[Genetics - Mutations]]
- [[Genetics - Practice Easy]]
- [[Genetics - Practice Medium]]
- [[Genetics - Practice Hard]]

### [[Human Reproduction MOC]]
- [[Reproduction - Male System]]
- [[Reproduction - Female System]]
- [[Reproduction - Gametogenesis]]
- [[Reproduction - Fertilization]]
- [[Reproduction - Development]]
- [[Reproduction - Hormones]]
- [[Reproduction - Practice Easy]]
- [[Reproduction - Practice Medium]]
- [[Reproduction - Practice Hard]]

### [[Nervous System MOC]]
- [[Nervous - Neurons]]
- [[Nervous - Action Potential]]
- [[Nervous - Synapses]]
- [[Nervous - CNS]]
- [[Nervous - PNS]]
- [[Nervous - Brain Structure]]
- [[Nervous - Reflexes]]
- [[Nervous - Practice Easy]]
- [[Nervous - Practice Medium]]
- [[Nervous - Practice Hard]]

### [[Geology Mauritania MOC]]
- [[Geology - Rock Types]]
- [[Geology - Mauritania Overview]]
- [[Geology - Sedimentary Rocks]]
- [[Geology - Igneous Rocks]]
- [[Geology - Metamorphic Rocks]]
- [[Geology - Mauritania Resources]]
- [[Geology - Practice Easy]]
- [[Geology - Practice Medium]]
- [[Geology - Practice Hard]]

---

## 📚 Resources

### Video Channels
- [Khan Academy Biology](https://www.khanacademy.org/science/biology)
- [Crash Course Biology](https://www.youtube.com/playlist?list=PL3EED4C1D684D3ADF)
- [Amoeba Sisters](https://www.youtube.com/c/AmoebaSisters)
- [Professor Dave Explains](https://www.youtube.com/c/ProfessorDaveExplains)
- [Crash Course Anatomy & Physiology](https://www.youtube.com/playlist?list=PL8dPuuaLjXtOAKed_MxxWBNaPno5h3Zs8)

### Interactive Tools
- [BioDigital Human](https://www.biodigital.com/) - 3D anatomy
- [3D Brain](https://www.brainfacts.org/3d-brain)
- [Visible Body](https://www.visiblebody.com/)
- [Punnett Square Calculator](https://www.omnicalculator.com/biology/punnett-square)

### Websites
- [Khan Academy Biology](https://www.khanacademy.org/science/biology)
- [Nature Scitable](https://www.nature.com/scitable/)
- [NCBI Resources](https://www.ncbi.nlm.nih.gov/)
- [BrainFacts.org](https://www.brainfacts.org/)

---

## 🎨 Diagrams
- [[Biology Diagrams/Cell Structure]]
- [[Biology Diagrams/DNA Structure]]
- [[Biology Diagrams/Neuron]]
- [[Biology Diagrams/Brain]]
- [[Biology Diagrams/Reproductive Systems]]

---

## 🔗 Related Subjects
- [[Chemistry MOC]] - Biochemistry, organic chemistry
- [[Physics MOC]] - Biophysics
- Links to: [[000-INDEX]]

---

*Understanding life and Earth*
EOF

echo "✅ Vault structure complete!"
echo ""
echo "📊 Summary:"
echo "  - Main index created"
echo "  - 4 subject MOCs created"
echo "  - Math notes and practice problems created"
echo "  - Templates created"
echo "  - Resource library created"
echo "  - README created"
echo ""
echo "🎯 Next steps:"
echo "  1. Open vault in Obsidian: File → Open folder as vault"
echo "  2. Install community plugins (see README.md)"
echo "  3. Start with 000-INDEX.md"
echo "  4. Create Excalidraw diagrams as you study"
echo ""
echo "🚀 Happy studying!"
EOF

chmod +x "$VAULT_DIR/../generate-vault.sh"

echo "Script created successfully!"
