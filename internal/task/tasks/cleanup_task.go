package tasks

import (
	"context"
	"log/slog"
	"time"

	"gin_demo/pkg/task"
	"github.com/redis/go-redis/v9"
)

// CleanupTask 数据清理任务
type CleanupTask struct {
	redis redis.UniversalClient
}

// NewCleanupTask 创建清理任务
func NewCleanupTask(redis redis.UniversalClient) task.Task {
	return &CleanupTask{
		redis: redis,
	}
}

func (t *CleanupTask) Name() string {
	return "cleanup_task"
}

func (t *CleanupTask) Spec() string {
	// 每天凌晨 2 点执行
	return "0 0 2 * * *"
}

func (t *CleanupTask) Timeout() time.Duration {
	return 10 * time.Minute
}

func (t *CleanupTask) Run(ctx context.Context) error {
	slog.Info("CleanupTask: Starting cleanup...")
	
	// 示例：清理过期的 Redis 缓存（这里只是示例）
	// 实际应用中可以清理过期数据、日志等
	
	// 1. 清理临时缓存
	pattern := "temp:*"
	iter := t.redis.Scan(ctx, 0, pattern, 100).Iterator()
	
	count := 0
	for iter.Next(ctx) {
		key := iter.Val()
		// 检查 TTL，如果没有过期时间则删除
		ttl, err := t.redis.TTL(ctx, key).Result()
		if err != nil {
			continue
		}
		
		if ttl == -1 { // 没有设置过期时间
			if err := t.redis.Del(ctx, key).Err(); err != nil {
				slog.Warn("Failed to delete key", "key", key, "error", err)
			} else {
				count++
			}
		}
	}
	
	if err := iter.Err(); err != nil {
		slog.Error("CleanupTask: Scan error", "error", err)
		return err
	}
	
	slog.Info("CleanupTask: Completed", "cleaned_keys", count)
	return nil
}
