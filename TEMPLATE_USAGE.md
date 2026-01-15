# 🚀 如何使用此模板

这是一个**生产级 Go Web API 脚手架**，开箱即用，适合快速启动新项目。

---

## 📦 模板特性

### ✅ 完整的技术栈
- **Web 框架**: Gin v1.11.0
- **数据库**: MySQL 8.0+ (支持主从) / PostgreSQL 15+
- **缓存**: Redis 7+ 哨兵模式 (1主+2从+3哨兵)
- **依赖注入**: Wire
- **SQL 生成**: sqlc
- **配置管理**: Viper (多环境)
- **监控**: Prometheus + pprof
- **容器化**: Docker + Docker Compose

### ✅ 生产级功能
- JWT 认证 + RBAC 权限控制
- 多级缓存 (防击穿/穿透/雪崩)
- 分布式定时任务
- 请求限流
- 安全防护 (CORS + HSTS + CSP)
- 结构化日志
- 健康检查
- 性能分析

### ✅ 完善的文档
- 51 个文档文件
- API、架构、配置、部署文档齐全
- 包含 6 个优化报告
- 提供故障排查手册

---

## 🎯 适用场景

✅ **REST API 服务**
✅ **微服务项目**
✅ **企业级应用后端**
✅ **SaaS 平台后端**
✅ **学习 Go Web 最佳实践**
✅ **团队项目脚手架**

---

## 🚀 快速开始（4 步启动）

### 方式 1: 使用 GitHub 模板（推荐）

#### 1️⃣ 创建你的项目

```bash
# 点击 GitHub 上的 "Use this template" 按钮
# 或使用命令行:
git clone https://github.com/yourusername/go-web-scaffold.git my-project
cd my-project
```

#### 2️⃣ 重命名模块（⚠️ 必须执行）

```bash
# 使用自动重命名脚本（推荐）
chmod +x scripts/rename-module.sh
./scripts/rename-module.sh github.com/yourname/my-project

# 示例:
./scripts/rename-module.sh github.com/mycompany/awesome-api

# 脚本会自动:
# - 更新 go.mod 模块名
# - 替换所有 Go 文件中的导入路径
# - 更新 Makefile 和文档
# - 运行 go mod tidy
```

**⚠️ 重要**: 此步骤必须在开始开发前完成，否则导入路径会出错！

#### 3️⃣ 初始化项目

```bash
# 自动初始化（安装工具、启动环境、执行迁移）
make init

# 或手动执行
make tools        # 安装开发工具
make dev          # 启动 Docker 环境
make migrate-up   # 执行数据库迁移
```

#### 4️⃣ 运行项目

```bash
make run
```

✅ 访问 http://localhost:8080/health 验证服务启动

---

### 方式 2: 完全手动

如果你想完全理解每一步：

```bash
# 1. 克隆项目
git clone https://github.com/yourusername/go-web-scaffold.git my-project
cd my-project

# 2. 手动重命名模块
# 编辑 go.mod，将第一行改为:
# module github.com/yourname/my-project

# 3. 批量替换导入路径
find . -type f -name "*.go" ! -path "./vendor/*" -exec sed -i '' 's/gin_demo/my-project/g' {} +

# 或使用更精确的替换（推荐）
grep -rl "gin_demo" --include="*.go" . | xargs sed -i '' 's|gin_demo|my-project|g'

# 4. 更新依赖
go mod tidy

# 5. 重新生成代码
make generate

# 6. 验证
make build

# 7. 启动
make init && make run
```

---

## 🔧 自定义项目

### 1. 修改模块名称

**自动方式**（推荐）:
```bash
chmod +x scripts/rename-module.sh
./scripts/rename-module.sh github.com/yourname/yourproject
```

**手动方式**:
```bash
# 1. 修改 go.mod 第一行
module github.com/yourname/yourproject

# 2. 批量替换所有导入路径
find . -type f -name "*.go" -exec sed -i '' 's/gin_demo/yourproject/g' {} +

# 3. 重新整理依赖
go mod tidy
```

### 2. 配置环境

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置
vim .env

