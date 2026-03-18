//! Combined test containers

pub mod postgres;
pub mod redis;

pub use postgres::TestPostgres;
pub use redis::TestRedis;

pub struct TestInfrastructure {
    pub postgres: TestPostgres,
    pub redis: TestRedis,
}

impl TestInfrastructure {
    pub async fn new() -> Self {
        let postgres = TestPostgres::new().await;
        let redis = TestRedis::new().await;

        Self { postgres, redis }
    }

    pub async fn database_url(&self) -> String {
        self.postgres.connection_string().await
    }

    pub async fn redis_url(&self) -> String {
        self.redis.connection_string().await
    }
}
