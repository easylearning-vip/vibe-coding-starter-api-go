package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/logger"
)

// AuthMiddleware JWT 认证中间件
type AuthMiddleware struct {
	config *config.Config
	cache  cache.Cache
	logger logger.Logger
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(
	config *config.Config,
	cache cache.Cache,
	logger logger.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		config: config,
		cache:  cache,
		logger: logger,
	}
}

// JWTClaims JWT 声明
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// RequireAuth 需要认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			m.logger.Warn("Missing authorization token", "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Authorization token required",
			})
			c.Abort()
			return
		}

		claims, err := m.validateToken(token)
		if err != nil {
			m.logger.Warn("Invalid token", "error", err, "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 检查 token 是否被撤销
		if m.isTokenRevoked(token) {
			m.logger.Warn("Revoked token used", "user_id", claims.UserID)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Token has been revoked",
			})
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("token", token)

		m.logger.Debug("User authenticated",
			"user_id", claims.UserID,
			"username", claims.Username,
			"path", c.Request.URL.Path)

		c.Next()
	}
}

// RequireRole 需要特定角色的中间件
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, requiredRole := range roles {
			if role == requiredRole || role == model.UserRoleAdmin {
				c.Next()
				return
			}
		}

		m.logger.Warn("Insufficient permissions",
			"user_role", role,
			"required_roles", roles,
			"path", c.Request.URL.Path)

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "Insufficient permissions",
		})
		c.Abort()
	}
}

// RequirePermission 需要特定权限的中间件
func (m *AuthMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		if !m.hasPermission(userRole.(string), permission) {
			m.logger.Warn("Permission denied",
				"user_role", userRole,
				"permission", permission,
				"path", c.Request.URL.Path)

			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Permission denied",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth 可选认证中间件
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := m.validateToken(token)
		if err != nil {
			// 可选认证，token 无效时不阻止请求
			c.Next()
			return
		}

		if !m.isTokenRevoked(token) {
			// 设置用户信息到上下文
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)
			c.Set("user_role", claims.Role)
			c.Set("token", token)
		}

		c.Next()
	}
}

// GenerateToken 生成 JWT token
func (m *AuthMiddleware) GenerateToken(user *model.User) (string, error) {
	now := time.Now()
	expirationTime := now.Add(time.Duration(m.config.JWT.Expiration) * time.Second)

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.config.JWT.Issuer,
			Subject:   user.Email,
			ID:        m.generateJTI(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.config.JWT.Secret))
	if err != nil {
		m.logger.Error("Failed to generate token", "error", err)
		return "", err
	}

	// 缓存 token 信息用于撤销检查
	m.cacheToken(tokenString, claims)

	return tokenString, nil
}

// RevokeToken 撤销 token
func (m *AuthMiddleware) RevokeToken(token string) error {
	claims, err := m.validateToken(token)
	if err != nil {
		return err
	}

	// 将 token 添加到黑名单
	blacklistKey := "token_blacklist:" + claims.ID
	expiration := time.Until(claims.ExpiresAt.Time)

	if err := m.cache.Set(context.Background(), blacklistKey, "revoked", expiration); err != nil {
		m.logger.Error("Failed to revoke token", "error", err)
		return err
	}

	m.logger.Info("Token revoked", "user_id", claims.UserID, "jti", claims.ID)
	return nil
}

// RefreshToken 刷新 token
func (m *AuthMiddleware) RefreshToken(token string) (string, error) {
	claims, err := m.validateToken(token)
	if err != nil {
		return "", err
	}

	// 检查 token 是否即将过期（在过期前 1 小时内可以刷新）
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", jwt.ErrTokenNotValidYet
	}

	// 创建新的 token
	user := &model.User{
		BaseModel: model.BaseModel{ID: claims.UserID},
		Username:  claims.Username,
		Email:     claims.Email,
		Role:      claims.Role,
	}

	newToken, err := m.GenerateToken(user)
	if err != nil {
		return "", err
	}

	// 撤销旧 token
	m.RevokeToken(token)

	return newToken, nil
}

// extractToken 从请求中提取 token
func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	// 从 Authorization header 提取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// 从查询参数提取
	if token := c.Query("token"); token != "" {
		return token
	}

	// 从 Cookie 提取
	if token, err := c.Cookie("access_token"); err == nil {
		return token
	}

	return ""
}

// validateToken 验证 token
func (m *AuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// isTokenRevoked 检查 token 是否被撤销
func (m *AuthMiddleware) isTokenRevoked(token string) bool {
	claims, err := m.validateToken(token)
	if err != nil {
		return true
	}

	blacklistKey := "token_blacklist:" + claims.ID
	exists, err := m.cache.Exists(context.Background(), blacklistKey)
	if err != nil {
		m.logger.Error("Failed to check token blacklist", "error", err)
		return false
	}

	return exists > 0
}

// hasPermission 检查用户是否有指定权限
func (m *AuthMiddleware) hasPermission(role, permission string) bool {
	// 管理员拥有所有权限
	if role == model.UserRoleAdmin {
		return true
	}

	// 这里可以实现更复杂的权限检查逻辑
	// 例如从数据库或配置文件中读取角色权限映射
	rolePermissions := map[string][]string{
		model.UserRoleUser: {
			"user:read",
			"user:update_self",
			"article:read",
			"file:read",
		},
		model.UserRoleAdmin: {
			"*", // 管理员拥有所有权限
		},
	}

	permissions, exists := rolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == "*" || p == permission {
			return true
		}
	}

	return false
}

// cacheToken 缓存 token 信息
func (m *AuthMiddleware) cacheToken(token string, claims *JWTClaims) {
	tokenKey := "token:" + claims.ID
	expiration := time.Until(claims.ExpiresAt.Time)

	tokenInfo := map[string]interface{}{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
		"role":     claims.Role,
		"token":    token,
	}

	if err := m.cache.Set(context.Background(), tokenKey, fmt.Sprintf("%v", tokenInfo), expiration); err != nil {
		m.logger.Error("Failed to cache token", "error", err)
	}
}

// generateJTI 生成 JWT ID
func (m *AuthMiddleware) generateJTI() string {
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), generateRandomString(8))
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
