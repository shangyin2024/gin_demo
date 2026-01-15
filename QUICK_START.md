# 🚀 快速开始指南

**5 分钟启动你的 Go Web 项目！**

---

## ⚡ 超快开始（推荐新手）

### 1. 下载项目

```bash
# 方式 1: 在 GitHub 上点击 "Use this template" 按钮（推荐）

# 方式 2: 克隆项目
git clone https://github.com/yourusername/go-web-scaffold.git my-project
cd my-project
```

### 2. 重命名模块 ⚠️ **必须执行！**

```bash
# 让脚本可执行
chmod +x scripts/rename-module.sh

# 运行重命名脚本（替换为你的模块名）
./scripts/rename-module.sh github.com/yourname/my-project

# 示例:
./scripts/rename-module.sh github.com/acme/awesome-api
```

**为什么必须重命名？**
- Go 使用模块路径管理依赖
- 不重命名会导致导入路径错误
- 脚本会自动更新所有相关文件

### 3. 一键启动

```bash
# 安装工具 + 启动环境 + 运行应用
make init && make run
```

### 4. 验证

```bash
# 健康检查
curl http://localhost:8080/health

# 注册用户
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"password123"}'
```

✅ **完成！你的 API 已经运行了！**

---

## 📖 详细步骤说明

### 步骤 1: 获取代码

#### 方式 A: 使用 GitHub 模板（推荐）

1. 访问 https://github.com/yourusername/go-web-scaffold
2. 点击绿色的 **"Use this template"** 按钮
3. 创建你的新仓库
4. 克隆到本地

#### 方式 B: 直接克隆

```bash
git clone https://github.com/yourusername/go-web-scaffold.git my-awesome-api
cd my-awesome-api
```

---

### 步骤 2: 重命名项目 ⚠️ 必须！

**为什么需要重命名？**

Go 的导入路径基于模块名。如果不重命名，你的项目会使用 `gin_demo`，导致：
- 导入路径混乱
- 无法正确引用自己的包
- 与其他项目冲突

**自动重命名（推荐）：**

```bash
# 1. 给脚本执行权限
chmod +x scripts/rename-module.sh

# 2. 运行脚本
./scripts/rename-module.sh github.com/yourname/yourproject

# 实际示例:
./scripts/rename-module.sh github.com/acme/awesome-api
```

脚本会自动：
- ✅ 更新 `go.mod` 模块名
- ✅ 替换所有 Go 文件的导入路径
- ✅ 更新 `Makefile` 和文档
- ✅ 运行 `go mod tidy`
- ✅ 验证编译

**手动重命名（了解原理）：**

```bash
# 1. 编辑 go.mod
vim go.mod
# 将第一行改为: module github.com/yourname/yourproject

# 2. 批量替换导入路径
find . -type f -name "*.go" ! -path "./vendor/*" \
  -exec sed -i '' 's/gin_demo/yourproject/g' {} +

# 3. 整理依赖
go mod tidy

# 4. 验证
go build ./...
```

---

### 步骤 3: 配置环境

#### 快速配置（使用默认值）

```bash
# 复制环境变量模板
cp .env.example .env

# 使用默认配置（开发环境已配置好）
# MySQL: localhost:3306
# Redis: localhost:6379
```

#### 自定义配置（可选）

```bash
# 编辑配置文件
vim config.yaml

# 或通过环境变量覆盖
export DATABASE_PASSWORD=your_password
export JWT_SECRET=your-secret-key
```

---

### 步骤 4: 启动服务

#### 方式 A: 一键启动（推荐）

```bash
# 安装工具 + 启动 Docker + 运行迁移 + 启动应用
make init && make run
```

这个命令会：
1. 安装 `sqlc`, `wire`, `golangci-lint` 等工具
2. 启动 MySQL 和 Redis (Docker)
3. 执行数据库迁移
4. 生成代码 (sqlc + wire)
5. 运行应用

#### 方式 B: 分步启动（理解流程）

```bash
# 1. 安装开发工具
make tools

# 2. 启动 MySQL 和 Redis
make dev

# 3. 执行数据库迁移
make migrate-up

# 4. 生成代码
make generate

# 5. 运行应用
make run
```

