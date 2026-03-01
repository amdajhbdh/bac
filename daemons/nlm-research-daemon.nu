#!/usr/bin/env nu
# NLM Research Daemon - MCP Powered
# Automates NotebookLM research using MCP tools

use std log
use lib/nlm_sync.nu

# Main research function
export def research [...topics: string] {
    print "=========================================="
    print "NotebookLM Research Automation (MCP)"
    print "=========================================="
    print ""
    
    # Default topics if none provided
    let researchTopics = if ($topics | is-empty) {
        print "No topics provided, using defaults..."
        [
            {topic: "BAC Math Terminale C", subject: "Math", description: "Mathematics curriculum for BAC"},
            {topic: "BAC Physique Terminale C", subject: "PC", description: "Physics curriculum for BAC"},
            {topic: "BAC Sciences Naturelles Terminale D", subject: "SN", description: "Science curriculum for BAC"}
        ]
    } else {
        $topics | each { |t| {topic: $t, subject: (nlm_sync categorize_pdf $t), description: $"Research on ($t)"} }
    }
    
    for item in $researchTopics {
        print $""
        print $"=== Researching: ($item.topic) ==="
        
        # Check if notebook exists
        # In a real tool call, we would use notebook_list here
        # But we'll handle creation in the main execution if needed
        print $"Target subject: ($item.subject)"
        
        # We'll use the MCP tools in the actual implementation phase
        # notebook_create(title=$item.topic)
        # source_add(notebook_id=$id, source_type="url", url=...)
    }
}

# Sync PDFs to their respective notebooks
export def sync_all_pdfs [] {
    print "=========================================="
    print "Syncing local PDFs to NotebookLM"
    print "=========================================="
    
    # Identify new PDFs
    let pdfs = (ls db/pdfs | where { |f| $f.name | str ends-with ".pdf" })
    let synced = (nlm_sync get_synced_files)
    
    for pdf in $pdfs {
        let name = ($pdf.name | path basename)
        if ($name in $synced) {
            continue
        }
        
        let subject = (nlm_sync categorize_pdf $name)
        let notebook_id = (nlm_sync get_notebook_for_subject $subject)
        
        if $notebook_id != null {
            print $"Syncing: ($name) to ($subject) notebook..."
            # source_add(notebook_id=$notebook_id, source_type="file", file_path=$pdf.name, wait=true)
            # nlm_sync mark_as_synced $name
        } else {
            print $"Warning: No notebook ID for ($subject). Create it first!"
        }
    }
}

# Create a Study Pack (Full suite of artifacts)
export def study_pack [notebook_id: string, topic: string] {
    print $"Generating full Study Pack for: ($topic)"
    
    # studio_create(notebook_id=$notebook_id, artifact_type="audio", audio_format="deep_dive", confirm=true)
    # studio_create(notebook_id=$notebook_id, artifact_type="infographic", detail_level="detailed", confirm=true)
    # studio_create(notebook_id=$notebook_id, artifact_type="mind_map", title=$topic, confirm=true)
    # studio_create(notebook_id=$notebook_id, artifact_type="quiz", question_count=10, difficulty="medium", confirm=true)
    # studio_create(notebook_id=$notebook_id, artifact_type="flashcards", difficulty="hard", confirm=true)
    
    print "Study pack generation started. Use 'studio_status' to track progress."
}
