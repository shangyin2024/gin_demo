# Validator - 参数验证工具

通用的请求参数验证工具包，用于统一的参数验证逻辑。

## 功能

1. **ID 验证** - 验证路径参数中的 ID（必须 > 0）
2. **分页验证** - 验证并规范化分页参数
3. **偏移量计算** - 计算数据库分页偏移量

## 使用示例

### ID 验证

```go
import "gin_demo/pkg/validator"

func (h *Handler) GetUser(c *gin.Context) {
    // 验证 ID（自动检查 > 0）
    userID, valid := validator.ValidateID(c, "id")
    if !valid {
        response.Error(c, response.New(response.CodeInvalidParams, "无效的用户ID"))
        return
    }
    
    // 使用 userID...
}
```

### 分页验证

```go
// 验证并规范化分页参数
page, pageSize := validator.ValidatePagination(reqPage, reqPageSize)
// page: 默认1，最小1
// pageSize: 默认20，最小1，最大100

// 计算偏移量
offset := validator.CalculateOffset(page, pageSize)
```

## 设计原则

1. **简单实用** - 提供常用的验证函数
2. **统一标准** - 全项目使用相同的验证逻辑
3. **可扩展** - 可以根据需要添加新的验证函数
4. **零依赖** - 仅依赖 gin 和标准库

## 验证规则

### ID 规则
- 必须是整数
- 必须 > 0

### 分页规则
- `page`: 默认1，最小1
- `pageSize`: 默认20，最小1，最大100

## 未来扩展

可以添加更多验证函数：
- Email 验证
- Phone 验证
- UUID 验证
- 等等...
