package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hwh/hwhkit-go/pkg/middleware"
)

// RouterManager 路由管理器
type RouterManager struct {
	server     *Server
	engine     *gin.Engine
	middleware *middleware.MiddlewareManager
}

// NewRouterManager 创建路由管理器
func NewRouterManager(server *Server) *RouterManager {
	return &RouterManager{
		server:     server,
		engine:     server.GetEngine(),
		middleware: server.GetMiddleware(),
	}
}

// RouteGroup 路由组配置
type RouteGroup struct {
	Path        string
	Middlewares []gin.HandlerFunc
	Routes      []Route
	SubGroups   []RouteGroup
}

// Route 路由配置
type Route struct {
	Method      string
	Path        string
	Handlers    []gin.HandlerFunc
	Middlewares []gin.HandlerFunc
}

// RegisterRouteGroup 注册路由组
func (rm *RouterManager) RegisterRouteGroup(group RouteGroup) {
	rm.registerGroup(rm.engine.Group(group.Path), group)
}

// registerGroup 递归注册路由组
func (rm *RouterManager) registerGroup(ginGroup *gin.RouterGroup, group RouteGroup) {
	// 应用中间件
	for _, mw := range group.Middlewares {
		ginGroup.Use(mw)
	}
	
	// 注册路由
	for _, route := range group.Routes {
		handlers := append(route.Middlewares, route.Handlers...)
		
		switch route.Method {
		case "GET":
			ginGroup.GET(route.Path, handlers...)
		case "POST":
			ginGroup.POST(route.Path, handlers...)
		case "PUT":
			ginGroup.PUT(route.Path, handlers...)
		case "PATCH":
			ginGroup.PATCH(route.Path, handlers...)
		case "DELETE":
			ginGroup.DELETE(route.Path, handlers...)
		case "OPTIONS":
			ginGroup.OPTIONS(route.Path, handlers...)
		case "HEAD":
			ginGroup.HEAD(route.Path, handlers...)
		case "ANY":
			ginGroup.Any(route.Path, handlers...)
		}
	}
	
	// 注册子组
	for _, subGroup := range group.SubGroups {
		rm.registerGroup(ginGroup.Group(subGroup.Path), subGroup)
	}
}

// APIRouter API路由管理器
type APIRouter struct {
	routerManager *RouterManager
	server        *Server
}

// NewAPIRouter 创建API路由管理器
func NewAPIRouter(server *Server) *APIRouter {
	return &APIRouter{
		routerManager: NewRouterManager(server),
		server:        server,
	}
}

// SetupV1API 设置V1版本API路由
func (ar *APIRouter) SetupV1API() {
	v1 := ar.server.Group("/api/v1")
	
	if ar.server.middleware != nil {
		v1.Use(ar.server.middleware.API()...)
	}
	
	// 公共路由（无需认证）
	public := v1.Group("/public")
	ar.setupPublicRoutes(public)
	
	// 认证路由
	auth := v1.Group("/auth")
	ar.setupAuthRoutes(auth)
	
	// 用户路由（需要认证）
	user := v1.Group("/user")
	if ar.server.middleware != nil {
		user.Use(ar.server.middleware.JWT())
	}
	ar.setupUserRoutes(user)
	
	// 管理员路由（需要管理员权限）
	admin := v1.Group("/admin")
	if ar.server.middleware != nil {
		admin.Use(ar.server.middleware.Admin()...)
	}
	ar.setupAdminRoutes(admin)
}

// setupPublicRoutes 设置公共路由
func (ar *APIRouter) setupPublicRoutes(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"timestamp": gin.H{
				"unix": gin.H{
					"seconds": gin.H{
						"value": "1234567890",
					},
				},
			},
		})
	})
	
	router.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "1.0.0",
			"name":    "hwhkit-go",
		})
	})
}

// setupAuthRoutes 设置认证路由
func (ar *APIRouter) setupAuthRoutes(router *gin.RouterGroup) {
	router.POST("/login", ar.loginHandler)
	router.POST("/register", ar.registerHandler)
	router.POST("/refresh", ar.refreshTokenHandler)
	
	// 需要认证的认证路由
	authed := router.Group("/")
	if ar.server.middleware != nil {
		authed.Use(ar.server.middleware.JWT())
	}
	authed.POST("/logout", ar.logoutHandler)
	authed.GET("/me", ar.meHandler)
}

// setupUserRoutes 设置用户路由
func (ar *APIRouter) setupUserRoutes(router *gin.RouterGroup) {
	router.GET("/profile", ar.getUserProfileHandler)
	router.PUT("/profile", ar.updateUserProfileHandler)
	router.POST("/change-password", ar.changePasswordHandler)
}

