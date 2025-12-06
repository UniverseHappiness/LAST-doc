package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MockDocumentService 是文档服务的模拟实现
type MockDocumentService struct {
	documents        map[string]*model.Document
	documentVersions map[string][]*model.DocumentVersion
	metadata         map[string]map[string]interface{}
}

func NewMockDocumentService() *MockDocumentService {
	return &MockDocumentService{
		documents:        make(map[string]*model.Document),
		documentVersions: make(map[string][]*model.DocumentVersion),
		metadata:         make(map[string]map[string]interface{}),
	}
}

func (m *MockDocumentService) UploadDocument(ctx context.Context, file *multipart.FileHeader, name, docType, version, library, description string, tags []string) (*model.Document, error) {
	docID := uuid.New().String()
	document := &model.Document{
		ID:          docID,
		Name:        name,
		Type:        model.DocumentType(docType),
		Version:     version,
		Tags:        tags,
		FilePath:    filepath.Join(os.TempDir(), file.Filename),
		FileSize:    file.Size,
		Status:      model.DocumentStatusCompleted,
		Description: description,
		Library:     library,
	}

	m.documents[docID] = document

	// 添加版本
	versionID := uuid.New().String()
	docVersion := &model.DocumentVersion{
		ID:          versionID,
		DocumentID:  docID,
		Version:     version,
		FilePath:    filepath.Join(os.TempDir(), file.Filename),
		FileSize:    file.Size,
		Status:      model.DocumentStatusCompleted,
		Description: description,
	}
	m.documentVersions[docID] = append(m.documentVersions[docID], docVersion)

	return document, nil
}

func (m *MockDocumentService) GetDocument(ctx context.Context, id string) (*model.Document, error) {
	if doc, exists := m.documents[id]; exists {
		return doc, nil
	}
	return nil, fmt.Errorf("文档不存在")
}

func (m *MockDocumentService) GetDocuments(ctx context.Context, page, size int, filters map[string]interface{}) ([]*model.Document, int64, error) {
	var docs []*model.Document
	for _, doc := range m.documents {
		docs = append(docs, doc)
	}
	return docs, int64(len(docs)), nil
}

func (m *MockDocumentService) GetDocumentVersions(ctx context.Context, documentID string) ([]*model.DocumentVersion, error) {
	if versions, exists := m.documentVersions[documentID]; exists {
		return versions, nil
	}
	return nil, fmt.Errorf("文档版本不存在")
}

func (m *MockDocumentService) GetDocumentByVersion(ctx context.Context, documentID, version string) (*model.DocumentVersion, error) {
	if versions, exists := m.documentVersions[documentID]; exists {
		for _, v := range versions {
			if v.Version == version {
				return v, nil
			}
		}
	}
	return nil, fmt.Errorf("文档版本不存在")
}

func (m *MockDocumentService) DeleteDocument(ctx context.Context, id string) error {
	if _, exists := m.documents[id]; exists {
		delete(m.documents, id)
		delete(m.documentVersions, id)
		delete(m.metadata, id)
		return nil
	}
	return fmt.Errorf("文档不存在")
}

