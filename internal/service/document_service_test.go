package service

import (
	"context"
	"mime/multipart"
	"reflect"
	"testing"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/repository"
)

// TestDocumentService 测试文档服务的简单测试
func TestDocumentService(t *testing.T) {
	// 简单的测试，验证服务是否能正常创建
	// 这里只做基本测试，更复杂的测试可以在后续添加

	// 验证文档类型检查函数
	testCases := []struct {
		docType string
		expect  bool
	}{
		{"markdown", true},
		{"pdf", true},
		{"docx", true},
		{"swagger", true},
		{"openapi", true},
		{"java_doc", true},
		{"unknown", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := isValidDocumentType(model.DocumentType(tc.docType))
		if result != tc.expect {
			t.Errorf("文档类型 %s 预期 %v, 实际 %v", tc.docType, tc.expect, result)
		}
	}
}

// TestStorageService 测试存储服务的简单测试
func TestStorageService(t *testing.T) {
	// 简单的测试，验证存储服务是否能正常创建
	// 由于存储服务需要文件系统操作，这里只做基本测试

	// 创建临时目录
	tempDir := "/tmp/test_storage"
	storageService := NewLocalStorageService(tempDir)

	// 验证基础目录是否设置正确
	if storageService.GetBaseDir() != tempDir {
		t.Errorf("存储目录设置错误，预期 %s, 实际 %s", tempDir, storageService.GetBaseDir())
	}
}

// TestParserService 测试解析服务的简单测试
func TestParserService(t *testing.T) {
	// 简单的测试，验证解析服务是否能正常创建
	parserService := NewParserService()

	// 验证解析服务是否创建成功
	if parserService == nil {
		t.Error("解析服务创建失败")
	}
}

// TestMarkdownParser 测试Markdown解析器的简单测试
func TestMarkdownParser(t *testing.T) {
	parser := NewMarkdownParser()

	// 验证支持的文件扩展名
	extensions := parser.SupportedExtensions()
	if len(extensions) != 2 {
		t.Errorf("Markdown解析器支持的扩展名数量错误，预期 2, 实际 %d", len(extensions))
	}

	// 验证扩展名是否正确
	validExtensions := map[string]bool{
		".md":       true,
		".markdown": true,
	}

	for _, ext := range extensions {
		if !validExtensions[ext] {
			t.Errorf("Markdown解析器不支持扩展名: %s", ext)
		}
	}
}

// TestPDFParser 测试PDF解析器的简单测试
func TestPDFParser(t *testing.T) {
	parser := NewPDFParser()

	// 验证支持的文件扩展名
	extensions := parser.SupportedExtensions()
	if len(extensions) != 1 {
		t.Errorf("PDF解析器支持的扩展名数量错误，预期 1, 实际 %d", len(extensions))
	}

	// 验证扩展名是否正确
	if extensions[0] != ".pdf" {
		t.Errorf("PDF解析器不支持扩展名: %s", extensions[0])
	}
}

// TestDocxParser 测试DOCX解析器的简单测试
func TestDocxParser(t *testing.T) {
	parser := NewDocxParser()

	// 验证支持的文件扩展名
	extensions := parser.SupportedExtensions()
	if len(extensions) != 2 {
		t.Errorf("DOCX解析器支持的扩展名数量错误，预期 2, 实际 %d", len(extensions))
	}

	// 验证扩展名是否正确
	validExtensions := map[string]bool{
		".docx": true,
		".doc":  true,
	}

	for _, ext := range extensions {
		if !validExtensions[ext] {
			t.Errorf("DOCX解析器不支持扩展名: %s", ext)
		}
	}
}

// TestSwaggerParser 测试Swagger解析器的简单测试
func TestSwaggerParser(t *testing.T) {
	parser := NewSwaggerParser()

	// 验证支持的文件扩展名
	extensions := parser.SupportedExtensions()
	if len(extensions) != 3 {
		t.Errorf("Swagger解析器支持的扩展名数量错误，预期 3, 实际 %d", len(extensions))
	}

	// 验证扩展名是否正确
	validExtensions := map[string]bool{
		".json": true,
		".yaml": true,
		".yml":  true,
	}

	for _, ext := range extensions {
		if !validExtensions[ext] {
			t.Errorf("Swagger解析器不支持扩展名: %s", ext)
		}
	}
}

// TestOpenAPIParser 测试OpenAPI解析器的简单测试
func TestOpenAPIParser(t *testing.T) {
	parser := NewOpenAPIParser()

	// 验证支持的文件扩展名
	extensions := parser.SupportedExtensions()
	if len(extensions) != 3 {
		t.Errorf("OpenAPI解析器支持的扩展名数量错误，预期 3, 实际 %d", len(extensions))
	}

	// 验证扩展名是否正确
	validExtensions := map[string]bool{
		".json": true,
		".yaml": true,
		".yml":  true,
	}

	for _, ext := range extensions {
		if !validExtensions[ext] {
			t.Errorf("OpenAPI解析器不支持扩展名: %s", ext)
		}
	}
}

// TestJavaDocParser 测试JavaDoc解析器的简单测试
func TestJavaDocParser(t *testing.T) {
	parser := NewJavaDocParser()

	// 验证支持的文件扩展名
	extensions := parser.SupportedExtensions()
	if len(extensions) != 2 {
		t.Errorf("JavaDoc解析器支持的扩展名数量错误，预期 2, 实际 %d", len(extensions))
	}

	// 验证扩展名是否正确
	validExtensions := map[string]bool{
		".html": true,
		".htm":  true,
	}

	for _, ext := range extensions {
		if !validExtensions[ext] {
			t.Errorf("JavaDoc解析器不支持扩展名: %s", ext)
		}
	}
}

func TestNewDocumentService(t *testing.T) {
	type args struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	tests := []struct {
		name string
		args args
		want DocumentService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDocumentService(tt.args.documentRepo, tt.args.versionRepo, tt.args.metadataRepo, tt.args.storageService, tt.args.parserService, tt.args.baseStorageDir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDocumentService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_documentService_UploadDocument(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		ctx         context.Context
		file        *multipart.FileHeader
		name        string
		docType     string
		version     string
		library     string
		description string
		tags        []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Document
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			got, err := s.UploadDocument(tt.args.ctx, tt.args.file, tt.args.name, tt.args.docType, tt.args.version, tt.args.library, tt.args.description, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentService.UploadDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentService.UploadDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_documentService_GetDocument(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Document
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			got, err := s.GetDocument(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentService.GetDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentService.GetDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_documentService_GetDocuments(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		ctx     context.Context
		page    int
		size    int
		filters map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Document
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			got, got1, err := s.GetDocuments(tt.args.ctx, tt.args.page, tt.args.size, tt.args.filters)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentService.GetDocuments() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentService.GetDocuments() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("documentService.GetDocuments() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_documentService_GetDocumentVersions(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		ctx        context.Context
		documentID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.DocumentVersion
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			got, err := s.GetDocumentVersions(tt.args.ctx, tt.args.documentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentService.GetDocumentVersions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentService.GetDocumentVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_documentService_GetDocumentByVersion(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		ctx        context.Context
		documentID string
		version    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.DocumentVersion
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			got, err := s.GetDocumentByVersion(tt.args.ctx, tt.args.documentID, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentService.GetDocumentByVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentService.GetDocumentByVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_documentService_DeleteDocument(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			if err := s.DeleteDocument(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("documentService.DeleteDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_documentService_UpdateDocument(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		ctx     context.Context
		id      string
		updates map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			if err := s.UpdateDocument(tt.args.ctx, tt.args.id, tt.args.updates); (err != nil) != tt.wantErr {
				t.Errorf("documentService.UpdateDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_documentService_DeleteDocumentVersion(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		ctx        context.Context
		documentID string
		version    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			if err := s.DeleteDocumentVersion(tt.args.ctx, tt.args.documentID, tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("documentService.DeleteDocumentVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_documentService_GetDocumentMetadata(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		ctx        context.Context
		documentID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			got, err := s.GetDocumentMetadata(tt.args.ctx, tt.args.documentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentService.GetDocumentMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("documentService.GetDocumentMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_documentService_GetDocumentVersionCount(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		ctx        context.Context
		documentID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			got, err := s.GetDocumentVersionCount(tt.args.ctx, tt.args.documentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("documentService.GetDocumentVersionCount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("documentService.GetDocumentVersionCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_documentService_saveFile(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		file     *multipart.FileHeader
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			if err := s.saveFile(tt.args.file, tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("documentService.saveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_documentService_processDocument(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		documentID string
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
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			s.processDocument(tt.args.documentID)
		})
	}
}

func Test_documentService_processDocumentWithFile(t *testing.T) {
	type fields struct {
		documentRepo   repository.DocumentRepository
		versionRepo    repository.DocumentVersionRepository
		metadataRepo   repository.DocumentMetadataRepository
		storageService StorageService
		parserService  DocumentParserService
		baseStorageDir string
	}
	type args struct {
		documentID string
		version    string
		filePath   string
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
			s := &documentService{
				documentRepo:   tt.fields.documentRepo,
				versionRepo:    tt.fields.versionRepo,
				metadataRepo:   tt.fields.metadataRepo,
				storageService: tt.fields.storageService,
				parserService:  tt.fields.parserService,
				baseStorageDir: tt.fields.baseStorageDir,
			}
			s.processDocumentWithFile(tt.args.documentID, tt.args.version, tt.args.filePath)
		})
	}
}

func Test_isValidDocumentType(t *testing.T) {
	type args struct {
		docType model.DocumentType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidDocumentType(tt.args.docType); got != tt.want {
				t.Errorf("isValidDocumentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
