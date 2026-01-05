package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/repository"
)

// Token使用限制配置常量
const (
	// DefaultContentMaxLength 默认内容最大长度（字符数）
	DefaultContentMaxLength = 30000
	// WarningContentLength 内容警告阈值（字符数）
	WarningContentLength = 60000
	// DefaultSearchResultLength 默认搜索结果长度限制（字符数）
	DefaultSearchResultLength = 1000
	// SearchResultMaxLength 搜索结果最大长度限制（字符数）
	SearchResultMaxLength = 2000
)

// MCPService MCP服务接口
type MCPService interface {
	// 处理MCP请求
	HandleRequest(ctx context.Context, req *model.MCPRequest, apiKey string) (*model.MCPResponse, error)

	// 初始化MCP连接
	Initialize(ctx context.Context, params *model.MCPInitializeParams) (*model.MCPInitializeResult, error)

	// 获取工具列表
	ListTools(ctx context.Context, params *model.MCPToolListParams) (*model.MCPToolListResult, error)

	// 调用工具
	CallTool(ctx context.Context, params *model.MCPToolCallParams) (*model.MCPToolResult, error)

	// 验证API密钥
	ValidateAPIKey(apiKey string) (*model.MCPAPIKey, error)

	// 创建API密钥
	CreateAPIKey(ctx context.Context, name, userID string, expiresAt *time.Time) (*model.MCPAPIKey, error)

	// 获取API密钥列表
	GetAPIKeys(ctx context.Context, userID string) ([]model.MCPAPIKey, error)

	// 删除API密钥
	DeleteAPIKey(ctx context.Context, keyID string) error

	// 根据密钥ID获取API密钥
	ValidateAPIKeyByKeyID(keyID string) (*model.MCPAPIKey, error)

	// 更新API密钥最后使用时间
	UpdateAPIKeyLastUsed(ctx context.Context, keyID string) error
}

// mcpService MCP服务实现
type mcpService struct {
	db              *gorm.DB
	searchService   SearchService
	documentService DocumentService
	versionRepo     repository.DocumentVersionRepository
	indexRepo       repository.SearchIndexRepository
}

// NewMCPService 创建MCP服务实例
func NewMCPService(db *gorm.DB, searchService SearchService, documentService DocumentService, versionRepo repository.DocumentVersionRepository, indexRepo repository.SearchIndexRepository) MCPService {
	return &mcpService{
		db:              db,
		searchService:   searchService,
		documentService: documentService,
		versionRepo:     versionRepo,
		indexRepo:       indexRepo,
	}
}

// HandleRequest 处理MCP请求
func (s *mcpService) HandleRequest(ctx context.Context, req *model.MCPRequest, apiKey string) (*model.MCPResponse, error) {
	// 验证API密钥
	_, err := s.ValidateAPIKey(apiKey)
	if err != nil {
		return s.createErrorResponse(req.ID, req.Method, -32600, "Invalid API key", nil)
	}

	switch req.Method {
	case "initialize":
		return s.handleInitialize(ctx, req)
	case "tools/list":
		return s.handleListTools(ctx, req)
	case "tools/call":
		return s.handleCallTool(ctx, req)
	default:
		return s.createErrorResponse(req.ID, req.Method, -32601, "Method not found", nil)
	}
}

// Initialize 初始化MCP连接
func (s *mcpService) Initialize(ctx context.Context, params *model.MCPInitializeParams) (*model.MCPInitializeResult, error) {
	log.Printf("Initializing MCP connection for client: %s %s", params.ClientInfo.Name, params.ClientInfo.Version)

	result := &model.MCPInitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: model.MCPCapabilities{
			Tools: &model.MCPToolsCapability{
				ListChanged: true,
			},
		},
		ServerInfo: model.MCPServerInfo{
			Name:    "AI技术文档库",
			Version: "1.0.0",
		},
	}

	return result, nil
}

