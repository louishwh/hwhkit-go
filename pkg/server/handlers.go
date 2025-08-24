package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hwh/hwhkit-go/pkg/auth"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id"`
	Timestamp int64       `json:"timestamp"`
}

// PaginatedResponse 分页响应结构
type PaginatedResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	RequestID  string      `json:"request_id"`
	Timestamp  int64       `json:"timestamp"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

// Success 成功响应
func (s *Server) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		RequestID: c.GetString("request_id"),
		Timestamp: time.Now().Unix(),
	})
}

// Error 错误响应
func (s *Server) Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:      code,
		Message:   message,
		RequestID: c.GetString("request_id"),
		Timestamp: time.Now().Unix(),
	})
}

// PaginatedSuccess 分页成功响应
func (s *Server) PaginatedSuccess(c *gin.Context, data interface{}, total int64) {
	page := c.GetInt("page")
	pageSize := c.GetInt("page_size")
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	
	c.JSON(http.StatusOK, PaginatedResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
		RequestID: c.GetString("request_id"),
		Timestamp: time.Now().Unix(),
	})
}

// 路由处理器

// handleRoot 根路径处理器
func (s *Server) handleRoot(c *gin.Context) {
	s.Success(c, gin.H{
		"service": "hwhkit-go",
		"version": "1.0.0",
		"message": "Welcome to HWHKit-Go Enterprise Web Service Toolkit",
	})
}

// handleInfo 信息处理器
func (s *Server) handleInfo(c *gin.Context) {
	s.Success(c, gin.H{
		"service":     "hwhkit-go",
		"version":     "1.0.0",
		"description": "Enterprise Web Service Toolkit",
		"features": []string{
			"HTTP routing and middleware",
			"Database operations (MySQL/PostgreSQL)",
			"Redis cache and session storage",
			"JWT authentication",
			"RBAC authorization",
			"Template engine support",
			"CORS and API documentation",
			"Structured logging",
		},
	})
}

