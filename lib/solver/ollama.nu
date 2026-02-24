# ============================================================================
# OLLAMA LOCAL SOLVER
# Uses local Ollama AI for problem solving
# ============================================================================

use ../orchestrator.nu make-failure make-success

# ============================================================================
# MAIN OLLAMA SOLVER
# ============================================================================

export def solve-ollama [problem: string] {
    let config = (load-config)
    let endpoint = $config.ollama_endpoint
    let model = $config.ollama_model
    
    print "Trying Ollama (local)..."
    
    # Check if Ollama is running
    let status = (check-ollama-status $endpoint)
    if not $status {
        return (make-failure "Ollama not running")
    }
    
    # Build BAC-specific prompt
    let prompt = (build-bac-prompt $problem)
    
    # Call Ollama
    let result = (call-ollama $endpoint $model $prompt)
    
    if ($result.success) {
        let steps = (extract-steps $result.response)
        make-success $result.response "ollama" $steps
    } else {
        make-failure $result.error
    }
}

# ============================================================================
# CHECK OLLAMA STATUS
# ============================================================================

def check-ollama-status [endpoint: string] {
    let response = (^curl -s -X GET $"($endpoint)/api/tags" 2>/dev/null)
    if ($response == "") {
        return false
    }
    try {
        let json = ($response | from json)
        $json.models != null
    } catch {
        false
    }
}

# ============================================================================
# BUILD BAC-SPECIFIC PROMPT
# ============================================================================

def build-bac-prompt [problem: string] {
    $"You are an expert math and physics tutor for Mauritania's BAC C exam.
Your task is to solve the following problem with clear, step-by-step explanations in French.

Problem: ($problem)

Requirements:
1. Show ALL calculation steps
2. Explain each step in French
3. Use proper mathematical notation (LaTeX)
4. If it's a geometry problem, describe the construction
5. If it's a physics problem, identify the laws used
6. Give the final answer clearly

Format your response as:
**Solution:**
[Your solution here]

**Steps:**
1. [Step 1]
2. [Step 2]
3. [Step 3]
... 

Remember: You are helping a student prepare for their BAC C exam. Be clear, patient, and educational."
}

# ============================================================================
# CALL OLLAMA API
# ============================================================================

def call-ollama [endpoint: string, model: string, prompt: string] {
    let body = $'{"model": "($model)", "prompt": "($prompt)", "stream": false, "options": {"temperature": 0.3, "num_predict": 2048}}'
    
    let response = (^curl -s -X POST $"($endpoint)/api/generate" 
        -H "Content-Type: application/json" 
        -d $body 2>&1)
    
    if ($response == "") {
        return { success: false, error: "Empty response" }
    }
    
    try {
        let json = ($response | from json)
        if ($json.response | is-not-nothing) {
            { success: true, response: $json.response }
        } else {
            { success: false, error: "No response field" }
        }
    } catch {
        { success: false, error: $'Parse error: ($response)' }
    }
}

# ============================================================================
# EXTRACT STEPS FROM RESPONSE
# ============================================================================

def extract-steps [response: string] {
    let lines = ($response | lines)
    let steps = []
    
    for line in $lines {
        if ($line =~ '^\d+\.' or $line =~ '^[-•*]') {
            $steps ++= ($line | str trim)
        }
    }
    
    if ($steps | is-empty) {
        [$response]
    } else {
        $steps
    }
}

# ============================================================================
# LOAD CONFIG
# ============================================================================

def load-config [] {
    let db_url = "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb"
    let ollama_endpoint = "http://localhost:11434"
    let ollama_model = "minimax-m2.5:cloud"
    
    {
        db_url: $db_url
        ollama_endpoint: $ollama_endpoint
        ollama_model: $ollama_model
    }
}

# ============================================================================
# SAVE SOLUTION TO DATABASE
# ============================================================================

def save-solution [problem: string, solution: string, subject: string, solver: string] {
    let config = (load-config)
    let db_url = $config.db_url
    
    try {
        let escaped_problem = ($problem | str replace -a "'" "''")
        let escaped_solution = ($solution | str replace -a "'" "''")
        
        let sql = $"INSERT INTO problems (subject, content, source_type) VALUES ('($subject)', '($escaped_problem) - Solver: ($solver)', 'solver')"
        ^psql $db_url -c $sql 2>/dev/null
        
        print $"Saved to database (solver: ($solver))"
    } catch {
        print "Warning: Could not save to database"
    }
}

# ============================================================================
# GENERATE ANIMATION
# ============================================================================

def generate-animation [problem: string, solution: string, subject: string] {
    print ""
    print "Generating animation with Noon..."
    
    let config = (load-config)
    let noon_dir = $config.noon_dir
    
    let timestamp = (date now | format date "%Y%m%d_%H%M%S")
    let filename = $"solve_($timestamp).rs"
    let filepath = ($noon_dir | path join "examples" $filename)
    
    let code = $'// BAC Unified - Auto-generated Animation
// Problem: ($problem)
// Solver: AI
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Title
    let title = scene.text()
        .with_text("BAC Unified - Solution")
        .with_position(-2.0, 5.0)
        .with_color(Color::YELLOW)
        .with_font_size(32)
        .make();
    scene.play(title.show_creation()).run_time(1.0);
    
    // Problem
    let prob = scene.text()
        .with_text("($problem)")
        .with_position(-5.0, 3.0)
        .with_color(Color::WHITE)
        .with_font_size(16)
        .make();
    scene.play(prob.show_creation()).run_time(1.5);
    
    // Solution header
    let sol_header = scene.text()
        .with_text("Solution:")
        .with_position(-5.0, 0.5)
        .with_color(Color::GREEN)
        .with_font_size(18)
        .make();
    scene.play(sol_header.show_creation()).run_time(0.5);
    
    // First few lines of solution (truncated for display)
    let sol_lines = ($solution | str substring 0..200)
    let sol = scene.text()
        .with_text("($sol_lines)")
        .with_position(-5.0, -1.0)
        .with_color(Color::CYAN)
        .with_font_size(14)
        .make();
    scene.play(sol.show_creation()).run_time(2.0);
    
    scene
}

fn main() {
    noon::app(|_| scene(app.window_rect())).run();
}
'
    
    $code | save -f $filepath
    
    print $"Animation code saved to: ($filepath)"
    print "To build and run:"
    print $"  cd ($noon_dir) && cargo run --example ($filename | path parse | get stem)"
}

# Helper for is-not-nothing
def is-not-nothing [x: any] {
    if ($x | describe) == "nothing" {
        false
    } else {
        true
    }
}
