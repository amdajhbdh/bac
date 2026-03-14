#!/usr/bin/env nu
# YouTube Daemon - Enhanced with dynamic topic input

# Create resources directory
mkdir "resources/youtube"
mkdir "resources/youtube/audio"
mkdir "resources/youtube/video"

# Check for required tools
def checkTool [tool: string] {
    let path = (which $tool)
    if $path == null {
        print $"Error: ($tool) not found. Please install it first."
        exit 1
    }
    print $"($tool) found at: ($path)"
}

# Download audio from YouTube search
def downloadAudio [topic: string, limit: int] {
    let output_dir = "resources/youtube/audio"
    let search_query = $"ytsearch($limit):($topic)"
    
    print $"  Downloading audio for: ($topic)"
    
    let result = (^yt-dlp 
        -x 
        --audio-format mp3 
        --audio-quality 0 
        -o $"($output_dir)/%(title)s.%(ext)s"
        --yes-playlist
        $search_query
        2>&1)
    
    if $result == "" {
        print $"    ✓ Downloaded: ($topic)"
    } else {
        print $"    ! Warning: ($result)"
    }
}

# Download video from YouTube search  
def downloadVideo [topic: string, limit: int] {
    let output_dir = "resources/youtube/video"
    let search_query = $"ytsearch($limit):($topic)"
    
    print $"  Downloading video for: ($topic)"
    
    let result = (^yt-dlp 
        -f "bestvideo[height<=720]+bestaudio/best[height<=720]"
        --merge-output-format mp4
        -o $"($output_dir)/%(title)s.%(ext)s"
        --yes-playlist
        $search_query
        2>&1)
    
    if $result == "" {
        print $"    ✓ Downloaded: ($topic)"
    } else {
        print $"    ! Warning: ($result)"
    }
}

# Get playlist info without downloading
def getPlaylistInfo [topic: string] {
    let search_query = $"ytsearch5:($topic)"
    
    print $"  Getting info for: ($topic)"
    
    let result = (^yt-dlp 
        --dump-json 
        --no-download
        $search_query
        2>&1 | head -n 1 | from json)
    
    print $"    Title: ($result.title)"
    print $"    Duration: ($result.duration)"
    print $"    View count: ($result.view_count)"
}

# Main entry point
export def main [...topics: string] {
    print "=========================================="
    print "YouTube Educational Content Downloader"
    print "=========================================="
    print ""
    
    # Check for required tools
    checkTool "yt-dlp"
    
    # Default topics if none provided
    let topics = if ($topics | is-empty) {
        print "No topics provided, using defaults..."
        [
            "BAC math Terminale C intégrales",
            "BAC math Terminale C suites numériques",
            "BAC physique Terminale C mécanique",
            "BAC chimie Terminale C organique"
        ]
    } else {
        $topics
    }
    
    print $""
    print $"Downloading content for: ($topics | length) topics"
    print $""
    
    # Download audio for each topic
    print "--- Audio Downloads ---"
    for topic in $topics {
        downloadAudio $topic 3
    }
    
    print $""
    print "--- Video Downloads ---"
    for topic in $topics {
        downloadVideo $topic 2
    }
    
    print $""
    print "=========================================="
    print "Download Complete!"
    print "=========================================="
    print $""
    print "Files saved to:"
    print "  - Audio: resources/youtube/audio/"
    print "  - Video: resources/youtube/video/"
}

# Export additional commands
export def search [query: string, limit: int] {
    checkTool "yt-dlp"
    
    let search_query = $"ytsearch($limit):($query)"
    
    print $"Searching for: ($query)"
    print $""
    
    ^yt-dlp --dump-json --no-download $search_query | each { |video|
        print $"Title: ($video.title)"
        print $"Duration: ($video.duration_string)"
        print $"Views: ($video.view_count)"
        print $"Channel: ($video.channel)"
        print $""
    }
}

export def playlist [url: string, audio_only: bool] {
    checkTool "yt-dlp"
    
    print $"Downloading playlist: ($url)"
    print $""
    
    if $audio_only {
        ^yt-dlp -x --audio-format mp3 -o "resources/youtube/audio/%(title)s.%(ext)s" $url
    } else {
        ^yt-dlp -f "bestvideo[height<=720]+bestaudio/best[height<=720]" --merge-output-format mp4 -o "resources/youtube/video/%(title)s.%(ext)s" $url
    }
    
    print "Playlist download complete!"
}