---

### 步骤 5: 测试 API

#### 健康检查

```bash
curl http://localhost:8080/health
```

预期响应：
```json
{
  "status": "healthy",
  "database": "healthy",
  "redis": "healthy"
}
```

#### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123"
  }'
```

#### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "password123"
  }'
```

保存返回的 `token`，用于后续认证请求。

#### 获取用户信息（需要认证）

```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 🐛 常见问题

### Q1: 启动失败，端口被占用

```bash
# 检查端口占用
lsof -i :8080    # 应用端口
lsof -i :3306    # MySQL
lsof -i :6379    # Redis

# 修改端口（编辑 config.yaml）
vim config.yaml
# server.port: 8081  # 改为其他端口
```

### Q2: 数据库连接失败

```bash
# 检查 Docker 服务
docker-compose ps

# 查看 MySQL 日志
docker-compose logs mysql-master

# 重启服务
make dev-stop && make dev
```

### Q3: Redis 连接失败

```bash
# 检查 Redis 状态
docker-compose ps redis-master

# 查看 Redis 日志
docker-compose logs redis-master

# 测试连接
docker-compose exec redis-master redis-cli ping
```

### Q4: make 命令不存在

```bash
# macOS
brew install make

# Ubuntu/Debian
sudo apt-get install build-essential

# 或手动执行命令（查看 Makefile）
```

### Q5: sqlc 或 wire 未安装

```bash
# 自动安装所有工具
make tools

# 或手动安装
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/google/wire/cmd/wire@latest
```

### Q6: 编译失败

```bash
# 清理并重新构建
make clean
go mod tidy
make generate
make build
```

---

## 📚 下一步

### 学习项目结构

```bash
# 查看项目结构文档
cat docs/ARCHITECTURE.md

# 查看配置说明
cat docs/CONFIGURATION.md

# 查看 API 文档
cat docs/API.md
```

### 添加新功能

```bash
# 查看模板使用指南
cat TEMPLATE_USAGE.md

# 查看贡献指南
cat CONTRIBUTING.md
```

### 开发新模块

1. 创建数据库表（`db/migrations/`）
2. 定义 SQL 查询（`db/queries/`）
3. 生成代码（`make sqlc`）
4. 实现 Repository、Service、Handler
5. 注册到 Wire（`internal/wire/`）
6. 生成依赖注入代码（`make wire`）
7. 添加路由

详见：[TEMPLATE_USAGE.md](TEMPLATE_USAGE.md)

---

## 🎯 常用命令速查

```bash
# 开发
make dev              # 启动开发环境
make run              # 运行应用
make dev-stop         # 停止开发环境

# 代码生成
make generate         # 生成所有代码
make sqlc             # 生成数据库代码
make wire             # 生成依赖注入代码

# 测试
make test             # 运行测试
make test-cover       # 测试 + 覆盖率
make bench            # 性能测试

# 代码质量
make lint             # 代码检查
make fmt              # 格式化
make check            # 完整检查

# 数据库
make migrate-up       # 执行迁移
make migrate-down     # 回滚迁移
make db-console       # 数据库控制台

# 构建
make build            # 编译应用
make docker-build     # 构建 Docker 镜像

# 清理
make clean            # 清理构建产物
```

---

## 💡 提示

1. **⚠️ 必须先重命名模块**，否则导入路径会出错
2. **首次启动使用 `make init`**，后续可以直接 `make run`
3. **配置文件支持多环境**：`config.dev.yaml`, `config.prod.yaml`
4. **环境变量优先级最高**，可以覆盖配置文件
5. **查看完整文档索引**：`docs/INDEX.md`

---

## 🆘 获取帮助

- 📖 查看完整文档：`docs/INDEX.md`
- 🐛 报告问题：GitHub Issues
- 💬 讨论交流：GitHub Discussions
- 📚 详细指南：`TEMPLATE_USAGE.md`
- 🤝 贡献代码：`CONTRIBUTING.md`

---

**🎉 开始构建你的应用吧！**

```bash
# 记住这三步！
./scripts/rename-module.sh github.com/yourname/yourproject
make init
make run
```

Happy Coding! 🚀
