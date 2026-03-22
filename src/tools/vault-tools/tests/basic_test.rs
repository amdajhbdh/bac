//! Basic tests for vault-tools

use vault_tools::models::*;

/// Test frontmatter parsing
#[test]
fn test_parse_frontmatter() {
    let content = r#"---
title: Test Note
tags:
  - rust
  - testing
---
# Body content
"#;

    let result = parse_frontmatter(content);
    assert!(result.is_some());

    let (fm, body) = result.unwrap();
    assert_eq!(fm.get("title").and_then(|v| v.as_str()), Some("Test Note"));
    assert!(body.contains("Body content"));
}

/// Test frontmatter serialization
#[test]
fn test_serialize_frontmatter() {
    let mut fm = Frontmatter::new();
    fm.insert("title".to_string(), serde_json::json!("My Note"));
    fm.insert("count".to_string(), serde_json::json!(42));

    let yaml = serialize_frontmatter(&fm);
    assert!(yaml.contains("title: My Note"));
    assert!(yaml.contains("count: 42"));
}

/// Test wikilink extraction
#[test]
fn test_extract_wikilinks() {
    let content = "See [[Note A]] and [[Note B|display]] for details.";
    let links = extract_wikilinks(content);

    assert_eq!(links.len(), 2);
    assert_eq!(links[0].target, "Note A");
    assert_eq!(links[0].display, None);
    assert_eq!(links[1].target, "Note B");
    assert_eq!(links[1].display, Some("display".to_string()));
}

/// Test tag extraction
#[test]
fn test_extract_tags() {
    let content = "This is #rust and #testing content with #nested/tag too.";
    let tags = extract_tags(content);

    assert!(tags.contains(&"rust".to_string()));
    assert!(tags.contains(&"testing".to_string()));
    assert!(tags.contains(&"nested/tag".to_string()));
}

/// Test NoteContent creation
#[test]
fn test_note_content_creation() {
    let content = NoteContent::new(
        "test.md".to_string(),
        "Test".to_string(),
        "# Hello\n\nThis is [[Link]] content #tag".to_string(),
        Frontmatter::new(),
    );

    assert_eq!(content.title, "Test");
    assert!(content.links.len() == 1);
    assert!(content.tags.contains(&"tag".to_string()));
}

/// Test ReadResult success
#[test]
fn test_read_result_success() {
    let note = NoteContent::new(
        "test.md".to_string(),
        "Test".to_string(),
        "content".to_string(),
        Frontmatter::new(),
    );

    let result = ReadResult::success(note);
    assert!(result.success);
    assert!(result.note.is_some());
    assert!(result.error.is_none());
}

/// Test ReadResult not found
#[test]
fn test_read_result_not_found() {
    let result = ReadResult::not_found("missing.md".to_string());
    assert!(!result.success);
    assert!(result.note.is_none());
    assert!(result.error.unwrap().contains("missing.md"));
}

/// Test SearchResults
#[test]
fn test_search_results() {
    let results = vec![SearchResult {
        path: "note1.md".to_string(),
        title: "Note 1".to_string(),
        snippet: "snippet...".to_string(),
        score: 1.0,
    }];

    let search = SearchResults::success("query".to_string(), results);
    assert!(search.success);
    assert_eq!(search.total, 1);
}

/// Test Wikilink struct
#[test]
fn test_wikilink_struct() {
    let link = Wikilink {
        target: "Target Note".to_string(),
        display: Some("Display Text".to_string()),
    };

    assert_eq!(link.target, "Target Note");
    assert_eq!(link.display, Some("Display Text".to_string()));
}
