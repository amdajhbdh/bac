#!/usr/bin/env nu
# Process all PDFs in parallel using Nushell

mkdir db/pdfs/processed

let pdfs = (ls db/pdfs/*.pdf | get name)
let count = ($pdfs | length)

print $"Processing (($count)) PDFs in parallel..."
print ""

# Use par-each for parallel processing
$pdfs | par-each {|f|
    let filename = (basename $f)
    let output = $"db/pdfs/processed/($filename).json"
    
    # Run curl and capture result
    let result = (curl -s -X POST http://127.0.0.1:3000/pdf -F $'file=@($f)' -o $output --write-out "%{http_code}")
    
    # Check if file was created and has content
    if ($result == "200" and (ls $output | get size.0) > 10) {
        print $"✓ ($filename)"
    } else {
        print $"✗ ($filename)"
    }
}

print ""
print "========================================="
print "All PDFs processed!"
print "Output: db/pdfs/processed/"
print "========================================="

# Show results
let processed = (ls db/pdfs/processed/*.json | length)
print $"Total processed: ($processed) files"