// ListTools 获取工具列表
func (s *mcpService) ListTools(ctx context.Context, params *model.MCPToolListParams) (*model.MCPToolListResult, error) {
	tools := []model.MCPTool{
		{
			Name:        "search_documents",
			Description: "搜索技术文档，支持关键词搜索和语义搜索",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "搜索查询关键词",
					},
					"types": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
						"description": "文档类型过滤器，如 pdf, docx, markdown等",
					},
					"version": map[string]interface{}{
						"type":        "string",
						"description": "文档版本过滤器",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "返回结果数量限制，默认为10",
					},
					"content_length": map[string]interface{}{
						"type":        "integer",
						"description": fmt.Sprintf("每个搜索结果的内容片段最大字符数，默认为%d，最大为%d", DefaultSearchResultLength, SearchResultMaxLength),
					},
				},
				"required": []string{"query"},
			},
		},
		{
			Name:        "get_documents_by_library",
			Description: "根据所属库名称获取文档列表",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"library": map[string]interface{}{
						"type":        "string",
						"description": "库名称",
					},
					"page": map[string]interface{}{
						"type":        "integer",
						"description": "页码，默认为1",
					},
					"size": map[string]interface{}{
						"type":        "integer",
						"description": "每页数量，默认为10",
					},
				},
				"required": []string{"library"},
			},
		},
		{
			Name:        "get_document_content",
			Description: "获取指定文档的详细内容（可使用版本ID或文档ID，支持自定义位置范围）",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"document_id": map[string]interface{}{
						"type":        "string",
						"description": "文档ID或版本ID（推荐使用get_documents_by_library返回的版本ID）",
					},
					"version": map[string]interface{}{
						"type":        "string",
						"description": "文档版本号（仅在提供文档ID时有效，如果未指定则使用最新版本）",
					},
					"start_position": map[string]interface{}{
						"type":        "integer",
						"description": "起始位置（字符位置，从文档的第几个字符开始获取内容）。如果未指定，默认从第0个字符开始",
					},
					"end_position": map[string]interface{}{
						"type":        "integer",
						"description": "结束位置（字符位置，获取到文档的第几个字符为止）。如果未指定，默认获取到内容结束",
					},
					"query": map[string]interface{}{
						"type":        "string",
						"description": "搜索关键词（用于定位内容位置，将围绕关键词前后扩展内容。注意：如果指定了start_position和end_position，优先使用位置参数）",
					},
					"content_length": map[string]interface{}{
						"type":        "integer",
						"description": fmt.Sprintf("返回内容的最大字符数，默认为%d，最大为%d。当启用智能截断时，会在markdown标题边界处完整补全，避免在同一标题内容中间断开", DefaultContentMaxLength, WarningContentLength),
					},
					"smart_truncate": map[string]interface{}{
						"type":        "boolean",
						"description": "是否启用智能截断模式（按markdown标题完整补全），默认为true。当为true时，如果到达长度限制时会继续补全直到下一个大标题，确保不会在同一个标题的内容中间截断",
					},
				},
				"required": []string{"document_id"},
			},
		},
	}

	return &model.MCPToolListResult{
		Tools: tools,
	}, nil
}

// CallTool 调用工具
func (s *mcpService) CallTool(ctx context.Context, params *model.MCPToolCallParams) (*model.MCPToolResult, error) {
	switch params.Name {
	case "search_documents":
		return s.searchDocumentsTool(ctx, params.Arguments)
	case "get_documents_by_library":
		return s.getDocumentsByLibraryTool(ctx, params.Arguments)
	case "get_document_content":
		return s.getDocumentContentTool(ctx, params.Arguments)
	default:
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: fmt.Sprintf("未知工具: %s", params.Name),
				},
			},
			IsError: true,
		}, nil
	}
}

// ValidateAPIKey 验证API密钥
func (s *mcpService) ValidateAPIKey(apiKey string) (*model.MCPAPIKey, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API密钥不能为空")
	}

	var key model.MCPAPIKey
	err := s.db.Where("key = ? AND enabled = ?", apiKey, true).First(&key).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("无效的API密钥")
		}
		return nil, fmt.Errorf("查询API密钥失败: %v", err)
	}

	// 检查密钥是否过期
	if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
		return nil, fmt.Errorf("API密钥已过期")
	}

	// 检查用户是否激活
	var user model.User
	if err := s.db.Where("id = ?", key.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("用户账户已被禁用")
	}

	return &key, nil
}

// CreateAPIKey 创建API密钥
func (s *mcpService) CreateAPIKey(ctx context.Context, name, userID string, expiresAt *time.Time) (*model.MCPAPIKey, error) {
	// 验证用户是否存在
	var user model.User
	if err := s.db.Where("id = ? AND is_active = ?", userID, true).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在或已被禁用")
	}

	// 生成API密钥
	key := uuid.New().String()

	apiKey := &model.MCPAPIKey{
		Name:      name,
		Key:       key,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	err := s.db.Create(apiKey).Error
	if err != nil {
		return nil, fmt.Errorf("创建API密钥失败: %v", err)
	}

	return apiKey, nil
}

// GetAPIKeys 获取API密钥列表
func (s *mcpService) GetAPIKeys(ctx context.Context, userID string) ([]model.MCPAPIKey, error) {
	var keys []model.MCPAPIKey
	err := s.db.Where("user_id = ? AND enabled = ?", userID, true).Find(&keys).Error
	if err != nil {
		return nil, fmt.Errorf("获取API密钥列表失败: %v", err)
	}

	return keys, nil
}

// DeleteAPIKey 删除API密钥（软删除）
func (s *mcpService) DeleteAPIKey(ctx context.Context, keyID string) error {
	err := s.db.Model(&model.MCPAPIKey{}).Where("id = ?", keyID).Update("enabled", false).Error
	if err != nil {
		return fmt.Errorf("删除API密钥失败: %v", err)
	}

	return nil
}

// ValidateAPIKeyByKeyID 根据密钥ID获取API密钥
func (s *mcpService) ValidateAPIKeyByKeyID(keyID string) (*model.MCPAPIKey, error) {
	if keyID == "" {
		return nil, fmt.Errorf("API密钥ID不能为空")
	}

	var key model.MCPAPIKey
	err := s.db.Where("id = ? AND enabled = ?", keyID, true).First(&key).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("API密钥不存在")
		}
		return nil, fmt.Errorf("查询API密钥失败: %v", err)
	}

	// 检查密钥是否过期
	if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
		return nil, fmt.Errorf("API密钥已过期")
	}

	// 检查用户是否激活
	var user model.User
	if err := s.db.Where("id = ?", key.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("用户账户已被禁用")
	}

	return &key, nil
}

