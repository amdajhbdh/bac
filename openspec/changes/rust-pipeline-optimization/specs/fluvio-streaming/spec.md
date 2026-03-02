# Fluvio Streaming Specification

## Overview

Event-driven streaming pipeline using Fluvio for real-time document processing. Provides automatic backpressure, exactly-once processing, and WASM-based transforms.

## Problem Statement

- **Current**: Sequential processing, no streaming, blocks on each stage
- **Bottleneck**: Can't handle burst traffic, no backpressure, no ordering guarantees
- **Goal**: Real-time streaming with backpressure and ordering

## Solution Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Fluvio Streaming Pipeline                   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌───────┐ │
│  │ocr-input │───▶│ocr-output│───▶│solver-   │───▶│storage│ │
│  │ (topic)  │    │ (topic)  │    │output   │    │(topic)│ │
│  └──────────┘    └──────────┘    │ (topic) │    └───────┘ │
│                                   └──────────┘              │
│                                                              │
│  Partitions: 8 per topic                                     │
│  Replication: 3                                              │
│                                                              │
│  WASM Transforms:                                            │
│  - ocr-preprocess: Image validation, format conversion       │
│  - solver-enrich: Add context from memory lookup            │
│  - format-storage: Prepare for DB insert                   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## ADDED Requirements

### Requirement: Topics accept document input
The streaming system SHALL accept documents through input topics and process them in order.

#### Scenario: Single document processed through pipeline
- **WHEN** document is sent to `ocr-input` topic
- **THEN** document flows through all stages
- **AND** result is stored in database

#### Scenario: Documents processed in order
- **WHEN** documents A, B, C sent to topic in order
- **THEN** documents are processed in FIFO order
- **AND** results appear in order

### Requirement: Automatic backpressure
The system SHALL apply backpressure when consumers can't keep up with producers.

#### Scenario: Consumer slow triggers backpressure
- **WHEN** consumer processing time > 100ms while producer sends at 10ms intervals
- **THEN** Fluvio automatically slows producer
- **AND** no message loss occurs

#### Scenario: Consumer recovers after backpressure
- **WHEN** consumer catches up
- **THEN** producer speed returns to normal
- **AND** no backlog remains

### Requirement: WASM transforms preprocess data
The pipeline SHALL use WASM transforms for efficient preprocessing.

#### Scenario: Image validated before OCR
- **WHEN** image sent to pipeline
- **THEN** WASM transform validates format and size
- **AND** invalid images are marked with error

#### Scenario: Context enriched before solver
- **WHEN** OCR result ready for solver
- **THEN** transform adds similar problems from memory
- **AND** solver receives enriched context

### Requirement: Exactly-once processing
The system SHALL process each message exactly once (no duplicates, no drops).

#### Scenario: Consumer restarts after crash
- **WHEN** consumer crashes and restarts
- **THEN** already-processed messages are not re-processed
- **AND** in-progress messages are retried

### Requirement: Metrics and monitoring
The system SHALL expose pipeline metrics for monitoring.

#### Scenario: Pipeline health check
- **WHEN** client queries /health
- **THEN** returns topic lag, consumer position, error count

#### Scenario: Pipeline latency tracked
- **WHEN** message flows through pipeline
- **THEN** latency is tracked per stage
- **AND** metrics exposed for Grafana

## Topic Configuration

```yaml
topics:
  - name: ocr-input
    partitions: 8
    replication: 3
    retention: 24h
    compression: gzip
    
  - name: ocr-output
    partitions: 8
    replication: 3
    
  - name: solver-input
    partitions: 8
    replication: 3
    
  - name: solver-output
    partitions: 8
    replication: 3
    
  - name: storage-input
    partitions: 4
    replication: 3
```

## Message Schema

```json
// ocr-input message
{
  "id": "uuid",
  "input_type": "image|pdf|url",
  "data": "base64 or url",
  "metadata": {
    "source": "user upload",
    "language": "fra"
  },
  "timestamp": "2026-03-02T12:00:00Z"
}

// ocr-output message
{
  "id": "uuid",
  "parent_id": "uuid",
  "text": "extracted text",
  "confidence": 0.95,
  "source": "tesseract-lstm",
  "processing_time_ms": 500,
  "timestamp": "2026-03-02T12:00:01Z"
}

// solver-output message
{
  "id": "uuid", 
  "parent_id": "uuid",
  "problem": "original text",
  "solution": "solved answer",
  "subject": "math",
  "concepts": ["derivative", "calculus"],
  "confidence": 0.85,
  "model": "llama3.2:3b",
  "timestamp": "2026-03-02T12:00:03Z"
}
```

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| End-to-end latency | < 5s | p95 |
| Topic throughput | 1000 msg/s | per partition |
| Consumer lag | < 100 | messages behind |
| Backpressure response | < 1s | time to detect |
| Error rate | < 0.1% | failed messages |

## Acceptance Criteria

1. **Functional**: Documents flow through all pipeline stages
2. **Ordering**: FIFO processing maintained
3. **Backpressure**: Producer slows when consumer overwhelmed
4. **Exactly-once**: No duplicate processing after crash
5. **Monitoring**: Health and latency metrics exposed
