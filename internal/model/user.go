package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID           string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username     string     `json:"username" gorm:"uniqueIndex;not null"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string     `json:"-" gorm:"not null"`                   // 密码哈希，不在JSON中返回
	Role         string     `json:"role" gorm:"not null;default:'user'"` // 用户角色：admin, editor, user
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Avatar       string     `json:"avatar"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName 返回用户表名
func (User) TableName() string {
	return "users"
}

// UserLogin 用户登录请求
type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserRegister 用户注册请求
type UserRegister struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// UserUpdate 用户更新请求
type UserUpdate struct {
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
	Role      string `json:"role,omitempty"`
}

// UserResponse 用户响应（不包含敏感信息）
type UserResponse struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Avatar    string     `json:"avatar"`
	IsActive  bool       `json:"is_active"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse 转换为用户响应对象
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Avatar:    u.Avatar,
		IsActive:  u.IsActive,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// PasswordResetRequest 密码重置请求
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordReset 密码重置
type PasswordReset struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// ChangePassword 修改密码
type ChangePassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// PasswordResetToken 密码重置令牌模型
type PasswordResetToken struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 返回密码重置令牌表名
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// LoginResponse 登录响应
type LoginResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
	TokenType    string        `json:"token_type"`
}
