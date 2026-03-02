# Pipeline Orchestration Specification (MODIFIED)

## Overview

Hybrid Go + Rust pipeline orchestrator that coordinates all stages: OCR (Rust), Solver (Go), Storage (Go). Provides streaming progress updates and graceful error handling.

## MODIFIED Requirements

### Requirement: Hybrid pipeline execution
The orchestrator SHALL coordinate Rust OCR service and Go solver/storage with proper error handling.

#### Scenario: Full pipeline executes successfully
- **WHEN** document submitted to pipeline
- **THEN** stages execute in order: OCR → Solve → Store
- **AND** results returned with all metadata

#### Scenario: OCR fails, fallback to manual
- **WHEN** OCR service returns error
- **THEN** pipeline logs error
- **AND** returns partial result with error flag

#### Scenario: Solver fails, use fallback
- **WHEN** solver returns error
- **THEN** fallback solver is invoked
- **AND** result includes fallback indicator

### Requirement: Streaming progress updates
The pipeline SHALL stream progress updates to client.

#### Scenario: Progress streamed during execution
- **WHEN** pipeline processing document
- **THEN** progress events streamed to client
- **AND** client shows progress bar

#### Scenario: Stage completion events
- **WHEN** each stage completes
- **THEN** event emitted with stage name and duration
- **AND** event includes partial results if applicable

### Requirement: Concurrent stage execution
The orchestrator SHALL run independent stages concurrently when possible.

#### Scenario: OCR and memory lookup run in parallel
- **WHEN** document submitted
- **THEN** OCR starts immediately
- **AND** memory lookup starts immediately (no dependency)
- **AND** solver waits for both to complete

#### Scenario: Storage runs async
- **WHEN** solver completes
- **THEN** result returned to client immediately
- **AND** storage runs in background

### Requirement: Circuit breaker for Rust service
The orchestrator SHALL use circuit breaker pattern for Rust OCR calls.

#### Scenario: Circuit opens after failures
- **WHEN** OCR fails 5 times in 30 seconds
- **THEN** circuit opens
- **AND** subsequent calls fail fast

#### Scenario: Circuit half-open
- **WHEN** circuit open for 60 seconds
- **THEN** test call allowed through
- **AND** if success, circuit closes

### Requirement: Graceful shutdown
The orchestrator SHALL shutdown cleanly, draining in-flight requests.

#### Scenario: Shutdown during processing
- **WHEN** SIGTERM received
- **THEN** new requests rejected
- **AND** in-flight requests complete
- **AND** resources cleaned up

### Requirement: Metrics and observability
The orchestrator SHALL expose comprehensive pipeline metrics.

#### Scenario: Pipeline metrics accessible
- **WHEN** client queries /metrics
- **THEN** returns:
  - pipeline_requests_total (counter)
  - pipeline_stage_duration_seconds (histogram)
  - pipeline_errors_total (counter)
  - pipeline_queue_depth (gauge)
  - pipeline_active (gauge)

## ADDED Requirements

### Requirement: Backpressure propagation
The orchestrator SHALL propagate backpressure from slow stages.

#### Scenario: Slow storage causes backpressure
- **WHEN** storage queue > 500
- **THEN** solver stage slowed
- **AND** OCR stage slowed
- **AND** no message loss

### Requirement: Health checks
The orchestrator SHALL expose health endpoint.

#### Scenario: Health check returns status
- **WHEN** client GET /health
- **THEN** returns:
  - overall: "healthy" | "degraded" | "unhealthy"
  - stages: {ocr: "up", solver: "up", storage: "up"}
  - queue depths
  - last error

## Interface

```go
type PipelineOrchestrator struct {
    ocrClient  *ocr.Client
    solver     *solver.Service
    storage    *storage.Pipeline
    circuitBreaker *circuitbreaker.CircuitBreaker
}

type PipelineRequest struct {
    ID          string
    InputType   string  // "image", "pdf", "url"
    Data        []byte
    Language    string
    Options     PipelineOptions
    Progress    chan PipelineProgress
}

type PipelineProgress struct {
    Stage   string  // "ocr", "solve", "store"
    Status  string  // "started", "completed", "failed"
    Percent int     // 0-100
    Result  interface{}
    Error   error
}

func NewPipelineOrchestrator(cfg Config) (*PipelineOrchestrator, error)

func (p *PipelineOrchestrator) Process(ctx context.Context, req PipelineRequest) (*PipelineResult, error)

func (p *PipelineOrchestrator) ProcessStream(ctx context.Context, input io.Reader) (<-chan PipelineResult, error)

func (p *PipelineOrchestrator) Health() HealthStatus

func (p *PipelineOrchestrator) Shutdown(ctx context.Context) error
```

## Configuration

```yaml
pipeline:
  workers:
    ocr: 4
    solver: 8
    storage: 4
  timeouts:
    ocr: 30s
    solver: 60s
    storage: 30s
  retries:
    ocr: 3
    solver: 2
    storage: 3
  circuit_breaker:
    failure_threshold: 5
    recovery_timeout: 30s
    half_open_requests: 3
```

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| End-to-end latency | < 10s | p95 |
| Stage overlap | > 50% | time savings |
| Error recovery | < 5s | p95 |
| Shutdown drain | < 30s | max |
| Throughput | 100/min | sustained |

## Acceptance Criteria

1. **Hybrid**: Coordinates Rust OCR and Go solver
2. **Streaming**: Progress updates via channels
3. **Concurrent**: Parallel stage execution
4. **Resilient**: Circuit breaker, retries, fallbacks
5. **Observable**: Metrics and health checks
