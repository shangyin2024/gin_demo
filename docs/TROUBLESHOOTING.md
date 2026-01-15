# 🔧 故障排查手册

**版本**: v3.0.0  
**更新日期**: 2026-01-15

---

## 📋 快速诊断

### 服务状态检查

```bash
# 1. 检查服务是否运行
curl http://localhost:8080/health
curl http://localhost:8080/health/ready
curl http://localhost:8080/health/live

# 2. 检查进程
ps aux | grep gin-demo

# 3. 检查日志
tail -f /var/log/gin-demo/app.log

# 4. 检查监控指标
curl http://localhost:8080/metrics | head -100
```

---

## 🔥 常见问题排查

### 问题 1: 服务启动失败

#### 症状
```bash
$ ./gin-demo
Failed to load config: ...
exit status 1
```

#### 可能原因

**原因 A: 配置文件缺失或格式错误**

```bash
# 检查
ls -la config*.yaml
cat config.yaml | head -20

# 解决
1. 确保 config.yaml 存在
2. 验证 YAML 语法: yamllint config.yaml
3. 检查必填字段是否存在
```

**原因 B: 环境变量未设置**

```bash
# 检查
echo $APP_ENV
echo $JWT_SECRET
echo $DATABASE_PASSWORD

# 解决
export APP_ENV=prod
export JWT_SECRET=your-secret-key
export DATABASE_PASSWORD=your-db-password
```

**原因 C: 数据库连接失败**

```bash
# 检查数据库连接
psql -h localhost -U postgres -d gin_demo -c "SELECT 1"

# 解决
1. 确保数据库已启动: docker-compose up -d postgres
2. 检查数据库配置是否正确
3. 检查网络连通性: telnet localhost 5432
```

**原因 D: Redis 连接失败**

```bash
# 检查 Redis 连接
redis-cli -h localhost ping

# 解决
1. 启动 Redis: docker-compose up -d redis
2. 检查 Redis 配置
3. 检查网络连通性: telnet localhost 6379
```

---

### 问题 2: 接口返回 500 错误

#### 症状
```json
{
  "code": 50001,
  "message": "服务器内部错误"
}
```

#### 排查步骤

**步骤 1: 查看日志**

```bash
# 查看错误日志
tail -f /var/log/gin-demo/app.log | grep ERROR

# 或使用 Docker
docker-compose logs -f app | grep ERROR

# 或使用 Kubernetes
kubectl logs -f deployment/gin-demo | grep ERROR
```

**步骤 2: 检查数据库**

```bash
# 检查数据库连接
curl http://localhost:8080/health | jq '.checks.database'

# 检查慢查询
psql -c "SELECT * FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10"

# 检查连接数
psql -c "SELECT count(*) FROM pg_stat_activity"
```

**步骤 3: 检查缓存**

```bash
# 检查 Redis 状态
redis-cli info stats

# 检查缓存命中率
curl http://localhost:8080/metrics | grep cache_hits

# 检查缓存大小
redis-cli dbsize
```

**步骤 4: 检查监控指标**

```bash
# 错误率
curl http://localhost:8080/metrics | grep 'http_requests_total{.*status="5'

# 响应时间
curl http://localhost:8080/metrics | grep http_request_duration_seconds

# 数据库连接
curl http://localhost:8080/metrics | grep db_connections_current
```

---

### 问题 3: 响应时间慢

#### 症状
```
接口响应时间 > 1s
P99 延迟告警
```

#### 排查清单

**1. 检查慢查询**

```bash
# Prometheus 查询
rate(db_slow_queries_total{threshold="100ms"}[5m])

# 数据库慢查询日志
psql -c "SELECT query, mean_exec_time FROM pg_stat_statements WHERE mean_exec_time > 100 ORDER BY mean_exec_time DESC LIMIT 10"

# 应用日志
grep "Slow query" /var/log/gin-demo/app.log
```

**解决方案**:
```sql
-- 1. 添加索引
CREATE INDEX idx_users_created_at ON users(created_at);

-- 2. 优化查询
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';

-- 3. 分析索引使用情况
SELECT * FROM pg_stat_user_indexes WHERE schemaname = 'public';
```

**2. 检查缓存命中率**

```bash
# 缓存命中率
curl http://localhost:8080/metrics | grep -E "(cache_hits|cache_misses)"

# 计算命中率
echo "命中率 = hits / (hits + misses)"
```

