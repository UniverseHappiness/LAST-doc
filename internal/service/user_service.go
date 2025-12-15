package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/repository"
)

// UserService 用户服务接口
type UserService interface {
	// 用户注册
	Register(ctx context.Context, req *model.UserRegister) (*model.UserResponse, error)

	// 用户登录
	Login(ctx context.Context, req *model.UserLogin) (*model.LoginResponse, error)

	// 根据ID获取用户
	GetByID(ctx context.Context, id string) (*model.UserResponse, error)

	// 更新用户信息
	Update(ctx context.Context, id string, req *model.UserUpdate) (*model.UserResponse, error)

	// 删除用户
	Delete(ctx context.Context, id string) error

	// 获取用户列表
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.UserResponse, int64, error)

	// 修改密码
	ChangePassword(ctx context.Context, userID string, req *model.ChangePassword) error

	// 重置密码
	ResetPassword(ctx context.Context, req *model.PasswordReset) error

	// 请求密码重置
	RequestPasswordReset(ctx context.Context, req *model.PasswordResetRequest) error

	// 验证JWT令牌
	ValidateJWT(tokenString string) (*model.JWTClaims, error)

	// 生成JWT令牌
	GenerateJWT(user *model.User) (string, error)

	// 刷新JWT令牌
	RefreshJWT(ctx context.Context, tokenString string) (string, error)
}

// userService 用户服务实现
type userService struct {
	userRepo               repository.UserRepository
	passwordResetTokenRepo repository.PasswordResetTokenRepository
	db                     *gorm.DB
	jwtSecret              string
	jwtExpiration          time.Duration
	refreshTokenExpiration time.Duration
}

// NewUserService 创建用户服务实例
func NewUserService(
	userRepo repository.UserRepository,
	passwordResetTokenRepo repository.PasswordResetTokenRepository,
	db *gorm.DB,
	jwtSecret string,
	jwtExpiration time.Duration,
	refreshTokenExpiration time.Duration,
) UserService {
	return &userService{
		userRepo:               userRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		db:                     db,
		jwtSecret:              jwtSecret,
		jwtExpiration:          jwtExpiration,
		refreshTokenExpiration: refreshTokenExpiration,
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req *model.UserRegister) (*model.UserResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱失败: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建用户
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "user", // 默认角色
		IsActive:     true,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	return user.ToResponse(), nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *model.UserLogin) (*model.LoginResponse, error) {
	// 根据用户名或邮箱查找用户
	var user *model.User
	var err error

	if strings.Contains(req.Username, "@") {
		user, err = s.userRepo.GetByEmail(req.Username)
	} else {
		user, err = s.userRepo.GetByUsername(req.Username)
	}

	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 检查用户是否激活
	if !user.IsActive {
		return nil, fmt.Errorf("用户账户已被禁用")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("密码错误")
	}

	// 更新最后登录时间
	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		// 记录错误但不影响登录流程
		fmt.Printf("更新最后登录时间失败: %v\n", err)
	}

	// 生成JWT令牌
	accessToken, err := s.GenerateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %v", err)
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %v", err)
	}

	return &model.LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtExpiration.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(ctx context.Context, id string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %v", err)
	}
	return user.ToResponse(), nil
}

// Update 更新用户信息
func (s *userService) Update(ctx context.Context, id string, req *model.UserUpdate) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %v", err)
	}

	// 更新字段
	if req.Email != "" && req.Email != user.Email {
		// 检查邮箱是否已存在
		exists, err := s.userRepo.ExistsByEmail(req.Email)
		if err != nil {
			return nil, fmt.Errorf("检查邮箱失败: %v", err)
		}
		if exists {
			return nil, fmt.Errorf("邮箱已存在")
		}
		user.Email = req.Email
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("更新用户失败: %v", err)
	}

	return user.ToResponse(), nil
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id string) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}
	return nil
}

// List 获取用户列表
func (s *userService) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]model.UserResponse, int64, error) {
	users, total, err := s.userRepo.List(offset, limit, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %v", err)
	}

	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *user.ToResponse())
	}

	return userResponses, total, nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID string, req *model.ChangePassword) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %v", err)
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		return fmt.Errorf("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新密码
	user.PasswordHash = string(hashedPassword)
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}

	return nil
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(ctx context.Context, req *model.PasswordReset) error {
	// 验证令牌
	resetToken, err := s.passwordResetTokenRepo.GetByToken(req.Token)
	if err != nil {
		return fmt.Errorf("无效的重置令牌")
	}

	// 检查令牌是否过期
	if time.Now().After(resetToken.ExpiresAt) {
		// 删除过期令牌
		s.passwordResetTokenRepo.Delete(req.Token)
		return fmt.Errorf("重置令牌已过期")
	}

	// 获取用户
	user, err := s.userRepo.GetByID(resetToken.UserID)
	if err != nil {
		return fmt.Errorf("获取用户失败: %v", err)
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新密码
	user.PasswordHash = string(hashedPassword)
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}

	// 删除令牌
	err = s.passwordResetTokenRepo.Delete(req.Token)
	if err != nil {
		// 记录错误但不影响重置流程
		fmt.Printf("删除重置令牌失败: %v\n", err)
	}

	return nil
}

// RequestPasswordReset 请求密码重置
func (s *userService) RequestPasswordReset(ctx context.Context, req *model.PasswordResetRequest) error {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		// 为了安全，即使用户不存在也返回成功
		return nil
	}

	// 生成重置令牌
	token, err := s.generateResetToken()
	if err != nil {
		return fmt.Errorf("生成重置令牌失败: %v", err)
	}

	// 创建重置令牌记录
	resetToken := &model.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24小时后过期
	}

	err = s.passwordResetTokenRepo.Create(resetToken)
	if err != nil {
		return fmt.Errorf("保存重置令牌失败: %v", err)
	}

	// TODO: 发送重置邮件
	fmt.Printf("密码重置令牌已生成: %s, 用户: %s\n", token, user.Email)

	return nil
}

// ValidateJWT 验证JWT令牌
func (s *userService) ValidateJWT(tokenString string) (*model.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &model.JWTClaims{
			UserID:   claims["user_id"].(string),
			Username: claims["username"].(string),
			Role:     claims["role"].(string),
		}, nil
	}

	return nil, fmt.Errorf("无效的令牌")
}

// GenerateJWT 生成JWT令牌
func (s *userService) GenerateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(s.jwtExpiration).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// RefreshJWT 刷新JWT令牌
func (s *userService) RefreshJWT(ctx context.Context, tokenString string) (string, error) {
	// 验证当前令牌
	claims, err := s.ValidateJWT(tokenString)
	if err != nil {
		return "", fmt.Errorf("无效的令牌: %v", err)
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("获取用户失败: %v", err)
	}

	// 检查用户是否激活
	if !user.IsActive {
		return "", fmt.Errorf("用户账户已被禁用")
	}

	// 生成新令牌
	return s.GenerateJWT(user)
}

// generateResetToken 生成重置令牌
func (s *userService) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateRefreshToken 生成刷新令牌
func (s *userService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
