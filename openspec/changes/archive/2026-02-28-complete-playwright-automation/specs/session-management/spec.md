## ADDED Requirements

### Requirement: Persistent session per service
The system SHALL maintain separate persistent browser sessions for each AI service.

#### Scenario: New session requested
- **WHEN** automation starts for a service
- **THEN** system uses playwright-cli with --persistent flag and service-specific profile directory

#### Scenario: Session already exists
- **WHEN** automation starts for a service with existing profile
- **THEN** system reuses the existing authenticated session

### Requirement: Session health check
The system SHALL verify that a session is still authenticated before attempting automation.

#### Scenario: Session authenticated
- **WHEN** page loads and shows chat interface
- **THEN** system proceeds with automation

#### Scenario: Session expired
- **WHEN** page redirects to login or shows auth error
- **THEN** system returns authentication error, suggests manual re-authentication

### Requirement: Clean session state
The system SHALL clear conversation history before new problem to avoid context contamination.

#### Scenario: New conversation needed
- **WHEN** starting new problem solve
- **THEN** system creates new conversation/chat if possible, or clears existing

#### Scenario: Cannot clear history
- **WHEN** service doesn't support clearing history
- **THEN** system proceeds with existing history, notes in response metadata
