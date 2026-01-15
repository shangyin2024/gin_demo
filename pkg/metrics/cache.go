package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ============================================================================
// 缓存指标
// ============================================================================

var (
	// 缓存操作计数
	CacheOperations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_operations_total",
		Help: "Total number of cache operations",
	}, []string{"operation", "entity"}) // operation: get, set, delete; entity: user, article, etc.

	// 缓存命中率
	CacheHits = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "Total number of cache hits",
	}, []string{"entity"})

	// 缓存未命中
	CacheMisses = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "Total number of cache misses",
	}, []string{"entity"})

	// 缓存操作延迟
	CacheLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "cache_operation_duration_seconds",
		Help:    "Cache operation latency in seconds",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1}, // 1ms 到 1s
	}, []string{"operation", "entity"})

	// 缓存大小（字节）
	CacheSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cache_size_bytes",
		Help: "Current size of cache in bytes",
	}, []string{"entity"})

	// 缓存条目数
	CacheEntries = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cache_entries_total",
		Help: "Current number of entries in cache",
	}, []string{"entity"})

	// 缓存驱逐计数
	CacheEvictions = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_evictions_total",
		Help: "Total number of cache evictions",
	}, []string{"entity", "reason"}) // reason: ttl_expired, memory_pressure, manual
)

// ============================================================================
// 缓存错误指标
// ============================================================================

var (
	// 缓存错误
	CacheErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_errors_total",
		Help: "Total number of cache errors",
	}, []string{"operation", "error_type"}) // error_type: connection_error, timeout, serialization_error
)

// ============================================================================
// 辅助函数
// ============================================================================

// RecordCacheHit 记录缓存命中
func RecordCacheHit(entity string) {
	CacheHits.WithLabelValues(entity).Inc()
	CacheOperations.WithLabelValues("get", entity).Inc()
}

// RecordCacheMiss 记录缓存未命中
func RecordCacheMiss(entity string) {
	CacheMisses.WithLabelValues(entity).Inc()
	CacheOperations.WithLabelValues("get", entity).Inc()
}

// RecordCacheSet 记录缓存设置
func RecordCacheSet(entity string) {
	CacheOperations.WithLabelValues("set", entity).Inc()
}

// RecordCacheDelete 记录缓存删除
func RecordCacheDelete(entity string) {
	CacheOperations.WithLabelValues("delete", entity).Inc()
}

// RecordCacheError 记录缓存错误
func RecordCacheError(operation, errorType string) {
	CacheErrors.WithLabelValues(operation, errorType).Inc()
}

// RecordCacheEviction 记录缓存驱逐
func RecordCacheEviction(entity, reason string) {
	CacheEvictions.WithLabelValues(entity, reason).Inc()
}

// UpdateCacheSize 更新缓存大小
func UpdateCacheSize(entity string, sizeBytes float64) {
	CacheSize.WithLabelValues(entity).Set(sizeBytes)
}

// UpdateCacheEntries 更新缓存条目数
func UpdateCacheEntries(entity string, count float64) {
	CacheEntries.WithLabelValues(entity).Set(count)
}

// GetCacheHitRate 计算缓存命中率（用于展示，非指标）
// 实际使用时应该通过 PromQL 计算：rate(cache_hits_total[5m]) / rate(cache_operations_total{operation="get"}[5m])
func GetCacheHitRate() string {
	return "Use PromQL: rate(cache_hits_total[5m]) / rate(cache_operations_total{operation=\"get\"}[5m])"
}
