package health

import (
	"context"
	"gin_demo/pkg/health"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisChecker Redis 健康检查器
type RedisChecker struct {
	redis redis.UniversalClient
}

// NewRedisChecker 创建 Redis 检查器
func NewRedisChecker(rdb redis.UniversalClient) *RedisChecker {
	return &RedisChecker{redis: rdb}
}

// Name 返回组件名称
func (r *RedisChecker) Name() string {
	return "redis"
}

// Check 执行 Redis 检查
func (r *RedisChecker) Check(ctx context.Context) health.Check {
	start := time.Now()
	check := health.Check{
		Status: health.StatusOK,
	}

	// Ping Redis
	if err := r.redis.Ping(ctx).Err(); err != nil {
		check.Status = health.StatusError
		check.Message = err.Error()
		check.Duration = time.Since(start).String()
		return check
	}

	check.Duration = time.Since(start).String()
	return check
}

// IsCritical Redis 不是关键组件（可以降级服务）
func (r *RedisChecker) IsCritical() bool {
	return false
}
