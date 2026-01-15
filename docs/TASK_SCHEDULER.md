# 定时任务系统

> 基于 Cron + Redis 分布式锁的企业级定时任务调度系统

---

## 🎯 核心特性

### ✅ 已实现功能

1. **Cron 调度** - 基于 `robfig/cron/v3`
   - ✅ 支持标准 Cron 表达式
   - ✅ 支持秒级精度
   - ✅ 时区支持

2. **Redis 分布式锁** - 防止重复执行
   - ✅ 自动加锁/解锁
   - ✅ 锁超时保护
   - ✅ 多实例安全

3. **任务管理**
   - ✅ 任务注册
   - ✅ 超时控制
   - ✅ 错误处理
   - ✅ 日志记录

4. **优雅关闭**
   - ✅ 等待任务完成
   - ✅ 自动清理资源

---

## 📁 目录结构

```
internal/task/
├── manager.go           # 任务管理器
└── tasks/              # 具体任务实现
    ├── example_task.go  # 示例任务
    ├── cleanup_task.go  # 清理任务
    └── stats_task.go    # 统计任务

pkg/task/
├── scheduler.go        # 任务调度器（核心）
└── base.go            # 基础任务类
```

---

## 🚀 快速开始

### 1. 创建任务

```go
package tasks

import (
    "context"
    "log/slog"
    "time"
    
    "gin_demo/pkg/task"
)

// MyTask 自定义任务
type MyTask struct{}

func NewMyTask() task.Task {
    return &MyTask{}
}

// Name 任务名称（唯一标识）
func (t *MyTask) Name() string {
    return "my_task"
}

// Spec Cron 表达式
func (t *MyTask) Spec() string {
    // 每 5 分钟执行一次
    return "0 */5 * * * *"
}

// Timeout 超时时间
func (t *MyTask) Timeout() time.Duration {
    return 2 * time.Minute
}

// Run 执行任务
func (t *MyTask) Run(ctx context.Context) error {
    slog.Info("MyTask: Starting...")
    
    // 你的任务逻辑
    // ...
    
    slog.Info("MyTask: Completed")
    return nil
}
```

### 2. 注册任务

在 `internal/task/manager.go` 中注册：

```go
func registerTasks(scheduler *task.Scheduler, redis *redis.Client, db *sql.DB) {
    // 注册你的任务
    if err := scheduler.Register(tasks.NewMyTask()); err != nil {
        panic(err)
    }
    
    // ... 其他任务
}
```

### 3. 启动（自动）

任务调度器会在应用启动时自动启动：

```go
// main.go 中已集成
app.TaskManager.Start()
defer app.TaskManager.Stop()
```

---

## 📝 Cron 表达式

### 标准格式（支持秒）

```
┌─────────── 秒 (0 - 59)
│ ┌─────────── 分 (0 - 59)
│ │ ┌─────────── 时 (0 - 23)
│ │ │ ┌─────────── 日 (1 - 31)
│ │ │ │ ┌─────────── 月 (1 - 12)
│ │ │ │ │ ┌─────────── 周 (0 - 6) (0 = 周日)
│ │ │ │ │ │
* * * * * *
```

### 常用表达式示例

| 表达式 | 说明 |
|--------|------|
| `0 * * * * *` | 每分钟执行（第 0 秒） |
| `0 */5 * * * *` | 每 5 分钟执行 |
| `0 0 * * * *` | 每小时执行 |
| `0 0 2 * * *` | 每天凌晨 2 点执行 |
| `0 0 0 * * 0` | 每周日凌晨执行 |
| `0 0 0 1 * *` | 每月 1 号凌晨执行 |
| `0 30 9 * * 1-5` | 工作日上午 9:30 执行 |

### 特殊字符

- `*` - 任意值
- `,` - 列举值 (如: `1,3,5`)
- `-` - 范围 (如: `1-5`)
- `/` - 步长 (如: `*/5`)

---

## 🔒 分布式锁机制

### 工作原理

