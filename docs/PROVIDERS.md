# BAC Agent Provider Configuration

This document describes how to configure AI model providers for the BAC Agent.

## Overview

The BAC Agent supports **multiple AI providers** including free tiers. This guide covers configuration for all supported providers.

## Supported Providers

### Free Providers (Recommended)

| Provider | Key | Free Tier | Best For |
|----------|-----|-----------|----------|
| **OpenRouter** | `openrouter` | $5 credits on signup | 100+ models |
| **Groq** | `groq` | Free tier | Fast inference |
| **Together AI** | `together` | $25 credits | Llama, Qwen |
| **Qwen/DashScope** | `qwen` | Free | Chinese models |
| **GLM/Zhipu** | `glm` | Free | Chinese models |
| **Moonshot/Kimi** | `kimi` | Free credits | Long context |
| **MiniMax** | `minimax` | Free | Multi-modal |
| **Google Gemini** | `gemini` | 500 req/day | Gemini 2.5 |
| **Ollama** | `ollama` | Local | Privacy, no cost |

### Premium Providers

| Provider | Key | Notes |
|----------|-----|-------|
| OpenAI | `openai` | GPT-4o, GPT-5 |
| Anthropic | `anthropic` | Claude 4 |
| OpenAI Codex | `openai-codex` | Code-specific |

## Environment Variables

Set these in your `.env` file:

```bash
# Free Providers
export OPENROUTER_API_KEY="sk-or-v1-..."
export GROQ_API_KEY="gsk_..."
export TOGETHER_API_KEY="..."
export DASHSCOPE_API_KEY="..."        # Qwen
export ZHIPU_API_KEY="..."            # GLM
export KIMI_API_KEY="..."            # Moonshot
export MINIMAX_API_KEY="..."

# Google (Gemini)
export GEMINI_API_KEY="..."

# Premium
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."

# Local (Ollama)
export OLLAMA_HOST="http://localhost:11434"
```

## Configuration File

Edit `crates/agent/config.yaml`:

### Using OpenRouter (Recommended for Free)

```yaml
providers:
  openrouter:
    api_key: ${OPENROUTER_API_KEY}
    default_model: deepseek/deepseek-r1
    # Or: google/gemini-2.0-flash, anthropic/claude-3-haiku, etc.

# Set default provider
agent:
  provider: openrouter
```

### Using Groq (Fast Free)

```yaml
providers:
  groq:
    api_key: ${GROQ_API_KEY}
    default_model: llama-3.3-70b-versatile

agent:
  provider: groq
```

### Using Qwen (Free Chinese Models)

```yaml
providers:
  qwen:
    api_key: ${DASHSCOPE_API_KEY}
    default_model: qwen-plus

agent:
  provider: qwen
```

### Using Ollama (Local)

```yaml
providers:
  ollama:
    host: http://localhost:11434
    default_model: qwen2.5:72b

agent:
  provider: ollama
```

### Using Multiple Providers (Fallback)

```yaml
providers:
  primary:
    openrouter:
      api_key: ${OPENROUTER_API_KEY}
      default_model: deepseek/deepseek-r1
  fallback:
    groq:
      api_key: ${GROQ_API_KEY}
      default_model: llama-3.3-70b-versatile

agent:
  provider: primary
  fallback: fallback
```

## Model Recommendations

### Coding

| Model | Provider | Notes |
|-------|---------|-------|
| DeepSeek R1 | OpenRouter | Best for reasoning |
| Qwen Coder | Qwen | Great for code |
| Claude 4 | Anthropic | Premium option |

### Chinese Content

| Model | Provider | Notes |
|-------|---------|-------|
| Qwen3 | Qwen | Free, excellent |
| GLM-5 | GLM | Free |
| Kimi | Moonshot | Long context |

### Fast/Development

| Model | Provider | Notes |
|-------|---------|-------|
| Llama 3.3 70B | Groq | Very fast |
| Gemini 2.0 Flash | Google | Free, 500/day |
| Qwen Turbo | Qwen | Free |

## Verification

Verify your configuration:

```bash
# List available models for a provider
./crates/agent/target/release/bac-agent providers

# Test a specific provider
./crates/agent/target/release/bac-agent agent -m "Hello"
```

## Troubleshooting

### "API key not found"

Make sure the environment variable is set:
```bash
export OPENROUTER_API_KEY="your-key"
```

### "Rate limit exceeded"

Switch to a different provider or wait. Groq and Gemini have generous free tiers.

### "Connection timeout"

Check your network or use a local provider like Ollama.

## API Endpoints

| Service | Port | Purpose |
|---------|------|---------|
| agent | 8081 | AI agent (gateway) |
| api | 8080 | REST API |
| ocr | 8082 | Document OCR |
| streaming | 8083 | Real-time streaming |
| web | 3000 | Frontend |
| postgres | 5432 | Database |

## Quick Start

1. Get API key from your preferred provider
2. Add to `.env` file
3. Update `config.yaml`
4. Restart agent: `./scripts/bac-agent-daemon.sh restart`
5. Test: `./scripts/bac-agent-send.sh message "Hello!"`
