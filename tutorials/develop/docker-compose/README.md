# Docker Compose 开发环境搭建指南

本指南将帮助您使用 Docker Compose 快速搭建 Vibe Coding Starter 的开发环境，包括 MySQL 和 Redis 服务。

## 目录

- [前置要求](#前置要求)
- [Docker 和 Docker Compose 安装](#docker-和-docker-compose-安装)
- [环境配置](#环境配置)
- [启动开发环境](#启动开发环境)
- [连接数据库](#连接数据库)
- [管理工具](#管理工具)
- [常用命令](#常用命令)
- [故障排除](#故障排除)

## 前置要求

- 操作系统：Linux、macOS 或 Windows
- 至少 4GB 可用内存
- 至少 2GB 可用磁盘空间

## Docker 和 Docker Compose 安装

### Ubuntu/Debian

```bash
# 更新包索引
sudo apt update

# 安装必要的包
sudo apt install -y apt-transport-https ca-certificates curl gnupg lsb-release

# 添加 Docker 官方 GPG 密钥
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# 添加 Docker 仓库
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 安装 Docker Engine
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 将当前用户添加到 docker 组
sudo usermod -aG docker $USER

# 重新登录或运行以下命令使组更改生效
newgrp docker
```

### CentOS/RHEL

```bash
# 安装必要的包
sudo yum install -y yum-utils

# 添加 Docker 仓库
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

# 安装 Docker Engine
sudo yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 启动 Docker 服务
sudo systemctl start docker
sudo systemctl enable docker

# 将当前用户添加到 docker 组
sudo usermod -aG docker $USER
```

### macOS

```bash
# 使用 Homebrew 安装
brew install --cask docker

# 或者下载 Docker Desktop for Mac
# https://docs.docker.com/desktop/mac/install/
```

### Windows

下载并安装 Docker Desktop for Windows：
https://docs.docker.com/desktop/windows/install/

## 环境配置

### 1. 进入项目目录

```bash
cd vibe-coding-starter-go-api/dev-tutorial/docker-compose
```

### 2. 检查配置文件

确保以下文件存在：
- `docker-compose.dev.yml` - Docker Compose 配置文件
- `mysql/my.cnf` - MySQL 配置文件
- `redis/redis.conf` - Redis 配置文件

### 3. 创建必要的目录

```bash
# 创建日志目录
mkdir -p logs

# 创建上传目录
mkdir -p uploads
```

## 启动开发环境

### 1. 启动所有服务

```bash
# 启动 MySQL 和 Redis
docker compose -f docker-compose.dev.yml up -d

# 查看服务状态
docker compose -f docker-compose.dev.yml ps
```

### 2. 启动带管理工具的服务

```bash
# 启动包含 phpMyAdmin 和 Redis Commander 的完整环境
docker compose -f docker-compose.dev.yml --profile tools up -d
```

### 3. 启动 PostgreSQL（可选）

```bash
# 启动 PostgreSQL
docker compose -f docker-compose.dev.yml --profile postgres up -d
```

### 4. 查看日志

```bash
# 查看所有服务日志
docker compose -f docker-compose.dev.yml logs -f

# 查看特定服务日志
docker compose -f docker-compose.dev.yml logs -f mysql
docker compose -f docker-compose.dev.yml logs -f redis
```

## 连接数据库

### MySQL 连接信息

- **主机**: localhost
- **端口**: 3306
- **数据库**: vibe_coding_starter
- **用户名**: vibe_user
- **密码**: vibe_password
- **Root 密码**: rootpassword

### 使用本地 MySQL 客户端连接

#### 快速连接

```bash
# 使用普通用户连接
mysql -h 127.0.0.1 -P 3306 -u vibe_user -pvibe_password vibe_coding_starter

# 使用 root 用户连接
mysql -h 127.0.0.1 -P 3306 -u root -prootpassword vibe_coding_starter
```

#### 配置 ~/.my.cnf 实现快速连接

创建或编辑 `~/.my.cnf` 文件：

```ini
[docker]
host = 127.0.0.1
port = 3306
user = vibe_user
password = vibe_password
database = vibe_coding_starter
default-character-set = utf8mb4
protocol = tcp

[docker-root]
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
mysql --defaults-group-suffix=docker

# 使用 root 配置连接
mysql --defaults-group-suffix=docker-root

# 创建别名（添加到 ~/.bashrc 或 ~/.zshrc）
alias mysql-docker='mysql --defaults-group-suffix=docker'
alias mysql-docker-root='mysql --defaults-group-suffix=docker-root'
```

#### 连接测试

```bash
# 测试连接
mysql --defaults-group-suffix=docker -e "SELECT 'Docker MySQL 连接成功!' as result;"

# 查看数据库
mysql --defaults-group-suffix=docker -e "SHOW DATABASES;"

# 查看表
mysql --defaults-group-suffix=docker -e "SHOW TABLES;"
```

### Redis 连接信息

- **主机**: localhost
- **端口**: 6379
- **密码**: 无

### PostgreSQL 连接信息（如果启用）

- **主机**: localhost
- **端口**: 5432
- **数据库**: vibe_coding_starter
- **用户名**: postgres
- **密码**: postgres_password

## 管理工具

### phpMyAdmin

如果启用了 tools profile，可以通过以下地址访问 phpMyAdmin：

- **URL**: http://localhost:8080
- **用户名**: root
- **密码**: rootpassword

### Redis Commander

如果启用了 tools profile，可以通过以下地址访问 Redis Commander：

- **URL**: http://localhost:8081

## 常用命令

### 服务管理

```bash
# 启动服务
docker compose -f docker-compose.dev.yml up -d

# 停止服务
docker compose -f docker-compose.dev.yml down

# 重启服务
docker compose -f docker-compose.dev.yml restart

# 停止并删除数据卷
docker compose -f docker-compose.dev.yml down -v
```

### 数据库操作

```bash
# 连接到 MySQL
docker compose -f docker-compose.dev.yml exec mysql mysql -u vibe_user -p vibe_coding_starter

# 连接到 Redis
docker compose -f docker-compose.dev.yml exec redis redis-cli

# 备份 MySQL 数据库
docker compose -f docker-compose.dev.yml exec mysql mysqldump -u root -p vibe_coding_starter > backup.sql

# 恢复 MySQL 数据库
docker compose -f docker-compose.dev.yml exec -T mysql mysql -u root -p vibe_coding_starter < backup.sql
```

### 监控和调试

```bash
# 查看资源使用情况
docker compose -f docker-compose.dev.yml top

# 查看服务状态
docker compose -f docker-compose.dev.yml ps

# 进入容器
docker compose -f docker-compose.dev.yml exec mysql bash
docker compose -f docker-compose.dev.yml exec redis sh
```

## 故障排除

### 常见问题

#### 1. 端口冲突

如果遇到端口冲突，可以修改 `docker-compose.dev.yml` 中的端口映射：

```yaml
ports:
  - "3307:3306"  # 将 MySQL 映射到 3307 端口
  - "6380:6379"  # 将 Redis 映射到 6380 端口
```

#### 2. 权限问题

```bash
# 修复 Docker 权限问题
sudo chown -R $USER:$USER /var/run/docker.sock
```

#### 3. 内存不足

```bash
# 清理未使用的 Docker 资源
docker system prune -a

# 查看 Docker 资源使用情况
docker system df
```

#### 4. 数据库连接失败

```bash
# 检查 MySQL 服务状态
docker compose -f docker-compose.dev.yml logs mysql

# 重启 MySQL 服务
docker compose -f docker-compose.dev.yml restart mysql

# 检查网络连接
docker compose -f docker-compose.dev.yml exec mysql ping redis
```

### 重置环境

如果需要完全重置开发环境：

```bash
# 停止所有服务并删除数据卷
docker compose -f docker-compose.dev.yml down -v

# 删除所有相关镜像
docker compose -f docker-compose.dev.yml down --rmi all

# 清理系统
docker system prune -a

# 重新启动
docker compose -f docker-compose.dev.yml up -d
```

## 配置应用程序

使用 Docker Compose 环境时，请确保您的应用程序使用 `config-docker.yaml` 配置文件：

```bash
# 在项目根目录运行应用程序
cd ../../..
go run cmd/server/main.go -config configs/config-docker.yaml
```

## 下一步

环境搭建完成后，您可以：

1. 运行数据库迁移
2. 启动应用程序
3. 运行测试
4. 开始开发

更多信息请参考项目的主要文档。
