#!/usr/bin/env nu
# Web Crawler Daemon - Enhanced with scraping logic

# Create resources directory
mkdir "resources/web"
mkdir "resources/web/html"
mkdir "resources/web/pdf"
mkdir "resources/web/links"

# Check for required tools
def checkTool [tool: string] {
    let path = (which $tool)
    if $path == null {
        print $"Warning: ($tool) not found"
    }
}

# Fetch page content
def fetchPage [url: string] {
    print $"  Fetching: ($url)"
    
    let result = (^curl -s -L --max-time 30 -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" $url 2>&1)
    
    if ($result | str length) > 0 {
        print $"    ✓ Fetched ($result | str length) bytes"
        $result
    } else {
        print $"    ! Failed to fetch"
        ""
    }
}

# Extract PDF links from HTML
def extractPdfLinks [html: string] {
    let pattern = 'href="([^"]+\.pdf)"'
    let links = ($html | grep -oP $pattern | each { |m| $m | str replace -r 'href="|"' '' })
    
    # Also try other patterns
    let pattern2 = 'src="([^"]+\.pdf)"'
    let links2 = ($html | grep -oP $pattern2 | each { |m| $m | str replace -r 'src="|"' '' })
    
    # Combine and dedupe
    ($links | concat $links2 | uniq)
}

# Extract all links from HTML
def extractLinks [html: string, baseUrl: string] {
    let pattern = 'href="([^"]+)"'
    let links = ($html | grep -oP $pattern | each { |m| 
        let link = ($m | str replace -r 'href="|"' '')
        # Handle relative URLs
        if ($link | str starts-with "/") {
            $baseUrl + $link
        } else if not ($link | str starts-with "http") {
            $baseUrl + "/" + $link
        } else {
            $link
        }
    })
    
    $links | uniq
}

# Filter educational links
def filterEducational [links: list<string>] {
    $links | where { |l|
        let lower = ($l | str downcase)
        $lower | str contains "bac" or 
        $lower | str contains "math" or 
        $lower | str contains "physique" or 
        $lower | str contains "examen" or 
        $lower | str contains "exercice" or
        $lower | str contains "terminal"
    }
}

# Download PDF file
def downloadPdf [url: string, outputDir: string] {
    let filename = ($url | path basename)
    let safeName = ($filename | str replace -r '[^\w\-.]' '_')
    let outputPath = $"($outputDir)/($safeName)"
    
    print $"    Downloading: ($filename)"
    
    let result = (^curl -s -L -o $outputPath $url 2>&1)
    
    if ($result | str length) == 0 {
        let size = (ls -l $outputPath | get 0.size)
        print $"    ✓ Saved: ($filename) (($size) bytes)"
    } else {
        print $"    ! Download failed"
    }
}

# Save HTML content
def saveHtml [url: string, html: string] {
    let filename = ($url | path basename | default "index" | str replace -r '[^\w\-\.]' '_')
    let outputPath = $"resources/web/html/($filename).html"
    
    $html | save -f $outputPath
    print $"    ✓ Saved HTML: ($outputPath)"
}

# Save link report
def saveLinks [url: string, links: list<string>] {
    let filename = ($url | path basename | default "links" | str replace -r '[^\w\-\.]' '_')
    let outputPath = $"resources/web/links/($filename).txt"
    
    $links | save -f $outputPath
    print $"    ✓ Saved ($links | length) links"
}

# Main crawl function
export def crawl [...urls: string] {
    print "=========================================="
    print "Web Crawler - Educational Content Finder"
    print "=========================================="
    print ""
    
    # Default URLs if none provided
    let targetUrls = if ($urls | is-empty) {
        print "No URLs provided, using defaults..."
        [
            "https://www.education.gouv.mr",
            "https://fr.wikipedia.org/wiki/Baccalauréat"
        ]
    } else {
        $urls
    }
    
    print $"Crawling ($targetUrls | length) sites..."
    print ""
    
    let allPdfs = []
    
    for url in $targetUrls {
        print $"--- ($url) ---"
        
        let html = fetchPage $url
        if ($html | str length) > 0 {
            # Save HTML
            saveHtml $url $html
            
            # Extract and save all links
            let links = extractLinks $html $url
            saveLinks $url $links
            
            # Filter for educational content
            let eduLinks = filterEducational $links
            if not ($eduLinks | is-empty) {
                print $"  Found ($eduLinks | length) educational links"
            }
            
            # Extract PDF links
            let pdfs = extractPdfLinks $html
            if not ($pdfs | is-empty) {
                print $"  Found ($pdfs | length) PDF links"
                for pdf in $pdfs {
                    # Make absolute URL if needed
                    let fullUrl = if ($pdf | str starts-with "http") {
                        $pdf
                    } else {
                        $url + "/" + $pdf
                    }
                    downloadPdf $fullUrl "resources/web/pdf"
                }
            }
        }
        
        print ""
    }
    
    print "=========================================="
    print "Crawl Complete!"
    print "=========================================="
    print $""
    print "Saved to:"
    print "  - HTML: resources/web/html/"
    print "  - PDFs: resources/web/pdf/"
    print "  - Links: resources/web/links/"
}

# Quick link checker
export def check [url: string] {
    print $"Checking: ($url)"
    print ""
    
    let html = fetchPage $url
    if ($html | str length) > 0 {
        let links = extractLinks $html $url
        let pdfs = extractPdfLinks $html
        let eduLinks = filterEducational $links
        
        print $"Total links: ($links | length)"
        print $"PDF links: ($pdfs | length)"
        print $"Educational links: ($eduLinks | length)"
        print $""
        
        if not ($eduLinks | is-empty) {
            print "Educational links found:"
            for link in ($eduLinks | str join "\n") {
                print $"  - ($link)"
            }
        }
    }
}
