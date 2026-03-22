#!/bin/bash
# ============================================================
# Cloudflare Tunnel Setup for Cloud Shell OCR
# Exposes local Ollama server via a public HTTPS URL
# ============================================================

set -e

PORT=11434
TUNNEL_NAME="bac-ocr"

log() { echo "[$(date '+%H:%M:%S')] $*"; }

install_cloudflared() {
	log "Installing cloudflared..."
	if command -v cloudflared &>/dev/null; then
		log "cloudflared already installed: $(cloudflared --version)"
		return 0
	fi

	# Detect OS and install
	if command -v apt-get &>/dev/null; then
		wget -qO /tmp/cloudflared.deb https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
		dpkg -i /tmp/cloudflared.deb
		rm /tmp/cloudflared.deb
	elif command -v pacman &>/dev/null; then
		sudo pacman -S cloudflared --noconfirm
	else
		curl -fsSL https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb -o /tmp/cloudflared.deb
		dpkg -i /tmp/cloudflared.deb
		rm /tmp/cloudflared.deb
	fi

	log "cloudflared installed successfully"
}

start_tunnel() {
	log "Starting Cloudflare Tunnel to localhost:${PORT}..."

	# Kill existing tunnel
	pkill -f "cloudflared tunnel" 2>/dev/null || true
	sleep 1

	# Start tunnel (no authentication needed for quick tunnels)
	nohup cloudflared tunnel --no-autoupdate \
		--url "http://localhost:${PORT}" \
		--logfile /tmp/cloudflared.log \
		--metrics "localhost:53121" \
		>/tmp/cloudflared-tunnel.log 2>&1 &

	local pid=$!
	echo $pid >/tmp/cloudflared.pid
	log "cloudflared tunnel started (PID: $pid)"

	# Wait for tunnel to be ready
	local retries=30
	while [ $retries -gt 0 ]; do
		local url=$(grep -o 'https://[a-z0-9-]*\.trycloudflare\.com' /tmp/cloudflared.log 2>/dev/null | tail -1)
		if [ -n "$url" ]; then
			log "Tunnel ready: $url"
			echo "$url" >/tmp/cloudflared-url.txt
			return 0
		fi
		sleep 2
		retries=$((retries - 1))
	done

	log "WARNING: Tunnel URL not found. Check /tmp/cloudflared.log"
	return 1
}

show_info() {
	log ""
	log "=== Cloudflare Tunnel Setup Complete ==="
	log ""
	local url=$(cat /tmp/cloudflared-url.txt 2>/dev/null || echo "NOT AVAILABLE")
	log "Public URL: $url"
	log ""
	log "=== Usage ==="
	log "Set in scripts/ocr-config.toml:"
	log "  cloud_shell_url = \"$url\""
	log ""
	log "Or export:"
	log "  export CLOUD_SHELL_URL=\"$url\""
	log ""
	log "Logs: tail -f /tmp/cloudflared.log"
	log "PID: $(cat /tmp/cloudflared.pid 2>/dev/null || echo 'not running')"
	log ""
}

main() {
	log "=== Cloudflare Tunnel Setup Started ==="
	install_cloudflared
	start_tunnel
	show_info
	log "=== Done ==="
}

main "$@"
