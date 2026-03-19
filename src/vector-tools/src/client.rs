//! PostgreSQL pgvector client

use deadpool_postgres::{Config, Pool, PoolConfig, Runtime};
use tokio_postgres::NoTls;

use crate::error::VectorError;

/// Database configuration
#[derive(Clone)]
pub struct DatabaseConfig {
    pub url: String,
    pub max_connections: usize,
}

impl Default for DatabaseConfig {
    fn default() -> Self {
        Self {
            url: std::env::var("DATABASE_URL")
                .unwrap_or_else(|_| "postgresql://localhost/postgres".to_string()),
            max_connections: 10,
        }
    }
}

impl From<&config::Config> for DatabaseConfig {
    fn from(cfg: &config::Config) -> Self {
        Self {
            url: cfg
                .get_string("database.url")
                .unwrap_or_else(|_| std::env::var("DATABASE_URL").unwrap_or_default()),
            max_connections: cfg
                .get("database.max_connections")
                .unwrap_or(10),
        }
    }
}

/// pgvector client wrapper
#[derive(Clone)]
pub struct PgVectorClient {
    pool: Pool,
}

impl PgVectorClient {
    /// Create a new client from database URL
    pub async fn new(url: &str) -> Result<Self, VectorError> {
        Self::with_config(DatabaseConfig {
            url: url.to_string(),
            max_connections: 10,
        })
        .await
    }

    /// Create a new client from config
    pub async fn with_config(config: DatabaseConfig) -> Result<Self, VectorError> {
        let mut cfg = Config::new();
        cfg.host = Some(Self::parse_host(&config.url));
        cfg.port = Some(Self::parse_port(&config.url));
        cfg.user = Some(Self::parse_user(&config.url));
        cfg.password = Some(Self::parse_password(&config.url));
        cfg.dbname = Some(Self::parse_dbname(&config.url));
        
        // Set pool config with max connections
        cfg.pool = Some(PoolConfig {
            max_size: config.max_connections,
            ..Default::default()
        });

        let pool = cfg.create_pool(Some(Runtime::Tokio1), NoTls)?;

        Ok(Self { pool })
    }

    /// Get a connection from the pool
    pub async fn get_client(&self) -> Result<deadpool_postgres::Object, VectorError> {
        self.pool.get().await.map_err(VectorError::Pool)
    }

    /// Initialize the database schema
    pub async fn init_schema(&self) -> Result<(), VectorError> {
        let client = self.get_client().await?;

        // Enable pgvector extension
        client
            .execute("CREATE EXTENSION IF NOT EXISTS vector", &[])
            .await
            .map_err(VectorError::Database)?;

        // Create vectors table
        client
            .execute(
                r#"
                CREATE TABLE IF NOT EXISTS vectors (
                    id BIGSERIAL PRIMARY KEY,
                    embedding vector(1536),
                    content TEXT NOT NULL,
                    metadata JSONB,
                    created_at TIMESTAMPTZ DEFAULT NOW()
                )
                "#,
                &[],
            )
            .await
            .map_err(VectorError::Database)?;

        // Create HNSW index (if not exists)
        client
            .execute(
                r#"
                CREATE INDEX IF NOT EXISTS vectors_hnsw_idx 
                ON vectors USING hnsw (embedding vector_l2_ops)
                WITH (m = 16, ef_construction = 64)
                "#,
                &[],
            )
            .await
            .map_err(VectorError::Database)?;

        tracing::info!("Database schema initialized");
        Ok(())
    }

    /// Check database health
    pub async fn health_check(&self) -> Result<String, VectorError> {
        let client = self.get_client().await?;
        let row = client.query_one("SELECT version()", &[]).await.map_err(VectorError::Database)?;
        Ok(row.get::<_, String>(0))
    }

    // URL parsing helpers
    fn parse_host(url: &str) -> String {
        Self::extract_param(url, "host=").unwrap_or_else(|| "localhost".to_string())
    }

    fn parse_port(url: &str) -> u16 {
        Self::extract_param(url, "port=")
            .and_then(|s| s.parse().ok())
            .unwrap_or(5432)
    }

    fn parse_user(url: &str) -> String {
        Self::extract_param(url, "user=").unwrap_or_else(|| "postgres".to_string())
    }

    fn parse_password(url: &str) -> String {
        Self::extract_param(url, "password=").unwrap_or_default()
    }

    fn parse_dbname(url: &str) -> String {
        Self::extract_param(url, "dbname=").unwrap_or_else(|| "postgres".to_string())
    }

    fn extract_param(url: &str, prefix: &str) -> Option<String> {
        let after_slash = url.split('@').last()?;
        let params = after_slash.split('?').next()?;
        
        params
            .split('&')
            .find(|p| p.starts_with(prefix))
            .map(|p| p[prefix.len()..].to_string())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_dbname() {
        let url = "postgresql://user:pass@host:5432/mydb?sslmode=require";
        assert_eq!(PgVectorClient::parse_dbname(url), "mydb");
    }

    #[test]
    fn test_parse_user() {
        let url = "postgresql://user:pass@host:5432/mydb";
        assert_eq!(PgVectorClient::parse_user(url), "user");
    }

    #[test]
    fn test_parse_host() {
        let url = "postgresql://user:pass@myhost:5432/mydb";
        assert_eq!(PgVectorClient::parse_host(url), "myhost");
    }
}
