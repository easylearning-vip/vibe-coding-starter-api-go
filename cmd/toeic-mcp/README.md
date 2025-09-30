# TOEIC MCP Server

TOEIC MCP (Model Context Protocol) Server 是一个基于 Go 的 MCP 服务器实现，为 AI 助手提供 TOEIC 练习相关的工具和功能。

## 架构设计

本项目遵循清洁架构（Clean Architecture）设计模式，确保代码的可测试性、可维护性和可扩展性。

### 目录结构

```
cmd/toeic-mcp/
├── main.go                          # 应用入口，依赖注入配置
└── README.md                        # 本文档

internal/
├── handler/mcp/                     # MCP处理器层
│   ├── server.go                    # MCP服务器构建和工具注册
│   ├── auth.go                      # Token认证逻辑
│   ├── user_info.go                 # 用户信息工具处理器
│   └── part2_practice_set.go        # Part2练习集工具处理器
├── model/mcp/                       # MCP模型层
│   └── params.go                    # MCP工具参数定义
├── service/                         # 业务逻辑层（复用现有服务）
│   └── part2_practice_set.go       # Part2练习集服务
└── repository/                      # 数据访问层（复用现有仓储）
    ├── user.go                      # 用户仓储
    └── part2_*.go                   # Part2相关仓储
```

### 架构层次

#### 1. 应用层 (cmd/toeic-mcp/main.go)
- **职责**: 应用初始化、依赖注入、服务器启动
- **特点**: 
  - 使用 Uber FX 进行依赖注入
  - 配置管理（支持命令行参数和环境变量）
  - 生命周期管理（优雅启动和关闭）
  - 最小化业务逻辑

#### 2. MCP处理器层 (internal/handler/mcp/)
- **职责**: MCP协议处理、工具注册、请求路由
- **组件**:
  - `server.go`: MCP服务器构建器，负责创建和配置MCP服务器
  - `auth.go`: Token提取和验证逻辑
  - `user_info.go`: 用户信息工具的定义和处理
  - `part2_practice_set.go`: Part2练习集工具的定义和处理

#### 3. 模型层 (internal/model/mcp/)
- **职责**: MCP工具的参数和响应模型定义
- **特点**: 
  - 使用 JSON Schema 标签支持参数验证
  - 清晰的类型定义

#### 4. 服务层 (internal/service/)
- **职责**: 业务逻辑实现
- **特点**: 复用现有的业务服务，无需重复实现

#### 5. 仓储层 (internal/repository/)
- **职责**: 数据访问抽象
- **特点**: 复用现有的数据访问层

## 功能特性

### 认证机制
- 基于 Token 的用户认证
- 支持多种 Token 格式:
  - `X-User-Token: <token>`
  - `Authorization: Token <token>`
  - `Authorization: Bearer <token>`
- 每个请求独立验证，支持多用户并发访问

### 可用工具

#### 1. get_user_info
查看当前用户信息（ID和名称）

**参数**: 无

**返回示例**:
```
user_id=123
name=张三
```

#### 2. list_part2_sets
查看个人 Part2 练习集列表（支持分页）

**参数**:
- `page`: 页码（默认 1）
- `page_size`: 每页大小（默认 20，最大 100）
- `sort`: 排序字段（可选）
- `order`: 排序方向 asc|desc（可选）

**返回示例**:
```
total=50 page=1 size=20
#1 total=30 done=15 correct=12 acc=80.00%
#2 total=25 done=25 correct=20 acc=80.00%
```

#### 3. auto_generate_part2_set
自动生成 Part2 练习集

**参数**:
- `count`: 题目数量（必需，10-30）
- `difficulty_level_id`: 难度级别ID（可选）
- `scenario_id`: 场景ID（可选）

**返回示例**:
```
✅ 练习集生成成功

练习集ID: 123
题目总数: 15
已添加 15 道题目到练习集
```

#### 4. get_part2_set_details
查看 Part2 练习集明细列表（包含题目详情）

**参数**:
- `set_id`: 练习集ID（必需）

**返回示例**:
```
📋 练习集详情

练习集ID: 123
题目总数: 15
已完成: 5
正确数: 4
准确率: 26.67%

题目列表:
─────────────────────────────────────

题目 1 (Item ID: 456)
问题: Where is the meeting?
A. In the conference room
B. At the restaurant
C. In the office
```

#### 5. submit_part2_answer
提交 Part2 练习题答案

**参数**:
- `set_id`: 练习集ID（必需）
- `item_id`: 题目项ID（必需）
- `user_answer`: 用户答案（必需，A/B/C）

**返回示例**:
```
✅ 回答正确！

你的答案: A
正确答案: A
```

