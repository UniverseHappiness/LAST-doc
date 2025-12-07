package repository

import (
	"context"
	"log"

	"github.com/UniverseHappiness/LAST-doc/internal/model"

	"gorm.io/gorm"
)

// DocumentVersionRepository 文档版本仓库接口
type DocumentVersionRepository interface {
	Create(ctx context.Context, version *model.DocumentVersion) error
	GetByID(ctx context.Context, id string) (*model.DocumentVersion, error)
	GetByDocumentID(ctx context.Context, documentID string) ([]*model.DocumentVersion, error)
	GetByDocumentIDAndVersion(ctx context.Context, documentID, version string) (*model.DocumentVersion, error)
	GetLatestVersion(ctx context.Context, documentID string) (*model.DocumentVersion, error)
	GetVersionsByStatus(ctx context.Context, documentID string, status model.DocumentStatus) ([]*model.DocumentVersion, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	UpdateContent(ctx context.Context, documentID, version string, content string, status model.DocumentStatus) error
	UpdateStatus(ctx context.Context, documentID, version string, status model.DocumentStatus) error
	Delete(ctx context.Context, id string) error
	DeleteByDocumentID(ctx context.Context, documentID string) error
	Count(ctx context.Context, documentID string) (int64, error)
}

// documentVersionRepository 文档版本仓库实现
type documentVersionRepository struct {
	db *gorm.DB
}

// NewDocumentVersionRepository 创建文档版本仓库实例
func NewDocumentVersionRepository(db *gorm.DB) DocumentVersionRepository {
	return &documentVersionRepository{
		db: db,
	}
}

// Create 创建文档版本
func (r *documentVersionRepository) Create(ctx context.Context, version *model.DocumentVersion) error {
	log.Printf("DEBUG: 数据库层创建文档版本 - 文档ID: %s, 版本: %s, 版本记录ID: %s\n",
		version.DocumentID, version.Version, version.ID)

	err := r.db.WithContext(ctx).Create(version).Error
	if err != nil {
		log.Printf("DEBUG: 数据库层创建文档版本失败 - 文档ID: %s, 版本: %s, 错误: %v\n",
			version.DocumentID, version.Version, err)
		return err
	}

	log.Printf("DEBUG: 数据库层创建文档版本成功 - 文档ID: %s, 版本: %s, 版本记录ID: %s\n",
		version.DocumentID, version.Version, version.ID)
	return nil
}

// GetByID 根据ID获取文档版本
func (r *documentVersionRepository) GetByID(ctx context.Context, id string) (*model.DocumentVersion, error) {
	var version model.DocumentVersion
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetByDocumentID 根据文档ID获取所有版本
func (r *documentVersionRepository) GetByDocumentID(ctx context.Context, documentID string) ([]*model.DocumentVersion, error) {
	var versions []*model.DocumentVersion
	err := r.db.WithContext(ctx).
		Where("document_id = ?", documentID).
		Order("created_at DESC").
		Find(&versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

// GetByDocumentIDAndVersion 根据文档ID和版本号获取文档版本
func (r *documentVersionRepository) GetByDocumentIDAndVersion(ctx context.Context, documentID, version string) (*model.DocumentVersion, error) {
	var docVersion model.DocumentVersion
	err := r.db.WithContext(ctx).
		Where("document_id = ? AND version = ?", documentID, version).
		First(&docVersion).Error
	if err != nil {
		return nil, err
	}
	return &docVersion, nil
}

// GetLatestVersion 获取文档的最新版本
func (r *documentVersionRepository) GetLatestVersion(ctx context.Context, documentID string) (*model.DocumentVersion, error) {
	var version model.DocumentVersion
	err := r.db.WithContext(ctx).
		Where("document_id = ? AND status = ?", documentID, model.DocumentStatusCompleted).
		Order("created_at DESC").
		First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetVersionsByStatus 根据状态获取文档版本
func (r *documentVersionRepository) GetVersionsByStatus(ctx context.Context, documentID string, status model.DocumentStatus) ([]*model.DocumentVersion, error) {
	var versions []*model.DocumentVersion
	err := r.db.WithContext(ctx).
		Where("document_id = ? AND status = ?", documentID, status).
		Order("created_at DESC").
		Find(&versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}

// Update 更新文档版本
func (r *documentVersionRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&model.DocumentVersion{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// UpdateContent 更新文档版本内容和状态
func (r *documentVersionRepository) UpdateContent(ctx context.Context, documentID, version string, content string, status model.DocumentStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.DocumentVersion{}).
		Where("document_id = ? AND version = ?", documentID, version).
		Updates(map[string]interface{}{
			"content": content,
			"status":  status,
		}).Error
}

// UpdateStatus 更新文档版本状态
func (r *documentVersionRepository) UpdateStatus(ctx context.Context, documentID, version string, status model.DocumentStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.DocumentVersion{}).
		Where("document_id = ? AND version = ?", documentID, version).
		Update("status", status).Error
}

// Delete 删除文档版本
func (r *documentVersionRepository) Delete(ctx context.Context, id string) error {
	// 获取文档版本信息，以便删除对应的搜索索引
	var version model.DocumentVersion
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&version).Error; err != nil {
		return err
	}

	// 删除对应的搜索索引
	if err := r.db.WithContext(ctx).Table("search_indices").Where("document_id = ? AND version = ?", version.DocumentID, version.Version).Delete(nil).Error; err != nil {
		return err
	}

	// 删除文档版本
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.DocumentVersion{}).Error
}

// DeleteByDocumentID 根据文档ID删除所有版本
func (r *documentVersionRepository) DeleteByDocumentID(ctx context.Context, documentID string) error {
	// 先删除所有版本的搜索索引
	if err := r.db.WithContext(ctx).Table("search_indices").Where("document_id = ?", documentID).Delete(nil).Error; err != nil {
		return err
	}

	// 然后删除所有文档版本
	return r.db.WithContext(ctx).
		Where("document_id = ?", documentID).
		Delete(&model.DocumentVersion{}).Error
}

// Count 获取文档的版本数量
func (r *documentVersionRepository) Count(ctx context.Context, documentID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.DocumentVersion{}).
		Where("document_id = ?", documentID).
		Count(&count).Error
	return count, err
}
