# Audio Daemon - Generate NLM audio overviews
# Run: nu daemons/audio-daemon.nu

let resources-dir = "resources/audio"

# Ensure directories exist
mkdir $resources-dir

print "🎧 Audio Daemon Starting..."

# Get notebooks
let notebooks = (nlm notebook list)

print $"Found (($notebooks | length)) notebooks"
print ""

# Generate audio for each notebook
export def generate-audio [notebook-id: string, title: string] {
    print $"Generating audio for: ($title)"
    
    nlm audio create $notebook-id
    
    # Download audio
    nlm download audio $notebook-id --output $resources-dir
    
    print "  Audio generated!"
    print ""
}

# Main loop
for notebook in $notebooks {
    generate-audio $notebook.id $notebook.title
}

print "🎧 Audio Daemon Complete!"
