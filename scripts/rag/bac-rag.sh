#!/bin/bash
# BAC RAG - Query the Worker API from aichat or terminal
# Usage: ./bac-rag.sh "question" [subject]
# Examples:
#   ./bac-rag.sh "explain mitose" biology
#   ./bac-rag.sh "solve x^2-5x+6=0" math

API="https://bac-api.amdajhbdh.workers.dev"
QUERY="${1:-}"
SUBJECT="${2:-}"

if [ -z "$QUERY" ]; then
	echo "Usage: $0 \"question\" [subject]"
	echo "Subjects: biology, chemistry, math, physics"
	exit 1
fi

PAYLOAD="{\"query\": \"$QUERY\"}"
if [ -n "$SUBJECT" ]; then
	PAYLOAD="{\"query\": \"$QUERY\", \"subjects\": [\"$SUBJECT\"]}"
fi

curl -s -X POST "$API/rag/query" \
	-H "Content-Type: application/json" \
	-d "$PAYLOAD" | python3 -c "
import json, sys
d = json.load(sys.stdin)
if 'error' in d:
    print(f\"Error: {d['error']}\")
else:
    print(d.get('answer', 'No answer'))
    if d.get('sources'):
        print(f\"\nSources: {len(d['sources'])} chunks\")
    if d.get('fallback'):
        print(f\"(via {d.get('model', 'fallback')})\")
"
