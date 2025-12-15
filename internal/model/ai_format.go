package model

import (
	"time"
)

// AIStructuredContent AI结构化内容
type AIStructuredContent struct {
	ID           string                `json:"id"`
	DocumentID   string                `json:"document_id"`
	Version      string                `json:"version"`
	DocType      DocumentType          `json:"doc_type"`
	Segments     []*ContentSegment     `json:"segments"`
	CodeExamples []*CodeExample        `json:"code_examples"`
	Annotations  []*SemanticAnnotation `json:"annotations"`
	Relations    []*ContentRelation    `json:"relations"`
	CreatedAt    time.Time             `json:"created_at"`
}

// ContentSegment 内容段落
type ContentSegment struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Type     string `json:"type"`
	Position int    `json:"position"`
}

// CodeExample 代码示例
type CodeExample struct {
	ID          string `json:"id"`
	Language    string `json:"language"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Position    int    `json:"position"`
	IsInline    bool   `json:"is_inline"`
}

// SemanticAnnotation 语义标注
type SemanticAnnotation struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Context  string `json:"context"`
	Position int    `json:"position"`
}

// ContentRelation 内容关系
type ContentRelation struct {
	SourceID   string  `json:"source_id"`
	TargetID   string  `json:"target_id"`
	Type       string  `json:"type"`
	Confidence float64 `json:"confidence"`
}

// LLMOptimizedContent LLM优化内容
type LLMOptimizedContent struct {
	ID         string                 `json:"id"`
	DocumentID string                 `json:"document_id"`
	Version    string                 `json:"version"`
	Header     string                 `json:"header"`
	Content    string                 `json:"content"`
	Footer     string                 `json:"footer"`
	Metadata   map[string]interface{} `json:"metadata"`
	Options    LLMFormatOptions       `json:"options"`
	CreatedAt  time.Time              `json:"created_at"`
}

// LLMFormatOptions LLM格式选项
type LLMFormatOptions struct {
	MaxTokens       int          `json:"max_tokens"`
	PreserveCode    bool         `json:"preserve_code"`
	SummaryLevel    SummaryLevel `json:"summary_level"`
	IncludeMetadata bool         `json:"include_metadata"`
}

// SummaryLevel 摘要级别
type SummaryLevel string

const (
	SummaryLevelBrief    SummaryLevel = "brief"
	SummaryLevelMedium   SummaryLevel = "medium"
	SummaryLevelDetailed SummaryLevel = "detailed"
)

// MultiGranularityRepresentation 多粒度表示
type MultiGranularityRepresentation struct {
	ID           string                       `json:"id"`
	DocumentID   string                       `json:"document_id"`
	Version      string                       `json:"version"`
	Overview     string                       `json:"overview"`
	Sections     []*SectionRepresentation     `json:"sections"`
	Paragraphs   []*ParagraphRepresentation   `json:"paragraphs"`
	CodeSnippets []*CodeSnippetRepresentation `json:"code_snippets"`
	CreatedAt    time.Time                    `json:"created_at"`
}

// SectionRepresentation 章节表示
type SectionRepresentation struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Type                string   `json:"type"`
	Content             string   `json:"content"`
	Position            int      `json:"position"`
	RelatedCodeExamples []string `json:"related_code_examples"`
	RelatedAnnotations  []string `json:"related_annotations"`
}

// ParagraphRepresentation 段落表示
type ParagraphRepresentation struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	SegmentID   string `json:"segment_id"`
	SegmentType string `json:"segment_type"`
	Position    int    `json:"position"`
}

// CodeSnippetRepresentation 代码片段表示
type CodeSnippetRepresentation struct {
	ID              string   `json:"id"`
	Language        string   `json:"language"`
	Code            string   `json:"code"`
	Description     string   `json:"description"`
	Position        int      `json:"position"`
	IsInline        bool     `json:"is_inline"`
	RelatedSegments []string `json:"related_segments"`
}

// ContextInjectionOptions 上下文注入选项
type ContextInjectionOptions struct {
	MaxContextSize int      `json:"max_context_size"`
	IncludeCode    bool     `json:"include_code"`
	PriorityLevel  Priority `json:"priority_level"`
	Format         Format   `json:"format"`
}

// Priority 优先级
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Format 格式
type Format string

const (
	FormatMarkdown  Format = "markdown"
	FormatJSON      Format = "json"
	FormatPlainText Format = "plain_text"
)

// ContextInjectionResult 上下文注入结果
type ContextInjectionResult struct {
	ID               string                  `json:"id"`
	DocumentID       string                  `json:"document_id"`
	Version          string                  `json:"version"`
	Query            string                  `json:"query"`
	SelectedContent  []*ContentSelection     `json:"selected_content"`
	FormattedContext string                  `json:"formatted_context"`
	Options          ContextInjectionOptions `json:"options"`
	CreatedAt        time.Time               `json:"created_at"`
}

// ContentSelection 内容选择
type ContentSelection struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Content  string  `json:"content"`
	Title    string  `json:"title"`
	Score    float64 `json:"score"`
	Position int     `json:"position"`
}
