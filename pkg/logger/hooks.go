package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// CallerHook 调用者信息钩子
type CallerHook struct {
	Field     string
	Skip      int
	levels    []logrus.Level
	Formatter func(file, function string, line int) string
}

// NewCallerHook 创建调用者信息钩子
func NewCallerHook(levels ...logrus.Level) *CallerHook {
	hook := CallerHook{
		Field:  "caller",
		Skip:   5,
		levels: levels,
		Formatter: func(file, function string, line int) string {
			return fmt.Sprintf("%s:%d", file, line)
		},
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels 返回支持的日志级别
func (hook *CallerHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire 执行钩子
func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = hook.getCaller()
	return nil
}

// getCaller 获取调用者信息
func (hook *CallerHook) getCaller() string {
	if pc, file, line, ok := runtime.Caller(hook.Skip); ok {
		if hook.Formatter == nil {
			return fmt.Sprintf("%s:%d", file, line)
		}
		
		function := runtime.FuncForPC(pc).Name()
		return hook.Formatter(file, function, line)
	}
	return "unknown"
}

// RequestIDHook 请求ID钩子
type RequestIDHook struct {
	levels []logrus.Level
}

// NewRequestIDHook 创建请求ID钩子
func NewRequestIDHook(levels ...logrus.Level) *RequestIDHook {
	hook := RequestIDHook{
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels 返回支持的日志级别
func (hook *RequestIDHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire 执行钩子
func (hook *RequestIDHook) Fire(entry *logrus.Entry) error {
	// 如果存在请求ID，添加到日志中
	if requestID, exists := entry.Data["request_id"]; exists {
		entry.Data["request_id"] = requestID
	}
	return nil
}

// SensitiveDataHook 敏感数据过滤钩子
type SensitiveDataHook struct {
	levels       []logrus.Level
	SensitiveKeys []string
	Replacement   string
}

// NewSensitiveDataHook 创建敏感数据过滤钩子
func NewSensitiveDataHook(sensitiveKeys []string, levels ...logrus.Level) *SensitiveDataHook {
	hook := SensitiveDataHook{
		levels:        levels,
		SensitiveKeys: sensitiveKeys,
		Replacement:   "***",
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels 返回支持的日志级别
func (hook *SensitiveDataHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire 执行钩子
func (hook *SensitiveDataHook) Fire(entry *logrus.Entry) error {
	for _, key := range hook.SensitiveKeys {
		if _, exists := entry.Data[key]; exists {
			entry.Data[key] = hook.Replacement
		}
	}
	
	// 过滤消息中的敏感信息
	message := entry.Message
	for _, key := range hook.SensitiveKeys {
		if strings.Contains(strings.ToLower(message), strings.ToLower(key)) {
			message = strings.ReplaceAll(message, key, hook.Replacement)
		}
	}
	entry.Message = message
	
	return nil
}

// ContextHook 上下文信息钩子
type ContextHook struct {
	levels []logrus.Level
}

// NewContextHook 创建上下文信息钩子
func NewContextHook(levels ...logrus.Level) *ContextHook {
	hook := ContextHook{
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels 返回支持的日志级别
func (hook *ContextHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire 执行钩子
func (hook *ContextHook) Fire(entry *logrus.Entry) error {
	// 添加服务信息
	entry.Data["service"] = "hwhkit-go"
	
	// 添加环境信息
	if env := entry.Data["env"]; env == nil {
		entry.Data["env"] = "development"
	}
	
	return nil
}

// MetricsHook 指标收集钩子
type MetricsHook struct {
	levels   []logrus.Level
	counters map[logrus.Level]int64
}

// NewMetricsHook 创建指标收集钩子
func NewMetricsHook(levels ...logrus.Level) *MetricsHook {
	hook := MetricsHook{
		levels:   levels,
		counters: make(map[logrus.Level]int64),
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels 返回支持的日志级别
func (hook *MetricsHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire 执行钩子
func (hook *MetricsHook) Fire(entry *logrus.Entry) error {
	hook.counters[entry.Level]++
	return nil
}

// GetCounters 获取计数器
func (hook *MetricsHook) GetCounters() map[logrus.Level]int64 {
	counters := make(map[logrus.Level]int64)
	for level, count := range hook.counters {
		counters[level] = count
	}
	return counters
}

// Reset 重置计数器
func (hook *MetricsHook) Reset() {
	hook.counters = make(map[logrus.Level]int64)
}

// GetTotalCount 获取总计数
func (hook *MetricsHook) GetTotalCount() int64 {
	var total int64
	for _, count := range hook.counters {
		total += count
	}
	return total
}

// GetCountByLevel 根据级别获取计数
func (hook *MetricsHook) GetCountByLevel(level logrus.Level) int64 {
	return hook.counters[level]
}