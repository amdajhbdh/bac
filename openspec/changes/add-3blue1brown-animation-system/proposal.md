## Why

The current animation system in BAC Unified uses a placeholder Noon (Rust) engine that generates only basic text-based animations. This is insufficient for the educational goals of the platform - students need 3Blue1Brown-style visual explanations that animate mathematical concepts, physics simulations, and chemistry visualizations with professional quality. The current system cannot generate graphs, function plots, geometric proofs, or any meaningful visual learning content.

## What Changes

- Replace the empty Noon placeholder with **Manim Community Edition** (Python) as the primary animation engine
- Implement **AI-powered code generation** using Ollama to convert problem solutions into Manim animations
- Create a **template system** for math, physics, and chemistry animations
- Add **manim-physics** plugin for real physics simulations (mechanics, waves, electromagnetism)
- Add **manim-chemistry** plugin for molecular visualizations
- Implement a **Go-Python bridge** to integrate Manim rendering into the Go-based agent
- Build a **rendering pipeline** with quality presets, caching, and FFmpeg video processing
- Add optional **voice narration** using Coqui TTS for 3Blue1Brown-style explanations
- Extend the database schema for animation templates, rendering jobs, and metadata

## Capabilities

### New Capabilities

- **animation-engine**: Core Manim integration for rendering mathematical animations
- **animation-code-generator**: AI-powered conversion of problem solutions to Manim Python code
- **animation-templates**: Subject-specific templates (math, physics, chemistry) for common problem types
- **animation-rendering**: Complete rendering pipeline with quality presets and video output
- **animation-caching**: Multi-layer caching system (memory, disk, database)
- **animation-voice**: Optional voice narration synthesis using Coqui TTS
- **animation-api**: REST API endpoints for animation generation, management, and streaming
- **animation-cli**: CLI commands for animation generation and management

### Modified Capabilities

- None - this is a net-new capability area for the animation system

## Impact

- **New Code**: `src/agent/internal/animation/` - Go animation module with Python bridge
- **New Code**: `bin/manim_renderer.py` - Python wrapper for Manim rendering
- **New Dependencies**: Manim, manim-physics, manim-chemistry, FFmpeg, Coqui TTS
- **Database**: New tables for animation_templates, rendering_jobs, narration_voices
- **API**: New endpoints under `/animations`, `/templates`, `/voice`
- **CLI**: New commands in `lib/bac.nu` for animation management
