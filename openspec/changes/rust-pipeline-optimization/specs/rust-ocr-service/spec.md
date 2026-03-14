# Rust OCR Service Specification

## Overview

High-performance OCR service written in Rust using leptess (Tesseract bindings) with parallel processing via Rayon. Provides 5x faster OCR than Go+exec approach.

## Problem Statement

- **Current**: Go calls tesseract via exec, processes one image at a time
- **Bottleneck**: 3-5 seconds per image, no batching, no parallelism
- **Goal**: <1 second per image with batch processing

## Solution Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   Rust OCR Service                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │  gRPC API   │───▶│  Worker     │───▶│  Results    │     │
│  │  (Go ←→)    │    │  Pool       │    │  Aggregator │     │
│  └─────────────┘    │  (Rayon)    │    └─────────────┘     │
│                     └─────────────┘                         │
│                          │                                  │
│        ┌─────────────────┼─────────────────┐              │
│        ▼                 ▼                 ▼              │
│  ┌──────────┐      ┌──────────┐      ┌──────────┐        │
│  │ Tesseract │      │ Tesseract│      │ Tesseract│        │
│  │ (image 1)│      │ (image 2)│      │ (image N)│        │
│  └──────────┘      └──────────┘      └──────────┘        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## ADDED Requirements

### Requirement: OCR service accepts image input via gRPC
The Rust OCR service SHALL accept image data through a gRPC interface and return structured OCR results.

#### Scenario: Process single image successfully
- **WHEN** client sends image bytes via gRPC `ProcessImage` RPC
- **THEN** service returns OCRResult with text, confidence, and processing time

#### Scenario: Process PDF with multiple pages
- **WHEN** client sends PDF via gRPC `ProcessPDF` RPC
- **THEN** service processes all pages in parallel and returns aggregated results

#### Scenario: Image too large returns error
- **WHEN** client sends image larger than 50MB
- **THEN** service returns error with "image too large" message

### Requirement: Parallel processing with Rayon worker pool
The OCR service SHALL process multiple images concurrently using a configurable Rayon thread pool.

#### Scenario: Batch of 10 images processed in parallel
- **WHEN** client sends batch of 10 images
- **THEN** images are processed concurrently (up to worker pool size)
- **AND** total time is approximately 1/N of sequential time

#### Scenario: Worker pool scales based on load
- **WHEN** queue depth exceeds 100
- **THEN** service adds workers up to max_workers (8)
- **AND** when queue depth drops below 10, workers are removed

### Requirement: Early exit on high confidence
The OCR service SHALL return immediately when any OCR method achieves confidence > 0.9.

#### Scenario: High confidence triggers early exit
- **WHEN** Tesseract returns confidence > 0.9
- **THEN** remaining OCR methods are cancelled
- **AND** result is returned immediately

### Requirement: Fallback to multiple OCR engines
The service SHALL try multiple OCR engines (Tesseract LSTM, legacy) and return the best result.

#### Scenario: Primary engine fails, fallback succeeds
- **WHEN** Tesseract LSTM fails to extract text
- **THEN** service falls back to legacy Tesseract engine
- **AND** returns result from fallback engine

### Requirement: Performance metrics exposed
The service SHALL expose Prometheus metrics for monitoring.

#### Scenario: Metrics endpoint accessible
- **WHEN** client requests /metrics endpoint
- **THEN** service returns metrics in Prometheus format including:
  - ocr_requests_total (counter)
  - ocr_processing_duration_seconds (histogram)
  - ocr_batch_size (gauge)
  - ocr_errors_total (counter)

## Interface Definition

### gRPC Service

```protobuf
service OCRService {
  rpc ProcessImage(ImageRequest) returns (OCRResult);
  rpc ProcessPDF(PDFRequest) returns (PDFResult);
  rpc ProcessBatch(BatchRequest) returns (BatchResult);
  rpc GetMetrics(Empty) returns (Metrics);
}

message ImageRequest {
  bytes image_data = 1;
  string language = 2;  // "fra+eng" default
}

message OCRResult {
  string text = 1;
  double confidence = 2;
  string source = 3;     // "tesseract-lstm", "tesseract-legacy"
  int64 processing_time_ms = 4;
}

message PDFRequest {
  bytes pdf_data = 1;
  bool extract_pages_parallel = 2;  // default true
}

message PDFResult {
  repeated OCRResult pages = 1;
  int32 total_pages = 2;
}

message BatchRequest {
  repeated ImageRequest images = 1;
}

message BatchResult {
  repeated OCRResult results = 1;
  int64 total_time_ms = 2;
}
```

## Configuration

```yaml
ocr:
  workers: 8
  max_queue_size: 1000
  batch_timeout_ms: 50
  batch_size: 10
  early_exit_confidence: 0.9
  languages:
    - fra
    - eng
    - ara
```

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Single image OCR | < 1s | p95 latency |
| Batch 10 images | < 5s | total time |
| PDF 100 pages | < 30s | total time |
| Throughput | 100 img/min | sustained |
| Memory usage | < 2GB | peak RSS |
| Error rate | < 0.1% | failed requests |

## Acceptance Criteria

1. **Functional**: Service correctly extracts French/Arabic text from images
2. **Performance**: Single image < 1s, batch 10 < 5s (measured)
3. **Reliability**: Graceful degradation when OCR engines fail
4. **Monitoring**: Metrics exposed at /metrics
5. **Integration**: Go client can call service via gRPC
