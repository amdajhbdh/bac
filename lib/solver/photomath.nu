# ============================================================================
# PHOTOMATH WEB SCRAPER
# Uses free web access to get math solutions
# ============================================================================

use ../orchestrator.nu make-failure make-success

# ============================================================================
# MAIN PHOTOMATH SOLVER
# ============================================================================

export def solve-photomath [problem: string] {
    # Check if it's a basic math problem PhotoMath can handle
    if not $problem) {
        return (make (is-basic-math-failure "Not a basic math problem")
    }
    
    print "  Trying PhotoMath..."
    
    # Try the PhotoMath website
    let result = (try-photomath-web $problem)
    
    if ($result.success) {
        let steps = (parse-photomath-steps $result.response)
        make-success $result.response "photomath" $steps
    } else {
        make-failure "PhotoMath unavailable"
    }
}

# ============================================================================
# CHECK IF BASIC MATH
# ============================================================================

def is-basic-math [problem: string] {
    let p = ($problem | str downcase)
    
    # Basic arithmetic and algebra
    let basic_patterns = [
        "x+" "x-" "x*" "x/" 
        "equation" "resoudre" "solve"
        "+" "-" "*" "/" "="
        "factoriser" "developper"
        "integrale" "derivee"
    ]
    
    # Check if contains at least one pattern
    let is_basic = ($basic_patterns | any {|pat| $p =~ $pat })
    
    $is_basic
}

# ============================================================================
# TRY PHOTOMATH WEB
# ============================================================================

def try-photomath-web [problem: string] {
    # PhotoMath doesn't have a public API, so we search for the problem
    # using DuckDuckGo and try to find solutions
    
    let search_term = ($problem | str replace -a " " "+")
    let url = $"https://duckduckgo.com/?q=($search_term)+photomath+solution"
    
    # Return failure - would need Playwright for actual scraping
    { success: false, error: "PhotoMath scraping requires Playwright" }
}

# ============================================================================
# PARSE PHOTOMATH STEPS
# ============================================================================

def parse-photomath-steps [response: string] {
    # Parse the scraped response to extract steps
    let lines = ($response | lines)
    let steps = []
    
    for line in $lines {
        let trimmed = ($line | str trim)
        if ($trimmed != "" and $trimmed | str length > 10) {
            $steps ++= $trimmed
        }
    }
    
    $steps
}

# ============================================================================
# MATH SYMBOL CONVERSION
# ============================================================================

export def normalize-math [problem: string] {
    # Convert French math terms to standard notation
    let normalized = $problem
        | str replace -a "×" "*"
        | str replace -a "÷" "/"
        | str replace -a "√" "sqrt"
        | str replace -a "π" "pi"
        | str replace -a "≤" "<="
        | str replace -a "≥" ">="
        | str replace -a "≠" "!="
        | str replace -a "²" "^2"
        | str replace -a "³" "^3"
    
    $normalized
}

# ============================================================================
# EXTRACT EQUATION
# ============================================================================

export def extract-equation [text: string] {
    # Extract mathematical equation from text
    let equation_pattern = '[^a-zA-ZàâäéèêëïîôùûüçæœÀÂÄÉÈÊËÏÎÔÙÛÜÇÆŒ0-9\+\-\*\/\=\(\)\[\]\{\}\.\,\s\^√π]'
    
    let cleaned = ($text | str replace -a $equation_pattern "")
    
    $cleaned
}
