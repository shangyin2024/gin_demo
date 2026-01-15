# 配置文件说明

## 📁 配置文件结构

```
.
├── config.yaml         # 基础配置（默认）
├── config.dev.yaml     # 开发环境配置
├── config.test.yaml    # 测试环境配置
├── config.prod.yaml    # 生产环境配置
├── dbconfig.yml        # 数据库迁移配置
└── sqlc.yaml          # SQL 代码生成配置
```

---

## 🔧 配置优先级

配置文件加载顺序（后者覆盖前者）：

1. `config.yaml` - 基础配置
2. `config.{env}.yaml` - 环境配置（根据 `APP_ENV` 环境变量）
3. 环境变量 - 最高优先级

---

## 📝 配置项说明

### 1. 服务器配置（server）

```yaml
server:
  host: 0.0.0.0              # 监听地址
  port: 8080                 # 监听端口
  mode: debug                # 模式: debug, release, test
  read_timeout: 10s          # 读取超时
  write_timeout: 10s         # 写入超时
  idle_timeout: 60s          # 空闲连接超时
  max_request_body_size: 10485760  # 最大请求体大小（10MB）
```

### 2. 数据库配置（database）

```yaml
database:
  driver: mysql              # 数据库驱动: mysql 或 postgres
  host: localhost            # 数据库地址
  port: 3306                 # 端口（MySQL: 3306, PostgreSQL: 5432）
  user: root                 # 用户名
  password: password         # 密码（建议用环境变量 DATABASE_PASSWORD）
  dbname: gin_demo           # 数据库名
  sslmode: disable           # SSL 模式（仅 PostgreSQL）
  max_open_conns: 50         # 最大打开连接数
  max_idle_conns: 10         # 最大空闲连接数
  conn_max_lifetime: 3m      # 连接最大生命周期
  conn_max_idle_time: 5m     # 连接最大空闲时间
```

### 3. Redis 配置（redis）

#### 单机模式

```yaml
redis:
  host: localhost            # Redis 地址
  port: 6379                 # Redis 端口
  password: ""               # 密码（建议用环境变量 REDIS_PASSWORD）
  db: 0                      # 数据库编号
  max_retries: 3             # 最大重试次数
  pool_size: 10              # 连接池大小
  min_idle_conns: 5          # 最小空闲连接数
  sentinel_enabled: false    # 是否启用哨兵模式
```

#### 哨兵模式（生产环境推荐）

```yaml
redis:
  sentinel_enabled: true     # 启用哨兵模式
  sentinel_master: mymaster  # 主节点名称
  sentinel_addrs:            # 哨兵地址列表
    - localhost:26379
    - localhost:26380
    - localhost:26381
  password: ""               # Redis 密码
  db: 0
  max_retries: 3
  pool_size: 50
  min_idle_conns: 10
```

### 4. 日志配置（logger）

```yaml
logger:
  level: info                # 日志级别: debug, info, warn, error
  format: json               # 格式: json, text
  add_source: false          # 是否添加代码位置
  request_id_key: request_id # Request ID 的键名
```

### 5. JWT 配置（jwt）

```yaml
jwt:
  secret: your-secret-key    # JWT 密钥（建议用环境变量 JWT_SECRET）
  expiration: 24h            # Token 过期时间
```

### 6. CORS 配置（cors）

```yaml
cors:
  allowed_origins:           # 允许的来源
    - "http://localhost:3000"
    - "http://localhost:8080"
  allow_credentials: true    # 是否允许携带凭证
  max_age: 43200            # 预检请求缓存时间（秒）
```

### 7. 安全配置（security）

```yaml
security:
  # HTTP 安全头
  headers:
    enabled: true                    # 是否启用
    enable_hsts: false               # HSTS（需要 HTTPS）
    hsts_max_age: 31536000          # HSTS 有效期（秒）
    hsts_include_subdomains: true   # HSTS 包含子域名
    enable_csp: true                # 内容安全策略
    csp_policy: "default-src 'self';" # CSP 策略
    enable_frame_options: true      # X-Frame-Options
    frame_options: "DENY"           # DENY 或 SAMEORIGIN
  
  # 压缩
  enable_compression: true          # 启用 Gzip 压缩
  compression_level: 5              # 压缩级别 (-1 或 0-9)
  
  # TLS/HTTPS
  tls:
    enabled: false                  # 是否启用 TLS
    cert_file: ""                   # 证书文件路径
    key_file: ""                    # 密钥文件路径
    min_version: "1.2"              # 最低 TLS 版本
```

### 8. 缓存配置（cache）

```yaml
cache:
  default_ttl: 5m            # 默认过期时间
  user_ttl: 5m              # 用户缓存过期时间
  user_index_ttl: 10m       # 用户索引缓存过期时间
  user_count_ttl: 1m        # 用户统计缓存过期时间
  user_session_ttl: 30m     # 用户会话缓存过期时间
  content_ttl: 10m          # 内容缓存过期时间
  content_list_ttl: 2m      # 内容列表缓存过期时间
  stats_ttl: 1m             # 统计数据缓存过期时间
  not_found_ttl: 5m         # 不存在记录缓存过期时间
  enable_jitter: true       # 启用缓存抖动
  jitter_percent: 20        # 抖动百分比
```

---

## 🌍 环境变量覆盖

敏感信息建议通过环境变量设置：

```bash
# 数据库
export DATABASE_PASSWORD="your-db-password"

# Redis
export REDIS_PASSWORD="your-redis-password"

# JWT
export JWT_SECRET="your-jwt-secret-key"

# 环境选择
export APP_ENV="prod"  # dev, test, prod
```

---

## 📚 环境配置说明

### 开发环境（config.dev.yaml）

- **模式**: debug
- **数据库**: localhost:3306
- **Redis**: 单机模式
- **日志**: debug 级别，text 格式
- **安全**: 关闭 HSTS 和 TLS

启动方式：
```bash
export APP_ENV=dev
./bin/app
```

### 测试环境（config.test.yaml）

- **模式**: test
- **数据库**: gin_demo_test
- **Redis**: db=1（避免冲突）
- **日志**: debug 级别
- **安全**: 禁用安全头和压缩

启动方式：
```bash
export APP_ENV=test
go test ./...
```

### 生产环境（config.prod.yaml）

- **模式**: release
- **数据库**: MySQL 主从，连接池 100
- **Redis**: 哨兵模式（高可用）
- **日志**: info 级别，json 格式
- **安全**: 全部启用（HSTS, CSP, TLS）

启动方式：
```bash
export APP_ENV=prod
export DATABASE_PASSWORD="xxx"
export REDIS_PASSWORD="xxx"
export JWT_SECRET="xxx"
./bin/app
```

---

## 🔍 配置验证

启动时会自动验证配置：

```go
// 生产环境强制检查
if cfg.Server.Mode == "release" {
    - JWT Secret 不能为默认值
    - 数据库密码不能为空
    - Redis 密码建议设置
}
```

---

## 📖 相关文档

- [部署清单](./DEPLOYMENT-CHECKLIST.md)
- [MySQL 迁移指南](./MYSQL_MIGRATION.md)
- [故障排查](./TROUBLESHOOTING.md)
- [架构文档](./ARCHITECTURE.md)
