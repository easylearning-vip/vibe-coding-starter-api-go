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

// CategoryWithProductCount 包含产品数量的分类信息
type CategoryWithProductCount struct {
	*model.ProductCategory
	ProductCount int64 `json:"product_count"`
}

// UserRepository 用户仓储接口
type UserRepository interface {
	Repository[model.User, uint]
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
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

// ProductCategoryRepository ProductCategory仓储接口
type ProductCategoryRepository interface {
	Repository[model.ProductCategory, uint]
	// 在这里添加特定的查询方法
	GetByName(ctx context.Context, name string) (*model.ProductCategory, error)
	GetByParentID(ctx context.Context, parentID uint) ([]*model.ProductCategory, error)
	GetCategoryPath(ctx context.Context, categoryID uint) ([]*model.ProductCategory, error)
	CountProductsByCategory(ctx context.Context, categoryID uint) (int64, error)
	GetCategoriesWithProductCount(ctx context.Context, opts ListOptions) ([]*CategoryWithProductCount, int64, error)
	HasChildren(ctx context.Context, categoryID uint) (bool, error)
	HasProducts(ctx context.Context, categoryID uint) (bool, error)
}

// ProductRepository Product仓储接口
type ProductRepository interface {
	Repository[model.Product, uint]
	// 在这里添加特定的查询方法
	GetByName(ctx context.Context, name string) (*model.Product, error)
	GetBySKU(ctx context.Context, sku string) (*model.Product, error)
	GetByCategoryID(ctx context.Context, categoryID uint, opts ListOptions) ([]*model.Product, int64, error)
	GetByCategoryIDs(ctx context.Context, categoryIDs []uint, opts ListOptions) ([]*model.Product, int64, error)
	GetByPriceRange(ctx context.Context, minPrice float64, maxPrice float64, opts ListOptions) ([]*model.Product, int64, error)
	GetLowStockProducts(ctx context.Context, opts ListOptions) ([]*model.Product, int64, error)
	GetHotSellingProducts(ctx context.Context, minSales int, opts ListOptions) ([]*model.Product, int64, error)
	SearchByName(ctx context.Context, query string, opts ListOptions) ([]*model.Product, int64, error)
	SearchBySKU(ctx context.Context, query string, opts ListOptions) ([]*model.Product, int64, error)
	GetActiveProducts(ctx context.Context, opts ListOptions) ([]*model.Product, int64, error)
	GetProductsInStock(ctx context.Context, opts ListOptions) ([]*model.Product, int64, error)
	UpdateStock(ctx context.Context, productID uint, quantityChange int) error
	BatchUpdatePrices(ctx context.Context, updates map[uint]float64) error
}
