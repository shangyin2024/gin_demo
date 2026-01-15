package tasks

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gin_demo/pkg/task"
)

// ExampleTask 示例任务
type ExampleTask struct{}

// NewExampleTask 创建示例任务
func NewExampleTask() task.Task {
	return &ExampleTask{}
}

func (t *ExampleTask) Name() string {
	return "example_task"
}

func (t *ExampleTask) Spec() string {
	// 每分钟执行一次
	// Cron 格式: 秒 分 时 日 月 周
	return "0 * * * * *" // 每分钟的第 0 秒
}

func (t *ExampleTask) Timeout() time.Duration {
	return 30 * time.Second
}

func (t *ExampleTask) Run(ctx context.Context) error {
	slog.Info("ExampleTask: Starting...")
	
	// 模拟任务执行
	select {
	case <-time.After(2 * time.Second):
		slog.Info("ExampleTask: Completed successfully")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("task cancelled: %w", ctx.Err())
	}
}
