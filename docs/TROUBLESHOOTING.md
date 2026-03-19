# Troubleshooting Guide

Common issues and solutions for BAC Knowledge System.

## Installation Issues

### Rust Toolchain Not Found

```bash
# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Reload shell
source ~/.cargo/env

# Verify
rustc --version
```

### aichat Installation Fails

```bash
# Try alternative install methods
# Method 1: Homebrew (macOS)
brew install aichat

# Method 2: From source
git clone https://github.com/sigoden/aichat
cd aichat && cargo install --path .
```

### Node.js Dependencies Fail

```bash
# Clear cache and retry
rm -rf node_modules package-lock.json
npm install

# Or use npm ci for clean install
npm ci
```

---

## Build Errors

### Compilation Fails with Missing Dependencies

```bash
# Update dependencies
cargo update

# Clean and rebuild
cargo clean
cargo build --release
```

### Link Errors on macOS

```bash
# Install required libraries
brew install openssl pkg-config

# Set environment
export OPENSSL_LIB_DIR=$(brew --prefix openssl)/lib
export OPENSSL_INCLUDE_DIR=$(brew --prefix openssl)/include
```

### Link Errors on Linux

```bash
# Install development libraries
sudo apt install libssl-dev pkg-config
```

---

## Service Connection Issues

### Services Won't Start

```bash
# Check if ports are in use
lsof -i :3001  # gemini-tools
lsof -i :3002  # vector-tools
lsof -i :3003  # vault-tools

# Kill existing processes
kill $(lsof -t -i :3001)
```

### Database Connection Failed

```bash
# Test connection
psql $NEON_DB_URL -c "SELECT 1"

# Check environment
echo $NEON_DB_URL

# For local PostgreSQL
sudo systemctl start postgresql
```

### pgvector Not Available

```bash
# Enable pgvector extension
psql -c "CREATE EXTENSION IF NOT EXISTS vector;"

# Check version
psql -c "SELECT extversion FROM pg_extension WHERE extname = 'vector';"
```

---

## aichat Issues

### Tool Services Not Responding

```bash
# Verify service health
curl http://localhost:3001/health  # gemini-tools
curl http://localhost:3002/health  # vector-tools

# Check environment variables
echo $GEMINI_TOOLS_URL
echo $VECTOR_TOOLS_URL
```

### RAG Not Finding Files

```bash
# Verify vault path
ls -la $VAULT_PATH

# Check config
grep "vault" config/aichat/config.yaml
```

### Session Not Saving

```bash
# Check permissions
chmod 755 ~/.aichat/sessions

# Verify save path in config
grep "save_session" config/aichat/config.yaml
```

---

## OCR Problems

### Tesseract Not Found

```bash
# Install Tesseract
# macOS
brew install tesseract tesseract-lang

# Ubuntu/Debian
sudo apt install tesseract-ocr tesseract-ocr-ara tesseract-ocr-eng

# Verify
tesseract --version
```

### OCR Returns Empty Text

```bash
# Test with known image
tesseract test.jpg stdout

# Increase image quality
# - Use higher resolution images
# - Ensure good contrast
# - Convert to grayscale if needed
convert input.png -colorspace gray output.png
```

---

## Environment Issues

### .env Not Loading

```bash
# Check file exists
ls -la .env

# Use direnv for auto-loading
echo 'export $(cat .env | xargs)' >> .bashrc
source .bashrc

# Or manually export
set -a && source .env && set +a
```

### Wrong API Keys

```bash
# Verify key format
echo $GEMINI_API_KEY | head -c 10

# Test Gemini connection
curl -H "Authorization: Bearer $GEMINI_API_KEY" \
  https://generativelanguage.googleapis.com/v1/models
```

---

## Performance Issues

### Slow Semantic Search

```bash
# Check query time
# Rebuild HNSW index with better parameters
curl -X POST http://localhost:3002/rebuild \
  -H "Content-Type: application/json" \
  -d '{"m": 32, "ef_construction": 128}'

# Increase ef_search
curl -X POST http://localhost:3002/search \
  -H "Content-Type: application/json" \
  -d '{"query": [...], "ef_search": 100}'
```

### High Memory Usage

```bash
# Limit concurrent requests
# In config/aichat/config.yaml
max_concurrent: 5

# Use streaming for large outputs
stream: true
```

---

## Cloud Shell Issues

### SSH Connection Fails

```bash
# Test SSH manually
ssh -i ~/.ssh/gcp-key cloud-shell@instance

# Regenerate keys if needed
ssh-keygen -t ed25519 -f ~/.ssh/gcp-key -N ""
```

### Cloud Upload Fails

```bash
# Check disk quota
gcloud compute ssh instance -- df -h

# Clear old files
gcloud compute ssh instance -- rm -rf /tmp/bac-*
```

---

## Getting Help

### Debug Mode

```bash
# Enable debug logging
RUST_LOG=debug cargo run -p gemini-tools

# Full trace
RUST_LOG=trace cargo run -p vector-tools
```

### Check Service Status

```bash
# All services
curl http://localhost:8080/health

# Individual
curl http://localhost:3001/health
curl http://localhost:3002/health
```

### View Logs

```bash
# Rust services use tracing
# Check stdout/stderr

# Or configure logging to file
RUST_LOG=debug ./target/release/gemini-tools 2>&1 | tee logs/gemini.log
```

---

## Reset Everything

If all else fails:

```bash
# Stop all services
pkill -f "cargo run"

# Clean databases
dropdb bac
createdb bac
psql -c "CREATE EXTENSION vector;"

# Fresh build
cargo clean
cargo build --release

# Restart
./scripts/bac-agent-daemon.sh
```
