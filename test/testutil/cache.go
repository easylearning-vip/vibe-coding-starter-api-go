package testutil

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"vibe-coding-starter/pkg/cache"
	testConfig "vibe-coding-starter/test/config"
)

// TestCache 测试缓存包装器
type TestCache struct {
	client    *redis.Client
	config    *testConfig.TestConfig
	useMemory bool
}

// NewTestCache 创建测试缓存连接
func NewTestCache(t *testing.T) *TestCache {
	config := testConfig.NewTestConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: config.Cache.Password,
		DB:       config.Cache.Database,
		PoolSize: config.Cache.PoolSize,
	})

	// 测试连接，如果失败则使用内存缓存
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	useMemory := false
	if err := client.Ping(ctx).Err(); err != nil {
		t.Logf("Failed to connect to test Redis: %v", err)
		t.Logf("Using in-memory cache for testing")
		useMemory = true
		client.Close() // 关闭失败的连接
		client = nil
	}

	return &TestCache{
		client:    client,
		config:    config,
		useMemory: useMemory,
	}
}

// Clean 清理测试缓存数据
func (tc *TestCache) Clean(t *testing.T) {
	if tc.useMemory {
		// 内存缓存不需要清理
		return
	}
	ctx := context.Background()
	if err := tc.client.FlushDB(ctx).Err(); err != nil {
		t.Logf("Warning: Failed to flush test cache: %v", err)
	}
}

// Close 关闭缓存连接
func (tc *TestCache) Close() error {
	if tc.useMemory {
		return nil
	}
	return tc.client.Close()
}

// CreateTestCache 创建实现cache.Cache接口的测试缓存
func (tc *TestCache) CreateTestCache() cache.Cache {
	if tc.useMemory {
		return &memoryCache{data: make(map[string]cacheItem)}
	}
	return &testCacheAdapter{tc}
}

// testCacheAdapter 适配器，实现cache.Cache接口
type testCacheAdapter struct {
	*TestCache
}

// Set 设置缓存
func (tca *testCacheAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return tca.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func (tca *testCacheAdapter) Get(ctx context.Context, key string) (string, error) {
	return tca.client.Get(ctx, key).Result()
}

// Del 删除缓存
func (tca *testCacheAdapter) Del(ctx context.Context, keys ...string) error {
	return tca.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (tca *testCacheAdapter) Exists(ctx context.Context, keys ...string) (int64, error) {
	return tca.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (tca *testCacheAdapter) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return tca.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (tca *testCacheAdapter) TTL(ctx context.Context, key string) (time.Duration, error) {
	return tca.client.TTL(ctx, key).Result()
}

// Incr 递增
func (tca *testCacheAdapter) Incr(ctx context.Context, key string) (int64, error) {
	return tca.client.Incr(ctx, key).Result()
}

// Decr 递减
func (tca *testCacheAdapter) Decr(ctx context.Context, key string) (int64, error) {
	return tca.client.Decr(ctx, key).Result()
}

// HSet 设置哈希字段
func (tca *testCacheAdapter) HSet(ctx context.Context, key string, values ...interface{}) error {
	return tca.client.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段
func (tca *testCacheAdapter) HGet(ctx context.Context, key, field string) (string, error) {
	return tca.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func (tca *testCacheAdapter) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return tca.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (tca *testCacheAdapter) HDel(ctx context.Context, key string, fields ...string) error {
	return tca.client.HDel(ctx, key, fields...).Err()
}

// LPush 左推入列表
func (tca *testCacheAdapter) LPush(ctx context.Context, key string, values ...interface{}) error {
	return tca.client.LPush(ctx, key, values...).Err()
}

// RPush 右推入列表
func (tca *testCacheAdapter) RPush(ctx context.Context, key string, values ...interface{}) error {
	return tca.client.RPush(ctx, key, values...).Err()
}

// LPop 左弹出列表
func (tca *testCacheAdapter) LPop(ctx context.Context, key string) (string, error) {
	return tca.client.LPop(ctx, key).Result()
}

// RPop 右弹出列表
func (tca *testCacheAdapter) RPop(ctx context.Context, key string) (string, error) {
	return tca.client.RPop(ctx, key).Result()
}

// LLen 获取列表长度
func (tca *testCacheAdapter) LLen(ctx context.Context, key string) (int64, error) {
	return tca.client.LLen(ctx, key).Result()
}

// Health 健康检查
func (tca *testCacheAdapter) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return tca.client.Ping(ctx).Err()
}

// cacheItem 内存缓存项
type cacheItem struct {
	value      string
	expiration time.Time
}

// memoryCache 内存缓存实现
type memoryCache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

// Set 设置缓存
func (mc *memoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	var exp time.Time
	if expiration > 0 {
		exp = time.Now().Add(expiration)
	}

	mc.data[key] = cacheItem{
		value:      fmt.Sprintf("%v", value),
		expiration: exp,
	}
	return nil
}

// Get 获取缓存
func (mc *memoryCache) Get(ctx context.Context, key string) (string, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, exists := mc.data[key]
	if !exists {
		return "", fmt.Errorf("key not found")
	}

	if !item.expiration.IsZero() && time.Now().After(item.expiration) {
		delete(mc.data, key)
		return "", fmt.Errorf("key expired")
	}

	return item.value, nil
}

// Del 删除缓存
func (mc *memoryCache) Del(ctx context.Context, keys ...string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	for _, key := range keys {
		delete(mc.data, key)
	}
	return nil
}

// Exists 检查键是否存在
func (mc *memoryCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	count := int64(0)
	for _, key := range keys {
		if item, exists := mc.data[key]; exists {
			if item.expiration.IsZero() || time.Now().Before(item.expiration) {
				count++
			}
		}
	}
	return count, nil
}

// Expire 设置过期时间
func (mc *memoryCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if item, exists := mc.data[key]; exists {
		item.expiration = time.Now().Add(expiration)
		mc.data[key] = item
	}
	return nil
}

// TTL 获取剩余过期时间
func (mc *memoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, exists := mc.data[key]
	if !exists {
		return -2 * time.Second, nil // key不存在
	}

	if item.expiration.IsZero() {
		return -1 * time.Second, nil // 永不过期
	}

	ttl := time.Until(item.expiration)
	if ttl <= 0 {
		return -2 * time.Second, nil // 已过期
	}

	return ttl, nil
}

// Health 健康检查
func (mc *memoryCache) Health() error {
	return nil // 内存缓存总是健康的
}

// Close 关闭缓存
func (mc *memoryCache) Close() error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.data = make(map[string]cacheItem)
	return nil
}
