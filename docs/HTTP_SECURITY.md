# HTTP 安全性与传输优化

> 完善的 HTTP 安全头和传输优化配置

---

## 🔒 安全头概览

### 已实现的安全头

| 安全头 | 作用 | 配置值 |
|--------|------|--------|
| **X-Content-Type-Options** | 防止 MIME 类型嗅探 | `nosniff` |
| **X-Frame-Options** | 防止点击劫持 | `DENY` / `SAMEORIGIN` |
| **X-XSS-Protection** | 启用 XSS 过滤器 | `1; mode=block` |
| **Referrer-Policy** | 控制 Referer 发送 | `strict-origin-when-cross-origin` |
| **Content-Security-Policy** | 防止 XSS/注入攻击 | 可配置 |
| **Strict-Transport-Security** | 强制 HTTPS | `max-age=31536000` |
| **Permissions-Policy** | 控制浏览器 API | 可配置 |
| **X-Permitted-Cross-Domain-Policies** | 跨域策略 | `none` |
| **X-Download-Options** | 防止自动执行下载 | `noopen` |

---

## ⚙️ 配置说明

### 1. 安全头配置

在 `config.yaml` 中配置：

```yaml
security:
  headers:
    enabled: true
    
    # HSTS (HTTP Strict Transport Security)
    enable_hsts: true
    hsts_max_age: 31536000  # 1年
    hsts_include_subdomains: true
    
    # CSP (Content Security Policy)
    enable_csp: true
    csp_policy: "default-src 'self'; script-src 'self';"
    
    # Frame Options
    enable_frame_options: true
    frame_options: "DENY"  # DENY / SAMEORIGIN
```

### 2. 压缩传输配置

```yaml
security:
  # Gzip 压缩
  enable_compression: true
  compression_level: 5  # -1 (默认), 0-9
```

### 3. TLS/HTTPS 配置

```yaml
security:
  tls:
    enabled: true
    cert_file: "/path/to/cert.pem"
    key_file: "/path/to/key.pem"
    min_version: "1.2"  # 1.2 / 1.3
```

---

## 📊 安全等级对比

### 开发环境 (config.dev.yaml)

```yaml
security:
  headers:
    enabled: true
    enable_hsts: false  # ❌ 关闭（不需要 HTTPS）
    enable_csp: true    # ⚠️  宽松策略
    csp_policy: "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval';"
  
  enable_compression: true
  compression_level: 5
  
  tls:
    enabled: false  # ❌ 关闭
```

**特点**: 
- 🟡 适度安全，方便开发
- 允许内联脚本和 eval
- 不强制 HTTPS

### 生产环境 (config.prod.yaml)

```yaml
security:
  headers:
    enabled: true
    enable_hsts: true   # ✅ 启用
    hsts_max_age: 31536000
    enable_csp: true    # ✅ 严格策略
    csp_policy: "default-src 'self'; script-src 'self'; style-src 'self';"
  
  enable_compression: true
  compression_level: 6
  
  tls:
    enabled: true       # ✅ 强制 HTTPS
    min_version: "1.2"
```

**特点**:
- 🔴 严格安全策略
- 禁止内联脚本
- 强制 HTTPS
- 最高压缩级别

---

## 🛡️ 安全头详解

### 1. X-Content-Type-Options

**作用**: 防止浏览器 MIME 类型嗅探

**配置**: `nosniff`

**攻击场景**:
- 攻击者上传恶意 `.txt` 文件
- 浏览器误判为 JavaScript 并执行

**防护**:
```
X-Content-Type-Options: nosniff
```

---

### 2. X-Frame-Options

**作用**: 防止点击劫持（Clickjacking）

**配置选项**:
- `DENY` - 完全禁止嵌入 iframe
- `SAMEORIGIN` - 只允许同源嵌入

**攻击场景**:
- 恶意网站将你的页面嵌入透明 iframe
- 诱导用户点击

**防护**:
```
X-Frame-Options: DENY
```

---

### 3. X-XSS-Protection

**作用**: 启用浏览器 XSS 过滤器

**配置**: `1; mode=block`

**说明**:
- `1` = 启用
- `mode=block` = 完全阻止页面加载

**防护**:
```
X-XSS-Protection: 1; mode=block
```

---

### 4. Content-Security-Policy (CSP)

**作用**: 防止 XSS、注入攻击

**常用策略**:

#### 严格策略（生产）
```
default-src 'self'; 
script-src 'self'; 
style-src 'self'; 
img-src 'self' data:; 
font-src 'self';
```

#### 宽松策略（开发）
```
default-src 'self'; 
script-src 'self' 'unsafe-inline' 'unsafe-eval'; 
style-src 'self' 'unsafe-inline';
```

**CSP 指令说明**:
| 指令 | 说明 | 示例 |
|------|------|------|
| `default-src` | 默认策略 | `'self'` |
| `script-src` | JavaScript 来源 | `'self' https://cdn.example.com` |
| `style-src` | CSS 来源 | `'self' 'unsafe-inline'` |
| `img-src` | 图片来源 | `'self' data: https:` |
| `font-src` | 字体来源 | `'self'` |
| `connect-src` | AJAX/WebSocket | `'self' https://api.example.com` |

**特殊值**:
- `'self'` - 同源
- `'none'` - 禁止
- `'unsafe-inline'` - 允许内联
- `'unsafe-eval'` - 允许 eval

---

### 5. Strict-Transport-Security (HSTS)

**作用**: 强制浏览器使用 HTTPS

**配置**: `max-age=31536000; includeSubDomains`

