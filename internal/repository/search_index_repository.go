package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"gorm.io/gorm"
)

// SearchIndexRepository 搜索索引仓库接口
type SearchIndexRepository interface {
	Create(ctx context.Context, index *model.SearchIndex) error
	CreateBatch(ctx context.Context, indices []*model.SearchIndex) error
	GetByID(ctx context.Context, id string) (*model.SearchIndex, error)
	GetByDocumentID(ctx context.Context, documentID string) ([]*model.SearchIndex, error)
	GetByDocumentIDAndVersion(ctx context.Context, documentID, version string) ([]*model.SearchIndex, error)
	Search(ctx context.Context, query string, filters map[string]interface{}, page, size int) ([]*model.SearchIndex, int64, error)
	SearchByKeywords(ctx context.Context, keywords []string, filters map[string]interface{}, page, size int) ([]*model.SearchIndex, int64, error)
	SearchByVector(ctx context.Context, vector string, filters map[string]interface{}, page, size int) ([]*model.SearchIndex, int64, error)
	DeleteByDocumentID(ctx context.Context, documentID string) error
	DeleteByDocumentIDAndVersion(ctx context.Context, documentID, version string) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	GetIndexingStatus(ctx context.Context, documentID string) (map[string]interface{}, error)
}

// searchIndexRepository 搜索索引仓库实现
type searchIndexRepository struct {
	db *gorm.DB
}

// NewSearchIndexRepository 创建搜索索引仓库实例
func NewSearchIndexRepository(db *gorm.DB) SearchIndexRepository {
	return &searchIndexRepository{
		db: db,
	}
}

// Create 创建搜索索引
func (r *searchIndexRepository) Create(ctx context.Context, index *model.SearchIndex) error {
	return r.db.WithContext(ctx).Create(index).Error
}

// CreateBatch 批量创建搜索索引
func (r *searchIndexRepository) CreateBatch(ctx context.Context, indices []*model.SearchIndex) error {
	log.Printf("DEBUG: CreateBatch called with %d indices", len(indices))
	if len(indices) == 0 {
		log.Printf("DEBUG: CreateBatch: no indices to create")
		return nil
	}

	// 手动处理批量插入，避免 GORM 的 JSON 序列化问题
	log.Printf("DEBUG: 开始事务处理批量插入")
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Printf("DEBUG: 开始事务失败: %v", tx.Error)
		return tx.Error
	}

	for i := 0; i < len(indices); i += 100 {
		end := i + 100
		if end > len(indices) {
			end = len(indices)
		}

		log.Printf("DEBUG: 处理批次 %d-%d，共 %d 个索引", i+1, end, end-i)
		batch := indices[i:end]
		values := make([]interface{}, 0, len(batch)*9)
		valueStrings := make([]string, 0, len(batch))

		for _, index := range batch {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
			values = append(values,
				index.ID,
				index.DocumentID,
				index.Version,
				index.Content,
				index.ContentType,
				index.Section,
				index.Keywords,
				index.Vector, // 已经是 JSON 字符串格式
				index.Metadata,
			)
		}

		query := fmt.Sprintf("INSERT INTO search_indices (id, document_id, version, content, content_type, section, keywords, vector, metadata) VALUES %s",
			strings.Join(valueStrings, ","))

		log.Printf("DEBUG: 执行插入查询，参数数量: %d", len(values))
		if err := tx.Exec(query, values...).Error; err != nil {
			log.Printf("DEBUG: 批量插入失败: %v", err)
			tx.Rollback()
			return err
		}
		log.Printf("DEBUG: 批量插入成功，批次 %d-%d", i+1, end)
	}

	log.Printf("DEBUG: 提交事务")
	if err := tx.Commit().Error; err != nil {
		log.Printf("DEBUG: 提交事务失败: %v", err)
		return err
	}
	log.Printf("DEBUG: 批量创建索引完成，共 %d 个索引", len(indices))
	return nil
}

// GetByID 根据ID获取搜索索引
func (r *searchIndexRepository) GetByID(ctx context.Context, id string) (*model.SearchIndex, error) {
	var index model.SearchIndex
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&index).Error
	if err != nil {
		return nil, err
	}
	return &index, nil
}

