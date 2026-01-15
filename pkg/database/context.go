package database

import (
	"context"
	"time"
)

// DefaultQueryTimeout 默认查询超时时间
const DefaultQueryTimeout = 5 * time.Second

// WithTimeout 为 context 添加超时
func WithTimeout(ctx context.Context, timeout ...time.Duration) (context.Context, context.CancelFunc) {
	t := DefaultQueryTimeout
	if len(timeout) > 0 && timeout[0] > 0 {
		t = timeout[0]
	}
	return context.WithTimeout(ctx, t)
}

// WithQueryTimeout 为数据库查询添加默认超时
func WithQueryTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, DefaultQueryTimeout)
}