**解决方案**:
- 命中率 < 50%: 增加 TTL 或预热缓存
- 命中率 < 20%: 检查缓存键是否正确

**3. 检查并发连接**

```bash
# 数据库连接数
curl http://localhost:8080/metrics | grep db_connections_current

# Redis 连接数
redis-cli client list | wc -l
```

**解决方案**:
- 连接数 > 80%: 增加连接池大小
- 连接数波动大: 检查是否有连接泄漏

---

### 问题 4: 认证失败

#### 症状
```json
{
  "code": 10002,
  "message": "未授权"
}
```

#### 排查步骤

**1. 检查 Token 格式**

```bash
# Token 应该是这样的格式
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# 常见错误
❌ Authorization: eyJhbGc...  (缺少 Bearer 前缀)
❌ Bearer eyJhbGc...          (缺少 Authorization 键)
```

**2. 检查 Token 有效性**

```bash
# 使用 jwt.io 解码 Token
# 或使用命令行工具

# 检查过期时间
jwt decode <token> | jq '.exp'
date -r <timestamp>

# 检查签名
# Token 应该使用正确的 JWT_SECRET 签名
```

**3. 检查日志**

```bash
# 查看认证日志
grep "Token" /var/log/gin-demo/app.log | tail -20

# 查看认证失败原因
grep "unauthorized" /var/log/gin-demo/app.log
```

---

### 问题 5: 缓存雪崩

#### 症状
```
- 数据库 CPU 突然飙升
- 响应时间剧增
- 大量慢查询
- 缓存未命中率突增
```

#### 应急处理

```bash
# 1. 立即限流（降低数据库压力）
# 修改限流配置，降低 QPS

# 2. 检查 Redis 状态
redis-cli ping
redis-cli info stats

# 3. 检查缓存键
redis-cli keys "cache:*" | head -20

# 4. 重启 Redis（如果必要）
docker-compose restart redis

# 5. 缓存预热
curl -X POST http://localhost:8080/admin/cache/warmup
```

#### 长期解决

```go
// 1. 确保 Jitter 已启用
cache.enable_jitter: true
cache.jitter_percent: 20

// 2. 避免大 key 过期
// 分散过期时间

// 3. 缓存降级
// Redis 故障时直接查数据库
```

---

### 问题 6: 内存泄漏

#### 症状
```
- 内存持续增长
- OOM killed
- 垃圾回收频繁
```

#### 排查步骤

**1. 分析内存使用**

```bash
# 1. 获取 pprof 数据
curl http://localhost:8080/debug/pprof/heap > heap.prof

# 2. 分析内存
go tool pprof heap.prof
> top 10
> list <function_name>

# 3. 查看 goroutine
curl http://localhost:8080/debug/pprof/goroutine?debug=1 | grep goroutine
```

**2. 常见泄漏点**

```go
// A. 未关闭的 HTTP 连接
defer resp.Body.Close()

// B. 未关闭的数据库连接
defer rows.Close()

// C. 无界的 goroutine
// 检查 task scheduler 是否泄漏

// D. 无界的 map
// 检查 RateLimiter.limiters 是否无限增长
```

**3. 临时解决**

```bash
# 重启服务（临时缓解）
docker-compose restart app

# 或 K8s
kubectl rollout restart deployment/gin-demo
```

**4. 长期修复**

```go
// 添加 goroutine 池
import "github.com/panjf2000/ants/v2"

pool, _ := ants.NewPool(100)  // 最多100个并发
defer pool.Release()

pool.Submit(func() {
    // 任务逻辑
})
```

---

### 问题 7: 数据库死锁

#### 症状
```
ERROR: deadlock detected
```

#### 排查步骤

```sql
-- 1. 查看当前锁
SELECT * FROM pg_locks WHERE NOT granted;

-- 2. 查看阻塞关系
SELECT 
    blocked_locks.pid AS blocked_pid,
    blocking_locks.pid AS blocking_pid,
    blocked_activity.query AS blocked_query,
    blocking_activity.query AS blocking_query
FROM pg_locks blocked_locks
JOIN pg_stat_activity blocked_activity ON blocked_locks.pid = blocked_activity.pid
JOIN pg_locks blocking_locks ON blocked_locks.relation = blocking_locks.relation
JOIN pg_stat_activity blocking_activity ON blocking_locks.pid = blocking_activity.pid
WHERE NOT blocked_locks.granted;

-- 3. 杀死阻塞进程（慎重！）
SELECT pg_terminate_backend(<pid>);
```

