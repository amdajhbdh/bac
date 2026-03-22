#!/bin/bash
# Extract text from all PDFs in 03-Resources folder

set -e

VAULT_DIR="/home/med/Documents/bac/resources/notes"
RESOURCES_DIR="$VAULT_DIR/03-Resources"
EXTRACTED_DIR="$VAULT_DIR/05-Extracted/Resources"

echo "=== PDF Extraction Script ==="
echo "Vault: $VAULT_DIR"
echo "Resources: $RESOURCES_DIR"
echo "Output: $EXTRACTED_DIR"
echo ""

# Create output directory
mkdir -p "$EXTRACTED_DIR"

# Count PDFs
PDF_COUNT=$(find "$RESOURCES_DIR" -name "*.pdf" | wc -l)
echo "Found $PDF_COUNT PDF files in $RESOURCES_DIR"
echo ""

# Extract each PDF
find "$RESOURCES_DIR" -name "*.pdf" | while read pdf_file; do
	filename=$(basename "$pdf_file")
	output_file="$EXTRACTED_DIR/${filename%.*}.md"

	echo "Processing: $filename"

	# Check if already extracted
	if [ -f "$output_file" ]; then
		echo "  ⚠️  Already extracted, skipping"
		continue
	fi

	# Extract text using unpdf
	if cargo run --bin extract-pdf -- "$pdf_file" markdown >"$output_file" 2>/dev/null; then
		lines=$(wc -l <"$output_file")
		chars=$(wc -c <"$output_file")
		echo "  ✅ Extracted: $lines lines, $chars chars"
	else
		echo "  ❌ Failed to extract"
		rm -f "$output_file"
	fi
done

echo ""
echo "=== Extraction Complete ==="
echo "Output directory: $EXTRACTED_DIR"
