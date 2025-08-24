package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 配置结构
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	JWT      JWTConfig      `json:"jwt"`
	Log      LogConfig      `json:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `json:"port"`
	Mode         string `json:"mode"`         // debug, release, test
	ReadTimeout  int    `json:"read_timeout"` // 秒
	WriteTimeout int    `json:"write_timeout"`
	Host         string `json:"host"`
	EnableCORS   bool   `json:"enable_cors"`
	EnableSwagger bool  `json:"enable_swagger"`
	TemplateDir  string `json:"template_dir"`
	StaticDir    string `json:"static_dir"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string `json:"type"`             // mysql, postgres
	Host            string `json:"host"`
	Port            int    `json:"port"`
	User            string `json:"user"`
	Password        string `json:"password"`
	Name            string `json:"name"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"` // 分钟
	SSLMode         string `json:"ssl_mode"`
	Charset         string `json:"charset"`
	AutoMigrate     bool   `json:"auto_migrate"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns"`
	MaxRetries   int    `json:"max_retries"`
	DialTimeout  int    `json:"dial_timeout"`  // 秒
	ReadTimeout  int    `json:"read_timeout"`  // 秒
	WriteTimeout int    `json:"write_timeout"` // 秒
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret       string `json:"secret"`
	ExpireHours  int    `json:"expire_hours"`
	RefreshHours int    `json:"refresh_hours"`
	Issuer       string `json:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level"`       // debug, info, warn, error
	Format     string `json:"format"`      // json, text
	Output     string `json:"output"`      // console, file, both
	FilePath   string `json:"file_path"`   // 日志文件路径
	MaxSize    int    `json:"max_size"`    // MB
	MaxBackups int    `json:"max_backups"` // 保留的备份数量
	MaxAge     int    `json:"max_age"`     // 保留天数
	Compress   bool   `json:"compress"`    // 是否压缩
}

// ConfigManager 配置管理器
type ConfigManager struct {
	config     *Config
	configURL  string // 远程配置API地址
	httpClient *http.Client
}

// New 创建新的配置管理器
func New() *ConfigManager {
	cm := &ConfigManager{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	
	// 加载配置
	cm.Load()
	return cm
}

// NewWithURL 创建支持远程配置的配置管理器
func NewWithURL(configURL string) *ConfigManager {
	cm := &ConfigManager{
		configURL: configURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	
	// 加载配置
	cm.Load()
	return cm
}

// Load 加载配置
func (cm *ConfigManager) Load() error {
	// 首先尝试从远程加载
	if cm.configURL != "" {
		if err := cm.loadFromRemote(); err == nil {
			return nil
		}
		fmt.Printf("Failed to load config from remote: %v, fallback to local\n", err)
	}
	
	// 从本地加载
	return cm.loadFromLocal()
}

// loadFromLocal 从本地环境变量和.env文件加载配置
func (cm *ConfigManager) loadFromLocal() error {
	// 加载.env文件
	_ = godotenv.Load()
	
	config := &Config{
		Server: ServerConfig{
			Port:          getEnvAsInt("SERVER_PORT", 8080),
			Mode:          getEnv("SERVER_MODE", "debug"),
			ReadTimeout:   getEnvAsInt("SERVER_READ_TIMEOUT", 60),
			WriteTimeout:  getEnvAsInt("SERVER_WRITE_TIMEOUT", 60),
			Host:          getEnv("SERVER_HOST", "0.0.0.0"),
			EnableCORS:    getEnvAsBool("SERVER_ENABLE_CORS", true),
			EnableSwagger: getEnvAsBool("SERVER_ENABLE_SWAGGER", true),
			TemplateDir:   getEnv("SERVER_TEMPLATE_DIR", "templates"),
			StaticDir:     getEnv("SERVER_STATIC_DIR", "static"),
		},
		Database: DatabaseConfig{
			Type:            getEnv("DB_TYPE", "mysql"),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvAsInt("DB_PORT", 3306),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "hwhkit"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 60),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			Charset:         getEnv("DB_CHARSET", "utf8mb4"),
			AutoMigrate:     getEnvAsBool("DB_AUTO_MIGRATE", true),
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnvAsInt("REDIS_PORT", 6379),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
			MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
			DialTimeout:  getEnvAsInt("REDIS_DIAL_TIMEOUT", 5),
			ReadTimeout:  getEnvAsInt("REDIS_READ_TIMEOUT", 3),
			WriteTimeout: getEnvAsInt("REDIS_WRITE_TIMEOUT", 3),
		},
		JWT: JWTConfig{
			Secret:       getEnv("JWT_SECRET", "hwhkit-default-secret-change-in-production"),
			ExpireHours:  getEnvAsInt("JWT_EXPIRE_HOURS", 24),
			RefreshHours: getEnvAsInt("JWT_REFRESH_HOURS", 168), // 7天
			Issuer:       getEnv("JWT_ISSUER", "hwhkit-go"),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			Output:     getEnv("LOG_OUTPUT", "console"),
			FilePath:   getEnv("LOG_FILE_PATH", "logs/app.log"),
			MaxSize:    getEnvAsInt("LOG_MAX_SIZE", 100),
			MaxBackups: getEnvAsInt("LOG_MAX_BACKUPS", 10),
			MaxAge:     getEnvAsInt("LOG_MAX_AGE", 30),
			Compress:   getEnvAsBool("LOG_COMPRESS", true),
		},
	}
	
	cm.config = config
	return nil
}

// loadFromRemote 从远程API加载配置
func (cm *ConfigManager) loadFromRemote() error {
	resp, err := cm.httpClient.Get(cm.configURL)
	if err != nil {
		return fmt.Errorf("failed to fetch config from remote: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote config API returned status: %d", resp.StatusCode)
	}
	
	var config Config
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return fmt.Errorf("failed to decode remote config: %w", err)
	}
	
	cm.config = &config
	return nil
}

// Get 获取配置
func (cm *ConfigManager) Get() *Config {
	return cm.config
}

// Reload 重新加载配置
func (cm *ConfigManager) Reload() error {
	return cm.Load()
}

// GetServer 获取服务器配置
func (cm *ConfigManager) GetServer() *ServerConfig {
	return &cm.config.Server
}

// GetDatabase 获取数据库配置
func (cm *ConfigManager) GetDatabase() *DatabaseConfig {
	return &cm.config.Database
}

// GetRedis 获取Redis配置
func (cm *ConfigManager) GetRedis() *RedisConfig {
	return &cm.config.Redis
}

// GetJWT 获取JWT配置
func (cm *ConfigManager) GetJWT() *JWTConfig {
	return &cm.config.JWT
}

// GetLog 获取日志配置
func (cm *ConfigManager) GetLog() *LogConfig {
	return &cm.config.Log
}

// 辅助函数
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}