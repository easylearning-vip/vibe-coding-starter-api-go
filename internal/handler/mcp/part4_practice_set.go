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

// ListPart4SetsTool creates the MCP tool definition for listing Part4 practice sets
func ListPart4SetsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_part4_sets",
		Description: "List personal Part4 practice sets (with pagination support)",
	}
}

// HandleListPart4Sets handles the request to list Part4 practice sets
func HandleListPart4Sets(user *model.User, p4Service service.Part4PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part4ListParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part4ListParams) (*mcp.CallToolResult, any, error) {
		// Set default pagination parameters
		page, size := 1, 20
		if params != nil {
			if params.Page > 0 {
				page = params.Page
			}
			if params.PageSize > 0 {
				size = params.PageSize
			}
		}
		// Limit maximum page size
		if size > 100 {
			size = 100
		}

		// Call service layer to get practice set list
		sets, total, err := p4Service.ListSets(ctx, user.ID, &service.ListPracticeSetOptions{
			Page:     page,
			PageSize: size,
			Sort:     "id",
			Order:    "desc",
		})
		if err != nil {
			return nil, nil, fmt.Errorf("query failed: %v", err)
		}

		// Handle empty results
		if len(sets) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "No practice sets available"}},
			}, nil, nil
		}

		// Format output results
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

// AutoGeneratePart4SetTool creates the MCP tool definition for auto-generating Part4 practice sets
func AutoGeneratePart4SetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "auto_generate_part4_set",
		Description: "Auto-generate Part4 practice set",
	}
}

// HandleAutoGeneratePart4Set handles the request to auto-generate Part4 practice sets
func HandleAutoGeneratePart4Set(user *model.User, p4Service service.Part4PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.AutoGeneratePart4SetParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.AutoGeneratePart4SetParams) (*mcp.CallToolResult, any, error) {
		// Validate parameters
		if params == nil || params.Count == 0 {
			return nil, nil, fmt.Errorf("parameter error: count is required")
		}

		if params.Count < 10 || params.Count > 30 {
			return nil, nil, fmt.Errorf("parameter error: count must be between 10 and 30")
		}

		// Build request
		genReq := &service.GeneratePracticeSetRequest{
			TotalQuestions:    int32(params.Count),
			ScenarioId:        params.ScenarioID,
			DifficultyLevelId: params.DifficultyLevelID,
			Mode:              "sequential",
		}

		// Call service layer to generate practice set
		set, items, err := p4Service.AutoGenerateSet(ctx, user.ID, genReq)
		if err != nil {
			return nil, nil, fmt.Errorf("generation failed: %v", err)
		}

		// Format output results
		var b strings.Builder
		fmt.Fprintf(&b, "Practice set generated successfully\n\n")
		fmt.Fprintf(&b, "Set ID: %d\n", set.ID)
		fmt.Fprintf(&b, "Total Questions: %d\n", set.TotalQuestions)
		if set.ScenarioId.Valid {
			fmt.Fprintf(&b, "Scenario ID: %d\n", set.ScenarioId.Int32)
		}
		if set.DifficultyLevelId.Valid {
			fmt.Fprintf(&b, "Difficulty Level ID: %d\n", set.DifficultyLevelId.Int32)
		}
		fmt.Fprintf(&b, "\nAdded %d questions to practice set\n", len(items))
		fmt.Fprintf(&b, "\nUse get_part4_set_details to view question details")

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
		}, nil, nil
	}
}

// GetPart4SetDetailsTool creates the MCP tool definition for getting Part4 practice set details
func GetPart4SetDetailsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_part4_set_details",
		Description: "View Part4 practice set details (including question details)",
	}
}

