package service

import (
	"context"
	"fmt"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/logger"
)

// articleService 文章服务实现
type articleService struct {
	articleRepo repository.ArticleRepository
	cache       cache.Cache
	logger      logger.Logger
}

// NewArticleService 创建文章服务
func NewArticleService(
	articleRepo repository.ArticleRepository,
	cache cache.Cache,
	logger logger.Logger,
) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		cache:       cache,
		logger:      logger,
	}
}

// Create 创建文章
func (s *articleService) Create(ctx context.Context, req *CreateArticleRequest) (*model.Article, error) {
	article := &model.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		CoverImage: req.CoverImage,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		AuthorID:   req.AuthorID,
	}

	// 设置默认状态
	if article.Status == "" {
		article.Status = model.ArticleStatusDraft
	}

	// 创建文章
	if err := s.articleRepo.Create(ctx, article); err != nil {
		s.logger.Error("Failed to create article", "title", req.Title, "error", err)
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	s.logger.Info("Article created successfully", "article_id", article.ID, "title", article.Title)
	return article, nil
}

// GetByID 根据 ID 获取文章
func (s *articleService) GetByID(ctx context.Context, id uint) (*model.Article, error) {
	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get article by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	return article, nil
}

// GetBySlug 根据 slug 获取文章
func (s *articleService) GetBySlug(ctx context.Context, slug string) (*model.Article, error) {
	article, err := s.articleRepo.GetBySlug(ctx, slug)
	if err != nil {
		s.logger.Error("Failed to get article by slug", "slug", slug, "error", err)
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	return article, nil
}

// Update 更新文章
func (s *articleService) Update(ctx context.Context, id uint, req *UpdateArticleRequest) (*model.Article, error) {
	// 获取现有文章
	article, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get article for update", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	// 更新字段
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.Summary != "" {
		article.Summary = req.Summary
	}
	if req.CoverImage != "" {
		article.CoverImage = req.CoverImage
	}
	if req.CategoryID != nil {
		article.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		article.Status = req.Status
	}

	// 保存更新
	if err := s.articleRepo.Update(ctx, article); err != nil {
		s.logger.Error("Failed to update article", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	s.logger.Info("Article updated successfully", "article_id", id)
	return article, nil
}

// Delete 删除文章
func (s *articleService) Delete(ctx context.Context, id uint) error {
	// 检查文章是否存在
	_, err := s.articleRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get article for deletion", "id", id, "error", err)
		return fmt.Errorf("failed to get article: %w", err)
	}

	// 删除文章
	if err := s.articleRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete article", "id", id, "error", err)
		return fmt.Errorf("failed to delete article: %w", err)
	}

	s.logger.Info("Article deleted successfully", "article_id", id)
	return nil
}

// List 获取文章列表
func (s *articleService) List(ctx context.Context, opts repository.ListOptions) ([]*model.Article, int64, error) {
	articles, total, err := s.articleRepo.List(ctx, opts)
	if err != nil {
		s.logger.Error("Failed to get articles list", "error", err)
		return nil, 0, fmt.Errorf("failed to get articles: %w", err)
	}

	return articles, total, nil
}

// GetPublished 获取已发布的文章列表
func (s *articleService) GetPublished(ctx context.Context, opts repository.ListOptions) ([]*model.Article, int64, error) {
	articles, total, err := s.articleRepo.GetPublished(ctx, opts)
	if err != nil {
		s.logger.Error("Failed to get published articles", "error", err)
		return nil, 0, fmt.Errorf("failed to get published articles: %w", err)
	}

	return articles, total, nil
}

// Search 搜索文章
func (s *articleService) Search(ctx context.Context, query string, opts repository.ListOptions) ([]*model.Article, int64, error) {
	if query == "" {
		return s.List(ctx, opts)
	}

	articles, total, err := s.articleRepo.Search(ctx, query, opts)
	if err != nil {
		s.logger.Error("Failed to search articles", "query", query, "error", err)
		return nil, 0, fmt.Errorf("failed to search articles: %w", err)
	}

	return articles, total, nil
}

// IncrementViewCount 增加浏览次数
func (s *articleService) IncrementViewCount(ctx context.Context, articleID uint) error {
	if err := s.articleRepo.IncrementViewCount(ctx, articleID); err != nil {
		s.logger.Error("Failed to increment view count", "article_id", articleID, "error", err)
		return fmt.Errorf("failed to increment view count: %w", err)
	}

	return nil
}
