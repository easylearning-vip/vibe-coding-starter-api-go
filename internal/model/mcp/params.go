package mcp

// GetUserInfoParams 获取用户信息参数
type GetUserInfoParams struct{}

// Part2ListParams Part2练习集列表查询参数
type Part2ListParams struct {
	Page     int    `json:"page" jsonschema:"Page number (default 1)"`
	PageSize int    `json:"page_size" jsonschema:"Page size (default 20, max 100)"`
	Sort     string `json:"sort" jsonschema:"Sort column (optional)"`
	Order    string `json:"order" jsonschema:"Sort order asc|desc (optional)"`
}

// AutoGeneratePart2SetParams 自动生成Part2练习集参数
type AutoGeneratePart2SetParams struct {
	Count             int   `json:"count" jsonschema:"Number of questions to include (required, 10-30)"`
	DifficultyLevelID int32 `json:"difficulty_level_id" jsonschema:"Difficulty level ID (optional)"`
	ScenarioID        int32 `json:"scenario_id" jsonschema:"Scenario ID (optional)"`
}

// GetPart2SetDetailsParams 获取Part2练习集详情参数
type GetPart2SetDetailsParams struct {
	SetID uint `json:"set_id" jsonschema:"Practice set ID (required)"`
}

// SubmitPart2AnswerParams 提交Part2答案参数
type SubmitPart2AnswerParams struct {
	SetID      uint   `json:"set_id" jsonschema:"Practice set ID (required)"`
	ItemID     uint   `json:"item_id" jsonschema:"Practice set item ID (required)"`
	UserAnswer string `json:"user_answer" jsonschema:"User's answer choice (required, e.g., A, B, or C)"`
	IsCorrect  *bool  `json:"is_correct" jsonschema:"Whether the answer is correct (optional, if not provided, will compare with correct answer in database)"`
}

// Part2PromptParams 获取Part2 AI提示词参数
type Part2PromptParams struct {
	// No parameters needed - always returns the latest prompt
}

// Part3ListParams Part3练习集列表查询参数
type Part3ListParams struct {
	Page     int    `json:"page" jsonschema:"Page number (default 1)"`
	PageSize int    `json:"page_size" jsonschema:"Page size (default 20, max 100)"`
	Sort     string `json:"sort" jsonschema:"Sort column (optional)"`
	Order    string `json:"order" jsonschema:"Sort order asc|desc (optional)"`
}

// AutoGeneratePart3SetParams 自动生成Part3练习集参数
type AutoGeneratePart3SetParams struct {
	Count             int   `json:"count" jsonschema:"Number of questions to include (required, 10-30)"`
	DifficultyLevelID int32 `json:"difficulty_level_id" jsonschema:"Difficulty level ID (optional)"`
	ScenarioID        int32 `json:"scenario_id" jsonschema:"Scenario ID (optional)"`
}

// GetPart3SetDetailsParams 获取Part3练习集详情参数
type GetPart3SetDetailsParams struct {
	SetID uint `json:"set_id" jsonschema:"Practice set ID (required)"`
}

// SubmitPart3AnswerParams 提交Part3答案参数
type SubmitPart3AnswerParams struct {
	SetID      uint   `json:"set_id" jsonschema:"Practice set ID (required)"`
	ItemID     uint   `json:"item_id" jsonschema:"Practice set item ID (required)"`
	UserAnswer string `json:"user_answer" jsonschema:"User's answer choice (required, e.g., A, B, C, or D)"`
	IsCorrect  *bool  `json:"is_correct" jsonschema:"Whether the answer is correct (optional, if not provided, will compare with correct answer in database)"`
}

// Part3PromptParams 获取Part3 AI提示词参数
type Part3PromptParams struct {
	// No parameters needed - always returns the latest prompt
}

// Part4ListParams Part4练习集列表查询参数
type Part4ListParams struct {
	Page     int    `json:"page" jsonschema:"Page number (default 1)"`
	PageSize int    `json:"page_size" jsonschema:"Page size (default 20, max 100)"`
	Sort     string `json:"sort" jsonschema:"Sort column (optional)"`
	Order    string `json:"order" jsonschema:"Sort order asc|desc (optional)"`
}

// AutoGeneratePart4SetParams 自动生成Part4练习集参数
type AutoGeneratePart4SetParams struct {
	Count             int   `json:"count" jsonschema:"Number of questions to include (required, 10-30)"`
	DifficultyLevelID int32 `json:"difficulty_level_id" jsonschema:"Difficulty level ID (optional)"`
	ScenarioID        int32 `json:"scenario_id" jsonschema:"Scenario ID (optional)"`
}

// GetPart4SetDetailsParams 获取Part4练习集详情参数
type GetPart4SetDetailsParams struct {
	SetID uint `json:"set_id" jsonschema:"Practice set ID (required)"`
}

// SubmitPart4AnswerParams 提交Part4答案参数
type SubmitPart4AnswerParams struct {
	SetID      uint   `json:"set_id" jsonschema:"Practice set ID (required)"`
	ItemID     uint   `json:"item_id" jsonschema:"Practice set item ID (required)"`
	UserAnswer string `json:"user_answer" jsonschema:"User's answer choice (required, e.g., A, B, C, or D)"`
	IsCorrect  *bool  `json:"is_correct" jsonschema:"Whether the answer is correct (optional, if not provided, will compare with correct answer in database)"`
}

// Part4PromptParams 获取Part4 AI提示词参数
type Part4PromptParams struct {
	// No parameters needed - always returns the latest prompt
}
