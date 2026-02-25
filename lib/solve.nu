#!/usr/bin/env nu
# BAC UNIFIED SOLVER - Pure Nushell with AI + Smart Noon animations

export def solve [problem: string, --animate] {
    print $"Problem: ($problem)"
    print "Solving..."
    
    # Get solution from AI
    let result = (call-ollama $problem)
    
    if $result.success {
        print ""
        print "=== SOLUTION ==="
        print $result.solution
        print "==============="
        print $"Solver: ($result.solver)"
        
        # Generate AI-powered animation if requested
        if $animate {
            print ""
            print "Generating Noon animation (AI-powered)..."
            
            # Extract key info from solution for animation
            let info = (extract-solution-info $problem $result.solution)
            let fp = (create-smart-animation $problem $info)
            
            print $"Saved: ($fp)"
            print "Run: cd src/noon && cargo run --example" ($fp | path parse | get stem)
        }
    } else {
        print $"Error: ($result.error)"
    }
}

# ============================================================================
# OLLAMA SOLVER
# ============================================================================

def call-ollama [problem: string] {
    let prompt = "Tu es prof maths BAC C. Resous: " + $problem + " Donne etapes en francais."
    let body = "{\"model\":\"llama3.2:3b\",\"prompt\":\"" + $prompt + "\",\"stream\":false,\"options\":{\"temperature\":0.3}}"
    
    let out = (^curl -s -X POST "http://localhost:11434/api/generate" -H "Content-Type:application/json" -d $body)
    
    try {
        let j = ($out | from json)
        { success: true, solution: $j.response, solver: "llama3.2:3b" }
    } catch {
        { success: false, error: $out }
    }
}

# ============================================================================
# EXTRACT SOLUTION INFO FOR ANIMATION
# ============================================================================

def extract-solution-info [problem: string, solution: string] {
    # Extract key parts from solution for animation
    
    # Get first step
    let step1 = if ($solution =~ "etape 1" or $solution =~ "Étape 1" or $solution =~ "ETAPE 1") {
        let parts = ($solution | split row "Étape 1")
        if (($parts | length) > 1) {
            let content = $parts.1
            let lines = ($content | split row "\n")
            ($lines.0 | str substring 0..60)
        } else { "Analyser le probleme" }
    } else { "Analyser le probleme" }
    
    # Get answer (look for x = or answer)
    let answer = if ($solution =~ "x =") {
        let parts = ($solution | split row "x =")
        let line = ($parts.1 | split row "\n" | get 0 | str trim)
        "x = " + ($line | str substring 0..15)
    } else if ($solution =~ "reponse" or $solution =~ "Réponse") {
        let parts = ($solution | split row "Réponse")
        ($parts.1 | split row "\n" | get 0 | str substring 0..30 | str trim)
    } else { "Solution trouvee" }
    
    { step1: $step1, answer: $answer, problem: $problem }
}

# ============================================================================
# CREATE SMART ANIMATION (TEMPLATE-BASED ON EXTRACTED INFO)
# ============================================================================

def create-smart-animation [problem: string, info: record] {
    let short = if (($problem | str length) > 25) {
        ($problem | str substring 0..25) + "..."
    } else { $problem }
    
    let step1 = if (($info.step1 | str length) > 40) {
        ($info.step1 | str substring 0..40)
    } else { $info.step1 }
    
    let answer = $info.answer
    
    let name = ($problem | str replace -a " " "_" | str replace -a "+" "plus" | str replace -a "=" "eq" | str replace -a "[" "" | str replace -a "]" "" | str substring 0..12)
    let filename = "bac_" + $name + ".rs"
    let filepath = "src/noon/examples/" + $filename
    
    let code = "// BAC Unified - AI-Generated Animation
// Problem: " + $problem + "
// Solution steps extracted by AI

use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // === TITLE ===
    let title = scene.text()
        .with_text(\"BAC - Resolution\")
        .with_position(-2.5, 5.5)
        .with_color(Color::YELLOW)
        .with_font_size(40)
        .make();
    scene.play(title.show_creation()).run_time(1.0);
    
    // === PROBLEM ===
    let prob = scene.text()
        .with_text(\"" + $short + "\")
        .with_position(-5.0, 4.0)
        .with_color(Color::WHITE)
        .with_font_size(22)
        .make();
    scene.play(prob.show_creation()).run_time(1.5);
    scene.wait();
    
    // === STEP 1 ===
    let s1 = scene.text()
        .with_text(\"1. \" + \"" + $step1 + "\")
        .with_position(-5.0, 2.0)
        .with_color(Color::GREEN)
        .with_font_size(18)
        .make();
    scene.play(s1.show_creation()).run_time(1.2);
    
    // === STEP 2 (Calculation) ===
    let s2 = scene.text()
        .with_text(\"2. Calculer\")
        .with_position(-5.0, 0.5)
        .with_color(Color::GREEN)
        .with_font_size(18)
        .make();
    scene.play(s2.show_creation()).run_time(1.0);
    
    // === STEP 3 (Verification) ===
    let s3 = scene.text()
        .with_text(\"3. Verifier\")
        .with_position(-5.0, -1.0)
        .with_color(Color::GREEN)
        .with_font_size(18)
        .make();
    scene.play(s3.show_creation()).run_time(1.0);
    
    // === FINAL ANSWER ===
    let ans = scene.text()
        .with_text(\"" + $answer + "\")
        .with_position(0.0, -3.5)
        .with_color(Color::CYAN)
        .with_font_size(36)
        .make();
    scene.play(ans.show_creation()).run_time(2.0);
    
    scene
}

fn main() {
    noon::app(|_| scene(app.window_rect())).run();
}
"
    
    $code | save -f $filepath
    $filepath
}
