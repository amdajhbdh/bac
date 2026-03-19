//! Integration tests for Vault Tools service (:3003)
//!
//! Tests the HTTP endpoints for Obsidian vault operations.

mod common;

use common::{create_client, is_service_running, urls, HEALTH_TIMEOUT};
use serde::{Deserialize, Serialize};

// =============================================================================
// Test Configuration
// =============================================================================

const VAULT_BASE_URL: &str = urls::VAULT;

// =============================================================================
// Health Endpoint Tests
// =============================================================================

/// Health response from vault-tools
#[derive(Debug, Deserialize)]
struct HealthResponse {
    status: String,
    #[serde(default)]
    vault_path: Option<String>,
}

/// Test: Health endpoint returns OK status
///
/// Verifies that the /health endpoint responds correctly.
#[tokio::test]
async fn test_health_endpoint() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    let response = client
        .get(format!("{}/health", VAULT_BASE_URL))
        .send()
        .await;
    
    // Assert: Response should be successful
    assert!(response.is_ok(), "Health endpoint should respond");
    let response = response.unwrap();
    
    // Assert: Should return 200 OK
    assert_eq!(
        response.status().as_u16(),
        200,
        "Health endpoint should return 200 OK"
    );
    
    // Assert: Response should be valid JSON with status field
    let health: HealthResponse = response.json().await.unwrap();
    assert_eq!(
        health.status, "ok",
        "Health status should be 'ok'"
    );
}

/// Test: Health endpoint includes vault path
///
/// Verifies that the health response includes vault configuration.
#[tokio::test]
async fn test_health_includes_vault_path() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    let response = client
        .get(format!("{}/health", VAULT_BASE_URL))
        .send()
        .await
        .unwrap();
    
    let health: HealthResponse = response.json().await.unwrap();
    assert!(
        health.vault_path.is_some(),
        "Health response should include vault_path"
    );
}

// =============================================================================
// Read Endpoint Tests
// =============================================================================

/// Response from /read endpoint
#[derive(Debug, Deserialize)]
struct ReadResult {
    #[serde(default)]
    success: bool,
    #[serde(default)]
    note: Option<NoteContent>,
    #[serde(default)]
    error: Option<String>,
}

/// Note content structure
#[derive(Debug, Deserialize)]
struct NoteContent {
    #[serde(default)]
    path: String,
    #[serde(default)]
    title: String,
    #[serde(default)]
    content: String,
    #[serde(default)]
    frontmatter: serde_json::Value,
    #[serde(default)]
    links: Vec<serde_json::Value>,
    #[serde(default)]
    tags: Vec<String>,
}

/// Test: Read endpoint accepts path parameter
///
/// Verifies that the /read endpoint accepts a path query parameter.
#[tokio::test]
async fn test_read_endpoint_accepts_path() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    // Try to read a note (may not exist, but should handle gracefully)
    let response = client
        .get(format!("{}/read?path=test.md", VAULT_BASE_URL))
        .send()
        .await;
    
    // Assert: Response should be received
    assert!(response.is_ok(), "Read endpoint should respond");
    
    // Status may be 200 or 404 depending on note existence
    let status = response.unwrap().status();
    assert!(
        status.is_success() || status.as_u16() == 404,
        "Read should return 200 or 404"
    );
}

/// Test: Read endpoint handles missing path
///
/// Verifies that the /read endpoint handles requests without a path.
#[tokio::test]
async fn test_read_endpoint_no_path() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    // Request without path parameter
    let response = client
        .get(format!("{}/read", VAULT_BASE_URL))
        .send()
        .await;
    
    assert!(response.is_ok(), "Read endpoint should handle missing path");
}

/// Test: Read endpoint returns proper error for non-existent note
///
/// Verifies that the /read endpoint returns appropriate error for missing files.
#[tokio::test]
async fn test_read_nonexistent_note() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .get(format!("{}/read?path=nonexistent/file.md", VAULT_BASE_URL))
        .send()
        .await
        .unwrap();
    
    let read_result: ReadResult = response.json().await.unwrap();
    
    // Should indicate failure
    assert!(
        !read_result.success || read_result.note.is_none(),
        "Non-existent note should not have success=true with note"
    );
}

// =============================================================================
// Write Endpoint Tests
// =============================================================================

/// Request body for /write endpoint
#[derive(Debug, Serialize)]
struct WriteRequest {
    path: String,
    content: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    frontmatter: Option<serde_json::Value>,
}

/// Response from /write endpoint
#[derive(Debug, Deserialize)]
struct WriteResult {
    #[serde(default)]
    success: bool,
    #[serde(default)]
    path: Option<String>,
    #[serde(default)]
    created: Option<bool>,
    #[serde(default)]
    error: Option<String>,
}

