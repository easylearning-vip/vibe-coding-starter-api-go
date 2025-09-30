# TOEIC MCP Server - Quick Start Guide

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- MySQL/PostgreSQL/SQLite database
- Valid configuration file

### Start the Server

```bash
# From project root
cd workspace/vocabulary/vibe-coding-starter-api-go

# Using default config (configs/config.yaml)
go run ./cmd/toeic-mcp/

# Using custom config
go run ./cmd/toeic-mcp/ -c configs/config-k3d.yaml

# Custom host and port
go run ./cmd/toeic-mcp/ -host 127.0.0.1 -port 8080
```

### Build and Run

```bash
# Build
go build -o toeic-mcp ./cmd/toeic-mcp/

# Run
./toeic-mcp -c configs/config.yaml
```

## 🔧 Configuration

### Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-c` | `configs/config.yaml` | Configuration file path |
| `-host` | `0.0.0.0` | HTTP host to listen on |
| `-port` | `6275` | HTTP port to listen on |

### Environment Variables

```bash
# Override config file location
export CONFIG_FILE=configs/config-k3d.yaml
go run ./cmd/toeic-mcp/
```

## 🔐 Authentication

### Get Your Token

1. Login to the system
2. Get your user token from the API
3. Use it in MCP client configuration

### Token Formats Supported

```bash
# Header: X-User-Token
X-User-Token: your-token-here

# Header: Authorization with Token
Authorization: Token your-token-here

# Header: Authorization with Bearer
Authorization: Bearer your-token-here
```

## 🛠️ Available Tools

### 1. get_user_info

Get current user information.

**Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "get_user_info"
  }
}
```

**Response:**
```
user_id=123
name=John Doe
```

### 2. list_part2_sets

List Part2 practice sets with pagination.

**Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "list_part2_sets",
    "arguments": {
      "page": 1,
      "page_size": 20
    }
  }
}
```

**Response:**
```
total=50 page=1 size=20
#1 total=30 done=15 correct=12 acc=80.00%
#2 total=25 done=25 correct=20 acc=80.00%
```

## 🔌 Client Configuration

### Augment Code

Add to your Augment configuration:

```json
{
  "mcpServers": {
    "toeic-mcp": {
      "url": "http://localhost:6275",
      "headers": {
        "X-User-Token": "your-token-here"
      }
    }
  }
}
```

### Cursor

Add to your Cursor MCP settings:

```json
{
  "toeic-mcp": {
    "command": "http",
    "args": ["http://localhost:6275"],
    "env": {
      "X_USER_TOKEN": "your-token-here"
    }
  }
}
```

## 📊 Health Check

```bash
# Check if server is running
curl http://localhost:6275/health

# Test with token
curl -H "X-User-Token: your-token" http://localhost:6275
```

## 🐛 Troubleshooting

### Server Won't Start

**Problem:** Port already in use
```bash
# Solution: Use different port
go run ./cmd/toeic-mcp/ -port 6276
```

**Problem:** Config file not found
```bash
# Solution: Specify full path
go run ./cmd/toeic-mcp/ -c /full/path/to/config.yaml
```

### Authentication Fails

**Problem:** Invalid token
```bash
# Check token format
curl -v -H "X-User-Token: your-token" http://localhost:6275
```

**Problem:** User not active
```sql
-- Check user status in database
SELECT id, username, status FROM users WHERE token = 'your-token';
```

### Tool Calls Fail

**Problem:** Missing parameters
```json
// Ensure all required parameters are provided
{
  "name": "list_part2_sets",
  "arguments": {
    "page": 1,
    "page_size": 20
  }
}
```

## 📝 Development

### Project Structure

```
cmd/toeic-mcp/
├── main.go              # Application entry point
├── README.md            # Full documentation
├── REFACTORING.md       # Refactoring details
└── QUICK_START.md       # This file

internal/
├── handler/mcp/         # MCP handlers
│   ├── server.go        # Main handler
│   ├── auth.go          # Authentication
│   ├── user_info.go     # User info tool
│   └── part2_practice_set.go  # Part2 tool
└── model/mcp/           # MCP models
    └── params.go        # Parameter definitions
```

### Adding a New Tool

1. **Define parameters** in `internal/model/mcp/params.go`
2. **Create handler** in `internal/handler/mcp/new_tool.go`
3. **Register tool** in `internal/handler/mcp/server.go`
4. **Test** your implementation

Example:
```go
// 1. Define params
type MyToolParams struct {
    Query string `json:"query"`
}

// 2. Create handler
func MyToolTool() *mcp.Tool {
    return &mcp.Tool{
        Name: "my_tool",
        Description: "My tool description",
    }
}

func HandleMyTool(...) func(...) {
    return func(...) {
        // Implementation
    }
}

// 3. Register in server.go
mcp.AddTool(server, MyToolTool(), HandleMyTool(...))
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/handler/mcp/...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Lint code
go vet ./...

# Static analysis
golangci-lint run
```

## 📚 Resources

- [Full Documentation](./README.md)
- [Refactoring Details](./REFACTORING.md)
- [MCP Specification](https://modelcontextprotocol.io/)
- [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk)

## 🆘 Getting Help

1. Check the [README](./README.md) for detailed documentation
2. Review [REFACTORING.md](./REFACTORING.md) for architecture details
3. Check server logs for error messages
4. Verify database connectivity
5. Ensure user token is valid and active

## 🎯 Common Use Cases

### Development

```bash
# Start with local config
go run ./cmd/toeic-mcp/ -c configs/config.yaml

# Watch logs
tail -f logs/toeic-mcp.log
```

### Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o toeic-mcp ./cmd/toeic-mcp/

# Run with production config
./toeic-mcp -c /etc/toeic-mcp/config.yaml
```

### Docker

```bash
# Build image
docker build -t toeic-mcp -f cmd/toeic-mcp/Dockerfile .

# Run container
docker run -p 6275:6275 \
  -v /path/to/config.yaml:/config.yaml \
  toeic-mcp -c /config.yaml
```

## ✅ Checklist

Before deploying:

- [ ] Configuration file is correct
- [ ] Database is accessible
- [ ] Port is available
- [ ] User tokens are valid
- [ ] Logs directory exists
- [ ] Tests pass
- [ ] Code is formatted
- [ ] No lint errors

## 🔄 Updates

Check for updates regularly:

```bash
# Update dependencies
go get -u ./...
go mod tidy

# Rebuild
go build ./cmd/toeic-mcp/
```

---

**Need more help?** Check the [full documentation](./README.md) or review the [refactoring guide](./REFACTORING.md).

