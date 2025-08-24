package auth

import (
	"fmt"
	"testing"
)

func TestRBAC(t *testing.T) {
	rbac := NewRBAC()
	
	// 测试添加权限
	permission1 := &Permission{
		ID:          "user.read",
		Name:        "读取用户",
		Description: "读取用户信息的权限",
		Resource:    "user",
		Action:      "read",
	}
	
	permission2 := &Permission{
		ID:          "user.write",
		Name:        "写入用户",
		Description: "写入用户信息的权限",
		Resource:    "user",
		Action:      "write",
	}
	
	err := rbac.AddPermission(permission1)
	if err != nil {
		t.Fatalf("Failed to add permission: %v", err)
	}
	
	err = rbac.AddPermission(permission2)
	if err != nil {
		t.Fatalf("Failed to add permission: %v", err)
	}
	
	// 测试获取权限
	retrievedPermission, err := rbac.GetPermission("user.read")
	if err != nil {
		t.Fatalf("Failed to get permission: %v", err)
	}
	
	if retrievedPermission.ID != permission1.ID {
		t.Errorf("Expected permission ID %s, got %s", permission1.ID, retrievedPermission.ID)
	}
	
	// 测试添加角色
	role := &Role{
		ID:          "user",
		Name:        "普通用户",
		Description: "普通用户角色",
		Permissions: []Permission{},
	}
	
	err = rbac.AddRole(role)
	if err != nil {
		t.Fatalf("Failed to add role: %v", err)
	}
	
	// 测试为角色添加权限
	err = rbac.AddPermissionToRole("user", "user.read")
	if err != nil {
		t.Fatalf("Failed to add permission to role: %v", err)
	}
	
	err = rbac.AddPermissionToRole("user", "user.write")
	if err != nil {
		t.Fatalf("Failed to add permission to role: %v", err)
	}
	
	// 测试为用户分配角色
	userID := "123"
	err = rbac.AssignRoleToUser(userID, "user")
	if err != nil {
		t.Fatalf("Failed to assign role to user: %v", err)
	}
	
	// 测试检查用户角色
	if !rbac.HasRole(userID, "user") {
		t.Error("User should have 'user' role")
	}
	
	if rbac.HasRole(userID, "admin") {
		t.Error("User should not have 'admin' role")
	}
	
	// 测试检查用户权限
	if !rbac.HasPermission(userID, "user.read") {
		t.Error("User should have 'user.read' permission")
	}
	
	if !rbac.HasPermission(userID, "user.write") {
		t.Error("User should have 'user.write' permission")
	}
	
	// 测试资源权限检查
	if !rbac.HasResourcePermission(userID, "user", "read") {
		t.Error("User should have read permission on user resource")
	}
	
	if rbac.HasResourcePermission(userID, "admin", "read") {
		t.Error("User should not have read permission on admin resource")
	}
}

func TestRBACRoleManagement(t *testing.T) {
	rbac := NewRBAC()
	
	// 添加权限
	permissions := []*Permission{
		{ID: "p1", Name: "Permission 1", Resource: "resource1", Action: "read"},
		{ID: "p2", Name: "Permission 2", Resource: "resource2", Action: "write"},
		{ID: "p3", Name: "Permission 3", Resource: "resource3", Action: "delete"},
	}
	
	for _, p := range permissions {
		rbac.AddPermission(p)
	}
	
	// 添加角色
	role1 := &Role{
		ID:          "role1",
		Name:        "Role 1",
		Description: "Test role 1",
		Permissions: []Permission{},
	}
	
	role2 := &Role{
		ID:          "role2",
		Name:        "Role 2",
		Description: "Test role 2",
		Permissions: []Permission{},
	}
	
	rbac.AddRole(role1)
	rbac.AddRole(role2)
	
	// 为角色添加权限
	rbac.AddPermissionToRole("role1", "p1")
	rbac.AddPermissionToRole("role1", "p2")
	rbac.AddPermissionToRole("role2", "p2")
	rbac.AddPermissionToRole("role2", "p3")
	
	// 为用户分配多个角色
	userID := "user123"
	rbac.AssignRoleToUser(userID, "role1")
	rbac.AssignRoleToUser(userID, "role2")
	
	// 测试用户拥有的所有权限
	userPermissions := rbac.GetUserPermissions(userID)
	if len(userPermissions) != 3 {
		t.Errorf("Expected 3 permissions, got %d", len(userPermissions))
	}
	
	// 测试用户角色
	userRoles := rbac.GetUserRoles(userID)
	if len(userRoles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(userRoles))
	}
	
	// 测试 HasAnyRole
	if !rbac.HasAnyRole(userID, []string{"role1", "role3"}) {
		t.Error("User should have at least one of the roles")
	}
	
	// 测试 HasAllRoles
	if !rbac.HasAllRoles(userID, []string{"role1", "role2"}) {
		t.Error("User should have all specified roles")
	}
	
	if rbac.HasAllRoles(userID, []string{"role1", "role3"}) {
		t.Error("User should not have all specified roles")
	}
}

