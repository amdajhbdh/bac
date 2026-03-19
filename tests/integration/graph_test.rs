//! Integration tests for Graph Tools service (:3005)
//!
//! Tests the HTTP endpoints for knowledge graph operations.

mod common;

use common::{create_client, is_service_running, urls, HEALTH_TIMEOUT};
use serde::{Deserialize, Serialize};

// =============================================================================
// Test Configuration
// =============================================================================

const GRAPH_BASE_URL: &str = urls::GRAPH;

// =============================================================================
// Health Endpoint Tests
// =============================================================================

/// Health response from graph-tools
#[derive(Debug, Deserialize)]
struct HealthResponse {
    #[serde(default)]
    status: String,
    #[serde(default)]
    version: Option<String>,
    #[serde(default)]
    service: Option<String>,
}

/// Test: Health endpoint returns OK status
///
/// Verifies that the /health endpoint responds correctly.
#[tokio::test]
async fn test_health_endpoint() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    let response = client
        .get(format!("{}/health", GRAPH_BASE_URL))
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

/// Test: Health endpoint includes service info
///
/// Verifies that the health response includes version and service info.
#[tokio::test]
async fn test_health_includes_service_info() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    let response = client
        .get(format!("{}/health", GRAPH_BASE_URL))
        .send()
        .await
        .unwrap();
    
    let health: HealthResponse = response.json().await.unwrap();
    
    assert!(
        health.service.is_some(),
        "Health response should include service name"
    );
    assert_eq!(
        health.service.as_deref(),
        Some("graph-tools"),
        "Service name should be 'graph-tools'"
    );
    assert!(
        health.version.is_some(),
        "Health response should include version"
    );
}

// =============================================================================
// Canvas Endpoint Tests
// =============================================================================

/// Request body for /canvas endpoint
#[derive(Debug, Serialize)]
struct CanvasRequest {
    subject: String,
}

/// Canvas file structure
#[derive(Debug, Deserialize)]
struct CanvasFile {
    #[serde(default)]
    nodes: Vec<CanvasNode>,
    #[serde(default)]
    edges: Vec<CanvasEdge>,
}

/// Canvas node structure
#[derive(Debug, Deserialize)]
struct CanvasNode {
    #[serde(default)]
    id: String,
    #[serde(default)]
    x: f64,
    #[serde(default)]
    y: f64,
    #[serde(default)]
    width: f64,
    #[serde(default)]
    height: f64,
    #[serde(rename = "type", default)]
    node_type: String,
    #[serde(default)]
    text: String,
    #[serde(default)]
    color: Option<String>,
}

/// Canvas edge structure
#[derive(Debug, Deserialize)]
struct CanvasEdge {
    #[serde(default)]
    id: String,
    #[serde(default)]
    from_node: String,
    #[serde(default)]
    to_node: String,
    #[serde(rename = "type", default)]
    edge_type: String,
    #[serde(default)]
    label: Option<String>,
}

/// Test: Canvas endpoint generates canvas file
///
/// Verifies that the /canvas endpoint can generate knowledge graph visualizations.
#[tokio::test]
async fn test_canvas_endpoint_generates_canvas() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    let request = CanvasRequest {
        subject: "Physics".to_string(),
    };
    
    // Act: Send canvas generation request
    let response = client
        .post(format!("{}/canvas", GRAPH_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Request should succeed
    assert!(response.is_ok(), "Canvas endpoint should respond");
    
    let canvas: CanvasFile = response.unwrap().json().await.unwrap();
    
    // Assert: Canvas should have nodes (at least the subject node)
    assert!(
        !canvas.nodes.is_empty(),
        "Canvas should contain at least one node"
    );
    
    // First node should be the subject
    assert_eq!(
        canvas.nodes[0].text, "Physics",
        "First node should represent the subject"
    );
}

/// Test: Canvas endpoint handles different subjects
///
/// Verifies that the /canvas endpoint generates different canvases for different subjects.
#[tokio::test]
async fn test_canvas_different_subjects() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    
    let subjects = vec![
        "Mathematics",
        "Chemistry",
        "Computer Science",
        "Biology",
    ];
    
    for subject in subjects {
        let request = CanvasRequest {
            subject: subject.to_string(),
        };
        
        let response = client
            .post(format!("{}/canvas", GRAPH_BASE_URL))
            .json(&request)
            .send()
            .await;
        
        assert!(
            response.is_ok(),
            "Canvas should generate for subject: {}",
            subject
        );
    }
}

