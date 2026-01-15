# 代码重构 v3.0 - 项目结构优化

> 日期: 2026-01-13  
> 版本: v3.0

---

## 🎯 重构目标

将 `main.go` 从 203 行精简到 65 行，通过合理的代码组织提升项目的可维护性和可读性。

---

## ✅ 已完成的重构

### 1️⃣ **main.go 精简**

**重构前**: 203 行，包含大量初始化和配置逻辑

**重构后**: 65 行，只保留核心启动流程

```go
func main() {
    // 1. 加载配置
    cfg, err := config.Load()
    
    // 2. 初始化应用（通过 Wire）
    app, err := wire.InitApp(cfg)
    
    // 3. 初始化日志
    app.Initialize()
    
    // 4. 配置应用（中间件 + 路由）
    app.Setup()
    
    // 5. 启动应用
    app.Start()
    
    // 6. 优雅关闭
    app.Server.WaitForShutdown()
    app.Shutdown()
}
```

**改进**: 代码行数减少 68%，逻辑更清晰 ✅

---

### 2️⃣ **新增 Application 层**

创建 `internal/app/` 目录，统一管理应用级逻辑：

```
internal/app/
├── app.go          # 应用程序主类
├── server.go       # HTTP 服务器
├── routes.go       # 路由配置
└── handlers.go     # 处理器集合
```

#### app.go - 应用主类

```go
type Application struct {
    Config      *config.Config
    Server      *Server
    DB          *sql.DB
    Redis       *redis.Client
    TaskManager TaskManager
    Handlers    *Handlers
}

// 核心方法
func (app *Application) Initialize() error  // 初始化
func (app *Application) Setup()              // 配置
func (app *Application) Start() error        // 启动
func (app *Application) Shutdown()           // 关闭
```

**职责**: 统一管理应用生命周期

---

#### server.go - HTTP 服务器

```go
type Server struct {
    engine   *gin.Engine
    config   *config.Config
    srv      *http.Server
    handlers *Handlers
}

// 核心方法
func (s *Server) SetupMiddlewares()      // 配置中间件
func (s *Server) SetupRoutes()           // 配置路由
func (s *Server) Start() error           // 启动服务
func (s *Server) Shutdown() error        // 关闭服务
func (s *Server) WaitForShutdown()       // 等待信号
```

**改进**:
- ✅ 中间件配置集中管理
- ✅ 配置参数化（不再硬编码）
- ✅ 职责单一

---

#### routes.go - 路由配置（NEW）

```go
// 路由层次化组织
func setupRoutes(engine *gin.Engine, handlers *Handlers)

func setupSystemRoutes(engine *gin.Engine, handlers *Handlers)
    - /metrics          (Prometheus)
    - /health           (健康检查)
    - /health/ready     (Readiness)
    - /health/live      (Liveness)

func setupAPIRoutes(engine *gin.Engine, handlers *Handlers)
    - /api/v1/*

func setupUserRoutes(rg *gin.RouterGroup, handlers *Handlers)
    - POST   /users/register
    - POST   /users/login
    - GET    /users
    - GET    /users/:id
    - PUT    /users/:id
    - PUT    /users/:id/password
    - DELETE /users/:id
```

**改进**:
- ✅ 路由独立文件，便于扩展
- ✅ 层次化组织（系统路由 / API 路由）
- ✅ 模块化设计（用户 / 文章 / 评论...）

---

#### handlers.go - 处理器集合

```go
type Handlers struct {
    User   *user.Handler
    Health *health.Handler
    Auth   *middleware.AuthMiddleware
}

func NewHandlers(...) *Handlers
```

**职责**: 统一管理所有 HTTP 处理器

---

### 3️⃣ **Wire 依赖注入优化**

#### 新增 AppSet

```go
// internal/wire/app.go
var AppSet = wire.NewSet(
    app.NewHandlers,
    app.New,
)
```

#### 简化 InitApp

```go
func InitApp(cfg *config.Config) (*app.Application, error) {
    wire.Build(
        InfrastructureSet,  // 基础设施
        RepositorySet,      // 数据访问
        ServiceSet,         // 业务逻辑
        HandlerSet,         // HTTP 处理
        TaskSet,            // 定时任务
        AppSet,             // 应用层 ⭐ NEW
    )
    return nil, nil
}
```

**改进**: 依赖注入更加清晰和模块化 ✅

---

## 📊 重构成果对比

### 代码行数变化

