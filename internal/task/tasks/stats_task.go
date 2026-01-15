package tasks

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"gin_demo/pkg/task"
)

// StatsTask 统计任务
type StatsTask struct {
	db *sql.DB
}

// NewStatsTask 创建统计任务
func NewStatsTask(db *sql.DB) task.Task {
	return &StatsTask{
		db: db,
	}
}

func (t *StatsTask) Name() string {
	return "stats_task"
}

func (t *StatsTask) Spec() string {
	// 每小时执行一次
	return "0 0 * * * *"
}

func (t *StatsTask) Timeout() time.Duration {
	return 5 * time.Minute
}

func (t *StatsTask) Run(ctx context.Context) error {
	slog.Info("StatsTask: Calculating statistics...")
	
	// 示例：统计用户数量
	var userCount int64
	err := t.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		slog.Error("StatsTask: Failed to count users", "error", err)
		return err
	}
	
	slog.Info("StatsTask: Statistics calculated",
		"total_users", userCount,
	)
	
	// 实际应用中可以将统计结果写入数据库或发送到监控系统
	
	return nil
}
