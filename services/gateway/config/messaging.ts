// Auto-generated stub for messaging service configuration
// Reuses existing configuration patterns from src/ and adapters/

export const messagingConfig = {
  telegram: {
    botToken: process.env.TELEGRAM_BOT_TOKEN || '',
    defaultChatId: process.env.TELEGRAM_DEFAULT_CHAT_ID || '',
  },
  whatsapp: {
    accessToken: process.env.WHATSAPP_ACCESS_TOKEN || '',
    phoneNumberId: process.env.WHATSAPP_PHONE_NUMBER_ID || '',
    verifyToken: process.env.WHATSAPP_VERIFY_TOKEN || '',
  },
  // Reuse existing rate limiting pattern from src/security.ts
  rateLimit: {
    windowMs: 60 * 1000, // 1 minute
    maxMessages: 30, // limit each IP to 30 requests per windowMs
  },
};