#!/bin/bash
# Batch OCR extraction script for BAC vault
# Extracts text from all scanned images in the OCR directory

set -e

VAULT_ROOT="/home/med/Documents/bac"
INPUT_DIR="$VAULT_ROOT/resources/notes/03-Resources/OCR"
OUTPUT_DIR="$VAULT_ROOT/resources/notes/03-Resources/OCR/extracted-text"

echo "Starting batch OCR extraction..."
echo "Input directory: $INPUT_DIR"
echo "Output directory: $OUTPUT_DIR"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Find all image files
echo "Finding image files..."
find "$INPUT_DIR" -type f \( -name "*.png" -o -name "*.jpg" -o -name "*.jpeg" \) >/tmp/ocr_images.txt

total_images=$(wc -l </tmp/ocr_images.txt)
echo "Found $total_images image files"
echo ""

# Process each image
counter=0
while IFS= read -r image_file; do
	counter=$((counter + 1))
	echo "[$counter/$total_images] Processing: $(basename "$image_file")"

	# Create output filename
	image_name=$(basename "$image_file" .png)
	image_name=$(basename "$image_name" .jpg)
	image_name=$(basename "$image_name" .jpeg)
	output_file="$OUTPUT_DIR/${image_name}.txt"

	# Perform OCR
	cargo run -p bac-ocr --bin ocr-extract -- "$image_file" "$output_file" 2>&1 | grep -E "(Performing|Text saved|Error)"

	# Check if output file was created
	if [ -f "$output_file" ]; then
		char_count=$(wc -c <"$output_file")
		echo "  Extracted $char_count characters to $(basename "$output_file")"
	else
		echo "  ❌ Failed to extract text"
	fi

	echo ""
done </tmp/ocr_images.txt

# Clean up
rm /tmp/ocr_images.txt

echo "Batch OCR extraction complete!"
echo "Extracted text saved to: $OUTPUT_DIR"
