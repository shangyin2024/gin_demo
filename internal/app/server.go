package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin_demo/internal/app/middleware"
	"gin_demo/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// Server HTTP 服务器
type Server struct {
	engine   *gin.Engine
	config   *config.Config
	srv      *http.Server
	handlers *Handlers
}

// NewServer 创建 HTTP 服务器
func NewServer(cfg *config.Config) *Server {
	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 引擎
	engine := gin.New()

	// 限制请求体大小
	engine.MaxMultipartMemory = cfg.Server.MaxRequestBodySize

	return &Server{
		engine: engine,
		config: cfg,
	}
}

// SetupMiddlewares 配置中间件
func (s *Server) SetupMiddlewares() {
	s.engine.Use(
		middleware.Recovery(),              // 错误恢复
		middleware.Metrics(),               // Prometheus 指标
		s.configureSecurityMiddleware(),    // HTTP 安全头
		s.configureCompressionMiddleware(), // Gzip 压缩
		s.configureCORS(),                  // CORS
		s.configureRequestID(),             // Request ID
		middleware.Logger(),                // 日志
		middleware.RateLimit(middleware.NewRateLimiter(100, 200)), // 限流
	)
	
	// 开发环境启用 pprof 性能分析
	if s.config.Server.Mode == "debug" || s.config.Server.Mode == "test" {
		slog.Info("pprof enabled", "mode", s.config.Server.Mode)
		RegisterPprof(s.engine)
	}
}

// SetupRoutes 配置路由
func (s *Server) SetupRoutes() {
	if s.handlers == nil {
		slog.Error("handlers not initialized")
		os.Exit(1)
	}
	setupRoutes(s.engine, s.handlers)
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	s.srv = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}

	// 在协程中启动服务器
	go func() {
		slog.Info("Server starting", "address", addr, "mode", s.config.Server.Mode)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	return nil
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Server shutting down...")
	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		return err
	}

	slog.Info("Server exited")
	return nil
}

// WaitForShutdown 等待关闭信号
func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// configureSecurityMiddleware 配置安全中间件
func (s *Server) configureSecurityMiddleware() gin.HandlerFunc {
	if !s.config.Security.Headers.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	return middleware.Security(middleware.SecurityConfig{
		EnableHSTS:            s.config.Security.Headers.EnableHSTS,
		HSTSMaxAge:            s.config.Security.Headers.HSTSMaxAge,
		HSTSIncludeSubdomains: s.config.Security.Headers.HSTSIncludeSubdomains,
		EnableCSP:             s.config.Security.Headers.EnableCSP,
		CSPPolicy:             s.config.Security.Headers.CSPPolicy,
		EnableFrameOptions:    s.config.Security.Headers.EnableFrameOptions,
		FrameOptions:          s.config.Security.Headers.FrameOptions,
	})
}

// configureCompressionMiddleware 配置压缩中间件
func (s *Server) configureCompressionMiddleware() gin.HandlerFunc {
	if !s.config.Security.EnableCompression {
		return func(c *gin.Context) { c.Next() }
	}

	return middleware.Compress(middleware.CompressConfig{
		Level: s.config.Security.CompressionLevel,
	})
}

// configureCORS 配置 CORS
func (s *Server) configureCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: s.config.CORS.AllowedOrigins,
		// 注意：不支持 PATCH 方法，因为很多企业（如微信）的网络环境不支持
		// 统一使用 PUT 进行资源更新（完整更新和部分更新都用 PUT）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: s.config.CORS.AllowCredentials,
		MaxAge:           time.Duration(s.config.CORS.MaxAge) * time.Second,
	})
}

// configureRequestID 配置 Request ID
func (s *Server) configureRequestID() gin.HandlerFunc {
	return requestid.New(requestid.WithHandler(func(c *gin.Context, rid string) {
		ctx := context.WithValue(c.Request.Context(), "requestId", rid)
		c.Request = c.Request.WithContext(ctx)
	}))
}
