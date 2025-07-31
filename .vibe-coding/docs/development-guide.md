# Vibe Coding Starter Go API - 二次开发指南

## 项目概述

Vibe Coding Starter Go API 是一个现代化的 Go Web API 项目模板，采用清洁架构设计，集成了企业级开发所需的核心组件。项目完全由 AI 工具开发，零人工代码编写，为 Vibe Coding 快速开发提供最佳实践示范。

### 核心特性
- **清洁架构**: 分层设计，职责分离，易于测试和维护
- **依赖注入**: 使用 Uber FX 框架实现依赖注入
- **多数据库支持**: MySQL、PostgreSQL、SQLite
- **缓存系统**: Redis 集成
- **认证授权**: JWT Token 认证
- **API 文档**: Swagger 自动生成
- **容器化**: Docker 和 Kubernetes 支持
- **完整测试**: 单元测试、集成测试覆盖

## 项目架构

### 目录结构
```
vibe-coding-starter-api-go/
├── cmd/                    # 应用程序入口点
│   ├── server/            # HTTP 服务器
│   ├── migrate/           # 数据库迁移工具
│   └── hello/             # 示例命令
├── internal/              # 私有应用代码
│   ├── config/            # 配置管理
│   ├── handler/           # HTTP 处理器（控制器层）
│   ├── middleware/        # 中间件
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   ├── server/            # 服务器配置
│   └── service/           # 业务逻辑层
├── pkg/                   # 可重用的库代码
│   ├── cache/             # 缓存抽象
│   ├── database/          # 数据库抽象
│   ├── logger/            # 日志抽象
│   └── migration/         # 迁移工具
├── configs/               # 配置文件
├── migrations/            # 数据库迁移脚本
├── test/                  # 测试代码
├── deploy/                # 部署配置
├── tutorials/develop/          # 开发环境教程
└── tools/                 # 开发工具
```

### 技术栈
- **Web 框架**: Gin (高性能 HTTP 框架)
- **ORM**: GORM (Go 对象关系映射)
- **缓存**: Redis
- **数据库**: MySQL/PostgreSQL/SQLite
- **依赖注入**: Uber FX
- **日志**: Zap (结构化日志)
- **配置**: Viper (配置管理)
- **认证**: JWT
- **测试**: Testify
- **文档**: Swagger
- **容器**: Docker + Kubernetes

### 架构层次

#### 1. 表示层 (Handler)
- 处理 HTTP 请求和响应
- 参数验证和数据绑定
- 路由定义和中间件应用
- 位置: `internal/handler/`

#### 2. 业务逻辑层 (Service)
- 核心业务逻辑实现
- 数据处理和业务规则
- 跨领域服务协调
- 位置: `internal/service/`

#### 3. 数据访问层 (Repository)
- 数据库操作抽象
- 查询构建和执行
- 数据持久化
- 位置: `internal/repository/`

#### 4. 基础设施层 (Pkg)
- 数据库连接管理
- 缓存操作
- 日志记录
- 位置: `pkg/`

## 开发环境搭建

### 前置要求
- Go 1.23+
- Docker & Docker Compose
- Make (可选，用于构建脚本)

### 快速开始
```bash
# 1. 克隆项目
git clone <repository-url>
cd vibe-coding-starter-api-go

# 2. 安装依赖
go mod tidy

# 3. 启动开发环境 (Docker Compose)
make dev-docker

# 4. 运行数据库迁移
make migrate-up

# 5. 启动应用
make run-docker
```

### 开发环境选项

#### 选项 1: Docker Compose (推荐)
```bash
cd tutorials/develop/docker-compose
docker compose -f docker-compose.dev.yml up -d
```

#### 选项 2: K3D (Kubernetes)
```bash
cd tutorials/develop/k3d
k3d cluster create --config k3d-cluster.yaml
```

#### 选项 3: 本地开发
```bash
# 需要本地安装 MySQL/PostgreSQL 和 Redis
make run-local
```

## 开发规范

### 代码组织原则
1. **单一职责**: 每个模块只负责一个功能
2. **依赖倒置**: 高层模块不依赖低层模块
3. **接口隔离**: 使用接口定义契约
4. **开闭原则**: 对扩展开放，对修改关闭

### 命名规范
- **包名**: 小写，简短，有意义 (如 `user`, `article`)
- **文件名**: 小写，下划线分隔 (如 `user_service.go`)
- **结构体**: 大驼峰 (如 `UserService`)
- **方法/函数**: 大驼峰公开，小驼峰私有
- **常量**: 大写，下划线分隔 (如 `MAX_RETRY_COUNT`)

### 错误处理
```go
// 使用 fmt.Errorf 包装错误
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// 定义业务错误类型
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

### 日志规范
```go
// 使用结构化日志
logger.Info("User created successfully", 
    "user_id", user.ID,
    "username", user.Username,
)

