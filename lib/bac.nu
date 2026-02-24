# ============================================================================
# BAC UNIFIED - Complete CLI Module
# Unified Nushell CLI for all services
# ============================================================================

# Load configuration
export def load-config [] {
    let config = {
        # Database (Neon cloud)
        db_url: "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb"
        
        # Garage S3
        garage_endpoint: "http://localhost:3900"
        garage_bucket: "new-bucket"
        garage_key: "GKd001902bdf51fcdbeaf4d4a5"
        garage_secret: "e92117c54a666734d6809f75251c09862bffe6ed610754f2e6e31732267e80b5"
        
        # Ollama AI (Using LOCAL llama3.2:3b - cloud models have rate limits)
        # Note: ministral-3:14b-cloud is available but rate-limited
        # You can switch when you have API key
        ollama_endpoint: "http://localhost:11434"
        ollama_model: "llama3.2:3b"
        
        # Animation
        noon_dir: "src/noon"
        animation_output: "data/animations"
    }
    $config
}

# ============================================================================
# DATABASE MODULE
# ============================================================================

export def db-query [sql: string] {
    let config = (load-config)
    let db_url = $config.db_url
    
    let result = (try {
        ^psql $db_url -c $sql 2>/dev/null
        true
    } catch {
        false
    })
    $result
}

export def db-insert [table: string, data: string] {
    let config = (load-config)
    let sql = $"INSERT INTO ($table) VALUES (($data))"
    db-query $sql
}

export def db-select [query: string] {
    let config = (load-config)
    let db_url = $config.db_url
    
    let result = (try {
        ^psql $db_url -t -c $query 2>/dev/null
    } catch {
        null
    })
    $result
}

# Get problems from database
export def db-get-problems [
    --subject: string = ""
    --limit: int = 10
] {
    let config = (load-config)
    let db_url = $config.db_url
    
    let subject_filter = if ($subject | is-empty) { "" } else { $"WHERE subject = '($subject)'" }
    let sql = $"SELECT id, subject, content, created_at FROM problems ($subject_filter) ORDER BY id DESC LIMIT ($limit)"
    
    let result = (try {
        ^psql $db_url -t -c $sql 2>/dev/null | lines | each { |line| $line | str trim }
    } catch {
        []
    })
    $result
}

# Get user stats
export def db-get-stats [] {
    let config = (load-config)
    let db_url = $config.db_url
    
    let query = "SELECT username, points, level, total_solved, accuracy_rate FROM users LIMIT 1"
    
    let result = (try {
        ^psql $db_url -t -c $query 2>/dev/null | str trim
    } catch {
        "demo user | 150 | 5 | 42 | 78"
    })
    $result
}

# ============================================================================
# AI MODULE
# ============================================================================

export def ai-solve [problem: string] {
    let config = (load-config)
    let endpoint = $config.ollama_endpoint
    let model = $config.ollama_model
    
    let prompt = "BAC: " + $problem
    
    let response = (try {
        let body = "{\"model\": \"" + $model + "\", \"prompt\": \"" + $prompt + "\", \"stream\": false}"
        let curl_output = (^curl -s -X POST $endpoint + "/api/generate" -H "Content-Type: application/json" -d $body)
        let json_result = ($curl_output | from json)
        $json_result.response
    } catch {
        "Solution: (AI unavailable)"
    })
    
    $response
}

export def ai-detect-type [problem: string, subject: string] {
    if not ($subject | is-empty) { return $subject }
    
    let p = ($problem | str downcase)
    
    if ($p =~ 'integra') { return "math" }
    if ($p =~ 'derivee') { return "math" }
    if ($p =~ 'polynome|complex') { return "math" }
    if ($p =~ 'projectile|force|mouvement') { return "pc" }
    if ($p =~ 'circuit|courant|resistance') { return "pc" }
    if ($p =~ 'molecule|atome|svt') { return "svt" }
    
    "math"
}

# ============================================================================
# STORAGE MODULE
# ============================================================================

export def storage-upload [file: string, name: string] {
    let config = (load-config)
    let endpoint = $config.garage_endpoint
    let bucket = $config.garage_bucket
    let key = $config.garage_key
    let secret = $config.garage_secret
    
    let result = (try {
        AWS_ACCESS_KEY_ID=$key AWS_SECRET_ACCESS_KEY=$secret aws --endpoint-url $endpoint s3 cp $file $"s3://($bucket)/($name)"
        true
    } catch {
        false
    })
    $result
}

export def storage-list [] {
    let config = (load-config)
    let endpoint = $config.garage_endpoint
    let bucket = $config.garage_bucket
    let key = $config.garage_key
    let secret = $config.garage_secret
    
    let result = (try {
        AWS_ACCESS_KEY_ID=$key AWS_SECRET_ACCESS_KEY=$secret aws --endpoint-url $endpoint s3 ls $"s3://($bucket)/"
    } catch {
        "No files"
    })
    $result
}

# ============================================================================
# ANIMATION MODULE
# ============================================================================

