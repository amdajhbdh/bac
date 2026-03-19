//! Insert operations for vectors

use pgvector::Vector;

use crate::client::PgVectorClient;
use crate::error::VectorError;
use crate::models::{InsertRequest, InsertResult};

/// Insert a single vector with metadata
pub async fn insert(
    client: &PgVectorClient,
    request: &InsertRequest,
) -> Result<InsertResult, VectorError> {
    // Validate embedding
    if request.embedding.is_empty() {
        return Err(VectorError::Validation("Embedding cannot be empty".to_string()));
    }

    // Normalize embedding
    let embedding = normalize_embedding(&request.embedding);
    let vector = Vector::from(embedding);

    let db_client = client.get_client().await?;

    let row = if let Some(id) = request.id {
        let params: Vec<&(dyn tokio_postgres::types::ToSql + Sync)> = vec![
            &id,
            &vector,
            &request.content,
        ];
        db_client
            .query_one(
                r#"
                INSERT INTO vectors (id, embedding, content, metadata)
                VALUES ($1, $2, $3, '{}'::jsonb)
                ON CONFLICT (id) DO UPDATE SET
                    embedding = EXCLUDED.embedding,
                    content = EXCLUDED.content,
                    metadata = EXCLUDED.metadata
                RETURNING id
                "#,
                &params,
            )
            .await
            .map_err(VectorError::Database)?
    } else {
        let params: Vec<&(dyn tokio_postgres::types::ToSql + Sync)> = vec![
            &vector,
            &request.content,
        ];
        db_client
            .query_one(
                r#"
                INSERT INTO vectors (embedding, content, metadata)
                VALUES ($1, $2, '{}'::jsonb)
                RETURNING id
                "#,
                &params,
            )
            .await
            .map_err(VectorError::Database)?
    };

    let id: i64 = row.get("id");

    tracing::debug!("Inserted vector with id: {}", id);

    Ok(InsertResult {
        id,
        success: true,
    })
}

/// Normalize vector to unit length for cosine similarity
pub fn normalize_embedding(embedding: &[f32]) -> Vec<f32> {
    let magnitude: f32 = embedding.iter().map(|x| x * x).sum::<f32>().sqrt();
    if magnitude == 0.0 {
        return embedding.to_vec();
    }
    embedding.iter().map(|x| x / magnitude).collect()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_normalize_embedding() {
        let embedding = vec![3.0, 4.0];
        let normalized = normalize_embedding(&embedding);
        
        let magnitude: f32 = normalized.iter().map(|x| x * x).sum::<f32>().sqrt();
        assert!((magnitude - 1.0).abs() < 0.001);
    }

    #[test]
    fn test_normalize_zero_vector() {
        let embedding = vec![0.0, 0.0];
        let normalized = normalize_embedding(&embedding);
        assert_eq!(normalized, embedding);
    }
}
