package middleware

import (
	"gin_demo/internal/response"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流器（基于官方 golang.org/x/time/rate）
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit // 每秒产生的令牌数
	burst    int        // 桶容量
}

// NewRateLimiter 创建限流器
// r: 每秒允许的请求数（QPS）
// burst: 允许的突发请求数（桶容量）
//
// 示例:
//   - NewRateLimiter(10, 20)  // 每秒 10 个请求，最多突发 20 个
//   - NewRateLimiter(100, 200) // 每秒 100 个请求，最多突发 200 个
func NewRateLimiter(r float64, burst int) *RateLimiter {
	limiter := &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(r),
		burst:    burst,
	}

	// 定期清理过期的限流器
	go limiter.cleanup()

	return limiter
}

// getLimiter 获取或创建特定 key 的限流器
func (l *RateLimiter) getLimiter(key string) *rate.Limiter {
	l.mu.RLock()
	limiter, exists := l.limiters[key]
	l.mu.RUnlock()

	if !exists {
		l.mu.Lock()
		// Double-check
		limiter, exists = l.limiters[key]
		if !exists {
			limiter = rate.NewLimiter(l.rate, l.burst)
			l.limiters[key] = limiter
		}
		l.mu.Unlock()
	}

	return limiter
}

// Allow 检查是否允许请求
func (l *RateLimiter) Allow(key string) bool {
	limiter := l.getLimiter(key)
	return limiter.Allow()
}

// cleanup 定期清理不活跃的限流器（防止内存泄漏）
func (l *RateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		for key, limiter := range l.limiters {
			// 如果限流器的令牌已经满了，说明长时间没有请求
			if limiter.Tokens() == float64(l.burst) {
				delete(l.limiters, key)
			}
		}
		l.mu.Unlock()
	}
}

// RateLimit 限流中间件
//
// 使用示例:
//
//	// 方式 1: 基于 IP 限流（默认）
//	limiter := NewRateLimiter(10, 20)  // 每秒 10 个请求
//	router.Use(RateLimit(limiter))
//
//	// 方式 2: 自定义限流 Key
//	router.Use(RateLimitWithKeyFunc(limiter, func(c *gin.Context) string {
//	    // 根据用户 ID 限流
//	    return c.GetString("user_id")
//	}))
func RateLimit(limiter *RateLimiter) gin.HandlerFunc {
	return RateLimitWithKeyFunc(limiter, func(c *gin.Context) string {
		// 默认使用客户端 IP 作为限流 Key
		return c.ClientIP()
	})
}

// RateLimitWithKeyFunc 支持自定义 Key 的限流中间件
func RateLimitWithKeyFunc(limiter *RateLimiter, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)

		if !limiter.Allow(key) {
			response.ErrorWithCode(c, response.CodeTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
