package service

import (
	"testing"
)

// TestEstimateTokens 测试Token估算功能
func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
		min      int
		max      int
	}{
		{
			name:     "空文本",
			text:     "",
			expected: 1, // 至少返回1
			min:      1,
			max:      1,
		},
		{
			name:     "短文本",
			text:     "Hello world",
			expected: 3, // 简单估算
			min:      1,
			max:      10,
		},
		{
			name:     "中文文本",
			text:     "你好世界",
			expected: 2,
			min:      1,
			max:      5,
		},
		{
			name:     "中英文混合",
			text:     "Hello 你好 World 世界",
			expected: 5,
			min:      3,
			max:      10,
		},
		{
			name:     "长文本",
			text:     "This is a longer text with many words and spaces. It should have an estimated token count based on the character count divided by 4.",
			expected: 30,
			min:      20,
			max:      50,
		},
		{
			name:     "带换行和标点",
			text:     "Line 1\nLine 2\nLine 3. With punctuation!",
			expected: 8,
			min:      5,
			max:      15,
		},
	}

	service := &mcpService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.estimateTokens(tt.text)
			if result < tt.min || result > tt.max {
				t.Errorf("estimateTokens() = %v, want in range [%v, %v]", result, tt.min, tt.max)
			}
			t.Logf("Text: %q, Tokens: %d", tt.text, result)
		})
	}
}

// TestTruncateWithWarning 测试带警告的文本截断功能
func TestTruncateWithWarning(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		maxLen        int
		wantTruncated bool
		expectedLen   int
	}{
		{
			name:          "不截断 - 短文本",
			text:          "Hello world",
			maxLen:        20,
			wantTruncated: false,
			expectedLen:   11,
		},
		{
			name:          "截断 - 长文本",
			text:          "This is a very long text that should be truncated",
			maxLen:        10,
			wantTruncated: true,
			expectedLen:   10 + 34, // 10字符 + 截断提示
		},
		{
			name:          "边界情况 - 刚好等于",
			text:          "Exactly",
			maxLen:        7,
			wantTruncated: true, // 等于也要截断并添加提示
			expectedLen:   7 + 34,
		},
		{
			name:          "空文本",
			text:          "",
			maxLen:        10,
			wantTruncated: false,
			expectedLen:   0,
		},
		{
			name:          "包含换行符",
			text:          "Line 1\nLine 2\nLine 3",
			maxLen:        5,
			wantTruncated: true,
			expectedLen:   5 + 34,
		},
	}

	service := &mcpService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, truncated := service.truncateWithWarning(tt.text, tt.maxLen)
			if truncated != tt.wantTruncated {
				t.Errorf("truncateWithWarning() truncated = %v, want %v", truncated, tt.wantTruncated)
			}
			if len(result) != tt.expectedLen {
				t.Errorf("truncateWithWarning() length = %v, want %v", len(result), tt.expectedLen)
			}
			t.Logf("Text length: %d, Result length: %d, Truncated: %v", len(tt.text), len(result), truncated)
		})
	}
}

// TestTruncateText 测试文本截断功能
func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "不截断",
			text:     "Hello",
			maxLen:   10,
			expected: "Hello",
		},
		{
			name:     "截断并添加省略号",
			text:     "Hello world",
			maxLen:   5,
			expected: "Hello...",
		},
		{
			name:     "空文本",
			text:     "",
			maxLen:   10,
			expected: "",
		},
		{
			name:     "正好等于最大长度",
			text:     "Exactly",
			maxLen:   7,
			expected: "Exactly",
		},
	}

	service := &mcpService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.truncateText(tt.text, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateText() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestTokenUsageLogging 测试Token使用日志记录
func TestTokenUsageLogging(t *testing.T) {
	service := &mcpService{}
	// 这个测试验证Token估算和截断功能的日志记录
	longText := "This is a very long text that will be truncated. " +
		"It contains many words to simulate a real document. " +
		"We want to ensure that the token estimation and truncation logic works correctly. " +
		"This text should be long enough to trigger the truncation mechanism when maxLen is small. " +
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
		"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

	t.Run("Truncate with token estimation", func(t *testing.T) {
		maxLen := 100
		truncated, isTruncated := service.truncateWithWarning(longText, maxLen)
		tokens := service.estimateTokens(truncated)

		if !isTruncated {
			t.Error("Expected text to be truncated")
		}

		if len(truncated) > maxLen+50 { // 允许截断提示的长度
			t.Errorf("Truncated text too long: got %d, want <= %d", len(truncated), maxLen+50)
		}

		if tokens <= 0 {
			t.Error("Token estimation should be positive")
		}

		t.Logf("Original length: %d, Truncated length: %d, Estimated tokens: %d, Truncated: %v",
			len(longText), len(truncated), tokens, isTruncated)
	})
}

// GetDocumentContentToolWithTokenParams 模拟测试 getDocumentContentTool 的 Token 参数处理
func TestGetDocumentContentToolWithTokenParams(t *testing.T) {
	tests := []struct {
		name              string
		contentLength     float64
		expectedMaxLength int
	}{
		{
			name:              "使用默认长度",
			contentLength:     0,
			expectedMaxLength: DefaultContentMaxLength,
		},
		{
			name:              "自定义长度 - 合理范围",
			contentLength:     5000,
			expectedMaxLength: 5000,
		},
		{
			name:              "超出警告阈值",
			contentLength:     100000,
			expectedMaxLength: WarningContentLength,
		},
		{
			name:              "负值 - 使用默认值",
			contentLength:     -100,
			expectedMaxLength: DefaultContentMaxLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maxContentLength := DefaultContentMaxLength

			if tt.contentLength != 0 {
				length := int(tt.contentLength)
				if length > 0 && length <= WarningContentLength {
					maxContentLength = length
				} else if length > WarningContentLength {
					maxContentLength = WarningContentLength
				}
			}

			if maxContentLength != tt.expectedMaxLength {
				t.Errorf("maxContentLength = %v, want %v", maxContentLength, tt.expectedMaxLength)
			}

			t.Logf("Content length param: %.0f, Max length: %d", tt.contentLength, maxContentLength)
		})
	}
}

