# 🔄 API 版本管理策略

**当前版本**: v1  
**规划**: v2, v3...  
**策略**: 渐进式演进，向后兼容

---

## 📋 目录

1. [版本管理策略](#版本管理策略)
2. [当前 v1 现状](#当前-v1-现状)
3. [v2 规划建议](#v2-规划建议)
4. [实施方案](#实施方案)
5. [最佳实践](#最佳实践)

---

## 1. 版本管理策略

### 何时需要新版本

#### 必须升级（Breaking Changes）

```
✅ 需要新版本的情况:
1. 删除字段
   v1: {"user_id": 1, "old_field": "x"}
   v2: {"user_id": 1}  ← 删除了 old_field

2. 改变字段类型
   v1: {"created_at": "2024-01-01"}  (string)
   v2: {"created_at": 1704067200}     (unix timestamp)

3. 改变字段名
   v1: {"user_id": 1}
   v2: {"id": 1}  ← user_id 改为 id

4. 改变响应结构
   v1: [{"user": {...}}]
   v2: {"users": [...], "total": 100}

5. 改变业务逻辑
   v1: 注册时不发送验证邮件
   v2: 注册时必须验证邮件
```

#### 无需升级（兼容性变更）

```
❌ 不需要新版本的情况:
1. 添加新字段（可选）
   v1: {"user_id": 1}
   v1: {"user_id": 1, "new_field": "x"}  ← 向后兼容

2. 添加新端点
   v1: /users
   v1: /users/stats  ← 新增端点

3. 添加可选参数
   v1: GET /users
   v1: GET /users?filter=active  ← 可选参数

4. 扩展错误信息
5. 性能优化
6. 内部重构
```

---

## 2. 当前 v1 现状

### v1 API 清单

```
用户相关:
POST   /api/v1/users/register      注册
POST   /api/v1/users/login         登录
GET    /api/v1/users/me            获取当前用户
GET    /api/v1/users/:id           获取指定用户
PUT    /api/v1/users/:id           更新用户
DELETE /api/v1/users/:id           删除用户
POST   /api/v1/users/:id/password  修改密码
GET    /api/v1/users                用户列表

健康检查:
GET    /health                      健康检查
GET    /health/ready                就绪检查
GET    /health/live                 存活检查

监控:
GET    /metrics                     Prometheus 指标
```

### v1 响应格式

```json
// 成功响应
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}

// 错误响应
{
  "code": 10001,
  "message": "用户不存在",
  "data": null
}

// 列表响应
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "total": 100
  }
}
```

---

## 3. v2 规划建议

### v2 可能的改进

#### 1. 更 RESTful 的设计

```go
// v1 (当前)
POST /api/v1/users/:id/password  ← 不符合 REST 风格

// v2 (建议)
PATCH /api/v2/users/:id/credentials  ← 更语义化

// v1 (当前)
GET /api/v1/users?page=1&limit=10

// v2 (建议)
GET /api/v2/users?page[number]=1&page[size]=10  ← JSON:API 风格
```

#### 2. 更丰富的响应格式

```json
// v1 (当前)
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "total": 100
  }
}

// v2 (建议 - JSON:API 风格)
{
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 10
  },
  "links": {
    "self": "/api/v2/users?page=1",
    "next": "/api/v2/users?page=2",
    "prev": null,
    "first": "/api/v2/users?page=1",
    "last": "/api/v2/users?page=10"
  }
}
```

#### 3. 更灵活的字段选择

```bash
# v1 (当前)
GET /api/v1/users/1
# 返回所有字段

# v2 (建议 - Sparse Fieldsets)
GET /api/v2/users/1?fields[user]=username,email
# 只返回指定字段

{
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com"
  }
}
```

#### 4. 关联资源加载

```bash
# v1 (当前)
GET /api/v1/users/1
# 需要额外请求获取关联数据

# v2 (建议 - Include)
GET /api/v2/users/1?include=profile,posts
# 一次请求返回关联数据

{
  "data": {
    "id": 1,
    "username": "alice",
    "relationships": {
      "profile": {"data": {"type": "profile", "id": 1}},
      "posts": {"data": [{"type": "post", "id": 10}]}
    }
  },
  "included": [
    {"type": "profile", "id": 1, "attributes": {...}},
    {"type": "post", "id": 10, "attributes": {...}}
  ]
}
```

#### 5. 批量操作支持

```json
// v2 (新增)
POST /api/v2/users/batch
{
  "operations": [
    {"method": "PATCH", "path": "/users/1", "body": {...}},
    {"method": "DELETE", "path": "/users/2"},
    {"method": "POST", "path": "/users", "body": {...}}
  ]
}
```

#### 6. GraphQL 支持（可选）

```graphql
# v2 (可选 - GraphQL 端点)
POST /api/v2/graphql

query {
  user(id: 1) {
    username
    email
    posts(limit: 10) {
      title
      createdAt
    }
  }
}
```

---

## 4. 实施方案

### 方案 A: URL 路径版本（推荐）⭐⭐⭐⭐⭐

```
优点:
✅ 简单直观
✅ 易于路由
✅ 浏览器友好
✅ 缓存友好

缺点:
⚠️ URL 冗余

实现:
/api/v1/users
/api/v2/users
/api/v3/users
```

**代码结构**:
```
internal/app/
├── v1/
│   ├── user/
│   │   ├── handler.go
│   │   ├── dto.go
│   │   └── routes.go
│   └── routes.go
│
├── v2/
│   ├── user/
│   │   ├── handler.go
│   │   ├── dto.go
│   │   └── routes.go
│   └── routes.go
│
└── routes.go  # 主路由
```

```go
// internal/app/routes.go
func (s *Server) SetupRoutes() {
    // v1 路由
    v1 := s.engine.Group("/api/v1")
    v1routes.Setup(v1, s.handlers)
    
    // v2 路由 (未来)
    v2 := s.engine.Group("/api/v2")
    v2routes.Setup(v2, s.handlers)
}
```

### 方案 B: Header 版本

```
优点:
✅ URL 清洁
✅ 更 RESTful

缺点:
⚠️ 不直观
⚠️ 浏览器测试不便
⚠️ 缓存复杂

实现:
GET /api/users
Header: Accept: application/vnd.myapp.v2+json
```

### 方案 C: 查询参数版本

```
优点:
✅ 灵活
✅ 可选（默认最新版）

缺点:
⚠️ 不够规范
⚠️ 缓存困难

实现:
/api/users?version=2
```

### 推荐：方案 A（URL 路径版本）

最简单、最直观、最易维护。

---

## 5. 最佳实践

### 5.1 版本生命周期管理

```
版本状态:
1. Active (活跃) - v3 当前最新
2. Deprecated (废弃) - v2 计划下线
3. Sunset (下线) - v1 已停止支持

时间线示例:
┌──────────────────────────────────────────┐
│ v1 │███████████████▓▓▓▓▓▓▓░░░░ (Sunset)  │
│ v2 │          ███████████████▓▓▓▓ (Depr) │
│ v3 │                    ████████ (Active) │
└──────────────────────────────────────────┘
     2024      2025      2026      2027
```

### 5.2 废弃通知机制

```go
// 在响应头中添加废弃信息
c.Header("X-API-Deprecation-Date", "2026-12-31")
c.Header("X-API-Deprecation-Info", "https://docs.example.com/api/v1-deprecation")
c.Header("X-API-Sunset-Date", "2027-06-30")

// 在响应中添加警告
{
  "code": 0,
  "message": "success",
  "warnings": [
    {
      "code": "DEPRECATED",
      "message": "API v1 will be sunset on 2027-06-30. Please migrate to v2.",
      "url": "https://docs.example.com/migration-guide"
    }
  ],
  "data": {...}
}
```

### 5.3 版本支持策略

```
策略:
1. 同时支持 N 个版本 (N=2 或 3)
2. 新版本发布后，旧版本进入废弃期
3. 废弃期至少 6 个月
4. 提前 3 个月通知下线

示例:
2024-01-01: v1 发布 (Active)
2025-01-01: v2 发布 (Active), v1 废弃 (Deprecated)
2025-07-01: v1 下线 (Sunset)
2026-01-01: v3 发布 (Active), v2 废弃 (Deprecated)
2026-07-01: v2 下线 (Sunset)
```

### 5.4 向后兼容的变更

```json
// ✅ 可以做的变更（不需要新版本）
1. 添加新字段（可选）
{
  "id": 1,
  "username": "alice",
  "new_field": "value"  // 新增，客户端可忽略
}

2. 添加新端点
POST /api/v1/users/export  // 新增端点

3. 添加可选参数
GET /api/v1/users?sort_by=created_at  // 新增可选参数

4. 放宽验证规则
username: 3-20 字符 → 3-50 字符  // 更宽松

// ❌ 不能做的变更（需要新版本）
1. 删除字段
{
  "id": 1
  // "username": "alice" 被删除了
}

2. 改变字段类型
"created_at": "2024-01-01" → 1704067200

3. 加强验证规则
username: 3-50 字符 → 3-20 字符  // 更严格

4. 改变业务逻辑
注册不验证邮件 → 注册必须验证邮件
```

### 5.5 文档管理

```
文档结构:
docs/
├── api/
│   ├── v1/
│   │   ├── overview.md
│   │   ├── users.md
│   │   ├── auth.md
│   │   └── openapi.yaml
│   │
│   ├── v2/
│   │   ├── overview.md
│   │   ├── migration-from-v1.md  ⭐ 迁移指南
│   │   ├── users.md
│   │   └── openapi.yaml
│   │
│   └── README.md  # 版本索引
```

### 5.6 测试策略

```go
// 为每个版本维护独立的测试
internal/app/
├── v1/
│   └── user/
│       ├── handler.go
│       └── handler_test.go  ← v1 测试
│
├── v2/
│   └── user/
│       ├── handler.go
│       └── handler_test.go  ← v2 测试
│
└── compatibility_test.go  ← 跨版本兼容性测试
```

---

## 📝 当前行动建议

### 立即行动（不需要 v2）

```
✅ 可以做的优化:
1. 添加更多可选字段
   - user.bio (个人简介)
   - user.avatar_url (头像)
   - user.last_login_at (最后登录时间)

2. 添加新端点
   - GET /api/v1/users/me/stats (用户统计)
   - POST /api/v1/users/me/avatar (上传头像)
   - GET /api/v1/users/search (搜索用户)

3. 优化查询参数
   - GET /api/v1/users?sort_by=created_at
   - GET /api/v1/users?filter[status]=active
   - GET /api/v1/users?include=profile

4. 增强错误信息
   - 更详细的错误描述
   - 错误码分类
   - 错误解决建议
```

### 未来规划（需要 v2 时）

```
⏰ 触发 v2 的条件:
1. 重大架构调整（如引入微服务）
2. 业务模型变更（如用户体系重构）
3. 技术升级（如迁移到 GraphQL）
4. v1 设计缺陷累积过多

📅 建议时间点:
- 项目规模 > 10万用户
- API 调用 > 1000 QPS
- 团队规模 > 10人
- 累积技术债务严重
```

---

## 🎯 总结

### 核心原则

1. **保持简单**: 不要过早引入复杂的版本管理
2. **向后兼容**: 尽量通过兼容性变更解决问题
3. **提前规划**: 在 v1 中留好扩展空间
4. **清晰沟通**: 提前通知废弃和下线计划

### 当前建议

```
✅ 短期 (3-6个月):
- 保持 v1，持续迭代
- 添加向后兼容的新功能
- 完善文档和测试

⏰ 中期 (6-12个月):
- 评估是否需要 v2
- 如需要，开始设计 v2 API
- 准备迁移指南

🚀 长期 (12个月+):
- 发布 v2（如果需要）
- v1 进入废弃期
- 逐步迁移用户到 v2
```

---

**记住**: 版本管理的目标是**平衡创新和稳定性**，而不是追求版本号的增长。

**当前项目**: 保持 v1，专注于完善功能和性能，等业务需求明确后再考虑 v2。
