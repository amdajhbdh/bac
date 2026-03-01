# BAC UNIFIED - Simple Agent

export def agent [problem: string] {
    print ""
    print "=============================================================="
    print "  🧠 BAC UNIFIED HYPER-TUTOR"
    print "=============================================================="
    print ""
    print $"Problem: ($problem)"
    
    # Solve
    let prompt = "Resous en francais: " + $problem
    let body = "{\"model\":\"llama3.2:3b\",\"prompt\":\"" + $prompt + "\",\"stream\":false}"
    let out = (^curl -s -X POST "http://localhost:11434/api/generate" -H "Content-Type:application/json" -d $body)
    let result = try { ($out | from json | get response) } catch { "Error" }
    
    print ""
    print "=============================================================="
    print "  SOLUTION"
    print "=============================================================="
    print $result
    print ""
}
