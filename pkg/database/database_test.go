package database

import (
	"testing"

	"github.com/hwh/hwhkit-go/pkg/config"
	"gorm.io/gorm"
)

// TestUser 测试用户模型
type TestUser struct {
	BaseModel
	Name  string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex"`
}

func TestDatabaseManager(t *testing.T) {
	// 使用SQLite内存数据库进行测试
	cfg := &config.DatabaseConfig{
		Type: "sqlite",
		Name: ":memory:",
	}
	
	// 由于我们的代码目前不支持SQLite，这里用MySQL配置做单元测试
	// 在实际测试中，你可能需要使用测试数据库
	cfg = &config.DatabaseConfig{
		Type:            "mysql",
		Host:            "localhost",
		Port:            3306,
		User:            "root",
		Password:        "password",
		Name:            "test_db",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30,
		Charset:         "utf8mb4",
	}
	
	// 注意：这个测试需要实际的数据库连接
	// 在CI/CD环境中，你可能需要跳过这个测试或使用模拟数据库
	t.Skip("Skipping database integration test - requires actual database")
	
	manager, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}
	defer manager.Close()
	
	// 测试健康检查
	if err := manager.Health(); err != nil {
		t.Errorf("Database health check failed: %v", err)
	}
	
	// 测试获取统计信息
	stats := manager.GetStats()
	if stats == nil {
		t.Error("Failed to get database stats")
	}
}

func TestBaseRepository(t *testing.T) {
	// 这个测试也需要实际的数据库连接
	t.Skip("Skipping repository test - requires actual database")
	
	// 假设我们有一个可用的数据库连接
	var db *gorm.DB // 这里应该是实际的数据库连接
	
	repo := NewBaseRepository[TestUser](db)
	
	// 测试创建
	user := &TestUser{
		Name:  "Test User",
		Email: "test@example.com",
	}
	
	if err := repo.Create(user); err != nil {
		t.Errorf("Failed to create user: %v", err)
	}
	
	// 测试获取
	if user.ID > 0 {
		retrievedUser, err := repo.GetByID(user.ID)
		if err != nil {
			t.Errorf("Failed to get user by ID: %v", err)
		}
		
		if retrievedUser.Name != user.Name {
			t.Errorf("Expected name %s, got %s", user.Name, retrievedUser.Name)
		}
	}
}

func TestMigrator(t *testing.T) {
	// 这个测试也需要实际的数据库连接
	t.Skip("Skipping migrator test - requires actual database")
	
	var db *gorm.DB // 这里应该是实际的数据库连接
	
	migrator := NewMigrator(db)
	migrator.AddModel(&TestUser{})
	migrator.SetVersion("1.0.0")
	
	// 测试状态检查
	if err := migrator.Status(); err != nil {
		t.Errorf("Failed to check migration status: %v", err)
	}
	
	// 测试迁移
	if err := migrator.Migrate(); err != nil {
		t.Errorf("Failed to run migration: %v", err)
	}
}

func TestPaginationResult(t *testing.T) {
	// 测试分页结果结构
	result := PaginationResult[TestUser]{
		Data:       []*TestUser{},
		Total:      100,
		Page:       1,
		PageSize:   10,
		TotalPages: 10,
	}
	
	if result.Total != 100 {
		t.Errorf("Expected total 100, got %d", result.Total)
	}
	
	if result.TotalPages != 10 {
		t.Errorf("Expected total pages 10, got %d", result.TotalPages)
	}
}

func TestIndexDefinition(t *testing.T) {
	// 测试索引定义
	index := IndexDefinition{
		Name:    "idx_user_email",
		Table:   "users",
		Columns: []string{"email"},
		Unique:  true,
	}
	
	if index.Name != "idx_user_email" {
		t.Errorf("Expected index name idx_user_email, got %s", index.Name)
	}
	
	if !index.Unique {
		t.Error("Expected unique index")
	}
}

// 测试辅助函数
func TestJoinColumns(t *testing.T) {
	// 测试空列
	result := joinColumns([]string{})
	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}
	
	// 测试单列
	result = joinColumns([]string{"name"})
	if result != "name" {
		t.Errorf("Expected 'name', got %s", result)
	}
	
	// 测试多列
	result = joinColumns([]string{"name", "email", "created_at"})
	expected := "name, email, created_at"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}