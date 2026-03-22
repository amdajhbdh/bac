// WhatsApp messaging integration using ZeroClaw HTTP API
export class WhatsAppClient {
  private zeroclawUrl: string;
  private phoneNumberId?: string;

  constructor(zeroclawUrl: string, phoneNumberId?: string) {
    this.zeroclawUrl = zeroclawUrl;
    this.phoneNumberId = phoneNumberId;
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

  async sendMessage(to: string, message: string): Promise<void> {
    if (!this.phoneNumberId) {
      throw new Error('Phone number ID not configured');
    }
    await this.request('/zeroclaw/send', 'POST', {
      message: message,
      channel_id: 'whatsapp',
      recipient: to,
    });
  }

  async markAsRead(messageId: string): Promise<void> {
    // For marking as read, we'd need a specific endpoint
    // For now, we'll just acknowledge via ZeroClaw's webhook system
    // This is a simplified implementation
    console.log(`Marking message ${messageId} as read (via ZeroClaw)`);
  }
}