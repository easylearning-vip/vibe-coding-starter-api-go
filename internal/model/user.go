package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	BaseModel
	Username  string     `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,min=3,max=50"`
	Email     string     `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email"`
	Password  string     `gorm:"size:255;not null" json:"-" validate:"required,min=6"`
	Nickname  string     `gorm:"size:50" json:"nickname" validate:"max=50"`
	Avatar    string     `gorm:"size:255" json:"avatar" validate:"url"`
	Role      string     `gorm:"size:20;default:user" json:"role" validate:"oneof=admin user"`
	Status    string     `gorm:"size:20;default:active" json:"status" validate:"oneof=active inactive banned"`
	LastLogin *time.Time `json:"last_login"`
	Articles  []Article  `gorm:"foreignKey:AuthorID" json:"articles,omitempty"`
}

// UserRole 用户角色常量
const (
	UserRoleAdmin = "admin"
	UserRoleUser  = "user"
)

// UserStatus 用户状态常量
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusBanned   = "banned"
)

// TableName 获取表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM 钩子：创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := u.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 设置默认值
	if u.Role == "" {
		u.Role = UserRoleUser
	}
	if u.Status == "" {
		u.Status = UserStatusActive
	}

	// 加密密码
	if u.Password != "" {
		hashedPassword, err := u.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}

	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := u.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// 如果密码被修改，重新加密
	if tx.Statement.Changed("Password") && u.Password != "" {
		hashedPassword, err := u.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}

	return nil
}

// HashPassword 加密密码
func (u *User) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsAdmin 检查是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsActive 检查是否为活跃状态
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsBanned 检查是否被禁用
func (u *User) IsBanned() bool {
	return u.Status == UserStatusBanned
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
}

// ToPublic 转换为公开信息（不包含敏感信息）
func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		LastLogin: u.LastLogin,
	}
}

// PublicUser 公开用户信息
type PublicUser struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Nickname  string     `json:"nickname"`
	Avatar    string     `json:"avatar"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastLogin *time.Time `json:"last_login"`
}
