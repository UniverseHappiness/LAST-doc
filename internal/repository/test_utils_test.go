package repository

import (
	"testing"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
)

func TestCreateTestDocument(t *testing.T) {
	doc := CreateTestDocument()

	if doc.Name != "测试文档" {
		t.Errorf("CreateTestDocument() Name = %v, want 测试文档", doc.Name)
	}

	if doc.Type != model.DocumentTypeMarkdown {
		t.Errorf("CreateTestDocument() Type = %v, want %v", doc.Type, model.DocumentTypeMarkdown)
	}

	if doc.Version != "1.0.0" {
		t.Errorf("CreateTestDocument() Version = %v, want 1.0.0", doc.Version)
	}

	if len(doc.Tags) != 2 {
		t.Errorf("CreateTestDocument() Tags length = %d, want 2", len(doc.Tags))
	}

	if doc.FilePath != "/tmp/test.md" {
		t.Errorf("CreateTestDocument() FilePath = %v, want /tmp/test.md", doc.FilePath)
	}

	if doc.FileSize != 1024 {
		t.Errorf("CreateTestDocument() FileSize = %d, want 1024", doc.FileSize)
	}

	if doc.Status != model.DocumentStatusCompleted {
		t.Errorf("CreateTestDocument() Status = %v, want %v", doc.Status, model.DocumentStatusCompleted)
	}

	if doc.Library != "测试库" {
		t.Errorf("CreateTestDocument() Library = %v, want 测试库", doc.Library)
	}

	if doc.Description != "这是一个测试文档" {
		t.Errorf("CreateTestDocument() Description = %v, want 这是一个测试文档", doc.Description)
	}

	if doc.CreatedAt.IsZero() {
		t.Error("CreateTestDocument() CreatedAt should not be zero")
	}

	if doc.UpdatedAt.IsZero() {
		t.Error("CreateTestDocument() UpdatedAt should not be zero")
	}

	// Verify VersionCount is set
	doc.VersionCount = 0
}

func TestCreateTestDocumentMetadata(t *testing.T) {
	documentID := "test-doc-id"
	meta := CreateTestDocumentMetadata(documentID)

	if meta.DocumentID != documentID {
		t.Errorf("CreateTestDocumentMetadata() DocumentID = %v, want %v", meta.DocumentID, documentID)
	}

	if meta.Metadata == nil {
		t.Error("CreateTestDocumentMetadata() Metadata should not be nil")
	}

	author, ok := meta.Metadata["author"]
	if !ok || author != "测试作者" {
		t.Errorf("CreateTestDocumentMetadata() Metadata[author] = %v, want 测试作者", author)
	}

	date, ok := meta.Metadata["date"]
	if !ok || date != "2023-01-01" {
		t.Errorf("CreateTestDocumentMetadata() Metadata[date] = %v, want 2023-01-01", date)
	}

	_, ok = meta.Metadata["tags"]
	if !ok {
		t.Error("CreateTestDocumentMetadata() Metadata[tags] should exist")
	}

	if meta.CreatedAt.IsZero() {
		t.Error("CreateTestDocumentMetadata() CreatedAt should not be zero")
	}

	if meta.UpdatedAt.IsZero() {
		t.Error("CreateTestDocumentMetadata() UpdatedAt should not be zero")
	}
}

