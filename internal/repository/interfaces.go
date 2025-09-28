package repository

import (
	"context"

	"vibe-coding-starter/internal/model"
)

// Repository 通用仓储接口
type Repository[T any, ID comparable] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id ID) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id ID) error
	List(ctx context.Context, opts ListOptions) ([]*T, int64, error)
}

// ListOptions 列表查询选项
type ListOptions struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Sort     string                 `json:"sort"`
	Order    string                 `json:"order"`
	Filters  map[string]interface{} `json:"filters"`
	Search   string                 `json:"search"`
}

// UserRepository 用户仓储接口
type UserRepository interface {
	Repository[model.User, uint]
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
	// Token-based operations
	GetByToken(ctx context.Context, token string) (*model.User, error)
	UpdateToken(ctx context.Context, userID uint, token *string) error
}

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	Repository[model.Article, uint]
	GetBySlug(ctx context.Context, slug string) (*model.Article, error)
	GetByAuthor(ctx context.Context, authorID uint, opts ListOptions) ([]*model.Article, int64, error)
	GetByCategory(ctx context.Context, categoryID uint, opts ListOptions) ([]*model.Article, int64, error)
	GetByTag(ctx context.Context, tagID uint, opts ListOptions) ([]*model.Article, int64, error)
	GetPublished(ctx context.Context, opts ListOptions) ([]*model.Article, int64, error)
	Search(ctx context.Context, query string, opts ListOptions) ([]*model.Article, int64, error)
	IncrementViewCount(ctx context.Context, articleID uint) error
}

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	Repository[model.Category, uint]
	GetBySlug(ctx context.Context, slug string) (*model.Category, error)
	GetByName(ctx context.Context, name string) (*model.Category, error)
}

// TagRepository 标签仓储接口
type TagRepository interface {
	Repository[model.Tag, uint]
	GetBySlug(ctx context.Context, slug string) (*model.Tag, error)
	GetByName(ctx context.Context, name string) (*model.Tag, error)
	GetByNames(ctx context.Context, names []string) ([]*model.Tag, error)
}

// CommentRepository 评论仓储接口
type CommentRepository interface {
	Repository[model.Comment, uint]
	GetByArticle(ctx context.Context, articleID uint, opts ListOptions) ([]*model.Comment, int64, error)
	GetByAuthor(ctx context.Context, authorID uint, opts ListOptions) ([]*model.Comment, int64, error)
	GetReplies(ctx context.Context, parentID uint, opts ListOptions) ([]*model.Comment, int64, error)
}

// FileRepository 文件仓储接口
type FileRepository interface {
	Repository[model.File, uint]
	GetByHash(ctx context.Context, hash string) (*model.File, error)
	GetByOwner(ctx context.Context, ownerID uint, opts ListOptions) ([]*model.File, int64, error)
}

// DictCategoryRepository 数据字典分类仓储接口
type DictCategoryRepository interface {
	Repository[model.DictCategory, uint]
	GetByCode(ctx context.Context, code string) (*model.DictCategory, error)
}

// DictItemRepository 数据字典项仓储接口
type DictItemRepository interface {
	Repository[model.DictItem, uint]
	GetByCategory(ctx context.Context, categoryCode string) ([]*model.DictItem, error)
	GetActiveByCategory(ctx context.Context, categoryCode string) ([]*model.DictItem, error)
	GetByCategoryAndKey(ctx context.Context, categoryCode, itemKey string) (*model.DictItem, error)
}

// DictRepository 数据字典仓储接口（组合接口）
type DictRepository interface {
	// 分类相关方法
	GetAllCategories(ctx context.Context) ([]*model.DictCategory, error)
	GetCategoryByCode(ctx context.Context, code string) (*model.DictCategory, error)
	CreateCategory(ctx context.Context, category *model.DictCategory) error
	DeleteCategory(ctx context.Context, id uint) error

	// 字典项相关方法
	GetItemsByCategory(ctx context.Context, categoryCode string) ([]*model.DictItem, error)
	GetActiveItemsByCategory(ctx context.Context, categoryCode string) ([]*model.DictItem, error)
	GetItemByID(ctx context.Context, id uint) (*model.DictItem, error)
	CreateItem(ctx context.Context, item *model.DictItem) error
	UpdateItem(ctx context.Context, item *model.DictItem) error
	DeleteItem(ctx context.Context, id uint) error
}

