package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
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
}

// NewMCPService 创建MCP服务实例
func NewMCPService(db *gorm.DB, searchService SearchService, documentService DocumentService) MCPService {
	return &mcpService{
		db:              db,
		searchService:   searchService,
		documentService: documentService,
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
			Description: "获取指定文档的详细内容",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"document_id": map[string]interface{}{
						"type":        "string",
						"description": "文档ID",
					},
					"version": map[string]interface{}{
						"type":        "string",
						"description": "文档版本，如果未指定则使用最新版本",
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
		documents = append(documents, model.MCPSearchDocument{
			ID:          item.ID,
			Name:        item.Title,
			Type:        "", // 需要从文档信息中获取
			Version:     item.Version,
			Library:     item.Library, // 添加所属库
			Score:       float64(item.Score),
			Content:     item.Content,
			ContentType: item.ContentType,
			Section:     item.Section,
		})
	}

	// 构造结果文本
	resultText := fmt.Sprintf("搜索查询: %s\n找到 %d 个相关文档:\n\n", query, len(documents))
	for i, doc := range documents {
		resultText += fmt.Sprintf("%d. %s (版本: %s, 类型: %s)\n", i+1, doc.Name, doc.Version, doc.Type)
		resultText += fmt.Sprintf("   所属库: %s\n", doc.Library)
		resultText += fmt.Sprintf("   相关度: %.2f\n", doc.Score)
		resultText += fmt.Sprintf("   内容片段: %s...\n\n", s.truncateText(doc.Content, 200))
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

	// 调用文档服务获取文档列表
	filters := map[string]interface{}{}
	filters["library"] = library
	documents, total, err := s.documentService.GetDocuments(ctx, page, size, filters)
	if err != nil {
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

	// 构造结果文本
	resultText := fmt.Sprintf("库: %s\n找到 %d 个文档 (第 %d 页, 每页 %d 个):\n\n", library, total, page, size)
	for i, doc := range documents {
		resultText += fmt.Sprintf("%d. %s\n", (page-1)*size+i+1, doc.Name)
		resultText += fmt.Sprintf("   文档ID: %s\n", doc.ID)
		resultText += fmt.Sprintf("   类型: %s\n", doc.Type)
		resultText += fmt.Sprintf("   版本: %s\n", doc.Version)
		resultText += fmt.Sprintf("   描述: %s\n\n", doc.Description)
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

// getDocumentContentTool 获取文档内容工具
func (s *mcpService) getDocumentContentTool(ctx context.Context, args map[string]interface{}) (*model.MCPToolResult, error) {
	// 解析参数
	documentID, ok := args["document_id"].(string)
	if !ok || documentID == "" {
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: "文档ID不能为空",
				},
			},
			IsError: true,
		}, nil
	}

	version, _ := args["version"].(string)
	log.Printf("DEBUG: getDocumentContentTool called with documentID=%s, version='%s'", documentID, version)

	// 获取文档版本
	var docVersion *model.DocumentVersion
	var err error

	if version != "" {
		log.Printf("DEBUG: Fetching specific version: documentID=%s, version=%s", documentID, version)
		docVersion, err = s.documentService.GetDocumentByVersion(ctx, documentID, version)
		log.Printf("DEBUG: GetDocumentByVersion result: err=%v, docVersion=%v", err, docVersion != nil)
	} else {
		// 获取最新版本
		log.Printf("DEBUG: Attempting to get latest version for documentID=%s", documentID)
		// 注意：version为空字符串时，应改用GetLatestVersion，而不是GetDocumentByVersion
		versions, listErr := s.documentService.GetDocumentVersions(ctx, documentID)
		log.Printf("DEBUG: GetDocumentVersions result: err=%v, count=%d, versions=%v", listErr, len(versions), versions)

		if listErr != nil {
			log.Printf("ERROR: Failed to get document versions: %v", listErr)
			err = fmt.Errorf("failed to get document versions: %v", listErr)
			docVersion = nil
		} else if len(versions) == 0 {
			// 版本列表为空，说明文档版本已被删除或创建失败
			log.Printf("ERROR: No document versions found for documentID=%s", documentID)
			log.Printf("ERROR: This could indicate data inconsistency between search_indices and document_versions tables")
			err = fmt.Errorf("no document versions found for documentID %s (document may have been deleted or indexing failed)", documentID)
			docVersion = nil
		} else {
			docVersion = versions[0] // 获取第一个版本
			log.Printf("DEBUG: Selected first version: documentID=%s, version=%s", docVersion.DocumentID, docVersion.Version)
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

	// 关键修复：检查 docVersion 是否为 nil
	if docVersion == nil {
		log.Printf("ERROR: docVersion is nil after successful version lookup")
		return &model.MCPToolResult{
			Content: []interface{}{
				model.MCPTextContent{
					Type: "text",
					Text: fmt.Sprintf("文档版本对象为空，无法获取内容"),
				},
			},
			IsError: true,
		}, nil
	}

	log.Printf("DEBUG: Successfully retrieved docVersion: documentID=%s, version=%s, contentLength=%d",
		docVersion.DocumentID, docVersion.Version, len(docVersion.Content))

	content := docVersion.Content

	// 构造结果文本
	displayVersion := version
	if displayVersion == "" && docVersion != nil {
		displayVersion = docVersion.Version
	}
	resultText := fmt.Sprintf("文档ID: %s\n版本: %s\n\n内容:\n%s", documentID, displayVersion, content)

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
