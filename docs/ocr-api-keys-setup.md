# OCR API Keys Setup Guide

This guide covers all 18 OCR providers in the BAC OCR chain. For each provider: acquisition steps, free tier limits, and current configuration status.

## Quick Status

| Provider | Status | Free Tier | Cost/Call |
|---|---|---|---|
| Mistral OCR | ✅ Configured | No | $0.002 |
| OCR.space | ✅ Works (no key) | 500/day | Free |
| Cloudflare Workers AI | ✅ Via Worker | 10K/day | Free |
| Yandex Vision | ❌ Not configured | 1K/3mo | Free |
| LFM2.5-VL (Ollama) | ⏳ Cloud Shell pending | Unlimited | Free |
| Qwen2.5-VL (Ollama) | ⏳ Cloud Shell pending | Unlimited | Free |
| LLaVA (Ollama) | ⏳ Cloud Shell pending | Unlimited | Free |
| PaddleOCR | ⏳ Not installed | Unlimited | Free |
| EasyOCR | ⏳ Not installed | Unlimited | Free |
| Tesseract | ✅ Available | Unlimited | Free |
| GLM-OCR | ⏳ Not installed | - | - |
| Cloudmersive | ❌ Not configured | 800/mo | Free |
| ABBYY Cloud | ❌ Not configured | 500/mo | Free |
| Base64.ai | ❌ Not configured | Trial | $0.01 |
| OCRAPI.cloud | ❌ Not configured | 250/mo | Free |
| Nanonets | ❌ Not configured | $200 credits | $0.05 |
| Mathpix | ❌ Not configured | 10/mo | Free |
| Puter.js | ✅ Available | - | Free |

---

## Cloud Providers

### 1. Mistral OCR (Primary) ⭐
**Status:** Configured in `.env` and Cloudflare Worker secrets.

**Cost:** $0.002 per page (1000 pages = $2)

**Signup:** https://console.mistral.ai/

**API Key:** `BziODszFgZvBkRJOkAePgDcgegEqbffe` (already set)

**Best for:** French, Arabic, math equations, high accuracy.

---

### 2. OCR.space
**Status:** Works without API key (free `helloworld` tier).

**Limits:** 500 requests/day, 500KB per file.

**Signup:** https://ocr.space/ocrapi (optional, for higher limits)

**No key needed** — uses free `helloworld` API key internally.

**Best for:** Quick fallback, free tier, no setup.

---

### 3. Cloudflare Workers AI (Llama Vision)
**Status:** Integrated via Cloudflare Worker endpoint.

**Limits:** ~10K requests/day (Worker invocations).

**No key needed** — uses Cloudflare's built-in Workers AI.

**Endpoint:** `https://bac-api.amdajhbdh.workers.dev/rag/ocr`

**Best for:** Fast, free, integrated with the RAG pipeline.

---

### 4. Yandex Vision OCR
**Status:** Not configured — needs IAM token and folder ID.

**Limits:** 1,000 operations/3 months (free trial).

**Signup:** https://console.cloud.yandex.com/

**Steps:**
1. Create a Yandex Cloud account
2. Create a folder (e.g., `bac-ocr`)
3. Enable Yandex Vision API
4. Create a service account
5. Generate an IAM token

**Environment variables:**
```bash
YANDEX_IAM_TOKEN="t1.9eu3..."
YANDEX_FOLDER_ID="b1g..."
```

**Set in Worker secrets:**
```bash
wrangler secret put YANDEX_IAM_TOKEN
wrangler secret put YANDEX_FOLDER_ID
```

**Best for:** Russian/Arabic text, good free tier.

---

### 5. Cloudmersive
**Status:** Not configured.

**Limits:** 800 API calls/month free.

**Signup:** https://cloudmersive.com/ (free tier requires signup)

**API Key:** Get from dashboard after signup.

**Environment variable:**
```bash
CLOUDMERSIVE_API_KEY="your-key"
```

---

### 6. ABBYY Cloud OCR
**Status:** Not configured.

**Limits:** 500 pages/month free.

**Signup:** https://www.abbyy.com/cloud-ocr-sdk/

**Credentials:**
```bash
ABBYY_APP_ID="your-app-id"
ABBYY_PASSWORD="your-password"
```

---

### 7. Base64.ai
**Status:** Not configured.

**Limits:** Trial only (paid: $0.01/page).

**Signup:** https://base64.ai/

**API Key:** Get from dashboard.

```bash
BASE64_API_KEY="your-key"
```

---

### 8. OCRAPI.cloud
**Status:** Not configured.

