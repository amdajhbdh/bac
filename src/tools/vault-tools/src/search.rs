//! Search operations for vault notes

use crate::models::{SearchResult, SearchResults};
use std::fs;
use walkdir::WalkDir;

use super::read::get_vault_path;

/// Search vault for notes matching query
pub async fn search_vault(query: &str) -> SearchResults {
    let vault_path = get_vault_path();
    let query_lower = query.to_lowercase();

    let mut results: Vec<SearchResult> = Vec::new();

    for entry in WalkDir::new(&vault_path)
        .into_iter()
        .filter_map(|e| e.ok())
        .filter(|e| e.path().extension().map_or(false, |ext| ext == "md"))
    {
        let path = entry.path();
        
        if let Ok(content) = fs::read_to_string(path) {
            let title = extract_title_from_content(&content, path);
            
            // Calculate match score and create snippet
            if let Some((score, snippet)) = calculate_match(&content, &title, &query_lower) {
                let relative_path = path
                    .strip_prefix(&vault_path)
                    .map(|p| p.to_string_lossy().to_string())
                    .unwrap_or_else(|_| path.to_string_lossy().to_string());

                results.push(SearchResult {
                    path: relative_path,
                    title,
                    snippet,
                    score,
                });
            }
        }
    }

    // Sort by score descending
    results.sort_by(|a, b| b.score.partial_cmp(&a.score).unwrap_or(std::cmp::Ordering::Equal));
    
    // Limit results
    results.truncate(50);

    SearchResults::success(query.to_string(), results)
}

/// Extract title from note content
fn extract_title_from_content(content: &str, path: &std::path::Path) -> String {
    // Try frontmatter title first
    if let Some((fm, _)) = crate::models::parse_frontmatter(content) {
        if let Some(title) = fm.get("title") {
            if let Some(s) = title.as_str() {
                return s.to_string();
            }
        }
    }

    // Fall back to first H1 heading
    for line in content.lines() {
        let trimmed = line.trim();
        if trimmed.starts_with("# ") {
            return trimmed[2..].to_string();
        }
        if trimmed.starts_with("---") {
            break; // Reached frontmatter end without title
        }
    }

    // Fall back to filename
    path.file_stem()
        .map(|s| s.to_string_lossy().to_string())
        .unwrap_or_else(|| "Untitled".to_string())
}

/// Calculate match score and extract snippet
fn calculate_match(content: &str, title: &str, query: &str) -> Option<(f32, String)> {
    let content_lower = content.to_lowercase();
    
    // Check if query appears in content
    if !content_lower.contains(query) {
        return None;
    }

    let mut score = 0.0;

    // Title match is most important
    if title.to_lowercase().contains(query) {
        score += 10.0;
    }

    // Count occurrences
    let occurrences = content_lower.matches(query).count();
    score += (occurrences as f32) * 0.5;

    // H1 heading match bonus
    let capitalized_query = {
        let mut chars = query.chars();
        match chars.next() {
            None => String::new(),
            Some(first) => first.to_uppercase().to_string() + chars.as_str(),
        }
    };
    if content.starts_with(&format!("# {}", query)) || 
       content.starts_with(&format!("# {}", capitalized_query)) {
        score += 5.0;
    }

    // Frontmatter match bonus
    if let Some((fm, _)) = crate::models::parse_frontmatter(content) {
        for (_, v) in fm {
            if let Some(s) = v.as_str() {
                if s.to_lowercase().contains(query) {
                    score += 3.0;
                }
            }
        }
    }

    // Generate snippet
    let snippet = extract_snippet(content, query, 150);

    Some((score, snippet))
}

/// Extract a snippet around the first match
fn extract_snippet(content: &str, query: &str, max_len: usize) -> String {
    let content_lower = content.to_lowercase();
    let query_lower = query.to_lowercase();
    
    if let Some(pos) = content_lower.find(&query_lower) {
        let start = pos.saturating_sub(max_len / 2);
        let end = (pos + query.len() + max_len / 2).min(content.len());
        
        let mut snippet = content[start..end].to_string();
        
        if start > 0 {
            snippet = format!("...{}", snippet);
        }
        if end < content.len() {
            snippet = format!("{}...", snippet);
        }
        
        // Clean up newlines for snippet
        snippet = snippet
            .lines()
            .take(3)
            .collect::<Vec<_>>()
            .join(" ")
            .trim()
            .to_string();
        
        snippet
    } else {
        // Return first part of content
        content
            .lines()
            .take(3)
            .collect::<Vec<_>>()
            .join(" ")
            .trim()
            .chars()
            .take(max_len)
            .collect::<String>()
    }
}

/// Search by tag
pub async fn search_by_tag(tag: &str) -> SearchResults {
    let vault_path = get_vault_path();
    // NOTE: query format ready for future content search: let _query = format!("#{}", tag.trim_start_matches('#'));

    let mut results: Vec<SearchResult> = Vec::new();

    for entry in WalkDir::new(&vault_path)
        .into_iter()
        .filter_map(|e| e.ok())
        .filter(|e| e.path().extension().map_or(false, |ext| ext == "md"))
    {
        let path = entry.path();
        
        if let Ok(content) = fs::read_to_string(path) {
            let tags = crate::models::extract_tags(&content);
            
            if tags.iter().any(|t| t == tag || t == tag.trim_start_matches('#')) {
                let title = extract_title_from_content(&content, path);
                let snippet = extract_snippet(&content, tag, 150);
                
                let relative_path = path
                    .strip_prefix(&vault_path)
                    .map(|p| p.to_string_lossy().to_string())
                    .unwrap_or_else(|_| path.to_string_lossy().to_string());

                results.push(SearchResult {
                    path: relative_path,
                    title,
                    snippet,
                    score: 1.0,
                });
            }
        }
    }

    SearchResults::success(tag.to_string(), results)
}

/// Get all unique tags in vault
pub async fn get_all_tags() -> Vec<(String, usize)> {
    let vault_path = get_vault_path();
    let mut tag_counts: std::collections::HashMap<String, usize> = std::collections::HashMap::new();

    for entry in WalkDir::new(&vault_path)
        .into_iter()
        .filter_map(|e| e.ok())
        .filter(|e| e.path().extension().map_or(false, |ext| ext == "md"))
    {
        if let Ok(content) = fs::read_to_string(entry.path()) {
            let tags = crate::models::extract_tags(&content);
            for tag in tags {
                *tag_counts.entry(tag).or_insert(0) += 1;
            }
        }
    }

    let mut tags: Vec<_> = tag_counts.into_iter().collect();
    tags.sort_by(|a, b| b.1.cmp(&a.1));
    tags
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_extract_snippet() {
        let content = "This is a test document with some content. The query appears here.";
        let snippet = extract_snippet(content, "query", 50);
        assert!(snippet.contains("query"));
    }

    #[test]
    fn test_calculate_match() {
        let content = "# Test Title\n\nThis is test content with the keyword.";
        let (score, _) = calculate_match(content, "Test Title", "keyword").unwrap();
        assert!(score > 0.0);
    }

    #[test]
    fn test_calculate_match_no_match() {
        let content = "# Test Title\n\nThis is test content.";
        let result = calculate_match(content, "Test Title", "nonexistent");
        assert!(result.is_none());
    }
}
