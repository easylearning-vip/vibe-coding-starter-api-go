# Vibe Coding Starter Go API - 文档中心

欢迎来到 Vibe Coding Starter Go API 的文档中心！这里包含了项目二次开发所需的所有文档资源。

## 📚 文档导航

### 🚀 [开发指南](./development-guide.md)
**适合人群**: 初学者、AI 工程师、后端开发者

**内容概览**:
- 项目概述和核心特性
- 完整的目录结构说明
- 技术栈详细介绍
- 开发环境搭建指南
- 代码规范和最佳实践
- 开发流程和分支管理
- AI Prompting 工程指南
- 常见开发任务示例
- 部署和运维指南

**为什么要读**: 这是入门必读文档，帮助你快速理解项目结构和开发规范。

### 🏗️ [架构文档](./architecture.md)
**适合人群**: 架构师、高级开发者、技术负责人

**内容概览**:
- 系统架构概览图
- 清洁架构设计原理
- 依赖注入架构详解
- 数据流和模块依赖关系
- 核心组件说明
- 设计原则和扩展点
- 性能和安全架构
- 监控和可观测性
- 部署架构设计

**为什么要读**: 深入理解系统设计思路，为架构决策和系统扩展提供指导。

### 🔌 [API 设计文档](./api-design.md)
**适合人群**: 前端开发者、API 集成开发者、测试工程师

**内容概览**:
- RESTful API 设计原则
- 完整的路由结构
- 中间件架构设计
- 数据模型定义
- 请求响应格式规范
- 认证授权机制
- 错误处理和状态码
- 输入验证规则
- 性能优化策略

**为什么要读**: 了解 API 接口规范，便于前后端协作和第三方集成。

### 🧪 [测试指南](./testing-guide.md)
**适合人群**: 测试工程师、质量保证工程师、开发者

**内容概览**:
- 多层次测试策略
- 单元测试最佳实践
- 集成测试实现方法
- 端到端测试流程
- Mock 对象使用指南
- 测试工具链介绍
- 覆盖率报告分析
- CI/CD 测试集成

**为什么要读**: 掌握测试方法论，确保代码质量和系统稳定性。

## 🎯 快速开始指南

### 对于初学者
1. 📖 先阅读 [开发指南](./development-guide.md) 了解项目基础
2. 🏗️ 浏览 [架构文档](./architecture.md) 理解系统设计
3. 🔧 按照开发指南搭建环境并运行项目
4. 🧪 学习 [测试指南](./testing-guide.md) 编写第一个测试

### 对于前端开发者
1. 🔌 重点阅读 [API 设计文档](./api-design.md)
2. 📖 参考开发指南中的 "API 端点详细设计" 部分
3. 🔧 启动本地开发环境进行接口调试
4. 📊 访问 Swagger 文档: `http://localhost:8080/swagger/index.html`

### 对于 AI 工程师
1. 📖 阅读开发指南中的 "AI Prompting 工程指南" 部分
2. 🏗️ 理解架构文档中的模块依赖关系
3. 🔌 熟悉 API 设计规范和数据模型
4. 🧪 学习测试驱动开发方法

### 对于运维工程师
1. 📖 关注开发指南中的 "部署指南" 部分
2. 🏗️ 理解架构文档中的 "部署架构" 部分
3. 🔧 熟悉 Docker 和 Kubernetes 配置
4. 📊 了解监控和日志收集方案

## 🛠️ 开发环境快速搭建

```bash
# 1. 克隆项目
git clone <repository-url>
cd vibe-coding-starter-go-api

# 2. 启动开发环境
make dev-docker

# 3. 运行数据库迁移
make migrate-up

# 4. 启动应用
make run-docker

# 5. 访问 API 文档
open http://localhost:8080/swagger/index.html
```

## 📋 项目特性一览

### ✅ 核心功能
- [x] 用户注册和认证系统
- [x] JWT Token 认证
- [x] 文章管理 CRUD
- [x] 文件上传和管理
- [x] RESTful API 设计
- [x] Swagger API 文档

### ✅ 技术特性
- [x] 清洁架构设计
- [x] 依赖注入 (Uber FX)
- [x] 多数据库支持 (MySQL/PostgreSQL/SQLite)
- [x] Redis 缓存集成
- [x] 结构化日志 (Zap)
- [x] 配置管理 (Viper)
- [x] 中间件系统 (认证/CORS/限流/安全)

