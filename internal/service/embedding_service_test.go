package service

import (
	"context"
	"testing"
)

// TestNewMockEmbeddingService 测试创建mock嵌入服务
func TestNewMockEmbeddingService(t *testing.T) {
	service := NewMockEmbeddingService()
	if service == nil {
		t.Fatal("NewMockEmbeddingService() 返回 nil")
	}

	// 验证返回的类型
	if _, ok := service.(*mockEmbeddingService); !ok {
		t.Error("NewMockEmbeddingService() 返回的类型不正确")
	}
}

// TestMockEmbeddingService_EmptyContent 测试空内容
func TestMockEmbeddingService_EmptyContent(t *testing.T) {
	service := NewMockEmbeddingService()
	ctx := context.Background()

	tests := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "空字符串",
			content:     "",
			expectError: true,
		},
		{
			name:        "只有空格",
			content:     "   ",
			expectError: true,
		},
		{
			name:        "只有制表符",
			content:     "\t\t",
			expectError: true,
		},
		{
			name:        "只有换行符",
			content:     "\n\n",
			expectError: true,
		},
		{
			name:        "空格和换行混合",
			content:     " \t \n ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embedding, err := service.GenerateEmbedding(ctx, tt.content)

			if tt.expectError {
				if err == nil {
					t.Error("期望返回错误，但返回了nil")
				}
				if embedding != nil {
					t.Error("空内容应该返回nil向量")
				}
			} else {
				if err != nil {
					t.Errorf("不应该返回错误: %v", err)
				}
			}
		})
	}
}

// TestMockEmbeddingService_GenerateEmbedding 测试生成嵌入向量
func TestMockEmbeddingService_GenerateEmbedding(t *testing.T) {
	service := NewMockEmbeddingService()
	ctx := context.Background()

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "英文短文本",
			content: "This is a test",
		},
		{
			name:    "中文文本",
			content: "这是一段测试文本",
		},
		{
			name:    "数字和符号",
			content: "test123!@#",
		},
		{
			name:    "长文本",
			content: "This is a long text that contains many words and characters to test the embedding generation functionality.",
		},
		{
			name:    "特殊字符",
			content: "测试：abc@#$%^&*()",
		},
		{
			name:    "单个字符",
			content: "a",
		},
		{
			name:    "混合语言",
			content: "Hello 世界 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embedding, err := service.GenerateEmbedding(ctx, tt.content)

			if err != nil {
				t.Errorf("GenerateEmbedding() error = %v", err)
				return
			}

			if embedding == nil {
				t.Error("GenerateEmbedding() 返回 nil 向量")
				return
			}
		})
	}
}

// TestMockEmbeddingService_VectorDimensions 测试向量维度
func TestMockEmbeddingService_VectorDimensions(t *testing.T) {
	service := NewMockEmbeddingService()
	ctx := context.Background()

	content := "test content for dimension check"
	embedding, err := service.GenerateEmbedding(ctx, content)

	if err != nil {
		t.Fatalf("GenerateEmbedding() error = %v", err)
	}

	expectedDimensions := 384
	if len(embedding) != expectedDimensions {
		t.Errorf("向量维度 = %d, expected %d", len(embedding), expectedDimensions)
	}
}

// TestMockEmbeddingService_VectorRange 测试向量值范围
func TestMockEmbeddingService_VectorRange(t *testing.T) {
	service := NewMockEmbeddingService()
	ctx := context.Background()

	content := "test content for range check"
	embedding, err := service.GenerateEmbedding(ctx, content)

	if err != nil {
		t.Fatalf("GenerateEmbedding() error = %v", err)
	}

	// 验证向量值在合理范围内（归一化后向量长度应该接近1）
	var sum float32
	for _, v := range embedding {
		sum += v * v
	}

	normLength := sqrt(sum)

	// 归一化后的向量长度应该在 [0.999, 1.001] 范围内
	if normLength < 0.999 || normLength > 1.001 {
		t.Errorf("向量长度 = %f, expected ~1.0 (归一化后)", normLength)
	}
}

