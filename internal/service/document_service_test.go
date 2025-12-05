package service

import (
	"testing"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
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
