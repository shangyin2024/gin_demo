package app

import (
	"database/sql"
	"log/slog"

	"gin_demo/internal/config"
	"gin_demo/pkg/logger"

	"github.com/redis/go-redis/v9"
)

// Application 应用程序
type Application struct {
	Config      *config.Config
	Server      *Server
	DB          *sql.DB
	Redis       redis.UniversalClient
	TaskManager TaskManager
	Handlers    *Handlers // HTTP 处理器
}

// TaskManager 任务管理器接口
type TaskManager interface {
	Start()
	Stop()
	ListTasks() []string
}

// New 创建应用程序实例
func New(
	cfg *config.Config,
	db *sql.DB,
	redis redis.UniversalClient,
	handlers *Handlers,
	taskManager TaskManager,
) *Application {
	server := NewServer(cfg)
	server.handlers = handlers // 注入 handlers

	return &Application{
		Config:      cfg,
		Server:      server,
		DB:          db,
		Redis:       redis,
		TaskManager: taskManager,
		Handlers:    handlers,
	}
}

// Initialize 初始化应用程序
func (app *Application) Initialize() error {
	// 初始化日志
	logger.Setup(logger.Config{
		Level:        parseLogLevel(app.Config.Logger.Level),
		IsJSON:       app.Config.Logger.IsJSON,
		AddSource:    app.Config.Logger.AddSource,
		RequestIDKey: app.Config.Logger.RequestIDKey,
	})
	slog.Info("Logger initialized")

	return nil
}

// Setup 配置应用程序
func (app *Application) Setup() {
	// 配置中间件
	app.Server.SetupMiddlewares()

	// 配置路由
	app.Server.SetupRoutes()

	slog.Info("Application setup completed")
}

// Start 启动应用程序
func (app *Application) Start() error {
	// 启动定时任务
	app.TaskManager.Start()
	slog.Info("Task scheduler started", "tasks", app.TaskManager.ListTasks())

	// 启动 HTTP 服务器
	if err := app.Server.Start(); err != nil {
		return err
	}

	return nil
}

// Shutdown 关闭应用程序
func (app *Application) Shutdown() {
	// 停止定时任务
	app.TaskManager.Stop()

	// 关闭 HTTP 服务器
	app.Server.Shutdown()

	// 清理资源
	app.Cleanup()
}

// Cleanup 清理资源
func (app *Application) Cleanup() {
	if app.DB != nil {
		if err := app.DB.Close(); err != nil {
			slog.Error("Failed to close database", "error", err)
		} else {
			slog.Debug("Database connection closed")
		}
	}
	
	if app.Redis != nil {
		if err := app.Redis.Close(); err != nil {
			slog.Error("Failed to close redis", "error", err)
		} else {
			slog.Debug("Redis connection closed")
		}
	}
	
	slog.Info("Resources cleaned up")
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
