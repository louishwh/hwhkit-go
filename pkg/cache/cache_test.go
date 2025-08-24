package cache

import (
	"testing"
	"time"

	"github.com/hwh/hwhkit-go/pkg/config"
)

func TestCacheManager(t *testing.T) {
	// 注意：这个测试需要实际的Redis连接
	// 在CI/CD环境中，你可能需要跳过这个测试或使用模拟Redis
	t.Skip("Skipping cache integration test - requires actual Redis")
	
	cfg := &config.RedisConfig{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5,
		ReadTimeout:  3,
		WriteTimeout: 3,
	}
	
	manager, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	defer manager.Close()
	
	// 测试健康检查
	if err := manager.Health(); err != nil {
		t.Errorf("Cache health check failed: %v", err)
	}
	
	// 测试设置和获取
	testKey := "test_key"
	testValue := "test_value"
	
	if err := manager.Set(testKey, testValue, time.Minute); err != nil {
		t.Errorf("Failed to set cache: %v", err)
	}
	
	value, err := manager.Get(testKey)
	if err != nil {
		t.Errorf("Failed to get cache: %v", err)
	}
	
	if value != testValue {
		t.Errorf("Expected %s, got %s", testValue, value)
	}
	
	// 测试删除
	if err := manager.Delete(testKey); err != nil {
		t.Errorf("Failed to delete cache: %v", err)
	}
	
	// 检查是否已删除
	exists, err := manager.Exists(testKey)
	if err != nil {
		t.Errorf("Failed to check existence: %v", err)
	}
	
	if exists {
		t.Error("Key should not exist after deletion")
	}
}

func TestCacheManagerJSON(t *testing.T) {
	t.Skip("Skipping cache JSON test - requires actual Redis")
	
	cfg := &config.RedisConfig{
		Host: "localhost",
		Port: 6379,
	}
	
	manager, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	defer manager.Close()
	
	// 测试JSON操作
	type TestStruct struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}
	
	testData := TestStruct{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}
	
	testKey := "test_json_key"
	
	// 设置JSON
	if err := manager.SetJSON(testKey, testData, time.Minute); err != nil {
		t.Errorf("Failed to set JSON cache: %v", err)
	}
	
	// 获取JSON
	var result TestStruct
	if err := manager.GetJSON(testKey, &result); err != nil {
		t.Errorf("Failed to get JSON cache: %v", err)
	}
	
	if result.Name != testData.Name {
		t.Errorf("Expected name %s, got %s", testData.Name, result.Name)
	}
	
	if result.Age != testData.Age {
		t.Errorf("Expected age %d, got %d", testData.Age, result.Age)
	}
}

func TestSessionManager(t *testing.T) {
	t.Skip("Skipping session test - requires actual Redis")
	
	cfg := &config.RedisConfig{
		Host: "localhost",
		Port: 6379,
	}
	
	cacheManager, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	defer cacheManager.Close()
	
	sessionManager := NewSessionManager(cacheManager, "test_session", time.Hour)
	
	// 测试创建会话
	userID := "user123"
	session, err := sessionManager.CreateSession(userID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
	
	if session.UserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, session.UserID)
	}
	
	if session.ID == "" {
		t.Error("Session ID should not be empty")
	}
	
	// 测试获取会话
	retrievedSession, err := sessionManager.GetSession(session.ID)
	if err != nil {
		t.Errorf("Failed to get session: %v", err)
	}
	
	if retrievedSession.ID != session.ID {
		t.Errorf("Expected session ID %s, got %s", session.ID, retrievedSession.ID)
	}
	
	// 测试设置会话数据
	testKey := "test_data"
	testValue := "test_value"
	
	if err := sessionManager.SetSessionData(session.ID, testKey, testValue); err != nil {
		t.Errorf("Failed to set session data: %v", err)
	}
	
	// 测试获取会话数据
	value, err := sessionManager.GetSessionData(session.ID, testKey)
	if err != nil {
		t.Errorf("Failed to get session data: %v", err)
	}
	
	if value != testValue {
		t.Errorf("Expected %s, got %s", testValue, value)
	}
	
	// 测试删除会话
	if err := sessionManager.DeleteSession(session.ID); err != nil {
		t.Errorf("Failed to delete session: %v", err)
	}
	
	// 检查会话是否已删除
	_, err = sessionManager.GetSession(session.ID)
	if err == nil {
		t.Error("Session should not exist after deletion")
	}
}

func TestSessionManagerStats(t *testing.T) {
	t.Skip("Skipping session stats test - requires actual Redis")
	
	cfg := &config.RedisConfig{
		Host: "localhost",
		Port: 6379,
	}
	
	cacheManager, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create cache manager: %v", err)
	}
	defer cacheManager.Close()
	
	sessionManager := NewSessionManager(cacheManager, "test_session_stats", time.Hour)
	
	// 创建几个测试会话
	userIDs := []string{"user1", "user2", "user3"}
	for _, userID := range userIDs {
		_, err := sessionManager.CreateSession(userID)
		if err != nil {
			t.Errorf("Failed to create session for user %s: %v", userID, err)
		}
	}
	
	// 获取统计信息
	stats, err := sessionManager.GetStats()
	if err != nil {
		t.Errorf("Failed to get session stats: %v", err)
	}
	
	if stats.TotalSessions != int64(len(userIDs)) {
		t.Errorf("Expected %d total sessions, got %d", len(userIDs), stats.TotalSessions)
	}
	
	if stats.ActiveSessions != int64(len(userIDs)) {
		t.Errorf("Expected %d active sessions, got %d", len(userIDs), stats.ActiveSessions)
	}
}

// 基准测试
func BenchmarkCacheSet(b *testing.B) {
	b.Skip("Skipping cache benchmark - requires actual Redis")
	
	cfg := &config.RedisConfig{
		Host: "localhost",
		Port: 6379,
	}
	
	manager, err := New(cfg)
	if err != nil {
		b.Fatalf("Failed to create cache manager: %v", err)
	}
	defer manager.Close()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "bench_key"
		value := "bench_value"
		manager.Set(key, value, time.Minute)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	b.Skip("Skipping cache benchmark - requires actual Redis")
	
	cfg := &config.RedisConfig{
		Host: "localhost",
		Port: 6379,
	}
	
	manager, err := New(cfg)
	if err != nil {
		b.Fatalf("Failed to create cache manager: %v", err)
	}
	defer manager.Close()
	
	// 预设数据
	key := "bench_key"
	value := "bench_value"
	manager.Set(key, value, time.Hour)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Get(key)
	}
}