package service

import (
	"context"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetProfile(ctx context.Context, userID uint) (*model.User, error)
	UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) (*model.User, error)
	ChangePassword(ctx context.Context, userID uint, req *ChangePasswordRequest) error
	GetUsers(ctx context.Context, opts repository.ListOptions) ([]*model.User, int64, error)
	DeleteUser(ctx context.Context, userID uint) error
}

// ArticleService 文章服务接口
type ArticleService interface {
	Create(ctx context.Context, req *CreateArticleRequest) (*model.Article, error)
	GetByID(ctx context.Context, id uint) (*model.Article, error)
	GetBySlug(ctx context.Context, slug string) (*model.Article, error)
	Update(ctx context.Context, id uint, req *UpdateArticleRequest) (*model.Article, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts repository.ListOptions) ([]*model.Article, int64, error)
	GetPublished(ctx context.Context, opts repository.ListOptions) ([]*model.Article, int64, error)
	Search(ctx context.Context, query string, opts repository.ListOptions) ([]*model.Article, int64, error)
	IncrementViewCount(ctx context.Context, articleID uint) error
}

// FileService 文件服务接口
type FileService interface {
	Upload(ctx context.Context, req *UploadRequest) (*model.File, error)
	GetByID(ctx context.Context, id uint) (*model.File, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, opts repository.ListOptions) ([]*model.File, int64, error)
	GetByOwner(ctx context.Context, ownerID uint, opts repository.ListOptions) ([]*model.File, int64, error)
	Download(ctx context.Context, id uint) (*DownloadResponse, error)
}

// DictService 数据字典服务接口
type DictService interface {
	GetDictCategories(ctx context.Context) ([]*model.DictCategory, error)
	GetDictItems(ctx context.Context, categoryCode string) ([]*model.DictItem, error)
	GetDictItemByKey(ctx context.Context, categoryCode, itemKey string) (*model.DictItem, error)
	CreateDictCategory(ctx context.Context, req *CreateCategoryRequest) (*model.DictCategory, error)
	DeleteDictCategory(ctx context.Context, id uint) error
	CreateDictItem(ctx context.Context, req *CreateItemRequest) (*model.DictItem, error)
	UpdateDictItem(ctx context.Context, id uint, req *UpdateItemRequest) (*model.DictItem, error)
	DeleteDictItem(ctx context.Context, id uint) error
	InitDefaultDictData(ctx context.Context) error
	ClearDefaultDictData(ctx context.Context) error
}

// 请求和响应结构体

// 用户相关
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname" validate:"max=50"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User  *model.PublicUser `json:"user"`
	Token string            `json:"token"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" validate:"min=3,max=50"`
	Nickname string `json:"nickname" validate:"max=50"`
	Avatar   string `json:"avatar" validate:"url"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// 文章相关
type CreateArticleRequest struct {
	Title      string `json:"title" validate:"required,max=200"`
	Content    string `json:"content" validate:"required"`
	Summary    string `json:"summary" validate:"max=500"`
	CoverImage string `json:"cover_image" validate:"url"`
	CategoryID *uint  `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
	Status     string `json:"status" validate:"oneof=draft published"`
	AuthorID   uint   `json:"author_id,omitempty"` // 作者ID，由服务器设置
}

type UpdateArticleRequest struct {
	Title      string `json:"title" validate:"max=200"`
	Content    string `json:"content"`
	Summary    string `json:"summary" validate:"max=500"`
	CoverImage string `json:"cover_image" validate:"url"`
	CategoryID *uint  `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
	Status     string `json:"status" validate:"oneof=draft published archived"`
}

// 文件相关
type UploadRequest struct {
	FileName    string `json:"file_name" validate:"required"`
	FileSize    int64  `json:"file_size" validate:"required"`
	MimeType    string `json:"mime_type" validate:"required"`
	FileData    []byte `json:"file_data" validate:"required"`
	IsPublic    bool   `json:"is_public"`
	StorageType string `json:"storage_type" validate:"oneof=local s3 oss"`
}

type DownloadResponse struct {
	File     *model.File `json:"file"`
	FileData []byte      `json:"file_data"`
	URL      string      `json:"url"`
}

// 数据字典相关
type CreateCategoryRequest struct {
	Code        string `json:"code" validate:"required,max=50"`
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type CreateItemRequest struct {
	CategoryCode string `json:"category_code" validate:"required,max=50"`
	ItemKey      string `json:"item_key" validate:"required,max=50"`
	ItemValue    string `json:"item_value" validate:"required,max=200"`
	Description  string `json:"description"`
	SortOrder    int    `json:"sort_order"`
	IsActive     *bool  `json:"is_active"`
}

type UpdateItemRequest struct {
	ItemValue   string `json:"item_value" validate:"max=200"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	IsActive    *bool  `json:"is_active"`
}