func TestRBACPermissionRemoval(t *testing.T) {
	rbac := NewRBAC()
	
	// 添加权限和角色
	permission := &Permission{
		ID:       "test.permission",
		Name:     "Test Permission",
		Resource: "test",
		Action:   "read",
	}
	
	rbac.AddPermission(permission)
	
	role := &Role{
		ID:          "test.role",
		Name:        "Test Role",
		Description: "Test role",
		Permissions: []Permission{},
	}
	
	rbac.AddRole(role)
	rbac.AddPermissionToRole("test.role", "test.permission")
	
	userID := "testuser"
	rbac.AssignRoleToUser(userID, "test.role")
	
	// 验证用户有权限
	if !rbac.HasPermission(userID, "test.permission") {
		t.Error("User should have permission before removal")
	}
	
	// 从角色中移除权限
	err := rbac.RemovePermissionFromRole("test.role", "test.permission")
	if err != nil {
		t.Fatalf("Failed to remove permission from role: %v", err)
	}
	
	// 验证用户不再有权限
	if rbac.HasPermission(userID, "test.permission") {
		t.Error("User should not have permission after removal")
	}
}

func TestRBACRoleRemoval(t *testing.T) {
	rbac := NewRBAC()
	
	// 添加角色
	role := &Role{
		ID:          "test.role",
		Name:        "Test Role",
		Description: "Test role",
		Permissions: []Permission{},
	}
	
	rbac.AddRole(role)
	
	userID := "testuser"
	rbac.AssignRoleToUser(userID, "test.role")
	
	// 验证用户有角色
	if !rbac.HasRole(userID, "test.role") {
		t.Error("User should have role before removal")
	}
	
	// 移除角色
	err := rbac.RemoveRole("test.role")
	if err != nil {
		t.Fatalf("Failed to remove role: %v", err)
	}
	
	// 验证用户不再有角色
	if rbac.HasRole(userID, "test.role") {
		t.Error("User should not have role after removal")
	}
}

func TestCreateDefaultRolesAndPermissions(t *testing.T) {
	rbac := NewRBAC()
	
	err := CreateDefaultRolesAndPermissions(rbac)
	if err != nil {
		t.Fatalf("Failed to create default roles and permissions: %v", err)
	}
	
	// 验证默认权限
	_, err = rbac.GetPermission("user.read")
	if err != nil {
		t.Error("Should have default user.read permission")
	}
	
	_, err = rbac.GetPermission("admin.all")
	if err != nil {
		t.Error("Should have default admin.all permission")
	}
	
	// 验证默认角色
	_, err = rbac.GetRole("user")
	if err != nil {
		t.Error("Should have default user role")
	}
	
	_, err = rbac.GetRole("admin")
	if err != nil {
		t.Error("Should have default admin role")
	}
	
	// 测试分配默认角色
	userID := "testuser"
	rbac.AssignRoleToUser(userID, "admin")
	
	// 管理员应该有所有权限
	if !rbac.HasResourcePermission(userID, "any", "any") {
		t.Error("Admin should have permission on any resource")
	}
}

