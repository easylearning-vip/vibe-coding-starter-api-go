package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/internal/service"
	"vibe-coding-starter/pkg/logger"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileService service.FileService
	logger      logger.Logger
}

// NewFileHandler 创建文件处理器
func NewFileHandler(
	fileService service.FileService,
	logger logger.Logger,
) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		logger:      logger,
	}
}

// Upload 上传文件
// @Summary 上传文件
// @Description 上传文件到服务器
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "文件"
// @Param is_public formData bool false "是否公开" default(false)
// @Success 201 {object} model.File
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/files/upload [post]
func (h *FileHandler) Upload(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get uploaded file", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "file_required",
			Message: "File is required",
		})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		h.logger.Error("Failed to open uploaded file", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "file_open_failed",
			Message: "Failed to open file",
		})
		return
	}
	defer src.Close()

	// 读取文件数据
	fileData := make([]byte, file.Size)
	if _, err := src.Read(fileData); err != nil {
		h.logger.Error("Failed to read file data", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "file_read_failed",
			Message: "Failed to read file",
		})
		return
	}

	// 获取其他参数
	isPublic := c.DefaultPostForm("is_public", "false") == "true"

	// 创建上传请求
	req := &service.UploadRequest{
		FileName:    file.Filename,
		FileSize:    file.Size,
		MimeType:    file.Header.Get("Content-Type"),
		FileData:    fileData,
		IsPublic:    isPublic,
		StorageType: "local", // 默认本地存储
	}

	uploadedFile, err := h.fileService.Upload(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to upload file", "filename", file.Filename, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "upload_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, uploadedFile)
}

// GetByID 根据 ID 获取文件信息
// @Summary 获取文件信息
// @Description 根据 ID 获取文件信息
// @Tags files
// @Accept json
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} model.File
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/files/{id} [get]
func (h *FileHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid file ID",
		})
		return
	}

	file, err := h.fileService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get file", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "file_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, file)
}

// Download 下载文件
// @Summary 下载文件
// @Description 下载指定文件
// @Tags files
// @Accept json
// @Produce application/octet-stream
// @Param id path int true "文件ID"
// @Success 200 {file} binary
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/files/{id}/download [get]
func (h *FileHandler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid file ID",
		})
		return
	}

	response, err := h.fileService.Download(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to download file", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "download_failed",
			Message: err.Error(),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+response.File.OriginalName)
	c.Header("Content-Type", response.File.MimeType)
	c.Header("Content-Length", strconv.FormatInt(response.File.Size, 10))

	// 返回文件数据
	c.Data(http.StatusOK, response.File.MimeType, response.FileData)
}

// List 获取文件列表
// @Summary 获取文件列表
// @Description 获取文件列表
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} ListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/files [get]
func (h *FileHandler) List(c *gin.Context) {
	opts := h.parseListOptions(c)
	files, total, err := h.fileService.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("Failed to get files", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "get_files_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Data:  files,
		Total: total,
		Page:  opts.Page,
		Size:  opts.PageSize,
	})
}

// Delete 删除文件
// @Summary 删除文件
// @Description 删除指定文件
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文件ID"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/files/{id} [delete]
func (h *FileHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid file ID",
		})
		return
	}

	if err := h.fileService.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("Failed to delete file", "id", id, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "File deleted successfully",
	})
}

// RegisterRoutes 注册路由
func (h *FileHandler) RegisterRoutes(r *gin.RouterGroup) {
	files := r.Group("/files")
	{
		files.GET("/:id", h.GetByID)
		files.GET("/:id/download", h.Download)

		// 需要认证的路由
		authenticated := files.Group("")
		authenticated.Use(AuthMiddleware())
		{
			authenticated.POST("/upload", h.Upload)
			authenticated.GET("", h.List)
			authenticated.DELETE("/:id", h.Delete)
		}
	}
}

// 辅助方法
func (h *FileHandler) parseListOptions(c *gin.Context) repository.ListOptions {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	return repository.ListOptions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		Sort:     c.DefaultQuery("sort", "created_at"),
		Order:    c.DefaultQuery("order", "desc"),
	}
}
