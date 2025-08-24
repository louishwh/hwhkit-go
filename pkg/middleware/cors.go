package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string      // 允许的来源
	AllowMethods     []string      // 允许的HTTP方法
	AllowHeaders     []string      // 允许的请求头
	ExposeHeaders    []string      // 暴露的响应头
	AllowCredentials bool          // 是否允许凭证
	MaxAge           time.Duration // 预检请求缓存时间
}

// DefaultCORSConfig 默认CORS配置
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}

// CORS 创建CORS中间件
func CORS(config ...*CORSConfig) gin.HandlerFunc {
	var cfg *CORSConfig
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	} else {
		cfg = DefaultCORSConfig()
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 检查是否允许该来源
		if cfg.AllowCredentials && !isOriginAllowed(origin, cfg.AllowOrigins) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		
		// 设置CORS头
		if len(cfg.AllowOrigins) == 1 && cfg.AllowOrigins[0] == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if isOriginAllowed(origin, cfg.AllowOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		if len(cfg.AllowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
		}
		
		if len(cfg.AllowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ", "))
		}
		
		if len(cfg.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))
		}
		
		if cfg.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		if cfg.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", formatDuration(cfg.MaxAge))
		}
		
		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// isOriginAllowed 检查来源是否被允许
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
		// 支持简单的通配符匹配
		if strings.HasSuffix(allowed, "*") {
			prefix := strings.TrimSuffix(allowed, "*")
			if strings.HasPrefix(origin, prefix) {
				return true
			}
		}
	}
	return false
}

// formatDuration 格式化时间为秒数
func formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	return strings.Trim(strings.Replace(strings.Replace(fmt.Sprintf("%d", seconds), "m", "", -1), "s", "", -1), " ")
}

// CORSWithDomains 创建允许特定域名的CORS中间件
func CORSWithDomains(domains ...string) gin.HandlerFunc {
	config := DefaultCORSConfig()
	config.AllowOrigins = domains
	config.AllowCredentials = true
	return CORS(config)
}

// CORSWithConfig 使用自定义配置创建CORS中间件
func CORSWithConfig(config *CORSConfig) gin.HandlerFunc {
	return CORS(config)
}