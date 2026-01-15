package task

import (
	"database/sql"
	"log/slog"

	"gin_demo/internal/task/tasks"
	"gin_demo/pkg/task"
	"github.com/redis/go-redis/v9"
)

// Manager 任务管理器
type Manager struct {
	scheduler *task.Scheduler
}

// NewManager 创建任务管理器
func NewManager(redis redis.UniversalClient, db *sql.DB) *Manager {
	// 创建调度器
	scheduler := task.NewScheduler(task.Config{
		Redis:      redis,
		LockPrefix: "task:lock:",
	})
	
	// 注册所有任务
	registerTasks(scheduler, redis, db)
	
	return &Manager{
		scheduler: scheduler,
	}
}

// registerTasks 注册所有任务
func registerTasks(scheduler *task.Scheduler, redis redis.UniversalClient, db *sql.DB) {
	taskList := []task.Task{
		tasks.NewExampleTask(),
		tasks.NewCleanupTask(redis),
		tasks.NewStatsTask(db),
		// 在这里添加更多任务...
	}

	for _, t := range taskList {
		if err := scheduler.Register(t); err != nil {
			// 注册失败时记录错误并跳过该任务，不影响其他任务
			slog.Error("Failed to register task", "task", t.Name(), "error", err)
		}
	}
}

// Start 启动任务调度
func (m *Manager) Start() {
	m.scheduler.Start()
}

// Stop 停止任务调度
func (m *Manager) Stop() {
	m.scheduler.Stop()
}

// ListTasks 列出所有任务
func (m *Manager) ListTasks() []string {
	return m.scheduler.ListTasks()
}
