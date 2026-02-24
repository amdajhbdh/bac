# ============================================================================
# EXTERNAL AI SOLVER - ChatGPT/Grok/Claude via Playwright
# Uses free AI services without API keys
# ============================================================================

use ../orchestrator.nu make-failure make-success

# ============================================================================
# MAIN EXTERNAL SOLVER
# ============================================================================

export def solve-external [problem: string, subject: string] {
    print "Trying external AI solvers..."
    
    # Try Grok first (often free, fast)
    let grok_result = (solve-grok $problem)
    if ($grok_result.success) {
        return $grok_result
    }
    
    # Try ChatGPT (may require login)
    let chatgpt_result = (solve-chatgpt $problem)
    if ($chatgpt_result.success) {
        return $chatgpt_result
    }
    
    # Try Claude (may require login)
    let claude_result = (solve-claude $problem)
    if ($claude_result.success) {
        return $claude_result
    }
    
    make-failure "All external AI solvers failed"
}

# ============================================================================
# GROK SOLVER (xAI - often has free tier)
# ============================================================================

export def solve-grok [problem: string] {
    print "  Trying Grok..."
    
    let prompt = (build-grok-prompt $problem)
    
    # Try using Grok API if available, otherwise use web
    let result = (try-grok-api $prompt)
    
    if ($result.success) {
        let steps = (extract-steps $result.response)
        make-success $result.response "grok" $steps
    } else {
        # Fallback: try web search + scraping
        let web_result = (try-grok-web $problem)
        if ($web_result.success) {
            make-success $web_result.response "grok-web" $web_result.steps
        } else {
            make-failure "Grok unavailable"
        }
    }
}

def build-grok-prompt [problem: string] {
    $"Solve this step by step in French for BAC C exam:

($problem)

Provide clear steps and final answer."
}

def try-grok-api [prompt: string] {
    # Try Grok API if key available
    let api_key = (get-env "XAI_API_KEY")
    
    if ($api_key | is-empty) {
        return { success: false, error: "No API key" }
    }
    
    let response = (^curl -s -X POST "https://api.x.ai/v1/chat/completions"
        -H $"Authorization: Bearer ($api_key)"
        -H "Content-Type: application/json"
        -d '{"model": "grok-2-1212", "messages": [{"role": "user", "content": "($prompt)"}], "temperature": 0.3}'
        2>&1)
    
    try {
        let json = ($response | from json)
        if ($json.choices != null) {
            { success: true, response: $json.choices.0.message.content }
        } else {
            { success: false, error: "No response" }
        }
    } catch {
        { success: false, error: $'Curl failed: ($response)' }
    }
}

def try-grok-web [problem: string] {
    # Use DuckDuckGo to search for similar problems
    let search_query = ($problem | str replace -a " " "+")
    let url = $"https://duckduckgo.com/?q=($search_query)+BAC+math+solution"
    
    # This would use Playwright - we'll create a Python wrapper
    # For now, return failure to try next solver
    { success: false, error: "Web scraping requires Playwright" }
}

# ============================================================================
# CHATGPT SOLVER
# ============================================================================

export def solve-chatgpt [problem: string] {
    print "  Trying ChatGPT..."
    
    let prompt = (build-chatgpt-prompt $problem)
    
    # Try using free ChatGPT web via Playwright
    # This requires authenticated session
    let web_result = (try-chatgpt-web $problem)
    
    if ($web_result.success) {
        make-success $web_result.response "chatgpt-web" $web_result.steps
    } else {
        make-failure "ChatGPT unavailable (requires login)"
    }
}

def build-chatgpt-prompt [problem: string] {
    $"Tu es un professeur de mathématiques pour le BAC C Mauritanie.
Résous ce problème étape par étape en français:

($problem)

Donne:
1. Les étapes de résolution
2. Les calculs
3. La réponse finale"
}

def try-chatgpt-web [problem: string] {
    # Would use Playwright to:
    # 1. Open chatgpt.com
    # 2. Login (if credentials provided)
    # 3. Send problem
    # 4. Get response
    
    # For now, return failure
    { success: false, error: "Requires Playwright automation" }
}

# ============================================================================
# CLAUDE SOLVER
# ============================================================================

export def solve-claude [problem: string] {
    print "  Trying Claude..."
    
    # Try Claude API if available
    let api_key = (get-env "ANTHROPIC_API_KEY")
    
    if ($api_key | is-empty) {
        return (make-failure "No Anthropic API key")
    }
    
    let prompt = (build-claude-prompt $problem)
    
    let response = (^curl -s -X POST "https://api.anthropic.com/v1/messages"
        -H $"x-api-key: ($api_key)"
        -H "anthropic-version: 2023-06-01"
        -H "Content-Type: application/json"
        -d '{"model": "claude-3-haiku-20240307", "max_tokens": 2048, "messages": [{"role": "user", "content": "($prompt)"}]}'
        2>&1)
    
    try {
        let json = ($response | from json)
        if ($json.content != null) {
            let steps = (extract-steps $json.content.0.text)
            make-success $json.content.0.text "claude-api" $steps
        } else {
            make-failure "Claude API error"
        }
    } catch {
        make-failure "Claude unavailable"
    }
}

def build-claude-prompt [problem: string] {
    $"Resolve this BAC C math/physics problem step by step in French:

($problem)

Show all steps clearly."
}

# ============================================================================
# EXTRACT STEPS
# ============================================================================

def extract-steps [response: string] {
    let lines = ($response | lines)
    let steps = []
    
    for line in $lines {
        let trimmed = ($line | str trim)
        if ($trimmed =~ '^\d+[\.\)]' or $trimmed =~ '^[-•*]') {
            $steps ++= $trimmed
        }
    }
    
    if ($steps | is-empty) {
        [$response]
    } else {
        $steps
    }
}

# ============================================================================
# HELPER: GET ENV OR DEFAULT
# ============================================================================

def get-env [name: string] {
    let value = (env | get $name)
    if ($value | describe) == "nothing" {
        ""
    } else {
        $value
    }
}