func TestPolicyEvaluator(t *testing.T) {
	rbac := NewRBAC()
	
	// 创建测试权限和角色
	permissions := []*Permission{
		{ID: "read.user", Name: "Read User", Resource: "user", Action: "read"},
		{ID: "write.user", Name: "Write User", Resource: "user", Action: "write"},
		{ID: "read.admin", Name: "Read Admin", Resource: "admin", Action: "read"},
	}
	
	for _, p := range permissions {
		rbac.AddPermission(p)
	}
	
	userRole := &Role{
		ID:          "user",
		Name:        "User",
		Description: "Regular user",
		Permissions: []Permission{*permissions[0]}, // 只有读取用户权限
	}
	
	adminRole := &Role{
		ID:          "admin",
		Name:        "Admin",
		Description: "Administrator",
		Permissions: []Permission{*permissions[0], *permissions[1], *permissions[2]}, // 所有权限
	}
	
	rbac.AddRole(userRole)
	rbac.AddRole(adminRole)
	
	// 分配角色
	regularUserID := "user123"
	adminUserID := "admin123"
	
	rbac.AssignRoleToUser(regularUserID, "user")
	rbac.AssignRoleToUser(adminUserID, "admin")
	
	// 创建策略评估器
	evaluator := NewPolicyEvaluator(rbac)
	
	// 测试策略
	readUserPolicy := &Policy{
		Resource: "user",
		Actions:  []string{"read"},
		Roles:    []string{"user", "admin"},
	}
	
	writeUserPolicy := &Policy{
		Resource: "user",
		Actions:  []string{"write"},
		Roles:    []string{"admin"},
	}
	
	adminOnlyPolicy := &Policy{
		Resource: "admin",
		Actions:  []string{"read"},
		Roles:    []string{"admin"},
	}
	
	// 测试普通用户
	if !evaluator.EvaluatePolicy(regularUserID, readUserPolicy) {
		t.Error("Regular user should be able to read user data")
	}
	
	if evaluator.EvaluatePolicy(regularUserID, writeUserPolicy) {
		t.Error("Regular user should not be able to write user data")
	}
	
	if evaluator.EvaluatePolicy(regularUserID, adminOnlyPolicy) {
		t.Error("Regular user should not access admin resources")
	}
	
	// 测试管理员用户
	if !evaluator.EvaluatePolicy(adminUserID, readUserPolicy) {
		t.Error("Admin should be able to read user data")
	}
	
	if !evaluator.EvaluatePolicy(adminUserID, writeUserPolicy) {
		t.Error("Admin should be able to write user data")
	}
	
	if !evaluator.EvaluatePolicy(adminUserID, adminOnlyPolicy) {
		t.Error("Admin should access admin resources")
	}
}

func TestResourceMatcher(t *testing.T) {
	matcher := NewResourceMatcher()
	
	// 测试精确匹配
	if !matcher.Match("user", "user") {
		t.Error("Should match exactly")
	}
	
	if matcher.Match("user", "admin") {
		t.Error("Should not match different resources")
	}
	
	// 测试通配符匹配
	if !matcher.Match("*", "user") {
		t.Error("Wildcard should match any resource")
	}
	
	if !matcher.Match("*", "admin") {
		t.Error("Wildcard should match any resource")
	}
	
	// 测试前缀匹配
	if !matcher.Match("user.*", "user.profile") {
		t.Error("Should match with prefix wildcard")
	}
	
	if !matcher.Match("user.*", "user.settings") {
		t.Error("Should match with prefix wildcard")
	}
	
	if matcher.Match("user.*", "admin.settings") {
		t.Error("Should not match different prefix")
	}
}

func TestRBACListMethods(t *testing.T) {
	rbac := NewRBAC()
	
	// 添加测试数据
	permissions := []*Permission{
		{ID: "p1", Name: "Permission 1", Resource: "r1", Action: "read"},
		{ID: "p2", Name: "Permission 2", Resource: "r2", Action: "write"},
	}
	
	roles := []*Role{
		{ID: "r1", Name: "Role 1", Description: "Test role 1"},
		{ID: "r2", Name: "Role 2", Description: "Test role 2"},
	}
	
	for _, p := range permissions {
		rbac.AddPermission(p)
	}
	
	for _, r := range roles {
		rbac.AddRole(r)
	}
	
	// 测试列表方法
	allPermissions := rbac.ListPermissions()
	if len(allPermissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(allPermissions))
	}
	
	allRoles := rbac.ListRoles()
	if len(allRoles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(allRoles))
	}
	
	// 测试用户角色详情
	userID := "testuser"
	rbac.AssignRoleToUser(userID, "r1")
	rbac.AssignRoleToUser(userID, "r2")
	
	userRoleDetails := rbac.GetRolesByUser(userID)
	if len(userRoleDetails) != 2 {
		t.Errorf("Expected 2 role details, got %d", len(userRoleDetails))
	}
}

// 基准测试
func BenchmarkRBACHasPermission(b *testing.B) {
	rbac := NewRBAC()
	
	// 设置测试数据
	for i := 0; i < 100; i++ {
		permission := &Permission{
			ID:       fmt.Sprintf("permission_%d", i),
			Name:     fmt.Sprintf("Permission %d", i),
			Resource: fmt.Sprintf("resource_%d", i),
			Action:   "read",
		}
		rbac.AddPermission(permission)
	}
	
	role := &Role{
		ID:          "test_role",
		Name:        "Test Role",
		Description: "Test role with many permissions",
		Permissions: []Permission{},
	}
	rbac.AddRole(role)
	
	for i := 0; i < 100; i++ {
		rbac.AddPermissionToRole("test_role", fmt.Sprintf("permission_%d", i))
	}
	
	userID := "test_user"
	rbac.AssignRoleToUser(userID, "test_role")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbac.HasPermission(userID, "permission_50")
	}
}