| 文件 | 重构前 | 重构后 | 变化 |
|------|--------|--------|------|
| **main.go** | 203 行 | 65 行 | -68% ⭐ |
| app.go | - | 128 行 | +128 (新增) |
| server.go | - | 162 行 | +162 (新增) |
| routes.go | - | 61 行 | +61 (新增) |
| handlers.go | - | 28 行 | +28 (新增) |

### 文件结构变化

**重构前**:
```
.
├── main.go (203 行，功能混杂)
└── internal/
    └── wire/ (依赖注入)
```

**重构后**:
```
.
├── main.go (65 行，只保留启动逻辑)
└── internal/
    ├── app/          ⭐ NEW
    │   ├── app.go       (应用主类)
    │   ├── server.go    (HTTP 服务器)
    │   ├── routes.go    (路由配置)
    │   └── handlers.go  (处理器集合)
    └── wire/
        └── app.go       (App 层注入)
```

---

## 🎨 设计模式应用

### 1. 分层架构（Layered Architecture）

```
┌─────────────────────────────────────┐
│        main.go (Entry Point)        │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      internal/app (Application)     │
│  - app.go      (生命周期管理)       │
│  - server.go   (HTTP 服务器)        │
│  - routes.go   (路由配置)           │
│  - handlers.go (处理器集合)         │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│    Handler → Service → Repository   │
└─────────────────────────────────────┘
```

### 2. 依赖注入（Dependency Injection）

通过 Wire 自动注入所有依赖，避免硬编码。

### 3. 单一职责原则（SRP）

- `main.go` - 只负责启动
- `app.go` - 只负责应用生命周期
- `server.go` - 只负责 HTTP 服务
- `routes.go` - 只负责路由配置

---

## 🚀 扩展性提升

### 添加新路由模块

现在只需要 3 步：

```go
// 1. 在 routes.go 中添加设置函数
func setupArticleRoutes(rg *gin.RouterGroup, handlers *Handlers) {
    articles := rg.Group("/articles")
    {
        articles.GET("", handlers.Article.List)
        articles.POST("", handlers.Article.Create)
        // ...
    }
}

// 2. 在 setupAPIRoutes 中调用
func setupAPIRoutes(engine *gin.Engine, handlers *Handlers) {
    v1 := engine.Group("/api/v1")
    {
        setupUserRoutes(v1, handlers)
        setupArticleRoutes(v1, handlers)  // ⭐ 添加这一行
    }
}

// 3. 在 handlers.go 中添加处理器
type Handlers struct {
    User    *user.Handler
    Article *article.Handler  // ⭐ 添加这一行
    // ...
}
```

---

## 📈 代码质量提升

### 可读性

- ✅ main.go 精简，启动流程一目了然
- ✅ 职责分离，每个文件职责单一
- ✅ 层次清晰，便于理解

### 可维护性

- ✅ 模块化设计，修改影响范围小
- ✅ 配置集中管理，便于调整
- ✅ 路由独立文件，便于扩展

### 可测试性

- ✅ 依赖注入，便于 Mock
- ✅ 接口清晰，便于单元测试
- ✅ 职责单一，便于测试隔离

---

## 🎓 最佳实践

### 1. 启动流程标准化

```go
func main() {
    // 1. 加载配置
    cfg := loadConfig()
    
    // 2. 初始化应用
    app := initApp(cfg)
    
    // 3. 配置应用
    app.Setup()
    
    // 4. 启动应用
    app.Start()
    
    // 5. 优雅关闭
    app.WaitForShutdown()
    app.Shutdown()
}
```

### 2. 错误处理统一化

```go
if err := app.Initialize(); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to initialize: %v\n", err)
    os.Exit(1)
}
```

### 3. 日志输出结构化

```go
slog.Info("Server starting", 
    "address", addr, 
    "mode", cfg.Server.Mode)
```

---

## 🔄 后续优化建议

### 短期

1. ✅ 添加更多路由模块（文章、评论等）
2. ✅ 完善单元测试覆盖
3. ✅ 添加集成测试

### 中期

1. 考虑添加 Graceful Restart
2. 添加配置热重载
3. 优化启动性能

### 长期

1. 考虑微服务拆分
2. 添加服务发现
3. 实现配置中心

---

## 📚 相关文档

- [项目架构](./ARCHITECTURE.md)
- [Wire 使用指南](./DEPENDENCY_INJECTION.md)
- [路由设计规范](./ROUTING.md)

---

**main.go 已精简 68%，项目结构更加清晰！** 🎉
