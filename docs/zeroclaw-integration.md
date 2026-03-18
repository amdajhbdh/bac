# ZeroClaw Integration Guide

This document describes how ZeroClaw is integrated into the BAC Unified platform for autonomous messaging and automation.

## Overview

ZeroClaw is a Rust-based autonomous agent runtime that provides:
- **Messaging Channels**: Telegram, WhatsApp, Discord, Slack, Matrix, etc.
- **Gateway Server**: Webhook and WebSocket endpoints
- **Cron Scheduler**: Scheduled task execution
- **Memory System**: Persistent conversation memory
- **Provider Agnostic**: Works with OpenAI, Anthropic, OpenRouter, etc.

## Architecture

```
BAC Platform
    ├── Automation Scripts (bac-automation.sh)
    ├── ZeroClaw Integration Scripts
    │   ├── zeroclaw-build.sh
    │   ├── zeroclaw-send.sh
    │   └── zeroclaw-daemon.sh
    └── ZeroClaw Runtime
        ├── Gateway (webhooks/websockets)
        ├── Channels (Telegram, WhatsApp, etc.)
        ├── Scheduler (cron jobs)
        └── Memory (SQLite)
```

## Installation

### 1. Build ZeroClaw

```bash
# Build from source
./scripts/zeroclaw-build.sh

# Or install from cargo
cargo install zeroclawlabs
```

### 2. Configure ZeroClaw

```bash
# Run onboarding wizard
./src/zeroclaw/target/release/zeroclaw onboard

# Or configure manually
# Edit: ~/.zeroclaw/config.toml
```

### 3. Set up Channels

```bash
# Add Telegram channel
./src/zeroclaw/target/release/zeroclaw channel add telegram

# Add WhatsApp channel
./src/zeroclaw/target/release/zeroclaw channel add whatsapp
```

## Usage

### Sending Messages

Use the integration script:

```bash
# Send to all configured channels
./scripts/zeroclaw-send.sh message "🚀 BAC Automation started"

# Send to specific channel
./scripts/zeroclaw-send.sh telegram "Study reminder: Time for math!"
./scripts/zeroclaw-send.sh whatsapp "Your study session is ready!"
```

### Managing Daemon

```bash
# Start daemon
./scripts/zeroclaw-daemon.sh start

# Check status
./scripts/zeroclaw-daemon.sh status

# View logs
./scripts/zeroclaw-daemon.sh logs

# Stop daemon
./scripts/zeroclaw-daemon.sh stop
```

### Automation Integration

The `bac-automation.sh` script now uses ZeroClaw:

```bash
# Run automation (sends notifications via ZeroClaw)
./scripts/bac-automation.sh
```

## Configuration

### Environment Variables

Create/update `.env` file:

```bash
# ZeroClaw Provider
ZEROCLAW_PROVIDER=openrouter
ZEROCLAW_API_KEY=sk-...
ZEROCLAW_MODEL=openrouter/auto

# Telegram
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_DEFAULT_CHAT_ID=your_chat_id

# WhatsApp
WHATSAPP_ACCESS_TOKEN=your_access_token
WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id
WHATSAPP_DEFAULT_TO=recipient_phone_number
```

### ZeroClaw Config

ZeroClaw configuration is stored in `~/.zeroclaw/config.toml`:

```toml
[provider]
type = "openrouter"
api_key = "sk-..."
model = "openrouter/auto"

[channels.telegram]
enabled = true
bot_token = "your_bot_token"
default_chat_id = "your_chat_id"

[channels.whatsapp]
enabled = true
access_token = "your_access_token"
phone_number_id = "your_phone_number_id"
default_to = "recipient_phone_number"

[gateway]
host = "127.0.0.1"
port = 42617
require_pairing = false
```

## Service Deployment

### Docker Compose

```bash
# Start ZeroClaw service
docker-compose -f deploy/zeroclaw-service.yml up -d

# View logs
docker-compose -f deploy/zeroclaw-service.yml logs -f zeroclaw
```

### Systemd Service

```bash
# Install service
sudo cp deploy/zeroclaw.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable zeroclaw.service
sudo systemctl start zeroclaw.service

# Check status
sudo systemctl status zeroclaw.service
```

## Automation Examples

### Daily Study Briefing

```bash
# Cron job: 9 AM every weekday
0 9 * * 1-5 /home/user/bac/scripts/bac-automation.sh
```

### Obsidian Sync Notification

```bash
# After syncing notes
./scripts/zeroclaw-send.sh message "📝 Notes synced from Obsidian vault"
```

### Claude Code Analysis

```bash
# After code analysis
./scripts/zeroclaw-send.sh message "🔍 Code analysis complete. Check dashboard for details."
```

## Troubleshooting

### ZeroClaw not found

```bash
# Build zeroclaw
./scripts/zeroclaw-build.sh

# Or add to PATH
export PATH="$HOME/bac/src/zeroclaw/target/release:$PATH"
```

### Messages not sending

```bash
# Check zeroclaw status
./scripts/zeroclaw-daemon.sh status

# Check channel configuration
./src/zeroclaw/target/release/zeroclaw channel list

# Test message
./scripts/zeroclaw-send.sh message "Test message"
```

### Gateway not accessible

```bash
# Check if gateway is running
curl http://127.0.0.1:42617/health

# Restart daemon
./scripts/zeroclaw-daemon.sh restart
```

## Migration from TypeScript Wrappers

The previous TypeScript wrappers (`src/messaging/telegram.ts`, `src/messaging/whatsapp.ts`) have been archived to `src/messaging/archive/` because:

1. **ZeroClaw already implements** these channels in Rust
2. **Better performance** - Rust binary vs Node.js runtime
3. **Unified interface** - Single CLI for all channels
4. **Autonomous runtime** - Daemon mode with scheduler

To migrate existing code:

```bash
# Old: Direct API calls
curl -X POST "https://api.telegram.org/bot${TOKEN}/sendMessage" ...

# New: Use ZeroClaw
./scripts/zeroclaw-send.sh message "Your message"
```

## References

- [ZeroClaw GitHub](https://github.com/zeroclaw-labs/zeroclaw)
- [ZeroClaw Documentation](https://zeroclawlabs.ai/docs)
- [BAC Automation Script](../scripts/bac-automation.sh)
- [Integration Scripts](../scripts/zeroclaw-*.sh)
