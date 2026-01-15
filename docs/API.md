# API 文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **响应格式**: JSON

---

## 响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": 10001,
  "message": "参数错误",
  "error": "详细错误信息（仅开发环境）"
}
```

### 状态码说明

| Code | 说明 |
|------|------|
| 0 | 成功 |
| 10001 | 参数错误 |
| 10002 | 未授权 |
| 10003 | 禁止访问 |
| 10004 | 资源不存在 |
| 10005 | 资源已存在 |
| 10006 | 内部错误 |
| 10007 | 密码错误 |

---

## 接口列表

### 1. 健康检查

**接口地址**: `GET /health`

**描述**: 检查服务健康状态

**请求参数**: 无

**响应示例**:

```json
{
  "status": "ok",
  "time": 1704096000
}
```

**curl 示例**:

```bash
curl http://localhost:8080/health
```

---

### 2. 用户注册

**接口地址**: `POST /api/v1/users/register`

**描述**: 注册新用户

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名（3-50字符） |
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码（6-50字符） |

**请求示例**:

```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "password123"
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "avatar": "",
    "status": 1,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

**错误响应**:

```json
{
  "code": 10005,
  "message": "用户已存在"
}
```

**curl 示例**:

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123"
  }'
```

---

### 3. 用户登录

**接口地址**: `POST /api/v1/users/login`

**描述**: 用户登录

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| email | string | 是 | 邮箱地址 |
| password | string | 是 | 密码 |

**请求示例**:

```json
{
  "email": "alice@example.com",
  "password": "password123"
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "avatar": "",
    "status": 1,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

**错误响应**:

```json
// 用户不存在
{
  "code": 10004,
  "message": "用户不存在"
}

// 密码错误
{
  "code": 10007,
  "message": "密码错误"
}
```

**curl 示例**:

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "password123"
  }'
```

---

### 4. 获取用户信息

**接口地址**: `GET /api/v1/users/:id`

**描述**: 根据用户 ID 获取用户信息

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int64 | 是 | 用户 ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "alice",
    "email": "alice@example.com",
    "avatar": "https://example.com/avatar.jpg",
    "status": 1,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

**错误响应**:

```json
{
  "code": 10004,
  "message": "用户不存在"
}
```

**curl 示例**:

```bash
curl http://localhost:8080/api/v1/users/1
```

---

### 5. 更新用户信息

**接口地址**: `PUT /api/v1/users/:id`

**描述**: 更新用户信息

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int64 | 是 | 用户 ID |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| username | string | 是 | 用户名（3-50字符） |
| email | string | 是 | 邮箱地址 |
| avatar | string | 否 | 头像 URL |

**请求示例**:

```json
{
  "username": "alice_new",
  "email": "alice_new@example.com",
  "avatar": "https://example.com/new_avatar.jpg"
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

**错误响应**:

```json
// 用户不存在
{
  "code": 10004,
  "message": "用户不存在"
}

// 邮箱已被占用
{
  "code": 10005,
  "message": "邮箱已被占用"
}
```

**curl 示例**:

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice_new",
    "email": "alice_new@example.com",
    "avatar": "https://example.com/new_avatar.jpg"
  }'
```

---

### 6. 修改密码

**接口地址**: `PUT /api/v1/users/:id/password`

**描述**: 修改用户密码

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int64 | 是 | 用户 ID |

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| old_password | string | 是 | 旧密码 |
| new_password | string | 是 | 新密码（6-50字符） |

**请求示例**:

```json
{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

**错误响应**:

```json
// 旧密码错误
{
  "code": 10007,
  "message": "旧密码错误"
}
```

**curl 示例**:

```bash
curl -X PUT http://localhost:8080/api/v1/users/1/password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword456"
  }'
```

---

### 7. 删除用户

**接口地址**: `DELETE /api/v1/users/:id`

**描述**: 删除用户（软删除）

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | int64 | 是 | 用户 ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

**错误响应**:

```json
{
  "code": 10004,
  "message": "用户不存在"
}
```

**curl 示例**:

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

---

### 8. 用户列表

**接口地址**: `GET /api/v1/users`

**描述**: 获取用户列表（分页）

**查询参数**:

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| size | int | 否 | 10 | 每页数量（1-100） |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "alice",
        "email": "alice@example.com",
        "avatar": "",
        "status": 1,
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z"
      },
      {
        "id": 2,
        "username": "bob",
        "email": "bob@example.com",
        "avatar": "",
        "status": 1,
        "created_at": "2024-01-01T11:00:00Z",
        "updated_at": "2024-01-01T11:00:00Z"
      }
    ],
    "page": 1,
    "size": 10,
    "total": 2
  }
}
```

**curl 示例**:

```bash
# 第一页，每页 10 条
curl http://localhost:8080/api/v1/users?page=1&size=10

