// Auto-generated integration tests for messaging middleware

import { describe, it, expect, beforeEach, vi } from 'vitest';
import express from 'express';
import request from 'supertest';
import { MessagingMiddleware } from '../src/api-gateway/src/middleware/messaging';

describe('MessagingMiddleware Integration', () => {
  let app: express.Express;
  let middleware: MessagingMiddleware;

  beforeEach(() => {
    app = express();
    middleware = new MessagingMiddleware();
    app.use(express.json());
    app.use((req, res, next) => middleware.logMessages(req, res, next));
    app.use((req, res, next) => middleware.rateLimit(req, res, next));

    app.post('/webhook', (req, res) => {
      res.json({ status: 'ok' });
    });

    vi.useFakeTimers();
  });

  it('should allow requests within rate limit', async () => {
    const response = await request(app).post('/webhook').send({ test: true });
    expect(response.status).toBe(200);
  });

  it('should rate limit excessive requests', async () => {
    // Send many requests quickly
    const promises = Array.from({ length: 40 }, () =>
      request(app).post('/webhook').send({ test: true })
    );
    const responses = await Promise.all(promises);

    const blocked = responses.filter(r => r.status === 429);
    expect(blockged).toBeGreaterThan(0);
  });
});