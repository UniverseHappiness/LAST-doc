package service

import (
	"context"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/repository"
)

// MockDocumentRepository 模拟DocumentRepository
type MockDocumentRepository struct {
	mock.Mock
	CreateFunc  func(ctx context.Context, document *model.Document) error
	GetByIDFunc func(ctx context.Context, id string) (*model.Document, error)
	ListFunc    func(ctx context.Context, page, size int, filters map[string]interface{}) ([]*model.Document, int64, error)
	UpdateFunc  func(ctx context.Context, id string, updates map[string]interface{}) error
	DeleteFunc  func(ctx context.Context, id string) error
}

func (m *MockDocumentRepository) Create(ctx context.Context, document *model.Document) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, document)
	}
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetByID(ctx context.Context, id string) (*model.Document, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Document), args.Error(1)
}

func (m *MockDocumentRepository) List(ctx context.Context, page, size int, filters map[string]interface{}) ([]*model.Document, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, page, size, filters)
	}
	args := m.Called(ctx, page, size, filters)
	return args.Get(0).([]*model.Document), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, updates)
	}
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockDocumentVersionRepository 模拟DocumentVersionRepository
type MockDocumentVersionRepository struct {
	mock.Mock
	CreateFunc                    func(ctx context.Context, version *model.DocumentVersion) error
	GetByDocumentIDAndVersionFunc func(ctx context.Context, documentID, version string) (*model.DocumentVersion, error)
	GetLatestVersionFunc          func(ctx context.Context, documentID string) (*model.DocumentVersion, error)
	GetByDocumentIDFunc           func(ctx context.Context, documentID string) ([]*model.DocumentVersion, error)
}

func (m *MockDocumentVersionRepository) Create(ctx context.Context, version *model.DocumentVersion) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, version)
	}
	args := m.Called(ctx, version)
	return args.Error(0)
}

func (m *MockDocumentVersionRepository) GetByDocumentIDAndVersion(ctx context.Context, documentID, version string) (*model.DocumentVersion, error) {
	if m.GetByDocumentIDAndVersionFunc != nil {
		return m.GetByDocumentIDAndVersionFunc(ctx, documentID, version)
	}
	args := m.Called(ctx, documentID, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DocumentVersion), args.Error(1)
}

func (m *MockDocumentVersionRepository) GetLatestVersion(ctx context.Context, documentID string) (*model.DocumentVersion, error) {
	if m.GetLatestVersionFunc != nil {
		return m.GetLatestVersionFunc(ctx, documentID)
	}
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DocumentVersion), args.Error(1)
}

func (m *MockDocumentVersionRepository) GetByDocumentID(ctx context.Context, documentID string) ([]*model.DocumentVersion, error) {
	if m.GetByDocumentIDFunc != nil {
		return m.GetByDocumentIDFunc(ctx, documentID)
	}
	args := m.Called(ctx, documentID)
	return args.Get(0).([]*model.DocumentVersion), args.Error(1)
}

// MockStorageService 模拟StorageService
type MockStorageService struct {
	mock.Mock
	SaveFileFunc         func(ctx context.Context, file *multipart.FileHeader, path string) error
	GenerateFilePathFunc func(documentID, fileName string) string
}

func (m *MockStorageService) SaveFile(ctx context.Context, file *multipart.FileHeader, path string) error {
	if m.SaveFileFunc != nil {
		return m.SaveFileFunc(ctx, file, path)
	}
	args := m.Called(ctx, file, path)
	return args.Error(0)
}

func (m *MockStorageService) GenerateFilePath(documentID, fileName string) string {
	if m.GenerateFilePathFunc != nil {
		return m.GenerateFilePathFunc(documentID, fileName)
	}
	args := m.Called(documentID, fileName)
	return args.String(0)
}

// TestDocumentService_UploadDocumentSuccess 测试成功上传文档
func TestDocumentService_UploadDocumentSuccess(t *testing.T) {
	// 准备mock
	mockDocRepo := new(MockDocumentRepository)
	mockVersionRepo := new(MockDocumentVersionRepository)
	mockStorage := new(MockStorageService)

	// 设置mock预期行为
	mockDocRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*model.Document{}, int64(0), nil)

	mockDocRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil)

	mockVersionRepo.On("GetByDocumentIDAndVersion", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, repository.ErrRecordNotFound)

	mockVersionRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil)

	mockStorage.On("SaveFile", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	mockStorage.On("GenerateFilePath", mock.Anything, mock.Anything).
		Return("/test/path")

	// 创建服务
	service := NewDocumentService(
		mockDocRepo,
		mockVersionRepo,
		new(MockDocumentRepository), // Mock for metadata repo
		mockStorage,
		nil,                    // parser service
		new(MockSearchService), // Mock for search service
		"/test/storage",
	)

	// 执行测试
	ctx := context.Background()
	fileHeader := &multipart.FileHeader{
		Filename: "test.pdf",
		Size:     1024,
	}

	doc, err := service.UploadDocument(ctx, fileHeader, "Test Doc", "pdf", "default", "1.0", "test library", "Test Description", []string{"tag1", "tag2"})

	// 验证
	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, "Test Doc", doc.Name)
	assert.Equal(t, model.DocumentType("pdf"), doc.Type)
	assert.Equal(t, model.DocumentVersion("1.0"), doc.Version)

	// 验证mock调用
	mockDocRepo.AssertCalled(t, "List")
	mockDocRepo.AssertCalled(t, "Create")
	mockVersionRepo.AssertCalled(t, "Create")
	mockStorage.AssertCalled(t, "SaveFile")
}

