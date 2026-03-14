# Contributing to BAC Unified

Thank you for your interest in contributing to BAC Unified!

## Code of Conduct

By participating in this project, you agree to follow our [Code of Conduct](CODE_OF_CONDUCT.md) (if exists) or uphold respectful and inclusive communication standards.

## How Can I Contribute?

### Reporting Bugs

1. Check existing issues to avoid duplicates
2. Use issue templates if available
3. Provide clear reproduction steps
4. Include relevant environment details

### Suggesting Features

1. Search existing proposals first
2. Use clear, descriptive titles
3. Explain the use case and expected behavior
4. Consider implementation complexity

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes following our code style
4. Add tests for new functionality
5. Ensure all tests pass
6. Update documentation if needed
7. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.22+
- Rust (latest stable)
- PostgreSQL 15+ with pgvector extension
- Node.js 18+ (for some tools)

### Local Development

```bash
# Clone the repository
git clone https://github.com/bac-unified/bac.git
cd bac

# Install Go dependencies
cd src/agent
go mod download

# Install Rust dependencies
cd ../noon
cargo fetch

# Set up environment
cp .env.example .env
# Edit .env with your credentials

# Run database migrations
psql $NEON_DB_URL -f sql/schema.sql

# Run the agent
cd ../src/agent
go run cmd/main.go

# Run the API
cd ../src/api
go run cmd/main.go
```

## Testing

### Running Tests

```bash
# All tests
go test ./...

# Single package
go test -v ./internal/solver/

# With coverage
go test -cover ./...

# With race detector
go test -race ./...
```

### Writing Tests

- Place tests in same package with `_test.go` suffix
- Use table-driven test patterns
- Follow naming convention: `TestFunctionName`

## Coding Standards

### Go

- Use `log/slog` for logging (never `fmt.Print*`)
- Wrap errors with context: `fmt.Errorf("doing thing: %w", err)`
- Group imports: stdlib, external, internal
- Use meaningful variable names
- Keep functions small and focused

### Commit Messages

- Use clear, descriptive messages
- Start with verb: "Add", "Fix", "Update", "Remove"
- Reference issues: "Fixes #123" or "Closes #456"

Example:
```
Add NLM caching with PostgreSQL

Implements SQLite-free cache using PostgreSQL with pgvector
for semantic similarity. Includes cache invalidation and
fallback to Ollama when NLM is unavailable.

Fixes #42
```

## Pull Request Process

1. Ensure all tests pass locally
2. Update documentation for any changes
3. PRs should target the `main` branch
4. Include a clear description of changes
5. Respond to review feedback promptly

## Review Criteria

PRs are reviewed based on:

- [ ] Code quality and style
- [ ] Test coverage
- [ ] Documentation updates
- [ ] No security vulnerabilities
- [ ] Performance implications considered

## Project Structure

```
bac/
├── src/
│   ├── agent/          # Go CLI agent
│   ├── api/            # REST API
│   └── noon/           # Rust animation engine
├── sql/                # Database schemas
├── config/             # Configuration files
├── doc/                # Documentation
└── openspec/           # Change management
```

## Getting Help

- Open an issue for bugs or feature requests
- Join discussions (if community exists)
- Check existing documentation in `doc/`

## Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes (for significant contributions)

---

*Last Updated: 2026-03-01*
