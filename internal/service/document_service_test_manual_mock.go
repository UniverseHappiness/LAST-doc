package service

import (
	"context"
	"io"
	"mime/multipart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/repository"
)

// MockSearchService 模拟SearchService
type MockSearchService struct {
	mock.Mock
	BuildIndexFunc           func(ctx context.Context, documentID, version string) error
	BuildIndexBatchFunc      func(ctx context.Context, indices []*model.SearchIndex) error
	SearchFunc               func(ctx context.Context, request *model.SearchRequest) (*model.SearchResponse, error)
	GetIndexingStatusFunc    func(ctx context.Context, documentID string) (map[string]interface{}, error)
	DeleteIndexFunc          func(ctx context.Context, documentID string) error
	DeleteIndexByVersionFunc func(ctx context.Context, documentID, version string) error
	ClearCacheFunc           func() error
}

func (m *MockSearchService) BuildIndex(ctx context.Context, documentID, version string) error {
	if m.BuildIndexFunc != nil {
		return m.BuildIndexFunc(ctx, documentID, version)
	}
	args := m.Called(ctx, documentID, version)
	return args.Error(0)
}

func (m *MockSearchService) BuildIndexBatch(ctx context.Context, indices []*model.SearchIndex) error {
	if m.BuildIndexBatchFunc != nil {
		return m.BuildIndexBatchFunc(ctx, indices)
	}
	args := m.Called(ctx, indices)
	return args.Error(0)
}

func (m *MockSearchService) Search(ctx context.Context, request *model.SearchRequest) (*model.SearchResponse, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, request)
	}
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SearchResponse), args.Error(1)
}

func (m *MockSearchService) GetIndexingStatus(ctx context.Context, documentID string) (map[string]interface{}, error) {
	if m.GetIndexingStatusFunc != nil {
		return m.GetIndexingStatusFunc(ctx, documentID)
	}
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockSearchService) DeleteIndex(ctx context.Context, documentID string) error {
	if m.DeleteIndexFunc != nil {
		return m.DeleteIndexFunc(ctx, documentID)
	}
	args := m.Called(ctx, documentID)
	return args.Error(0)
}

func (m *MockSearchService) DeleteIndexByVersion(ctx context.Context, documentID, version string) error {
	if m.DeleteIndexByVersionFunc != nil {
		return m.DeleteIndexByVersionFunc(ctx, documentID, version)
	}
	args := m.Called(ctx, documentID, version)
	return args.Error(0)
}

func (m *MockSearchService) ClearCache() error {
	if m.ClearCacheFunc != nil {
		return m.ClearCacheFunc()
	}
	args := m.Called()
	return args.Error(0)
}

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

func (m *MockDocumentRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentRepository) GetByLibrary(ctx context.Context, library string, page, size int) ([]*model.Document, int64, error) {
	args := m.Called(ctx, library, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Document), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentRepository) GetByLibraryAndVersion(ctx context.Context, library, version string, page, size int) ([]*model.Document, int64, error) {
	args := m.Called(ctx, library, version, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Document), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentRepository) GetByName(ctx context.Context, name string, page, size int) ([]*model.Document, int64, error) {
	args := m.Called(ctx, name, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Document), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentRepository) GetByTag(ctx context.Context, tag string, page, size int) ([]*model.Document, int64, error) {
	args := m.Called(ctx, tag, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Document), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentRepository) GetByType(ctx context.Context, docType model.DocumentType, page, size int) ([]*model.Document, int64, error) {
	args := m.Called(ctx, docType, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Document), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentRepository) GetByVersion(ctx context.Context, version string, page, size int) ([]*model.Document, int64, error) {
	args := m.Called(ctx, version, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Document), args.Get(1).(int64), args.Error(2)
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

func (m *MockDocumentVersionRepository) Count(ctx context.Context, documentID string) (int64, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentVersionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentVersionRepository) DeleteByDocumentID(ctx context.Context, documentID string) error {
	args := m.Called(ctx, documentID)
	return args.Error(0)
}

func (m *MockDocumentVersionRepository) GetByID(ctx context.Context, id string) (*model.DocumentVersion, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DocumentVersion), args.Error(1)
}

func (m *MockDocumentVersionRepository) GetVersionsByStatus(ctx context.Context, documentID string, status model.DocumentStatus) ([]*model.DocumentVersion, error) {
	args := m.Called(ctx, documentID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.DocumentVersion), args.Error(1)
}

func (m *MockDocumentVersionRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDocumentVersionRepository) UpdateByDocumentIDAndVersion(ctx context.Context, documentID, version string, updates map[string]interface{}) error {
	args := m.Called(ctx, documentID, version, updates)
	return args.Error(0)
}

