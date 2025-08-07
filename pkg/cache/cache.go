package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/pkg/logger"
)

// Cache 缓存接口
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Health() error
	Close() error
}

// redisCache Redis 缓存实现
type redisCache struct {
	client *redis.Client
	logger logger.Logger
}

// memoryCache 内存缓存实现
type memoryCache struct {
	data   map[string]cacheItem
	mutex  sync.RWMutex
	logger logger.Logger
}

type cacheItem struct {
	value      string
	expiration time.Time
}

// New 创建新的缓存实例
func New(cfg *config.Config, log logger.Logger) (Cache, error) {
	switch cfg.Cache.Driver {
	case "redis":
		return newRedisCache(cfg, log)
	case "memory":
		return newMemoryCache(log), nil
	default:
		return nil, fmt.Errorf("unsupported cache driver: %s", cfg.Cache.Driver)
	}
}

// newRedisCache 创建 Redis 缓存实例
func newRedisCache(cfg *config.Config, log logger.Logger) (Cache, error) {
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Cache.GetRedisAddress(),
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.Database,
		PoolSize: cfg.Cache.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Redis connected successfully",
		"host", cfg.Cache.Host,
		"port", cfg.Cache.Port,
		"database", cfg.Cache.Database,
	)

	return &redisCache{
		client: rdb,
		logger: log,
	}, nil
}

// newMemoryCache 创建内存缓存实例
func newMemoryCache(log logger.Logger) Cache {
	log.Info("Using in-memory cache for testing")
	return &memoryCache{
		data:   make(map[string]cacheItem),
		logger: log,
	}
}

// Set 设置缓存
func (r *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		r.logger.Error("Failed to set cache", "key", key, "error", err)
		return err
	}
	return nil
}

// Get 获取缓存
func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key not found: %s", key)
		}
		r.logger.Error("Failed to get cache", "key", key, "error", err)
		return "", err
	}
	return val, nil
}

// Del 删除缓存
func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		r.logger.Error("Failed to delete cache", "keys", keys, "error", err)
		return err
	}
	return nil
}

// Exists 检查键是否存在
func (r *redisCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	count, err := r.client.Exists(ctx, keys...).Result()
	if err != nil {
		r.logger.Error("Failed to check cache existence", "keys", keys, "error", err)
		return 0, err
	}
	return count, nil
}

// Expire 设置过期时间
func (r *redisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		r.logger.Error("Failed to set cache expiration", "key", key, "error", err)
		return err
	}
	return nil
}

// TTL 获取剩余过期时间
func (r *redisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to get cache TTL", "key", key, "error", err)
		return 0, err
	}
	return ttl, nil
}

// Health 检查缓存健康状态
func (r *redisCache) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Ping(ctx).Err()
}

// Close 关闭缓存连接
func (r *redisCache) Close() error {
	return r.client.Close()
}

// Memory cache implementation methods

// Set 设置缓存
func (m *memoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var exp time.Time
	if expiration > 0 {
		exp = time.Now().Add(expiration)
	}

	m.data[key] = cacheItem{
		value:      fmt.Sprintf("%v", value),
		expiration: exp,
	}
	return nil
}

// Get 获取缓存
func (m *memoryCache) Get(ctx context.Context, key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return "", fmt.Errorf("key not found")
	}

	// 检查是否过期
	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		delete(m.data, key)
		return "", fmt.Errorf("key expired")
	}

	return item.value, nil
}

// Del 删除缓存
func (m *memoryCache) Del(ctx context.Context, keys ...string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}

// Exists 检查缓存是否存在
func (m *memoryCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var count int64
	for _, key := range keys {
		if item, exists := m.data[key]; exists {
			// 检查是否过期
			if item.expiration.IsZero() || time.Now().Before(item.expiration) {
				count++
			}
		}
	}
	return count, nil
}

// Expire 设置过期时间
func (m *memoryCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if item, exists := m.data[key]; exists {
		item.expiration = time.Now().Add(expiration)
		m.data[key] = item
	}
	return nil
}

// TTL 获取剩余过期时间
func (m *memoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return -2 * time.Second, nil // Key does not exist
	}

	if item.expiration.IsZero() {
		return -1 * time.Second, nil // Key exists but has no expiration
	}

	ttl := time.Until(item.expiration)
	if ttl <= 0 {
		return -2 * time.Second, nil // Key has expired
	}

	return ttl, nil
}

// Health 健康检查
func (m *memoryCache) Health() error {
	return nil // Memory cache is always healthy
}

// Close 关闭缓存连接
func (m *memoryCache) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data = make(map[string]cacheItem)
	return nil
}
