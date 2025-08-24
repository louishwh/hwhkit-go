package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hwh/hwhkit-go/pkg/logger"
)

// LoggerConfig 日志中间件配置
type LoggerConfig struct {
	Logger         *logger.Manager // 日志管理器
	SkipPaths      []string        // 跳过记录的路径
	LogRequestBody bool            // 是否记录请求体
	LogResponseBody bool           // 是否记录响应体
	MaxBodySize    int64           // 最大记录的请求/响应体大小
}

// DefaultLoggerConfig 默认日志配置
func DefaultLoggerConfig(log *logger.Manager) *LoggerConfig {
	return &LoggerConfig{
		Logger:          log,
		SkipPaths:       []string{"/health", "/metrics"},
		LogRequestBody:  false,
		LogResponseBody: false,
		MaxBodySize:     1024 * 1024, // 1MB
	}
}

// responseWriter 自定义响应写入器，用于捕获响应数据
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入响应数据
func (w *responseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// WriteString 写入字符串响应数据
func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// Logger 创建日志记录中间件
func Logger(config ...*LoggerConfig) gin.HandlerFunc {
	var cfg *LoggerConfig
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	} else {
		// 如果没有提供配置，使用默认的空配置
		cfg = &LoggerConfig{
			SkipPaths:       []string{"/health", "/metrics"},
			LogRequestBody:  false,
			LogResponseBody: false,
			MaxBodySize:     1024 * 1024,
		}
	}

	return func(c *gin.Context) {
		// 检查是否跳过记录
		if shouldSkipPath(c.Request.URL.Path, cfg.SkipPaths) {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		// 读取请求体
		var requestBody string
		if cfg.LogRequestBody && c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil && int64(len(bodyBytes)) <= cfg.MaxBodySize {
				requestBody = string(bodyBytes)
				// 重新设置请求体，以供后续处理使用
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 创建自定义响应写入器
		var responseBody string
		var writer *responseWriter
		if cfg.LogResponseBody {
			writer = &responseWriter{
				ResponseWriter: c.Writer,
				body:          bytes.NewBuffer(nil),
			}
			c.Writer = writer
		}

		// 处理请求
		c.Next()

		// 获取响应体
		if cfg.LogResponseBody && writer != nil && int64(writer.body.Len()) <= cfg.MaxBodySize {
			responseBody = writer.body.String()
		}

		// 计算处理时间
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		// 构建完整的URL
		fullPath := path
		if rawQuery != "" {
			fullPath = path + "?" + rawQuery
		}

		// 构建日志字段
		fields := logger.Fields{
			"ip":         clientIP,
			"method":     method,
			"path":       fullPath,
			"status":     statusCode,
			"latency":    latency.String(),
			"user_agent": userAgent,
		}

		// 添加请求体字段
		if requestBody != "" {
			fields["request_body"] = requestBody
		}

		// 添加响应体字段
		if responseBody != "" {
			fields["response_body"] = responseBody
		}

		// 添加用户信息（如果存在）
		if userID, exists := GetUserID(c); exists {
			fields["user_id"] = userID
		}
		if username, exists := GetUsername(c); exists {
			fields["username"] = username
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		// 记录日志
		logMessage := fmt.Sprintf("%s %s %d %s", method, fullPath, statusCode, latency)

		if cfg.Logger != nil {
			// 根据状态码决定日志级别
			switch {
			case statusCode >= 500:
				cfg.Logger.WithFields(fields).Error(logMessage)
			case statusCode >= 400:
				cfg.Logger.WithFields(fields).Warn(logMessage)
			default:
				cfg.Logger.WithFields(fields).Info(logMessage)
			}
		} else {
			// 如果没有日志管理器，使用gin的默认日志
			fmt.Printf("[GIN] %s\n", logMessage)
		}
	}
}

// LoggerWithManager 使用日志管理器创建日志中间件
func LoggerWithManager(log *logger.Manager) gin.HandlerFunc {
	config := DefaultLoggerConfig(log)
	return Logger(config)
}

// LoggerWithConfig 使用自定义配置创建日志中间件
func LoggerWithConfig(config *LoggerConfig) gin.HandlerFunc {
	return Logger(config)
}

// RequestLogger 专门记录请求的中间件
func RequestLogger(log *logger.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 记录请求开始
		if log != nil {
			log.WithFields(logger.Fields{
				"ip":         c.ClientIP(),
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"user_agent": c.Request.UserAgent(),
			}).Info("Request started")
		}
		
		c.Next()
		
		// 记录请求结束
		latency := time.Since(start)
		if log != nil {
			log.WithFields(logger.Fields{
				"ip":      c.ClientIP(),
				"method":  c.Request.Method,
				"path":    c.Request.URL.Path,
				"status":  c.Writer.Status(),
				"latency": latency.String(),
			}).Info("Request completed")
		}
	}
}

// ErrorLogger 错误记录中间件
func ErrorLogger(log *logger.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// 记录错误
		if len(c.Errors) > 0 && log != nil {
			for _, err := range c.Errors {
				log.WithFields(logger.Fields{
					"ip":     c.ClientIP(),
					"method": c.Request.Method,
					"path":   c.Request.URL.Path,
					"type":   err.Type,
				}).Error(err.Error())
			}
		}
	}
}

// AccessLogger 访问日志中间件（简化版本）
func AccessLogger(log *logger.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		
		if log != nil {
			latency := time.Since(start)
			log.WithFields(logger.Fields{
				"timestamp": start.Format(time.RFC3339),
				"client_ip": c.ClientIP(),
				"method":    c.Request.Method,
				"uri":       c.Request.RequestURI,
				"status":    c.Writer.Status(),
				"latency":   latency.Nanoseconds(),
				"size":      c.Writer.Size(),
			}).Info("Access log")
		}
	}
}