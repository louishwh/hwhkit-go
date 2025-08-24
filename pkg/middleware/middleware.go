package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hwh/hwhkit-go/pkg/auth"
	"github.com/hwh/hwhkit-go/pkg/logger"
)

// MiddlewareManager 中间件管理器
type MiddlewareManager struct {
	authManager *auth.Manager
	logger      *logger.Manager
}

// NewMiddlewareManager 创建中间件管理器
func NewMiddlewareManager(authManager *auth.Manager, log *logger.Manager) *MiddlewareManager {
	return &MiddlewareManager{
		authManager: authManager,
		logger:      log,
	}
}

// DefaultCORS 默认CORS中间件
func (m *MiddlewareManager) DefaultCORS() gin.HandlerFunc {
	return CORS()
}

// CORS 自定义CORS中间件
func (m *MiddlewareManager) CORS(config *CORSConfig) gin.HandlerFunc {
	return CORSWithConfig(config)
}

// JWT JWT认证中间件
func (m *MiddlewareManager) JWT() gin.HandlerFunc {
	return JWTWithManager(m.authManager)
}

// JWTOptional 可选JWT认证中间件
func (m *MiddlewareManager) JWTOptional() gin.HandlerFunc {
	config := DefaultJWTConfig(m.authManager)
	return JWTOptional(config)
}

// RequireRole 角色验证中间件
func (m *MiddlewareManager) RequireRole(roles ...string) gin.HandlerFunc {
	return RequireRole(m.authManager, roles...)
}

// Logger 日志记录中间件
func (m *MiddlewareManager) Logger() gin.HandlerFunc {
	return LoggerWithManager(m.logger)
}

// RequestLogger 请求日志中间件
func (m *MiddlewareManager) RequestLogger() gin.HandlerFunc {
	return RequestLogger(m.logger)
}

// ErrorLogger 错误日志中间件
func (m *MiddlewareManager) ErrorLogger() gin.HandlerFunc {
	return ErrorLogger(m.logger)
}

// RateLimit 限流中间件
func (m *MiddlewareManager) RateLimit(rate, burst int) gin.HandlerFunc {
	return RateLimitByIP(rate, burst)
}

// RateLimitByUser 用户限流中间件
func (m *MiddlewareManager) RateLimitByUser(rate, burst int) gin.HandlerFunc {
	return RateLimitByUser(rate, burst)
}

// Common 通用中间件组合
func (m *MiddlewareManager) Common() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.DefaultCORS(),
		m.Logger(),
		m.ErrorLogger(),
		gin.Recovery(),
	}
}

// Protected 受保护路由中间件组合
func (m *MiddlewareManager) Protected() []gin.HandlerFunc {
	middlewares := m.Common()
	middlewares = append(middlewares, m.JWT())
	return middlewares
}

// Admin 管理员路由中间件组合
func (m *MiddlewareManager) Admin() []gin.HandlerFunc {
	middlewares := m.Protected()
	middlewares = append(middlewares, m.RequireRole("admin"))
	return middlewares
}

// API API路由中间件组合
func (m *MiddlewareManager) API() []gin.HandlerFunc {
	middlewares := m.Common()
	middlewares = append(middlewares, m.RateLimit(100, 10)) // 100 req/s, burst 10
	return middlewares
}

// PublicAPI 公共API路由中间件组合
func (m *MiddlewareManager) PublicAPI() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		m.DefaultCORS(),
		m.RequestLogger(),
		m.RateLimit(50, 5), // 50 req/s, burst 5
		gin.Recovery(),
	}
}

// 便捷函数，直接使用而无需创建管理器

// DefaultMiddlewares 默认中间件组合
func DefaultMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		CORS(),
		Logger(),
		gin.Recovery(),
	}
}

// ProtectedMiddlewares 受保护路由中间件组合（需要提供认证管理器）
func ProtectedMiddlewares(authManager *auth.Manager) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		CORS(),
		Logger(),
		JWTWithManager(authManager),
		gin.Recovery(),
	}
}

// APIMiddlewares API中间件组合
func APIMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		CORS(),
		Logger(),
		RateLimitByIP(100, 10),
		gin.Recovery(),
	}
}
