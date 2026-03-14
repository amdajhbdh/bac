---
name: bac-ocr
description: AI-powered OCR pipeline for extracting text from PDFs using Tesseract + AI correction. Use when processing study materials, teacher notes, or exam papers. Triggers: ocr, extract, pdf, tesseract, scan, document
---

# BAC AI-OCR Pipeline Skill

Process PDFs with AI-powered OCR that corrects errors and structures content.

## Quick Commands

```bash
# Extract single PDF
./bin/ai-ocr extract path/to/file.pdf

# Batch process
./bin/ai-ocr batch 03-Resources/

# With cloud AI (for better results)
USE_CLOUD_MODEL=true CLOUD_MODEL=minimax-m2.5:cloud ./bin/ai-ocr extract file.pdf

# Verbose mode
./bin/ai-ocr --verbose extract file.pdf
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OLLAMA_URL` | http://localhost:11434 | Ollama API endpoint |
| `OLLAMA_MODEL` | llama3.2:3b | Local model |
| `USE_CLOUD_MODEL` | false | Use cloud AI |
| `CLOUD_MODEL` | minimax-m2.5:cloud | Cloud model |

## Available Cloud Models

- `minimax-m2.5:cloud` - Best for French/Arabic
- `deepseek-v3.2:cloud` - Good reasoning
- `glm-5:cloud` - Multilingual
- `kimi-k2.5:cloud` - Fast

## Output

- Extracted text in `05-Extracted/`
- Knowledge graphs in `knowledge-graphs/`
- Front matter with tags

## Tips

- Use cloud models for handwritten content
- Local Ollama for privacy/speed
- Verbose mode shows real-time progress
