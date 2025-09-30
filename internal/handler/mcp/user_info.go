package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"vibe-coding-starter/internal/model"
	mcpModel "vibe-coding-starter/internal/model/mcp"
)

// GetUserInfoTool creates the MCP tool definition for getting user information
func GetUserInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_user_info",
		Description: "Get current user information (ID and name)",
	}
}

// HandleGetUserInfo handles the request to get user information
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
