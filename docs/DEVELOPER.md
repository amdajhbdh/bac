# BAC Knowledge System - Developer Guide

> **Target Audience**: Developers contributing to the BAC Knowledge System  
> **Prerequisites**: Rust, PostgreSQL, basic CLI familiarity

---

## Development Environment Setup

### 1. Install Dependencies

```bash
# Rust toolchain (if not installed)
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source ~/.cargo/env

# Verify installation
rustc --version
cargo --version

# Additional tools
npm install -g @anthropic-ai/claude-cli  # Optional: Claude for AI assistance
```

### 2. Clone and Configure

```bash
# Clone repository
git clone https://github.com/your-repo/bac.git
cd bac

# Copy environment template
cp .env.example .env

# Edit with your credentials
nano .env
```

Required environment variables:

```bash
# Database
NEON_DB_URL=postgresql://user:pass@host/db?sslmode=require

# AI
GEMINI_API_KEY=your_gemini_api_key

# Services (local defaults)
GEMINI_TOOLS_URL=http://localhost:3001
VECTOR_TOOLS_URL=http://localhost:3002
VAULT_TOOLS_URL=http://localhost:3003
```

### 3. Database Setup

```bash
# Connect to Neon database
psql $NEON_DB_URL

# Create extension
CREATE EXTENSION IF NOT EXISTS vector;

# Verify
SELECT extversion FROM pg_extension WHERE extname = 'vector';

# Run schema (if exists)
# \i sql/schema.sql
```

### 4. Build Services

```bash
# Build all services
cargo build --release

# Build specific service
cargo build -p gemini-tools --release

# Build with debug info
cargo build
```

### 5. Run Tests

```bash
# Run all tests
cargo test

# Run with output
cargo test -- --nocapture

# Run specific test
cargo test test_name

# Run with coverage
cargo tarpaulin --out Html
```

---

## Project Structure

```
bac/
├── src/
│   ├── api/                  # API gateway (if exists)
│   ├── gemini-tools/         # Gemini API wrapper
│   │   ├── src/
│   │   │   ├── main.rs       # Entry point
│   │   │   ├── handlers.rs   # HTTP handlers
│   │   │   ├── models.rs     # Data structures
│   │   │   └── lib.rs        # Library code
│   │   ├── Cargo.toml
│   │   └── tests/
│   ├── vector-tools/         # pgvector operations
│   ├── vault-tools/          # Obsidian operations
│   ├── cloud-tools/          # Cloud Shell integration
│   └── graph-tools/          # Knowledge graph
├── config/
│   └── aichat/               # aichat configuration
├── resources/
│   └── notes/                # Obsidian vault
├── sql/                      # Database schemas
├── scripts/                  # Automation scripts
├── tests/                    # Integration tests
└── docs/                     # Documentation
```

---

## Adding New Tools

### Step 1: Create the Crate

```bash
# Create new tool crate
cargo new src/my-new-tool
cd src/my-new-tool
```

### Step 2: Define Dependencies

```toml
# Cargo.toml
[package]
name = "my-new-tool"
version = "0.1.0"
edition = "2021"

[dependencies]
axum = "0.7"
serde = { version = "1", features = ["derive"] }
serde_json = "1"
tokio = { version = "1", features = ["full"] }
tracing = "0.1"
anyhow = "1"
```

### Step 3: Implement the Service

```rust
// src/main.rs
use axum::{routing::post, Router, Json};
use serde::{Deserialize, Serialize};

#[derive(Deserialize)]
struct Request {
    // Your request fields
}

#[derive(Serialize)]
struct Response {
    // Your response fields
}

async fn handle_action(Json(payload): Json<Request>) -> Json<Response> {
    // Your implementation
    Json(Response {
        // Your fields
    })
}

async fn health() -> &'static str {
    "healthy"
}

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/health", get(health))
        .route("/action", post(handle_action))
        .into_make_service();

    let listener = tokio::net::TcpListener::bind("0.0.0.0:3006")
        .await
        .unwrap();

    println!("my-new-tool running on :3006");
    axum::serve(listener, app).await.unwrap();
}
```

### Step 4: Add to Workspace

```toml
# Root Cargo.toml
[workspace]
members = [
    "src/gemini-tools",
    "src/vector-tools",
    "src/vault-tools",
    "src/cloud-tools",
    "src/graph-tools",
    "src/my-new-tool",  # Add this
]
```

