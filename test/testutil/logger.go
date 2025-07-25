package testutil

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"vibe-coding-starter/pkg/logger"
)

// TestLogger 测试日志包装器
type TestLogger struct {
	*zap.Logger
	t *testing.T
}

// NewTestLogger 创建测试日志器
func NewTestLogger(t *testing.T) *TestLogger {
	// 创建测试用的zap配置
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	config.Development = true
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	zapLogger, err := config.Build()
	if err != nil {
		t.Fatalf("Failed to create test logger: %v", err)
	}

	return &TestLogger{
		Logger: zapLogger,
		t:      t,
	}
}

// CreateTestLogger 创建实现logger.Logger接口的测试日志器
func (tl *TestLogger) CreateTestLogger() logger.Logger {
	return &testLoggerAdapter{tl}
}

// testLoggerAdapter 适配器，实现logger.Logger接口
type testLoggerAdapter struct {
	*TestLogger
}

// Debug 调试日志
func (tla *testLoggerAdapter) Debug(msg string, fields ...interface{}) {
	tla.Logger.Debug(msg, convertFields(fields...)...)
}

// Info 信息日志
func (tla *testLoggerAdapter) Info(msg string, fields ...interface{}) {
	tla.Logger.Info(msg, convertFields(fields...)...)
}

// Warn 警告日志
func (tla *testLoggerAdapter) Warn(msg string, fields ...interface{}) {
	tla.Logger.Warn(msg, convertFields(fields...)...)
}

// Error 错误日志
func (tla *testLoggerAdapter) Error(msg string, fields ...interface{}) {
	tla.Logger.Error(msg, convertFields(fields...)...)
}

// Fatal 致命错误日志
func (tla *testLoggerAdapter) Fatal(msg string, fields ...interface{}) {
	tla.Logger.Fatal(msg, convertFields(fields...)...)
}

// With 添加字段
func (tla *testLoggerAdapter) With(fields ...interface{}) logger.Logger {
	newLogger := tla.Logger.With(convertFields(fields...)...)
	return &testLoggerAdapter{
		TestLogger: &TestLogger{
			Logger: newLogger,
			t:      tla.t,
		},
	}
}

// Sync 同步日志
func (tla *testLoggerAdapter) Sync() error {
	return tla.Logger.Sync()
}

// convertFields 转换字段为zap字段
func convertFields(fields ...interface{}) []zap.Field {
	if len(fields)%2 != 0 {
		return []zap.Field{zap.Any("invalid_fields", fields)}
	}

	zapFields := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		value := fields[i+1]
		zapFields = append(zapFields, zap.Any(key, value))
	}

	return zapFields
}

// Close 关闭日志器
func (tl *TestLogger) Close() error {
	return tl.Logger.Sync()
}
