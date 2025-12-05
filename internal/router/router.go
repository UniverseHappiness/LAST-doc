package router

import (
	"github.com/UniverseHappiness/LAST-doc/internal/handler"

	"github.com/gin-gonic/gin"
)

// Router 路由器
type Router struct {
	documentHandler *handler.DocumentHandler
}

// NewRouter 创建路由器实例
func NewRouter(documentHandler *handler.DocumentHandler) *Router {
	return &Router{
		documentHandler: documentHandler,
	}
}

// SetupRoutes 设置路由
func (r *Router) SetupRoutes() *gin.Engine {
	// 创建Gin引擎
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 使用中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API版本分组
	v1 := router.Group("/api/v1")
	{
		// 文档管理路由
		documents := v1.Group("/documents")
		{
			// 上传文档
			documents.POST("", r.documentHandler.UploadDocument)

			// 获取文档列表
			documents.GET("", r.documentHandler.GetDocuments)

			// 获取指定文档
			documents.GET("/:id", r.documentHandler.GetDocument)

			// 更新文档
			documents.PUT("/:id", r.documentHandler.UpdateDocument)

			// 删除文档
			documents.DELETE("/:id", r.documentHandler.DeleteDocument)

			// 获取文档版本列表
			documents.GET("/:id/versions", r.documentHandler.GetDocumentVersions)

			// 获取指定版本的文档
			documents.GET("/:id/versions/:version", r.documentHandler.GetDocumentByVersion)

			// 删除指定版本的文档
			documents.DELETE("/:id/versions/:version", r.documentHandler.DeleteDocumentVersion)

			// 下载文档
			documents.GET("/:id/download", r.documentHandler.DownloadDocument)

			// 获取文档元数据
			documents.GET("/:id/metadata", r.documentHandler.GetDocumentMetadata)

			// 下载文档版本
			documents.GET("/:id/versions/:version/download", r.documentHandler.DownloadDocumentVersion)
		}
	}

	return router
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
