# NLM Research Daemon - Auto-discover new sources
# Run: nu daemons/nlm-research-daemon.nu

let topics-file = "config/topics.yaml"

# Get topics from config
let topics = (open $topics_file | grep -oP "^  - .*" | each { |line| $line | str replace -r "^  - " "" })

print "🔬 NLM Research Daemon Starting..."
print $"Found (($topics | length)) topics to research"
print ""

# Check authentication
print "Checking NLM authentication..."
nlm login --check

# Research function
export def research-topic [query: string] {
    print $"Researching: ($query)"
    
    # Start research task
    nlm research start $query --mode deep
    
    print "  Research started..."
    print ""
}

# Main loop
for topic in $topics {
    research-topic $topic
}

print "🔬 NLM Research Daemon Complete!"
