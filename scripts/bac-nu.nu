#!/usr/bin/env nu
# BAC Study CLI — Nushell native
# Usage: ./bac-nu.nu [command] [args]

const WORKER = "https://bac-api.amdajhbdh.workers.dev"

def curl [method: string, path: string, data?: any, timeout: int = 60] {
    let cmd = [
        "curl", "-s", "-X", $method, $"($WORKER)($path)"
    ]
    let cmd = if $data != null {
        [...$cmd, "-H", "Content-Type: application/json", "-d", ($data | to json)]
    } else { $cmd }
    
    let result = (^curl ...$cmd)
    if ($result | is-empty) {
        return null
    }
    
    do { $result | from json } catch { null }
}

# Query the knowledge vault
export def query [question: string, --subjects (-s): string?, --json (-j): bool = false] {
    let data = { query: $question }
    if $subjects != null {
        $data = ($data | merge { subjects: ($subjects | split row ",") })
    }
    
    let r = (curl POST "/rag/query" $data)
    
    if $r == null {
        print $"Error: Cannot reach worker"
        return null
    }
    
    if ($r | has key error) {
        print $"Error: $($r.error)"
        return null
    }
    
    if $json {
        print ($r | to json)
        return
    }
    
    print $"(ansi green_bold)Question:(ansi reset) ($question)"
    print $"(ansi cyan)Subjects:(ansi reset) (($r.subjects | to json -r))"
    print ""
    print $r.answer
    print ""
    print "(ansi yellow_bold)Sources:(ansi reset)"
    for $s in ($r.sources | each { |i| $s }) {
        let idx = ($r.sources | index-of $s | $in + 1)
        print $"  [($idx)] ($s.source) - (ansi dim)($s.subject)(ansi reset) (score: ($s.score | into string -d 2))"
    }
}

# Solve a question
export def solve [question: string, --subject (-s): string?, --json (-j): bool = false] {
    let data = { question: $question }
    if $subject != null {
        $data = ($data | merge { subject: $subject })
    }
    
    print "(ansi yellow)Solving...(ansi reset)"
    let r = (curl POST "/rag/solve" $data 120)
    
    if $r == null {
        print $"Error: Cannot solve"
        return null
    }
    
    if $json {
        print ($r | to json)
        return
    }
    
    print ""
    print "(ansi green_bold)Solution:(ansi reset)"
    print $r.solution
}

# Search vault
export def search [query: string, --subjects (-s): string?, --top (-k): int = 5, --json (-j): bool = false] {
    let data = { query: $query, topK: $top }
    if $subjects != null {
        $data = ($data | merge { subjects: ($subjects | split row ",") })
    }
    
    let r = (curl POST "/rag/search" $data)
    
    if $r == null {
        print $"Error: Search failed"
        return null
    }
    
    if $json {
        print ($r | to json)
        return
    }
    
    let chunks = $r.chunks
    if ($chunks | is-empty) {
        print "No results found"
        return null
    }
    
    print $"(ansi green_bold)Results for:(ansi reset) ($query)"
    print ""
    for $c in $chunks {
        print $"(ansi cyan)($c.metadata.subject)(ansi reset) | score: ($c.score | into string -d 3)"
        print $"  $($c.metadata.text | str slice ..200)..."
        print ""
    }
}

# Get system status
export def status [--json (-j): bool = false] {
    let r = (curl GET "/rag/status")
    let s = (curl GET "/rag/subjects")
    
    if $r == null {
        print $"Error: Cannot reach worker"
        return null
    }
    
    if $json {
        print ({status: $r, subjects: $s} | to json)
        return
    }
    
    print "(ansi green_bold)System Status(ansi reset)"
    print $"  Total Vectors: ($r.total_vectors)"
    print $"  Dimension: ($r.dimension)"
    print ""
    print "(ansi cyan)By Subject(ansi reset)"
    for $subj in ["biology", "chemistry", "math", "physics"] {
        let key = $"subject_($subj)"
        print $"  ($subj): ($s.$key)"
    }
}

# Practice question
export def practice [ --subject (-s): string?, --json (-j): bool = false ] {
    let url = if $subject != null {
        $"/rag/practice?subject=($subject)"
    } else {
        "/rag/practice"
    }
    
    let r = (curl GET url)
    
    if $r == null {
        print $"Error: No questions available"
        return null
    }
    
    if $json {
        print ($r | to json)
        return
    }
    
    let q = $r.question
    print "(ansi green_bold)Question #(ansi reset)($q.id)"
    print $q.text
    print $""
    print $"(ansi dim)Topic: ($q.topic_tags)(ansi reset)"
}
    
    let r = (curl GET $url)
    
    if $r == null {
        print $"Error: No questions available"
        return null
    }
    
    let q = $r.question
    print "(ansi green_bold)Question #(ansi reset)($q.id)"
    print $q.text
    print $""
    print $"(ansi dim)Topic: ($q.topic_tags)(ansi reset)"
}

# Track progress
export def track [ --user (-u): int = 1 ] {
    let r = (curl GET $"/rag/track?user_id=($user)")
    
    if $r == null {
        print $"Error: Cannot get stats"
        return null
    }
    
    let stats = $r.stats
    
    print "(ansi green_bold)Study Progress(ansi reset)"
    print $"  Streak: ($stats.streak) days"
    print $"  Total Reviews: ($stats.total_reviews)"
    print $"  Accuracy: ($stats.accuracy | into string -d 1)%"
    print $"  Due Today: ($stats.due_today)"
}

# Main entry point
export def main [] {
    print "(ansi green)BAC Study CLI - Nushell Edition(ansi reset)"
    print "Usage: bac-nu.nu <command> [options]"
    print ""
    print "Commands:"
    print "  query <question>     Query the vault"
    print "  solve <question>     Solve a problem"
    print "  search <query>       Search vault"
    print "  status               Show system status"
    print "  practice             Get practice question"
    print "  track                Show progress"
}