# cd /home/med/Documents/bac
# mkdir -p db/pdfs/processed
# echo "Processing all PDFs in db/pdfs/..."
# echo ""
# count=0
# success=0
# failed=0
# for f in db/pdfs/*.pdf; do
#   filename=$(basename "$f")
#   count=$((count + 1))
#
#   echo "[$count] Processing: $filename"
#
#   result=$(curl -s -X POST http://127.0.0.1:3000/pdf -F "file=@$f" > "db/pdfs/processed/$filename.json" 2>&1)
#
#   if [ -f "db/pdfs/processed/$filename.json" ]; then
#     # Check if it's valid JSON
#     if jq -e '.' "db/pdfs/processed/$filename.json" > /dev/null 2>&1; then
#       success=$((success + 1))
#       echo "  ✓ Done"
#     else
#       failed=$((failed + 1))
#       echo "  ✗ Failed (invalid response)"
#     fi
#   else
#     failed=$((failed + 1))
#     echo "  ✗ Failed"
#   fi
# done
# echo ""
# echo "========================================="
# echo "Processed: $count PDFs"
# echo "Success: $success"
# echo "Failed: $failed"
# echo "Output: db/pdfs/processed/"
# echo "========================================="
# Impl in nu with better nu features and err handling and data mngement
cd /home/med/Documents/bac
let pdf_dir = "db/pdfs";
let output_dir = "db/pdfs/processed";
  # not rust but nu script
echo "Processing all PDFs in $pdf_dir..."
let count = 0;
let success = 0;
let failed = 0;
ls $"($pdf_dir)/*.pdf" | each while { |f|
  let filename = ($f | path basename);
  let count = $count + 1;
  echo "[$count] Processing: $filename";
  
  # let us use http instead of curl for better error handling and job for spawning
  # we will also use a job to process multiple PDFs in parallel
  let result = http post "http://localhost:3000/pdf" -f "file=@$f" 2>&1 | save $"($output_dir)/$filename.json";
  if ($result | path exists) {
    # Check if it's valid JSON
    if $result | from json > /dev/null 2>&1 {
      let success = $success + 1;
      echo "  ✓ Done";
    } else {
      let failed = $failed + 1;
      echo "  ✗ Failed (invalid response)";
    }
  } else {
    let failed = $failed + 1;
    echo "  ✗ Failed";
  }

};
