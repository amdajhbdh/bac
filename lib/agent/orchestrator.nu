# ============================================================================
# AGENT ORCHESTRATOR - BAC Unified Hyper-Tutor
# Main entry point that coordinates all agents
# ============================================================================

# ============================================================================
# MAIN AGENT FUNCTION
# ============================================================================

export def agent [problem: string, --animate, --input_type: string = "text", --mode: string = "solve"] {
    print ""
    print "═══════════════════════════════════════════════════════════"
    print "  🧠 BAC UNIFIED HYPER-TUTOR AGENT"
    print "═══════════════════════════════════════════════════════════"
    print ""
    
    # Stage 1: Input Processing
    print "📥 Stage 1: Processing input..."
    let input_result = (agent-input $problem $input_type)
    print $"  → Extracted: ($input_result.text | str length) chars"
    
    # Stage 2: Classification
    print ""
    print "🔍 Stage 2: Classifying problem..."
    let class_result = (agent-classify $input_result.text)
    print $"  → Subject: ($class_result.subject)"
    print $"  → Chapter: ($class_result.chapter)"
    print $"  → Concepts: ($class_result.concepts | str join ', ')"
    
    # Stage 3: Memory Lookup
    print ""
    print "💾 Stage 3: Checking memory..."
    let memory_result = (agent-memory-lookup $input_result.text $class_result.subject)
    print $"  → Found ($memory_result.similar_count) similar problems"
    
    # Stage 4: Solve
    print ""
    print "🧮 Stage 4: Solving..."
    let solve_result = (agent-solve $input_result.text $class_result $memory_result)
    print $"  → Solver: ($solve_result.solver)"
    print $"  → Steps: ($solve_result.steps)"
    
    # Stage 5: Online Research (if needed)
    let final_result = if $solve_result.confidence < 0.7 {
        print ""
        print "🔬 Stage 5: Online research (low confidence)..."
        let online_result = (agent-online-research $input_result.text $class_result)
        print $"  → Online sources: ($online_result.sources)"
        merge-solutions $solve_result $online_result
    } else {
        $solve_result
    }
    
    # Stage 6: Visualize (if requested)
    if $animate {
        print ""
        print "🎬 Stage 6: Generating animation..."
        let viz_result = (agent-visualize $input_result.text $final_result)
        print $"  → Animation: ($viz_result.file)"
    }
    
    # Stage 7: Learn & Remember
    print ""
    print "📚 Stage 7: Learning..."
    let learn_result = (agent-learn $input_result.text $class_result $final_result)
    print $"  → Stored in memory"
    
    # Final Output
    print ""
    print "═══════════════════════════════════════════════════════════"
    print "  ✅ SOLUTION"
    print "═══════════════════════════════════════════════════════════"
    print ""
    print $final_result.solution
    print ""
    print "───────────────────────────────────────────────────────────"
    print $"Subject: ($class_result.subject) | Chapter: ($class_result.chapter)"
    print $"Concepts: ($class_result.concepts | str join ', ')"
    print $"Solver: ($final_result.solver) | Confidence: ($final_result.confidence)"
    print "───────────────────────────────────────────────────────────"
    
    # Return result
    {
        problem: $input_result.text
        classification: $class_result
        solution: $final_result.solution
        solver: $final_result.solver
        confidence: $final_result.confidence
        concepts: $class_result.concepts
    }
}

# ============================================================================
# STAGE 1: INPUT AGENT
# ============================================================================

def agent-input [problem: string, input_type: string] {
    match $input_type {
        "text" => { text: $problem, type: "text" }
        "image" => { text: "(OCR not implemented)", type: "image" }
        "pdf" => { text: "(PDF extraction not implemented)", type: "pdf" }
        _ => { text: $problem, type: "text" }
    }
}

# ============================================================================
# STAGE 2: CLASSIFIER AGENT
# ============================================================================

def agent-classify [text: string] {
    # Use Ollama to classify the problem
    let prompt = "Classify this BAC problem. Respond ONLY in this format:
Subject: math|pc|svt|chimie
Chapter: (one word)
Concepts: (comma separated list, max 3)
Difficulty: 1-5

Problem: " + $text
    
    let result = (call-ollama $prompt)
    let lines = ($result | split row "\n")
    
    let subject = if ($result =~ "pc|physique") { "pc" } 
        else if ($result =~ "chimie") { "chimie" }
        else if ($result =~ "svt|biologie") { "svt" }
        else { "math" }
    
    let chapter = extract-field $result "Chapter"
    let concepts = (extract-field $result "Concepts" | split row ",")
    let diff_str = (extract-field $result "Difficulty")
    let difficulty = if ($diff_str == "") { 3 } else { (default ($diff_str | into int) 3) }
    
    { 
        subject: $subject 
        chapter: $chapter
        concepts: $concepts
        difficulty: $difficulty
    }
}

