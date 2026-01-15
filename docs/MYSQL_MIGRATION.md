# 🔄 MySQL 迁移指南

**当前状态**: PostgreSQL 15+  
**目标**: MySQL 8.0+  
**迁移策略**: 最小改动，保持代码通用性

---

## 📋 迁移检查清单

### 1. SQL 语法差异 ⚠️

#### 占位符不同
```sql
-- PostgreSQL (当前)
SELECT * FROM users WHERE id = $1 AND email = $2

-- MySQL (需要改为)
SELECT * FROM users WHERE id = ? AND email = ?
```

**解决方案**: 修改 `sqlc.yaml` 配置
```yaml
version: "2"
sql:
  - engine: "mysql"  # 改为 mysql
    queries: "internal/repository/queries/"
    schema: "db/schema/"
    gen:
      go:
        package: "repository"
        out: "internal/repository"
```

#### AUTO_INCREMENT vs SERIAL
```sql
-- PostgreSQL
id BIGSERIAL PRIMARY KEY

-- MySQL
id BIGINT AUTO_INCREMENT PRIMARY KEY
```

#### RETURNING 子句不支持
```sql
-- PostgreSQL (当前)
INSERT INTO users (...) VALUES (...) RETURNING *;

-- MySQL (需要改为)
INSERT INTO users (...) VALUES (...);
-- 然后使用 LAST_INSERT_ID()
```

#### JSON 类型
```sql
-- PostgreSQL
data JSONB  -- 更高效

-- MySQL 8.0+
data JSON   -- 支持，但性能略低
```

---

## 🔧 代码适配

### 1. 驱动更换

```go
// go.mod
// 移除 PostgreSQL
- github.com/lib/pq v1.10.9

// 添加 MySQL
+ github.com/go-sql-driver/mysql v1.7.1
```

### 2. 连接字符串

```go
// internal/wire/infrastructure.go

// PostgreSQL (当前)
dsn := fmt.Sprintf(
    "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
    cfg.Database.Host,
    cfg.Database.Port,
    cfg.Database.User,
    cfg.Database.Password,
    cfg.Database.Name,
    cfg.Database.SSLMode,
)

// MySQL (修改为)
dsn := fmt.Sprintf(
    "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
    cfg.Database.User,
    cfg.Database.Password,
    cfg.Database.Host,
    cfg.Database.Port,
    cfg.Database.Name,
)

db, err := sql.Open("mysql", dsn)  // 改为 mysql
```

### 3. 数据库迁移工具

```yaml
# dbconfig.yml

# PostgreSQL (当前)
development:
  dialect: postgres
  datasource: host=localhost port=5432 ...

# MySQL (修改为)
development:
  dialect: mysql
  datasource: root:password@tcp(localhost:3306)/gin_demo?parseTime=true
  dir: db/migrations
  table: schema_migrations
```

---

## 🔍 MySQL 特有优化

### 1. 字符集设置 ⭐⭐⭐⭐⭐

```sql
-- 创建数据库时指定
CREATE DATABASE gin_demo
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

-- 创建表时
CREATE TABLE users (
  ...
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**为什么是 utf8mb4**:
- MySQL 的 `utf8` 只支持 3 字节（不支持 emoji 😀）
- `utf8mb4` 支持 4 字节（完整 Unicode）

### 2. 存储引擎选择

```sql
-- 推荐 InnoDB (默认)
ENGINE=InnoDB

优势:
✅ 支持事务 (ACID)
✅ 支持外键
✅ 行级锁（并发好）
✅ 崩溃恢复
```

### 3. 索引优化

```sql
-- MySQL 索引特性
CREATE INDEX idx_users_email ON users(email);  -- B+Tree

-- 全文索引 (MySQL 5.6+)
CREATE FULLTEXT INDEX idx_users_search ON users(username, bio);

-- 使用
SELECT * FROM users WHERE MATCH(username, bio) AGAINST('keyword');
```

### 4. 连接配置优化

```go
// MySQL 特有参数
db.SetConnMaxLifetime(time.Minute * 3)  // 连接最大生命周期
db.SetMaxOpenConns(50)                   // 最大打开连接数
db.SetMaxIdleConns(10)                   // 最大空闲连接数

// MySQL 连接字符串参数
dsn := "user:pass@tcp(host:3306)/db?" +
    "charset=utf8mb4" +              // 字符集
    "&parseTime=True" +               // 解析 TIME 类型
    "&loc=Local" +                    // 时区
    "&timeout=10s" +                  // 连接超时
    "&readTimeout=30s" +              // 读超时
    "&writeTimeout=30s" +             // 写超时
    "&maxAllowedPacket=67108864"      // 最大包大小 (64MB)
```

---

## 📊 性能对比

| 特性 | PostgreSQL | MySQL (InnoDB) | 建议 |
|------|-----------|----------------|------|
| **并发写入** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐☆ | PG 略优 |
| **复杂查询** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐☆ | PG 优化器更好 |
| **JSON 支持** | ⭐⭐⭐⭐⭐ (JSONB) | ⭐⭐⭐⭐☆ (JSON) | PG 更强 |
| **全文搜索** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐☆☆ | PG 更强 |
| **简单查询** | ⭐⭐⭐⭐☆ | ⭐⭐⭐⭐⭐ | MySQL 略快 |
| **生态系统** | ⭐⭐⭐⭐☆ | ⭐⭐⭐⭐⭐ | MySQL 更广 |
| **学习曲线** | ⭐⭐⭐☆☆ | ⭐⭐⭐⭐☆ | MySQL 更简单 |

---

## 🚀 迁移步骤

### Phase 1: 准备阶段（1-2天）

```bash
# 1. 安装 MySQL 8.0+
docker run -d \
  --name mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=gin_demo \
  mysql:8.0 \
  --character-set-server=utf8mb4 \
  --collation-server=utf8mb4_unicode_ci

