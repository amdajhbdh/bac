---
source: Termux GitHub Wiki, Rust Documentation
library: Termux
package: termux
topic: cross-compilation
fetched: 2026-03-18T21:56:00Z
official_docs: https://github.com/termux/termux-packages/wiki
---

# Cross-Compilation for Android via Termux

## Overview

Cross-compilation allows building software for Android devices from a Linux host system. Termux provides tools and configurations to facilitate this process.

## Rust Cross-Compilation

### Setting Up Rust Targets

To cross-compile Rust applications for Android, add the appropriate target:

```bash
# Add Android targets
rustup target add aarch64-linux-android
rustup target add armv7-linux-androideabi
rustup target add i686-linux-android
rustup target add x86_64-linux-android
```

### Cargo Configuration

Create or modify `~/.cargo/config.toml`:

```toml
[target.aarch64-linux-android]
rustflags = ["-C", "link-arg=--target=aarch64-linux-android"]

[target.armv7-linux-androideabi]
rustflags = ["-C", "link-arg=--target=arm-linux-androideabi"]

[target.i686-linux-android]
rustflags = ["-C", "link-arg=--target=i686-linux-android"]

[target.x86_64-linux-android]
rustflags = ["-C", "link-arg=--target=x86_64-linux-android"]
```

### Building for Android

```bash
# Build for specific target
cargo build --target aarch64-linux-android --release

# Or set default target in .cargo/config.toml
[build]
target = "aarch64-linux-android"
```

## Android NDK Integration

### NDK Setup

Termux build system automatically handles NDK setup. For manual configuration:

1. Download Android NDK from Google
2. Set environment variables:
```bash
export ANDROID_NDK_HOME=/path/to/ndk
export PATH=$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/linux-x86_64/bin:$PATH
```

### Using NDK with Cargo

For projects requiring C dependencies:

```toml
[target.aarch64-linux-android]
linker = "/path/to/ndk/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android24-clang"
rustflags = ["-C", "link-arg=--target=aarch64-linux-android"]
```

## Termux Package Cross-Compilation

### Build Script Configuration

In your package's `build.sh`:

```bash
TERMUX_PKG_HOMEPAGE="https://example.com"
TERMUX_PKG_DESCRIPTION="Cross-compiled Rust package"
TERMUX_PKG_LICENSE="MIT"
TERMUX_PKG_MAINTAINER="@termux"
TERMUX_PKG_VERSION="1.0.0"
TERMUX_PKG_SRCURL="https://github.com/example/repo/archive/v${TERMUX_PKG_VERSION}.tar.gz"
TERMUX_PKG_SHA256="abc123..."

termux_step_pre_configure() {
    termux_setup_rust
    # Additional Rust configuration if needed
    export CARGO_TARGET_AARCH64_LINUX_ANDROID_RUSTFLAGS="-C link-arg=--target=aarch64-linux-android"
}
```

### Architecture-Specific Builds

Termux supports multiple architectures:
- `aarch64` (ARM 64-bit) - Default
- `arm` (ARM 32-bit)
- `i686` (x86 32-bit)
- `x86_64` (x86 64-bit)

Build for specific architecture:
```bash
./build-package.sh -a arm <package-name>
```

## Cross-Compilation Tools

### Using `cross` Tool

The `cross` tool simplifies cross-compilation:

```bash
# Install cross
cargo install cross

# Build with cross
cross build --target aarch64-linux-android --release
```

### Docker-Based Cross-Compilation

Termux's Docker environment includes all necessary toolchains:

```bash
# Start Docker environment
./scripts/run-docker.sh

# Inside container, build package
./build-package.sh <package-name>
```

## Common Issues and Solutions

### Linker Errors

If encountering linker errors:
1. Ensure NDK is properly configured
2. Check `rustflags` in `.cargo/config.toml`
3. Verify target triple matches NDK version

### Missing Dependencies

For packages with C dependencies:
1. Ensure NDK sysroot is accessible
2. Set `PKG_CONFIG_PATH` appropriately
3. Use `termux_setup_rust` in build script

### Build Time Issues

- **Slow builds**: Use `TERMUX_PKG_MAKE_PROCESSES` to limit parallel jobs
- **Memory issues**: Reduce parallel compilation with `cargo build -j 1`

## Best Practices for Cross-Compilation

1. **Use Termux Docker environment** for consistent builds
2. **Test on actual devices** when possible
3. **Pin Rust version** in `rust-toolchain.toml`
4. **Use `cross` tool** for complex cross-compilation scenarios
5. **Handle C dependencies** carefully with NDK paths
6. **Consider API level** compatibility (Android 7+ for Termux)

## Example: Building OCR System for Android

For your BAC exam preparation OCR system:

```bash
# In build.sh
TERMUX_PKG_HOMEPAGE="https://github.com/your-org/bac-ocr"
TERMUX_PKG_DESCRIPTION="OCR system for BAC exam preparation"
TERMUX_PKG_LICENSE="GPL-3.0"
TERMUX_PKG_MAINTAINER="@termux"
TERMUX_PKG_VERSION="1.0.0"
TERMUX_PKG_SRCURL="https://github.com/your-org/bac-ocr/archive/v${TERMUX_PKG_VERSION}.tar.gz"
TERMUX_PKG_SHA256="..."

termux_step_pre_configure() {
    termux_setup_rust
    # Ensure Tesseract dependencies are available
    TERMUX_PKG_DEPENDS="tesseract, leptonica"
}
```

Build command:
```bash
./build-package.sh -I bac-ocr
```