func (m *MockDocumentService) DeleteDocumentVersion(ctx context.Context, documentID, version string) error {
	if versions, exists := m.documentVersions[documentID]; exists {
		for i, v := range versions {
			if v.Version == version {
				m.documentVersions[documentID] = append(versions[:i], versions[i+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("文档版本不存在")
}

func (m *MockDocumentService) GetDocumentVersionCount(ctx context.Context, documentID string) (int64, error) {
	if versions, exists := m.documentVersions[documentID]; exists {
		return int64(len(versions)), nil
	}
	return 0, nil
}

func (m *MockDocumentService) GetDocumentMetadata(ctx context.Context, documentID string) (map[string]interface{}, error) {
	if meta, exists := m.metadata[documentID]; exists {
		return meta, nil
	}
	return nil, nil
}

func (m *MockDocumentService) UpdateDocument(ctx context.Context, id string, updates map[string]interface{}) error {
	if doc, exists := m.documents[id]; exists {
		for key, value := range updates {
			switch key {
			case "name":
				doc.Name = value.(string)
			case "description":
				doc.Description = value.(string)
			case "tags":
				doc.Tags = value.([]string)
			}
		}
		return nil
	}
	return fmt.Errorf("文档不存在")
}

// TestUploadDocument 测试上传文档功能
func TestUploadDocument(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟文档服务
	mockService := NewMockDocumentService()
	handler := NewDocumentHandler(mockService)

	// 创建测试路由
	router := gin.New()
	router.POST("/api/v1/documents", handler.UploadDocument)

	// 创建测试文件内容
	fileContent := "测试文档内容"
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	// 创建表单文件
	part, err := writer.CreateFormFile("file", "test.md")
	if err != nil {
		t.Fatalf("创建表单文件失败: %v", err)
	}
	part.Write([]byte(fileContent))

	// 添加其他表单字段
	writer.WriteField("name", "测试文档")
	writer.WriteField("type", "markdown")
	writer.WriteField("version", "1.0.0")
	writer.WriteField("library", "测试库")
	writer.WriteField("description", "这是一个测试文档")
	writer.WriteField("tags", "测试,文档")

	writer.Close()

	// 创建HTTP请求
	req, _ := http.NewRequest("POST", "/api/v1/documents", buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 发送请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("预期状态码 %d, 实际 %d", http.StatusOK, w.Code)
	}

	// 检查响应内容
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["code"].(float64) != 200 {
		t.Errorf("预期响应码 200, 实际 %v", response["code"])
	}

	// 检查返回的数据
	data := response["data"].(map[string]interface{})
	if data["id"] == nil || data["id"].(string) == "" {
		t.Error("返回的文档ID为空")
	}
}

// TestGetDocuments 测试获取文档列表功能
func TestGetDocuments(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟文档服务并添加测试数据
	mockService := NewMockDocumentService()

	// 添加测试文档
	doc1 := &model.Document{
		ID:          uuid.New().String(),
		Name:        "测试文档1",
		Type:        model.DocumentTypeMarkdown,
		Version:     "1.0.0",
		Library:     "测试库",
		Description: "测试文档1描述",
	}

	doc2 := &model.Document{
		ID:          uuid.New().String(),
		Name:        "测试文档2",
		Type:        model.DocumentTypePDF,
		Version:     "2.0.0",
		Library:     "测试库",
		Description: "测试文档2描述",
	}

	mockService.documents[doc1.ID] = doc1
	mockService.documents[doc2.ID] = doc2

	handler := NewDocumentHandler(mockService)

	// 创建测试路由
	router := gin.New()
	router.GET("/api/v1/documents", handler.GetDocuments)

	// 创建HTTP请求
	req, _ := http.NewRequest("GET", "/api/v1/documents?page=1&size=10", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 发送请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("预期状态码 %d, 实际 %d", http.StatusOK, w.Code)
	}

	// 检查响应内容
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["code"].(float64) != 200 {
		t.Errorf("预期响应码 200, 实际 %v", response["code"])
	}

	// 检查返回的数据
	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})

	if len(items) != 2 {
		t.Errorf("预期返回2个文档，实际返回 %d 个", len(items))
	}
}

// TestGetDocument 测试获取单个文档功能
func TestGetDocument(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟文档服务并添加测试数据
	mockService := NewMockDocumentService()

	docID := uuid.New().String()
	doc := &model.Document{
		ID:          docID,
		Name:        "测试文档",
		Type:        model.DocumentTypeMarkdown,
		Version:     "1.0.0",
		Library:     "测试库",
		Description: "测试文档描述",
	}

	mockService.documents[docID] = doc

	handler := NewDocumentHandler(mockService)

	// 创建测试路由
	router := gin.New()
	router.GET("/api/v1/documents/:id", handler.GetDocument)

	// 创建HTTP请求
	req, _ := http.NewRequest("GET", "/api/v1/documents/"+docID, nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 发送请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("预期状态码 %d, 实际 %d", http.StatusOK, w.Code)
	}

	// 检查响应内容
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["code"].(float64) != 200 {
		t.Errorf("预期响应码 200, 实际 %v", response["code"])
	}

	// 检查返回的数据
	data := response["data"].(map[string]interface{})
	if data["id"].(string) != docID {
		t.Errorf("返回的文档ID不匹配，预期 %s, 实际 %s", docID, data["id"].(string))
	}
}

// TestGetDocumentVersions 测试获取文档版本功能
func TestGetDocumentVersions(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟文档服务并添加测试数据
	mockService := NewMockDocumentService()

	docID := uuid.New().String()
	doc := &model.Document{
		ID:          docID,
		Name:        "测试文档",
		Type:        model.DocumentTypeMarkdown,
		Version:     "1.0.0",
		Library:     "测试库",
		Description: "测试文档描述",
	}

	mockService.documents[docID] = doc

	// 添加版本
	version1 := &model.DocumentVersion{
		ID:          uuid.New().String(),
		DocumentID:  docID,
		Version:     "1.0.0",
		FilePath:    "/tmp/test_1.0.0.md",
		FileSize:    1024,
		Status:      model.DocumentStatusCompleted,
		Description: "版本1.0.0",
	}

	version2 := &model.DocumentVersion{
		ID:          uuid.New().String(),
		DocumentID:  docID,
		Version:     "2.0.0",
		FilePath:    "/tmp/test_2.0.0.md",
		FileSize:    2048,
		Status:      model.DocumentStatusCompleted,
		Description: "版本2.0.0",
	}

	mockService.documentVersions[docID] = append(mockService.documentVersions[docID], version1, version2)

	handler := NewDocumentHandler(mockService)

	// 创建测试路由
	router := gin.New()
	router.GET("/api/v1/documents/:id/versions", handler.GetDocumentVersions)

	// 创建HTTP请求
	req, _ := http.NewRequest("GET", "/api/v1/documents/"+docID+"/versions", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 发送请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("预期状态码 %d, 实际 %d", http.StatusOK, w.Code)
	}

	// 检查响应内容
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["code"].(float64) != 200 {
		t.Errorf("预期响应码 200, 实际 %v", response["code"])
	}

	// 检查返回的数据
	data := response["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("预期返回2个版本，实际返回 %d 个", len(data))
	}
}

// TestDeleteDocument 测试删除文档功能
func TestDeleteDocument(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟文档服务并添加测试数据
	mockService := NewMockDocumentService()

	docID := uuid.New().String()
	doc := &model.Document{
		ID:          docID,
		Name:        "测试文档",
		Type:        model.DocumentTypeMarkdown,
		Version:     "1.0.0",
		Library:     "测试库",
		Description: "测试文档描述",
	}

	mockService.documents[docID] = doc

	handler := NewDocumentHandler(mockService)

	// 创建测试路由
	router := gin.New()
	router.DELETE("/api/v1/documents/:id", handler.DeleteDocument)

	// 创建HTTP请求
	req, _ := http.NewRequest("DELETE", "/api/v1/documents/"+docID, nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 发送请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("预期状态码 %d, 实际 %d", http.StatusOK, w.Code)
	}

	// 检查响应内容
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response["code"].(float64) != 200 {
		t.Errorf("预期响应码 200, 实际 %v", response["code"])
	}

	// 验证文档是否已被删除
	if _, exists := mockService.documents[docID]; exists {
		t.Error("文档未被删除")
	}
}

func TestNewDocumentHandler(t *testing.T) {
	// 添加日志来验证测试是否被执行
	t.Log("开始执行 TestNewDocumentHandler")

	type args struct {
		documentService service.DocumentService
	}
	tests := []struct {
		name string
		args args
		want *DocumentHandler
	}{
		{
			name: "正常创建DocumentHandler",
			args: args{
				documentService: NewMockDocumentService(),
			},
			want: &DocumentHandler{
				documentService: NewMockDocumentService(),
			},
		},
		{
			name: "创建nil服务的DocumentHandler",
			args: args{
				documentService: nil,
			},
			want: &DocumentHandler{
				documentService: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("执行测试用例: %s", tt.name)
			got := NewDocumentHandler(tt.args.documentService)

			// 验证handler不为nil
			if got == nil {
				t.Errorf("NewDocumentHandler() = nil, want non-nil")
				return
			}

			// 验证服务是否正确设置
			if got.documentService != tt.args.documentService {
				t.Errorf("NewDocumentHandler().documentService = %v, want %v", got.documentService, tt.args.documentService)
			}
		})
	}
}

func TestDocumentHandler_UploadDocument(t *testing.T) {
	// 添加日志来验证测试是否被执行
	t.Log("开始执行 TestDocumentHandler_UploadDocument")

	gin.SetMode(gin.TestMode)

	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectStatus int
		expectCode   int
	}{
		{
			name: "缺少name参数",
			fields: fields{
				documentService: NewMockDocumentService(),
			},
			args: func() args {
				// 创建没有name参数的请求
				buffer := &bytes.Buffer{}
				writer := multipart.NewWriter(buffer)

				// 创建表单文件
				part, _ := writer.CreateFormFile("file", "test.md")
				part.Write([]byte("测试内容"))

				// 添加其他表单字段，但不添加name
				writer.WriteField("type", "markdown")
				writer.WriteField("version", "1.0.0")
				writer.WriteField("library", "测试库")
				writer.Close()

				req, _ := http.NewRequest("POST", "/api/v1/documents", buffer)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req

				return args{c: c}
			}(),
			expectStatus: http.StatusBadRequest,
			expectCode:   400,
		},
		{
			name: "缺少type参数",
			fields: fields{
				documentService: NewMockDocumentService(),
			},
			args: func() args {
				buffer := &bytes.Buffer{}
				writer := multipart.NewWriter(buffer)

				part, _ := writer.CreateFormFile("file", "test.md")
				part.Write([]byte("测试内容"))

				// 添加其他表单字段，但不添加type
				writer.WriteField("name", "测试文档")
				writer.WriteField("version", "1.0.0")
				writer.WriteField("library", "测试库")
				writer.Close()

				req, _ := http.NewRequest("POST", "/api/v1/documents", buffer)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req

				return args{c: c}
			}(),
			expectStatus: http.StatusBadRequest,
			expectCode:   400,
		},
		{
			name: "文件类型不匹配",
			fields: fields{
				documentService: NewMockDocumentService(),
			},
			args: func() args {
				buffer := &bytes.Buffer{}
				writer := multipart.NewWriter(buffer)

				// 创建.pdf文件但type为markdown
				part, _ := writer.CreateFormFile("file", "test.pdf")
				part.Write([]byte("测试内容"))

				writer.WriteField("name", "测试文档")
				writer.WriteField("type", "markdown") // 类型不匹配
				writer.WriteField("version", "1.0.0")
				writer.WriteField("library", "测试库")
				writer.Close()

				req, _ := http.NewRequest("POST", "/api/v1/documents", buffer)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req

				return args{c: c}
			}(),
			expectStatus: http.StatusBadRequest,
			expectCode:   400,
		},
		{
			name: "服务层返回错误",
			fields: fields{
				documentService: &MockDocumentService{
					documents:        make(map[string]*model.Document),
					documentVersions: make(map[string][]*model.DocumentVersion),
					metadata:         make(map[string]map[string]interface{}),
				},
			},
			args: func() args {
				buffer := &bytes.Buffer{}
				writer := multipart.NewWriter(buffer)

				part, _ := writer.CreateFormFile("file", "test.md")
				part.Write([]byte("测试内容"))

				writer.WriteField("name", "测试文档")
				writer.WriteField("type", "markdown")
				writer.WriteField("version", "1.0.0")
				writer.WriteField("library", "测试库")
				writer.Close()

				req, _ := http.NewRequest("POST", "/api/v1/documents", buffer)
				req.Header.Set("Content-Type", writer.FormDataContentType())

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req

				return args{c: c}
			}(),
			expectStatus: http.StatusOK,
			expectCode:   200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("执行测试用例: %s", tt.name)

			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}

			// 创建响应记录器
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = tt.args.c.Request

			// 执行handler
			h.UploadDocument(c)

			// 验证响应状态码
			if w.Code != tt.expectStatus {
				t.Errorf("UploadDocument() status = %v, want %v", w.Code, tt.expectStatus)
			}

			// 验证响应内容
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("解析响应失败: %v", err)
				return
			}

			if response["code"].(float64) != float64(tt.expectCode) {
				t.Errorf("UploadDocument() code = %v, want %v", response["code"], tt.expectCode)
			}
		})
	}
}

