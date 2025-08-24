package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hwh/hwhkit-go/pkg/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志管理器
type Logger struct {
	logger *logrus.Logger
	config *config.LogConfig
}

// New 创建新的日志管理器
func New(cfg *config.LogConfig) (*Logger, error) {
	logger := logrus.New()
	
	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
	
	// 设置日志格式
	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
	
	// 设置输出
	if err := setOutput(logger, cfg); err != nil {
		return nil, err
	}
	
	return &Logger{
		logger: logger,
		config: cfg,
	}, nil
}

// setOutput 设置日志输出
func setOutput(logger *logrus.Logger, cfg *config.LogConfig) error {
	switch strings.ToLower(cfg.Output) {
	case "console":
		logger.SetOutput(os.Stdout)
	case "file":
		writer, err := createFileWriter(cfg)
		if err != nil {
			return err
		}
		logger.SetOutput(writer)
	case "both":
		fileWriter, err := createFileWriter(cfg)
		if err != nil {
			return err
		}
		multiWriter := io.MultiWriter(os.Stdout, fileWriter)
		logger.SetOutput(multiWriter)
	default:
		logger.SetOutput(os.Stdout)
	}
	return nil
}

// createFileWriter 创建文件写入器
func createFileWriter(cfg *config.LogConfig) (io.Writer, error) {
	// 确保日志目录存在
	logDir := filepath.Dir(cfg.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}
	
	// 使用 lumberjack 进行日志轮转
	return &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,    // MB
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,     // 天数
		Compress:   cfg.Compress,
	}, nil
}

// GetLogger 获取logrus实例
func (l *Logger) GetLogger() *logrus.Logger {
	return l.logger
}

// Debug 记录调试日志
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf 记录格式化调试日志
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Info 记录信息日志
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof 记录格式化信息日志
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warn 记录警告日志
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf 记录格式化警告日志
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error 记录错误日志
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf 记录格式化错误日志
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatal 记录致命错误日志并退出程序
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf 记录格式化致命错误日志并退出程序
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// Panic 记录恐慌日志并触发panic
func (l *Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

// Panicf 记录格式化恐慌日志并触发panic
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

// WithField 添加字段
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

// WithFields 添加多个字段
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

// WithError 添加错误字段
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.logger.WithError(err)
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level logrus.Level) {
	l.logger.SetLevel(level)
}

// GetLevel 获取当前日志级别
func (l *Logger) GetLevel() logrus.Level {
	return l.logger.GetLevel()
}

// AddHook 添加钩子
func (l *Logger) AddHook(hook logrus.Hook) {
	l.logger.AddHook(hook)
}

// IsLevelEnabled 检查级别是否启用
func (l *Logger) IsLevelEnabled(level logrus.Level) bool {
	return l.logger.IsLevelEnabled(level)
}

// 全局日志实例
var defaultLogger *Logger

// SetDefault 设置默认日志实例
func SetDefault(logger *Logger) {
	defaultLogger = logger
}

// GetDefault 获取默认日志实例
func GetDefault() *Logger {
	return defaultLogger
}

// 全局函数，使用默认日志实例

// Debug 使用默认日志实例记录调试日志
func Debug(args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(args...)
	}
}

// Debugf 使用默认日志实例记录格式化调试日志
func Debugf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debugf(format, args...)
	}
}

// Info 使用默认日志实例记录信息日志
func Info(args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(args...)
	}
}

// Infof 使用默认日志实例记录格式化信息日志
func Infof(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Infof(format, args...)
	}
}

// Warn 使用默认日志实例记录警告日志
func Warn(args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(args...)
	}
}

// Warnf 使用默认日志实例记录格式化警告日志
func Warnf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warnf(format, args...)
	}
}

// Error 使用默认日志实例记录错误日志
func Error(args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(args...)
	}
}

// Errorf 使用默认日志实例记录格式化错误日志
func Errorf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Errorf(format, args...)
	}
}

// Fatal 使用默认日志实例记录致命错误日志并退出程序
func Fatal(args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatal(args...)
	}
}

// Fatalf 使用默认日志实例记录格式化致命错误日志并退出程序
func Fatalf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatalf(format, args...)
	}
}

// WithField 使用默认日志实例添加字段
func WithField(key string, value interface{}) *logrus.Entry {
	if defaultLogger != nil {
		return defaultLogger.WithField(key, value)
	}
	return logrus.NewEntry(logrus.StandardLogger())
}

// WithFields 使用默认日志实例添加多个字段
func WithFields(fields logrus.Fields) *logrus.Entry {
	if defaultLogger != nil {
		return defaultLogger.WithFields(fields)
	}
	return logrus.NewEntry(logrus.StandardLogger())
}

// WithError 使用默认日志实例添加错误字段
func WithError(err error) *logrus.Entry {
	if defaultLogger != nil {
		return defaultLogger.WithError(err)
	}
	return logrus.NewEntry(logrus.StandardLogger())
}