// UpdateAPIKeyLastUsed 更新API密钥最后使用时间
func (s *mcpService) UpdateAPIKeyLastUsed(ctx context.Context, keyID string) error {
	now := time.Now()
	err := s.db.Model(&model.MCPAPIKey{}).Where("id = ?", keyID).Update("last_used", now).Error
	if err != nil {
		return fmt.Errorf("更新API密钥最后使用时间失败: %v", err)
	}

	return nil
}

// handleInitialize 处理初始化请求
func (s *mcpService) handleInitialize(ctx context.Context, req *model.MCPRequest) (*model.MCPResponse, error) {
	// 解析参数
	var params model.MCPInitializeParams
	if req.Params != nil {
		paramBytes, err := json.Marshal(req.Params)
		if err != nil {
			return s.createErrorResponse(req.ID, "initialize", -32602, "Invalid params", nil)
		}

		err = json.Unmarshal(paramBytes, &params)
		if err != nil {
			return s.createErrorResponse(req.ID, "initialize", -32602, "Invalid params", nil)
		}
	}

	// 调用初始化方法
	result, err := s.Initialize(ctx, &params)
	if err != nil {
		return s.createErrorResponse(req.ID, "initialize", -32603, "Internal error", err.Error())
	}

	return s.createSuccessResponse(req.ID, "initialize", result)
}

// handleListTools 处理工具列表请求
func (s *mcpService) handleListTools(ctx context.Context, req *model.MCPRequest) (*model.MCPResponse, error) {
	// 解析参数
	var params model.MCPToolListParams
	if req.Params != nil {
		paramBytes, err := json.Marshal(req.Params)
		if err != nil {
			return s.createErrorResponse(req.ID, "tools/list", -32602, "Invalid params", nil)
		}

		err = json.Unmarshal(paramBytes, &params)
		if err != nil {
			return s.createErrorResponse(req.ID, "tools/list", -32602, "Invalid params", nil)
		}
	}

	// 调用工具列表方法
	result, err := s.ListTools(ctx, &params)
	if err != nil {
		return s.createErrorResponse(req.ID, "tools/list", -32603, "Internal error", err.Error())
	}

	return s.createSuccessResponse(req.ID, "tools/list", result)
}

// handleCallTool 处理工具调用请求
func (s *mcpService) handleCallTool(ctx context.Context, req *model.MCPRequest) (*model.MCPResponse, error) {
	// 解析参数
	var params model.MCPToolCallParams
	if req.Params != nil {
		paramBytes, err := json.Marshal(req.Params)
		if err != nil {
			return s.createErrorResponse(req.ID, "tools/call", -32602, "Invalid params", nil)
		}

		err = json.Unmarshal(paramBytes, &params)
		if err != nil {
			return s.createErrorResponse(req.ID, "tools/call", -32602, "Invalid params", nil)
		}
	}

	// 调用工具方法
	result, err := s.CallTool(ctx, &params)
	if err != nil {
		return s.createErrorResponse(req.ID, "tools/call", -32603, "Internal error", err.Error())
	}

	return s.createSuccessResponse(req.ID, "tools/call", result)
}

