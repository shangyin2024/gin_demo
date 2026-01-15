# ADR-XXX: [简短的决策标题]

**日期**: YYYY-MM-DD  
**状态**: [提议中 | 已采纳 | 已废弃 | 已替代]  
**决策者**: [姓名/团队]  
**替代者**: [如果被替代，指向新的 ADR]

---

## 背景

描述需要做决策的背景和问题。

**示例**:
> 项目需要实现用户认证功能，需要选择认证方案。当前用户量约 10万，预计增长到 100万。

---

## 决策

明确说明做出的决策。

**示例**:
> 使用 JWT（JSON Web Token）进行用户认证，而不使用 Session。

---

## 考虑的方案

列出所有考虑过的方案及其优缺点。

### 方案 1: JWT

**优点**:
- 无状态，易于横向扩展
- 包含用户信息，减少查询
- 跨域友好

**缺点**:
- 无法主动失效（需要黑名单）
- Token 体积较大
- 密钥泄露风险

### 方案 2: Session + Redis

**优点**:
- 可以主动失效
- Token 体积小
- 灵活性高

**缺点**:
- 需要 Redis 存储
- 查询开销
- 跨域复杂

### 方案 3: OAuth 2.0

**优点**:
- 标准协议
- 支持第三方登录
- 安全性高

**缺点**:
- 复杂度高
- 实现成本高
- 对于内部系统过度设计

---

## 决策原因

详细说明为什么选择这个方案。

**示例**:
> 1. 项目需要横向扩展到多实例，JWT 的无状态特性完美契合
> 2. 用户量在可控范围内，JWT 无法主动失效的问题可以通过黑名单解决
> 3. 实现简单，开发效率高
> 4. 性能优秀，无额外的 Redis 查询

---

## 后果

### 正面后果
- 服务无状态，易于扩展
- 性能优秀
- 实现简单

### 负面后果
- 需要实现 Token 黑名单
- Token 泄露后难以撤销
- 需要定期更换密钥

### 风险缓解
1. 实现 Token 黑名单（Redis）
2. 设置较短的过期时间（24小时）
3. 使用 HTTPS 防止 Token 泄露
4. 添加 Token 刷新机制

---

## 实现细节

### 依赖
```go
github.com/golang-jwt/jwt/v5
```

### 配置
```yaml
jwt:
  secret: your-secret-key
  expiration: 24h
```

### 代码示例
```go
// 生成 Token
token, err := jwtManager.GenerateToken(userID)

// 验证 Token
claims, err := jwtManager.ValidateToken(token)
```

---

## 相关 ADR

- [ADR-002: 使用 Redis 作为缓存](./002-use-redis-cache.md)
- [ADR-005: RBAC 权限设计](./005-rbac-design.md)

---

## 参考资料

- [JWT 规范](https://jwt.io/)
- [RFC 7519](https://tools.ietf.org/html/rfc7519)
- [OWASP JWT Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_for_Java_Cheat_Sheet.html)

---

## 历史

| 日期 | 修改 | 修改人 |
|------|------|--------|
| 2024-01-01 | 创建 | Alice |
| 2024-06-15 | 添加黑名单机制 | Bob |
| 2026-01-15 | 扩展为 RBAC | Charlie |