# 2. 备份 PostgreSQL 数据（如果有）
pg_dump -h localhost -U postgres gin_demo > backup.sql

# 3. 更新依赖
go get github.com/go-sql-driver/mysql@latest
go mod tidy
```

### Phase 2: SQL 迁移（2-3天）

```bash
# 1. 转换 Schema
# db/schema/001_users.sql

-- PostgreSQL 版本
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- MySQL 版本
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_email (email),
    INDEX idx_users_username (username),
    INDEX idx_users_status (status),
    INDEX idx_users_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

# 2. 转换查询 (queries/*.sql)
# 将 $1, $2 改为 ?, ?
find internal/repository/queries -name "*.sql" -exec sed -i 's/\$[0-9]/?/g' {} \;

# 3. 重新生成代码
sqlc generate
```

### Phase 3: 代码适配（1天）

```go
// 1. 更新导入
import (
    _ "github.com/go-sql-driver/mysql"  // MySQL 驱动
)

// 2. 更新配置
// config.yaml
database:
  driver: mysql  # 改为 mysql
  host: localhost
  port: 3306     # MySQL 端口
  name: gin_demo
  user: root
  password: password
  ssl_mode: ""   # MySQL 不需要
  
  # MySQL 特有配置
  charset: utf8mb4
  parse_time: true
  loc: Local
  
  max_open_conns: 50
  max_idle_conns: 10
  conn_max_lifetime: 180  # 秒

// 3. 更新数据库初始化
func provideDatabase(cfg *config.Config) (*sql.DB, error) {
    var dsn string
    
    switch cfg.Database.Driver {
    case "mysql":
        dsn = fmt.Sprintf(
            "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
            cfg.Database.User,
            cfg.Database.Password,
            cfg.Database.Host,
            cfg.Database.Port,
            cfg.Database.Name,
            cfg.Database.Charset,
            cfg.Database.ParseTime,
            cfg.Database.Loc,
        )
    case "postgres":
        dsn = fmt.Sprintf(
            "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
            cfg.Database.Host,
            cfg.Database.Port,
            cfg.Database.User,
            cfg.Database.Password,
            cfg.Database.Name,
            cfg.Database.SSLMode,
        )
    default:
        return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
    }
    
    db, err := sql.Open(cfg.Database.Driver, dsn)
    if err != nil {
        return nil, err
    }
    
    // 连接池配置
    db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
    db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
    db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)
    
    return db, db.Ping()
}
```

### Phase 4: 测试验证（2-3天）

```bash
# 1. 单元测试
APP_ENV=test go test -v ./internal/domain/service/...

# 2. 集成测试
APP_ENV=test go test -v ./internal/repository/...

# 3. HTTP 测试
APP_ENV=test go test -v ./internal/app/handler/...

# 4. 手动测试
make run
curl http://localhost:8080/health
```

---

## ⚠️ 常见问题

### 1. 时间类型处理

```go
// MySQL 需要在连接字符串中添加 parseTime=True
dsn := "...?parseTime=True"

// 否则 time.Time 会被解析为 []byte
```

### 2. 布尔类型

```sql
-- PostgreSQL 有原生 BOOLEAN
is_active BOOLEAN DEFAULT TRUE

-- MySQL 使用 TINYINT(1)
is_active TINYINT(1) DEFAULT 1

-- Go 代码中统一使用 bool 即可
```

### 3. LIMIT 语法

```sql
-- 两者相同
SELECT * FROM users LIMIT 10 OFFSET 20;  -- ✅ 都支持

-- MySQL 简写
SELECT * FROM users LIMIT 20, 10;  -- ✅ MySQL 特有
```

### 4. 事务隔离级别

```go
// MySQL 默认: REPEATABLE READ (可重复读)
// PostgreSQL 默认: READ COMMITTED (读已提交)

// 建议: 显式设置
tx, err := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelReadCommitted,  // 统一使用
})
```

---

## 📚 推荐资源

### MySQL 最佳实践
- [MySQL 8.0 Reference Manual](https://dev.mysql.com/doc/refman/8.0/en/)
- [High Performance MySQL](https://www.oreilly.com/library/view/high-performance-mysql/9781492080503/)
- [MySQL Internals Manual](https://dev.mysql.com/doc/internals/en/)

### 性能优化
- [MySQL Performance Tuning](https://dev.mysql.com/doc/refman/8.0/en/optimization.html)
- [MySQL 索引优化](https://dev.mysql.com/doc/refman/8.0/en/optimization-indexes.html)
- [InnoDB 引擎优化](https://dev.mysql.com/doc/refman/8.0/en/innodb-optimization.html)

---

## ✅ 验收标准

迁移完成后，确保以下检查通过：

- [ ] 所有测试通过 (`go test ./...`)
- [ ] 应用正常启动
- [ ] 健康检查通过 (`/health`)
- [ ] 用户注册正常
- [ ] 用户登录正常
- [ ] 数据查询正常
- [ ] 缓存工作正常
- [ ] 事务工作正常
- [ ] 性能测试通过
- [ ] 压测结果满意

---

**预计总工作量**: 6-8 天  
**风险等级**: 中等  
**建议**: 在测试环境完整验证后再迁移生产环境
