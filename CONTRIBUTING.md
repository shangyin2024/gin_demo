# 🤝 贡献指南

感谢你对这个项目感兴趣！我们欢迎任何形式的贡献。

---

## 📋 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发流程](#开发流程)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [Pull Request 流程](#pull-request-流程)

---

## 行为准则

我们致力于为所有参与者提供一个友好、安全和热情的环境。请遵守以下原则：

- 使用友善和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表现出同理心

---

## 如何贡献

### 🐛 报告 Bug

1. 在 [Issues](https://github.com/yourusername/go-web-scaffold/issues) 页面搜索，确认问题未被报告
2. 使用 Bug Report 模板创建新 Issue
3. 提供详细的复现步骤和环境信息
4. 如可能，提供最小化的复现示例

### ✨ 提出新功能

1. 先在 [Discussions](https://github.com/yourusername/go-web-scaffold/discussions) 讨论
2. 获得反馈后，使用 Feature Request 模板创建 Issue
3. 详细说明功能的用途和价值
4. 等待维护者反馈

### 📚 改进文档

文档改进非常受欢迎！包括但不限于：

- 修正拼写或语法错误
- 补充缺失的说明
- 添加使用示例
- 翻译文档

---

## 开发流程

### 1. Fork 项目

```bash
# 1. 在 GitHub 上点击 Fork 按钮

# 2. 克隆你的 Fork
git clone https://github.com/YOUR_USERNAME/go-web-scaffold.git
cd go-web-scaffold

# 3. 添加上游仓库
git remote add upstream https://github.com/ORIGINAL_OWNER/go-web-scaffold.git
```

### 2. 创建分支

```bash
# 从 main 分支创建新分支
git checkout -b feature/your-feature-name

# 分支命名规范:
# - feature/xxx   新功能
# - bugfix/xxx    Bug 修复
# - docs/xxx      文档更新
# - refactor/xxx  代码重构
# - test/xxx      测试相关
```

### 3. 设置开发环境

```bash
# 安装依赖
make tools

# 启动开发环境
make dev

# 运行测试确保环境正常
make test
```

### 4. 开发和测试

```bash
# 编写代码...

# 生成代码（如修改了 SQL 或 Wire）
make generate

# 运行测试
make test

# 检查代码质量
make check

# 本地运行验证
make run
```

### 5. 提交代码

```bash
# 添加变更
git add .

# 提交（遵循提交规范）
git commit -m "feat: add awesome feature"

# 推送到你的 Fork
git push origin feature/your-feature-name
```

### 6. 创建 Pull Request

1. 在 GitHub 上打开你的 Fork
2. 点击 "Compare & pull request"
3. 填写 PR 描述（使用模板）
4. 等待代码审查

---

## 代码规范

### Go 代码风格

遵循官方 Go 代码规范：

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

### 关键原则

#### 1. 命名规范

```go
// ✅ 好的命名
func GetUserByID(ctx context.Context, id int64) (*User, error)
var userCache *cache.Manager
const MaxRetryCount = 3

// ❌ 不好的命名
func get_user(id int64) *User  // 使用驼峰而非下划线
var cache1 *cache.Manager       // 避免无意义的数字后缀
const max_retry = 3             // 常量使用大驼峰
```

#### 2. 错误处理

```go
// ✅ 正确的错误处理
user, err := service.GetUser(ctx, id)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// ❌ 不处理错误
user, _ := service.GetUser(ctx, id)  // 永远不要忽略错误
```

#### 3. 上下文传递

```go
// ✅ 第一个参数传递 context
func GetUser(ctx context.Context, id int64) (*User, error)

// ❌ 不传递或放在其他位置
func GetUser(id int64) (*User, error)
func GetUser(id int64, ctx context.Context) (*User, error)
```

#### 4. 接口设计

```go
// ✅ 小而专注的接口
type UserRepository interface {
    GetByID(ctx context.Context, id int64) (*User, error)
    Create(ctx context.Context, user *User) error
}

// ❌ 过大的接口
type Repository interface {
    GetByID(ctx context.Context, id int64) (interface{}, error)
    Create(ctx context.Context, data interface{}) error
    Update(ctx context.Context, data interface{}) error
    Delete(ctx context.Context, id int64) error
    // ... 10+ 个方法
}
```

### 项目特定规范

#### 三层架构

```
Handler  → 处理 HTTP 请求，参数验证
Service  → 实现业务逻辑，权限校验
Repository → 数据访问，缓存管理
```

#### 依赖注入

```go
// 使用 Wire 进行依赖注入
// 1. 在 internal/wire/ 中定义 Provider
// 2. 添加到相应的 Wire Set
// 3. 运行 make wire 生成代码
```

#### 数据库访问

```go
// 使用 sqlc 生成类型安全的数据库代码
// 1. 在 db/queries/ 中定义 SQL
// 2. 运行 make sqlc 生成代码
// 3. 在 Repository 层使用生成的代码
```

---

## 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

### 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档变更
- `style`: 代码格式（不影响代码运行）
- `refactor`: 重构（既不是新增功能，也不是修复 bug）
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动
- `revert`: 回退

### 示例

```bash
# 新功能
git commit -m "feat(user): add user registration API"

# Bug 修复
git commit -m "fix(auth): resolve JWT token expiration issue"

# 文档
git commit -m "docs(readme): update installation instructions"

# 重构
git commit -m "refactor(repository): improve cache strategy"

# 性能优化
git commit -m "perf(database): optimize user query with index"

# 测试
git commit -m "test(service): add unit tests for user service"
```

### 完整示例

```
feat(user): add password reset functionality

- Implement password reset token generation
- Add email notification for password reset
- Create reset password API endpoint

Closes #123
```

---

## Pull Request 流程

### PR 检查清单

提交 PR 前，请确认：

- [ ] 代码遵循项目编码规范
- [ ] 已执行 `make lint` 并通过
- [ ] 已添加/更新相关测试
- [ ] 所有测试通过 (`make test`)
- [ ] 已更新相关文档
- [ ] 代码已自测
- [ ] PR 标题清晰明确
- [ ] 填写了 PR 描述模板

### 代码审查

所有 PR 需要经过代码审查：

1. **自动检查**: CI 会自动运行测试和代码检查
2. **人工审查**: 维护者会审查代码并提供反馈
3. **修改和更新**: 根据反馈修改代码
4. **合并**: 审查通过后，维护者会合并 PR

### 审查标准

代码审查关注：

- 功能正确性
- 代码质量和可读性
- 测试覆盖率
- 文档完整性
- 性能影响
- 安全性

---

## 测试要求

### 单元测试

```go
// 为新功能添加单元测试
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    service := setupTestService(t)
    user := &domain.User{
        Username: "testuser",
        Email:    "test@example.com",
    }

    // Act
    result, err := service.CreateUser(context.Background(), user)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, user.Username, result.Username)
}
```

### 集成测试

```go
// 为复杂场景添加集成测试
func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    db := setupTestDB(t)
    defer db.Close()
    
    // 测试完整的数据库操作流程
}
```

### 运行测试

```bash
# 单元测试
make test

# 包含集成测试
make test-all

# 查看覆盖率
make test-cover
```

---

## 文档规范

### 代码注释

```go
// GetUserByID retrieves a user by their unique identifier.
// It first checks the cache, then falls back to the database.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - id: User ID to retrieve
//
// Returns:
//   - *User: The user object if found
//   - error: Error if user not found or database error
func GetUserByID(ctx context.Context, id int64) (*User, error) {
    // Implementation
}
```

### Markdown 文档

- 使用清晰的标题层级
- 提供代码示例
- 添加链接引用
- 使用表格组织信息
- 包含截图（如适用）

---

## 版本发布

版本号遵循 [Semantic Versioning](https://semver.org/)：

```
MAJOR.MINOR.PATCH

例如: v1.2.3

- MAJOR: 不兼容的 API 变更
- MINOR: 向后兼容的功能新增
- PATCH: 向后兼容的问题修正
```

---

## 常见问题

### Q: 我的 PR 什么时候会被审查？

A: 我们会尽快审查，通常在 1-3 个工作日内。复杂的 PR 可能需要更长时间。

### Q: 如何同步上游的更新？

```bash
# 获取上游更新
git fetch upstream

# 合并到本地分支
git checkout main
git merge upstream/main

# 推送到你的 Fork
git push origin main
```

### Q: 我的 PR 被拒绝了怎么办？

不要气馁！阅读反馈，理解原因，必要时可以：
- 修改代码重新提交
- 在 Issue/Discussion 中讨论
- 寻求帮助和建议

### Q: 可以同时提交多个功能吗？

建议每个 PR 只关注一个功能或修复。这样更容易审查和合并。

---

## 获取帮助

- 📖 查看 [文档](docs/INDEX.md)
- 💬 在 [Discussions](https://github.com/yourusername/go-web-scaffold/discussions) 提问
- 🐛 在 [Issues](https://github.com/yourusername/go-web-scaffold/issues) 报告问题

---

## 感谢

感谢你花时间为这个项目做出贡献！每一个贡献都很重要。🎉

---

**Happy Contributing! 🚀**