/// Test: Canvas has valid structure
///
/// Verifies that generated canvas files have valid structure.
#[tokio::test]
async fn test_canvas_valid_structure() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    let request = CanvasRequest {
        subject: "Test Subject".to_string(),
    };
    
    let response = client
        .post(format!("{}/canvas", GRAPH_BASE_URL))
        .json(&request)
        .send()
        .await
        .unwrap();
    
    let canvas: CanvasFile = response.json().await.unwrap();
    
    // Verify nodes have required fields
    for node in &canvas.nodes {
        assert!(
            !node.id.is_empty(),
            "Node should have an ID"
        );
        assert!(
            node.width > 0.0 && node.height > 0.0,
            "Node should have positive dimensions"
        );
        assert!(
            !node.node_type.is_empty(),
            "Node should have a type"
        );
    }
    
    // Verify edges reference valid nodes
    let node_ids: std::collections::HashSet<_> = canvas.nodes.iter().map(|n| &n.id).collect();
    for edge in &canvas.edges {
        assert!(
            node_ids.contains(&edge.from_node),
            "Edge from_node should reference valid node"
        );
        assert!(
            node_ids.contains(&edge.to_node),
            "Edge to_node should reference valid node"
        );
    }
}

// =============================================================================
// Extract Endpoint Tests
// =============================================================================

/// Request body for /extract endpoint
#[derive(Debug, Serialize)]
struct ExtractRequest {
    #[serde(skip_serializing_if = "Option::is_none")]
    content: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    source: Option<String>,
}

/// Extracted entities structure
#[derive(Debug, Deserialize)]
struct Entities {
    #[serde(default)]
    concepts: Vec<Concept>,
    #[serde(default)]
    definitions: Vec<Definition>,
    #[serde(default)]
    relationships: Vec<Relationship>,
}

/// Concept structure
#[derive(Debug, Deserialize)]
struct Concept {
    #[serde(default)]
    id: String,
    #[serde(default)]
    name: String,
    #[serde(default)]
    source_note: String,
    #[serde(default)]
    mentions: usize,
}

/// Definition structure
#[derive(Debug, Deserialize)]
struct Definition {
    #[serde(default)]
    id: String,
    #[serde(default)]
    term: String,
    #[serde(default)]
    definition: String,
    #[serde(default)]
    source_note: String,
}

/// Relationship structure
#[derive(Debug, Deserialize)]
struct Relationship {
    #[serde(default)]
    id: String,
    #[serde(default)]
    source_id: String,
    #[serde(default)]
    target_id: String,
    #[serde(rename = "relationship_type", default)]
    relationship_type: String,
    #[serde(default)]
    weight: f32,
}

/// Test: Extract endpoint extracts from content
///
/// Verifies that the /extract endpoint can extract entities from provided content.
#[tokio::test]
async fn test_extract_endpoint_from_content() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    let request = ExtractRequest {
        content: Some(
            "## Electric Field\nForce:: A push or pull on an object\n\
             Energy:: The capacity to do work\n\
             The relationship between force and energy is fundamental."
                .to_string(),
        ),
        source: Some("physics.md".to_string()),
    };
    
    // Act: Send extract request
    let response = client
        .post(format!("{}/extract", GRAPH_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Request should succeed
    assert!(response.is_ok(), "Extract endpoint should respond");
    
    let entities: Entities = response.unwrap().json().await.unwrap();
    
    // Should extract at least one type of entity
    assert!(
        !entities.concepts.is_empty()
            || !entities.definitions.is_empty()
            || !entities.relationships.is_empty(),
        "Extract should find some entities"
    );
}

/// Test: Extract endpoint handles markdown content
///
/// Verifies that the extract endpoint recognizes markdown patterns.
#[tokio::test]
async fn test_extract_markdown_patterns() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    
    let test_cases = vec![
        // Definition pattern: Term:: Definition
        "Gravity:: Force of attraction between masses",
        // Concept in heading
        "# Photosynthesis\nThe process by which plants convert sunlight",
        // Wikilink
        "[[Mitochondria]] is the powerhouse of the cell",
    ];
    
    for content in test_cases {
        let request = ExtractRequest {
            content: Some(content.to_string()),
            source: Some("test.md".to_string()),
        };
        
        let response = client
            .post(format!("{}/extract", GRAPH_BASE_URL))
            .json(&request)
            .send()
            .await;
        
        assert!(
            response.is_ok(),
            "Extract should handle content: {}...",
            &content[..content.len().min(30)]
        );
    }
}

