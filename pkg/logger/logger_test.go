package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hwh/hwhkit-go/pkg/config"
	"github.com/sirupsen/logrus"
)

func TestLoggerCreation(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	if logger == nil {
		t.Fatal("Logger should not be nil")
	}
	
	if logger.GetLevel() != logrus.InfoLevel {
		t.Errorf("Expected log level Info, got %v", logger.GetLevel())
	}
}

func TestLoggerLevels(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "debug",
		Format: "text",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// 测试各种日志级别
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")
	
	// 测试格式化日志
	logger.Debugf("Debug message with param: %s", "test")
	logger.Infof("Info message with param: %d", 123)
	logger.Warnf("Warning message with param: %v", true)
	logger.Errorf("Error message with param: %f", 3.14)
}

func TestLoggerWithFields(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// 测试单个字段
	logger.WithField("user_id", "123").Info("User logged in")
	
	// 测试多个字段
	logger.WithFields(logrus.Fields{
		"user_id":   "123",
		"action":    "login",
		"timestamp": time.Now(),
	}).Info("User action logged")
	
	// 测试错误字段
	err = &TestError{Message: "test error"}
	logger.WithError(err).Error("An error occurred")
}

type TestError struct {
	Message string
}

func (e *TestError) Error() string {
	return e.Message
}

func TestLoggerFileOutput(t *testing.T) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "hwhkit_test_logs")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)
	
	logFile := filepath.Join(tempDir, "test.log")
	
	cfg := &config.LogConfig{
		Level:      "info",
		Format:     "text",
		Output:     "file",
		FilePath:   logFile,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
	}
	
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	logger.Info("Test log message")
	
	// 检查文件是否创建
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

func TestGlobalLogger(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "debug",
		Format: "text",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// 设置全局日志实例
	SetDefault(logger)
	
	// 测试全局函数
	Debug("Global debug message")
	Info("Global info message")
	Warn("Global warning message")
	Error("Global error message")
	
	// 测试格式化全局函数
	Debugf("Global debug with param: %s", "test")
	Infof("Global info with param: %d", 456)
	
	// 测试带字段的全局函数
	WithField("global", true).Info("Global message with field")
	WithFields(logrus.Fields{
		"key1": "value1",
		"key2": "value2",
	}).Info("Global message with fields")
}

func TestCallerHook(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// 添加调用者钩子
	hook := NewCallerHook(logrus.InfoLevel, logrus.ErrorLevel)
	logger.AddHook(hook)
	
	logger.Info("Message with caller info")
	logger.Error("Error with caller info")
}

func TestSensitiveDataHook(t *testing.T) {
	// 使用字节缓冲区捕获输出
	var buf bytes.Buffer
	
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "text",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// 重定向输出到缓冲区
	logger.GetLogger().SetOutput(&buf)
	
	// 添加敏感数据过滤钩子
	sensitiveKeys := []string{"password", "token", "secret"}
	hook := NewSensitiveDataHook(sensitiveKeys)
	logger.AddHook(hook)
	
	// 记录包含敏感数据的日志
	logger.WithFields(logrus.Fields{
		"user":     "john",
		"password": "secret123",
		"token":    "abc123",
	}).Info("User authentication")
	
	output := buf.String()
	
	// 检查敏感数据是否被替换
	if !bytes.Contains([]byte(output), []byte("***")) {
		t.Error("Sensitive data was not masked")
	}
	
	if bytes.Contains([]byte(output), []byte("secret123")) {
		t.Error("Password was not masked")
	}
	
	if bytes.Contains([]byte(output), []byte("abc123")) {
		t.Error("Token was not masked")
	}
}

func TestMetricsHook(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "debug",
		Format: "text",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// 添加指标钩子
	hook := NewMetricsHook()
	logger.AddHook(hook)
	
	// 记录不同级别的日志
	logger.Debug("Debug message")
	logger.Info("Info message 1")
	logger.Info("Info message 2")
	logger.Warn("Warning message")
	logger.Error("Error message")
	
	// 检查计数器
	counters := hook.GetCounters()
	
	if counters[logrus.DebugLevel] != 1 {
		t.Errorf("Expected 1 debug log, got %d", counters[logrus.DebugLevel])
	}
	
	if counters[logrus.InfoLevel] != 2 {
		t.Errorf("Expected 2 info logs, got %d", counters[logrus.InfoLevel])
	}
	
	if counters[logrus.WarnLevel] != 1 {
		t.Errorf("Expected 1 warn log, got %d", counters[logrus.WarnLevel])
	}
	
	if counters[logrus.ErrorLevel] != 1 {
		t.Errorf("Expected 1 error log, got %d", counters[logrus.ErrorLevel])
	}
	
	totalCount := hook.GetTotalCount()
	if totalCount != 5 {
		t.Errorf("Expected total count 5, got %d", totalCount)
	}
	
	// 测试重置
	hook.Reset()
	if hook.GetTotalCount() != 0 {
		t.Error("Counters were not reset")
	}
}

func TestContextHook(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// 添加上下文钩子
	hook := NewContextHook()
	logger.AddHook(hook)
	
	logger.Info("Message with context")
}

func TestLoggerSetLevel(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "text",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	
	// 测试设置级别
	logger.SetLevel(logrus.DebugLevel)
	if logger.GetLevel() != logrus.DebugLevel {
		t.Errorf("Expected debug level, got %v", logger.GetLevel())
	}
	
	// 测试级别检查
	if !logger.IsLevelEnabled(logrus.DebugLevel) {
		t.Error("Debug level should be enabled")
	}
	
	if !logger.IsLevelEnabled(logrus.InfoLevel) {
		t.Error("Info level should be enabled")
	}
	
	logger.SetLevel(logrus.WarnLevel)
	if logger.IsLevelEnabled(logrus.DebugLevel) {
		t.Error("Debug level should not be enabled")
	}
	
	if logger.IsLevelEnabled(logrus.InfoLevel) {
		t.Error("Info level should not be enabled")
	}
	
	if !logger.IsLevelEnabled(logrus.WarnLevel) {
		t.Error("Warn level should be enabled")
	}
}

// 基准测试
func BenchmarkLogger(b *testing.B) {
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message")
	}
}

func BenchmarkLoggerWithFields(b *testing.B) {
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	}
	
	logger, err := New(cfg)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}
	
	fields := logrus.Fields{
		"user_id": "123",
		"action":  "login",
		"ip":      "192.168.1.1",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(fields).Info("Benchmark message with fields")
	}
}