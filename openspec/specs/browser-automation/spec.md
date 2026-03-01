## ADDED Requirements

### Requirement: Type problem into input
The system SHALL type text into the chat input field of the AI service.

#### Scenario: Input field present
- **WHEN** page loads and input field is visible
- **THEN** system types the problem text into the input field

#### Scenario: Input field not ready
- **WHEN** page loads but input field is not immediately visible
- **THEN** system waits up to 10 seconds for input field, fails gracefully if not found

### Requirement: Submit prompt
The system SHALL click the submit/send button to send the prompt to the AI.

#### Scenario: Submit button exists
- **WHEN** problem is typed into input
- **THEN** system clicks the submit button to send the message

#### Scenario: Submit via keyboard
- **WHEN** submit button click fails
- **THEN** system tries pressing Enter to submit

### Requirement: Wait for response
The system SHALL wait for the AI to generate a response before extracting content.

#### Scenario: Response appears
- **WHEN** prompt is submitted
- **THEN** system waits for response text to appear, up to 60 seconds

#### Scenario: Response timeout
- **WHEN** response takes longer than 60 seconds
- **THEN** system returns partial result or timeout error

### Requirement: Take snapshot for extraction
The system SHALL capture a page snapshot to extract response content.

#### Scenario: Snapshot requested
- **WHEN** system needs to extract response
- **THEN** playwright-cli snapshot is called to capture current DOM state
