package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// StringArray 自定义字符串数组类型，用于处理 PostgreSQL 数组
type StringArray []string

// DocumentType 定义文档类型
type DocumentType string

const (
	DocumentTypeMarkdown DocumentType = "markdown"
	DocumentTypePDF      DocumentType = "pdf"
	DocumentTypeDocx     DocumentType = "docx"
	DocumentTypeSwagger  DocumentType = "swagger"
	DocumentTypeOpenAPI  DocumentType = "openapi"
	DocumentTypeJavaDoc  DocumentType = "java_doc"
)

// DocumentStatus 定义文档状态
type DocumentStatus string

const (
	DocumentStatusUploading  DocumentStatus = "uploading"
	DocumentStatusProcessing DocumentStatus = "processing"
	DocumentStatusCompleted  DocumentStatus = "completed"
	DocumentStatusFailed     DocumentStatus = "failed"
)

// Document 定义文档模型
type Document struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null;index"`
	Type         DocumentType   `json:"type" gorm:"not null;index"`
	Version      string         `json:"version" gorm:"not null;index"`
	Tags         StringArray    `json:"tags" gorm:"type:character varying[]"`
	FilePath     string         `json:"file_path" gorm:"not null"`
	FileSize     int64          `json:"file_size" gorm:"not null"`
	Status       DocumentStatus `json:"status" gorm:"not null;index"`
	Description  string         `json:"description"`
	Library      string         `json:"library" gorm:"index"`     // 所属库，用于版本过滤
	Content      string         `json:"content" gorm:"type:text"` // 解析后的内容摘要
	VersionCount int64          `json:"version_count" gorm:"-"`   // 版本数量，不存储到数据库
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// DocumentVersion 定义文档版本模型
type DocumentVersion struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	DocumentID  string         `json:"document_id" gorm:"not null;index"`
	Version     string         `json:"version" gorm:"not null;index"`
	FilePath    string         `json:"file_path" gorm:"not null"`
	FileSize    int64          `json:"file_size" gorm:"not null"`
	Status      DocumentStatus `json:"status" gorm:"not null"`
	Description string         `json:"description"`
	Content     string         `json:"content" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定DocumentVersion模型的表名
func (DocumentVersion) TableName() string {
	return "document_versions"
}

// DocumentMetadata 定义文档元数据模型
type DocumentMetadata struct {
	ID         string                 `json:"id" gorm:"primaryKey"`
	DocumentID string                 `json:"document_id" gorm:"not null;index"`
	Metadata   map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt  time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// Value 实现 driver.Valuer 接口，用于将 StringArray 转换为数据库可识别的格式
func (t StringArray) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}

	if len(t) == 0 {
		return "{}", nil
	}

	// 将 StringArray 转换为 PostgreSQL 数组格式的字符串
	// 格式: {"tag1","tag2","tag3"}
	var result string
	result = "{"
	for i, item := range t {
		if i > 0 {
			result += ","
		}
		// 转义字符串中的特殊字符
		escaped := strings.ReplaceAll(item, "\"", "\\\"")
		result += "\"" + escaped + "\""
	}
	result += "}"

	return result, nil
}

// Scan 实现 sql.Scanner 接口，用于从数据库读取数据到 StringArray
func (t *StringArray) Scan(value interface{}) error {
	if value == nil {
		*t = nil
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}

	// 处理 PostgreSQL 数组格式，例如: {"tag1","tag2","tag3"}
	if len(str) >= 2 && str[0] == '{' && str[len(str)-1] == '}' {
		// 移除大括号
		content := str[1 : len(str)-1]
		if content == "" {
			*t = StringArray{}
			return nil
		}

		// 分割字符串
		items := strings.Split(content, ",")
		result := make(StringArray, len(items))

		for i, item := range items {
			// 移除引号和转义字符
			if len(item) >= 2 && item[0] == '"' && item[len(item)-1] == '"' {
				unquoted := item[1 : len(item)-1]
				// 处理转义字符
				result[i] = strings.ReplaceAll(unquoted, "\\\"", "\"")
			} else {
				result[i] = item
			}
		}

		*t = result
		return nil
	}

	// 如果不是 PostgreSQL 数组格式，尝试 JSON 格式
	var result StringArray
	if err := json.Unmarshal([]byte(str), &result); err == nil {
		*t = result
		return nil
	}

	return fmt.Errorf("failed to parse tags: %s", str)
}

// TableName 指定Document模型的表名
func (Document) TableName() string {
	return "documents"
}

// TableName 指定DocumentMetadata模型的表名
func (DocumentMetadata) TableName() string {
	return "document_metadata"
}