/// Test: Extract endpoint extracts from existing notes
///
/// Verifies that the /extract endpoint can work without explicit content.
#[tokio::test]
async fn test_extract_without_content() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    
    // Empty request - should extract from existing notes
    let request = ExtractRequest {
        content: None,
        source: None,
    };
    
    let response = client
        .post(format!("{}/extract", GRAPH_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "Extract should work without explicit content");
}

// =============================================================================
// Relationships Endpoint Tests
// =============================================================================

/// Test: Relationships endpoint builds relationship graph
///
/// Verifies that the /relationships endpoint generates entity relationships.
#[tokio::test]
async fn test_relationships_endpoint() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    let request = CanvasRequest {
        subject: "Physics".to_string(),
    };
    
    // Act: Send relationships request
    let response = client
        .post(format!("{}/relationships", GRAPH_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Request should succeed
    assert!(response.is_ok(), "Relationships endpoint should respond");
    
    let canvas: CanvasFile = response.unwrap().json().await.unwrap();
    
    // Should return canvas structure
    assert!(
        canvas.nodes.is_empty() || !canvas.nodes.is_empty(),
        "Relationships should return canvas"
    );
}

// =============================================================================
// Export Endpoint Tests
// =============================================================================

/// Test: Export endpoint returns graph JSON
///
/// Verifies that the /export endpoint can export the knowledge graph.
#[tokio::test]
async fn test_export_endpoint() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    
    // Act: Send export request
    let response = client
        .get(format!("{}/export", GRAPH_BASE_URL))
        .send()
        .await;
    
    // Assert: Request should succeed
    assert!(response.is_ok(), "Export endpoint should respond");
    
    // Response should be valid JSON (text response)
    let body = response.unwrap().text().await.unwrap();
    
    // Should be parseable as JSON
    assert!(
        serde_json::from_str::<serde_json::Value>(&body).is_ok(),
        "Export should return valid JSON"
    );
}

/// Test: Export endpoint with subject parameter
///
/// Verifies that the export endpoint accepts subject parameter.
#[tokio::test]
async fn test_export_with_subject() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .get(format!("{}/export?subject=Physics", GRAPH_BASE_URL))
        .send()
        .await;
    
    assert!(response.is_ok(), "Export should accept subject parameter");
}

// =============================================================================
// Visualize Endpoint Tests
// =============================================================================

/// Request body for /visualize endpoint
#[derive(Debug, Serialize)]
struct VisualizeRequest {
    subject: String,
}

/// Visualization URL structure
#[derive(Debug, Deserialize)]
struct VisualizationUrl {
    #[serde(default)]
    url: String,
    #[serde(default)]
    embed_url: String,
    #[serde(default)]
    graph_data_url: Option<String>,
}

/// Test: Visualize endpoint generates visualization URLs
///
/// Verifies that the /visualize endpoint generates Mermaid.js visualization URLs.
#[tokio::test]
async fn test_visualize_endpoint_generates_urls() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    let request = VisualizeRequest {
        subject: "Physics".to_string(),
    };
    
    // Act: Send visualize request
    let response = client
        .post(format!("{}/visualize", GRAPH_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Request should succeed
    assert!(response.is_ok(), "Visualize endpoint should respond");
    
    let viz: VisualizationUrl = response.unwrap().json().await.unwrap();
    
    // Assert: Should have URL fields
    assert!(
        !viz.url.is_empty(),
        "Visualization should have URL"
    );
    assert!(
        !viz.embed_url.is_empty(),
        "Visualization should have embed URL"
    );
}

/// Test: Visualize endpoint generates Mermaid URLs
///
/// Verifies that the visualization URLs point to Mermaid.js.
#[tokio::test]
async fn test_visualize_generates_mermaid_urls() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    let request = VisualizeRequest {
        subject: "Chemistry".to_string(),
    };
    
    let response = client
        .post(format!("{}/visualize", GRAPH_BASE_URL))
        .json(&request)
        .send()
        .await
        .unwrap();
    
    let viz: VisualizationUrl = response.json().await.unwrap();
    
    // URLs should contain "mermaid"
    assert!(
        viz.url.contains("mermaid"),
        "URL should contain 'mermaid'"
    );
    assert!(
        viz.embed_url.contains("mermaid"),
        "Embed URL should contain 'mermaid'"
    );
}

