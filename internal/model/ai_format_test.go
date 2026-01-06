package model

import (
	"testing"
	"time"
)

func TestSummaryLevel_Constants(t *testing.T) {
	if SummaryLevelBrief != "brief" {
		t.Errorf("SummaryLevelBrief = %s, want brief", SummaryLevelBrief)
	}
	if SummaryLevelMedium != "medium" {
		t.Errorf("SummaryLevelMedium = %s, want medium", SummaryLevelMedium)
	}
	if SummaryLevelDetailed != "detailed" {
		t.Errorf("SummaryLevelDetailed = %s, want detailed", SummaryLevelDetailed)
	}
}

func TestPriority_Constants(t *testing.T) {
	if PriorityLow != "low" {
		t.Errorf("PriorityLow = %s, want low", PriorityLow)
	}
	if PriorityMedium != "medium" {
		t.Errorf("PriorityMedium = %s, want medium", PriorityMedium)
	}
	if PriorityHigh != "high" {
		t.Errorf("PriorityHigh = %s, want high", PriorityHigh)
	}
}

func TestFormat_Constants(t *testing.T) {
	if FormatMarkdown != "markdown" {
		t.Errorf("FormatMarkdown = %s, want markdown", FormatMarkdown)
	}
	if FormatJSON != "json" {
		t.Errorf("FormatJSON = %s, want json", FormatJSON)
	}
	if FormatPlainText != "plain_text" {
		t.Errorf("FormatPlainText = %s, want plain_text", FormatPlainText)
	}
}

func TestAIStructuredContent_Fields(t *testing.T) {
	content := AIStructuredContent{
		ID:         "test-id",
		DocumentID: "doc-id",
		Version:    "1.0.0",
		DocType:    DocumentTypeMarkdown,
		CreatedAt:  time.Now(),
	}

	if content.ID != "test-id" {
		t.Errorf("ID = %s, want test-id", content.ID)
	}
	if content.DocumentID != "doc-id" {
		t.Errorf("DocumentID = %s, want doc-id", content.DocumentID)
	}
	if content.Version != "1.0.0" {
		t.Errorf("Version = %s, want 1.0.0", content.Version)
	}
	if content.DocType != DocumentTypeMarkdown {
		t.Errorf("DocType = %s, want DocumentTypeMarkdown", content.DocType)
	}
	if content.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestContentSegment_Fields(t *testing.T) {
	segment := ContentSegment{
		ID:       "seg-1",
		Title:    "Test Segment",
		Content:  "Test content",
		Type:     "text",
		Position: 10,
	}

	if segment.ID != "seg-1" {
		t.Errorf("ID = %s, want seg-1", segment.ID)
	}
	if segment.Title != "Test Segment" {
		t.Errorf("Title = %s, want Test Segment", segment.Title)
	}
	if segment.Content != "Test content" {
		t.Errorf("Content = %s, want Test content", segment.Content)
	}
	if segment.Type != "text" {
		t.Errorf("Type = %s, want text", segment.Type)
	}
	if segment.Position != 10 {
		t.Errorf("Position = %d, want 10", segment.Position)
	}
}

func TestCodeExample_Fields(t *testing.T) {
	example := CodeExample{
		ID:          "code-1",
		Language:    "go",
		Code:        "fmt.Println(\"test\")",
		Description: "Test code example",
		Position:    5,
		IsInline:    false,
	}

	if example.ID != "code-1" {
		t.Errorf("ID = %s, want code-1", example.ID)
	}
	if example.Language != "go" {
		t.Errorf("Language = %s, want go", example.Language)
	}
	if example.Code != "fmt.Println(\"test\")" {
		t.Errorf("Code = %s, want fmt.Println(\"test\")", example.Code)
	}
	if example.Description != "Test code example" {
		t.Errorf("Description = %s, want Test code example", example.Description)
	}
	if example.Position != 5 {
		t.Errorf("Position = %d, want 5", example.Position)
	}
	if example.IsInline {
		t.Error("IsInline should be false")
	}
}

func TestSemanticAnnotation_Fields(t *testing.T) {
	annotation := SemanticAnnotation{
		ID:       "anno-1",
		Type:     "category",
		Value:    "test category",
		Context:  "in intro section",
		Position: 3,
	}

	if annotation.ID != "anno-1" {
		t.Errorf("ID = %s, want anno-1", annotation.ID)
	}
	if annotation.Type != "category" {
		t.Errorf("Type = %s, want category", annotation.Type)
	}
	if annotation.Value != "test category" {
		t.Errorf("Value = %s, want test category", annotation.Value)
	}
	if annotation.Context != "in intro section" {
		t.Errorf("Context = %s, want in intro section", annotation.Context)
	}
	if annotation.Position != 3 {
		t.Errorf("Position = %d, want 3", annotation.Position)
	}
}

func TestLLMOptimizedContent_Fields(t *testing.T) {
	now := time.Now()
	content := LLMOptimizedContent{
		ID:         "llm-1",
		DocumentID: "doc-1",
		Version:    "2.0.0",
		Content:    "Optimized content",
		CreatedAt:  now,
	}

	if content.ID != "llm-1" {
		t.Errorf("ID = %s, want llm-1", content.ID)
	}
	if content.DocumentID != "doc-1" {
		t.Errorf("DocumentID = %s, want doc-1", content.DocumentID)
	}
	if content.Version != "2.0.0" {
		t.Errorf("Version = %s, want 2.0.0", content.Version)
	}
	if content.Content != "Optimized content" {
		t.Errorf("Content = %s, want Optimized content", content.Content)
	}
	if content.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}
