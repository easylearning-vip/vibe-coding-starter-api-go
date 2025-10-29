# TOEIC MCP 快速开始指南

5 分钟快速配置 TOEIC MCP 服务器并开始使用。

## 前置要求

- Go 1.21 或更高版本
- MySQL 数据库（已配置）
- Claude Desktop、Cursor 或其他 MCP 客户端

## 步骤 1: 启动服务器

```bash
# 进入项目目录
cd /path/to/vibe-coding-starter-api-go

# 启动 MCP 服务器
go run ./cmd/toeic-mcp/main.go -c config/dev.yaml
```

您应该看到类似的输出：

```
INFO  Starting TOEIC MCP server (streamable HTTP) addr=0.0.0.0:6275 auth_header="Authorization: Token <token> or X-User-Token: <token>"
```

## 步骤 2: 获取用户 Token

连接到数据库并获取您的用户 Token：

```sql
-- 查询您的 Token
SELECT id, username, email, token FROM users WHERE username = 'your_username';

-- 如果 token 为空，生成一个新的
UPDATE users SET token = 'my_secure_token_123' WHERE username = 'your_username';
```

## 步骤 3: 配置客户端

### 选项 A: Claude Desktop（推荐）

1. 找到配置文件：
   - **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
   - **Linux**: `~/.config/Claude/claude_desktop_config.json`

2. 编辑配置文件：

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "my_secure_token_123"
      }
    }
  }
}
```

3. 完全退出并重启 Claude Desktop

### 选项 B: Cursor

1. 在项目根目录创建 `.cursor/mcp.json`：

```bash
mkdir -p .cursor
cat > .cursor/mcp.json << 'EOF'
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "my_secure_token_123"
      }
    }
  }
}
EOF
```

2. 重启 Cursor

### 选项 C: 使用环境变量（更安全）

1. 设置环境变量：

```bash
# Linux/macOS
export TOEIC_USER_TOKEN="my_secure_token_123"

# Windows PowerShell
$env:TOEIC_USER_TOKEN="my_secure_token_123"
```

2. 使用环境变量的配置：

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "${env:TOEIC_USER_TOKEN}"
      }
    }
  }
}
```

## 步骤 4: 测试连接

### 在 Claude Desktop 中测试

打开 Claude Desktop，输入：

```
请使用 TOEIC MCP 工具列出所有 Part 2 练习集
```

或者：

```
帮我获取 TOEIC Part 2 的练习题目
```

### 使用 curl 测试

```bash
# 测试初始化
curl -H "X-User-Token: my_secure_token_123" \
     -H "Content-Type: application/json" \
     -X POST http://localhost:6275/ \
     -d '{
       "jsonrpc": "2.0",
       "id": 1,
       "method": "initialize",
       "params": {
         "protocolVersion": "2024-11-05",
         "capabilities": {},
         "clientInfo": {
           "name": "test",
           "version": "1.0"
         }
       }
     }'

# 测试工具列表
curl -H "X-User-Token: my_secure_token_123" \
     -H "Content-Type: application/json" \
     -X POST http://localhost:6275/ \
     -d '{
       "jsonrpc": "2.0",
       "id": 2,
       "method": "tools/list",
       "params": {}
     }'
```

## 可用工具

成功连接后，您可以使用以下工具：

### Part 2 工具
- `list_part2_sets` - 列出所有 Part 2 练习集
- `get_part2_set_details` - 获取练习集详情和题目
- `submit_part2_answer` - 提交答案并获取评分
- `auto_generate_part2_set` - 自动生成新的练习集
- `get_part2_prompt` - 获取 Part 2 AI 提示词

### Part 3 工具
- `list_part3_sets` - 列出所有 Part 3 练习集
- `get_part3_set_details` - 获取练习集详情和对话
- `submit_part3_answer` - 提交答案并获取评分
- `auto_generate_part3_set` - 自动生成新的练习集
- `get_part3_prompt` - 获取 Part 3 AI 提示词

### Part 4 工具
- `list_part4_sets` - 列出所有 Part 4 练习集
- `get_part4_set_details` - 获取练习集详情和独白
- `submit_part4_answer` - 提交答案并获取评分
- `auto_generate_part4_set` - 自动生成新的练习集
- `get_part4_prompt` - 获取 Part 4 AI 提示词

### 用户信息工具
- `get_user_info` - 获取当前用户信息

## 使用示例

### 示例 1: 获取练习集列表

在 Claude 中：
```
请列出所有 Part 2 练习集
```

### 示例 2: 开始练习

```
请获取 ID 为 1 的 Part 2 练习集详情，我想开始练习
```

### 示例 3: 提交答案

```
我选择答案 A，请帮我提交 Part 2 练习集 1 的第 1 题答案
```

### 示例 4: 使用 AI 提示词

```
请给我 Part 2 的 AI 练习提示词，我想让 AI 帮我练习听力
```

## 常见问题

### Q: 连接失败怎么办？

A: 检查以下几点：
1. 服务器是否正在运行（查看终端输出）
2. 端口 6275 是否被占用
3. Token 是否正确
4. 配置文件格式是否正确

### Q: 认证失败怎么办？

A: 
1. 确认数据库中的 Token 与配置文件中的一致
2. 检查用户状态是否为 active
3. 查看服务器日志中的错误信息

### Q: 工具列表为空？

A:
1. 确认认证成功
2. 检查服务器日志，确认工具注册成功
3. 尝试重启客户端

### Q: 如何查看详细日志？

A:
```bash
# 设置日志级别为 debug
LOG_LEVEL=debug go run ./cmd/toeic-mcp/main.go -c config/dev.yaml
```

## 下一步

- 📖 阅读[完整配置教程](./MCP_CONFIGURATION_GUIDE.md)了解更多配置选项
- 🔧 查看[故障排查指南](./MCP_CONFIGURATION_GUIDE.md#故障排查)解决问题
- 📝 查看 [Part 2 提示词模板](./part2_prompt.md)了解如何使用 AI 练习
- 📝 查看 [Part 3 提示词模板](./part3_prompt.md)了解如何使用 AI 练习
- 📝 查看 [Part 4 提示词模板](./part4_prompt.md)了解如何使用 AI 练习

## 获取帮助

- GitHub Issues: https://github.com/alvinzane/doubao-toeic/issues
- Email: awary@qq.com

---

**祝您学习愉快！** 🎉

