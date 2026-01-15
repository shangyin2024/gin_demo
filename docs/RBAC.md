# RBAC 权限控制系统

本项目实现了完整的 RBAC（基于角色的访问控制）系统，支持角色和细粒度权限管理。

## 📚 目录

- [角色定义](#角色定义)
- [权限定义](#权限定义)
- [使用方法](#使用方法)
- [最佳实践](#最佳实践)
- [示例代码](#示例代码)

---

## 🎭 角色定义

系统预定义了 5 种角色，按权限从低到高排列：

| 角色 | 常量 | 级别 | 说明 |
|------|------|------|------|
| 游客 | `RoleGuest` | 0 | 未登录用户，只能查看公开内容 |
| 普通用户 | `RoleUser` | 40 | 已注册用户，可以操作自己的数据 |
| 版主 | `RoleModerator` | 60 | 可以审核内容，管理用户行为 |
| 管理员 | `RoleAdmin` | 80 | 可以管理用户和内容，但不能修改系统配置 |
| 超级管理员 | `RoleSuperAdmin` | 100 | 拥有所有权限，包括系统配置 |

### 权限继承规则

- **超级管理员**：拥有所有权限
- **管理员**：拥有大部分权限（除了系统配置）
- **版主**：拥有内容审核和用户查看权限
- **普通用户**：只能读写自己的内容
- **游客**：只能读取公开内容

---

## 🔐 权限定义

系统支持细粒度权限控制，权限格式为 `资源:操作`：

### 用户权限

```go
PermissionUserRead   = "user:read"    // 读取用户信息
PermissionUserWrite  = "user:write"   // 修改用户信息
PermissionUserDelete = "user:delete"  // 删除用户
```

### 内容权限

```go
PermissionContentRead   = "content:read"    // 读取内容
PermissionContentWrite  = "content:write"   // 创建/修改内容
PermissionContentDelete = "content:delete"  // 删除内容
PermissionContentAudit  = "content:audit"   // 审核内容
```

### 系统权限

```go
PermissionSystemConfig  = "system:config"   // 修改系统配置
PermissionSystemMonitor = "system:monitor"  // 查看系统监控
```

---

## 🚀 使用方法

### 1. 生成包含角色的 Token

```go
import "gin_demo/pkg/auth"

// 创建 RBAC JWT 管理器
jwtManager := auth.NewRBACJWTManager(secret, expiration)

// 生成包含角色的 Token（管理员）
token, err := jwtManager.GenerateToken(
    userID,              // 用户 ID
    auth.RoleAdmin,      // 角色
    // 可选：额外的细粒度权限
    auth.PermissionSystemMonitor,
)
```

### 2. 在路由中应用权限控制

#### 方式 A: 角色检查（推荐用于粗粒度控制）

```go
import (
    "gin_demo/internal/app/middleware"
    "gin_demo/pkg/auth"
)

// 需要管理员角色
admin := router.Group("/admin")
admin.Use(handlers.Auth.Handle())                              // 先认证
admin.Use(middleware.RequireRole(auth.RoleAdmin, auth.RoleSuperAdmin))  // 再检查角色
{
    admin.GET("/users", handler.ListUsers)
}

// 需要超级管理员角色
superAdmin := router.Group("/system")
superAdmin.Use(handlers.Auth.Handle())
superAdmin.Use(middleware.RequireSuperAdmin())  // 快捷方法
{
    superAdmin.POST("/config", handler.UpdateConfig)
}
```

#### 方式 B: 权限检查（推荐用于细粒度控制）

```go
// 需要特定权限
users := router.Group("/users")
users.Use(handlers.Auth.Handle())
users.Use(middleware.RequirePermission(
    auth.PermissionUserWrite,
    auth.PermissionUserDelete,
))
{
    users.DELETE("/:id", handler.DeleteUser)
}

// 需要任意一个权限即可
content := router.Group("/content")
content.Use(handlers.Auth.Handle())
content.Use(middleware.RequireAnyPermission(
    auth.PermissionContentWrite,
    auth.PermissionContentAudit,
))
{
    content.PUT("/:id", handler.UpdateContent)
}
```

### 3. 在 Handler 内部进行权限检查

```go
func (h *Handler) UpdateUser(c *gin.Context) {
    // 获取 RBAC Claims
    claims := middleware.GetRBACClaims(c)
    if claims == nil {
        response.Error(c, response.ErrUnauthorized)
        return
    }

    // 检查角色
    if !claims.HasRole(auth.RoleAdmin) {
        response.Error(c, response.ErrForbidden)
        return
    }

    // 检查权限
    if !claims.HasPermission(auth.PermissionUserWrite) {
        response.Error(c, response.ErrForbidden)
        return
    }

    // 执行操作...
}
```

### 4. 辅助方法

```go
// 在 Handler 中快速检查
func (h *Handler) SomeAction(c *gin.Context) {
    // 检查是否是管理员
    if !middleware.IsAdmin(c) {
        response.Error(c, response.ErrForbidden)
        return
    }

    // 获取当前用户角色
    role := middleware.GetUserRole(c)
    
    // 检查是否有权限
    if !middleware.HasPermission(c, auth.PermissionSystemConfig) {
        response.Error(c, response.ErrForbidden)
        return
    }
}
```

---

## 💡 最佳实践

### 1. 路由保护策略

```go
// ✅ 推荐：先认证，再授权
admin := router.Group("/admin")
admin.Use(handlers.Auth.Handle())              // 1. 认证
admin.Use(middleware.RequireAdmin())           // 2. 授权

// ❌ 错误：顺序颠倒
admin := router.Group("/admin")
admin.Use(middleware.RequireAdmin())           // 错误：此时还没认证
admin.Use(handlers.Auth.Handle())
```

### 2. 多层权限控制

```go
// 允许多个角色访问
users := router.Group("/users")
users.Use(handlers.Auth.Handle())
users.Use(middleware.RequireRole(
    auth.RoleAdmin,
    auth.RoleSuperAdmin,
    auth.RoleModerator,  // 版主也可以访问
))

// 同时检查角色和权限
sensitive := router.Group("/sensitive")
sensitive.Use(handlers.Auth.Handle())
sensitive.Use(middleware.RequireRole(auth.RoleAdmin, auth.RoleSuperAdmin))
sensitive.Use(middleware.RequirePermission(auth.PermissionSystemConfig))
```

### 3. 动态权限判断

```go
func (h *Handler) UpdateUser(c *gin.Context) {
    claims := middleware.GetRBACClaims(c)
    targetUserID := c.Param("id")
    
    // 普通用户只能修改自己的信息
    if claims.Role == auth.RoleUser {
        if claims.UserID != targetUserID {
            response.Error(c, response.ErrForbidden)
            return
        }
    }
    
    // 管理员可以修改任何用户
    // ...
}
```

### 4. 错误消息优化

```go
// ✅ 推荐：提供清晰的错误信息
if !claims.HasPermission(auth.PermissionUserDelete) {
    response.Error(c, response.New(
        response.CodeForbidden,
        "权限不足：删除用户需要 user:delete 权限",
    ))
    return
}

// ❌ 避免：泄露系统信息
if !claims.HasPermission(auth.PermissionSystemConfig) {
    response.Error(c, response.New(
        response.CodeForbidden,
        "权限不足",  // 不要暴露具体需要什么权限
    ))
    return
}
```

---

## 📝 示例代码

### 完整的用户管理路由

```go
func setupUserRoutes(rg *gin.RouterGroup, handlers *Handlers) {
    users := rg.Group("/users")
    {
        // 公开路由
        users.POST("/register", handlers.User.Register)
        users.POST("/login", handlers.User.Login)

        // 普通用户路由
        profile := users.Group("/me")
        profile.Use(handlers.Auth.Handle())  // 只需要认证
        {
            profile.GET("", handlers.User.GetProfile)
            profile.PUT("", handlers.User.UpdateProfile)
        }

        // 管理员路由
        admin := users.Group("")
        admin.Use(handlers.Auth.Handle())
        admin.Use(middleware.RequireAdmin())
        {
            admin.GET("", handlers.User.ListUsers)
            admin.GET("/:id", handlers.User.GetUser)
            admin.PUT("/:id", handlers.User.UpdateUser)
        }

        // 超级管理员路由
        superAdmin := users.Group("")
        superAdmin.Use(handlers.Auth.Handle())
        superAdmin.Use(middleware.RequireSuperAdmin())
        {
            superAdmin.DELETE("/:id", handlers.User.DeleteUser)
        }
    }
}
```

### Handler 中的权限检查

```go
func (h *Handler) DeleteUser(c *gin.Context) {
    // 1. 获取 Claims
    claims := middleware.GetRBACClaims(c)
    if claims == nil {
        response.Error(c, response.ErrUnauthorized)
        return
    }

    // 2. 检查是否是超级管理员
    if !claims.IsSuperAdmin() {
        response.Error(c, response.New(
            response.CodeForbidden,
            "只有超级管理员可以删除用户",
        ))
        return
    }

    // 3. 防止自我删除
    targetUserID := c.Param("id")
    if claims.UserID == targetUserID {
        response.Error(c, response.New(
            response.CodeForbidden,
            "不能删除自己的账号",
        ))
        return
    }

    // 4. 执行删除
    err := h.userService.DeleteUser(c.Request.Context(), targetUserID)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, gin.H{"message": "用户已删除"})
}
```

---

## 🔧 配置与初始化

### Wire 依赖注入配置

```go
// wire.go
var HandlerSet = wire.NewSet(
    // ...
    provideRBACJWTManager,  // 添加 RBAC JWT Manager
    middleware.NewRBACMiddleware,
    // ...
)

func provideRBACJWTManager(cfg *config.Config) *auth.RBACJWTManager {
    return auth.NewRBACJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)
}
```

### 用户登录时设置角色

```go
func (h *Handler) Login(c *gin.Context) {
    // ... 验证用户 ...

    // 根据用户信息确定角色
    role := determineUserRole(user)
    
    // 生成包含角色的 Token
    token, err := h.rbacJWTManager.GenerateToken(
        user.ID,
        role,
        // 可选：额外权限
    )
    
    // ...
}

func determineUserRole(user *User) auth.Role {
    if user.IsSuperAdmin {
        return auth.RoleSuperAdmin
    }
    if user.IsAdmin {
        return auth.RoleAdmin
    }
    return auth.RoleUser
}
```

---

## 📊 权限矩阵

| 操作 | guest | user | moderator | admin | super_admin |
|------|-------|------|-----------|-------|-------------|
| 注册/登录 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 查看自己信息 | ❌ | ✅ | ✅ | ✅ | ✅ |
| 修改自己信息 | ❌ | ✅ | ✅ | ✅ | ✅ |
| 查看他人信息 | ❌ | ❌ | ✅ | ✅ | ✅ |
| 修改他人信息 | ❌ | ❌ | ❌ | ✅ | ✅ |
| 删除用户 | ❌ | ❌ | ❌ | ❌ | ✅ |
| 审核内容 | ❌ | ❌ | ✅ | ✅ | ✅ |
| 系统配置 | ❌ | ❌ | ❌ | ❌ | ✅ |

---

## 🔍 调试技巧

### 1. 打印当前用户权限

```go
func (h *Handler) Debug(c *gin.Context) {
    claims := middleware.GetRBACClaims(c)
    if claims != nil {
        c.JSON(200, gin.H{
            "user_id":     claims.UserID,
            "role":        claims.Role,
            "permissions": claims.Permissions,
            "role_level":  claims.GetRoleLevel(),
        })
    }
}
```

### 2. 日志记录

```go
func (h *Handler) SensitiveAction(c *gin.Context) {
    claims := middleware.GetRBACClaims(c)
    
    slog.InfoContext(c.Request.Context(), "Sensitive action requested",
        "user_id", claims.UserID,
        "role", claims.Role,
        "permissions", claims.Permissions,
        "action", "delete_user",
    )
    
    // ...
}
```

---

## 📚 相关文档

- [JWT 认证文档](./JWT.md)
- [API 文档](./API.md)
- [架构设计](./ARCHITECTURE.md)