// DepartmentRepository Department仓储接口
type DepartmentRepository interface {
	Repository[model.Department, uint]
	// 在这里添加特定的查询方法
	GetByName(ctx context.Context, name string) (*model.Department, error)
	GetByParentId(ctx context.Context, parentId uint) ([]*model.Department, error)
	GetByCode(ctx context.Context, code string) (*model.Department, error)
	GetChildrenTree(ctx context.Context, parentId uint) ([]*model.Department, error)
}

// DifficultyLevelRepository DifficultyLevel仓储接口
type DifficultyLevelRepository interface {
	Repository[model.DifficultyLevel, uint]
	// 在这里添加特定的查询方法
	GetByName(ctx context.Context, name string) (*model.DifficultyLevel, error)
}

// ScenarioRepository Scenario仓储接口
type ScenarioRepository interface {
	Repository[model.Scenario, uint]
	// 在这里添加特定的查询方法
	GetByName(ctx context.Context, name string) (*model.Scenario, error)
}

// TestRepository Test仓储接口
type TestRepository interface {
	Repository[model.Test, uint]
	// 在这里添加特定的查询方法
	GetByName(ctx context.Context, name string) (*model.Test, error)
}

// Part2QuestionRepository Part2Question仓储接口
type Part2QuestionRepository interface {
	Repository[model.Part2Question, uint]
	// 在这里添加特定的查询方法
	GetByQuestionText(ctx context.Context, questionText string) (*model.Part2Question, error)
	DeleteByTestID(ctx context.Context, testID uint) error
}

// Part3ConversationRepository Part3Conversation仓储接口
type Part3ConversationRepository interface {
	Repository[model.Part3Conversation, uint]
	// 在这里添加特定的查询方法
	GetByName(ctx context.Context, name string) (*model.Part3Conversation, error)
	DeleteByTestID(ctx context.Context, testID uint) error
	// 预加载答案选项的列表查询
	ListWithAnswerOptions(ctx context.Context, opts ListOptions) ([]*model.Part3Conversation, int64, error)
}

// Part3AnswerOptionRepository Part3AnswerOption仓储接口
type Part3AnswerOptionRepository interface {
	Repository[model.Part3AnswerOption, uint]
	// 在这里添加特定的查询方法
	GetByName(ctx context.Context, name string) (*model.Part3AnswerOption, error)
	DeleteByConversationIDs(ctx context.Context, conversationIDs []uint) error
}

// Part4TalkRepository Part4Talk仓储接口
type Part4TalkRepository interface {
	Repository[model.Part4Talk, uint]
	// 在这里添加特定的查询方法
	GetByName(ctx context.Context, name string) (*model.Part4Talk, error)
	DeleteByTestID(ctx context.Context, testID uint) error
	// 预加载答案选项的列表查询
	ListWithAnswerOptions(ctx context.Context, opts ListOptions) ([]*model.Part4Talk, int64, error)
}

// Part4AnswerOptionRepository Part4AnswerOption仓储接口
type Part4AnswerOptionRepository interface {
	Repository[model.Part4AnswerOption, uint]
	// 在这里添加特定的查询方法
	GetByName(ctx context.Context, name string) (*model.Part4AnswerOption, error)
	DeleteByTalkIDs(ctx context.Context, talkIDs []uint) error
}

// Part2PracticeItemRepository 个人Part2练习记录仓储接口

// Part3PracticeItemRepository 个人Part3练习记录仓储接口

// Part4PracticeItemRepository 个人Part4练习记录仓储接口

// Practice Set repositories (Master/Detail) for Part2/Part3/Part4
// Part2
type Part2PracticeSetRepository interface {
	Repository[model.Part2PracticeSet, uint]
}
type Part2PracticeSetItemRepository interface {
	Repository[model.Part2PracticeSetItem, uint]
}

// Part3
type Part3PracticeSetRepository interface {
	Repository[model.Part3PracticeSet, uint]
}
type Part3PracticeSetItemRepository interface {
	Repository[model.Part3PracticeSetItem, uint]
}

// Part4
type Part4PracticeSetRepository interface {
	Repository[model.Part4PracticeSet, uint]
}
type Part4PracticeSetItemRepository interface {
	Repository[model.Part4PracticeSetItem, uint]
}
