package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"vibe-coding-starter/internal/model"
	mcpModel "vibe-coding-starter/internal/model/mcp"
)

// GetUserInfoTool 创建获取用户信息的MCP工具定义
func GetUserInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_user_info",
		Description: "查看当前用户信息（id与名称）",
	}
}

// HandleGetUserInfo 处理获取用户信息的请求
func HandleGetUserInfo(user *model.User) func(ctx context.Context, req *mcp.CallToolRequest, _ *mcpModel.GetUserInfoParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, _ *mcpModel.GetUserInfoParams) (*mcp.CallToolResult, any, error) {
		name := user.Nickname
		if strings.TrimSpace(name) == "" {
			name = user.Username
		}
		text := fmt.Sprintf("user_id=%d\nname=%s", user.ID, name)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: text}},
		}, nil, nil
	}
}

