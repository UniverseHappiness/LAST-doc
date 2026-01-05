package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Vector 自定义类型，用于处理 jsonb 字段
type Vector []float32

// Value 实现 driver.Valuer 接口，用于将 Vector 转换为数据库可接受的值
func (v Vector) Value() (driver.Value, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

// Scan 实现 sql.Scanner 接口，用于从数据库读取数据到 Vector
func (v *Vector) Scan(value interface{}) error {
	if value == nil {
		*v = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, *v)
}

// SearchIndex 定义搜索索引模型
type SearchIndex struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	DocumentID    string    `json:"document_id" gorm:"not null;index"`
	Version       string    `json:"version" gorm:"not null;index"`
	Content       string    `json:"content" gorm:"type:text;not null"`
	ContentType   string    `json:"content_type" gorm:"not null;index"` // text, code, etc.
	Section       string    `json:"section" gorm:"index"`               // 文档章节
	Keywords      string    `json:"keywords" gorm:"type:text"`          // 关键词
	Vector        string    `json:"vector" gorm:"type:jsonb"`           // 语义向量，以JSON字符串格式存储
	Embedding     []float32 `json:"embedding" gorm:"-"`                 // 真实嵌入向量，使用pgvector扩展（暂时禁用GORM自动迁移）
	Metadata      string    `json:"metadata" gorm:"type:jsonb"`         // 额外元数据
	Score         float32   `json:"score"`                              // 搜索相关度得分
	StartPosition int       `json:"start_position"`                     // 片段在原文档中的起始位置（字符数）
	EndPosition   int       `json:"end_position"`                       // 片段在原文档中的结束位置（字符数）
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// SearchRequest 定义搜索请求模型
type SearchRequest struct {
	Query      string                 `json:"query" binding:"required"`
	Filters    map[string]interface{} `json:"filters"`
	Page       int                    `json:"page" binding:"min=1"`
	Size       int                    `json:"size" binding:"min=1,max=100"`
	SearchType string                 `json:"searchType" binding:"required"` // keyword, semantic, hybrid
}

// SearchResponse 定义搜索响应模型
type SearchResponse struct {
	Total int64          `json:"total"`
	Items []SearchResult `json:"items"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// SearchResult 定义搜索结果模型
type SearchResult struct {
	ID          string                 `json:"id"`
	DocumentID  string                 `json:"document_id"`
	Version     string                 `json:"version"`
	Title       string                 `json:"title"`
	Library     string                 `json:"library"` // 所属库
	Content     string                 `json:"content"`
	Snippet     string                 `json:"snippet"`
	Score       float32                `json:"score"`
	ContentType string                 `json:"content_type"`
	Section     string                 `json:"section"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// TableName 指定SearchIndex模型的表名
func (SearchIndex) TableName() string {
	return "search_indices"
}