// HandleGetPart4SetDetails handles the request to get Part4 practice set details
func HandleGetPart4SetDetails(user *model.User, p4Service service.Part4PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.GetPart4SetDetailsParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.GetPart4SetDetailsParams) (*mcp.CallToolResult, any, error) {
		// Validate parameters
		if params == nil || params.SetID == 0 {
			return nil, nil, fmt.Errorf("parameter error: set_id is required")
		}

		// Get practice set basic information
		set, err := p4Service.GetSetByID(ctx, user.ID, params.SetID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get practice set: %v", err)
		}

		// Get practice set question list
		questions, err := p4Service.ListSetQuestions(ctx, user.ID, params.SetID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get question list: %v", err)
		}

		// Format output results
		var b strings.Builder
		fmt.Fprintf(&b, "Practice Set Details\n\n")
		fmt.Fprintf(&b, "Set ID: %d\n", set.ID)
		fmt.Fprintf(&b, "Total Questions: %d\n", set.TotalQuestions)
		fmt.Fprintf(&b, "Completed: %d\n", set.CompletedCount)
		fmt.Fprintf(&b, "Correct: %d\n", set.CorrectCount)
		fmt.Fprintf(&b, "Accuracy: %.2f%%\n\n", set.Accuracy*100)

		if len(questions) == 0 {
			fmt.Fprintf(&b, "No questions available")
		} else {
			fmt.Fprintf(&b, "Question List:\n\n")
			for i, q := range questions {
				// Display talk content
				fmt.Fprintf(&b, "(Item ID: %d) Talk %d:\n", q.ItemID, i+1)
				if q.Talk.Title.Valid && q.Talk.Title.String != "" {
					fmt.Fprintf(&b, "Title: %s\n", q.Talk.Title.String)
				}
				fmt.Fprintf(&b, "Content: %s\n\n", q.Talk.Content)

				// Display the specific question for this item
				var questionText string
				switch q.QuestionIndex {
				case 1:
					questionText = q.Talk.Question1
				case 2:
					questionText = q.Talk.Question2
				case 3:
					questionText = q.Talk.Question3
				}
				fmt.Fprintf(&b, "Question %d: %s\n", q.QuestionIndex, questionText)

				// Find and display answer options for this question
				for _, opt := range q.Talk.AnswerOptions {
					if opt.QuestionNumber == q.QuestionIndex {
						fmt.Fprintf(&b, "A. %s\n", opt.OptionA)
						fmt.Fprintf(&b, "B. %s\n", opt.OptionB)
						fmt.Fprintf(&b, "C. %s\n", opt.OptionC)
						fmt.Fprintf(&b, "D. %s\n", opt.OptionD)
						break
					}
				}

				if i < len(questions)-1 {
					fmt.Fprintf(&b, "\n")
				}
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
		}, nil, nil
	}
}

// SubmitPart4AnswerTool creates the MCP tool definition for submitting Part4 answers
func SubmitPart4AnswerTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "submit_part4_answer",
		Description: "Submit Part4 practice answer. Optionally provide is_correct parameter for AI assistant to judge answer correctness instead of comparing with database.",
	}
}

// HandleSubmitPart4Answer handles the request to submit Part4 answers
func HandleSubmitPart4Answer(user *model.User, p4Service service.Part4PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.SubmitPart4AnswerParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.SubmitPart4AnswerParams) (*mcp.CallToolResult, any, error) {
		// Validate parameters
		if params == nil {
			return nil, nil, fmt.Errorf("parameter error: missing required parameters")
		}
		if params.SetID == 0 {
			return nil, nil, fmt.Errorf("parameter error: set_id is required")
		}
		if params.ItemID == 0 {
			return nil, nil, fmt.Errorf("parameter error: item_id is required")
		}
		if params.UserAnswer == "" {
			return nil, nil, fmt.Errorf("parameter error: user_answer is required")
		}

		// Validate answer format
		answer := strings.ToUpper(strings.TrimSpace(params.UserAnswer))
		if answer != "A" && answer != "B" && answer != "C" && answer != "D" {
			return nil, nil, fmt.Errorf("parameter error: user_answer must be A, B, C, or D")
		}

		// Get question details to verify answer
		itemDetail, err := p4Service.GetItemDetail(ctx, user.ID, params.ItemID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get question: %v", err)
		}

		// Check if question belongs to the specified practice set
		if itemDetail.SetID != params.SetID {
			return nil, nil, fmt.Errorf("question does not belong to the specified practice set")
		}

		// Determine if answer is correct
		var correctAnswer string
		var isCorrect bool

		if params.IsCorrect != nil {
			// Use AI assistant's judgment if provided
			isCorrect = *params.IsCorrect
		} else {
			// Compare with correct answer in database
			// Find the answer option for this question
			for _, opt := range itemDetail.Talk.AnswerOptions {
				if opt.QuestionNumber == itemDetail.QuestionIndex && opt.CorrectAnswer.Valid {
					correctAnswer = strings.ToUpper(strings.TrimSpace(opt.CorrectAnswer.String))
					break
				}
			}
			isCorrect = answer == correctAnswer
		}

		// Get correct answer for display (if available)
		if correctAnswer == "" {
			for _, opt := range itemDetail.Talk.AnswerOptions {
				if opt.QuestionNumber == itemDetail.QuestionIndex && opt.CorrectAnswer.Valid {
					correctAnswer = strings.ToUpper(strings.TrimSpace(opt.CorrectAnswer.String))
					break
				}
			}
		}

		// Submit answer
		submitReq := &service.SubmitPart4AnswerRequest{
			SelectedAnswer: answer,
			IsCorrect:      isCorrect,
		}
		if err := p4Service.SubmitAnswer(ctx, user.ID, params.ItemID, submitReq); err != nil {
			return nil, nil, fmt.Errorf("failed to submit answer: %v", err)
		}

		// Format output results
		var b strings.Builder
		if isCorrect {
			fmt.Fprintf(&b, "Correct answer!\n\n")
		} else {
			fmt.Fprintf(&b, "Incorrect answer\n\n")
		}

		fmt.Fprintf(&b, "Your answer: %s\n", answer)
		if correctAnswer != "" {
			fmt.Fprintf(&b, "Correct answer: %s\n", correctAnswer)
		} else {
			fmt.Fprintf(&b, "Correct answer: (not set)\n")
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
		}, nil, nil
	}
}

