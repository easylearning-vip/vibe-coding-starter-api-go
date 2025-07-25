# Vibe Coding Starter API - K8s 部署资源

本目录包含将 Vibe Coding Starter API 部署到 k3d 开发环境的所有 Kubernetes 资源清单。

## 文件说明

- `configmap.yaml` - 应用配置文件
- `service.yaml` - 服务定义
- `deployment.yaml` - 部署配置
- `ingress.yaml` - 入口路由配置
- `deploy.sh` - 自动化部署脚本
- `k8s-dev-manual.md` - 详细的手动部署指南

## 快速部署

### 方式一：使用自动化脚本（推荐）

```bash
# 进入部署目录
cd vibe-coding-starter-go-api/deploy/k8s

# 执行部署脚本
./deploy.sh

# 清理部署
./deploy.sh clean
```

### 方式二：手动部署

```bash
# 1. 设置本地镜像仓库
docker run -d --name local-registry --restart=always -p 5555:5000 registry:2
docker network connect k3d-vibe-dev local-registry

# 2. 构建和推送镜像
cd vibe-coding-starter-go-api
docker build -t localhost:5555/vibe-coding-starter-api:latest .
docker push localhost:5555/vibe-coding-starter-api:latest

# 3. 部署 Kubernetes 资源
cd deploy/k8s
kubectl apply -f configmap.yaml
kubectl apply -f service.yaml
kubectl apply -f deployment.yaml
kubectl apply -f ingress.yaml

# 4. 等待部署就绪
kubectl wait --for=condition=ready pod -l app=vibe-api -n vibe-dev --timeout=300s

# 5. 配置本地访问
echo '127.0.0.1 api.vibe-dev.com' | sudo tee -a /etc/hosts
```

## 验证部署

```bash
# 查看部署状态
kubectl get all -n vibe-dev -l app=vibe-api

# 测试 API 访问
curl http://api.vibe-dev.com:8000/health

# 查看应用日志
kubectl logs -f deployment/vibe-api-deployment -n vibe-dev
```

## 访问地址

- **API 服务**: http://api.vibe-dev.com:8000
- **健康检查**: http://api.vibe-dev.com:8000/health
- **API 文档**: http://api.vibe-dev.com:8000/swagger/index.html (如果有)

## 故障排除

### 常见问题

1. **镜像拉取失败**
   ```bash
   # 检查本地仓库
   curl http://localhost:5555/v2/_catalog

   # 重新推送镜像
   docker push localhost:5555/vibe-coding-starter-api:latest
   ```

2. **Pod 启动失败**
   ```bash
   # 查看 Pod 详情
   kubectl describe pod -l app=vibe-api -n vibe-dev
   
   # 查看日志
   kubectl logs -l app=vibe-api -n vibe-dev
   ```

3. **数据库连接失败**
   ```bash
   # 检查 MySQL 服务
   kubectl get svc mysql -n vibe-dev
   
   # 测试连接
   kubectl exec -it deployment/vibe-api-deployment -n vibe-dev -- nc -zv mysql.vibe-dev.svc.cluster.local 3306
   ```

4. **域名访问失败**
   ```bash
   # 检查 hosts 文件
   grep api.vibe-dev.com /etc/hosts
   
   # 检查 Ingress
   kubectl get ingress -n vibe-dev
   ```

### 清理部署

```bash
# 使用脚本清理
./deploy.sh clean

# 或手动清理
kubectl delete -f .
```

## 配置说明

### 应用配置

应用使用 ConfigMap 中的配置文件，主要配置项：

- **数据库**: 连接到 k8s 内的 MySQL 服务
- **缓存**: 连接到 k8s 内的 Redis 服务
- **端口**: 容器内监听 8080 端口
- **健康检查**: `/health` 端点

### 资源限制

- **CPU**: 请求 100m，限制 500m
- **内存**: 请求 128Mi，限制 512Mi
- **副本数**: 2 个实例

### 安全配置

- 使用非 root 用户运行
- 只读根文件系统
- 删除所有 Linux capabilities

## 更多信息

详细的部署说明请参考 `k8s-dev-manual.md` 文件。
