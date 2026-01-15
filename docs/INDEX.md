# 📚 文档索引

> **最后更新**: 2026-01-15  
> **项目版本**: v3.0  
> **生产就绪度**: 98%

---

## 🚀 快速开始

新手必读文档（按顺序阅读）：

1. [README.md](../README.md) - 项目介绍和快速开始
2. [CONFIGURATION.md](./CONFIGURATION.md) - 配置文件详细说明
3. [API.md](./API.md) - API 接口文档
4. [DEPLOYMENT-CHECKLIST.md](./DEPLOYMENT-CHECKLIST.md) - 部署清单

---

## 📖 核心文档

### 架构与设计

| 文档 | 说明 | 重要性 |
|------|------|--------|
| [ARCHITECTURE.md](./ARCHITECTURE.md) | 项目架构设计 | ⭐⭐⭐⭐⭐ |
| [SENIOR_ARCHITECT_ANALYSIS.md](./SENIOR_ARCHITECT_ANALYSIS.md) | 资深架构师深度分析 | ⭐⭐⭐⭐⭐ |
| [ADR-TEMPLATE.md](./ADR-TEMPLATE.md) | 架构决策记录模板 | ⭐⭐⭐⭐ |

### API 与接口

| 文档 | 说明 | 重要性 |
|------|------|--------|
| [API.md](./API.md) | RESTful API 接口文档 | ⭐⭐⭐⭐⭐ |
| [API_VERSIONING.md](./API_VERSIONING.md) | API 版本管理策略 | ⭐⭐⭐⭐ |
| [swagger.md](./swagger.md) | Swagger 文档说明 | ⭐⭐⭐ |

### 功能模块

| 文档 | 说明 | 重要性 |
|------|------|--------|
| [RBAC.md](./RBAC.md) | 角色权限控制 | ⭐⭐⭐⭐⭐ |
| [HTTP_SECURITY.md](./HTTP_SECURITY.md) | HTTP 安全防护 | ⭐⭐⭐⭐⭐ |
| [TASK_SCHEDULER.md](./TASK_SCHEDULER.md) | 定时任务调度 | ⭐⭐⭐⭐ |
| [TASK_FEATURE.md](./TASK_FEATURE.md) | 任务系统功能 | ⭐⭐⭐ |

### 数据库与缓存

| 文档 | 说明 | 重要性 |
|------|------|--------|
| [MYSQL_MIGRATION.md](./MYSQL_MIGRATION.md) | MySQL 迁移指南 | ⭐⭐⭐⭐⭐ |
| [reports/MySQL_Redis哨兵迁移完成报告.md](./reports/MySQL_Redis哨兵迁移完成报告.md) | MySQL + Redis 哨兵迁移报告 | ⭐⭐⭐⭐⭐ |
| [../pkg/cache/README.md](../pkg/cache/README.md) | 缓存管理器使用指南 | ⭐⭐⭐⭐ |
| [../pkg/database/README.md](../pkg/database/README.md) | 数据库使用指南 | ⭐⭐⭐⭐ |

---

## 🛠️ 开发指南

### 代码质量

| 文档 | 说明 | 重要性 |
|------|------|--------|
| [CODE_QUALITY_REPORT.md](./CODE_QUALITY_REPORT.md) | 代码质量报告 | ⭐⭐⭐⭐ |
| [CODE_REVIEW.md](./CODE_REVIEW.md) | 代码审查指南 | ⭐⭐⭐⭐ |
| [SECURITY_IMPROVEMENTS.md](./SECURITY_IMPROVEMENTS.md) | 安全改进记录 | ⭐⭐⭐⭐ |

### 模块说明

| 文档 | 说明 | 重要性 |
|------|------|--------|
| [../pkg/README.md](../pkg/README.md) | pkg 包总览 | ⭐⭐⭐⭐ |
| [../pkg/task/README.md](../pkg/task/README.md) | 任务调度器 | ⭐⭐⭐ |
| [../pkg/validator/README.md](../pkg/validator/README.md) | 数据验证器 | ⭐⭐⭐ |
| [../internal/app/middleware/README.md](../internal/app/middleware/README.md) | 中间件说明 | ⭐⭐⭐⭐ |
| [../internal/response/README.md](../internal/response/README.md) | 响应格式化 | ⭐⭐⭐ |

---

## 🚀 运维指南

### 部署与监控

| 文档 | 说明 | 重要性 |
|------|------|--------|
| [DEPLOYMENT-CHECKLIST.md](./DEPLOYMENT-CHECKLIST.md) | 部署前检查清单 | ⭐⭐⭐⭐⭐ |
| [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) | 故障排查手册 | ⭐⭐⭐⭐⭐ |
| [CONFIGURATION.md](./CONFIGURATION.md) | 配置文件详解 | ⭐⭐⭐⭐⭐ |

---

## 📜 版本历史与报告

| 文档 | 说明 | 版本 |
|------|------|------|
| [../CHANGELOG.md](../CHANGELOG.md) | 完整变更日志 | v1.0 - v4.0 |
| [reports/README.md](./reports/README.md) | 📊 项目报告汇总 | 所有报告 |

---

## 📦 归档文档

历史版本文档已移至 [archive/](./archive/) 目录：

- CLEANUP_SUMMARY.md
- CODE_REFACTORING.md
- FINAL_IMPROVEMENTS.md
- FINAL_REFACTORING.md
- FINAL_SUMMARY.md
- IMPROVEMENTS_SUMMARY.md
- OPTIMIZATION_V2.md
- REFACTORING_COMPLETE.md
- REFACTORING_V3.md
- V3_COMPLETE.md
- README_V3.md
- README_V4.md
- 优化总结.md
- 最终分析报告.md
- 架构分析总结报告.md

---

## 🎯 推荐阅读路径

### 路径 1: 新开发者入门

```
README.md → CONFIGURATION.md → API.md → ARCHITECTURE.md
→ pkg/README.md → RBAC.md
```

### 路径 2: 架构师/Tech Lead

```
SENIOR_ARCHITECT_ANALYSIS.md → ARCHITECTURE.md → ADR-TEMPLATE.md
→ MYSQL_MIGRATION.md → API_VERSIONING.md
```

### 路径 3: 运维/SRE

```
DEPLOYMENT-CHECKLIST.md → CONFIGURATION.md → TROUBLESHOOTING.md
→ MySQL_Redis哨兵迁移完成报告.md → HTTP_SECURITY.md
```

### 路径 4: 质量保障/测试

```
CODE_QUALITY_REPORT.md → API.md → TROUBLESHOOTING.md
→ SECURITY_IMPROVEMENTS.md
```

---

## 🔍 文档搜索提示

- **架构相关**: ARCHITECTURE, SENIOR_ARCHITECT_ANALYSIS, ADR
- **API 相关**: API.md, API_VERSIONING, swagger
- **安全相关**: SECURITY, RBAC, HTTP_SECURITY
- **数据库相关**: MYSQL_MIGRATION, database/README
- **缓存相关**: cache/README
- **运维相关**: DEPLOYMENT, TROUBLESHOOTING, CONFIGURATION
- **代码质量**: CODE_QUALITY, CODE_REVIEW

---

## 💡 文档贡献

欢迎补充和完善文档：

1. 使用 [ADR-TEMPLATE.md](./ADR-TEMPLATE.md) 记录重要决策
2. 保持文档与代码同步
3. 添加清晰的示例代码
4. 标注文档更新时间

---

**项目状态**：✅ 生产就绪（98%）  
**技术栈**：Go 1.25 + Gin + MySQL + Redis 哨兵  
**部署方式**：Docker Compose / Kubernetes
