package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"vibe-coding-starter/internal/model"
	mcpModel "vibe-coding-starter/internal/model/mcp"
	"vibe-coding-starter/internal/service"
)

// ListPart2SetsTool 创建列出Part2练习集的MCP工具定义
func ListPart2SetsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_part2_sets",
		Description: "查看个人 Part2 练习集列表（支持分页）",
	}
}

// HandleListPart2Sets 处理列出Part2练习集的请求
func HandleListPart2Sets(user *model.User, p2Service service.Part2PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part2ListParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part2ListParams) (*mcp.CallToolResult, any, error) {
		// 设置默认分页参数
		page, size := 1, 20
		if params != nil {
			if params.Page > 0 {
				page = params.Page
			}
			if params.PageSize > 0 {
				size = params.PageSize
			}
		}
		// 限制最大页面大小
		if size > 100 {
			size = 100
		}

		// 调用服务层获取练习集列表
		sets, total, err := p2Service.ListSets(ctx, user.ID, &service.ListPracticeSetOptions{
			Page:     page,
			PageSize: size,
			Sort:     "id",
			Order:    "desc",
		})
		if err != nil {
			return nil, nil, fmt.Errorf("查询失败: %v", err)
		}

		// 处理空结果
		if len(sets) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "暂无练习集"}},
			}, nil, nil
		}

		// 格式化输出结果
		var b strings.Builder
		fmt.Fprintf(&b, "total=%d page=%d size=%d\n", total, page, size)
		for _, s := range sets {
			fmt.Fprintf(&b, "#%d total=%d done=%d correct=%d acc=%.2f%%\n",
				s.ID, s.TotalQuestions, s.CompletedCount, s.CorrectCount, s.Accuracy*100)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
		}, nil, nil
	}
}

// AutoGeneratePart2SetTool 创建自动生成Part2练习集的MCP工具定义
func AutoGeneratePart2SetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "auto_generate_part2_set",
		Description: "自动生成 Part2 练习集",
	}
}

// HandleAutoGeneratePart2Set 处理自动生成Part2练习集的请求
func HandleAutoGeneratePart2Set(user *model.User, p2Service service.Part2PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.AutoGeneratePart2SetParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.AutoGeneratePart2SetParams) (*mcp.CallToolResult, any, error) {
		// 验证参数
		if params == nil || params.Count == 0 {
			return nil, nil, fmt.Errorf("参数错误: count 是必需的")
		}

		if params.Count < 10 || params.Count > 30 {
			return nil, nil, fmt.Errorf("参数错误: count 必须在 10-30 之间")
		}

		// 构建请求
		genReq := &service.GeneratePracticeSetRequest{
			TotalQuestions:    int32(params.Count),
			ScenarioId:        params.ScenarioID,
			DifficultyLevelId: params.DifficultyLevelID,
			Mode:              "sequential",
		}

		// 调用服务层生成练习集
		set, items, err := p2Service.AutoGenerateSet(ctx, user.ID, genReq)
		if err != nil {
			return nil, nil, fmt.Errorf("生成失败: %v", err)
		}

		// 格式化输出结果
		var b strings.Builder
		fmt.Fprintf(&b, "✅ 练习集生成成功\n\n")
		fmt.Fprintf(&b, "练习集ID: %d\n", set.ID)
		fmt.Fprintf(&b, "题目总数: %d\n", set.TotalQuestions)
		if set.ScenarioId.Valid {
			fmt.Fprintf(&b, "场景ID: %d\n", set.ScenarioId.Int32)
		}
		if set.DifficultyLevelId.Valid {
			fmt.Fprintf(&b, "难度级别ID: %d\n", set.DifficultyLevelId.Int32)
		}
		fmt.Fprintf(&b, "\n已添加 %d 道题目到练习集\n", len(items))
		fmt.Fprintf(&b, "\n使用 get_part2_set_details 查看题目详情")

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
		}, nil, nil
	}
}

// GetPart2SetDetailsTool 创建获取Part2练习集详情的MCP工具定义
func GetPart2SetDetailsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_part2_set_details",
		Description: "查看 Part2 练习集明细列表（包含题目详情）",
	}
}

