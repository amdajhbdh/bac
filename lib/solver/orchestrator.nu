# ============================================================================
# BAC UNIFIED SOLVER ORCHESTRATOR
# Routes problems to the best available solver
# Priority: Ollama → PhotoMath → Wolfram → External AI
# ============================================================================

# ============================================================================
# SOLVER ROUTING LOGIC
# ============================================================================

export def solve [
    problem: string
    --subject: string = ""
    --animate: bool = false
    --force-external: bool = false
] {
    let subject = if ($subject | is-empty) { (detect-subject $problem) } else { $subject }
    
    print $"%(char-esc)[1;36m🎯 Solving: ($problem)%(char-reset)"
    print $"Subject: ($subject)"
    print ""
    
    # Try solvers in order of preference
    let result = if $force-external {
        solve-external $problem $subject
    } else {
        # Step 1: Try local Ollama
        let ollama_result = (solve-ollama $problem)
        if ($ollama_result.success) {
            $ollama_result
        } else {
            # Step 2: Try PhotoMath for basic math
            let photomath_result = (solve-photomath $problem)
            if ($photomath_result.success) {
                $photomath_result
            } else {
                # Step 3: Try Wolfram Alpha
                let wolfram_result = (solve-wolfram $problem)
                if ($wolfram_result.success) {
                    $wolfram_result
                } else {
                    # Step 4: Try external AI via Playwright
                    solve-external $problem $subject
                }
            }
        }
    }
    
    # Save to database
    if $result.success {
        save-solution $problem $result.solution $subject $result.solver
    }
    
    # Generate animation if requested
    if $animate and $result.success {
        generate-animation $problem $result.solution $subject
    }
    
    $result
}

# ============================================================================
# SUBJECT DETECTION
# ============================================================================

export def detect-subject [problem: string] {
    let p = ($problem | str downcase)
    
    # Math keywords
    if ($p =~ 'integra|derivee|limite|derive|primitiv') { return "math" }
    if ($p =~ 'matrice|determinant|polynome|complex|equation') { return "math" }
    if ($p =~ 'geometri|angle|cercle|triangle|parallele') { return "math" }
    
    # Physics keywords
    if ($p =~ 'force|masse|vitesse|acceleration|energie') { return "pc" }
    if ($p =~ 'courant|tension|resistance|champ|induction') { return "pc" }
    if ($p =~ 'optique|ondes|son|reflexion|refraction') { return "pc" }
    
    # Chemistry keywords
    if ($p =~ 'reaction|chimique|molecule|atome|ion') { return "chimie" }
    if ($p =~ 'acide|base|ph|concentration|ionique') { return "chimie" }
    
    # SVT keywords
    if ($p =~ 'cellule|adn|gene|cromosome|heredite') { return "svt" }
    if ($p =~ 'systeme|nerveux|cardiaque|cerveau|hormone') { return "svt" }
    
    "math"  # default
}

# ============================================================================
# SOLVER RESULTS STRUCTURE
# ============================================================================

# Returns: { success: bool, solution: string, solver: string, steps: list }
export def make-result [success: bool, solution: string, solver: string, steps: list] {
    {
        success: $success
        solution: $solution
        solver: $solver
        steps: $steps
    }
}

# ============================================================================
# SOLVER RESULTS
# ============================================================================

export def make-success [solution: string, solver: string, steps: list] {
    make-result true $solution $solver $steps
}

export def make-failure [reason: string] {
    make-result false $"Solver failed: ($reason)" "none" []
}
