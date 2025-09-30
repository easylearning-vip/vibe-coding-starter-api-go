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
}
