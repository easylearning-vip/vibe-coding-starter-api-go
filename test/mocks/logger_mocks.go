package mocks

import (
	"github.com/stretchr/testify/mock"

	"vibe-coding-starter/pkg/logger"
)

// MockLogger 日志模拟
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	args := []interface{}{msg}
	args = append(args, fields...)
	m.Called(args...)
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
	args := []interface{}{msg}
	args = append(args, fields...)
	m.Called(args...)
}

func (m *MockLogger) Warn(msg string, fields ...interface{}) {
	args := []interface{}{msg}
	args = append(args, fields...)
	m.Called(args...)
}

func (m *MockLogger) Error(msg string, fields ...interface{}) {
	args := []interface{}{msg}
	args = append(args, fields...)
	m.Called(args...)
}

func (m *MockLogger) Fatal(msg string, fields ...interface{}) {
	args := []interface{}{msg}
	args = append(args, fields...)
	m.Called(args...)
}

func (m *MockLogger) With(fields ...interface{}) logger.Logger {
	args := m.Called(fields...)
	return args.Get(0).(logger.Logger)
}

func (m *MockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}
