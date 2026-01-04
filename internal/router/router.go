package router

import (
	"github.com/UniverseHappiness/LAST-doc/internal/handler"
	"github.com/UniverseHappiness/LAST-doc/internal/middleware"
	"github.com/UniverseHappiness/LAST-doc/internal/service"

	"github.com/gin-gonic/gin"
)

// Router 路由器
type Router struct {
	documentHandler   *handler.DocumentHandler
	searchHandler     *handler.SearchHandler
	aiFormatHandler   *handler.AIFormatHandler
	mcpHandler        *handler.MCPHandler
	userHandler       *handler.UserHandler
	monitorHandler    *handler.MonitorHandler
	healthHandler     *handler.HealthHandler
	backupHandler     *handler.BackupHandler
	metricsHandler    *handler.MetricsHandler
	authMiddleware    *middleware.AuthMiddleware
	loggingMiddleware *middleware.LoggingMiddleware
}

// NewRouter 创建路由器实例
func NewRouter(documentHandler *handler.DocumentHandler, searchHandler *handler.SearchHandler, aiFormatHandler *handler.AIFormatHandler, mcpHandler *handler.MCPHandler, userHandler *handler.UserHandler, monitorHandler *handler.MonitorHandler, healthHandler *handler.HealthHandler, backupHandler *handler.BackupHandler, userService service.UserService, monitorService service.MonitorService) *Router {
	return &Router{
		documentHandler:   documentHandler,
		searchHandler:     searchHandler,
		aiFormatHandler:   aiFormatHandler,
		mcpHandler:        mcpHandler,
		userHandler:       userHandler,
		monitorHandler:    monitorHandler,
		healthHandler:     healthHandler,
		backupHandler:     backupHandler,
		metricsHandler:    handler.NewMetricsHandler(),
		authMiddleware:    middleware.NewAuthMiddleware(userService),
		loggingMiddleware: middleware.NewLoggingMiddleware(monitorService),
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
	router.Use(middleware.CORS())
	router.Use(r.loggingMiddleware.LogRequest()) // 添加日志记录中间件

	// 健康检查路由（用于Kubernetes探针）
	router.GET("/health", r.healthHandler.CheckHealth)
	router.GET("/health/live", r.healthHandler.LivenessProbe)
	router.GET("/health/ready", r.healthHandler.ReadinessProbe)
	router.GET("/health/circuit-breakers", r.healthHandler.CircuitBreakers)

	// Prometheus metrics 端点
	router.GET("/metrics", r.metricsHandler.ServeMetrics)

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

			// 获取最新版本的文档
			documents.GET("/:id/versions/latest", r.documentHandler.GetLatestVersion)

			// 获取指定版本的文档
			documents.GET("/:id/versions/:version", r.documentHandler.GetDocumentByVersion)

			// 更新指定版本的文档
			documents.PUT("/:id/versions/:version", r.documentHandler.UpdateDocumentVersion)

			// 删除指定版本的文档
			documents.DELETE("/:id/versions/:version", r.documentHandler.DeleteDocumentVersion)

			// 下载文档
			documents.GET("/:id/download", r.documentHandler.DownloadDocument)

			// 获取文档元数据
			documents.GET("/:id/metadata", r.documentHandler.GetDocumentMetadata)

			// 下载文档版本
			documents.GET("/:id/versions/:version/download", r.documentHandler.DownloadDocumentVersion)

			// 为所有缺少索引的文档构建搜索索引
			documents.POST("/build-missing-indexes", r.documentHandler.BuildAllMissingIndexes)
		}

		// 搜索路由
		search := v1.Group("/search")
		{
			// 搜索文档 (POST方式)
			search.POST("", r.searchHandler.Search)

			// 搜索文档 (GET方式)
			search.GET("", r.searchHandler.SearchGet)

			// 构建文档索引
			search.POST("/documents/:id/versions/:version/index", r.searchHandler.BuildIndex)

			// 获取索引状态
			search.GET("/documents/:id/index/status", r.searchHandler.GetIndexingStatus)

			// 删除索引
			search.DELETE("/documents/:id/index", r.searchHandler.DeleteIndex)

			// 删除指定版本的索引
			search.DELETE("/documents/:id/versions/:version/index", r.searchHandler.DeleteIndexByVersion)
		}

		// AI友好格式路由
		aiFormat := v1.Group("/ai-format")
		{
			// 获取文档的结构化内容
			aiFormat.GET("/documents/:id/versions/:version/structured", r.aiFormatHandler.StructuredContent)

			// 生成LLM优化格式
			aiFormat.POST("/documents/:id/versions/:version/llm", r.aiFormatHandler.GenerateLLMFormat)

			// 生成多粒度文档表示
			aiFormat.GET("/documents/:id/versions/:version/multigranularity", r.aiFormatHandler.GenerateMultiGranularityRepresentation)

			// 注入上下文
			aiFormat.POST("/documents/:id/versions/:version/context", r.aiFormatHandler.InjectContext)

			// 获取AI友好格式列表
			aiFormat.GET("/documents/:id/versions/:version", r.aiFormatHandler.GetAIFriendlyFormats)
		}

		// MCP协议路由
		router.POST("/mcp", r.mcpHandler.HandleMCPRequest)
		router.GET("/mcp", r.mcpHandler.SendMessage) // 也支持GET方式

		// MCP配置和管理路由
		mcp := v1.Group("/mcp")
		{
			// 获取MCP配置
			mcp.GET("/config", r.mcpHandler.GetMCPConfig)

			// 测试MCP连接
			mcp.GET("/test", r.mcpHandler.TestMCPConnection)

			// API密钥管理
			keys := mcp.Group("/keys")
			keys.Use(r.authMiddleware.RequireAuth()) // 需要认证
			{
				keys.POST("", r.mcpHandler.CreateAPIKey)
				keys.GET("", r.mcpHandler.GetAPIKeys)
				keys.DELETE("/:id", r.mcpHandler.DeleteAPIKey)
			}
		}

		// 用户认证路由
		auth := v1.Group("/auth")
		{
			// 用户注册
			auth.POST("/register", r.userHandler.Register)

			// 用户登录
			auth.POST("/login", r.userHandler.Login)

			// 刷新令牌
			auth.POST("/refresh", r.userHandler.RefreshToken)

			// 请求密码重置
			auth.POST("/password/reset-request", r.userHandler.RequestPasswordReset)

			// 重置密码
			auth.POST("/password/reset", r.userHandler.ResetPassword)
		}

		// 用户管理路由
		users := v1.Group("/users")
		users.Use(r.authMiddleware.RequireAuth()) // 需要认证
		{
			// 当前用户信息
			users.GET("/profile", r.userHandler.GetProfile)
			users.PUT("/profile", r.userHandler.UpdateProfile)
			users.POST("/change-password", r.userHandler.ChangePassword)

			// 管理员功能
			admin := users.Group("")
			admin.Use(r.authMiddleware.RequireAdmin()) // 需要管理员权限
			{
				admin.GET("", r.userHandler.ListUsers)
				admin.GET("/:id", r.userHandler.GetUser)
				admin.PUT("/:id", r.userHandler.UpdateUser)
				admin.DELETE("/:id", r.userHandler.DeleteUser)
			}
		}

		// 系统监控路由（仅管理员）
		monitor := v1.Group("/monitor")
		monitor.Use(r.authMiddleware.RequireAuth())  // 需要认证
		monitor.Use(r.authMiddleware.RequireAdmin()) // 需要管理员权限
		{
			// 获取当前系统指标
			monitor.GET("/metrics/current", r.monitorHandler.GetCurrentMetrics)

			// 获取指标历史数据
			monitor.GET("/metrics/history", r.monitorHandler.GetMetricsHistory)

			// 获取指标报告
			monitor.GET("/metrics/report", r.monitorHandler.GetMetricsReport)

			// 获取系统状态
			monitor.GET("/status", r.monitorHandler.GetSystemStatus)

			// 获取日志列表
			monitor.GET("/logs", r.monitorHandler.GetLogs)

			// 获取日志统计
			monitor.GET("/logs/stats", r.monitorHandler.GetLogStats)

			// 获取性能报告
			monitor.GET("/performance", r.monitorHandler.GetPerformanceReport)

			// 清理旧数据
			monitor.POST("/cleanup", r.monitorHandler.CleanupOldData)
		}

		// 备份管理路由（仅管理员）
		backup := v1.Group("/backup")
		backup.Use(r.authMiddleware.RequireAuth())  // 需要认证
		backup.Use(r.authMiddleware.RequireAdmin()) // 需要管理员权限
		{
			// 创建备份
			backup.POST("/create", r.backupHandler.CreateBackup)

			// 获取备份列表
			backup.GET("/list", r.backupHandler.ListBackups)

			// 恢复备份
			backup.POST("/restore/:backupId", r.backupHandler.RestoreBackup)

			// 删除备份
			backup.DELETE("/:backupId", r.backupHandler.DeleteBackup)
		}
	}

	return router
}
