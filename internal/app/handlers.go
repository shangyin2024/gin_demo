package app

import (
	"gin_demo/internal/app/handler/health"
	"gin_demo/internal/app/handler/user"
	"gin_demo/internal/app/middleware"
)

// Handlers 所有 HTTP 处理器
type Handlers struct {
	User   *user.Handler
	Health *health.Handler
	Auth   *middleware.AuthMiddleware
}

// NewHandlers 创建处理器集合
func NewHandlers(
	userHandler *user.Handler,
	healthHandler *health.Handler,
	authMiddleware *middleware.AuthMiddleware,
) *Handlers {
	return &Handlers{
		User:   userHandler,
		Health: healthHandler,
		Auth:   authMiddleware,
	}
}