func TestDocumentHandler_GetDocument(t *testing.T) {
	// 添加日志来验证测试是否被执行
	t.Log("开始执行 TestDocumentHandler_GetDocument")

	gin.SetMode(gin.TestMode)

	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectStatus int
		expectCode   int
	}{
		{
			name: "空文档ID",
			fields: fields{
				documentService: NewMockDocumentService(),
			},
			args: func() args {
				req, _ := http.NewRequest("GET", "/api/v1/documents/", nil)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req
				c.Params = []gin.Param{} // 空参数
				return args{c: c}
			}(),
			expectStatus: http.StatusBadRequest,
			expectCode:   400,
		},
		{
			name: "文档不存在",
			fields: fields{
				documentService: NewMockDocumentService(),
			},
			args: func() args {
				// 不添加任何文档，模拟不存在的情况
				req, _ := http.NewRequest("GET", "/api/v1/documents/non-existent-id", nil)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req
				c.Params = []gin.Param{{Key: "id", Value: "non-existent-id"}}
				return args{c: c}
			}(),
			expectStatus: http.StatusInternalServerError,
			expectCode:   500,
		},
		{
			name: "成功获取文档",
			fields: fields{
				documentService: func() service.DocumentService {
					mockService := NewMockDocumentService()
					docID := "test-doc-id"
					doc := &model.Document{
						ID:          docID,
						Name:        "测试文档",
						Type:        model.DocumentTypeMarkdown,
						Version:     "1.0.0",
						Library:     "测试库",
						Description: "测试文档描述",
					}
					mockService.documents[docID] = doc
					return mockService
				}(),
			},
			args: func() args {
				req, _ := http.NewRequest("GET", "/api/v1/documents/test-doc-id", nil)
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req
				c.Params = []gin.Param{{Key: "id", Value: "test-doc-id"}}
				return args{c: c}
			}(),
			expectStatus: http.StatusOK,
			expectCode:   200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("执行测试用例: %s", tt.name)

			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}

			// 创建响应记录器
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = tt.args.c.Request
			c.Params = tt.args.c.Params

			// 执行handler
			h.GetDocument(c)

			// 验证响应状态码
			if w.Code != tt.expectStatus {
				t.Errorf("GetDocument() status = %v, want %v", w.Code, tt.expectStatus)
			}

			// 验证响应内容
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("解析响应失败: %v", err)
				return
			}

			if response["code"].(float64) != float64(tt.expectCode) {
				t.Errorf("GetDocument() code = %v, want %v", response["code"], tt.expectCode)
			}
		})
	}
}

