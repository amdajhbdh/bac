#!/bin/bash
# BAC Agent setup script for BAC Unified platform

set -euo pipefail

echo "Setting up BAC Agent for BAC Unified platform..."

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
AGENT_DIR="$PROJECT_ROOT/crates/agent"

# Check if BAC Agent is already built
if [ ! -f "$AGENT_DIR/target/release/bac-agent" ]; then
	echo "Building BAC Agent..."
	cd "$AGENT_DIR"
	cargo build --release
	cd "$PROJECT_ROOT"
fi

# Create .env file from example if it doesn't exist
if [ ! -f "$AGENT_DIR/.env" ]; then
	if [ -f "$AGENT_DIR/.env.example" ]; then
		cp "$AGENT_DIR/.env.example" "$AGENT_DIR/.env"
		echo "Created .env from .env.example"
		echo "Please edit $AGENT_DIR/.env to add your API keys"
	else
		echo "Warning: .env.example not found"
	fi
fi

# Start BAC Agent daemon in background
echo "Starting BAC Agent daemon..."
cd "$AGENT_DIR"
./target/release/bac-agent daemon &
BAC_AGENT_PID=$!
echo $BAC_AGENT_PID >"$PROJECT_ROOT/.tmp/bac-agent.pid"
cd "$PROJECT_ROOT"

# Wait a moment for startup
sleep 3

# Check if BAC Agent is responding
echo "Checking BAC Agent status..."
if curl -s http://127.0.0.1:42617/health >/dev/null 2>&1; then
	echo "✓ BAC Agent is running and healthy"
else
	echo "✓ BAC Agent started (health check may need more time)"
fi

# Set up some default cron jobs (optional)
echo "Setting up default cron jobs..."
curl -X POST http://127.0.0.1:42617/bac-agent/cron \
	-H "Content-Type: application/json" \
	-d '{"expression":"0 9 * * *","command":"echo \"Good morning! This is your daily briefing from BAC.\""}' 2>/dev/null || echo "Note: Cron setup requires running daemon"

echo ""
echo "=========================================="
echo "BAC Agent setup complete!"
echo "=========================================="
echo "PID: $BAC_AGENT_PID (stored in .tmp/bac-agent.pid)"
echo ""
echo "To manage the agent:"
echo "  $PROJECT_ROOT/scripts/bac-agent-daemon.sh start"
echo "  $PROJECT_ROOT/scripts/bac-agent-daemon.sh stop"
echo "  $PROJECT_ROOT/scripts/bac-agent-daemon.sh status"
echo ""
echo "To send messages:"
echo "  $PROJECT_ROOT/scripts/bac-agent-send.sh message \"Hello!\""
