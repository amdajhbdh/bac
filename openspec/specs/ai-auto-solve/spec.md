## ADDED Requirements

### Requirement: Auto-solve with AI service
The system SHALL automatically send a problem to a configured AI service and return the solution text.

#### Scenario: Successful solve
- **WHEN** user requests online research for a problem
- **THEN** system opens the configured AI service in Chrome, types the problem, submits it, and extracts the response

#### Scenario: Service rate limited
- **WHEN** AI service returns rate limit error
- **THEN** system waits with exponential backoff and retries up to 3 times

#### Scenario: All services fail
- **WHEN** all AI services fail to respond
- **THEN** system returns failure with error message, allows manual browser interaction

### Requirement: Multi-service fallback
The system SHALL try multiple AI services in order when one fails.

#### Scenario: Primary service fails
- **WHEN** DeepSeek fails (rate limit, error)
- **THEN** system automatically tries Grok, then Claude, then ChatGPT

#### Scenario: Successful fallback
- **WHEN** first service fails but second succeeds
- **THEN** system returns the successful result with service identifier

### Requirement: Solution extraction
The system SHALL extract the AI's response text from the chat interface.

#### Scenario: Standard response
- **WHEN** AI service returns a text response
- **THEN** system extracts the response text from the message area

#### Scenario: Response with formatting
- **WHEN** AI response contains markdown or code blocks
- **THEN** system extracts raw text, preserving line breaks
