# Vibe Coding Starter API 本地开发环境搭建手册

欢迎使用 Vibe Coding Starter API 开发环境搭建手册！本手册提供了两种主流的本地开发环境搭建方式：**Docker Compose** 和 **k3d**，帮助您快速启动包含 MySQL 和 Redis 的完整开发环境。

## 🚀 快速选择

| 特性 | Docker Compose | k3d |
|------|----------------|-----|
| **学习曲线** | 简单易用 | 需要 Kubernetes 基础 |
| **资源消耗** | 较低 | 中等 |
| **启动速度** | 快速 | 中等 |
| **扩展性** | 有限 | 优秀 |
| **生产环境相似度** | 中等 | 高 |
| **适用场景** | 日常开发、快速原型 | 云原生开发、学习 K8s |

### 推荐选择

- **新手开发者** 或 **快速开发**: 选择 [Docker Compose](#docker-compose-方式)
- **云原生开发** 或 **学习 Kubernetes**: 选择 [k3d 方式](#k3d-方式)

## 📋 环境概览

两种方式都将为您提供：

### 数据库服务
- **MySQL 8.0.33**: 主要数据库，包含示例数据
- **Redis 7**: 缓存和会话存储
- **PostgreSQL 15** (可选): 备选数据库

### 管理工具
- **phpMyAdmin**: MySQL 数据库管理界面
- **Redis Commander**: Redis 数据管理界面

### 网络配置
- 所有服务通过内部网络互联
- 外部端口映射便于本地访问
- 健康检查确保服务可用性

## 🐳 Docker Compose 方式

Docker Compose 是最简单快速的开发环境搭建方式，适合大多数开发场景。

### 特点
- ✅ 配置简单，一键启动
- ✅ 资源消耗低
- ✅ 启动速度快
- ✅ 适合日常开发
- ✅ 包含管理工具

### 快速开始

```bash
# 1. 进入 Docker Compose 目录
cd vibe-coding-starter-api-go/tutorials/develop/docker-compose

# 2. 启动基础服务（MySQL + Redis）
docker compose -f docker-compose.dev.yml up -d

# 3. 启动完整服务（包含管理工具）
docker compose -f docker-compose.dev.yml --profile tools up -d

# 4. 查看服务状态
docker compose -f docker-compose.dev.yml ps
```

### 服务访问

| 服务 | 地址 | 用户名 | 密码 |
|------|------|--------|------|
| MySQL | localhost:3306 | vibe_user | vibe_password |
| Redis | localhost:6379 | - | - |
| phpMyAdmin | http://localhost:8080 | root | rootpassword |
| Redis Commander | http://localhost:8081 | - | - |

### 应用配置

使用 Docker Compose 环境时，应用程序应使用以下配置文件：

```bash
# 启动应用程序
go run cmd/server/main.go -config configs/config-docker.yaml
```

### 详细文档

完整的 Docker Compose 使用指南请参考：[Docker Compose 详细文档](docker-compose/README.md)

### MySQL 客户端连接

如需使用本地 MySQL 客户端连接，请参考：[MySQL 客户端连接指南](mysql-client-guide.md)

## ☸️ k3d 方式

k3d 提供了真正的 Kubernetes 开发环境，适合云原生应用开发和学习。

### 特点
- ✅ 真实的 Kubernetes 环境
- ✅ 支持云原生开发模式
- ✅ 优秀的扩展性
- ✅ 接近生产环境
- ✅ 学习 Kubernetes 的最佳选择

### 快速开始

```bash
# 1. 进入 k3d 目录
cd vibe-coding-starter-api-go/tutorials/develop/k3d

# 2. 使用配置文件创建 k3d 集群
k3d cluster create --config k3d-cluster.yaml

# 3. 验证集群状态
kubectl cluster-info
kubectl get nodes

# 4. 部署命名空间和基础配置
kubectl apply -f manifests/namespace.yaml

# 5. 部署 MySQL 服务
kubectl apply -f manifests/mysql.yaml

# 6. 部署 Redis 服务
kubectl apply -f manifests/redis.yaml

# 7. 等待所有 Pod 就绪
kubectl wait --for=condition=ready pod --all -n vibe-dev --timeout=300s

# 8. 查看部署状态
kubectl get all -n vibe-dev
```

### 服务访问

| 服务 | 集群内地址 | 外部地址 | 用户名 | 密码 |
|------|------------|----------|--------|------|
| MySQL | mysql.vibe-dev.svc.cluster.local:3306 | localhost:30306 | vibe_user | vibe_password |
| Redis | redis.vibe-dev.svc.cluster.local:6379 | localhost:30379 | - | - |

### 应用配置

使用 k3d 环境时，应用程序应使用以下配置文件：

```bash
# 启动应用程序
go run cmd/server/main.go -config configs/config-k3d.yaml
```

### 详细文档

完整的 k3d 使用指南请参考：[k3d 详细文档](k3d/README.md)

### MySQL 客户端连接

如需使用本地 MySQL 客户端连接，请参考：[MySQL 客户端连接指南](mysql-client-guide.md)

## 🔧 配置文件说明

项目提供了针对不同环境的配置文件：

### configs/config-docker.yaml
- 适用于 Docker Compose 环境
- 数据库连接：localhost:3306
- Redis 连接：localhost:6379
- 开发模式优化配置

### configs/config-k3d.yaml
- 适用于 k3d 环境
- 支持集群内和外部连接
- 数据库连接：localhost:30306
- Redis 连接：localhost:30379
- Kubernetes 原生配置

### configs/config.yaml
- 默认配置文件
- 使用 SQLite 数据库
- 适合快速测试

## 🛠️ 开发工作流

### 1. 环境准备

选择并启动开发环境：

```bash
# Docker Compose 方式
cd tutorials/develop/docker-compose
docker compose -f docker-compose.dev.yml up -d

# 或 k3d 方式
cd tutorials/develop/k3d
k3d cluster create --config k3d-cluster.yaml
kubectl apply -f manifests/namespace.yaml
kubectl apply -f manifests/mysql.yaml
kubectl apply -f manifests/redis.yaml
```

### 2. 数据库迁移

```bash
# 回到项目根目录
cd ../..

# 运行数据库迁移 (注意使用正确的命令语法)
go run cmd/migrate/main.go -c configs/config-docker.yaml up
# 或者使用 k3d 配置
go run cmd/migrate/main.go -c configs/config-k3d.yaml up

# 查看迁移状态
go run cmd/migrate/main.go -c configs/config-k3d.yaml version

# 其他迁移命令
go run cmd/migrate/main.go -c configs/config-k3d.yaml down    # 回滚最后一个迁移
go run cmd/migrate/main.go -c configs/config-k3d.yaml fresh   # 重新创建所有表
go run cmd/migrate/main.go -c configs/config-k3d.yaml drop    # 删除所有表
```

**迁移说明：**
- 使用 `-c` 或 `--config` 参数指定配置文件
- `up` 命令执行所有待执行的迁移
- `version` 命令显示当前迁移版本
- 迁移文件位于 `migrations/mysql/` 目录
- 包含初始数据库结构和示例数据

### 3. 启动应用

```bash
# 安装依赖
go mod tidy

# 使用对应的配置文件启动应用
go run cmd/server/main.go -c configs/config-docker.yaml
# 或
go run cmd/server/main.go -c configs/config-k3d.yaml
```

### 4. 开发和测试

```bash
# 运行测试
go test ./...

# 运行特定测试
go test ./test/...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔍 故障排除

### 常见问题

#### 1. 端口冲突
如果遇到端口冲突，可以修改配置文件中的端口映射。

#### 2. 权限问题
确保 Docker 有足够的权限，可能需要将用户添加到 docker 组。

#### 3. 内存不足
确保系统有足够的内存，建议至少 4GB 可用内存。

#### 4. 网络连接问题
检查防火墙设置，确保相关端口未被阻止。

### 获取帮助

- 查看详细的故障排除指南：
  - [Docker Compose 故障排除](docker-compose/README.md#故障排除)
  - [k3d 故障排除](k3d/README.md#故障排除)

## 📚 进阶使用

### 数据持久化

两种环境都配置了数据持久化：
- MySQL 数据存储在 Docker 卷中
- Redis 数据支持 AOF 持久化
- 重启服务不会丢失数据

### 性能监控

可以启用监控组件：
- Prometheus 指标收集
- Grafana 可视化面板
- 应用性能监控

### 扩展服务

可以根据需要添加其他服务：
- Elasticsearch
- RabbitMQ
- MinIO (S3 兼容存储)

## 🤝 贡献

如果您在使用过程中发现问题或有改进建议，欢迎：

1. 提交 Issue
2. 发起 Pull Request
3. 完善文档

## 📄 许可证

本项目采用 MIT 许可证，详情请参考 LICENSE 文件。

---

**祝您开发愉快！** 🎉

如果您有任何问题，请查看详细文档或联系项目维护者。
