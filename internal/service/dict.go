package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/logger"
)

// dictService 数据字典服务实现
type dictService struct {
	dictRepo repository.DictRepository
	cache    cache.Cache
	logger   logger.Logger
}

// NewDictService 创建数据字典服务
func NewDictService(
	dictRepo repository.DictRepository,
	cache cache.Cache,
	logger logger.Logger,
) DictService {
	return &dictService{
		dictRepo: dictRepo,
		cache:    cache,
		logger:   logger,
	}
}

// GetDictCategories 获取所有字典分类（带缓存）
func (s *dictService) GetDictCategories(ctx context.Context) ([]*model.DictCategory, error) {
	// 尝试从缓存获取
	cacheKey := "dict_categories:all"
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var categories []*model.DictCategory
		if err := json.Unmarshal([]byte(cached), &categories); err == nil {
			s.logger.Debug("Dict categories retrieved from cache", "count", len(categories))
			return categories, nil
		}
	}

	// 从数据库获取
	categories, err := s.dictRepo.GetAllCategories(ctx)
	if err != nil {
		s.logger.Error("Failed to get dict categories", "error", err)
		return nil, fmt.Errorf("failed to get dict categories: %w", err)
	}

	// 缓存结果
	if data, err := json.Marshal(categories); err == nil {
		s.cache.Set(ctx, cacheKey, string(data), time.Hour*24) // 缓存24小时
	}

	s.logger.Debug("Dict categories retrieved from database", "count", len(categories))
	return categories, nil
}

// GetDictItems 获取字典项（带缓存）
func (s *dictService) GetDictItems(ctx context.Context, categoryCode string) ([]*model.DictItem, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("dict_items:%s", categoryCode)
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var items []*model.DictItem
		if err := json.Unmarshal([]byte(cached), &items); err == nil {
			s.logger.Debug("Dict items retrieved from cache", "category_code", categoryCode, "count", len(items))
			return items, nil
		}
	}

	// 从数据库获取
	items, err := s.dictRepo.GetActiveItemsByCategory(ctx, categoryCode)
	if err != nil {
		s.logger.Error("Failed to get dict items", "category_code", categoryCode, "error", err)
		return nil, fmt.Errorf("failed to get dict items: %w", err)
	}

	// 缓存结果
	if data, err := json.Marshal(items); err == nil {
		s.cache.Set(ctx, cacheKey, string(data), time.Hour*24) // 缓存24小时
	}

	s.logger.Debug("Dict items retrieved from database", "category_code", categoryCode, "count", len(items))
	return items, nil
}

// GetDictItemByKey 获取特定字典项
func (s *dictService) GetDictItemByKey(ctx context.Context, categoryCode, itemKey string) (*model.DictItem, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("dict_item:%s:%s", categoryCode, itemKey)
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var item model.DictItem
		if err := json.Unmarshal([]byte(cached), &item); err == nil {
			s.logger.Debug("Dict item retrieved from cache", "category_code", categoryCode, "item_key", itemKey)
			return &item, nil
		}
	}

	// 从数据库获取
	item, err := s.dictRepo.GetItemsByCategory(ctx, categoryCode)
	if err != nil {
		s.logger.Error("Failed to get dict items for key lookup", "category_code", categoryCode, "item_key", itemKey, "error", err)
		return nil, fmt.Errorf("failed to get dict item: %w", err)
	}

	// 查找特定的键值
	for _, dictItem := range item {
		if dictItem.ItemKey == itemKey {
			// 缓存结果
			if data, err := json.Marshal(dictItem); err == nil {
				s.cache.Set(ctx, cacheKey, string(data), time.Hour*24)
			}
			s.logger.Debug("Dict item found", "category_code", categoryCode, "item_key", itemKey)
			return dictItem, nil
		}
	}

	return nil, fmt.Errorf("dict item not found with category %s and key %s", categoryCode, itemKey)
}

