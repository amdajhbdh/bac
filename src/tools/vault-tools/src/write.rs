//! Write operations for vault notes

use crate::models::{Frontmatter, WriteRequest, WriteResult};
use anyhow::{Context, Result};
use std::fs;
use std::path::PathBuf;

use super::read::{get_vault_path, resolve_note_path};

/// Write a note to the vault
pub async fn write_note(request: WriteRequest) -> WriteResult {
    let vault_path = get_vault_path();
    let full_path = resolve_note_path(&vault_path, &request.path);

    // Check if file exists to determine if this is a create or update
    let created = !full_path.exists();

    match write_note_to_path(&full_path, &request.content, request.frontmatter.as_ref()) {
        Ok(()) => {
            let relative_path = full_path
                .strip_prefix(&vault_path)
                .map(|p| p.to_string_lossy().to_string())
                .unwrap_or_else(|_| request.path.clone());

            WriteResult::success(relative_path, created)
        }
        Err(e) => WriteResult::error(format!("Failed to write note: {}", e)),
    }
}

/// Write note content to a specific path
fn write_note_to_path(path: &PathBuf, content: &str, frontmatter: Option<&Frontmatter>) -> Result<()> {
    // Ensure parent directory exists
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent).context("Creating parent directories")?;
    }

    // Build the full markdown content
    let markdown = build_markdown(content, frontmatter);

    fs::write(path, markdown).context("Writing note file")?;
    Ok(())
}

/// Build markdown with frontmatter
fn build_markdown(content: &str, frontmatter: Option<&Frontmatter>) -> String {
    match frontmatter {
        Some(fm) if !fm.is_empty() => {
            let fm_str = crate::models::serialize_frontmatter(fm);
            if fm_str.is_empty() {
                content.to_string()
            } else {
                format!("{}\n\n{}", fm_str, content)
            }
        }
        _ => content.to_string(),
    }
}

/// Update frontmatter of an existing note
pub async fn update_frontmatter(path: &str, updates: Frontmatter) -> WriteResult {
    let vault_path = get_vault_path();
    let full_path = resolve_note_path(&vault_path, path);

    match fs::read_to_string(&full_path) {
        Ok(content) => {
            let (existing_fm, body) = crate::models::parse_frontmatter(&content)
                .unwrap_or_else(|| (Frontmatter::new(), content.as_str()));

            // Merge updates with existing frontmatter
            let mut merged = existing_fm;
            for (k, v) in updates {
                merged.insert(k, v);
            }

            match write_note_to_path(&full_path, body.trim(), Some(&merged)) {
                Ok(()) => {
                    let relative_path = full_path
                        .strip_prefix(&vault_path)
                        .map(|p| p.to_string_lossy().to_string())
                        .unwrap_or_else(|_| path.to_string());
                    WriteResult::success(relative_path, false)
                }
                Err(e) => WriteResult::error(format!("Failed to write note: {}", e)),
            }
        }
        Err(e) => {
            let error_msg = e.to_string();
            if error_msg.contains("No such file") || error_msg.contains("not found") {
                WriteResult::error(format!("Note not found: {}", path))
            } else {
                WriteResult::error(format!("Failed to read note: {}", e))
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_build_markdown_with_frontmatter() {
        let mut fm = Frontmatter::new();
        fm.insert("title".to_string(), serde_json::json!("Test"));
        fm.insert("tags".to_string(), serde_json::json!(["rust", "test"]));

        let markdown = build_markdown("Hello world", Some(&fm));
        assert!(markdown.starts_with("---"));
        assert!(markdown.contains("title: Test"));
        assert!(markdown.contains("Hello world"));
    }

    #[test]
    fn test_build_markdown_without_frontmatter() {
        let markdown = build_markdown("Hello world", None);
        assert_eq!(markdown, "Hello world");
    }
}
