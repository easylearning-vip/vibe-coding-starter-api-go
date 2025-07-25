package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockCache 缓存模拟
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Del(ctx context.Context, keys ...string) error {
	// 将可变参数转换为interface{}切片
	callArgs := make([]interface{}, len(keys)+1)
	callArgs[0] = ctx
	for i, key := range keys {
		callArgs[i+1] = key
	}
	args := m.Called(callArgs...)
	return args.Error(0)
}

func (m *MockCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	args := m.Called(ctx, keys)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	args := m.Called(ctx, key, expiration)
	return args.Error(0)
}

func (m *MockCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockCache) Incr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCache) Decr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCache) HSet(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockCache) HGet(ctx context.Context, key, field string) (string, error) {
	args := m.Called(ctx, key, field)
	return args.String(0), args.Error(1)
}

func (m *MockCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockCache) HDel(ctx context.Context, key string, fields ...string) error {
	args := m.Called(ctx, key, fields)
	return args.Error(0)
}

func (m *MockCache) LPush(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockCache) RPush(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockCache) LPop(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) RPop(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) LLen(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCache) Health() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCache) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockDatabase 数据库模拟
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetDB() *gorm.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*gorm.DB)
}

func (m *MockDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabase) AutoMigrate(dst ...interface{}) error {
	args := m.Called(dst)
	return args.Error(0)
}

func (m *MockDatabase) Health() error {
	args := m.Called()
	return args.Error(0)
}