#### 预防措施

```go
// 1. 事务尽量短
// 2. 操作顺序一致（避免循环等待）
// 3. 使用合适的隔离级别

// 错误示例（可能死锁）
tx1: UPDATE users SET ... WHERE id = 1;
tx1: UPDATE orders SET ... WHERE user_id = 1;

tx2: UPDATE orders SET ... WHERE user_id = 1;  ← 等待 tx1
tx2: UPDATE users SET ... WHERE id = 1;        ← 死锁！

// 正确示例（顺序一致）
tx1: UPDATE users ...   → UPDATE orders ...
tx2: UPDATE users ...   → UPDATE orders ...
```

---

## 📊 监控指标解读

### 关键指标阈值

| 指标 | 正常 | 警告 | 危险 |
|------|------|------|------|
| P99 延迟 | < 200ms | 200-500ms | > 500ms |
| 错误率 | < 0.1% | 0.1-1% | > 1% |
| 缓存命中率 | > 80% | 60-80% | < 60% |
| 慢查询占比 | < 5% | 5-10% | > 10% |
| 数据库连接 | < 60% | 60-80% | > 80% |
| Redis 内存 | < 70% | 70-90% | > 90% |

### 告警处理流程

```
告警触发
    ↓
查看仪表盘（确认）
    ↓
检查日志（定位）
    ↓
临时缓解（限流/重启）
    ↓
根因分析（代码/配置）
    ↓
修复上线（PR + 部署）
    ↓
验证解决（监控观察）
    ↓
记录文档（事后复盘）
```

---

## 🆘 紧急故障处理

### 场景 1: 服务完全不可用

```bash
# 1. 立即检查
curl http://localhost:8080/health  # 超时或连接拒绝

# 2. 检查进程
ps aux | grep gin-demo  # 进程不存在？

# 3. 查看日志
tail -100 /var/log/gin-demo/app.log

# 4. 尝试重启
systemctl restart gin-demo
# 或
docker-compose restart app

# 5. 如果无法启动，回滚
docker-compose down
docker-compose up -d app:v2.2.0  # 上一个稳定版本

# 6. 通知团队
# 发送告警通知
```

### 场景 2: 数据库连接耗尽

```bash
# 症状
ERROR: sorry, too many clients already

# 立即处理
# 1. 检查连接数
psql -c "SELECT count(*) FROM pg_stat_activity"

# 2. 查看连接来源
psql -c "SELECT client_addr, count(*) FROM pg_stat_activity GROUP BY client_addr"

# 3. 杀死空闲连接（慎重！）
psql -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE state = 'idle' AND state_change < now() - interval '5 minutes'"

# 4. 增加连接池限制（临时）
docker exec postgres psql -U postgres -c "ALTER SYSTEM SET max_connections = 200"
docker restart postgres
```

### 场景 3: 内存溢出（OOM）

```bash
# 症状
FATAL: kernel killed process (OOM)

# 排查
# 1. 查看系统内存
free -h

# 2. 查看进程内存
ps aux --sort=-%mem | head -10

# 3. 查看应用内存
curl http://localhost:8080/debug/pprof/heap

# 临时解决
# 1. 重启服务
docker-compose restart app

# 2. 增加内存限制
# docker-compose.yml
services:
  app:
    mem_limit: 1g  # 增加到 1GB

# 长期解决
# 1. 分析内存泄漏
go tool pprof http://localhost:8080/debug/pprof/heap

# 2. 修复代码
# 3. 添加监控告警
```

---

## 🔍 性能问题排查

### 定位性能瓶颈

```bash
# 1. 整体性能分析
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof
go tool pprof cpu.prof
> top 20
> list <function_name>

# 2. 数据库性能
# 查看慢查询
curl http://localhost:8080/metrics | grep db_slow_queries

# 查看查询延迟分布
curl http://localhost:8080/metrics | grep db_query_duration_seconds

# 3. 缓存性能
# 查看缓存命中率
curl http://localhost:8080/metrics | grep cache_hits

# 查看缓存延迟
curl http://localhost:8080/metrics | grep cache_operation_duration

# 4. 网络延迟
# 使用 tcpdump 抓包分析
tcpdump -i any -w capture.pcap port 8080
```

### 优化建议

**数据库优化**:
```sql
-- 1. 添加索引
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);

-- 2. 分析查询计划
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';

-- 3. 更新统计信息
ANALYZE users;

-- 4. 清理膨胀
VACUUM FULL users;
```

