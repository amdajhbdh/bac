#!/usr/bin/env bash
# OCR Pipeline for BAC Study Vault
# Converts all PDFs to searchable markdown

set -e

VAULT_DIR="/home/med/Documents/bac/resources/notes"
OUTPUT_DIR="$VAULT_DIR/03-Resources/OCR"

mkdir -p "$OUTPUT_DIR"

ocr_pdf() {
	local pdf="$1"
	local basename=$(basename "$pdf" .pdf)
	local output="$OUTPUT_DIR/$basename"

	echo "Processing: $pdf"

	# Convert PDF to images
	pdftoppm -r 300 -png "$pdf" "$output" >/dev/null 2>&1

	# OCR each image and combine
	local text_output="$output.md"
	>"$text_output"

	for img in "$output"*.png; do
		if [ -f "$img" ]; then
			tesseract -l ara+fra+eng "$img" >>"$text_output" 2>/dev/null
			echo "" >>"$text_output"
		fi
	done

	# Add front matter
	local frontmatter="---
title: $(basename "$pdf" .pdf)
date: $(date +%Y-%m-%d)
tags: [ocr, $(dirname "$pdf" | xargs basename)]
source: $pdf
---

"
	local temp=$(mktemp)
	echo "$frontmatter" >"$temp"
	cat "$text_output" >>"$temp"
	mv "$temp" "$text_output"

	# Cleanup images
	rm -f "$output"*.png

	echo "  -> $text_output"
}

export -f ocr_pdf

# Find all PDFs and process
find "$VAULT_DIR" -name "*.pdf" -type f | while read pdf; do
	ocr_pdf "$pdf"
done

echo "OCR complete! Output in: $OUTPUT_DIR"
