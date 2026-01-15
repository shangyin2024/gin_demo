package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		// 构建日志字段
		fields := []interface{}{
			"method", c.Request.Method,
			"path", path,
			"status", statusCode,
			"latency", latency.String(),
			"ip", c.ClientIP(),
		}

		if query != "" {
			fields = append(fields, "query", query)
		}

		// 根据状态码选择日志级别
		switch {
		case statusCode >= 500:
			slog.ErrorContext(c.Request.Context(), "Server error", fields...)
		case statusCode >= 400:
			slog.WarnContext(c.Request.Context(), "Client error", fields...)
		default:
			slog.InfoContext(c.Request.Context(), "Request completed", fields...)
		}
	}
}
