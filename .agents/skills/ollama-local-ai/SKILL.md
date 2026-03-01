# Skill: Ollama Local AI

## Purpose

Ollama for local AI inference in BAC Unified.

## When to use

- Solving math problems locally
- Generating embeddings
- Offline-first AI processing

## Priority Order

1. **Ollama** - Local (primary)
2. **NotebookLM** - Cloud API (secondary)
3. **Cloud APIs** - Fallback only

## Go Integration

```go
import "github.com/ollama/ollama/api"

func SolveWithOllama(ctx context.Context, problem string) (string, error) {
    client, err := api.ClientFromEnvironment()
    if err != nil {
        return "", err
    }

    req := &api.GenerateRequest{
        Model:  "math-coder",
        Prompt: "Solve step by step: " + problem,
        Stream: &[]bool{false}[0],
    }

    var resp api.GenerateResponse
    err = client.Generate(ctx, req, &resp)
    return resp.Response, err
}
```

## Embeddings

```go
func GetEmbedding(ctx context.Context, text string) ([]float64, error) {
    client, _ := api.ClientFromEnvironment()
    
    req := &api.EmbedRequest{
        Model:  "nomic-embed-text",
        Input:  text,
    }
    
    var resp api.EmbedResponse
    err := client.Embed(ctx, req, &resp)
    return resp.Embeddings[0], err
}
```

## Available Models

| Model | Purpose | Size |
|-------|---------|------|
| `llama3` | General | 4.7GB |
| `math-coder` | Math solving | 4.7GB |
| `codellama` | Code | 3.8GB |
| `nomic-embed-text` | Embeddings | 274MB |

## Commands

```bash
# Run Ollama
ollama serve

# Pull model
ollama pull llama3
ollama pull math-coder
ollama pull nomic-embed-text

# List models
ollama list

# Run interactively
ollama run llama3 "Explain derivatives"
```

## Environment

| Variable | Default |
|----------|---------|
| `OLLAMA_HOST` | localhost:11434 |
| `OLLAMA_MODEL` | llama3 |

## Configuration

```yaml
# For math problems
models:
  - name: math-coder
    temperature: 0.1
    top_p: 0.9
```

## Error Handling

```go
if errors.Is(err, context.DeadlineExceeded) {
    // Fallback to cloud
    return cloudSolve(ctx, problem)
}

if strings.Contains(err.Error(), "not found") {
    // Pull model
    return pullAndRetry(ctx, model, problem)
}
```

## Anti-Patterns

- ❌ Using cloud APIs before trying Ollama
- ❌ Not handling model loading errors
- ❌ Blocking on slow inference
- ❌ Not using context cancellation
