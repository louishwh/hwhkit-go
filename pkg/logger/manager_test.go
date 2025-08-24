package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/hwh/hwhkit-go/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	}
	
	manager, err := New(cfg)
	require.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Equal(t, logrus.InfoLevel, manager.logger.GetLevel())
}

func TestNewWithInvalidLevel(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "invalid",
		Format: "json",
		Output: "console",
	}
	
	_, err := New(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid log level")
}

func TestSetFormatter(t *testing.T) {
	tests := []struct {
		name   string
		format string
		valid  bool
	}{
		{"json formatter", "json", true},
		{"text formatter", "text", true},
		{"invalid formatter", "invalid", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.LogConfig{
				Level:  "info",
				Format: tt.format,
				Output: "console",
			}
			
			manager, err := New(cfg)
			if tt.valid {
				require.NoError(t, err)
				assert.NotNil(t, manager)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestLoggingMethods(t *testing.T) {
	// 创建一个缓冲区来捕获日志输出
	var buf bytes.Buffer
	
	cfg := &config.LogConfig{
		Level:  "debug",
		Format: "text",
		Output: "console",
	}
	
	manager, err := New(cfg)
	require.NoError(t, err)
	
	// 设置输出到缓冲区
	manager.logger.SetOutput(&buf)
	
	// 测试不同级别的日志
	manager.Debug("debug message")
	manager.Info("info message")
	manager.Warn("warn message")
	manager.Error("error message")
	
	output := buf.String()
	assert.Contains(t, output, "debug message")
	assert.Contains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
}

func TestFormattedLogging(t *testing.T) {
	var buf bytes.Buffer
	
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "text",
		Output: "console",
	}
	
	manager, err := New(cfg)
	require.NoError(t, err)
	manager.logger.SetOutput(&buf)
	
	manager.Infof("formatted message: %s, number: %d", "test", 42)
	
	output := buf.String()
	assert.Contains(t, output, "formatted message: test, number: 42")
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer
	
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	}
	
	manager, err := New(cfg)
	require.NoError(t, err)
	manager.logger.SetOutput(&buf)
	
	manager.WithFields(Fields{
		"user_id": 123,
		"action":  "login",
	}).Info("user action")
	
	output := buf.String()
	assert.Contains(t, output, "user_id")
	assert.Contains(t, output, "123")
	assert.Contains(t, output, "action")
	assert.Contains(t, output, "login")
}

func TestWithField(t *testing.T) {
	var buf bytes.Buffer
	
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "json",
		Output: "console",
	}
	
	manager, err := New(cfg)
	require.NoError(t, err)
	manager.logger.SetOutput(&buf)
	
	manager.WithField("request_id", "req-123").Info("processing request")
	
	output := buf.String()
	assert.Contains(t, output, "request_id")
	assert.Contains(t, output, "req-123")
}

func TestWithError(t *testing.T) {
	var buf bytes.Buffer
	
	cfg := &config.LogConfig{
		Level:  "error",
		Format: "json",
		Output: "console",
	}
	
	manager, err := New(cfg)
	require.NoError(t, err)
	manager.logger.SetOutput(&buf)
	
	testErr := assert.AnError
	manager.WithError(testErr).Error("operation failed")
	
	output := buf.String()
	assert.Contains(t, output, "error")
	assert.Contains(t, output, testErr.Error())
}

func TestSetLevel(t *testing.T) {
	cfg := &config.LogConfig{
		Level:  "info",
		Format: "text",
		Output: "console",
	}
	
	manager, err := New(cfg)
	require.NoError(t, err)
	
	// 初始级别应该是info
	assert.Equal(t, "info", manager.GetLevel())
	
	// 修改级别为debug
	err = manager.SetLevel("debug")
	require.NoError(t, err)
	assert.Equal(t, "debug", manager.GetLevel())
	
	// 测试无效级别
	err = manager.SetLevel("invalid")
	assert.Error(t, err)
}

func TestFileOutput(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
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
	
	manager, err := New(cfg)
	require.NoError(t, err)
	
	manager.Info("test log message")
	
	// 检查文件是否存在
	_, err = os.Stat(logFile)
	assert.NoError(t, err)
}

func TestCreateLogDir(t *testing.T) {
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "logs", "app")
	logFile := filepath.Join(logDir, "test.log")
	
	cfg := &config.LogConfig{
		Level:    "info",
		Format:   "text",
		Output:   "file",
		FilePath: logFile,
	}
	
	manager, err := New(cfg)
	require.NoError(t, err)
	
	// 检查目录是否被创建
	_, err = os.Stat(logDir)
	assert.NoError(t, err)
}