// handleHealth 健康检查处理器
func (s *Server) handleHealth(c *gin.Context) {
	health := gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"services":  gin.H{},
	}
	
	// 检查数据库
	if err := s.db.Health(); err != nil {
		health["services"].(gin.H)["database"] = gin.H{
			"status": "error",
			"error":  err.Error(),
		}
	} else {
		health["services"].(gin.H)["database"] = gin.H{
			"status": "ok",
		}
	}
	
	// 检查缓存
	if err := s.cache.Health(); err != nil {
		health["services"].(gin.H)["cache"] = gin.H{
			"status": "error",
			"error":  err.Error(),
		}
	} else {
		health["services"].(gin.H)["cache"] = gin.H{
			"status": "ok",
		}
	}
	
	s.Success(c, health)
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	Roles    []string `json:"roles"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	Email string `json:"email" binding:"email"`
	Name  string `json:"name"`
}

// handleLogin 登录处理器
func (s *Server) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	
	// 模拟用户提供器（实际应用中应该从数据库获取）
	userProvider := func(username string) (*auth.User, error) {
		// 这里应该从数据库查询用户
		// 为了演示，使用硬编码用户
		if username == "admin" {
			hashedPassword, _ := s.authService.GetPasswordManager().HashPassword("admin123")
			return &auth.User{
				ID:       "1",
				Username: "admin",
				Email:    "admin@example.com",
				Password: hashedPassword,
				Roles:    []string{"admin"},
				IsActive: true,
			}, nil
		}
		return nil, &AuthError{Message: "user not found"}
	}
	
	tokenPair, err := s.authService.Login(req.Username, req.Password, userProvider)
	if err != nil {
		s.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	
	s.Success(c, tokenPair)
}

// handleRegister 注册处理器
func (s *Server) handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	
	// 模拟用户创建器（实际应用中应该保存到数据库）
	userCreator := func(user *auth.User) error {
		// 这里应该保存用户到数据库
		user.ID = "generated-id-" + strconv.FormatInt(time.Now().Unix(), 10)
		s.logger.Infof("User created: %s", user.Username)
		return nil
	}
	
	roles := req.Roles
	if len(roles) == 0 {
		roles = []string{"user"}
	}
	
	tokenPair, err := s.authService.Register(req.Username, req.Email, req.Password, roles, userCreator)
	if err != nil {
		s.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	
	s.Success(c, tokenPair)
}

// handleRefreshToken 刷新令牌处理器
func (s *Server) handleRefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	
	tokenPair, err := s.authService.GetJWTManager().RefreshToken(req.RefreshToken)
	if err != nil {
		s.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	
	s.Success(c, tokenPair)
}

// handleProfile 获取个人资料处理器
func (s *Server) handleProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	username := c.GetString("username")
	email := c.GetString("user_email")
	roles := c.GetStringSlice("user_roles")
	
	profile := gin.H{
		"id":       userID,
		"username": username,
		"email":    email,
		"roles":    roles,
	}
	
	s.Success(c, profile)
}

// handleUpdateProfile 更新个人资料处理器
func (s *Server) handleUpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	
	userID := c.GetString("user_id")
	
	// 这里应该更新数据库中的用户信息
	s.logger.Infof("Updating profile for user %s", userID)
	
	s.Success(c, gin.H{
		"message": "Profile updated successfully",
	})
}

// handleLogout 登出处理器
func (s *Server) handleLogout(c *gin.Context) {
	// 在实际应用中，可能需要将令牌加入黑名单
	userID := c.GetString("user_id")
	s.logger.Infof("User %s logged out", userID)
	
	s.Success(c, gin.H{
		"message": "Logged out successfully",
	})
}

// handleChangePassword 修改密码处理器
func (s *Server) handleChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	
	userID := c.GetString("user_id")
	
	// 模拟用户提供器和更新器
	userProvider := func(id string) (*auth.User, error) {
		// 这里应该从数据库查询用户
		return &auth.User{
			ID:       id,
			Username: "user",
			Password: "hashed_old_password",
		}, nil
	}
	
	userUpdater := func(user *auth.User) error {
		// 这里应该更新数据库中的用户密码
		s.logger.Infof("Password updated for user %s", user.ID)
		return nil
	}
	
	err := s.authService.ChangePassword(userID, req.OldPassword, req.NewPassword, userProvider, userUpdater)
	if err != nil {
		s.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	
	s.Success(c, gin.H{
		"message": "Password changed successfully",
	})
}

// handleListUsers 列出用户处理器（管理员）
func (s *Server) handleListUsers(c *gin.Context) {
	// 使用分页中间件解析的参数
	page := c.GetInt("page")
	pageSize := c.GetInt("page_size")
	offset := c.GetInt("offset")
	
	// 模拟用户列表（实际应用中应该从数据库查询）
	users := []gin.H{
		{"id": "1", "username": "admin", "email": "admin@example.com", "roles": []string{"admin"}},
		{"id": "2", "username": "user1", "email": "user1@example.com", "roles": []string{"user"}},
		{"id": "3", "username": "user2", "email": "user2@example.com", "roles": []string{"user"}},
	}
	
	// 模拟分页
	total := int64(len(users))
	start := offset
	end := offset + pageSize
	if end > len(users) {
		end = len(users)
	}
	if start > len(users) {
		start = len(users)
	}
	
	paginatedUsers := users[start:end]
	
	s.logger.Infof("Listing users: page=%d, pageSize=%d, total=%d", page, pageSize, total)
	
	s.PaginatedSuccess(c, paginatedUsers, total)
}

// handleStats 统计信息处理器（管理员）
func (s *Server) handleStats(c *gin.Context) {
	dbStats := s.db.GetStats()
	cacheStats := s.cache.GetStats()
	
	stats := gin.H{
		"database": dbStats,
		"cache":    cacheStats,
		"server": gin.H{
			"uptime":  time.Since(time.Now()).String(), // 这里应该记录服务器启动时间
			"version": "1.0.0",
		},
	}
	
	s.Success(c, stats)
}

// AuthError 认证错误
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}