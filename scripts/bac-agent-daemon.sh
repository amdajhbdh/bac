#!/usr/bin/env bash
# Manage BAC Agent daemon service
# Controls the autonomous runtime (gateway + channels + scheduler)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Find bac-agent binary
find_bac_agent() {
	# Check in project build directory (crates/agent)
	if [ -f "$PROJECT_ROOT/crates/agent/target/release/bac-agent" ]; then
		echo "$PROJECT_ROOT/crates/agent/target/release/bac-agent"
		return 0
	fi

	# Check in project debug directory
	if [ -f "$PROJECT_ROOT/crates/agent/target/debug/bac-agent" ]; then
		echo "$PROJECT_ROOT/crates/agent/target/debug/bac-agent"
		return 0
	fi

	# Check in PATH
	if command -v bac-agent &>/dev/null; then
		echo "bac-agent"
		return 0
	fi

	# Check common install locations
	if [ -f "$HOME/.cargo/bin/bac-agent" ]; then
		echo "$HOME/.cargo/bin/bac-agent"
		return 0
	fi

	echo ""
	return 1
}

BAC_AGENT_BIN=$(find_bac_agent)

if [ -z "$BAC_AGENT_BIN" ]; then
	echo "Error: bac-agent binary not found"
	echo "Please run: $PROJECT_ROOT/scripts/bac-agent-build.sh"
	exit 1
fi

echo "Using bac-agent: $BAC_AGENT_BIN"

# PID file for daemon management
PID_FILE="$PROJECT_ROOT/.tmp/bac-agent-daemon.pid"
LOG_FILE="$PROJECT_ROOT/.tmp/bac-agent-daemon.log"

# Ensure tmp directory exists
mkdir -p "$(dirname "$PID_FILE")"

# Function to start daemon
start_daemon() {
	if [ -f "$PID_FILE" ]; then
		local pid=$(cat "$PID_FILE")
		if kill -0 "$pid" 2>/dev/null; then
			echo "BAC Agent daemon already running (PID: $pid)"
			return 0
		else
			echo "Stale PID file found, removing..."
			rm -f "$PID_FILE"
		fi
	fi

	echo "Starting BAC Agent daemon..."

	# Start daemon in background
	nohup "$BAC_AGENT_BIN" daemon >>"$LOG_FILE" 2>&1 &
	local pid=$!

	echo $pid >"$PID_FILE"
	echo "✅ BAC Agent daemon started (PID: $pid)"
	echo "Log file: $LOG_FILE"

	# Wait a moment and check if it's running
	sleep 2
	if kill -0 "$pid" 2>/dev/null; then
		echo "✅ Daemon is running"
	else
		echo "❌ Daemon failed to start"
		echo "Check log: $LOG_FILE"
		rm -f "$PID_FILE"
		return 1
	fi
}

# Function to stop daemon
stop_daemon() {
	if [ ! -f "$PID_FILE" ]; then
		echo "BAC Agent daemon is not running (no PID file)"
		return 0
	fi

	local pid=$(cat "$PID_FILE")

	if ! kill -0 "$pid" 2>/dev/null; then
		echo "BAC Agent daemon is not running (stale PID file)"
		rm -f "$PID_FILE"
		return 0
	fi

	echo "Stopping BAC Agent daemon (PID: $pid)..."
	kill "$pid"

	# Wait for process to exit
	local timeout=10
	while kill -0 "$pid" 2>/dev/null && [ $timeout -gt 0 ]; do
		sleep 1
		timeout=$((timeout - 1))
	done

	if kill -0 "$pid" 2>/dev/null; then
		echo "Force killing daemon..."
		kill -9 "$pid"
	fi

	rm -f "$PID_FILE"
	echo "✅ BAC Agent daemon stopped"
}

# Function to check status
check_status() {
	if [ ! -f "$PID_FILE" ]; then
		echo "BAC Agent daemon: NOT RUNNING"
		return 1
	fi

	local pid=$(cat "$PID_FILE")

	if kill -0 "$pid" 2>/dev/null; then
		echo "BAC Agent daemon: RUNNING (PID: $pid)"

		# Try to get status from bac-agent
		if "$BAC_AGENT_BIN" status &>/dev/null; then
			echo "Status check: OK"
		else
			echo "Status check: Unable to connect"
		fi

		return 0
	else
		echo "BAC Agent daemon: NOT RUNNING (stale PID file)"
		rm -f "$PID_FILE"
		return 1
	fi
}

# Function to restart daemon
restart_daemon() {
	echo "Restarting BAC Agent daemon..."
	stop_daemon
	sleep 2
	start_daemon
}

# Function to view logs
view_logs() {
	if [ ! -f "$LOG_FILE" ]; then
		echo "No log file found: $LOG_FILE"
		return 1
	fi

	echo "=== BAC Agent Daemon Logs ==="
	echo "File: $LOG_FILE"
	echo ""
	tail -50 "$LOG_FILE"
}

# Function to follow logs
follow_logs() {
	if [ ! -f "$LOG_FILE" ]; then
		echo "No log file found: $LOG_FILE"
		return 1
	fi

	echo "Following BAC Agent daemon logs..."
	echo "Press Ctrl+C to stop"
	echo "File: $LOG_FILE"
	echo ""
	tail -f "$LOG_FILE"
}

# Main execution
case "${1:-}" in
"start")
	start_daemon
	;;
"stop")
	stop_daemon
	;;
"restart")
	restart_daemon
	;;
"status")
	check_status
	;;
"logs")
	view_logs
	;;
"follow")
	follow_logs
	;;
"help" | "-h" | "--help")
	echo "BAC Agent Daemon Manager"
	echo ""
	echo "Usage:"
	echo "  $0 start    - Start the daemon"
	echo "  $0 stop     - Stop the daemon"
	echo "  $0 restart  - Restart the daemon"
	echo "  $0 status   - Check daemon status"
	echo "  $0 logs     - View recent logs"
	echo "  $0 follow   - Follow logs in real-time"
	echo ""
	echo "The daemon runs the full BAC Agent runtime:"
	echo "  - Gateway server (webhooks, websockets)"
	echo "  - All configured channels (Telegram, WhatsApp, etc.)"
	echo "  - Cron scheduler"
	echo "  - Heartbeat monitor"
	;;
*)
	echo "Unknown command: $1"
	echo "Use: $0 help"
	exit 1
	;;
esac