// CreateDictCategory 创建字典分类
func (s *dictService) CreateDictCategory(ctx context.Context, req *CreateCategoryRequest) (*model.DictCategory, error) {
	// 验证分类代码是否已存在
	if existing, err := s.dictRepo.GetCategoryByCode(ctx, req.Code); err == nil && existing != nil {
		return nil, fmt.Errorf("category with code %s already exists", req.Code)
	}

	category := &model.DictCategory{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	if err := s.dictRepo.CreateCategory(ctx, category); err != nil {
		s.logger.Error("Failed to create dict category", "code", req.Code, "error", err)
		return nil, fmt.Errorf("failed to create dict category: %w", err)
	}

	s.logger.Info("Dict category created successfully", "code", req.Code, "name", req.Name)
	return category, nil
}

// DeleteDictCategory 删除字典分类
func (s *dictService) DeleteDictCategory(ctx context.Context, id uint) error {
	// 获取所有分类来找到要删除的分类
	categories, err := s.dictRepo.GetAllCategories(ctx)
	if err != nil {
		s.logger.Error("Failed to get categories", "error", err)
		return fmt.Errorf("failed to get categories: %w", err)
	}

	var targetCategory *model.DictCategory
	for _, cat := range categories {
		if cat.ID == id {
			targetCategory = cat
			break
		}
	}

	if targetCategory == nil {
		return fmt.Errorf("category with id %d not found", id)
	}

	// 检查分类下是否还有字典项
	items, err := s.dictRepo.GetItemsByCategory(ctx, targetCategory.Code)
	if err != nil {
		s.logger.Error("Failed to check category items", "category_code", targetCategory.Code, "error", err)
		return fmt.Errorf("failed to check category items: %w", err)
	}

	if len(items) > 0 {
		return fmt.Errorf("cannot delete category %s: it contains %d items", targetCategory.Name, len(items))
	}

	// 删除分类
	if err := s.dictRepo.DeleteCategory(ctx, id); err != nil {
		s.logger.Error("Failed to delete dict category", "id", id, "error", err)
		return fmt.Errorf("failed to delete dict category: %w", err)
	}

	// 清除分类缓存
	cacheKey := "dict_categories:all"
	if err := s.cache.Del(ctx, cacheKey); err != nil {
		s.logger.Error("Failed to clear categories cache", "error", err)
	}

	s.logger.Info("Dict category deleted successfully", "id", id, "code", targetCategory.Code, "name", targetCategory.Name)
	return nil
}

// CreateDictItem 创建字典项
func (s *dictService) CreateDictItem(ctx context.Context, req *CreateItemRequest) (*model.DictItem, error) {
	// 验证分类是否存在
	if _, err := s.dictRepo.GetCategoryByCode(ctx, req.CategoryCode); err != nil {
		return nil, fmt.Errorf("category with code %s does not exist", req.CategoryCode)
	}

	item := &model.DictItem{
		CategoryCode: req.CategoryCode,
		ItemKey:      req.ItemKey,
		ItemValue:    req.ItemValue,
		Description:  req.Description,
		SortOrder:    req.SortOrder,
		IsActive:     req.IsActive,
	}

	if err := s.dictRepo.CreateItem(ctx, item); err != nil {
		s.logger.Error("Failed to create dict item", "category_code", req.CategoryCode, "item_key", req.ItemKey, "error", err)
		return nil, fmt.Errorf("failed to create dict item: %w", err)
	}

	// 清除相关缓存
	s.clearCache(ctx, req.CategoryCode, req.ItemKey)

	s.logger.Info("Dict item created successfully", "category_code", req.CategoryCode, "item_key", req.ItemKey)
	return item, nil
}

// UpdateDictItem 更新字典项
func (s *dictService) UpdateDictItem(ctx context.Context, id uint, req *UpdateItemRequest) (*model.DictItem, error) {
	// 直接通过ID获取字典项
	targetItem, err := s.dictRepo.GetItemByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get dict item by ID", "id", id, "error", err)
		return nil, fmt.Errorf("dict item not found with id %d", id)
	}

	// 更新字段
	if req.ItemValue != "" {
		targetItem.ItemValue = req.ItemValue
	}
	targetItem.Description = req.Description
	targetItem.SortOrder = req.SortOrder
	if req.IsActive != nil {
		targetItem.IsActive = req.IsActive
	}

	if err := s.dictRepo.UpdateItem(ctx, targetItem); err != nil {
		s.logger.Error("Failed to update dict item", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update dict item: %w", err)
	}

	// 清除相关缓存
	s.clearCache(ctx, targetItem.CategoryCode, targetItem.ItemKey)

	s.logger.Info("Dict item updated successfully", "id", id, "category_code", targetItem.CategoryCode, "item_key", targetItem.ItemKey)
	return targetItem, nil
}

// DeleteDictItem 删除字典项
func (s *dictService) DeleteDictItem(ctx context.Context, id uint) error {
	// 直接通过ID获取字典项信息用于清除缓存
	targetItem, err := s.dictRepo.GetItemByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get dict item by ID", "id", id, "error", err)
		return fmt.Errorf("dict item not found with id %d", id)
	}

	if err := s.dictRepo.DeleteItem(ctx, id); err != nil {
		s.logger.Error("Failed to delete dict item", "id", id, "error", err)
		return fmt.Errorf("failed to delete dict item: %w", err)
	}

	// 清除相关缓存
	s.clearCache(ctx, targetItem.CategoryCode, targetItem.ItemKey)

	s.logger.Info("Dict item deleted successfully", "id", id, "category_code", targetItem.CategoryCode, "item_key", targetItem.ItemKey)
	return nil
}

