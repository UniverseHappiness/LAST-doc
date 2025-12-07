package router

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/UniverseHappiness/LAST-doc/internal/handler"
	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/gin-gonic/gin"
)

// MockDocumentService 是文档服务的模拟实现
type MockDocumentService struct{}

// MockSearchService 是搜索服务的模拟实现
type MockSearchService struct{}

func (m *MockDocumentService) UploadDocument(ctx context.Context, file *multipart.FileHeader, name, docType, version, library, description string, tags []string) (*model.Document, error) {
	return &model.Document{
		ID:       "test-id",
		Name:     name,
		FilePath: "/tmp/test-document.txt", // 添加文件路径以避免重定向
	}, nil
}

func (m *MockDocumentService) GetDocument(ctx context.Context, id string) (*model.Document, error) {
	return &model.Document{
		ID:       id,
		Name:     "测试文档",
		FilePath: "/tmp/test-document.txt", // 添加文件路径以避免重定向
	}, nil
}

func (m *MockDocumentService) GetDocuments(ctx context.Context, page, size int, filters map[string]interface{}) ([]*model.Document, int64, error) {
	return []*model.Document{
		{ID: "test-id-1", Name: "测试文档1"},
		{ID: "test-id-2", Name: "测试文档2"},
	}, 2, nil
}

func (m *MockDocumentService) GetDocumentVersions(ctx context.Context, documentID string) ([]*model.DocumentVersion, error) {
	return []*model.DocumentVersion{
		{ID: "version-id-1", DocumentID: documentID, Version: "1.0.0"},
		{ID: "version-id-2", DocumentID: documentID, Version: "2.0.0"},
	}, nil
}

func (m *MockDocumentService) GetDocumentByVersion(ctx context.Context, documentID, version string) (*model.DocumentVersion, error) {
	return &model.DocumentVersion{
		ID:         "version-id",
		DocumentID: documentID,
		Version:    version,
		FilePath:   "/tmp/test-document-version.txt", // 添加文件路径以避免重定向
	}, nil
}

func (m *MockDocumentService) DeleteDocument(ctx context.Context, id string) error {
	return nil
}

func (m *MockDocumentService) DeleteDocumentVersion(ctx context.Context, documentID, version string) error {
	return nil
}

func (m *MockDocumentService) GetDocumentVersionCount(ctx context.Context, documentID string) (int64, error) {
	return 2, nil
}

func (m *MockDocumentService) GetDocumentMetadata(ctx context.Context, documentID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"author": "测试作者",
		"date":   "2023-01-01",
	}, nil
}

func (m *MockDocumentService) UpdateDocument(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}

// MockSearchService 方法实现
func (m *MockSearchService) BuildIndex(ctx context.Context, documentID, version string) error {
	return nil
}

func (m *MockSearchService) BuildIndexBatch(ctx context.Context, indices []*model.SearchIndex) error {
	return nil
}

func (m *MockSearchService) Search(ctx context.Context, request *model.SearchRequest) (*model.SearchResponse, error) {
	return &model.SearchResponse{
		Total: 1,
		Items: []model.SearchResult{
			{
				ID:          "test-search-id",
				DocumentID:  "test-doc-id",
				Version:     "1.0.0",
				Title:       "测试文档",
				Content:     "测试内容",
				Snippet:     "测试片段",
				Score:       0.9,
				ContentType: "text",
				Section:     "测试章节",
			},
		},
		Page: 1,
		Size: 10,
	}, nil
}

func (m *MockSearchService) GetIndexingStatus(ctx context.Context, documentID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"document_id": documentID,
		"total":       int64(10),
		"indexed":     int64(5),
		"progress":    float64(50),
		"status":      "indexing",
	}, nil
}

func (m *MockSearchService) DeleteIndex(ctx context.Context, documentID string) error {
	return nil
}

func (m *MockSearchService) DeleteIndexByVersion(ctx context.Context, documentID, version string) error {
	return nil
}

// TestRouter_SetupRoutes 测试路由设置功能
func TestRouter_SetupRoutes(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟服务
	mockDocService := &MockDocumentService{}
	mockSearchService := &MockSearchService{}
	documentHandler := handler.NewDocumentHandler(mockDocService)
	searchHandler := handler.NewSearchHandler(mockSearchService)

	// 创建路由器
	router := NewRouter(documentHandler, searchHandler)
	r := router.SetupRoutes()

	// 测试健康检查端点
	performRequest(t, r, "GET", "/health", http.StatusOK)

	// 测试文档管理API端点
	performRequest(t, r, "GET", "/api/v1/documents", http.StatusOK)
	performRequest(t, r, "GET", "/api/v1/documents/test-id", http.StatusOK)
	performRequest(t, r, "GET", "/api/v1/documents/test-id/versions", http.StatusOK)
	performRequest(t, r, "GET", "/api/v1/documents/test-id/versions/1.0.0", http.StatusOK)
	performRequest(t, r, "GET", "/api/v1/documents/test-id/metadata", http.StatusOK)
	performRequest(t, r, "GET", "/api/v1/documents/test-id/download", http.StatusOK)
	performRequest(t, r, "GET", "/api/v1/documents/test-id/versions/1.0.0/download", http.StatusOK)
}

