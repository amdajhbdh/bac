package cache

import (
	"context"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var logger = slog.New(slog.NewJSONHandler(nil, nil))

const (
	DefaultDimension = 1536
	DefaultTopK      = 10
)

type Vector []float32

func (v Vector) Dim() int {
	return len(v)
}

func (v Vector) CosineSim(u Vector) float64 {
	if len(v) != len(u) {
		return 0
	}

	dot := float64(0)
	normV := float64(0)
	normU := float64(0)

	for i := range v {
		dot += float64(v[i]) * float64(u[i])
		normV += float64(v[i]) * float64(v[i])
		normU += float64(u[i]) * float64(u[i])
	}

	if normV == 0 || normU == 0 {
		return 0
	}

	return dot / (math.Sqrt(normV) * math.Sqrt(normU))
}

type SearchResult struct {
	ID         string  `json:"id"`
	Similarity float64 `json:"similarity"`
	Problem    string  `json:"problem"`
	Solution   string  `json:"solution"`
}

type USearchCache struct {
	dimension int
	maxSize   int
	vectors   map[string]Vector
	problems  map[string]SearchResult
	lru       []string
	mu        sync.RWMutex

	pool  *pgxpool.Pool
	redis *redis.Client
}

type Config struct {
	Dimension   int
	MaxSize     int
	RedisAddr   string
	PostgresURL string
}

func DefaultConfig() Config {
	return Config{
		Dimension: DefaultDimension,
		MaxSize:   100000,
		RedisAddr: "localhost:6379",
	}
}

func NewCache(cfg Config) (*USearchCache, error) {
	cache := &USearchCache{
		dimension: cfg.Dimension,
		maxSize:   cfg.MaxSize,
		vectors:   make(map[string]Vector),
		problems:  make(map[string]SearchResult),
		lru:       make([]string, 0, cfg.MaxSize),
	}

	if cfg.RedisAddr != "" {
		cache.redis = redis.NewClient(&redis.Options{
			Addr: cfg.RedisAddr,
		})
	}

	if cfg.PostgresURL != "" {
		pool, err := pgxpool.New(context.Background(), cfg.PostgresURL)
		if err != nil {
			logger.Error("postgres connection failed", "error", err)
		} else {
			cache.pool = pool
		}
	}

	logger.Info("USearch cache initialized",
		"dimension", cfg.Dimension,
		"max_size", cfg.MaxSize,
	)

	return cache, nil
}

func (c *USearchCache) Add(ctx context.Context, id string, vector Vector, problem, solution string) error {
	if len(vector) != c.dimension {
		return &CacheError{Op: "add", Reason: "invalid dimension"}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.vectors[id]; !exists {
		if len(c.vectors) >= c.maxSize {
			c.evictOne()
		}
		c.lru = append(c.lru, id)
	}

	c.vectors[id] = vector
	c.problems[id] = SearchResult{
		ID:       id,
		Problem:  problem,
		Solution: solution,
	}

	return nil
}

func (c *USearchCache) Search(ctx context.Context, query Vector, topK int) ([]SearchResult, error) {
	if len(query) != c.dimension {
		return nil, &CacheError{Op: "search", Reason: "invalid query dimension"}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	type scoredResult struct {
		id    string
		score float64
	}

	results := make([]scoredResult, 0, len(c.vectors))

	for id, vec := range c.vectors {
		sim := Vector(query).CosineSim(vec)
		results = append(results, scoredResult{id: id, score: sim})
	}

	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].score > results[i].score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	if topK > len(results) {
		topK = len(results)
	}

	output := make([]SearchResult, topK)
	for i := 0; i < topK; i++ {
		pr := c.problems[results[i].id]
		pr.Similarity = results[i].score
		output[i] = pr
	}

	return output, nil
}

func (c *USearchCache) evictOne() {
	if len(c.lru) == 0 {
		return
	}

	oldest := c.lru[0]
	c.lru = c.lru[1:]
	delete(c.vectors, oldest)
	delete(c.problems, oldest)
}

func (c *USearchCache) Warm(ctx context.Context) error {
	if c.pool == nil {
		return &CacheError{Op: "warm", Reason: "no postgres connection"}
	}

	logger.Info("warming cache from postgres...")

	rows, err := c.pool.Query(ctx, `
		SELECT id, question_text, solution_text
		FROM questions
		WHERE question_vector IS NOT NULL
		LIMIT 10000
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var problem, solution string

		if err := rows.Scan(&id, &problem, &solution); err != nil {
			continue
		}

		vector := generateHashVector(problem + solution)
		c.Add(ctx, string(rune(id)), vector, problem, solution)
		count++
	}

	logger.Info("cache warmed", "count", count)
	return nil
}

func generateHashVector(text string) Vector {
	emb := make(Vector, DefaultDimension)

	h1 := hashString(text)
	h2 := hashString(text + "salt1")
	h3 := hashString(text + "salt2")

	for i := 0; i < DefaultDimension; i++ {
		bitPos := i % 64
		bit1 := (h1 >> uint(bitPos)) & 1
		bit2 := (h2 >> uint(bitPos)) & 1
		bit3 := (h3 >> uint(bitPos)) & 1

		emb[i] = float32(bit1)*0.5 + float32(bit2)*0.3 + float32(bit3)*0.2
		if emb[i] == 0 {
			emb[i] = 0.1
		}
	}

	return emb
}

func hashString(s string) uint64 {
	var h uint64 = 14695981039346656037
	for _, c := range s {
		h ^= uint64(c)
		h *= 1099511628211
	}
	return h
}

func (c *USearchCache) Stats() USearchStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return USearchStats{
		VectorCount: len(c.vectors),
		MaxSize:     c.maxSize,
		Dimension:   c.dimension,
	}
}

type USearchStats struct {
	VectorCount int
	MaxSize     int
	Dimension   int
}

type CacheError struct {
	Op     string
	Reason string
}

func (e *CacheError) Error() string {
	return "cache error: " + e.Op + " - " + e.Reason
}

func NewCacheError(op, reason string) *CacheError {
	return &CacheError{Op: op, Reason: reason}
}

func (c *USearchCache) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
	if c.redis != nil {
		c.redis.Close()
	}
}

type MultiTierCache struct {
	hot   *USearchCache
	redis *redis.Client
	pg    *pgxpool.Pool
}

func NewMultiTierCache(cfg Config) (*MultiTierCache, error) {
	hot, err := NewCache(cfg)
	if err != nil {
		return nil, err
	}

	mtc := &MultiTierCache{
		hot: hot,
	}

	if cfg.RedisAddr != "" {
		mtc.redis = redis.NewClient(&redis.Options{
			Addr: cfg.RedisAddr,
		})
	}

	if cfg.PostgresURL != "" {
		pool, err := pgxpool.New(context.Background(), cfg.PostgresURL)
		if err == nil {
			mtc.pg = pool
		}
	}

	return mtc, nil
}

func (m *MultiTierCache) Search(ctx context.Context, query Vector, topK int) ([]SearchResult, error) {
	results, err := m.hot.Search(ctx, query, topK)
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return results, nil
	}

	if m.redis != nil {
		key := "vector:" + string(rune(time.Now().UnixNano()))
		m.redis.Get(ctx, key)
	}

	return results, nil
}

func (m *MultiTierCache) Warm(ctx context.Context) error {
	return m.hot.Warm(ctx)
}

func (m *MultiTierCache) Close() {
	m.hot.Close()
	if m.redis != nil {
		m.redis.Close()
	}
	if m.pg != nil {
		m.pg.Close()
	}
}
