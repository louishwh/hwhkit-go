# HWHKit-Go

一个功能齐全的Go工具库，提供了Web开发中常用的各种功能模块，包括配置管理、数据库操作、缓存、认证、中间件、日志记录等。

## 特性

- 🔧 **配置管理**: 支持环境变量和远程配置
- 🗄️ **数据库支持**: MySQL, PostgreSQL with GORM
- 🚀 **缓存支持**: Redis with connection pooling
- 🔐 **JWT认证**: 完整的JWT令牌管理
- 🌐 **HTTP服务器**: 基于Gin的HTTP服务器封装
- 🛡️ **中间件**: CORS, 认证, 日志, 限流等
- 📝 **日志管理**: 基于logrus的结构化日志
- 🔨 **工具函数**: 字符串, JSON, 时间, HTTP等工具

## 快速开始

### 安装

```bash
go mod init your-project
go get github.com/hwh/hwhkit-go
```

### 基本使用

```go
package main

import (
    "github.com/hwh/hwhkit-go/pkg/config"
    "github.com/hwh/hwhkit-go/pkg/logger"
    "github.com/hwh/hwhkit-go/pkg/server"
)

func main() {
    // 创建配置管理器
    configManager := config.New()
    cfg := configManager.Get()
    
    // 创建日志管理器
    logManager, _ := logger.New(configManager.GetLog())
    
    // 创建服务器
    serverConfig := &server.ServerConfig{
        Config: cfg,
        Logger: logManager,
    }
    
    httpServer, _ := server.New(serverConfig)
    
    // 添加路由
    httpServer.GET("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello, World!"})
    })
    
    // 启动服务器
    httpServer.StartWithGracefulShutdown()
}
```

## 模块详解

### 1. 配置管理 (Config)

支持从环境变量和.env文件加载配置，也支持远程配置API。

```go
// 基本使用
configManager := config.New()
cfg := configManager.Get()

// 使用远程配置
configManager := config.NewWithURL("https://config-api.example.com/config")
```

**支持的配置项:**
- 服务器配置 (端口, 模式, 超时等)
- 数据库配置 (MySQL, PostgreSQL)
- Redis配置
- JWT配置
- 日志配置

### 2. 日志管理 (Logger)

基于logrus的结构化日志管理，支持多种输出格式和目标。

```go
// 创建日志管理器
logManager, err := logger.New(&config.LogConfig{
    Level:  "info",
    Format: "json",
    Output: "both", // console, file, both
})

// 使用日志
logManager.Info("Application started")
logManager.WithFields(logger.Fields{
    "user_id": 123,
    "action":  "login",
}).Info("User logged in")

// 错误日志
logManager.WithError(err).Error("Database connection failed")
```

### 3. 数据库管理 (Database)

基于GORM的数据库管理，支持MySQL和PostgreSQL。

```go
// 创建数据库管理器
dbManager, err := database.New(&config.DatabaseConfig{
    Type:     "mysql",
    Host:     "localhost",
    Port:     3306,
    User:     "root",
    Password: "password",
    Name:     "myapp",
})

// 获取GORM实例
db := dbManager.GetDB()

// 自动迁移
dbManager.Migrate(&User{}, &Post{})

// 事务操作
err := dbManager.Transaction(func(tx *gorm.DB) error {
    // 事务内的操作
    return nil
})
```

### 4. 缓存管理 (Cache)

Redis缓存管理，支持连接池和各种数据类型操作。

```go
// 创建缓存管理器
cacheManager, err := cache.New(&config.RedisConfig{
    Host:     "localhost",
    Port:     6379,
    Password: "",
    DB:       0,
})

// 基本操作
cacheManager.Set("key", "value", time.Hour)
value, err := cacheManager.Get("key")

// JSON操作
cacheManager.SetJSON("user:123", user, time.Hour)
var user User
err := cacheManager.GetJSON("user:123", &user)

// 列表操作
cacheManager.LPush("queue", "item1", "item2")
item, err := cacheManager.RPop("queue")
```

### 5. JWT认证 (Auth)

