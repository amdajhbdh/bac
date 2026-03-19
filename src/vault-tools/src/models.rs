//! Data models for vault-tools

use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Note frontmatter (YAML metadata)
pub type Frontmatter = HashMap<String, serde_json::Value>;

/// Note content with frontmatter
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NoteContent {
    pub path: String,
    pub title: String,
    pub content: String,
    pub frontmatter: Frontmatter,
    pub links: Vec<Wikilink>,
    pub tags: Vec<String>,
}

impl NoteContent {
    pub fn new(path: String, title: String, content: String, frontmatter: Frontmatter) -> Self {
        let links = extract_wikilinks(&content);
        let tags = extract_tags(&content);
        Self {
            path,
            title,
            content,
            frontmatter,
            links,
            tags,
        }
    }
}

/// Wikilink reference
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Wikilink {
    pub target: String,
    pub display: Option<String>,
}

/// Read operation result
#[derive(Debug, Serialize, Deserialize)]
pub struct ReadResult {
    pub success: bool,
    pub note: Option<NoteContent>,
    pub error: Option<String>,
}

impl ReadResult {
    pub fn success(note: NoteContent) -> Self {
        Self {
            success: true,
            note: Some(note),
            error: None,
        }
    }

    pub fn not_found(path: String) -> Self {
        Self {
            success: false,
            note: None,
            error: Some(format!("Note not found: {}", path)),
        }
    }

    pub fn error(msg: String) -> Self {
        Self {
            success: false,
            note: None,
            error: Some(msg),
        }
    }
}

/// Write request
#[derive(Debug, Deserialize)]
pub struct WriteRequest {
    pub path: String,
    pub content: String,
    pub frontmatter: Option<Frontmatter>,
}

/// Write operation result
#[derive(Debug, Serialize, Deserialize)]
pub struct WriteResult {
    pub success: bool,
    pub path: Option<String>,
    pub created: bool,
    pub error: Option<String>,
}

impl WriteResult {
    pub fn success(path: String, created: bool) -> Self {
        Self {
            success: true,
            path: Some(path),
            created,
            error: None,
        }
    }

    pub fn error(msg: String) -> Self {
        Self {
            success: false,
            path: None,
            created: false,
            error: Some(msg),
        }
    }
}

/// Link request
#[derive(Debug, Deserialize)]
pub struct LinkRequest {
    pub source: String,
    pub target: String,
    pub link_type: Option<String>,
}

/// Link operation result
#[derive(Debug, Serialize, Deserialize)]
pub struct LinkResult {
    pub success: bool,
    pub source: Option<String>,
    pub target: Option<String>,
    pub error: Option<String>,
}

impl LinkResult {
    pub fn success(source: String, target: String) -> Self {
        Self {
            success: true,
            source: Some(source),
            target: Some(target),
            error: None,
        }
    }

    pub fn error(msg: String) -> Self {
        Self {
            success: false,
            source: None,
            target: None,
            error: Some(msg),
        }
    }
}

/// MOC update request
#[derive(Debug, Deserialize)]
pub struct MocUpdateRequest {
    pub subject: String,
    pub entries: Option<Vec<MocEntry>>,
}

/// MOC entry
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MocEntry {
    pub title: String,
    pub path: String,
    pub description: Option<String>,
}

/// MOC update result
#[derive(Debug, Serialize, Deserialize)]
pub struct MocUpdateResult {
    pub success: bool,
    pub path: Option<String>,
    pub error: Option<String>,
}

impl MocUpdateResult {
    pub fn success(path: String) -> Self {
        Self {
            success: true,
            path: Some(path),
            error: None,
        }
    }

    pub fn error(msg: String) -> Self {
        Self {
            success: false,
            path: None,
            error: Some(msg),
        }
    }
}

/// Search result item
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SearchResult {
    pub path: String,
    pub title: String,
    pub snippet: String,
    pub score: f32,
}

/// Search results
#[derive(Debug, Serialize, Deserialize)]
pub struct SearchResults {
    pub success: bool,
    pub query: String,
    pub results: Vec<SearchResult>,
    pub total: usize,
    pub error: Option<String>,
}

impl SearchResults {
    pub fn success(query: String, results: Vec<SearchResult>) -> Self {
        let total = results.len();
        Self {
            success: true,
            query,
            total,
            results,
            error: None,
        }
    }

    pub fn empty(query: String) -> Self {
        Self {
            success: true,
            query,
            total: 0,
            results: vec![],
            error: None,
        }
    }

    pub fn error(msg: String) -> Self {
        Self {
            success: false,
            query: String::new(),
            total: 0,
            results: vec![],
            error: Some(msg),
        }
    }
}

/// Health check response
#[derive(Debug, Serialize)]
pub struct HealthResponse {
    pub status: String,
    pub vault_path: String,
}

// ============================================================================
// Utility functions for parsing Obsidian markdown
// ============================================================================

use once_cell::sync::Lazy;
use regex::Regex;

static WIKILINK_REGEX: Lazy<Regex> =
    Lazy::new(|| Regex::new(r"\[\[([^\]|]+)(?:\|([^\]]+))?\]\]").unwrap());

