# Task - 定时任务调度器

> 基于 Cron + Redis 分布式锁的企业级任务调度系统

---

## 🌟 特性

- ✅ **Cron 调度** - 支持秒级精度的 Cron 表达式
- ✅ **分布式锁** - Redis 实现，防止多实例重复执行
- ✅ **超时控制** - 每个任务可独立设置超时时间
- ✅ **优雅关闭** - 等待运行中的任务完成
- ✅ **错误处理** - 完善的错误记录和处理机制
- ✅ **简单易用** - 清晰的接口设计

---

## 📦 安装

```bash
go get github.com/robfig/cron/v3
```

---

## 🚀 快速开始

### 1. 创建任务

```go
import (
    "context"
    "time"
    "gin_demo/pkg/task"
)

type MyTask struct{}

func (t *MyTask) Name() string {
    return "my_task"
}

func (t *MyTask) Spec() string {
    return "0 */5 * * * *" // 每 5 分钟
}

func (t *MyTask) Timeout() time.Duration {
    return 2 * time.Minute
}

func (t *MyTask) Run(ctx context.Context) error {
    // 你的任务逻辑
    return nil
}
```

### 2. 创建调度器并注册任务

```go
// 创建调度器
scheduler := task.NewScheduler(task.Config{
    Redis:      redisClient,
    LockPrefix: "task:lock:",
})

// 注册任务
if err := scheduler.Register(&MyTask{}); err != nil {
    log.Fatal(err)
}

// 启动调度器
scheduler.Start()
defer scheduler.Stop()
```

---

## 📝 Cron 表达式

### 格式

```
秒 分 时 日 月 周
*  *  *  *  *  *
```

### 示例

```go
"0 * * * * *"      // 每分钟
"0 */5 * * * *"    // 每 5 分钟
"0 0 * * * *"      // 每小时
"0 0 2 * * *"      // 每天凌晨 2 点
"0 30 9 * * 1-5"   // 工作日上午 9:30
```

---

## 🔒 分布式锁

### 原理

使用 Redis `SET NX EX` 实现分布式锁：

```
实例 A 尝试执行任务
  ↓
获取 Redis 锁
  ↓
成功? 
  ├─ 是 → 执行任务 → 释放锁
  └─ 否 → 跳过（其他实例正在执行）
```

### 安全保证

- **原子操作** - SET NX EX 保证原子性
- **自动过期** - 防止死锁
- **自动释放** - defer 确保锁释放

---

## 🛠️ API 文档

### Task 接口

```go
type Task interface {
    // Name 返回任务唯一名称
    Name() string
    
    // Spec 返回 Cron 表达式
    Spec() string
    
    // Run 执行任务
    Run(ctx context.Context) error
    
    // Timeout 返回任务超时时间
    Timeout() time.Duration
}
```

### Scheduler

```go
// 创建调度器
func NewScheduler(config Config) *Scheduler

// 注册任务
func (s *Scheduler) Register(task Task) error

// 启动调度器
func (s *Scheduler) Start()

// 停止调度器（等待运行中的任务）
func (s *Scheduler) Stop()

// 列出所有已注册任务
func (s *Scheduler) ListTasks() []string
```

### BaseTask

快速创建简单任务：

```go
task := task.NewBaseTask(
    "simple_task",           // 名称
    "0 */10 * * * *",       // Cron
    1 * time.Minute,        // 超时
    func(ctx context.Context) error {
        // 任务逻辑
        return nil
    },
)
```

---

## 📊 日志输出

```
INFO Task registered name=my_task spec="0 */5 * * * *"
INFO Task scheduler started tasks=3
INFO Task started task=my_task
INFO Task completed task=my_task duration=1.2s
```

---

## ⚠️ 注意事项

1. **幂等性** - 任务应设计为幂等的
2. **超时** - 合理设置超时时间
3. **错误处理** - 任务失败会记录日志但不重试
4. **Context** - 长任务应检查 context 取消

---

## 📖 完整文档

详见 [定时任务系统文档](../../docs/TASK_SCHEDULER.md)

---

**通用、可靠、易用的任务调度系统！** ⏰
