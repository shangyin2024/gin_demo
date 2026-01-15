package main

import (
	"fmt"
	"os"

	"gin_demo/internal/config"
	"gin_demo/internal/wire"
)

// @title           Gin Demo API
// @version         2.2
// @description     企业级 RESTful API 骨架项目
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化应用（通过 Wire 依赖注入）
	app, err := wire.InitApp(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize app: %v\n", err)
		os.Exit(1)
	}

	// 3. 初始化日志系统
	if err := app.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize: %v\n", err)
		os.Exit(1)
	}

	// 4. 配置应用（中间件 + 路由）
	app.Setup()

	// 5. 启动应用（HTTP 服务器 + 定时任务）
	if err := app.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start app: %v\n", err)
		os.Exit(1)
	}

	// 6. 等待关闭信号并优雅关闭
	app.Server.WaitForShutdown()
	app.Shutdown()
}
