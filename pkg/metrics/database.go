package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ============================================================================
// 数据库指标
// ============================================================================

var (
	// 数据库查询延迟
	DBQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "db_query_duration_seconds",
		Help:    "Database query latency in seconds",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}, // 1ms 到 10s
	}, []string{"operation", "table"}) // operation: select, insert, update, delete

	// 慢查询计数
	DBSlowQueries = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "db_slow_queries_total",
		Help: "Total number of slow database queries",
	}, []string{"operation", "table", "threshold"}) // threshold: 100ms, 500ms, 1s, 5s

	// 数据库连接数
	DBConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_connections_current",
		Help: "Current number of database connections",
	}, []string{"state"}) // state: open, in_use, idle

	// 数据库错误
	DBErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "db_errors_total",
		Help: "Total number of database errors",
	}, []string{"operation", "error_type"}) // error_type: connection_error, query_error, timeout, deadlock

	// 事务指标
	DBTransactions = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "db_transactions_total",
		Help: "Total number of database transactions",
	}, []string{"status"}) // status: committed, rolled_back

	// 事务延迟
	DBTransactionDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "db_transaction_duration_seconds",
		Help:    "Database transaction latency in seconds",
		Buckets: []float64{.01, .05, .1, .25, .5, 1, 2.5, 5, 10}, // 10ms 到 10s
	})
)

// ============================================================================
// 慢查询阈值
// ============================================================================

const (
	// 慢查询阈值（毫秒）
	SlowQueryThreshold100ms  = 100
	SlowQueryThreshold500ms  = 500
	SlowQueryThreshold1s     = 1000
	SlowQueryThreshold5s     = 5000
)

// ============================================================================
// 辅助函数
// ============================================================================

// RecordDBQuery 记录数据库查询
func RecordDBQuery(operation, table string, durationSeconds float64) {
	DBQueryDuration.WithLabelValues(operation, table).Observe(durationSeconds)
	
	// 检查是否是慢查询
	durationMs := durationSeconds * 1000
	if durationMs >= SlowQueryThreshold100ms {
		DBSlowQueries.WithLabelValues(operation, table, "100ms").Inc()
	}
	if durationMs >= SlowQueryThreshold500ms {
		DBSlowQueries.WithLabelValues(operation, table, "500ms").Inc()
	}
	if durationMs >= SlowQueryThreshold1s {
		DBSlowQueries.WithLabelValues(operation, table, "1s").Inc()
	}
	if durationMs >= SlowQueryThreshold5s {
		DBSlowQueries.WithLabelValues(operation, table, "5s").Inc()
	}
}

// RecordDBError 记录数据库错误
func RecordDBError(operation, errorType string) {
	DBErrors.WithLabelValues(operation, errorType).Inc()
}

// RecordDBTransaction 记录事务
func RecordDBTransaction(committed bool, duration float64) {
	status := "committed"
	if !committed {
		status = "rolled_back"
	}
	DBTransactions.WithLabelValues(status).Inc()
	DBTransactionDuration.Observe(duration)
}

// UpdateDBConnections 更新数据库连接数
func UpdateDBConnections(openConns, inUse, idle int) {
	DBConnections.WithLabelValues("open").Set(float64(openConns))
	DBConnections.WithLabelValues("in_use").Set(float64(inUse))
	DBConnections.WithLabelValues("idle").Set(float64(idle))
}
