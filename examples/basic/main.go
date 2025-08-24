package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hwh/hwhkit-go/pkg/auth"
	"github.com/hwh/hwhkit-go/pkg/cache"
	"github.com/hwh/hwhkit-go/pkg/config"
	"github.com/hwh/hwhkit-go/pkg/database"
	"github.com/hwh/hwhkit-go/pkg/logger"
	"github.com/hwh/hwhkit-go/pkg/middleware"
	"github.com/hwh/hwhkit-go/pkg/server"
	"github.com/hwh/hwhkit-go/pkg/utils"
)

func main() {
	// 1. 创建配置管理器
	configManager := config.New()
	cfg := configManager.Get()
	
	fmt.Printf("Server will run on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	
	// 2. 创建日志管理器
	logManager, err := logger.New(configManager.GetLog())
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	
	logManager.Info("Application starting...")
	
	// 3. 创建数据库管理器
	dbManager, err := database.New(configManager.GetDatabase())
	if err != nil {
		logManager.Fatalf("Failed to create database manager: %v", err)
	}
	
	logManager.Info("Database connected successfully")
	
	// 4. 创建缓存管理器
	cacheManager, err := cache.New(configManager.GetRedis())
	if err != nil {
		logManager.Warnf("Failed to create cache manager: %v", err)
		// 缓存不是必须的，继续运行
	} else {
		logManager.Info("Cache connected successfully")
	}
	
	// 5. 创建认证管理器
	authManager := auth.New(configManager.GetJWT())
	
	// 6. 创建HTTP服务器
	serverConfig := &server.ServerConfig{
		Config:   cfg,
		Logger:   logManager,
		Database: dbManager,
		Cache:    cacheManager,
		Auth:     authManager,
	}
	
	httpServer, err := server.New(serverConfig)
	if err != nil {
		logManager.Fatalf("Failed to create server: %v", err)
	}
	
	// 7. 设置API路由
	apiRouter := server.NewAPIRouter(httpServer)
	apiRouter.SetupV1API()
	
	// 8. 添加自定义路由
	setupCustomRoutes(httpServer, logManager)
	
	// 9. 启动服务器（支持优雅关闭）
	logManager.Info("Starting server with graceful shutdown support...")
	if err := httpServer.StartWithGracefulShutdown(); err != nil {
		logManager.Fatalf("Server failed: %v", err)
	}
	
	logManager.Info("Application stopped")
}

// setupCustomRoutes 设置自定义路由
func setupCustomRoutes(srv *server.Server, log *logger.Manager) {
	// 演示工具函数的使用
	srv.GET("/demo/utils", func(c *gin.Context) {
		// 字符串工具示例
		original := "hello_world"
		camelCase := utils.Str.CamelCase(original)
		snakeCase := utils.Str.SnakeCase("HelloWorld")
		randomStr := utils.Str.RandomString(10)
		
		// JSON工具示例
		data := map[string]interface{}{
			"name": "John",
			"age":  30,
		}
		jsonStr, _ := utils.JSON.ToPrettyJSON(data)
		
		// 时间工具示例
		now := utils.Time.Now()
		formatted := utils.Time.FormatNowDateTime()
		
		// HTTP工具示例
		httpUtil := utils.HTTP
		isValidURL := httpUtil.IsValidURL("https://example.com")
		
		response := gin.H{
			"string_examples": gin.H{
				"original":    original,
				"camel_case":  camelCase,
				"snake_case":  snakeCase,
				"random":      randomStr,
			},
			"json_example": jsonStr,
			"time_examples": gin.H{
				"now":       now,
				"formatted": formatted,
			},
			"http_examples": gin.H{
				"is_valid_url": isValidURL,
			},
		}
		
		c.JSON(200, response)
	})
	
	// 演示缓存的使用
	srv.GET("/demo/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		cacheManager := srv.GetCache()
		
		if cacheManager == nil {
			c.JSON(500, gin.H{"error": "Cache not available"})
			return
		}
		
		// 尝试从缓存获取
		value, err := cacheManager.Get(key)
		if err != nil {
			// 缓存中没有，设置一个值
			newValue := fmt.Sprintf("cached_value_for_%s_at_%s", key, utils.Time.FormatNowDateTime())
			err = cacheManager.Set(key, newValue, 300*time.Second) // 5分钟过期
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			value = newValue
		}
		
		c.JSON(200, gin.H{
			"key":   key,
			"value": value,
			"from":  "cache",
		})
	})
	
	// 演示JWT认证的使用
	srv.POST("/demo/auth/login", func(c *gin.Context) {
		// 简单的演示登录（实际项目中需要验证用户名密码）
		authManager := srv.GetAuth()
		
		// 生成令牌对
		tokenPair, err := authManager.GenerateTokenPair(123, "demo_user", "demo@example.com", "user")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(200, gin.H{
			"message": "Login successful",
			"tokens":  tokenPair,
		})
	})
	
	// 需要认证的演示路由
	protected := srv.Group("/demo/protected")
	if srv.GetMiddleware() != nil {
		protected.Use(srv.GetMiddleware().JWT())
	}
	
	protected.GET("/profile", func(c *gin.Context) {
		// 获取用户信息
		if claims, exists := middleware.GetClaims(c); exists {
			c.JSON(200, gin.H{
				"message": "This is a protected route",
				"user": gin.H{
					"id":       claims.UserID,
					"username": claims.Username,
					"email":    claims.Email,
					"role":     claims.Role,
				},
			})
		} else {
			c.JSON(401, gin.H{"error": "No user information found"})
		}
	})
	
	// 演示数据库的使用（需要先定义模型）
	srv.GET("/demo/db/stats", func(c *gin.Context) {
		dbManager := srv.GetDatabase()
		if dbManager == nil {
			c.JSON(500, gin.H{"error": "Database not available"})
			return
		}
		
		stats := dbManager.GetStats()
		c.JSON(200, gin.H{
			"message": "Database statistics",
			"stats":   stats,
		})
	})
	
	log.Info("Custom routes setup completed")
}