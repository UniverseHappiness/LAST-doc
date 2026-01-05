package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/repository"
	"github.com/google/uuid"
)

// SearchService 搜索服务接口
type SearchService interface {
	BuildIndex(ctx context.Context, documentID, version string) error
	BuildIndexBatch(ctx context.Context, indices []*model.SearchIndex) error
	Search(ctx context.Context, request *model.SearchRequest) (*model.SearchResponse, error)
	GetIndexingStatus(ctx context.Context, documentID string) (map[string]interface{}, error)
	DeleteIndex(ctx context.Context, documentID string) error
	DeleteIndexByVersion(ctx context.Context, documentID, version string) error
	ClearCache() error
}

// searchService 搜索服务实现
type searchService struct {
	indexRepo        repository.SearchIndexRepository
	documentRepo     repository.DocumentRepository
	versionRepo      repository.DocumentVersionRepository
	cacheService     CacheService
	embeddingService EmbeddingService
	indexingEnabled  bool
}

// NewSearchService 创建搜索服务实例
func NewSearchService(
	indexRepo repository.SearchIndexRepository,
	documentRepo repository.DocumentRepository,
	versionRepo repository.DocumentVersionRepository,
	cacheService CacheService,
	embeddingService EmbeddingService,
	indexingEnabled bool,
) SearchService {
	return &searchService{
		indexRepo:        indexRepo,
		documentRepo:     documentRepo,
		versionRepo:      versionRepo,
		cacheService:     cacheService,
		embeddingService: embeddingService,
		indexingEnabled:  indexingEnabled,
	}
}

// BuildIndex 构建文档索引
func (s *searchService) BuildIndex(ctx context.Context, documentID, version string) error {
	// 获取文档版本信息
	docVersion, err := s.versionRepo.GetByDocumentIDAndVersion(ctx, documentID, version)
	if err != nil {
		return fmt.Errorf("failed to get document version: %v", err)
	}

	// 检查文档是否已解析完成
	if docVersion.Status != model.DocumentStatusCompleted {
		return fmt.Errorf("document is not ready for indexing, status: %s", docVersion.Status)
	}

	// 获取文档信息
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("failed to get document: %v", err)
	}

	// 解析文档内容并构建索引（会自动处理更新逻辑）
	log.Printf("DEBUG: 开始解析文档内容并构建索引 - 文档ID: %s, 版本: %s", documentID, version)
	indices, err := s.parseAndBuildIndices(document, docVersion)
	if err != nil {
		return fmt.Errorf("failed to parse and build indices: %v", err)
	}
	log.Printf("DEBUG: 解析完成，生成了 %d 个索引 - 文档ID: %s, 版本: %s", len(indices), documentID, version)

	// parseAndBuildIndices已经处理了索引的删除和创建，这里不需要再创建
	log.Printf("Successfully built %d indices for document %s version %s", len(indices), documentID, version)
	return nil
}

// BuildIndexBatch 批量构建索引
func (s *searchService) BuildIndexBatch(ctx context.Context, indices []*model.SearchIndex) error {
	if len(indices) == 0 {
		return nil
	}
	return s.indexRepo.CreateBatch(ctx, indices)
}

