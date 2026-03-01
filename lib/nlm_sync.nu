# BAC Unified - NotebookLM Sync Module
# Automatically syncs local PDFs to NotebookLM notebooks

use std log

# Configuration
const SYNC_DB = "data/nlm_sync.json"
const PDF_DIR = "db/pdfs"

# Initialize sync database if it doesn't exist
def init_db [] {
    if not (SYNC_DB | path exists) {
        mkdir (SYNC_DB | path dirname)
        "{ "synced_files": [], "notebooks": {} }" | save SYNC_DB
    }
}

# Get synced files list
def get_synced_files [] {
    init_db
    open SYNC_DB | get synced_files
}

# Add file to synced list
def mark_as_synced [file: string] {
    let db = (open SYNC_DB)
    let new_synced = ($db.synced_files | append $file | uniq)
    $db | upsert synced_files $new_synced | save -f SYNC_DB
}

# Categorize PDF based on filename
export def categorize_pdf [filename: string] {
    let fn = ($filename | str downcase)
    
    if ($fn =~ 'math' or $fn =~ 'matrice' or $fn =~ 'arithmétique' or $fn =~ 'intégrale' or $fn =~ 'complexe' or $fn =~ 'suite' or $fn =~ '7c') {
        return "Math"
    }
    
    if ($fn =~ 'pc' or $fn =~ 'physique' or $fn =~ 'chimie' or $fn =~ 'mécanique' or $fn =~ 'dynamique' or $fn =~ 'projectile' or $fn =~ 'ressort' or $fn =~ 'électrique' or $fn =~ 'induction' or $fn =~ 'cinétique') {
        return "PC"
    }
    
    if ($fn =~ 'sn' or $fn =~ 'science' or $fn =~ 'biologie' or $fn =~ 'hormone' or $fn =~ 'génétique' or $fn =~ 'système nerveux' or $fn =~ '7d') {
        return "SN"
    }
    
    return "General"
}

# Get or create notebook for a subject
export def get_notebook_for_subject [subject: string] {
    # This function would ideally use notebook_list and notebook_create
    # For now, we'll return a placeholder or look it up in the sync DB
    let db = (open SYNC_DB)
    if ($db.notebooks | has $subject) {
        return ($db.notebooks | get $subject)
    }
    
    # If not in DB, we'll need to create it (logic handled in calling script or later)
    return null
}

# Set notebook ID for a subject
export def set_notebook_id [subject: string, id: string] {
    let db = (open SYNC_DB)
    let new_notebooks = ($db.notebooks | upsert $subject $id)
    $db | upsert notebooks $new_notebooks | save -f SYNC_DB
}

# Sync all new PDFs
export def sync_pdfs [] {
    init_db
    let synced = (get_synced_files)
    let pdfs = (ls PDF_DIR | where { |f| $f.name | str ends-with ".pdf" })
    
    print $"Found ($pdfs | length) PDFs"
    
    for pdf in $pdfs {
        let name = ($pdf.name | path basename)
        if ($name in $synced) {
            # print $"Skipping already synced: ($name)"
            continue
        }
        
        let subject = (categorize_pdf $name)
        print $"New PDF: ($name) -> Subject: ($subject)"
        
        # Actual sync logic will be called by the daemon using MCP tools
        # For now, we just identify them
    }
}
