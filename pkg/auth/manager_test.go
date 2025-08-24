package auth

import (
	"testing"
	"time"

	"github.com/hwh/hwhkit-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestConfig() *config.JWTConfig {
	return &config.JWTConfig{
		Secret:       "test-secret-key",
		ExpireHours:  1,
		RefreshHours: 24,
		Issuer:       "test-app",
	}
}

func TestNew(t *testing.T) {
	cfg := getTestConfig()
	manager := New(cfg)
	
	assert.NotNil(t, manager)
	assert.Equal(t, cfg, manager.config)
}

func TestGenerateToken(t *testing.T) {
	manager := New(getTestConfig())
	
	token, err := manager.GenerateToken(123, "testuser", "test@example.com", "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// 验证生成的令牌
	claims, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, int64(123), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, "test-app", claims.Issuer)
}

func TestGenerateRefreshToken(t *testing.T) {
	manager := New(getTestConfig())
	
	token, err := manager.GenerateRefreshToken(123, "testuser")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// 验证刷新令牌的特殊Subject
	claims, err := manager.GetTokenClaims(token)
	require.NoError(t, err)
	assert.Equal(t, "refresh:123", claims.Subject)
}

func TestGenerateTokenPair(t *testing.T) {
	manager := New(getTestConfig())
	
	tokenPair, err := manager.GenerateTokenPair(123, "testuser", "test@example.com", "user")
	require.NoError(t, err)
	
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.Equal(t, "Bearer", tokenPair.TokenType)
	assert.Greater(t, tokenPair.ExpiresAt, time.Now().Unix())
	
	// 验证访问令牌
	claims, err := manager.ValidateToken(tokenPair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, int64(123), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	
	// 验证刷新令牌
	refreshClaims, err := manager.GetTokenClaims(tokenPair.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, "refresh:123", refreshClaims.Subject)
}

func TestValidateToken(t *testing.T) {
	manager := New(getTestConfig())
	
	// 生成有效令牌
	token, err := manager.GenerateToken(123, "testuser", "test@example.com", "admin")
	require.NoError(t, err)
	
	// 验证有效令牌
	claims, err := manager.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, int64(123), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	
	// 测试无效令牌
	_, err = manager.ValidateToken("invalid.token.here")
	assert.Error(t, err)
	
	// 测试空令牌
	_, err = manager.ValidateToken("")
	assert.Error(t, err)
}

func TestValidateTokenWithWrongSecret(t *testing.T) {
	manager1 := New(getTestConfig())
	
	// 使用不同密钥的管理器
	cfg2 := getTestConfig()
	cfg2.Secret = "different-secret"
	manager2 := New(cfg2)
	
	// 用manager1生成令牌
	token, err := manager1.GenerateToken(123, "testuser", "test@example.com", "admin")
	require.NoError(t, err)
	
	// 用manager2验证令牌应该失败
	_, err = manager2.ValidateToken(token)
	assert.Error(t, err)
}

func TestRefreshToken(t *testing.T) {
	manager := New(getTestConfig())
	
	// 生成原始令牌对
	originalPair, err := manager.GenerateTokenPair(123, "testuser", "test@example.com", "user")
	require.NoError(t, err)
	
	// 使用刷新令牌生成新的令牌对
	newPair, err := manager.RefreshToken(originalPair.RefreshToken)
	require.NoError(t, err)
	
	assert.NotEmpty(t, newPair.AccessToken)
	assert.NotEmpty(t, newPair.RefreshToken)
	assert.NotEqual(t, originalPair.AccessToken, newPair.AccessToken)
	
	// 验证新的访问令牌
	claims, err := manager.ValidateToken(newPair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, int64(123), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
}

func TestRefreshTokenWithAccessToken(t *testing.T) {
	manager := New(getTestConfig())
	
	// 生成访问令牌
	accessToken, err := manager.GenerateToken(123, "testuser", "test@example.com", "user")
	require.NoError(t, err)
	
	// 尝试用访问令牌刷新应该失败
	_, err = manager.RefreshToken(accessToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a refresh token")
}

func TestExtractUserID(t *testing.T) {
	manager := New(getTestConfig())
	
	token, err := manager.GenerateToken(456, "testuser", "test@example.com", "user")
	require.NoError(t, err)
	
	userID, err := manager.ExtractUserID(token)
	require.NoError(t, err)
	assert.Equal(t, int64(456), userID)
}

func TestExtractUsername(t *testing.T) {
	manager := New(getTestConfig())
	
	token, err := manager.GenerateToken(123, "johndoe", "john@example.com", "user")
	require.NoError(t, err)
	
	username, err := manager.ExtractUsername(token)
	require.NoError(t, err)
	assert.Equal(t, "johndoe", username)
}

func TestExtractRole(t *testing.T) {
	manager := New(getTestConfig())
	
	token, err := manager.GenerateToken(123, "testuser", "test@example.com", "moderator")
	require.NoError(t, err)
	
	role, err := manager.ExtractRole(token)
	require.NoError(t, err)
	assert.Equal(t, "moderator", role)
}

func TestIsTokenExpired(t *testing.T) {
	// 创建一个过期时间很短的配置
	cfg := getTestConfig()
	cfg.ExpireHours = 0 // 立即过期
	manager := New(cfg)
	
	token, err := manager.GenerateToken(123, "testuser", "test@example.com", "user")
	require.NoError(t, err)
	
	// 稍等一下确保令牌过期
	time.Sleep(1 * time.Second)
	
	expired, err := manager.IsTokenExpired(token)
	require.NoError(t, err)
	assert.True(t, expired)
}

func TestValidateRole(t *testing.T) {
	manager := New(getTestConfig())
	
	// 生成管理员令牌
	adminToken, err := manager.GenerateToken(123, "admin", "admin@example.com", "admin")
	require.NoError(t, err)
	
	// 生成用户令牌
	userToken, err := manager.GenerateToken(456, "user", "user@example.com", "user")
	require.NoError(t, err)
	
	// 测试管理员权限验证
	err = manager.ValidateRole(adminToken, "admin")
	assert.NoError(t, err)
	
	err = manager.ValidateRole(adminToken, "user", "admin")
	assert.NoError(t, err)
	
	// 测试用户权限验证失败
	err = manager.ValidateRole(userToken, "admin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient permissions")
	
	// 测试用户权限验证成功
	err = manager.ValidateRole(userToken, "user")
	assert.NoError(t, err)
	
	// 测试无权限要求
	err = manager.ValidateRole(userToken)
	assert.NoError(t, err)
}

func TestGetTokenExpiration(t *testing.T) {
	manager := New(getTestConfig())
	
	beforeGeneration := time.Now()
	token, err := manager.GenerateToken(123, "testuser", "test@example.com", "user")
	require.NoError(t, err)
	afterGeneration := time.Now()
	
	expiration, err := manager.GetTokenExpiration(token)
	require.NoError(t, err)
	
	// 过期时间应该在生成前+1小时到生成后+1小时之间
	expectedMin := beforeGeneration.Add(time.Hour)
	expectedMax := afterGeneration.Add(time.Hour)
	
	assert.True(t, expiration.After(expectedMin) || expiration.Equal(expectedMin))
	assert.True(t, expiration.Before(expectedMax) || expiration.Equal(expectedMax))
}

func TestGetTokenClaims(t *testing.T) {
	// 创建过期令牌
	cfg := getTestConfig()
	cfg.ExpireHours = 0
	manager := New(cfg)
	
	token, err := manager.GenerateToken(123, "testuser", "test@example.com", "user")
	require.NoError(t, err)
	
	time.Sleep(1 * time.Second)
	
	// ValidateToken应该失败
	_, err = manager.ValidateToken(token)
	assert.Error(t, err)
	
	// 但GetTokenClaims应该成功
	claims, err := manager.GetTokenClaims(token)
	require.NoError(t, err)
	assert.Equal(t, int64(123), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
}