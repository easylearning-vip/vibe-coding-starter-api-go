package cache

import (
	"context"
	"fmt"
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

// New 创建新的缓存实例
func New(cfg *config.Config, log logger.Logger) (Cache, error) {
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
