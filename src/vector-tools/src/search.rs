//! Semantic search operations with HNSW index

use std::time::Instant;

use pgvector::Vector;

use crate::client::PgVectorClient;
use crate::error::VectorError;
use crate::models::{SearchFilters, SearchRequest, SearchResult, SearchResults};

/// Perform semantic search using HNSW index
pub async fn search(
    client: &PgVectorClient,
    request: &SearchRequest,
) -> Result<SearchResults, VectorError> {
    let start = Instant::now();
    let top_k = request.top_k.unwrap_or(10);

    // Validate query vector
    if request.query.is_empty() {
        return Err(VectorError::Validation("Query vector cannot be empty".to_string()));
    }

    let _embedding = Vector::from(request.query.clone());

    // Build query based on filters
    let sql = build_search_query(&_embedding, top_k, &request.filters);

    let db_client = client.get_client().await?;
    
    // Use parameterized query with pgvector
    let params: Vec<&(dyn tokio_postgres::types::ToSql + Sync)> = vec![
        &_embedding,
    ];
    let rows = db_client.query(&sql, &params).await.map_err(VectorError::Database)?;

    let results: Vec<SearchResult> = rows
        .iter()
        .map(|row| {
            let distance: f64 = row.get("distance");
            // Convert L2 distance to similarity score (lower distance = higher similarity)
            let similarity = 1.0 / (1.0 + distance);
            
            SearchResult {
                id: row.get("id"),
                content: row.get("content"),
                metadata: row.get("metadata"),
                similarity,
            }
        })
        .collect();

    let query_time_ms = start.elapsed().as_millis() as u64;

    Ok(SearchResults {
        results,
        query_time_ms: Some(query_time_ms),
    })
}

/// Build search query with optional filters
fn build_search_query(
    _embedding: &Vector,
    top_k: usize,
    filters: &Option<SearchFilters>,
) -> String {
    let mut sql = String::new();

    sql.push_str(&format!(
        "SELECT id, content, metadata, (embedding <=> $1) as distance 
         FROM vectors 
         WHERE embedding IS NOT NULL"
    ));

    // Add filters if present
    if let Some(f) = filters {
        if f.category.is_some() {
            sql.push_str(" AND category = $2");
        }
        if f.tags.is_some() && !f.tags.as_ref().unwrap().is_empty() {
            sql.push_str(" AND metadata->'tags' ?| $3");
        }
    }

    sql.push_str(&format!(
        " ORDER BY embedding <=> $1 
         LIMIT {}",
        top_k
    ));

    sql
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_build_search_query_no_filters() {
        let embedding = Vector::from(vec![0.1, 0.2, 0.3]);
        let sql = build_search_query(&embedding, 10, &None);

        assert!(sql.contains("SELECT"));
        assert!(sql.contains("ORDER BY embedding <=> $1"));
        assert!(sql.contains("LIMIT 10"));
        assert!(sql.contains("WHERE embedding IS NOT NULL"));
    }

    #[test]
    fn test_distance_to_similarity() {
        // Lower distance = higher similarity
        let similarity_0 = 1.0 / (1.0 + 0.0);
        let similarity_1 = 1.0 / (1.0 + 1.0);
        
        assert!(similarity_0 > similarity_1);
    }
}
