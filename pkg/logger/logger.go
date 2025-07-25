package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"vibe-coding-starter/internal/config"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	With(fields ...interface{}) Logger
	Sync() error
}

// zapLogger Zap 日志实现
type zapLogger struct {
	logger *zap.SugaredLogger
}

// New 创建新的日志实例
func New(cfg *config.Config) (Logger, error) {
	// 配置日志级别
	level, err := zapcore.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// 配置编码器
	var encoderConfig zapcore.EncoderConfig
	if cfg.Logger.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		// 使用开发配置，更适合inline格式
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		encoderConfig.ConsoleSeparator = " | "
	}

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	// 配置编码器
	var encoder zapcore.Encoder
	if cfg.Logger.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置输出
	var writeSyncer zapcore.WriteSyncer
	if cfg.Logger.Output == "file" {
		// 确保日志目录存在
		if err := os.MkdirAll("logs", 0755); err != nil {
			return nil, err
		}

		file, err := os.OpenFile(cfg.Logger.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		writeSyncer = zapcore.AddSync(file)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 创建核心
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 创建 logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &zapLogger{
		logger: logger.Sugar(),
	}, nil
}

// Debug 记录调试日志
func (l *zapLogger) Debug(msg string, fields ...interface{}) {
	l.logger.Debugw(msg, fields...)
}

// Info 记录信息日志
func (l *zapLogger) Info(msg string, fields ...interface{}) {
	l.logger.Infow(msg, fields...)
}

// Warn 记录警告日志
func (l *zapLogger) Warn(msg string, fields ...interface{}) {
	l.logger.Warnw(msg, fields...)
}

// Error 记录错误日志
func (l *zapLogger) Error(msg string, fields ...interface{}) {
	l.logger.Errorw(msg, fields...)
}

// Fatal 记录致命错误日志
func (l *zapLogger) Fatal(msg string, fields ...interface{}) {
	l.logger.Fatalw(msg, fields...)
}

// With 添加字段
func (l *zapLogger) With(fields ...interface{}) Logger {
	return &zapLogger{
		logger: l.logger.With(fields...),
	}
}

// Sync 同步日志
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
