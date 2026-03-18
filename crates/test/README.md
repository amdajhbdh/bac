# BAC Test Infrastructure

Testcontainers-based testing infrastructure for BAC services.

## Features

- PostgreSQL testcontainer
- Redis testcontainer
- Combined infrastructure for full integration tests

## Usage

```rust
use bac_test::TestPostgres;

#[tokio::test]
async fn test_with_postgres() {
    let pg = TestPostgres::new().await;
    let conn_string = pg.connection_string().await;
    
    // Use connection string with your database client
}
```

## Docker Requirement

Testcontainers requires Docker to be running. Make sure Docker is installed and the daemon is running:

```bash
docker ps
```

## Available Containers

| Container | Default Port | Database |
|----------|-------------|----------|
| PostgreSQL | 5432 | bac_test |
| Redis | 6379 | - |
