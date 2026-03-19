//! Wikilink handling for Obsidian notes

use crate::models::{LinkRequest, LinkResult, Wikilink};
use anyhow::{Context, Result};
use regex::Regex;
use std::fs;
use std::path::PathBuf;

use super::read::{get_vault_path, resolve_note_path};

/// Create a wikilink between two notes
pub async fn create_link(request: LinkRequest) -> LinkResult {
    let vault_path = get_vault_path();
    let source_path = resolve_note_path(&vault_path, &request.source);

    if !source_path.exists() {
        return LinkResult::error(format!("Source note not found: {}", request.source));
    }

    // Build the wikilink
    let link_text = build_wikilink(&request.target, &request.link_type);

    match append_link_to_note(&source_path, &link_text) {
        Ok(()) => LinkResult::success(request.source, request.target),
        Err(e) => LinkResult::error(format!("Failed to create link: {}", e)),
    }
}

/// Build wikilink text from target and type
fn build_wikilink(target: &str, link_type: &Option<String>) -> String {
    match link_type.as_deref() {
        Some("alias") | Some("display") => {
            let display = Path::new(target)
                .file_stem()
                .map(|s| s.to_string_lossy().to_string())
                .unwrap_or_else(|| target.to_string());
            format!("[[{}|{}]]", target, display)
        }
        _ => format!("[[{}]]", target),
    }
}

/// Append a wikilink to a note
fn append_link_to_note(path: &PathBuf, link: &str) -> Result<()> {
    let content = fs::read_to_string(path).context("Reading source note")?;
    let new_content = format!("{}\n{}", content.trim_end(), link);
    fs::write(path, new_content).context("Writing updated note")?;
    Ok(())
}

/// Insert a wikilink at a specific position in content
pub fn insert_wikilink(content: &str, target: &str, position: usize, display: Option<&str>) -> String {
    let link = match display {
        Some(d) => format!("[[{}|{}]]", target, d),
        None => format!("[[{}]]", target),
    };

    let mut result = content.to_string();
    result.insert_str(position, &link);
    result
}

/// Parse all wikilinks from content
pub fn parse_wikilinks(content: &str) -> Vec<Wikilink> {
    crate::models::extract_wikilinks(content)
}

/// Get unique target notes from wikilinks
pub fn get_linked_notes(content: &str) -> Vec<String> {
    parse_wikilinks(content)
        .into_iter()
        .map(|w| w.target)
        .collect::<std::collections::HashSet<_>>()
        .into_iter()
        .collect()
}

/// Rename all references to a note in content
pub fn rename_links(content: &str, old_target: &str, new_target: &str) -> String {
    let regex_pattern = format!(r"\[\[{}\s*(?:\|[^\]]+)?\]\]", regex::escape(old_target));
    let re = Regex::new(&regex_pattern).unwrap();
    
    re.replace_all(content, |caps: &regex::Captures| {
        if let Some(display) = caps.get(0).and_then(|m| {
            m.as_str()
                .split('|')
                .nth(1)
                .map(|s| s.trim_end_matches(']'))
        }) {
            format!("[[{}|{}]]", new_target, display)
        } else {
            format!("[[{}]]", new_target)
        }
    }).to_string()
}

use std::path::Path;

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_build_wikilink_simple() {
        let link = build_wikilink("Test Note", &None);
        assert_eq!(link, "[[Test Note]]");
    }

    #[test]
    fn test_build_wikilink_with_display() {
        let link = build_wikilink("Test Note", &Some("display".to_string()));
        assert!(link.contains("|"));
        assert!(link.contains("Test Note"));
    }

    #[test]
    fn test_parse_wikilinks() {
        let content = "See [[Note A]] and [[Note B|Display Text]] for details.";
        let links = parse_wikilinks(content);
        
        assert_eq!(links.len(), 2);
        assert_eq!(links[0].target, "Note A");
        assert_eq!(links[0].display, None);
        assert_eq!(links[1].target, "Note B");
        assert_eq!(links[1].display, Some("Display Text".to_string()));
    }

    #[test]
    fn test_rename_links() {
        let content = "See [[Old Note]] and [[Old Note|alias]] for info.";
        let result = rename_links(content, "Old Note", "New Note");
        
        assert!(result.contains("[[New Note]]"));
        assert!(result.contains("[[New Note|alias]]"));
    }
}
