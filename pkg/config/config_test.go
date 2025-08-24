package config

import (
	"os"
	"testing"
)

func TestConfigManager_LoadFromLocal(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_TYPE", "postgres")
	os.Setenv("JWT_SECRET", "test-secret")
	
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_TYPE")
		os.Unsetenv("JWT_SECRET")
	}()
	
	cm := New()
	config := cm.Get()
	
	if config.Server.Port != 9090 {
		t.Errorf("Expected server port 9090, got %d", config.Server.Port)
	}
	
	if config.Database.Type != "postgres" {
		t.Errorf("Expected database type postgres, got %s", config.Database.Type)
	}
	
	if config.JWT.Secret != "test-secret" {
		t.Errorf("Expected JWT secret test-secret, got %s", config.JWT.Secret)
	}
}

func TestConfigManager_GetMethods(t *testing.T) {
	cm := New()
	
	// 测试各个Get方法
	serverConfig := cm.GetServer()
	if serverConfig == nil {
		t.Error("GetServer() returned nil")
	}
	
	dbConfig := cm.GetDatabase()
	if dbConfig == nil {
		t.Error("GetDatabase() returned nil")
	}
	
	redisConfig := cm.GetRedis()
	if redisConfig == nil {
		t.Error("GetRedis() returned nil")
	}
	
	jwtConfig := cm.GetJWT()
	if jwtConfig == nil {
		t.Error("GetJWT() returned nil")
	}
	
	logConfig := cm.GetLog()
	if logConfig == nil {
		t.Error("GetLog() returned nil")
	}
}

func TestGetEnvFunctions(t *testing.T) {
	// 测试 getEnv
	os.Setenv("TEST_STRING", "test_value")
	defer os.Unsetenv("TEST_STRING")
	
	value := getEnv("TEST_STRING", "default")
	if value != "test_value" {
		t.Errorf("Expected test_value, got %s", value)
	}
	
	defaultValue := getEnv("NON_EXISTENT", "default")
	if defaultValue != "default" {
		t.Errorf("Expected default, got %s", defaultValue)
	}
	
	// 测试 getEnvAsInt
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")
	
	intValue := getEnvAsInt("TEST_INT", 456)
	if intValue != 123 {
		t.Errorf("Expected 123, got %d", intValue)
	}
	
	defaultInt := getEnvAsInt("NON_EXISTENT_INT", 456)
	if defaultInt != 456 {
		t.Errorf("Expected 456, got %d", defaultInt)
	}
	
	// 测试 getEnvAsBool
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")
	
	boolValue := getEnvAsBool("TEST_BOOL", false)
	if !boolValue {
		t.Error("Expected true, got false")
	}
	
	defaultBool := getEnvAsBool("NON_EXISTENT_BOOL", false)
	if defaultBool {
		t.Error("Expected false, got true")
	}
}