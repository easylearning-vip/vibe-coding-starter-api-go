package test

import (
	"context"
	"database/sql"
	"testing"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/service"
)

// TestGeneratedModels tests that all generated models can be instantiated
func TestGeneratedModels(t *testing.T) {
	// Test DifficultyLevel model
	difficultyLevel := &model.DifficultyLevel{}
	if difficultyLevel == nil {
		t.Error("Failed to create DifficultyLevel model")
	}

	// Test Scenario model
	scenario := &model.Scenario{}
	if scenario == nil {
		t.Error("Failed to create Scenario model")
	}

	// Test Test model
	testModel := &model.Test{}
	if testModel == nil {
		t.Error("Failed to create Test model")
	}

	// Test Part2Question model
	part2Question := &model.Part2Question{}
	if part2Question == nil {
		t.Error("Failed to create Part2Question model")
	}

	// Test Part3Conversation model
	part3Conversation := &model.Part3Conversation{}
	if part3Conversation == nil {
		t.Error("Failed to create Part3Conversation model")
	}

	// Test Part3AnswerOption model
	part3AnswerOption := &model.Part3AnswerOption{}
	if part3AnswerOption == nil {
		t.Error("Failed to create Part3AnswerOption model")
	}

	// Test Part4Talk model
	part4Talk := &model.Part4Talk{}
	if part4Talk == nil {
		t.Error("Failed to create Part4Talk model")
	}

	// Test Part4AnswerOption model
	part4AnswerOption := &model.Part4AnswerOption{}
	if part4AnswerOption == nil {
		t.Error("Failed to create Part4AnswerOption model")
	}

	t.Log("All generated models can be instantiated successfully")
}

// TestGeneratedServiceRequests tests that all generated service request structs can be instantiated
func TestGeneratedServiceRequests(t *testing.T) {
	ctx := context.Background()
	_ = ctx // Use context to avoid unused import error

	// Test Part2Question service requests
	createPart2QuestionReq := &service.CreatePart2QuestionRequest{
		TestId:            1,
		QuestionNumber:    1,
		QuestionText:      "Test Question",
		OptionA:           "Option A",
		OptionB:           "Option B",
		OptionC:           "Option C",
		CorrectAnswer:     sql.NullString{String: "A", Valid: true},
		ScenarioId:        sql.NullInt32{Int32: 1, Valid: true},
		DifficultyLevelId: sql.NullInt32{Int32: 1, Valid: true},
	}
	if createPart2QuestionReq == nil {
		t.Error("Failed to create CreatePart2QuestionRequest")
	}

	// Test Part3Conversation service requests
	createPart3ConversationReq := &service.CreatePart3ConversationRequest{
		TestId:             1,
		ConversationNumber: 1,
		Title:              sql.NullString{String: "Test Conversation", Valid: true},
		Content:            "Test Content",
		Question1:          "Question 1",
		Question2:          "Question 2",
		Question3:          "Question 3",
		ScenarioId:         sql.NullInt32{Int32: 1, Valid: true},
		DifficultyLevelId:  sql.NullInt32{Int32: 1, Valid: true},
	}
	if createPart3ConversationReq == nil {
		t.Error("Failed to create CreatePart3ConversationRequest")
	}

	// Test Part3AnswerOption service requests
	createPart3AnswerOptionReq := &service.CreatePart3AnswerOptionRequest{
		ConversationId: 1,
		QuestionNumber: 1,
		OptionA:        "Option A",
		OptionB:        "Option B",
		OptionC:        "Option C",
		OptionD:        "Option D",
		CorrectAnswer:  sql.NullString{String: "A", Valid: true},
	}
	if createPart3AnswerOptionReq == nil {
		t.Error("Failed to create CreatePart3AnswerOptionRequest")
	}

	// Test Part4Talk service requests
	createPart4TalkReq := &service.CreatePart4TalkRequest{
		TestId:            1,
		TalkNumber:        1,
		Title:             sql.NullString{String: "Test Talk", Valid: true},
		Content:           "Test Content",
		Question1:         "Question 1",
		Question2:         "Question 2",
		Question3:         "Question 3",
		ScenarioId:        sql.NullInt32{Int32: 1, Valid: true},
		DifficultyLevelId: sql.NullInt32{Int32: 1, Valid: true},
	}
	if createPart4TalkReq == nil {
		t.Error("Failed to create CreatePart4TalkRequest")
	}

	// Test Part4AnswerOption service requests
	createPart4AnswerOptionReq := &service.CreatePart4AnswerOptionRequest{
		TalkId:         1,
		QuestionNumber: 1,
		OptionA:        "Option A",
		OptionB:        "Option B",
		OptionC:        "Option C",
		OptionD:        "Option D",
		CorrectAnswer:  sql.NullString{String: "A", Valid: true},
	}
	if createPart4AnswerOptionReq == nil {
		t.Error("Failed to create CreatePart4AnswerOptionRequest")
	}

	t.Log("All generated service request structs can be instantiated successfully")
}

// TestTableNames tests that all models return correct table names
func TestTableNames(t *testing.T) {
	tests := []struct {
		model     interface{ TableName() string }
		expected  string
		modelName string
	}{
		{&model.DifficultyLevel{}, "difficulty_levels", "DifficultyLevel"},
		{&model.Scenario{}, "scenarios", "Scenario"},
		{&model.Test{}, "tests", "Test"},
		{&model.Part2Question{}, "part2_questions", "Part2Question"},
		{&model.Part3Conversation{}, "part3_conversations", "Part3Conversation"},
		{&model.Part3AnswerOption{}, "part3_answer_options", "Part3AnswerOption"},
		{&model.Part4Talk{}, "part4_talks", "Part4Talk"},
		{&model.Part4AnswerOption{}, "part4_answer_options", "Part4AnswerOption"},
	}

	for _, tt := range tests {
		t.Run(tt.modelName, func(t *testing.T) {
			if got := tt.model.TableName(); got != tt.expected {
				t.Errorf("%s.TableName() = %v, want %v", tt.modelName, got, tt.expected)
			}
		})
	}

	t.Log("All generated models return correct table names")
}
