package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/hwh/hwhkit-go/pkg/config"
)

func TestJWTManager(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:       "test-secret-key",
		ExpireHours:  24,
		RefreshHours: 168,
		Issuer:       "test-issuer",
	}
	
	jwtManager := NewJWTManager(cfg)
	
	user := &User{
		ID:       "123",
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"user", "admin"},
		IsActive: true,
	}
	
	// 测试生成访问令牌
	token, err := jwtManager.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	if token == "" {
		t.Error("Token should not be empty")
	}
	
	// 测试验证令牌
	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	
	if claims.UserID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, claims.UserID)
	}
	
	if claims.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, claims.Username)
	}
	
	if claims.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, claims.Email)
	}
	
	if len(claims.Roles) != len(user.Roles) {
		t.Errorf("Expected %d roles, got %d", len(user.Roles), len(claims.Roles))
	}
}

func TestJWTManagerTokenPair(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:       "test-secret-key",
		ExpireHours:  1,
		RefreshHours: 24,
		Issuer:       "test-issuer",
	}
	
	jwtManager := NewJWTManager(cfg)
	
	user := &User{
		ID:       "123",
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"user"},
		IsActive: true,
	}
	
	// 测试生成令牌对
	tokenPair, err := jwtManager.GenerateTokenPair(user)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}
	
	if tokenPair.AccessToken == "" {
		t.Error("Access token should not be empty")
	}
	
	if tokenPair.RefreshToken == "" {
		t.Error("Refresh token should not be empty")
	}
	
	if tokenPair.TokenType != "Bearer" {
		t.Errorf("Expected token type Bearer, got %s", tokenPair.TokenType)
	}
	
	if tokenPair.ExpiresIn != 3600 {
		t.Errorf("Expected expires in 3600, got %d", tokenPair.ExpiresIn)
	}
	
	// 测试刷新令牌
	newTokenPair, err := jwtManager.RefreshToken(tokenPair.RefreshToken)
	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}
	
	if newTokenPair.AccessToken == tokenPair.AccessToken {
		t.Error("New access token should be different from old one")
	}
}

func TestJWTManagerExtractToken(t *testing.T) {
	jwtManager := NewJWTManager(&config.JWTConfig{})
	
	// 测试正确的格式
	authHeader := "Bearer abc123def456"
	token, err := jwtManager.ExtractTokenFromHeader(authHeader)
	if err != nil {
		t.Errorf("Failed to extract token: %v", err)
	}
	
	if token != "abc123def456" {
		t.Errorf("Expected token abc123def456, got %s", token)
	}
	
	// 测试错误的格式
	invalidHeaders := []string{
		"",
		"abc123",
		"Basic abc123",
		"Bearer",
		"bearer abc123",
	}
	
	for _, header := range invalidHeaders {
		_, err := jwtManager.ExtractTokenFromHeader(header)
		if err == nil {
			t.Errorf("Should fail for invalid header: %s", header)
		}
	}
}

func TestJWTManagerRoles(t *testing.T) {
	jwtManager := NewJWTManager(&config.JWTConfig{})
	
	claims := &Claims{
		UserID:   "123",
		Username: "testuser",
		Roles:    []string{"user", "admin", "moderator"},
	}
	
	// 测试 HasRole
	if !jwtManager.HasRole(claims, "admin") {
		t.Error("Should have admin role")
	}
	
	if jwtManager.HasRole(claims, "guest") {
		t.Error("Should not have guest role")
	}
	
	// 测试 HasAnyRole
	if !jwtManager.HasAnyRole(claims, []string{"guest", "admin"}) {
		t.Error("Should have at least one of the roles")
	}
	
	if jwtManager.HasAnyRole(claims, []string{"guest", "visitor"}) {
		t.Error("Should not have any of the roles")
	}
	
	// 测试 HasAllRoles
	if !jwtManager.HasAllRoles(claims, []string{"user", "admin"}) {
		t.Error("Should have all specified roles")
	}
	
	if jwtManager.HasAllRoles(claims, []string{"user", "guest"}) {
		t.Error("Should not have all specified roles")
	}
}

func TestPasswordManager(t *testing.T) {
	pm := NewPasswordManager(10) // 使用较低的cost以加快测试
	
	password := "TestPassword123!"
	
	// 测试密码哈希
	hash, err := pm.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	if hash == "" {
		t.Error("Hash should not be empty")
	}
	
	if hash == password {
		t.Error("Hash should be different from original password")
	}
	
	// 测试密码验证
	if !pm.CheckPassword(password, hash) {
		t.Error("Password should be valid")
	}
	
	if pm.CheckPassword("WrongPassword", hash) {
		t.Error("Wrong password should not be valid")
	}
}

