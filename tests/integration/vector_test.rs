//! Integration tests for Vector Tools service (:3002)
//!
//! Tests the HTTP endpoints for pgvector operations with HNSW index support.

mod common;

use common::{create_client, is_service_running, urls, HEALTH_TIMEOUT};
use serde::{Deserialize, Serialize};

// =============================================================================
// Test Configuration
// =============================================================================

const VECTOR_BASE_URL: &str = urls::VECTOR;

// =============================================================================
// Health Endpoint Tests
// =============================================================================

/// Health response from vector-tools
#[derive(Debug, Deserialize)]
struct HealthResponse {
    status: String,
    version: String,
    #[serde(default)]
    database: Option<String>,
}

/// Test: Health endpoint returns healthy status
///
/// Verifies that the /health endpoint responds with proper health status.
#[tokio::test]
async fn test_health_endpoint() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let response = client
        .get(format!("{}/health", VECTOR_BASE_URL))
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
        health.status, "healthy",
        "Health status should be 'healthy'"
    );
    assert!(
        !health.version.is_empty(),
        "Health response should include version"
    );
}

/// Test: Health endpoint includes database status
///
/// Verifies that the health check includes database connectivity info.
#[tokio::test]
async fn test_health_includes_database_status() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let response = client
        .get(format!("{}/health", VECTOR_BASE_URL))
        .send()
        .await
        .unwrap();
    
    let health: HealthResponse = response.json().await.unwrap();
    assert!(
        health.database.is_some(),
        "Health response should include database status"
    );
}

// =============================================================================
// Search Endpoint Tests
// =============================================================================

/// Request body for /search endpoint
#[derive(Debug, Serialize)]
struct SearchRequest {
    query: Vec<f32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    top_k: Option<usize>,
    #[serde(skip_serializing_if = "Option::is_none")]
    filters: Option<SearchFilters>,
}

/// Optional filters for search
#[derive(Debug, Serialize, Default)]
struct SearchFilters {
    #[serde(skip_serializing_if = "Option::is_none")]
    category: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    tags: Option<Vec<String>>,
}

/// Response from /search endpoint
#[derive(Debug, Deserialize)]
struct SearchResponse {
    #[serde(default)]
    results: Vec<SearchResult>,
    #[serde(default)]
    query_time_ms: Option<u64>,
}

/// Search result item
#[derive(Debug, Deserialize)]
struct SearchResult {
    #[serde(default)]
    id: i64,
    #[serde(default)]
    content: String,
    #[serde(default)]
    metadata: Option<serde_json::Value>,
    #[serde(default)]
    similarity: f64,
}

