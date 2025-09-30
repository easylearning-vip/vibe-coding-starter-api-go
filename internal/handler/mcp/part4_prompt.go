package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"vibe-coding-starter/internal/model"
	mcpModel "vibe-coding-starter/internal/model/mcp"
	"vibe-coding-starter/internal/service"
)

// Part4PromptTool creates the MCP tool definition for getting Part4 AI prompt
func Part4PromptTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "part4_prompt",
		Description: "Get the latest Part4 AI assistant prompt for using TOEIC MCP tools. This prompt provides instructions on how to use the MCP tools effectively.",
	}
}

// HandlePart4Prompt handles the request to get Part4 AI prompt
func HandlePart4Prompt(user *model.User, promptService service.ToeicAiPromptService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part4PromptParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part4PromptParams) (*mcp.CallToolResult, any, error) {
		// Get the latest prompt (ordered by ID DESC, get the first one)
		prompts, total, err := promptService.List(ctx, &service.ListToeicAiPromptOptions{
			Page:     1,
			PageSize: 1,
			Sort:     "id",
			Order:    "desc",
		})
		if err != nil || total == 0 {
			return nil, nil, fmt.Errorf("failed to get Part4 prompt: no prompts available")
		}

		// Return the Part4 prompt from the latest record
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: prompts[0].Part4Prompt}},
		}, nil, nil
	}
}

