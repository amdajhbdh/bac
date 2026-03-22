//! Basic tests for vector-tools

#[cfg(test)]
mod tests {
    use vector_tools::init;
    use vector_tools::models::*;

    #[test]
    fn test_init() {
        // Just verify the library initializes without panic
        init();
    }

    #[test]
    fn test_search_request_default() {
        let request = SearchRequest::default();
        assert_eq!(request.top_k, Some(10));
        assert!(request.filters.is_none());
    }

    #[test]
    fn test_search_filters() {
        let filters = SearchFilters {
            category: Some("tech".to_string()),
            tags: Some(vec!["rust".to_string(), "pgvector".to_string()]),
            metadata: None,
        };

        assert_eq!(filters.category, Some("tech".to_string()));
        assert_eq!(filters.tags.as_ref().unwrap().len(), 2);
    }

    #[test]
    fn test_insert_request() {
        let request = InsertRequest {
            id: Some(1),
            embedding: vec![0.1, 0.2, 0.3],
            content: "Test content".to_string(),
            metadata: Some(serde_json::json!({"key": "value"})),
        };

        assert_eq!(request.id, Some(1));
        assert_eq!(request.embedding.len(), 3);
    }

    #[test]
    fn test_batch_insert_request_validation() {
        // Valid request
        let valid = BatchInsertRequest {
            records: vec![BatchRecord {
                id: None,
                embedding: vec![0.1, 0.2],
                content: "test".to_string(),
                metadata: None,
            }],
        };
        assert!(valid.validate().is_ok());

        // Empty records
        let empty = BatchInsertRequest { records: vec![] };
        assert!(empty.validate().is_err());

        // Empty embedding
        let empty_emb = BatchInsertRequest {
            records: vec![BatchRecord {
                id: None,
                embedding: vec![],
                content: "test".to_string(),
                metadata: None,
            }],
        };
        assert!(empty_emb.validate().is_err());
    }

    #[test]
    fn test_rebuild_request_default() {
        let request = RebuildRequest::default();
        assert_eq!(request.m, Some(16));
        assert_eq!(request.ef_construction, Some(64));
    }

    #[test]
    fn test_health_response() {
        let response = HealthResponse {
            status: "healthy".to_string(),
            version: "1.0.0".to_string(),
            database: "PostgreSQL 15".to_string(),
        };

        assert_eq!(response.status, "healthy");
    }

    #[test]
    fn test_delete_result() {
        let result = DeleteResult {
            id: 42,
            success: true,
        };

        assert_eq!(result.id, 42);
        assert!(result.success);
    }

    #[test]
    fn test_search_result_serialization() {
        let result = SearchResult {
            id: 1,
            content: "test".to_string(),
            metadata: Some(serde_json::json!({"source": "test"})),
            similarity: 0.95,
        };

        let json = serde_json::to_string(&result).unwrap();
        assert!(json.contains("\"id\":1"));
        assert!(json.contains("\"similarity\":0.95"));
    }
}