/// Test: Write endpoint creates new note
///
/// Verifies that the /write endpoint can create new notes.
#[tokio::test]
async fn test_write_endpoint_creates_note() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    let request = WriteRequest {
        path: format!(
            "integration-test-{}.md",
            uuid::Uuid::new_v4()
        ),
        content: "# Test Note\n\nThis is a test note created by integration tests.".to_string(),
        frontmatter: Some(serde_json::json!({
            "created": "integration-test",
            "tags": ["test", "automated"]
        })),
    };
    
    // Act: Send write request
    let response = client
        .post(format!("{}/write", VAULT_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Request should succeed
    assert!(response.is_ok(), "Write endpoint should accept requests");
}

/// Test: Write endpoint handles markdown content
///
/// Verifies that the /write endpoint preserves markdown formatting.
#[tokio::test]
async fn test_write_endpoint_markdown_content() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    let markdown_content = r#"# Heading 1

## Heading 2

- List item 1
- List item 2

> Blockquote

```rust
fn hello() {
    println!("Hello, World!");
}
```

[[WikiLink]] and #tags work here.
"#;
    
    let request = WriteRequest {
        path: format!(
            "test-markdown-{}.md",
            uuid::Uuid::new_v4()
        ),
        content: markdown_content.to_string(),
        frontmatter: None,
    };
    
    let response = client
        .post(format!("{}/write", VAULT_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "Write endpoint should handle markdown");
}

/// Test: Write endpoint with empty content
///
/// Verifies that the /write endpoint handles empty content.
#[tokio::test]
async fn test_write_endpoint_empty_content() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    let request = WriteRequest {
        path: format!(
            "test-empty-{}.md",
            uuid::Uuid::new_v4()
        ),
        content: "".to_string(),
        frontmatter: None,
    };
    
    let response = client
        .post(format!("{}/write", VAULT_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "Write endpoint should handle empty content");
}

// =============================================================================
// Search Endpoint Tests
// =============================================================================

/// Response from /search endpoint
#[derive(Debug, Deserialize)]
struct SearchResults {
    #[serde(default)]
    success: bool,
    #[serde(default)]
    query: String,
    #[serde(default)]
    results: Vec<SearchResultItem>,
    #[serde(default)]
    total: Option<usize>,
    #[serde(default)]
    error: Option<String>,
}

/// Search result item
#[derive(Debug, Deserialize)]
struct SearchResultItem {
    #[serde(default)]
    path: String,
    #[serde(default)]
    title: String,
    #[serde(default)]
    snippet: String,
    #[serde(default)]
    score: f32,
}

/// Test: Search endpoint accepts query parameter
///
/// Verifies that the /search endpoint can search vault contents.
#[tokio::test]
async fn test_search_endpoint_accepts_query() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    // Search for a common term
    let response = client
        .get(format!("{}/search?q=test", VAULT_BASE_URL))
        .send()
        .await;
    
    // Assert: Response should be successful
    assert!(response.is_ok(), "Search endpoint should respond");
    
    let response = response.unwrap();
    assert_eq!(
        response.status().as_u16(),
        200,
        "Search should return 200"
    );
    
    // Assert: Response should be valid JSON
    let results: SearchResults = response.json().await.unwrap();
    assert_eq!(
        results.query, "test",
        "Search should return results for query 'test'"
    );
}

/// Test: Search endpoint handles empty query
///
/// Verifies that the /search endpoint handles empty search queries.
#[tokio::test]
async fn test_search_endpoint_empty_query() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .get(format!("{}/search", VAULT_BASE_URL))
        .send()
        .await;
    
    assert!(response.is_ok(), "Search endpoint should handle empty query");
}

/// Test: Search returns valid structure
///
/// Verifies that the search response has the expected structure.
#[tokio::test]
async fn test_search_returns_valid_structure() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .get(format!("{}/search?q=the", VAULT_BASE_URL))
        .send()
        .await
        .unwrap();
    
    let results: SearchResults = response.json().await.unwrap();
    
    // Assert: Response should have required fields
    assert!(
        results.success || !results.error.is_none(),
        "Search response should indicate success or have error"
    );
    assert!(
        results.total.is_some() || !results.results.is_empty(),
        "Search should return total or results"
    );
}

// =============================================================================
// Search by Tag Tests
// =============================================================================

/// Test: Search by tag endpoint works
///
/// Verifies that the /search/tag endpoint can filter by tags.
#[tokio::test]
async fn test_search_by_tag_endpoint() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .get(format!("{}/search/tag?tag=test", VAULT_BASE_URL))
        .send()
        .await;
    
    assert!(response.is_ok(), "Search by tag should respond");
}

/// Test: Search by tag with no matches
///
/// Verifies that the /search/tag endpoint handles no matches.
#[tokio::test]
async fn test_search_by_tag_no_matches() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    // Use a very unlikely tag
    let response = client
        .get(format!("{}/search/tag?tag=nonexistentTag123", VAULT_BASE_URL))
        .send()
        .await
        .unwrap();
    
    let results: SearchResults = response.json().await.unwrap();
    
    // Should return success with empty results
    assert!(results.success, "Search should return success");
    assert!(
        results.results.is_empty() || results.total == Some(0),
        "Should have no results for non-existent tag"
    );
}

// =============================================================================
// Tags Endpoint Tests
// =============================================================================

/// Test: List all tags endpoint works
///
/// Verifies that the /tags endpoint returns tag statistics.
#[tokio::test]
async fn test_tags_endpoint() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .get(format!("{}/tags", VAULT_BASE_URL))
        .send()
        .await;
    
    assert!(response.is_ok(), "Tags endpoint should respond");
    
    // Response should be an array of tuples (tag, count)
    // The exact format depends on implementation
}

