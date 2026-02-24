# ============================================================================
# BAC UNIFIED SOLVER - Nushell Interface
# Calls Python solver with fallbacks
# ============================================================================

# ============================================================================
# SOLVE COMMAND
# ============================================================================

export def solve [
    problem: string
    --subject: string = ""
    --animate: bool = false
] {
    print $"%(char-esc)[1;36m🎯 Solving...%(char-reset)"
    
    let detected_subject = if ($subject == "") {
        detect-subject $problem
    } else {
        $subject
    }
    
    print $"Subject: ($detected_subject)"
    print ""
    
    # Call Python solver
    let result = (call-python-solver $problem)
    
    if $result.success {
        print $"%(char-esc)[1;32mSolution:%(char-reset)"
        print $result.solution
        print ""
        print $"%(char-esc)[1;33mSolver: ($result.solver)%(char-reset)"
        
        # Save to database
        save-solution $problem $result.solution $detected_subject
        
        # Generate animation if requested
        if $animate {
            generate-animation $problem $result.solution $detected_subject
        }
    } else {
        print $"%(char-esc)[1;31mError: ($result.error)%(char-reset)"
    }
}

# ============================================================================
# CALL PYTHON SOLVER
# ============================================================================

def call-python-solver [problem: string] {
    let script = $"python src/python/solver/solve.py '($problem)'"
    
    let output = (^nu -c $script 2>&1)
    
    # Parse output
    let lines = ($output | lines)
    let solution_lines = []
    let solver = "unknown"
    let in_solution = false
    
    for line in $lines {
        if ($line =~ "^====") {
            $in_solution = not $in_solution
            continue
        }
        if $in_solution and ($line | str length > 0) {
            $solution_lines ++= $line
        }
        if ($line =~ "Solver:") {
            $solver = ($line | str replace -a "Solver:" "" | str trim)
        }
    }
    
    let solution = ($solution_lines | str join "\n")
    
    if ($solution | str length > 10) {
        { success: true, solution: $solution, solver: $solver }
    } else {
        { success: false, error: $output }
    }
}

# ============================================================================
# DETECT SUBJECT
# ============================================================================

def detect-subject [problem: string] {
    let p = ($problem | str downcase)
    
    if ($p =~ 'force|courant|tension|physique|mecanique|optique|energie') { return "pc" }
    if ($p =~ 'chimie|reaction|molecule|atome|acide|base') { return "chimie" }
    if ($p =~ 'cellule|adn|svt|biologie|genetique|hormone') { return "svt" }
    if ($p =~ 'integrale|derivee|matrice|complexe|polynome|equation') { return "math" }
    
    "math"
}

# ============================================================================
# SAVE SOLUTION TO DATABASE
# ============================================================================

def save-solution [problem: string, solution: string, subject: string] {
    let db_url = "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb"
    
    try {
        let escaped_problem = ($problem | str replace -a "'" "''")
        
        let sql = $"INSERT INTO problems (subject, content, source_type) VALUES ('($subject)', '($escaped_problem)', 'solver')"
        ^psql $db_url -c $sql 2>/dev/null
        
        print $"Saved to database"
    } catch {
        print "Warning: Database unavailable"
    }
}

# ============================================================================
# GENERATE ANIMATION
# ============================================================================

def generate-animation [problem: string, solution: string, subject: string] {
    print ""
    print "Generating animation..."
    
    let noon_dir = "src/noon"
    let timestamp = (date now | format date "%Y%m%d_%H%M%S")
    let filename = $"solve_($timestamp).rs"
    let filepath = ($noon_dir | path join "examples" $filename)
    
    let code = $'// BAC Unified - Solution Animation
// Problem: ($problem)
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    let title = scene.text()
        .with_text("BAC - Solution")
        .with_position(-2.0, 5.0)
        .with_color(Color::YELLOW)
        .with_font_size(32)
        .make();
    scene.play(title.show_creation()).run_time(1.0);
    
    scene
}

fn main() {
    noon::app(|_| scene(app.window_rect())).run();
}
'
    
    $code | save -f $filepath
    
    print $"Animation saved: ($filepath)"
}
