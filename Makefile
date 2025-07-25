# Vibe Coding Starter API - Makefile
# 用于构建、测试和部署应用程序

.PHONY: help build test clean docker-build docker-push k8s-deploy k8s-clean dev-setup

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
APP_NAME := vibe-coding-starter-api
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
REGISTRY := localhost:5000
IMAGE_TAG := $(REGISTRY)/$(APP_NAME):latest

# Go 相关变量
GO_VERSION := 1.23
GOOS := linux
GOARCH := amd64
CGO_ENABLED := 0

# 构建标志
LDFLAGS := -w -s -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)
BUILD_FLAGS := -a -installsuffix cgo -ldflags="$(LDFLAGS)"

help: ## 显示帮助信息
	@echo "Vibe Coding Starter API - 可用命令:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# 开发环境设置
dev-setup: ## 设置开发环境
	@echo "设置开发环境..."
	go mod tidy
	go mod download
	@echo "开发环境设置完成"

# 构建相关
build: ## 构建应用程序
	@echo "构建应用程序..."
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o bin/server cmd/server/main.go
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(BUILD_FLAGS) -o bin/migrate cmd/migrate/main.go
	@echo "构建完成: bin/server, bin/migrate"

build-local: ## 构建本地版本
	@echo "构建本地版本..."
	go build -o bin/server-local cmd/server/main.go
	go build -o bin/migrate-local cmd/migrate/main.go
	@echo "本地构建完成"

# 测试相关
test: ## 运行所有测试
	@echo "运行测试..."
	go test -v -race -short ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试覆盖率..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告生成: coverage.html"

# Docker 相关
docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像..."
	docker build -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest .
	docker tag $(APP_NAME):latest $(IMAGE_TAG)
	@echo "Docker 镜像构建完成: $(IMAGE_TAG)"

docker-push: docker-build ## 推送 Docker 镜像到本地仓库
	@echo "推送 Docker 镜像..."
	docker push $(IMAGE_TAG)
	@echo "镜像推送完成: $(IMAGE_TAG)"

docker-run: ## 运行 Docker 容器
	@echo "运行 Docker 容器..."
	docker run --rm -p 8080:8080 --name $(APP_NAME) $(IMAGE_TAG)

# K8s 相关
k8s-deploy: docker-push ## 部署到 k8s 开发环境
	@echo "部署到 k8s 开发环境..."
	cd deploy/k8s && ./deploy.sh

k8s-clean: ## 清理 k8s 部署
	@echo "清理 k8s 部署..."
	cd deploy/k8s && ./deploy.sh clean

k8s-status: ## 查看 k8s 部署状态
	@echo "查看 k8s 部署状态..."
	kubectl get all -n vibe-dev -l app=vibe-api

k8s-logs: ## 查看 k8s 应用日志
	@echo "查看应用日志..."
	kubectl logs -f deployment/vibe-api-deployment -n vibe-dev

# 开发服务器
dev-docker: ## 使用 Docker Compose 启动开发环境
	@echo "启动 Docker 开发环境..."
	cd dev-tutorial/docker-compose && docker compose -f docker-compose.dev.yml up -d

dev-k3d: ## 使用 k3d 启动开发环境
	@echo "启动 k3d 开发环境..."
	cd dev-tutorial/k3d && k3d cluster create --config k3d-cluster.yaml
	kubectl apply -f dev-tutorial/k3d/manifests/

run-local: build-local ## 运行本地开发服务器
	@echo "启动本地开发服务器..."
	./bin/server-local -c configs/config.yaml

run-docker: ## 使用 Docker 配置运行本地服务器
	@echo "使用 Docker 配置运行服务器..."
	go run cmd/server/main.go -c configs/config-docker.yaml

run-k3d: ## 使用 k3d 配置运行本地服务器
	@echo "使用 k3d 配置运行服务器..."
	go run cmd/server/main.go -c configs/config-k3d.yaml

# 数据库迁移
migrate-up: ## 执行数据库迁移
	@echo "执行数据库迁移..."
	go run cmd/migrate/main.go -c configs/config.yaml up

migrate-down: ## 回滚数据库迁移
	@echo "回滚数据库迁移..."
	go run cmd/migrate/main.go -c configs/config.yaml down

migrate-version: ## 查看迁移版本
	@echo "查看迁移版本..."
	go run cmd/migrate/main.go -c configs/config.yaml version

# 清理
clean: ## 清理构建文件
	@echo "清理构建文件..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	docker rmi $(APP_NAME):latest $(APP_NAME):$(VERSION) $(IMAGE_TAG) 2>/dev/null || true
	@echo "清理完成"

# 格式化和检查
fmt: ## 格式化代码
	@echo "格式化代码..."
	go fmt ./...
	goimports -w .

lint: ## 运行代码检查
	@echo "运行代码检查..."
	golangci-lint run

# 完整的开发流程
dev-full: dev-setup test docker-build k8s-deploy ## 完整的开发部署流程
	@echo "完整开发部署流程完成"

# 快速部署
quick-deploy: docker-push k8s-deploy ## 快速重新部署
	@echo "快速重新部署完成"