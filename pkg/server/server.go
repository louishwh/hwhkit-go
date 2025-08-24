package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hwh/hwhkit-go/pkg/auth"
	"github.com/hwh/hwhkit-go/pkg/cache"
	"github.com/hwh/hwhkit-go/pkg/config"
	"github.com/hwh/hwhkit-go/pkg/database"
	"github.com/hwh/hwhkit-go/pkg/logger"
	"github.com/hwh/hwhkit-go/pkg/middleware"
)

// Server HTTP服务器
type Server struct {
	engine      *gin.Engine
	httpServer  *http.Server
	config      *config.Config
	logger      *logger.Manager
	db          *database.Manager
	cache       *cache.Manager
	auth        *auth.Manager
	middleware  *middleware.MiddlewareManager
}

// ServerConfig 服务器配置选项
type ServerConfig struct {
	Config     *config.Config
	Logger     *logger.Manager
	Database   *database.Manager
	Cache      *cache.Manager
	Auth       *auth.Manager
}

// New 创建新的HTTP服务器
func New(cfg *ServerConfig) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("server config is required")
	}
	
	if cfg.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	
	// 设置Gin模式
	gin.SetMode(cfg.Config.Server.Mode)
	
	// 创建Gin引擎
	engine := gin.New()
	
	// 创建服务器实例
	server := &Server{
		engine: engine,
		config: cfg.Config,
		logger: cfg.Logger,
		db:     cfg.Database,
		cache:  cfg.Cache,
		auth:   cfg.Auth,
	}
	
	// 创建中间件管理器
	if cfg.Auth != nil && cfg.Logger != nil {
		server.middleware = middleware.NewMiddlewareManager(cfg.Auth, cfg.Logger)
	}
	
	// 配置HTTP服务器
	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Config.Server.Host, cfg.Config.Server.Port),
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Config.Server.WriteTimeout) * time.Second,
	}
	
	// 设置默认中间件
	server.setupDefaultMiddlewares()
	
	// 设置基础路由
	server.setupBasicRoutes()
	
	return server, nil
}

// setupDefaultMiddlewares 设置默认中间件
func (s *Server) setupDefaultMiddlewares() {
	// 如果有中间件管理器，使用它
	if s.middleware != nil {
		for _, mw := range s.middleware.Common() {
			s.engine.Use(mw)
		}
	} else {
		// 使用基础中间件
		s.engine.Use(gin.Recovery())
		if s.logger != nil {
			s.engine.Use(middleware.LoggerWithManager(s.logger))
		}
		if s.config.Server.EnableCORS {
			s.engine.Use(middleware.CORS())
		}
	}
}

// setupBasicRoutes 设置基础路由
func (s *Server) setupBasicRoutes() {
	// 健康检查
	s.engine.GET("/health", s.healthHandler)
	s.engine.GET("/health/live", s.livenessHandler)
	s.engine.GET("/health/ready", s.readinessHandler)
	
	// 信息路由
	s.engine.GET("/info", s.infoHandler)
	
	// 指标路由（如果需要）
	s.engine.GET("/metrics", s.metricsHandler)
}

// GetEngine 获取Gin引擎
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

// GetConfig 获取配置
func (s *Server) GetConfig() *config.Config {
	return s.config
}

// GetLogger 获取日志管理器
func (s *Server) GetLogger() *logger.Manager {
	return s.logger
}

// GetDatabase 获取数据库管理器
func (s *Server) GetDatabase() *database.Manager {
	return s.db
}

// GetCache 获取缓存管理器
func (s *Server) GetCache() *cache.Manager {
	return s.cache
}

// GetAuth 获取认证管理器
func (s *Server) GetAuth() *auth.Manager {
	return s.auth
}

// GetMiddleware 获取中间件管理器
func (s *Server) GetMiddleware() *middleware.MiddlewareManager {
	return s.middleware
}

// Group 创建路由组
func (s *Server) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return s.engine.Group(relativePath, handlers...)
}

// Use 添加中间件
func (s *Server) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return s.engine.Use(middleware...)
}

// GET 注册GET路由
func (s *Server) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.engine.GET(relativePath, handlers...)
}

// POST 注册POST路由
func (s *Server) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.engine.POST(relativePath, handlers...)
}

// PUT 注册PUT路由
func (s *Server) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.engine.PUT(relativePath, handlers...)
}

// PATCH 注册PATCH路由
func (s *Server) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.engine.PATCH(relativePath, handlers...)
}

// DELETE 注册DELETE路由
func (s *Server) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.engine.DELETE(relativePath, handlers...)
}

// OPTIONS 注册OPTIONS路由
func (s *Server) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.engine.OPTIONS(relativePath, handlers...)
}

// HEAD 注册HEAD路由
func (s *Server) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.engine.HEAD(relativePath, handlers...)
}

// Any 注册任意方法路由
func (s *Server) Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return s.engine.Any(relativePath, handlers...)
}

// Static 注册静态文件服务
func (s *Server) Static(relativePath, root string) gin.IRoutes {
	return s.engine.Static(relativePath, root)
}

