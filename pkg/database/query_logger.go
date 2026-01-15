package database

import (
	"context"
	"log/slog"
	"time"

	"gin_demo/pkg/metrics"
)

// SlowQueryThreshold 慢查询阈值（默认100ms）
var SlowQueryThreshold = 100 * time.Millisecond

// LogQuery 记录查询（带慢查询检测）
func LogQuery(ctx context.Context, operation, table, query string, duration time.Duration, err error) {
	// 记录 Prometheus 指标
	metrics.RecordDBQuery(operation, table, duration.Seconds())

	// 记录错误
	if err != nil {
		metrics.RecordDBError(operation, classifyDBError(err))
	}

	// 慢查询日志
	if duration > SlowQueryThreshold {
		slog.WarnContext(ctx, "🐌 Slow query detected",
			"operation", operation,
			"table", table,
			"duration", duration.String(),
			"threshold", SlowQueryThreshold.String(),
			"query", truncateQuery(query, 200), // 截断长查询
		)
	} else {
		// 正常查询只记录 Debug 级别
		slog.DebugContext(ctx, "Database query executed",
			"operation", operation,
			"table", table,
			"duration", duration.String(),
		)
	}
}

// LogTransaction 记录事务
func LogTransaction(ctx context.Context, committed bool, duration time.Duration, err error) {
	// 记录 Prometheus 指标
	metrics.RecordDBTransaction(committed, duration.Seconds())

	// 记录日志
	status := "committed"
	if !committed {
		status = "rolled_back"
	}

	if err != nil {
		slog.ErrorContext(ctx, "Transaction failed",
			"status", status,
			"duration", duration.String(),
			"error", err,
		)
	} else if duration > SlowQueryThreshold {
		slog.WarnContext(ctx, "Slow transaction detected",
			"status", status,
			"duration", duration.String(),
			"threshold", SlowQueryThreshold.String(),
		)
	} else {
		slog.DebugContext(ctx, "Transaction completed",
			"status", status,
			"duration", duration.String(),
		)
	}
}

// WithQueryLogging 包装查询函数，自动记录慢查询
func WithQueryLogging(ctx context.Context, operation, table string, queryFn func() error) error {
	start := time.Now()
	err := queryFn()
	duration := time.Since(start)

	LogQuery(ctx, operation, table, "", duration, err)
	return err
}

// WithTransactionLogging 包装事务函数，自动记录
func WithTransactionLogging(ctx context.Context, txFn func() error) error {
	start := time.Now()
	err := txFn()
	duration := time.Since(start)

	committed := err == nil
	LogTransaction(ctx, committed, duration, err)
	return err
}

// classifyDBError 分类数据库错误
func classifyDBError(err error) string {
	if err == nil {
		return "none"
	}

	errStr := err.Error()
	
	// 简单的错误分类（可根据实际数据库扩展）
	switch {
	case contains(errStr, "connection"):
		return "connection_error"
	case contains(errStr, "timeout"):
		return "timeout"
	case contains(errStr, "deadlock"):
		return "deadlock"
	case contains(errStr, "constraint"):
		return "constraint_violation"
	case contains(errStr, "syntax"):
		return "syntax_error"
	default:
		return "query_error"
	}
}

// truncateQuery 截断长查询（用于日志）
func truncateQuery(query string, maxLen int) string {
	if len(query) <= maxLen {
		return query
	}
	return query[:maxLen] + "..."
}

// contains 检查字符串包含（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && hasSubstring(s, substr)))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
