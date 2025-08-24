package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hwh/hwhkit-go/pkg/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Manager 日志管理器
type Manager struct {
	logger *logrus.Logger
	config *config.LogConfig
}

// Fields 日志字段类型
type Fields map[string]interface{}

// New 创建新的日志管理器
func New(cfg *config.LogConfig) (*Manager, error) {
	manager := &Manager{
		logger: logrus.New(),
		config: cfg,
	}
	
	if err := manager.configure(); err != nil {
		return nil, err
	}
	
	return manager, nil
}

// configure 配置日志管理器
func (m *Manager) configure() error {
	// 设置日志级别
	level, err := logrus.ParseLevel(m.config.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	m.logger.SetLevel(level)
	
	// 设置日志格式
	if err := m.setFormatter(); err != nil {
		return err
	}
	
	// 设置输出
	if err := m.setOutput(); err != nil {
		return err
	}
	
	// 设置调用者信息
	m.logger.SetReportCaller(true)
	
	return nil
}

// setFormatter 设置日志格式
func (m *Manager) setFormatter() error {
	switch strings.ToLower(m.config.Format) {
	case "json":
		m.logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := filepath.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	case "text":
		m.logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := filepath.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	default:
		return fmt.Errorf("unsupported log format: %s", m.config.Format)
	}
	
	return nil
}

// setOutput 设置日志输出
func (m *Manager) setOutput() error {
	switch strings.ToLower(m.config.Output) {
	case "console":
		m.logger.SetOutput(os.Stdout)
		
	case "file":
		if err := m.createLogDir(); err != nil {
			return err
		}
		
		fileWriter := &lumberjack.Logger{
			Filename:   m.config.FilePath,
			MaxSize:    m.config.MaxSize,
			MaxBackups: m.config.MaxBackups,
			MaxAge:     m.config.MaxAge,
			Compress:   m.config.Compress,
		}
		m.logger.SetOutput(fileWriter)
		
	case "both":
		if err := m.createLogDir(); err != nil {
			return err
		}
		
		fileWriter := &lumberjack.Logger{
			Filename:   m.config.FilePath,
			MaxSize:    m.config.MaxSize,
			MaxBackups: m.config.MaxBackups,
			MaxAge:     m.config.MaxAge,
			Compress:   m.config.Compress,
		}
		
		multiWriter := io.MultiWriter(os.Stdout, fileWriter)
		m.logger.SetOutput(multiWriter)
		
	default:
		return fmt.Errorf("unsupported log output: %s", m.config.Output)
	}
	
	return nil
}

// createLogDir 创建日志目录
func (m *Manager) createLogDir() error {
	if m.config.FilePath == "" {
		return fmt.Errorf("log file path is empty")
	}
	
	dir := filepath.Dir(m.config.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	
	return nil
}

// GetLogger 获取原始的logrus实例
func (m *Manager) GetLogger() *logrus.Logger {
	return m.logger
}

// WithFields 添加字段
func (m *Manager) WithFields(fields Fields) *logrus.Entry {
	return m.logger.WithFields(logrus.Fields(fields))
}

// WithField 添加单个字段
func (m *Manager) WithField(key string, value interface{}) *logrus.Entry {
	return m.logger.WithField(key, value)
}

// WithError 添加错误字段
func (m *Manager) WithError(err error) *logrus.Entry {
	return m.logger.WithError(err)
}

// Debug 记录调试日志
func (m *Manager) Debug(args ...interface{}) {
	m.logger.Debug(args...)
}

// Debugf 记录格式化调试日志
func (m *Manager) Debugf(format string, args ...interface{}) {
	m.logger.Debugf(format, args...)
}

// Info 记录信息日志
func (m *Manager) Info(args ...interface{}) {
	m.logger.Info(args...)
}

// Infof 记录格式化信息日志
func (m *Manager) Infof(format string, args ...interface{}) {
	m.logger.Infof(format, args...)
}

// Warn 记录警告日志
func (m *Manager) Warn(args ...interface{}) {
	m.logger.Warn(args...)
}

// Warnf 记录格式化警告日志
func (m *Manager) Warnf(format string, args ...interface{}) {
	m.logger.Warnf(format, args...)
}

// Error 记录错误日志
func (m *Manager) Error(args ...interface{}) {
	m.logger.Error(args...)
}

// Errorf 记录格式化错误日志
func (m *Manager) Errorf(format string, args ...interface{}) {
	m.logger.Errorf(format, args...)
}

// Fatal 记录致命错误日志并退出程序
func (m *Manager) Fatal(args ...interface{}) {
	m.logger.Fatal(args...)
}

// Fatalf 记录格式化致命错误日志并退出程序
func (m *Manager) Fatalf(format string, args ...interface{}) {
	m.logger.Fatalf(format, args...)
}

// Panic 记录panic日志并触发panic
func (m *Manager) Panic(args ...interface{}) {
	m.logger.Panic(args...)
}

// Panicf 记录格式化panic日志并触发panic
func (m *Manager) Panicf(format string, args ...interface{}) {
	m.logger.Panicf(format, args...)
}

// SetLevel 动态设置日志级别
func (m *Manager) SetLevel(level string) error {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	m.logger.SetLevel(logLevel)
	m.config.Level = level
	return nil
}

// GetLevel 获取当前日志级别
func (m *Manager) GetLevel() string {
	return m.logger.GetLevel().String()
}