### Step 5: Configure aichat

```yaml
# config/aichat/config.yaml
tools:
  - name: my-new-tool
    url: http://localhost:3006
    actions:
      - name: action
        description: "Description of what this action does"
        parameters:
          - name: field1
            type: string
            required: true
          - name: field2
            type: number
            required: false
```

### Step 6: Document and Test

```bash
# Test the service
cargo run -p my-new-tool

# Test health
curl http://localhost:3006/health

# Test action
curl -X POST http://localhost:3006/action \
  -H "Content-Type: application/json" \
  -d '{"field1": "value"}'
```

---

## Creating Custom Agents

### Agent Structure

Agents in aichat are defined by:

1. **Instructions file** (`config/aichat/agents/<name>/instructions.md`)
2. **Configuration** (`config/aichat/config.yaml`)

### Example: Subject-Specific Agent

```markdown
---
name: physics-tutor
model: gemini:gemini-3.0-flash
temperature: 0.4
tools:
  - gemini-tools
  - vector-tools
  - vault-tools
---

# Physics Tutor Agent

You are a physics tutor specializing in BAC-level physics.

## Specialization
- Classical mechanics
- Thermodynamics
- Electromagnetism
- Optics

## Response Format
Always use LaTeX for formulas:
- Inline: $E = mc^2$
- Display: $$\int_0^\infty f(x) dx$$

## Tool Usage
- Use vector-tools.search for finding relevant physics notes
- Use vault-tools.read for reviewing existing content
- Use gemini-tools.generate for creating new explanations
```

### Using Custom Agents

```bash
# Start aichat with specific agent
aichat --config config/aichat/config.yaml --agent physics-tutor

# Or switch in-session
> /agent physics-tutor
```

---

## Writing Tests

### Unit Tests

```rust
// src/gemini-tools/src/lib.rs
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_embedding_normalization() {
        let vector = vec![1.0, 2.0, 3.0];
        let normalized = normalize(&vector);
        
        // Check L2 norm is 1.0
        let magnitude: f32 = normalized.iter()
            .map(|x| x * x)
            .sum::<f32>()
            .sqrt();
        
        assert!((magnitude - 1.0).abs() < 0.001);
    }

    #[test]
    fn test_content_extraction() {
        let content = "# Title\n\nSome **bold** text.";
        let extracted = extract_headings(content);
        
        assert_eq!(extracted.len(), 1);
        assert_eq!(extracted[0], "Title");
    }
}
```

### Integration Tests

```rust
// tests/integration_test.rs
use reqwest;

#[tokio::test]
async fn test_gemini_tools_health() {
    let response = reqwest::get("http://localhost:3001/health")
        .await
        .unwrap();
    
    assert_eq!(response.status(), 200);
}

#[tokio::test]
async fn test_vector_tools_search() {
    let client = reqwest::Client::new();
    let response = client
        .post("http://localhost:3002/search")
        .json(&serde_json::json!({
            "query": vec![0.1; 768],
            "top_k": 5
        }))
        .send()
        .await
        .unwrap();
    
    assert_eq!(response.status(), 200);
    let body: serde_json::Value = response.json().await.unwrap();
    assert!(body["results"].is_array());
}
```

### Running Tests

```bash
# Run all tests
cargo test

# Run with output
cargo test -- --nocapture

# Run tests for specific crate
cargo test -p gemini-tools

# Run tests matching pattern
cargo test test_search

# Run with coverage
cargo tarpaulin
```

---

## Debugging Services

### Enable Debug Logging

```bash
# All services
RUST_LOG=debug cargo run -p <service>

# Specific service with trace
RUST_LOG=trace cargo run -p gemini-tools

# Log to file
RUST_LOG=debug cargo run -p vector-tools 2>&1 | tee logs/vector.log
```

### Check Service Health

```bash
# All services
curl http://localhost:8080/health  # If gateway exists

# Individual
curl http://localhost:3001/health  # gemini-tools
curl http://localhost:3002/health  # vector-tools
curl http://localhost:3003/health  # vault-tools
```

### Inspect Network

```bash
# Check ports in use
lsof -i :3001
lsof -i :3002

# Test connectivity
nc -zv localhost 3001
nc -zv localhost 3002
```

### Database Debugging

