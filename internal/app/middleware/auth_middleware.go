package middleware

import (
	"gin_demo/internal/response"
	"gin_demo/pkg/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件（推荐使用的标准结构体风格）
type AuthMiddleware struct {
	jwtManager *auth.DefaultJWTManager
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(jwtManager *auth.DefaultJWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// Handle 处理认证（标准接口）
func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := authenticateRequest(c, m.jwtManager)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		// 将用户 ID 存入 context
		c.Set(UserIDKey, claims.UserID)

		c.Next()
	}
}

// HandleOptional 可选认证（Token 存在时验证，不存在时也放行）
func (m *AuthMiddleware) HandleOptional() gin.HandlerFunc {
	return OptionalAuth(m.jwtManager)
}
