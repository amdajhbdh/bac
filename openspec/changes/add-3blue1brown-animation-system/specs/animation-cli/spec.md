## ADDED Requirements

### Requirement: CLI provides animation generation command
The CLI SHALL provide a command to generate animations from the terminal.

#### Scenario: bac animate <problem>
- **WHEN** user runs "bac animate" with problem text
- **THEN** CLI SHALL generate animation and display video path

#### Scenario: bac animate with options
- **WHEN** user runs "bac animate --subject math --quality medium"
- **THEN** CLI SHALL use specified subject and quality

### Requirement: CLI provides animation listing
The CLI SHALL provide commands to list and manage animations.

#### Scenario: bac animations list
- **WHEN** user runs "bac animations list"
- **THEN** CLI SHALL display table of animations with status

#### Scenario: bac animate get <id>
- **WHEN** user runs "bac animate get" with ID
- **THEN** CLI SHALL display animation details and video URL

### Requirement: CLI provides template commands
The CLI SHALL provide commands to list available templates.

#### Scenario: bac templates list
- **WHEN** user runs "bac templates list"
- **THEN** CLI SHALL display available templates with descriptions

#### Scenario: bac templates show <id>
- **WHEN** user runs "bac templates show" with template ID
- **THEN** CLI SHALL display template details and code preview

### Requirement: CLI supports animation with narration
The CLI SHALL support generating animations with voice narration.

#### Scenario: bac animate --audio
- **WHEN** user runs "bac animate" with --audio flag
- **THEN** CLI SHALL generate animation with voice narration