# 或直接编辑配置文件
vim config.yaml        # 开发环境
vim config.prod.yaml   # 生产环境
```

### 3. 修改数据库表结构

```bash
# 1. 创建迁移文件
cat > db/migrations/002_your_migration.sql << 'EOF'
-- +migrate Up
CREATE TABLE your_table (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +migrate Down
DROP TABLE IF EXISTS your_table;
EOF

# 2. 执行迁移
make migrate-up
```

### 4. 添加新的 SQL 查询

```bash
# 1. 编辑 SQL 查询文件
vim db/queries/your_queries.sql

# 示例：
# -- name: GetItemByID :one
# SELECT * FROM items WHERE id = ? LIMIT 1;
#
# -- name: ListItems :many
# SELECT * FROM items ORDER BY created_at DESC LIMIT ? OFFSET ?;

# 2. 生成代码
make sqlc
```

### 5. 添加新的业务模块

```bash
# 1. 创建目录结构
mkdir -p internal/app/handler/item
mkdir -p internal/domain/service/item
mkdir -p internal/repository

# 2. 参考 user 模块实现
# - handler: 处理 HTTP 请求
# - service: 实现业务逻辑
# - repository: 数据访问

# 3. 在 wire 中注册
vim internal/wire/handler.go
vim internal/wire/service.go
vim internal/wire/repository.go

# 4. 重新生成依赖注入代码
make wire
```

---

## 📚 核心概念

### 三层架构

```
HTTP Request
    ↓
┌─────────────────────┐
│   Handler Layer     │  ← 路由、验证、响应
└─────────────────────┘
    ↓
┌─────────────────────┐
│   Service Layer     │  ← 业务逻辑、权限
└─────────────────────┘
    ↓
┌─────────────────────┐
│  Repository Layer   │  ← 数据访问、缓存
└─────────────────────┘
    ↓
Database / Cache
```

### 依赖注入流程

```go
// 1. 定义 Provider (internal/wire/xxx.go)
func provideUserService(repo repository.UserRepositoryInterface) *service.UserService {
    return service.NewUserService(repo)
}

// 2. 添加到 Wire 集合
var ServiceSet = wire.NewSet(
    provideUserService,
    // ... 其他 services
)

// 3. 运行 wire 生成代码
// make wire

// 4. 在 main.go 中使用
app := wire.InitializeApplication(cfg)
```

### 缓存策略

```go
// Repository 层自动处理缓存
func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*User, error) {
    // 1. 尝试从缓存获取
    // 2. 缓存未命中，查询数据库
    // 3. 写入缓存
    // 4. 返回结果
}

// 支持的缓存操作：
// - Get: 获取单个对象
// - MGet: 批量获取
// - Set: 设置缓存
// - Delete: 删除缓存
// - GetWithFallback: 缓存穿透保护
```

---

## 🛠️ 常用开发命令

### 开发环境

```bash
make dev              # 启动开发环境（Docker + MySQL + Redis）
make run              # 运行应用
make dev-stop         # 停止开发环境
make logs             # 查看日志
```

### 代码生成

```bash
make generate         # 生成所有代码（sqlc + wire）
make sqlc             # 生成数据库访问代码
make wire             # 生成依赖注入代码
```

### 数据库管理

```bash
make migrate-up       # 执行数据库迁移
make migrate-down     # 回滚迁移
make migrate-status   # 查看迁移状态
make db-console       # 进入数据库控制台
```

### 测试与质量

```bash
make test             # 运行测试
make test-cover       # 测试 + 覆盖率
make bench            # 性能基准测试
make lint             # 代码检查
make check            # 完整质量检查
```

### 性能分析

```bash
make pprof            # 启动 pprof 服务
make bench-cpu        # CPU 性能分析
make bench-mem        # 内存性能分析
```

### 构建与部署

```bash
make build            # 编译应用
make docker-build     # 构建 Docker 镜像
make docker-run       # 运行 Docker 容器
```

---

## 📁 关键文件说明

### 配置文件

| 文件 | 说明 |
|------|------|
| `config.yaml` | 默认配置 |
| `config.dev.yaml` | 开发环境配置 |
| `config.test.yaml` | 测试环境配置 |
| `config.prod.yaml` | 生产环境配置 |
| `.env.example` | 环境变量模板 |

### 代码生成

| 文件 | 说明 |
|------|------|
| `sqlc.yaml` | sqlc 配置 (SQL → Go 代码) |
| `dbconfig.yml` | 数据库迁移配置 |
| `internal/wire/wire.go` | Wire 依赖注入配置 |

### Docker

| 文件 | 说明 |
|------|------|
| `Dockerfile` | 应用镜像 (多阶段构建) |
| `docker-compose.yml` | 服务编排 (MySQL + Redis 哨兵) |

### 文档

| 目录/文件 | 说明 |
|----------|------|
| `docs/INDEX.md` | 文档索引（必读） |
| `docs/CONFIGURATION.md` | 配置详解 |
| `docs/API.md` | API 接口文档 |
| `docs/ARCHITECTURE.md` | 架构设计 |
| `docs/DEPLOYMENT-CHECKLIST.md` | 部署检查清单 |

---

## 🎯 典型使用场景

### 场景 1: 创建新的 REST API

```bash
# 1. 定义数据表
vim db/migrations/002_create_products.sql

