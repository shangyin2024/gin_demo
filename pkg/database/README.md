# database - 数据库连接工具

> 支持 PostgreSQL 和 MySQL 的通用数据库连接工具

---

## 📦 特性

- ✅ 支持 PostgreSQL 和 MySQL
- ✅ 统一的配置接口
- ✅ 自动连接池管理
- ✅ 连接健康检查
- ✅ 零业务依赖，纯工具性质

---

## 🚀 使用方式

### 方式 1: 自动选择数据库类型（推荐）

```go
import "gin_demo/pkg/database"

// 配置
cfg := database.CommonConfig{
    Type:            database.TypeMySQL,  // 或 database.TypePostgreSQL
    Host:            "localhost",
    Port:            3306,
    User:            "root",
    Password:        "password",
    DBName:          "mydb",
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    ConnMaxLifetime: 5 * time.Minute,
    ConnMaxIdleTime: 10 * time.Minute,
    
    // MySQL 特定配置
    Charset:   "utf8mb4",
    ParseTime: true,
    Loc:       "Local",
}

// 创建连接
db, err := database.New(cfg)
if err != nil {
    log.Fatal(err)
}
defer database.Close(db)
```

### 方式 2: PostgreSQL 专用

```go
import "gin_demo/pkg/database"

cfg := database.Config{
    Host:            "localhost",
    Port:            5432,
    User:            "postgres",
    Password:        "postgres",
    DBName:          "mydb",
    SSLMode:         "disable",
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    ConnMaxLifetime: 5 * time.Minute,
    ConnMaxIdleTime: 10 * time.Minute,
}

db, err := database.NewPostgres(cfg)
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### 方式 3: MySQL 专用

```go
import "gin_demo/pkg/database"

cfg := database.MySQLConfig{
    Host:            "localhost",
    Port:            3306,
    User:            "root",
    Password:        "password",
    DBName:          "mydb",
    Charset:         "utf8mb4",
    ParseTime:       true,
    Loc:             "Local",
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    ConnMaxLifetime: 5 * time.Minute,
    ConnMaxIdleTime: 10 * time.Minute,
}

db, err := database.NewMySQL(cfg)
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### 方式 4: 从 DSN 创建 MySQL 连接

```go
dsn := "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

db, err := database.NewMySQLFromDSN(dsn, database.MySQLConfig{
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    ConnMaxLifetime: 5 * time.Minute,
})
```

---

## ⚙️ 配置说明

### 通用配置 (CommonConfig)

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `Type` | Type | 数据库类型（postgres/mysql） | 必填 |
| `Host` | string | 主机地址 | 必填 |
| `Port` | int | 端口号 | 必填 |
| `User` | string | 用户名 | 必填 |
| `Password` | string | 密码 | 必填 |
| `DBName` | string | 数据库名 | 必填 |
| `MaxOpenConns` | int | 最大打开连接数 | 25 |
| `MaxIdleConns` | int | 最大空闲连接数 | 5 |
| `ConnMaxLifetime` | time.Duration | 连接最大生命周期 | 5m |
| `ConnMaxIdleTime` | time.Duration | 连接最大空闲时间 | 10m |

### PostgreSQL 特定配置

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `SSLMode` | string | SSL 模式 | disable |

SSL 模式选项：
- `disable` - 不使用 SSL
- `require` - 必须使用 SSL
- `verify-ca` - 验证 CA 证书
- `verify-full` - 完全验证

### MySQL 特定配置

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `Charset` | string | 字符集 | utf8mb4 |
| `ParseTime` | bool | 是否解析时间类型 | true |
| `Loc` | string | 时区 | Local |

---

## 🔧 连接池配置建议

### 小型应用（<100 并发）

```go
MaxOpenConns:    10
MaxIdleConns:    3
ConnMaxLifetime: 5 * time.Minute
ConnMaxIdleTime: 10 * time.Minute
```

### 中型应用（100-1000 并发）

```go
MaxOpenConns:    25
MaxIdleConns:    5
ConnMaxLifetime: 5 * time.Minute
ConnMaxIdleTime: 10 * time.Minute
```

### 大型应用（>1000 并发）

