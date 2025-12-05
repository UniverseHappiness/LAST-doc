package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"ai-doc-library/internal/model"
)

// DocumentRepository 文档仓库接口
type DocumentRepository interface {
	Create(ctx context.Context, document *model.Document) error
	GetByID(ctx context.Context, id string) (*model.Document, error)
	List(ctx context.Context, page, size int, filters map[string]interface{}) ([]*model.Document, int64, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	GetByLibrary(ctx context.Context, library string, page, size int) ([]*model.Document, int64, error)
	GetByType(ctx context.Context, docType model.DocumentType, page, size int) ([]*model.Document, int64, error)
	GetByTag(ctx context.Context, tag string, page, size int) ([]*model.Document, int64, error)
	GetByVersion(ctx context.Context, version string, page, size int) ([]*model.Document, int64, error)
	GetByLibraryAndVersion(ctx context.Context, library, version string, page, size int) ([]*model.Document, int64, error)
	GetByName(ctx context.Context, name string, page, size int) ([]*model.Document, int64, error)
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
}

// documentRepository 文档仓库实现
type documentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository 创建文档仓库实例
func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{
		db: db,
	}
}

// Create 创建文档
func (r *documentRepository) Create(ctx context.Context, document *model.Document) error {
	return r.db.WithContext(ctx).Create(document).Error
}

// GetByID 根据ID获取文档
func (r *documentRepository) GetByID(ctx context.Context, id string) (*model.Document, error) {
	var document model.Document
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// List 获取文档列表
func (r *documentRepository) List(ctx context.Context, page, size int, filters map[string]interface{}) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Document{})

	// 应用过滤条件
	if filters != nil {
		if library, ok := filters["library"]; ok {
			query = query.Where("library = ?", library)
		}
		if docType, ok := filters["type"]; ok {
			query = query.Where("type = ?", docType)
		}
		if version, ok := filters["version"]; ok {
			query = query.Where("version = ?", version)
		}
		if status, ok := filters["status"]; ok {
			query = query.Where("status = ?", status)
		}
		if name, ok := filters["name"]; ok {
			query = query.Where("name LIKE ?", "%"+name.(string)+"%")
		}
		if tags, ok := filters["tags"]; ok {
			if tagList, ok := tags.([]string); ok && len(tagList) > 0 {
				for _, tag := range tagList {
					query = query.Where("? = ANY(tags)", tag)
				}
			}
		}
		if startTime, ok := filters["start_time"]; ok {
			query = query.Where("created_at >= ?", startTime)
		}
		if endTime, ok := filters["end_time"]; ok {
			query = query.Where("created_at <= ?", endTime)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	// 排序
	query = query.Order("created_at DESC")

	// 执行查询
	if err := query.Find(&documents).Error; err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// Update 更新文档
func (r *documentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.WithContext(ctx).Model(&model.Document{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除文档
func (r *documentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Document{}).Error
}

// GetByLibrary 根据库获取文档
func (r *documentRepository) GetByLibrary(ctx context.Context, library string, page, size int) ([]*model.Document, int64, error) {
	return r.List(ctx, page, size, map[string]interface{}{"library": library})
}

// GetByType 根据类型获取文档
func (r *documentRepository) GetByType(ctx context.Context, docType model.DocumentType, page, size int) ([]*model.Document, int64, error) {
	return r.List(ctx, page, size, map[string]interface{}{"type": docType})
}

// GetByTag 根据标签获取文档
func (r *documentRepository) GetByTag(ctx context.Context, tag string, page, size int) ([]*model.Document, int64, error) {
	return r.List(ctx, page, size, map[string]interface{}{"tags": []string{tag}})
}

// GetByVersion 根据版本获取文档
func (r *documentRepository) GetByVersion(ctx context.Context, version string, page, size int) ([]*model.Document, int64, error) {
	return r.List(ctx, page, size, map[string]interface{}{"version": version})
}

// GetByLibraryAndVersion 根据库和版本获取文档
func (r *documentRepository) GetByLibraryAndVersion(ctx context.Context, library, version string, page, size int) ([]*model.Document, int64, error) {
	return r.List(ctx, page, size, map[string]interface{}{"library": library, "version": version})
}

// GetByName 根据名称获取文档
func (r *documentRepository) GetByName(ctx context.Context, name string, page, size int) ([]*model.Document, int64, error) {
	return r.List(ctx, page, size, map[string]interface{}{"name": name})
}

// Count 获取文档总数
func (r *documentRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Document{})

	// 应用过滤条件
	if filters != nil {
		if library, ok := filters["library"]; ok {
			query = query.Where("library = ?", library)
		}
		if docType, ok := filters["type"]; ok {
			query = query.Where("type = ?", docType)
		}
		if version, ok := filters["version"]; ok {
			query = query.Where("version = ?", version)
		}
		if status, ok := filters["status"]; ok {
			query = query.Where("status = ?", status)
		}
		if name, ok := filters["name"]; ok {
			query = query.Where("name LIKE ?", "%"+name.(string)+"%")
		}
		if tags, ok := filters["tags"]; ok {
			if tagList, ok := tags.([]string); ok && len(tagList) > 0 {
				for _, tag := range tagList {
					query = query.Where("? = ANY(tags)", tag)
				}
			}
		}
		if startTime, ok := filters["start_time"]; ok {
			query = query.Where("created_at >= ?", startTime)
		}
		if endTime, ok := filters["end_time"]; ok {
			query = query.Where("created_at <= ?", endTime)
		}
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}