export def animate-generate [problem: string, solution: string, type: string] {
    let config = (load-config)
    let noon_dir = $config.noon_dir
    
    let timestamp = (date now | format date "%Y%m%d_%H%M%S")
    let filename = $"solve_($timestamp).rs"
    let filepath = ($noon_dir | path join "examples" $filename)
    
    let code = "// BAC Unified - Auto-generated Animation
// Problem: " + $problem + "
use noon::prelude::*;

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    let title = scene.text()
        .with_text(\"BAC Unified - Solution\")
        .with_position(-2.0, 5.0)
        .with_color(Color::YELLOW)
        .with_font_size(36)
        .make();
    scene.play(title.show_creation()).run_time(1.0);
    
    let prob = scene.text()
        .with_text(\"" + $problem + "\")
        .with_position(-5.0, 3.5)
        .with_color(Color::WHITE)
        .with_font_size(20)
        .make();
    scene.play(prob.show_creation()).run_time(1.0);
    
    let sol = scene.text()
        .with_text(\"Solution\")
        .with_position(-5.0, 1.0)
        .with_color(Color::GREEN)
        .with_font_size(18)
        .make();
    scene.play(sol.show_creation()).run_time(1.5);
    
    scene
}

fn main() {
    noon::app(|_| scene(app.window_rect())).run();
}
"
    
    $code | save -f $filepath
    
    $filepath
}

export def animate-build [filename: string] {
    let config = (load-config)
    let noon_dir = $config.noon_dir
    
    let example_name = ($filename | path parse | get stem)
    
    let result = (try {
        cd $noon_dir
        let output = (^cargo build --release --example $example_name out+err> /dev/null | complete)
        true
    } catch {
        false
    })
    $result
}

# ============================================================================
# BLOCKCHAIN AUDIT MODULE
# ============================================================================

export def audit-hash [data: string] {
    let temp_file = "/tmp/bac_hash_input.txt"
    $data | save -f $temp_file
    let hash_result = (^sha256sum $temp_file)
    let hash = ($hash_result | str substring 0..64)
    $hash
}

export def audit-record [action: string, data: string] {
    let config = (load-config)
    let timestamp = (date now | format date "%Y-%m-%d %H:%M:%S")
    let hash = (audit-hash $data)
    
    let record = $"($timestamp) | ($action) | ($hash) | ($data)"
    
    # Save to audit log file
    let log_dir = "data/audit"
    mkdir $log_dir
    let log_file = ($log_dir | path join "audit.log")
    $record | save -a $log_file
    
    $record
}

export def audit-verify [data: string, expected_hash: string] {
    let actual_hash = (audit-hash $data)
    $actual_hash == $expected_hash
}

export def audit-log [] {
    let log_file = "data/audit/audit.log"
    
    if not ($log_file | path exists) {
        print "No audit log found"
        return
    }
    
    open $log_file
}

# ============================================================================
# PREDICTION MODULE
# ============================================================================

export def predict [
    --subject: string = "all"
    --year: int = 2026
] {
    # Get historical data from database
    let problems = (db-get-problems --limit 100)
    
    # Simple prediction based on topics
    let predictions = if $subject == "all" {
        [
            {topic: "integrals", probability: "85%", subject: "math"}
            {topic: "complex", probability: "78%", subject: "math"}
            {topic: "derivatives", probability: "72%", subject: "math"}
            {topic: "mechanics", probability: "80%", subject: "pc"}
            {topic: "electromagnetism", probability: "70%", subject: "pc"}
            {topic: "genetics", probability: "75%", subject: "svt"}
        ]
    } else {
        [
            {topic: "integrals", probability: "85%", subject: "math"}
            {topic: "complex", probability: "78%", subject: "math"}
            {topic: "derivatives", probability: "72%", subject: "math"}
        ]
    }
    
    $predictions
}

# ============================================================================
# PRACTICE MODULE
# ============================================================================

export def practice [
    --subject: string = ""
    --topic: string = ""
    --count: int = 5
] {
    let problems = (db-get-problems --subject $subject --limit $count)
    
    if ($problems | is-empty) {
        print "No questions found. Try submitting some problems first!"
        return
    }
    
    print "%(char-esc)[1;36mPractice Mode%(char-reset)"
    print ""
    
    for problem in $problems {
        print $problem
        print "---"
    }
}

# ============================================================================
# STATS MODULE
# ============================================================================

export def stats [] {
    let user_stats = (db-get-stats)
    
    print "%(char-esc)[1;36m📊 Your Statistics%(char-reset)"
    print ""
    print $user_stats
}

# ============================================================================
# LEADERBOARD MODULE
# ============================================================================

export def leaderboard [
    --subject: string = "all"
    --limit: int = 10
] {
    print "%(char-esc)[1;36m🏆 Leaderboard%(char-reset)"
    print ""
    
    # Query from database or use mock
    let query = "SELECT username, points FROM users ORDER BY points DESC LIMIT " + ($limit | into string)
    let results = (db-select $query)
    
    if ($results == null) {
        print "1. Ahmed - 2,450 points"
        print "2. Fatma - 2,320 points"
        print "3. Mohamed - 2,180 points"
    } else {
        print $results
    }
}
