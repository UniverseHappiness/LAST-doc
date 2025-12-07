package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/service"

	"github.com/gin-gonic/gin"
)

// DocumentHandler 文档处理器
type DocumentHandler struct {
	documentService service.DocumentService
}

// NewDocumentHandler 创建文档处理器实例
func NewDocumentHandler(documentService service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

// UploadDocument 上传文档
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	// 获取表单数据
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档名称不能为空",
		})
		return
	}

	docType := c.PostForm("type")
	log.Printf("DEBUG: 后端接收到的文档类型 - docType: '%s'\n", docType)
	if docType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档类型不能为空",
		})
		return
	}

	version := c.PostForm("version")
	if version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档版本不能为空",
		})
		return
	}

	library := c.PostForm("library")
	if library == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "所属库不能为空",
		})
		return
	}

	description := c.PostForm("description")
	tagsStr := c.PostForm("tags")
	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		// 去除前后空格
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	// 添加调试日志
	log.Printf("DEBUG: 前端传递的标签数据 - tagsStr: '%s', 处理后的 tags: %+v (类型: %T)\n", tagsStr, tags, tags)

	// 获取上传的文件
	log.Printf("DEBUG: 后端开始获取上传文件")
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("DEBUG: 后端获取上传文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取上传文件失败: " + err.Error(),
		})
		return
	}
	log.Printf("DEBUG: 后端成功获取上传文件 - 文件名: %s, 文件大小: %d, MIME类型: %s", file.Filename, file.Size, file.Header.Get("Content-Type"))

	// 验证文件类型
	if !isValidFileType(file.Filename, model.DocumentType(docType)) {
		log.Printf("DEBUG: 文件类型不匹配 - 文件名: %s, 选择类型: %s\n", file.Filename, docType)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件类型与文档类型不匹配",
		})
		return
	}

	// 调用服务层上传文档
	document, err := h.documentService.UploadDocument(context.Background(), file, name, docType, version, library, description, tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传文档失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":     document.ID,
			"name":   document.Name,
			"status": document.Status,
		},
		"message": "上传成功",
	})
}

// GetDocument 获取文档
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	document, err := h.documentService.GetDocument(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    document,
		"message": "获取成功",
	})
}

// GetDocuments 获取文档列表
func (h *DocumentHandler) GetDocuments(c *gin.Context) {
	// 解析分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	// 解析过滤条件
	filters := make(map[string]interface{})

	if library := c.Query("library"); library != "" {
		filters["library"] = library
	}

	if docType := c.Query("type"); docType != "" {
		filters["type"] = model.DocumentType(docType)
	}

	if version := c.Query("version"); version != "" {
		filters["version"] = version
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = model.DocumentStatus(status)
	}

	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}

	if tags := c.Query("tags"); tags != "" {
		filters["tags"] = strings.Split(tags, ",")
	}

	// 调用服务层获取文档列表
	documents, total, err := h.documentService.GetDocuments(context.Background(), page, size, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档列表失败: " + err.Error(),
		})
		return
	}

	// 为每个文档添加版本数量
	for _, doc := range documents {
		versionCount, err := h.documentService.GetDocumentVersionCount(context.Background(), doc.ID)
		if err == nil {
			doc.VersionCount = versionCount
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total": total,
			"items": documents,
			"page":  page,
			"size":  size,
		},
		"message": "获取成功",
	})
}

// GetDocumentVersions 获取文档版本列表
func (h *DocumentHandler) GetDocumentVersions(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	versions, err := h.documentService.GetDocumentVersions(context.Background(), documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档版本失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    versions,
		"message": "获取成功",
	})
}

// GetDocumentByVersion 根据版本获取文档
func (h *DocumentHandler) GetDocumentByVersion(c *gin.Context) {
	documentID := c.Param("id")
	version := c.Param("version")

	// 添加调试日志
	log.Printf("DEBUG: 获取文档版本 - 文档ID: %s, 版本: %s\n", documentID, version)

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

	documentVersion, err := h.documentService.GetDocumentByVersion(context.Background(), documentID, version)
	if err != nil {
		log.Printf("DEBUG: 获取文档版本失败 - 文档ID: %s, 版本: %s, 错误: %v\n", documentID, version, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档版本失败: " + err.Error(),
		})
		return
	}

	log.Printf("DEBUG: 成功获取文档版本 - 文档ID: %s, 版本: %s, 内容长度: %d, 状态: %s\n",
		documentID, version, len(documentVersion.Content), documentVersion.Status)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    documentVersion,
		"message": "获取成功",
	})
}

// DeleteDocument 删除文档
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	if err := h.documentService.DeleteDocument(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除文档失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// UpdateDocument 更新文档
func (h *DocumentHandler) UpdateDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	// 解析更新数据
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求数据格式错误: " + err.Error(),
		})
		return
	}

	// 不允许更新ID和创建时间
	delete(updates, "id")
	delete(updates, "created_at")

	if err := h.documentService.UpdateDocument(context.Background(), id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新文档失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// DeleteDocumentVersion 删除文档版本
func (h *DocumentHandler) DeleteDocumentVersion(c *gin.Context) {
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

	if err := h.documentService.DeleteDocumentVersion(context.Background(), documentID, version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除文档版本失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// DownloadDocument 下载文档
func (h *DocumentHandler) DownloadDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	// 获取文档信息
	document, err := h.documentService.GetDocument(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档失败: " + err.Error(),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+document.Name)
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件
	c.File(document.FilePath)
}

// DownloadDocumentVersion 下载文档版本
func (h *DocumentHandler) DownloadDocumentVersion(c *gin.Context) {
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

	// 获取文档版本信息
	documentVersion, err := h.documentService.GetDocumentByVersion(context.Background(), documentID, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档版本失败: " + err.Error(),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=v"+version+"_"+documentID)
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件
	c.File(documentVersion.FilePath)
}

// GetDocumentMetadata 获取文档元数据
func (h *DocumentHandler) GetDocumentMetadata(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文档ID不能为空",
		})
		return
	}

	metadata, err := h.documentService.GetDocumentMetadata(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档元数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    metadata,
		"message": "获取成功",
	})
}

// BuildAllMissingIndexes 为所有缺少索引的文档构建搜索索引
func (h *DocumentHandler) BuildAllMissingIndexes(c *gin.Context) {
	log.Printf("DEBUG: 开始为所有缺少索引的文档构建搜索索引")

	if err := h.documentService.BuildAllMissingIndexes(context.Background()); err != nil {
		log.Printf("DEBUG: 构建搜索索引失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "构建搜索索引失败: " + err.Error(),
		})
		return
	}

	log.Printf("DEBUG: 成功为所有缺少索引的文档构建搜索索引")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "搜索索引构建成功",
	})
}

// isValidFileType 验证文件类型是否与文档类型匹配
func isValidFileType(filename string, docType model.DocumentType) bool {
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex == -1 {
		return false // 没有扩展名
	}
	ext := strings.ToLower(filename[dotIndex:])

	switch docType {
	case model.DocumentTypeMarkdown:
		return ext == ".md" || ext == ".markdown"
	case model.DocumentTypePDF:
		return ext == ".pdf"
	case model.DocumentTypeDocx:
		return ext == ".docx" || ext == ".doc"
	case model.DocumentTypeSwagger, model.DocumentTypeOpenAPI:
		return ext == ".json" || ext == ".yaml" || ext == ".yml"
	case model.DocumentTypeJavaDoc:
		return ext == ".html" || ext == ".htm"
	default:
		return false
	}
}
