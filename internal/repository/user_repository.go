package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	// 创建用户
	Create(user *model.User) error

	// 根据ID获取用户
	GetByID(id string) (*model.User, error)

	// 根据用户名获取用户
	GetByUsername(username string) (*model.User, error)

	// 根据邮箱获取用户
	GetByEmail(email string) (*model.User, error)

	// 更新用户
	Update(user *model.User) error

	// 删除用户
	Delete(id string) error

	// 获取用户列表
	List(offset, limit int, filters map[string]interface{}) ([]model.User, int64, error)

	// 更新最后登录时间
	UpdateLastLogin(id string) error

	// 检查用户名是否存在
	ExistsByUsername(username string) (bool, error)

	// 检查邮箱是否存在
	ExistsByEmail(email string) (bool, error)
}

// userRepository 用户仓库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

// List 获取用户列表
func (r *userRepository) List(offset, limit int, filters map[string]interface{}) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "role":
			query = query.Where("role = ?", value)
		case "is_active":
			query = query.Where("is_active = ?", value)
		case "search":
			searchTerm := value.(string)
			query = query.Where("username ILIKE ? OR email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
				"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")
		}
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(id string) error {
	now := time.Now()
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("last_login", now).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// PasswordResetTokenRepository 密码重置令牌仓库接口
type PasswordResetTokenRepository interface {
	// 创建密码重置令牌
	Create(token *model.PasswordResetToken) error

	// 根据令牌获取密码重置令牌
	GetByToken(token string) (*model.PasswordResetToken, error)

	// 删除令牌
	Delete(token string) error

	// 清理过期令牌
	CleanExpiredTokens() error
}

// passwordResetTokenRepository 密码重置令牌仓库实现
type passwordResetTokenRepository struct {
	db *gorm.DB
}

// NewPasswordResetTokenRepository 创建密码重置令牌仓库实例
func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{
		db: db,
	}
}

// Create 创建密码重置令牌
func (r *passwordResetTokenRepository) Create(token *model.PasswordResetToken) error {
	return r.db.Create(token).Error
}

// GetByToken 根据令牌获取密码重置令牌
func (r *passwordResetTokenRepository) GetByToken(token string) (*model.PasswordResetToken, error) {
	var resetToken model.PasswordResetToken
	err := r.db.Where("token = ?", token).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

// Delete 删除令牌
func (r *passwordResetTokenRepository) Delete(token string) error {
	return r.db.Delete(&model.PasswordResetToken{}, "token = ?", token).Error
}

// CleanExpiredTokens 清理过期令牌
func (r *passwordResetTokenRepository) CleanExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&model.PasswordResetToken{}).Error
}
