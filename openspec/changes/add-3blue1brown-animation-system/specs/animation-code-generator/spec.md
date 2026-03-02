## ADDED Requirements

### Requirement: Code generator produces valid Manim Python code
The code generator SHALL produce syntactically valid Manim Python code from problem descriptions and solutions.

#### Scenario: Generate code from math problem
- **WHEN** a math problem with solution is provided
- **THEN** the generator SHALL produce Python code with a valid Manim Scene class

#### Scenario: Generate code from physics problem
- **WHEN** a physics problem is provided
- **THEN** the generator SHALL produce code using manim-physics classes

#### Scenario: Generate code from chemistry problem
- **WHEN** a chemistry problem is provided
- **THEN** the generator SHALL produce code using manim-chemistry classes

### Requirement: Code generator uses template matching
The code generator SHALL match problems to appropriate templates based on keywords and subject.

#### Scenario: Match function graphing problem
- **WHEN** problem contains "graph", "function", "courbe", "f(x)"
- **THEN** the generator SHALL use the function_graph template

#### Scenario: Match equation solving problem
- **WHEN** problem contains "solve", "résoudre", "equation"
- **THEN** the generator SHALL use the equation_solving template

#### Scenario: Fallback to generic template
- **WHEN** no keyword match is found
- **THEN** the generator SHALL use a generic animated explanation template

### Requirement: Code generator uses AI enhancement
The code generator SHALL use Ollama to enhance or generate custom Manim code when templates are insufficient.

#### Scenario: AI generates custom animation
- **WHEN** template match confidence is low
- **THEN** the generator SHALL call Ollama with a prompt to generate custom code

#### Scenario: AI code validation
- **WHEN** Ollama returns generated code
- **THEN** the validator SHALL check for required imports and Scene class

#### Scenario: Fallback on AI failure
- **WHEN** Ollama fails or returns invalid code
- **THEN** the generator SHALL fall back to template-based code

### Requirement: Code generator extracts variables from problems
The code generator SHALL extract relevant variables (functions, equations, values) from problem text.

#### Scenario: Extract mathematical function
- **WHEN** problem contains "f(x) = x^2" or similar
- **THEN** the extractor SHALL identify "x**2" as the function expression

#### Scenario: Extract geometric parameters
- **WHEN** problem describes a circle with radius
- **THEN** the extractor SHALL identify radius value