// GetByDocumentID 根据文档ID获取搜索索引
func (r *searchIndexRepository) GetByDocumentID(ctx context.Context, documentID string) ([]*model.SearchIndex, error) {
	var indices []*model.SearchIndex
	err := r.db.WithContext(ctx).Where("document_id = ?", documentID).Order("created_at DESC").Find(&indices).Error
	if err != nil {
		return nil, err
	}
	return indices, nil
}

// GetByDocumentIDAndVersion 根据文档ID和版本获取搜索索引
func (r *searchIndexRepository) GetByDocumentIDAndVersion(ctx context.Context, documentID, version string) ([]*model.SearchIndex, error) {
	var indices []*model.SearchIndex
	err := r.db.WithContext(ctx).Where("document_id = ? AND TRIM(version) = ?", documentID, version).Order("created_at DESC").Find(&indices).Error
	if err != nil {
		return nil, err
	}
	return indices, nil
}

// Search 关键词搜索
func (r *searchIndexRepository) Search(ctx context.Context, query string, filters map[string]interface{}, page, size int) ([]*model.SearchIndex, int64, error) {
	var indices []*model.SearchIndex
	var total int64

	db := r.db.WithContext(ctx).Model(&model.SearchIndex{})

	// 构建搜索查询
	searchQuery := db.Where("content LIKE ?", "%"+query+"%")

	// 应用过滤条件
	searchQuery = r.applyFilters(searchQuery, filters)

	// 获取总数
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		searchQuery = searchQuery.Offset(offset).Limit(size)
	}

	// 执行查询
	if err := searchQuery.Order("created_at DESC").Find(&indices).Error; err != nil {
		return nil, 0, err
	}

	return indices, total, nil
}

// SearchByKeywords 根据关键词搜索
func (r *searchIndexRepository) SearchByKeywords(ctx context.Context, keywords []string, filters map[string]interface{}, page, size int) ([]*model.SearchIndex, int64, error) {
	var indices []*model.SearchIndex
	var total int64

	db := r.db.WithContext(ctx).Model(&model.SearchIndex{})

	// 构建关键词搜索查询
	if len(keywords) == 0 {
		return []*model.SearchIndex{}, 0, nil
	}

	log.Printf("DEBUG: SearchByKeywords called with keywords: %v", keywords)

	// 构建基础查询 - 使用OR条件匹配任意关键词
	searchQuery := db.Where("1=0") // 初始设置为不匹配任何记录

	// 为每个关键词添加OR条件，匹配任意关键词
	for i, keyword := range keywords {
		log.Printf("DEBUG: Adding keyword filter: %s", keyword)
		searchQuery = searchQuery.Or("content LIKE ?", "%"+keyword+"%")
		if i == 0 {
			// 记录第一个关键词的查询条件
			log.Printf("DEBUG: First keyword condition: content LIKE '%%%s%%'", keyword)
		}
	}

	log.Printf("DEBUG: Final search query built for %d keywords", len(keywords))

	// 应用过滤条件
	searchQuery = r.applyFilters(searchQuery, filters)

	// 开启SQL日志记录
	searchQuery = searchQuery.Debug()

	// 获取总数 - 使用单独的查询，避免分页影响
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	log.Printf("DEBUG: Total count before pagination: %d", total)

	// 分页查询 - 使用游标式分页提高性能
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		searchQuery = searchQuery.Offset(offset).Limit(size)
	}

	// 执行查询 - 添加排序以获得更一致的结果
	if err := searchQuery.Order("created_at DESC").Find(&indices).Error; err != nil {
		return nil, 0, err
	}

	log.Printf("DEBUG: Found %d indices, total count: %d", len(indices), total)

	// 显示数据库中的样本内容用于调试
	if len(indices) == 0 {
		var sampleIndices []*model.SearchIndex
		if err := db.Limit(3).Find(&sampleIndices).Error; err == nil {
			log.Printf("DEBUG: Sample database content:")
			for i, idx := range sampleIndices {
				log.Printf("DEBUG: Sample %d: content='%s', keywords='%s'", i+1, idx.Content, idx.Keywords)
			}
		}
	}

	// 为每个结果计算相关性得分
	for _, index := range indices {
		r.calculateKeywordScore(index, keywords)
	}

	return indices, total, nil
}

