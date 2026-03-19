---
source: Termux GitHub Wiki, Rust Documentation
library: Termux
package: termux
topic: mobile-development
fetched: 2026-03-18T21:56:00Z
official_docs: https://github.com/termux/termux-packages/wiki
---

# Mobile Development Workflows with Rust and Termux

## Overview

Developing mobile applications with Rust using Termux provides a powerful combination of performance and portability. This guide covers best practices for building, testing, and deploying Rust applications on Android via Termux.

## Development Environment Setup

### Prerequisites

1. **Termux App** (v0.118.0 or higher)
   - Install from F-Droid or GitHub releases
   - Avoid Play Store versions (deprecated)

2. **Essential Packages**
```bash
pkg update
pkg install rust cargo git make clang
pkg install tesseract leptonica  # For OCR applications
```

3. **Storage Permissions**
   - Grant storage access to Termux
   - Use `termux-setup-storage` command

### Project Structure

Organize your Rust project for mobile deployment:

```
bac-ocr/
├── Cargo.toml
├── src/
│   ├── main.rs
│   └── lib.rs
├── assets/          # Static resources
├── tests/           # Integration tests
└── termux/
    └── build.sh     # Termux package build script
```

## Building Rust Applications for Termux

### Cargo Configuration

Create `.cargo/config.toml` in your project:

```toml
[build]
target = "aarch64-linux-android"

[target.aarch64-linux-android]
linker = "aarch64-linux-android24-clang"
rustflags = ["-C", "link-arg=--target=aarch64-linux-android"]

[termux]
# Termux-specific configuration
prefix = "/data/data/com.termux/files/usr"
```

### Build Commands

```bash
# Debug build
cargo build --target aarch64-linux-android

# Release build (optimized)
cargo build --target aarch64-linux-android --release

# With specific features
cargo build --target aarch64-linux-android --features "tesseract"
```

### Cross-Compilation from Linux

If developing on Linux host:

```bash
# Add Android targets
rustup target add aarch64-linux-android

# Set up NDK (if needed)
export ANDROID_NDK_HOME=/path/to/ndk

# Build
cargo build --target aarch64-linux-android --release
```

## Package Management in Termux

### Installing Dependencies

```bash
# System packages
pkg install tesseract
pkg install leptonica
pkg install imagemagick  # For image processing

# Rust crates (add to Cargo.toml)
[dependencies]
tesseract = "0.15"
image = "0.24"
```

### Managing Termux Packages

Create a Termux package for your application:

```bash
# In termux/build.sh
TERMUX_PKG_HOMEPAGE="https://github.com/your-org/bac-ocr"
TERMUX_PKG_DESCRIPTION="OCR system for BAC exam preparation"
TERMUX_PKG_LICENSE="GPL-3.0"
TERMUX_PKG_MAINTAINER="@termux"
TERMUX_PKG_VERSION="1.0.0"
TERMUX_PKG_SRCURL="https://github.com/your-org/bac-ocr/archive/v${TERMUX_PKG_VERSION}.tar.gz"
TERMUX_PKG_SHA256="..."
TERMUX_PKG_DEPENDS="tesseract, leptonica, libjpeg-turbo, libpng"

termux_step_pre_configure() {
    termux_setup_rust
}
```

### Building and Installing

```bash
# Build package
./build-package.sh -I bac-ocr

# Install on device
dpkg -i output/aarch64/bac-ocr_1.0.0_aarch64.deb
```

## Testing on Android Devices

### Unit Tests

```bash
# Run tests in Termux
cargo test --target aarch64-linux-android

# With specific features
cargo test --target aarch64-linux-android --features "testing"
```

### Integration Testing

Create test scripts in `tests/` directory:

```rust
// tests/integration_test.rs
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_ocr_processing() {
        // Test OCR functionality
        let result = process_image("test.jpg");
        assert!(result.is_ok());
    }
}
```

### Performance Testing

```bash
# Profile with cargo-flamegraph (if available)
cargo flamegraph --target aarch64-linux-android --release

# Benchmark with criterion
cargo bench --target aarch64-linux-android
```

## Deployment Strategies

### Direct Installation

1. Build `.deb` package
2. Transfer to device via SSH or file sharing
3. Install with `dpkg -i package.deb`

### Repository Distribution

1. Set up package repository
2. Add repository to Termux sources
3. Users can install with `pkg install`

### App Bundles

