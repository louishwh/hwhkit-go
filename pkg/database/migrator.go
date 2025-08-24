package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Migrator 数据库迁移器
type Migrator struct {
	db      *gorm.DB
	models  []interface{}
	seeds   []SeedFunc
	version string
}

// SeedFunc 种子数据函数类型
type SeedFunc func(*gorm.DB) error

// NewMigrator 创建新的迁移器
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:      db,
		models:  make([]interface{}, 0),
		seeds:   make([]SeedFunc, 0),
		version: "1.0.0",
	}
}

// AddModel 添加需要迁移的模型
func (m *Migrator) AddModel(model interface{}) *Migrator {
	m.models = append(m.models, model)
	return m
}

// AddModels 批量添加需要迁移的模型
func (m *Migrator) AddModels(models ...interface{}) *Migrator {
	m.models = append(m.models, models...)
	return m
}

// AddSeed 添加种子数据
func (m *Migrator) AddSeed(seed SeedFunc) *Migrator {
	m.seeds = append(m.seeds, seed)
	return m
}

// SetVersion 设置迁移版本
func (m *Migrator) SetVersion(version string) *Migrator {
	m.version = version
	return m
}

// Migrate 执行迁移
func (m *Migrator) Migrate() error {
	log.Printf("Starting database migration (version: %s)...", m.version)
	
	// 自动迁移模型
	if len(m.models) > 0 {
		log.Printf("Migrating %d models...", len(m.models))
		if err := m.db.AutoMigrate(m.models...); err != nil {
			return fmt.Errorf("failed to migrate models: %w", err)
		}
		log.Println("Models migrated successfully")
	}
	
	// 执行种子数据
	if len(m.seeds) > 0 {
		log.Printf("Running %d seed functions...", len(m.seeds))
		for i, seed := range m.seeds {
			log.Printf("Running seed %d/%d...", i+1, len(m.seeds))
			if err := seed(m.db); err != nil {
				return fmt.Errorf("failed to run seed %d: %w", i+1, err)
			}
		}
		log.Println("Seeds executed successfully")
	}
	
	log.Println("Database migration completed successfully")
	return nil
}

// Rollback 回滚迁移（删除表）
func (m *Migrator) Rollback() error {
	log.Printf("Starting database rollback (version: %s)...", m.version)
	
	// 删除表（按相反顺序）
	for i := len(m.models) - 1; i >= 0; i-- {
		model := m.models[i]
		if err := m.db.Migrator().DropTable(model); err != nil {
			log.Printf("Warning: failed to drop table for model %T: %v", model, err)
		} else {
			log.Printf("Dropped table for model %T", model)
		}
	}
	
	log.Println("Database rollback completed")
	return nil
}

// Status 检查迁移状态
func (m *Migrator) Status() error {
	log.Printf("Checking migration status (version: %s)...", m.version)
	
	for _, model := range m.models {
		hasTable := m.db.Migrator().HasTable(model)
		log.Printf("Model %T: table exists = %v", model, hasTable)
		
		if hasTable {
			// 检查列信息
			tableName := m.db.Statement.Table
			columns, err := m.db.Migrator().ColumnTypes(tableName)
			if err != nil {
				log.Printf("  Warning: failed to get column info: %v", err)
			} else {
				log.Printf("  Columns: %d", len(columns))
			}
		}
	}
	
	return nil
}

// CreateIndexes 创建索引
func (m *Migrator) CreateIndexes(indexes []IndexDefinition) error {
	log.Printf("Creating %d indexes...", len(indexes))
	
	for _, index := range indexes {
		if err := m.createIndex(index); err != nil {
			return fmt.Errorf("failed to create index %s: %w", index.Name, err)
		}
		log.Printf("Created index: %s", index.Name)
	}
	
	return nil
}

// IndexDefinition 索引定义
type IndexDefinition struct {
	Name    string
	Table   string
	Columns []string
	Unique  bool
}

// createIndex 创建单个索引
func (m *Migrator) createIndex(index IndexDefinition) error {
	sql := fmt.Sprintf("CREATE")
	if index.Unique {
		sql += " UNIQUE"
	}
	sql += fmt.Sprintf(" INDEX %s ON %s (%s)",
		index.Name,
		index.Table,
		joinColumns(index.Columns),
	)
	
	return m.db.Exec(sql).Error
}

// joinColumns 连接列名
func joinColumns(columns []string) string {
	if len(columns) == 0 {
		return ""
	}
	
	result := columns[0]
	for i := 1; i < len(columns); i++ {
		result += ", " + columns[i]
	}
	return result
}

// GetDB 获取数据库实例
func (m *Migrator) GetDB() *gorm.DB {
	return m.db
}

// 预定义的一些常用种子数据函数

// CreateAdminUser 创建管理员用户种子数据
func CreateAdminUser(username, email, password string) SeedFunc {
	return func(db *gorm.DB) error {
		// 这里只是示例，实际使用时需要根据具体的用户模型来实现
		log.Printf("Creating admin user: %s (%s)", username, email)
		// 实现创建管理员用户的逻辑
		return nil
	}
}

// CreateDefaultRoles 创建默认角色种子数据
func CreateDefaultRoles(roles []string) SeedFunc {
	return func(db *gorm.DB) error {
		log.Printf("Creating default roles: %v", roles)
		// 实现创建默认角色的逻辑
		return nil
	}
}