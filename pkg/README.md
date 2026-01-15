# pkg - 可复用的公共工具包

> **设计原则**: pkg 目录包含的是 **纯工具性质的代码**，可以被其他项目直接复用，不应包含任何业务逻辑。

---

## 📋 目录结构

```
pkg/
├── auth/           # 认证工具（JWT）
├── cache/          # 缓存管理（Redis）
├── database/       # 数据库连接工具
├── errors/         # 错误处理工具
├── ginx/           # Gin 框架扩展
├── health/         # 健康检查接口（纯接口定义）
└── logger/         # 日志工具
```

---

## 🎯 设计原则

### ✅ 应该放在 pkg 的内容

1. **通用工具函数** - 可在任何项目中复用
2. **第三方库封装** - 如 Redis、JWT 的通用封装
3. **接口定义** - 通用的接口和数据结构
4. **标准库扩展** - 对 Go 标准库的增强

### ❌ 不应该放在 pkg 的内容

1. **业务逻辑** - 与具体业务相关的代码
2. **领域模型** - 业务实体和聚合根
3. **具体实现** - 依赖业务规则的具体实现
4. **internal 引用** - 不应引用 `internal/` 下的任何代码

---

## 📦 各包说明

### auth/ - 认证工具

**功能**: JWT Token 生成、验证、刷新

**通用性**: ✅ 完全通用，可用于任何需要 JWT 的项目

```go
jwtManager := auth.NewJWTManager(secret, 24*time.Hour)
token, _ := jwtManager.GenerateToken(userID, username, email)
```

**依赖**: `github.com/golang-jwt/jwt/v5`

---

### cache/ - 缓存管理

**功能**: Redis 缓存操作、防击穿/穿透/雪崩

**通用性**: ✅ 通用的缓存模式，适用于任何项目

```go
manager := cache.NewManager(redisClient)
user, _ := cache.TakeByID(ctx, manager, "user", userID, 5*time.Minute, queryFunc)
```

**依赖**: `github.com/redis/go-redis/v9`

---

### database/ - 数据库工具

**功能**: PostgreSQL 连接池配置

**通用性**: ✅ 通用的数据库连接工具

```go
db, _ := database.NewPostgres(config)
```

**依赖**: `database/sql`, `github.com/lib/pq`

---

### errors/ - 错误处理

**功能**: 业务错误码、错误包装

**通用性**: ✅ 通用的错误处理模式

```go
err := errors.New(errors.CodeNotFound, "用户不存在")
```

**依赖**: 仅标准库

---

### ginx/ - Gin 扩展

**功能**: 统一响应、参数绑定、分页工具

**通用性**: ✅ Gin 框架的通用扩展

```go
ginx.Success(c, data)
ginx.Error(c, err)
pagination := ginx.GetPagination(c)
```

**依赖**: `github.com/gin-gonic/gin`

---

### health/ - 健康检查接口 ⭐

**功能**: 定义健康检查的通用接口

**通用性**: ✅ 纯接口定义，无具体实现

**特点**:
- 只定义接口和数据结构
- 具体实现在 `internal/health/` 中
- 符合 "接口定义与实现分离" 原则

```go
// pkg/health/ 定义接口
type Checker interface {
    Check(ctx context.Context) HealthStatus
}

// internal/health/ 提供具体实现
type DatabaseChecker struct { ... }
```

**依赖**: 仅标准库

---

### logger/ - 日志工具

**功能**: 结构化日志配置

**通用性**: ✅ 基于标准库 slog 的通用配置

```go
logger.Setup(config)
```

**依赖**: `log/slog`

---

## 🔍 如何判断代码是否应该放在 pkg

### 检查清单

- [ ] 代码是否可以在其他项目中直接使用？
- [ ] 代码是否依赖 `internal/` 目录？
- [ ] 代码是否包含业务逻辑？
- [ ] 代码是否依赖领域模型？

**如果回答**:
- ✅ 可以在其他项目使用 → 放 `pkg/`
- ❌ 依赖 internal 或业务逻辑 → 放 `internal/`

### 重构案例：健康检查

**之前（不合理）**:
```go
// pkg/health/checker.go
type Checker struct {
    db    *sql.DB
    redis *redis.Client
}

func (c *Checker) Check() {
    // 硬编码数据库和 Redis 检查
}
```
**问题**: 具体实现放在 pkg 中，不够通用

**之后（合理）**:
```go
// pkg/health/checker.go - 只定义接口
type Checker interface {
    Check(ctx context.Context) HealthStatus
}

type ComponentChecker interface {
    Name() string
    Check(ctx context.Context) Check
    IsCritical() bool
}

// internal/health/database_checker.go - 具体实现
type DatabaseChecker struct { ... }
```
**优点**: 
- pkg 保持纯接口定义
- 具体实现在 internal 中
- 其他项目可以实现自己的 Checker

---

## 🚀 使用示例

### 在其他项目中使用

```go
import (
    "your-project/pkg/auth"
    "your-project/pkg/cache"
    "your-project/pkg/ginx"
)

// 完全可以直接复用
jwtManager := auth.NewJWTManager(secret, expiration)
cacheManager := cache.NewManager(redisClient)
ginx.Success(c, data)
```

### 扩展示例

```go
// 实现自己的健康检查器
type MongoChecker struct {
    client *mongo.Client
}

func (m *MongoChecker) Name() string { return "mongodb" }
func (m *MongoChecker) Check(ctx context.Context) health.Check { ... }
func (m *MongoChecker) IsCritical() bool { return true }

// 使用通用的 MultiChecker
checker := health.NewMultiChecker("1.0.0", 
    NewMongoChecker(mongoClient),
    NewRedisChecker(redisClient),
)
```

---

## 📝 最佳实践

1. **保持通用性** - 代码应该适用于不同的业务场景
2. **最小依赖** - 只依赖标准库和常用第三方库
3. **接口优先** - 定义清晰的接口，而非硬编码实现
4. **文档完善** - 提供清晰的使用示例
5. **单元测试** - 每个工具都应有完整的测试

---

## 🔄 重构指南

如果发现 pkg 中的代码有业务侵入：

1. **识别问题**
   - 代码是否引用了 `internal/`？
   - 代码是否依赖特定业务逻辑？

2. **提取接口**
   - 在 `pkg/` 中定义通用接口
   - 保留数据结构和常量定义

3. **移动实现**
   - 将具体实现移到 `internal/`
   - 保持接口在 `pkg/` 中

4. **更新依赖**
   - 更新 Wire 配置
   - 更新导入路径
   - 重新生成依赖注入代码

---

## ✅ 验证清单

运行以下命令验证 pkg 的纯净性：

```bash
# 检查是否有 pkg 引用 internal
grep -r "gin_demo/internal" pkg/
# 应该返回：无匹配结果

# 检查是否有业务相关的导入
grep -r "repository\|service\|domain" pkg/
# 应该返回：无匹配结果
```

---

**维护者注意**: 任何往 pkg 添加代码的 PR 都应该严格审查，确保符合上述设计原则。
