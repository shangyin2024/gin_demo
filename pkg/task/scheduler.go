package task

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

// Task 任务接口
type Task interface {
	// Name 任务名称（唯一标识）
	Name() string
	
	// Spec Cron 表达式
	Spec() string
	
	// Run 执行任务
	Run(ctx context.Context) error
	
	// Timeout 任务超时时间
	Timeout() time.Duration
}

// Scheduler 任务调度器
type Scheduler struct {
	cron       *cron.Cron
	redis      redis.UniversalClient
	tasks      map[string]Task
	mu         sync.RWMutex
	lockPrefix string // Redis 锁前缀
	lockTTL    time.Duration
}

// Config 调度器配置
type Config struct {
	// Redis 客户端
	Redis redis.UniversalClient
	
	// 锁前缀
	LockPrefix string
	
	// 锁 TTL（默认为任务超时时间）
	LockTTL time.Duration
	
	// 时区
	Location *time.Location
}

// NewScheduler 创建任务调度器
func NewScheduler(config Config) *Scheduler {
	if config.LockPrefix == "" {
		config.LockPrefix = "task:lock:"
	}
	
	if config.LockTTL == 0 {
		config.LockTTL = 5 * time.Minute
	}
	
	if config.Location == nil {
		config.Location = time.Local
	}
	
	// 创建 cron 调度器（支持秒级）
	cronOptions := []cron.Option{
		cron.WithLocation(config.Location),
		cron.WithSeconds(), // 支持秒级 cron
	}
	
	return &Scheduler{
		cron:       cron.New(cronOptions...),
		redis:      config.Redis,
		tasks:      make(map[string]Task),
		lockPrefix: config.LockPrefix,
		lockTTL:    config.LockTTL,
	}
}

// Register 注册任务
func (s *Scheduler) Register(task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	name := task.Name()
	if _, exists := s.tasks[name]; exists {
		return fmt.Errorf("task %s already registered", name)
	}
	
	// 添加到 cron
	_, err := s.cron.AddFunc(task.Spec(), func() {
		s.runTask(task)
	})
	if err != nil {
		return fmt.Errorf("failed to add task %s: %w", name, err)
	}
	
	s.tasks[name] = task
	slog.Info("Task registered", "name", name, "spec", task.Spec())
	
	return nil
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.cron.Start()
	slog.Info("Task scheduler started", "tasks", len(s.tasks))
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	slog.Info("Task scheduler stopped")
}

// runTask 执行任务（带分布式锁）
func (s *Scheduler) runTask(task Task) {
	ctx := context.Background()
	name := task.Name()
	
	// 1. 尝试获取分布式锁
	lockKey := s.lockPrefix + name
	locked, err := s.acquireLock(ctx, lockKey, task.Timeout())
	if err != nil {
		slog.Error("Failed to acquire lock", "task", name, "error", err)
		return
	}
	
	if !locked {
		slog.Debug("Task already running on another instance", "task", name)
		return
	}
	
	// 2. 确保释放锁
	defer func() {
		if err := s.releaseLock(ctx, lockKey); err != nil {
			slog.Error("Failed to release lock", "task", name, "error", err)
		}
	}()
	
	// 3. 执行任务（带超时）
	taskCtx, cancel := context.WithTimeout(ctx, task.Timeout())
	defer cancel()
	
	start := time.Now()
	slog.Info("Task started", "task", name)
	
	if err := task.Run(taskCtx); err != nil {
		slog.Error("Task failed",
			"task", name,
			"error", err,
			"duration", time.Since(start),
		)
		return
	}
	
	slog.Info("Task completed",
		"task", name,
		"duration", time.Since(start),
	)
}

// acquireLock 获取分布式锁（基于 Redis）
func (s *Scheduler) acquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if ttl == 0 {
		ttl = s.lockTTL
	}
	
	// 使用 SET NX EX 原子操作
	result, err := s.redis.SetNX(ctx, key, time.Now().Unix(), ttl).Result()
	if err != nil {
		return false, err
	}
	
	return result, nil
}

// releaseLock 释放分布式锁
func (s *Scheduler) releaseLock(ctx context.Context, key string) error {
	return s.redis.Del(ctx, key).Err()
}

// ListTasks 列出所有已注册的任务
func (s *Scheduler) ListTasks() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	names := make([]string, 0, len(s.tasks))
	for name := range s.tasks {
		names = append(names, name)
	}
	
	return names
}

// GetTask 获取任务
func (s *Scheduler) GetTask(name string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	task, exists := s.tasks[name]
	return task, exists
}
