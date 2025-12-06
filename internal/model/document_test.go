package model

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStringArray_Value 测试StringArray的Value方法
func TestStringArray_Value(t *testing.T) {
	tests := []struct {
		name     string
		input    StringArray
		expected driver.Value
		wantErr  bool
	}{
		{
			name:     "nil数组",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "空数组",
			input:    StringArray{},
			expected: "{}",
			wantErr:  false,
		},
		{
			name:     "单个标签",
			input:    StringArray{"tag1"},
			expected: "{\"tag1\"}",
			wantErr:  false,
		},
		{
			name:     "多个标签",
			input:    StringArray{"tag1", "tag2", "tag3"},
			expected: "{\"tag1\",\"tag2\",\"tag3\"}",
			wantErr:  false,
		},
		{
			name:     "包含特殊字符的标签",
			input:    StringArray{"tag\"with\"quotes", "tag,with,commas"},
			expected: "{\"tag\\\"with\\\"quotes\",\"tag,with,commas\"}",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("StringArray.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("StringArray.Value() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestStringArray_Scan 测试StringArray的Scan方法
func TestStringArray_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected StringArray
		wantErr  bool
	}{
		{
			name:     "nil值",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "PostgreSQL数组格式 - 空数组",
			input:    "{}",
			expected: StringArray{},
			wantErr:  false,
		},
		{
			name:     "PostgreSQL数组格式 - 单个元素",
			input:    "{\"tag1\"}",
			expected: StringArray{"tag1"},
			wantErr:  false,
		},
		{
			name:     "PostgreSQL数组格式 - 多个元素",
			input:    "{\"tag1\",\"tag2\",\"tag3\"}",
			expected: StringArray{"tag1", "tag2", "tag3"},
			wantErr:  false,
		},
		{
			name:     "PostgreSQL数组格式 - 带转义字符",
			input:    "{\"tag\\\"with\\\"quotes\"}",
			expected: StringArray{"tag\"with\"quotes"},
			wantErr:  false,
		},
		{
			name:     "JSON格式",
			input:    "[\"tag1\",\"tag2\"]",
			expected: StringArray{"tag1", "tag2"},
			wantErr:  false,
		},
		{
			name:     "字节切片输入",
			input:    []byte("{\"tag1\",\"tag2\"}"),
			expected: StringArray{"tag1", "tag2"},
			wantErr:  false,
		},
		{
			name:     "无效类型",
			input:    123,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result StringArray
			err := result.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringArray.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestDocument_TableName 测试Document模型的TableName方法
func TestDocument_TableName(t *testing.T) {
	doc := Document{}
	tableName := doc.TableName()
	assert.Equal(t, "documents", tableName)
}

// TestDocumentVersion_TableName 测试DocumentVersion模型的TableName方法
func TestDocumentVersion_TableName(t *testing.T) {
	version := DocumentVersion{}
	tableName := version.TableName()
	assert.Equal(t, "document_versions", tableName)
}

// TestDocumentMetadata_TableName 测试DocumentMetadata模型的TableName方法
func TestDocumentMetadata_TableName(t *testing.T) {
	metadata := DocumentMetadata{}
	tableName := metadata.TableName()
	assert.Equal(t, "document_metadata", tableName)
}

// TestDocument_Validation 测试Document模型的验证
func TestDocument_Validation(t *testing.T) {
	tests := []struct {
		name    string
		doc     Document
		isValid bool
	}{
		{
			name: "有效文档",
			doc: Document{
				ID:        "test-id",
				Name:      "测试文档",
				Type:      DocumentTypeMarkdown,
				Version:   "1.0.0",
				FilePath:  "/path/to/file.md",
				FileSize:  1024,
				Status:    DocumentStatusCompleted,
				Library:   "test-library",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			isValid: true,
		},
		{
			name: "空ID",
			doc: Document{
				ID:        "",
				Name:      "测试文档",
				Type:      DocumentTypeMarkdown,
				Version:   "1.0.0",
				FilePath:  "/path/to/file.md",
				FileSize:  1024,
				Status:    DocumentStatusCompleted,
				Library:   "test-library",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			isValid: false,
		},
		{
			name: "空名称",
			doc: Document{
				ID:        "test-id",
				Name:      "",
				Type:      DocumentTypeMarkdown,
				Version:   "1.0.0",
				FilePath:  "/path/to/file.md",
				FileSize:  1024,
				Status:    DocumentStatusCompleted,
				Library:   "test-library",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			isValid: false,
		},
		{
			name: "无效类型",
			doc: Document{
				ID:        "test-id",
				Name:      "测试文档",
				Type:      "invalid-type",
				Version:   "1.0.0",
				FilePath:  "/path/to/file.md",
				FileSize:  1024,
				Status:    DocumentStatusCompleted,
				Library:   "test-library",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			isValid: false,
		},
		{
			name: "空版本",
			doc: Document{
				ID:        "test-id",
				Name:      "测试文档",
				Type:      DocumentTypeMarkdown,
				Version:   "",
				FilePath:  "/path/to/file.md",
				FileSize:  1024,
				Status:    DocumentStatusCompleted,
				Library:   "test-library",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			isValid: false,
		},
		{
			name: "空文件路径",
			doc: Document{
				ID:        "test-id",
				Name:      "测试文档",
				Type:      DocumentTypeMarkdown,
				Version:   "1.0.0",
				FilePath:  "",
				FileSize:  1024,
				Status:    DocumentStatusCompleted,
				Library:   "test-library",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			isValid: false,
		},
		{
			name: "无效状态",
			doc: Document{
				ID:        "test-id",
				Name:      "测试文档",
				Type:      DocumentTypeMarkdown,
				Version:   "1.0.0",
				FilePath:  "/path/to/file.md",
				FileSize:  1024,
				Status:    "invalid-status",
				Library:   "test-library",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validateDocument(tt.doc)
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

// TestDocumentVersion_Validation 测试DocumentVersion模型的验证
func TestDocumentVersion_Validation(t *testing.T) {
	tests := []struct {
		name    string
		version DocumentVersion
		isValid bool
	}{
		{
			name: "有效版本",
			version: DocumentVersion{
				ID:         "version-id",
				DocumentID: "doc-id",
				Version:    "1.0.0",
				FilePath:   "/path/to/file.md",
				FileSize:   1024,
				Status:     DocumentStatusCompleted,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			isValid: true,
		},
		{
			name: "空ID",
			version: DocumentVersion{
				ID:         "",
				DocumentID: "doc-id",
				Version:    "1.0.0",
				FilePath:   "/path/to/file.md",
				FileSize:   1024,
				Status:     DocumentStatusCompleted,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			isValid: false,
		},
		{
			name: "空文档ID",
			version: DocumentVersion{
				ID:         "version-id",
				DocumentID: "",
				Version:    "1.0.0",
				FilePath:   "/path/to/file.md",
				FileSize:   1024,
				Status:     DocumentStatusCompleted,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			isValid: false,
		},
		{
			name: "空版本号",
			version: DocumentVersion{
				ID:         "version-id",
				DocumentID: "doc-id",
				Version:    "",
				FilePath:   "/path/to/file.md",
				FileSize:   1024,
				Status:     DocumentStatusCompleted,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validateDocumentVersion(tt.version)
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

// TestDocumentMetadata_Validation 测试DocumentMetadata模型的验证
func TestDocumentMetadata_Validation(t *testing.T) {
	tests := []struct {
		name     string
		metadata DocumentMetadata
		isValid  bool
	}{
		{
			name: "有效元数据",
			metadata: DocumentMetadata{
				ID:         "meta-id",
				DocumentID: "doc-id",
				Metadata:   map[string]interface{}{"author": "test-author", "date": "2023-01-01"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			isValid: true,
		},
		{
			name: "空ID",
			metadata: DocumentMetadata{
				ID:         "",
				DocumentID: "doc-id",
				Metadata:   map[string]interface{}{"author": "test-author"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			isValid: false,
		},
		{
			name: "空文档ID",
			metadata: DocumentMetadata{
				ID:         "meta-id",
				DocumentID: "",
				Metadata:   map[string]interface{}{"author": "test-author"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			isValid: false,
		},
		{
			name: "nil元数据",
			metadata: DocumentMetadata{
				ID:         "meta-id",
				DocumentID: "doc-id",
				Metadata:   nil,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			isValid: true, // nil metadata is allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validateDocumentMetadata(tt.metadata)
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

// TestDocument_JSONSerialization 测试Document模型的JSON序列化
func TestDocument_JSONSerialization(t *testing.T) {
	now := time.Now()
	doc := Document{
		ID:           "test-id",
		Name:         "测试文档",
		Type:         DocumentTypeMarkdown,
		Version:      "1.0.0",
		Tags:         StringArray{"tag1", "tag2"},
		FilePath:     "/path/to/file.md",
		FileSize:     1024,
		Status:       DocumentStatusCompleted,
		Description:  "测试描述",
		Library:      "test-library",
		Content:      "测试内容",
		VersionCount: 5,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(doc)
	require.NoError(t, err)

	// 反序列化
	var deserializedDoc Document
	err = json.Unmarshal(jsonData, &deserializedDoc)
	require.NoError(t, err)

	// 验证字段
	assert.Equal(t, doc.ID, deserializedDoc.ID)
	assert.Equal(t, doc.Name, deserializedDoc.Name)
	assert.Equal(t, doc.Type, deserializedDoc.Type)
	assert.Equal(t, doc.Version, deserializedDoc.Version)
	assert.Equal(t, doc.Tags, deserializedDoc.Tags)
	assert.Equal(t, doc.FilePath, deserializedDoc.FilePath)
	assert.Equal(t, doc.FileSize, deserializedDoc.FileSize)
	assert.Equal(t, doc.Status, deserializedDoc.Status)
	assert.Equal(t, doc.Description, deserializedDoc.Description)
	assert.Equal(t, doc.Library, deserializedDoc.Library)
	assert.Equal(t, doc.Content, deserializedDoc.Content)
	assert.Equal(t, doc.VersionCount, deserializedDoc.VersionCount)
	assert.Equal(t, doc.CreatedAt.Unix(), deserializedDoc.CreatedAt.Unix())
	assert.Equal(t, doc.UpdatedAt.Unix(), deserializedDoc.UpdatedAt.Unix())
}

// TestDocumentVersion_JSONSerialization 测试DocumentVersion模型的JSON序列化
func TestDocumentVersion_JSONSerialization(t *testing.T) {
	now := time.Now()
	version := DocumentVersion{
		ID:          "version-id",
		DocumentID:  "doc-id",
		Version:     "1.0.0",
		FilePath:    "/path/to/file.md",
		FileSize:    1024,
		Status:      DocumentStatusCompleted,
		Description: "版本描述",
		Content:     "版本内容",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(version)
	require.NoError(t, err)

	// 反序列化
	var deserializedVersion DocumentVersion
	err = json.Unmarshal(jsonData, &deserializedVersion)
	require.NoError(t, err)

	// 验证字段
	assert.Equal(t, version.ID, deserializedVersion.ID)
	assert.Equal(t, version.DocumentID, deserializedVersion.DocumentID)
	assert.Equal(t, version.Version, deserializedVersion.Version)
	assert.Equal(t, version.FilePath, deserializedVersion.FilePath)
	assert.Equal(t, version.FileSize, deserializedVersion.FileSize)
	assert.Equal(t, version.Status, deserializedVersion.Status)
	assert.Equal(t, version.Description, deserializedVersion.Description)
	assert.Equal(t, version.Content, deserializedVersion.Content)
	assert.Equal(t, version.CreatedAt.Unix(), deserializedVersion.CreatedAt.Unix())
	assert.Equal(t, version.UpdatedAt.Unix(), deserializedVersion.UpdatedAt.Unix())
}

// TestDocumentMetadata_JSONSerialization 测试DocumentMetadata模型的JSON序列化
func TestDocumentMetadata_JSONSerialization(t *testing.T) {
	now := time.Now()
	metadata := DocumentMetadata{
		ID:         "meta-id",
		DocumentID: "doc-id",
		Metadata:   map[string]interface{}{"author": "test-author", "date": "2023-01-01", "count": 42},
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(metadata)
	require.NoError(t, err)

	// 反序列化
	var deserializedMetadata DocumentMetadata
	err = json.Unmarshal(jsonData, &deserializedMetadata)
	require.NoError(t, err)

	// 验证字段
	assert.Equal(t, metadata.ID, deserializedMetadata.ID)
	assert.Equal(t, metadata.DocumentID, deserializedMetadata.DocumentID)
	assert.Equal(t, metadata.Metadata["author"], deserializedMetadata.Metadata["author"])
	assert.Equal(t, metadata.Metadata["date"], deserializedMetadata.Metadata["date"])
	assert.Equal(t, metadata.Metadata["count"], deserializedMetadata.Metadata["count"])
	assert.Equal(t, metadata.CreatedAt.Unix(), deserializedMetadata.CreatedAt.Unix())
	assert.Equal(t, metadata.UpdatedAt.Unix(), deserializedMetadata.UpdatedAt.Unix())
}

// validateDocument 验证Document模型的有效性
func validateDocument(doc Document) bool {
	return doc.ID != "" &&
		doc.Name != "" &&
		doc.Type != "" &&
		doc.Version != "" &&
		doc.FilePath != "" &&
		doc.FileSize > 0 &&
		doc.Status != "" &&
		doc.Library != ""
}

// validateDocumentVersion 验证DocumentVersion模型的有效性
func validateDocumentVersion(version DocumentVersion) bool {
	return version.ID != "" &&
		version.DocumentID != "" &&
		version.Version != "" &&
		version.FilePath != "" &&
		version.FileSize > 0 &&
		version.Status != ""
}

// validateDocumentMetadata 验证DocumentMetadata模型的有效性
func validateDocumentMetadata(metadata DocumentMetadata) bool {
	return metadata.ID != "" && metadata.DocumentID != ""
}
