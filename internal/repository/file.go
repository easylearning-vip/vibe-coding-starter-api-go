package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// fileRepository 文件仓储实现
type fileRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewFileRepository 创建文件仓储
func NewFileRepository(db database.Database, logger logger.Logger) FileRepository {
	return &fileRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建文件记录
func (r *fileRepository) Create(ctx context.Context, file *model.File) error {
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		r.logger.Error("Failed to create file", "error", err)
		return fmt.Errorf("failed to create file: %w", err)
	}
	return nil
}

// GetByID 根据 ID 获取文件
func (r *fileRepository) GetByID(ctx context.Context, id uint) (*model.File, error) {
	var file model.File
	if err := r.db.WithContext(ctx).
		Preload("Owner").
		First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("file not found with id %d", id)
		}
		r.logger.Error("Failed to get file by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return &file, nil
}

// Update 更新文件记录
func (r *fileRepository) Update(ctx context.Context, file *model.File) error {
	if err := r.db.WithContext(ctx).Save(file).Error; err != nil {
		r.logger.Error("Failed to update file", "id", file.ID, "error", err)
		return fmt.Errorf("failed to update file: %w", err)
	}
	return nil
}

// Delete 删除文件记录
func (r *fileRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.File{}, id).Error; err != nil {
		r.logger.Error("Failed to delete file", "id", id, "error", err)
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// List 获取文件列表
func (r *fileRepository) List(ctx context.Context, opts ListOptions) ([]*model.File, int64, error) {
	var files []*model.File
	var total int64

	query := r.db.WithContext(ctx).Model(&model.File{}).
		Preload("Owner")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("name LIKE ? OR original_name LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count files", "error", err)
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
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
	if err := query.Find(&files).Error; err != nil {
		r.logger.Error("Failed to list files", "error", err)
		return nil, 0, fmt.Errorf("failed to list files: %w", err)
	}

	return files, total, nil
}

// GetByHash 根据哈希值获取文件
func (r *fileRepository) GetByHash(ctx context.Context, hash string) (*model.File, error) {
	var file model.File
	if err := r.db.WithContext(ctx).
		Preload("Owner").
		Where("hash = ?", hash).
		First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("file not found with hash %s", hash)
		}
		r.logger.Error("Failed to get file by hash", "hash", hash, "error", err)
		return nil, fmt.Errorf("failed to get file: %w", err)
	}
	return &file, nil
}

// GetByOwner 根据所有者获取文件列表
func (r *fileRepository) GetByOwner(ctx context.Context, ownerID uint, opts ListOptions) ([]*model.File, int64, error) {
	var files []*model.File
	var total int64

	query := r.db.WithContext(ctx).Model(&model.File{}).
		Where("owner_id = ?", ownerID).
		Preload("Owner")

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("name LIKE ? OR original_name LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count files by owner", "owner_id", ownerID, "error", err)
		return nil, 0, fmt.Errorf("failed to count files: %w", err)
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
	if err := query.Find(&files).Error; err != nil {
		r.logger.Error("Failed to get files by owner", "owner_id", ownerID, "error", err)
		return nil, 0, fmt.Errorf("failed to get files: %w", err)
	}

	return files, total, nil
}

// applyFilters 应用过滤器
func (r *fileRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		switch key {
		case "storage_type":
			query = query.Where("storage_type = ?", value)
		case "mime_type":
			query = query.Where("mime_type = ?", value)
		case "extension":
			query = query.Where("extension = ?", value)
		case "is_public":
			query = query.Where("is_public = ?", value)
		case "owner_id":
			query = query.Where("owner_id = ?", value)
		case "size_min":
			query = query.Where("size >= ?", value)
		case "size_max":
			query = query.Where("size <= ?", value)
		case "created_after":
			query = query.Where("created_at >= ?", value)
		case "created_before":
			query = query.Where("created_at <= ?", value)
		}
	}

	return query
}