**详细文档**: 查看 [NEW_TOOLS.md](./NEW_TOOLS.md) 了解新工具的完整文档

## 使用方法

### 启动服务器

```bash
# 使用默认配置
go run ./cmd/toeic-mcp/

# 指定配置文件
go run ./cmd/toeic-mcp/ -c configs/config.yaml

# 指定主机和端口
go run ./cmd/toeic-mcp/ -host 0.0.0.0 -port 6275

# 使用环境变量指定配置
CONFIG_FILE=configs/config.yaml go run ./cmd/toeic-mcp/
```

### 构建可执行文件

```bash
# 构建
go build -o toeic-mcp ./cmd/toeic-mcp/

# 运行
./toeic-mcp -c configs/config.yaml
```

### 客户端配置

在 AI 工具（如 Augment、Cursor）中配置 MCP 服务器：

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "your-user-token-here"
      }
    }
  }
}
```

## 开发指南

### 添加新工具

1. **定义参数模型** (internal/model/mcp/params.go)
```go
type NewToolParams struct {
    Field1 string `json:"field1" jsonschema:"Description of field1"`
    Field2 int    `json:"field2" jsonschema:"Description of field2"`
}
```

2. **创建工具处理器** (internal/handler/mcp/new_tool.go)
```go
package mcp

import (
    "context"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    mcpModel "vibe-coding-starter/internal/model/mcp"
)

// NewToolTool 创建工具定义
func NewToolTool() *mcp.Tool {
    return &mcp.Tool{
        Name:        "new_tool",
        Description: "工具描述",
    }
}

// HandleNewTool 处理工具请求
func HandleNewTool(/* dependencies */) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.NewToolParams) (*mcp.CallToolResult, any, error) {
    return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.NewToolParams) (*mcp.CallToolResult, any, error) {
        // 实现业务逻辑
        return &mcp.CallToolResult{
            Content: []mcp.Content{&mcp.TextContent{Text: "result"}},
        }, nil, nil
    }
}
```

3. **注册工具** (internal/handler/mcp/server.go)
```go
func (h *MCPHandler) registerTools(server *mcp.Server, user *model.User) {
    // ... 现有工具 ...
    
    // 注册新工具
    mcp.AddTool(server, NewToolTool(), HandleNewTool(/* dependencies */))
}
```

4. **添加依赖** (如果需要新的服务或仓储)
```go
// 在 MCPHandler 结构体中添加
type MCPHandler struct {
    // ... 现有字段 ...
    newService service.NewService
}

// 在 NewMCPHandler 中注入
func NewMCPHandler(
    // ... 现有参数 ...
    newService service.NewService,
) *MCPHandler {
    return &MCPHandler{
        // ... 现有字段 ...
        newService: newService,
    }
}
```

5. **更新依赖注入** (cmd/toeic-mcp/main.go)
```go
fx.Provide(
    // ... 现有提供者 ...
    service.NewNewService,
),
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/handler/mcp/...

# 运行测试并显示覆盖率
go test -cover ./...
```

## 设计原则

### 1. 依赖倒置原则 (DIP)
- 高层模块（Handler）不依赖低层模块（Repository）
- 都依赖于抽象接口（Service Interface）
- 通过依赖注入实现解耦

### 2. 单一职责原则 (SRP)
- 每个模块只有一个变化的理由
- Handler 只负责 MCP 协议处理
- Service 只负责业务逻辑
- Repository 只负责数据访问

### 3. 开闭原则 (OCP)
- 对扩展开放：可以轻松添加新工具
- 对修改关闭：添加新功能不需要修改现有代码

### 4. 接口隔离原则 (ISP)
- 使用小而专注的接口
- 避免臃肿的接口定义

## 性能考虑

- **无状态设计**: 每个请求独立处理，支持水平扩展
- **连接池**: 复用数据库连接
- **并发安全**: 使用 context 传递请求上下文
- **优雅关闭**: 支持正在处理的请求完成后再关闭

## 安全性

- **Token 验证**: 每个请求都验证用户 Token
- **用户隔离**: 每个用户只能访问自己的数据
- **输入验证**: 使用 JSON Schema 验证参数
- **错误处理**: 不泄露敏感信息

## 故障排查

### 常见问题

1. **Token 认证失败**
   - 检查 Token 格式是否正确
   - 确认用户状态是否为 Active
   - 查看服务器日志获取详细错误信息

2. **服务启动失败**
   - 检查配置文件路径是否正确
   - 确认数据库连接配置
   - 查看端口是否被占用

3. **工具调用失败**
   - 检查参数格式是否符合 JSON Schema
   - 查看服务器日志获取详细错误
   - 确认用户有权限访问相关资源

## 参考资料

- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Uber FX Documentation](https://uber-go.github.io/fx/)

