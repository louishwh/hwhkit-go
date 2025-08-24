package server

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// TemplateManager 模板管理器
type TemplateManager struct {
	templateDir string
	funcMap     template.FuncMap
}

// NewTemplateManager 创建模板管理器
func NewTemplateManager(templateDir string) *TemplateManager {
	tm := &TemplateManager{
		templateDir: templateDir,
		funcMap:     make(template.FuncMap),
	}
	
	// 添加默认模板函数
	tm.addDefaultFunctions()
	
	return tm
}

// addDefaultFunctions 添加默认模板函数
func (tm *TemplateManager) addDefaultFunctions() {
	tm.funcMap["formatDate"] = func(t time.Time, format string) string {
		return t.Format(format)
	}
	
	tm.funcMap["formatDateTime"] = func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	}
	
	tm.funcMap["formatTime"] = func(t time.Time) string {
		return t.Format("15:04:05")
	}
	
	tm.funcMap["now"] = func() time.Time {
		return time.Now()
	}
	
	tm.funcMap["add"] = func(a, b int) int {
		return a + b
	}
	
	tm.funcMap["sub"] = func(a, b int) int {
		return a - b
	}
	
	tm.funcMap["mul"] = func(a, b int) int {
		return a * b
	}
	
	tm.funcMap["div"] = func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a / b
	}
	
	tm.funcMap["eq"] = func(a, b interface{}) bool {
		return a == b
	}
	
	tm.funcMap["ne"] = func(a, b interface{}) bool {
		return a != b
	}
	
	tm.funcMap["lt"] = func(a, b int) bool {
		return a < b
	}
	
	tm.funcMap["le"] = func(a, b int) bool {
		return a <= b
	}
	
	tm.funcMap["gt"] = func(a, b int) bool {
		return a > b
	}
	
	tm.funcMap["ge"] = func(a, b int) bool {
		return a >= b
	}
	
	tm.funcMap["contains"] = func(s, substr string) bool {
		return strings.Contains(s, substr)
	}
	
	tm.funcMap["upper"] = func(s string) string {
		return strings.ToUpper(s)
	}
	
	tm.funcMap["lower"] = func(s string) string {
		return strings.ToLower(s)
	}
	
	tm.funcMap["title"] = func(s string) string {
		return strings.Title(s)
	}
	
	tm.funcMap["trim"] = func(s string) string {
		return strings.TrimSpace(s)
	}
	
	tm.funcMap["len"] = func(v interface{}) int {
		switch v := v.(type) {
		case string:
			return len(v)
		case []interface{}:
			return len(v)
		case map[string]interface{}:
			return len(v)
		default:
			return 0
		}
	}
	
	tm.funcMap["slice"] = func(start, end int, items []interface{}) []interface{} {
		if start < 0 {
			start = 0
		}
		if end > len(items) {
			end = len(items)
		}
		if start >= end {
			return []interface{}{}
		}
		return items[start:end]
	}
	
	tm.funcMap["range"] = func(start, end int) []int {
		if start >= end {
			return []int{}
		}
		result := make([]int, end-start)
		for i := range result {
			result[i] = start + i
		}
		return result
	}
	
	tm.funcMap["default"] = func(defaultValue, value interface{}) interface{} {
		if value == nil || value == "" {
			return defaultValue
		}
		return value
	}
}

// AddFunction 添加自定义模板函数
func (tm *TemplateManager) AddFunction(name string, fn interface{}) {
	tm.funcMap[name] = fn
}

// GetFuncMap 获取函数映射
func (tm *TemplateManager) GetFuncMap() template.FuncMap {
	return tm.funcMap
}

// LoadTemplates 加载模板
func (tm *TemplateManager) LoadTemplates(engine *gin.Engine) {
	engine.SetFuncMap(tm.funcMap)
	
	// 加载所有模板文件
	pattern := filepath.Join(tm.templateDir, "**/*")
	engine.LoadHTMLGlob(pattern)
}

// SetupTemplateRoutes 设置模板路由
func (s *Server) SetupTemplateRoutes() {
	// 创建模板管理器
	templateManager := NewTemplateManager(s.config.Server.TemplateDir)
	
	// 加载模板
	templateManager.LoadTemplates(s.engine)
	
	// 设置静态文件服务
	s.engine.Static("/static", s.config.Server.StaticDir)
	
	// 前后端不分离的页面路由
	pages := s.engine.Group("/")
	{
		pages.GET("/", s.handleHomePage)
		pages.GET("/login", s.handleLoginPage)
		pages.GET("/register", s.handleRegisterPage)
		pages.GET("/dashboard", s.handleDashboardPage)
		pages.GET("/profile", s.handleProfilePage)
	}
	
	// 表单处理路由
	forms := s.engine.Group("/forms")
	{
		forms.POST("/login", s.handleLoginForm)
		forms.POST("/register", s.handleRegisterForm)
		forms.POST("/logout", s.handleLogoutForm)
	}
}

// 页面处理器