# 2. 定义 SQL 查询
vim db/queries/products.sql

# 3. 生成代码
make generate

# 4. 实现业务逻辑
# - internal/repository/product_repository.go
# - internal/domain/service/product_service.go
# - internal/app/handler/product/handler.go

# 5. 注册路由
vim internal/app/router.go

# 6. 测试
make test
make run
```

### 场景 2: 添加定时任务

```go
// 1. 创建任务文件
// internal/task/tasks/my_task.go
package tasks

type MyTask struct {
    redis redis.UniversalClient
}

func (t *MyTask) Run() {
    // 任务逻辑
}

// 2. 注册任务
// internal/task/manager.go
func (m *Manager) registerTasks() {
    m.cron.AddFunc("0 */5 * * * *", m.tasks.myTask.Run)
}

// 3. Wire 注入
// internal/wire/task.go
```

### 场景 3: 实现 RBAC 权限控制

```go
// 已实现！直接使用：

// 1. 在路由上添加权限中间件
api.GET("/admin/users", 
    middleware.RequireRoles("admin"),
    handler.ListUsers)

// 2. 或检查权限
api.GET("/users/:id",
    middleware.RequirePermissions("user:read"),
    handler.GetUser)

// 详见: docs/RBAC.md
```

---

## 🔐 安全最佳实践

### 生产环境配置

```bash
# 1. 使用强密钥
export JWT_SECRET=$(openssl rand -base64 32)

# 2. 启用 HTTPS
export SERVER_TLS_ENABLED=true
export SERVER_TLS_CERT=/path/to/cert.pem
export SERVER_TLS_KEY=/path/to/key.pem

# 3. 限制 CORS
export CORS_ORIGINS=https://yourdomain.com

# 4. 关闭调试模式
export SERVER_MODE=release

# 5. 配置 Redis 密码
export REDIS_PASSWORD=your-strong-password
```

### 环境变量管理

```bash
# 开发环境
cp .env.example .env
vim .env

# 生产环境（不要提交 .env 到 Git）
# 使用 Kubernetes Secrets / AWS Secrets Manager 等
```

---

## 📊 性能调优

### 数据库优化

```yaml
database:
  max_open_conns: 25      # 根据负载调整
  max_idle_conns: 25
  conn_max_lifetime: 5m
  conn_max_idle_time: 10m
```

### Redis 优化

```yaml
redis:
  pool_size: 100          # 连接池大小
  min_idle_conns: 10      # 最小空闲连接
  max_retries: 3          # 重试次数
  sentinel_enabled: true  # 生产环境建议开启
```

### 应用优化

```bash
# 1. 查看性能指标
curl http://localhost:9090/metrics

# 2. 实时性能分析
make pprof
# 访问 http://localhost:6060/debug/pprof/

# 3. CPU 分析
go tool pprof http://localhost:6060/debug/pprof/profile