// StaticFile 注册单个静态文件
func (s *Server) StaticFile(relativePath, filepath string) gin.IRoutes {
	return s.engine.StaticFile(relativePath, filepath)
}

// StaticFS 注册文件系统
func (s *Server) StaticFS(relativePath string, fs http.FileSystem) gin.IRoutes {
	return s.engine.StaticFS(relativePath, fs)
}

// LoadHTMLGlob 加载HTML模板
func (s *Server) LoadHTMLGlob(pattern string) {
	s.engine.LoadHTMLGlob(pattern)
}

// LoadHTMLFiles 加载HTML文件
func (s *Server) LoadHTMLFiles(files ...string) {
	s.engine.LoadHTMLFiles(files...)
}

// Start 启动服务器
func (s *Server) Start() error {
	if s.logger != nil {
		s.logger.Infof("Starting server on %s", s.httpServer.Addr)
	}
	
	// 在goroutine中启动服务器
	errChan := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start server: %w", err)
		}
	}()
	
	// 等待服务器启动或错误
	select {
	case err := <-errChan:
		return err
	case <-time.After(100 * time.Millisecond):
		if s.logger != nil {
			s.logger.Infof("Server started successfully on %s", s.httpServer.Addr)
		}
		return nil
	}
}

// StartWithGracefulShutdown 启动服务器并支持优雅关闭
func (s *Server) StartWithGracefulShutdown() error {
	// 启动服务器
	if err := s.Start(); err != nil {
		return err
	}
	
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	if s.logger != nil {
		s.logger.Info("Shutting down server...")
	}
	
	// 优雅关闭
	return s.Shutdown()
}

// Shutdown 关闭服务器
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if s.logger != nil {
		s.logger.Info("Server shutdown initiated")
	}
	
	// 关闭HTTP服务器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		if s.logger != nil {
			s.logger.Errorf("Server forced to shutdown: %v", err)
		}
		return err
	}
	
	// 关闭数据库连接
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			if s.logger != nil {
				s.logger.Errorf("Failed to close database: %v", err)
			}
		}
	}
	
	// 关闭缓存连接
	if s.cache != nil {
		if err := s.cache.Close(); err != nil {
			if s.logger != nil {
				s.logger.Errorf("Failed to close cache: %v", err)
			}
		}
	}
	
	if s.logger != nil {
		s.logger.Info("Server shutdown completed")
	}
	
	return nil
}

// health handler
func (s *Server) healthHandler(c *gin.Context) {
	status := gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0", // 可以从配置中获取
	}
	
	// 检查数据库健康状态
	if s.db != nil {
		if err := s.db.Health(); err != nil {
			status["database"] = "error"
			status["database_error"] = err.Error()
		} else {
			status["database"] = "ok"
		}
	}
	
	// 检查缓存健康状态
	if s.cache != nil {
		if err := s.cache.Health(); err != nil {
			status["cache"] = "error"
			status["cache_error"] = err.Error()
		} else {
			status["cache"] = "ok"
		}
	}
	
	// 如果有组件错误，返回503
	if status["database"] == "error" || status["cache"] == "error" {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}
	
	c.JSON(http.StatusOK, status)
}

// liveness handler
func (s *Server) livenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
		"timestamp": time.Now().Unix(),
	})
}

// readiness handler
func (s *Server) readinessHandler(c *gin.Context) {
	ready := true
	status := gin.H{
		"status": "ready",
		"timestamp": time.Now().Unix(),
		"checks": gin.H{},
	}
	
	// 检查数据库
	if s.db != nil {
		if err := s.db.Health(); err != nil {
			ready = false
			status["checks"].(gin.H)["database"] = gin.H{
				"status": "not_ready",
				"error":  err.Error(),
			}
		} else {
			status["checks"].(gin.H)["database"] = gin.H{
				"status": "ready",
			}
		}
	}
	
	// 检查缓存
	if s.cache != nil {
		if err := s.cache.Health(); err != nil {
			ready = false
			status["checks"].(gin.H)["cache"] = gin.H{
				"status": "not_ready",
				"error":  err.Error(),
			}
		} else {
			status["checks"].(gin.H)["cache"] = gin.H{
				"status": "ready",
			}
		}
	}
	
	if !ready {
		status["status"] = "not_ready"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}
	
	c.JSON(http.StatusOK, status)
}

// info handler
func (s *Server) infoHandler(c *gin.Context) {
	info := gin.H{
		"name":        "hwhkit-go",
		"version":     "1.0.0",
		"environment": s.config.Server.Mode,
		"timestamp":   time.Now().Unix(),
		"uptime":      time.Since(time.Now()).String(), // 这里应该用服务器启动时间
	}
	
	c.JSON(http.StatusOK, info)
}

// metrics handler
func (s *Server) metricsHandler(c *gin.Context) {
	metrics := gin.H{
		"timestamp": time.Now().Unix(),
	}
	
	// 添加数据库统计
	if s.db != nil {
		metrics["database"] = s.db.GetStats()
	}
	
	// 添加缓存统计
	if s.cache != nil {
		metrics["cache"] = s.cache.GetStats()
	}
	
	c.JSON(http.StatusOK, metrics)
}