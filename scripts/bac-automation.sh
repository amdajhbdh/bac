#!/usr/bin/env bash

# BAC Study Automation Script
# Integrates ZeroClaw messaging, Obsidian sync, and Claude Code operations

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Load environment variables
export $(grep -v '^#' "$PROJECT_ROOT/.env" 2>/dev/null | xargs)

echo "== BAC Study Automation =="
echo "Project root: $PROJECT_ROOT"
echo "Timestamp: $(date -Iseconds)"

# Find zeroclaw binary
find_zeroclaw() {
	# Check in project build directory
	if [ -f "$PROJECT_ROOT/src/zeroclaw/target/release/zeroclaw" ]; then
		echo "$PROJECT_ROOT/src/zeroclaw/target/release/zeroclaw"
		return 0
	fi

	# Check in project debug directory
	if [ -f "$PROJECT_ROOT/src/zeroclaw/target/debug/zeroclaw" ]; then
		echo "$PROJECT_ROOT/src/zeroclaw/target/debug/zeroclaw"
		return 0
	fi

	# Check in PATH
	if command -v zeroclaw &>/dev/null; then
		echo "zeroclaw"
		return 0
	fi

	# Check common install locations
	if [ -f "$HOME/.cargo/bin/zeroclaw" ]; then
		echo "$HOME/.cargo/bin/zeroclaw"
		return 0
	fi

	echo ""
	return 1
}

ZEROCRAW_BIN=$(find_zeroclaw)

# Function to send message via ZeroClaw
send_message() {
	local message="$1"

	if [ -n "$ZEROCRAW_BIN" ]; then
		echo "Sending via ZeroClaw: $message"
		# Use zeroclaw agent to send message to all configured channels
		"$ZEROCRAW_BIN" agent -m "$message" 2>/dev/null || {
			echo "Warning: ZeroClaw not available, using fallback"
			send_fallback "$message"
		}
	else
		send_fallback "$message"
	fi
}

# Fallback: direct API calls if zeroclaw not available
send_fallback() {
	local message="$1"

	# Telegram fallback
	if [[ -n "${TELEGRAM_BOT_TOKEN:-}" && -n "${TELEGRAM_DEFAULT_CHAT_ID:-}" ]]; then
		curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
			-d chat_id="$TELEGRAM_DEFAULT_CHAT_ID" \
			-d text="$message" \
			-d parse_mode="Markdown" >/dev/null || true
	fi

	# WhatsApp fallback
	if [[ -n "${WHATSAPP_ACCESS_TOKEN:-}" && -n "${WHATSAPP_PHONE_NUMBER_ID:-}" ]]; then
		curl -s -X POST "https://graph.facebook.com/v18.0/${WHATSAPP_PHONE_NUMBER_ID}/messages" \
			-H "Authorization: Bearer ${WHATSAPP_ACCESS_TOKEN}" \
			-H "Content-Type: application/json" \
			-d '{"messaging_product":"whatsapp","to":"'"${WHATSAPP_DEFAULT_TO:-}"'","type":"text","text":{"body":"'"$message"'"}}' >/dev/null || true
	fi
}

# Main tasks
main() {
	send_message "🚀 BAC Automation started\n• Timestamp: $(date '+%Y-%m-%d %H:%M:%S')\n• System: Initializing study briefings..."

	# TODO: Add actual automation tasks here
	# 1. Generate study briefings
	# 2. Process new notes from resources/notes/
	# 3. Trigger Claude Code analysis
	# 4. Update status dashboards

	echo "Automation tasks complete"
	send_message "✅ BAC Automation completed successfully"

	# Optional: trigger Obsidian sync (separate script)
	# node "$SCRIPT_DIR/sync-obsidian.js"
}

main "$@"
