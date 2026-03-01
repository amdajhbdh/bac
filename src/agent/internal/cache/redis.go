package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

type SessionData struct {
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	Device    string    `json:"device"`
	IP        string    `json:"ip"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RateLimitData struct {
	Count   int       `json:"count"`
	ResetAt time.Time `json:"reset_at"`
}

func NewRedisCache() (*RedisCache, error) {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Warn("redis connection failed", "error", err)
		return nil, err
	}

	ttl := 24 * time.Hour
	if ttlStr := os.Getenv("REDIS_CACHE_TTL"); ttlStr != "" {
		if parsed, err := time.ParseDuration(ttlStr); err == nil {
			ttl = parsed
		}
	}

	slog.Info("redis cache connected", "url", url)
	return &RedisCache{client: client, ttl: ttl}, nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}
		data = string(bytes)
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *RedisCache) GetSession(ctx context.Context, token string) (*SessionData, error) {
	key := fmt.Sprintf("session:%s", token)
	data, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}

	var session SessionData
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	if session.ExpiresAt.Before(time.Now()) {
		r.Delete(ctx, key)
		return nil, nil
	}

	return &session, nil
}

func (r *RedisCache) SetSession(ctx context.Context, token string, session *SessionData) error {
	key := fmt.Sprintf("session:%s", token)
	ttl := session.ExpiresAt.Sub(time.Now())
	if ttl < 0 {
		ttl = 24 * time.Hour
	}
	return r.Set(ctx, key, session, ttl)
}

func (r *RedisCache) DeleteSession(ctx context.Context, token string) error {
	key := fmt.Sprintf("session:%s", token)
	return r.Delete(ctx, key)
}

func (r *RedisCache) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	rateKey := fmt.Sprintf("ratelimit:%s", key)

	data, err := r.Get(ctx, rateKey)
	if err != nil {
		return false, err
	}

	if data == "" {
		rateData := RateLimitData{
			Count:   1,
			ResetAt: time.Now().Add(window),
		}
		err := r.Set(ctx, rateKey, &rateData, window)
		return true, err
	}

	var rateData RateLimitData
	if err := json.Unmarshal([]byte(data), &rateData); err != nil {
		return false, fmt.Errorf("unmarshal rate data: %w", err)
	}

	if rateData.Count >= limit {
		return false, nil
	}

	rateData.Count++
	err = r.Set(ctx, rateKey, &rateData, time.Until(rateData.ResetAt))
	return err == nil, err
}

func (r *RedisCache) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	return r.client.Publish(ctx, channel, data).Err()
}

func (r *RedisCache) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.client.Subscribe(ctx, channel)
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

type CacheStats struct {
	Keys   int64
	Hits   int64
	Misses int64
}

func (r *RedisCache) GetStats(ctx context.Context) (*CacheStats, error) {
	keys, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	return &CacheStats{
		Keys: int64(len(keys)),
	}, nil
}
