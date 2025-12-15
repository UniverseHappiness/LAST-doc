package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/service"

	"github.com/gin-gonic/gin"
)

// AIFormatHandler AI格式处理器
type AIFormatHandler struct {
	aiFormatService service.AIFriendlyFormatService
	documentService service.DocumentService
}

// NewAIFormatHandler 创建AI格式处理器实例
func NewAIFormatHandler(aiFormatService service.AIFriendlyFormatService, documentService service.DocumentService) *AIFormatHandler {
	return &AIFormatHandler{
		aiFormatService: aiFormatService,
		documentService: documentService,
	}
}

// StructuredContent 结构化文档内容
func (h *AIFormatHandler) StructuredContent(c *gin.Context) {
	documentID := c.Param("id")
	version := c.Param("version")

	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "版本号不能为空",
		})
		return
	}

	// 获取文档版本内容
	documentVersion, err := h.documentService.GetDocumentByVersion(context.Background(), documentID, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档版本失败: " + err.Error(),
		})
		return
	}

	// 获取文档信息以获取类型
	document, err := h.documentService.GetDocument(context.Background(), documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档信息失败: " + err.Error(),
		})
		return
	}

	// 结构化文档内容
	structuredContent, err := h.aiFormatService.StructuredContent(context.Background(), documentID, version, documentVersion.Content, document.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "结构化文档内容失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    structuredContent,
		"message": "结构化成功",
	})
}

// GenerateLLMFormat 生成LLM优化格式
func (h *AIFormatHandler) GenerateLLMFormat(c *gin.Context) {
	documentID := c.Param("id")
	version := c.Param("version")

	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "版本号不能为空",
		})
		return
	}

	// 解析LLM格式选项
	var options model.LLMFormatOptions
	if err := c.ShouldBindJSON(&options); err != nil {
		// 如果解析失败，使用默认选项
		options = model.LLMFormatOptions{
			MaxTokens:       4000,
			PreserveCode:    true,
			SummaryLevel:    model.SummaryLevelMedium,
			IncludeMetadata: true,
		}
	}

	// 获取文档版本内容
	documentVersion, err := h.documentService.GetDocumentByVersion(context.Background(), documentID, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档版本失败: " + err.Error(),
		})
		return
	}

	// 获取文档信息以获取类型
	document, err := h.documentService.GetDocument(context.Background(), documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档信息失败: " + err.Error(),
		})
		return
	}

	// 结构化文档内容
	structuredContent, err := h.aiFormatService.StructuredContent(context.Background(), documentID, version, documentVersion.Content, document.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "结构化文档内容失败: " + err.Error(),
		})
		return
	}

	// 生成LLM优化格式
	llmContent, err := h.aiFormatService.GenerateLLMFormat(context.Background(), structuredContent, &options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成LLM优化格式失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    llmContent,
		"message": "生成成功",
	})
}

// GenerateMultiGranularityRepresentation 生成多粒度文档表示
func (h *AIFormatHandler) GenerateMultiGranularityRepresentation(c *gin.Context) {
	documentID := c.Param("id")
	version := c.Param("version")

	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "版本号不能为空",
		})
		return
	}

	// 获取文档版本内容
	documentVersion, err := h.documentService.GetDocumentByVersion(context.Background(), documentID, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档版本失败: " + err.Error(),
		})
		return
	}

	// 获取文档信息以获取类型
	document, err := h.documentService.GetDocument(context.Background(), documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档信息失败: " + err.Error(),
		})
		return
	}

	// 结构化文档内容
	structuredContent, err := h.aiFormatService.StructuredContent(context.Background(), documentID, version, documentVersion.Content, document.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "结构化文档内容失败: " + err.Error(),
		})
		return
	}

	// 生成多粒度文档表示
	representation, err := h.aiFormatService.GenerateMultiGranularityRepresentation(context.Background(), structuredContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成多粒度文档表示失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    representation,
		"message": "生成成功",
	})
}

// InjectContext 注入上下文
func (h *AIFormatHandler) InjectContext(c *gin.Context) {
	documentID := c.Param("id")
	version := c.Param("version")

	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "版本号不能为空",
		})
		return
	}

	// 解析请求参数
	var request struct {
		Query   string                        `json:"query"`
		Options model.ContextInjectionOptions `json:"options"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数格式错误: " + err.Error(),
		})
		return
	}

	if request.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "查询内容不能为空",
		})
		return
	}

	// 注入上下文
	result, err := h.aiFormatService.InjectContext(context.Background(), documentID, version, request.Query, &request.Options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "注入上下文失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    result,
		"message": "注入成功",
	})
}

// GetAIFriendlyFormats 获取AI友好格式列表
func (h *AIFormatHandler) GetAIFriendlyFormats(c *gin.Context) {
	documentID := c.Param("id")
	version := c.Param("version")

	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "版本号不能为空",
		})
		return
	}

	// 解析格式类型参数
	formatType := c.Query("type")
	maxTokens, _ := strconv.Atoi(c.DefaultQuery("max_tokens", "4000"))
	includeCode := c.DefaultQuery("include_code", "true") == "true"
	summaryLevel := model.SummaryLevel(c.DefaultQuery("summary_level", "medium"))

	// 获取文档版本内容
	documentVersion, err := h.documentService.GetDocumentByVersion(context.Background(), documentID, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档版本失败: " + err.Error(),
		})
		return
	}

	// 获取文档信息以获取类型
	document, err := h.documentService.GetDocument(context.Background(), documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档信息失败: " + err.Error(),
		})
		return
	}

	// 结构化文档内容
	structuredContent, err := h.aiFormatService.StructuredContent(context.Background(), documentID, version, documentVersion.Content, document.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "结构化文档内容失败: " + err.Error(),
		})
		return
	}

	// 根据请求的格式类型返回相应的格式
	var result interface{}

	switch formatType {
	case "structured":
		result = structuredContent

	case "llm":
		options := &model.LLMFormatOptions{
			MaxTokens:       maxTokens,
			PreserveCode:    includeCode,
			SummaryLevel:    summaryLevel,
			IncludeMetadata: true,
		}
		result, err = h.aiFormatService.GenerateLLMFormat(context.Background(), structuredContent, options)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "生成LLM优化格式失败: " + err.Error(),
			})
			return
		}

	case "multigranularity":
		result, err = h.aiFormatService.GenerateMultiGranularityRepresentation(context.Background(), structuredContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "生成多粒度文档表示失败: " + err.Error(),
			})
			return
		}

	default:
		// 默认返回所有格式
		llmOptions := &model.LLMFormatOptions{
			MaxTokens:       maxTokens,
			PreserveCode:    includeCode,
			SummaryLevel:    summaryLevel,
			IncludeMetadata: true,
		}

		llmContent, err := h.aiFormatService.GenerateLLMFormat(context.Background(), structuredContent, llmOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "生成LLM优化格式失败: " + err.Error(),
			})
			return
		}

		multiRep, err := h.aiFormatService.GenerateMultiGranularityRepresentation(context.Background(), structuredContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "生成多粒度文档表示失败: " + err.Error(),
			})
			return
		}

		result = gin.H{
			"structured":       structuredContent,
			"llm_optimized":    llmContent,
			"multigranularity": multiRep,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    result,
		"message": "获取成功",
	})
}
