package repository

import (
	"context"
	"encoding/json"

	"ai-doc-library/internal/model"

	"gorm.io/gorm"
)

// DocumentMetadataRepository 文档元数据仓库接口
type DocumentMetadataRepository interface {
	Create(ctx context.Context, metadata *model.DocumentMetadata) error
	GetByID(ctx context.Context, id string) (*model.DocumentMetadata, error)
	GetByDocumentID(ctx context.Context, documentID string) ([]*model.DocumentMetadata, error)
	GetLatestMetadata(ctx context.Context, documentID string) (*model.DocumentMetadata, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	DeleteByDocumentID(ctx context.Context, documentID string) error
	Count(ctx context.Context, documentID string) (int64, error)
}

// documentMetadataRepository 文档元数据仓库实现
type documentMetadataRepository struct {
	db *gorm.DB
}

// NewDocumentMetadataRepository 创建文档元数据仓库实例
func NewDocumentMetadataRepository(db *gorm.DB) DocumentMetadataRepository {
	return &documentMetadataRepository{
		db: db,
	}
}

// Create 创建文档元数据
func (r *documentMetadataRepository) Create(ctx context.Context, metadata *model.DocumentMetadata) error {
	return r.db.WithContext(ctx).Create(metadata).Error
}

// GetByID 根据ID获取文档元数据
func (r *documentMetadataRepository) GetByID(ctx context.Context, id string) (*model.DocumentMetadata, error) {
	var metadata model.DocumentMetadata
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&metadata).Error
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

// GetByDocumentID 根据文档ID获取所有元数据
func (r *documentMetadataRepository) GetByDocumentID(ctx context.Context, documentID string) ([]*model.DocumentMetadata, error) {
	// 使用原生SQL查询来处理JSONB字段
	rows, err := r.db.WithContext(ctx).
		Table("document_metadata").
		Select("id, document_id, metadata, created_at, updated_at").
		Where("document_id = ?", documentID).
		Order("created_at DESC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadataList []*model.DocumentMetadata
	for rows.Next() {
		var metadata model.DocumentMetadata
		var metadataBytes []byte

		// 扫描行数据，将JSONB字段作为[]byte处理
		err := rows.Scan(&metadata.ID, &metadata.DocumentID, &metadataBytes, &metadata.CreatedAt, &metadata.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// 解析JSON字节数组为map[string]interface{}
		if len(metadataBytes) > 0 {
			var metadataMap map[string]interface{}
			if err := json.Unmarshal(metadataBytes, &metadataMap); err != nil {
				return nil, err
			}
			metadata.Metadata = metadataMap
		}

		metadataList = append(metadataList, &metadata)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metadataList, nil
}

// GetLatestMetadata 获取文档的最新元数据
func (r *documentMetadataRepository) GetLatestMetadata(ctx context.Context, documentID string) (*model.DocumentMetadata, error) {
	var metadata model.DocumentMetadata
	err := r.db.WithContext(ctx).
		Where("document_id = ?", documentID).
		Order("created_at DESC").
		First(&metadata).Error
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

// Update 更新文档元数据
func (r *documentMetadataRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&model.DocumentMetadata{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete 删除文档元数据
func (r *documentMetadataRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.DocumentMetadata{}).Error
}

// DeleteByDocumentID 根据文档ID删除所有元数据
func (r *documentMetadataRepository) DeleteByDocumentID(ctx context.Context, documentID string) error {
	return r.db.WithContext(ctx).
		Where("document_id = ?", documentID).
		Delete(&model.DocumentMetadata{}).Error
}

// Count 获取文档的元数据数量
func (r *documentMetadataRepository) Count(ctx context.Context, documentID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.DocumentMetadata{}).
		Where("document_id = ?", documentID).
		Count(&count).Error
	return count, err
}
