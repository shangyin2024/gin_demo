package health

import (
	"context"
	"time"
)

// Status 健康状态
type Status string

const (
	// StatusOK 正常
	StatusOK Status = "ok"
	// StatusDegraded 降级（部分服务不可用）
	StatusDegraded Status = "degraded"
	// StatusError 错误（核心服务不可用）
	StatusError Status = "error"
)

// HealthStatus 健康检查响应
type HealthStatus struct {
	Status    Status           `json:"status"`
	Timestamp int64            `json:"timestamp"`
	Version   string           `json:"version"`
	Checks    map[string]Check `json:"checks"`
}

// Check 单个检查项结果
type Check struct {
	Status   Status `json:"status"`
	Message  string `json:"message,omitempty"`
	Duration string `json:"duration"`
}

// Checker 健康检查器接口（通用）
type Checker interface {
	// Check 执行健康检查
	Check(ctx context.Context) HealthStatus

	// IsHealthy 简单的健康检查（只返回 true/false）
	IsHealthy(ctx context.Context) bool
}

// ComponentChecker 单个组件检查器接口
type ComponentChecker interface {
	// Name 返回组件名称（如 "database", "redis"）
	Name() string

	// Check 执行检查
	Check(ctx context.Context) Check

	// IsCritical 是否是关键组件（如果失败，整个服务不可用）
	IsCritical() bool
}

// MultiChecker 通用的多组件健康检查器
type MultiChecker struct {
	version  string
	checkers []ComponentChecker
}

// NewMultiChecker 创建多组件健康检查器
func NewMultiChecker(version string, checkers ...ComponentChecker) *MultiChecker {
	return &MultiChecker{
		version:  version,
		checkers: checkers,
	}
}

// Check 执行所有组件的健康检查
func (m *MultiChecker) Check(ctx context.Context) HealthStatus {
	status := HealthStatus{
		Status:    StatusOK,
		Timestamp: time.Now().Unix(),
		Version:   m.version,
		Checks:    make(map[string]Check),
	}

	for _, checker := range m.checkers {
		check := checker.Check(ctx)
		status.Checks[checker.Name()] = check

		// 根据检查结果更新整体状态
		if check.Status == StatusError {
			if checker.IsCritical() {
				status.Status = StatusError
			} else if status.Status == StatusOK {
				status.Status = StatusDegraded
			}
		}
	}

	return status
}

// IsHealthy 简单的健康检查
func (m *MultiChecker) IsHealthy(ctx context.Context) bool {
	status := m.Check(ctx)
	return status.Status != StatusError
}

// AddChecker 添加检查器
func (m *MultiChecker) AddChecker(checker ComponentChecker) {
	m.checkers = append(m.checkers, checker)
}