static TAG_REGEX: Lazy<Regex> = Lazy::new(|| Regex::new(r"#([a-zA-Z][a-zA-Z0-9_/-]*)").unwrap());

/// Extract wikilinks from markdown content
pub fn extract_wikilinks(content: &str) -> Vec<Wikilink> {
    WIKILINK_REGEX
        .captures_iter(content)
        .map(|cap| Wikilink {
            target: cap[1].trim().to_string(),
            display: cap.get(2).map(|m| m.as_str().trim().to_string()),
        })
        .collect()
}

/// Extract tags from markdown content
pub fn extract_tags(content: &str) -> Vec<String> {
    TAG_REGEX
        .captures_iter(content)
        .map(|cap| cap[1].to_string())
        .collect::<Vec<_>>()
        .into_iter()
        .collect::<std::collections::HashSet<_>>()
        .into_iter()
        .collect()
}

/// Parse frontmatter from markdown content
pub fn parse_frontmatter(content: &str) -> Option<(Frontmatter, &str)> {
    let content = content.trim_start();
    if !content.starts_with("---") {
        return None;
    }

    let rest = &content[3..];
    let end_pos = rest.find("\n---")?;

    let yaml_str = &rest[..end_pos];
    let remaining = &rest[end_pos + 4..];

    match serde_yaml::from_str::<serde_yaml::Value>(yaml_str) {
        Ok(value) => {
            let fm = match value {
                serde_yaml::Value::Mapping(m) => m
                    .into_iter()
                    .filter_map(|(k, v)| {
                        let key = match k {
                            serde_yaml::Value::String(s) => s,
                            serde_yaml::Value::Number(n) => n.to_string(),
                            _ => return None,
                        };
                        let value = json_from_yaml(&v);
                        Some((key, value))
                    })
                    .collect(),
                _ => Frontmatter::new(),
            };
            Some((fm, remaining))
        }
        Err(_) => None,
    }
}

/// Convert YAML value to JSON value
fn json_from_yaml(yaml: &serde_yaml::Value) -> serde_json::Value {
    match yaml {
        serde_yaml::Value::Null => serde_json::Value::Null,
        serde_yaml::Value::Bool(b) => serde_json::Value::Bool(*b),
        serde_yaml::Value::Number(n) => {
            if let Some(i) = n.as_i64() {
                serde_json::Value::Number(i.into())
            } else if let Some(f) = n.as_f64() {
                serde_json::Number::from_f64(f)
                    .map(serde_json::Value::Number)
                    .unwrap_or(serde_json::Value::Null)
            } else {
                serde_json::Value::Null
            }
        }
        serde_yaml::Value::String(s) => serde_json::Value::String(s.clone()),
        serde_yaml::Value::Sequence(seq) => {
            serde_json::Value::Array(seq.iter().map(json_from_yaml).collect())
        }
        serde_yaml::Value::Mapping(map) => {
            let obj: std::collections::HashMap<String, serde_json::Value> = map
                .iter()
                .filter_map(|(k, v)| {
                    let key = match k {
                        serde_yaml::Value::String(s) => s.clone(),
                        serde_yaml::Value::Number(n) => n.to_string(),
                        _ => return None,
                    };
                    Some((key, json_from_yaml(v)))
                })
                .collect();
            serde_json::Value::Object(obj.into_iter().collect())
        }
        serde_yaml::Value::Tagged(tagged) => json_from_yaml(&tagged.value),
    }
}

/// Serialize frontmatter to markdown
pub fn serialize_frontmatter(frontmatter: &Frontmatter) -> String {
    if frontmatter.is_empty() {
        return String::new();
    }

    // Build YAML mapping manually
    let mut yaml_parts = Vec::new();
    for (key, value) in frontmatter {
        let yaml_val = yaml_from_json(value);
        yaml_parts.push(format!("{}: {}", key, yaml_val));
    }

    if yaml_parts.is_empty() {
        String::new()
    } else {
        format!("---\n{}\n---", yaml_parts.join("\n"))
    }
}

/// Convert JSON value to YAML string
fn yaml_from_json(value: &serde_json::Value) -> String {
    match value {
        serde_json::Value::Null => "null".to_string(),
        serde_json::Value::Bool(b) => b.to_string(),
        serde_json::Value::Number(n) => n.to_string(),
        serde_json::Value::String(s) => {
            if s.contains(':')
                || s.contains('#')
                || s.contains('\n')
                || s.starts_with('[')
                || s.starts_with('{')
            {
                format!("\"{}\"", s.replace('"', "\\\""))
            } else {
                s.clone()
            }
        }
        serde_json::Value::Array(arr) => {
            let items: Vec<String> = arr.iter().map(yaml_from_json).collect();
            format!("[{}]", items.join(", "))
        }
        serde_json::Value::Object(obj) => {
            let items: Vec<String> = obj
                .iter()
                .map(|(k, v)| format!("{}: {}", k, yaml_from_json(v)))
                .collect();
            format!("{{{}}}", items.join(", "))
        }
    }
}
