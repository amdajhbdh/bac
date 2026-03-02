## ADDED Requirements

### Requirement: REST API provides animation generation endpoint
The REST API SHALL provide an endpoint to generate animations from problems.

#### Scenario: POST /animations/generate
- **WHEN** POST request with problem, solution, subject, quality is received
- **THEN** the API SHALL queue animation and return job ID

#### Scenario: Sync generation
- **WHEN** sync=true in request
- **THEN** the API SHALL wait for completion and return video URL

### Requirement: REST API provides animation retrieval
The REST API SHALL provide endpoints to retrieve generated animations.

#### Scenario: GET /animations/:id
- **WHEN** GET request with valid animation ID
- **THEN** the API SHALL return animation metadata including video URL

#### Scenario: GET /animations/:id/stream
- **WHEN** GET request for streaming
- **THEN** the API SHALL stream video file with appropriate content-type

#### Scenario: GET /animations
- **WHEN** GET request without ID
- **THEN** the API SHALL return paginated list of animations

### Requirement: REST API provides template management
The REST API SHALL provide endpoints to list and manage animation templates.

#### Scenario: GET /templates
- **WHEN** GET request to /templates
- **THEN** the API SHALL return list of available templates

#### Scenario: GET /templates/:id
- **WHEN** GET request with template ID
- **THEN** the API SHALL return template details including code template

### Requirement: REST API provides queue status
The REST API SHALL provide endpoints to check rendering queue status.

#### Scenario: GET /queue/status
- **WHEN** GET request to queue status
- **THEN** the API SHALL return queue depth, processing count, average wait time

### Requirement: REST API returns errors properly
The REST API SHALL return appropriate HTTP status codes and error messages.

#### Scenario: Invalid request body
- **WHEN** POST with missing required fields
- **THEN** API SHALL return 400 with validation errors

#### Scenario: Animation not found
- **WHEN** GET /animations/:id with unknown ID
- **THEN** API SHALL return 404

#### Scenario: Internal error
- **WHEN** Rendering fails internally
- **THEN** API SHALL return 500 with error message