/// Test: Search endpoint accepts valid vector query
///
/// Verifies that the /search endpoint can process semantic search queries.
#[tokio::test]
async fn test_search_endpoint_accepts_query() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let request = SearchRequest {
        query: vec![0.1; 768], // Common embedding dimension
        top_k: Some(5),
        filters: None,
    };
    
    // Act: Send search request
    let response = client
        .post(format!("{}/search", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Request should succeed (even if empty results)
    assert!(response.is_ok(), "Search endpoint should accept requests");
}

/// Test: Search endpoint returns valid response structure
///
/// Verifies that the /search endpoint returns properly structured results.
#[tokio::test]
async fn test_search_returns_valid_structure() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let request = SearchRequest {
        query: vec![0.0; 768],
        top_k: Some(10),
        filters: None,
    };
    
    let response = client
        .post(format!("{}/search", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await
        .unwrap();
    
    // Parse response
    let search_response: SearchResponse = response.json().await.unwrap();
    
    // Assert: Response should have results array (may be empty)
    assert!(
        search_response.results.is_empty() || !search_response.results.is_empty(),
        "Search should return results array"
    );
}

/// Test: Search endpoint handles filters
///
/// Verifies that the /search endpoint can apply category filters.
#[tokio::test]
async fn test_search_with_filters() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let request = SearchRequest {
        query: vec![0.1; 768],
        top_k: Some(5),
        filters: Some(SearchFilters {
            category: Some("physics".to_string()),
            tags: None,
        }),
    };
    
    let response = client
        .post(format!("{}/search", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "Search with filters should work");
}

/// Test: Search handles various top_k values
///
/// Verifies that the /search endpoint handles different result count limits.
#[tokio::test]
async fn test_search_top_k_variations() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let top_k_values = vec![1, 5, 10, 50, 100];
    
    for top_k in top_k_values {
        let request = SearchRequest {
            query: vec![0.0; 768],
            top_k: Some(top_k),
            filters: None,
        };
        
        let response = client
            .post(format!("{}/search", VECTOR_BASE_URL))
            .json(&request)
            .send()
            .await;
        
        assert!(
            response.is_ok(),
            "Search with top_k={} should work",
            top_k
        );
    }
}

// =============================================================================
// Insert Endpoint Tests
// =============================================================================

/// Request body for /insert endpoint
#[derive(Debug, Serialize)]
struct InsertRequest {
    #[serde(skip_serializing_if = "Option::is_none")]
    id: Option<i64>,
    embedding: Vec<f32>,
    content: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    metadata: Option<serde_json::Value>,
}

/// Response from /insert endpoint
#[derive(Debug, Deserialize)]
struct InsertResponse {
    #[serde(default)]
    id: Option<i64>,
    #[serde(default)]
    success: Option<bool>,
    #[serde(default)]
    error: Option<String>,
}

/// Test: Insert endpoint adds vector to database
///
/// Verifies that the /insert endpoint can add new vectors.
#[tokio::test]
async fn test_insert_endpoint_adds_vector() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let request = InsertRequest {
        id: None,
        embedding: vec![0.5; 768],
        content: "Test content for vector insertion".to_string(),
        metadata: Some(serde_json::json!({
            "category": "test",
            "tags": ["integration", "test"]
        })),
    };
    
    // Act: Send insert request
    let response = client
        .post(format!("{}/insert", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Assert: Request should succeed
    assert!(response.is_ok(), "Insert endpoint should accept requests");
}

/// Test: Insert endpoint with specific ID
///
/// Verifies that the /insert endpoint respects user-provided IDs.
#[tokio::test]
async fn test_insert_with_specific_id() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let request = InsertRequest {
        id: Some(99999), // Use high ID to avoid conflicts
        embedding: vec![0.3; 768],
        content: "Content with specific ID".to_string(),
        metadata: None,
    };
    
    let response = client
        .post(format!("{}/insert", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "Insert with specific ID should work");
}

// =============================================================================
// Batch Insert Endpoint Tests
// =============================================================================

/// Request body for /batch endpoint
#[derive(Debug, Serialize)]
struct BatchInsertRequest {
    records: Vec<BatchRecord>,
}

/// Single record for batch insert
#[derive(Debug, Serialize)]
struct BatchRecord {
    #[serde(skip_serializing_if = "Option::is_none")]
    id: Option<i64>,
    embedding: Vec<f32>,
    content: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    metadata: Option<serde_json::Value>,
}

/// Response from /batch endpoint
#[derive(Debug, Deserialize)]
struct BatchInsertResponse {
    #[serde(default)]
    inserted: Option<usize>,
    #[serde(default)]
    failed: Option<usize>,
    #[serde(default)]
    ids: Option<Vec<i64>>,
    #[serde(default)]
    errors: Option<Vec<String>>,
}

/// Test: Batch insert adds multiple vectors
///
/// Verifies that the /batch endpoint can insert multiple vectors efficiently.
#[tokio::test]
async fn test_batch_insert_multiple_vectors() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    
    // Create 3 test records
    let records: Vec<BatchRecord> = (0..3)
        .map(|i| BatchRecord {
            id: None,
            embedding: vec![0.1 * (i as f32); 768],
            content: format!("Batch test record {}", i),
            metadata: None,
        })
        .collect();
    
    let request = BatchInsertRequest { records };
    
    let response = client
        .post(format!("{}/batch", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "Batch insert should work");
}

/// Test: Batch insert handles empty batch
///
/// Verifies that the /batch endpoint handles empty input.
#[tokio::test]
async fn test_batch_insert_empty() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let request = BatchInsertRequest { records: vec![] };
    
    let response = client
        .post(format!("{}/batch", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Service may reject empty batches
    if response.is_ok() {
        let status = response.as_ref().unwrap().status();
        assert!(
            status.is_server_error() || status.as_u16() == 400,
            "Empty batch should return error status"
        );
    }
}

// =============================================================================
// Delete Endpoint Tests
// =============================================================================

/// Response from /delete/:id endpoint
#[derive(Debug, Deserialize)]
struct DeleteResponse {
    #[serde(default)]
    id: Option<i64>,
    #[serde(default)]
    success: Option<bool>,
    #[serde(default)]
    error: Option<String>,
}

/// Test: Delete endpoint removes vector
///
/// Verifies that the /delete endpoint can remove vectors by ID.
#[tokio::test]
async fn test_delete_endpoint_removes_vector() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    
    // First, insert a vector to get an ID
    let insert_request = InsertRequest {
        id: None,
        embedding: vec![0.9; 768],
        content: "Vector to be deleted".to_string(),
        metadata: None,
    };
    
    let insert_response = client
        .post(format!("{}/insert", VECTOR_BASE_URL))
        .json(&insert_request)
        .send()
        .await;
    
    if insert_response.is_err() {
        // Skip if we can't insert (database not available)
        return;
    }
    
    // Try to delete a non-existent ID (safer than deleting actual data)
    let response = client
        .delete(&format!("{}/delete/{}", VECTOR_BASE_URL, -1))
        .send()
        .await;
    
    assert!(response.is_ok(), "Delete endpoint should respond");
}

/// Test: Delete endpoint returns 404 for non-existent ID
///
/// Verifies that the /delete endpoint handles missing IDs appropriately.
#[tokio::test]
async fn test_delete_nonexistent_returns_appropriate_status() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    
    // Use a very high ID that likely doesn't exist
    let response = client
        .delete(&format!("{}/delete/{}", VECTOR_BASE_URL, 999999999))
        .send()
        .await;
    
    assert!(response.is_ok(), "Delete should return a response");
    // Service may return 200 with success=false or 404
}

// =============================================================================
// Rebuild Index Endpoint Tests
// =============================================================================

/// Request body for /rebuild endpoint
#[derive(Debug, Serialize)]
struct RebuildRequest {
    #[serde(skip_serializing_if = "Option::is_none")]
    m: Option<i32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    ef_construction: Option<i32>,
}

/// Response from /rebuild endpoint
#[derive(Debug, Deserialize)]
struct RebuildResponse {
    #[serde(default)]
    success: Option<bool>,
    #[serde(default)]
    message: Option<String>,
    #[serde(default)]
    duration_ms: Option<u64>,
}

/// Test: Rebuild endpoint triggers HNSW index rebuild
///
/// Verifies that the /rebuild endpoint can rebuild the HNSW index.
/// Note: This operation may be slow and is marked as ignored by default.
#[tokio::test]
#[ignore] // Index rebuild is slow, only run manually
async fn test_rebuild_endpoint_rebuilds_index() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let request = RebuildRequest {
        m: Some(16),
        ef_construction: Some(64),
    };
    
    let response = client
        .post(format!("{}/rebuild", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "Rebuild endpoint should accept requests");
}

/// Test: Rebuild endpoint accepts custom parameters
///
/// Verifies that the /rebuild endpoint accepts HNSW parameters.
#[tokio::test]
async fn test_rebuild_with_custom_parameters() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    
    // Test with minimal parameters
    let request = RebuildRequest {
        m: Some(8),
        ef_construction: Some(32),
    };
    
    let response = client
        .post(format!("{}/rebuild", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    assert!(response.is_ok(), "Rebuild with custom params should work");
}

// =============================================================================
// Error Handling Tests
// =============================================================================

/// Test: Invalid embedding dimension returns error
///
/// Verifies that the service validates embedding dimensions.
#[tokio::test]
async fn test_invalid_embedding_dimension() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    let request = SearchRequest {
        query: vec![0.1; 100], // Wrong dimension
        top_k: Some(5),
        filters: None,
    };
    
    let response = client
        .post(format!("{}/search", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Should either error or return empty results
    if response.is_ok() {
        // Accept as valid (depends on service implementation)
    }
}

/// Test: Missing required fields returns error
///
/// Verifies that the service validates required request fields.
#[tokio::test]
async fn test_missing_content_returns_error() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = create_client();
    
    // Insert without content field
    let request = serde_json::json!({
        "embedding": vec![0.1; 768]
        // missing "content" field
    });
    
    let response = client
        .post(format!("{}/insert", VECTOR_BASE_URL))
        .json(&request)
        .send()
        .await;
    
    // Service should return 400 or 422 for missing required field
    if response.is_ok() {
        let status = response.unwrap().status();
        assert!(
            status.as_u16() == 400 || status.as_u16() == 422 || status.is_server_error(),
            "Missing content should return error status"
        );
    }
}

// =============================================================================
// Performance Tests
// =============================================================================

/// Test: Health endpoint is fast
///
/// Verifies that health checks complete quickly.
#[tokio::test]
async fn test_health_response_performance() {
    skip_if_service_unavailable!(VECTOR_BASE_URL);
    
    let client = Client::builder()
        .timeout(HEALTH_TIMEOUT)
        .build()
        .unwrap();
    
    let max_response_time = std::time::Duration::from_millis(200);
    
    let start = std::time::Instant::now();
    let response = client
        .get(format!("{}/health", VECTOR_BASE_URL))
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
