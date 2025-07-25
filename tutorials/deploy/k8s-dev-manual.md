# Vibe Coding Starter API - K8s 开发环境手动部署指南

本指南将帮助您将 Vibe Coding Starter API 应用手动部署到 k3d 开发环境中进行验证。

## 前置条件

- k3d 集群已创建并运行（使用 `dev-tutorial/k3d` 下的配置）
- MySQL 和 Redis 服务已部署在 `vibe-dev` 命名空间
- kubectl 已配置并可访问 k3d 集群
- Docker 已安装并运行

## 重要说明

本部署使用 ConfigMap 来配置应用程序连接到集群内的 MySQL 和 Redis 服务：
- **MySQL**: `mysql.vibe-dev.svc.cluster.local:3306`
- **Redis**: `redis.vibe-dev.svc.cluster.local:6379`

应用程序通过 `CONFIG_FILE` 环境变量读取挂载的配置文件 `/app/config/config.yaml`。

## 部署概览

本部署包含以下 Kubernetes 资源：
- **ConfigMap**: 应用配置
- **Service**: 服务暴露
- **Deployment**: 应用部署
- **Ingress**: 外部访问路由

最终通过域名 `api.vibe-dev.com:8000` 访问 API 服务。

## 重要更新

本教程已更新以支持 **Traefik Ingress Controller**（k3d 默认）：
- ✅ 自动启用 Traefik（移除了 `--disable=traefik` 参数）
- ✅ 使用 Traefik 注解替代 nginx 注解
- ✅ 通过 k3d 端口映射 `localhost:8000` 访问服务
- ✅ 修复了部署脚本的工作目录问题

## 快速部署（推荐）

如果您想要快速部署而不需要了解详细步骤，可以使用自动化部署脚本：

```bash
# 进入部署目录
cd vibe-coding-starter-go-api/deploy/k8s

# 执行自动化部署
./deploy.sh

# 查看部署状态
kubectl get all -n vibe-dev -l app=vibe-api
```

以下是手动部署的详细步骤，用于学习和故障排除。

## 步骤 1: 使用 k3d 内置镜像仓库

k3d 集群已经配置了内置的镜像仓库 `vibe-registry`，无需额外设置。

### 1.1 验证 k3d 镜像仓库

```bash
# 查看 k3d 镜像仓库状态
k3d registry list

# 验证仓库可访问性
curl http://localhost:5555/v2/_catalog
```

**注意**: k3d 内置仓库在主机上通过 `localhost:5555` 访问，在集群内通过 `vibe-registry:5555` 访问。

## 步骤 2: 构建和推送应用镜像

### 2.1 构建应用镜像

```bash
# 进入项目根目录
cd vibe-coding-starter-go-api

# 构建 Docker 镜像
docker build -t localhost:5555/vibe-coding-starter-api:latest .

# 推送镜像到 k3d 仓库
docker push localhost:5555/vibe-coding-starter-api:latest

# 验证镜像已推送
curl http://localhost:5555/v2/vibe-coding-starter-api/tags/list
```

**重要**: 应用程序已更新以支持通过 `CONFIG_FILE` 环境变量读取配置文件，确保能正确连接到集群内的 MySQL 和 Redis 服务。

## 步骤 3: 创建 Kubernetes 资源清单

### 3.1 创建 ConfigMap

创建文件 `deploy/k8s/configmap.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: vibe-api-config
  namespace: vibe-dev
  labels:
    app: vibe-api
    environment: development
data:
  config.yaml: |
    # Vibe Coding Starter k8s 部署配置
    server:
      host: "0.0.0.0"
      port: 8080
      mode: "debug"
      read_timeout: 30
      write_timeout: 30
      idle_timeout: 60

    # 数据库配置 - 连接到 k8s 内的 MySQL
    database:
      driver: "mysql"
      host: "mysql.vibe-dev.svc.cluster.local"
      port: 3306
      username: "vibe_user"
      password: "vibe_password"
      database: "vibe_coding_starter"
      charset: "utf8mb4"
      max_idle_conns: 10
      max_open_conns: 100
      conn_max_lifetime: 3600

    # 缓存配置 - 连接到 k8s 内的 Redis
    cache:
      driver: "redis"
      host: "redis.vibe-dev.svc.cluster.local"
      port: 6379
      password: ""
      database: 0
      pool_size: 10

    # 日志配置
    logger:
      level: "debug"
      format: "console"
      output: "stdout"

    # JWT 配置
    jwt:
      secret: "vibe-k8s-dev-secret-key"
      issuer: "vibe-coding-starter-k8s"
      expiration: 86400

    # 监控配置
    monitoring:
      enabled: true
      metrics_path: "/metrics"
      health_path: "/health"

    # 开发配置
    development:
      auto_migrate: true
      seed_data: true
      debug_sql: true
```