// TestDocumentService_UploadDocument_DuplicateVersion 测试重复版本号
func TestDocumentService_UploadDocument_DuplicateVersion(t *testing.T) {
	// 准备mock
	mockDocRepo := new(MockDocumentRepository)
	mockVersionRepo := new(MockDocumentVersionRepository)
	mockStorage := new(MockStorageService)

	existingDoc := &model.Document{
		ID: "existing-doc-id",
	}

	existingVersion := &model.DocumentVersion{
		ID: "existing-version-id",
	}

	mockDocRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*model.Document{existingDoc}, int64(1), nil)

	mockVersionRepo.On("GetByDocumentIDAndVersion", mock.Anything, "existing-doc-id", "1.0").
		Return(existingVersion, nil)

	mockStorage.On("SaveFile", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	// 创建服务
	service := NewDocumentService(
		mockDocRepo,
		mockVersionRepo,
		new(MockDocumentRepository),
		mockStorage,
		nil,
		new(MockSearchService),
		"/test/storage",
	)

	// 执行测试
	ctx := context.Background()
	fileHeader := &multipart.FileHeader{
		Filename: "test.pdf",
	}

	doc, err := service.UploadDocument(ctx, fileHeader, "Test Doc", "pdf", "default", "1.0", "test library", "Test Description", []string{})

	// 验证应该返回错误（版本已存在）
	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.Contains(t, err.Error(), "版本号 1.0 已存在")
}

// TestDocumentService_UploadDocument_InvalidType 测试无效文档类型
func TestDocumentService_UploadDocument_InvalidType(t *testing.T) {
	// 准备mock
	mockDocRepo := new(MockDocumentRepository)
	mockStorage := new(MockStorageService)

	mockStorage.On("SaveFile", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	// 创建服务
	service := NewDocumentService(
		mockDocRepo,
		new(MockDocumentVersionRepository),
		new(MockDocumentRepository),
		mockStorage,
		nil,
		new(MockSearchService),
		"/test/storage",
	)

	// 执行测试 - 使用无效类型
	ctx := context.Background()
	fileHeader := &multipart.FileHeader{
		Filename: "test.txt",
	}

	doc, err := service.UploadDocument(ctx, fileHeader, "Test Doc", "invalid_type", "default", "1.0", "test library", "Test Description", []string{})

	// 验证应该返回错误
	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.Contains(t, err.Error(), "invalid document type")
}

// TestDocumentService_GetDocumentWithVersions 测试获取文档并更新为最新版本
func TestDocumentService_GetDocumentWithVersions(t *testing.T) {
	// 准备mock
	mockDocRepo := new(MockDocumentRepository)
	mockVersionRepo := new(MockDocumentVersionRepository)

	latestVersion := &model.DocumentVersion{
		ID:          "version-id",
		DocumentID:  "doc-id",
		Version:     "2.0",
		Content:     "test content",
		Status:      model.DocumentStatusReady,
		Description: "version description",
		FilePath:    "/test/path",
		FileSize:    2048,
		UpdatedAt:   mock.Anything,
	}

	document := &model.Document{
		ID:      "doc-id",
		Name:    "Test Doc",
		Version: "1.0",
	}

	mockDocRepo.On("GetByID", mock.Anything, "doc-id").
		Return(document, nil)

	mockVersionRepo.On("GetLatestVersion", mock.Anything, "doc-id").
		Return(latestVersion, nil)

	// 创建服务
	service := NewDocumentService(
		mockDocRepo,
		mockVersionRepo,
		new(MockDocumentRepository),
		new(MockStorageService),
		nil,
		new(MockSearchService),
		"/test/storage",
	)

	// 执行测试
	ctx := context.Background()
	doc, err := service.GetDocument(ctx, "doc-id")

	// 验证
	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, "2.0", doc.Version) // 应该更新为最新版本
	assert.Equal(t, "test content", doc.Content)
	assert.Equal(t, "version description", doc.Description)
	assert.Equal(t, int64(2048), doc.FileSize)

	// 验证mock调用
	mockDocRepo.AssertCalled(t, "GetByID")
	mockVersionRepo.AssertCalled(t, "GetLatestVersion")
}

// TestDocumentService_DeleteDocument 测试删除文档
func TestDocumentService_DeleteDocument(t *testing.T) {
	// 准备mock
	mockDocRepo := new(MockDocumentRepository)

	mockDocRepo.On("Delete", mock.Anything, "doc-id").
		Return(nil)

	// 创建服务
	service := NewDocumentService(
		mockDocRepo,
		new(MockDocumentVersionRepository),
		new(MockDocumentRepository),
		new(MockStorageService),
		nil,
		new(MockSearchService),
		"/test/storage",
	)

	// 执行测试
	ctx := context.Background()
	err := service.DeleteDocument(ctx, "doc-id")

	// 验证
	assert.NoError(t, err)
	mockDocRepo.AssertCalled(t, "Delete", mock.Anything, "doc-id")
}

// TestDocumentService_DeleteDocumentError 测试删除文档失败
func TestDocumentService_DeleteDocumentError(t *testing.T) {
	// 准备mock
	mockDocRepo := new(MockDocumentRepository)

	expectedError := repository.ErrRecordNotFound
	mockDocRepo.On("Delete", mock.Anything, "doc-id").
		Return(expectedError)

	// 创建服务
	service := NewDocumentService(
		mockDocRepo,
		new(MockDocumentVersionRepository),
		new(MockDocumentRepository),
		new(MockStorageService),
		nil,
		new(MockSearchService),
		"/test/storage",
	)

	// 执行测试
	ctx := context.Background()
	err := service.DeleteDocument(ctx, "doc-id")

	// 验证
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
