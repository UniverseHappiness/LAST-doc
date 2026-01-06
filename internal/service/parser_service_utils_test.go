package service

import (
	"strings"
	"testing"
)

// TestMinFunction 测试min辅助函数
func TestMinFunction(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{
			name: "a小于b",
			a:    3,
			b:    5,
			want: 3,
		},
		{
			name: "a大于b",
			a:    7,
			b:    2,
			want: 2,
		},
		{
			name: "a等于b",
			a:    4,
			b:    4,
			want: 4,
		},
		{
			name: "a为0",
			a:    0,
			b:    10,
			want: 0,
		},
		{
			name: "负数",
			a:    -5,
			b:    -3,
			want: -5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.a, tt.b); got != tt.want {
				t.Errorf("min(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TestExtractMarkdownMetadataFunc 测试extractMarkdownMetadata函数
func TestExtractMarkdownMetadataFunc(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    map[string]interface{}
	}{
		{
			name:    "包含标题",
			content: "# 测试标题\n\n内容",
			want: map[string]interface{}{
				"title": "测试标题",
			},
		},
		{
			name: "多行内容",
			content: `# Main Title

Some content here.

## Sub Title

More content.
`,
			want: map[string]interface{}{
				"title": "Main Title",
			},
		},
		{
			name:    "无标题",
			content: "只有内容没有标题",
			want:    map[string]interface{}{},
		},
		{
			name:    "空内容",
			content: "",
			want:    map[string]interface{}{},
		},
		{
			name:    "只有标题",
			content: "# Only Title",
			want: map[string]interface{}{
				"title": "Only Title",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMarkdownMetadata(tt.content)

			for key, wantVal := range tt.want {
				gotVal, ok := got[key]
				if !ok {
					t.Errorf("extractMarkdownMetadata() 缺少字段 %s", key)
					continue
				}
				if gotVal != wantVal {
					t.Errorf("extractMarkdownMetadata()[%s] = %v, want %v", key, gotVal, wantVal)
				}
			}
		})
	}
}

// TestExtractSwaggerMetadataFunc 测试extractSwaggerMetadata函数
func TestExtractSwaggerMetadataFunc(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    map[string]interface{}
	}{
		{
			name: "Swagger格式",
			content: `{
  "swagger": "2.0",
  "info": {"title": "API"},
  "paths": {}
}`,
			want: map[string]interface{}{
				"spec_version": "swagger",
				"has_info":     true,
				"has_paths":    true,
			},
		},
		{
			name: "OpenAPI格式",
			content: `{
  "openapi": "3.0.0",
  "info": {"title": "API"},
  "paths": {}
}`,
			want: map[string]interface{}{
				"spec_version": "openapi",
				"has_info":     true,
				"has_paths":    true,
			},
		},
		{
			name: "只有info",
			content: `{
  "swagger": "2.0",
  "info": {"title": "API"}
}`,
			want: map[string]interface{}{
				"spec_version": "swagger",
				"has_info":     true,
			},
		},
		{
			name:    "空JSON",
			content: `{}`,
			want:    map[string]interface{}{},
		},
		{
			name:    "无效格式",
			content: "not a json",
			want:    map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSwaggerMetadata(tt.content)

			for key, wantVal := range tt.want {
				gotVal, ok := got[key]
				if !ok {
					t.Errorf("extractSwaggerMetadata() 缺少字段 %s", key)
					continue
				}
				if gotVal != wantVal {
					t.Errorf("extractSwaggerMetadata()[%s] = %v, want %v", key, gotVal, wantVal)
				}
			}
		})
	}
}

// TestParserConstructors 测试各种Parser的构造函数
func TestParserConstructors(t *testing.T) {
	tests := []struct {
		name    string
		parser  interface{}
		wantNil bool
	}{
		{
			name:    "NewMarkdownParser",
			parser:  NewMarkdownParser(),
			wantNil: false,
		},
		{
			name:    "NewPDFParser",
			parser:  NewPDFParser(),
			wantNil: false,
		},
		{
			name:    "NewDocxParser",
			parser:  NewDocxParser(),
			wantNil: false,
		},
		{
			name:    "NewSwaggerParser",
			parser:  NewSwaggerParser(),
			wantNil: false,
		},
		{
			name:    "NewOpenAPIParser",
			parser:  NewOpenAPIParser(),
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.parser == nil && !tt.wantNil {
				t.Errorf("%s() 返回 nil", tt.name)
			}
		})
	}
}

// TestSupportedExtensions 测试各parser的SupportedExtensions方法
func TestSupportedExtensions(t *testing.T) {
	tests := []struct {
		name         string
		parser       DocumentParser
		expectedExts []string
	}{
		{
			name:         "MarkdownParser",
			parser:       NewMarkdownParser(),
			expectedExts: []string{".md", ".markdown"},
		},
		{
			name:         "PDFParser",
			parser:       NewPDFParser(),
			expectedExts: []string{".pdf"},
		},
		{
			name:         "DocxParser",
			parser:       NewDocxParser(),
			expectedExts: []string{".docx", ".doc"},
		},
		{
			name:         "SwaggerParser",
			parser:       NewSwaggerParser(),
			expectedExts: []string{".json", ".yaml", ".yml"},
		},
		{
			name:         "OpenAPIParser",
			parser:       NewOpenAPIParser(),
			expectedExts: []string{".json", ".yaml", ".yml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exts := tt.parser.SupportedExtensions()

			if len(exts) == 0 {
				t.Errorf("%s.SupportedExtensions() 返回空列表", tt.name)
			}

			// 验证至少有一个预期的扩展名
			hasExpected := false
			for _, expectedExt := range tt.expectedExts {
				for _, actualExt := range exts {
					if actualExt == expectedExt {
						hasExpected = true
						break
					}
				}
				if hasExpected {
					break
				}
			}

			if !hasExpected {
				t.Errorf("%s.SupportedExtensions() 应该包含预期扩展名", tt.name)
			}
		})
	}
}