### 3.2 创建 Service

创建文件 `deploy/k8s/service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: vibe-api-service
  namespace: vibe-dev
  labels:
    app: vibe-api
    environment: development
spec:
  selector:
    app: vibe-api
  ports:
    - name: http
      port: 8080
      targetPort: 8080
      protocol: TCP
  type: ClusterIP
```

### 3.3 创建 Deployment

创建文件 `deploy/k8s/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vibe-api-deployment
  namespace: vibe-dev
  labels:
    app: vibe-api
    environment: development
spec:
  replicas: 2
  selector:
    matchLabels:
      app: vibe-api
  template:
    metadata:
      labels:
        app: vibe-api
        environment: development
    spec:
      containers:
      - name: vibe-api
        image: localhost:5555/vibe-coding-starter-api:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: CONFIG_FILE
          value: "/app/config/config.yaml"
        volumeMounts:
        - name: config-volume
          mountPath: /app/config
          readOnly: true
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
      volumes:
      - name: config-volume
        configMap:
          name: vibe-api-config
      restartPolicy: Always
```

### 3.4 创建 Ingress

创建文件 `deploy/k8s/ingress.yaml`:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: vibe-api-ingress
  namespace: vibe-dev
  labels:
    app: vibe-api
    environment: development
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/ssl-redirect: "false"
spec:
  rules:
  - host: api.vibe-dev.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: vibe-api-service
            port:
              number: 8080
```

## 步骤 4: 部署应用到 k8s

### 4.1 应用所有资源清单

```bash
# 进入部署目录
cd deploy/k8s

# 部署 ConfigMap
kubectl apply -f configmap.yaml

# 部署 Service
kubectl apply -f service.yaml

# 部署 Deployment
kubectl apply -f deployment.yaml

# 部署 Ingress
kubectl apply -f ingress.yaml
```

### 4.2 验证部署状态

```bash
# 查看所有资源
kubectl get all -n vibe-dev -l app=vibe-api

# 查看 Pod 状态
kubectl get pods -n vibe-dev -l app=vibe-api

# 查看 Pod 日志
kubectl logs -f deployment/vibe-api-deployment -n vibe-dev

# 等待 Pod 就绪
kubectl wait --for=condition=ready pod -l app=vibe-api -n vibe-dev --timeout=300s

# 测试应用健康状态
kubectl run test-api --image=curlimages/curl --rm -it --restart=Never -n vibe-dev -- \
  curl -s http://vibe-api-service:8080/health
```

**预期输出**: 应该看到包含 `"status":"healthy"`, `"database":{"status":"healthy"}`, `"cache":{"status":"healthy"}` 的 JSON 响应，表示应用成功连接到 MySQL 和 Redis。

## 步骤 5: 配置本地访问

### 5.1 配置 hosts 文件

```bash
# 编辑 hosts 文件（Linux/macOS）
sudo vim /etc/hosts

# 或者使用 echo 命令添加
echo "127.0.0.1 api.vibe-dev.com" | sudo tee -a /etc/hosts

# Windows 用户编辑 C:\Windows\System32\drivers\etc\hosts
# 添加以下行：
# 127.0.0.1 api.vibe-dev.com
```

### 5.2 验证 Traefik Ingress 访问

```bash
# 测试健康检查端点（通过 Traefik）
curl http://api.vibe-dev.com:8000/health

# 预期输出：
# {"status":"healthy","timestamp":"...","version":"1.0.0","services":{"cache":{"status":"healthy"},"database":{"status":"healthy"}}}

# 测试 API 端点
curl http://api.vibe-dev.com:8000/api/v1/users

# 预期输出（需要认证）：
# {"error":"unauthorized","message":"Authorization token required"}

# 测试不存在的端点
curl http://api.vibe-dev.com:8000/nonexistent

# 预期输出：
# 404 page not found
```

### 5.3 验证 Ingress 配置

```bash
# 查看 Ingress 状态
kubectl get ingress -n vibe-dev

# 查看 Ingress 详细信息
kubectl describe ingress vibe-api-ingress -n vibe-dev

