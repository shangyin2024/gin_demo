package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Role 用户角色类型
type Role string

const (
	// RoleGuest 游客（未登录）
	RoleGuest Role = "guest"
	// RoleUser 普通用户
	RoleUser Role = "user"
	// RoleModerator 版主/审核员
	RoleModerator Role = "moderator"
	// RoleAdmin 管理员
	RoleAdmin Role = "admin"
	// RoleSuperAdmin 超级管理员
	RoleSuperAdmin Role = "super_admin"
)

// Permission 权限类型
type Permission string

const (
	// 用户权限
	PermissionUserRead   Permission = "user:read"
	PermissionUserWrite  Permission = "user:write"
	PermissionUserDelete Permission = "user:delete"

	// 内容权限
	PermissionContentRead   Permission = "content:read"
	PermissionContentWrite  Permission = "content:write"
	PermissionContentDelete Permission = "content:delete"
	PermissionContentAudit  Permission = "content:audit"

	// 系统权限
	PermissionSystemConfig Permission = "system:config"
	PermissionSystemMonitor Permission = "system:monitor"
)

// RBACClaims JWT 声明（包含角色和权限）
type RBACClaims struct {
	UserID      int64        `json:"user_id"`
	Role        Role         `json:"role"`                   // 用户角色
	Permissions []Permission `json:"permissions,omitempty"`  // 细粒度权限（可选）
	jwt.RegisteredClaims
}

// RBACJWTManager RBAC JWT 管理器
type RBACJWTManager struct {
	secret     string
	expiration time.Duration
}

// NewRBACJWTManager 创建 RBAC JWT 管理器
func NewRBACJWTManager(secret string, expiration time.Duration) *RBACJWTManager {
	return &RBACJWTManager{
		secret:     secret,
		expiration: expiration,
	}
}

// GenerateToken 生成包含角色信息的 JWT Token
func (m *RBACJWTManager) GenerateToken(userID int64, role Role, permissions ...Permission) (string, error) {
	now := time.Now()
	claims := RBACClaims{
		UserID:      userID,
		Role:        role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT Token
func (m *RBACJWTManager) ValidateToken(tokenString string) (*RBACClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RBACClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		if jwt.ErrTokenExpired != nil {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*RBACClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken 刷新 Token（保持原有角色和权限）
func (m *RBACJWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil && err != ErrExpiredToken {
		return "", err
	}

	// 即使 Token 过期，只要签名有效，就允许刷新
	return m.GenerateToken(claims.UserID, claims.Role, claims.Permissions...)
}

// HasRole 检查是否拥有指定角色
func (c *RBACClaims) HasRole(role Role) bool {
	return c.Role == role
}

// HasAnyRole 检查是否拥有任意指定角色
func (c *RBACClaims) HasAnyRole(roles ...Role) bool {
	for _, role := range roles {
		if c.Role == role {
			return true
		}
	}
	return false
}

// HasPermission 检查是否拥有指定权限
func (c *RBACClaims) HasPermission(permission Permission) bool {
	// 超级管理员拥有所有权限
	if c.Role == RoleSuperAdmin {
		return true
	}

	// 检查显式权限
	for _, p := range c.Permissions {
		if p == permission {
			return true
		}
	}

	// 检查基于角色的隐式权限
	return c.hasImplicitPermission(permission)
}

// HasAllPermissions 检查是否拥有所有指定权限
func (c *RBACClaims) HasAllPermissions(permissions ...Permission) bool {
	for _, permission := range permissions {
		if !c.HasPermission(permission) {
			return false
		}
	}
	return true
}

// HasAnyPermission 检查是否拥有任意指定权限
func (c *RBACClaims) HasAnyPermission(permissions ...Permission) bool {
	for _, permission := range permissions {
		if c.HasPermission(permission) {
			return true
		}
	}
	return false
}

// hasImplicitPermission 基于角色的隐式权限（权限继承）
func (c *RBACClaims) hasImplicitPermission(permission Permission) bool {
	switch c.Role {
	case RoleSuperAdmin:
		return true // 超级管理员拥有所有权限

	case RoleAdmin:
		// 管理员拥有大部分权限（除了系统配置）
		return permission != PermissionSystemConfig

	case RoleModerator:
		// 版主拥有内容审核权限
		return permission == PermissionContentRead ||
			permission == PermissionContentAudit ||
			permission == PermissionUserRead

	case RoleUser:
		// 普通用户只能读取和写入内容
		return permission == PermissionContentRead ||
			permission == PermissionContentWrite ||
			permission == PermissionUserRead

	case RoleGuest:
		// 游客只能读取内容
		return permission == PermissionContentRead

	default:
		return false
	}
}

// IsAdmin 是否是管理员（admin 或 super_admin）
func (c *RBACClaims) IsAdmin() bool {
	return c.Role == RoleAdmin || c.Role == RoleSuperAdmin
}

// IsSuperAdmin 是否是超级管理员
func (c *RBACClaims) IsSuperAdmin() bool {
	return c.Role == RoleSuperAdmin
}

// GetRoleLevel 获取角色级别（用于权限比较）
func (c *RBACClaims) GetRoleLevel() int {
	switch c.Role {
	case RoleSuperAdmin:
		return 100
	case RoleAdmin:
		return 80
	case RoleModerator:
		return 60
	case RoleUser:
		return 40
	case RoleGuest:
		return 0
	default:
		return -1
	}
}

// HasHigherRoleThan 是否拥有比指定角色更高的权限
func (c *RBACClaims) HasHigherRoleThan(role Role) bool {
	targetLevel := roleLevel(role)
	return c.GetRoleLevel() > targetLevel
}

// roleLevel 获取角色级别（辅助函数）
func roleLevel(role Role) int {
	switch role {
	case RoleSuperAdmin:
		return 100
	case RoleAdmin:
		return 80
	case RoleModerator:
		return 60
	case RoleUser:
		return 40
	case RoleGuest:
		return 0
	default:
		return -1
	}
}
