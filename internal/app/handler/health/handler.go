package health

import (
	"gin_demo/pkg/health"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler 健康检查处理器
type Handler struct {
	checker health.Checker
	// 缓存机制，防止频繁检查
	cache    *healthCache
	cacheTTL time.Duration
}

// healthCache 健康检查缓存
type healthCache struct {
	mu       sync.RWMutex
	status   health.HealthStatus
	cachedAt time.Time
}

// NewHandler 创建健康检查处理器
func NewHandler(checker health.Checker) *Handler {
	return &Handler{
		checker:  checker,
		cache:    &healthCache{},
		cacheTTL: 5 * time.Second, // 缓存 5 秒，防止频繁检查
	}
}

// Health 健康检查（带缓存）
//
// @Summary 健康检查
// @Description 检查应用及其依赖服务的健康状态
// @Tags 系统监控
// @Accept json
// @Produce json
// @Success 200 {object} health.HealthStatus "服务健康"
// @Failure 503 {object} health.HealthStatus "服务不可用"
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	status := h.getCachedStatus(c)

	// 根据健康状态设置 HTTP 状态码
	httpStatus := http.StatusOK
	if status.Status == health.StatusError {
		httpStatus = http.StatusServiceUnavailable
	} else if status.Status == health.StatusDegraded {
		httpStatus = http.StatusOK // 降级状态仍返回 200
	}

	c.JSON(httpStatus, status)
}

// Ready 就绪检查（用于 Kubernetes readiness probe）
//
// @Summary 就绪检查
// @Description 检查服务是否准备好接收流量（用于 Kubernetes Readiness Probe）
// @Tags 系统监控
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "服务就绪"
// @Failure 503 {object} map[string]string "服务未就绪"
// @Router /health/ready [get]
func (h *Handler) Ready(c *gin.Context) {
	if h.checker.IsHealthy(c.Request.Context()) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
		})
	}
}

// Live 存活检查（用于 Kubernetes liveness probe）
//
// @Summary 存活检查
// @Description 检查服务是否存活（用于 Kubernetes Liveness Probe）
// @Tags 系统监控
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "服务存活"
// @Router /health/live [get]
func (h *Handler) Live(c *gin.Context) {
	// 存活检查只检查应用本身是否运行，不检查依赖服务
	// 这样可以避免因依赖服务故障导致容器被重启
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// getCachedStatus 获取缓存的健康状态（防止频繁检查）
func (h *Handler) getCachedStatus(c *gin.Context) health.HealthStatus {
	h.cache.mu.RLock()
	if time.Since(h.cache.cachedAt) < h.cacheTTL {
		status := h.cache.status
		h.cache.mu.RUnlock()
		return status
	}
	h.cache.mu.RUnlock()

	// 缓存过期，重新检查
	h.cache.mu.Lock()
	defer h.cache.mu.Unlock()

	// Double-check
	if time.Since(h.cache.cachedAt) < h.cacheTTL {
		return h.cache.status
	}

	// 执行健康检查
	status := h.checker.Check(c.Request.Context())
	h.cache.status = status
	h.cache.cachedAt = time.Now()

	return status
}