func TestCreateTestDocumentVersion(t *testing.T) {
	documentID := "test-doc-id"
	ver := CreateTestDocumentVersion(documentID)

	if ver.DocumentID != documentID {
		t.Errorf("CreateTestDocumentVersion() DocumentID = %v, want %v", ver.DocumentID, documentID)
	}

	if ver.Version != "1.0.0" {
		t.Errorf("CreateTestDocumentVersion() Version = %v, want 1.0.0", ver.Version)
	}

	if ver.FilePath != "/tmp/test_1.0.0.md" {
		t.Errorf("CreateTestDocumentVersion() FilePath = %v, want /tmp/test_1.0.0.md", ver.FilePath)
	}

	if ver.FileSize != 1024 {
		t.Errorf("CreateTestDocumentVersion() FileSize = %d, want 1024", ver.FileSize)
	}

	if ver.Status != model.DocumentStatusCompleted {
		t.Errorf("CreateTestDocumentVersion() Status = %v, want %v", ver.Status, model.DocumentStatusCompleted)
	}

	if ver.Description != "版本1.0.0" {
		t.Errorf("CreateTestDocumentVersion() Description = %v, want 版本1.0.0", ver.Description)
	}

	if ver.CreatedAt.IsZero() {
		t.Error("CreateTestDocumentVersion() CreatedAt should not be zero")
	}

	if ver.UpdatedAt.IsZero() {
		t.Error("CreateTestDocumentVersion() UpdatedAt should not be zero")
	}
}

func TestDocumentType_Constants(t *testing.T) {
	types := []model.DocumentType{
		model.DocumentTypeMarkdown,
		model.DocumentTypePDF,
		model.DocumentTypeDocx,
		model.DocumentTypeSwagger,
		model.DocumentTypeOpenAPI,
		model.DocumentTypeJavaDoc,
	}

	for _, typ := range types {
		if typ == "" {
			t.Errorf("DocumentType constant should not be empty")
		}
	}
}

func TestDocumentCategory_Constants(t *testing.T) {
	categories := []model.DocumentCategory{
		model.CategoryCode,
		model.CategoryDocument,
	}

	for _, cat := range categories {
		if cat == "" {
			t.Errorf("DocumentCategory constant should not be empty")
		}
	}
}

func TestDocumentStatus_Constants(t *testing.T) {
	statuses := []model.DocumentStatus{
		model.DocumentStatusUploading,
		model.DocumentStatusProcessing,
		model.DocumentStatusCompleted,
		model.DocumentStatusFailed,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("DocumentStatus constant should not be empty")
		}
	}
}

func TestDocument_TableName(t *testing.T) {
	doc := model.Document{}

	if doc.TableName() != "documents" {
		t.Errorf("Document.TableName() = %v, want documents", doc.TableName())
	}
}

func TestDocumentVersion_TableName(t *testing.T) {
	ver := model.DocumentVersion{}

	if ver.TableName() != "document_versions" {
		t.Errorf("DocumentVersion.TableName() = %v, want document_versions", ver.TableName())
	}
}

func TestDocumentMetadata_TableName(t *testing.T) {
	meta := model.DocumentMetadata{}

	if meta.TableName() != "document_metadata" {
		t.Errorf("DocumentMetadata.TableName() = %v, want document_metadata", meta.TableName())
	}
}

func TestStringArray_Value(t *testing.T) {
	arr := model.StringArray{"tag1", "tag2"}

	value, err := arr.Value()
	if err != nil {
		t.Errorf("StringArray.Value() error = %v", err)
	}

	strValue, ok := value.(string)
	if !ok {
		t.Errorf("StringArray.Value() result should be string")
	}

	if strValue != `{"tag1","tag2"}` {
		t.Errorf("StringArray.Value() = %v, want {\"tag1\",\"tag2\"}", strValue)
	}
}

func TestStringArray_Scan(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  model.StringArray
	}{
		{
			name:  "empty array",
			input: `{}`,
			want:  model.StringArray{},
		},
		{
			name:  "single element",
			input: `{"tag1"}`,
			want:  model.StringArray{"tag1"},
		},
		{
			name:  "multiple elements",
			input: `{"tag1","tag2","tag3"}`,
			want:  model.StringArray{"tag1", "tag2", "tag3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result model.StringArray
			err := result.Scan(tt.input)

			if err != nil {
				t.Errorf("StringArray.Scan() error = %v, want nil", err)
			}

			if len(result) != len(tt.want) {
				t.Errorf("StringArray.Scan() length = %d, want %d", len(result), len(tt.want))
			}

			for i, item := range result {
				if item != tt.want[i] {
					t.Errorf("StringArray.Scan() [%d] = %v, want %v", i, item, tt.want[i])
				}
			}
		})
	}
}
