# 多阶段构建 Dockerfile for Vibe Coding Starter

# 构建阶段
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置 Alpine 镜像源为阿里云镜像
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 设置 Go 模块代理为中国镜像
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.BuildTime=$(date +%Y-%m-%d_%H:%M:%S)" \
    -a -installsuffix cgo \
    -o main cmd/server/main.go

# 运行阶段
FROM alpine:latest


# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件和迁移文件
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

# 创建必要的目录
RUN mkdir -p logs uploads && \
    chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 启动应用程序
CMD ["./main"]
