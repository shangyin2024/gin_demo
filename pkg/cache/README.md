# Cache Manager 使用文档

## 概述

一个生产级的 Redis 缓存管理器，提供以下特性：

- ✅ **防缓存击穿** - 使用 singleflight 合并并发请求
- ✅ **防缓存穿透** - 缓存空结果（NotFoundPlaceholder）
- ✅ **防缓存雪崩** - 随机过期时间（Jitter）
- ✅ **主键 + 索引** - 支持 ID 和索引查询
- ✅ **泛型支持** - 类型安全
- ✅ **自动清理** - 写操作自动删除相关缓存

---

## 快速开始

### 1. 初始化

```go
import (
    "gin_demo/pkg/cache"
    "github.com/redis/go-redis/v9"
)

// 创建 Redis 客户端
rdb := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// 创建缓存管理器
cacheManager := cache.NewManager(rdb)
```

### 2. 主键查询（TakeByID）

```go
user, err := cache.TakeByID(ctx, cacheManager, "user", userID, 5*time.Minute,
    func(ctx context.Context) (User, error) {
        // 从数据库查询
        return db.GetUserByID(ctx, userID)
    })
```

**流程：**
1. 查缓存 `cache:user:123`
2. 命中 → 返回数据
3. 未命中 → 调用 `queryFn` 从 DB 查询
4. 使用 `singleflight` 防止并发击穿
5. 写入缓存（带 Jitter）

### 3. 索引查询（TakeByIndex）

```go
user, err := cache.TakeByIndex[User, int64](
    ctx, cacheManager, "user", "email", email, 5*time.Minute,
    // 查询 ID
    func(ctx context.Context) (int64, error) {
        return db.GetUserIDByEmail(ctx, email)
    },
    // 通过 ID 查询数据
    func(ctx context.Context, id int64) (User, error) {
        return db.GetUserByID(ctx, id)
    },
    // ID 转换器
    func(idStr string) (int64, error) {
        return strconv.ParseInt(idStr, 10, 64)
    },
)
```

**流程：**
1. 查索引缓存 `cache:user:email:alice@example.com` → ID
2. 拿到 ID 后，走主键缓存逻辑
3. 最终返回完整数据

### 4. 更新操作（ExecByID）

```go
err := cacheManager.ExecByID(ctx, "user", userID, func(ctx context.Context) error {
    return db.UpdateUser(ctx, userID, newData)
})
```

**流程：**
1. 执行数据库更新
2. 成功后删除缓存 `cache:user:123`

### 5. 更新操作（ExecByIDWithIndexes）

```go
indexes := []string{
    cacheManager.BuildIndexKey("user", "email", oldEmail),
    cacheManager.BuildIndexKey("user", "email", newEmail),
}

err := cacheManager.ExecByIDWithIndexes(ctx, "user", userID, indexes, 
    func(ctx context.Context) error {
        return db.UpdateUserEmail(ctx, userID, newEmail)
    })
```

**流程：**
1. 执行数据库更新
2. 成功后删除：
   - 主键缓存 `cache:user:123`
   - 旧索引 `cache:user:email:old@example.com`
   - 新索引 `cache:user:email:new@example.com`

---

## 三大防护机制

### 1. 防缓存击穿（Cache Breakdown）

**问题：** 热点 Key 过期，大量请求同时查 DB

**解决：** 使用 `singleflight`

```go
// 100 个并发请求同一个 Key
// 只有 1 个请求真正查 DB
// 其他 99 个等待并共享结果
raw, err, _ := sfGroup.Do(key, func() (any, error) {
    return queryDB()
})
```

### 2. 防缓存穿透（Cache Penetration）

**问题：** 恶意请求不存在的数据，每次都查 DB

**解决：** 缓存空结果

```go
if errors.Is(err, sql.ErrNoRows) {
    // 数据不存在，缓存占位符（短 TTL）
    _ = rdb.Set(ctx, key, NotFoundPlaceholder, DefaultNotFoundTTL).Err()
}
```

### 3. 防缓存雪崩（Cache Avalanche）

**问题：** 大量 Key 同时过期，DB 压力激增

**解决：** 随机过期时间

```go
func (m *Manager) getJitterTTL(baseTTL time.Duration) time.Duration {
    // 20% 范围波动 + 0~30秒噪声
    jitter := rand.Int63n(int64(baseTTL) / 5)
    noise := time.Duration(rand.Int63n(30)) * time.Second
    return baseTTL + time.Duration(jitter) + noise
}

// 5分钟基础 TTL → 实际 5~7分钟
```

---

## 缓存 Key 设计

### 主键缓存

```
cache:user:123
cache:order:456
```

### 索引缓存

```
cache:user:email:alice@example.com  → "123"
cache:user:phone:13800138000        → "456"
```

---

## 最佳实践

### ✅ DO

1. **使用合适的 TTL**
   ```go
   // 热点数据：较长 TTL
   TakeByID(ctx, m, "user", id, 30*time.Minute, queryFn)
   
   // 普通数据：中等 TTL
   TakeByID(ctx, m, "order", id, 5*time.Minute, queryFn)
   ```

2. **索引查询复用主键缓存**
   ```go
   // ✅ 正确：索引查到 ID 后，调用 GetUserByID 复用主键缓存
   func (r *Repo) GetUserByEmail(ctx, email) (User, error) {
       return TakeByIndex(ctx, m, "user", "email", email, 5*time.Minute,
           func(ctx) (int64, error) { return getIDByEmail(ctx, email) },
           func(ctx, id) (User, error) { return r.GetUserByID(ctx, id) },
           parseID,
       )
   }
   ```

3. **写操作清理索引**
   ```go
   // ✅ 更新 Email 时，清理旧索引 + 新索引
   indexes := []string{
       m.BuildIndexKey("user", "email", oldEmail),
       m.BuildIndexKey("user", "email", newEmail),
   }
   m.ExecByIDWithIndexes(ctx, "user", id, indexes, updateFn)
   ```

### ❌ DON'T

1. **不要缓存所有数据**
   ```go
   // ❌ 不要缓存低频、大体积数据
   TakeByID(ctx, m, "log", id, 1*time.Hour, queryFn)
   ```

2. **不要忘记清理索引**
   ```go
   // ❌ 更新 Email 但不清理索引缓存
   m.ExecByID(ctx, "user", id, updateEmailFn) // 旧索引还在！
   ```

3. **不要使用过长 TTL**
   ```go
   // ❌ 24 小时太长，数据可能不一致
   TakeByID(ctx, m, "user", id, 24*time.Hour, queryFn)
   ```

---

## 性能优化

### 批量预热

```go
func (r *Repo) WarmupUsers(ctx context.Context, userIDs []int64) error {
    for _, id := range userIDs {
        go func(id int64) {
            _, _ = r.GetUserByID(ctx, id) // 触发缓存填充
        }(id)
    }
    return nil
}
```

### 监控指标

```go
// 可以添加 Prometheus 指标
var (
    cacheHits   = prometheus.NewCounter(...)
    cacheMisses = prometheus.NewCounter(...)
)
```

---

## 完整示例

参考 `example.go` 文件，包含：
- 主键查询
- 索引查询
- 更新操作（清理主键）
- 更新操作（清理主键 + 索引）
- 删除操作

---

## 总结

| 特性 | 实现 |
|------|------|
| 防击穿 | singleflight |
| 防穿透 | NotFoundPlaceholder |
| 防雪崩 | getJitterTTL |
| 类型安全 | 泛型 |
| 索引支持 | TakeByIndex |
| 自动清理 | ExecByID |

**这是一个生产级的缓存实现！** 🎯