```bash
# Connect to Neon
psql $NEON_DB_URL

# Check tables
\dt

# Check indexes
\di

# Explain query
EXPLAIN ANALYZE SELECT * FROM notes 
WHERE embedding <=> '[0.1, 0.2, ...]'
ORDER BY embedding <=> '[0.1, 0.2, ...]'
LIMIT 5;
```

### Common Issues

| Issue | Debug Command | Solution |
|-------|---------------|----------|
| Service won't start | `lsof -i :3001` | Kill existing process |
| API timeout | Check `RUST_LOG=debug` | Increase timeout |
| Embedding mismatch | `SELECT embedding_dim FROM pg_attribute` | Check dimensions |
| HNSW slow | `EXPLAIN ANALYZE` query | Tune m/ef params |

---

## Performance Optimization

### Profiling

```bash
# CPU profiling
cargo flamegraph -p gemini-tools -- --release

# Memory profiling
cargo heaptrack target/release/gemini-tools

# Benchmark critical paths
cargo bench
```

### Index Tuning

```sql
-- Default HNSW parameters
CREATE INDEX ON notes USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- High quality (slower build, faster search)
CREATE INDEX ON notes USING hnsw (embedding vector_cosine_ops)
WITH (m = 32, ef_construction = 128);

-- Low memory (for large datasets)
CREATE INDEX ON notes USING hnsw (embedding vector_cosine_ops)
WITH (m = 8, ef_construction = 32);
```

### Batch Operations

```rust
// Batch insert for efficiency
async fn batch_insert(records: Vec<Record>) -> Result<()> {
    let batch_size = 100;
    
    for chunk in records.chunks(batch_size) {
        // Insert chunk
        insert_chunk(chunk).await?;
    }
    
    Ok(())
}
```

### Connection Pooling

```rust
// Use connection pool for database
use sqlx::PgPool;

let pool = PgPool::builder()
    .max_size(10)
    .min_idle(Some(5))
    .acquire_timeout(std::time::Duration::from_secs(30))
    .build(&database_url)
    .await?;
```

---

## Contributing Guidelines

### Code Style

- Run `cargo fmt` before committing
- Run `cargo clippy` and address warnings
- Maximum line length: 100 characters

```bash
# Format code
cargo fmt

# Lint
cargo clippy -- -D warnings

# Both
cargo fmt && cargo clippy
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new tool for OCR processing
fix: resolve vector search timeout issue
docs: update API documentation
test: add integration tests for vault-tools
refactor: simplify embedding normalization
```

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make changes with tests
4. Run formatting and linting
5. Commit with clear messages
6. Push and create PR
7. Address review feedback

### Code Review Checklist

- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No clippy warnings
- [ ] Code is idiomatic Rust
- [ ] Error handling is explicit
- [ ] No sensitive data in commits

---

## Release Process

### Version Bumping

```bash
# Update version in Cargo.toml files
# Follow SemVer: MAJOR.MINOR.PATCH

# Tag release
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

### Building Releases

```bash
# Build all services
cargo build --release --workspace

# Build single service
cargo build --release -p gemini-tools

# Cross-compile (example: Linux x86_64)
cargo build --release --target x86_64-unknown-linux-musl
```

### Deployment

```bash
# Stop existing services
pkill -f "cargo run"

# Deploy new binaries
cp target/release/gemini-tools /usr/local/bin/
cp target/release/vector-tools /usr/local/bin/
# ...

# Restart services
./scripts/bac-agent-daemon.sh
```

---

## Useful Commands

```bash
# Development
cargo watch -x build          # Auto-rebuild on changes
cargo test -p <crate>         # Test specific crate
cargo doc --open              # Generate docs

# Maintenance
cargo outdated                # Check for updates
cargo audit                  # Security audit
cargo tree -i <dep>          # Dependency tree

# Cleanup
cargo clean                   # Remove build artifacts
cargo sweep -f 30 -r 30       # Remove old downloads
```

---

## Getting Help

- **Documentation**: See `/docs/` directory
- **Issues**: Open at GitHub repository
- **Discussions**: Use GitHub Discussions
- **Chat**: Join project Discord/Slack

### Debug Checklist

```markdown
## Debug Session Checklist
- [ ] Service running?
- [ ] Ports accessible?
- [ ] Environment variables set?
- [ ] Database connected?
- [ ] Logs show errors?
- [ ] Network connectivity?
- [ ] API keys valid?
```

---

*Last updated: 2026-03-19*
