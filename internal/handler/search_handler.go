package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/service"

	"github.com/gin-gonic/gin"
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	searchService service.SearchService
}

// NewSearchHandler 创建搜索处理器实例
func NewSearchHandler(searchService service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// Search 执行搜索
func (h *SearchHandler) Search(c *gin.Context) {
	// 解析搜索请求
	var request model.SearchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Size <= 0 || request.Size > 100 {
		request.Size = 10
	}
	if request.SearchType == "" {
		request.SearchType = "keyword"
	}

	// 执行搜索
	response, err := h.searchService.Search(context.Background(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    response,
		"message": "搜索成功",
	})
}

// SearchGet 通过GET方法执行搜索
func (h *SearchHandler) SearchGet(c *gin.Context) {
	// 解析查询参数
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "搜索关键词不能为空",
		})
		return
	}

	// 解析分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	// 解析搜索类型
	searchType := c.DefaultQuery("search_type", "keyword")

	// 解析过滤条件
	filters := make(map[string]any)
	if documentID := c.Query("document_id"); documentID != "" {
		filters["document_id"] = documentID
	}
	if version := c.Query("version"); version != "" {
		filters["version"] = version
	}
	if contentType := c.Query("content_type"); contentType != "" {
		filters["content_type"] = contentType
	}
	if section := c.Query("section"); section != "" {
		filters["section"] = section
	}

	// 构建搜索请求
	request := &model.SearchRequest{
		Query:      query,
		Filters:    filters,
		Page:       page,
		Size:       size,
		SearchType: searchType,
	}

	// 执行搜索
	response, err := h.searchService.Search(context.Background(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    response,
		"message": "搜索成功",
	})
}

// BuildIndex 构建文档索引
func (h *SearchHandler) BuildIndex(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	version := c.Param("version")
	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "版本号不能为空",
		})
		return
	}

	// 构建索引
	if err := h.searchService.BuildIndex(context.Background(), documentID, version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "构建索引失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "索引构建成功",
	})
}

// GetIndexingStatus 获取索引状态
func (h *SearchHandler) GetIndexingStatus(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	// 获取索引状态
	status, err := h.searchService.GetIndexingStatus(context.Background(), documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取索引状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    status,
		"message": "获取成功",
	})
}

// DeleteIndex 删除索引
func (h *SearchHandler) DeleteIndex(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	// 删除索引
	if err := h.searchService.DeleteIndex(context.Background(), documentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除索引失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// DeleteIndexByVersion 删除指定版本的索引
func (h *SearchHandler) DeleteIndexByVersion(c *gin.Context) {
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

	// 删除索引
	if err := h.searchService.DeleteIndexByVersion(context.Background(), documentID, version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除索引失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