### ✅ 开发体验
- [x] 完整的单元测试覆盖
- [x] 集成测试和端到端测试
- [x] Docker 容器化支持
- [x] Kubernetes 部署配置
- [x] Makefile 构建脚本
- [x] 开发环境一键启动

### ✅ 生产就绪
- [x] 健康检查端点
- [x] 优雅关闭机制
- [x] 错误处理和恢复
- [x] 请求限流和安全防护
- [x] 数据库迁移管理
- [x] 监控和日志集成

## 🔧 常用命令速查

### 开发命令
```bash
make dev-setup          # 设置开发环境
make build              # 构建应用
make test               # 运行测试
make test-coverage      # 生成覆盖率报告
make run-local          # 本地运行
make run-docker         # Docker 环境运行
```

### 数据库命令
```bash
make migrate-up         # 执行数据库迁移
make migrate-down       # 回滚数据库迁移
make migrate-version    # 查看迁移版本
```

### Docker 命令
```bash
make docker-build       # 构建 Docker 镜像
make docker-run         # 运行 Docker 容器
make dev-docker         # 启动开发环境
```

### Kubernetes 命令
```bash
make k8s-deploy         # 部署到 K8s
make k8s-status         # 查看部署状态
make k8s-logs           # 查看应用日志
make k8s-clean          # 清理部署
```

## 🤖 AI 辅助开发

### 项目上下文提示词
```
这是一个采用清洁架构的 Go Web API 项目：
- 使用 Gin 框架和 GORM ORM
- 采用依赖注入模式 (Uber FX)
- 分层架构：Handler -> Service -> Repository
- 支持多种数据库和 Redis 缓存
- 完整的测试覆盖和容器化部署
```

### 功能开发提示词模板
```
请为 [功能名称] 实现完整的 CRUD 功能，包括：
1. 数据模型定义 (internal/model/)
2. Repository 接口和实现 (internal/repository/)
3. Service 业务逻辑 (internal/service/)
4. Handler HTTP 处理 (internal/handler/)
5. 路由注册和中间件配置
6. 完整的单元测试和集成测试
7. Swagger API 文档注释

要求遵循项目现有的代码风格和架构模式。
```

## 📞 获取帮助

### 文档问题
- 如果文档有不清楚的地方，请查看对应的源代码
- 参考 `test/` 目录下的测试用例了解使用方法
- 查看 `dev-tutorial/` 目录下的教程文档

### 开发问题
- 查看项目的 Issue 和 PR 历史
- 参考 `tools/` 目录下的开发工具
- 使用 `make help` 查看所有可用命令

### 部署问题
- 查看 `deploy/` 目录下的部署配置
- 参考 `dev-tutorial/` 中的环境搭建指南
- 检查 Docker 和 Kubernetes 配置文件

## 🎉 贡献指南

### 代码贡献
1. Fork 项目并创建功能分支
2. 遵循项目的代码规范和架构模式
3. 编写完整的测试用例
4. 更新相关文档
5. 提交 Pull Request

### 文档贡献
1. 发现文档错误或不完整的地方
2. 创建 Issue 或直接提交 PR
3. 保持文档的简洁性和准确性
4. 确保示例代码可以正常运行

## 📄 许可证

本项目采用 MIT 许可证，详情请查看 LICENSE 文件。

---

## 🌟 项目亮点

- **🤖 AI 友好**: 完全由 AI 工具开发，零人工代码编写
- **📚 文档完善**: 详细的开发文档和 API 文档
- **🧪 测试完整**: 多层次测试策略，高覆盖率
- **🚀 生产就绪**: 企业级特性，可直接用于生产环境
- **🔧 开发友好**: 一键启动开发环境，丰富的开发工具
- **📦 容器化**: 完整的 Docker 和 Kubernetes 支持

这个项目不仅是一个 Go API 模板，更是 Vibe Coding 开发方法论的最佳实践示范。无论你是初学者还是经验丰富的开发者，都能从中获得价值并快速上手进行二次开发。

**开始你的 Vibe Coding 之旅吧！** 🚀