# 第二页，每页 20 条
curl http://localhost:8080/api/v1/users?page=2&size=20
```

---

## 错误处理

### HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

### 业务错误码

所有业务错误都返回 HTTP 200，通过响应体中的 `code` 字段区分：

```json
{
  "code": 10001,  // 业务错误码
  "message": "参数错误",
  "error": "username is required"  // 详细错误（仅开发环境）
}
```

---

## 限流说明

### 限流策略

- **算法**: 令牌桶（Token Bucket）
- **限制**: 每秒 100 个请求
- **桶容量**: 200 个令牌
- **Key**: 客户端 IP

### 超限响应

```json
{
  "code": 429,
  "message": "请求过于频繁，请稍后再试"
}
```

HTTP 状态码: `429 Too Many Requests`

---

## 日志追踪

### Request ID

每个请求都会自动生成一个唯一的 Request ID，用于日志追踪。

**响应头**:

```
X-Request-ID: 01234567-89ab-cdef-0123-456789abcdef
```

**日志示例**:

```json
{
  "time": "2024-01-01 10:00:00",
  "level": "INFO",
  "msg": "Request completed",
  "request_id": "01234567-89ab-cdef-0123-456789abcdef",
  "method": "GET",
  "path": "/api/v1/users/1",
  "status": 200,
  "latency": "10ms"
}
```

---

## 测试示例

### Postman Collection

**导入 URL**: (待补充)

### 完整测试流程

```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "test123456"
  }'

# 2. 登录
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "test123456"
  }'

# 3. 获取用户信息
curl http://localhost:8080/api/v1/users/1

# 4. 更新用户信息
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser_updated",
    "email": "test_updated@example.com",
    "avatar": "https://example.com/avatar.jpg"
  }'

# 5. 修改密码
curl -X PUT http://localhost:8080/api/v1/users/1/password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "test123456",
    "new_password": "newpassword123"
  }'

# 6. 获取用户列表
curl http://localhost:8080/api/v1/users?page=1&size=10

# 7. 删除用户
curl -X DELETE http://localhost:8080/api/v1/users/1
```

---

## 性能指标

### 预期性能

- **QPS**: 5000+ (单机)
- **响应时间**: 
  - P50: < 10ms
  - P95: < 50ms
  - P99: < 100ms
- **缓存命中率**: > 90%

### 性能测试

使用 Apache Bench (ab) 进行压测：

```bash
# 测试注册接口
ab -n 1000 -c 100 -p register.json -T application/json \
  http://localhost:8080/api/v1/users/register

# 测试查询接口
ab -n 10000 -c 100 \
  http://localhost:8080/api/v1/users/1
```

---

## 最佳实践

### 1. 参数验证

客户端应进行基本的参数验证：

```javascript
// ✅ 正确
if (username.length < 3 || username.length > 50) {
  alert("用户名长度为 3-50 字符");
  return;
}

// ❌ 错误：直接发送请求，依赖服务端验证
fetch('/api/v1/users/register', { ... })
```

### 2. 错误处理

客户端应根据业务错误码进行友好提示：

```javascript
const response = await fetch('/api/v1/users/login', { ... });
const data = await response.json();

if (data.code !== 0) {
  switch (data.code) {
    case 10004:
      alert("用户不存在");
      break;
    case 10007:
      alert("密码错误");
      break;
    default:
      alert(data.message);
  }
}
```

### 3. 限流处理

遇到 429 错误时，应进行退避重试：

```javascript
async function requestWithRetry(url, options, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    const response = await fetch(url, options);
    
    if (response.status !== 429) {
      return response;
    }
    
    // 指数退避
    await sleep(Math.pow(2, i) * 1000);
  }
  
  throw new Error("请求过于频繁");
}
```

---

## 常见问题

### Q1: 如何处理超时？

**A**: 默认超时时间为 10 秒。客户端应设置合理的超时时间：

```javascript
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), 5000);

fetch(url, {
  signal: controller.signal
}).finally(() => clearTimeout(timeoutId));
```

### Q2: 如何追踪请求？

**A**: 使用响应头中的 `X-Request-ID` 进行日志追踪。

### Q3: 分页最大限制是多少？

**A**: 单页最大 100 条记录。

---

## 更新日志

### v1.0.0 (2024-01-01)

- ✅ 用户注册、登录
- ✅ 用户 CRUD
- ✅ 密码修改
- ✅ 用户列表（分页）

---

## 联系方式

- **GitHub Issues**: https://github.com/yourusername/gin-demo/issues
- **Email**: your-email@example.com

---

**Happy Coding! 🎉**