**重要提示**:
- ⚠️ **只在 HTTPS 下启用**
- ⚠️ **测试后再启用** - 一旦设置，无法在客户端撤销

**防护**:
```
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

---

### 6. Referrer-Policy

**作用**: 控制 Referer 头发送

**配置选项**:
- `no-referrer` - 不发送
- `strict-origin` - 只发送源
- `strict-origin-when-cross-origin` - 跨域只发送源（推荐）

**防护**:
```
Referrer-Policy: strict-origin-when-cross-origin
```

---

### 7. Permissions-Policy

**作用**: 控制浏览器功能访问

**配置示例**:
```
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

**常用功能**:
- `geolocation` - 地理位置
- `microphone` - 麦克风
- `camera` - 摄像头
- `payment` - 支付 API
- `usb` - USB 访问

---

## 🗜️ Gzip 压缩

### 工作原理

1. 客户端发送 `Accept-Encoding: gzip`
2. 服务器压缩响应
3. 客户端解压缩

### 压缩级别

| 级别 | 说明 | 压缩率 | CPU 消耗 |
|------|------|--------|----------|
| -1 | 默认 | 中 | 中 |
| 0 | 不压缩 | 0% | 极低 |
| 1 | 最快 | 低 | 低 |
| 5 | 推荐 | 中 | 中 |
| 9 | 最佳 | 高 | 高 |

### 排除策略

**自动排除**:
- 已压缩文件: `.jpg`, `.png`, `.gif`, `.zip`, `.gz`
- 视频/音频: `.mp4`, `.mp3`, `.avi`
- Prometheus 指标: `/metrics`

**压缩效果**:
- JSON 响应: 70-80% 压缩率
- HTML: 60-70% 压缩率
- CSS/JS: 50-60% 压缩率

---

## 🔐 HTTPS/TLS 配置

### 生成自签名证书（开发）

```bash
# 生成私钥
openssl genrsa -out key.pem 2048

# 生成证书
openssl req -new -x509 -key key.pem -out cert.pem -days 365
```

### 使用 Let's Encrypt（生产）

```bash
# 安装 certbot
sudo apt install certbot

# 获取证书
sudo certbot certonly --standalone -d yourdomain.com
```

### TLS 版本说明

| 版本 | 状态 | 说明 |
|------|------|------|
| TLS 1.0 | ❌ 已弃用 | 不安全 |
| TLS 1.1 | ❌ 已弃用 | 不安全 |
| TLS 1.2 | ✅ 推荐 | 安全，兼容性好 |
| TLS 1.3 | ✅ 最佳 | 最安全，性能最好 |

---

## 📝 使用示例

### 在代码中使用

```go
// main.go
import "gin_demo/internal/app/middleware"

// 方式 1: 使用配置
engine.Use(middleware.Security(middleware.SecurityConfig{
    EnableHSTS:  true,
    HSTSMaxAge:  31536000,
    EnableCSP:   true,
    CSPPolicy:   "default-src 'self';",
}))

// 方式 2: 简化版（开发）
engine.Use(middleware.SecureHeaders())

// 方式 3: 从配置加载（推荐）
engine.Use(configureSecurityMiddleware(cfg))
```

### 测试安全头

```bash
# 检查响应头
curl -I http://localhost:8080/api/v1/users

# 预期输出
HTTP/1.1 200 OK
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Content-Security-Policy: default-src 'self'; ...
```

### 在线检测工具

- [Security Headers](https://securityheaders.com/) - 安全头检测
- [SSL Labs](https://www.ssllabs.com/ssltest/) - TLS 配置检测
- [CSP Evaluator](https://csp-evaluator.withgoogle.com/) - CSP 策略检测

---

## 🎯 安全检查清单

### 部署前检查

- [ ] 启用所有安全头
- [ ] 配置严格的 CSP 策略
- [ ] 启用 HSTS（仅 HTTPS）
- [ ] 配置 TLS 1.2+
- [ ] 启用 Gzip 压缩
- [ ] 测试所有路由的响应头
- [ ] 使用在线工具检测
- [ ] 配置 CORS 白名单
- [ ] 限制请求体大小
- [ ] 启用 Rate Limiting

### 定期审计

- [ ] 每月检查安全头配置
- [ ] 每季度更新 TLS 证书
- [ ] 每季度审查 CSP 策略
- [ ] 监控安全漏洞公告

---

## 🚨 常见错误

### 1. CSP 策略过严

**问题**: 页面无法加载 CDN 资源

**解决**:
```yaml
csp_policy: "default-src 'self'; script-src 'self' https://cdn.jsdelivr.net;"
```

### 2. HSTS 误配置

**问题**: 在 HTTP 环境下启用 HSTS

**解决**: 只在生产环境（HTTPS）启用
```yaml
enable_hsts: false  # 开发环境
```

### 3. Frame Options 冲突

**问题**: 需要在 iframe 中嵌入

**解决**: 使用 `SAMEORIGIN`
```yaml
frame_options: "SAMEORIGIN"
```

---

## 📖 参考资料

- [OWASP Secure Headers Project](https://owasp.org/www-project-secure-headers/)
- [MDN - CSP](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)
- [MDN - HSTS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security)
- [Security Headers 最佳实践](https://securityheaders.com/)

---

## 🎓 最佳实践

1. **渐进式加固**: 从宽松到严格
2. **测试为先**: 在开发环境充分测试
3. **监控告警**: 关注 CSP 违规报告
4. **定期更新**: 跟进安全标准更新
5. **文档记录**: 记录所有安全配置

---

**现在项目已具备企业级 HTTP 安全性！** 🔒
