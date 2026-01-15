# Swagger API 文档

本项目使用 **swaggo/swag** 生成 Swagger API 文档。

## 快速开始

### 1. 安装 swag 工具

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. 添加 API 注释

在 `main.go` 中添加通用信息：

```go
// @title           Gin Demo API
// @version         1.0
// @description     RESTful API 骨架项目，基于 Gin + PostgreSQL + Redis
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
```

在 Handler 方法中添加详细注释：

```go
// Register godoc
// @Summary      用户注册
// @Description  创建新用户账号
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "注册信息"
// @Success      200      {object}  response.Response{data=UserResponse}
// @Failure      400      {object}  response.Response
// @Router       /users/register [post]
func (h *Handler) Register(c *gin.Context) {
    // ...
}
```

### 3. 生成文档

```bash
# 生成 Swagger 文档
swag init

# 格式化注释
swag fmt
```

### 4. 集成到项目

在 `main.go` 中添加：

```go
import (
    _ "gin_demo/docs" // Swagger docs
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// 添加 Swagger 路由
engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

### 5. 访问文档

```bash
# 启动服务后访问
http://localhost:8080/swagger/index.html
```

## 常用注释标签

### 通用注释（main.go）

| 标签 | 说明 | 示例 |
|------|------|------|
| `@title` | API 标题 | `@title Gin Demo API` |
| `@version` | API 版本 | `@version 1.0` |
| `@description` | API 描述 | `@description 这是一个示例API` |
| `@host` | 服务器地址 | `@host localhost:8080` |
| `@BasePath` | 基础路径 | `@BasePath /api/v1` |
| `@securityDefinitions` | 安全定义 | 见上方示例 |

### 路由注释（handler）

| 标签 | 说明 | 示例 |
|------|------|------|
| `@Summary` | 简短摘要 | `@Summary 获取用户信息` |
| `@Description` | 详细描述 | `@Description 根据ID获取用户详情` |
| `@Tags` | 分组标签 | `@Tags users` |
| `@Accept` | 接受格式 | `@Accept json` |
| `@Produce` | 返回格式 | `@Produce json` |
| `@Param` | 参数定义 | `@Param id path int true "用户ID"` |
| `@Success` | 成功响应 | `@Success 200 {object} Response` |
| `@Failure` | 失败响应 | `@Failure 400 {object} ErrorResponse` |
| `@Router` | 路由路径 | `@Router /users/{id} [get]` |
| `@Security` | 安全要求 | `@Security BearerAuth` |

### 参数类型

- `path` - 路径参数（/:id）
- `query` - 查询参数（?name=xxx）
- `header` - 请求头
- `body` - 请求体
- `formData` - 表单数据

## 完整示例

```go
// ListUsers godoc
// @Summary      用户列表
// @Description  获取用户列表（分页）
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page      query     int  false  "页码"  default(1)
// @Param        page_size query     int  false  "每页数量"  default(20)
// @Success      200       {object}  response.Response{data=response.ListResponse}
// @Failure      401       {object}  response.Response
// @Router       /users [get]
// @Security     BearerAuth
func (h *Handler) ListUsers(c *gin.Context) {
    // 实现...
}
```

## 参考资料

- [Swag 官方文档](https://github.com/swaggo/swag)
- [Swagger 注释规范](https://swagger.io/docs/specification/about/)
- [示例项目](https://github.com/swaggo/swag/tree/master/example)

## 注意事项

1. **结构体注释** - 确保所有 DTO 结构体都有注释
2. **类型匹配** - 确保 Swagger 注释中的类型与实际代码一致
3. **定期更新** - 修改 API 后记得重新生成文档
4. **安全性** - 生产环境可考虑关闭 Swagger UI 或添加认证
