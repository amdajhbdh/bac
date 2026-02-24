#!/usr/bin/env nu
# ============================================================================
# BAC UNIFIED - Main CLI
# Unified command-line interface for all services
# ============================================================================

# Load unified module
use lib/bac.nu

# Main entry
export def main [...args] {
    if ($args | is-empty) {
        show-help
    } else {
        let cmd = $args.0
        let rest = ($args | drop 1)
        
        match $cmd {
            "solve" => { cmd-solve $rest }
            "submit" => { cmd-submit $rest }
            "predict" => { cmd-predict $rest }
            "practice" => { cmd-practice $rest }
            "stats" => { cmd-stats }
            "leaderboard" => { cmd-leaderboard $rest }
            "audit" => { cmd-audit }
            "storage" => { 
                let storage_args = ($args | skip 1)
                cmd-storage $storage_args 
            }
            "test" => { cmd-test }
            _ => { 
                print $"Unknown command: ($cmd)"
                show-help
            }
        }
    }
}

# ============================================================================
# HELP
# ============================================================================

def show-help [] {
    print "%(char-esc)[1;36m🦅 BAC Unified%(char-reset)"
    print "National BAC Exam Preparation System"
    print ""
    print "Commands:"
    print "  solve <problem>       Solve a problem with AI"
    print "  solve <problem> -a    Solve and generate animation"
    print "  submit <file>        Submit image/PDF for OCR"
    print "  practice             Practice with questions"
    print "  predict              Predict exam questions"
    print "  stats                View your statistics"
    print "  leaderboard          View rankings"
    print "  audit                View blockchain audit log"
    print "  storage              Storage status"
    print ""
    print "Options:"
    print "  --subject <s>        Subject: math, pc, svt"
    print "  --animate, -a        Generate animation"
}

# ============================================================================
# SOLVE COMMAND - Uses Python solver with local Ollama
# ============================================================================

def cmd-solve [args: list] {
    # Parse args - problem is all non-flag arguments
    let flags = ["--subject", "-s", "--animate", "-a"]
    let problem_args = ($args | where {|x| not ($flags | any {|f| $x == $f}) })
    let problem = ($problem_args | str join " ")
    
    if ($problem | is-empty) {
        print "Usage: bac solve <problem> [--animate] [--subject <subject>]"
        return
    }
    
    let animate = ($args | any {|x| $x == "--animate" or $x == "-a"})
    let subject_arg = (get-subject-flag $args)
    let subject = if ($subject_arg | is-empty) { "math" } else { $subject_arg }
    
    print "%(char-esc)[1;36m🦅 Solving with AI...%(char-reset)"
    print $"Problem: ($problem)"
    print $"Subject: ($subject)"
    print ""
    
    # Call Python solver
    let result = (python-solve $problem)
    
    if $result.success {
        print "%(char-esc)[1;32mSolution:%(char-reset)"
        print $result.solution
        print ""
        print $"%(char-esc)[1;33m✓ Solved by: ($result.solver)%(char-reset)"
        
        # Save to database
        (save-to-database $problem $result.solution $subject)
        
        # Record audit
        (bac audit-record "solve" $problem)
        
        # Generate animation if requested
        if $animate {
            print ""
            print "Generating animation..."
            let filepath = (bac animate-generate $problem $result.solution $subject)
            print $"Animation saved: ($filepath)"
        }
    } else {
        print $"%(char-esc)[1;31mError: ($result.error)%(char-reset)"
        print ""
        print "Falling back to basic solver..."
        
        # Fallback to basic AI
        let solution = (bac ai-solve $problem)
        print ""
        print "%(char-esc)[1;32mSolution:%(char-reset)"
        print $solution
    }
}

# ============================================================================
# PYTHON SOLVER WRAPPER
# ============================================================================

def python-solve [problem: string] {
    # Call Python solver directly
    let result = (^python src/python/solver/solve.py $problem)
    
    # Parse result - simple extraction
    let output = ($result | into string)
    let lines = ($output | split row "\n")
    
    # Find solution section (between ===)
    mut $in_sol = false
    mut $sol_lines = []
    mut $solver_name = "llama3.2:3b"
    
    for ln in $lines {
        if ($ln =~ "====") {
            $in_sol = not $in_sol
        } else if $in_sol {
            $sol_lines = ($sol_lines | append $ln)
        }
        if ($ln =~ "Solver:") {
            $solver_name = ($ln | str replace -a "Solver:" "" | str trim)
        }
    }
    
    let solution = ($sol_lines | str join "\n")
    
    if (($solution | str length) > 5) {
        { success: true, solution: $solution, solver: $solver_name }
    } else {
        { success: false, error: $output }
    }
}

def get-subject-flag [args: list] {
    if ($args | is-empty) { return "" }
    let len = ($args | length)
    let max_idx = ($len - 1)
    
    for i in 0..$max_idx {
        let arg = ($args | get $i)
        if ($arg == "--subject" or $arg == "-s") {
            return ($args | get ($i + 1))
        }
    }
    ""
}

# ============================================================================
# SUBMIT COMMAND
# ============================================================================

def cmd-submit [args: list] {
    let file = ($args | get 0 | default "")
    
    if ($file | is-empty) {
        print "Usage: bac submit <file> [--subject <subject>]"
        return
    }
    
    print "%(char-esc)[1;36m📤 Submitting...%(char-reset)"
    print $"File: ($file)"
    print ""
    
    # Check file exists
    if not ($file | path exists) {
        print "Error: File not found"
        return
    }
    
    print "Step 1: OCR processing..."
    print "(OCR not implemented - use external tool)"
    print ""
    
    print "Step 2: Extracting problem text..."
    let problem_text = "Sample problem from " + $file
    print $"Extracted: ($problem_text)"
    print ""
    
    # Solve the problem
    print "Step 3: Solving..."
    let solution = (bac ai-solve $problem_text)
    print $solution
    
    # Save to database
    print ""
    (save-to-database $problem_text $solution "unknown")
    
    # Record audit
    (audit-record "submit" $file)
    
    # Upload to storage if file is valid
    let ext = ($file | path parse | get extension)
    if (["jpg", "png", "pdf"] | any {|e| $e == $ext}) {
        print ""
        print "Uploading to storage..."
        let filename = ($file | path parse | get file_name)
        let result = (storage-upload $file $filename)
        if $result {
            print "Uploaded successfully"
        }
    }
}

