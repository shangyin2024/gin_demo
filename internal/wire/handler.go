package wire

import (
	"gin_demo/internal/app/handler/health"
	"gin_demo/internal/app/handler/user"
	"gin_demo/internal/app/middleware"

	"github.com/google/wire"
)

// HandlerSet Handler 层 Provider 集合
var HandlerSet = wire.NewSet(
	user.NewHandler,
	health.NewHandler,
	middleware.NewAuthMiddleware,
)
