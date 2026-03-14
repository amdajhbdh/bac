#!/bin/bash
# Test script for BAC Unified API using Podman

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "=== BAC Unified API - Podman Testing ==="
echo ""

# Check if podman is installed
if ! command -v podman &>/dev/null; then
	echo "Error: podman is not installed"
	exit 1
fi

# Check if podman-compose is installed
if ! command -v podman-compose &>/dev/null; then
	echo "Installing podman-compose..."
	pip install podman-compose
fi

echo "1. Starting containers..."
podman-compose -f podman-compose.test.yaml down 2>/dev/null || true
podman-compose -f podman-compose.test.yaml up -d

echo ""
echo "2. Waiting for services to be healthy..."
echo ""

# Wait for postgres
echo "Waiting for PostgreSQL..."
for i in {1..30}; do
	if podman exec postgres pg_isready -U bac -d bac_unified &>/dev/null; then
		echo "PostgreSQL is ready!"
		break
	fi
	sleep 1
done

# Wait for elasticsearch
echo "Waiting for Elasticsearch..."
for i in {1..30}; do
	if curl -s http://localhost:9200 &>/dev/null; then
		echo "Elasticsearch is ready!"
		break
	fi
	sleep 1
done

# Wait for API
echo "Waiting for API..."
for i in {1..30}; do
	if curl -s http://localhost:8080/health &>/dev/null; then
		echo "API is ready!"
		break
	fi
	sleep 1
done

echo ""
echo "3. Running tests..."

# Test health endpoint
echo "Testing /health endpoint..."
curl -s http://localhost:8080/health | jq .

# Test API endpoints
echo ""
echo "Testing /api/v1/subjects..."
curl -s http://localhost:8080/api/v1/subjects | jq .

echo ""
echo "Testing /api/v1/questions..."
curl -s http://localhost:8080/api/v1/questions | jq .

echo ""
echo "Testing /api/v1/leaderboard..."
curl -s http://localhost:8080/api/v1/leaderboard | jq .

echo ""
echo "=== Tests Complete ==="
echo ""
echo "API Documentation: http://localhost:8080/swagger/index.html"
echo ""
echo "To stop containers: podman-compose -f podman-compose.test.yaml down"
