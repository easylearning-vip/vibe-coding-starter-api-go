package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/logger"
)

// Permission 权限定义
type Permission struct {
	Resource string `json:"resource"` // 资源名称，如 user, article, file
	Action   string `json:"action"`   // 操作名称，如 create, read, update, delete
	Scope    string `json:"scope"`    // 权限范围，如 own, all
}

// String 返回权限的字符串表示
func (p Permission) String() string {
	if p.Scope != "" {
		return fmt.Sprintf("%s:%s:%s", p.Resource, p.Action, p.Scope)
	}
	return fmt.Sprintf("%s:%s", p.Resource, p.Action)
}

// PermissionMiddleware 权限控制中间件
type PermissionMiddleware struct {
	config *config.Config
	cache  cache.Cache
	logger logger.Logger
}

// NewPermissionMiddleware 创建权限控制中间件
func NewPermissionMiddleware(
	config *config.Config,
	cache cache.Cache,
	logger logger.Logger,
) *PermissionMiddleware {
	return &PermissionMiddleware{
		config: config,
		cache:  cache,
		logger: logger,
	}
}

// RequirePermissions 需要指定权限的中间件
func (m *PermissionMiddleware) RequirePermissions(permissions ...Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		userRole, _ := c.Get("user_role")
		role := userRole.(string)

		// 检查用户是否有所需权限
		for _, permission := range permissions {
			if !m.hasPermission(userID.(uint), role, permission, c) {
				m.logger.Warn("Permission denied",
					"user_id", userID,
					"permission", permission.String(),
					"path", c.Request.URL.Path)

				c.JSON(http.StatusForbidden, gin.H{
					"error":   "forbidden",
					"message": fmt.Sprintf("Permission denied: %s", permission.String()),
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireAnyPermission 需要任一权限的中间件
func (m *PermissionMiddleware) RequireAnyPermission(permissions ...Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		userRole, _ := c.Get("user_role")
		role := userRole.(string)

		// 检查用户是否有任一权限
		for _, permission := range permissions {
			if m.hasPermission(userID.(uint), role, permission, c) {
				c.Next()
				return
			}
		}

		m.logger.Warn("No required permissions",
			"user_id", userID,
			"permissions", permissions,
			"path", c.Request.URL.Path)

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "Insufficient permissions",
		})
		c.Abort()
	}
}

// RequireOwnership 需要资源所有权的中间件
func (m *PermissionMiddleware) RequireOwnership(resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		userRole, _ := c.Get("user_role")
		role := userRole.(string)

		// 管理员跳过所有权检查
		if role == model.UserRoleAdmin {
			c.Next()
			return
		}

		// 检查资源所有权
		if !m.checkOwnership(userID.(uint), resourceType, c) {
			m.logger.Warn("Ownership check failed",
				"user_id", userID,
				"resource_type", resourceType,
				"path", c.Request.URL.Path)

			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "You can only access your own resources",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasPermission 检查用户是否有指定权限
func (m *PermissionMiddleware) hasPermission(userID uint, role string, permission Permission, c *gin.Context) bool {
	// 管理员拥有所有权限
	if role == model.UserRoleAdmin {
		return true
	}

	// 从缓存获取用户权限
	userPermissions := m.getUserPermissions(userID, role)

	// 检查精确匹配
	for _, p := range userPermissions {
		if p.String() == permission.String() {
			return true
		}
	}

	// 检查通配符权限
	for _, p := range userPermissions {
		if m.matchesWildcard(p, permission) {
			return true
		}
	}

	// 检查基于上下文的权限（如资源所有权）
	if permission.Scope == "own" {
		return m.checkOwnership(userID, permission.Resource, c)
	}

	return false
}

// getUserPermissions 获取用户权限列表
func (m *PermissionMiddleware) getUserPermissions(userID uint, role string) []Permission {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("user_permissions:%d", userID)
	if cached, err := m.cache.Get(context.Background(), cacheKey); err == nil {
		var permissions []Permission
		if err := json.Unmarshal([]byte(cached), &permissions); err == nil {
			return permissions
		}
	}

	// 获取角色默认权限
	permissions := m.getRolePermissions(role)

	// 这里可以添加用户特定权限的逻辑
	// 例如从数据库查询用户的额外权限

	// 缓存权限列表
	if data, err := json.Marshal(permissions); err == nil {
		m.cache.Set(context.Background(), cacheKey, string(data), time.Hour)
	}

	return permissions
}

// getRolePermissions 获取角色权限
func (m *PermissionMiddleware) getRolePermissions(role string) []Permission {
	rolePermissions := map[string][]Permission{
		model.UserRoleUser: {
			{Resource: "user", Action: "read", Scope: "own"},
			{Resource: "user", Action: "update", Scope: "own"},
			{Resource: "article", Action: "read"},
			{Resource: "article", Action: "create"},
			{Resource: "article", Action: "update", Scope: "own"},
			{Resource: "article", Action: "delete", Scope: "own"},
			{Resource: "file", Action: "read"},
			{Resource: "file", Action: "upload"},
			{Resource: "file", Action: "delete", Scope: "own"},
		},
		model.UserRoleAdmin: {
			{Resource: "*", Action: "*"}, // 管理员拥有所有权限
		},
	}

	if permissions, exists := rolePermissions[role]; exists {
		return permissions
	}

	return []Permission{}
}

// matchesWildcard 检查权限是否匹配通配符
func (m *PermissionMiddleware) matchesWildcard(userPerm, requiredPerm Permission) bool {
	// 检查资源通配符
	if userPerm.Resource == "*" {
		return true
	}

	// 检查操作通配符
	if userPerm.Resource == requiredPerm.Resource && userPerm.Action == "*" {
		return true
	}

	// 检查范围通配符
	if userPerm.Resource == requiredPerm.Resource &&
		userPerm.Action == requiredPerm.Action &&
		userPerm.Scope == "*" {
		return true
	}

	return false
}

// checkOwnership 检查资源所有权
func (m *PermissionMiddleware) checkOwnership(userID uint, resourceType string, c *gin.Context) bool {
	switch resourceType {
	case "user":
		return m.checkUserOwnership(userID, c)
	case "article":
		return m.checkArticleOwnership(userID, c)
	case "file":
		return m.checkFileOwnership(userID, c)
	default:
		return false
	}
}

// checkUserOwnership 检查用户资源所有权
func (m *PermissionMiddleware) checkUserOwnership(userID uint, c *gin.Context) bool {
	// 从路径参数获取目标用户 ID
	targetUserID := c.Param("id")
	if targetUserID == "" {
		targetUserID = c.Param("user_id")
	}

	// 如果没有指定目标用户，默认为当前用户
	if targetUserID == "" {
		return true
	}

	// 检查是否为同一用户
	if targetUserID == fmt.Sprintf("%d", userID) {
		return true
	}

	return false
}

// checkArticleOwnership 检查文章资源所有权
func (m *PermissionMiddleware) checkArticleOwnership(userID uint, c *gin.Context) bool {
	articleID := c.Param("id")
	if articleID == "" {
		return false
	}

	// 这里应该查询数据库检查文章所有者
	// 为了示例，我们使用缓存检查
	cacheKey := fmt.Sprintf("article_owner:%s", articleID)
	if ownerID, err := m.cache.Get(context.Background(), cacheKey); err == nil {
		return ownerID == fmt.Sprintf("%d", userID)
	}

	// 如果缓存中没有，应该查询数据库
	// 这里返回 false 作为安全默认值
	return false
}

// checkFileOwnership 检查文件资源所有权
func (m *PermissionMiddleware) checkFileOwnership(userID uint, c *gin.Context) bool {
	fileID := c.Param("id")
	if fileID == "" {
		return false
	}

	// 这里应该查询数据库检查文件所有者
	// 为了示例，我们使用缓存检查
	cacheKey := fmt.Sprintf("file_owner:%s", fileID)
	if ownerID, err := m.cache.Get(context.Background(), cacheKey); err == nil {
		return ownerID == fmt.Sprintf("%d", userID)
	}

	// 如果缓存中没有，应该查询数据库
	// 这里返回 false 作为安全默认值
	return false
}

// ClearUserPermissions 清除用户权限缓存
func (m *PermissionMiddleware) ClearUserPermissions(userID uint) error {
	cacheKey := fmt.Sprintf("user_permissions:%d", userID)
	return m.cache.Del(context.Background(), cacheKey)
}

// 预定义的权限常量
var (
	// 用户权限
	PermUserRead      = Permission{Resource: "user", Action: "read"}
	PermUserCreate    = Permission{Resource: "user", Action: "create"}
	PermUserUpdate    = Permission{Resource: "user", Action: "update"}
	PermUserUpdateOwn = Permission{Resource: "user", Action: "update", Scope: "own"}
	PermUserDelete    = Permission{Resource: "user", Action: "delete"}
	PermUserDeleteOwn = Permission{Resource: "user", Action: "delete", Scope: "own"}

	// 文章权限
	PermArticleRead      = Permission{Resource: "article", Action: "read"}
	PermArticleCreate    = Permission{Resource: "article", Action: "create"}
	PermArticleUpdate    = Permission{Resource: "article", Action: "update"}
	PermArticleUpdateOwn = Permission{Resource: "article", Action: "update", Scope: "own"}
	PermArticleDelete    = Permission{Resource: "article", Action: "delete"}
	PermArticleDeleteOwn = Permission{Resource: "article", Action: "delete", Scope: "own"}

	// 文件权限
	PermFileRead      = Permission{Resource: "file", Action: "read"}
	PermFileUpload    = Permission{Resource: "file", Action: "upload"}
	PermFileDelete    = Permission{Resource: "file", Action: "delete"}
	PermFileDeleteOwn = Permission{Resource: "file", Action: "delete", Scope: "own"}

	// 管理权限
	PermAdminAll = Permission{Resource: "*", Action: "*"}
)

// 便捷的权限检查函数
func RequireUserRead() gin.HandlerFunc {
	return NewPermissionMiddleware(nil, nil, nil).RequirePermissions(PermUserRead)
}

func RequireUserUpdate() gin.HandlerFunc {
	return NewPermissionMiddleware(nil, nil, nil).RequirePermissions(PermUserUpdate)
}

func RequireUserUpdateOwn() gin.HandlerFunc {
	return NewPermissionMiddleware(nil, nil, nil).RequirePermissions(PermUserUpdateOwn)
}

func RequireArticleCreate() gin.HandlerFunc {
	return NewPermissionMiddleware(nil, nil, nil).RequirePermissions(PermArticleCreate)
}

func RequireArticleUpdateOwn() gin.HandlerFunc {
	return NewPermissionMiddleware(nil, nil, nil).RequirePermissions(PermArticleUpdateOwn)
}

func RequireFileUpload() gin.HandlerFunc {
	return NewPermissionMiddleware(nil, nil, nil).RequirePermissions(PermFileUpload)
}
