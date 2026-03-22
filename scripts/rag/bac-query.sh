#!/bin/bash
# BAC RAG Query Tool
# Usage: ./bac-query.sh "your question in French"

QUERY="${1:-}"
if [ -z "$QUERY" ]; then
	echo "Usage: ./bac-query.sh \"question\""
	exit 1
fi

curl -s -X POST "https://bac-api.amdajhbdh.workers.dev/rag/query" \
	-H "Content-Type: application/json" \
	-d "{\"query\": \"$QUERY\"}" | jq -r '
    "=== Réponse ===",
    .answer,
    "",
    "=== Sources ===",
    (.sources[:3] | .[] | "\(.subject): \(.source)")
  '
