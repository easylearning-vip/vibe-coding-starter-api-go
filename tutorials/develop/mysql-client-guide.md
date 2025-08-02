# MySQL 客户端连接指南

本指南将帮助您配置本地 MySQL 客户端，快速连接到 k3d 环境中的 MySQL 数据库。支持命令行工具、图形化客户端和无密码快速登录。

## 📋 目录

- [连接信息](#连接信息)
- [命令行客户端](#命令行客户端)


## 🔗 连接信息

根据 `configs/config.yaml` 配置，k3d 环境中的 MySQL 连接信息如下：

### 基本连接参数
```yaml
# 数据库连接信息
Host: 127.0.0.1
Port: 3306 (k3d NodePort 映射)
Username: vibe_user
Password: vibe_password
Database: vibe_coding_starter
Charset: utf8mb4
```

### k3d 集群信息
```yaml
# k3d 特定配置
Cluster Name: vibe-dev
Namespace: vibe-dev
MySQL Service: mysql.vibe-dev.svc.cluster.local
Internal Port: 3306
NodePort: 30306
```

## 💻 命令行客户端

### 1. 安装 MySQL 客户端

#### Ubuntu/Debian
```bash
# 安装 MySQL 客户端
sudo apt update
sudo apt install mysql-client

# 或者安装完整的 MySQL（包含客户端）
sudo apt install mysql-server
```

### 2. 基本连接命令

#### 标准连接
```bash
mysql -h 127.0.0.1 -P 3306 -u vibe_user -p vibe_coding_starter
# 输入密码：vibe_password
```

#### 一行命令连接（不推荐生产环境）
```bash
mysql -h 127.0.0.1 -P 3306 -u vibe_user -pvibe_password vibe_coding_starter
```

#### 使用完整参数
```bash
mysql \
  --host=127.0.0.1 \
  --port=3306 \
  --user=vibe_user \
  --password=vibe_password \
  --database=vibe_coding_starter \
  --default-character-set=utf8mb4
```

### 3. 快速无密码登录配置

创建 MySQL 配置文件实现无密码登录：

#### 创建配置文件
```bash
# 创建 MySQL 配置目录（如果不存在）
mkdir -p ~/.mysql

# 创建配置文件
cat > ~/.mysql/vibe-dev.cnf << 'EOF'
[client]
host = 127.0.0.1
port = 3306
user = vibe_user
password = vibe_password
database = vibe_coding_starter
default-character-set = utf8mb4

[mysql]
prompt = "vibe-dev> "
auto-rehash
EOF

# 设置文件权限（重要：保护密码安全）
chmod 600 ~/.mysql/vibe-dev.cnf
```

#### 使用配置文件连接
```bash
# 使用配置文件快速连接
mysql --defaults-file=~/.mysql/vibe-dev.cnf

# 创建别名方便使用
echo 'alias mysql-vibe="mysql --defaults-file=~/.mysql/vibe-dev.cnf"' >> ~/.bashrc
# 或者 zsh 用户
echo 'alias mysql-vibe="mysql --defaults-file=~/.mysql/vibe-dev.cnf"' >> ~/.zshrc

# 重新加载配置
source ~/.bashrc  # 或 source ~/.zshrc

# 现在可以直接使用
mysql-vibe
```