**Limits:** 250 requests/month free.

**Signup:** https://ocrapi.cloud/

```bash
OCRAPI_KEY="your-key"
```

---

### 9. Nanonets
**Status:** Not configured.

**Limits:** $200 free credits.

**Signup:** https://nanonets.com/

```bash
NANONETS_API_KEY="your-key"
```

---

### 10. Mathpix
**Status:** Not configured.

**Limits:** 10 pages/month free (SNIPPED plan).

**Signup:** https://mathpix.com/

**Best for:** Math equations with LaTeX output.

```bash
MATHPIX_APP_ID="your-app-id"
MATHPIX_APP_KEY="your-app-key"
```

---

## Cloud Shell Providers (via Ollama)

These require setting up the Cloud Shell with Ollama. See `scripts/cloud-shell-ocr.sh`.

### 11. LFM2.5-VL (Liquid AI) ⭐
**Model:** `LiquidAI/LFM2.5-VL-1.6B-GGUF` (Q4_0 = 696MB)

**OCRBench:** 822 (highest in chain)

**Setup:**
```bash
# In Cloud Shell:
bash scripts/cloud-shell-ocr.sh
```

**Best for:** French + Arabic text, best open-source accuracy.

---

### 12. Qwen2.5-VL
**Model:** `qwen2.5vl:latest` via Ollama

**Setup:** Pulled automatically by `cloud-shell-ocr.sh`

**Best for:** General-purpose OCR, multilingual.

---

### 13. LLaVA
**Model:** `llava:latest` via Ollama

**Setup:** Pulled automatically by `cloud-shell-ocr.sh`

**Best for:** Fallback vision model, stable.

---

### 14. GLM-OCR (Zhipu AI)
**Model:** `glm-ocr:latest` via Ollama

**Setup:** Pulled by `cloud-shell-ocr.sh`

**Best for:** Math equations, Chinese text.

---

## Local Python Providers

### 15. Tesseract OCR
**Status:** ✅ Likely available (check with `tesseract --version`)

**Install:**
```bash
# Linux/Cloud Shell:
sudo apt install tesseract-ocr tesseract-ocr-fra tesseract-ocr-ara tesseract-ocr-eng

# Termux:
pkg in tesseract tesseract-lang
```

**Best for:** Clean printed documents, French/Arabic/English.

---

### 16. PaddleOCR
**Status:** Not installed — needs setup.

**Install:**
```bash
pip install paddlepaddle paddleocr
# or for GPU:
pip install paddlepaddle-gpu paddleocr
```

**Best for:** Best open-source accuracy, handles complex layouts.

---

### 17. EasyOCR
**Status:** Not installed — needs setup.

**Install:**
```bash
pip install easyocr
# Downloads ~1.7GB model on first run
```

**Best for:** Deep learning OCR, 80+ languages including Arabic.

---

### 18. Puter.js (Browser OCR)
**Status:** ✅ Available (runs in browser/CLI).

**No install needed** — uses `puter-js` npm package.

**Best for:** Web deployment, no server needed.

---

## Cloudflare Tunnel Setup

To expose Cloud Shell Ollama for remote access:

```bash
# In Cloud Shell:
bash scripts/cloudflared-tunnel.sh

# This outputs a public URL like:
# https://xxxx.trycloudflare.com
```

Set this URL in `ocr-config.toml` as `cloud_shell_url`.

---

## Environment Variables

All keys should be set in `.env` (gitignored):

```bash
# Core OCR
MISTRAL_API_KEY="BziODszFgZvBkRJOkAePgDcgegEqbffe"

# Optional
YANDEX_IAM_TOKEN=""
YANDEX_FOLDER_ID=""
CLOUDMERSIVE_API_KEY=""
ABBYY_APP_ID=""
ABBYY_PASSWORD=""
BASE64_API_KEY=""
OCRAPI_KEY=""
NANONETS_API_KEY=""
MATHPIX_APP_ID=""
MATHPIX_APP_KEY=""
```

For Cloudflare Worker, set secrets with:
```bash
wrangler secret put MISTRAL_API_KEY
wrangler secret put YANDEX_IAM_TOKEN
wrangler secret put YANDEX_FOLDER_ID
```

---

## Security Best Practices

1. **Never commit API keys** — `.env` is gitignored
2. **Use `.env.example`** for template, never real keys
3. **Rotate keys regularly** — especially paid tier keys
4. **Use least-privilege** — create dedicated API keys per service
5. **Monitor usage** — set up billing alerts on cloud consoles
6. **Cloudflare Worker secrets** — for keys used in deployed Worker code