// performRequest 执行HTTP请求并验证响应状态码
func performRequest(t *testing.T, r *gin.Engine, method, path string, expectedStatus int) {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != expectedStatus {
		t.Errorf("请求 %s %s 预期状态码 %d, 实际 %d", method, path, expectedStatus, w.Code)
	}
}

// TestCORSMiddleware 测试CORS中间件
func TestCORSMiddleware(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟服务
	mockDocService := &MockDocumentService{}
	mockSearchService := &MockSearchService{}
	documentHandler := handler.NewDocumentHandler(mockDocService)
	searchHandler := handler.NewSearchHandler(mockSearchService)

	// 创建路由器
	router := NewRouter(documentHandler, searchHandler)
	r := router.SetupRoutes()

	// 创建OPTIONS请求
	req, _ := http.NewRequest("OPTIONS", "/api/v1/documents", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证CORS头
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("CORS头 Access-Control-Allow-Origin 设置错误")
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("CORS头 Access-Control-Allow-Methods 设置错误")
	}

	if w.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("CORS头 Access-Control-Allow-Headers 设置错误")
	}

	// 验证OPTIONS请求返回204状态码
	if w.Code != http.StatusNoContent {
		t.Errorf("OPTIONS请求预期状态码 %d, 实际 %d", http.StatusNoContent, w.Code)
	}
}

// TestAPIVersionGrouping 测试API版本分组
func TestAPIVersionGrouping(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟服务
	mockDocService := &MockDocumentService{}
	mockSearchService := &MockSearchService{}
	documentHandler := handler.NewDocumentHandler(mockDocService)
	searchHandler := handler.NewSearchHandler(mockSearchService)

	// 创建路由器
	router := NewRouter(documentHandler, searchHandler)
	r := router.SetupRoutes()

	// 测试API v1分组路由
	performRequest(t, r, "GET", "/api/v1/documents", http.StatusOK)
	performRequest(t, r, "GET", "/api/v1/documents/test-id", http.StatusOK)
	performRequest(t, r, "GET", "/api/v1/documents/test-id/versions", http.StatusOK)

	// 测试不存在的API版本
	performRequest(t, r, "GET", "/api/v2/documents", http.StatusNotFound)
	performRequest(t, r, "GET", "/api/documents", http.StatusNotFound)
}

// TestNewRouter 测试创建路由器实例
func TestNewRouter(t *testing.T) {
	// 创建模拟服务
	mockDocService := &MockDocumentService{}
	mockSearchService := &MockSearchService{}
	documentHandler := handler.NewDocumentHandler(mockDocService)
	searchHandler := handler.NewSearchHandler(mockSearchService)

	// 测试创建路由器
	router := NewRouter(documentHandler, searchHandler)

	if router == nil {
		t.Fatal("NewRouter 返回了 nil")
	}

	// 检查返回的类型是否正确
	if router == nil {
		t.Fatal("NewRouter 返回了 nil")
	}

	// 检查文档处理器是否正确设置
	if router.documentHandler != documentHandler {
		t.Error("路由器的文档处理器未正确设置")
	}

	// 检查搜索处理器是否正确设置
	if router.searchHandler != searchHandler {
		t.Error("路由器的搜索处理器未正确设置")
	}

	// 测试路由器可以正常设置路由
	r := router.SetupRoutes()
	if r == nil {
		t.Fatal("SetupRoutes 返回了 nil")
	}
}

// Test_corsMiddleware 测试CORS中间件
func Test_corsMiddleware(t *testing.T) {
	// 测试CORS中间件不为nil
	middleware := corsMiddleware()
	if middleware == nil {
		t.Fatal("corsMiddleware 返回了 nil")
	}

	// 验证返回的是函数类型
	if reflect.TypeOf(middleware).Kind() != reflect.Func {
		t.Errorf("corsMiddleware 返回的类型不是函数，实际 %T", middleware)
	}

	// 创建测试Gin引擎来测试中间件功能
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware)

	// 添加测试路由
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// 测试OPTIONS请求
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证CORS头
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("CORS头 Access-Control-Allow-Origin 设置错误")
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("CORS头 Access-Control-Allow-Methods 设置错误")
	}

	if w.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("CORS头 Access-Control-Allow-Headers 设置错误")
	}

	// 验证OPTIONS请求返回204状态码
	if w.Code != http.StatusNoContent {
		t.Errorf("OPTIONS请求预期状态码 %d, 实际 %d", http.StatusNoContent, w.Code)
	}
}
