# TOEIC MCP 服务器配置完整教程

本教程详细介绍如何配置和使用 TOEIC MCP 服务器，包括认证配置、多种客户端支持和常见问题解决。

## 目录

- [快速开始](#快速开始)
- [认证配置](#认证配置)
- [客户端配置](#客户端配置)
  - [Claude Desktop](#claude-desktop)
  - [Cursor](#cursor)
  - [VS Code](#vs-code)
  - [Cline](#cline)
- [高级配置](#高级配置)
- [故障排查](#故障排查)

---

## 快速开始

### 1. 启动 MCP 服务器

```bash
# 使用配置文件启动
go run ./cmd/toeic-mcp/main.go -c config/dev.yaml

# 或使用环境变量
export CONFIG_FILE=/path/to/config.yaml
go run ./cmd/toeic-mcp/main.go

# 自定义主机和端口
go run ./cmd/toeic-mcp/main.go -c config/dev.yaml --host 0.0.0.0 --port 6275
```

### 2. 获取用户 Token

TOEIC MCP 使用基于 Token 的认证。您需要从数据库中获取或生成用户 Token：

```sql
-- 查询用户 Token
SELECT id, username, email, token FROM users WHERE username = 'your_username';

-- 如果没有 Token，需要生成并更新
UPDATE users SET token = 'your_generated_token' WHERE id = 1;
```

### 3. 验证服务器运行

```bash
# 测试服务器连接（无需认证）
curl http://localhost:6275/

# 测试带认证的请求
curl -H "X-User-Token: YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -X POST http://localhost:6275/ \
     -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}'
```

---

## 认证配置

TOEIC MCP 服务器支持三种 Token 传递方式，按优先级排序：

### 方式 1: X-User-Token Header（推荐）

这是最简单直接的方式，推荐用于所有客户端：

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "your_token_here"
      }
    }
  }
}
```

### 方式 2: Authorization Token

使用标准的 Authorization header，格式为 `Token <token>`：

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "Authorization": "Token your_token_here"
      }
    }
  }
}
```

### 方式 3: Authorization Bearer

使用 Bearer token 格式：

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "Authorization": "Bearer your_token_here"
      }
    }
  }
}
```

### 使用环境变量（推荐用于生产环境）

为了安全起见，不要在配置文件中硬编码 Token，而是使用环境变量：

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

然后在环境中设置：

```bash
# Linux/macOS
export TOEIC_USER_TOKEN="your_token_here"

# Windows PowerShell
$env:TOEIC_USER_TOKEN="your_token_here"

# Windows CMD
set TOEIC_USER_TOKEN=your_token_here
```

---

## 客户端配置

### Claude Desktop

Claude Desktop 使用 `claude_desktop_config.json` 配置文件。

#### 配置文件位置

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

#### 基本配置

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

#### 完整配置示例

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "${env:TOEIC_USER_TOKEN}",
        "Content-Type": "application/json"
      },
      "note": "TOEIC Practice Server - Provides Part 2/3/4 practice tools"
    }
  }
}
```

#### 重启 Claude Desktop

配置完成后，需要完全退出并重启 Claude Desktop：

1. 退出 Claude Desktop（不是最小化）
2. 重新启动应用
3. 在对话中测试：`请列出所有可用的 TOEIC 练习集`

---

### Cursor

Cursor 使用 `.cursor/mcp.json` 配置文件。

#### 配置文件位置

在项目根目录创建 `.cursor/mcp.json`：

```bash
mkdir -p .cursor
touch .cursor/mcp.json
```

#### 配置示例

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

#### 使用方法

1. 保存配置文件
2. 重启 Cursor
3. 在 Composer 或 Chat 中使用：`@toeic-mcp 获取 Part 2 练习集列表`

---

### VS Code

VS Code 使用 `settings.json` 配置 MCP 服务器。

#### 配置文件位置

- **用户设置**: `~/.config/Code/User/settings.json` (Linux/macOS)
- **工作区设置**: `.vscode/settings.json`

#### 配置示例

```json
{
  "github.copilot.chat.mcp.servers": {
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

#### 使用方法

1. 打开命令面板 (`Ctrl+Shift+P` 或 `Cmd+Shift+P`)
2. 输入 "Reload Window" 重新加载
3. 在 GitHub Copilot Chat 中使用 MCP 工具

---

### Cline

Cline 使用 `cline_mcp_settings.json` 配置文件。

#### 配置文件位置

- **VS Code Extension**: `~/.vscode/extensions/saoudrizwan.claude-dev-*/cline_mcp_settings.json`

#### 配置示例

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "Authorization": "Bearer ${env:TOEIC_USER_TOKEN}"
      },
      "alwaysAllow": [
        "list_part2_sets",
        "get_part2_set_details",
        "list_part3_sets",
        "get_part3_set_details",
        "list_part4_sets",
        "get_part4_set_details"
      ]
    }
  }
}
```

---

## 高级配置

### 多用户配置

如果需要为不同用户配置不同的 Token：

```json
{
  "mcpServers": {
    "toeic-mcp-user1": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "${env:TOEIC_USER1_TOKEN}"
      }
    },
    "toeic-mcp-user2": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "${env:TOEIC_USER2_TOKEN}"
      }
    }
  }
}
```

### 远程服务器配置

如果 MCP 服务器部署在远程主机：

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "https://your-domain.com/mcp",
      "headers": {
        "X-User-Token": "${env:TOEIC_USER_TOKEN}"
      }
    }
  }
}
```

**安全建议**：
- 使用 HTTPS 而不是 HTTP
- 配置防火墙规则限制访问
- 定期轮换 Token
- 使用反向代理（如 Nginx）添加额外的安全层

### 代理配置

如果需要通过代理访问：

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "type": "streamable-http",
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "${env:TOEIC_USER_TOKEN}"
      },
      "proxy": "http://proxy.example.com:8080"
    }
  }
}
```

