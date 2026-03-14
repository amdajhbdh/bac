# ============================================================================
# WOLFRAM ALPHA SOLVER
# Uses Wolfram Alpha for step-by-step math solutions
# ============================================================================

use ../orchestrator.nu make-failure make-success

# ============================================================================
# MAIN WOLFRAM SOLVER
# ============================================================================

export def solve-wolfram [problem: string] {
    # Check if Wolfram can handle this problem
    if not (can-wolfram-solve $problem) {
        return (make-failure "Not suitable for Wolfram Alpha")
    }
    
    print "  Trying Wolfram Alpha..."
    
    # Try Wolfram Alpha API
    let result = (try-wolfram-api $problem)
    
    if ($result.success) {
        let steps = (parse-wolfram-steps $result.response)
        make-success $result.response "wolfram" $steps
    } else {
        # Try free web version
        let web_result = (try-wolfram-web $problem)
        if ($web_result.success) {
            make-success $web_result.response "wolfram-web" $web_result.steps
        } else {
            make-failure "Wolfram Alpha unavailable"
        }
    }
}

# ============================================================================
# CHECK IF WOLFRAM CAN SOLVE
# ============================================================================

def can-wolfram-solve [problem: string] {
    let p = ($problem | str downcase)
    
    # Wolfram is good at these
    let wolfram_patterns = [
        "integra" "derivee" "limite"
        "equation" "resoudre" "solve"
        "matrice" "determinant"
        "complexe" "polynome"
        "limite" "serie"
        "differentielle" "differential"
    ]
    
    ($wolfram_patterns | any {|pat| $p =~ $pat })
}

# ============================================================================
# TRY WOLFRAM API
# ============================================================================

def try-wolfram-api [problem: string] {
    let appid = (get-wolfram-appid)
    
    if ($appid | is-empty) {
        return { success: false, error: "No Wolfram AppID" }
    }
    
    let encoded_problem = (url-encode $problem)
    let url = $"http://api.wolframalpha.com/v2/query?appid=($appid)&input=($encoded_problem)&format=plaintext&podstate=Step-by-step+solution"
    
    let response = (^curl -s $url 2>&1)
    
    if ($response == "") {
        return { success: false, error: "Empty response" }
    }
    
    try {
        # Parse XML response
        let result = (parse-wolfram-xml $response)
        
        if ($result != "") {
            { success: true, response: $result }
        } else {
            { success: false, error: "No solution found" }
        }
    } catch {
        { success: false, error: "Parse error" }
    }
}

# ============================================================================
# TRY WOLFRAM WEB
# ============================================================================

def try-wolfram-web [problem: string] {
    # Would use Playwright to scrape wolframalpha.com
    # For now, return failure
    { success: false, error: "Requires Playwright" }
}

# ============================================================================
# PARSE WOLFRAM XML RESPONSE
# ============================================================================

def parse-wolfram-xml [xml: string] {
    # Extract plain text from Wolfram XML
    # Look for <plaintext> tags
    
    let lines = ($xml | lines)
    let results = []
    let in_result = false
    
    for line in $lines {
        if ($line =~ '<plaintext>') {
            $in_result = true
        }
        if $in_result {
            let content = ($line | str replace -a '<plaintext>' '' | str replace -a '</plaintext>' '' | str trim)
            if ($content != "") {
                $results ++= $content
            }
        }
        if ($line =~ '</plaintext>') {
            $in_result = false
        }
    }
    
    ($results | str join "\n")
}

# ============================================================================
# URL ENCODE
# ============================================================================

def url-encode [text: string] {
    # Simple URL encoding
    let encoded = $text
        | str replace -a " " "%20"
        | str replace -a "+" "%2B"
        | str replace -a "=" "%3D"
        | str replace -a "(" "%28"
        | str replace -a ")" "%29"
    
    $encoded
}

# ============================================================================
# GET WOLFRAM APPID
# ============================================================================

def get-wolfram-appid [] {
    let value = (env | get "WOLFRAM_APPID")
    if ($value | describe) == "nothing" {
        ""
    } else {
        $value
    }
}