// HandleGetPart2SetDetails 处理获取Part2练习集详情的请求
func HandleGetPart2SetDetails(user *model.User, p2Service service.Part2PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.GetPart2SetDetailsParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.GetPart2SetDetailsParams) (*mcp.CallToolResult, any, error) {
		// 验证参数
		if params == nil || params.SetID == 0 {
			return nil, nil, fmt.Errorf("参数错误: set_id 是必需的")
		}

		// 获取练习集基本信息
		set, err := p2Service.GetSetByID(ctx, user.ID, params.SetID)
		if err != nil {
			return nil, nil, fmt.Errorf("获取练习集失败: %v", err)
		}

		// 获取练习集题目列表
		questions, err := p2Service.ListSetQuestions(ctx, user.ID, params.SetID)
		if err != nil {
			return nil, nil, fmt.Errorf("获取题目列表失败: %v", err)
		}

		// 格式化输出结果
		var b strings.Builder
		fmt.Fprintf(&b, "📋 练习集详情\n\n")
		fmt.Fprintf(&b, "练习集ID: %d\n", set.ID)
		fmt.Fprintf(&b, "题目总数: %d\n", set.TotalQuestions)
		fmt.Fprintf(&b, "已完成: %d\n", set.CompletedCount)
		fmt.Fprintf(&b, "正确数: %d\n", set.CorrectCount)
		fmt.Fprintf(&b, "准确率: %.2f%%\n\n", set.Accuracy*100)

		if len(questions) == 0 {
			fmt.Fprintf(&b, "暂无题目")
		} else {
			fmt.Fprintf(&b, "题目列表:\n")
			fmt.Fprintf(&b, "─────────────────────────────────────\n\n")
			for i, q := range questions {
				fmt.Fprintf(&b, "题目 %d (Item ID: %d)\n", i+1, q.ItemID)
				fmt.Fprintf(&b, "问题: %s\n", q.Question.QuestionText)
				fmt.Fprintf(&b, "A. %s\n", q.Question.OptionA)
				fmt.Fprintf(&b, "B. %s\n", q.Question.OptionB)
				fmt.Fprintf(&b, "C. %s\n", q.Question.OptionC)
				if i < len(questions)-1 {
					fmt.Fprintf(&b, "\n─────────────────────────────────────\n\n")
				}
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
		}, nil, nil
	}
}

// SubmitPart2AnswerTool 创建提交Part2答案的MCP工具定义
func SubmitPart2AnswerTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "submit_part2_answer",
		Description: "提交 Part2 练习题答案",
	}
}

// HandleSubmitPart2Answer 处理提交Part2答案的请求
func HandleSubmitPart2Answer(user *model.User, p2Service service.Part2PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.SubmitPart2AnswerParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.SubmitPart2AnswerParams) (*mcp.CallToolResult, any, error) {
		// 验证参数
		if params == nil {
			return nil, nil, fmt.Errorf("参数错误: 缺少必需参数")
		}
		if params.SetID == 0 {
			return nil, nil, fmt.Errorf("参数错误: set_id 是必需的")
		}
		if params.ItemID == 0 {
			return nil, nil, fmt.Errorf("参数错误: item_id 是必需的")
		}
		if params.UserAnswer == "" {
			return nil, nil, fmt.Errorf("参数错误: user_answer 是必需的")
		}

		// 验证答案格式
		answer := strings.ToUpper(strings.TrimSpace(params.UserAnswer))
		if answer != "A" && answer != "B" && answer != "C" {
			return nil, nil, fmt.Errorf("参数错误: user_answer 必须是 A、B 或 C")
		}

		// 获取题目详情以验证答案
		itemDetail, err := p2Service.GetItemDetail(ctx, user.ID, params.ItemID)
		if err != nil {
			return nil, nil, fmt.Errorf("获取题目失败: %v", err)
		}

		// 检查题目是否属于指定的练习集
		if itemDetail.SetID != params.SetID {
			return nil, nil, fmt.Errorf("题目不属于指定的练习集")
		}

		// 判断答案是否正确
		var correctAnswer string
		if itemDetail.Question.CorrectAnswer.Valid {
			correctAnswer = strings.ToUpper(strings.TrimSpace(itemDetail.Question.CorrectAnswer.String))
		}
		isCorrect := answer == correctAnswer

		// 提交答案
		submitReq := &service.SubmitPart2AnswerRequest{
			SelectedAnswer: answer,
			IsCorrect:      isCorrect,
		}
		if err := p2Service.SubmitAnswer(ctx, user.ID, params.ItemID, submitReq); err != nil {
			return nil, nil, fmt.Errorf("提交答案失败: %v", err)
		}

		// 格式化输出结果
		var b strings.Builder
		if isCorrect {
			fmt.Fprintf(&b, "✅ 回答正确！\n\n")
		} else {
			fmt.Fprintf(&b, "❌ 回答错误\n\n")
		}

		fmt.Fprintf(&b, "你的答案: %s\n", answer)
		if correctAnswer != "" {
			fmt.Fprintf(&b, "正确答案: %s\n", correctAnswer)
		} else {
			fmt.Fprintf(&b, "正确答案: (未设置)\n")
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
		}, nil, nil
	}
}