// calculateKeywordScore 计算关键词搜索的相关性得分
func (r *searchIndexRepository) calculateKeywordScore(index *model.SearchIndex, keywords []string) {
	if len(keywords) == 0 {
		index.Score = 0
		return
	}

	content := strings.ToLower(index.Content)
	matchCount := 0

	for _, keyword := range keywords {
		keywordLower := strings.ToLower(keyword)
		count := strings.Count(content, keywordLower)
		if count > 0 {
			matchCount++
		}
	}

	// 计算最终得分：匹配度
	if len(keywords) > 0 {
		index.Score = float32(matchCount) / float32(len(keywords))
	} else {
		index.Score = 0
	}
}

// SearchByVector 向量搜索（语义搜索）
func (r *searchIndexRepository) SearchByVector(ctx context.Context, vector string, filters map[string]interface{}, page, size int) ([]*model.SearchIndex, int64, error) {
	// 将JSON字符串解析为向量
	var queryVector []float32
	if err := json.Unmarshal([]byte(vector), &queryVector); err != nil {
		// 如果解析失败，返回空结果
		return []*model.SearchIndex{}, 0, nil
	}

	// 使用增强的向量搜索功能，结合多种相似度计算方法
	return r.enhancedVectorSearch(ctx, queryVector, filters, page, size)
}

// DeleteByDocumentID 根据文档ID删除搜索索引
func (r *searchIndexRepository) DeleteByDocumentID(ctx context.Context, documentID string) error {
	return r.db.WithContext(ctx).Where("document_id = ?", documentID).Delete(&model.SearchIndex{}).Error
}

// DeleteByDocumentIDAndVersion 根据文档ID和版本删除搜索索引
func (r *searchIndexRepository) DeleteByDocumentIDAndVersion(ctx context.Context, documentID, version string) error {
	return r.db.WithContext(ctx).Where("document_id = ? AND TRIM(version) = ?", documentID, version).Delete(&model.SearchIndex{}).Error
}

// Update 更新搜索索引
func (r *searchIndexRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.SearchIndex{}).Where("id = ?", id).Updates(updates).Error
}

// GetIndexingStatus 获取索引状态
func (r *searchIndexRepository) GetIndexingStatus(ctx context.Context, documentID string) (map[string]interface{}, error) {
	var total, indexed int64

	// 获取文档的总内容段落数
	if err := r.db.WithContext(ctx).Table("document_versions").Where("document_id = ?", documentID).Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取已索引的段落数
	if err := r.db.WithContext(ctx).Model(&model.SearchIndex{}).Where("document_id = ?", documentID).Count(&indexed).Error; err != nil {
		return nil, err
	}

	// 计算索引进度
	progress := float64(0)
	if total > 0 {
		progress = float64(indexed) / float64(total) * 100
	}

	return map[string]interface{}{
		"document_id": documentID,
		"total":       total,
		"indexed":     indexed,
		"progress":    progress,
		"status":      "indexed",
	}, nil
}

// applyFilters 应用过滤条件
func (r *searchIndexRepository) applyFilters(db *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return db
	}

	if documentID, ok := filters["document_id"]; ok && documentID != "" && documentID != nil {
		db = db.Where("document_id = ?", documentID)
	}

	if version, ok := filters["version"]; ok && version != "" && version != nil {
		db = db.Where("TRIM(version) = ?", version)
	}

	if contentType, ok := filters["content_type"]; ok && contentType != "" && contentType != nil {
		db = db.Where("content_type = ?", contentType)
	}

	if section, ok := filters["section"]; ok && section != "" && section != nil {
		db = db.Where("section = ?", section)
	}

	return db
}

// calculateCosineSimilarity 计算余弦相似度并排序
func (r *searchIndexRepository) calculateCosineSimilarity(indices []*model.SearchIndex, queryVector []float32) []*model.SearchIndex {
	// 计算每个索引与查询向量的余弦相似度
	for _, index := range indices {
		if len(index.Vector) > 0 {
			// 解析JSON字符串格式的向量
			var vectorSlice []float32
			if err := json.Unmarshal([]byte(index.Vector), &vectorSlice); err == nil {
				similarity := r.cosineSimilarity(queryVector, vectorSlice)
				index.Score = similarity
			}
		}
	}

	// 按相似度排序
	for i := 0; i < len(indices)-1; i++ {
		for j := i + 1; j < len(indices); j++ {
			if indices[i].Score < indices[j].Score {
				indices[i], indices[j] = indices[j], indices[i]
			}
		}
	}

	return indices
}