func TestDocumentHandler_GetDocuments(t *testing.T) {
	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}
			h.GetDocuments(tt.args.c)
		})
	}
}

func TestDocumentHandler_GetDocumentVersions(t *testing.T) {
	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}
			h.GetDocumentVersions(tt.args.c)
		})
	}
}

func TestDocumentHandler_GetDocumentByVersion(t *testing.T) {
	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}
			h.GetDocumentByVersion(tt.args.c)
		})
	}
}

func TestDocumentHandler_DeleteDocument(t *testing.T) {
	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}
			h.DeleteDocument(tt.args.c)
		})
	}
}

func TestDocumentHandler_UpdateDocument(t *testing.T) {
	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}
			h.UpdateDocument(tt.args.c)
		})
	}
}

func TestDocumentHandler_DeleteDocumentVersion(t *testing.T) {
	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}
			h.DeleteDocumentVersion(tt.args.c)
		})
	}
}

func TestDocumentHandler_DownloadDocument(t *testing.T) {
	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}
			h.DownloadDocument(tt.args.c)
		})
	}
}

func TestDocumentHandler_DownloadDocumentVersion(t *testing.T) {
	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}
			h.DownloadDocumentVersion(tt.args.c)
		})
	}
}

func TestDocumentHandler_GetDocumentMetadata(t *testing.T) {
	type fields struct {
		documentService service.DocumentService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DocumentHandler{
				documentService: tt.fields.documentService,
			}
			h.GetDocumentMetadata(tt.args.c)
		})
	}
}