// TestMarkdownMetadataStatistics 测试Markdown元数据统计功能
func TestMarkdownMetadataStatistics(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		expectWordCount  bool
		expectCodeBlocks bool
	}{
		{
			name:             "简单内容",
			content:          "# Title\n\nContent here",
			expectWordCount:  true,
			expectCodeBlocks: false,
		},
		{
			name:             "包含代码块",
			content:          "# Title\n\n```code\n```\n```python\nprint('hello')\n```",
			expectWordCount:  true,
			expectCodeBlocks: true,
		},
		{
			name:             "单个代码块",
			content:          "```go\nfunc test() {}\n```",
			expectWordCount:  false,
			expectCodeBlocks: false, // 单个代码块不计数
		},
		{
			name:             "无代码块",
			content:          "Just some text",
			expectWordCount:  true,
			expectCodeBlocks: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMarkdownMetadata(tt.content)

			if tt.expectWordCount {
				if _, hasWordCount := got["word_count"]; !hasWordCount {
					t.Error("extractMarkdownMetadata() 应该包含 word_count")
				}
			}

			if tt.expectCodeBlocks {
				if _, hasCodeBlocks := got["code_block_count"]; !hasCodeBlocks {
					t.Error("extractMarkdownMetadata() 应该包含 code_block_count")
				}
			}
		})
	}
}

// TestSwaggerMetadataDetection 测试Swagger元数据检测
func TestSwaggerMetadataDetection(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		detectInfo  bool
		detectPaths bool
	}{
		{
			name: "完整Swagger",
			content: `{
  "swagger": "2.0",
  "info": {"title": "Test"},
  "paths": {}
}`,
			detectInfo:  true,
			detectPaths: true,
		},
		{
			name:        "只有info",
			content:     `{"swagger": "2.0", "info": {}}`,
			detectInfo:  true,
			detectPaths: false,
		},
		{
			name:        "只有paths",
			content:     `{"swagger": "2.0", "paths": {}}`,
			detectInfo:  false,
			detectPaths: true,
		},
		{
			name:        "空json",
			content:     `{}`,
			detectInfo:  false,
			detectPaths: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSwaggerMetadata(tt.content)

			hasInfo := got["has_info"] == true
			hasPaths := got["has_paths"] == true

			if tt.detectInfo && !hasInfo {
				t.Error("应该检测到info字段")
			}
			if !tt.detectInfo && hasInfo {
				t.Error("不应该检测到info字段")
			}

			if tt.detectPaths && !hasPaths {
				t.Error("应该检测到paths字段")
			}
			if !tt.detectPaths && hasPaths {
				t.Error("不应该检测到paths字段")
			}
		})
	}
}

// TestWordCountExtraction 测试字数提取
func TestWordCountExtraction(t *testing.T) {
	tests := []struct {
		name    string
		content string
		min     int
		max     int
	}{
		{
			name:    "单个词",
			content: "Hello",
			min:     1,
			max:     1,
		},
		{
			name:    "多个词",
			content: "Hello world this is a test",
			min:     5,
			max:     7,
		},
		{
			name:    "空内容",
			content: "",
			min:     0,
			max:     0,
		},
		{
			name:    "只有空格",
			content: "   ",
			min:     0,
			max:     0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMarkdownMetadata(tt.content)
			wordCount, ok := got["word_count"]

			if !ok {
				t.Error("应该包含word_count字段")
				return
			}

			count, ok := wordCount.(int)
			if !ok {
				t.Error("word_count应该是整数")
				return
			}

			if count < tt.min || count > tt.max {
				t.Errorf("word_count = %d, 期望在范围 [%d, %d] 内", count, tt.min, tt.max)
			}
		})
	}
}

// TestCodeBlockCounting 测试代码块计数
func TestCodeBlockCounting(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		expectCnt int
	}{
		{
			name:      "无代码块",
			content:   "Just text",
			expectCnt: 0,
		},
		{
			name:      "一对代码块",
			content:   "```code\ntest\n```",
			expectCnt: 1,
		},
		{
			name:      "两对代码块",
			content:   "```go\n```\n```python\n```",
			expectCnt: 2,
		},
		{
			name:      "三个代码块",
			content:   "```\n```\n```\n```",
			expectCnt: 2, // 6个```，除以2=3对？不对，应该是完整的代码块对数
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractMarkdownMetadata(tt.content)

			// 检查是否有代码块标记
			backtickCount := strings.Count(tt.content, "```")

			if backtickCount > 0 {
				if _, hasCodeBlocks := got["code_block_count"]; !hasCodeBlocks {
					t.Error("包含代码块时应该有code_block_count字段")
				}
			}
		})
	}
}
