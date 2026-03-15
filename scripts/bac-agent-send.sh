#!/usr/bin/env bash
# Send messages via BAC Agent
# Uses bac-agent agent -m "message" to send messages through configured channels

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Find bac-agent binary
find_bac_agent() {
	# Check in project build directory
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

# Check if bac-agent is configured
if [ ! -f "$HOME/.bac-agent/config.toml" ]; then
	echo "Warning: BAC Agent not configured yet"
	echo "Run: $BAC_AGENT_BIN onboard"
	echo "Or set up channels first: $BAC_AGENT_BIN channel add telegram"
fi

# Function to send message
send_message() {
	local message="$1"
	local channel="${2:-}" # Optional channel filter

	echo "Sending message via BAC Agent..."
	echo "Message: $message"

	if [ -n "$channel" ]; then
		echo "Channel: $channel"
	fi

	# Send message using bac-agent agent
	"$BAC_AGENT_BIN" agent -m "$message"
}

# Function to send to specific channel
send_to_channel() {
	local channel="$1"
	local message="$2"

	echo "Sending to $channel: $message"
	send_message "$message"
}

# Main execution
case "${1:-}" in
"telegram")
	send_to_channel "telegram" "${2:-}"
	;;
"whatsapp")
	send_to_channel "whatsapp" "${2:-}"
	;;
"message" | "" | "-m")
	if [ -z "${2:-}" ]; then
		echo "Usage: $0 message \"Your message here\""
		echo "       $0 telegram \"Your message here\""
		echo "       $0 whatsapp \"Your message here\""
		exit 1
	fi
	send_message "$2"
	;;
"help" | "-h" | "--help")
	echo "BAC Agent Message Sender"
	echo ""
	echo "Usage:"
	echo "  $0 message \"Your message\"     # Send to all configured channels"
	echo "  $0 telegram \"Your message\"    # Send via Telegram"
	echo "  $0 whatsapp \"Your message\"    # Send via WhatsApp"
	echo ""
	echo "Examples:"
	echo "  $0 message \"🚀 BAC Automation started\""
	echo "  $0 telegram \"Study reminder: Time for math!\""
	;;
*)
	# Default: treat as message
	send_message "$1"
	;;
esac
