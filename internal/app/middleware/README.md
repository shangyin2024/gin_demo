# Middleware 中间件文档

本项目的所有中间件都采用**统一的结构体封装风格**，便于依赖注入和测试。

## 📁 文件结构

```
middleware/
├── auth.go              # JWT 基础认证（兼容性函数）
├── auth_middleware.go   # JWT 认证中间件（推荐）
├── rbac.go              # RBAC 权限控制中间件
├── logger.go            # 日志中间件
├── recovery.go          # 错误恢复中间件
├── ratelimit.go         # 限流中间件
├── metrics.go           # Prometheus 指标中间件
├── compress.go          # Gzip 压缩中间件
├── security.go          # HTTP 安全头中间件
└── README.md            # 本文档
```

## 🎨 中间件风格规范

### 1. 标准结构体风格（推荐）

所有需要依赖注入的中间件都采用结构体封装：

```go
// 定义中间件结构体
type MyMiddleware struct {
    dependency1 *SomeDependency
    dependency2 *AnotherDependency
}

// 构造函数（用于 Wire 依赖注入）
func NewMyMiddleware(dep1 *SomeDependency, dep2 *AnotherDependency) *MyMiddleware {
    return &MyMiddleware{
        dependency1: dep1,
        dependency2: dep2,
    }
}

// Handle 方法返回 gin.HandlerFunc
func (m *MyMiddleware) Handle() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
        c.Next()
    }
}
```

**使用示例**：
```go
// Wire 注入
authMiddleware := middleware.NewAuthMiddleware(jwtManager)

// 路由中使用
router.Use(authMiddleware.Handle())
```

### 2. 函数式风格（无依赖注入场景）

对于不需要外部依赖的简单中间件，可以使用函数式风格：

```go
// Logger 日志中间件（无需依赖注入）
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
        c.Next()
    }
}
```

**使用示例**：
```go
router.Use(middleware.Logger())
router.Use(middleware.Recovery())
```

### 3. 配置式风格（需要配置参数）

对于需要配置的中间件，提供配置结构体：

```go
type SecurityConfig struct {
    EnableHSTS bool
    // ... 其他配置
}

func Security(config ...SecurityConfig) gin.HandlerFunc {
    cfg := DefaultSecurityConfig
    if len(config) > 0 {
        cfg = config[0]
    }
    
    return func(c *gin.Context) {
        // 使用配置
        c.Next()
    }
}
```

---

## 📚 中间件使用指南

### 全局中间件（推荐顺序）

在 `server.go` 中的 `SetupMiddlewares` 方法中配置：

```go
func (s *Server) SetupMiddlewares() {
    s.engine.Use(
        middleware.Recovery(),              // 1. 错误恢复（最先）
        middleware.Metrics(),               // 2. Prometheus 指标
        s.configureSecurityMiddleware(),    // 3. HTTP 安全头
        s.configureCompressionMiddleware(), // 4. Gzip 压缩
        s.configureCORS(),                  // 5. CORS
        s.configureRequestID(),             // 6. Request ID
        middleware.Logger(),                // 7. 日志
        middleware.RateLimit(...),          // 8. 限流（最后）
    )
}
```

**顺序说明**：
1. **Recovery** 必须最先，捕获所有 panic
2. **Metrics** 尽早记录，包含所有后续中间件的耗时
3. **Security/Compression/CORS** 在业务逻辑前处理
4. **RequestID** 为后续日志提供追踪标识
5. **Logger** 记录完整的请求信息
6. **RateLimit** 最后，避免记录被拒绝的请求

### 路由级中间件

#### 1. 基础认证

```go
// 需要登录
protected := router.Group("/api")
protected.Use(handlers.Auth.Handle())
{
    protected.GET("/profile", handler.GetProfile)
}
```

#### 2. 角色权限控制

```go
// 需要管理员角色
admin := router.Group("/admin")
admin.Use(handlers.Auth.Handle())                                    // 先认证
admin.Use(middleware.RequireRole(auth.RoleAdmin, auth.RoleSuperAdmin))  // 再授权
{
    admin.GET("/users", handler.ListUsers)
}
```

#### 3. 细粒度权限控制

```go
// 需要特定权限
sensitive := router.Group("/sensitive")
sensitive.Use(handlers.Auth.Handle())
sensitive.Use(middleware.RequirePermission(
    auth.PermissionUserWrite,
    auth.PermissionUserDelete,
))
{
    sensitive.POST("/batch", handler.BatchOperation)
}
```

