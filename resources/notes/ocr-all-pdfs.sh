#!/bin/bash
# OCR all PDFs in 03-Resources/

RESOURCES="/home/med/Documents/bac/resources/notes/Study-Vault/03-Resources"
OUTPUT="/home/med/Documents/bac/resources/notes/Study-Vault/05-Extracted/OCR"

mkdir -p "$OUTPUT"

cd "$RESOURCES"

echo "🔍 Starting OCR for all PDFs..."
echo ""

for pdf in *.pdf; do
    echo "📄 Processing: $pdf"
    
    # Extract filename without extension
    name="${pdf%.pdf}"
    outfile="$OUTPUT/${name}.txt"
    
    # Try pdftotext first (faster)
    if pdftotext "$pdf" "$outfile" 2>/dev/null; then
        size=$(wc -c < "$outfile")
        if [ "$size" -gt 1000 ]; then
            echo "   ✅ Extracted with pdftotext ($size bytes)"
            continue
        fi
    fi
    
    # If pdftotext failed or extracted too little, use OCR
    echo "   🔄 Running OCR (this may take a while)..."
    ocrmypdf --force-ocr --language ara+eng+fra --output-type txt "$pdf" /dev/null --sidecar "$outfile" 2>/dev/null
    
    if [ -f "$outfile" ]; then
        size=$(wc -c < "$outfile")
        echo "   ✅ OCR complete ($size bytes)"
    else
        echo "   ❌ OCR failed"
    fi
    
    echo ""
done

echo "✅ OCR complete! Output in: $OUTPUT"
echo ""
echo "📊 Results:"
ls -lh "$OUTPUT" | tail -n +2 | awk '{print "   " $9 " - " $5}'

