#!/bin/bash
# BAC Study - Start All Services
# This script starts all the microservices

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Ports
PORT_API=8081
PORT_GEMINI=3001
PORT_VECTOR=3002
PORT_VAULT=3003
PORT_CLOUD=3004
PORT_GRAPH=3005
PORT_RAG=5000
PORT_OCR=8000

echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}              BAC Study - Starting Services${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"

# Function to check if a port is in use
check_port() {
	if command -v lsof &>/dev/null; then
		lsof -i :$1 &>/dev/null
	elif command -v ss &>/dev/null; then
		ss -tuln | grep -q ":$1 "
	fi
}

# Start a service in background
start_service() {
	local name=$1
	local port=$2
	local cmd=$3

	echo -ne "Starting $name on port $port... "

	if check_port $port; then
		echo -e "${YELLOW}already running${NC}"
		return 0
	fi

	# Run in background
	cd "$PROJECT_DIR"
	$cmd >/tmp/bac-$name.log 2>&1 &
	local pid=$!

	# Wait a bit and check if started
	sleep 2
	if kill -0 $pid 2>/dev/null; then
		echo -e "${GREEN}✓ (PID: $pid)${NC}"
		echo $pid >/tmp/bac-$name.pid
	else
		echo -e "${RED}✗ failed${NC}"
		echo -e "${RED}Check log: /tmp/bac-$name.log${NC}"
	fi
}

# Build all services first
echo -e "\n${GREEN}[Building Services]${NC}"
cd "$PROJECT_DIR"

if [ -f "Cargo.toml" ]; then
	cargo build --release 2>/dev/null || echo -e "${YELLOW}Build warnings (continuing)${NC}"
fi

# Start services
echo -e "\n${GREEN}[Starting Services]${NC}"

# API Gateway
start_service "api" $PORT_API "cd src/api && go run main.go"

# Tool services (if built)
# start_service "gemini-tools" $PORT_GEMINI "cargo run -p gemini-tools -- serve --port $PORT_GEMINI"
# start_service "vector-tools" $PORT_VECTOR "cargo run -p vector-tools -- serve --port $PORT_VECTOR"
# start_service "vault-tools" $PORT_VAULT "cargo run -p vault-tools -- serve --port $PORT_VAULT"
# start_service "cloud-tools" $PORT_CLOUD "cargo run -p cloud-tools -- serve --port $PORT_CLOUD"
# start_service "graph-tools" $PORT_GRAPH "cargo run -p graph-tools -- serve --port $PORT_GRAPH"

# RAG Pipeline
start_service "rag" $PORT_RAG "cargo run -p bac-rag -- serve --port $PORT_RAG"

# OCR Server (if Python available)
if command -v python3 &>/dev/null; then
	start_service "ocr" $PORT_OCR "python3 scripts/ocr/local-ocr-server.py"
fi

# Summary
echo -e "\n${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}              Services Status${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "API Gateway:     http://localhost:$PORT_API"
echo -e "RAG Pipeline:    http://localhost:$PORT_RAG"
echo -e "OCR Server:      http://localhost:$PORT_OCR"
echo ""
echo "Logs directory: /tmp/bac-*.log"
echo "Stop all: ./scripts/stop-services.sh"
echo ""

# Save PIDs
echo $! >/tmp/bac-master.pid