For complex applications:
- Create multiple packages (main, plugins, data)
- Use subpackages for optional components
- Manage dependencies carefully

## Best Practices for Mobile Development

### Performance Optimization

1. **Binary Size**
   - Use `cargo build --release`
   - Strip symbols: `strip target/aarch64-linux-android/release/binary`
   - Enable LTO in Cargo.toml:

```toml
[profile.release]
lto = true
codegen-units = 1
panic = 'abort'
```

2. **Memory Usage**
   - Use efficient data structures
   - Implement proper error handling
   - Avoid unnecessary allocations

3. **Battery Efficiency**
   - Minimize CPU usage
   - Use background processing wisely
   - Implement proper wake lock management

### User Experience

1. **CLI Applications**
   - Provide clear help messages
   - Use color output sparingly
   - Implement progress indicators

2. **GUI Applications** (via X11)
   - Use GTK or other mobile-friendly UI frameworks
   - Test touch interactions
   - Optimize for different screen sizes

### Security Considerations

1. **File Permissions**
   - Respect Android sandbox
   - Use appropriate storage locations
   - Avoid world-writable files

2. **Data Protection**
   - Encrypt sensitive data
   - Use secure storage locations
   - Implement proper authentication

## OCR-Specific Considerations

### Tesseract Integration

For your BAC exam preparation OCR system:

```toml
[dependencies]
tesseract = "0.15"
image = "0.24"
anyhow = "1.0"
```

```rust
use tesseract::Tesseract;

fn process_image(path: &str) -> Result<String, anyhow::Error> {
    let mut t = Tesseract::new(None, Some("eng"))?;
    t.set_image(path)?;
    t.recognize()?;
    Ok(t.get_text()?)
}
```

### Image Processing

```rust
use image::io::Reader as ImageReader;

fn preprocess_image(path: &str) -> Result<(), anyhow::Error> {
    let img = ImageReader::open(path)?.decode()?;
    // Apply preprocessing (grayscale, thresholding, etc.)
    Ok(())
}
```

### Language Support

For French/Arabic text (Mauritania context):

```rust
// Initialize with multiple languages
let mut t = Tesseract::new(None, Some("fra+ara+eng"))?;
```

## Debugging and Troubleshooting

### Common Issues

1. **Linker Errors**
   - Check NDK configuration
   - Verify target triple
   - Ensure all dependencies are available

2. **Runtime Errors**
   - Use `RUST_BACKTRACE=1` for detailed traces
   - Check file permissions
   - Verify resource availability

3. **Performance Issues**
   - Profile with `perf` or `flamegraph`
   - Check memory usage with `top`
   - Optimize hot paths

### Debug Tools

```bash
# Install debugging tools
pkg install gdb
pkg install valgrind  # If available

# Run with debugger
gdb ./target/aarch64-linux-android/release/bac-ocr
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Build Termux Package

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Rust
        uses: actions-rs/toolchain@v1
        with:
          toolchain: stable
          target: aarch64-linux-android
      - name: Build
        run: cargo build --target aarch64-linux-android --release
```

### Automated Testing

```bash
# In CI pipeline
cargo test --target aarch64-linux-android
cargo clippy --target aarch64-linux-android
cargo fmt -- --check
```

## Resources and References

### Official Documentation

- Termux Wiki: https://github.com/termux/termux-packages/wiki
- Rust Documentation: https://doc.rust-lang.org/
- Android NDK: https://developer.android.com/ndk

### Community Resources

- Termux GitHub Discussions
- Rust Users Forum
- Android Developer Community

### Tools and Utilities

- `cargo-ndk`: Simplifies NDK integration
- `cross`: Cross-compilation made easy
- `termux-setup-storage`: Storage permission setup

## Example: Complete BAC OCR Workflow

```bash
# 1. Setup environment
pkg update
pkg install rust cargo git tesseract leptonica

# 2. Clone repository
git clone https://github.com/your-org/bac-ocr
cd bac-ocr

# 3. Build application
cargo build --target aarch64-linux-android --release

# 4. Test with sample images
./target/aarch64-linux-android/release/bac-ocr process exams/

# 5. Create Termux package
cd termux
./build-package.sh -I bac-ocr

# 6. Install on device
dpkg -i ../output/aarch64/bac-ocr_*.deb

# 7. Run application
bac-ocr --help
```

This workflow provides a complete foundation for developing, testing, and deploying Rust applications on Android via Termux, specifically tailored for your BAC exam preparation OCR system.
