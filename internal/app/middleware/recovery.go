package middleware

import (
	"gin_demo/internal/response"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 错误恢复中间件（增强版，基于 Gin 内置）
//
// 对比 Gin 内置 gin.Recovery()，增加了：
// - 结构化日志记录（slog）
// - 统一错误响应格式
// - Request ID 传递
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 信息和堆栈
				slog.ErrorContext(c.Request.Context(), "Panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)

				// 返回统一错误响应
				response.ErrorWithCode(c, response.CodeInternalError, "服务器内部错误")
				c.Abort()
			}
		}()

		c.Next()
	}
}

// RecoveryGin 使用 Gin 内置 Recovery（推荐生产环境）
//
// 使用示例:
//
//	engine.Use(middleware.RecoveryGin())
func RecoveryGin() gin.HandlerFunc {
	return gin.Recovery()
}
