# AGENTS.md - Content Creation and Development Rules

This file contains rules for agents working in this repository. It covers both **content creation** (educational materials) and **code development** (Python utilities).

---

## Part 1: Python Development Rules

### Build/Lint/Test Commands

**Python Version**: 3.11+

**Linting** (uses ruff):
```bash
# Lint all Python files
ruff check .

# Fix auto-fixable issues
ruff check . --fix

# Lint specific file
ruff check agent-system.py
```

**Testing**:
```bash
# Run pytest (if tests exist)
pytest

# Run a single test
pytest tests/test_file.py::test_function_name

# Run with verbose output
pytest -v

# Run tests matching a pattern
pytest -k "test_pattern"
```

**Type Checking** (optional, using mypy):
```bash
mypy .
```

**Format Code**:
```bash
# Format with ruff
ruff format .
```

### Code Style Guidelines

**Imports**: Group: stdlib → third-party → local; avoid wildcards; sort alphabetically

**Formatting**: 100 chars max, 4 spaces, trailing commas, one blank line between top-level

**Types**: Use type hints, prefer `Optional[X]`, concrete types for public APIs

**Naming**: Classes `PascalCase`, funcs/vars `snake_case`, constants `SCREAMING_SNAKE_CASE`, private `_prefix`

**Error Handling**: Specific exceptions, meaningful messages, never silent catch

```python
def process_file(file_path: Path, output_dir: Path) -> List[str]:
    """Process a single file and return results."""
    if not file_path.exists():
        raise FileNotFoundError(f"Input file not found: {file_path}")
    results: List[str] = []
    try:
        content = file_path.read_text(encoding="utf-8")
        results = parse_content(content)
    except UnicodeDecodeError as e:
        raise ValueError(f"Invalid encoding in {file_path}") from e
    return results
```

---

## Part 2: Mathematical Notation Standards

### ✅ **REQUIRED: Use LaTeX for ALL mathematical expressions**

**Inline math**: Use `$...$` for inline expressions
```markdown
The pH is calculated as $pH = -\log[H_3O^+]$
```

**Display math**: Use `$$...$$` for centered equations
```markdown
$$K_a = \frac{[H_3O^+][A^-]}{[AH]}$$
```

### ❌ **FORBIDDEN: ASCII/Unicode math symbols**

**Subscripts (FORBIDDEN):** H₂O, CO₂, NH₃, x₀, x₁ → Use `$H_2O$`, `$CO_2$`, `$NH_3$`, `$x_0$`, `$x_1$`

**Superscripts (FORBIDDEN):** x², x³, e^x, 10⁻⁴, Ca²⁺ → Use `$x^2$`, `$x^3$`, `$e^x$`, `$10^{-4}$`, `$Ca^{2+}$`

**Greek letters (FORBIDDEN):** α, β, γ, δ, θ, λ, μ, π, σ, φ, ω → Use `$\alpha$`, `$\beta$`, etc.

**Special symbols (FORBIDDEN):** ≤, ≥, ≠, ±, ×, ÷, ∑, ∫, √, ∞ → Use `$\leq$`, `$\geq$`, `$\neq$`, `$\pm$`, etc.

### **Common LaTeX Conversions**

| ASCII/Unicode | LaTeX |
|---------------|-------|
| H₂SO₄ | `$H_2SO_4$` |
| 10⁻¹⁴ | `$10^{-14}$` |
| ∫₀¹ f(x)dx | `$\int_0^1 f(x)dx$` |

---

## Part 3: Content Organization Rules

**File Structure**: `Topic-Subtopic.md`, include level `(7ème/Terminale)`, add source `*Source: ...*`

**Section Headers**:
```markdown
# Main Topic - Subtopic (Level)
## Concepts Clés / Key Concepts  
### Specific Topic
```

**Exercise Format**:
```markdown
## Exercice N: Title

### Énoncé
[Problem with LaTeX math]

### Méthode de résolution
1. **Step 1**: Description
   $$equation$$

### Résultat
$$final\_answer$$
```

---

## Part 4: Quality Standards

### **Mathematical Rigor**
- Always show units: `$C = 0.1 \text{ mol/L}$`
- Use proper notation: `$[H_3O^+]$` not `[H3O+]`
- Include equilibrium arrows: `$\rightleftharpoons$` for reversible reactions
- Use proper chemical formulas: `$CH_3COOH$` not `CH3COOH`

### **Language Standards**
- **French**: Use for Mauritanian BAC content
- **English**: Use for international/general concepts
- **Bilingual**: Provide both when helpful

### **Content Levels**
- **Basic (6ème)**: Simple concepts, basic calculations
- **Intermediate (7ème C/D)**: Complex problems, multiple steps
- **Advanced (7ème/Terminale)**: University-prep level, theoretical depth

---

## Part 5: Automated Checks

1. **Math Check**: Scan for ASCII math symbols
2. **LaTeX Validation**: Ensure all `$...$` and `$$...$$` are properly closed
3. **Level Verification**: Confirm content matches stated difficulty
4. **Source Attribution**: Include original source information
5. **Code Lint**: Run `ruff check .` on any modified Python files

```bash
grep -r "[₀-₉⁺⁻²³¹⁴⁵⁆⁷⁸⁹]" *.md
sed -i 's/H₂O/$H_2O$/g' *.md
```

---

## Part 6: Agent Responsibilities

- **Content Agents**: Convert ASCII math to LaTeX, maintain source attribution
- **QA Agents**: Scan for non-LaTeX notation, verify equation formatting
- **Code Agents**: Follow Part 1 guidelines, run linting, add type hints
- **Update Agents**: Apply formatting rules, maintain consistency

---

*This document must be followed by all content creation and modification agents*
