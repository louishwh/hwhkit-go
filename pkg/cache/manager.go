package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hwh/hwhkit-go/pkg/config"
	"github.com/redis/go-redis/v9"
)

// Manager Redis缓存管理器
type Manager struct {
	client *redis.Client
	config *config.RedisConfig
	ctx    context.Context
}

// New 创建新的Redis缓存管理器
func New(cfg *config.RedisConfig) (*Manager, error) {
	ctx := context.Background()
	
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	})
	
	// 测试连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	return &Manager{
		client: rdb,
		config: cfg,
		ctx:    ctx,
	}, nil
}

// GetClient 获取Redis客户端
func (m *Manager) GetClient() *redis.Client {
	return m.client
}

// Close 关闭Redis连接
func (m *Manager) Close() error {
	return m.client.Close()
}

// Health 检查Redis健康状态
func (m *Manager) Health() error {
	return m.client.Ping(m.ctx).Err()
}

// Set 设置缓存值
func (m *Manager) Set(key string, value interface{}, expiration time.Duration) error {
	return m.client.Set(m.ctx, key, value, expiration).Err()
}

// Get 获取缓存值
func (m *Manager) Get(key string) (string, error) {
	return m.client.Get(m.ctx, key).Result()
}

// GetBytes 获取缓存值（字节）
func (m *Manager) GetBytes(key string) ([]byte, error) {
	return m.client.Get(m.ctx, key).Bytes()
}

// SetJSON 设置JSON缓存
func (m *Manager) SetJSON(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return m.client.Set(m.ctx, key, jsonData, expiration).Err()
}

// GetJSON 获取JSON缓存
func (m *Manager) GetJSON(key string, dest interface{}) error {
	jsonData, err := m.client.Get(m.ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, dest)
}

// Delete 删除缓存
func (m *Manager) Delete(keys ...string) error {
	return m.client.Del(m.ctx, keys...).Err()
}

// Exists 检查键是否存在
func (m *Manager) Exists(key string) (bool, error) {
	count, err := m.client.Exists(m.ctx, key).Result()
	return count > 0, err
}

// Expire 设置过期时间
func (m *Manager) Expire(key string, expiration time.Duration) error {
	return m.client.Expire(m.ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (m *Manager) TTL(key string) (time.Duration, error) {
	return m.client.TTL(m.ctx, key).Result()
}

// Increment 递增
func (m *Manager) Increment(key string) (int64, error) {
	return m.client.Incr(m.ctx, key).Result()
}

// IncrementBy 按指定值递增
func (m *Manager) IncrementBy(key string, value int64) (int64, error) {
	return m.client.IncrBy(m.ctx, key, value).Result()
}

// Decrement 递减
func (m *Manager) Decrement(key string) (int64, error) {
	return m.client.Decr(m.ctx, key).Result()
}

// DecrementBy 按指定值递减
func (m *Manager) DecrementBy(key string, value int64) (int64, error) {
	return m.client.DecrBy(m.ctx, key, value).Result()
}

// HSet 设置哈希字段
func (m *Manager) HSet(key string, field string, value interface{}) error {
	return m.client.HSet(m.ctx, key, field, value).Err()
}

// HGet 获取哈希字段
func (m *Manager) HGet(key string, field string) (string, error) {
	return m.client.HGet(m.ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func (m *Manager) HGetAll(key string) (map[string]string, error) {
	return m.client.HGetAll(m.ctx, key).Result()
}

// HDel 删除哈希字段
func (m *Manager) HDel(key string, fields ...string) error {
	return m.client.HDel(m.ctx, key, fields...).Err()
}

// LPush 从左侧推入列表
func (m *Manager) LPush(key string, values ...interface{}) error {
	return m.client.LPush(m.ctx, key, values...).Err()
}

// RPush 从右侧推入列表
func (m *Manager) RPush(key string, values ...interface{}) error {
	return m.client.RPush(m.ctx, key, values...).Err()
}

// LPop 从左侧弹出列表元素
func (m *Manager) LPop(key string) (string, error) {
	return m.client.LPop(m.ctx, key).Result()
}

// RPop 从右侧弹出列表元素
func (m *Manager) RPop(key string) (string, error) {
	return m.client.RPop(m.ctx, key).Result()
}

// LLen 获取列表长度
func (m *Manager) LLen(key string) (int64, error) {
	return m.client.LLen(m.ctx, key).Result()
}

// LRange 获取列表范围
func (m *Manager) LRange(key string, start, stop int64) ([]string, error) {
	return m.client.LRange(m.ctx, key, start, stop).Result()
}

// SAdd 添加集合成员
func (m *Manager) SAdd(key string, members ...interface{}) error {
	return m.client.SAdd(m.ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (m *Manager) SMembers(key string) ([]string, error) {
	return m.client.SMembers(m.ctx, key).Result()
}

// SIsMember 检查是否为集合成员
func (m *Manager) SIsMember(key string, member interface{}) (bool, error) {
	return m.client.SIsMember(m.ctx, key, member).Result()
}

// SRem 移除集合成员
func (m *Manager) SRem(key string, members ...interface{}) error {
	return m.client.SRem(m.ctx, key, members...).Err()
}

// SCard 获取集合成员数量
func (m *Manager) SCard(key string) (int64, error) {
	return m.client.SCard(m.ctx, key).Result()
}

// Keys 获取匹配模式的键
func (m *Manager) Keys(pattern string) ([]string, error) {
	return m.client.Keys(m.ctx, pattern).Result()
}

// FlushDB 清空当前数据库
func (m *Manager) FlushDB() error {
	return m.client.FlushDB(m.ctx).Err()
}

// FlushAll 清空所有数据库
func (m *Manager) FlushAll() error {
	return m.client.FlushAll(m.ctx).Err()
}

// Pipeline 开始管道操作
func (m *Manager) Pipeline() redis.Pipeliner {
	return m.client.Pipeline()
}

// Transaction 开始事务操作
func (m *Manager) Transaction(fn func(redis.Pipeliner) error) error {
	pipe := m.client.TxPipeline()
	if err := fn(pipe); err != nil {
		return err
	}
	_, err := pipe.Exec(m.ctx)
	return err
}

// GetStats 获取Redis统计信息
func (m *Manager) GetStats() map[string]interface{} {
	info := m.client.Info(m.ctx)
	if info.Err() != nil {
		return map[string]interface{}{
			"error": info.Err().Error(),
		}
	}
	
	poolStats := m.client.PoolStats()
	
	return map[string]interface{}{
		"pool_hits":         poolStats.Hits,
		"pool_misses":       poolStats.Misses,
		"pool_timeouts":     poolStats.Timeouts,
		"pool_total_conns":  poolStats.TotalConns,
		"pool_idle_conns":   poolStats.IdleConns,
		"pool_stale_conns":  poolStats.StaleConns,
		"redis_info":        info.Val(),
	}
}