完整的JWT令牌管理系统。

```go
// 创建认证管理器
authManager := auth.New(&config.JWTConfig{
    Secret:       "your-secret-key",
    ExpireHours:  24,
    RefreshHours: 168,
    Issuer:       "your-app",
})

// 生成令牌对
tokenPair, err := authManager.GenerateTokenPair(
    userID,
    username,
    email,
    role,
)

// 验证令牌
claims, err := authManager.ValidateToken(token)

// 刷新令牌
newTokenPair, err := authManager.RefreshToken(refreshToken)
```

### 6. 中间件 (Middleware)

提供各种Gin中间件。

```go
import "github.com/hwh/hwhkit-go/pkg/middleware"

// CORS中间件
engine.Use(middleware.CORS())

// JWT认证中间件
engine.Use(middleware.JWTWithManager(authManager))

// 日志中间件
engine.Use(middleware.LoggerWithManager(logManager))

// 限流中间件
engine.Use(middleware.RateLimitByIP(100, 10)) // 100 req/s, burst 10

// 角色验证
adminRoutes := engine.Group("/admin")
adminRoutes.Use(middleware.RequireRole(authManager, "admin"))
```

### 7. HTTP服务器 (Server)

基于Gin的HTTP服务器封装。

```go
// 创建服务器
serverConfig := &server.ServerConfig{
    Config:   cfg,
    Logger:   logManager,
    Database: dbManager,
    Cache:    cacheManager,
    Auth:     authManager,
}

httpServer, err := server.New(serverConfig)

// 设置API路由
apiRouter := server.NewAPIRouter(httpServer)
apiRouter.SetupV1API()

// 添加自定义路由
httpServer.GET("/custom", handler)

// 路由组
api := httpServer.Group("/api")
api.GET("/users", getUsersHandler)

// 启动服务器
httpServer.StartWithGracefulShutdown()
```

### 8. 工具函数 (Utils)

提供各种常用工具函数。

```go
import "github.com/hwh/hwhkit-go/pkg/utils"

// 字符串工具
utils.Str.CamelCase("hello_world")    // "helloWorld"
utils.Str.SnakeCase("HelloWorld")     // "hello_world"
utils.Str.RandomString(10)            // 随机字符串

// JSON工具
jsonStr, _ := utils.JSON.ToPrettyJSON(data)
utils.JSON.FromJSON(jsonStr, &result)

// 时间工具
utils.Time.FormatNowDateTime()        // "2023-12-25 10:30:45"
utils.Time.AddDays(time.Now(), 7)     // 7天后

// HTTP工具
resp, err := utils.HTTP.Get("https://api.example.com", headers)
```

## 环境配置

复制 `.env.example` 到 `.env` 并根据需要修改配置:

```bash
cp .env.example .env
```

主要配置项：

```env
# 服务器
SERVER_PORT=8080
SERVER_MODE=debug

# 数据库
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_NAME=myapp

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
```

## API文档

### 健康检查

- `GET /health` - 整体健康状态
- `GET /health/live` - 存活检查
- `GET /health/ready` - 就绪检查

### 认证 API

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/refresh` - 刷新令牌
- `POST /api/v1/auth/logout` - 用户登出

### 用户 API

- `GET /api/v1/user/profile` - 获取用户资料
- `PUT /api/v1/user/profile` - 更新用户资料

### 管理员 API

- `GET /api/v1/admin/users` - 获取用户列表
- `GET /api/v1/admin/stats` - 获取统计信息

## 项目结构

```
hwhkit-go/
├── pkg/                    # 核心包
│   ├── auth/              # JWT认证
│   ├── cache/             # Redis缓存
│   ├── config/            # 配置管理
│   ├── database/          # 数据库管理
│   ├── logger/            # 日志管理
│   ├── middleware/        # Gin中间件
│   ├── server/            # HTTP服务器
│   └── utils/             # 工具函数
├── examples/              # 示例代码
│   └── basic/            # 基本使用示例
├── tests/                 # 测试文件
├── docs/                  # 文档
├── go.mod
├── go.sum
└── README.md
```

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License