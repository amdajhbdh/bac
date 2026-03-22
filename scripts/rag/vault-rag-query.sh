#!/bin/bash
# Query all 4 BAC Notes RAGs individually and present results
# Usage: bash vault-rag-query.sh "your question"

QUERY="${1:-}"

if [ -z "$QUERY" ]; then
	echo "Usage: $0 \"question\""
	exit 1
fi

echo "========================================"
echo "  BAC Notes — Searching all RAGs"
echo "========================================"
echo "Query: $QUERY"
echo ""

for RAG in biology chemistry math physics; do
	echo "━━━ $RAG ━━━"
	echo "$QUERY" | aichat --rag "$RAG" 2>/dev/null
	echo ""
done
