package health

import (
	"context"
	"database/sql"
	"gin_demo/pkg/health"
	"time"
)

// DatabaseChecker 数据库健康检查器
type DatabaseChecker struct {
	db *sql.DB
}

// NewDatabaseChecker 创建数据库检查器
func NewDatabaseChecker(db *sql.DB) *DatabaseChecker {
	return &DatabaseChecker{db: db}
}

// Name 返回组件名称
func (d *DatabaseChecker) Name() string {
	return "database"
}

// Check 执行数据库检查
func (d *DatabaseChecker) Check(ctx context.Context) health.Check {
	start := time.Now()
	check := health.Check{
		Status: health.StatusOK,
	}

	// Ping 数据库
	if err := d.db.PingContext(ctx); err != nil {
		check.Status = health.StatusError
		check.Message = err.Error()
		check.Duration = time.Since(start).String()
		return check
	}

	// 获取数据库统计信息
	stats := d.db.Stats()
	check.Duration = time.Since(start).String()

	// 检查连接池是否耗尽
	if stats.MaxOpenConnections > 0 && stats.OpenConnections >= stats.MaxOpenConnections {
		check.Status = health.StatusDegraded
		check.Message = "connection pool exhausted"
	}

	return check
}

// IsCritical 数据库是关键组件
func (d *DatabaseChecker) IsCritical() bool {
	return true
}
