# Web Crawler Daemon - Scrape Mauritanian education sites
# Run: nu daemons/webcrawler-daemon.nu

let resources-dir = "resources/web"

# Ensure directories exist
mkdir $resources-dir

# Target sites
const TARGET_SITES = [
    "https://www.minedu.mr",
    "https://www.bacmr.net",
    "https://www.mesrs.mr",
]

print "🕷️ Web Crawler Daemon Starting..."
print $"Found (($TARGET_SITES | length)) sites to crawl"
print ""

# Crawl function
export def crawl-site [url: string] {
    print $"Crawling: ($url)"
    
    # Use curl to fetch and grep for PDF links
    let cmd = $"curl -s ($url) | grep -oP 'https?://[^']+\\.pdf' | head -20"
    
    print $"  Found PDF links would be extracted here"
    print ""
}

# Main loop
for site in $TARGET_SITES {
    crawl-site $site
}

print "🕷️ Web Crawler Daemon Complete!"
