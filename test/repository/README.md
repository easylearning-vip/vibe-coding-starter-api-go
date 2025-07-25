# Repository 单元测试

本目录包含了 vibe-coding-starter 项目中所有 Repository 层的单元测试。

## 测试概述

Repository 层测试覆盖了以下6个核心仓储：

- **UserRepository**: 用户数据访问
- **ArticleRepository**: 文章数据访问  
- **CategoryRepository**: 分类数据访问
- **TagRepository**: 标签数据访问
- **CommentRepository**: 评论数据访问
- **FileRepository**: 文件数据访问

## 测试环境

测试使用 Docker 容器提供隔离的数据库环境：

- **MySQL 5.7**: 主数据库 (端口: 3307)
- **Redis 7**: 缓存数据库 (端口: 6380)

## 快速开始

### 方法1: 使用自动化脚本 (推荐)

```bash
# 运行完整的Repository测试 (包含Docker环境管理)
./test/scripts/run_repository_tests.sh
```

### 方法2: 使用Makefile命令

```bash
# 使用Docker环境运行Repository测试
make test-repository-docker

# 或者只运行Repository测试 (需要手动管理环境)
make test-repository
```

### 方法3: 手动运行

```bash
# 1. 启动测试环境
make -f Makefile.test test-setup

# 2. 运行测试
go test -v -race -short ./test/repository/...

# 3. 清理环境
make -f Makefile.test test-teardown
```

## 测试结构

每个Repository测试都遵循相同的结构：

```go
type XxxRepositoryTestSuite struct {
    suite.Suite
    db     *testutil.TestDatabase
    cache  *testutil.TestCache  
    logger *testutil.TestLogger
    repo   repository.XxxRepository
    ctx    context.Context
}
```

### 测试生命周期

1. **SetupSuite**: 初始化测试环境和依赖
2. **SetupTest**: 每个测试前清理数据
3. **测试执行**: 运行具体测试用例
4. **TearDownSuite**: 清理测试环境

## 测试用例类型

### 基础CRUD操作
- `TestCreate`: 创建实体
- `TestGetByID`: 按ID查询
- `TestUpdate`: 更新实体
- `TestDelete`: 删除实体

### 查询操作
- `TestList`: 列表查询
- `TestListWithFilters`: 条件过滤
- `TestListWithSearch`: 搜索查询
- `TestPagination`: 分页查询

### 业务查询
- `TestGetByEmail`: 按邮箱查询用户
- `TestGetBySlug`: 按slug查询文章/分类
- `TestGetByAuthor`: 按作者查询文章
- `TestGetByCategory`: 按分类查询文章

### 错误处理
- `TestGetByIDNotFound`: 记录不存在
- `TestCreateDuplicateXxx`: 重复数据检测

## 运行单个测试

```bash
# 运行特定Repository的测试
go test -v -run TestUserRepository ./test/repository/

# 运行特定测试用例
go test -v -run TestUserRepository/TestCreate ./test/repository/

# 运行所有Create测试
go test -v -run TestCreate ./test/repository/
```

## 测试配置

测试配置位于 `test/config/test_config.go`：

```go
// 数据库配置
Database: config.DatabaseConfig{
    Driver:   "mysql",
    Host:     "localhost",
    Port:     3307,
    Username: "test_user", 
    Password: "test_password",
    Database: "vibe_coding_starter_test",
    Charset:  "utf8mb4",
}

// 缓存配置  
Cache: config.CacheConfig{
    Driver:   "redis",
    Host:     "localhost", 
    Port:     6380,
    Database: 1,
}
```

## 测试工具

### TestDatabase
提供数据库测试支持：
- 自动迁移表结构
- 测试间数据清理
- 连接池管理

### TestCache  
提供缓存测试支持：
- Redis连接管理
- 测试间缓存清理

### TestLogger
提供日志测试支持：
- 结构化日志输出
- 测试日志隔离

## 故障排除

### Docker相关问题

```bash
# 检查Docker状态
sudo docker ps

# 查看容器日志
sudo docker logs vibe-mysql-test
sudo docker logs vibe-redis-test

# 重启测试环境
make -f Makefile.test test-teardown
make -f Makefile.test test-setup
```

### 数据库连接问题

```bash
# 测试MySQL连接
sudo docker exec vibe-mysql-test mysqladmin ping -h localhost

# 测试Redis连接  
sudo docker exec vibe-redis-test redis-cli ping
```

### 权限问题

如果遇到Docker权限问题：

```bash
# 将用户添加到docker组
sudo usermod -aG docker $USER

# 重新登录或运行
newgrp docker
```

## 测试最佳实践

1. **数据隔离**: 每个测试用例使用独立的数据
2. **并发安全**: 使用 `-race` 标志检测竞态条件
3. **错误验证**: 验证错误类型和消息
4. **边界测试**: 测试边界条件和异常情况
5. **性能考虑**: 避免不必要的数据库操作

## 贡献指南

添加新的Repository测试时：

1. 创建对应的测试文件 `xxx_repository_test.go`
2. 实现 `XxxRepositoryTestSuite` 结构
3. 添加所有CRUD和业务方法的测试
4. 确保测试覆盖错误情况
5. 运行测试确保通过

## 相关文档

- [Repository Test Report](../../REPOSITORY_TEST_REPORT.md): 详细测试报告
- [Makefile.test](../../Makefile.test): 测试环境管理
- [Docker Compose](../../docker-compose.test.yml): 测试服务配置
