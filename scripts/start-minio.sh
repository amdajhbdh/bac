#!/usr/bin/env bash
# Start MinIO server for BAC Unified

set -e

# Configuration
CONTAINER_NAME="bac-minio"
MINIO_PORT=9000
MINIO_PORT_API=9001
DATA_DIR="$(pwd)/data/minio"

# Create data directory
mkdir -p "$DATA_DIR"

# Check if container already running
if podman ps --format "{{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
    echo "MinIO already running"
    exit 0
fi

# Check if container exists (stopped)
if podman ps -a --format "{{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
    echo "Starting existing MinIO container..."
    podman start "$CONTAINER_NAME"
else
    echo "Creating new MinIO container..."
    podman run -d \
        --name "$CONTAINER_NAME" \
        -p "$MINIO_PORT:9000" \
        -p "$MINIO_PORT_API:9001" \
        -v "$DATA_DIR:/data" \
        -e MINIO_ROOT_USER="bacunified" \
        -e MINIO_ROOT_PASSWORD="bacpassword123" \
        docker.io/minio/minio server /data --console-address ":9001"
    
    echo "Waiting for MinIO to start..."
    sleep 5
fi

echo "MinIO started!"
echo "  API: http://localhost:$MINIO_PORT"
echo "  Console: http://localhost:$MINIO_PORT_API"
echo "  User: bacunified"
echo "  Password: bacpassword123"

# Create bucket
podman exec "$CONTAINER_NAME" mc alias set local http://localhost:9000 bacunified bacpassword123 2>/dev/null || true
podman exec "$CONTAINER_NAME" mc mb local/bac-unified 2>/dev/null || true
podman exec "$CONTAINER_NAME" mc anonymous set download local/bac-unified 2>/dev/null || true

echo "Bucket 'bac-unified' created"
