package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/database"
	"vibe-coding-starter/pkg/logger"
)

// userRepository 用户仓储实现
type userRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db database.Database, logger logger.Logger) UserRepository {
	return &userRepository{
		db:     db.GetDB(),
		logger: logger,
	}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.Error("Failed to create user", "error", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID 根据 ID 获取用户
func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with id %d", id)
		}
		r.logger.Error("Failed to get user by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		r.logger.Error("Failed to update user", "id", user.ID, "error", err)
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		r.logger.Error("Failed to delete user", "id", id, "error", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, opts ListOptions) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// 应用过滤器
	query = r.applyFilters(query, opts.Filters)

	// 应用搜索
	if opts.Search != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count users", "error", err)
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
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
	if err := query.Find(&users).Error; err != nil {
		r.logger.Error("Failed to list users", "error", err)
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		r.logger.Error("Failed to get user by email", "email", email, "error", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		r.logger.Error("Failed to get user by username", "username", username, "error", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("last_login", gorm.Expr("NOW()")).Error; err != nil {
		r.logger.Error("Failed to update last login", "user_id", userID, "error", err)
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// applyFilters 应用过滤器
func (r *userRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		switch key {
		case "role":
			query = query.Where("role = ?", value)
		case "status":
			query = query.Where("status = ?", value)
		case "created_after":
			query = query.Where("created_at >= ?", value)
		case "created_before":
			query = query.Where("created_at <= ?", value)
		}
	}

	return query
}