```go
MaxOpenConns:    100
MaxIdleConns:    20
ConnMaxLifetime: 3 * time.Minute
ConnMaxIdleTime: 5 * time.Minute
```

---

## 📊 连接池监控

```go
// 获取连接池状态
stats := db.Stats()

fmt.Printf("OpenConnections: %d\n", stats.OpenConnections)
fmt.Printf("InUse: %d\n", stats.InUse)
fmt.Printf("Idle: %d\n", stats.Idle)
fmt.Printf("WaitCount: %d\n", stats.WaitCount)
fmt.Printf("WaitDuration: %s\n", stats.WaitDuration)
```

---

## 🎯 最佳实践

### 1. 使用连接池

```go
// ✅ 正确：复用同一个 *sql.DB
var db *sql.DB

func init() {
    var err error
    db, err = database.New(cfg)
    if err != nil {
        panic(err)
    }
}

// ❌ 错误：每次都创建新连接
func Query() {
    db, _ := database.New(cfg)  // 不要这样做！
    defer db.Close()
}
```

### 2. 使用 Context 超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

rows, err := db.QueryContext(ctx, "SELECT * FROM users")
```

### 3. 使用 Prepared Statements

```go
stmt, err := db.Prepare("SELECT * FROM users WHERE id = ?")
defer stmt.Close()

var user User
err = stmt.QueryRow(1).Scan(&user.ID, &user.Name)
```

### 4. 正确处理错误

```go
rows, err := db.Query("SELECT * FROM users")
if err != nil {
    return fmt.Errorf("query failed: %w", err)
}
defer rows.Close()

for rows.Next() {
    // ...
}

// 检查遍历过程中的错误
if err := rows.Err(); err != nil {
    return fmt.Errorf("rows iteration failed: %w", err)
}
```

---

## 🔄 数据库切换

从 PostgreSQL 切换到 MySQL：

```go
// 之前（PostgreSQL）
cfg := database.Config{
    Host:     "localhost",
    Port:     5432,
    User:     "postgres",
    Password: "postgres",
    DBName:   "mydb",
    SSLMode:  "disable",
}
db, _ := database.NewPostgres(cfg)

// 之后（MySQL）- 使用统一接口
cfg := database.CommonConfig{
    Type:      database.TypeMySQL,  // 只需修改这里
    Host:      "localhost",
    Port:      3306,
    User:      "root",
    Password:  "password",
    DBName:    "mydb",
    Charset:   "utf8mb4",
    ParseTime: true,
}
db, _ := database.New(cfg)
```

---

## 🐛 常见问题

### Q: 连接池耗尽怎么办？

**A**: 检查是否有连接泄漏

```go
// 确保释放连接
rows, err := db.Query("...")
if err != nil {
    return err
}
defer rows.Close()  // 重要！

// 或者使用 QueryRow（自动关闭）
err := db.QueryRow("...").Scan(&var)
```

### Q: MySQL 中文乱码？

**A**: 使用 utf8mb4 字符集

```go
cfg := database.MySQLConfig{
    Charset: "utf8mb4",  // 支持完整的 Unicode
    // ...
}
```

### Q: 时区问题？

**A**: MySQL 配置时区

```go
cfg := database.MySQLConfig{
    Loc: "Asia/Shanghai",  // 或 "Local"
    // ...
}
```

### Q: 如何调试 SQL 查询？

**A**: 使用日志

```go
import "log"

// 方式 1: 手动打印
log.Printf("Executing: %s with args: %v", query, args)
result, err := db.Exec(query, args...)

// 方式 2: 使用第三方库
// github.com/gchaincl/sqlhooks
```

---

## 📚 依赖

- PostgreSQL: `github.com/lib/pq`
- MySQL: `github.com/go-sql-driver/mysql`

---

## ✅ 测试清单

```bash
# PostgreSQL
□ 连接测试
□ 查询测试
□ 事务测试
□ 连接池测试

# MySQL
□ 连接测试
□ 查询测试
□ 事务测试
□ 连接池测试
□ 字符集测试
□ 时区测试
```

---

## 🔗 相关链接

- [Go database/sql 文档](https://pkg.go.dev/database/sql)
- [PostgreSQL 驱动文档](https://github.com/lib/pq)
- [MySQL 驱动文档](https://github.com/go-sql-driver/mysql)