---

## 故障排查

### 问题 1: 连接失败

**症状**: 客户端无法连接到 MCP 服务器

**解决方案**:

```bash
# 1. 检查服务器是否运行
curl http://localhost:6275/

# 2. 检查端口是否被占用
lsof -i :6275  # Linux/macOS
netstat -ano | findstr :6275  # Windows

# 3. 检查防火墙设置
sudo ufw status  # Linux
```

### 问题 2: 认证失败

**症状**: 服务器返回 401 Unauthorized 或认证错误

**解决方案**:

```bash
# 1. 验证 Token 是否正确
curl -H "X-User-Token: YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -X POST http://localhost:6275/ \
     -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}'

# 2. 检查数据库中的 Token
# 连接到数据库并查询
SELECT id, username, token, status FROM users WHERE token = 'YOUR_TOKEN';

# 3. 检查用户状态是否为 active
```

### 问题 3: 环境变量未生效

**症状**: 使用 `${env:VARIABLE}` 但 Token 未被替换

**解决方案**:

```bash
# 1. 确认环境变量已设置
echo $TOEIC_USER_TOKEN  # Linux/macOS
echo %TOEIC_USER_TOKEN%  # Windows CMD
echo $env:TOEIC_USER_TOKEN  # Windows PowerShell

# 2. 在启动客户端前设置环境变量
export TOEIC_USER_TOKEN="your_token"
# 然后启动 Claude Desktop 或其他客户端

# 3. 对于 macOS，可能需要在 launchd 中设置
launchctl setenv TOEIC_USER_TOKEN "your_token"
```

### 问题 4: 工具不可用

**症状**: MCP 连接成功但工具列表为空

**解决方案**:

1. 检查服务器日志，确认工具注册成功
2. 发送 `tools/list` 请求测试：

```bash
curl -H "X-User-Token: YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -X POST http://localhost:6275/ \
     -d '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
```

### 问题 5: 客户端配置未生效

**症状**: 修改配置文件后没有变化

**解决方案**:

1. **Claude Desktop**: 完全退出应用（检查系统托盘），然后重启
2. **Cursor**: 重启 Cursor 应用
3. **VS Code**: 运行 "Reload Window" 命令
4. **Cline**: 重新加载 VS Code 窗口

### 调试技巧

#### 启用详细日志

修改服务器启动命令，增加日志级别：

```bash
# 设置日志级别为 debug
LOG_LEVEL=debug go run ./cmd/toeic-mcp/main.go -c config/dev.yaml
```

#### 使用 MCP Inspector

MCP Inspector 是官方的调试工具：

```bash
# 安装
npm install -g @modelcontextprotocol/inspector

# 使用
mcp-inspector http://localhost:6275 \
  --header "X-User-Token: YOUR_TOKEN"
```

#### 查看服务器日志

服务器会输出详细的认证和请求日志：

```
INFO  Starting TOEIC MCP server (streamable HTTP) addr=0.0.0.0:6275
WARN  missing token for MCP request
WARN  invalid user token for MCP err="record not found"
INFO  User authenticated user_id=1 username=testuser
```

---

## 安全最佳实践

1. **不要在代码仓库中提交 Token**
   - 使用 `.gitignore` 排除配置文件
   - 使用环境变量存储敏感信息

2. **定期轮换 Token**
   ```sql
   UPDATE users SET token = 'new_token' WHERE id = 1;
   ```

3. **使用 HTTPS**
   - 生产环境必须使用 HTTPS
   - 配置 SSL/TLS 证书

4. **限制访问**
   - 使用防火墙规则
   - 仅允许必要的 IP 地址访问

5. **监控和审计**
   - 记录所有认证尝试
   - 监控异常访问模式
   - 定期审查访问日志

---

## 相关资源

- [Model Context Protocol 官方文档](https://modelcontextprotocol.io/)
- [MCP Transports 规范](https://modelcontextprotocol.io/docs/concepts/transports)
- [TOEIC MCP 项目主页](https://github.com/alvinzane/doubao-toeic)
- [Go SDK 文档](https://github.com/modelcontextprotocol/go-sdk)

---

## 获取帮助

如果遇到问题：

1. 查看本文档的[故障排查](#故障排查)部分
2. 检查服务器日志输出
3. 在 GitHub 提交 Issue: https://github.com/alvinzane/doubao-toeic/issues
4. 联系维护者: awary@qq.com

---

**最后更新**: 2025-09-30
**版本**: 1.0.0

