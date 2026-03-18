# Messaging Service Integration

BAC Unified's messaging system provides real‑time communication via Telegram and WhatsApp, with webhook routing, rate limiting, and scheduled automation.

## Architecture

- **Adapters**: `src/messaging/telegram.ts` & `src/messaging/whatsapp.ts`
- **Middleware**: `src/api-gateway/src/middleware/messaging.ts`
- **Configuration**: `services/gateway/config/messaging.ts`
- **Automation**: `scripts/bac-automation.sh` & `scripts/sync-obsidian.js`
- **CI/CD**: `.github/workflows/bac-study.yml`

## Setup

1. Add environment variables to `.env`:
   ```bash
   TELEGRAM_BOT_TOKEN=your_bot_token
   TELEGRAM_DEFAULT_CHAT_ID=your_chat_id
   WHATSAPP_ACCESS_TOKEN=your_access_token
   WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id
   ```

2. Install dependencies:
   ```bash
   cd .opencode && npm ci
   ```

3. Run tests:
   ```bash
   npm test -- --scope messaging
   ```

## Usage

### Telegram
```typescript
import { TelegramBot } from 'src/messaging/telegram';
const bot = new TelegramBot(process.env.TELEGRAM_BOT_TOKEN, process.env.TELEGRAM_DEFAULT_CHAT_ID);
await bot.sendMessage('Hello from BAC!');
```

### WhatsApp
```typescript
import { WhatsAppClient } from 'src/messaging/whatsapp';
const client = new WhatsAppClient(process.env.WHATSAPP_ACCESS_TOKEN, process.env.WHATSAPP_PHONE_NUMBER_ID);
await client.sendMessage('+1234567890', 'Hello from BAC!');
```

## Automation

- **Local**: `./scripts/bac-automation.sh` (run manually or via cron)
- **CI**: Push to `master` or modify `resources/notes/**` triggers the BAC Study workflow
- **Obsidian Sync**: `./scripts/sync-obsidian.js` (runs in CI, can be scheduled locally)

## Verification

```bash
# Unit tests
npm test -- tests/messaging

# Integration tests
npm test -- tests/integration

# Smoke test
npm run test:messaging
```

## Key Design Decisions

- Reused existing security and lifecycle patterns from `src/security.ts` and `src/lifecycle.ts`
- Middleware uses in-memory rate limiting (suitable for single‑instance deployment)
- All messaging adapters expose minimal, consistent interfaces
- Configuration centralised under `services/gateway/config/` for environment-based overrides