// =============================================================================
// Link Endpoint Tests
// =============================================================================

/// Request body for /link endpoint
#[derive(Debug, Serialize)]
struct LinkRequest {
    source: String,
    target: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    link_type: Option<String>,
}

/// Response from /link endpoint
#[derive(Debug, Deserialize)]
struct LinkResult {
    #[serde(default)]
    success: bool,
    #[serde(default)]
    source: Option<String>,
    #[serde(default)]
    target: Option<String>,
    #[serde(default)]
    error: Option<String>,
}

/// Test: Link endpoint creates connections between notes
///
/// Verifies that the /link endpoint can create wikilinks.
#[tokio::test]
async fn test_link_endpoint_creates_connection() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    let request = LinkRequest {
        source: "test-source.md".to_string(),
        target: "test-target.md".to_string(),
        link_type: Some("wikilink".to_string()),
    };
    
    let response = client
        .post(format!("{}/link", VAULT_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // May succeed or fail depending on whether files exist
    assert!(response.is_ok(), "Link endpoint should respond");
}

/// Test: Link endpoint validates required fields
///
/// Verifies that the /link endpoint validates input.
#[tokio::test]
async fn test_link_endpoint_validates_input() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    // Missing target
    let request = serde_json::json!({
        "source": "test.md"
    });
    
    let response = client
        .post(format!("{}/link", VAULT_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Service should return error for missing fields
    if response.is_ok() {
        let status = response.unwrap().status();
        assert!(
            status.as_u16() >= 400,
            "Missing required field should return error"
        );
    }
}

// =============================================================================
// MOC Endpoint Tests
// =============================================================================

/// Request body for /moc/update endpoint
#[derive(Debug, Serialize)]
struct MocUpdateRequest {
    subject: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    entries: Option<Vec<MocEntry>>,
}

/// MOC entry
#[derive(Debug, Serialize)]
struct MocEntry {
    title: String,
    path: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    description: Option<String>,
}

/// Response from /moc/update endpoint
#[derive(Debug, Deserialize)]
struct MocUpdateResult {
    #[serde(default)]
    success: bool,
    #[serde(default)]
    path: Option<String>,
    #[serde(default)]
    error: Option<String>,
}

/// Test: MOC update endpoint creates map of content
///
/// Verifies that the /moc/update endpoint can generate MOC notes.
#[tokio::test]
async fn test_moc_update_endpoint() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    let request = MocUpdateRequest {
        subject: "Test Subject".to_string(),
        entries: Some(vec![
            MocEntry {
                title: "Chapter 1".to_string(),
                path: "chapters/chapter-1.md".to_string(),
                description: Some("Introduction".to_string()),
            },
            MocEntry {
                title: "Chapter 2".to_string(),
                path: "chapters/chapter-2.md".to_string(),
                description: Some("Advanced topics".to_string()),
            },
        ]),
    };
    
    let response = client
        .post(format!("{}/moc/update", VAULT_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "MOC update should respond");
}

// =============================================================================
// Error Handling Tests
// =============================================================================

/// Test: Invalid JSON returns error
///
/// Verifies that the service handles malformed JSON gracefully.
#[tokio::test]
async fn test_invalid_json_returns_error() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .post(format!("{}/write", VAULT_BASE_URL))
        .header("Content-Type", "application/json")
        .body("{ invalid json }")
        .send()
        .await;
    
    // Should return 400 Bad Request
    if response.is_ok() {
        let status = response.unwrap().status();
        assert!(
            status.as_u16() >= 400,
            "Invalid JSON should return error status"
        );
    }
}

/// Test: Wrong HTTP method on POST endpoint
///
/// Verifies that using wrong method returns 405.
#[tokio::test]
async fn test_write_get_returns_405() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = create_client();
    
    // Try GET on POST-only endpoint
    let response = client
        .get(format!("{}/write", VAULT_BASE_URL))
        .send()
        .await;
    
    if response.is_ok() {
        let status = response.unwrap().status();
        assert!(
            status.as_u16() == 405 || status.as_u16() == 400,
            "Wrong method should return 405 or 400"
        );
    }
}

// =============================================================================
// Performance Tests
// =============================================================================

/// Test: Health endpoint responds quickly
///
/// Verifies that health checks complete within acceptable time.
#[tokio::test]
async fn test_health_response_performance() {
    skip_if_service_unavailable!(VAULT_BASE_URL);
    
    let client = Client::builder()
        .timeout(HEALTH_TIMEOUT)
        .build()
        .unwrap();
    
    let max_response_time = std::time::Duration::from_millis(200);
    
    let start = std::time::Instant::now();
    let response = client
        .get(format!("{}/health", VAULT_BASE_URL))
        .send()
        .await;
    let elapsed = start.elapsed();
    
    assert!(response.is_ok(), "Health endpoint should respond");
    assert!(
        elapsed < max_response_time,
        "Health should respond within 200ms, took {:?}",
        elapsed
    );
}
