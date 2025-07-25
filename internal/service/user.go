package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"vibe-coding-starter/internal/config"
	"vibe-coding-starter/internal/model"
	"vibe-coding-starter/internal/repository"
	"vibe-coding-starter/pkg/cache"
	"vibe-coding-starter/pkg/logger"
)

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
	cache    cache.Cache
	logger   logger.Logger
	config   *config.Config
}

// NewUserService 创建用户服务
func NewUserService(
	userRepo repository.UserRepository,
	cache cache.Cache,
	logger logger.Logger,
	config *config.Config,
) UserService {
	return &userService{
		userRepo: userRepo,
		cache:    cache,
		logger:   logger,
		config:   config,
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req *RegisterRequest) (*model.User, error) {
	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.logger.Error("Failed to check existing user by email", "email", req.Email, "error", err)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// 检查用户名是否已存在
	existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.logger.Error("Failed to check existing user by username", "username", req.Username, "error", err)
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", req.Username)
	}

	// 创建新用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 密码会在 BeforeCreate 钩子中加密
		Nickname: req.Nickname,
		Role:     model.UserRoleUser,
		Status:   model.UserStatusActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", "email", req.Email, "error", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("User registered successfully", "user_id", user.ID, "email", user.Email)
	return user, nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	s.logger.Debug("Login attempt", "username", req.Username)

	// 获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		s.logger.Error("Failed to get user by username", "username", req.Username, "error", err)
		return nil, fmt.Errorf("invalid username or password")
	}

	s.logger.Debug("User found", "user_id", user.ID, "username", user.Username, "status", user.Status)

	// 检查用户状态
	if !user.IsActive() {
		s.logger.Warn("User account is not active", "user_id", user.ID, "username", user.Username, "status", user.Status)
		return nil, fmt.Errorf("user account is not active")
	}

	s.logger.Debug("Checking password", "user_id", user.ID)

	// 验证密码
	if !user.CheckPassword(req.Password) {
		s.logger.Warn("Invalid password attempt", "user_id", user.ID, "username", user.Username)
		return nil, fmt.Errorf("invalid username or password")
	}

	s.logger.Debug("Password verified successfully", "user_id", user.ID)

	// 生成 JWT Token
	token, err := s.generateJWTToken(user)
	if err != nil {
		s.logger.Error("Failed to generate JWT token", "user_id", user.ID, "error", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.logger.Error("Failed to update last login", "user_id", user.ID, "error", err)
		// 不返回错误，因为这不是关键操作
	}

	s.logger.Info("User logged in successfully", "user_id", user.ID, "email", user.Email)

	return &LoginResponse{
		User:  user.ToPublic(),
		Token: token,
	}, nil
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user profile", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return user, nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) (*model.User, error) {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user for update", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 检查用户名是否已被其他用户使用
	if req.Username != "" && req.Username != user.Username {
		existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
		if err != nil && err != gorm.ErrRecordNotFound {
			s.logger.Error("Failed to check existing username", "username", req.Username, "error", err)
			return nil, fmt.Errorf("failed to check username: %w", err)
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, fmt.Errorf("username %s is already taken", req.Username)
		}
		user.Username = req.Username
	}

	// 更新其他字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("Failed to update user profile", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	s.logger.Info("User profile updated successfully", "user_id", userID)
	return user, nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID uint, req *ChangePasswordRequest) error {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user for password change", "user_id", userID, "error", err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		s.logger.Warn("Invalid old password attempt", "user_id", userID)
		return fmt.Errorf("invalid old password")
	}

	// 设置新密码
	user.Password = req.NewPassword // 密码会在 BeforeUpdate 钩子中加密

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("Failed to update user password", "user_id", userID, "error", err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.logger.Info("User password changed successfully", "user_id", userID)
	return nil
}

// GetUsers 获取用户列表
func (s *userService) GetUsers(ctx context.Context, opts repository.ListOptions) ([]*model.User, int64, error) {
	users, total, err := s.userRepo.List(ctx, opts)
	if err != nil {
		s.logger.Error("Failed to get users list", "error", err)
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, userID uint) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user for deletion", "user_id", userID, "error", err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 删除用户
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.Error("Failed to delete user", "user_id", userID, "error", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Info("User deleted successfully", "user_id", userID)
	return nil
}

// generateJWTToken 生成 JWT Token
func (s *userService) generateJWTToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Duration(s.config.JWT.Expiration) * time.Second).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      s.config.JWT.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}
