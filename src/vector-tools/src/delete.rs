//! Delete operations

use crate::client::PgVectorClient;
use crate::error::VectorError;
use crate::models::DeleteResult;

/// Delete a vector by ID
pub async fn delete_vector(client: &PgVectorClient, id: i64) -> Result<DeleteResult, VectorError> {
    let db_client = client.get_client().await?;

    let result = db_client
        .execute("DELETE FROM vectors WHERE id = $1", &[&id])
        .await
        .map_err(VectorError::Database)?;

    if result == 0 {
        return Err(VectorError::NotFound(format!("Vector with id {} not found", id)));
    }

    tracing::debug!("Deleted vector with id: {}", id);

    Ok(DeleteResult {
        id,
        success: true,
    })
}

/// Rebuild the HNSW index
pub async fn rebuild_index(
    client: &PgVectorClient,
    m: Option<i32>,
    ef_construction: Option<i32>,
) -> Result<(bool, String, u64), VectorError> {
    use std::time::Instant;

    let start = Instant::now();
    let m = m.unwrap_or(16);
    let ef_construction = ef_construction.unwrap_or(64);

    let db_client = client.get_client().await?;

    // Drop existing HNSW index
    db_client
        .execute("DROP INDEX IF EXISTS vectors_hnsw_idx", &[])
        .await
        .map_err(VectorError::Database)?;

    // Create new HNSW index with specified parameters
    let sql = format!(
        r#"
        CREATE INDEX vectors_hnsw_idx 
        ON vectors USING hnsw (embedding vector_l2_ops)
        WITH (m = {}, ef_construction = {})
        "#,
        m, ef_construction
    );

    db_client
        .execute(&sql, &[])
        .await
        .map_err(VectorError::Database)?;

    let duration_ms = start.elapsed().as_millis() as u64;
    let message = format!(
        "HNSW index rebuilt with m={}, ef_construction={}",
        m, ef_construction
    );

    tracing::info!("{}", message);

    Ok((true, message, duration_ms))
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_delete_result() {
        use crate::models::DeleteResult;
        
        let result = DeleteResult {
            id: 1,
            success: true,
        };
        
        assert_eq!(result.id, 1);
        assert!(result.success);
    }
}