// searchDocumentsTool 搜索文档工具
func (s *mcpService) searchDocumentsTool(ctx context.Context, args map[string]interface{}) (*model.MCPToolResult, error) {
	// 解析参数
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: "搜索查询不能为空",
				},
			},
			IsError: true,
		}, nil
	}

	var types []string
	if typesArg, ok := args["types"].([]interface{}); ok {
		for _, t := range typesArg {
			if typeStr, ok := t.(string); ok {
				types = append(types, typeStr)
			}
		}
	}

	version, _ := args["version"].(string)

	limit := 10 // 默认限制
	if limitArg, ok := args["limit"].(float64); ok {
		limit = int(limitArg)
	}

	// 解析内容长度限制参数，默认使用 DefaultSearchResultLength
	resultContentLength := DefaultSearchResultLength
	if contentLengthArg, ok := args["content_length"].(float64); ok {
		length := int(contentLengthArg)
		// 确保长度在合理范围内
		if length > 0 && length <= SearchResultMaxLength {
			resultContentLength = length
		} else if length > SearchResultMaxLength {
			log.Printf("WARNING: Requested content_length %d exceeds max limit %d, using max limit", length, SearchResultMaxLength)
			resultContentLength = SearchResultMaxLength
		}
	}

	// 构建搜索请求
	filters := make(map[string]interface{})
	if len(types) > 0 {
		filters["types"] = types
	}
	if version != "" {
		filters["version"] = version
	}

	searchRequest := &model.SearchRequest{
		Query:      query,
		Filters:    filters,
		Page:       1,
		Size:       limit,
		SearchType: "hybrid", // 默认使用混合搜索
	}

	// 调用搜索服务
	searchResult, err := s.searchService.Search(ctx, searchRequest)
	if err != nil {
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: fmt.Sprintf("搜索失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// 转换搜索结果为MCP格式
	var documents []model.MCPSearchDocument
	for _, item := range searchResult.Items {
		// 计算关键词在文档内容中的实际位置
		actualStartPos := 0
		actualEndPos := len(item.Content)

		if query != "" && item.Content != "" {
			// 在内容中查找查询词的位置
			queryLower := strings.ToLower(query)
			contentLower := strings.ToLower(item.Content)

			// 查找第一个匹配位置
			matchIndex := strings.Index(contentLower, queryLower)
			if matchIndex != -1 {
				// 找到匹配，计算精确的起始和结束位置
				actualStartPos = matchIndex
				actualEndPos = matchIndex + len(query)
			}
		}

		// 获取文档名称（优先使用Title，如果Title为空则使用Section）
		documentName := item.Title
		if documentName == "" {
			documentName = item.Section
		}

		documents = append(documents, model.MCPSearchDocument{
			ID:            item.ID,         // 搜索索引ID，每个代码片段的唯一标识
			DocumentID:    item.DocumentID, // 文档ID，用于获取完整文档
			Name:          documentName,
			Type:          item.ContentType, // 直接使用内容类型
			Version:       item.Version,
			Library:       item.Library, // 添加所属库
			Score:         float64(item.Score),
			Content:       item.Content,
			Snippet:       item.Snippet, // 包含查询关键词的上下文片段
			ContentType:   item.ContentType,
			Section:       item.Section,
			StartPosition: actualStartPos, // 关键词在内容中的实际起始位置
			EndPosition:   actualEndPos,   // 关键词在内容中的实际结束位置
		})
	}

	// 构造结果文本
	resultText := fmt.Sprintf("搜索查询: %s\n找到 %d 个相关文档:\n\n", query, len(documents))
	totalTokens := 0
	for i, doc := range documents {
		resultText += fmt.Sprintf("%d. %s (版本: %s, 类型: %s)\n", i+1, doc.Name, doc.Version, doc.Type)
		resultText += fmt.Sprintf("   文档ID: %s\n", doc.ID) // 添加文档ID，方便后续调用get_document_content
		resultText += fmt.Sprintf("   所属库: %s\n", doc.Library)
		resultText += fmt.Sprintf("   相关度: %.2f\n", doc.Score)
		// 显示位置信息
		if doc.StartPosition > 0 || doc.EndPosition > 0 {
			resultText += fmt.Sprintf("   起始位置: %d 字符\n", doc.StartPosition)
			resultText += fmt.Sprintf("   结束位置: %d 字符\n", doc.EndPosition)
		}
		// 使用Snippet（包含查询关键词的上下文片段）
		snippet := doc.Snippet
		// 使用配置的长度限制截断内容
		truncatedSnippet := s.truncateText(snippet, resultContentLength)
		tokens := s.estimateTokens(truncatedSnippet)
		totalTokens += tokens
		resultText += fmt.Sprintf("   内容片段: %s\n   估算Token数: %d\n\n", truncatedSnippet, tokens)
	}
	// 记录总Token使用情况
	log.Printf("INFO: Search results - totalDocuments: %d, totalEstimatedTokens: %d, avgTokensPerDoc: %d",
		len(documents), totalTokens, totalTokens/len(documents))

	return &model.MCPToolResult{
		Content: []interface{}{
			model.MCPTextContent{
				Type: "text",
				Text: resultText,
			},
		},
		IsError: false,
	}, nil
}

// getDocumentsByLibraryTool 根据库获取文档列表工具
func (s *mcpService) getDocumentsByLibraryTool(ctx context.Context, args map[string]interface{}) (*model.MCPToolResult, error) {
	// 解析参数
	library, ok := args["library"].(string)
	if !ok || library == "" {
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: "库名称不能为空",
				},
			},
			IsError: true,
		}, nil
	}

	page := 1
	if pageArg, ok := args["page"].(float64); ok && pageArg > 0 {
		page = int(pageArg)
	}

	size := 10
	if sizeArg, ok := args["size"].(float64); ok && sizeArg > 0 {
		size = int(sizeArg)
	}

	log.Printf("DEBUG: getDocumentsByLibraryTool - library=%s, page=%d, size=%d", library, page, size)

	// 首先获取该库的所有文档
	filters := map[string]interface{}{}
	filters["library"] = library
	documents, _, err := s.documentService.GetDocuments(ctx, 1, 100, filters)

	if err != nil {
		log.Printf("ERROR: getDocumentsByLibraryTool - Failed to get documents: %v", err)
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: fmt.Sprintf("获取文档列表失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// 收集所有版本
	var allVersions []*model.DocumentVersion
	for _, doc := range documents {
		versions, err := s.documentService.GetDocumentVersions(ctx, doc.ID)
		if err != nil {
			log.Printf("WARNING: getDocumentsByLibraryTool - Failed to get versions for document %s: %v", doc.ID, err)
			continue
		}
		allVersions = append(allVersions, versions...)
	}

	// 计算总数和分页
	total := int64(len(allVersions))
	totalPages := (total + int64(size) - 1) / int64(size)

	// 计算偏移量
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if startIndex >= len(allVersions) {
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: fmt.Sprintf("页码超出范围，共 %d 页", totalPages),
				},
			},
			IsError: true,
		}, nil
	}
	if endIndex > len(allVersions) {
		endIndex = len(allVersions)
	}

	// 获取当前页的版本
	pageVersions := allVersions[startIndex:endIndex]

	log.Printf("DEBUG: getDocumentsByLibraryTool - totalVersions=%d, pageVersions=%d, totalPages=%d", total, len(pageVersions), totalPages)

	// 构造结果文本
	resultText := fmt.Sprintf("库: %s\n找到 %d 个文档版本 (第 %d 页, 每页 %d 个):\n\n", library, total, page, size)

	for i, version := range pageVersions {
		// 获取版本对应的文档信息
		doc, err := s.documentService.GetDocument(ctx, version.DocumentID)
		if err != nil {
			log.Printf("WARNING: getDocumentsByLibraryTool - Failed to get document %s: %v", version.DocumentID, err)
			resultText += fmt.Sprintf("%d. 文档版本\n", (page-1)*size+i+1)
			resultText += fmt.Sprintf("   版本ID: %s\n", version.ID)
			resultText += fmt.Sprintf("   版本号: %s\n", version.Version)
			resultText += fmt.Sprintf("   状态: %s\n", version.Status)
			resultText += fmt.Sprintf("   文件大小: %d bytes\n", version.FileSize)
			resultText += fmt.Sprintf("   文档信息获取失败\n\n")
		} else {
			resultText += fmt.Sprintf("%d. %s (版本: %s)\n", (page-1)*size+i+1, doc.Name, version.Version)
			resultText += fmt.Sprintf("   文档ID: %s\n", doc.ID)
			resultText += fmt.Sprintf("   版本ID: %s\n", version.ID)
			resultText += fmt.Sprintf("   版本号: %s\n", version.Version)
			resultText += fmt.Sprintf("   类型: %s\n", doc.Type)
			resultText += fmt.Sprintf("   状态: %s\n", version.Status)
			resultText += fmt.Sprintf("   文件大小: %d bytes\n", version.FileSize)
			resultText += fmt.Sprintf("   创建时间: %s\n", version.CreatedAt.Format("2006-01-02 15:04:05"))
			resultText += fmt.Sprintf("   描述: %s\n\n", doc.Description)
		}
	}

	// 添加分页信息提示
	if totalPages > 1 {
		resultText += fmt.Sprintf("\n分页信息: 共 %d 页，当前第 %d 页\n", totalPages, page)
		if page < int(totalPages) {
			resultText += fmt.Sprintf("提示: 可以使用 page=%d 查看下一页\n", page+1)
		}
	}

	log.Printf("INFO: getDocumentsByLibraryTool - Successfully returned %d versions for library %s", len(pageVersions), library)

	return &model.MCPToolResult{
		Content: []interface{}{
			model.MCPTextContent{
				Type: "text",
				Text: resultText,
			},
		},
		IsError: false,
	}, nil
}

