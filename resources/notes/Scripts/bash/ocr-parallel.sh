#!/usr/bin/env bash
# High-Performance Parallel OCR Pipeline
# Uses all CPU cores for parallel processing

set -e

VAULT_DIR="/home/med/Documents/bac/resources/notes"
OUTPUT_DIR="$VAULT_DIR/03-Resources/OCR"
mkdir -p "$OUTPUT_DIR"

# Find all PDFs
find "$VAULT_DIR" \( -path "*/03-Resources/*.pdf" -o -path "*/07-Assets/PDFs/*.pdf" -o -path "*/mauritanian-bac/**/*.pdf" \) -type f >/tmp/pdfs_to_ocr.txt

TOTAL=$(wc -l </tmp/pdfs_to_ocr.txt)
echo "Found $TOTAL PDFs to OCR"
echo "Using $(nproc) parallel workers"

ocr_single() {
	local pdf="$1"
	local basename=$(basename "$pdf" .pdf)
	local output="$OUTPUT_DIR/$basename"

	# Skip if already processed
	if [ -f "$output.md" ]; then
		echo "SKIP: $basename (already done)"
		return 0
	fi

	echo "OCR: $basename"

	# Convert PDF to images (single page per image)
	pdftoppm -r 200 -png "$pdf" "$output" 2>/dev/null

	# OCR all images and combine
	>"$output.md"
	for img in "$output"-[0-9]*.png; do
		[ -f "$img" ] || continue
		tesseract -l ara+fra+eng "$img" stdout 2>/dev/null >>"$output.md"
		rm -f "$img" # Cleanup immediately
	done

	# Add front matter
	if [ -s "$output.md" ]; then
		local frontmatter="---
title: $basename
date: $(date +%Y-%m-%d)
tags: [ocr, extracted]
source: ${pdf#$VAULT_DIR/}
---

"
		local temp=$(mktemp)
		echo "$frontmatter" >"$temp"
		cat "$output.md" >>"$temp"
		mv "$temp" "$output.md"
		echo "DONE: $basename"
	else
		rm -f "$output.md"
		echo "FAIL: $basename (empty)"
	fi
}

export -f ocr_single
export OUTPUT_DIR VAULT_DIR

# Run in parallel using all cores
cat /tmp/pdfs_to_ocr.txt | xargs -P $(nproc) -I {} bash -c 'ocr_single "$@"' _ {}

echo "=== OCR Complete ==="
ls -la "$OUTPUT_DIR" | wc -l
echo "files in output directory"
