# OpenCrust Integration for BAC Unified
# Use pre-built binary + simple scripts

## Quick Start

### 1. Install OpenCrust
```bash
curl -fsSL https://raw.githubusercontent.com/opencrust-org/opencrust/main/install.sh | sh
```

### 2. Configure
```bash
cp deploy/opencrust/config.toml ~/.config/opencrust/config.toml
# Edit with your Telegram/WhatsApp tokens
```

### 3. Run
```bash
opencrust start
```

## Commands Available (Telegram/WhatsApp)
- `vault sync` - Pull/push vault
- `vault stats` - Show statistics
- `vault search <query>` - Search notes
- `ocr <file>` - Process PDF
- `help` - Show all commands

## Configuration
Edit `deploy/opencrust/config.toml`:
```toml
[model]
provider = "openai"
model = "llama3.2:3b"
base_url = "http://localhost:11434/v1"

[channels.telegram]
enabled = true

[channels.whatsapp]
enabled = true
```
