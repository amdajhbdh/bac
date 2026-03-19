---
source: Termux GitHub Wiki
library: Termux
package: termux
topic: build-systems
fetched: 2026-03-18T21:56:00Z
official_docs: https://github.com/termux/termux-packages/wiki
---

# Termux Build Systems and Package Management

## Overview

Termux provides a comprehensive build environment for creating packages that run on Android devices. The build system is based on Debian's package management but adapted for the Android environment.

## Build Environment Setup

### Docker Container (Recommended)

The official Docker image ensures reproducible builds:

```bash
./scripts/run-docker.sh
```

This creates a container with the name `termux-package-builder` by default, mounting the `termux-packages` repo at `/home/builder/termux-packages`.

### Host OS Setup

For Ubuntu:
```bash
./scripts/setup-ubuntu.sh
```

For Arch Linux:
```bash
./scripts/setup-archlinux.sh
```

### On-Device Builds

Termux supports building packages directly on Android devices:

```bash
./scripts/setup-termux.sh
```

**Note:** Not all packages support on-device builds. Check `$TERMUX_PKG_ON_DEVICE_BUILD_NOT_SUPPORTED` in the package's `build.sh`.

## Package Build Process

### Directory Structure

- `./build-package.sh` - Main build script
- `./packages/` - Main channel packages
- `./root-packages/` - Root-required packages
- `./x11-packages/` - X11/GUI packages
- `./output/` - Built package files (.deb)

### Building Packages

Basic build command:
```bash
./build-package.sh <package-name>
```

Common flags:
- `-I` - Download dependencies instead of building them
- `-f` - Force rebuild even if already built
- `-a <arch>` - Specify architecture (aarch64, arm, i686, x86_64)
- `-o <dir>` - Output directory for built packages

### Package Channels

1. **main** - Standard terminal packages
2. **root** - Packages requiring root access
3. **x11** - X11/GUI packages

## Package Build Script Variables

### Required Variables

- `TERMUX_PKG_HOMEPAGE` - Package homepage URL
- `TERMUX_PKG_DESCRIPTION` - One-line description
- `TERMUX_PKG_LICENSE` - SPDX license identifier
- `TERMUX_PKG_MAINTAINER` - Package maintainer (@termux or individual)
- `TERMUX_PKG_VERSION` - Package version
- `TERMUX_PKG_SRCURL` - Source code URL
- `TERMUX_PKG_SHA256` - SHA-256 checksum of source

### Optional Variables

- `TERMUX_PKG_DEPENDS` - Runtime dependencies
- `TERMUX_PKG_BUILD_DEPENDS` - Build-time dependencies
- `TERMUX_PKG_REVISION` - Package revision number
- `TERMUX_PKG_EXTRA_CONFIGURE_ARGS` - Extra configure arguments
- `TERMUX_PKG_NO_STRIP` - Disable binary stripping
- `TERMUX_PKG_ON_DEVICE_BUILD_NOT_SUPPORTED` - Disable on-device builds

## Build Steps

The build process follows these steps:

1. **Setup Variables** - Initialize environment and paths
2. **Get Dependencies** - Download or build dependencies
3. **Get Source** - Download and extract source code
4. **Patch Package** - Apply any patches
5. **Configure** - Run configure script or CMake/Meson
6. **Make** - Compile source code
7. **Install** - Install to staging directory
8. **Massage** - Strip binaries, fix permissions
9. **Create Package** - Generate .deb file

## Rust Support in Termux

### Setup Rust in build.sh

To build Rust packages, use the `termux_setup_rust` utility function:

```bash
termux_step_pre_configure() {
    termux_setup_rust
}
```

### Rust Configuration

For cross-compilation to Android targets, add to `~/.cargo/config.toml`:

```toml
[target.aarch64-linux-android]
rustflags = "-Clink-arg=--target=aarch64-linux-android"

[target.armv7-linux-androideabi]
rustflags = "-Clink-arg=--target=arm-linux-androideabi"
```

### Building Rust Packages

Example build.sh for a Rust package:

```bash
TERMUX_PKG_HOMEPAGE=https://example.com
TERMUX_PKG_DESCRIPTION="Example Rust package"
TERMUX_PKG_LICENSE="MIT"
TERMUX_PKG_MAINTAINER="@termux"
TERMUX_PKG_VERSION=1.0.0
TERMUX_PKG_SRCURL=https://github.com/example/repo/archive/v${TERMUX_PKG_VERSION}.tar.gz
TERMUX_PKG_SHA256=abc123...

termux_step_pre_configure() {
    termux_setup_rust
}
```

## Package Management

### Installing Packages

```bash
pkg install <package-name>
```

### Updating Packages

```bash
pkg upgrade
```

### Changing Mirrors

If experiencing repository errors:

```bash
termux-change-repo
```

Then select a different mirror for each repository channel.

### Repository Channels

- Main: `https://packages.termux.dev/apt/termux-main`
- Root: `https://packages.termux.dev/apt/termux-root`
- X11: `https://packages.termux.dev/apt/termux-x11`

## Best Practices

1. **Use Docker** for consistent, reproducible builds
2. **Test on device** before submitting packages
3. **Follow coding guidelines** for consistency
4. **Use appropriate licenses** (SPDX identifiers)
5. **Document dependencies** clearly
6. **Consider on-device build limitations**
7. **Use version pinning** for stable builds
