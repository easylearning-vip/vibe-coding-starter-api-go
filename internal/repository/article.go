package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// articleRepository 文章仓储实现
type articleRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewArticleRepository 创建文章仓储
func NewArticleRepository(db database.Database, logger logger.Logger) ArticleRepository {
	return &articleRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建文章
func (r *articleRepository) Create(ctx context.Context, article *model.Article) error {
	if err := r.db.WithContext(ctx).Create(article).Error; err != nil {
		r.logger.Error("Failed to create article", "error", err)
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}

// GetByID 根据 ID 获取文章
func (r *articleRepository) GetByID(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	if err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("article not found with id %d", id)
		}
		r.logger.Error("Failed to get article by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get article: %w", err)
	}
	return &article, nil
}

// Update 更新文章
func (r *articleRepository) Update(ctx context.Context, article *model.Article) error {
	if err := r.db.WithContext(ctx).Save(article).Error; err != nil {
		r.logger.Error("Failed to update article", "id", article.ID, "error", err)
		return fmt.Errorf("failed to update article: %w", err)
	}
	return nil
}

// Delete 删除文章
func (r *articleRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Article{}, id).Error; err != nil {
		r.logger.Error("Failed to delete article", "id", id, "error", err)
		return fmt.Errorf("failed to delete article: %w", err)
	}
	return nil
}

// List 获取文章列表
func (r *articleRepository) List(ctx context.Context, opts ListOptions) ([]*model.Article, int64, error) {
	var articles []*model.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Article{}).
		Preload("Author").
		Preload("Category").
		Preload("Tags")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("title LIKE ? OR content LIKE ? OR summary LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count articles", "error", err)
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// 应用分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	// 执行查询
	if err := query.Find(&articles).Error; err != nil {
		r.logger.Error("Failed to list articles", "error", err)
		return nil, 0, fmt.Errorf("failed to list articles: %w", err)
	}

	return articles, total, nil
}

// GetBySlug 根据 slug 获取文章
func (r *articleRepository) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	var article model.Article
	if err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Where("slug = ?", slug).
		First(&article).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("article not found with slug %s", slug)
		}
		r.logger.Error("Failed to get article by slug", "slug", slug, "error", err)
		return nil, fmt.Errorf("failed to get article: %w", err)
	}
	return &article, nil
}

// GetByAuthor 根据作者获取文章列表
func (r *articleRepository) GetByAuthor(ctx context.Context, authorID uint, opts ListOptions) ([]*model.Article, int64, error) {
	var articles []*model.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Article{}).
		Where("author_id = ?", authorID).
		Preload("Author").
		Preload("Category").
		Preload("Tags")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count articles by author", "author_id", authorID, "error", err)
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	// 应用排序和分页
	query = r.applySortAndPagination(query, opts)

	// 执行查询
	if err := query.Find(&articles).Error; err != nil {
		r.logger.Error("Failed to get articles by author", "author_id", authorID, "error", err)
		return nil, 0, fmt.Errorf("failed to get articles: %w", err)
	}

	return articles, total, nil
}

// GetByCategory 根据分类获取文章列表
func (r *articleRepository) GetByCategory(ctx context.Context, categoryID uint, opts ListOptions) ([]*model.Article, int64, error) {
	var articles []*model.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Article{}).
		Where("category_id = ?", categoryID).
		Preload("Author").
		Preload("Category").
		Preload("Tags")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count articles by category", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	// 应用排序和分页
	query = r.applySortAndPagination(query, opts)

	// 执行查询
	if err := query.Find(&articles).Error; err != nil {
		r.logger.Error("Failed to get articles by category", "category_id", categoryID, "error", err)
		return nil, 0, fmt.Errorf("failed to get articles: %w", err)
	}

	return articles, total, nil
}

// GetByTag 根据标签获取文章列表
func (r *articleRepository) GetByTag(ctx context.Context, tagID uint, opts ListOptions) ([]*model.Article, int64, error) {
	var articles []*model.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Article{}).
		Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Where("article_tags.tag_id = ?", tagID).
		Preload("Author").
		Preload("Category").
		Preload("Tags")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count articles by tag", "tag_id", tagID, "error", err)
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	// 应用排序和分页
	query = r.applySortAndPagination(query, opts)

	// 执行查询
	if err := query.Find(&articles).Error; err != nil {
		r.logger.Error("Failed to get articles by tag", "tag_id", tagID, "error", err)
		return nil, 0, fmt.Errorf("failed to get articles: %w", err)
	}

	return articles, total, nil
}

// GetPublished 获取已发布的文章列表
func (r *articleRepository) GetPublished(ctx context.Context, opts ListOptions) ([]*model.Article, int64, error) {
	var articles []*model.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Article{}).
		Where("status = ?", model.ArticleStatusPublished).
		Preload("Author").
		Preload("Category").
		Preload("Tags")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("title LIKE ? OR content LIKE ? OR summary LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count published articles", "error", err)
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	// 应用排序和分页
	query = r.applySortAndPagination(query, opts)

	// 执行查询
	if err := query.Find(&articles).Error; err != nil {
		r.logger.Error("Failed to get published articles", "error", err)
		return nil, 0, fmt.Errorf("failed to get articles: %w", err)
	}

	return articles, total, nil
}

// Search 搜索文章
func (r *articleRepository) Search(ctx context.Context, query string, opts ListOptions) ([]*model.Article, int64, error) {
	var articles []*model.Article
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&model.Article{}).
		Where("title LIKE ? OR content LIKE ? OR summary LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Preload("Author").
		Preload("Category").
		Preload("Tags")

	// 应用过滤器
	dbQuery = r.applyFilters(dbQuery, opts.Filters)

	// 获取总数
	if err := dbQuery.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count search results", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	// 应用排序和分页
	dbQuery = r.applySortAndPagination(dbQuery, opts)

	// 执行查询
	if err := dbQuery.Find(&articles).Error; err != nil {
		r.logger.Error("Failed to search articles", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to search articles: %w", err)
	}

	return articles, total, nil
}

// IncrementViewCount 增加浏览次数
func (r *articleRepository) IncrementViewCount(ctx context.Context, articleID uint) error {
	if err := r.db.WithContext(ctx).Model(&model.Article{}).
		Where("id = ?", articleID).
		Update("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
		r.logger.Error("Failed to increment view count", "article_id", articleID, "error", err)
		return fmt.Errorf("failed to increment view count: %w", err)
	}
	return nil
}

// applyFilters 应用过滤器
func (r *articleRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		switch key {
		case "status":
			query = query.Where("status = ?", value)
		case "author_id":
			query = query.Where("author_id = ?", value)
		case "category_id":
			query = query.Where("category_id = ?", value)
		case "created_after":
			query = query.Where("created_at >= ?", value)
		case "created_before":
			query = query.Where("created_at <= ?", value)
		case "published_after":
			query = query.Where("published_at >= ?", value)
		case "published_before":
			query = query.Where("published_at <= ?", value)
		}
	}

	return query
}

// applySortAndPagination 应用排序和分页
func (r *articleRepository) applySortAndPagination(query *gorm.DB, opts ListOptions) *gorm.DB {
	// 应用排序
	if opts.Sort != "" {
		order := "ASC"
		if opts.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// 应用分页
	if opts.Page > 0 && opts.PageSize > 0 {
		offset := (opts.Page - 1) * opts.PageSize
		query = query.Offset(offset).Limit(opts.PageSize)
	}

	return query
}
