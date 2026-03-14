# 🎬 Manim Animations for Study Vault

**Auto-generated animations based on actual exercise content**

## 📊 Available Animations

### Biology (4 scenes)
- `HormoneRegulation` - FSH/LH hormone levels visualization
- `NeuronSynapse` - Synaptic transmission with neurotoxins
- `GeneticCrossing` - Drosophila genetic crossing
- `ImmuneResponse` - B lymphocyte to plasma cell differentiation

### Physics (5 scenes)
- `ProjectileMotionDetailed` - Projectile with velocity components
- `SpringOscillationDetailed` - Spring oscillation with energy
- `ElectricFieldLines` - Electric field between charges
- `DynamicsForces` - Newton's F=ma demonstration
- `MagneticForce` - Laplace force (coming soon)

### Mathematics (3 scenes)
- `ComplexNumbersPlane` - Complex numbers on complex plane
- `MatrixOperations` - Matrix multiplication step-by-step
- `IntegralCalculation` - Definite integral with Riemann sums

### Chemistry (2 scenes)
- `ChemicalReactionMechanism` - Reaction with energy diagram
- `OrganicChemistry` - Alcohol functional groups

## 🚀 Quick Start

### Install Manim
```bash
pip install manim
```

### Generate Single Animation
```bash
cd ~/Documents/bac/resources/notes/Study-Vault/Animations
manim -pql generate_all_animations.py HormoneRegulation
```

### Generate All Animations
```bash
# Low quality (fast)
manim -ql generate_all_animations.py

# High quality (slow)
manim -qh generate_all_animations.py
```

## 📝 Command Options

- `-ql` : Low quality (480p, fast)
- `-qm` : Medium quality (720p)
- `-qh` : High quality (1080p)
- `-qk` : 4K quality (2160p)
- `-p` : Preview after rendering
- `-s` : Save last frame as image

## 🎯 Generate Specific Subjects

### Biology Only
```bash
manim -pql generate_all_animations.py HormoneRegulation NeuronSynapse GeneticCrossing ImmuneResponse
```

### Physics Only
```bash
manim -pql generate_all_animations.py ProjectileMotionDetailed SpringOscillationDetailed ElectricFieldLines DynamicsForces
```

### Math Only
```bash
manim -pql generate_all_animations.py ComplexNumbersPlane MatrixOperations IntegralCalculation
```

### Chemistry Only
```bash
manim -pql generate_all_animations.py ChemicalReactionMechanism OrganicChemistry
```

## 📂 Output Location

Videos will be saved in:
```
Animations/media/videos/generate_all_animations/[quality]/
```

## 🔧 Troubleshooting

### Error: "manim: command not found"
```bash
pip install --upgrade manim
```

### Error: "LaTeX not found"
```bash
# Ubuntu/Debian
sudo apt install texlive texlive-latex-extra

# macOS
brew install --cask mactex
```

### Slow rendering?
Use `-ql` for low quality during testing:
```bash
manim -ql generate_all_animations.py SceneName
```

## 📚 Linking Videos to Notes

After generating, videos are automatically linked in:
- `Organized-Notes/[Subject]/[Topic].md`

Example:
```markdown
## Animation
![Hormone Regulation](../../Animations/media/videos/.../HormoneRegulation.mp4)
```

## 🎨 Customization

Edit `generate_all_animations.py` to:
- Change colors
- Adjust timing
- Add more scenes
- Modify animations

## 📊 Progress Tracking

- [ ] Biology animations (4/4)
- [ ] Physics animations (5/7)
- [ ] Math animations (3/5)
- [ ] Chemistry animations (2/4)

**Total: 14/20 scenes ready**

## 🚀 Batch Generation Script

Create `generate_all.sh`:
```bash
#!/bin/bash
cd ~/Documents/bac/resources/notes/Study-Vault/Animations

echo "Generating all animations..."

# Biology
manim -ql generate_all_animations.py HormoneRegulation
manim -ql generate_all_animations.py NeuronSynapse
manim -ql generate_all_animations.py GeneticCrossing
manim -ql generate_all_animations.py ImmuneResponse

# Physics
manim -ql generate_all_animations.py ProjectileMotionDetailed
manim -ql generate_all_animations.py SpringOscillationDetailed
manim -ql generate_all_animations.py ElectricFieldLines
manim -ql generate_all_animations.py DynamicsForces

# Math
manim -ql generate_all_animations.py ComplexNumbersPlane
manim -ql generate_all_animations.py MatrixOperations
manim -ql generate_all_animations.py IntegralCalculation

# Chemistry
manim -ql generate_all_animations.py ChemicalReactionMechanism
manim -ql generate_all_animations.py OrganicChemistry

echo "✓ All animations generated!"
```

Run with: `bash generate_all.sh`

---

**Note:** Animations are based on actual exercises from your study materials. Debug and customize as needed!
