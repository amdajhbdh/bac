//! Vector Tools - pgvector operations service (:3002)
//!
//! HTTP service for semantic search and vector operations using pgvector with HNSW index.

use vector_tools::service;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let port = std::env::var("PORT")
        .unwrap_or_else(|_| "3002".to_string())
        .parse::<u16>()?;

    service::run(port).await?;

    Ok(())
}
