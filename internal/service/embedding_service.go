package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// EmbeddingService 嵌入向量服务接口
type EmbeddingService interface {
	GenerateEmbedding(ctx context.Context, content string) ([]float32, error)
}

// openAIEmbeddingService OpenAI 嵌入向量服务实现
type openAIEmbeddingService struct {
	client *openai.Client
	model  openai.EmbeddingModel
}

// NewOpenAIEmbeddingService 创建 OpenAI 嵌入向量服务实例
func NewOpenAIEmbeddingService(apiKey, modelStr string) EmbeddingService {
	// 如果未提供 API Key，尝试从环境变量获取
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	// 如果仍未提供 API Key，使用空的客户端（用于测试或模拟）
	config := openai.DefaultConfig(apiKey)

	// 可以自定义 OpenAI API 基础 URL，以支持兼容的服务
	if baseURL := os.Getenv("OPENAI_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	client := openai.NewClientWithConfig(config)

	// 默认模型
	var model openai.EmbeddingModel
	if modelStr == "" {
		model = openai.AdaEmbeddingV2
	} else {
		model = openai.EmbeddingModel(modelStr)
	}

	return &openAIEmbeddingService{
		client: client,
		model:  model,
	}
}

// GenerateEmbedding 生成文本的嵌入向量
func (s *openAIEmbeddingService) GenerateEmbedding(ctx context.Context, content string) ([]float32, error) {
	// 检查内容是否为空
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content is empty")
	}

	// 截断过长的内容，OpenAI API 有输入限制
	if len(content) > 8192 {
		content = content[:8192]
		log.Printf("Warning: Content truncated to 8192 characters for embedding generation")
	}

	// 创建嵌入请求
	req := openai.EmbeddingRequest{
		Input: []string{content},
		Model: s.model,
	}

	// 调用 OpenAI API
	resp, err := s.client.CreateEmbeddings(ctx, req)
	if err != nil {
		log.Printf("Error generating embedding: %v", err)
		return nil, fmt.Errorf("failed to generate embedding: %v", err)
	}

	// 检查响应
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	// 返回嵌入向量
	return resp.Data[0].Embedding, nil
}

// mockEmbeddingService 模拟嵌入向量服务实现（用于测试或当 OpenAI 服务不可用时）
type mockEmbeddingService struct{}

// NewMockEmbeddingService 创建模拟嵌入向量服务实例
func NewMockEmbeddingService() EmbeddingService {
	return &mockEmbeddingService{}
}

// GenerateEmbedding 生成模拟的文本嵌入向量
func (s *mockEmbeddingService) GenerateEmbedding(ctx context.Context, content string) ([]float32, error) {
	// 检查内容是否为空
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content is empty")
	}

	// 生成模拟的 384 维向量（与 sentence-transformers/all-MiniLM-L6-v2 的维度相同）
	vector := make([]float32, 384)

	// 使用简单哈希生成模拟向量
	for i := 0; i < len(vector); i++ {
		// 基于内容和位置生成伪随机值
		hash := simpleHash(fmt.Sprintf("%s-%d", content, i))
		// 转换为 -1 到 1 之间的浮点数
		vector[i] = float32(hash%1000)/500.0 - 1.0
	}

	// 归一化向量
	norm := float32(0)
	for _, v := range vector {
		norm += v * v
	}
	if norm > 0 {
		norm = sqrt(norm)
		for i := range vector {
			vector[i] /= norm
		}
	}

	return vector, nil
}
