//! Batch insert operations

use pgvector::Vector;

use crate::client::PgVectorClient;
use crate::error::VectorError;
use crate::insert::normalize_embedding;
use crate::models::{BatchInsertRequest, BatchRecord, BatchResult};

/// Perform batch insert of vectors
pub async fn batch_insert(
    client: &PgVectorClient,
    request: &BatchInsertRequest,
) -> Result<BatchResult, VectorError> {
    // Validate request
    request
        .validate()
        .map_err(|e| VectorError::Validation(e))?;

    let mut ids = Vec::new();
    let mut errors = Vec::new();
    let mut inserted = 0;
    let mut failed = 0;

    // Process in smaller batches to avoid overwhelming the connection
    let batch_size = 50;
    let records = &request.records;

    for chunk in records.chunks(batch_size) {
        match process_batch_chunk(client, chunk).await {
            Ok(chunk_ids) => {
                inserted += chunk_ids.len();
                ids.extend(chunk_ids);
            }
            Err(e) => {
                failed += chunk.len();
                errors.push(e.to_string());
            }
        }
    }

    Ok(BatchResult {
        inserted,
        failed,
        ids,
        errors,
    })
}

/// Process a chunk of records
async fn process_batch_chunk(
    client: &PgVectorClient,
    records: &[BatchRecord],
) -> Result<Vec<i64>, VectorError> {
    let mut db_client = client.get_client().await?;
    let transaction = db_client.transaction().await.map_err(VectorError::Database)?;

    let mut ids = Vec::with_capacity(records.len());

    for record in records {
        let normalized = normalize_embedding(&record.embedding);
        let vector = Vector::from(normalized);

        let row = if let Some(id) = record.id {
            let params: Vec<&(dyn tokio_postgres::types::ToSql + Sync)> = vec![
                &id,
                &vector,
                &record.content,
            ];
            transaction
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
                &record.content,
            ];
            transaction
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

        ids.push(row.get("id"));
    }

    transaction.commit().await.map_err(VectorError::Database)?;

    tracing::debug!("Batch inserted {} vectors", ids.len());
    Ok(ids)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_validate_empty_records() {
        let request = BatchInsertRequest { records: vec![] };
        assert!(request.validate().is_err());
    }

    #[test]
    fn test_validate_empty_embedding() {
        let request = BatchInsertRequest {
            records: vec![BatchRecord {
                id: None,
                embedding: vec![],
                content: "test".to_string(),
                metadata: None,
            }],
        };
        assert!(request.validate().is_err());
    }

    #[test]
    fn test_validate_valid_request() {
        let request = BatchInsertRequest {
            records: vec![BatchRecord {
                id: None,
                embedding: vec![0.1, 0.2, 0.3],
                content: "test".to_string(),
                metadata: None,
            }],
        };
        assert!(request.validate().is_ok());
    }
}
