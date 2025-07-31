# Vibe Coding Starter 开发环境安装脚本

本目录包含了用于快速搭建 Vibe Coding Starter 开发环境的自动化脚本。

## 脚本概览

### 🚀 一键安装脚本

#### `setup-dev-environment.sh` - 完整环境安装
适用于**全新的 Ubuntu 系统**，从零开始安装所有依赖。

**功能：**
- 自动检测系统要求
- 安装 Docker Engine
- 安装 k3d 和 kubectl
- 创建 k3d 开发集群
- 部署 MySQL 和 Redis 服务
- 配置 shell 补全和便捷别名
- 验证所有服务状态

**使用方法：**
```bash
# 下载并运行（推荐）
curl -fsSL https://raw.githubusercontent.com/easylearning-vip/vibe-coding-starter-api-go/main/tutorials/develop/scripts/setup-dev-environment.sh | bash

# 或者本地运行
cd vibe-coding-starter-api-go/tutorials/develop/scripts
./setup-dev-environment.sh
```

#### `quick-start.sh` - 快速启动
适用于**已安装基础组件**的系统，快速创建和启动开发环境。

**前提条件：**
- 已安装 Docker
- 已安装 k3d
- 已安装 kubectl

**使用方法：**
```bash
cd vibe-coding-starter-api-go/tutorials/develop/scripts
./quick-start.sh
```

### 🔧 单独组件安装脚本

#### `install-docker.sh` - Docker 安装
安装 Docker Engine 和 Docker Compose。

**支持系统：**
- Ubuntu/Debian
- CentOS/RHEL/Fedora
- macOS (通过 Homebrew)

**使用方法：**
```bash
./install-docker.sh
```

#### `install-k3d.sh` - k3d 工具链安装
安装 k3d、kubectl 和 Helm。

**支持系统：**
- Linux (x86_64, arm64)
- macOS (x86_64, arm64)

**使用方法：**
```bash
./install-k3d.sh
```

## 系统要求

### 最低要求
- **操作系统**: Ubuntu 18.04+ / Debian 10+
- **内存**: 4GB RAM
- **磁盘空间**: 10GB 可用空间
- **网络**: 稳定的互联网连接

### 推荐配置
- **操作系统**: Ubuntu 20.04+ / Debian 11+
- **内存**: 8GB+ RAM
- **磁盘空间**: 20GB+ 可用空间
- **CPU**: 2+ 核心

## 快速开始

### 方案一：全新系统（推荐）

如果您使用的是全新的 Ubuntu 系统：

```bash
# 1. 克隆项目（如果还没有）
git clone <repository-url>
cd vibe-coding-starter-api-go/tutorials/develop/scripts

# 2. 运行一键安装脚本
./setup-dev-environment.sh

# 3. 重新加载 shell 配置
source ~/.bashrc

# 4. 验证安装
vibe-status
```

### 方案二：已有 Docker 环境

如果您已经安装了 Docker：

```bash
# 1. 安装 k3d 工具链
./install-k3d.sh

# 2. 快速启动开发环境
./quick-start.sh

# 3. 验证服务
kubectl get all -n vibe-dev
```

### 方案三：分步安装

如果您希望分步安装：

```bash
# 1. 安装 Docker
./install-docker.sh

# 2. 重新登录或运行
newgrp docker

# 3. 安装 k3d 工具链
./install-k3d.sh

# 4. 手动创建集群
cd ../k3d
k3d cluster create --config k3d-cluster.yaml
kubectl apply -f manifests/
```

## 安装后验证

### 检查服务状态
```bash
# 查看所有服务
kubectl get all -n vibe-dev

# 查看 Pod 状态
kubectl get pods -n vibe-dev

# 查看服务端点
kubectl get svc -n vibe-dev
```

### 测试数据库连接
```bash
# MySQL 连接测试
mysql -h localhost -P 3306 -u vibe_user -pvibe_password vibe_coding_starter -e "SELECT 'MySQL 连接成功!' as result;"

# Redis 连接测试
redis-cli -h localhost -p 6379 ping
```

### 查看服务日志
```bash
# MySQL 日志
kubectl logs -f deployment/mysql -n vibe-dev

# Redis 日志
kubectl logs -f deployment/redis -n vibe-dev
```

## 便捷别名

安装脚本会自动创建以下便捷别名：

```bash
k          # kubectl 的简写
kgp        # kubectl get pods
kgs        # kubectl get svc
kgn        # kubectl get nodes
vibe-dev   # 切换到 vibe-dev 命名空间
vibe-status # 查看所有服务状态
vibe-mysql  # 快速连接 MySQL
vibe-redis  # 快速连接 Redis
vibe-logs-mysql # 查看 MySQL 日志
vibe-logs-redis # 查看 Redis 日志
```

## 数据库连接信息

### MySQL
- **Host**: localhost
- **Port**: 3306
- **Database**: vibe_coding_starter
- **Username**: vibe_user
- **Password**: vibe_password
- **Root Password**: rootpassword

### Redis
- **Host**: localhost
- **Port**: 6379
- **Password**: (无)

## 常见问题

### 1. Docker 权限问题
```bash
# 将用户添加到 docker 组
sudo usermod -aG docker $USER

# 重新登录或运行
newgrp docker
```

### 2. 端口冲突
如果端口 3306 或 6379 被占用：
```bash
# 查看端口占用
sudo netstat -tlnp | grep :3306
sudo netstat -tlnp | grep :6379

# 停止冲突服务或修改 k3d-cluster.yaml 中的端口映射
```

### 3. 集群创建失败
```bash
# 删除现有集群
k3d cluster delete vibe-dev

# 清理 Docker 资源
docker system prune -f

# 重新创建
./quick-start.sh
```

### 4. 服务启动失败
```bash
# 查看 Pod 详情
kubectl describe pod <pod-name> -n vibe-dev

# 查看事件
kubectl get events -n vibe-dev --sort-by='.lastTimestamp'

# 重新部署
kubectl delete -f manifests/
kubectl apply -f manifests/
```

## 管理命令

### 集群管理
```bash
# 启动集群
k3d cluster start vibe-dev

# 停止集群
k3d cluster stop vibe-dev

# 删除集群
k3d cluster delete vibe-dev

# 查看集群列表
k3d cluster list
```

### 服务管理
```bash
# 重启服务
kubectl rollout restart deployment/mysql -n vibe-dev
kubectl rollout restart deployment/redis -n vibe-dev

# 扩缩容（不推荐用于数据库）
kubectl scale deployment/redis --replicas=2 -n vibe-dev

# 更新配置
kubectl apply -f manifests/
```

## 卸载

### 完全卸载
```bash
# 1. 删除 k3d 集群
k3d cluster delete vibe-dev

# 2. 卸载 k3d
sudo rm -f /usr/local/bin/k3d

# 3. 卸载 kubectl
sudo rm -f /usr/local/bin/kubectl

# 4. 卸载 Docker（可选）
sudo apt remove docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo rm -rf /var/lib/docker
sudo rm -rf /etc/docker

# 5. 清理配置文件
rm -rf ~/.local/share/k3d
rm -rf ~/.kube
```

### 仅删除开发环境
```bash
# 删除集群但保留工具
k3d cluster delete vibe-dev

# 清理存储
rm -rf ~/.local/share/k3d/vibe-dev-storage
```

## 技术支持

如果遇到问题，请：

1. 查看本文档的常见问题部分
2. 检查 [k3d 官方文档](https://k3d.io/)
3. 查看项目的 Issues 页面
4. 联系技术支持团队

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这些脚本！
