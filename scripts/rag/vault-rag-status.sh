#!/bin/bash
# Check status of all BAC Notes RAGs
echo "========================================"
echo "  BAC Notes Vault RAG Status"
echo "========================================"
echo ""

total_files=0
total_vectors=0
built=0

for RAG in biology chemistry math physics; do
	RAG_PATH="$HOME/.config/aichat/rags/$RAG.yaml"
	if [ -f "$RAG_PATH" ]; then
		SIZE=$(stat -c%s "$RAG_PATH" 2>/dev/null || stat -f%z "$RAG_PATH" 2>/dev/null)
		FILES=$(python3 -c "
import yaml
with open('$RAG_PATH') as f:
    rag = yaml.safe_load(f)
print(len(rag.get('files', {})))
" 2>/dev/null || echo "?")
		VECTS=$(python3 -c "
import yaml
with open('$RAG_PATH') as f:
    rag = yaml.safe_load(f)
print(len(rag.get('vectors', {})))
" 2>/dev/null || echo "?")
		if [ "$VECTS" != "0" ] && [ "$VECTS" != "?" ]; then
			echo "  [✓] $RAG — $FILES files, $VECTS vectors ($((SIZE / 1024))KB)"
			built=$((built + 1))
			total_files=$((total_files + FILES))
			total_vectors=$((total_vectors + VECTS))
		else
			echo "  [✗] $RAG — empty (needs rebuild)"
		fi
	else
		echo "  [✗] $RAG — not configured"
	fi
done

echo ""
echo "Built: $built/4"
echo "Total: $total_files files, $total_vectors vectors"
echo ""
echo "Query with:"
echo "  aichat --rag <name> \"question\""
echo "  bash scripts/vault-rag-query.sh \"question\"  # queries all 4"
