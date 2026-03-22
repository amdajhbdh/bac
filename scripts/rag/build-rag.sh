#!/bin/bash
# Build a single subject RAG with Ollama (memory-safe)
set -e

SUBJECT="$1"
RAG_FILE="$HOME/.config/aichat/rags/$SUBJECT.yaml"
LOG="/tmp/aichat-${SUBJECT}.log"

if [ -z "$SUBJECT" ]; then
	echo "Usage: $0 <subject>  (biology|math|physics|chemistry)"
	exit 1
fi

if [ ! -f "$RAG_FILE" ]; then
	echo "RAG config not found: $RAG_FILE"
	exit 1
fi

echo "Building $SUBJECT RAG..."
echo "Log: $LOG"

# Run with nohup and capture PID
nohup aichat --rebuild-rag --rag "$SUBJECT" >"$LOG" 2>&1 &
PID=$!

echo "PID: $PID"

# Monitor loop
while kill -0 $PID 2>/dev/null; do
	LAST=$(tail -3 "$LOG" 2>/dev/null | tail -1)
	MEM=$(ps -p $PID -o rss= 2>/dev/null | awk '{printf "%.0f MB", $1/1024}' || echo "?")
	echo "[$(date '+%H:%M:%S')] PID=$PID MEM=$MEM | $LAST"
	sleep 30
done

# Check result
if grep -q "✓ Saved" "$LOG"; then
	echo "✅ $SUBJECT RAG built successfully!"
	echo "Size: $(du -h "$RAG_FILE" | cut -f1)"
else
	echo "❌ $SUBJECT RAG build failed."
	echo "Last 10 lines of log:"
	tail -10 "$LOG"
fi