// getDocumentContentTool 获取文档内容工具
func (s *mcpService) getDocumentContentTool(ctx context.Context, args map[string]interface{}) (*model.MCPToolResult, error) {
	// 解析参数 - document_id 可以是文档ID或版本ID
	documentID, ok := args["document_id"].(string)
	if !ok || documentID == "" {
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: "文档ID或版本ID不能为空",
				},
			},
			IsError: true,
		}, nil
	}

	version, _ := args["version"].(string)

	// 解析内容长度限制参数，默认使用 DefaultContentMaxLength
	maxContentLength := DefaultContentMaxLength
	if contentLengthArg, ok := args["content_length"].(float64); ok {
		length := int(contentLengthArg)
		// 确保长度在合理范围内
		if length > 0 && length <= WarningContentLength {
			maxContentLength = length
		} else if length > WarningContentLength {
			log.Printf("WARNING: Requested content_length %d exceeds warning threshold %d, using warning threshold", length, WarningContentLength)
			maxContentLength = WarningContentLength
		}
	}

	// 解析智能截断参数，默认为true
	smartTruncate := true
	if smartTruncateArg, ok := args["smart_truncate"].(bool); ok {
		smartTruncate = smartTruncateArg
	}

	// 解析自定义位置参数
	customStartPos := -1
	customEndPos := -1

	if startPosArg, ok := args["start_position"].(float64); ok {
		customStartPos = int(startPosArg)
	}

	if endPosArg, ok := args["end_position"].(float64); ok {
		customEndPos = int(endPosArg)
	}

	// 解析查询关键词
	query, _ := args["query"].(string)

	log.Printf("DEBUG: getDocumentContentTool called with documentID=%s (可能为版本ID), version='%s', maxContentLength=%d, smartTruncate=%v, startPos=%d, endPos=%d, query='%s'",
		documentID, version, maxContentLength, smartTruncate, customStartPos, customEndPos, query)

	// 首先尝试作为搜索索引ID（片段ID）获取片段内容
	docIndex, err := s.indexRepo.GetByID(ctx, documentID)
	if err == nil && docIndex != nil {
		// 成功作为搜索索引ID获取（片段）
		log.Printf("DEBUG: Retrieved search index by ID: %s (fragment content)", documentID)
		content := docIndex.Content
		originalLength := len(content)

		// 确定要获取的内容范围
		actualStartPos := 0
		actualEndPos := originalLength

		// 优先使用自定义位置参数
		if customStartPos >= 0 {
			actualStartPos = customStartPos
			// 如果指定了起始位置但未指定结束位置，使用到内容末尾
			if customEndPos < 0 {
				actualEndPos = originalLength
			} else {
				actualEndPos = customEndPos
			}
		}

		// 验证位置参数的有效性
		if actualStartPos < 0 {
			actualStartPos = 0
		}
		if actualEndPos > originalLength {
			actualEndPos = originalLength
		}
		if actualEndPos <= actualStartPos {
			actualEndPos = actualStartPos + 1
		}

		// 提取指定位置的内容
		extractEnd := actualEndPos
		if extractEnd > originalLength {
			extractEnd = originalLength
		}
		extractedContent := content[actualStartPos:extractEnd]

		// 检查内容长度并记录警告
		if len(extractedContent) > WarningContentLength {
			log.Printf("WARNING: Document fragment content length %d exceeds warning threshold %d", len(extractedContent), WarningContentLength)
		}

		// 检查是否需要智能截断（根据长度）
		truncatedContent := extractedContent
		isTruncated := false

		if len(extractedContent) > maxContentLength {
			if smartTruncate {
				// 智能截断：确保在markdown标题边界处截断
				truncatedContent, isTruncated = s.smartTruncateByHeading(extractedContent, 0, maxContentLength)
			} else {
				// 直接截断
				truncatedContent = extractedContent[0:maxContentLength] + "\n\n[内容已截断 - 仅显示前" + fmt.Sprintf("%d", maxContentLength) + "个字符]"
				isTruncated = true
			}
		}

		tokens := s.estimateTokens(truncatedContent)

		positionInfo := ""
		if customStartPos >= 0 {
			positionInfo = fmt.Sprintf("\n- 自定义起始位置: %d 字符\n- 自定义结束位置: %d 字符", actualStartPos, actualEndPos)
		}

		resultText := fmt.Sprintf("片段ID: %s (搜索索引)\n文档ID: %s\n版本: %s\n章节: %s\n\n元数据:\n- 原始长度: %d 字符\n- 返回长度: %d 字符%s\n- 估算Token数: %d\n- 是否截断: %v\n\n内容:\n%s",
			documentID, docIndex.DocumentID, docIndex.Version, docIndex.Section, originalLength, len(truncatedContent), positionInfo, tokens, isTruncated, truncatedContent)

		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: resultText,
				},
			},
			IsError: false,
		}, nil
	}

	// 如果不是搜索索引ID，按原有逻辑获取文档版本
	var docVersion *model.DocumentVersion
	var isVersionID bool

	// 首先尝试作为版本ID获取
	docVersion, err = s.versionRepo.GetByID(ctx, documentID)
	if err == nil && docVersion != nil {
		// 成功作为版本ID获取
		isVersionID = true
		log.Printf("DEBUG: Retrieved document by version ID: %s", documentID)
	}

	// 如果不是版本ID或失败，尝试作为文档ID获取
	if !isVersionID {
		if version != "" {
			log.Printf("DEBUG: Fetching specific version: documentID=%s, version=%s", documentID, version)
			docVersion, err = s.documentService.GetDocumentByVersion(ctx, documentID, version)
			log.Printf("DEBUG: GetDocumentByVersion result: err=%v, docVersion=%v", err, docVersion != nil)
		} else {
			// 获取最新版本
			log.Printf("DEBUG: Attempting to get latest version for documentID=%s", documentID)
			versions, listErr := s.documentService.GetDocumentVersions(ctx, documentID)
			log.Printf("DEBUG: GetDocumentVersions result: err=%v, count=%d", listErr, len(versions))

			if listErr != nil {
				log.Printf("ERROR: Failed to get document versions: %v", listErr)
				err = fmt.Errorf("failed to get document versions: %v", listErr)
				docVersion = nil
			} else if len(versions) == 0 {
				log.Printf("ERROR: No document versions found for documentID=%s", documentID)
				err = fmt.Errorf("no document versions found for documentID %s", documentID)
				docVersion = nil
			} else {
				docVersion = versions[0]
				log.Printf("DEBUG: Selected latest version: documentID=%s, version=%s", docVersion.DocumentID, docVersion.Version)
			}
		}
	}

	if err != nil {
		log.Printf("ERROR: Failed to get document version: %v", err)
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: fmt.Sprintf("获取文档版本失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// 检查 docVersion 是否为 nil
	if docVersion == nil {
		log.Printf("ERROR: docVersion is nil after successful version lookup")
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: "文档版本对象为空，无法获取内容",
				},
			},
			IsError: true,
		}, nil
	}

	log.Printf("DEBUG: Successfully retrieved docVersion: versionID=%s, documentID=%s, version=%s, contentLength=%d",
		docVersion.ID, docVersion.DocumentID, docVersion.Version, len(docVersion.Content))

	content := docVersion.Content
	originalLength := len(content)

	// 检查内容长度并记录警告
	if originalLength > WarningContentLength {
		log.Printf("WARNING: Document content length %d exceeds warning threshold %d", originalLength, WarningContentLength)
	}

	// 截断内容以控制Token使用
	truncatedContent := content
	isTruncated := false

	// 如果指定了查询关键词，先定位关键词位置
	var startFrom int = 0
	if query != "" && smartTruncate {
		// 查找关键词位置
		if pos := strings.Index(strings.ToLower(content), strings.ToLower(query)); pos != -1 {
			// 计算起始点：从关键词位置往前扩展（一般往前扩展maxContentLength/2）
			halfLength := maxContentLength / 2
			startFrom = pos - halfLength
			if startFrom < 0 {
				startFrom = 0
			}
			log.Printf("DEBUG: Found query at position %d, starting from %d", pos, startFrom)
		} else {
			log.Printf("DEBUG: Query '%s' not found in content", query)
			startFrom = 0
		}
	}

	// 根据是否启用智能截断来处理内容
	if len(content) > maxContentLength {
		if smartTruncate {
			// 智能截断：确保在markdown标题边界处截断
			truncatedContent, isTruncated = s.smartTruncateByHeading(content, startFrom, maxContentLength)
		} else {
			// 直接截断
			truncatedContent = content[startFrom:startFrom+maxContentLength] + "\n\n[内容已截断 - 仅显示前" + fmt.Sprintf("%d", maxContentLength) + "个字符]"
			isTruncated = true
		}
	} else if startFrom > 0 {
		// 如果内容未超过限制但有起始位置偏移
		endPos := startFrom + maxContentLength
		if endPos > len(content) {
			endPos = len(content)
		}
		truncatedContent, isTruncated = s.smartTruncateByHeading(content, startFrom, endPos-startFrom)
	}

	// 估算Token数量并记录
	charCount := len(truncatedContent)
	estimatedTokens := charCount / 4
	spaceCount := 0
	for _, char := range truncatedContent {
		if char == ' ' || char == '\n' || char == '\t' {
			spaceCount++
		}
	}
	estimatedTokens += spaceCount / 10
	if estimatedTokens < 1 {
		estimatedTokens = 1
	}

	log.Printf("INFO: Document content processed - originalLength: %d, returnedLength: %d, estimatedTokens: %d, truncated: %v",
		originalLength, len(truncatedContent), estimatedTokens, isTruncated)

	// 构造结果文本
	displayVersion := version
	if displayVersion == "" {
		displayVersion = docVersion.Version
	}

	var resultText string
	if isVersionID {
		resultText = fmt.Sprintf("版本ID: %s\n文档ID: %s\n版本号: %s\n\n元数据:\n- 原始长度: %d 字符\n- 返回长度: %d 字符\n- 估算Token数: %d\n- 是否截断: %v\n\n内容:\n%s",
			docVersion.ID, docVersion.DocumentID, docVersion.Version, originalLength, len(truncatedContent), estimatedTokens, isTruncated, truncatedContent)
	} else {
		resultText = fmt.Sprintf("文档ID: %s\n版本: %s\n\n元数据:\n- 原始长度: %d 字符\n- 返回长度: %d 字符\n- 估算Token数: %d\n- 是否截断: %v\n\n内容:\n%s",
			documentID, displayVersion, originalLength, len(truncatedContent), estimatedTokens, isTruncated, truncatedContent)
	}

	return &model.MCPToolResult{
		Content: []interface{}{
			model.MCPTextContent{
				Type: "text",
				Text: resultText,
			},
		},
		IsError: false,
	}, nil
}

