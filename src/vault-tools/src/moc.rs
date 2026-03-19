//! Map of Content (MOC) operations

use crate::models::{MocEntry, MocUpdateRequest, MocUpdateResult};
use anyhow::{Context, Result};
use std::fs;
use std::path::PathBuf;

use super::read::{get_vault_path, resolve_note_path};

/// Update a Map of Content file
pub async fn update_moc(request: MocUpdateRequest) -> MocUpdateResult {
    let vault_path = get_vault_path();
    let moc_path = find_or_create_moc_path(&vault_path, &request.subject);

    match update_moc_content(&moc_path, &request.subject, request.entries) {
        Ok(()) => {
            let relative_path = moc_path
                .strip_prefix(&vault_path)
                .map(|p| p.to_string_lossy().to_string())
                .unwrap_or_else(|_| request.subject.clone());
            MocUpdateResult::success(relative_path)
        }
        Err(e) => MocUpdateResult::error(format!("Failed to update MOC: {}", e)),
    }
}

/// Find or create MOC file path for a subject
fn find_or_create_moc_path(vault_path: &PathBuf, subject: &str) -> PathBuf {
    // Try common MOC naming patterns
    let patterns = [
        format!("MOC - {}.md", subject),
        format!("MOC - {}/index.md", subject),
        format!("{}/index.md", subject),
        format!("{}.md", subject),
    ];

    for pattern in &patterns {
        let path = vault_path.join(pattern);
        if path.exists() {
            return path;
        }
    }

    // Default: MOC - Subject.md
    vault_path.join(format!("MOC - {}.md", subject))
}

/// Update MOC file content
fn update_moc_content(path: &PathBuf, subject: &str, entries: Option<Vec<MocEntry>>) -> Result<()> {
    let content = build_moc_content(subject, entries);
    
    // Ensure parent directory exists
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent).context("Creating MOC directory")?;
    }

    fs::write(path, content).context("Writing MOC file")?;
    Ok(())
}

/// Build MOC markdown content
fn build_moc_content(subject: &str, entries: Option<Vec<MocEntry>>) -> String {
    let mut content = String::new();

    // Frontmatter
    content.push_str("---\n");
    content.push_str("type: map-of-content\n");
    content.push_str("created: ");
    content.push_str(&std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .map(|d| d.as_secs().to_string())
        .unwrap_or_default());
    content.push_str("\n---\n\n");

    // Title
    content.push_str("# ");
    content.push_str(subject);
    content.push_str("\n\n");

    // Entries
    if let Some(items) = entries {
        for entry in items {
            content.push_str("## ");
            content.push_str(&entry.title);
            content.push('\n');
            
            content.push_str("![[");
            content.push_str(&entry.path);
            content.push_str("]]\n");
            
            if let Some(desc) = entry.description {
                content.push('\n');
                content.push_str(&desc);
                content.push('\n');
            }
            
            content.push('\n');
        }
    } else {
        content.push_str("_No entries yet. Populate this MOC with linked notes._\n");
    }

    content
}

/// Add entry to existing MOC
pub async fn add_to_moc(subject: &str, entry: MocEntry) -> MocUpdateResult {
    let vault_path = get_vault_path();
    let moc_path = find_or_create_moc_path(&vault_path, subject);

    match fs::read_to_string(&moc_path) {
        Ok(content) => {
            let new_entry = format!(
                "\n## {}\n![[{}]]\n",
                entry.title,
                entry.path
            );
            let updated = format!("{}\n{}", content.trim_end(), new_entry);
            
            match fs::write(&moc_path, updated) {
                Ok(()) => {
                    let relative_path = moc_path
                        .strip_prefix(&vault_path)
                        .map(|p| p.to_string_lossy().to_string())
                        .unwrap_or_else(|_| subject.to_string());
                    MocUpdateResult::success(relative_path)
                }
                Err(e) => MocUpdateResult::error(format!("Failed to write MOC: {}", e)),
            }
        }
        Err(_) => {
            // MOC doesn't exist, create it
            update_moc(MocUpdateRequest {
                subject: subject.to_string(),
                entries: Some(vec![entry]),
            }).await
        }
    }
}

/// Get entries from MOC file
pub fn get_moc_entries(path: &str) -> Result<Vec<MocEntry>> {
    let vault_path = get_vault_path();
    let full_path = resolve_note_path(&vault_path, path);
    
    let content = fs::read_to_string(&full_path).context("Reading MOC file")?;
    
    // Extract entries from markdown
    let mut entries = Vec::new();
    let mut current_title: Option<String> = None;
    let mut current_path: Option<String> = None;
    
    for line in content.lines() {
        let trimmed = line.trim();
        
        if trimmed.starts_with("## ") {
            if let (Some(title), Some(note_path)) = (current_title.take(), current_path.take()) {
                entries.push(MocEntry {
                    title,
                    path: note_path,
                    description: None,
                });
            }
            current_title = Some(trimmed[3..].to_string());
        } else if trimmed.starts_with("![[") && trimmed.ends_with("]]") {
            let note_path = &trimmed[3..trimmed.len() - 2];
            current_path = Some(note_path.to_string());
        }
    }
    
    // Don't forget the last entry
    if let (Some(title), Some(note_path)) = (current_title, current_path) {
        entries.push(MocEntry {
            title,
            path: note_path,
            description: None,
        });
    }
    
    Ok(entries)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_build_moc_content() {
        let entries = vec![
            MocEntry {
                title: "Getting Started".to_string(),
                path: "getting-started.md".to_string(),
                description: Some("Start here".to_string()),
            },
        ];
        
        let content = build_moc_content("Test MOC", Some(entries));
        assert!(content.contains("# Test MOC"));
        assert!(content.contains("## Getting Started"));
        assert!(content.contains("![[getting-started.md]]"));
    }

    #[test]
    fn test_build_moc_content_empty() {
        let content = build_moc_content("Empty MOC", None);
        assert!(content.contains("# Empty MOC"));
        assert!(content.contains("_No entries yet"));
    }
}