// Search 执行搜索
func (s *searchService) Search(ctx context.Context, request *model.SearchRequest) (*model.SearchResponse, error) {
	// 添加调试日志
	log.Printf("DEBUG: Search called with query: %s, type: %s, page: %d, size: %d",
		request.Query, request.SearchType, request.Page, request.Size)

	// 生成缓存键
	cacheKey := searchCacheKey(request.Query, request.SearchType, request.Filters, request.Page, request.Size)

	// 尝试从缓存获取结果
	if cachedResult, found := s.cacheService.Get(cacheKey); found {
		if response, ok := cachedResult.(*model.SearchResponse); ok {
			log.Printf("Search result found in cache for query: %s", request.Query)
			return response, nil
		}
	}

	var indices []*model.SearchIndex
	var total int64
	var err error

	startTime := time.Now()

	// 根据搜索类型执行不同的搜索策略
	switch request.SearchType {
	case "keyword":
		keywords := s.extractKeywords(request.Query)
		indices, total, err = s.indexRepo.SearchByKeywords(ctx, keywords, request.Filters, request.Page, request.Size)
		// 为每个结果计算相关性得分
		for _, index := range indices {
			index.Score = s.calculateRelevanceScore(index, request.Query, "keyword")
		}
	case "semantic":
		// 生成查询向量
		queryVector := s.generateQueryVector(request.Query)
		indices, total, err = s.indexRepo.SearchByVector(ctx, queryVector, request.Filters, request.Page, request.Size)
		// 向量搜索结果已经包含相似度得分，但可以进一步优化
		for _, index := range indices {
			index.Score = s.calculateRelevanceScore(index, request.Query, "semantic")
		}
	case "hybrid":
		// 混合搜索：先关键词搜索，再语义搜索，然后合并结果
		var keywordIndices []*model.SearchIndex
		var semanticIndices []*model.SearchIndex
		var keywordTotal, semanticTotal int64

		// 关键词搜索
		keywords := s.extractKeywords(request.Query)
		keywordIndices, keywordTotal, err = s.indexRepo.SearchByKeywords(ctx, keywords, request.Filters, request.Page, request.Size)
		if err != nil {
			return nil, fmt.Errorf("keyword search failed: %v", err)
		}
		// 为关键词搜索结果计算相关性得分
		for _, index := range keywordIndices {
			index.Score = s.calculateRelevanceScore(index, request.Query, "keyword")
		}

		// 语义搜索
		queryVector := s.generateQueryVector(request.Query)
		semanticIndices, semanticTotal, err = s.indexRepo.SearchByVector(ctx, queryVector, request.Filters, request.Page, request.Size)
		if err != nil {
			return nil, fmt.Errorf("semantic search failed: %v", err)
		}
		// 为语义搜索结果计算相关性得分
		for _, index := range semanticIndices {
			index.Score = s.calculateRelevanceScore(index, request.Query, "semantic")
		}

		// 合并结果
		indices, total = s.mergeSearchResults(keywordIndices, semanticIndices, keywordTotal, semanticTotal)
	default:
		// 默认使用关键词搜索
		keywords := s.extractKeywords(request.Query)
		indices, total, err = s.indexRepo.SearchByKeywords(ctx, keywords, request.Filters, request.Page, request.Size)
		// 为每个结果计算相关性得分
		for _, index := range indices {
			index.Score = s.calculateRelevanceScore(index, request.Query, "keyword")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}

	// 计算搜索耗时
	duration := time.Since(startTime)
	log.Printf("Search completed in %v for query: %s", duration, request.Query)

	// 转换为搜索结果，传递查询词以便在片段中显示上下文
	log.Printf("DEBUG: Found %d indices, total count: %d", len(indices), total)
	results := s.convertToSearchResultsWithQuery(indices, request.Query)
	log.Printf("DEBUG: Converted to %d search results", len(results))

	response := &model.SearchResponse{
		Total: total,
		Items: results,
		Page:  request.Page,
		Size:  request.Size,
	}

	// 智能缓存策略
	s.applyCacheStrategy(cacheKey, response, duration, request)

	return response, nil
}

// GetIndexingStatus 获取索引状态
func (s *searchService) GetIndexingStatus(ctx context.Context, documentID string) (map[string]interface{}, error) {
	return s.indexRepo.GetIndexingStatus(ctx, documentID)
}

// DeleteIndex 删除索引
func (s *searchService) DeleteIndex(ctx context.Context, documentID string) error {
	return s.indexRepo.DeleteByDocumentID(ctx, documentID)
}

// DeleteIndexByVersion 删除指定版本的索引
func (s *searchService) DeleteIndexByVersion(ctx context.Context, documentID, version string) error {
	return s.indexRepo.DeleteByDocumentIDAndVersion(ctx, documentID, version)
}

// ClearCache 清空缓存
func (s *searchService) ClearCache() error {
	return s.cacheService.Clear()
}

// parseAndBuildIndices 解析文档内容并构建索引
func (s *searchService) parseAndBuildIndices(document *model.Document, docVersion *model.DocumentVersion) ([]*model.SearchIndex, error) {
	// 直接使用整个文档内容，不再分段处理
	content := docVersion.Content
	contentLength := len(content)

	// 生成向量
	vectorSlice := s.generateContentVector(content)
	embeddingSlice := s.generateEmbedding(content) // 生成真实嵌入向量

	// 将向量转换为JSON字符串
	vectorJSON, err := json.Marshal(vectorSlice)
	if err != nil {
		log.Printf("Error marshaling vector to JSON: %v", err)
		vectorJSON = []byte("[]")
	}

	// 计算位置信息：从文档开始到结束
	startPos := 0
	endPos := contentLength

	// 创建单个索引条目
	newID := generateID()
	index := &model.SearchIndex{
		ID:            newID,
		DocumentID:    document.ID,
		Version:       docVersion.Version,
		Content:       content,
		ContentType:   "text",
		Section:       document.Name,
		Keywords:      "",                 // 不使用关键词
		Vector:        string(vectorJSON), // 传统向量，以JSON字符串格式存储
		Metadata:      s.buildMetadataWithPosition(document, docVersion, startPos, endPos),
		StartPosition: startPos, // 记录起始位置
		EndPosition:   endPos,   // 记录结束位置
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 只有在 pgvector 扩展可用且生成了有效嵌入向量时才设置 Embedding 字段
	if embeddingSlice != nil {
		index.Embedding = embeddingSlice
	}

	indices := []*model.SearchIndex{index}

	// 删除该文档版本的所有现有索引
	if err := s.indexRepo.DeleteByDocumentIDAndVersion(context.Background(), document.ID, docVersion.Version); err != nil {
		log.Printf("Error deleting existing indices: %v", err)
	}

	// 创建新索引
	log.Printf("DEBUG: 创建索引 - 文档ID: %s, 版本: %s, 位置: %d-%d", document.ID, docVersion.Version, startPos, endPos)
	if err := s.indexRepo.Create(context.Background(), index); err != nil {
		log.Printf("DEBUG: 创建索引失败 - 文档ID: %s, 版本: %s, 错误: %v", document.ID, docVersion.Version, err)
		return nil, fmt.Errorf("failed to create index: %v", err)
	}
	log.Printf("DEBUG: 创建索引成功 - 文档ID: %s, 版本: %s, 位置: %d-%d", document.ID, docVersion.Version, startPos, endPos)

	log.Printf("Successfully built index for document %s version %s with position %d-%d", document.ID, docVersion.Version, startPos, endPos)

	return indices, nil
}

// extractKeywords 从查询中提取关键词，严格匹配
func (s *searchService) extractKeywords(query string) []string {
	// 去除首尾空格，保持原始大小写
	query = strings.TrimSpace(query)
	if query == "" {
		return []string{}
	}

	// 直接按空格分词，不做复杂处理
	words := strings.Fields(query)

	// 过滤空字符串
	var keywords []string
	for _, word := range words {
		if word != "" {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// generateQueryVector 生成查询向量
func (s *searchService) generateQueryVector(query string) string {
	// 生成真实的嵌入向量
	embedding := s.generateEmbedding(query)
	if embedding == nil {
		// 如果嵌入生成失败，使用备用方法
		return s.generateFallbackVector(query)
	}

	// 将向量转换为JSON字符串
	vectorJSON, err := json.Marshal(embedding)
	if err != nil {
		log.Printf("Error marshaling vector to JSON: %v", err)
		return "[]"
	}
	return string(vectorJSON)
}

// generateFallbackVector 生成备用向量（当嵌入服务不可用时）
func (s *searchService) generateFallbackVector(query string) string {
	// 备用向量生成方法
	keywords := s.extractKeywords(query)
	vector := make([]float32, 100) // 假设向量维度为100

	// 简单的哈希映射
	for _, keyword := range keywords {
		hash := simpleHash(keyword)
		// 确保索引为正数
		index := hash % 100
		if index < 0 {
			index = -index
		}
		index = index % 100 // 再次确保在范围内
		vector[index] += 1.0
	}

	// 归一化
	norm := float32(0)
	for _, v := range vector {
		norm += v * v
	}
	if norm > 0 {
		norm = sqrt(norm)
		for i := range vector {
			vector[i] /= norm
		}
	}

	// 将向量转换为JSON字符串
	vectorJSON, err := json.Marshal(vector)
	if err != nil {
		log.Printf("Error marshaling fallback vector to JSON: %v", err)
		return "[]"
	}
	return string(vectorJSON)
}

// generateContentVector 生成内容向量
func (s *searchService) generateContentVector(content string) []float32 {
	// 生成真实的嵌入向量
	embedding := s.generateEmbedding(content)
	if embedding == nil {
		// 如果嵌入生成失败，使用备用方法
		return s.generateFallbackContentVector(content)
	}
	return embedding
}

// generateFallbackContentVector 生成备用内容向量（当嵌入服务不可用时）
func (s *searchService) generateFallbackContentVector(content string) []float32 {
	// 备用向量生成方法
	vector := make([]float32, 100) // 假设向量维度为100

	// 简单的哈希映射
	for i := 0; i < len(vector); i++ {
		// 基于内容位置生成伪随机值
		hash := simpleHash(fmt.Sprintf("%s-%d", content, i))
		// 转换为 -1 到 1 之间的浮点数
		vector[i] = float32(hash%1000)/500.0 - 1.0
	}

	// 归一化
	norm := float32(0)
	for _, v := range vector {
		norm += v * v
	}
	if norm > 0 {
		norm = sqrt(norm)
		for i := range vector {
			vector[i] /= norm
		}
	}

	return vector
}

// generateEmbedding 生成真实嵌入向量
func (s *searchService) generateEmbedding(content string) []float32 {
	// 如果内容为空，返回 nil
	if strings.TrimSpace(content) == "" {
		return nil
	}

	// 使用嵌入服务生成向量
	ctx := context.Background()
	embedding, err := s.embeddingService.GenerateEmbedding(ctx, content)
	if err != nil {
		log.Printf("Error generating embedding with service: %v, falling back to mock service", err)

		// 如果嵌入服务失败，使用模拟服务
		mockService := NewMockEmbeddingService()
		embedding, err = mockService.GenerateEmbedding(ctx, content)
		if err != nil {
			log.Printf("Error generating embedding with mock service: %v", err)
			return nil
		}
	}

	return embedding
}

// mergeSearchResults 合并搜索结果
func (s *searchService) mergeSearchResults(keywordIndices, semanticIndices []*model.SearchIndex, keywordTotal, semanticTotal int64) ([]*model.SearchIndex, int64) {
	// 高级结果合并策略，使用加权评分
	// 关键词搜索权重 0.6，语义搜索权重 0.4

	// 创建结果映射，避免重复
	resultMap := make(map[string]*model.SearchIndex)

	// 添加关键词搜索结果
	for _, idx := range keywordIndices {
		// 关键词搜索结果使用权重 0.6
		idx.Score = idx.Score * 0.6
		resultMap[idx.ID] = idx
	}

	// 添加语义搜索结果
	for _, idx := range semanticIndices {
		if existing, ok := resultMap[idx.ID]; ok {
			// 如果已存在，合并分数，语义搜索权重 0.4
			existing.Score = existing.Score + (idx.Score * 0.4)
		} else {
			// 语义搜索结果使用权重 0.4
			idx.Score = idx.Score * 0.4
			resultMap[idx.ID] = idx
		}
	}

	// 转换为切片
	var results []*model.SearchIndex
	for _, idx := range resultMap {
		results = append(results, idx)
	}

	// 按分数排序
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Score < results[j].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// 计算总数
	// 使用两个搜索结果中较大的总数
	total := keywordTotal
	if semanticTotal > keywordTotal {
		total = semanticTotal
	}

	return results, total
}

// calculateRelevanceScore 计算相关性得分
func (s *searchService) calculateRelevanceScore(index *model.SearchIndex, query string, searchType string) float32 {
	var score float32

	switch searchType {
	case "keyword":
		// 关键词搜索：基于匹配度和出现频率计算得分
		keywords := s.extractKeywords(query)
		content := strings.ToLower(index.Content)
		matchCount := 0

		for _, keyword := range keywords {
			keywordLower := strings.ToLower(keyword)
			count := strings.Count(content, keywordLower)
			if count > 0 {
				matchCount++
			}
		}

		if len(keywords) > 0 {
			score = float32(matchCount) / float32(len(keywords))
		}

	case "semantic":
		// 语义搜索：使用向量相似度作为得分
		score = index.Score
	}

	return score
}

// convertToSearchResults 转换为搜索结果
func (s *searchService) convertToSearchResults(indices []*model.SearchIndex) []model.SearchResult {
	return s.convertToSearchResultsWithQuery(indices, "")
}

// convertToSearchResultsWithQuery 转换为搜索结果（带查询词）
func (s *searchService) convertToSearchResultsWithQuery(indices []*model.SearchIndex, query string) []model.SearchResult {
	var results []model.SearchResult

	for _, idx := range indices {
		// 计算片段在原始文档中的位置
		startPos := 0
		endPos := len(idx.Content)

		// 如果搜索索引中存储了位置信息，使用它
		if idx.StartPosition > 0 {
			startPos = idx.StartPosition
		}
		if idx.EndPosition > 0 && idx.EndPosition <= len(idx.Content) {
			endPos = idx.EndPosition
		}

		// 生成内容片段，优先显示包含查询词的上下文
		snippet := s.generateSnippetWithQuery(idx.Content, query, 300)

		// 解析元数据
		var metadata map[string]interface{}
		library := ""
		if idx.Metadata != "" {
			if err := json.Unmarshal([]byte(idx.Metadata), &metadata); err != nil {
				log.Printf("Error unmarshaling metadata: %v", err)
				metadata = make(map[string]interface{})
			} else {
				// 从元数据中提取库信息
				if libraryVal, ok := metadata["document_library"].(string); ok {
					library = libraryVal
				}
			}
		} else {
			metadata = make(map[string]interface{})
		}

		// 在元数据中包含位置信息
		if startPos > 0 || endPos > 0 {
			metadata["start_position"] = startPos
			metadata["end_position"] = endPos
		}

		result := model.SearchResult{
			ID:          idx.ID,
			DocumentID:  idx.DocumentID,
			Version:     idx.Version,
			Library:     library,
			Content:     idx.Content,
			Snippet:     snippet,
			Score:       idx.Score,
			ContentType: idx.ContentType,
			Section:     idx.Section,
			Metadata:    metadata,
		}

		results = append(results, result)
	}

	return results
}

// generateSnippet 生成内容片段
func (s *searchService) generateSnippet(content string, maxLength int) string {
	return s.generateSnippetWithQuery(content, "", maxLength)
}

// generateSnippetWithQuery 生成包含查询词的内容片段
func (s *searchService) generateSnippetWithQuery(content, query string, maxLength int) string {
	if len(content) <= maxLength {
		return content
	}

	// 如果没有查询词，返回开头部分
	if query == "" {
		return s.truncateToSmartLength(content, maxLength)
	}

	// 在内容中查找查询词
	lowerContent := strings.ToLower(content)
	lowerQuery := strings.ToLower(query)

	// 尝试找到查询词的位置
	idx := strings.Index(lowerContent, lowerQuery)
	if idx == -1 {
		// 找不到查询词，返回开头部分
		return s.truncateToSmartLength(content, maxLength)
	}

	// 智能计算片段起始位置，确保包含完整的markdown代码块
	start := s.findOptimalStartPosition(content, idx, maxLength)

	// 计算结束位置
	end := start + maxLength
	if end > len(content) {
		end = len(content)
		start = end - maxLength
		if start < 0 {
			start = 0
		}
	}

	snippet := content[start:end]

	// 添加省略号标记截断
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(content) {
		snippet = snippet + "..."
	}

	return snippet
}

// findOptimalStartPosition 查找最优的片段起始位置
// 确保包含完整的markdown代码块（避免在代码块中间截断）
func (s *searchService) findOptimalStartPosition(content string, matchIndex, maxLength int) int {
	// 默认从匹配位置开始，向前50%
	halfLength := maxLength / 2
	start := matchIndex - halfLength

	// 检查向前是否包含未闭合的markdown代码块
	if start > 0 {
		// 向前查找最近的换行符，确保不在代码块中间截断
		safeStart := s.findSafeBreakPoint(content, start, -1, 20)
		if safeStart != -1 {
			start = safeStart
		}
	}

	// 确保不超出内容长度
	if start < 0 {
		start = 0
	}

	// 检查后边界，如果包含未闭合的代码块，向后扩展
	if s.hasUnclosedCodeBlock(content, start+maxLength) {
		// 尝试向后扩展以包含完整的代码块
		extendedEnd := s.findCompleteCodeBlockEnd(content, start+maxLength)
		if extendedEnd != -1 && (extendedEnd-start) <= maxLength*2 {
			// 可以向后扩展，但不超出maxLength的2倍
			return start
		}
	}

	return start
}

// findSafeBreakPoint 在指定范围内查找安全的断点（换行符）
func (s *searchService) findSafeBreakPoint(content string, start int, direction, maxLookback int) int {
	contentLen := len(content)

	if direction < 0 {
		// 向前查找
		begin := start + direction
		end := start
		for i := begin; i >= 0 && i >= end-maxLookback; i-- {
			if content[i] == '\n' {
				return i + 1 // 返回换行符后面的位置
			}
		}
	} else {
		// 向后查找（不使用，因为要向后扩展）
		for i := start; i < contentLen && i < start+direction+maxLookback; i++ {
			if content[i] == '\n' {
				return i + 1 // 返回换行符后面的位置
			}
		}
	}

	return -1
}

// hasUnclosedCodeBlock 检查指定位置后是否有未闭合的代码块
func (s *searchService) hasUnclosedCodeBlock(content string, endPos int) bool {
	slice := content[:endPos]

	// 检查是否包含未闭合的代码块标记
	// 代码块开始：```, ```bash, ```python, ```json, ```javascript, ```code, ```text,
	// 代码块开始：~~~python, ~~~bash, ~~~shell,
	codeStartPatterns := []string{
		"```python", "```bash", "```shell", "```json",
		"```javascript", "```js", "```code", "```text",
		"~~~python", "~~~bash", "~~~shell",
	}

	// 检查每个代码块开始标记
	for i := len(slice) - 1; i >= 0; i-- {
		// 查找最近的代码块开始
		hasStart := false
		hasEnd := false

		for _, pattern := range codeStartPatterns {
			if strings.HasPrefix(slice[i:], pattern) || strings.Contains(slice[i:], pattern) {
				hasStart = true
				// 查找对应的结束标记
				endPattern := strings.TrimPrefix(pattern, "```")
				if idx := strings.Index(slice[i:], endPattern+"```"); idx != -1 {
					hasEnd = true
					break
				}
			}
		}

		// 如果有开始但没有结束，说明有未闭合的代码块
		if hasStart && !hasEnd {
			return true
		}
	}

	return false
}

// findCompleteCodeBlockEnd 查找最近的代码块结束标记
func (s *searchService) findCompleteCodeBlockEnd(content string, start int) int {
	slice := content[start:]

	// 代码块结束模式
	endPatterns := []string{"```", "~~~"}

	minEndPos := -1

	for i := 0; i < len(slice); i++ {
		// 跳过前几个字符，避免匹配代码块开始后的结束标记
		if i > 10 {
			break
		}

		// 查找代码块结束
		for _, pattern := range endPatterns {
			if strings.HasPrefix(slice[i:], pattern) && i > 5 {
				// 确保不是代码块开始后的结束标记
				if !strings.HasPrefix(slice[i:], "```") ||
					(i > 10 && strings.Count(slice[:i+10], "```") > 0) {
					if minEndPos == -1 || i < minEndPos {
						minEndPos = i
					}
				}
			}
		}
	}

	if minEndPos != -1 {
		return start + minEndPos + 3 // +3 是为了包含 ``` 或 ~~~ 结束标记
	}

	return -1
}

// truncateToSmartLength 智能截断到合适长度
func (s *searchService) truncateToSmartLength(content string, maxLength int) string {
	if len(content) <= maxLength {
		return content
	}
	return content[:maxLength] + "..."
}

// buildMetadata 构建元数据（简化版本，不使用分段）
func (s *searchService) buildMetadata(document *model.Document, docVersion *model.DocumentVersion) string {
	metadata := map[string]interface{}{
		"document_name":    document.Name,
		"document_type":    document.Type,
		"document_library": document.Library,
		"version":          docVersion.Version,
		"section_type":     "text",
		"section_title":    document.Name,
	}

	// 使用json.Marshal进行正确的JSON序列化
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		log.Printf("Error marshaling metadata to JSON: %v", err)
		return "{}"
	}
	return string(metadataJSON)
}

// generateID 生成唯一ID
func generateID() string {
	return fmt.Sprintf("%s", uuid.New().String())
}

// simpleHash 简单哈希函数
func simpleHash(s string) int {
	hash := 5381
	for _, c := range s {
		hash = ((hash << 5) + hash) + int(c)
	}
	return hash
}

// sqrt 计算平方根
func sqrt(x float32) float32 {
	var z float32 = 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

// applyCacheStrategy 智能缓存策略
func (s *searchService) applyCacheStrategy(cacheKey string, response *model.SearchResponse, duration time.Duration, request *model.SearchRequest) {
	// 根据搜索类型和性能决定是否缓存
	shouldCache := false
	ttl := 5 * time.Minute // 默认TTL

	switch request.SearchType {
	case "keyword":
		// 关键词搜索速度快，可以缓存
		shouldCache = true
		ttl = 5 * time.Minute
	case "semantic":
		// 语义搜索计算成本高，应该缓存更长时间
		shouldCache = true
		ttl = 10 * time.Minute
	case "hybrid":
		// 混合搜索计算成本最高，应该缓存最长时间
		shouldCache = true
		ttl = 15 * time.Minute
	}

	// 如果搜索结果为空，缓存时间较短
	if response.Total == 0 || len(response.Items) == 0 {
		ttl = 2 * time.Minute
	}

	// 如果搜索耗时较长，应该缓存更长时间
	if duration > 500*time.Millisecond {
		ttl = ttl * 2
	}

	if shouldCache {
		if err := s.cacheService.Set(cacheKey, response, ttl); err != nil {
			log.Printf("Failed to cache search result: %v", err)
		} else {
			log.Printf("Search result cached for query: %s, TTL: %v", request.Query, ttl)
		}
	}
}

// buildSimpleMetadata 构建简化的元数据
func (s *searchService) buildSimpleMetadata(document *model.Document, docVersion *model.DocumentVersion) string {
	metadata := map[string]interface{}{
		"document_name":    document.Name,
		"document_type":    document.Type,
		"document_library": document.Library,
		"version":          docVersion.Version,
	}

	// 使用json.Marshal进行正确的JSON序列化
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		log.Printf("Error marshaling metadata to JSON: %v", err)
		return "{}"
	}
	return string(metadataJSON)
}

// buildMetadataWithPosition 构建包含位置信息的元数据
func (s *searchService) buildMetadataWithPosition(document *model.Document, docVersion *model.DocumentVersion, startPos, endPos int) string {
	metadata := map[string]interface{}{
		"document_name":    document.Name,
		"document_type":    document.Type,
		"document_library": document.Library,
		"version":          docVersion.Version,
		"start_position":   float64(startPos), // JSON中数字默认为float64
		"end_position":     float64(endPos),
		"content_length":   float64(endPos - startPos),
	}

	// 使用json.Marshal进行正确的JSON序列化
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		log.Printf("Error marshaling metadata to JSON: %v", err)
		return "{}"
	}
	return string(metadataJSON)
}
