#!/usr/bin/env node

// Auto-generated messaging service server entry point
// Runs Express server for Telegram/WhatsApp webhooks

import express from 'express';
import cors from 'cors';
import { TelegramBot } from './telegram';
import { WhatsAppClient } from './whatsapp';
import { messagingConfig } from '../services/gateway/config/messaging';

const app = express();
app.use(cors());
app.use(express.json());

const telegramBot = new TelegramBot(
  messagingConfig.telegram.botToken,
  messagingConfig.telegram.defaultChatId
);
const whatsappClient = new WhatsAppClient(
  messagingConfig.whatsapp.accessToken,
  messagingConfig.whatsapp.phoneNumberId
);

// Health check
app.get('/health', (req, res) => {
  res.json({ service: 'messaging', status: 'ok' });
});

// Telegram webhook endpoint
app.post('/webhook/telegram', async (req, res) => {
  const { message, chat_id } = req.body;

  if (!message) {
    return res.status(400).json({ error: 'message is required' });
  }

  try {
    await telegramBot.sendMessage(message);
    res.json({ status: 'sent', platform: 'telegram' });
  } catch (err: any) {
    res.status(500).json({ error: err.message });
  }
});

// WhatsApp webhook endpoint
app.post('/webhook/whatsapp', async (req, res) => {
  const { message, to } = req.body;

  if (!message) {
    return res.status(400).json({ error: 'message is required' });
  }

  try {
    await whatsappClient.sendMessage(to || '', message);
    res.json({ status: 'sent', platform: 'whatsapp' });
  } catch (err: any) {
    res.status(500).json({ error: err.message });
  }
});

// Start server
const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
  console.log(`Messaging service listening on port ${PORT}`);
});