package app

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// setupRoutes 配置所有路由
func setupRoutes(engine *gin.Engine, handlers *Handlers) {
	// 系统路由（无需认证）
	setupSystemRoutes(engine, handlers)

	// API 路由
	setupAPIRoutes(engine, handlers)
}

// setupSystemRoutes 配置系统路由
func setupSystemRoutes(engine *gin.Engine, handlers *Handlers) {
	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 健康检查
	health := engine.Group("/health")
	{
		health.GET("", handlers.Health.Health)      // 完整健康检查
		health.GET("/ready", handlers.Health.Ready) // Readiness Probe
		health.GET("/live", handlers.Health.Live)   // Liveness Probe
	}
}

// setupAPIRoutes 配置 API 路由
func setupAPIRoutes(engine *gin.Engine, handlers *Handlers) {
	// API v1
	v1 := engine.Group("/api/v1")
	{
		setupAPIv1Routes(v1, handlers)
	}

	// 未来可以添加 v2
	// v2 := engine.Group("/api/v2")
	// {
	//     setupAPIv2Routes(v2, handlers)
	// }
}
