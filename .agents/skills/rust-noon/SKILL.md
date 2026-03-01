# Skill: Rust + Noon Animation

## Purpose

Rust development with Noon animation engine for BAC Unified.

## When to use

- Working on `src/noon/` for animation generation
- Creating visual explanations for math/physics problems
- Building the animation CLI

## Project Structure

```
src/noon/
├── noon/                 # Main library
│   ├── src/
│   │   └── lib.rs
│   └── Cargo.toml
├── examples/             # Example animations
├── assets/               # Static assets
└── Cargo.toml            # Workspace
```

## Commands

```bash
cd src/noon

# Build
cargo build --release

# Run
cargo run --example hello

# Test
cargo test

# Check
cargo check

# Format
cargo fmt

# Lint
cargo clippy -- -D warnings
```

## Noon Animation Concepts

Noon is a Rust-based animation engine for creating educational visualizations.

### Basic Animation

```rust
use noon::prelude::*;

fn main() {
    // Create animation context
    let mut ctx = AnimationContext::new();
    
    // Add shapes
    ctx.add(Circle::new(100.0, 100.0, 50.0));
    
    // Animate
    ctx.animate(Duration::from_secs(2), |t| {
        // Interpolation logic
    });
    
    // Render
    ctx.render("output.mp4");
}
```

## Environment

| Tool | Purpose |
|------|---------|
| Rust 1.70+ | Compiler |
| wgpu | Graphics (GPU) |
| ffmpeg | Video encoding |

## Integration with Go

The Go agent calls Noon for animation generation:

```go
// From src/agent/internal/animation/animation.go
cmd := exec.Command("cargo", "run", "--manifest-path", "src/noon/Cargo.toml", "--example", name)
```

## Anti-Patterns

- ❌ Blocking main thread with animations
- ❌ Not handling GPU errors gracefully
- ❌ Large animations without streaming
