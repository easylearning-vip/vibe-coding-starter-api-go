package mcp

import (
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// MCPHandler MCP服务处理器
type MCPHandler struct {
	logger        logger.Logger
	userRepo      repository.UserRepository
	p2Service     service.Part2PracticeSetService
	promptService service.ToeicAiPromptService
}

// NewMCPHandler 创建MCP处理器
func NewMCPHandler(
	logger logger.Logger,
	userRepo repository.UserRepository,
	p2Service service.Part2PracticeSetService,
	promptService service.ToeicAiPromptService,
) *MCPHandler {
	return &MCPHandler{
		logger:        logger,
		userRepo:      userRepo,
		p2Service:     p2Service,
		promptService: promptService,
	}
}

// BuildHTTPHandler 构建MCP Streamable HTTP处理器
// 实现基于每个请求的token认证
func (h *MCPHandler) BuildHTTPHandler() http.Handler {
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		// 提取用户token
		token := ExtractToken(r)
		if token == "" {
			h.logger.Warn("missing token for MCP request")
			return nil
		}

		// 验证token并获取用户信息
		ctx := r.Context()
		user, err := h.userRepo.GetByToken(ctx, token)
		if err != nil || user == nil || !user.IsActive() {
			h.logger.Warn("invalid user token for MCP", "err", err)
			return nil
		}

		// 创建绑定到该用户的MCP服务器
		server := mcp.NewServer(&mcp.Implementation{
			Name:    "toeic-mcp",
			Version: "0.1.0",
		}, nil)

		// 注册工具
		h.registerTools(server, user)

		return server
	}, &mcp.StreamableHTTPOptions{Stateless: false})

	return handler
}

// registerTools 注册所有MCP工具
func (h *MCPHandler) registerTools(server *mcp.Server, user *model.User) {
	// 注册获取用户信息工具
	mcp.AddTool(server, GetUserInfoTool(), HandleGetUserInfo(user))

	// 注册Part2 AI提示词工具
	mcp.AddTool(server, Part2PromptTool(), HandlePart2Prompt(user, h.promptService))

	// 注册Part2练习集相关工具
	mcp.AddTool(server, ListPart2SetsTool(), HandleListPart2Sets(user, h.p2Service))
	mcp.AddTool(server, AutoGeneratePart2SetTool(), HandleAutoGeneratePart2Set(user, h.p2Service))
	mcp.AddTool(server, GetPart2SetDetailsTool(), HandleGetPart2SetDetails(user, h.p2Service))
	mcp.AddTool(server, SubmitPart2AnswerTool(), HandleSubmitPart2Answer(user, h.p2Service))
}