# ============================================================================
# PREDICT COMMAND
# ============================================================================

def cmd-predict [args: list] {
    let subject = (get-subject-flag $args)
    let year = (get-year-flag $args)
    
    print "%(char-esc)[1;36m🎯 Exam Predictions %(char-reset)" ($year)
    print ""
    
    let predictions = (bac predict --subject $subject --year $year)
    
    print "High Probability Topics:"
    for pred in $predictions {
        print $"  ($pred.subject) - ($pred.topic): ($pred.probability)"
    }
}

def get-year-flag [args: list] {
    if ($args | is-empty) { return 2026 }
    for i in 0..($args | length) {
        let arg = ($args | get $i)
        if $arg == "--year" or $arg == "-y" {
            return ($args | get ($i + 1) | into int)
        }
    }
    2026
}

# ============================================================================
# PRACTICE COMMAND
# ============================================================================

def cmd-practice [args: list] {
    let subject = (get-subject-flag $args)
    let topic = (get-topic-flag $args)
    let count = (get-count-flag $args)
    
    print "%(char-esc)[1;36m📚 Practice Mode%(char-reset)"
    print ""
    
    bac practice --subject $subject --topic $topic --count $count
}

def get-topic-flag [args: list] {
    if ($args | is-empty) { return "" }
    for i in 0..($args | length) {
        let arg = ($args | get $i)
        if $arg == "--topic" or $arg == "-t" {
            return ($args | get ($i + 1))
        }
    }
    ""
}

def get-count-flag [args: list] {
    if ($args | is-empty) { return 5 }
    for i in 0..($args | length) {
        let arg = ($args | get $i)
        if $arg == "--count" or $arg == "-c" {
            return ($args | get ($i + 1) | into int)
        }
    }
    5
}

# ============================================================================
# STATS COMMAND
# ============================================================================

def cmd-stats [] {
    bac stats
}

# ============================================================================
# LEADERBOARD COMMAND
# ============================================================================

def cmd-leaderboard [args: list] {
    let subject = (get-subject-flag $args)
    let limit = (get-limit-flag $args)
    
    bac leaderboard --subject $subject --limit $limit
}

def get-limit-flag [args: list] {
    if ($args | is-empty) { return 10 }
    for i in 0..($args | length) {
        let arg = ($args | get $i)
        if $arg == "--limit" or $arg == "-l" {
            return ($args | get ($i + 1) | into int)
        }
    }
    10
}

# ============================================================================
# AUDIT COMMAND
# ============================================================================

def cmd-audit [] {
    print "%(char-esc)[1;36m📋 Blockchain Audit Log%(char-reset)"
    print ""
    bac audit-log
}

# ============================================================================
# STORAGE COMMAND
# ============================================================================

def cmd-storage [args: list] {
    let action = if ($args | is-empty) { "list" } else { $args.0 }
    
    print "%(char-esc)[1;36m💾 Storage Status%(char-reset)"
    
    if $action == "list" {
        print "Files in Garage:"
        bac storage-list
    } else if $action == "status" {
        let config = (bac load-config)
        print $"Endpoint: ($config.garage_endpoint)"
        print $"Bucket: ($config.garage_bucket)"
    } else {
        print "Usage: bac storage [list|status]"
    }
}

# ============================================================================
# TEST COMMAND
# ============================================================================

def cmd-test [] {
    print "%(char-esc)[1;36m🧪 Testing Services%(char-reset)"
    print ""
    
    # Test database
    print "Database: "
    let db_test = (bac db-query "SELECT 1")
    if $db_test {
        print "  ✓ Connected"
    } else {
        print "  ✗ Failed"
    }
    
    # Test storage
    print "Storage: "
    let storage_test = (bac storage-list)
    if not ($storage_test == "No files") {
        print "  ✓ Connected"
    } else {
        print "  ~ Empty (OK)"
    }
    
    # Test AI
    print "AI: "
    let ai_test = (bac ai-solve "test")
    if not ($ai_test == "Solution: (AI unavailable)") {
        print "  ✓ Connected"
    } else {
        print "  ✗ Failed"
    }
    
    print ""
    print "All systems checked!"
}

# ============================================================================
# SAVE TO DATABASE
# ============================================================================

def save-to-database [problem: string, solution: string, subject: string] {
    let config = (bac load-config)
    let db_url = $config.db_url
    
    let result = (try {
        let escaped_problem = ($problem | str replace -a "'" "''")
        
        let sql = "INSERT INTO problems (subject, content, source_type) VALUES ('" + $subject + "', '" + $escaped_problem + "', 'cli')"
        ^psql $db_url -c $sql 2>/dev/null
        
        let session_sql = "INSERT INTO study_sessions (subject, topic, duration_minutes, score, notes) VALUES ('" + $subject + "', 'problem solving', 2, 100, 'AI solved')"
        ^psql $db_url -c $session_sql 2>/dev/null
        
        print "%(char-esc)[1;32m✓ Saved to database%(char-reset)"
        "ok"
    } catch {
        print "%(char-esc)[1;33m⚠ Database unavailable%(char-reset)"
        null
    })
}
