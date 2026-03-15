#!/usr/bin/env bash
# Build bac-agent from source
# This script compiles the BAC agent binary for use in BAC automation

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
AGENT_DIR="$PROJECT_ROOT/crates/agent"

echo "== Building BAC Agent =="
echo "Directory: $AGENT_DIR"
echo "Timestamp: $(date -Iseconds)"

cd "$AGENT_DIR"

# Check if Rust is installed
if ! command -v cargo &>/dev/null; then
	echo "Error: Rust/Cargo is not installed"
	echo "Install Rust: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
	exit 1
fi

# Build release binary
echo "Building release binary..."
cargo build --release

# Verify build
if [ -f "target/release/bac-agent" ]; then
	echo "✅ Build successful: target/release/bac-agent"
	ls -lh target/release/bac-agent

	# Optional: Add to PATH for this session
	export PATH="$AGENT_DIR/target/release:$PATH"
	echo "Added to PATH: $AGENT_DIR/target/release"
else
	echo "❌ Build failed: binary not found"
	exit 1
fi

echo "BAC Agent build complete!"
