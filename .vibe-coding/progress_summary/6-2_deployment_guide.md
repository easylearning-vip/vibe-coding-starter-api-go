# Vibe Coding Starter API - 部署文档

## 概述
本文档提供了 Vibe Coding Starter API 的完整部署指南，包括 Docker 容器化、Kubernetes 部署和生产环境配置。

## 系统要求

### 硬件要求
- **CPU**: 最少 2 核，推荐 4 核
- **内存**: 最少 4GB，推荐 8GB
- **存储**: 最少 20GB 可用空间
- **网络**: 稳定的互联网连接

### 软件要求
- **Docker**: 20.10 或更高版本
- **Kubernetes**: 1.25 或更高版本
- **kubectl**: 与 Kubernetes 版本匹配
- **Helm**: 3.0 或更高版本（可选）

## 架构概述

### 应用架构
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend UI   │────│   API Gateway   │────│   Backend API   │
│  (React + AntD) │    │   (Nginx/Ingress)│    │   (Go + Gin)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                               ┌─────────────────┐
                                               │   MySQL DB      │
                                               │   (RDS/MySQL)   │
                                               └─────────────────┘
                                                        │
                                               ┌─────────────────┐
                                               │   Redis Cache   │
                                               │   (ElastiCache) │
                                               └─────────────────┘
```

### 组件说明
- **Backend API**: Go 语言编写的 RESTful API
- **Frontend UI**: React + Ant Design Pro 管理界面
- **Database**: MySQL 数据库存储业务数据
- **Cache**: Redis 缓存提高性能
- **Load Balancer**: Nginx 或 Kubernetes Ingress

## 部署方式

### 1. Docker 容器化部署

#### 1.1 构建 Docker 镜像
```bash
# 进入项目目录
cd /workspace/vibe-coding-starter-api-go

# 构建 Docker 镜像
docker build -t vibe-coding-starter-api:latest .

# 或者使用 Makefile
make docker-build
```

#### 1.2 运行 Docker 容器
```bash
# 运行容器
docker run -d \
  --name vibe-api \
  -p 8080:8080 \
  -e DB_HOST=mysql \
  -e DB_PORT=3306 \
  -e DB_USER=root \
  -e DB_PASSWORD=secret \
  -e DB_NAME=vibe_db \
  -e REDIS_HOST=redis \
  -e REDIS_PORT=6379 \
  -e JWT_SECRET=your-secret-key \
  vibe-coding-starter-api:latest
```

#### 1.3 Docker Compose 部署
```bash
# 使用 Docker Compose
cd tutorials/develop/docker-compose
docker compose -f docker-compose.dev.yml up -d
```

### 2. Kubernetes 部署

#### 2.1 准备 Kubernetes 集群
```bash
# 使用 k3d 创建开发集群
cd tutorials/develop/k3d
k3d cluster create --config k3d-cluster.yaml

# 或者使用现有集群
kubectl config use-context your-cluster
```

#### 2.2 部署到 Kubernetes
```bash
# 构建并推送镜像
make docker-build
make docker-push

# 部署到 Kubernetes
make k8s-deploy
```

#### 2.3 验证部署
```bash
# 检查部署状态
kubectl get all -n vibe-dev

# 查看应用日志
kubectl logs -f deployment/vibe-api-deployment -n vibe-dev

# 检查服务状态
kubectl get svc -n vibe-dev
```

## 环境配置

### 开发环境配置
```yaml
# configs/config.yaml
server:
  port: 8080
  debug: true

database:
  host: localhost
  port: 3306
  user: root
  password: secret
  name: vibe_db
  type: mysql

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: dev-secret-key
  expire_hours: 24
```

### 生产环境配置
```yaml
# configs/config-prod.yaml
server:
  port: 8080
  debug: false

database:
  host: prod-mysql.example.com
  port: 3306
  user: vibe_user
  password: ${DB_PASSWORD}
  name: vibe_prod
  type: mysql
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: prod-redis.example.com
  port: 6379
  password: ${REDIS_PASSWORD}
  db: 0

jwt:
  secret: ${JWT_SECRET}
  expire_hours: 168

logging:
  level: info
  format: json
```

## 安全配置

### 1. 环境变量
```bash
# 敏感信息通过环境变量传递
export DB_PASSWORD="your-secure-password"
export REDIS_PASSWORD="your-redis-password"
export JWT_SECRET="your-jwt-secret-key"
```

### 2. Kubernetes Secrets
```yaml
# secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: vibe-api-secrets
  namespace: vibe-dev
type: Opaque
data:
  db-password: <base64-encoded-password>
  redis-password: <base64-encoded-password>
  jwt-secret: <base64-encoded-secret>
```

### 3. 网络安全
- 配置防火墙规则
- 使用 HTTPS/TLS 加密
- 配置 CORS 策略
- 启用 API 限流

## 监控和日志

### 1. 应用监控
- **健康检查**: `/api/v1/health` 端点
- **指标收集**: 集成 Prometheus
- **日志聚合**: 使用 ELK 或类似方案

### 2. 数据库监控
- **MySQL**: 性能监控、慢查询日志
- **Redis**: 内存使用、连接数监控

### 3. 系统监控
- **CPU/内存使用率**
- **磁盘空间**
- **网络流量**

## 备份和恢复

### 1. 数据库备份
```bash
# MySQL 备份
mysqldump -h mysql-host -u vibe_user -p vibe_prod > backup.sql

# 定时备份脚本
0 2 * * * /usr/bin/mysqldump -h mysql-host -u vibe_user -p vibe_prod | gzip > /backup/vibe_$(date +\%Y\%m\%d).sql.gz
```

### 2. 恢复数据
```bash
# 恢复数据库
mysql -h mysql-host -u vibe_user -p vibe_prod < backup.sql
```

## 性能优化

### 1. 数据库优化
- 创建适当的索引
- 配置连接池
- 启用查询缓存

### 2. 应用优化
- 启用 Gzip 压缩
- 配置 HTTP 缓存
- 使用 CDN 加速静态资源

### 3. 基础设施优化
- 使用负载均衡
- 配置自动扩展
- 优化存储性能

## 故障排除

### 1. 常见问题
- **数据库连接失败**: 检查数据库配置和网络连接
- **Redis 连接失败**: 检查 Redis 服务状态
- **JWT 验证失败**: 检查 JWT 密钥配置
- **内存不足**: 增加 Kubernetes Pod 资源限制

### 2. 调试命令
```bash
# 查看应用日志
kubectl logs deployment/vibe-api-deployment -n vibe-dev

# 进入容器调试
kubectl exec -it deployment/vibe-api-deployment -n vibe-dev -- /bin/bash

# 检查资源使用情况
kubectl top pod -n vibe-dev
```

## 部署检查清单

### 部署前检查
- [ ] 代码编译通过
- [ ] 单元测试通过
- [ ] 安全检查通过
- [ ] 配置文件准备就绪
- [ ] 数据库迁移脚本准备

### 部署后检查
- [ ] 应用启动成功
- [ ] 健康检查正常
- [ ] 数据库连接正常
- [ ] Redis 连接正常
- [ ] API 端点响应正常
- [ ] 日志输出正常

### 生产环境检查
- [ ] HTTPS 证书配置
- [ ] 监控系统集成
- [ ] 备份策略配置
- [ ] 告警规则配置
- [ ] 性能基线建立

## 联系信息

如有部署问题或需要技术支持，请联系：
- **技术支持**: support@example.com
- **开发团队**: dev@example.com
- **运维团队**: ops@example.com

---

*文档版本: 1.0*  
*最后更新: 2025-08-15*  
*维护者: Vibe Coding Starter Team*