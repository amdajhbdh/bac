# ============================================================================
# VISION/OCR PIPELINE
# Extracts text from images and PDFs
# Uses Tesseract (free) + optional Mathpix
# ============================================================================

# ============================================================================
# MAIN OCR FUNCTION
# ============================================================================

export def ocr [
    file: string
    --math: bool = false
] {
    print $"%(char-esc)[1;36m📷 OCR Processing: ($file)%(char-reset)"
    
    # Check file exists
    if not ($file | path exists) {
        return { success: false, error: "File not found" }
    }
    
    let ext = ($file | path parse | get extension | str downcase)
    
    match $ext {
        "jpg" | "jpeg" | "png" | "bmp" => {
            if $math {
                ocr-math-image $file
            } else {
                ocr-image $file
            }
        }
        "pdf" => {
            ocr-pdf $file
        }
        _ => {
            { success: false, error: $"Unsupported format: ($ext)" }
        }
    }
}

# ============================================================================
# IMAGE OCR (Tesseract)
# ============================================================================

def ocr-image [file: string] {
    print "Running Tesseract OCR..."
    
    let output_base = ($file | path parse | get stem)
    let output_file = $"/tmp/($output_base).txt"
    
    let result = (^tesseract $file $output_base 2>&1)
    
    if ($result == "") {
        # Success - read output
        if ($output_file | path exists) {
            let text = (open $output_file)
            { success: true, text: $text, type: "text" }
        } else {
            { success: false, error: "Output file not created" }
        }
    } else {
        { success: false, error: $'Tesseract failed: ($result)' }
    }
}

# ============================================================================
# MATH IMAGE OCR (Tesseract with math mode)
# ============================================================================

def ocr-math-image [file: string] {
    print "Running Math OCR (Tesseract with math mode)..."
    
    # Try tesseract with math PSM mode
    let output_base = ($file | path parse |get stem)
    let result = (^tesseract $file $output_base -c "tessedit_char_whitelist=0123456789+-*/=()[]{}<>xyzabcfrsincolexpqrt" --psm 6 2>&1)
    
    let output_file = $"/tmp/($output_base).txt"
    
    if ($output_file | path exists) {
        let text = (open $output_file)
        let cleaned = (clean-math-text $text)
        { success: true, text: $cleaned, type: "math" }
    } else {
        # Fallback to regular OCR
        ocr-image $file
    }
}

# ============================================================================
# PDF OCR
# ============================================================================

def ocr-pdf [file: string] {
    print "Converting PDF to text..."
    
    # Try pdftotext first (fastest)
    let output_base = ($file | path parse | get stem)
    let output_file = $"/tmp/($output_base).txt"
    
    let result = (^pdftotext -layout $file $output_file 2>&1)
    
    if ($result == "" and ($output_file | path exists)) {
        let text = (open $output_file)
        { success: true, text: $text, type: "pdf" }
    } else {
        # Fallback: convert each page
        convert-pdf-pages $file
    }
}

# ============================================================================
# CONVERT PDF PAGES
# ============================================================================

def convert-pdf-pages [file: string] {
    # Use pdftoppm to convert pages to images, then OCR
    print "Converting PDF pages to images..."
    
    let output_base = ($file | path parse | get stem)
    let temp_dir = $"/tmp/($output_base)_pages"
    
    ^mkdir -p $temp_dir
    
    # Convert PDF to images
    let convert_result = (^pdftoppm -r 150 -png $file $"($temp_dir)/page" 2>&1)
    
    if ($convert_result != "") {
        return { success: false, error: "PDF conversion failed" }
    }
    
    # OCR each page
    let pages = (ls $temp_dir)
    let all_text = ""
    
    for page in $pages {
        if ($page.name =~ '.png$') {
            let page_text = (ocr-image $page.name)
            if $page_text.success {
                $all_text = $all_text + "\n" + $page_text.text
            }
        }
    }
    
    if ($all_text != "") {
        { success: true, text: $all_text, type: "pdf" }
    } else {
        { success: false, error: "No text extracted" }
    }
}

# ============================================================================
# CLEAN MATH TEXT
# ============================================================================

def clean-math-text [text: string] {
    let cleaned = $text
        | str replace -a "\n" " "
        | str replace -a "  +" " + "
        | str replace -a "  -" " - "
        | str replace -a "\\s+" " "
        | str trim
    
    $cleaned
}

# ============================================================================
# EXTRACT MATH EQUATIONS
# ============================================================================

export def extract-equations [text: string] {
    # Extract potential math equations from text
    let lines = ($text | lines)
    let equations = []
    
    for line in $lines {
        let trimmed = ($line | str trim)
        # Look for lines with = or math symbols
        if ($trimmed =~ '=' or $trimmed =~ '[+\-*/]' | str length > 3) {
            $equations ++= $trimmed
        }
    }
    
    $equations
}

# ============================================================================
# DETECT IF IMAGE HAS MATH
# ============================================================================

export def has-math [file: string] {
    # Simple heuristic: check file name or do quick OCR
    let basename = ($file | path parse | get stem | str downcase)
    
    let math_keywords = ["math", "equation", "积分", "derivee", "integra"]
    
    ($math_keywords | any {|k| $basename =~ $k })
}