// createSuccessResponse 创建成功响应
func (s *mcpService) createSuccessResponse(id interface{}, method string, result interface{}) (*model.MCPResponse, error) {
	return &model.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Result:  result,
	}, nil
}

// createErrorResponse 创建错误响应
func (s *mcpService) createErrorResponse(id interface{}, method string, code int, message string, data interface{}) (*model.MCPResponse, error) {
	return &model.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Error: &model.MCPError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}, nil
}

// truncateText 截断文本
func (s *mcpService) truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}

// estimateTokens 估算文本的token数量
// 这是一个简化估算方法，假设平均每个token约为4个字符（适用于中英文混合）
// 对于更精确的估算，可以使用专门的tokenizer库
func (s *mcpService) estimateTokens(text string) int {
	// 简化的token估算：字符数 / 4 + 空格和标点符号的额外计数
	charCount := len(text)
	// 基础估算
	estimatedTokens := charCount / 4
	// 考虑空格和标点符号会增加token数量
	spaceCount := 0
	for _, char := range text {
		if char == ' ' || char == '\n' || char == '\t' {
			spaceCount++
		}
	}
	estimatedTokens += spaceCount / 10
	// 确保至少返回1个token
	if estimatedTokens < 1 {
		estimatedTokens = 1
	}
	return estimatedTokens
}

