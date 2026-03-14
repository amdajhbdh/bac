#!/bin/bash
# OCR all PDF files in ~/Downloads/bac folders

SOURCE_DIR="$HOME/Downloads"
OUTPUT_DIR="$SOURCE_DIR/bac-ocr"

mkdir -p "$OUTPUT_DIR"

for folder in math-sm pc-sm mauritania mauritania-extra; do
	folder_path="$SOURCE_DIR/$folder"
	if [ -d "$folder_path" ]; then
		output_folder="$OUTPUT_DIR/$folder"
		mkdir -p "$output_folder"

		echo "=== Processing $folder ==="
		for pdf in "$folder_path"/*.pdf; do
			if [ -f "$pdf" ]; then
				filename=$(basename "$pdf" .pdf)
				txt_file="$output_folder/$filename.txt"

				if [ -f "$txt_file" ] && [ -s "$txt_file" ]; then
					echo "Skipping: $filename (already processed)"
					continue
				fi

				echo "OCR: $filename"
				pdftoppm -png -r 300 "$pdf" /tmp/ocr_page >/dev/null 2>&1

				for img in /tmp/ocr_page-*.png; do
					if [ -f "$img" ]; then
						tesseract "$img" stdout >>"$txt_file" 2>/dev/null
						rm "$img"
					fi
				done

				if [ -f "$txt_file" ] && [ -s "$txt_file" ]; then
					echo "Done: $filename"
				else
					echo "Failed: $filename"
				fi
			fi
		done
	fi
done

echo "=== OCR Complete ==="
echo "Output: $OUTPUT_DIR"
find "$OUTPUT_DIR" -name "*.txt" | wc -l
echo "text files created"
