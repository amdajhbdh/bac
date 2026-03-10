#!/bin/bash

ASSETS_DIR="/home/med/Documents/bac/resources/notes/Study-Vault/Assets"
OUTPUT_DIR="/home/med/Documents/bac/resources/notes/Study-Vault/Extracted"

count=0
total=$(find "$ASSETS_DIR" -name "*.pdf" | wc -l)

echo "Extracting $total PDFs..."

for pdf in "$ASSETS_DIR"/*.pdf; do
    ((count++))
    filename=$(basename "$pdf" .pdf)
    
    # Categorize by subject
    if [[ "$filename" =~ (MATHS|arithmétique|intégrale|matrices|nombre complexe|JoussourMaths) ]]; then
        category="Math"
    elif [[ "$filename" =~ (PC|dynamique|projectile|ressort|champ électrique|force de la place|mécanique|énergie) ]]; then
        category="Physics"
    elif [[ "$filename" =~ (chimie|alcools|cinétique) ]]; then
        category="Chemistry"
    elif [[ "$filename" =~ (SN|Brassage|RFD|hormones|nerveux|Joussour SN) ]]; then
        category="Biology"
    elif [[ "$filename" =~ FRANÇAIS ]]; then
        category="French"
    else
        category="General"
    fi
    
    output_file="$OUTPUT_DIR/$category/${filename}.md"
    
    echo "[$count/$total] Processing: $filename -> $category"
    
    # Extract text using pdftotext
    pdftotext -layout "$pdf" - 2>/dev/null | sed 's/\r$//' > "$output_file"
    
    # Add metadata header
    sed -i "1i# $filename\n\n**Source:** $(basename "$pdf")\n**Category:** $category\n**Extracted:** $(date +%Y-%m-%d)\n\n---\n" "$output_file"
done

echo "Extraction complete! Files saved to: $OUTPUT_DIR"
