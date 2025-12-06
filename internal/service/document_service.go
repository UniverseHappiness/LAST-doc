package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/repository"

	"github.com/google/uuid"
)

// DocumentService 文档服务接口
type DocumentService interface {
	UploadDocument(ctx context.Context, file *multipart.FileHeader, name, docType, version, library, description string, tags []string) (*model.Document, error)
	GetDocument(ctx context.Context, id string) (*model.Document, error)
	GetDocuments(ctx context.Context, page, size int, filters map[string]interface{}) ([]*model.Document, int64, error)
	GetDocumentVersions(ctx context.Context, documentID string) ([]*model.DocumentVersion, error)
	GetDocumentByVersion(ctx context.Context, documentID, version string) (*model.DocumentVersion, error)
	DeleteDocument(ctx context.Context, id string) error
	DeleteDocumentVersion(ctx context.Context, documentID, version string) error
	GetDocumentVersionCount(ctx context.Context, documentID string) (int64, error)
	GetDocumentMetadata(ctx context.Context, documentID string) (map[string]interface{}, error)
	UpdateDocument(ctx context.Context, id string, updates map[string]interface{}) error
}

// documentService 文档服务实现
type documentService struct {
	documentRepo   repository.DocumentRepository
	versionRepo    repository.DocumentVersionRepository
	metadataRepo   repository.DocumentMetadataRepository
	storageService StorageService
	parserService  DocumentParserService
	baseStorageDir string
}

// NewDocumentService 创建文档服务实例
func NewDocumentService(
	documentRepo repository.DocumentRepository,
	versionRepo repository.DocumentVersionRepository,
	metadataRepo repository.DocumentMetadataRepository,
	storageService StorageService,
	parserService DocumentParserService,
	baseStorageDir string,
) DocumentService {
	return &documentService{
		documentRepo:   documentRepo,
		versionRepo:    versionRepo,
		metadataRepo:   metadataRepo,
		storageService: storageService,
		parserService:  parserService,
		baseStorageDir: baseStorageDir,
	}
}

