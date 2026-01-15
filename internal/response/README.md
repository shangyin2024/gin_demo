# response - API 响应处理

> 统一的 API 响应格式、错误码定义和分页工具

---

## 📦 包含内容

### 1. 错误码 (errors.go)

项目的业务错误码定义：

```go
const (
    CodeOK              Code = 0      // 成功
    
    // 客户端错误 (10xxx)
    CodeInvalidParams   Code = 10001  // 参数错误
    CodeUnauthorized    Code = 10002  // 未授权
    CodeForbidden       Code = 10003  // 禁止访问
    CodeNotFound        Code = 10004  // 资源不存在
    CodeAlreadyExists   Code = 10005  // 资源已存在
    CodeTooManyRequests Code = 10006  // 请求过于频繁
    CodeInvalidPassword Code = 10007  // 密码错误
    
    // 服务端错误 (50xxx)
    CodeInternalError   Code = 50001  // 内部错误
    CodeDatabaseError   Code = 50002  // 数据库错误
    CodeCacheError      Code = 50003  // 缓存错误
)
```

**快捷函数**:

```go
// 创建错误
err := response.New(response.CodeNotFound, "用户不存在")

// 包装错误
err := response.Wrap(dbErr, response.CodeDatabaseError, "数据库查询失败")

// 预定义错误
err := response.ErrUnauthorized
```

---

### 2. 响应格式 (response.go)

统一的 JSON 响应结构：

```go
type Response struct {
    Code    Code   `json:"code"`            // 业务状态码
    Message string `json:"message"`         // 提示信息
    Data    any    `json:"data,omitempty"`  // 响应数据
    Error   string `json:"error,omitempty"` // 错误详情（仅开发环境）
}
```

**响应函数**:

```go
// 成功响应
response.Success(c, user)
// {"code":0,"message":"success","data":{...}}

// 错误响应（自动映射 HTTP 状态码）
response.Error(c, err)
// {"code":10004,"message":"资源不存在"}

// 指定错误码
response.ErrorWithCode(c, response.CodeUnauthorized, "请先登录")
// {"code":10002,"message":"请先登录"}
```

**HTTP 状态码映射**:

| 错误码 | HTTP 状态 |
|--------|-----------|
| 0 | 200 OK |
| 10001 | 400 Bad Request |
| 10002 | 401 Unauthorized |
| 10003 | 403 Forbidden |
| 10004 | 404 Not Found |
| 10005 | 409 Conflict |
| 10006 | 429 Too Many Requests |
| 50001+ | 500 Internal Server Error |

---

### 3. 分页 (pagination.go)

统一的分页处理：

```go
// 分页请求
type PaginationRequest struct {
    Page     int // 页码（从 1 开始）
    PageSize int // 每页数量（默认 10，最大 100）
}

// 分页响应
type PaginationResponse struct {
    Page       int   // 当前页码
    PageSize   int   // 每页数量
    Total      int64 // 总记录数
    TotalPages int   // 总页数
}

// 列表响应
type ListResponse[T any] struct {
    Items      []T
    Pagination PaginationResponse
}
```

**使用示例**:

```go
// 获取分页参数
pagination := response.GetPagination(c)
// GET /users?page=2&size=20

// 查询数据
users, total, _ := service.ListUsers(ctx, 
    pagination.GetLimit(),   // 20
    pagination.GetOffset())  // 20

// 返回响应
resp := response.NewListResponse(users, 
    response.NewPaginationResponse(pagination.Page, pagination.PageSize, total))
response.Success(c, resp)
```

**响应格式**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "pagination": {
      "page": 2,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

---

## 🎯 使用示例

### Handler 完整示例

```go
package user

import (
    "gin_demo/internal/response"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    service UserService
}

// 获取用户（单个）
func (h *Handler) GetUser(c *gin.Context) {
    userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.Error(c, response.NewWithError(
            response.CodeInvalidParams, "无效的用户ID", err))
        return
    }
    
    user, err := h.service.GetUser(c.Request.Context(), userID)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    response.Success(c, user)
}

// 创建用户
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, response.NewWithError(
            response.CodeInvalidParams, "参数错误", err))
        return
    }
    
    user, err := h.service.CreateUser(c.Request.Context(), req)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    response.Success(c, user)
}

// 获取用户列表（带分页）
func (h *Handler) ListUsers(c *gin.Context) {
    // 获取分页参数
    pagination := response.GetPagination(c)
    
    // 查询数据
    users, total, err := h.service.ListUsers(
        c.Request.Context(),
        pagination.GetLimit(),
        pagination.GetOffset(),
    )
    if err != nil {
        response.Error(c, err)
        return
    }
    
    // 构建响应
    resp := response.NewListResponse(users, 
        response.NewPaginationResponse(
            pagination.Page, 
            pagination.PageSize, 
            total))
    
    response.Success(c, resp)
}
```

### Service 层使用错误码

```go
package service

import "gin_demo/internal/response"

func (s *UserService) GetUser(ctx context.Context, id int64) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, response.ErrNotFound
        }
        return nil, response.Wrap(err, response.CodeDatabaseError, "查询用户失败")
    }
    return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    // 检查用户是否已存在
    exists, _ := s.repo.ExistsByEmail(ctx, input.Email)
    if exists {
        return nil, response.ErrAlreadyExists
    }
    
    // 创建用户
    user, err := s.repo.Create(ctx, input)
    if err != nil {
        return nil, response.Wrap(err, response.CodeDatabaseError, "创建用户失败")
    }
    
    return user, nil
}
```

---

## 📐 设计原则

### 为什么在 internal/response？

1. **项目特定**: 响应格式和错误码是项目定义的，不是通用的
2. **紧密耦合**: 错误码、响应格式、HTTP 映射紧密关联
3. **业务相关**: 与项目的 API 设计直接相关

### 统一管理的好处

- ✅ 所有 API 返回格式一致
- ✅ 错误码集中管理，避免冲突
- ✅ HTTP 状态码映射统一
- ✅ 前端对接更简单
- ✅ 便于添加全局日志、监控

---

## ✅ 总结

`internal/response` 包提供了完整的 API 响应解决方案：

- **错误码**: 统一的业务错误码定义
- **响应格式**: 标准化的 JSON 响应结构
- **分页工具**: 开箱即用的分页功能
- **HTTP 映射**: 自动的状态码映射

**让 API 开发更简单、更一致！** 🚀
