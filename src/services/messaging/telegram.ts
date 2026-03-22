// Telegram messaging integration using ZeroClaw HTTP API
export class TelegramBot {
  private zeroclawUrl: string;
  private chatId?: string;

  constructor(zeroclawUrl: string, chatId?: string) {
    this.zeroclawUrl = zeroclawUrl;
    this.chatId = chatId;
  }

  private async request<T>(endpoint: string, method: string = 'POST', body?: Record<string, any>): Promise<T> {
    const url = `${this.zeroclawUrl}${endpoint}`;
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    const response = await fetch(url, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!response.ok) {
      throw new Error(`ZeroClaw API error: ${response.status}`);
    }
    return response.json() as Promise<T>;
  }

  async sendMessage(text: string): Promise<void> {
    if (!this.chatId) {
      throw new Error('Chat ID not configured');
    }
    await this.request('/zeroclaw/send', 'POST', {
      message: text,
      channel_id: 'telegram',
      recipient: this.chatId,
    });
  }

  async getUpdates(): Promise<any> {
    // For getting updates, we'd typically use webhooks, but for compatibility:
    return this.request('/zeroclaw/channels/telegram/updates', 'GET');
  }
}