// UploadDocument 上传文档
func (s *documentService) UploadDocument(ctx context.Context, file *multipart.FileHeader, name, docType, version, library, description string, tags []string) (*model.Document, error) {
	// 验证文档类型
	documentType := model.DocumentType(docType)
	if !isValidDocumentType(documentType) {
		return nil, fmt.Errorf("invalid document type: %s", docType)
	}

	// 检查是否已有同库文档（仅通过Library判断）
	existingDocs, _, err := s.documentRepo.List(ctx, 1, 100, map[string]interface{}{
		"library": library,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check existing documents: %v", err)
	}

	// 生成文档ID
	var documentID string

	// 如果找到同库的文档，使用其ID
	if len(existingDocs) > 0 {
		documentID = existingDocs[0].ID
	} else {
		// 否则创建新文档
		documentID = uuid.New().String()
	}

	// 创建文档存储目录
	storageDir := filepath.Join(s.baseStorageDir, documentID)
	log.Printf("DEBUG: 创建文档存储目录 - 路径: %s\n", storageDir)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %v", err)
	}

	// 保存文件
	filePath := filepath.Join(storageDir, file.Filename)
	log.Printf("DEBUG: 保存文件 - 路径: %s\n", filePath)
	if err := s.saveFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}
	log.Printf("DEBUG: 文件保存成功 - 大小: %d 字节\n", file.Size)

	// 如果是新文档，创建文档记录
	var document *model.Document
	if len(existingDocs) == 0 {
		document = &model.Document{
			ID:          documentID,
			Name:        name,
			Type:        documentType,
			Version:     version,
			Tags:        tags,
			FilePath:    filePath,
			FileSize:    file.Size,
			Status:      model.DocumentStatusProcessing,
			Description: description,
			Library:     library,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// 保存文档记录
		if err := s.documentRepo.Create(ctx, document); err != nil {
			// 删除已保存的文件
			os.Remove(filePath)
			return nil, fmt.Errorf("failed to create document record: %v", err)
		}
	} else {
		// 获取现有文档
		document = existingDocs[0]
	}

	// 检查版本号是否已存在
	log.Printf("DEBUG: 检查版本号唯一性 - 文档ID: %s, 版本: %s\n", documentID, version)
	existingVersion, err := s.versionRepo.GetByDocumentIDAndVersion(ctx, documentID, version)
	if err == nil && existingVersion != nil {
		log.Printf("DEBUG: 版本号已存在 - 文档ID: %s, 版本: %s, 现有版本ID: %s\n",
			documentID, version, existingVersion.ID)
		// 删除已保存的文件
		os.Remove(filePath)
		// 如果是新文档，还需要删除文档记录
		if len(existingDocs) == 0 {
			s.documentRepo.Delete(ctx, documentID)
		}
		return nil, fmt.Errorf("版本号 %s 已存在，请使用不同的版本号", version)
	}
	log.Printf("DEBUG: 版本号唯一性检查通过 - 文档ID: %s, 版本: %s\n", documentID, version)

	// 创建文档版本记录
	documentVersion := &model.DocumentVersion{
		ID:          uuid.New().String(),
		DocumentID:  document.ID, // 使用文档的ID
		Version:     version,
		FilePath:    filePath,
		FileSize:    file.Size,
		Status:      model.DocumentStatusProcessing,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	log.Printf("DEBUG: 准备创建文档版本记录 - 文档ID: %s, 版本: %s, 版本记录ID: %s\n",
		documentID, version, documentVersion.ID)
	if err := s.versionRepo.Create(ctx, documentVersion); err != nil {
		log.Printf("DEBUG: 创建文档版本记录失败 - 文档ID: %s, 版本: %s, 错误: %v\n",
			documentID, version, err)
		// 删除已保存的文件
		os.Remove(filePath)
		// 如果是新文档，还需要删除文档记录
		if len(existingDocs) == 0 {
			s.documentRepo.Delete(ctx, documentID)
		}
		return nil, fmt.Errorf("failed to create document version record: %v", err)
	}
	log.Printf("DEBUG: 文档版本记录创建成功 - 文档ID: %s, 版本: %s, 版本记录ID: %s\n",
		documentID, version, documentVersion.ID)

	// 异步处理文档解析
	log.Printf("DEBUG: 开始异步处理文档解析 - 文档ID: %s, 版本: %s, 文件路径: %s\n", documentID, version, filePath)
	go s.processDocumentWithFile(documentID, version, filePath)

	return document, nil
}

// GetDocument 获取文档
func (s *documentService) GetDocument(ctx context.Context, id string) (*model.Document, error) {
	return s.documentRepo.GetByID(ctx, id)
}

// GetDocuments 获取文档列表
func (s *documentService) GetDocuments(ctx context.Context, page, size int, filters map[string]interface{}) ([]*model.Document, int64, error) {
	return s.documentRepo.List(ctx, page, size, filters)
}

// GetDocumentVersions 获取文档版本列表
func (s *documentService) GetDocumentVersions(ctx context.Context, documentID string) ([]*model.DocumentVersion, error) {
	return s.versionRepo.GetByDocumentID(ctx, documentID)
}

// GetDocumentByVersion 根据版本获取文档
func (s *documentService) GetDocumentByVersion(ctx context.Context, documentID, version string) (*model.DocumentVersion, error) {
	return s.versionRepo.GetByDocumentIDAndVersion(ctx, documentID, version)
}

// DeleteDocument 删除文档
func (s *documentService) DeleteDocument(ctx context.Context, id string) error {
	// 获取文档信息
	document, err := s.documentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 删除文件
	if err := os.RemoveAll(filepath.Dir(document.FilePath)); err != nil {
		return fmt.Errorf("failed to delete document files: %v", err)
	}

	// 删除文档记录
	if err := s.documentRepo.Delete(ctx, id); err != nil {
		return err
	}

	// 删除文档版本记录
	if err := s.versionRepo.DeleteByDocumentID(ctx, id); err != nil {
		return err
	}

	// 删除文档元数据
	if err := s.metadataRepo.DeleteByDocumentID(ctx, id); err != nil {
		return err
	}

	return nil
}

// UpdateDocument 更新文档
func (s *documentService) UpdateDocument(ctx context.Context, id string, updates map[string]interface{}) error {
	return s.documentRepo.Update(ctx, id, updates)
}

// DeleteDocumentVersion 删除文档版本
func (s *documentService) DeleteDocumentVersion(ctx context.Context, documentID, version string) error {
	// 获取文档版本信息
	docVersion, err := s.versionRepo.GetByDocumentIDAndVersion(ctx, documentID, version)
	if err != nil {
		return err
	}

	// 删除版本文件
	if err := os.Remove(docVersion.FilePath); err != nil {
		return fmt.Errorf("failed to delete version file: %v", err)
	}

	// 删除文档版本记录
	if err := s.versionRepo.Delete(ctx, docVersion.ID); err != nil {
		return err
	}

	return nil
}

// GetDocumentMetadata 获取文档元数据
func (s *documentService) GetDocumentMetadata(ctx context.Context, documentID string) (map[string]interface{}, error) {
	metadataList, err := s.metadataRepo.GetByDocumentID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	if len(metadataList) == 0 {
		return nil, nil
	}

	// 返回最新的元数据
	return metadataList[0].Metadata, nil
}

// GetDocumentVersionCount 获取文档的版本数量
func (s *documentService) GetDocumentVersionCount(ctx context.Context, documentID string) (int64, error) {
	return s.versionRepo.Count(ctx, documentID)
}

// saveFile 保存文件
func (s *documentService) saveFile(file *multipart.FileHeader, filePath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

// processDocument 处理文档（解析、提取元数据等）
func (s *documentService) processDocument(documentID string) {
	log.Printf("DEBUG: 进入processDocument函数 - 文档ID: %s\n", documentID)
	ctx := context.Background()

	// 获取文档信息
	log.Printf("DEBUG: 获取文档信息 - 文档ID: %s\n", documentID)
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		log.Printf("DEBUG: 获取文档信息失败 - 文档ID: %s, 错误: %v\n", documentID, err)
		return
	}
	log.Printf("DEBUG: 成功获取文档信息 - 文档ID: %s, 名称: %s, 版本: %s, 状态: %s\n", documentID, document.Name, document.Version, document.Status)

	// 解析文档
	log.Printf("DEBUG: 开始解析文档 - 文档ID: %s, 文件路径: %s, 类型: %s\n", documentID, document.FilePath, document.Type)
	content, metadata, err := s.parserService.ParseDocument(ctx, document.FilePath, document.Type)
	if err != nil {
		log.Printf("DEBUG: 解析文档失败 - 文档ID: %s, 错误: %v\n", documentID, err)
		// 更新文档状态为失败
		s.documentRepo.Update(ctx, documentID, map[string]interface{}{
			"status": model.DocumentStatusFailed,
		})
		s.versionRepo.UpdateStatus(ctx, documentID, document.Version, model.DocumentStatusFailed)
		return
	}
	log.Printf("DEBUG: 文档解析成功 - 文档ID: %s, 内容长度: %d\n", documentID, len(content))
	// log.Printf("DEBUG: 文档数据 content: %s\n", content)
	log.Printf("DEBUG: 文档数据 metadata: %v\n", metadata)

	// 提取元数据
	log.Printf("DEBUG: 开始提取元数据 - 文档ID: %s\n", documentID)

	// 更新文档内容
	log.Printf("DEBUG: 更新文档内容和状态 - 文档ID: %s\n", documentID)
	s.documentRepo.Update(ctx, documentID, map[string]interface{}{
		"content": content,
		"status":  model.DocumentStatusCompleted,
	})

	// 更新文档版本内容
	log.Printf("DEBUG: 更新文档版本内容和状态 - 文档ID: %s, 版本: %s\n", documentID, document.Version)
	s.versionRepo.UpdateContent(ctx, documentID, document.Version, content, model.DocumentStatusCompleted)

	// 保存元数据
	if len(metadata) > 0 {
		log.Printf("DEBUG: 保存文档元数据 - 文档ID: %s\n", documentID)
		docMetadata := &model.DocumentMetadata{
			ID:         uuid.New().String(),
			DocumentID: documentID,
			Metadata:   metadata,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		s.metadataRepo.Create(ctx, docMetadata)
	}
	log.Printf("DEBUG: 文档处理完成 - 文档ID: %s\n", documentID)
}

// processDocumentWithFile 处理带指定文件路径的文档（用于版本处理）
func (s *documentService) processDocumentWithFile(documentID, version, filePath string) {
	log.Printf("DEBUG: 进入processDocumentWithFile函数 - 文档ID: %s, 版本: %s, 文件路径: %s\n", documentID, version, filePath)
	ctx := context.Background()

	// 获取文档信息
	log.Printf("DEBUG: 获取文档信息 - 文档ID: %s\n", documentID)
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		log.Printf("DEBUG: 获取文档信息失败 - 文档ID: %s, 错误: %v\n", documentID, err)
		return
	}
	log.Printf("DEBUG: 成功获取文档信息 - 文档ID: %s, 名称: %s, 版本: %s, 状态: %s\n", documentID, document.Name, document.Version, document.Status)

	// 根据文件扩展名确定文档类型，而不是使用原始文档的类型
	fileType := s.detectFileTypeFromFile(filePath)
	log.Printf("DEBUG: 检测到文件类型 - 文档ID: %s, 文件路径: %s, 检测类型: %s, 原始类型: %s\n", documentID, filePath, fileType, document.Type)

	// 解析文档
	log.Printf("DEBUG: 开始解析文档 - 文档ID: %s, 文件路径: %s, 类型: %s\n", documentID, filePath, fileType)

	// 添加文档类型特定的诊断日志
	switch document.Type {
	case model.DocumentTypePDF:
		log.Printf("DEBUG: 检测到PDF文档，将使用PDF解析器 - 文档ID: %s\n", documentID)
	case model.DocumentTypeDocx:
		log.Printf("DEBUG: 检测到DOCX文档，将使用DOCX解析器 - 文档ID: %s\n", documentID)
	default:
		log.Printf("DEBUG: 检测到其他类型文档 - 文档ID: %s, 类型: %s\n", documentID, document.Type)
	}

	content, metadata, err := s.parserService.ParseDocument(ctx, filePath, fileType)
	if err != nil {
		log.Printf("DEBUG: 解析文档失败 - 文档ID: %s, 错误: %v\n", documentID, err)
		// 更新文档状态为失败
		s.documentRepo.Update(ctx, documentID, map[string]interface{}{
			"status": model.DocumentStatusFailed,
		})
		s.versionRepo.UpdateStatus(ctx, documentID, version, model.DocumentStatusFailed)
		return
	}
	log.Printf("DEBUG: 文档解析成功 - 文档ID: %s, 内容长度: %d, 元数据键数量: %d\n", documentID, len(content), len(metadata))
	log.Printf("DEBUG: PDF解析后数据存储位置:\n")
	log.Printf("  - 原始PDF文件: %s\n", filePath)
	log.Printf("  - 文档内容(数据库): documents表.content字段, 文档ID: %s\n", documentID)
	log.Printf("  - 版本内容(数据库): document_versions表.content字段, 文档ID: %s, 版本: %s\n", documentID, version)
	log.Printf("  - 元数据(数据库): document_metadata表.metadata字段, 文档ID: %s\n", documentID)

	// 更新文档内容
	log.Printf("DEBUG: 更新文档内容和状态 - 文档ID: %s\n", documentID)
	s.documentRepo.Update(ctx, documentID, map[string]interface{}{
		"content": content,
		"status":  model.DocumentStatusCompleted,
	})

	// 更新文档版本内容
	log.Printf("DEBUG: 更新文档版本内容和状态 - 文档ID: %s, 版本: %s\n", documentID, version)
	s.versionRepo.UpdateContent(ctx, documentID, version, content, model.DocumentStatusCompleted)

	// 保存元数据
	if len(metadata) > 0 {
		log.Printf("DEBUG: 保存文档元数据 - 文档ID: %s\n", documentID)
		docMetadata := &model.DocumentMetadata{
			ID:         uuid.New().String(),
			DocumentID: documentID,
			Metadata:   metadata,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		s.metadataRepo.Create(ctx, docMetadata)
	}
	log.Printf("DEBUG: 文档处理完成 - 文档ID: %s, 版本: %s\n", documentID, version)
}

// isValidDocumentType 验证文档类型是否有效
func isValidDocumentType(docType model.DocumentType) bool {
	switch docType {
	case model.DocumentTypeMarkdown,
		model.DocumentTypePDF,
		model.DocumentTypeDocx,
		model.DocumentTypeSwagger,
		model.DocumentTypeOpenAPI,
		model.DocumentTypeJavaDoc:
		return true
	default:
		return false
	}
}

// detectFileTypeFromFile 根据文件扩展名检测文件类型
func (s *documentService) detectFileTypeFromFile(filePath string) model.DocumentType {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".md", ".markdown":
		return model.DocumentTypeMarkdown
	case ".pdf":
		return model.DocumentTypePDF
	case ".docx", ".doc":
		return model.DocumentTypeDocx
	case ".json", ".yaml", ".yml":
		// 简单判断是Swagger还是OpenAPI或普通JSON/YAML
		if data, err := os.ReadFile(filePath); err == nil {
			content := string(data)
			if strings.Contains(content, "\"swagger\"") || strings.Contains(content, "\"openapi\"") {
				return model.DocumentTypeSwagger
			}
			return model.DocumentTypeOpenAPI
		}
		return model.DocumentTypeOpenAPI
	case ".html", ".htm":
		return model.DocumentTypeJavaDoc
	default:
		return model.DocumentTypeMarkdown // 默认返回markdown类型
	}
}

// StorageService 存储服务接口
type StorageService interface {
	SaveFile(ctx context.Context, file *multipart.FileHeader, path string) error
	DeleteFile(ctx context.Context, path string) error
	GetFile(ctx context.Context, path string) ([]byte, error)
}