// clearCache 清除相关缓存
func (s *dictService) clearCache(ctx context.Context, categoryCode, itemKey string) {
	// 清除分类缓存
	categoryKey := fmt.Sprintf("dict_items:%s", categoryCode)
	if err := s.cache.Del(ctx, categoryKey); err != nil {
		s.logger.Error("Failed to clear category cache", "key", categoryKey, "error", err)
	}

	// 清除特定项缓存
	if itemKey != "" {
		itemCacheKey := fmt.Sprintf("dict_item:%s:%s", categoryCode, itemKey)
		if err := s.cache.Del(ctx, itemCacheKey); err != nil {
			s.logger.Error("Failed to clear item cache", "key", itemCacheKey, "error", err)
		}
	}

	s.logger.Debug("Dict cache cleared", "category_code", categoryCode, "item_key", itemKey)
}

// boolPtr 创建bool指针的辅助函数
func boolPtr(b bool) *bool {
	return &b
}

// InitDefaultDictData 初始化默认字典数据
func (s *dictService) InitDefaultDictData(ctx context.Context) error {
	s.logger.Info("Starting to initialize default dictionary data")

	// 定义默认分类
	categories := []CreateCategoryRequest{
		{Code: "article_status", Name: "文章状态", Description: "文章的发布状态管理", SortOrder: 1}, // article_status
		{Code: "comment_status", Name: "评论状态", Description: "评论的审核状态管理", SortOrder: 2}, // comment_status
		{Code: "user_role", Name: "用户角色", Description: "用户权限角色管理", SortOrder: 3},       // user_role
		{Code: "user_status", Name: "用户状态", Description: "用户账户状态管理", SortOrder: 4},     // user_status
		{Code: "storage_type", Name: "存储类型", Description: "文件存储类型管理", SortOrder: 5},    // storage_type
	}

	// 创建分类（如果不存在）
	for _, categoryReq := range categories {
		if _, err := s.dictRepo.GetCategoryByCode(ctx, categoryReq.Code); err != nil {
			// 分类不存在，创建它
			if _, err := s.CreateDictCategory(ctx, &categoryReq); err != nil {
				s.logger.Error("Failed to create default category", "code", categoryReq.Code, "error", err)
				return fmt.Errorf("failed to create category %s: %w", categoryReq.Code, err)
			}
		}
	}

	// 定义默认字典项
	items := []CreateItemRequest{
		// article_status: draft(草稿), published(已发布), archived(已归档)
		{CategoryCode: "article_status", ItemKey: "draft", ItemValue: "草稿", Description: "文章草稿状态，未发布", SortOrder: 1, IsActive: boolPtr(true)},
		{CategoryCode: "article_status", ItemKey: "published", ItemValue: "已发布", Description: "文章已发布状态，对外可见", SortOrder: 2, IsActive: boolPtr(true)},
		{CategoryCode: "article_status", ItemKey: "archived", ItemValue: "已归档", Description: "文章已归档状态，不再显示", SortOrder: 3, IsActive: boolPtr(true)},

		// comment_status: pending(待审核), approved(已批准), rejected(已拒绝)
		{CategoryCode: "comment_status", ItemKey: "pending", ItemValue: "待审核", Description: "评论待审核状态，暂不显示", SortOrder: 1, IsActive: boolPtr(true)},
		{CategoryCode: "comment_status", ItemKey: "approved", ItemValue: "已批准", Description: "评论已批准状态，对外可见", SortOrder: 2, IsActive: boolPtr(true)},
		{CategoryCode: "comment_status", ItemKey: "rejected", ItemValue: "已拒绝", Description: "评论已拒绝状态，不予显示", SortOrder: 3, IsActive: boolPtr(true)},

		// user_role: admin(管理员), user(普通用户)
		{CategoryCode: "user_role", ItemKey: "admin", ItemValue: "管理员", Description: "系统管理员，拥有所有权限", SortOrder: 1, IsActive: boolPtr(true)},
		{CategoryCode: "user_role", ItemKey: "user", ItemValue: "普通用户", Description: "普通用户，拥有基本权限", SortOrder: 2, IsActive: boolPtr(true)},

		// user_status: active(活跃), inactive(非活跃), banned(已禁用)
		{CategoryCode: "user_status", ItemKey: "active", ItemValue: "活跃", Description: "用户账户正常活跃状态", SortOrder: 1, IsActive: boolPtr(true)},
		{CategoryCode: "user_status", ItemKey: "inactive", ItemValue: "非活跃", Description: "用户账户非活跃状态", SortOrder: 2, IsActive: boolPtr(true)},
		{CategoryCode: "user_status", ItemKey: "banned", ItemValue: "已禁用", Description: "用户账户已被禁用", SortOrder: 3, IsActive: boolPtr(true)},

		// storage_type: local(本地存储), s3(AWS S3), oss(阿里云OSS)
		{CategoryCode: "storage_type", ItemKey: "local", ItemValue: "本地存储", Description: "文件存储在本地服务器", SortOrder: 1, IsActive: boolPtr(true)},
		{CategoryCode: "storage_type", ItemKey: "s3", ItemValue: "AWS S3", Description: "文件存储在Amazon S3", SortOrder: 2, IsActive: boolPtr(true)},
		{CategoryCode: "storage_type", ItemKey: "oss", ItemValue: "阿里云OSS", Description: "文件存储在阿里云对象存储", SortOrder: 3, IsActive: boolPtr(true)},
	}

	// 创建字典项（如果不存在）
	for _, itemReq := range items {
		// 检查项是否已存在
		if _, err := s.GetDictItemByKey(ctx, itemReq.CategoryCode, itemReq.ItemKey); err != nil {
			// 项不存在，创建它
			if _, err := s.CreateDictItem(ctx, &itemReq); err != nil {
				s.logger.Error("Failed to create default dict item", "category_code", itemReq.CategoryCode, "item_key", itemReq.ItemKey, "error", err)
				return fmt.Errorf("failed to create dict item %s.%s: %w", itemReq.CategoryCode, itemReq.ItemKey, err)
			}
		}
	}

	s.logger.Info("Default dictionary data initialized successfully")
	return nil
}