// TestMockEmbeddingService_VectorUniqueness 测试向量唯一性
func TestMockEmbeddingService_VectorUniqueness(t *testing.T) {
	service := NewMockEmbeddingService()
	ctx := context.Background()

	content1 := "first content"
	content2 := "second content"

	embedding1, err1 := service.GenerateEmbedding(ctx, content1)
	embedding2, err2 := service.GenerateEmbedding(ctx, content2)

	if err1 != nil || err2 != nil {
		t.Fatalf("GenerateEmbedding() error: err1=%v, err2=%v", err1, err2)
	}

	// 不同内容应该生成不同的向量
	if vectorEqual(embedding1, embedding2) {
		t.Error("不同内容应该生成不同的向量")
	}
}

// TestMockEmbeddingService_VectorConsistency 测试向量一致性
func TestMockEmbeddingService_VectorConsistency(t *testing.T) {
	service := NewMockEmbeddingService()
	ctx := context.Background()

	content := "consistent content"

	// 多次生成相同内容的向量
	embedding1, err1 := service.GenerateEmbedding(ctx, content)
	embedding2, err2 := service.GenerateEmbedding(ctx, content)
	embedding3, err3 := service.GenerateEmbedding(ctx, content)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("GenerateEmbedding() error: err1=%v, err2=%v, err3=%v", err1, err2, err3)
	}

	// 相同内容应该生成相同的向量
	if !vectorEqual(embedding1, embedding2) {
		t.Error("相同内容应该生成相同的向量（第1次和第2次）")
	}

	if !vectorEqual(embedding2, embedding3) {
		t.Error("相同内容应该生成相同的向量（第2次和第3次）")
	}
}

// TestMockEmbeddingService_LongContent 测试长内容
func TestMockEmbeddingService_LongContent(t *testing.T) {
	service := NewMockEmbeddingService()
	ctx := context.Background()

	// 创建8193字符的长内容（超过8192限制）
	longContent := ""
	for i := 0; i < 8193; i++ {
		longContent += "a"
	}

	embedding, err := service.GenerateEmbedding(ctx, longContent)

	// mock服务没有长度限制，应该能正常生成
	if err != nil {
		t.Errorf("GenerateEmbedding() 不应该因为内容长度而失败: %v", err)
	}

	if embedding == nil {
		t.Error("GenerateEmbedding() 返回 nil 向量")
	}

	// 验证向量维度仍然正确
	if len(embedding) != 384 {
		t.Errorf("向量维度 = %d, expected 384", len(embedding))
	}
}

// TestMockEmbeddingService_CaseSensitive 测试大小写敏感性
func TestMockEmbeddingService_CaseSensitive(t *testing.T) {
	service := NewMockEmbeddingService()
	ctx := context.Background()

	content1 := "Test Content"
	content2 := "test content"

	embedding1, err1 := service.GenerateEmbedding(ctx, content1)
	embedding2, err2 := service.GenerateEmbedding(ctx, content2)

	if err1 != nil || err2 != nil {
		t.Fatalf("GenerateEmbedding() error: err1=%v, err2=%v", err1, err2)
	}

	// 大小写不同应该生成不同的向量
	if vectorEqual(embedding1, embedding2) {
		t.Error("大小写不同应该生成不同的向量")
	}
}

// TestMockEmbeddingService_Utf8Content 测试UTF-8内容
func TestMockEmbeddingService_Utf8Content(t *testing.T) {
	service := NewMockEmbeddingService()
	ctx := context.Background()

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "日文",
			content: "こんにちは",
		},
		{
			name:    "韩文",
			content: "안녕하세요",
		},
		{
			name:    "阿拉伯文",
			content: "مرحبا",
		},
		{
			name:    "俄文",
			content: "Привет",
		},
		{
			name:    "emoji",
			content: "Hello 🌍 World 🚀",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embedding, err := service.GenerateEmbedding(ctx, tt.content)

			if err != nil {
				t.Errorf("GenerateEmbedding() error = %v", err)
				return
			}

			if embedding == nil {
				t.Error("GenerateEmbedding() 返回 nil 向量")
				return
			}

			// 验证向量维度正确
			if len(embedding) != 384 {
				t.Errorf("向量维度 = %d, expected 384", len(embedding))
			}
		})
	}
}

