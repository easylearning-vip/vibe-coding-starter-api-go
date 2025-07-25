package config

import (
	"fmt"
	"os"
	"strconv"

	"vibe-coding-starter/internal/config"
)

// TestConfig 测试配置
type TestConfig struct {
	*config.Config
}

// NewTestConfig 创建测试配置
func NewTestConfig() *TestConfig {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         8080,
			Mode:         "test",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Database: config.DatabaseConfig{
			Driver:          "sqlite",
			Host:            "",
			Port:            0,
			Username:        "",
			Password:        "",
			Database:        ":memory:", // 使用内存数据库
			Charset:         "",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: 3600, // 1 hour in seconds
		},
		Cache: config.CacheConfig{
			Driver:   "redis",
			Host:     getEnvOrDefault("TEST_REDIS_HOST", "localhost"),
			Port:     getEnvIntOrDefault("TEST_REDIS_PORT", 6380),
			Password: "",
			Database: 1, // 使用不同的数据库
			PoolSize: 10,
		},
		Logger: config.LoggerConfig{
			Level:      "debug",
			Format:     "json",
			Output:     "stdout",
			Filename:   "",
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 3,
			Compress:   true,
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret-key-for-testing-only",
			Issuer:     "vibe-coding-starter-test",
			Expiration: 86400, // 24 hours in seconds
		},
	}

	return &TestConfig{Config: cfg}
}

// GetDSN 获取测试数据库连接字符串
func (c *TestConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.Charset,
	)
}

// GetRedisAddr 获取测试Redis地址
func (c *TestConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Cache.Host, c.Cache.Port)
}

// getEnvOrDefault 获取环境变量或默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault 获取环境变量整数值或默认值
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
