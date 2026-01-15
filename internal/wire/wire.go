//go:build wireinject
// +build wireinject

package wire

import (
	"gin_demo/internal/app"
	"gin_demo/internal/config"

	"github.com/google/wire"
)

// InitApp 初始化应用（Wire 会自动生成实现）
func InitApp(cfg *config.Config) (*app.Application, error) {
	wire.Build(
		// 基础设施层（数据库、Redis、缓存、JWT）
		InfrastructureSet,

		// Repository 层（数据访问）
		RepositorySet,

		// Service 层（业务逻辑）
		ServiceSet,

		// Handler 层（HTTP 处理）
		HandlerSet,

		// Task 层（定时任务）
		TaskSet,

		// App 层
		AppSet,
	)
	return nil, nil
}