**缓存优化**:
```go
// 1. 增加 TTL
cache.user_ttl: 10m  // 从 5m 增加到 10m

// 2. 预热热点数据
func (m *Manager) WarmupHotData() {
    hotUsers := []int64{1, 2, 3, 4, 5}
    for _, id := range hotUsers {
        go repo.GetUserByID(ctx, id)
    }
}

// 3. 批量操作
// 使用 Pipeline 减少网络往返
```

**应用优化**:
```go
// 1. 减少数据库查询
// 使用缓存

// 2. 并发处理
// 使用 goroutine + WaitGroup

// 3. 连接复用
// 使用全局 HTTP Client
```

---

## 📚 日志分析

### 日志级别说明

```
DEBUG: 详细的调试信息（开发环境）
INFO:  正常的业务流程（生产环境默认）
WARN:  警告信息（需要关注但不影响功能）
ERROR: 错误信息（影响功能，需要立即处理）
```

### 日志查询示例

```bash
# 1. 查看最近的错误
tail -1000 /var/log/gin-demo/app.log | grep ERROR

# 2. 统计错误类型
grep ERROR /var/log/gin-demo/app.log | cut -d'"' -f4 | sort | uniq -c | sort -rn

# 3. 追踪特定请求（通过 Request ID）
grep "request_id=abc-123" /var/log/gin-demo/app.log

# 4. 查看特定用户的操作
grep "user_id=12345" /var/log/gin-demo/app.log

# 5. 查看慢查询
grep "Slow query" /var/log/gin-demo/app.log
```

---

## 🛠️ 实用工具脚本

### 1. 健康检查脚本

```bash
#!/bin/bash
# scripts/health_check.sh

echo "🔍 Checking service health..."

# 检查服务
if curl -sf http://localhost:8080/health > /dev/null; then
    echo "✅ Service is healthy"
else
    echo "❌ Service is down!"
    exit 1
fi

# 检查数据库
if curl -sf http://localhost:8080/health | jq -e '.checks.database.status == "ok"' > /dev/null; then
    echo "✅ Database is healthy"
else
    echo "❌ Database connection failed!"
    exit 1
fi

# 检查 Redis
if curl -sf http://localhost:8080/health | jq -e '.checks.redis.status == "ok"' > /dev/null; then
    echo "✅ Redis is healthy"
else
    echo "⚠️  Redis connection failed (degraded mode)"
fi

echo "✅ All checks passed!"
```

### 2. 缓存清理脚本

```bash
#!/bin/bash
# scripts/clear_cache.sh

echo "🗑️  Clearing cache..."

# 清理所有用户缓存
redis-cli --scan --pattern "cache:user:*" | xargs redis-cli del

# 清理统计缓存
redis-cli del cache:user:count:total

# 清理索引缓存
redis-cli --scan --pattern "cache:user:email:*" | xargs redis-cli del
redis-cli --scan --pattern "cache:user:username:*" | xargs redis-cli del

echo "✅ Cache cleared!"
```

### 3. 数据库备份脚本

```bash
#!/bin/bash
# scripts/backup_database.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/gin-demo"
BACKUP_FILE="$BACKUP_DIR/gin_demo_$DATE.sql"

echo "📦 Starting database backup..."

# 创建备份目录
mkdir -p $BACKUP_DIR

# 执行备份
pg_dump -h localhost -U postgres -d gin_demo > $BACKUP_FILE

# 压缩
gzip $BACKUP_FILE

# 保留最近 30 天的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

echo "✅ Backup completed: $BACKUP_FILE.gz"
```

---

## 📞 获取帮助

### 内部资源
- 📖 [架构文档](./ARCHITECTURE.md)
- 📖 [API 文档](./API.md)
- 📖 [RBAC 文档](./RBAC.md)

### 外部资源
- [Gin 问题排查](https://github.com/gin-gonic/gin/issues)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)
- [Redis 命令参考](https://redis.io/commands)
- [Go 性能优化](https://go.dev/doc/diagnostics)

### 联系支持
- 📧 技术支持: tech-support@example.com
- 💬 Slack: #gin-demo-support
- 📱 紧急热线: +86-xxx-xxxx-xxxx (仅生产故障)

---

**提示**: 
- 遇到问题先查看日志和监控
- 记录每次故障的原因和解决方案
- 定期更新本文档
- 建立故障知识库
