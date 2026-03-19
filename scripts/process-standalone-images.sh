#!/bin/bash
# Process standalone images in the BAC vault

set -e

VAULT_ROOT="/home/med/Documents/bac"
INPUT_DIR="$VAULT_ROOT/resources/notes"
OUTPUT_DIR="$VAULT_ROOT/resources/notes/07-Assets/processed-images"

echo "Starting standalone image processing..."
echo "Input directory: $INPUT_DIR"
echo "Output directory: $OUTPUT_DIR"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Find all image files (excluding virtual environment and site-packages)
echo "Finding image files..."
find "$INPUT_DIR" -type f \( -name "*.png" -o -name "*.jpg" -o -name "*.jpeg" -o -name "*.gif" -o -name "*.bmp" -o -name "*.webp" \) ! -path "*/.venv/*" ! -path "*/site-packages/*" >/tmp/image_files.txt

total_images=$(wc -l </tmp/image_files.txt)
echo "Found $total_images image files"
echo ""

# Process each image
counter=0
while IFS= read -r image_file; do
	counter=$((counter + 1))
	echo "[$counter/$total_images] Processing: $(basename "$image_file")"

	# Get the relative path from input directory
	relative_path="${image_file#$INPUT_DIR/}"

	# Create output subdirectory structure
	output_subdir="$OUTPUT_DIR/$(dirname "$relative_path")"
	mkdir -p "$output_subdir"

	# Copy or process the image
	cp "$image_file" "$output_subdir/"

	echo "  Copied to: $output_subdir/"
done </tmp/image_files.txt

# Clean up
rm /tmp/image_files.txt

echo ""
echo "Standalone image processing complete!"
echo "Processed images saved to: $OUTPUT_DIR"
