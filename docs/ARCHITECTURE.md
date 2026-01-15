# 架构设计文档

## 目录

- [概述](#概述)
- [架构原则](#架构原则)
- [分层架构](#分层架构)
- [核心组件](#核心组件)
- [数据流](#数据流)
- [设计模式](#设计模式)
- [技术选型](#技术选型)

---

## 概述

本项目采用经典的**三层架构**，结合**领域驱动设计（DDD）**的部分思想，构建了一个清晰、可维护、可扩展的 Go Web API 应用。

### 核心目标

- ✅ **清晰的职责分离** - 每一层只关注自己的职责
- ✅ **易于测试** - 接口化设计，方便 Mock
- ✅ **高性能** - 缓存优化、连接池、并发控制
- ✅ **可维护** - 代码结构清晰，注释完整
- ✅ **可扩展** - 易于添加新功能和模块

---

## 架构原则

### 1. 单一职责原则（SRP）

每个模块、每个类只负责一件事：

- **Handler** - 只处理 HTTP 请求和响应
- **Service** - 只处理业务逻辑
- **Repository** - 只处理数据访问

### 2. 依赖倒置原则（DIP）

高层模块不依赖低层模块，都依赖抽象：

```go
// Service 依赖 Repository 接口，而非具体实现
type UserService interface {
    Register(ctx context.Context, ...) (User, error)
}

type userService struct {
    userRepo *UserRepository // 依赖注入
}
```

### 3. 接口隔离原则（ISP）

定义小而专注的接口：

```go
// ✅ 正确：小接口
type UserReader interface {
    GetUserByID(ctx context.Context, userID int64) (User, error)
}

// ❌ 错误：大接口
type UserRepository interface {
    // 包含 20 个方法...
}
```

### 4. 开闭原则（OCP）

对扩展开放，对修改关闭：

- 新增功能通过**新增代码**实现，而非修改已有代码
- 使用**中间件**扩展功能
- 使用**策略模式**实现可替换的算法

---

## 分层架构

```
┌────────────────────────────────────────────────────────┐
│                    Presentation Layer                  │
│                  (Handler + Middleware)                │
│  - HTTP 请求处理                                        │
│  - 参数验证和绑定                                       │
│  - 响应格式化                                           │
│  - 错误处理                                             │
└────────────────────────────────────────────────────────┘
                            ↓
┌────────────────────────────────────────────────────────┐
│                     Business Layer                     │
│                       (Service)                        │
│  - 业务逻辑                                             │
│  - 权限校验                                             │
│  - 事务管理                                             │
│  - 调用多个 Repository                                  │
└────────────────────────────────────────────────────────┘
                            ↓
┌────────────────────────────────────────────────────────┐
│                   Data Access Layer                    │
│                    (Repository)                        │
│  - 数据库访问                                           │
│  - 缓存管理                                             │
│  - SQL 生成（sqlc）                                     │
│  - 事务支持                                             │
└────────────────────────────────────────────────────────┘
                            ↓
┌────────────────────────────────────────────────────────┐
│                Infrastructure Layer                    │
│                (Database + Cache)                      │
│  - PostgreSQL                                          │
│  - Redis                                               │
└────────────────────────────────────────────────────────┘
```

### Handler Layer（表示层）

**职责：**
- 接收 HTTP 请求
- 参数验证和绑定
- 调用 Service
- 格式化响应
- 错误处理

**示例：**

```go
func (h *UserHandler) Register(c *gin.Context) {
    // 1. 参数绑定
    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, dto.Error(dto.CodeInvalidParams, "参数错误"))
        return
    }

    // 2. 调用 Service
    user, err := h.userService.Register(c.Request.Context(), 
        req.Username, req.Email, req.Password)
    
    // 3. 错误处理
    if err != nil {
        if errors.Is(err, service.ErrUserExists) {
            c.JSON(409, dto.Error(dto.CodeAlreadyExists, "用户已存在"))
            return
        }
        c.JSON(500, dto.Error(dto.CodeInternalError, "注册失败"))
        return
    }

    // 4. 成功响应
    c.JSON(200, dto.Success(convertToUserResponse(user)))
}
```

### Service Layer（业务层）

**职责：**
- 实现业务逻辑
- 权限校验
- 调用多个 Repository
- 事务管理
- 错误封装

**示例：**

```go
func (s *userService) Register(ctx context.Context, username, email, password string) (User, error) {
    // 1. 业务规则验证
    if username == "" || email == "" || password == "" {
        return User{}, ErrInvalidInput
    }

    // 2. 检查用户是否存在
    existingUser, _ := s.userRepo.GetUserByEmail(ctx, email)
    if existingUser.ID > 0 {
        return User{}, ErrUserExists
    }

    // 3. 密码加密（业务逻辑）
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

    // 4. 调用 Repository
    return s.userRepo.CreateUser(ctx, CreateUserParams{
        Username: username,
        Email:    email,
        Password: string(hashedPassword),
    })
}
```

### Repository Layer（数据访问层）

**职责：**
- 数据库访问
- 缓存管理
- SQL 查询执行
- 数据映射

**示例：**

```go
func (r *UserRepository) GetUserByID(ctx context.Context, userID int64) (User, error) {
    // 1. 先查缓存
    return cache.TakeByID(ctx, r.cache, "user", userID, 5*time.Minute,
        // 2. 缓存未命中，查数据库
        func(ctx context.Context) (User, error) {
            return r.queries.GetUserByID(ctx, userID)
        })
}
```

---

## 核心组件

### 1. 配置管理（Config）

**位置：** `config/config.go`

**职责：**
- 从环境变量加载配置
- 验证配置有效性
- 提供配置访问接口

**特点：**
- 支持默认值
- 类型安全
- 统一管理

### 2. 日志系统（Logger）

**位置：** `pkg/logger/log.go`

**特点：**
- 基于 `log/slog`
- 结构化日志
- 自动包含 `request_id`
- 支持 JSON 格式
- 支持源码位置

**示例：**

```go
slog.InfoContext(ctx, "User created",
    "user_id", 123,
    "username", "alice",
)

// 输出：
// {"time":"2024-01-01 10:00:00","level":"INFO","msg":"User created",
//  "request_id":"xxx","user_id":123,"username":"alice"}
```

### 3. 缓存管理（Cache Manager）

**位置：** `pkg/cache/manager.go`

**特点：**
- 主键 + 索引双层缓存
- 三层防护（防击穿/穿透/雪崩）
- 自动过期和清理
- 泛型支持

**核心方法：**

```go
// 主键查询
TakeByID[T](ctx, manager, "user", 123, 5*time.Minute, queryFn)

// 索引查询
TakeByIndex[T, ID](ctx, manager, "user", "email", email, 5*time.Minute, ...)

// 写操作（自动清理缓存）
ExecByID(ctx, "user", 123, execFn)
```

### 4. 数据库连接（Database）

**位置：** `pkg/database/postgres.go`

**特点：**
- 连接池优化
- 自动 Ping 测试
- 超时控制
- 优雅关闭

**配置：**

```go
MaxOpenConns:    25,              // 最大打开连接数
MaxIdleConns:    5,               // 最大空闲连接数
ConnMaxLifetime: 5 * time.Minute, // 连接最大生命周期
ConnMaxIdleTime: 1 * time.Minute, // 连接最大空闲时间
```

### 5. 中间件（Middleware）

**位置：** `internal/app/middleware/`

#### CORS 中间件

```go
engine.Use(middleware.CORS())
```

#### 日志中间件

```go
engine.Use(middleware.Logger())
```

#### 限流中间件

```go
limiter := middleware.NewTokenBucketLimiter(100, 200)
engine.Use(middleware.RateLimit(limiter))
```

#### Recovery 中间件

```go
engine.Use(middleware.Recovery())
```

---

## 数据流

### 读操作流程

```
1. 客户端请求
   ↓
2. Middleware (RequestID, Logger, RateLimit, etc.)
   ↓
3. Handler (参数验证)
   ↓
4. Service (业务逻辑)
   ↓
5. Repository (查缓存)
   ├─ 缓存命中 → 返回数据
   └─ 缓存未命中
       ↓
      查数据库
       ↓
      写入缓存
       ↓
      返回数据
   ↓
6. 响应客户端
```

### 写操作流程

```
1. 客户端请求
   ↓
2. Middleware
   ↓
3. Handler (参数验证)
   ↓
4. Service (业务逻辑、权限校验)
   ↓
5. Repository
   ├─ 开始事务（可选）
   ├─ 写入数据库
   ├─ 清理缓存
   └─ 提交事务
   ↓
6. 响应客户端
```

---

## 设计模式

### 1. 依赖注入（DI）

```go
// main.go
cacheManager := cache.NewManager(rdb)
userRepo := repository.NewUserRepository(db, cacheManager)
userService := service.NewUserService(userRepo)
userHandler := handler.NewUserHandler(userService)
```

**优点：**
- 易于测试（可 Mock）
- 降低耦合
- 提高可维护性

### 2. 仓储模式（Repository Pattern）

```go
type UserRepository struct {
    queries *Queries
    cache   *cache.Manager
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID int64) (User, error) {
    // 封装数据访问逻辑
}
```

**优点：**
- 统一数据访问
- 易于切换存储
- 支持缓存

### 3. 策略模式（Strategy Pattern）

```go
type RateLimiter interface {
    Allow(key string) bool
}

// 可以有多种限流策略
type TokenBucketLimiter struct { ... }
type LeakyBucketLimiter struct { ... }
```

### 4. 装饰器模式（Decorator Pattern）

中间件就是装饰器模式的应用：

```go
engine.Use(middleware1)  // 装饰 1
engine.Use(middleware2)  // 装饰 2
```

---

## 技术选型

### 为什么选择 Gin？

- ✅ 高性能（基于 httprouter）
- ✅ 丰富的中间件生态
- ✅ 简单易用的 API
- ✅ 完善的文档

### 为什么选择 sqlc？

- ✅ 类型安全（编译时检查）
- ✅ 无运行时开销
- ✅ 标准 SQL 语法
- ✅ 自动生成 Go 代码

### 为什么选择 PostgreSQL？

- ✅ 成熟稳定
- ✅ 功能丰富（JSON、全文搜索等）
- ✅ ACID 支持
- ✅ 开源免费

### 为什么选择 Redis？

- ✅ 高性能（内存存储）
- ✅ 丰富的数据结构
- ✅ 支持持久化
- ✅ 集群支持

### 为什么选择 slog？

- ✅ Go 标准库（1.21+）
- ✅ 结构化日志
- ✅ 高性能
- ✅ 零依赖

---

## 扩展性

### 添加新模块

1. 创建数据库表（`db/migrations/`）
2. 定义 SQL 查询（`db/queries/`）
3. 生成 sqlc 代码（`sqlc generate`）
4. 实现 Repository（`internal/repository/`）
5. 实现 Service（`internal/domain/service/`）
6. 实现 Handler（`internal/app/handler/`）
7. 注册路由（`main.go`）

### 添加新中间件

```go
// internal/app/middleware/my_middleware.go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        c.Next()
        // 后置处理
    }
}

// main.go
engine.Use(middleware.MyMiddleware())
```

---

## 性能优化

### 数据库层

- ✅ 连接池优化
- ✅ 索引优化
- ✅ 使用 Prepared Statements
- ✅ 批量操作

### 缓存层

- ✅ 主键 + 索引缓存
- ✅ 缓存预热
- ✅ 自动过期
- ✅ 防击穿/穿透/雪崩

### 并发层

- ✅ singleflight 防击穿
- ✅ 协程池
- ✅ 超时控制
- ✅ 限流保护

---

## 安全性

### 密码安全

- ✅ bcrypt 加密
- ✅ 盐值随机
- ✅ 不返回密码字段

### SQL 注入防护

- ✅ 使用 sqlc 生成的参数化查询
- ✅ 输入验证

### 限流保护

- ✅ 令牌桶算法
- ✅ 基于 IP 限流
- ✅ 自动清理过期数据

---

## 监控和日志

### 结构化日志

```go
slog.InfoContext(ctx, "Request completed",
    "method", "GET",
    "path", "/api/v1/users/123",
    "status", 200,
    "latency", "10ms",
)
```

### 追踪

- ✅ Request ID 自动生成
- ✅ 日志自动关联
- ✅ 分布式追踪支持（可扩展）

---

## 测试策略

### 单元测试

```go
func TestUserService_Register(t *testing.T) {
    // Mock Repository
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    // 测试
    user, err := service.Register(ctx, "alice", "alice@example.com", "password")
    assert.NoError(t, err)
    assert.Equal(t, "alice", user.Username)
}
```

### 集成测试

```go
func TestUserAPI_Register(t *testing.T) {
    // 启动测试服务器
    router := setupRouter()
    
    // 发送请求
    resp := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/users/register", body)
    router.ServeHTTP(resp, req)
    
    // 验证响应
    assert.Equal(t, 200, resp.Code)
}
```

---

## 总结

本项目采用**清晰的分层架构**、**接口化设计**、**依赖注入**等最佳实践，构建了一个**高性能**、**可维护**、**可扩展**的 Go Web API 应用。

**核心优势：**

- ✅ 清晰的职责分离
- ✅ 易于测试和维护
- ✅ 高性能（缓存 + 优化）
- ✅ 完善的错误处理
- ✅ 结构化日志
- ✅ 生产级可用

---

**参考资料：**

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://www.domainlanguage.com/ddd/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
