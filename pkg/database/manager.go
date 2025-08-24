package database

import (
	"fmt"
	"time"

	"github.com/hwh/hwhkit-go/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Manager 数据库管理器
type Manager struct {
	db     *gorm.DB
	config *config.DatabaseConfig
}

// New 创建新的数据库管理器
func New(cfg *config.DatabaseConfig) (*Manager, error) {
	manager := &Manager{
		config: cfg,
	}
	
	if err := manager.connect(); err != nil {
		return nil, err
	}
	
	return manager, nil
}

// connect 连接数据库
func (m *Manager) connect() error {
	var dialector gorm.Dialector
	var dsn string
	
	switch m.config.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			m.config.User,
			m.config.Password,
			m.config.Host,
			m.config.Port,
			m.config.Name,
			m.config.Charset,
		)
		dialector = mysql.Open(dsn)
		
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
			m.config.Host,
			m.config.User,
			m.config.Password,
			m.config.Name,
			m.config.Port,
			m.config.SSLMode,
		)
		dialector = postgres.Open(dsn)
		
	default:
		return fmt.Errorf("unsupported database type: %s", m.config.Type)
	}
	
	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	
	// 连接数据库
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// 获取底层的 sql.DB 对象来配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	
	// 配置连接池
	sqlDB.SetMaxOpenConns(m.config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(m.config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(m.config.ConnMaxLifetime) * time.Minute)
	
	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	
	m.db = db
	return nil
}

// GetDB 获取GORM数据库实例
func (m *Manager) GetDB() *gorm.DB {
	return m.db
}

// Close 关闭数据库连接
func (m *Manager) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate 执行数据库迁移
func (m *Manager) Migrate(models ...interface{}) error {
	return m.db.AutoMigrate(models...)
}

// Transaction 执行事务
func (m *Manager) Transaction(fn func(*gorm.DB) error) error {
	return m.db.Transaction(fn)
}

// Health 检查数据库健康状态
func (m *Manager) Health() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetStats 获取数据库连接统计信息
func (m *Manager) GetStats() map[string]interface{} {
	sqlDB, err := m.db.DB()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}
	
	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration":       stats.WaitDuration.String(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}