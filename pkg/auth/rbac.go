package auth

import (
	"errors"
	"strings"
)

// Permission 权限
type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
}

// Role 角色
type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}

// RBAC RBAC权限管理器
type RBAC struct {
	roles       map[string]*Role
	permissions map[string]*Permission
	userRoles   map[string][]string // userID -> roleNames
}

// NewRBAC 创建RBAC管理器
func NewRBAC() *RBAC {
	return &RBAC{
		roles:       make(map[string]*Role),
		permissions: make(map[string]*Permission),
		userRoles:   make(map[string][]string),
	}
}

// AddPermission 添加权限
func (rbac *RBAC) AddPermission(permission *Permission) error {
	if permission.ID == "" {
		return errors.New("permission ID cannot be empty")
	}
	if permission.Name == "" {
		return errors.New("permission name cannot be empty")
	}
	
	rbac.permissions[permission.ID] = permission
	return nil
}

// GetPermission 获取权限
func (rbac *RBAC) GetPermission(permissionID string) (*Permission, error) {
	permission, exists := rbac.permissions[permissionID]
	if !exists {
		return nil, errors.New("permission not found")
	}
	return permission, nil
}

// RemovePermission 移除权限
func (rbac *RBAC) RemovePermission(permissionID string) error {
	if _, exists := rbac.permissions[permissionID]; !exists {
		return errors.New("permission not found")
	}
	
	// 从所有角色中移除该权限
	for _, role := range rbac.roles {
		rbac.removePermissionFromRole(role, permissionID)
	}
	
	delete(rbac.permissions, permissionID)
	return nil
}

// AddRole 添加角色
func (rbac *RBAC) AddRole(role *Role) error {
	if role.ID == "" {
		return errors.New("role ID cannot be empty")
	}
	if role.Name == "" {
		return errors.New("role name cannot be empty")
	}
	
	rbac.roles[role.ID] = role
	return nil
}

// GetRole 获取角色
func (rbac *RBAC) GetRole(roleID string) (*Role, error) {
	role, exists := rbac.roles[roleID]
	if !exists {
		return nil, errors.New("role not found")
	}
	return role, nil
}

// RemoveRole 移除角色
func (rbac *RBAC) RemoveRole(roleID string) error {
	if _, exists := rbac.roles[roleID]; !exists {
		return errors.New("role not found")
	}
	
	// 从所有用户中移除该角色
	for userID, roles := range rbac.userRoles {
		rbac.userRoles[userID] = rbac.removeRoleFromSlice(roles, roleID)
	}
	
	delete(rbac.roles, roleID)
	return nil
}

// AddPermissionToRole 为角色添加权限
func (rbac *RBAC) AddPermissionToRole(roleID, permissionID string) error {
	role, exists := rbac.roles[roleID]
	if !exists {
		return errors.New("role not found")
	}
	
	permission, exists := rbac.permissions[permissionID]
	if !exists {
		return errors.New("permission not found")
	}
	
	// 检查权限是否已存在
	for _, p := range role.Permissions {
		if p.ID == permissionID {
			return nil // 权限已存在，不重复添加
		}
	}
	
	role.Permissions = append(role.Permissions, *permission)
	return nil
}

// RemovePermissionFromRole 从角色中移除权限
func (rbac *RBAC) RemovePermissionFromRole(roleID, permissionID string) error {
	role, exists := rbac.roles[roleID]
	if !exists {
		return errors.New("role not found")
	}
	
	rbac.removePermissionFromRole(role, permissionID)
	return nil
}

// AssignRoleToUser 为用户分配角色
func (rbac *RBAC) AssignRoleToUser(userID, roleID string) error {
	if _, exists := rbac.roles[roleID]; !exists {
		return errors.New("role not found")
	}
	
	userRoles := rbac.userRoles[userID]
	
	// 检查角色是否已分配
	for _, role := range userRoles {
		if role == roleID {
			return nil // 角色已分配，不重复添加
		}
	}
	
	rbac.userRoles[userID] = append(userRoles, roleID)
	return nil
}

// RemoveRoleFromUser 从用户中移除角色
func (rbac *RBAC) RemoveRoleFromUser(userID, roleID string) error {
	userRoles, exists := rbac.userRoles[userID]
	if !exists {
		return errors.New("user has no roles")
	}
	
	rbac.userRoles[userID] = rbac.removeRoleFromSlice(userRoles, roleID)
	return nil
}

// GetUserRoles 获取用户角色
func (rbac *RBAC) GetUserRoles(userID string) []string {
	return rbac.userRoles[userID]
}

// GetUserPermissions 获取用户所有权限
func (rbac *RBAC) GetUserPermissions(userID string) []Permission {
	var permissions []Permission
	userRoles := rbac.userRoles[userID]
	
	for _, roleID := range userRoles {
		if role, exists := rbac.roles[roleID]; exists {
			permissions = append(permissions, role.Permissions...)
		}
	}
	
	// 去重
	return rbac.deduplicatePermissions(permissions)
}

// HasPermission 检查用户是否拥有指定权限
func (rbac *RBAC) HasPermission(userID, permissionID string) bool {
	userPermissions := rbac.GetUserPermissions(userID)
	
	for _, permission := range userPermissions {
		if permission.ID == permissionID {
			return true
		}
	}
	
	return false
}

// HasResourcePermission 检查用户是否拥有资源权限
func (rbac *RBAC) HasResourcePermission(userID, resource, action string) bool {
	userPermissions := rbac.GetUserPermissions(userID)
	
	for _, permission := range userPermissions {
		if permission.Resource == resource && permission.Action == action {
			return true
		}
		// 支持通配符
		if permission.Resource == "*" || permission.Action == "*" {
			return true
		}
	}
	
	return false
}

