# YouTube Daemon - Search and download educational content
# Run: nu daemons/youtube-daemon.nu

let resources-dir = "resources/youtube"
let topics-file = "config/topics.yaml"

# Ensure directories exist
mkdir $resources-dir

# Get topics from config
let topics = (open $topics_file | grep -oP "^  - .*" | each { |line| $line | str replace -r "^  - " "" })

print "🦅 YouTube Daemon Starting..."
print $"Found (($topics | length)) topics to search"
print ""

# Search and download function
export def search-topic [topic: string, max-results: int = 5] {
    print $"Searching: ($topic)"
    
    let safe-topic = ($topic | str replace " " "+")
    
    # Use yt-dlp to search and get info
    let cmd = $"yt-dlp --match-filter \"duration<1800\" --max-downloads ($max-results) -x --audio-format mp3 -o '($resources-dir)/%(title)s.%(ext)s' \"ytsearch50:($topic) bac\""
    
    print $"  Command: ($cmd)"
    print ""
}

# Main loop
for topic in $topics {
    search-topic $topic 3
}

print "🦅 YouTube Daemon Complete!"