// setupAdminRoutes 设置管理员路由
func (ar *APIRouter) setupAdminRoutes(router *gin.RouterGroup) {
	router.GET("/users", ar.listUsersHandler)
	router.GET("/users/:id", ar.getUserHandler)
	router.PUT("/users/:id", ar.updateUserHandler)
	router.DELETE("/users/:id", ar.deleteUserHandler)
	
	router.GET("/stats", ar.getStatsHandler)
	router.GET("/logs", ar.getLogsHandler)
}

// 认证相关处理器
func (ar *APIRouter) loginHandler(c *gin.Context) {
	// TODO: 实现登录逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "Login endpoint - implementation needed",
	})
}

func (ar *APIRouter) registerHandler(c *gin.Context) {
	// TODO: 实现注册逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "Register endpoint - implementation needed",
	})
}

func (ar *APIRouter) refreshTokenHandler(c *gin.Context) {
	// TODO: 实现刷新令牌逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "Refresh token endpoint - implementation needed",
	})
}

func (ar *APIRouter) logoutHandler(c *gin.Context) {
	// TODO: 实现登出逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout endpoint - implementation needed",
	})
}

func (ar *APIRouter) meHandler(c *gin.Context) {
	// TODO: 实现获取当前用户信息逻辑
	if claims, exists := middleware.GetClaims(c); exists {
		c.JSON(http.StatusOK, gin.H{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"email":    claims.Email,
			"role":     claims.Role,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No user information found",
		})
	}
}

// 用户相关处理器
func (ar *APIRouter) getUserProfileHandler(c *gin.Context) {
	// TODO: 实现获取用户档案逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "Get user profile endpoint - implementation needed",
	})
}

func (ar *APIRouter) updateUserProfileHandler(c *gin.Context) {
	// TODO: 实现更新用户档案逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "Update user profile endpoint - implementation needed",
	})
}

func (ar *APIRouter) changePasswordHandler(c *gin.Context) {
	// TODO: 实现修改密码逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "Change password endpoint - implementation needed",
	})
}

// 管理员相关处理器
func (ar *APIRouter) listUsersHandler(c *gin.Context) {
	// TODO: 实现列出用户逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "List users endpoint - implementation needed",
		"users":   []interface{}{},
	})
}

func (ar *APIRouter) getUserHandler(c *gin.Context) {
	// TODO: 实现获取指定用户逻辑
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Get user endpoint - implementation needed",
		"user_id": userID,
	})
}

func (ar *APIRouter) updateUserHandler(c *gin.Context) {
	// TODO: 实现更新用户逻辑
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Update user endpoint - implementation needed",
		"user_id": userID,
	})
}

func (ar *APIRouter) deleteUserHandler(c *gin.Context) {
	// TODO: 实现删除用户逻辑
	userID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete user endpoint - implementation needed",
		"user_id": userID,
	})
}

func (ar *APIRouter) getStatsHandler(c *gin.Context) {
	// TODO: 实现获取统计信息逻辑
	stats := gin.H{
		"total_users":    0,
		"active_users":   0,
		"total_requests": 0,
	}
	
	// 如果有数据库管理器，可以添加数据库统计
	if ar.server.db != nil {
		stats["database"] = ar.server.db.GetStats()
	}
	
	// 如果有缓存管理器，可以添加缓存统计
	if ar.server.cache != nil {
		stats["cache"] = ar.server.cache.GetStats()
	}
	
	c.JSON(http.StatusOK, stats)
}

func (ar *APIRouter) getLogsHandler(c *gin.Context) {
	// TODO: 实现获取日志逻辑
	c.JSON(http.StatusOK, gin.H{
		"message": "Get logs endpoint - implementation needed",
		"logs":    []interface{}{},
	})
}

// Builder 路由构建器
type Builder struct {
	server *Server
	groups []RouteGroup
}

// NewBuilder 创建路由构建器
func NewBuilder(server *Server) *Builder {
	return &Builder{
		server: server,
		groups: make([]RouteGroup, 0),
	}
}

// AddGroup 添加路由组
func (b *Builder) AddGroup(group RouteGroup) *Builder {
	b.groups = append(b.groups, group)
	return b
}

// Build 构建并注册所有路由
func (b *Builder) Build() {
	rm := NewRouterManager(b.server)
	for _, group := range b.groups {
		rm.RegisterRouteGroup(group)
	}
}