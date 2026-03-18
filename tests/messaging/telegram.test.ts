// Auto-generated unit tests for Telegram messaging

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { TelegramBot } from '../../src/messaging/telegram';

describe('TelegramBot', () => {
  let bot: TelegramBot;

  beforeEach(() => {
    bot = new TelegramBot('test-token', '123456');
    vi.stubGlobal('fetch', vi.fn());
  });

  it('should send message successfully', async () => {
    const mockFetch = vi.mocked(fetch);
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ ok: true, result: {} }),
    } as Response);

    await bot.sendMessage('Hello, world!');

    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('sendMessage'),
      expect.objectContaining({
        body: expect.stringContaining('chat_id=123456'),
      })
    );
  });

  it('should throw error if chat ID is missing', async () => {
    const badBot = new TelegramBot('test-token');

    await expect(badBot.sendMessage('Hello')).rejects.toThrow('Chat ID not configured');
  });
});