# 查看 Traefik 服务
kubectl get svc -n kube-system | grep traefik
curl http://api.vibe-dev.com:8000/swagger/index.html
```

## 故障排除

### 常见问题

1. **镜像拉取失败**
   ```bash
   # 检查镜像是否存在
   curl http://localhost:5555/v2/vibe-coding-starter-api/tags/list

   # 重新构建和推送镜像
   docker build -t localhost:5555/vibe-coding-starter-api:latest .
   docker push localhost:5555/vibe-coding-starter-api:latest
   ```

2. **Pod 启动失败**
   ```bash
   # 查看 Pod 详细信息
   kubectl describe pod -l app=vibe-api -n vibe-dev
   
   # 查看 Pod 日志
   kubectl logs -l app=vibe-api -n vibe-dev
   ```

3. **数据库连接失败**
   ```bash
   # 检查 MySQL 服务状态
   kubectl get svc mysql -n vibe-dev
   
   # 测试数据库连接
   kubectl exec -it deployment/vibe-api-deployment -n vibe-dev -- sh
   # 在容器内测试
   nc -zv mysql.vibe-dev.svc.cluster.local 3306
   ```

4. **Ingress 访问失败**
   ```bash
   # 检查 Ingress 状态
   kubectl get ingress -n vibe-dev
   
   # 检查 Ingress 控制器
   kubectl get pods -n kube-system | grep traefik
   ```

### 清理部署

```bash
# 删除所有 API 相关资源
kubectl delete -f deploy/k8s/

# 或者单独删除
kubectl delete ingress vibe-api-ingress -n vibe-dev
kubectl delete deployment vibe-api-deployment -n vibe-dev
kubectl delete service vibe-api-service -n vibe-dev
kubectl delete configmap vibe-api-config -n vibe-dev
```

## 步骤 6: 测试和验证

### 6.1 运行自动化测试

```bash
# 运行 API 测试脚本
cd deploy/k8s
./test-api.sh
```

### 6.2 手动测试 API

```bash
# 基础健康检查
curl http://api.vibe-dev.com:8000/health

# API 健康检查
curl http://api.vibe-dev.com:8000/api/v1/health

# 查看 Prometheus 指标
curl http://api.vibe-dev.com:8000/metrics

# 测试用户注册（示例）
curl -X POST http://api.vibe-dev.com:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123"}'
```

## 使用 Makefile 简化操作

项目根目录的 Makefile 提供了便捷的命令：

```bash
# 查看所有可用命令
make help

# 完整的开发部署流程
make dev-full

# 快速重新部署
make quick-deploy

# 查看 k8s 状态
make k8s-status

# 查看应用日志
make k8s-logs

# 清理部署
make k8s-clean
```

## 监控和调试

### 查看应用状态

```bash
# 查看所有资源
kubectl get all -n vibe-dev -l app=vibe-api

# 查看 Pod 详情
kubectl describe pod -l app=vibe-api -n vibe-dev

# 实时查看日志
kubectl logs -f deployment/vibe-api-deployment -n vibe-dev

# 查看事件
kubectl get events -n vibe-dev --sort-by='.lastTimestamp'
```

### 进入容器调试

```bash
# 进入应用容器
kubectl exec -it deployment/vibe-api-deployment -n vibe-dev -- sh

# 在容器内测试数据库连接
kubectl exec -it deployment/vibe-api-deployment -n vibe-dev -- nc -zv mysql.vibe-dev.svc.cluster.local 3306

# 在容器内测试 Redis 连接
kubectl exec -it deployment/vibe-api-deployment -n vibe-dev -- nc -zv redis.vibe-dev.svc.cluster.local 6379
```

## 扩展和优化

### 扩展副本数

```bash
# 扩展到 3 个副本
kubectl scale deployment vibe-api-deployment --replicas=3 -n vibe-dev

# 查看扩展状态
kubectl get deployment vibe-api-deployment -n vibe-dev
```

### 更新应用

```bash
# 重新构建和推送镜像
make docker-push

# 重启部署以拉取新镜像
kubectl rollout restart deployment/vibe-api-deployment -n vibe-dev

# 查看滚动更新状态
kubectl rollout status deployment/vibe-api-deployment -n vibe-dev
```

## 下一步

部署成功后，您可以：
1. 运行 `./test-api.sh` 验证 API 功能
2. 使用 `make k8s-logs` 查看应用日志
3. 通过 `http://api.vibe-dev.com:8000/metrics` 查看监控指标
4. 进行 API 开发和调试
5. 扩展部署配置以支持生产环境

更多详细信息请参考：
- `deploy/k8s/README.md` - K8s 部署资源说明
- `dev-tutorial/k3d/README.md` - k3d 环境详细指南
- 项目根目录的 `Makefile` - 构建和部署命令
