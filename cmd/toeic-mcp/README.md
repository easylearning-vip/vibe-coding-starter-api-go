# TOEIC MCP Server

基于 Model Context Protocol (MCP) 的 TOEIC 练习服务器，为 Claude Code 和其他 MCP 客户端提供 TOEIC 学习工具。

## 功能特性

- ✅ **TOEIC Part 2 练习**: 完整的听力练习管理系统
- ✅ **答案评分**: 即时反馈和详细解析
- ✅ **AI 提示词**: 内置 TOEIC 学习提示词库
- ✅ **多用户支持**: 基于 Token 的用户认证
- ✅ **Streamable HTTP**: 支持双向实时通信

## 快速开始

### 1. 启动服务器

```bash
go run ./cmd/toeic-mcp/main.go -c config/dev.yaml
```

### 2. 配置 Claude CLI

```bash
claude mcp add toeic-mcp http://localhost:6275 \
  --transport http \
  --headers "X-User-Token: YOUR_TOKEN"
```

### 3. 开始使用

在 Claude Code 中:
```
> 请使用 TOEIC MCP 工具获取所有 Part 2 练习集
```

详细步骤请查看 [快速开始指南](./QUICKSTART.md)

## 文档

- 🚀 [快速开始](./QUICKSTART.md) - 5 分钟快速配置
- 📖 [完整配置教程](./MCP_CONFIGURATION_GUIDE.md) - 详细的配置和使用说明
- 🔐 [认证配置详解](./AUTHENTICATION.md) - Token 认证机制和安全配置
- 📝 [Part 2 提示词模板](./part2_prompt.md) - AI 练习提示词
- 📝 [Part 3 提示词模板](./part3_prompt.md) - AI 练习提示词
- 📝 [Part 4 提示词模板](./part4_prompt.md) - AI 练习提示词

## 配置示例

### 基本配置

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "YOUR_TOKEN_HERE"
      }
    }
  }
}
```

更多配置示例:
- [基本配置](./mcp-config-example.json) - 使用硬编码 Token
- [环境变量配置](./mcp-config-with-env.json) - 使用环境变量（推荐）
- [Bearer Token 配置](./mcp-config-bearer.json) - 使用 Bearer 认证
- [远程服务器配置](./mcp-config-remote.json) - HTTPS 远程服务器

## 可用工具

| 工具名称 | 功能描述 |
|---------|---------|
| `list_part2_sets` | 列出所有 Part 2 练习集 |
| `get_part2_set_details` | 获取练习集详情和题目 |
| `submit_part2_answer` | 提交答案并获取评分 |
| `list_toeic_prompts` | 列出所有 AI 提示词 |
| `get_toeic_prompt` | 获取指定提示词内容 |

## 认证方式

服务器支持三种 Token 传递方式:

```bash
# 方式 1: X-User-Token Header (推荐)
X-User-Token: your_token

# 方式 2: Authorization Token
Authorization: Token your_token

# 方式 3: Authorization Bearer
Authorization: Bearer your_token
```

## 命令行参数

```bash
# 查看帮助
go run ./cmd/toeic-mcp/main.go -h

# 常用参数
-c string       # 配置文件路径
--host string   # 监听主机 (默认: 0.0.0.0)
--port int      # 监听端口 (默认: 6275)
```

## 环境变量

```bash
# 通过环境变量指定配置文件
export CONFIG_FILE=/path/to/config.yaml
go run ./cmd/toeic-mcp/main.go
```

## 架构说明

### 技术栈

- **语言**: Go 1.21+
- **框架**: Uber FX (依赖注入)
- **数据库**: MySQL (通过 GORM)
- **传输协议**: Streamable HTTP (MCP)
- **认证**: Token-based (每请求验证)

### 目录结构

```
cmd/toeic-mcp/
├── main.go                    # 服务器入口
├── README.md                  # 本文件
├── QUICKSTART.md              # 快速开始
├── CLAUDE_CLI_SETUP.md        # 完整配置教程
├── part2_prompt.md            # Part 2 提示词
├── mcp-config-example.json    # 配置示例
└── mcp-config-with-env.json   # 环境变量配置示例

internal/handler/mcp/
├── server.go                  # MCP 服务器核心
├── auth.go                    # 认证逻辑
├── tools_part2.go             # Part 2 工具
└── tools_prompt.go            # 提示词工具
```

## 开发计划

- [ ] Part 3/4 听力练习支持
- [ ] Part 5/6 语法练习支持
- [ ] Part 7 阅读练习支持
- [ ] 学习进度追踪
- [ ] 错题本功能
- [ ] 模拟考试模式

## 故障排查

### 连接问题

```bash
# 检查服务器状态
curl http://localhost:6275/health

# 检查端口占用
lsof -i :6275
```

### 认证问题

```bash
# 测试 Token
curl -H "X-User-Token: YOUR_TOKEN" http://localhost:6275/

# 查看服务器日志
# 日志会显示认证失败的详细信息
```

### Claude CLI 问题

```bash
# 重新加载配置
claude mcp reload

# 查看服务器状态
claude mcp list toeic-mcp

# 查看详细日志
claude mcp logs toeic-mcp
```

## 安全建议

1. 🔒 不要在代码仓库中提交包含 Token 的配置文件
2. 🔄 定期轮换用户 Token
3. 🌐 生产环境使用 HTTPS
4. 🛡️ 配置防火墙规则限制访问
5. 📊 监控异常访问日志

## 相关资源

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [Claude CLI 文档](https://docs.claude.com/en/docs/claude-code/mcp)
- [项目主页](https://github.com/alvinzane/doubao-toeic)

## 许可证

[项目许可证信息]

## 贡献

欢迎提交 Issue 和 Pull Request!

## 联系方式

- GitHub: [@alvinzane](https://github.com/alvinzane)
- Email: awary@qq.com