// handleHomePage 首页处理器
func (s *Server) handleHomePage(c *gin.Context) {
	data := gin.H{
		"title":   "HWHKit-Go - Enterprise Web Service Toolkit",
		"message": "Welcome to HWHKit-Go",
		"year":    time.Now().Year(),
	}
	
	c.HTML(http.StatusOK, "index.html", data)
}

// handleLoginPage 登录页面处理器
func (s *Server) handleLoginPage(c *gin.Context) {
	data := gin.H{
		"title": "Login - HWHKit-Go",
		"error": c.Query("error"),
	}
	
	c.HTML(http.StatusOK, "auth/login.html", data)
}

// handleRegisterPage 注册页面处理器
func (s *Server) handleRegisterPage(c *gin.Context) {
	data := gin.H{
		"title": "Register - HWHKit-Go",
		"error": c.Query("error"),
	}
	
	c.HTML(http.StatusOK, "auth/register.html", data)
}

// handleDashboardPage 仪表板页面处理器
func (s *Server) handleDashboardPage(c *gin.Context) {
	// 检查用户是否已登录（通过session或cookie）
	userID := s.getUserFromSession(c)
	if userID == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	data := gin.H{
		"title":  "Dashboard - HWHKit-Go",
		"user":   s.getUserInfo(userID),
		"stats":  s.getDashboardStats(),
	}
	
	c.HTML(http.StatusOK, "dashboard.html", data)
}

// handleProfilePage 个人资料页面处理器
func (s *Server) handleProfilePage(c *gin.Context) {
	userID := s.getUserFromSession(c)
	if userID == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	data := gin.H{
		"title": "Profile - HWHKit-Go",
		"user":  s.getUserInfo(userID),
	}
	
	c.HTML(http.StatusOK, "profile.html", data)
}

// 表单处理器

// handleLoginForm 登录表单处理器
func (s *Server) handleLoginForm(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	
	// 验证用户
	if s.authenticateUser(username, password) {
		// 创建会话
		sessionID := s.createUserSession(username)
		c.SetCookie("session_id", sessionID, 3600*24, "/", "", false, true)
		
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}
	
	c.Redirect(http.StatusFound, "/login?error=invalid_credentials")
}

// handleRegisterForm 注册表单处理器
func (s *Server) handleRegisterForm(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")
	
	// 验证输入
	if password != confirmPassword {
		c.Redirect(http.StatusFound, "/register?error=passwords_mismatch")
		return
	}
	
	// 创建用户
	if err := s.createUser(username, email, password); err != nil {
		c.Redirect(http.StatusFound, "/register?error=user_creation_failed")
		return
	}
	
	c.Redirect(http.StatusFound, "/login?message=registration_successful")
}

// handleLogoutForm 登出表单处理器
func (s *Server) handleLogoutForm(c *gin.Context) {
	sessionID, _ := c.Cookie("session_id")
	if sessionID != "" {
		s.destroyUserSession(sessionID)
		c.SetCookie("session_id", "", -1, "/", "", false, true)
	}
	
	c.Redirect(http.StatusFound, "/")
}

// 辅助方法

// getUserFromSession 从会话获取用户
func (s *Server) getUserFromSession(c *gin.Context) string {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		return ""
	}
	
	// 从缓存中获取会话信息
	userID, err := s.cache.Get("session:" + sessionID)
	if err != nil {
		return ""
	}
	
	return userID
}

// createUserSession 创建用户会话
func (s *Server) createUserSession(username string) string {
	sessionID := generateSessionID()
	
	// 在实际应用中，这里应该存储用户ID而不是用户名
	s.cache.Set("session:"+sessionID, username, time.Hour*24)
	
	return sessionID
}

// destroyUserSession 销毁用户会话
func (s *Server) destroyUserSession(sessionID string) {
	s.cache.Delete("session:" + sessionID)
}

// authenticateUser 验证用户
func (s *Server) authenticateUser(username, password string) bool {
	// 这里应该从数据库验证用户
	// 为了演示，使用硬编码
	return username == "admin" && password == "admin123"
}

// createUser 创建用户
func (s *Server) createUser(username, email, password string) error {
	// 这里应该将用户信息保存到数据库
	s.logger.Infof("Creating user: %s (%s)", username, email)
	return nil
}

// getUserInfo 获取用户信息
func (s *Server) getUserInfo(userID string) gin.H {
	// 这里应该从数据库获取用户信息
	return gin.H{
		"id":       userID,
		"username": userID, // 简化处理
		"email":    userID + "@example.com",
	}
}

// getDashboardStats 获取仪表板统计
func (s *Server) getDashboardStats() gin.H {
	return gin.H{
		"total_users":    100,
		"active_users":   50,
		"total_requests": 1000,
		"response_time":  "25ms",
	}
}

// generateSessionID 生成会话ID
func generateSessionID() string {
	return generateRequestID() // 复用请求ID生成逻辑
}

import (
	"net/http"
	"strings"
	"time"
)