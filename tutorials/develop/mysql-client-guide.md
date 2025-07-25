# MySQL 客户端连接指南

本指南将帮助您配置本地 MySQL 客户端，快速连接到 Docker Compose 和 k3d 环境中的 MySQL 数据库。

## 目录

- [安装 MySQL 客户端](#安装-mysql-客户端)
- [配置 MySQL 客户端](#配置-mysql-客户端)
- [Docker Compose 环境连接](#docker-compose-环境连接)
- [k3d 环境连接](#k3d-环境连接)
- [连接测试](#连接测试)
- [常用操作](#常用操作)
- [故障排除](#故障排除)

## 安装 MySQL 客户端

### Ubuntu/Debian

```bash
# 安装 MySQL 客户端
sudo apt update
sudo apt install mysql-client

# 验证安装
mysql --version
```

### CentOS/RHEL

```bash
# 安装 MySQL 客户端
sudo yum install mysql

# 或者使用 dnf (较新版本)
sudo dnf install mysql

# 验证安装
mysql --version
```

### macOS

```bash
# 使用 Homebrew 安装
brew install mysql-client

# 添加到 PATH (如果需要)
echo 'export PATH="/opt/homebrew/opt/mysql-client/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# 验证安装
mysql --version
```

### Windows

1. 下载 MySQL Installer: https://dev.mysql.com/downloads/installer/
2. 选择 "Custom" 安装类型
3. 只选择 "MySQL Command Line Client"
4. 完成安装

## 配置 MySQL 客户端

### 创建用户配置文件

创建 `~/.my.cnf` 文件来存储不同环境的连接配置：

```bash
# 创建配置文件
touch ~/.my.cnf
chmod 600 ~/.my.cnf  # 设置安全权限
```

### 配置文件内容

编辑 `~/.my.cnf` 文件，添加以下内容：

```ini
# MySQL 客户端配置文件
# 文件位置: ~/.my.cnf

[client]
# 默认配置
default-character-set = utf8mb4
protocol = tcp

# Docker Compose 环境配置
[docker]
host = 127.0.0.1
port = 3306
user = vibe_user
password = vibe_password
database = vibe_coding_starter
default-character-set = utf8mb4
protocol = tcp

# Docker Compose Root 用户配置
[docker-root]
host = 127.0.0.1
port = 3306
user = root
password = rootpassword
database = vibe_coding_starter
default-character-set = utf8mb4
protocol = tcp

# k3d 环境配置 - 使用固定 NodePort (通过 k3d 端口映射 3306->30306)
[k3d]
host = 127.0.0.1
port = 3306
user = vibe_user
password = vibe_password
database = vibe_coding_starter
default-character-set = utf8mb4
protocol = tcp

# k3d Root 用户配置
[k3d-root]
host = 127.0.0.1
port = 3306
user = root
password = rootpassword
database = vibe_coding_starter
default-character-set = utf8mb4
protocol = tcp

```

## Docker Compose 环境连接

### 启动 Docker Compose 环境

```bash
# 进入 Docker Compose 目录
cd vibe-coding-starter-go-api/dev-tutorial/docker-compose

# 启动服务
docker-compose -f docker-compose.dev.yml up -d

# 等待服务启动
sleep 30

# 检查服务状态
docker-compose -f docker-compose.dev.yml ps
```

### 连接方式

#### 方式一：使用配置文件连接

```bash
# 使用 docker 配置连接
mysql --defaults-group-suffix=docker

# 或者使用 root 用户连接
mysql --defaults-group-suffix=docker-root
```

#### 方式二：命令行参数连接

```bash
# 普通用户连接
mysql -h 127.0.0.1 -P 3306 -u vibe_user -pvibe_password vibe_coding_starter

# Root 用户连接
mysql -h 127.0.0.1 -P 3306 -u root -prootpassword vibe_coding_starter
```

#### 方式三：创建连接别名

在 `~/.bashrc` 或 `~/.zshrc` 中添加：

```bash
# Docker Compose MySQL 连接别名
alias mysql-docker='mysql --defaults-group-suffix=docker'
alias mysql-docker-root='mysql --defaults-group-suffix=docker-root'

# 重新加载配置
source ~/.bashrc  # 或 source ~/.zshrc
```

使用别名连接：

```bash
# 使用别名连接
mysql-docker
mysql-docker-root
```

## k3d 环境连接

### 启动 k3d 环境

```bash
# 进入 k3d 目录
cd vibe-coding-starter-go-api/dev-tutorial/k3d

# 启动环境
./scripts/setup-k3d-dev.sh start

# 或者手动启动
k3d cluster create vibe-dev --agents 1 --port "3306:30306@server:0" --port "6379:30379@server:0"
kubectl apply -f manifests/namespace.yaml
kubectl apply -f manifests/mysql.yaml
kubectl apply -f manifests/redis.yaml
```

### 连接方式

#### 方式一：使用配置文件连接

```bash
# 使用 k3d 配置连接
mysql --defaults-group-suffix=k3d

# 或者使用 root 用户连接
mysql --defaults-group-suffix=k3d-root
```

#### 方式二：创建连接别名

在 `~/.bashrc` 或 `~/.zshrc` 中添加：

```bash
# k3d MySQL 连接别名
alias mysql-k3d='mysql --defaults-group-suffix=k3d'
alias mysql-k3d-root='mysql --defaults-group-suffix=k3d-root'
```

使用别名连接：

```bash
# 使用别名连接
mysql-k3d
mysql-k3d-root
```

## 连接测试

### 基本连接测试

#### Docker Compose 环境测试

```bash
# 测试 Docker Compose 连接
echo "测试 Docker Compose 连接..."
mysql --defaults-group-suffix=docker -e "SELECT 'Docker Compose MySQL 连接成功!' as result;"
```

#### k3d 环境测试

```bash
# 测试 k3d 连接 (推荐使用直接连接)
echo "测试 k3d 连接..."

# 方式1: 直接命令行连接 (推荐)
mysql -h 127.0.0.1 -P 3306 -u vibe_user -pvibe_password vibe_coding_starter -e "SELECT 'k3d NodePort 连接成功!' as result;"

# 方式2: 使用别名连接
alias mysql-k3d='mysql -h 127.0.0.1 -P 3306 -u vibe_user -pvibe_password vibe_coding_starter'
mysql-k3d -e "SELECT 'k3d 别名连接成功!' as result;"

# 方式3: 使用配置文件连接 (如果配置正确)
mysql --defaults-group-suffix=k3d -e "SELECT 'k3d 配置文件连接成功!' as result;"
```

**k3d 端口映射说明**：
- k3d 集群通过端口映射 `3306:30306` 提供外部访问
- 本地 3306 端口直接映射到 k3d 集群的 MySQL NodePort 30306
- 因此可以直接使用 `localhost:3306` 连接到 k3d 中的 MySQL

### 数据库结构测试

```bash
# 查看数据库
mysql --defaults-group-suffix=docker -e "SHOW DATABASES;"

# 查看表结构
mysql --defaults-group-suffix=docker -e "SHOW TABLES;"

# 查看用户数据
mysql --defaults-group-suffix=docker -e "SELECT username, email, is_admin FROM users LIMIT 5;"
```

## 常用操作

### 数据库管理

```bash
# 创建数据库备份
mysqldump --defaults-group-suffix=docker vibe_coding_starter > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据库
mysql --defaults-group-suffix=docker vibe_coding_starter < backup_20231201_120000.sql

# 查看数据库大小
mysql --defaults-group-suffix=docker -e "
SELECT 
    table_schema AS 'Database',
    ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
FROM information_schema.tables 
WHERE table_schema = 'vibe_coding_starter'
GROUP BY table_schema;"
```

### 用户管理

```bash
# 查看当前用户权限
mysql --defaults-group-suffix=docker-root -e "SHOW GRANTS FOR 'vibe_user'@'%';"

# 创建新用户 (使用 root 连接)
mysql --defaults-group-suffix=docker-root -e "
CREATE USER 'dev_user'@'%' IDENTIFIED BY 'dev_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON vibe_coding_starter.* TO 'dev_user'@'%';
FLUSH PRIVILEGES;"
```

### 监控和调试

```bash
# 查看连接状态
mysql --defaults-group-suffix=docker-root -e "SHOW PROCESSLIST;"

# 查看数据库状态
mysql --defaults-group-suffix=docker-root -e "SHOW STATUS LIKE 'Connections';"

# 查看慢查询
mysql --defaults-group-suffix=docker-root -e "SHOW VARIABLES LIKE 'slow_query%';"
```

