## ADDED Requirements

### Requirement: Template system loads templates from YAML files
The template system SHALL load animation templates from YAML files in a defined directory structure.

#### Scenario: Load math templates
- **WHEN** template manager initializes
- **THEN** it SHALL load all YAML files from templates/math/ directory

#### Scenario: Load physics templates
- **WHEN** template manager initializes
- **THEN** it SHALL load all YAML files from templates/physics/ directory

#### Scenario: Load chemistry templates
- **WHEN** template manager initializes
- **THEN** it SHALL load all YAML files from templates/chemistry/ directory

### Requirement: Template matches problems based on keywords
The template system SHALL match incoming problems to appropriate templates using keyword matching.

#### Scenario: Match by keyword
- **WHEN** problem text contains template keywords
- **THEN** the matcher SHALL return the highest-matching template

#### Scenario: Match by subject
- **WHEN** problem specifies subject (math/physics/chemistry)
- **THEN** the matcher SHALL filter templates by that subject

#### Scenario: No match found
- **WHEN** no template matches the problem
- **THEN** the system SHALL return a generic template

### Template: Function Graph Animation

#### Scenario: Render function graph template
- **WHEN** function_graph template is selected with function "x**2"
- **THEN** the generated code SHALL create axes, plot the function, and add labels

### Template: Equation Solving Animation

#### Scenario: Render equation solving template
- **WHEN** equation_solving template is selected with equation "2x + 5 = 15"
- **THEN** the generated code SHALL show step-by-step solving with highlighted changes

### Template: Physics Mechanics Animation

#### Scenario: Render physics mechanics template
- **WHEN** mechanics template is selected with projectile motion
- **THEN** the generated code SHALL use SpaceScene with gravity simulation

### Template: Chemistry Molecular Animation

#### Scenario: Render chemistry molecule template
- **WHEN** molecule template is selected with "H2O"
- **THEN** the generated code SHALL create 2D or 3D molecular structure