```
实例 A                     Redis                     实例 B
   │                         │                          │
   ├─ 尝试获取锁 ────────►   │                          │
   │  SET task:lock:xxx NX  │                          │
   │                         ├─ 返回 true ───────────► │
   │                         │                          │
   ├─ 执行任务 ──────────►   │                          │
   │                         │   ◄──── 尝试获取锁 ─────┤
   │                         ├─ 返回 false ────────────┤
   │                         │   (锁已被占用)           │
   │                         │                          │
   ├─ 释放锁 ────────────►   │                          │
   │  DEL task:lock:xxx     │                          │
```

### 关键参数

```go
// 锁前缀
LockPrefix: "task:lock:"

// 锁 TTL（默认为任务超时时间）
LockTTL: 5 * time.Minute

// Redis Key 格式
// task:lock:{task_name}
```

### 安全保证

1. **原子操作** - `SET NX EX` 保证原子性
2. **自动过期** - TTL 防止死锁
3. **自动清理** - defer 确保释放锁

---

## 📊 内置任务示例

### 1. ExampleTask - 示例任务

**用途**: 演示任务系统用法

**执行频率**: 每分钟

**代码**:
```go
func (t *ExampleTask) Spec() string {
    return "0 * * * * *" // 每分钟的第 0 秒
}
```

---

### 2. CleanupTask - 清理任务

**用途**: 清理过期数据

**执行频率**: 每天凌晨 2 点

**功能**:
- 清理没有 TTL 的临时缓存
- 清理过期日志
- 清理临时文件

**代码**:
```go
func (t *CleanupTask) Spec() string {
    return "0 0 2 * * *" // 每天 02:00:00
}
```

---

### 3. StatsTask - 统计任务

**用途**: 计算统计数据

**执行频率**: 每小时

**功能**:
- 统计用户数量
- 统计 API 调用次数
- 生成报表数据

**代码**:
```go
func (t *StatsTask) Spec() string {
    return "0 0 * * * *" // 每小时的 00:00
}
```

---

## 🛠️ 高级用法

### 1. 使用 BaseTask 快速创建任务

```go
// 简单任务可以直接使用 BaseTask
func init() {
    myTask := task.NewBaseTask(
        "simple_task",           // 名称
        "0 */10 * * * *",       // 每 10 分钟
        1 * time.Minute,        // 超时 1 分钟
        func(ctx context.Context) error {
            // 你的逻辑
            log.Println("Simple task running")
            return nil
        },
    )
    
    scheduler.Register(myTask)
}
```

### 2. 任务间依赖

```go
type DependentTask struct {
    otherService *SomeService
}

func (t *DependentTask) Run(ctx context.Context) error {
    // 可以注入其他服务
    data, err := t.otherService.GetData(ctx)
    if err != nil {
        return err
    }
    
    // 处理数据
    return t.processData(data)
}
```

### 3. 动态调整执行时间

```go
type DynamicTask struct {
    config *Config
}

func (t *DynamicTask) Spec() string {
    // 从配置读取
    return t.config.TaskSchedule
}
```

### 4. 条件执行

```go
func (t *ConditionalTask) Run(ctx context.Context) error {
    // 检查条件
    if !t.shouldRun() {
        slog.Info("Task skipped: condition not met")
        return nil
    }
    
    // 执行任务
    return t.doWork(ctx)
}
```

---

## 📈 监控与调试

### 日志输出

任务系统会自动记录关键事件：

```
# 任务注册
INFO Task registered name=example_task spec="0 * * * * *"

# 任务启动
INFO Task scheduler started tasks=3

# 任务执行
INFO Task started task=example_task
INFO Task completed task=example_task duration=2.1s

# 任务失败
ERROR Task failed task=cleanup_task error="connection refused" duration=5s

# 锁冲突
DEBUG Task already running on another instance task=stats_task
```

### 查看已注册任务

```go
tasks := app.TaskManager.ListTasks()
// ["example_task", "cleanup_task", "stats_task"]
```

### 添加 Prometheus 指标

可以扩展添加任务执行指标：

```go
var (
    taskDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "task_duration_seconds",
            Help: "Task execution duration",
        },
        []string{"task_name"},
    )
    
    taskErrors = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "task_errors_total",
            Help: "Total task errors",
        },
        []string{"task_name"},
    )
)
```

