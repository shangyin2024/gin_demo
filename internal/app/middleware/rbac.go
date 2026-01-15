package middleware

import (
	"gin_demo/internal/response"
	"gin_demo/pkg/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// RBACClaimsKey RBAC Claims 在 context 中的键名
	RBACClaimsKey = "rbac_claims"
)

// RBACMiddleware RBAC 认证中间件
type RBACMiddleware struct {
	jwtManager *auth.RBACJWTManager
}

// NewRBACMiddleware 创建 RBAC 认证中间件
func NewRBACMiddleware(jwtManager *auth.RBACJWTManager) *RBACMiddleware {
	return &RBACMiddleware{
		jwtManager: jwtManager,
	}
}

// Handle 处理认证（提取并验证 Token 中的角色信息）
func (m *RBACMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization header
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Error(c, response.New(response.CodeUnauthorized, "未提供认证信息"))
			c.Abort()
			return
		}

		// 2. 检查 Bearer 前缀
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Error(c, response.New(response.CodeUnauthorized, "认证格式错误"))
			c.Abort()
			return
		}

		// 3. 提取 Token
		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == "" {
			response.Error(c, response.New(response.CodeUnauthorized, "Token 不能为空"))
			c.Abort()
			return
		}

		// 4. 验证 Token 并提取 Claims
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			switch err {
			case auth.ErrExpiredToken:
				response.Error(c, response.New(response.CodeUnauthorized, "Token 已过期"))
			case auth.ErrInvalidToken:
				response.Error(c, response.New(response.CodeUnauthorized, "Token 无效"))
			default:
				response.Error(c, response.NewWithError(response.CodeUnauthorized, "Token 验证失败", err))
			}
			c.Abort()
			return
		}

		// 5. 将 Claims 存入 context（包含用户ID、角色和权限）
		c.Set(UserIDKey, claims.UserID)
		c.Set(RBACClaimsKey, claims)

		c.Next()
	}
}

// RequireRole 要求指定角色（中间件）
func RequireRole(roles ...auth.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetRBACClaims(c)
		if claims == nil {
			response.Error(c, response.New(response.CodeUnauthorized, "未认证"))
			c.Abort()
			return
		}

		if !claims.HasAnyRole(roles...) {
			response.Error(c, response.New(response.CodeForbidden, "权限不足：需要角色 "+rolesToString(roles)))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission 要求指定权限（中间件）
func RequirePermission(permissions ...auth.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetRBACClaims(c)
		if claims == nil {
			response.Error(c, response.New(response.CodeUnauthorized, "未认证"))
			c.Abort()
			return
		}

		if !claims.HasAllPermissions(permissions...) {
			response.Error(c, response.New(response.CodeForbidden, "权限不足：缺少必要权限"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 要求任意权限（中间件）
func RequireAnyPermission(permissions ...auth.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetRBACClaims(c)
		if claims == nil {
			response.Error(c, response.New(response.CodeUnauthorized, "未认证"))
			c.Abort()
			return
		}

		if !claims.HasAnyPermission(permissions...) {
			response.Error(c, response.New(response.CodeForbidden, "权限不足：缺少任意所需权限"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin 要求管理员权限（中间件）
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(auth.RoleAdmin, auth.RoleSuperAdmin)
}

// RequireSuperAdmin 要求超级管理员权限（中间件）
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole(auth.RoleSuperAdmin)
}

// GetRBACClaims 从 context 中获取 RBAC Claims
func GetRBACClaims(c *gin.Context) *auth.RBACClaims {
	claims, exists := c.Get(RBACClaimsKey)
	if !exists {
		return nil
	}
	rbacClaims, ok := claims.(*auth.RBACClaims)
	if !ok {
		return nil
	}
	return rbacClaims
}

// GetUserRole 从 context 中获取用户角色
func GetUserRole(c *gin.Context) auth.Role {
	claims := GetRBACClaims(c)
	if claims == nil {
		return auth.RoleGuest
	}
	return claims.Role
}

// HasPermission 检查当前用户是否有指定权限
func HasPermission(c *gin.Context, permission auth.Permission) bool {
	claims := GetRBACClaims(c)
	if claims == nil {
		return false
	}
	return claims.HasPermission(permission)
}

// HasRole 检查当前用户是否有指定角色
func HasRole(c *gin.Context, role auth.Role) bool {
	claims := GetRBACClaims(c)
	if claims == nil {
		return false
	}
	return claims.HasRole(role)
}

// IsAdmin 检查当前用户是否是管理员
func IsAdmin(c *gin.Context) bool {
	claims := GetRBACClaims(c)
	if claims == nil {
		return false
	}
	return claims.IsAdmin()
}

// rolesToString 将角色列表转换为字符串（用于错误消息）
func rolesToString(roles []auth.Role) string {
	if len(roles) == 0 {
		return ""
	}
	roleStrs := make([]string, len(roles))
	for i, role := range roles {
		roleStrs[i] = string(role)
	}
	return strings.Join(roleStrs, ", ")
}