func TestPasswordStrength(t *testing.T) {
	pm := NewPasswordManager(10)
	
	// 测试有效密码
	validPasswords := []string{
		"Password123!",
		"MySecure@Pass1",
		"Strong#Password456",
	}
	
	for _, password := range validPasswords {
		if err := pm.ValidatePasswordStrength(password); err != nil {
			t.Errorf("Valid password %s should pass validation: %v", password, err)
		}
	}
	
	// 测试无效密码
	invalidPasswords := []string{
		"short",        // 太短
		"password",     // 没有大写字母
		"PASSWORD",     // 没有小写字母
		"Password",     // 没有数字
		"Password123",  // 没有特殊字符
	}
	
	for _, password := range invalidPasswords {
		if err := pm.ValidatePasswordStrength(password); err == nil {
			t.Errorf("Invalid password %s should fail validation", password)
		}
	}
}

func TestAuthService(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:       "test-secret",
		ExpireHours:  1,
		RefreshHours: 24,
		Issuer:       "test",
	}
	
	authService := NewAuthService(cfg)
	
	// 模拟用户数据存储
	users := make(map[string]*User)
	
	userProvider := func(username string) (*User, error) {
		if user, exists := users[username]; exists {
			return user, nil
		}
		return nil, &TestError{Message: "user not found"}
	}
	
	userCreator := func(user *User) error {
		user.ID = "generated-id-123"
		users[user.Username] = user
		return nil
	}
	
	userUpdater := func(user *User) error {
		users[user.Username] = user
		return nil
	}
	
	// 测试注册
	tokenPair, err := authService.Register("testuser", "test@example.com", "Password123!", []string{"user"}, userCreator)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	
	if tokenPair == nil {
		t.Fatal("Token pair should not be nil")
	}
	
	// 测试登录
	tokenPair, err = authService.Login("testuser", "Password123!", userProvider)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	
	if tokenPair == nil {
		t.Fatal("Token pair should not be nil")
	}
	
	// 测试错误密码登录
	_, err = authService.Login("testuser", "WrongPassword", userProvider)
	if err == nil {
		t.Error("Should fail with wrong password")
	}
	
	// 测试修改密码
	err = authService.ChangePassword("generated-id-123", "Password123!", "NewPassword456!", userProvider, userUpdater)
	if err != nil {
		t.Errorf("Failed to change password: %v", err)
	}
	
	// 验证新密码
	_, err = authService.Login("testuser", "NewPassword456!", userProvider)
	if err != nil {
		t.Errorf("Failed to login with new password: %v", err)
	}
}

type TestError struct {
	Message string
}

func (e *TestError) Error() string {
	return e.Message
}

func TestTokenExpiration(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:       "test-secret",
		ExpireHours:  0, // 立即过期
		RefreshHours: 1,
		Issuer:       "test",
	}
	
	jwtManager := NewJWTManager(cfg)
	
	user := &User{
		ID:       "123",
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"user"},
		IsActive: true,
	}
	
	// 生成立即过期的令牌
	token, err := jwtManager.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// 等待一小段时间确保令牌过期
	time.Sleep(time.Millisecond * 100)
	
	// 验证过期的令牌
	_, err = jwtManager.ValidateToken(token)
	if err == nil {
		t.Error("Expired token should not be valid")
	}
	
	if !strings.Contains(err.Error(), "expired") {
		t.Errorf("Error should mention expiration, got: %v", err)
	}
}

// 基准测试
func BenchmarkJWTGenerate(b *testing.B) {
	cfg := &config.JWTConfig{
		Secret:       "test-secret-key",
		ExpireHours:  24,
		RefreshHours: 168,
		Issuer:       "test-issuer",
	}
	
	jwtManager := NewJWTManager(cfg)
	
	user := &User{
		ID:       "123",
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"user", "admin"},
		IsActive: true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jwtManager.GenerateToken(user)
	}
}

func BenchmarkJWTValidate(b *testing.B) {
	cfg := &config.JWTConfig{
		Secret:       "test-secret-key",
		ExpireHours:  24,
		RefreshHours: 168,
		Issuer:       "test-issuer",
	}
	
	jwtManager := NewJWTManager(cfg)
	
	user := &User{
		ID:       "123",
		Username: "testuser",
		Email:    "test@example.com",
		Roles:    []string{"user", "admin"},
		IsActive: true,
	}
	
	token, _ := jwtManager.GenerateToken(user)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jwtManager.ValidateToken(token)
	}
}

func BenchmarkPasswordHash(b *testing.B) {
	pm := NewPasswordManager(10)
	password := "TestPassword123!"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pm.HashPassword(password)
	}
}