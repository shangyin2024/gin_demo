package middleware

import (
	"gin_demo/internal/response"
	"gin_demo/pkg/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// AuthorizationHeader Authorization 请求头
	AuthorizationHeader = "Authorization"
	// BearerPrefix Bearer 前缀
	BearerPrefix = "Bearer "
	// UserIDKey 用户 ID 在 context 中的键名
	UserIDKey = "user_id"
)

// ============================================================================
// 认证中间件（内部函数，由 AuthMiddleware 使用）
// ============================================================================

// authenticateRequest 认证请求（内部函数）
func authenticateRequest(c *gin.Context, jwtManager *auth.DefaultJWTManager) (*auth.DefaultClaims, error) {
	// 1. 获取 Authorization header
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader == "" {
		return nil, response.New(response.CodeUnauthorized, "未提供认证信息")
	}

	// 2. 检查 Bearer 前缀
	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return nil, response.New(response.CodeUnauthorized, "认证格式错误")
	}

	// 3. 提取 Token
	tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
	if tokenString == "" {
		return nil, response.New(response.CodeUnauthorized, "Token 不能为空")
	}

	// 4. 验证 Token
	claims, err := jwtManager.ValidateToken(tokenString)
	if err != nil {
		switch err {
		case auth.ErrExpiredToken:
			return nil, response.New(response.CodeUnauthorized, "Token 已过期")
		case auth.ErrInvalidToken:
			return nil, response.New(response.CodeUnauthorized, "Token 无效")
		default:
			return nil, response.NewWithError(response.CodeUnauthorized, "Token 验证失败", err)
		}
	}

	return claims, nil
}

// ============================================================================
// 兼容性函数（保持向后兼容）
// ============================================================================

// Auth JWT 认证中间件（函数式，保留用于向后兼容）
// 推荐使用 AuthMiddleware 结构体版本
func Auth(jwtManager *auth.DefaultJWTManager) gin.HandlerFunc {
	middleware := NewAuthMiddleware(jwtManager)
	return middleware.Handle()
}

// OptionalAuth 可选认证中间件（Token 存在时验证，不存在时也放行）
func OptionalAuth(jwtManager *auth.DefaultJWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := jwtManager.ValidateToken(tokenString)
		if err == nil {
			c.Set(UserIDKey, claims.UserID)
		}

		c.Next()
	}
}

// GetUserID 从 context 中获取用户 ID
// 返回 0 表示未找到（通常在已认证路由中不会出现）
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0
	}
	id, ok := userID.(int64)
	if !ok {
		return 0
	}
	return id
}
