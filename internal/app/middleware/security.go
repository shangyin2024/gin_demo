package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// SecurityConfig 安全配置
type SecurityConfig struct {
	// 是否启用 HSTS (HTTP Strict Transport Security)
	EnableHSTS bool
	// HSTS 最大时间（秒）
	HSTSMaxAge int
	// 是否包含子域名
	HSTSIncludeSubdomains bool
	// 是否启用 CSP (Content Security Policy)
	EnableCSP bool
	// CSP 策略
	CSPPolicy string
	// 是否启用 X-Frame-Options
	EnableFrameOptions bool
	// Frame Options 值: DENY, SAMEORIGIN
	FrameOptions string
}

// DefaultSecurityConfig 默认安全配置
var DefaultSecurityConfig = SecurityConfig{
	EnableHSTS:            true,
	HSTSMaxAge:            31536000, // 1 年
	HSTSIncludeSubdomains: true,
	EnableCSP:             true,
	CSPPolicy:             "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';",
	EnableFrameOptions:    true,
	FrameOptions:          "DENY",
}

// Security HTTP 安全头中间件
func Security(config ...SecurityConfig) gin.HandlerFunc {
	cfg := DefaultSecurityConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *gin.Context) {
		// 1. X-Content-Type-Options
		// 防止浏览器进行 MIME 类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")

		// 2. X-Frame-Options
		// 防止点击劫持攻击
		if cfg.EnableFrameOptions {
			c.Header("X-Frame-Options", cfg.FrameOptions)
		}

		// 3. X-XSS-Protection
		// 启用浏览器的 XSS 过滤器
		c.Header("X-XSS-Protection", "1; mode=block")

		// 4. Referrer-Policy
		// 控制 Referer 头的发送
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 5. Content-Security-Policy (CSP)
		// 防止 XSS 和其他注入攻击
		if cfg.EnableCSP {
			c.Header("Content-Security-Policy", cfg.CSPPolicy)
		}

		// 6. Strict-Transport-Security (HSTS)
		// 强制使用 HTTPS
		if cfg.EnableHSTS && c.Request.TLS != nil {
			hstsValue := "max-age=" + strconv.Itoa(cfg.HSTSMaxAge)
			if cfg.HSTSIncludeSubdomains {
				hstsValue += "; includeSubDomains"
			}
			c.Header("Strict-Transport-Security", hstsValue)
		}

		// 7. Permissions-Policy (替代 Feature-Policy)
		// 控制浏览器功能和 API 的访问
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// 8. X-Permitted-Cross-Domain-Policies
		// 控制跨域策略文件
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		// 9. X-Download-Options
		// 防止 IE 自动执行下载的文件
		c.Header("X-Download-Options", "noopen")

		c.Next()
	}
}

// SecureHeaders 简化版安全头（用于开发环境）
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}
