#!/bin/bash
# BAC Study - Stop All Services

echo "Stopping BAC services..."

# Kill all service PIDs
for service in api rag ocr gemini-tools vector-tools vault-tools cloud-tools graph-tools; do
	if [ -f "/tmp/bac-$service.pid" ]; then
		pid=$(cat /tmp/bac-$service.pid)
		if kill -0 $pid 2>/dev/null; then
			kill $pid 2>/dev/null && echo "Stopped $service (PID: $pid)"
		fi
		rm /tmp/bac-$service.pid
	fi
done

# Also kill any remaining bac-* processes
pkill -f "bac-rag" 2>/dev/null
pkill -f "bac-api" 2>/dev/null
pkill -f "local-ocr-server" 2>/dev/null

echo "All services stopped"