# 4. 内存分析
go tool pprof http://localhost:6060/debug/pprof/heap
```

---

## 🚀 部署建议

### Docker 部署（推荐）

```bash
# 1. 构建镜像
make docker-build

# 2. 推送到仓库
docker tag gin-demo:latest your-registry/gin-demo:v1.0.0
docker push your-registry/gin-demo:v1.0.0

# 3. 部署
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes 部署

```bash
# 参考 docs/K8S_DEPLOYMENT.md
kubectl apply -f k8s/
```

### 二进制部署

```bash
# 1. 编译
make build-linux

# 2. 上传到服务器
scp bin/gin-demo-linux-amd64 user@server:/opt/app/

# 3. 运行
./bin/gin-demo-linux-amd64
```

---

## 🐛 故障排查

### 常见问题

**1. 编译失败**
```bash
# 清理并重新构建
make clean
go mod tidy
make generate
make build
```

**2. 数据库连接失败**
```bash
# 检查配置
cat config.yaml | grep -A 10 database

# 测试连接
make db-console
```

**3. Redis 连接失败**
```bash
# 检查 Redis 状态
docker-compose ps redis-master

# 测试连接
make redis-console
```

**4. Wire 生成错误**
```bash
# 重新生成
rm internal/wire/wire_gen.go
make wire
```

详见: `docs/TROUBLESHOOTING.md`

---

## 📖 推荐阅读顺序

### 新手入门
1. `README.md` - 项目概览
2. `TEMPLATE_USAGE.md` (本文) - 模板使用
3. `docs/INDEX.md` - 文档索引
4. `docs/CONFIGURATION.md` - 配置说明

### 开发者
1. `docs/ARCHITECTURE.md` - 架构设计
2. `docs/API.md` - API 文档
3. `docs/DATABASE.md` - 数据库设计
4. `pkg/README.md` - 包设计原则

### 运维人员
1. `docs/DEPLOYMENT-CHECKLIST.md` - 部署清单
2. `docs/TROUBLESHOOTING.md` - 故障排查
3. `docs/MONITORING.md` - 监控告警

---

## 🎓 学习资源

### 官方文档
- [Gin](https://gin-gonic.com/)
- [Wire](https://github.com/google/wire)
- [sqlc](https://sqlc.dev/)
- [Viper](https://github.com/spf13/viper)

### 最佳实践
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide)

---

## 💡 提示与技巧

### 快速添加新模块

```bash
# 使用脚本（如果提供）
./scripts/create-module.sh item

# 手动创建
mkdir -p internal/app/handler/item
mkdir -p internal/domain/service/item
# 复制 user 模块作为模板
```

### 代码生成工作流

```bash
# 完整流程
1. 编辑 db/migrations/*.sql    # 修改表结构
2. 编辑 db/queries/*.sql       # 修改 SQL 查询
3. make migrate-up             # 执行迁移
4. make sqlc                   # 生成数据访问代码
5. 实现 business logic          # Service/Handler
6. 编辑 internal/wire/*.go     # 注册依赖
7. make wire                   # 生成注入代码
8. make test                   # 测试
9. make run                    # 运行
```

### 环境切换

```bash
# 开发环境
export APP_ENV=dev
make run

# 测试环境
export APP_ENV=test
make run

# 生产环境
export APP_ENV=prod
make run
```

---

## 🤝 贡献指南

如果你对这个模板有改进建议：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing`)
5. 创建 Pull Request

---

## ⭐ Star 历史

如果这个模板对你有帮助，请给个 Star！⭐

---

## 📞 获取帮助

- 📖 查看文档: `docs/INDEX.md`
- 🐛 报告问题: [GitHub Issues](https://github.com/yourusername/go-web-scaffold/issues)
- 💬 讨论交流: [GitHub Discussions](https://github.com/yourusername/go-web-scaffold/discussions)

---

## 📝 更新日志

查看 `CHANGELOG.md` 了解版本更新历史。

---

## 📜 许可证

MIT License - 详见 `LICENSE` 文件

---

**🎉 开始构建你的 Go Web 应用吧！**

```bash
make init && make run
```

Happy Coding! 🚀
