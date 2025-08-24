package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hwh/hwhkit-go/pkg/auth"
)

// JWTConfig JWT中间件配置
type JWTConfig struct {
	AuthManager    *auth.Manager // JWT管理器
	TokenLookup    string        // 令牌查找方式: "header:Authorization", "query:token", "cookie:token"
	TokenHeadName  string        // 令牌头部名称，默认为"Bearer"
	SkipPaths      []string      // 跳过验证的路径
	ErrorHandler   func(*gin.Context, error) // 错误处理函数
	SuccessHandler func(*gin.Context, *auth.Claims) // 成功处理函数
}

// DefaultJWTConfig 默认JWT配置
func DefaultJWTConfig(authManager *auth.Manager) *JWTConfig {
	return &JWTConfig{
		AuthManager:   authManager,
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		SkipPaths:     []string{},
		ErrorHandler: func(c *gin.Context, err error) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
			c.Abort()
		},
		SuccessHandler: func(c *gin.Context, claims *auth.Claims) {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
			c.Set("claims", claims)
		},
	}
}

// JWT 创建JWT认证中间件
func JWT(config *JWTConfig) gin.HandlerFunc {
	if config == nil {
		panic("JWT middleware requires a config")
	}
	
	if config.AuthManager == nil {
		panic("JWT middleware requires an auth manager")
	}
	
	return func(c *gin.Context) {
		// 检查是否跳过验证
		if shouldSkipPath(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}
		
		// 提取令牌
		token, err := extractToken(c, config)
		if err != nil {
			config.ErrorHandler(c, fmt.Errorf("failed to extract token: %w", err))
			return
		}
		
		// 验证令牌
		claims, err := config.AuthManager.ValidateToken(token)
		if err != nil {
			config.ErrorHandler(c, fmt.Errorf("invalid token: %w", err))
			return
		}
		
		// 调用成功处理函数
		config.SuccessHandler(c, claims)
		
		c.Next()
	}
}

// JWTWithManager 使用认证管理器创建JWT中间件
func JWTWithManager(authManager *auth.Manager) gin.HandlerFunc {
	config := DefaultJWTConfig(authManager)
	return JWT(config)
}

// JWTWithConfig 使用自定义配置创建JWT中间件
func JWTWithConfig(config *JWTConfig) gin.HandlerFunc {
	return JWT(config)
}

// JWTOptional 创建可选的JWT认证中间件（令牌存在时验证，不存在时跳过）
func JWTOptional(config *JWTConfig) gin.HandlerFunc {
	if config == nil {
		panic("JWT middleware requires a config")
	}
	
	if config.AuthManager == nil {
		panic("JWT middleware requires an auth manager")
	}
	
	return func(c *gin.Context) {
		// 尝试提取令牌
		token, err := extractToken(c, config)
		if err != nil {
			// 没有令牌时继续执行
			c.Next()
			return
		}
		
		// 验证令牌
		claims, err := config.AuthManager.ValidateToken(token)
		if err != nil {
			// 令牌无效时继续执行
			c.Next()
			return
		}
		
		// 调用成功处理函数
		config.SuccessHandler(c, claims)
		
		c.Next()
	}
}

// RequireRole 创建角色验证中间件
func RequireRole(authManager *auth.Manager, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取令牌
		claimsValue, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "No authentication claims found",
			})
			c.Abort()
			return
		}
		
		claims, ok := claimsValue.(*auth.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized", 
				"message": "Invalid authentication claims",
			})
			c.Abort()
			return
		}
		
		// 检查角色权限
		userRole := claims.Role
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}
		
		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": fmt.Sprintf("Required role: %v, got: %s", roles, userRole),
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RequireAnyRole 创建任意角色验证中间件（只要有其中一个角色即可）
func RequireAnyRole(authManager *auth.Manager, roles ...string) gin.HandlerFunc {
	return RequireRole(authManager, roles...)
}

// RequireAllRoles 创建全部角色验证中间件（需要拥有所有指定角色）
func RequireAllRoles(authManager *auth.Manager, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsValue, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "No authentication claims found",
			})
			c.Abort()
			return
		}
		
		claims, ok := claimsValue.(*auth.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid authentication claims",
			})
			c.Abort()
			return
		}
		
		// 这里简化实现，假设用户只有一个角色
		// 在实际项目中，可能需要支持多角色
		userRole := claims.Role
		for _, role := range roles {
			if userRole != role {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "Forbidden",
					"message": fmt.Sprintf("Required all roles: %v, got: %s", roles, userRole),
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

// extractToken 提取令牌
func extractToken(c *gin.Context, config *JWTConfig) (string, error) {
	parts := strings.Split(config.TokenLookup, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token lookup format")
	}
	
	method := parts[0]
	key := parts[1]
	
	switch method {
	case "header":
		authHeader := c.GetHeader(key)
		if authHeader == "" {
			return "", fmt.Errorf("authorization header not found")
		}
		
		// 检查Bearer前缀
		if config.TokenHeadName != "" {
			prefix := config.TokenHeadName + " "
			if !strings.HasPrefix(authHeader, prefix) {
				return "", fmt.Errorf("authorization header format must be %s {token}", config.TokenHeadName)
			}
			return strings.TrimPrefix(authHeader, prefix), nil
		}
		return authHeader, nil
		
	case "query":
		token := c.Query(key)
		if token == "" {
			return "", fmt.Errorf("token not found in query parameters")
		}
		return token, nil
		
	case "cookie":
		token, err := c.Cookie(key)
		if err != nil {
			return "", fmt.Errorf("token not found in cookie: %w", err)
		}
		return token, nil
		
	default:
		return "", fmt.Errorf("unsupported token lookup method: %s", method)
	}
}

// shouldSkipPath 检查是否应该跳过路径验证
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if matchPath(path, skipPath) {
			return true
		}
	}
	return false
}

// matchPath 路径匹配（支持简单的通配符）
func matchPath(path, pattern string) bool {
	if pattern == path {
		return true
	}
	
	// 支持*通配符
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(path, prefix)
	}
	
	return false
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) (int64, bool) {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(int64); ok {
			return id, true
		}
	}
	return 0, false
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name, true
		}
	}
	return "", false
}

// GetUserRole 从上下文获取用户角色
func GetUserRole(c *gin.Context) (string, bool) {
	if role, exists := c.Get("role"); exists {
		if r, ok := role.(string); ok {
			return r, true
		}
	}
	return "", false
}

// GetClaims 从上下文获取完整的JWT声明
func GetClaims(c *gin.Context) (*auth.Claims, bool) {
	if claims, exists := c.Get("claims"); exists {
		if c, ok := claims.(*auth.Claims); ok {
			return c, true
		}
	}
	return nil, false
}