logger.Error("Failed to create user",
    "error", err,
    "username", req.Username,
)
```

## 开发流程

### 1. 功能开发流程
```
需求分析 → 设计接口 → 编写测试 → 实现功能 → 代码审查 → 部署测试
```

### 2. 分支管理
- `main`: 主分支，生产环境代码
- `develop`: 开发分支，集成测试
- `feature/*`: 功能分支
- `hotfix/*`: 紧急修复分支

### 3. 提交规范
```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式调整
refactor: 重构
test: 测试相关
chore: 构建/工具相关
```

### 4. 代码审查检查点
- [ ] 代码符合项目规范
- [ ] 单元测试覆盖率 > 80%
- [ ] 无安全漏洞
- [ ] 性能影响评估
- [ ] 文档更新完整

## 测试策略

### 测试层次
1. **单元测试**: 测试单个函数/方法
2. **集成测试**: 测试模块间交互
3. **端到端测试**: 测试完整业务流程

### 测试工具
- **测试框架**: `testify/suite`
- **Mock 工具**: 手动 Mock 实现
- **测试数据库**: SQLite (内存模式)
- **测试缓存**: 内存实现

### 运行测试
```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行特定包的测试
go test -v ./internal/service/...
```

### 测试示例
```go
func TestUserService_Create(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}

type UserServiceTestSuite struct {
    suite.Suite
    service UserService
    mockRepo *MockUserRepository
}

func (s *UserServiceTestSuite) SetupTest() {
    s.mockRepo = &MockUserRepository{}
    s.service = NewUserService(s.mockRepo)
}

func (s *UserServiceTestSuite) TestCreate_Success() {
    // Given
    req := &CreateUserRequest{
        Username: "testuser",
        Email:    "test@example.com",
    }
    
    // When
    user, err := s.service.Create(context.Background(), req)
    
    // Then
    s.NoError(err)
    s.Equal("testuser", user.Username)
}
```

## AI Prompting 工程指南

### 项目上下文提示词
```
这是一个 Go Web API 项目，使用清洁架构设计：
- 使用 Gin 框架和 GORM ORM
- 采用依赖注入 (Uber FX)
- 分层架构：Handler -> Service -> Repository
- 支持 MySQL/PostgreSQL/SQLite 数据库
- 集成 Redis 缓存和 JWT 认证
- 完整的测试覆盖和 Docker 部署支持
```

### 功能开发提示词模板
```
请为 [功能名称] 实现以下组件：

1. 数据模型 (internal/model/)
2. Repository 接口和实现 (internal/repository/)
3. Service 接口和实现 (internal/service/)
4. Handler 实现 (internal/handler/)
5. 路由注册
6. 单元测试

要求：
- 遵循项目现有的代码风格和架构模式
- 包含完整的错误处理和日志记录
- 提供 Swagger 文档注释
- 编写对应的单元测试
```

### 调试提示词
```
项目出现 [错误描述]，请帮助分析：

环境信息：
- Go 版本: 1.23
- 数据库: [MySQL/PostgreSQL/SQLite]
- 部署方式: [Docker/K8s/本地]

错误日志：
[粘贴错误日志]

请提供：
1. 问题根因分析
2. 解决方案
3. 预防措施
```

## 常见开发任务

### 添加新的 API 端点
1. 定义数据模型 (`internal/model/`)
2. 创建 Repository 接口和实现
3. 实现 Service 业务逻辑
4. 创建 Handler 处理 HTTP 请求
5. 注册路由
6. 编写测试
7. 更新 API 文档

### 数据库迁移
```bash
# 创建新迁移
go run cmd/migrate/main.go create add_user_table

# 执行迁移
make migrate-up

# 回滚迁移
make migrate-down
```

### 添加中间件
```go
// 在 internal/middleware/ 中实现
func NewCustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
        c.Next()
    }
}

// 在路由中使用
router.Use(NewCustomMiddleware())
```

## 部署指南

### Docker 部署
```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

### Kubernetes 部署
```bash
# 部署到 K8s
make k8s-deploy

# 查看状态
make k8s-status

# 查看日志
make k8s-logs
```

## 故障排查

### 常见问题
1. **数据库连接失败**: 检查配置文件和网络连接
2. **Redis 连接失败**: 确认 Redis 服务状态
3. **JWT Token 无效**: 检查密钥配置和过期时间
4. **测试失败**: 确认测试数据库配置

### 日志查看
```bash
# 应用日志
kubectl logs -f deployment/vibe-api-deployment -n vibe-dev

# 数据库日志
docker logs vibe-mysql

# Redis 日志
docker logs vibe-redis
```

## 性能优化

### 数据库优化
- 合理使用索引
- 避免 N+1 查询问题
- 使用连接池
- 定期分析慢查询

### 缓存策略
- 热点数据缓存
- 查询结果缓存
- 会话缓存
- 合理设置过期时间

### API 性能
- 请求限流
- 响应压缩
- 分页查询
- 异步处理

## 安全考虑

### 认证授权
- JWT Token 管理
- 角色权限控制
- API 访问限制
- 敏感数据加密

### 输入验证
- 参数验证
- SQL 注入防护
- XSS 防护
- CSRF 防护

## 扩展指南

### 添加新的数据库支持
1. 在 `pkg/database/` 中添加驱动支持
2. 更新配置结构
3. 添加相应的迁移脚本
4. 更新文档

### 集成第三方服务
1. 在 `pkg/` 中创建抽象接口
2. 实现具体的服务客户端
3. 通过依赖注入集成
4. 添加配置选项
5. 编写测试

---

## 总结

本项目采用现代化的 Go 开发最佳实践，通过清洁架构和依赖注入实现了高度可测试和可维护的代码结构。无论是初学者学习还是 AI 辅助开发，都能快速上手并进行二次开发。

遵循本指南的规范和流程，可以确保代码质量和项目的长期可维护性。