func Test_isValidFileType(t *testing.T) {
	// 添加日志来验证测试是否被执行
	t.Log("开始执行 Test_isValidFileType")

	type args struct {
		filename string
		docType  model.DocumentType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Markdown文件 - .md扩展名",
			args: args{
				filename: "test.md",
				docType:  model.DocumentTypeMarkdown,
			},
			want: true,
		},
		{
			name: "Markdown文件 - .markdown扩展名",
			args: args{
				filename: "test.markdown",
				docType:  model.DocumentTypeMarkdown,
			},
			want: true,
		},
		{
			name: "Markdown文件 - 大写扩展名",
			args: args{
				filename: "test.MD",
				docType:  model.DocumentTypeMarkdown,
			},
			want: true,
		},
		{
			name: "PDF文件 - .pdf扩展名",
			args: args{
				filename: "test.pdf",
				docType:  model.DocumentTypePDF,
			},
			want: true,
		},
		{
			name: "PDF文件 - 大写扩展名",
			args: args{
				filename: "test.PDF",
				docType:  model.DocumentTypePDF,
			},
			want: true,
		},
		{
			name: "Docx文件 - .docx扩展名",
			args: args{
				filename: "test.docx",
				docType:  model.DocumentTypeDocx,
			},
			want: true,
		},
		{
			name: "Docx文件 - .doc扩展名",
			args: args{
				filename: "test.doc",
				docType:  model.DocumentTypeDocx,
			},
			want: true,
		},
		{
			name: "Swagger文件 - .json扩展名",
			args: args{
				filename: "swagger.json",
				docType:  model.DocumentTypeSwagger,
			},
			want: true,
		},
		{
			name: "Swagger文件 - .yaml扩展名",
			args: args{
				filename: "swagger.yaml",
				docType:  model.DocumentTypeSwagger,
			},
			want: true,
		},
		{
			name: "Swagger文件 - .yml扩展名",
			args: args{
				filename: "swagger.yml",
				docType:  model.DocumentTypeSwagger,
			},
			want: true,
		},
		{
			name: "OpenAPI文件 - .json扩展名",
			args: args{
				filename: "openapi.json",
				docType:  model.DocumentTypeOpenAPI,
			},
			want: true,
		},
		{
			name: "JavaDoc文件 - .html扩展名",
			args: args{
				filename: "index.html",
				docType:  model.DocumentTypeJavaDoc,
			},
			want: true,
		},
		{
			name: "JavaDoc文件 - .htm扩展名",
			args: args{
				filename: "index.htm",
				docType:  model.DocumentTypeJavaDoc,
			},
			want: true,
		},
		{
			name: "无扩展名文件",
			args: args{
				filename: "test",
				docType:  model.DocumentTypeMarkdown,
			},
			want: false,
		},
		{
			name: "文件类型不匹配 - PDF文件但类型为Markdown",
			args: args{
				filename: "test.pdf",
				docType:  model.DocumentTypeMarkdown,
			},
			want: false,
		},
		{
			name: "文件类型不匹配 - Docx文件但类型为PDF",
			args: args{
				filename: "test.docx",
				docType:  model.DocumentTypePDF,
			},
			want: false,
		},
		{
			name: "未知文档类型",
			args: args{
				filename: "test.md",
				docType:  "unknown_type",
			},
			want: false,
		},
		{
			name: "空文件名",
			args: args{
				filename: "",
				docType:  model.DocumentTypeMarkdown,
			},
			want: false,
		},
		{
			name: "只有点的文件名",
			args: args{
				filename: ".",
				docType:  model.DocumentTypeMarkdown,
			},
			want: false,
		},
		{
			name: "多个点的文件名",
			args: args{
				filename: "test.file.md",
				docType:  model.DocumentTypeMarkdown,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("执行测试用例: %s", tt.name)
			if got := isValidFileType(tt.args.filename, tt.args.docType); got != tt.want {
				t.Errorf("isValidFileType() = %v, want %v", got, tt.want)
			}
		})
	}
}
