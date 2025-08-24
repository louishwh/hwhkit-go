package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hwh/hwhkit-go/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

// JWTManager JWT管理器
type JWTManager struct {
	config *config.JWTConfig
}

// Claims JWT声明
type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// User 用户信息
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"-"` // 不返回密码
	Roles    []string `json:"roles"`
	IsActive bool     `json:"is_active"`
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	return &JWTManager{
		config: cfg,
	}
}

// GenerateToken 生成访问令牌
func (jm *JWTManager) GenerateToken(user *User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(jm.config.ExpireHours) * time.Hour)
	
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jm.config.Issuer,
			Subject:   user.ID,
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.config.Secret))
}

// GenerateRefreshToken 生成刷新令牌
func (jm *JWTManager) GenerateRefreshToken(user *User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(jm.config.RefreshHours) * time.Hour)
	
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jm.config.Issuer,
			Subject:   user.ID,
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.config.Secret))
}

// GenerateTokenPair 生成令牌对
func (jm *JWTManager) GenerateTokenPair(user *User) (*TokenPair, error) {
	accessToken, err := jm.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	
	refreshToken, err := jm.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(jm.config.ExpireHours * 3600), // 转换为秒
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken 验证令牌
func (jm *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jm.config.Secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

// RefreshToken 刷新令牌
func (jm *JWTManager) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := jm.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	
	// 检查是否为有效的刷新令牌（通常刷新令牌的过期时间更长）
	now := time.Now()
	if claims.ExpiresAt.Time.Before(now) {
		return nil, errors.New("refresh token expired")
	}
	
	// 创建用户对象用于生成新令牌
	user := &User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Roles:    claims.Roles,
	}
	
	return jm.GenerateTokenPair(user)
}

// ExtractTokenFromHeader 从请求头中提取令牌
func (jm *JWTManager) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}
	
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}
	
	return parts[1], nil
}

// GetUserFromToken 从令牌中获取用户信息
func (jm *JWTManager) GetUserFromToken(tokenString string) (*User, error) {
	claims, err := jm.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	return &User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Roles:    claims.Roles,
		IsActive: true, // 假设令牌有效则用户活跃
	}, nil
}

// HasRole 检查用户是否具有指定角色
func (jm *JWTManager) HasRole(claims *Claims, role string) bool {
	for _, userRole := range claims.Roles {
		if userRole == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否具有任意指定角色
func (jm *JWTManager) HasAnyRole(claims *Claims, roles []string) bool {
	for _, role := range roles {
		if jm.HasRole(claims, role) {
			return true
		}
	}
	return false
}

// HasAllRoles 检查用户是否具有所有指定角色
func (jm *JWTManager) HasAllRoles(claims *Claims, roles []string) bool {
	for _, role := range roles {
		if !jm.HasRole(claims, role) {
			return false
		}
	}
	return true
}

// IsTokenExpired 检查令牌是否过期
func (jm *JWTManager) IsTokenExpired(claims *Claims) bool {
	return claims.ExpiresAt.Time.Before(time.Now())
}

// GetTokenRemainingTime 获取令牌剩余时间
func (jm *JWTManager) GetTokenRemainingTime(claims *Claims) time.Duration {
	return time.Until(claims.ExpiresAt.Time)
}

// PasswordManager 密码管理器
type PasswordManager struct {
	cost int
}

// NewPasswordManager 创建密码管理器
func NewPasswordManager(cost int) *PasswordManager {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &PasswordManager{
		cost: cost,
	}
}

// HashPassword 哈希密码
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost)
	return string(bytes), err
}

// CheckPassword 检查密码
func (pm *PasswordManager) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength 验证密码强度
func (pm *PasswordManager) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false
	
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		default:
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}
	
	return nil
}

// AuthService 认证服务
type AuthService struct {
	jwtManager      *JWTManager
	passwordManager *PasswordManager
	config          *config.JWTConfig
}

// NewAuthService 创建认证服务
func NewAuthService(cfg *config.JWTConfig) *AuthService {
	return &AuthService{
		jwtManager:      NewJWTManager(cfg),
		passwordManager: NewPasswordManager(bcrypt.DefaultCost),
		config:          cfg,
	}
}

// GetJWTManager 获取JWT管理器
func (as *AuthService) GetJWTManager() *JWTManager {
	return as.jwtManager
}

// GetPasswordManager 获取密码管理器
func (as *AuthService) GetPasswordManager() *PasswordManager {
	return as.passwordManager
}

// Login 用户登录
func (as *AuthService) Login(username, password string, userProvider func(string) (*User, error)) (*TokenPair, error) {
	// 获取用户信息
	user, err := userProvider(username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}
	
	// 验证密码
	if !as.passwordManager.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid password")
	}
	
	// 生成令牌对
	return as.jwtManager.GenerateTokenPair(user)
}

// Register 用户注册
func (as *AuthService) Register(username, email, password string, roles []string, userCreator func(*User) error) (*TokenPair, error) {
	// 验证密码强度
	if err := as.passwordManager.ValidatePasswordStrength(password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}
	
	// 哈希密码
	hashedPassword, err := as.passwordManager.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	
	// 创建用户
	user := &User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Roles:    roles,
		IsActive: true,
	}
	
	// 调用用户创建器
	if err := userCreator(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	// 生成令牌对
	return as.jwtManager.GenerateTokenPair(user)
}

// ChangePassword 修改密码
func (as *AuthService) ChangePassword(userID, oldPassword, newPassword string, userProvider func(string) (*User, error), userUpdater func(*User) error) error {
	// 获取用户信息
	user, err := userProvider(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	
	// 验证旧密码
	if !as.passwordManager.CheckPassword(oldPassword, user.Password) {
		return errors.New("invalid old password")
	}
	
	// 验证新密码强度
	if err := as.passwordManager.ValidatePasswordStrength(newPassword); err != nil {
		return fmt.Errorf("new password validation failed: %w", err)
	}
	
	// 哈希新密码
	hashedPassword, err := as.passwordManager.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}
	
	// 更新用户密码
	user.Password = hashedPassword
	return userUpdater(user)
}