// TestMockEmbeddingService_ContextCancellation 测试上下文取消
func TestMockEmbeddingService_ContextCancellation(t *testing.T) {
	service := NewMockEmbeddingService()

	ctx, cancel := context.WithCancel(context.Background())

	// 立即取消上下文
	cancel()

	content := "test content"
	_, err := service.GenerateEmbedding(ctx, content)

	// mock服务是同步的，所以上下文取消可能不会立即生效
	// 这里只是测试不会panic
	if err != nil {
		// 如果因为上下文取消而失败，这是可以接受的
		t.Logf("上下文取消后生成嵌入向量错误: %v", err)
	}
}

// vectorEqual 辅助函数：比较两个向量是否相等
func vectorEqual(v1, v2 []float32) bool {
	if len(v1) != len(v2) {
		return false
	}

	for i := range v1 {
		if v1[i] != v2[i] {
			return false
		}
	}

	return true
}

// TestSimpleHash 测试简单哈希函数
func TestSimpleHash(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "空字符串",
			input: "",
		},
		{
			name:  "常见字符串",
			input: "hello world",
		},
		{
			name:  "数字字符串",
			input: "12345",
		},
		{
			name:  "特殊字符",
			input: "!@#$%^&*()",
		},
		{
			name:  "长字符串",
			input: "This is a very long string to test the hash function",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := simpleHash(tt.input)

			// 哈希值应该是整数
			_ = hash

			// 简单测试：相同输入应该产生相同哈希
			hash2 := simpleHash(tt.input)
			if hash != hash2 {
				t.Error("相同输入应该产生相同哈希")
			}
		})
	}
}

// TestSimpleHash_Uniqueness 测试哈希唯一性
func TestSimpleHash_Uniqueness(t *testing.T) {
	tests := []struct {
		name       string
		input1     string
		input2     string
		expectSame bool
	}{
		{
			name:       "不同字符串",
			input1:     "hello",
			input2:     "world",
			expectSame: false,
		},
		{
			name:       "大小写不同",
			input1:     "Hello",
			input2:     "hello",
			expectSame: false,
		},
		{
			name:       "相同字符串",
			input1:     "same",
			input2:     "same",
			expectSame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1 := simpleHash(tt.input1)
			hash2 := simpleHash(tt.input2)

			same := hash1 == hash2
			if same != tt.expectSame {
				t.Errorf("哈希值相同性 = %v, expected %v", same, tt.expectSame)
			}
		})
	}
}

// TestSqrt 测试平方根计算
func TestSqrt(t *testing.T) {
	tests := []struct {
		name  string
		input float32
		delta float32
	}{
		{
			name:  "零",
			input: 0,
			delta: 0.001,
		},
		{
			name:  "1",
			input: 1,
			delta: 0.001,
		},
		{
			name:  "4",
			input: 4,
			delta: 0.001,
		},
		{
			name:  "9",
			input: 9,
			delta: 0.001,
		},
		{
			name:  "0.25",
			input: 0.25,
			delta: 0.001,
		},
		{
			name:  "100",
			input: 100,
			delta: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sqrt(tt.input)
			expected := float32(0)

			// 计算预期值（简单情况）
			if tt.name == "零" {
				expected = 0
			} else if tt.name == "1" {
				expected = 1
			} else if tt.name == "4" {
				expected = 2
			} else if tt.name == "9" {
				expected = 3
			} else if tt.name == "0.25" {
				expected = 0.5
			} else if tt.name == "100" {
				expected = 10
			}

			if result < expected-tt.delta || result > expected+tt.delta {
				t.Errorf("sqrt(%f) = %f, expected %f (±%f)", tt.input, result, expected, tt.delta)
			}
		})
	}
}

// TestSqrt_Precision 测试平方根精度
func TestSqrt_Precision(t *testing.T) {
	// 测试 sqrt(2) 的精度
	input := float32(2.0)
	result := sqrt(input)
	expected := float32(1.41421356) // 精确的 sqrt(2)
	delta := float32(0.0001)

	if result < expected-delta || result > expected+delta {
		t.Errorf("sqrt(2) = %f, expected %f (±%f)", result, expected, delta)
	}
}