// smartTruncateByHeading 按markdown标题智能截断文本
// 参数: text - 完整文本, startFrom - 起始位置, maxLen - 最大长度限制
// 返回: 截断后的文本, 是否被截断
// 说明:
//  1. 从startFrom位置开始获取最多maxLen长度的内容
//  2. 如果在长度限制处截断，但当前仍在某个markdown标题（#、##等）的内容下方，
//     则继续补全直到遇到下一个同级或更高级别的标题
//  3. 确保不会在同一个标题的内容中间截断
func (s *mcpService) smartTruncateByHeading(text string, startFrom, maxLen int) (string, bool) {
	if len(text) <= startFrom {
		return "", false
	}

	// 计算理论上的结束位置
	endPos := startFrom + maxLen
	if endPos > len(text) {
		endPos = len(text)
	}

	// 如果内容未超过最大长度，直接返回
	if endPos == len(text) {
		return text[startFrom:], false
	}

	// 需要判断是否应该继续扩展
	result := text[startFrom:endPos]
	shouldExtend := false

	// 检查截断位置是否在某个标题下方的内容中
	// 扫描从startFrom到endPos之间的内容，找到最后一个markdown标题
	lastHeadingLevel := 0
	lastHeadingPos := -1

	// 标记是否在代码块中，避免将代码注释误认为markdown标题
	inCodeBlock := false

	for i := startFrom; i < endPos; i++ {
		// 检测代码块边界（```）
		if i+2 < len(text) && text[i:i+3] == "```" {
			if inCodeBlock {
				inCodeBlock = false
			} else if !inCodeBlock {
				// 检查不是行内代码（前面不是`）
				if i == 0 || text[i-1] != '`' {
					inCodeBlock = true
				}
			}
			// 跳过 ```
			i += 2
			continue
		}

		// 检测行内代码（`...`）
		if text[i] == '`' && (i == 0 || text[i-1] != '`') {
			if i+1 < len(text) && text[i+1] == '`' && (i+2 >= len(text) || text[i+2] != '`') {
				// 这是一个单反引号，不处理
			} else {
				// 简单处理：跳过反引号
				continue
			}
		}

		// 检测markdown标题（# 开头，后面有空格），但不在代码块内
		if !inCodeBlock {
			// 必须是行首（前面是换行符或就是文档开头）
			if i == 0 || text[i-1] == '\n' {
				if text[i] == '#' {
					// 计算标题级别
					level := 0
					for j := i; j < len(text) && text[j] == '#'; j++ {
						level++
					}
					// 确保后面有空格（是真正的标题）且不超过6级
					if i+level < len(text) && (text[i+level] == ' ' || text[i+level] == '\t') && level <= 6 {
						lastHeadingLevel = level
						lastHeadingPos = i
					}
				}
			}
		}
	}

	// 如果找到了标题，检查截断位置是否在该标题内容下方
	if lastHeadingPos != -1 && endPos > lastHeadingPos {
		// 重新查找下一个同级或更高级别的标题（从endPos开始，但跳过代码块）
		nextHeadingPos := -1
		inCodeBlock := false

		for i := endPos; i < len(text); i++ {
			// 检测代码块边界（```）
			if i+2 < len(text) && text[i:i+3] == "```" {
				if inCodeBlock {
					inCodeBlock = false
				} else if !inCodeBlock {
					if i == 0 || text[i-1] != '`' {
						inCodeBlock = true
					}
				}
				i += 2
				continue
			}

			// 检测markdown标题，但不在代码块内
			if !inCodeBlock {
				if i == 0 || text[i-1] == '\n' {
					if text[i] == '#' {
						// 计算标题级别
						level := 0
						for j := i; j < len(text) && text[j] == '#'; j++ {
							level++
						}
						// 确保后面有空格且不超过6级
						if i+level < len(text) && (text[i+level] == ' ' || text[i+level] == '\t') && level <= 6 {
							// 如果是同级或更高级别的标题（level <= lastHeadingLevel）
							if level <= lastHeadingLevel {
								nextHeadingPos = i
								break
							}
						}
					}
				}
			}
		}

		// 如果找到了下一个标题，扩展到该标题之前
		if nextHeadingPos != -1 {
			shouldExtend = true
			endPos = nextHeadingPos
		} else {
			// 如果没有找到下一个标题，尝试扩展到文档末尾
			// 但限制在maxLen的1.5倍以内，避免扩展过多
			softLimit := startFrom + maxLen*3/2
			if softLimit < len(text) {
				endPos = softLimit
			} else {
				endPos = len(text)
				shouldExtend = true
			}
		}
	}

	// 获取最终结果
	result = text[startFrom:endPos]
	isTruncated := endPos < len(text)

	// 添加截断提示
	if isTruncated {
		actualLength := len(result)
		extendMsg := ""
		if shouldExtend {
			extendMsg = "（已补全到标题边界）"
		}
		result += "\n\n[内容已截断 - 已显示 " + fmt.Sprintf("%d", actualLength) +
			" 个字符" + extendMsg + "]"
	}

	log.Printf("DEBUG: smartTruncateByHeading - startFrom=%d, requestedLen=%d, actualLen=%d, truncated=%v",
		startFrom, maxLen, len(result), isTruncated)

	return result, isTruncated
}

// truncateWithWarning 带截断警告的文本处理
// 参数: text - 原始文本, maxLen - 最大长度限制
// 返回: 截断后的文本, 是否被截断, 原始长度
func (s *mcpService) truncateWithWarning(text string, maxLen int) (string, bool) {
	if len(text) <= maxLen {
		return text, false
	}
	// 截断文本并添加省略号
	truncated := text[:maxLen] + "\n\n[内容已截断 - 仅显示前" + fmt.Sprintf("%d", maxLen) + "个字符]"
	return truncated, true
}
