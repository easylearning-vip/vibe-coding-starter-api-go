package test

import (
	"testing"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/pkg/logger"
)

func TestLoggerLevels(t *testing.T) {
	// 加载配置
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 创建日志实例
	log, err := logger.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 测试不同日志级别
	t.Log("Testing different log levels:")
	
	log.Debug("This is a DEBUG message", "key1", "value1", "number", 123)
	log.Info("This is an INFO message", "user", "admin", "action", "login")
	log.Warn("This is a WARN message", "latency", 200, "threshold", 100)
	log.Error("This is an ERROR message", "error", "test error", "request_id", "req-123")
	
	// 测试带字段的日志
	contextLogger := log.With("service", "user-service", "version", "1.0.0")
	contextLogger.Info("Service started successfully", "port", 8080)
	contextLogger.Debug("Processing request", "method", "POST", "path", "/api/users")
	contextLogger.Warn("High memory usage detected", "usage", "85%")
	contextLogger.Error("Database connection failed", "error", "connection timeout")

	// 同步日志
	if err := log.Sync(); err != nil {
		t.Logf("Failed to sync logger: %v", err)
	}
}

func TestLoggerFormats(t *testing.T) {
	// 测试JSON格式
	t.Log("Testing JSON format:")
	cfg, _ := config.New()
	cfg.Logger.Format = "json"
	
	jsonLogger, err := logger.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create JSON logger: %v", err)
	}
	
	jsonLogger.Info("JSON format test", "format", "json", "readable", false)
	
	// 测试Console格式
	t.Log("Testing Console format:")
	cfg.Logger.Format = "console"
	
	consoleLogger, err := logger.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create console logger: %v", err)
	}
	
	consoleLogger.Info("Console format test", "format", "console", "readable", true)
	
	// 同步日志
	jsonLogger.Sync()
	consoleLogger.Sync()
}
