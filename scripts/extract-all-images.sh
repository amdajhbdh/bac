#!/bin/bash
# Batch image extraction script for BAC vault
# Extracts images from all PDF files in the vault

set -e

VAULT_ROOT="/home/med/Documents/bac"
OUTPUT_DIR="$VAULT_ROOT/resources/notes/07-Assets/extracted-images"
PDF_DIR="$VAULT_ROOT/resources/notes/03-Resources"

echo "Starting batch image extraction..."
echo "Output directory: $OUTPUT_DIR"
echo "PDF directory: $PDF_DIR"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Find all PDF files and extract images
find "$PDF_DIR" -name "*.pdf" -type f | while read -r pdf_file; do
	echo "Processing: $(basename "$pdf_file")"

	# Create subdirectory for this PDF
	pdf_name=$(basename "$pdf_file" .pdf)
	pdf_output_dir="$OUTPUT_DIR/$pdf_name"
	mkdir -p "$pdf_output_dir"

	# Extract images
	cargo run -p bac-ocr --bin extract-images -- "$pdf_file" "$pdf_output_dir" pdf 2>&1 | grep -E "(Extracting|Extracted|Error)"

	echo ""
done

echo "Batch image extraction complete!"
echo "Images saved to: $OUTPUT_DIR"
