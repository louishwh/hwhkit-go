package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hwh/hwhkit-go/pkg/config"
)

// Claims JWT声明结构
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	TokenType    string `json:"token_type"`
}

// Manager JWT认证管理器
type Manager struct {
	config        *config.JWTConfig
	signingMethod jwt.SigningMethod
}

// New 创建新的JWT认证管理器
func New(cfg *config.JWTConfig) *Manager {
	return &Manager{
		config:        cfg,
		signingMethod: jwt.SigningMethodHS256,
	}
}

// GenerateToken 生成访问令牌
func (m *Manager) GenerateToken(userID int64, username, email, role string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(m.config.ExpireHours) * time.Hour)
	
	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  []string{m.config.Issuer},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	
	token := jwt.NewWithClaims(m.signingMethod, claims)
	return token.SignedString([]byte(m.config.Secret))
}

// GenerateRefreshToken 生成刷新令牌
func (m *Manager) GenerateRefreshToken(userID int64, username string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(m.config.RefreshHours) * time.Hour)
	
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			Subject:   fmt.Sprintf("refresh:%d", userID),
			Audience:  []string{m.config.Issuer},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	
	token := jwt.NewWithClaims(m.signingMethod, claims)
	return token.SignedString([]byte(m.config.Secret))
}

// GenerateTokenPair 生成令牌对
func (m *Manager) GenerateTokenPair(userID int64, username, email, role string) (*TokenPair, error) {
	accessToken, err := m.GenerateToken(userID, username, email, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	
	refreshToken, err := m.GenerateRefreshToken(userID, username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	
	expiresAt := time.Now().Add(time.Duration(m.config.ExpireHours) * time.Hour).Unix()
	
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken 验证令牌
func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.config.Secret), nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	
	// 验证发行者
	if claims.Issuer != m.config.Issuer {
		return nil, errors.New("invalid token issuer")
	}
	
	return claims, nil
}

// RefreshToken 刷新令牌
func (m *Manager) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	// 验证刷新令牌
	claims, err := m.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	
	// 检查是否为刷新令牌
	if claims.Subject == "" || claims.Subject[:8] != "refresh:" {
		return nil, errors.New("not a refresh token")
	}
	
	// 生成新的令牌对
	return m.GenerateTokenPair(claims.UserID, claims.Username, claims.Email, claims.Role)
}

// ExtractUserID 从令牌中提取用户ID
func (m *Manager) ExtractUserID(tokenString string) (int64, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ExtractUsername 从令牌中提取用户名
func (m *Manager) ExtractUsername(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}

// ExtractRole 从令牌中提取角色
func (m *Manager) ExtractRole(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}

// IsTokenExpired 检查令牌是否过期
func (m *Manager) IsTokenExpired(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.config.Secret), nil
	})
	
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return true, nil
		}
		return false, err
	}
	
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return false, errors.New("invalid token claims")
	}
	
	return claims.ExpiresAt.Before(time.Now()), nil
}

// GetTokenClaims 获取令牌声明信息（包括过期的令牌）
func (m *Manager) GetTokenClaims(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.config.Secret), nil
	}, jwt.WithoutClaimsValidation())
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	
	return claims, nil
}

// ValidateRole 验证用户角色权限
func (m *Manager) ValidateRole(tokenString string, requiredRoles ...string) error {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return err
	}
	
	if len(requiredRoles) == 0 {
		return nil // 无角色要求
	}
	
	userRole := claims.Role
	for _, role := range requiredRoles {
		if userRole == role {
			return nil
		}
	}
	
	return fmt.Errorf("insufficient permissions: required one of %v, got %s", requiredRoles, userRole)
}

// GetTokenExpiration 获取令牌过期时间
func (m *Manager) GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := m.GetTokenClaims(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	
	if claims.ExpiresAt == nil {
		return time.Time{}, errors.New("token has no expiration time")
	}
	
	return claims.ExpiresAt.Time, nil
}