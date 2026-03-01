#!/bin/bash
# Zen Browser Automation Wrapper for BAC Agent

set -e

PROVIDER="$1"
PROBLEM="$2"
PROFILE_DIR="$HOME/.bac-agent/auth/zen-$PROVIDER"

mkdir -p "$PROFILE_DIR"

case "$PROVIDER" in
deepseek)
	URL="https://chat.deepseek.com"
	;;
grok)
	URL="https://grok.com"
	;;
claude)
	URL="https://claude.ai"
	;;
chatgpt)
	URL="https://chat.openai.com"
	;;
*)
	echo "Unknown provider: $PROVIDER"
	exit 1
	;;
esac

# Open Zen with the URL and profile
zen --new-window "$URL" \
	--profile "$PROFILE_DIR" \
	2>/dev/null &

echo "Opened $PROVIDER at $URL"
echo "Profile: $PROFILE_DIR"

# Wait for browser to start
sleep 2

# Take screenshot for debugging
import -window root "$HOME/bac-$PROVIDER-$(date +%s).png" 2>/dev/null || true

echo "Browser opened. Please interact manually or use screenshot for automation."
