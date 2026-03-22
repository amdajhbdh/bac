#!/bin/bash
# ============================================================
# Cloud Shell OCR Setup Script
# Installs Ollama + vision models for OCR fallback chain
# Usage: curl -fsSL https://.../cloud-shell-ocr.sh | bash
# ============================================================

set -e

LOGFILE="/tmp/ocr-setup.log"
PORT=11434
export OLLAMA_HOST="0.0.0.0:${PORT}"

log() { echo "[$(date '+%H:%M:%S')] $*" | tee -a "$LOGFILE"; }

detect_os() {
	if command -v apt-get &>/dev/null; then
		echo "debian"
	elif command -v pacman &>/dev/null; then
		echo "arch"
	elif command -v brew &>/dev/null; then
		echo "macos"
	else
		echo "unknown"
	fi
}

install_ollama() {
	log "Installing Ollama..."
	if command -v ollama &>/dev/null; then
		log "Ollama already installed: $(ollama --version 2>/dev/null || echo 'version unknown')"
		return 0
	fi

	local os=$(detect_os)
	case "$os" in
	debian)
		curl -fsSL https://ollama.ai/install.sh | sh
		;;
	macos)
		brew install ollama
		;;
	*)
		curl -fsSL https://ollama.ai/install.sh | sh
		;;
	esac

	log "Ollama installed successfully"
}

pull_vision_models() {
	log "Pulling vision models (this may take 5-15 minutes)..."

	local models=(
		"liquidai/lfm2.5-vl-1.6b:q4_0"
		"qwen2.5vl:latest"
		"llava:latest"
	)

	for model in "${models[@]}"; do
		local name=$(echo "$model" | cut -d: -f1)
		if ollama list 2>/dev/null | grep -q "^$name"; then
			log "Model already exists: $name"
		else
			log "Pulling: $model"
			if ollama pull "$model" 2>&1 | tee -a "$LOGFILE"; then
				log "Pulled: $model"
			else
				log "WARNING: Failed to pull $model"
			fi
		fi
	done
}

start_server() {
	log "Starting llama-server on port ${PORT}..."

	# Kill existing server
	pkill -f "llama-server.*${PORT}" 2>/dev/null || true
	sleep 1

	# Start in background
	nohup llama-server \
		-m ~/.ollama/models/manifests/registry.ollama.ai/library/liquidai/lfm2.5-vl-1.6b \
		-port "$PORT" \
		--host 0.0.0.0 \
		-ngl 0 \
		>/tmp/llama-server.log 2>&1 &

	local pid=$!
	echo $pid >/tmp/llama-server.pid
	log "llama-server started (PID: $pid) on 0.0.0.0:${PORT}"

	# Wait for server to start
	local retries=10
	while [ $retries -gt 0 ]; do
		if curl -s http://localhost:${PORT}/health &>/dev/null || curl -s http://localhost:${PORT}/api/tags &>/dev/null; then
			log "Server is ready!"
			return 0
		fi
		sleep 2
		retries=$((retries - 1))
	done

	log "WARNING: Server may not be fully ready. Check /tmp/llama-server.log"
	return 0
}

start_ollama() {
	log "Starting Ollama daemon..."

	pkill -f "ollama serve" 2>/dev/null || true
	sleep 1

	nohup ollama serve --host 0.0.0.0 >/tmp/ollama.log 2>&1 &
	local pid=$!
	echo $pid >/tmp/ollama.pid
	log "Ollama daemon started (PID: $pid)"

	sleep 3
}

test_ocr() {
	log "Testing OCR with a simple prompt..."

	# Create a test image with text using ImageMagick or fallback to echo
	local test_prompt="What text is in this image? Just say 'Hello world' if you can see it."

	if command -v curl &>/dev/null; then
		# Test with a URL if available
		curl -s -X POST http://localhost:${PORT}/api/generate \
			-H "Content-Type: application/json" \
			-d "{\"model\":\"liquidai/lfm2.5-vl-1.6b\",\"prompt\":\"$test_prompt\",\"stream\":false}" |
			head -c 200 || log "Server not responding yet"
	fi
}

show_usage() {
	log ""
	log "=== Cloud Shell OCR Setup Complete ==="
	log ""
	log "Server running on: 0.0.0.0:${PORT}"
	log "Health check: curl http://localhost:${PORT}/health"
	log "List models: curl http://localhost:${PORT}/api/tags"
	log "Logs: tail -f /tmp/llama-server.log"
	log ""
	log "=== Available Vision Models ==="
	log "  - liquidai/lfm2.5-vl-1.6b (best for French/Arabic)"
	log "  - qwen2.5vl (new multimodal)"
	log "  - llava (stable vision)"
	log ""
	log "=== Exposing for External Access ==="
	log "Option 1: Cloud Shell Web Preview"
	log "  Click 'Web preview' button in Cloud Shell, use port ${PORT}"
	log ""
	log "Option 2: Cloudflare Tunnel (run separately):"
	log "  cloudflared tunnel --url http://localhost:${PORT}"
	log ""
	log "=== OCR API Usage ==="
	log "curl -X POST http://localhost:${PORT}/api/generate \\"
	log "  -H 'Content-Type: application/json' \\"
	log "  -d '{\"model\":\"liquidai/lfm2.5-vl-1.6b\",\"images\":[\"<base64_image>\"],\"prompt\":\"Extract all text from this image\"}'"
	log ""
}

main() {
	log "=== Cloud Shell OCR Setup Started ==="
	log "Log file: $LOGFILE"

	install_ollama
	start_ollama
	pull_vision_models
	start_server
	show_usage

	log "=== Setup Complete ==="
}

main "$@"