// HasRole 检查用户是否拥有指定角色
func (rbac *RBAC) HasRole(userID, roleID string) bool {
	userRoles := rbac.userRoles[userID]
	
	for _, role := range userRoles {
		if role == roleID {
			return true
		}
	}
	
	return false
}

// HasAnyRole 检查用户是否拥有任意指定角色
func (rbac *RBAC) HasAnyRole(userID string, roleIDs []string) bool {
	for _, roleID := range roleIDs {
		if rbac.HasRole(userID, roleID) {
			return true
		}
	}
	return false
}

// HasAllRoles 检查用户是否拥有所有指定角色
func (rbac *RBAC) HasAllRoles(userID string, roleIDs []string) bool {
	for _, roleID := range roleIDs {
		if !rbac.HasRole(userID, roleID) {
			return false
		}
	}
	return true
}

// ListRoles 列出所有角色
func (rbac *RBAC) ListRoles() []*Role {
	var roles []*Role
	for _, role := range rbac.roles {
		roles = append(roles, role)
	}
	return roles
}

// ListPermissions 列出所有权限
func (rbac *RBAC) ListPermissions() []*Permission {
	var permissions []*Permission
	for _, permission := range rbac.permissions {
		permissions = append(permissions, permission)
	}
	return permissions
}

// GetRolesByUser 获取用户的所有角色详情
func (rbac *RBAC) GetRolesByUser(userID string) []*Role {
	var roles []*Role
	userRoles := rbac.userRoles[userID]
	
	for _, roleID := range userRoles {
		if role, exists := rbac.roles[roleID]; exists {
			roles = append(roles, role)
		}
	}
	
	return roles
}

// 辅助方法

// removePermissionFromRole 从角色中移除权限
func (rbac *RBAC) removePermissionFromRole(role *Role, permissionID string) {
	for i, permission := range role.Permissions {
		if permission.ID == permissionID {
			role.Permissions = append(role.Permissions[:i], role.Permissions[i+1:]...)
			break
		}
	}
}

// removeRoleFromSlice 从角色切片中移除指定角色
func (rbac *RBAC) removeRoleFromSlice(roles []string, roleID string) []string {
	for i, role := range roles {
		if role == roleID {
			return append(roles[:i], roles[i+1:]...)
		}
	}
	return roles
}

// deduplicatePermissions 去重权限
func (rbac *RBAC) deduplicatePermissions(permissions []Permission) []Permission {
	seen := make(map[string]bool)
	var result []Permission
	
	for _, permission := range permissions {
		if !seen[permission.ID] {
			seen[permission.ID] = true
			result = append(result, permission)
		}
	}
	
	return result
}

// PolicyEvaluator 策略评估器
type PolicyEvaluator struct {
	rbac *RBAC
}

// NewPolicyEvaluator 创建策略评估器
func NewPolicyEvaluator(rbac *RBAC) *PolicyEvaluator {
	return &PolicyEvaluator{
		rbac: rbac,
	}
}

// Policy 访问策略
type Policy struct {
	Resource    string   `json:"resource"`
	Actions     []string `json:"actions"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// EvaluatePolicy 评估策略
func (pe *PolicyEvaluator) EvaluatePolicy(userID string, policy *Policy) bool {
	// 检查角色
	if len(policy.Roles) > 0 {
		if !pe.rbac.HasAnyRole(userID, policy.Roles) {
			return false
		}
	}
	
	// 检查权限
	if len(policy.Permissions) > 0 {
		for _, permissionID := range policy.Permissions {
			if !pe.rbac.HasPermission(userID, permissionID) {
				return false
			}
		}
	}
	
	// 检查资源和动作
	for _, action := range policy.Actions {
		if !pe.rbac.HasResourcePermission(userID, policy.Resource, action) {
			return false
		}
	}
	
	return true
}

// 预定义角色和权限

// CreateDefaultRolesAndPermissions 创建默认角色和权限
func CreateDefaultRolesAndPermissions(rbac *RBAC) error {
	// 创建基本权限
	permissions := []*Permission{
		{ID: "user.read", Name: "读取用户", Resource: "user", Action: "read"},
		{ID: "user.write", Name: "写入用户", Resource: "user", Action: "write"},
		{ID: "user.delete", Name: "删除用户", Resource: "user", Action: "delete"},
		{ID: "admin.all", Name: "管理员权限", Resource: "*", Action: "*"},
	}
	
	for _, permission := range permissions {
		if err := rbac.AddPermission(permission); err != nil {
			return err
		}
	}
	
	// 创建基本角色
	roles := []*Role{
		{
			ID:          "user",
			Name:        "普通用户",
			Description: "普通用户角色，只能读取用户信息",
			Permissions: []Permission{*permissions[0]},
		},
		{
			ID:          "admin",
			Name:        "管理员",
			Description: "管理员角色，拥有所有权限",
			Permissions: []Permission{*permissions[3]},
		},
	}
	
	for _, role := range roles {
		if err := rbac.AddRole(role); err != nil {
			return err
		}
	}
	
	return nil
}

// ResourceMatcher 资源匹配器
type ResourceMatcher struct{}

// NewResourceMatcher 创建资源匹配器
func NewResourceMatcher() *ResourceMatcher {
	return &ResourceMatcher{}
}

// Match 匹配资源
func (rm *ResourceMatcher) Match(pattern, resource string) bool {
	// 简单的通配符匹配
	if pattern == "*" {
		return true
	}
	
	if pattern == resource {
		return true
	}
	
	// 支持前缀匹配，如 "user.*" 匹配 "user.profile"
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(resource, prefix)
	}
	
	return false
}