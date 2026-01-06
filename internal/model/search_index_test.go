package model

import (
	"testing"
)

func TestVector_Value_Nil(t *testing.T) {
	var v Vector
	value, err := v.Value()
	if err != nil {
		t.Errorf("Value() error = %v, want nil", err)
	}
	if value != nil {
		t.Error("Value() should return nil for nil Vector")
	}
}

func TestVector_Value_Empty(t *testing.T) {
	v := Vector{}
	value, err := v.Value()
	if err != nil {
		t.Errorf("Value() error = %v, want nil", err)
	}
	if value == nil {
		t.Error("Value() should not return nil for empty Vector")
	}
}

func TestVector_Scan_Nil(t *testing.T) {
	var v Vector
	err := v.Scan(nil)
	if err != nil {
		t.Errorf("Scan() error = %v, want nil", err)
	}
	if v != nil {
		t.Error("Scan() should set vector to nil for nil value")
	}
}

func TestVector_Scan_ValidJSON(t *testing.T) {
	j := `[0.1,0.2,0.3]`
	var v Vector
	err := v.Scan([]byte(j))
	if err != nil {
		t.Errorf("Scan() error = %v, want nil", err)
	}
	if len(v) != 3 {
		t.Errorf("Scan() length = %d, want 3", len(v))
	}
	if v[0] != 0.1 {
		t.Errorf("Scan() v[0] = %f, want 0.1", v[0])
	}
}

func TestSearchIndex_TableName(t *testing.T) {
	index := SearchIndex{}
	if index.TableName() != "search_indices" {
		t.Errorf("TableName() = %s, want search_indices", index.TableName())
	}
}

func TestSearchIndex_Fields(t *testing.T) {
	index := SearchIndex{
		ID:            "test-id",
		DocumentID:    "doc-id",
		Version:       "1.0.0",
		Content:       "test content",
		ContentType:   "text",
		Section:       "section-1",
		Keywords:      "test,key",
		Vector:        `[0.1,0.2]`,
		Metadata:      `{"key":"value"}`,
		Score:         0.95,
		StartPosition: 0,
		EndPosition:   100,
	}

	if index.ID != "test-id" {
		t.Errorf("ID = %s, want test-id", index.ID)
	}
	if index.DocumentID != "doc-id" {
		t.Errorf("DocumentID = %s, want doc-id", index.DocumentID)
	}
	if index.Version != "1.0.0" {
		t.Errorf("Version = %s, want 1.0.0", index.Version)
	}
	if index.Content != "test content" {
		t.Errorf("Content = %s, want test content", index.Content)
	}
	if index.ContentType != "text" {
		t.Errorf("ContentType = %s, want text", index.ContentType)
	}
	if index.Section != "section-1" {
		t.Errorf("Section = %s, want section-1", index.Section)
	}
	if index.Keywords != "test,key" {
		t.Errorf("Keywords = %s, want test,key", index.Keywords)
	}
	if index.Score != 0.95 {
		t.Errorf("Score = %f, want 0.95", index.Score)
	}
	if index.StartPosition != 0 {
		t.Errorf("StartPosition = %d, want 0", index.StartPosition)
	}
	if index.EndPosition != 100 {
		t.Errorf("EndPosition = %d, want 100", index.EndPosition)
	}
}

func TestSearchRequest_Fields(t *testing.T) {
	filters := map[string]interface{}{
		"library": "test-lib",
		"type":    "markdown",
	}

	req := SearchRequest{
		Query:      "test query",
		Filters:    filters,
		Page:       1,
		Size:       10,
		SearchType: "keyword",
	}

	if req.Query != "test query" {
		t.Errorf("Query = %s, want test query", req.Query)
	}
	if req.Filters["library"] != "test-lib" {
		t.Error("Filters not set correctly")
	}
	if req.Page != 1 {
		t.Errorf("Page = %d, want 1", req.Page)
	}
	if req.Size != 10 {
		t.Errorf("Size = %d, want 10", req.Size)
	}
	if req.SearchType != "keyword" {
		t.Errorf("SearchType = %s, want keyword", req.SearchType)
	}
}

func TestSearchResponse_Fields(t *testing.T) {
	items := []SearchResult{
		{ID: "1", DocumentID: "doc-1", Score: 0.9},
		{ID: "2", DocumentID: "doc-2", Score: 0.8},
	}

	resp := SearchResponse{
		Total: 100,
		Items: items,
		Page:  1,
		Size:  10,
	}

	if resp.Total != 100 {
		t.Errorf("Total = %d, want 100", resp.Total)
	}
	if len(resp.Items) != 2 {
		t.Errorf("Items length = %d, want 2", len(resp.Items))
	}
	if resp.Page != 1 {
		t.Errorf("Page = %d, want 1", resp.Page)
	}
	if resp.Size != 10 {
		t.Errorf("Size = %d, want 10", resp.Size)
	}
}

func TestSearchResult_Fields(t *testing.T) {
	result := SearchResult{
		ID:          "result-id",
		DocumentID:  "doc-123",
		Version:     "1.0.0",
		Title:       "Test Document",
		Library:     "test-lib",
		Content:     "Document content",
		Snippet:     "Snippet...",
		Score:       0.95,
		ContentType: "text",
		Section:     "Introduction",
		Metadata: map[string]interface{}{
			"author": "test author",
			"date":   "2024-01-01",
		},
	}

	if result.ID != "result-id" {
		t.Errorf("ID = %s, want result-id", result.ID)
	}
	if result.DocumentID != "doc-123" {
		t.Errorf("DocumentID = %s, want doc-123", result.DocumentID)
	}
	if result.Version != "1.0.0" {
		t.Errorf("Version = %s, want 1.0.0", result.Version)
	}
	if result.Title != "Test Document" {
		t.Errorf("Title = %s, want Test Document", result.Title)
	}
	if result.Library != "test-lib" {
		t.Errorf("Library = %s, want test-lib", result.Library)
	}
	if result.Content != "Document content" {
		t.Errorf("Content = %s, want Document content", result.Content)
	}
	if result.Snippet != "Snippet..." {
		t.Errorf("Snippet = %s, want Snippet...", result.Snippet)
	}
	if result.Score != 0.95 {
		t.Errorf("Score = %f, want 0.95", result.Score)
	}
	if result.ContentType != "text" {
		t.Errorf("ContentType = %s, want text", result.ContentType)
	}
	if result.Section != "Introduction" {
		t.Errorf("Section = %s, want Introduction", result.Section)
	}
	if result.Metadata["author"] != "test author" {
		t.Error("Metadata not set correctly")
	}
}