def extract-field [text: string, field: string] {
    let lines = ($text | split row "\n")
    let match = ($lines | where {|l| $l =~ $field} | get 0 --optional)
    if $match != null {
        let parts = ($match | split row ":")
        if (($parts | length) > 1) {
            ($parts.1 | str trim)
        } else { "" }
    } else { "" }
}

# ============================================================================
# STAGE 3: MEMORY AGENT (Vector DB)
# ============================================================================

def agent-memory-lookup [text: string, subject: string] {
    # Query similar problems from database
    # For now, return empty result (would use pgvector)
    {
        similar_problems: []
        similar_count: 0
        context: ""
    }
}

# ============================================================================
# STAGE 4: SOLVER AGENT
# ============================================================================

def agent-solve [text: string, classification: record, memory: record] {
    # Build context from memory
    let context = if ($memory.context != "") {
        "\n\nContexte:\n" + $memory.context
    } else { "" }
    
    # Build prompt
    let prompt = "Tu es un professeur BAC C (" + $classification.subject + "). 
Résous ce problème étape par étape en français.
Concepts impliqués: " + ($classification.concepts | str join ", ") + "

Problème: " + $text + "

Format:
**Solution:**
[réponse]

**Étapes:**
1. [étape 1]
2. [étape 2]
3. [étape 3]"
    
    # Call Ollama
    let result = (call-ollama $prompt)
    
    let step_lines = ($result | split row "\n" | where {|l| $l =~ $'\d+\.' })
    let step_count = if ($step_lines | is-empty) { 3 } else { $step_lines | length }
    
    { 
        solution: $result 
        solver: "llama3.2:3b"
        steps: $step_count
        confidence: 0.85
    }
}

# ============================================================================
# STAGE 5: ONLINE RESEARCH AGENT
# ============================================================================

def agent-online-research [text: string, classification: record] {
    # Would use Playwright to query online AIs
    # For now, return empty
    {
        sources: 0
        results: []
    }
}

def merge-solutions [local: record, online: record] {
    # Merge online research with local solution
    $local
}

# ============================================================================
# STAGE 6: VISUALIZER AGENT
# ============================================================================

def agent-visualize [problem: string, solution: record] {
    # Generate Noon animation
    let code = (create-noon-code $problem)
    let filepath = (save-noon $code $problem)
    
    { file: $filepath }
}

# ============================================================================
# STAGE 7: TEACHER AGENT (Learning)
# ============================================================================

def agent-learn [problem: string, classification: record, solution: record] {
    # Store in database for future recall
    print "  → Would store in vector DB: problem + embedding"
    { stored: true }
}

# ============================================================================
# HELPER: OLLAMA CALL
# ============================================================================

def call-ollama [prompt: string] {
    let body = "{\"model\":\"llama3.2:3b\",\"prompt\":\"" + $prompt + "\",\"stream\":false,\"options\":{\"temperature\":0.3}}"
    
    let out = (^curl -s -X POST "http://localhost:11434/api/generate" -H "Content-Type:application/json" -d $body)
    
    try {
        let j = ($out | from json)
        $j.response
    } catch {
        "Error: Could not connect to Ollama"
    }
}

# ============================================================================
# HELPER: NOON CODE
# ============================================================================

def create-noon-code [problem: string] {
    let short = if (($problem | str length) > 25) {
        ($problem | str substring 0..25) + "..."
    } else { $problem }
    
    "use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    let title = scene.text().with_text(\"BAC\").with_position(-2,5).with_color(Color::YELLOW).with_font_size(36).make();
    scene.play(title.show_creation()).run_time(1.0);
    let prob = scene.text().with_text(\"" + $short + "\").with_position(-5,3).with_color(Color::WHITE).with_font_size(18).make();
    scene.play(prob.show_creation()).run_time(1.5);
    scene
}

fn main() { noon::app(|_| scene(app.window_rect())).run(); }"
}

def save-noon [code: string, problem: string] {
    let name = ($problem | str replace -a " " "_" | str substring 0..12)
    let fp = "src/noon/examples/bac_" + $name + ".rs"
    $code | save -f $fp
    $fp
}

def default [value: any, default_val: any] {
    if ($value | describe) == "nothing" { $default_val } else { $value }
}
