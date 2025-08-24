package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hwh/hwhkit-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         8080,
			Mode:         "test",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
	}
	
	serverConfig := &ServerConfig{
		Config: cfg,
	}
	
	server, err := New(serverConfig)
	require.NoError(t, err)
	assert.NotNil(t, server)
	assert.Equal(t, cfg, server.config)
	assert.NotNil(t, server.engine)
	assert.NotNil(t, server.httpServer)
}

func TestNewWithNilConfig(t *testing.T) {
	_, err := New(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server config is required")
}

func TestNewWithNilInnerConfig(t *testing.T) {
	serverConfig := &ServerConfig{}
	_, err := New(serverConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config is required")
}

func TestServerRoutes(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         8080,
			Mode:         gin.TestMode,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
	}
	
	serverConfig := &ServerConfig{
		Config: cfg,
	}
	
	server, err := New(serverConfig)
	require.NoError(t, err)
	
	// 测试健康检查路由
	t.Run("Health endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		server.engine.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
		assert.NotNil(t, response["timestamp"])
	})
	
	// 测试存活检查路由
	t.Run("Liveness endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health/live", nil)
		server.engine.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "alive", response["status"])
	})
	
	// 测试就绪检查路由
	t.Run("Readiness endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health/ready", nil)
		server.engine.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "ready", response["status"])
	})
	
	// 测试信息路由
	t.Run("Info endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info", nil)
		server.engine.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "hwhkit-go", response["name"])
		assert.Equal(t, gin.TestMode, response["environment"])
	})
	
	// 测试指标路由
	t.Run("Metrics endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/metrics", nil)
		server.engine.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotNil(t, response["timestamp"])
	})
}

func TestServerMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         8080,
			Mode:         gin.TestMode,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
	}
	
	serverConfig := &ServerConfig{
		Config: cfg,
	}
	
	server, err := New(serverConfig)
	require.NoError(t, err)
	
	// 测试添加路由
	server.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})
	
	server.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "created"})
	})
	
	// 测试GET路由
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	server.engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 测试POST路由
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", nil)
	server.engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestServerGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         8080,
			Mode:         gin.TestMode,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
	}
	
	serverConfig := &ServerConfig{
		Config: cfg,
	}
	
	server, err := New(serverConfig)
	require.NoError(t, err)
	
	// 创建路由组
	api := server.Group("/api")
	api.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"users": []string{}})
	})
	
	// 测试组路由
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users", nil)
	server.engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response["users"])
}

func TestAPIRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         8080,
			Mode:         gin.TestMode,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
	}
	
	serverConfig := &ServerConfig{
		Config: cfg,
	}
	
	server, err := New(serverConfig)
	require.NoError(t, err)
	
	// 设置API路由
	apiRouter := NewAPIRouter(server)
	apiRouter.SetupV1API()
	
	// 测试公共路由
	t.Run("Public ping endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/public/ping", nil)
		server.engine.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "pong", response["message"])
	})
	
	// 测试版本路由
	t.Run("Version endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/public/version", nil)
		server.engine.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "1.0.0", response["version"])
		assert.Equal(t, "hwhkit-go", response["name"])
	})
}

func TestRouteBuilder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         8080,
			Mode:         gin.TestMode,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
	}
	
	serverConfig := &ServerConfig{
		Config: cfg,
	}
	
	server, err := New(serverConfig)
	require.NoError(t, err)
	
	// 使用构建器创建路由
	builder := NewBuilder(server)
	builder.AddGroup(RouteGroup{
		Path: "/test",
		Routes: []Route{
			{
				Method: "GET",
				Path:   "/hello",
				Handlers: []gin.HandlerFunc{
					func(c *gin.Context) {
						c.JSON(http.StatusOK, gin.H{"message": "hello from builder"})
					},
				},
			},
		},
	})
	builder.Build()
	
	// 测试构建器创建的路由
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test/hello", nil)
	server.engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "hello from builder", response["message"])
}