package task

import (
	"context"
	"time"
)

// BaseTask 任务基类（提供默认实现）
type BaseTask struct {
	name    string
	spec    string
	timeout time.Duration
	handler func(ctx context.Context) error
}

// NewBaseTask 创建基础任务
func NewBaseTask(name, spec string, timeout time.Duration, handler func(ctx context.Context) error) *BaseTask {
	return &BaseTask{
		name:    name,
		spec:    spec,
		timeout: timeout,
		handler: handler,
	}
}

// Name 实现 Task 接口
func (t *BaseTask) Name() string {
	return t.name
}

// Spec 实现 Task 接口
func (t *BaseTask) Spec() string {
	return t.spec
}

// Timeout 实现 Task 接口
func (t *BaseTask) Timeout() time.Duration {
	if t.timeout == 0 {
		return 5 * time.Minute // 默认 5 分钟
	}
	return t.timeout
}

// Run 实现 Task 接口
func (t *BaseTask) Run(ctx context.Context) error {
	if t.handler == nil {
		return nil
	}
	return t.handler(ctx)
}