#### 4. 可选认证

```go
// Token 可选（有则验证，无则放行）
public := router.Group("/public")
public.Use(authMiddleware.HandleOptional())
{
    public.GET("/articles", handler.ListArticles)  // 登录用户可看更多
}
```

---

## 🔧 自定义中间件开发

### 模板

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "log/slog"
)

// MyMiddleware 自定义中间件
type MyMiddleware struct {
    config MyConfig
}

// MyConfig 中间件配置
type MyConfig struct {
    Enabled bool
    Value   string
}

// NewMyMiddleware 创建中间件实例（用于 Wire 依赖注入）
func NewMyMiddleware(config MyConfig) *MyMiddleware {
    return &MyMiddleware{
        config: config,
    }
}

// Handle 返回 gin.HandlerFunc
func (m *MyMiddleware) Handle() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        if !m.config.Enabled {
            c.Next()
            return
        }

        // 业务逻辑
        slog.InfoContext(c.Request.Context(), "Middleware executed")

        // 调用下一个中间件
        c.Next()

        // 后置处理（可选）
        slog.DebugContext(c.Request.Context(), "Middleware completed")
    }
}
```

### Wire 集成

```go
// wire/handler.go
var HandlerSet = wire.NewSet(
    // ... 其他 Provider
    middleware.NewMyMiddleware,
    // ...
)

// wire/infrastructure.go
func provideMyMiddlewareConfig(cfg *config.Config) middleware.MyConfig {
    return middleware.MyConfig{
        Enabled: cfg.MyFeature.Enabled,
        Value:   cfg.MyFeature.Value,
    }
}
```

---

## 🧪 中间件测试

### 单元测试示例

```go
func TestAuthMiddleware_Handle(t *testing.T) {
    // 1. 准备测试数据
    jwtManager := auth.NewDefaultJWTManager("test-secret", 1*time.Hour)
    token, _ := jwtManager.GenerateToken(123)

    // 2. 创建中间件
    authMiddleware := middleware.NewAuthMiddleware(jwtManager)

    // 3. 创建测试上下文
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("GET", "/test", nil)
    c.Request.Header.Set("Authorization", "Bearer "+token)

    // 4. 执行中间件
    handler := authMiddleware.Handle()
    handler(c)

    // 5. 断言
    assert.False(t, c.IsAborted())
    userID := middleware.GetUserID(c)
    assert.Equal(t, int64(123), userID)
}
```

---

## 📊 中间件对比

| 中间件 | 风格 | 依赖注入 | 配置 | 用途 |
|--------|------|----------|------|------|
| AuthMiddleware | 结构体 | ✅ | ❌ | JWT 认证 |
| RBACMiddleware | 结构体 | ✅ | ❌ | 角色权限控制 |
| Logger | 函数式 | ❌ | ❌ | 请求日志 |
| Recovery | 函数式 | ❌ | ❌ | Panic 恢复 |
| RateLimiter | 结构体 | ❌ | ✅ | 限流 |
| Security | 配置式 | ❌ | ✅ | 安全头 |
| Compress | 配置式 | ❌ | ✅ | Gzip 压缩 |
| Metrics | 函数式 | ❌ | ❌ | Prometheus |

---

## 💡 最佳实践

### 1. 中间件命名

- 结构体命名：`XxxMiddleware`（如 `AuthMiddleware`）
- 函数命名：`Xxx`（如 `Logger`）
- Handle 方法：统一使用 `Handle()`

### 2. 错误处理

```go
// ✅ 推荐：使用统一的响应格式
response.Error(c, response.ErrUnauthorized)
c.Abort()

// ❌ 避免：直接使用 c.JSON
c.JSON(401, gin.H{"error": "unauthorized"})
c.Abort()
```

### 3. Context 传递

```go
// ✅ 推荐：使用有意义的键名常量
const UserIDKey = "user_id"
c.Set(UserIDKey, userID)

// ❌ 避免：硬编码字符串
c.Set("uid", userID)
```

### 4. 日志记录

```go
// ✅ 推荐：使用结构化日志
slog.InfoContext(c.Request.Context(), "User authenticated",
    "user_id", userID,
    "role", role,
)

// ❌ 避免：使用 fmt.Println
fmt.Println("User", userID, "authenticated")
```

---

## 🔗 相关文档

- [RBAC 权限控制](../../../docs/RBAC.md)
- [JWT 认证](../../../docs/JWT.md)
- [架构设计](../../../docs/ARCHITECTURE.md)
