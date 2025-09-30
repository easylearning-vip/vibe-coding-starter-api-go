package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"vibe-coding-starter/internal/model"
	mcpModel "vibe-coding-starter/internal/model/mcp"
	"vibe-coding-starter/internal/service"
)

// Part2PromptTool creates the MCP tool definition for getting Part2 AI prompt
func Part2PromptTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "part2_prompt",
		Description: "Get the latest Part2 AI assistant prompt for using TOEIC MCP tools. This prompt provides instructions on how to use the MCP tools effectively.",
	}
}

// HandlePart2Prompt handles the request to get Part2 AI prompt
func HandlePart2Prompt(user *model.User, promptService service.ToeicAiPromptService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part2PromptParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part2PromptParams) (*mcp.CallToolResult, any, error) {
		// Get the latest prompt (ordered by ID DESC, get the first one)
		prompts, total, err := promptService.List(ctx, &service.ListToeicAiPromptOptions{
			Page:     1,
			PageSize: 1,
			Sort:     "id",
			Order:    "desc",
		})
		if err != nil || total == 0 {
			return nil, nil, fmt.Errorf("failed to get Part2 prompt: no prompts available")
		}

		// Return the Part2 prompt from the latest record
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: prompts[0].Part2Prompt}},
		}, nil, nil
	}
}
