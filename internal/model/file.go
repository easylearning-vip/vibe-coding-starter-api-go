package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// File 文件模型
type File struct {
	BaseModel
	Name          string `gorm:"size:255;not null" json:"name" validate:"required"`
	OriginalName  string `gorm:"size:255;not null" json:"original_name" validate:"required"`
	Path          string `gorm:"size:500;not null" json:"path" validate:"required"`
	URL           string `gorm:"size:500" json:"url"`
	Size          int64  `gorm:"not null" json:"size" validate:"required"`
	MimeType      string `gorm:"size:100;not null" json:"mime_type" validate:"required"`
	Extension     string `gorm:"size:10;not null" json:"extension" validate:"required"`
	Hash          string `gorm:"uniqueIndex;size:64;not null" json:"hash" validate:"required"`
	StorageType   string `gorm:"size:20;default:local" json:"storage_type" validate:"oneof=local s3 oss"`
	OwnerID       uint   `gorm:"not null" json:"owner_id" validate:"required"`
	Owner         User   `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	IsPublic      bool   `gorm:"default:false" json:"is_public"`
	DownloadCount int    `gorm:"default:0" json:"download_count"`
}

// StorageType 存储类型常量
const (
	StorageTypeLocal = "local"
	StorageTypeS3    = "s3"
	StorageTypeOSS   = "oss"
)

// TableName 获取表名
func (File) TableName() string {
	return "files"
}

// BeforeCreate GORM 钩子：创建前
func (f *File) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := f.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 设置默认存储类型
	if f.StorageType == "" {
		f.StorageType = StorageTypeLocal
	}

	// 从原始文件名提取扩展名
	if f.Extension == "" && f.OriginalName != "" {
		f.Extension = strings.ToLower(filepath.Ext(f.OriginalName))
		if f.Extension != "" && f.Extension[0] == '.' {
			f.Extension = f.Extension[1:] // 移除点号
		}
	}

	return nil
}

// IsImage 检查是否为图片文件
func (f *File) IsImage() bool {
	imageExtensions := []string{"jpg", "jpeg", "png", "gif", "bmp", "webp", "svg"}
	for _, ext := range imageExtensions {
		if f.Extension == ext {
			return true
		}
	}
	return false
}

// IsVideo 检查是否为视频文件
func (f *File) IsVideo() bool {
	videoExtensions := []string{"mp4", "avi", "mov", "wmv", "flv", "webm", "mkv"}
	for _, ext := range videoExtensions {
		if f.Extension == ext {
			return true
		}
	}
	return false
}

// IsAudio 检查是否为音频文件
func (f *File) IsAudio() bool {
	audioExtensions := []string{"mp3", "wav", "flac", "aac", "ogg", "wma"}
	for _, ext := range audioExtensions {
		if f.Extension == ext {
			return true
		}
	}
	return false
}

// IsDocument 检查是否为文档文件
func (f *File) IsDocument() bool {
	docExtensions := []string{"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "rtf"}
	for _, ext := range docExtensions {
		if f.Extension == ext {
			return true
		}
	}
	return false
}

// GetFileType 获取文件类型
func (f *File) GetFileType() string {
	switch {
	case f.IsImage():
		return "image"
	case f.IsVideo():
		return "video"
	case f.IsAudio():
		return "audio"
	case f.IsDocument():
		return "document"
	default:
		return "other"
	}
}

// GetSizeFormatted 获取格式化的文件大小
func (f *File) GetSizeFormatted() string {
	const unit = 1024
	if f.Size < unit {
		return fmt.Sprintf("%d B", f.Size)
	}
	div, exp := int64(unit), 0
	for n := f.Size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(f.Size)/float64(div), "KMGTPE"[exp])
}

// IncrementDownloadCount 增加下载次数
func (f *File) IncrementDownloadCount() {
	f.DownloadCount++
}

// GetFullURL 获取完整的文件 URL
func (f *File) GetFullURL(baseURL string) string {
	if f.URL != "" {
		// 如果已经是完整 URL，直接返回
		if strings.HasPrefix(f.URL, "http://") || strings.HasPrefix(f.URL, "https://") {
			return f.URL
		}
		// 拼接基础 URL
		return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(f.URL, "/")
	}
	return ""
}

// CanAccess 检查用户是否可以访问文件
func (f *File) CanAccess(userID uint) bool {
	// 公开文件任何人都可以访问
	if f.IsPublic {
		return true
	}
	// 文件所有者可以访问
	return f.OwnerID == userID
}
