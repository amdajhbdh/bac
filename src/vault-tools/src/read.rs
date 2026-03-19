//! Read operations for vault notes

use crate::models::{NoteContent, ReadResult};
use anyhow::{Context, Result};
use std::fs;
use std::path::{Path, PathBuf};

/// Get the configured vault path
pub fn get_vault_path() -> PathBuf {
    std::env::var("VAULT_PATH")
        .or_else(|_| std::env::var("OBSIDIAN_VAULT_PATH"))
        .map(PathBuf::from)
        .unwrap_or_else(|_| {
            dirs::document_dir()
                .unwrap_or_else(|| PathBuf::from("."))
                .join("Obsidian")
        })
}

/// Resolve a note path relative to vault
pub fn resolve_note_path(vault_path: &Path, path: &str) -> PathBuf {
    let path = path.trim();
    if path.starts_with('/') || path.contains(':') {
        // Absolute path
        PathBuf::from(path)
    } else {
        // Relative to vault
        vault_path.join(path)
    }
    .with_extension("md")
}

/// Read a note from the vault
pub async fn read_note(path: &str) -> ReadResult {
    let vault_path = get_vault_path();
    let full_path = resolve_note_path(&vault_path, path);

    match read_note_from_path(&full_path) {
        Ok(content) => {
            let relative_path = full_path
                .strip_prefix(&vault_path)
                .map(|p| p.to_string_lossy().to_string())
                .unwrap_or_else(|_| path.to_string());

            let title = extract_title(&content, &relative_path);
            let (frontmatter, body) = crate::models::parse_frontmatter(&content)
                .unwrap_or_else(|| (crate::models::Frontmatter::new(), content.as_str()));

            ReadResult::success(NoteContent::new(
                relative_path,
                title,
                body.trim().to_string(),
                frontmatter,
            ))
        }
        Err(e) => {
            let error_msg = e.to_string();
            if error_msg.contains("No such file") || error_msg.contains("not found") {
                ReadResult::not_found(path.to_string())
            } else {
                ReadResult::error(format!("Failed to read note: {}", e))
            }
        }
    }
}

/// Read note content from a specific path
fn read_note_from_path(path: &Path) -> Result<String> {
    fs::read_to_string(path).with_context(|| format!("Reading note: {}", path.display()))
}

/// Extract title from frontmatter or filename
fn extract_title(content: &str, path: &str) -> String {
    // Try frontmatter first
    if let Some((fm, _)) = crate::models::parse_frontmatter(content) {
        if let Some(title) = fm.get("title") {
            if let Some(s) = title.as_str() {
                return s.to_string();
            }
        }
    }

    // Fall back to filename
    Path::new(path)
        .file_stem()
        .map(|s| s.to_string_lossy().to_string())
        .unwrap_or_else(|| path.to_string())
}

/// Check if a note exists
pub fn note_exists(path: &str) -> bool {
    let vault_path = get_vault_path();
    let full_path = resolve_note_path(&vault_path, path);
    full_path.exists()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_extract_title_from_filename() {
        let title = extract_title("", "test-note.md");
        assert_eq!(title, "test-note");
    }

    #[test]
    fn test_extract_title_from_frontmatter() {
        let content = r#"---
title: Custom Title
---
Body content"#;
        let title = extract_title(content, "test.md");
        assert_eq!(title, "Custom Title");
    }
}