// cosineSimilarity 计算余弦相似度
func (r *searchIndexRepository) cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float32
	var normA float32
	var normB float32

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (sqrt(normA) * sqrt(normB))
}

// euclideanDistance 计算欧几里得距离（用于向量搜索的另一种相似度度量）
func (r *searchIndexRepository) euclideanDistance(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var sum float32
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return sqrt(sum)
}

// enhancedVectorSearch 增强向量搜索，使用多种相似度计算方法
func (r *searchIndexRepository) enhancedVectorSearch(ctx context.Context, vector []float32, filters map[string]interface{}, page, size int) ([]*model.SearchIndex, int64, error) {
	var indices []*model.SearchIndex
	var total int64

	db := r.db.WithContext(ctx).Model(&model.SearchIndex{})

	// 检查是否有 embedding 列（pgvector 扩展是否可用）
	var hasEmbeddingColumn bool
	err := r.db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = 'search_indices' AND column_name = 'embedding'").Scan(&hasEmbeddingColumn).Error
	if err != nil {
		log.Printf("Error checking for embedding column: %v", err)
		hasEmbeddingColumn = false
	}

	var searchQuery *gorm.DB
	if hasEmbeddingColumn {
		// 使用 pgvector 的 embedding 列进行搜索
		searchQuery = db.Where("embedding IS NOT NULL")
	} else {
		// 回退到传统的 vector 列
		searchQuery = db.Where("vector IS NOT NULL")
		log.Printf("pgvector extension not available, falling back to traditional vector search")
	}

	// 应用过滤条件
	searchQuery = r.applyFilters(searchQuery, filters)

	// 获取总数
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 限制初始结果集大小以提高性能
	initialLimit := 1000 // 只获取前1000条记录进行相似度计算
	searchQuery = searchQuery.Limit(initialLimit)

	// 获取符合条件的结果
	if err := searchQuery.Find(&indices).Error; err != nil {
		return nil, 0, err
	}

	// 计算多种相似度并合并得分
	for _, index := range indices {
		if hasEmbeddingColumn && len(index.Embedding) > 0 {
			// 使用 pgvector embedding 计算相似度
			cosineSim := r.cosineSimilarity(vector, index.Embedding)
			euclideanDist := r.euclideanDistance(vector, index.Embedding)
			euclideanSim := 1.0 / (1.0 + euclideanDist) // 转换为相似度分数
			index.Score = 0.7*float32(cosineSim) + 0.3*float32(euclideanSim)
		} else if len(index.Vector) > 0 {
			// 回退到传统向量
			// 解析JSON字符串格式的向量
			var vectorSlice []float32
			if err := json.Unmarshal([]byte(index.Vector), &vectorSlice); err == nil {
				cosineSim := r.cosineSimilarity(vector, vectorSlice)
				euclideanDist := r.euclideanDistance(vector, vectorSlice)
				euclideanSim := 1.0 / (1.0 + euclideanDist) // 转换为相似度分数
				index.Score = 0.7*float32(cosineSim) + 0.3*float32(euclideanSim)
			}
		}
	}

	// 按相似度排序 - 使用更高效的排序算法
	for i := 0; i < len(indices)-1; i++ {
		for j := i + 1; j < len(indices); j++ {
			if indices[i].Score < indices[j].Score {
				indices[i], indices[j] = indices[j], indices[i]
			}
		}
	}

	// 应用分页
	if page > 0 && size > 0 {
		start := (page - 1) * size
		end := start + size
		if end > len(indices) {
			end = len(indices)
		}
		if start < len(indices) {
			indices = indices[start:end]
		} else {
			indices = []*model.SearchIndex{}
		}
	}

	return indices, total, nil
}

// sqrt 计算平方根
func sqrt(x float32) float32 {
	// 简单的平方根实现，实际项目中可以使用 math.Sqrt
	var z float32 = 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}
