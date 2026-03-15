# BAC Unified - Just Commands
# Run with: just <recipe>

# Set default recipe
default:
    @just --list

# ===========================================
# BUILD COMMANDS
# ===========================================

# Build all services
build:
    cargo build --release -p bac-agent
    cargo build --release -p bac-api
    cargo build --release -p bac-ocr
    cargo build --release -p bac-streaming

# Build just the agent
build-agent:
    cargo build --release -p bac-agent

# Build in debug mode
debug:
    cargo build -p bac-agent

# ===========================================
# RUN COMMANDS
# ===========================================

# Start the agent daemon
start:
    ./crates/agent/target/release/bac-agent daemon &

# Stop the agent daemon
stop:
    pkill -f "bac-agent daemon" || true

# Restart the agent
restart: stop start

# Check agent status
status:
    ./crates/agent/target/release/bac-agent status

# ===========================================
# DOCKER COMMANDS
# ===========================================

# Build all Docker images
docker-build:
    docker-compose build

# Start all services
docker-up:
    docker-compose up -d

# Stop all services
docker-down:
    docker-compose down

# View logs
docker-logs:
    docker-compose logs -f

# ===========================================
# CLEANUP COMMANDS
# ===========================================

# Remove old deprecated directories (manual)
clean-old:
    rm -rf src/zeroclaw src/ai-agent src/api src/gateway src/ocr-service src/ocr-pipeline src/ai-ocr

# Clean build artifacts
clean:
    cargo clean

# ===========================================
# TEST COMMANDS
# ===========================================

# Run tests
test:
    cargo test

# Run agent tests
test-agent:
    cargo test -p bac-agent

# ===========================================
# PROVIDER COMMANDS
# ===========================================

# List available providers
providers:
    ./crates/agent/target/release/bac-agent providers

# Test agent with message
test-message MSG:
    ./crates/agent/target/release/bac-agent agent -m "{{MSG}}"

# ===========================================
# SETUP COMMANDS
# ===========================================

# Initial setup
setup: build
    @echo "Run: just start to start the agent"
    @echo "Configure API keys in .env"