---

## ⚠️ 注意事项

### 1. 任务幂等性

任务应该设计为幂等的，防止重复执行导致问题：

```go
func (t *MyTask) Run(ctx context.Context) error {
    // ✅ 好的做法：检查是否已处理
    if t.isProcessed(recordID) {
        return nil
    }
    
    // 处理逻辑
    // ...
    
    // 标记为已处理
    t.markProcessed(recordID)
    return nil
}
```

### 2. 超时设置

合理设置任务超时时间，防止任务卡死：

```go
func (t *MyTask) Timeout() time.Duration {
    // 根据任务实际执行时间设置
    // 建议设置为预期时间的 2-3 倍
    return 5 * time.Minute
}
```

### 3. 错误处理

任务失败会记录日志，但不会重试。如需重试机制，请自行实现：

```go
func (t *MyTask) Run(ctx context.Context) error {
    maxRetries := 3
    var err error
    
    for i := 0; i < maxRetries; i++ {
        err = t.doWork(ctx)
        if err == nil {
            return nil
        }
        
        slog.Warn("Task retry", "attempt", i+1, "error", err)
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    return fmt.Errorf("task failed after %d retries: %w", maxRetries, err)
}
```

### 4. 长时间运行任务

对于长时间运行的任务，注意检查 context 取消：

```go
func (t *LongTask) Run(ctx context.Context) error {
    items := t.getItems()
    
    for _, item := range items {
        // 检查是否取消
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        
        // 处理单个项目
        if err := t.processItem(item); err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## 🔧 配置建议

### 开发环境

```go
// 任务执行频率可以更频繁，便于测试
"0 */1 * * * *"  // 每分钟

// 超时时间可以更短
Timeout: 30 * time.Second
```

### 生产环境

```go
// 根据业务需求设置合理频率
"0 0 2 * * *"    // 每天凌晨 2 点（避开高峰）

// 超时时间要充足
Timeout: 10 * time.Minute

// 锁 TTL 要大于任务超时时间
LockTTL: 15 * time.Minute
```

---

## 🚨 故障排查

### 问题 1: 任务不执行

**检查**:
1. Cron 表达式是否正确
2. 任务是否已注册
3. 日志中是否有错误

**解决**:
```bash
# 查看日志
grep "Task registered" app.log
grep "Task started" app.log
```

### 问题 2: 任务重复执行

**原因**: Redis 锁未生效

**检查**:
1. Redis 连接是否正常
2. 锁 TTL 是否过短

**解决**:
```go
// 增加锁 TTL
LockTTL: 10 * time.Minute
```

### 问题 3: 任务超时

**原因**: 任务执行时间超过超时设置

**解决**:
1. 优化任务逻辑
2. 增加超时时间
3. 拆分为多个子任务

---

## 📚 API 参考

### Task 接口

```go
type Task interface {
    Name() string                     // 任务名称
    Spec() string                     // Cron 表达式
    Run(ctx context.Context) error   // 执行逻辑
    Timeout() time.Duration          // 超时时间
}
```

### Scheduler 方法

```go
// 注册任务
func (s *Scheduler) Register(task Task) error

// 启动调度器
func (s *Scheduler) Start()

// 停止调度器
func (s *Scheduler) Stop()

// 列出所有任务
func (s *Scheduler) ListTasks() []string
```

---

## 🎓 最佳实践

1. **任务命名**: 使用清晰的名称，如 `cleanup_expired_sessions`
2. **日志记录**: 在任务开始和结束时记录日志
3. **错误处理**: 不要让异常中断任务调度
4. **超时控制**: 总是设置合理的超时时间
5. **幂等设计**: 任务应该可以安全地重复执行
6. **资源清理**: 使用 defer 确保资源释放
7. **监控告警**: 对关键任务设置监控和告警

---

## 📖 参考资料

- [robfig/cron 官方文档](https://pkg.go.dev/github.com/robfig/cron/v3)
- [Cron 表达式指南](https://crontab.guru/)
- [Redis 分布式锁](https://redis.io/docs/manual/patterns/distributed-locks/)

---

**定时任务系统已就绪，支持分布式部署！** ⏰