func (m *MockDocumentVersionRepository) UpdateContent(ctx context.Context, documentID, version string, content string, status model.DocumentStatus) error {
	args := m.Called(ctx, documentID, version, content, status)
	return args.Error(0)
}

func (m *MockDocumentVersionRepository) UpdateStatus(ctx context.Context, documentID, version string, status model.DocumentStatus) error {
	args := m.Called(ctx, documentID, version, status)
	return args.Error(0)
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

func (m *MockStorageService) CopyFile(ctx context.Context, srcPath, dstPath string) error {
	args := m.Called(ctx, srcPath, dstPath)
	return args.Error(0)
}

func (m *MockStorageService) DeleteFile(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *MockStorageService) FileExists(ctx context.Context, path string) bool {
	args := m.Called(ctx, path)
	return args.Bool(0)
}

func (m *MockStorageService) GetFile(ctx context.Context, path string) ([]byte, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorageService) GetFileSize(ctx context.Context, path string) (int64, error) {
	args := m.Called(ctx, path)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorageService) GetFileStream(ctx context.Context, path string) (io.ReadCloser, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageService) SaveFileStream(ctx context.Context, path string, reader io.Reader) error {
	args := m.Called(ctx, path, reader)
	return args.Error(0)
}

func (m *MockStorageService) MoveFile(ctx context.Context, srcPath, dstPath string) error {
	args := m.Called(ctx, srcPath, dstPath)
	return args.Error(0)
}

func (m *MockStorageService) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockDocumentMetadataRepository 模拟DocumentMetadataRepository
type MockDocumentMetadataRepository struct {
	mock.Mock
	CreateFunc             func(ctx context.Context, metadata *model.DocumentMetadata) error
	GetByIDFunc            func(ctx context.Context, id string) (*model.DocumentMetadata, error)
	GetByDocumentIDFunc    func(ctx context.Context, documentID string) ([]*model.DocumentMetadata, error)
	GetLatestFunc          func(ctx context.Context, documentID string) (*model.DocumentMetadata, error)
	CountFunc              func(ctx context.Context, documentID string) (int64, error)
	UpdateFunc             func(ctx context.Context, id string, updates map[string]interface{}) error
	DeleteFunc             func(ctx context.Context, id string) error
	DeleteByDocumentIDFunc func(ctx context.Context, documentID string) error
}

func (m *MockDocumentMetadataRepository) Create(ctx context.Context, metadata *model.DocumentMetadata) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, metadata)
	}
	args := m.Called(ctx, metadata)
	return args.Error(0)
}

func (m *MockDocumentMetadataRepository) GetByID(ctx context.Context, id string) (*model.DocumentMetadata, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DocumentMetadata), args.Error(1)
}

func (m *MockDocumentMetadataRepository) GetByDocumentID(ctx context.Context, documentID string) ([]*model.DocumentMetadata, error) {
	if m.GetByDocumentIDFunc != nil {
		return m.GetByDocumentIDFunc(ctx, documentID)
	}
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.DocumentMetadata), args.Error(1)
}

func (m *MockDocumentMetadataRepository) GetLatest(ctx context.Context, documentID string) (*model.DocumentMetadata, error) {
	if m.GetLatestFunc != nil {
		return m.GetLatestFunc(ctx, documentID)
	}
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DocumentMetadata), args.Error(1)
}

func (m *MockDocumentMetadataRepository) GetLatestMetadata(ctx context.Context, documentID string) (*model.DocumentMetadata, error) {
	if m.GetLatestFunc != nil {
		return m.GetLatestFunc(ctx, documentID)
	}
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DocumentMetadata), args.Error(1)
}

func (m *MockDocumentMetadataRepository) Count(ctx context.Context, documentID string) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx, documentID)
	}
	args := m.Called(ctx, documentID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDocumentMetadataRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, updates)
	}
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockDocumentMetadataRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDocumentMetadataRepository) DeleteByDocumentID(ctx context.Context, documentID string) error {
	if m.DeleteByDocumentIDFunc != nil {
		return m.DeleteByDocumentIDFunc(ctx, documentID)
	}
	args := m.Called(ctx, documentID)
	return args.Error(0)
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
		new(MockDocumentMetadataRepository), // Mock for metadata repo
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
	assert.Equal(t, "1.0", doc.Version)

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
		new(MockDocumentMetadataRepository),
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
		new(MockDocumentMetadataRepository),
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
		Status:      model.DocumentStatusCompleted,
		Description: "version description",
		FilePath:    "/test/path",
		FileSize:    2048,
		UpdatedAt:   time.Now(),
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
		new(MockDocumentMetadataRepository),
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
		new(MockDocumentMetadataRepository),
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
		new(MockDocumentMetadataRepository),
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
