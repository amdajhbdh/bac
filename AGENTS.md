# BAC Unified - Developer Guide

<!--toc:start-->
- [BAC Unified - Developer Guide](#bac-unified-developer-guide)
  - [Project Vision](#project-vision)
  - [Project Overview](#project-overview)
  - [Build & Run Commands](#build-run-commands)
  - [Architecture](#architecture)
  - [Media-to-Animation Pipeline (Noon)](#media-to-animation-pipeline-noon)
  - [Code Style Guidelines](#code-style-guidelines)
    - [Nushell (CLI)](#nushell-cli)
    - [Rust (Noon)](#rust-noon)
    - [Python (API)](#python-api)
  - [Version Control](#version-control)
<!--toc:end-->

This file provides guidelines for AI agents working on the BAC Unified project.

## Project Vision

**BAC Unified** is a national exam preparation system for Mauritania's BAC C students:

- Country-wide unified question bank
- Multi-modal input (images, PDFs, URLs)
- AI-powered exercise solver with step-by-step solutions
- Question prediction based on historical patterns
- Animation visualizations using **noon** (Rust-based, faster than manim)
- Full audit trail with blockchain

## Project Overview

- **Primary Languages**: Nushell (CLI), Rust (Noon animations), Go/Python (API)
- **CLI Entry Point**: `bac.nu`
- **Animation Engine**: noon (<https://github.com/yongkyuns/noon>)
- **Database**: PostgreSQL + pgvector
- **Storage**: MinIO (local S3)

## Build & Run Commands

```bash
# Run the main CLI
nu bac.nu
nu bac.nu --help

# Validate Nushell syntax
nu --validate bac.nu

# Run Noon animations
cd src/noon && cargo build --release
cargo run --release --example math_graph

# Run Python services
cd src/api && python -m uvicorn main:app --reload

# Database
psql -U postgres -d bac_unified -f sql/schema.sql

# Run all tests
cargo test --manifest-path src/noon/Cargo.toml
```

## Architecture

```
User Input → OCR/Parse → AI Analysis → Solver → Noon Animation → Storage
                                    ↓
                            Question Bank (PostgreSQL)
                                    ↓
                            Prediction Engine
```

## Media-to-Animation Pipeline (Noon)

1. **Input**: Image/PDF/Text of exercise
2. **OCR**: Extract problem text (Tesseract, EasyOCR)
3. **AI Analysis**: Identify math concepts, diagram elements
4. **Solver**: Generate step-by-step solution
5. **Noon Animation**: Create animation from solution
   - Graph plotting
   - Geometry animations
   - Physics simulations
   - Formula reveals

## Code Style Guidelines

### Nushell (CLI)

- Files: `kebab-case`
- Commands: `snake_case`
- Indent: 4 spaces
- Max line: 100 chars

### Rust (Noon)

- Follow Rust standard conventions
- Use `cargo fmt`
- Add doc comments for public functions

### Python (API)

NEVER USE Python UNLESS FOR THINGS THAT IT's the most performant in

- Follow PEP 8
- Use type hints
- Use `black` formatter

## Version Control

- Use **jj** for version control
- Commit regularly with descriptive messages
- Update `progress.md` with changes

---

*Last Updated: 2026-02-23*
