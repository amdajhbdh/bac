// Auto-generated unit tests for WhatsApp messaging

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { WhatsAppClient } from '../../src/messaging/whatsapp';

describe('WhatsAppClient', () => {
  let client: WhatsAppClient;

  beforeEach(() => {
    client = new WhatsAppClient('test-token', '123456');
    vi.stubGlobal('fetch', vi.fn());
  });

  it('should send message successfully', async () => {
    const mockFetch = vi.mocked(fetch);
    mockFetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ messages: [{ id: 'msg-123' }] }),
    } as Response);

    await client.sendMessage('+1234567890', 'Hello!');

    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('graph.facebook.com'),
      expect.objectContaining({
        method: 'POST',
        body: expect.stringContaining('Hello!'),
      })
    );
  });

  it('should throw error if phone number ID is missing', async () => {
    const badClient = new WhatsAppClient('test-token');

    await expect(badClient.sendMessage('+1234567890', 'Hello')).rejects.toThrow(
      'Phone number ID not configured'
    );
  });
});