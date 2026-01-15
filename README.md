# 🚀 Gin Web API Scaffold

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Gin Version](https://img.shields.io/badge/Gin-1.11.0-00ADD8?style=flat)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**生产级 Go Web API 脚手架**

一个功能完整、开箱即用的 Go Web API 项目模板，集成现代化技术栈和企业级最佳实践

[特性](#-特性) • [快速开始](#-快速开始) • [文档](#-文档) • [架构](#️-架构设计) • [贡献](#-贡献)

</div>

---

## ✨ 特性

### 🏗️ 核心架构
- ✅ **DDD 分层架构** - Handler → Service → Repository 清晰分层
- ✅ **依赖注入** - Google Wire 自动生成，类型安全
- ✅ **类型安全数据库** - sqlc 生成类型安全的 SQL 代码
- ✅ **泛型支持** - BaseRepository 泛型实现，减少重复代码
- ✅ **接口化设计** - 易于测试和替换实现

### 💾 数据层
- ✅ **MySQL 8.0+** - 主流关系型数据库（支持 PostgreSQL 切换）
- ✅ **Redis 哨兵模式** - 1主+2从+3哨兵高可用架构
- ✅ **智能缓存** - 三层防护（防击穿/穿透/雪崩）
- ✅ **数据库迁移** - sql-migrate 版本化管理
- ✅ **事务支持** - 统一事务管理

### 🔐 安全认证
- ✅ **JWT 认证** - 无状态 Token 认证
- ✅ **RBAC 权限** - 基于角色的访问控制
- ✅ **密码加密** - bcrypt 安全加密
- ✅ **安全中间件** - CORS、HSTS、CSP、X-Frame-Options

### 📊 监控运维
- ✅ **Prometheus 指标** - HTTP、业务、缓存、数据库多维度监控
- ✅ **pprof 性能分析** - CPU、内存、协程实时分析
- ✅ **健康检查** - 数据库、Redis、磁盘多项检查
- ✅ **结构化日志** - Go slog 标准库，JSON 格式
- ✅ **请求追踪** - Request ID 全链路追踪

### 🛡️ 可靠性
- ✅ **限流保护** - 官方 rate 限流器
- ✅ **优雅关闭** - 信号处理和资源清理
- ✅ **错误恢复** - Panic 捕获和恢复
- ✅ **超时控制** - 请求级别超时
- ✅ **定时任务** - Cron 调度 + 分布式锁

### 🧪 开发体验
- ✅ **完整测试** - 单元测试、集成测试、基准测试
- ✅ **Makefile** - 30+ 开发命令，一键操作
- ✅ **Docker 支持** - 多阶段构建 + Compose 编排
- ✅ **多环境配置** - dev/test/prod 环境隔离
- ✅ **代码质量** - golangci-lint + pre-commit hook
- ✅ **完整文档** - 47+ 文档文件，全面覆盖

---

## 📦 技术栈

### 核心框架
| 技术 | 版本 | 用途 |
|------|------|------|
| [Gin](https://github.com/gin-gonic/gin) | v1.11.0 | HTTP Web 框架 |
| [MySQL](https://www.mysql.com/) | 8.0+ | 关系型数据库 |
| [Redis](https://redis.io/) | 7.0+ | 缓存 + 分布式锁 |
| [sqlc](https://sqlc.dev/) | v1.30.0 | SQL 代码生成 |
| [Wire](https://github.com/google/wire) | v0.7.0 | 依赖注入 |

### 核心库
| 库 | 用途 |
|------|------|
| [go-redis/v9](https://github.com/redis/go-redis) | Redis 客户端（支持哨兵/集群） |
| [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) | MySQL 驱动 |
| [Viper](https://github.com/spf13/viper) | 配置管理 |
| [slog](https://pkg.go.dev/log/slog) | 结构化日志 |
| [prometheus](https://github.com/prometheus/client_golang) | 监控指标 |
| [cron/v3](https://github.com/robfig/cron) | 定时任务 |
| [testify](https://github.com/stretchr/testify) | 测试工具 |

---

## 🏗️ 项目结构

```
gin-scaffold/
├── cmd/                          # 命令行入口
│   └── main.go                  # 应用主入口
├── internal/                     # 私有代码（不可被外部导入）
│   ├── app/                     # 应用层
│   │   ├── handler/             # HTTP 处理器
│   │   │   ├── user/           # 用户模块处理器
│   │   │   └── health/         # 健康检查
│   │   ├── middleware/          # 中间件
│   │   │   ├── auth.go         # JWT 认证
│   │   │   ├── rbac.go         # 权限控制
│   │   │   ├── rate_limit.go  # 限流
│   │   │   └── ...
│   │   ├── server.go           # 服务器配置
│   │   └── app.go              # 应用引导
│   ├── domain/                  # 领域层
│   │   ├── entity/             # 实体定义
│   │   └── service/            # 业务逻辑
│   ├── repository/              # 仓储层
│   │   ├── base_repository.go  # 泛型基础仓储
│   │   ├── user_repository.go  # 用户仓储
│   │   └── query/              # sqlc 生成代码
│   ├── config/                  # 配置加载
│   ├── health/                  # 健康检查实现
│   ├── task/                    # 定时任务
│   │   ├── manager.go          # 任务管理器
│   │   └── tasks/              # 具体任务
│   └── wire/                    # Wire 依赖注入
│       ├── wire.go             # Wire 定义
│       └── wire_gen.go         # 生成代码
├── pkg/                         # 公共库（可被外部导入）
│   ├── auth/                   # JWT 认证工具
│   ├── cache/                  # 缓存管理器
│   ├── database/               # 数据库连接
│   ├── errors/                 # 错误定义
│   ├── logger/                 # 日志工具
│   ├── metrics/                # Prometheus 指标
│   ├── task/                   # 任务调度器
│   ├── health/                 # 健康检查接口
│   └── validator/              # 参数验证
├── db/                          # 数据库相关
│   ├── migrations/             # 迁移脚本
│   │   └── 001_create_users_table.sql
│   └── queries/                # SQL 查询定义
│       └── users.sql
├── docs/                        # 文档
│   ├── INDEX.md                # 文档索引
│   ├── API.md                  # API 文档
│   ├── ARCHITECTURE.md         # 架构设计
│   ├── CONFIGURATION.md        # 配置说明
│   ├── DEPLOYMENT-CHECKLIST.md # 部署清单
│   ├── TROUBLESHOOTING.md      # 故障排查
│   └── reports/                # 项目报告
├── scripts/                     # 脚本工具
│   └── wait-for-it.sh          # 服务等待脚本
├── config/                      # 配置文件
│   ├── config.yaml             # 默认配置
│   ├── config.dev.yaml         # 开发环境
│   ├── config.test.yaml        # 测试环境
│   └── config.prod.yaml        # 生产环境
├── .env.example                 # 环境变量示例
├── docker-compose.yml           # Docker 编排
├── Dockerfile                   # Docker 镜像
├── Makefile                     # 开发命令
├── sqlc.yaml                    # sqlc 配置
├── dbconfig.yml                 # 数据库迁移配置
├── .golangci.yml               # golangci-lint 配置
├── .gitignore                  # Git 忽略文件
├── go.mod                       # Go 模块定义
├── LICENSE                      # MIT 许可证
└── README.md                    # 项目说明
```

---

## 🚀 快速开始

### 方式 1: 使用 GitHub 模板（推荐）

```bash
# 1. 点击 "Use this template" 按钮创建你的仓库
# 2. 克隆你的新仓库
git clone https://github.com/yourusername/your-project.git
cd your-project

# 3. 一键初始化（安装工具 + 启动环境 + 数据库迁移）
make init

# 4. 运行项目
make run
```

### 方式 2: 手动步骤

#### 1️⃣ 环境准备

**必需软件**:
- Go 1.21+ （推荐 1.25+）
- Docker & Docker Compose
- Make

**可选工具**:
```bash
# 安装开发工具
make tools

# 或手动安装
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/google/wire/cmd/wire@latest
go install github.com/rubenv/sql-migrate/...@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### 2️⃣ 配置项目

```bash
# 1. 复制环境变量配置
cp .env.example .env

# 2. 根据需要修改配置
vim .env
# 或直接编辑 config.yaml
vim config.yaml
```

**关键配置项**:
- `database.host` - 数据库地址
- `database.password` - 数据库密码
- `redis.sentinel_enabled` - 是否启用 Redis 哨兵（默认 false）
- `jwt.secret` - JWT 密钥（**生产环境必改**）
- `server.mode` - 运行模式（debug/test/release）

#### 3️⃣ 启动服务

```bash
# 使用 Docker Compose 启动依赖服务
docker-compose up -d

# 检查服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

这将启动：
- MySQL 8.0（端口 3306）
- Redis Master（端口 6379）
- Redis Slave 1（端口 6380）
- Redis Slave 2（端口 6381）
- Redis Sentinel 1-3（端口 26379-26381）

#### 4️⃣ 数据库迁移

```bash
# 执行迁移
make migrate-up

# 检查迁移状态
make migrate-status

# 回滚（如需要）
make migrate-down
```

#### 5️⃣ 生成代码

```bash
# 生成所有代码（sqlc + wire）
make generate

# 或分别生成
make sqlc    # 生成数据库访问代码
make wire    # 生成依赖注入代码
```

#### 6️⃣ 运行项目

```bash
# 开发模式运行
make run

# 或编译后运行
make build
./bin/app
```

服务启动在 `http://localhost:8080`

#### 7️⃣ 测试接口

```bash
# 健康检查
curl http://localhost:8080/health

# Prometheus 指标
curl http://localhost:8080/metrics

# 用户注册
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 用户登录
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

---

## 📖 文档

### 核心文档
- [📘 快速开始](docs/README.md) - 新手入门指南
- [📚 文档索引](docs/INDEX.md) - 所有文档导航
- [🏛️ 架构设计](docs/ARCHITECTURE.md) - 系统架构详解
- [⚙️ 配置说明](docs/CONFIGURATION.md) - 完整配置指南
- [📡 API 文档](docs/API.md) - 接口定义和示例

### 开发指南
- [🔧 开发指南](docs/DEVELOPMENT.md) - 本地开发流程
- [📋 RBAC 权限](docs/RBAC.md) - 权限系统说明
- [💾 数据库设计](docs/DATABASE.md) - 数据库架构
- [📦 缓存策略](docs/CACHE.md) - 缓存设计

### 运维指南
- [🚀 部署清单](docs/DEPLOYMENT-CHECKLIST.md) - 生产部署步骤
- [🔍 故障排查](docs/TROUBLESHOOTING.md) - 常见问题解决
- [📊 监控指南](docs/MONITORING.md) - Prometheus 监控

### 迁移报告
- [🔄 MySQL 迁移](docs/reports/MySQL_Redis哨兵迁移完成报告.md) - PostgreSQL → MySQL 迁移记录
- [📝 优化报告](docs/reports/优化完成报告.md) - 性能优化总结

---

## 🎯 Makefile 命令速查

### 开发常用
```bash
make help              # 查看所有命令
make init              # 一键初始化项目 ⭐
make dev               # 启动开发环境
make run               # 运行应用
make test              # 运行测试
make check             # 代码检查（fmt + vet + lint）
```

### 代码生成
```bash
make generate          # 生成所有代码
make sqlc              # 生成 sqlc 代码
make wire              # 生成 Wire 代码
```

### 数据库
```bash
make migrate-up        # 执行迁移
make migrate-down      # 回滚迁移
make migrate-status    # 查看状态
make db-console        # 进入数据库控制台
```

### 构建部署
```bash
make build             # 编译应用
make build-linux       # 编译 Linux 版本
make docker-build      # 构建 Docker 镜像
make docker-run        # 运行 Docker 容器
```

### 测试分析
```bash
make test              # 运行测试
make test-cover        # 测试覆盖率
make bench             # 性能基准测试
make pprof             # 性能分析
```

### 代码质量
```bash
make lint              # 代码检查
make fmt               # 格式化代码
make complexity        # 复杂度分析
make security          # 安全扫描
make vuln              # 漏洞检查
```

### 工具管理
```bash
make tools             # 安装开发工具
make deps              # 下载依赖
make tidy              # 整理依赖
make clean             # 清理构建产物
```

完整命令列表请运行 `make help`

---

## 🏛️ 架构设计

### 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                      Presentation Layer                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Handler    │  │  Middleware  │  │     DTO      │      │
│  │ (HTTP/JSON)  │  │  (Auth/RBAC) │  │  (Request/   │      │
│  │              │  │              │  │   Response)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                       Business Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Service    │  │    Entity    │  │  Business    │      │
│  │ (业务逻辑)    │  │  (领域模型)   │  │    Rules     │      │
│  │              │  │              │  │  (验证/计算)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                     Persistence Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Repository  │  │     Query    │  │     Cache    │      │
│  │ (数据访问)    │  │   (sqlc生成)  │  │  (Redis)     │      │
│  │              │  │              │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                     Infrastructure                           │
│             ┌──────────────┐  ┌──────────────┐              │
│             │    MySQL     │  │    Redis     │              │
│             │   (Master    │  │  (Sentinel)  │              │
│             │    /Slave)   │  │              │              │
│             └──────────────┘  └──────────────┘              │
└─────────────────────────────────────────────────────────────┘
```

### 依赖注入流程

```
Wire Provider Graph:
  Application
       ├── Server (Gin Engine)
       │     ├── Handlers (User, Health...)
       │     │     └── Services
       │     │           └── Repositories
       │     │                 ├── Database (*sql.DB)
       │     │                 └── Cache (Redis)
       │     └── Middlewares
       ├── TaskManager (Cron)
       │     ├── Redis (Lock)
       │     └── Tasks
       ├── Database (*sql.DB)
       ├── Redis (UniversalClient)
       ├── Config (Viper)
       └── Logger (slog)
```

### 缓存策略

#### 三层防护机制

```
请求 → [防击穿] → [防穿透] → [防雪崩] → 数据库
         ↓            ↓           ↓
    singleflight   空值缓存    随机TTL
    (合并请求)    (NotFound)  (Jitter)
```

**实现细节**:
1. **防击穿**: 使用 `singleflight.Group` 合并并发请求
2. **防穿透**: 缓存空结果（`cache:notfound` placeholder）
3. **防雪崩**: TTL 随机化（base ± 20% jitter）

#### 缓存 Key 设计

```
user:id:{id}                    # 主键缓存
user:email:{email}              # 唯一索引缓存
user:username:{username}        # 唯一索引缓存
user:list:{page}:{size}         # 列表缓存
user:count                      # 统计缓存
```

---

## 🧪 测试

### 测试覆盖

```
Module                Coverage    Files
─────────────────────────────────────────
handler/user         72.1%       ✅
domain/service       55.3%       ✅
repository           集成测试     ✅
middleware           -           ⚠️
pkg/validator        100%        ✅
```

### 运行测试

```bash
# 所有测试
make test

# 测试覆盖率
make test-cover

# 特定包
go test ./internal/repository

# 集成测试
go test -tags=integration ./...

# 基准测试
make bench
```

### 测试示例

```go
// 单元测试
func TestUserService_GetUser(t *testing.T) {
    // 使用 testify/mock
    mockRepo := new(MockUserRepository)
    mockRepo.On("GetByID", mock.Anything, int64(1)).
        Return(&entity.User{ID: 1}, nil)
    
    service := NewUserService(mockRepo)
    user, err := service.GetUser(context.Background(), 1)
    
    assert.NoError(t, err)
    assert.Equal(t, int64(1), user.ID)
}
```

---

## 🚀 部署

### Docker 部署

```bash
# 1. 构建镜像
make docker-build

# 2. 使用 Docker Compose 部署
docker-compose up -d

# 3. 查看日志
docker-compose logs -f app

# 4. 停止服务
docker-compose down
```

### 二进制部署

```bash
# 1. 编译 Linux 版本
make build-linux

# 2. 上传到服务器
scp bin/app-linux-amd64 user@server:/app/

# 3. 在服务器上运行
./app-linux-amd64
```

### 生产环境检查清单

详见 [部署清单](docs/DEPLOYMENT-CHECKLIST.md)

**核心检查项**:
- ✅ 修改 JWT 密钥
- ✅ 配置数据库连接（主从）
- ✅ 配置 Redis 哨兵
- ✅ 设置日志级别为 `info`
- ✅ 启用 HTTPS
- ✅ 配置监控告警
- ✅ 备份策略

---

## 📊 监控

### Prometheus 指标

访问 `http://localhost:8080/metrics` 查看所有指标

**内置指标**:
```
# HTTP 指标
http_requests_total              # 请求总数
http_request_duration_seconds    # 请求耗时

# 业务指标
business_user_registered_total   # 用户注册数
business_user_login_total        # 登录次数

# 缓存指标
cache_hits_total                 # 缓存命中
cache_misses_total               # 缓存未命中

# 数据库指标
db_query_duration_seconds        # 查询耗时
db_connections_current           # 当前连接数
```

### pprof 性能分析

仅在 `debug` 和 `test` 模式下启用

```bash
# 启动性能分析
make pprof

# 或手动访问
go tool pprof http://localhost:8080/debug/pprof/profile
go tool pprof http://localhost:8080/debug/pprof/heap
```

---

## 🔧 配置说明

### 配置优先级

```
环境变量 > 命令行参数 > 配置文件 > 默认值
```

### 环境变量映射

```bash
# 数据库
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USER=root
DATABASE_PASSWORD=secret
DATABASE_NAME=myapp

# Redis 哨兵模式
REDIS_SENTINEL_ENABLED=true
REDIS_SENTINEL_MASTER=mymaster
REDIS_SENTINEL_ADDRS=sentinel1:26379,sentinel2:26379,sentinel3:26379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24

# 服务器
SERVER_MODE=release
SERVER_PORT=8080
```

完整配置说明见 [CONFIGURATION.md](docs/CONFIGURATION.md)

---

## 🤝 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详情。

### 贡献方式

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

### 代码规范

```bash
# 提交前检查
make check

# 包括：
# - go fmt         格式化
# - go vet         静态检查
# - golangci-lint  代码规范
# - go test        测试
```

---

## 📝 更新日志

### v4.0.0 (2026-01-15)
- ✨ MySQL 8.0 支持，替换 PostgreSQL
- ✨ Redis 哨兵模式高可用架构
- ✨ Prometheus 多维度监控指标
- ✨ pprof 性能分析集成
- ✨ 性能基准测试套件
- 📚 文档全面重组（47+ 文档）
- 🐛 修复所有编译和测试错误

### v3.0.0
- ✨ Wire 依赖注入
- ✨ 泛型 BaseRepository
- ✨ 三层缓存防护
- ✨ RBAC 权限系统

---

## ❓ 常见问题

### 如何切换到 PostgreSQL？

1. 修改 `config.yaml`:
```yaml
database:
  driver: postgres
  port: 5432
```

2. 重新生成代码:
```bash
make generate
```

### Redis 单机模式如何配置？

```yaml
redis:
  sentinel_enabled: false
  host: localhost
  port: 6379
```

### 如何添加新的 API 接口？

1. 在 `db/queries/` 添加 SQL
2. 运行 `make sqlc` 生成代码
3. 在 `repository/` 实现数据访问
4. 在 `service/` 实现业务逻辑
5. 在 `handler/` 实现 HTTP 处理
6. 在 `wire/` 注册依赖
7. 运行 `make wire` 生成注入代码

详见 [开发指南](docs/DEVELOPMENT.md)

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 🌟 Star History

如果这个项目对你有帮助，请给个 ⭐️ Star！

---

## 🔗 相关链接

- [Gin 文档](https://gin-gonic.com/docs/)
- [sqlc 文档](https://docs.sqlc.dev/)
- [Wire 文档](https://github.com/google/wire/blob/main/docs/guide.md)
- [Go 官方文档](https://go.dev/doc/)

---

## 📮 联系方式

- **问题反馈**: [GitHub Issues](https://github.com/yourusername/gin-scaffold/issues)
- **功能建议**: [GitHub Discussions](https://github.com/yourusername/gin-scaffold/discussions)

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给个 Star！⭐**

Made with ❤️ by the Go community

</div>
