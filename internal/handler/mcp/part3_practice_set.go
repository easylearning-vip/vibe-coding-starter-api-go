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

// ListPart3SetsTool creates the MCP tool definition for listing Part3 practice sets
func ListPart3SetsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_part3_sets",
		Description: "List personal Part3 practice sets (with pagination support)",
	}
}

// HandleListPart3Sets handles the request to list Part3 practice sets
func HandleListPart3Sets(user *model.User, p3Service service.Part3PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part3ListParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.Part3ListParams) (*mcp.CallToolResult, any, error) {
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
		sets, total, err := p3Service.ListSets(ctx, user.ID, &service.ListPracticeSetOptions{
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

// AutoGeneratePart3SetTool creates the MCP tool definition for auto-generating Part3 practice sets
func AutoGeneratePart3SetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "auto_generate_part3_set",
		Description: "Auto-generate Part3 practice set",
	}
}

// HandleAutoGeneratePart3Set handles the request to auto-generate Part3 practice sets
func HandleAutoGeneratePart3Set(user *model.User, p3Service service.Part3PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.AutoGeneratePart3SetParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.AutoGeneratePart3SetParams) (*mcp.CallToolResult, any, error) {
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
		set, items, err := p3Service.AutoGenerateSet(ctx, user.ID, genReq)
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
		fmt.Fprintf(&b, "\nUse get_part3_set_details to view conversation and question details")

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
		}, nil, nil
	}
}

// GetPart3SetDetailsTool creates the MCP tool definition for getting Part3 practice set details
func GetPart3SetDetailsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_part3_set_details",
		Description: "View Part3 practice set details (including conversation and question details)",
	}
}

// HandleGetPart3SetDetails handles the request to get Part3 practice set details
func HandleGetPart3SetDetails(user *model.User, p3Service service.Part3PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.GetPart3SetDetailsParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.GetPart3SetDetailsParams) (*mcp.CallToolResult, any, error) {
		// Validate parameters
		if params == nil || params.SetID == 0 {
			return nil, nil, fmt.Errorf("parameter error: set_id is required")
		}

		// Get practice set basic information
		set, err := p3Service.GetSetByID(ctx, user.ID, params.SetID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get practice set: %v", err)
		}

		// Get practice set question list
		questions, err := p3Service.ListSetQuestions(ctx, user.ID, params.SetID)
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
				// Display conversation content (only once per conversation)
				if i == 0 || (i > 0 && questions[i-1].Conversation.ID != q.Conversation.ID) {
					fmt.Fprintf(&b, "=== Conversation %d ===\n", q.Conversation.ConversationNumber)
					if q.Conversation.Title.Valid && q.Conversation.Title.String != "" {
						fmt.Fprintf(&b, "Title: %s\n", q.Conversation.Title.String)
					}
					fmt.Fprintf(&b, "%s\n\n", q.Conversation.Content)
				}

				// Display question based on question index
				var questionText string
				var answerOptions *model.Part3AnswerOption
				switch q.QuestionIndex {
				case 1:
					questionText = q.Conversation.Question1
				case 2:
					questionText = q.Conversation.Question2
				case 3:
					questionText = q.Conversation.Question3
				}

				// Find the corresponding answer options
				for _, opt := range q.Conversation.AnswerOptions {
					if opt.QuestionNumber == q.QuestionIndex {
						answerOptions = &opt
						break
					}
				}

				fmt.Fprintf(&b, "(Item ID: %d) Question %d: %s\n", q.ItemID, i+1, questionText)
				if answerOptions != nil {
					fmt.Fprintf(&b, "A. %s\n", answerOptions.OptionA)
					fmt.Fprintf(&b, "B. %s\n", answerOptions.OptionB)
					fmt.Fprintf(&b, "C. %s\n", answerOptions.OptionC)
					fmt.Fprintf(&b, "D. %s\n", answerOptions.OptionD)
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

// SubmitPart3AnswerTool creates the MCP tool definition for submitting Part3 answers
func SubmitPart3AnswerTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "submit_part3_answer",
		Description: "Submit Part3 practice answer. Optionally provide is_correct parameter for AI assistant to judge answer correctness instead of comparing with database.",
	}
}

// HandleSubmitPart3Answer handles the request to submit Part3 answers
func HandleSubmitPart3Answer(user *model.User, p3Service service.Part3PracticeSetService) func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.SubmitPart3AnswerParams) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, params *mcpModel.SubmitPart3AnswerParams) (*mcp.CallToolResult, any, error) {
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
		itemDetail, err := p3Service.GetItemDetail(ctx, user.ID, params.ItemID)
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
			// Find the answer options for this question
			for _, opt := range itemDetail.Conversation.AnswerOptions {
				if opt.QuestionNumber == itemDetail.QuestionIndex {
					if opt.CorrectAnswer.Valid {
						correctAnswer = strings.ToUpper(strings.TrimSpace(opt.CorrectAnswer.String))
					}
					break
				}
			}
			isCorrect = answer == correctAnswer
		}

		// Get correct answer for display (if available)
		if correctAnswer == "" {
			for _, opt := range itemDetail.Conversation.AnswerOptions {
				if opt.QuestionNumber == itemDetail.QuestionIndex {
					if opt.CorrectAnswer.Valid {
						correctAnswer = strings.ToUpper(strings.TrimSpace(opt.CorrectAnswer.String))
					}
					break
				}
			}
		}

		// Submit answer
		submitReq := &service.SubmitPart3AnswerRequest{
			SelectedAnswer: answer,
			IsCorrect:      isCorrect,
		}
		if err := p3Service.SubmitAnswer(ctx, user.ID, params.ItemID, submitReq); err != nil {
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

