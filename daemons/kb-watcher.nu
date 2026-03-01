#!/usr/bin/env nu
# BAC Unified - Knowledge Base Watcher Daemon
# Automatically indexes new files in db/pdfs/ using Vector KB

use std log

# Config
const PDF_DIR = "db/pdfs"
const INDEXER = "scripts/process_kb.py"

# Watch for new files
export def main [] {
    print "=========================================="
    print "🚀 Knowledge Base Watcher Started"
    print $"📁 Monitoring: (PDF_DIR)"
    print "=========================================="
    
    # Watch for changes using fswatch (if available)
    # Fallback to polling for maximum compatibility
    loop {
        let files = (ls PDF_DIR | where { |f| $f.name | str ends-with ".pdf" })
        
        # We'll use the sync DB we created earlier to track indexed files
        # Check against nlm_sync.json (synced_files)
        let sync_db = (open data/nlm_sync.json)
        let indexed = $sync_db.synced_files
        
        for file in $files {
            let name = ($file.name | path basename)
            if ($name in $indexed) {
                continue
            }
            
            print $"🔔 New resource detected: ($name)"
            
            # 1. Extract metadata and index in Vector DB
            print "  → Extracting metadata and vectorizing..."
            let result = (try {
                ^python scripts/process_kb.py $file.name
                true
            } catch {
                false
            })
            
            if $result {
                print "  ✓ Vectorized successfully"
                # 2. Add to NotebookLM sync list
                let new_db = ($sync_db | upsert synced_files ($indexed | append $name | uniq))
                $new_db | save -f data/nlm_sync.json
            } else {
                print "  ✗ Error processing resource"
            }
        }
        
        # Poll every 30 seconds
        sleep 30sec
    }
}
