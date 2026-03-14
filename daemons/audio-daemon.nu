#!/usr/bin/env nu
# Audio Daemon - Enhanced with podcast creation workflows

# Create resources directories
mkdir "resources/audio"
mkdir "resources/audio/podcasts"
mkdir "resources/audio/transcripts"

# Check for required tools
def checkTools [] {
    let tools = {
        ffmpeg: (which ffmpeg),
        yt_dlp: (which yt-dlp),
        nlm: (which nlm)
    }
    
    for tool in $tools {
        if ($tool.1 | is-null) {
            print $"Warning: ($tool.0) not found"
        }
    }
    
    $tools
}

# Convert video to audio
def videoToAudio [inputFile: string, outputFile: string, format: string] {
    print $"Converting: ($inputFile)"
    
    let ext = (match $format {
        "mp3" => "mp3"
        "wav" => "wav"
        "aac" => "aac"
        _ => "mp3"
    })
    
    let output = (if $outputFile == "" {
        $"resources/audio/( $inputFile | path parse | get stem).($ext)"
    } else {
        $outputFile
    })
    
    let result = (^ffmpeg -i $inputFile -vn -ab 128k -ar 44100 -y $output 2>&1)
    
    if ($result | str contains "Error") {
        print $"  ! Error: ($result)"
    } else {
        print $"  ✓ Saved: ($output)"
    }
    
    $output
}

# Extract audio from YouTube
def extractYoutubeAudio [url: string, format: string] {
    print $"Extracting audio from YouTube: ($url)"
    
    let outputTemplate = "resources/audio/podcasts/%(title)s.%(ext)s"
    let result = (^yt-dlp -x --audio-format $format -o $outputTemplate $url 2>&1)
    
    if ($result | str contains "ERROR") {
        print $"  ! Error: ($result)"
    } else {
        print $"  ✓ Downloaded"
    }
}

# Generate TTS using nlm
def generateTTS [text: string, voice: string] {
    print $"Generating TTS: ($voice)"
    
    # Check if nlm is available
    let nlmTool = (which nlm)
    if $nlmTool | is-null {
        print "  ! nlm not available, using fallback"
        return ""
    }
    
    let result = (^nlm generate audio $text --voice $voice 2>&1)
    print $"  Result: ($result)"
    
    $result
}

# Create audio summary from text
def createAudioSummary [content: string, title: string] {
    print $"Creating audio summary: ($title)"
    
    # Save content to temp file
    let tempFile = $"/tmp/audio_summary_(date now | format date '%s').txt"
    $content | save -f $tempFile
    
    # Use nlm to generate audio
    let nlmTool = (which nlm)
    if $nlmTool | is-null {
        print "  ! nlm not available"
        rm $tempFile
        return
    }
    
    # Generate podcast-style summary
    let outputFile = $"resources/audio/podcasts/($title | str replace ' ' '_' | str downcase).mp3"
    let result = (^nlm generate podcast $tempFile -o $outputFile 2>&1)
    
    if ($result | str contains "Error") {
        print $"  ! Error: ($result)"
    } else {
        print $"  ✓ Saved: ($outputFile)"
    }
    
    rm $tempFile
}

# List audio files
def listAudio [dir: string] {
    let audioDir = (if $dir == "" { "resources/audio" } else { $dir })
    
    if not ($audioDir | path exists) {
        print $"Directory not found: ($audioDir)"
        return
    }
    
    print $"Audio files in ($audioDir):"
    ls -l $audioDir | where type == "file" | each { |f|
        let size_mb = (($f.size / 1024 / 1024) | math round | into string)
        print $"  - ($f.name) (($size_mb) MB)"
    }
}

# Main audio processing function
export def process [...inputs: string] {
    print "=========================================="
    print "Audio Processing Daemon"
    print "=========================================="
    print ""
    
    let tools = checkTools
    
    if ($inputs | is-empty) {
        print "No inputs provided. Available commands:"
        print "  audio process <file_or_url>..."
        print "  audio youtube <url>..."
        print "  audio list"
        print "  audio summary <text>"
        return
    }
    
    for input in $inputs {
        if ($input | str starts-with "http") {
            # YouTube or URL
            if ($input | str contains "youtube.com") or ($input | str contains "youtu.be") {
                extractYoutubeAudio $input "mp3"
            } else {
                print $"Downloading: ($input)"
                let result = (^curl -L -o $"/tmp/audio_temp.( $input | path parse | get extension)" $input 2>&1)
                print $"  Downloaded"
            }
        } else {
            # Local file
            if ($input | path exists) {
                videoToAudio $input "" "mp3"
            } else {
                print $"File not found: ($input)"
            }
        }
    }
    
    print ""
    print "Processing complete!"
}

# Process YouTube URLs
export def youtube [...urls: string] {
    print "=========================================="
    print "YouTube Audio Extraction"
    print "=========================================="
    print ""
    
    if ($urls | is-empty) {
        print "Usage: audio youtube <url>..."
        return
    }
    
    for url in $urls {
        extractYoutubeAudio $url "mp3"
    }
}

# List audio files
export def ls [dir: string] {
    listAudio $dir
}

# Create audio summary
export def summary [title: string, ...content: string] {
    print "=========================================="
    print "Audio Summary Generation"
    print "=========================================="
    print ""
    
    let text = ($content | str join " ")
    
    if ($text | str length) == 0 {
        print "No content provided"
        return
    }
    
    createAudioSummary $text $title
}

# Convert audio formats
export def convert [input: string, format: string] {
    print "=========================================="
    print $"Converting to ($format)"
    print "=========================================="
    print ""
    
    if not ($input | path exists) {
        print $"File not found: ($input)"
        return
    }
    
    videoToAudio $input "" $format
}

# TTS generation
export def speak [text: string, voice: string] {
    print "=========================================="
    print "Text-to-Speech Generation"
    print "=========================================="
    print ""
    
    let voiceName = (if $voice == "" { "default" } else { $voice })
    print $"Voice: ($voiceName)"
    print $"Text: ($text | str length) characters"
    print ""
    
    generateTTS $text $voiceName
}

# Main entry point
export def main [] {
    print "=========================================="
    print "Audio Daemon Ready"
    print "=========================================="
    print ""
    print "Available commands:"
    print "  process <files_or_urls>  - Process audio/video files"
    print "  youtube <urls>           - Extract audio from YouTube"
    print "  list [dir]               - List audio files"
    print "  summary <title> <text>   - Create audio summary"
    print "  convert <file> <format>  - Convert audio format"
    print "  speak <text> [voice]     - Text-to-speech"
    print ""
    
    let tools = checkTools
    print ""
    
    # Show available audio files
    listAudio "resources/audio"
}