// ClearDefaultDictData 清除默认字典数据
func (s *dictService) ClearDefaultDictData(ctx context.Context) error {
	s.logger.Info("Starting to clear default dictionary data")

	// 定义要删除的默认分类代码
	defaultCategories := []string{
		"article_status",
		"comment_status",
		"user_role",
		"user_status",
		"storage_type",
	}

	// 删除每个默认分类及其下的所有字典项
	for _, categoryCode := range defaultCategories {
		// 获取分类信息
		category, err := s.dictRepo.GetCategoryByCode(ctx, categoryCode)
		if err != nil {
			s.logger.Warn("Default category not found, skipping", "code", categoryCode)
			continue
		}

		// 获取该分类下的所有字典项
		items, err := s.dictRepo.GetItemsByCategory(ctx, categoryCode)
		if err != nil {
			s.logger.Error("Failed to get items for category", "category_code", categoryCode, "error", err)
			continue
		}

		// 删除所有字典项
		for _, item := range items {
			if err := s.dictRepo.DeleteItem(ctx, item.ID); err != nil {
				s.logger.Error("Failed to delete dict item", "id", item.ID, "category_code", categoryCode, "item_key", item.ItemKey, "error", err)
				return fmt.Errorf("failed to delete dict item %s.%s: %w", categoryCode, item.ItemKey, err)
			}
			s.logger.Debug("Dict item deleted", "id", item.ID, "category_code", categoryCode, "item_key", item.ItemKey)
		}

		// 删除分类
		if err := s.dictRepo.DeleteCategory(ctx, category.ID); err != nil {
			s.logger.Error("Failed to delete dict category", "id", category.ID, "code", categoryCode, "error", err)
			return fmt.Errorf("failed to delete dict category %s: %w", categoryCode, err)
		}
		s.logger.Debug("Dict category deleted", "id", category.ID, "code", categoryCode)
	}

	// 清除所有相关缓存
	cacheKeys := []string{
		"dict_categories:all",
	}

	// 清除分类缓存
	for _, key := range cacheKeys {
		if err := s.cache.Del(ctx, key); err != nil {
			s.logger.Error("Failed to clear cache", "key", key, "error", err)
		}
	}

	// 清除各分类的字典项缓存
	for _, categoryCode := range defaultCategories {
		categoryKey := fmt.Sprintf("dict_items:%s", categoryCode)
		if err := s.cache.Del(ctx, categoryKey); err != nil {
			s.logger.Error("Failed to clear category cache", "key", categoryKey, "error", err)
		}
	}

	s.logger.Info("Default dictionary data cleared successfully")
	return nil
}
