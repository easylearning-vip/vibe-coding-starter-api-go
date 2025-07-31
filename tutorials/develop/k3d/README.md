# k3d 开发环境搭建指南

本指南将帮助您使用 k3d 快速搭建 Vibe Coding Starter 的 Kubernetes 开发环境，包括 MySQL 和 Redis 服务。

## 目录

- [什么是 k3d](#什么是-k3d)
- [前置要求](#前置要求)
- [安装 k3d 和相关工具](#安装-k3d-和相关工具)
- [快速开始](#快速开始)
- [手动配置](#手动配置)
- [存储配置最佳实践](#存储配置最佳实践)
- [连接数据库](#连接数据库)
- [常用命令](#常用命令)
- [故障排除](#故障排除)
- [高级功能](#高级功能)

## 什么是 k3d

k3d 是一个轻量级的工具，可以在 Docker 中运行 k3s（轻量级 Kubernetes）集群。它非常适合：

- 本地开发和测试
- CI/CD 流水线
- 学习 Kubernetes
- 快速原型开发

## 前置要求

- 操作系统：Linux、macOS 或 Windows
- Docker 已安装并运行
- 至少 4GB 可用内存
- 至少 2GB 可用磁盘空间

## 安装 k3d 和相关工具

### 安装步骤

#### 1. 安装 k3d

```bash
# Linux/macOS
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

# 或者使用包管理器
# macOS
brew install k3d

# Ubuntu/Debian
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | TAG=v5.6.0 bash
```

#### 2. 安装 kubectl

```bash
# Linux
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/

# macOS
brew install kubectl

# 或者
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/
```

#### 3. 验证安装

```bash
# 检查 k3d 版本
k3d version

# 检查 kubectl 版本
kubectl version --client

# 检查 Docker 状态
docker info
```

## 快速开始

按照以下步骤快速搭建 k3d 开发环境：

### 1. 创建 k3d 集群

```bash
# 进入 k3d 目录
cd vibe-coding-starter-api-go/tutorials/develop/k3d

# 使用配置文件创建集群
k3d cluster create --config k3d-cluster.yaml

# 验证集群状态
kubectl cluster-info
kubectl get nodes
```

### 2. 部署服务

```bash
# 部署命名空间和基础配置
kubectl apply -f manifests/namespace.yaml

# 部署 MySQL 服务
kubectl apply -f manifests/mysql.yaml

# 部署 Redis 服务
kubectl apply -f manifests/redis.yaml

# 等待所有 Pod 就绪
kubectl wait --for=condition=ready pod --all -n vibe-dev --timeout=300s

# 查看部署状态
kubectl get all -n vibe-dev
```

### 3. 验证服务

```bash
# 检查 Pod 状态
kubectl get pods -n vibe-dev

# 检查服务状态
kubectl get svc -n vibe-dev

# 查看 Pod 日志（可选）
kubectl logs -l app=mysql -n vibe-dev
kubectl logs -l app=redis -n vibe-dev

```

## 详细配置说明

### k3d 集群配置

集群使用 `k3d-cluster.yaml` 配置文件创建，主要特性：

- **集群名称**: vibe-dev
- **节点配置**: 1 个 server 节点 + 2 个 agent 节点
- **端口映射**:
  - API Server: 6443
  - HTTP: 8080 → 80
  - HTTPS: 8443 → 443
  - NodePort 范围: 30000-32767
- **存储**: 本地路径存储（使用 `~/.local/share/k3d/vibe-dev-storage`，避免使用 `/tmp` 目录）
- **内置镜像仓库**: vibe-registry (localhost:5555)
- **自动 kubeconfig 管理**: 自动更新默认 kubeconfig 并切换上下文
- **Traefik Ingress Controller**: 已启用，支持域名路由和负载均衡

### 服务部署说明

#### 命名空间 (namespace.yaml)
- 创建 `vibe-dev` 命名空间
- 配置镜像拉取密钥
- 设置默认服务账户

#### MySQL 服务 (mysql.yaml)
- **镜像**: mysql:8.0.33
- **存储**: 5Gi 持久化存储
- **配置**: 优化的开发环境配置
- **初始化**: 自动创建数据库和测试数据
- **访问**: ClusterIP (集群内) + NodePort 30306 (外部)

#### Redis 服务 (redis.yaml)
- **镜像**: redis:7-alpine
- **存储**: 1Gi 持久化存储
- **配置**: 开发环境优化配置
- **持久化**: AOF 模式
- **访问**: ClusterIP (集群内) + NodePort 30379 (外部)

## 存储配置最佳实践

### 重要提醒：避免使用 /tmp 目录

本项目的 k3d 配置已经更新，不再使用 `/tmp/vibe-dev-storage` 作为存储路径，而是使用更安全的 `~/.local/share/k3d/vibe-dev-storage`。

### 为什么不使用 /tmp 目录？

1. **数据丢失风险**：`/tmp` 目录在系统重启时会被清空
2. **定期清理**：系统可能会定期清理 `/tmp` 目录中的旧文件
3. **权限问题**：可能出现权限冲突和安全问题
4. **备份困难**：备份脚本通常会跳过 `/tmp` 目录

### 推荐的存储配置

当前配置使用用户主目录下的持久化存储：

```yaml
volumes:
  - volume: ~/.local/share/k3d/vibe-dev-storage:/var/lib/rancher/k3s/storage
```

### 首次使用前的准备

```bash
# 创建存储目录
mkdir -p ~/.local/share/k3d/vibe-dev-storage

# 设置适当权限
chmod 755 ~/.local/share/k3d/vibe-dev-storage
```

### 更多详细信息

请参阅 [存储最佳实践文档](./STORAGE_BEST_PRACTICES.md) 了解：
- 详细的风险分析
- 其他存储方案选择
- 数据迁移指南
- 监控和维护建议

## 连接数据库

### 方式一：使用 NodePort（推荐）

服务已经配置了 NodePort，可以直接通过本地端口访问：

```bash
# MySQL 连接信息
Host: localhost
Port: 3306 (通过 k3d 端口映射)
Username: vibe_user
Password: vibe_password
Database: vibe_coding_starter

# Redis 连接信息
Host: localhost
Port: 6379 (通过 k3d 端口映射)
Password: (无)
```

#### 使用本地 MySQL 客户端连接

**快速连接**：

```bash
# 使用普通用户连接
mysql -h 127.0.0.1 -P 3306 -u vibe_user -pvibe_password vibe_coding_starter

# 使用 root 用户连接
mysql -h 127.0.0.1 -P 3306 -u root -prootpassword vibe_coding_starter
```

**配置 ~/.my.cnf 实现快速连接**：

创建或编辑 `~/.my.cnf` 文件：

```ini
[k3d]
host = 127.0.0.1
port = 3306
user = vibe_user
password = vibe_password
database = vibe_coding_starter
default-character-set = utf8mb4
protocol = tcp

[k3d-root]
host = 127.0.0.1
port = 3306
user = root
password = rootpassword
database = vibe_coding_starter
default-character-set = utf8mb4
protocol = tcp
```

然后使用配置连接：

```bash
# 使用配置文件连接
mysql --defaults-group-suffix=k3d

# 使用 root 配置连接
mysql --defaults-group-suffix=k3d-root

# 创建别名（添加到 ~/.bashrc 或 ~/.zshrc）
alias mysql-k3d='mysql --defaults-group-suffix=k3d'
alias mysql-k3d-root='mysql --defaults-group-suffix=k3d-root'
```

**连接测试**：

```bash
# 测试连接
mysql --defaults-group-suffix=k3d -e "SELECT 'k3d MySQL 连接成功!' as result;"

# 查看数据库
mysql --defaults-group-suffix=k3d -e "SHOW DATABASES;"

# 查看表和数据
mysql --defaults-group-suffix=k3d -e "SHOW TABLES; SELECT username, email FROM users LIMIT 3;"
```

### 方式二：使用端口转发

```bash
# MySQL 端口转发
kubectl port-forward svc/mysql 3306:3306 -n vibe-dev

# Redis 端口转发
kubectl port-forward svc/redis 6379:6379 -n vibe-dev

# 在另一个终端中连接
mysql -h localhost -P 3306 -u vibe_user -p
redis-cli -h localhost -p 6379
```

### 方式三：集群内连接

如果您的应用程序部署在 k3d 集群内，使用以下连接信息：

```yaml
# MySQL
Host: mysql.vibe-dev.svc.cluster.local
Port: 3306

# Redis  
Host: redis.vibe-dev.svc.cluster.local
Port: 6379
```

## 常用命令

### 集群管理

```bash
# 列出所有集群
k3d cluster list

# 启动集群
k3d cluster start vibe-dev

# 停止集群
k3d cluster stop vibe-dev

# 删除集群
k3d cluster delete vibe-dev
```

### 服务管理

```bash
# 查看所有资源
kubectl get all -n vibe-dev

# 查看 Pod 状态
kubectl get pods -n vibe-dev

# 查看服务状态
kubectl get svc -n vibe-dev

# 查看 Pod 日志
kubectl logs -f deployment/mysql -n vibe-dev
kubectl logs -f deployment/redis -n vibe-dev
```

### 数据库操作

```bash
# 连接到 MySQL Pod
kubectl exec -it deployment/mysql -n vibe-dev -- mysql -u vibe_user -p

# 连接到 Redis Pod
kubectl exec -it deployment/redis -n vibe-dev -- redis-cli

# 备份 MySQL 数据
kubectl exec deployment/mysql -n vibe-dev -- mysqldump -u root -p vibe_coding_starter > backup.sql
```

### 调试和监控

```bash
# 描述资源详情
kubectl describe pod <pod-name> -n vibe-dev

# 查看事件
kubectl get events -n vibe-dev --sort-by='.lastTimestamp'

# 查看资源使用情况
kubectl top nodes
kubectl top pods -n vibe-dev
```

## 故障排除

### 常见问题

#### 1. 集群创建失败

```bash
# 检查 Docker 状态
docker info

# 清理旧的集群
k3d cluster delete vibe-dev

# 重新创建
k3d cluster create --config k3d-cluster.yaml
```

#### 2. Pod 启动失败

```bash
# 查看 Pod 详情
kubectl describe pod <pod-name> -n vibe-dev

# 查看 Pod 日志
kubectl logs <pod-name> -n vibe-dev

# 检查镜像拉取
kubectl get events -n vibe-dev | grep -i pull
```

#### 3. 网络连接问题

```bash
# 检查服务端点
kubectl get endpoints -n vibe-dev

# 测试集群内连接
kubectl run test-pod --image=busybox --rm -it --restart=Never -- sh
# 在 Pod 内测试
nslookup mysql.vibe-dev.svc.cluster.local
telnet mysql.vibe-dev.svc.cluster.local 3306
```

#### 4. 存储问题

**注意**：如果您之前使用的是 `/tmp/vibe-dev-storage` 配置，请先参考 [存储最佳实践文档](./STORAGE_BEST_PRACTICES.md) 进行数据迁移。

如果 Redis Pod 处于 Pending 状态且 PVC 无法绑定：

```bash
# 查看 PVC 状态
kubectl get pvc -n vibe-dev

# 查看 Pod 详细信息
kubectl describe pod redis-0 -n vibe-dev

# 如果遇到 "WaitForFirstConsumer" 问题，可以临时使用 emptyDir：
# 1. 删除现有的 StatefulSet 和 PVC
kubectl delete statefulset redis -n vibe-dev
kubectl delete pvc redis-data-redis-0 -n vibe-dev

# 2. 临时修改 redis.yaml，将 volumeClaimTemplates 替换为：
#    volumes:
#    - name: redis-data
#      emptyDir: {}

# 3. 重新部署
kubectl apply -f manifests/redis.yaml

# 查看存储类
kubectl get storageclass

# 检查节点存储
kubectl describe node
```

### 重置环境

```bash
# 删除所有资源
kubectl delete namespace vibe-dev

# 删除集群
k3d cluster delete vibe-dev

# 清理 Docker 资源（可选）
docker system prune -a

# 重新创建环境
k3d cluster create --config k3d-cluster.yaml
kubectl apply -f manifests/namespace.yaml
kubectl apply -f manifests/mysql.yaml
kubectl apply -f manifests/redis.yaml
```

### 性能优化

```bash
# 调整资源限制
kubectl edit statefulset mysql -n vibe-dev
kubectl edit statefulset redis -n vibe-dev

# 查看资源使用情况
kubectl top pods -n vibe-dev
```

## 配置应用程序

使用 k3d 环境时，请确保您的应用程序使用 `config-k3d.yaml` 配置文件：

```bash
# 在项目根目录运行应用程序
cd ../../..
go run cmd/server/main.go -config configs/config-k3d.yaml
```
