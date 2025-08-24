package cache

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// SessionManager 会话管理器
type SessionManager struct {
	cache      *Manager
	prefix     string
	expiration time.Duration
}

// Session 会话数据结构
type Session struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	ExpiresAt time.Time              `json:"expires_at"`
}

// NewSessionManager 创建会话管理器
func NewSessionManager(cache *Manager, prefix string, expiration time.Duration) *SessionManager {
	if prefix == "" {
		prefix = "session"
	}
	if expiration == 0 {
		expiration = 24 * time.Hour // 默认24小时
	}
	
	return &SessionManager{
		cache:      cache,
		prefix:     prefix,
		expiration: expiration,
	}
}

// CreateSession 创建新会话
func (sm *SessionManager) CreateSession(userID string) (*Session, error) {
	sessionID, err := sm.generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}
	
	now := time.Now()
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		Data:      make(map[string]interface{}),
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(sm.expiration),
	}
	
	if err := sm.saveSession(session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}
	
	return session, nil
}

// GetSession 获取会话
func (sm *SessionManager) GetSession(sessionID string) (*Session, error) {
	key := sm.getSessionKey(sessionID)
	
	var session Session
	if err := sm.cache.GetJSON(key, &session); err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	
	// 检查会话是否过期
	if time.Now().After(session.ExpiresAt) {
		sm.DeleteSession(sessionID) // 删除过期会话
		return nil, fmt.Errorf("session expired")
	}
	
	return &session, nil
}

// UpdateSession 更新会话
func (sm *SessionManager) UpdateSession(session *Session) error {
	session.UpdatedAt = time.Now()
	return sm.saveSession(session)
}

// DeleteSession 删除会话
func (sm *SessionManager) DeleteSession(sessionID string) error {
	key := sm.getSessionKey(sessionID)
	return sm.cache.Delete(key)
}

// RefreshSession 刷新会话过期时间
func (sm *SessionManager) RefreshSession(sessionID string) error {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}
	
	session.ExpiresAt = time.Now().Add(sm.expiration)
	return sm.UpdateSession(session)
}

// SetSessionData 设置会话数据
func (sm *SessionManager) SetSessionData(sessionID string, key string, value interface{}) error {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}
	
	session.Data[key] = value
	return sm.UpdateSession(session)
}

// GetSessionData 获取会话数据
func (sm *SessionManager) GetSessionData(sessionID string, key string) (interface{}, error) {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	
	value, exists := session.Data[key]
	if !exists {
		return nil, fmt.Errorf("session data key not found: %s", key)
	}
	
	return value, nil
}

// RemoveSessionData 移除会话数据
func (sm *SessionManager) RemoveSessionData(sessionID string, key string) error {
	session, err := sm.GetSession(sessionID)
	if err != nil {
		return err
	}
	
	delete(session.Data, key)
	return sm.UpdateSession(session)
}

// GetUserSessions 获取用户的所有会话
func (sm *SessionManager) GetUserSessions(userID string) ([]*Session, error) {
	pattern := sm.prefix + ":*"
	keys, err := sm.cache.Keys(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get session keys: %w", err)
	}
	
	var userSessions []*Session
	for _, key := range keys {
		var session Session
		if err := sm.cache.GetJSON(key, &session); err != nil {
			continue // 跳过无法解析的会话
		}
		
		if session.UserID == userID && time.Now().Before(session.ExpiresAt) {
			userSessions = append(userSessions, &session)
		}
	}
	
	return userSessions, nil
}

// DeleteUserSessions 删除用户的所有会话
func (sm *SessionManager) DeleteUserSessions(userID string) error {
	sessions, err := sm.GetUserSessions(userID)
	if err != nil {
		return err
	}
	
	for _, session := range sessions {
		if err := sm.DeleteSession(session.ID); err != nil {
			return fmt.Errorf("failed to delete session %s: %w", session.ID, err)
		}
	}
	
	return nil
}

// CleanExpiredSessions 清理过期会话
func (sm *SessionManager) CleanExpiredSessions() error {
	pattern := sm.prefix + ":*"
	keys, err := sm.cache.Keys(pattern)
	if err != nil {
		return fmt.Errorf("failed to get session keys: %w", err)
	}
	
	now := time.Now()
	for _, key := range keys {
		var session Session
		if err := sm.cache.GetJSON(key, &session); err != nil {
			continue
		}
		
		if now.After(session.ExpiresAt) {
			sessionID := sm.extractSessionID(key)
			sm.DeleteSession(sessionID)
		}
	}
	
	return nil
}

// GetSessionCount 获取活跃会话数量
func (sm *SessionManager) GetSessionCount() (int64, error) {
	pattern := sm.prefix + ":*"
	keys, err := sm.cache.Keys(pattern)
	if err != nil {
		return 0, fmt.Errorf("failed to get session keys: %w", err)
	}
	
	var activeCount int64
	now := time.Now()
	
	for _, key := range keys {
		var session Session
		if err := sm.cache.GetJSON(key, &session); err != nil {
			continue
		}
		
		if now.Before(session.ExpiresAt) {
			activeCount++
		}
	}
	
	return activeCount, nil
}

// IsValidSession 检查会话是否有效
func (sm *SessionManager) IsValidSession(sessionID string) bool {
	_, err := sm.GetSession(sessionID)
	return err == nil
}

// 私有方法

// generateSessionID 生成会话ID
func (sm *SessionManager) generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// getSessionKey 获取会话的Redis键
func (sm *SessionManager) getSessionKey(sessionID string) string {
	return fmt.Sprintf("%s:%s", sm.prefix, sessionID)
}

// extractSessionID 从Redis键中提取会话ID
func (sm *SessionManager) extractSessionID(key string) string {
	prefix := sm.prefix + ":"
	if len(key) > len(prefix) {
		return key[len(prefix):]
	}
	return ""
}

// saveSession 保存会话到Redis
func (sm *SessionManager) saveSession(session *Session) error {
	key := sm.getSessionKey(session.ID)
	
	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}
	
	return sm.cache.Set(key, sessionData, sm.expiration)
}

// SessionStats 会话统计信息
type SessionStats struct {
	TotalSessions   int64              `json:"total_sessions"`
	ActiveSessions  int64              `json:"active_sessions"`
	ExpiredSessions int64              `json:"expired_sessions"`
	UserSessions    map[string]int64   `json:"user_sessions"`
}

// GetStats 获取会话统计信息
func (sm *SessionManager) GetStats() (*SessionStats, error) {
	pattern := sm.prefix + ":*"
	keys, err := sm.cache.Keys(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get session keys: %w", err)
	}
	
	stats := &SessionStats{
		TotalSessions: int64(len(keys)),
		UserSessions:  make(map[string]int64),
	}
	
	now := time.Now()
	for _, key := range keys {
		var session Session
		if err := sm.cache.GetJSON(key, &session); err != nil {
			continue
		}
		
		if now.Before(session.ExpiresAt) {
			stats.ActiveSessions++
			if session.UserID != "" {
				stats.UserSessions[session.UserID]++
			}
		} else {
			stats.ExpiredSessions++
		}
	}
	
	return stats, nil
}