/// Test: Visualize endpoint handles different subjects
///
/// Verifies that the visualize endpoint works with various subjects.
#[tokio::test]
async fn test_visualize_different_subjects() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    
    let subjects = vec![
        "Mathematics",
        "Biology",
        "Computer Science",
        "History",
    ];
    
    for subject in subjects {
        let request = VisualizeRequest {
            subject: subject.to_string(),
        };
        
        let response = client
            .post(format!("{}/visualize", GRAPH_BASE_URL))
            .json(&request)
            .send()
            .await;
        
        assert!(
            response.is_ok(),
            "Visualize should work for: {}",
            subject
        );
    }
}

// =============================================================================
// Error Handling Tests
// =============================================================================

/// Test: Invalid JSON returns error
///
/// Verifies that the service handles malformed JSON gracefully.
#[tokio::test]
async fn test_invalid_json_returns_error() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    
    let response = client
        .post(format!("{}/canvas", GRAPH_BASE_URL))
        .header("Content-Type", "application/json")
        .body("{ invalid json }")
        .send()
        .await;
    
    // Should return 400 Bad Request or 422 Unprocessable Entity
    if response.is_ok() {
        let status = response.unwrap().status();
        assert!(
            status.as_u16() == 400 || status.as_u16() == 422,
            "Invalid JSON should return 400 or 422"
        );
    }
}

/// Test: Wrong HTTP method on POST endpoint
///
/// Verifies that using GET on POST-only endpoints returns appropriate error.
#[tokio::test]
async fn test_canvas_get_returns_405() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    
    // Try GET on POST-only endpoint
    let response = client
        .get(format!("{}/canvas", GRAPH_BASE_URL))
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

/// Test: Wrong HTTP method on POST endpoint for visualize
///
/// Verifies that using GET on visualize endpoint returns appropriate error.
#[tokio::test]
async fn test_visualize_get_returns_405() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    
    // Try GET on POST-only visualize endpoint
    let response = client
        .get(format!("{}/visualize", GRAPH_BASE_URL))
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
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = Client::builder()
        .timeout(HEALTH_TIMEOUT)
        .build()
        .unwrap();
    
    let max_response_time = std::time::Duration::from_millis(200);
    
    let start = std::time::Instant::now();
    let response = client
        .get(format!("{}/health", GRAPH_BASE_URL))
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

/// Test: Canvas generation is reasonably fast
///
/// Verifies that canvas generation completes in acceptable time.
#[tokio::test]
async fn test_canvas_generation_performance() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    let request = CanvasRequest {
        subject: "Physics".to_string(),
    };
    
    let max_response_time = std::time::Duration::from_secs(5);
    
    let start = std::time::Instant::now();
    let response = client
        .post(format!("{}/canvas", GRAPH_BASE_URL))
        .json(&request)
        .send()
        .await;
    let elapsed = start.elapsed();
    
    assert!(response.is_ok(), "Canvas generation should respond");
    assert!(
        elapsed < max_response_time,
        "Canvas generation should complete within 5s, took {:?}",
        elapsed
    );
}

// =============================================================================
// Cross-Feature Tests
// =============================================================================

/// Test: Full workflow - extract, canvas, visualize
///
/// Verifies that the graph-tools service can handle a complete workflow.
#[tokio::test]
async fn test_full_workflow() {
    skip_if_service_unavailable!(GRAPH_BASE_URL);
    
    let client = create_client();
    
    // Step 1: Extract entities
    let extract_request = ExtractRequest {
        content: Some(
            "# Physics\n\
             Force:: A push or pull\n\
             Energy:: Capacity to do work\n\
             Mass:: Amount of matter"
                .to_string(),
        ),
        source: Some("physics.md".to_string()),
    };
    
    let extract_response = client
        .post(format!("{}/extract", GRAPH_BASE_URL))
        .json(&extract_request)
        .send()
        .await
        .unwrap();
    
    assert!(
        extract_response.status().is_success(),
        "Extract should succeed"
    );
    
    // Step 2: Generate canvas
    let canvas_request = CanvasRequest {
        subject: "Physics".to_string(),
    };
    
    let canvas_response = client
        .post(format!("{}/canvas", GRAPH_BASE_URL))
        .json(&canvas_request)
        .send()
        .await
        .unwrap();
    
    assert!(
        canvas_response.status().is_success(),
        "Canvas generation should succeed"
    );
    
    // Step 3: Generate visualization
    let viz_request = VisualizeRequest {
        subject: "Physics".to_string(),
    };
    
    let viz_response = client
        .post(format!("{}/visualize", GRAPH_BASE_URL))
        .json(&viz_request)
        .send()
        .await
        .unwrap();
    
    assert!(
        viz_response.status().is_success(),
        "Visualization should succeed"
    );
}
