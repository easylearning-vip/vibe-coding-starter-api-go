package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/logger"
)

// fileService 文件服务实现
type fileService struct {
	fileRepo repository.FileRepository
	cache    cache.Cache
	logger   logger.Logger
}

// NewFileService 创建文件服务
func NewFileService(
	fileRepo repository.FileRepository,
	cache cache.Cache,
	logger logger.Logger,
) FileService {
	return &fileService{
		fileRepo: fileRepo,
		cache:    cache,
		logger:   logger,
	}
}

// Upload 上传文件
func (s *fileService) Upload(ctx context.Context, req *UploadRequest) (*model.File, error) {
	// 计算文件哈希
	hash := s.calculateHash(req.FileData)

	// 检查文件是否已存在
	existingFile, err := s.fileRepo.GetByHash(ctx, hash)
	if err == nil && existingFile != nil {
		s.logger.Info("File already exists, returning existing file", "hash", hash, "file_id", existingFile.ID)
		return existingFile, nil
	}

	// 生成唯一文件名
	fileName := s.generateFileName(req.FileName)

	// 设置存储类型
	storageType := req.StorageType
	if storageType == "" {
		storageType = model.StorageTypeLocal
	}

	// 保存文件到存储
	filePath, fileURL, err := s.saveFile(req.FileData, fileName, storageType)
	if err != nil {
		s.logger.Error("Failed to save file", "file_name", req.FileName, "error", err)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 创建文件记录
	file := &model.File{
		Name:         fileName,
		OriginalName: req.FileName,
		Path:         filePath,
		URL:          fileURL,
		Size:         req.FileSize,
		MimeType:     req.MimeType,
		Hash:         hash,
		StorageType:  storageType,
		IsPublic:     req.IsPublic,
	}

	if err := s.fileRepo.Create(ctx, file); err != nil {
		// 如果数据库保存失败，删除已保存的文件
		s.deletePhysicalFile(filePath)
		s.logger.Error("Failed to create file record", "file_name", req.FileName, "error", err)
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	s.logger.Info("File uploaded successfully", "file_id", file.ID, "file_name", file.Name)
	return file, nil
}

// GetByID 根据 ID 获取文件
func (s *fileService) GetByID(ctx context.Context, id uint) (*model.File, error) {
	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get file by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return file, nil
}

// Delete 删除文件
func (s *fileService) Delete(ctx context.Context, id uint) error {
	// 获取文件信息
	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get file for deletion", "id", id, "error", err)
		return fmt.Errorf("failed to get file: %w", err)
	}

	// 删除数据库记录
	if err := s.fileRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete file record", "id", id, "error", err)
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	// 删除物理文件
	if err := s.deletePhysicalFile(file.Path); err != nil {
		s.logger.Warn("Failed to delete physical file", "path", file.Path, "error", err)
		// 不返回错误，因为数据库记录已删除
	}

	s.logger.Info("File deleted successfully", "file_id", id)
	return nil
}

// List 获取文件列表
func (s *fileService) List(ctx context.Context, opts repository.ListOptions) ([]*model.File, int64, error) {
	files, total, err := s.fileRepo.List(ctx, opts)
	if err != nil {
		s.logger.Error("Failed to get files list", "error", err)
		return nil, 0, fmt.Errorf("failed to get files: %w", err)
	}

	return files, total, nil
}

// GetByOwner 根据所有者获取文件列表
func (s *fileService) GetByOwner(ctx context.Context, ownerID uint, opts repository.ListOptions) ([]*model.File, int64, error) {
	files, total, err := s.fileRepo.GetByOwner(ctx, ownerID, opts)
	if err != nil {
		s.logger.Error("Failed to get files by owner", "owner_id", ownerID, "error", err)
		return nil, 0, fmt.Errorf("failed to get files: %w", err)
	}

	return files, total, nil
}

// Download 下载文件
func (s *fileService) Download(ctx context.Context, id uint) (*DownloadResponse, error) {
	// 获取文件信息
	file, err := s.fileRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get file for download", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	// 读取文件数据
	fileData, err := s.readFile(file.Path)
	if err != nil {
		s.logger.Error("Failed to read file data", "path", file.Path, "error", err)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 增加下载次数
	file.IncrementDownloadCount()
	if err := s.fileRepo.Update(ctx, file); err != nil {
		s.logger.Warn("Failed to update download count", "file_id", id, "error", err)
		// 不返回错误，因为文件下载成功
	}

	return &DownloadResponse{
		File:     file,
		FileData: fileData,
		URL:      file.URL,
	}, nil
}

// calculateHash 计算文件哈希
func (s *fileService) calculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// generateFileName 生成唯一文件名
func (s *fileService) generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d_%s%s", timestamp, generateRandomString(8), ext)
}

// saveFile 保存文件到存储
func (s *fileService) saveFile(data []byte, fileName, storageType string) (string, string, error) {
	switch storageType {
	case model.StorageTypeLocal:
		return s.saveToLocal(data, fileName)
	case model.StorageTypeS3:
		return s.saveToS3(data, fileName)
	case model.StorageTypeOSS:
		return s.saveToOSS(data, fileName)
	default:
		return "", "", fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// saveToLocal 保存到本地存储
func (s *fileService) saveToLocal(data []byte, fileName string) (string, string, error) {
	// 创建上传目录
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 文件路径
	filePath := filepath.Join(uploadDir, fileName)

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", "", fmt.Errorf("failed to write file: %w", err)
	}

	// 生成 URL
	fileURL := fmt.Sprintf("/uploads/%s", fileName)

	return filePath, fileURL, nil
}

// saveToS3 保存到 S3 存储 (占位符实现)
func (s *fileService) saveToS3(data []byte, fileName string) (string, string, error) {
	// TODO: 实现 S3 上传逻辑
	return "", "", fmt.Errorf("S3 storage not implemented yet")
}

// saveToOSS 保存到阿里云 OSS 存储 (占位符实现)
func (s *fileService) saveToOSS(data []byte, fileName string) (string, string, error) {
	// TODO: 实现 OSS 上传逻辑
	return "", "", fmt.Errorf("OSS storage not implemented yet")
}

// readFile 读取文件数据
func (s *fileService) readFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// deletePhysicalFile 删除物理文件
func (s *fileService) deletePhysicalFile(filePath string) error {
	return os.Remove(filePath)
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