// TestSearchDocumentsToolWithTokenParams 模拟测试 searchDocumentsTool 的 Token 参数处理
func TestSearchDocumentsToolWithTokenParams(t *testing.T) {
	tests := []struct {
		name              string
		contentLength     float64
		expectedMaxLength int
	}{
		{
			name:              "使用默认长度",
			contentLength:     0,
			expectedMaxLength: DefaultSearchResultLength,
		},
		{
			name:              "自定义长度 - 合理范围",
			contentLength:     300,
			expectedMaxLength: 300,
		},
		{
			name:              "超出最大限制",
			contentLength:     5000,
			expectedMaxLength: SearchResultMaxLength,
		},
		{
			name:              "负值 - 使用默认值",
			contentLength:     -50,
			expectedMaxLength: DefaultSearchResultLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultContentLength := DefaultSearchResultLength

			if tt.contentLength != 0 {
				length := int(tt.contentLength)
				if length > 0 && length <= SearchResultMaxLength {
					resultContentLength = length
				} else if length > SearchResultMaxLength {
					resultContentLength = SearchResultMaxLength
				}
			}

			if resultContentLength != tt.expectedMaxLength {
				t.Errorf("resultContentLength = %v, want %v", resultContentLength, tt.expectedMaxLength)
			}

			t.Logf("Content length param: %.0f, Result length: %d", tt.contentLength, resultContentLength)
		})
	}
}

// TestTokenConstants 验证常量配置
func TestTokenConstants(t *testing.T) {
	t.Run("验证常量值", func(t *testing.T) {
		if DefaultContentMaxLength <= 0 {
			t.Error("DefaultContentMaxLength should be positive")
		}
		if WarningContentLength <= DefaultContentMaxLength {
			t.Error("WarningContentLength should be greater than DefaultContentMaxLength")
		}
		if DefaultSearchResultLength <= 0 {
			t.Error("DefaultSearchResultLength should be positive")
		}
		if SearchResultMaxLength <= DefaultSearchResultLength {
			t.Error("SearchResultMaxLength should be greater than DefaultSearchResultLength")
		}

		t.Logf("DefaultContentMaxLength: %d", DefaultContentMaxLength)
		t.Logf("WarningContentLength: %d", WarningContentLength)
		t.Logf("DefaultSearchResultLength: %d", DefaultSearchResultLength)
		t.Logf("SearchResultMaxLength: %d", SearchResultMaxLength)
	})
}

// TestContentLengthValidation 测试内容长度验证逻辑
func TestContentLengthValidation(t *testing.T) {
	tests := []struct {
		name        string
		maxLength   int
		isValid     bool
		description string
	}{
		{
			name:        "有效长度 - 小",
			maxLength:   1000,
			isValid:     true,
			description: "小长度应该被接受",
		},
		{
			name:        "有效长度 - 中等",
			maxLength:   10000,
			isValid:     true,
			description: "中等长度应该被接受",
		},
		{
			name:        "有效长度 - 警告阈值",
			maxLength:   50000,
			isValid:     true,
			description: "警告阈值长度应该被接受",
		},
		{
			name:        "无效长度 - 零",
			maxLength:   0,
			isValid:     false,
			description: "零长度应该被拒绝",
		},
		{
			name:        "无效长度 - 负数",
			maxLength:   -100,
			isValid:     false,
			description: "负长度应该被拒绝",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.maxLength > 0 && tt.maxLength <= WarningContentLength

			if isValid != tt.isValid {
				t.Errorf("Validation result = %v, want %v (%s)", isValid, tt.isValid, tt.description)
			}
		})
	}
}
