#!/bin/bash
# Parallel OCR script
SOURCE_DIR="$HOME/Downloads"
OUTPUT_DIR="$SOURCE_DIR/bac-ocr"

process_pdf() {
	pdf=$1
	folder=$2
	filename=$(basename "$pdf" .pdf)
	output_folder="$OUTPUT_DIR/$folder"
	txt_file="$output_folder/$filename.txt"

	if [ -f "$txt_file" ] && [ -s "$txt_file" ]; then
		echo "SKIP: $filename"
		return
	fi

	rm -f /tmp/ocr_$$.png
	pdftoppm -png -r 150 "$pdf" /tmp/ocr_$$ >/dev/null 2>&1

	for img in /tmp/ocr_$$-*.png; do
		if [ -f "$img" ]; then
			tesseract "$img" stdout -l fra+eng 2>/dev/null >>"$txt_file"
			rm "$img"
		fi
	done

	if [ -f "$txt_file" ] && [ -s "$txt_file" ]; then
		echo "DONE: $filename"
	else
		echo "FAIL: $filename"
	fi
}

export -f process_pdf
export OUTPUT_DIR
export SOURCE_DIR

# Process remaining PDFs in parallel
for folder in pc-sm mauritania mauritania-extra; do
	folder_path="$SOURCE_DIR/$folder"
	if [ -d "$folder_path" ]; then
		mkdir -p "$OUTPUT_DIR/$folder"
		ls "$folder_path"/*.pdf | parallel -j 4 process_pdf {} "$folder"
	fi
done
