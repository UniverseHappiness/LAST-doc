package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
)

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DocumentParserService 解析服务接口
type DocumentParserService interface {
	ParseDocument(ctx context.Context, filePath string, docType model.DocumentType) (string, map[string]interface{}, error)
}

// parserService 解析服务实现
type parserService struct {
	parsers    map[model.DocumentType]DocumentParser
	grpcClient *GRPCClient
}

// NewParserService 创建解析服务实例
func NewParserService() DocumentParserService {
	service := &parserService{
		parsers:    make(map[model.DocumentType]DocumentParser),
		grpcClient: NewGRPCClient(),
	}

	// 注册各种文档类型的解析器
	service.RegisterParser(model.DocumentTypeMarkdown, NewMarkdownParser())
	service.RegisterParser(model.DocumentTypeSwagger, NewSwaggerParser())
	service.RegisterParser(model.DocumentTypeOpenAPI, NewOpenAPIParser())
	service.RegisterParser(model.DocumentTypeJavaDoc, NewJavaDocParser())

	// 从环境变量获取 gRPC 服务器地址
	grpcHost := os.Getenv("GRPC_SERVER_HOST")
	grpcPort := os.Getenv("GRPC_SERVER_PORT")

	if grpcHost == "" {
		grpcHost = "localhost"
	}
	if grpcPort == "" {
		grpcPort = "50051"
	}

	grpcAddr := grpcHost + ":" + grpcPort

	// 连接到gRPC服务
	if err := service.grpcClient.Connect(grpcAddr); err != nil {
		fmt.Printf("警告：连接gRPC服务失败(%s)，将使用本地解析器: %v\n", grpcAddr, err)
		// 如果连接失败，回退到本地解析器
		service.RegisterParser(model.DocumentTypePDF, NewPDFParser())
		service.RegisterParser(model.DocumentTypeDocx, NewDocxParser())
	} else {
		// 连接成功，使用gRPC解析器
		service.RegisterParser(model.DocumentTypePDF, NewPDFGRPCParser(service.grpcClient))
		service.RegisterParser(model.DocumentTypeDocx, NewDocxGRPCParser(service.grpcClient))
	}

	return service
}

// RegisterParser 注册文档解析器
func (s *parserService) RegisterParser(docType model.DocumentType, parser DocumentParser) {
	s.parsers[docType] = parser
}

// ParseDocument 解析文档
func (s *parserService) ParseDocument(ctx context.Context, filePath string, docType model.DocumentType) (string, map[string]interface{}, error) {
	parser, ok := s.parsers[docType]
	if !ok {
		return "", nil, fmt.Errorf("unsupported document type: %s", docType)
	}

	return parser.Parse(ctx, filePath)
}

// DocumentParser 文档解析器接口
type DocumentParser interface {
	Parse(ctx context.Context, filePath string) (string, map[string]interface{}, error)
	SupportedExtensions() []string
}

// markdownParser Markdown解析器
type markdownParser struct{}

// NewMarkdownParser 创建Markdown解析器
func NewMarkdownParser() DocumentParser {
	return &markdownParser{}
}

// Parse 解析Markdown文档
func (p *markdownParser) Parse(ctx context.Context, filePath string) (string, map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read markdown file: %v", err)
	}

	content := string(data)
	metadata := extractMarkdownMetadata(content)

	return content, metadata, nil
}

// SupportedExtensions 返回支持的文件扩展名
func (p *markdownParser) SupportedExtensions() []string {
	return []string{".md", ".markdown"}
}

// extractMarkdownMetadata 提取Markdown元数据
func extractMarkdownMetadata(content string) map[string]interface{} {
	metadata := make(map[string]interface{})

	lines := strings.Split(content, "\n")

	// 提取标题
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			metadata["title"] = strings.TrimSpace(line[2:])
			break
		}
	}

	// 统计字数
	metadata["word_count"] = len(strings.Fields(content))

	// 统计代码块数量
	codeBlockCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "```" {
			codeBlockCount++
		}
	}
	metadata["code_block_count"] = codeBlockCount / 2

	return metadata
}

// pdfParser PDF解析器
type pdfParser struct{}

// NewPDFParser 创建PDF解析器
func NewPDFParser() DocumentParser {
	return &pdfParser{}
}

// Parse 解析PDF文档
func (p *pdfParser) Parse(ctx context.Context, filePath string) (string, map[string]interface{}, error) {
	// 添加诊断日志
	fmt.Printf("[DEBUG] PDF解析开始 - 文件路径: %s\n", filePath)

	// 这里应该使用PDF解析库，如github.com/ledongthuc/pdf
	// 为了简化，这里只返回模拟数据
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[DEBUG] PDF文件读取失败 - 路径: %s, 错误: %v\n", filePath, err)
		return "", nil, fmt.Errorf("failed to read PDF file: %v", err)
	}

	// PDF文件是二进制格式，绝对不要直接输出内容到终端
	// 只输出文件大小信息，避免二进制内容导致终端乱码
	fmt.Printf("[DEBUG] PDF文件读取成功 - 路径: %s, 大小: %d 字节\n", filePath, len(data))

	// 模拟提取文本内容
	content := fmt.Sprintf("PDF文档内容（大小：%d字节）", len(data))
	metadata := map[string]interface{}{
		"file_size": len(data),
		"type":      "pdf",
		"pages":     10, // 模拟页数
	}

	fmt.Printf("[DEBUG] PDF解析完成 - 路径: %s, 内容长度: %d\n", filePath, len(content))
	return content, metadata, nil
}

// SupportedExtensions 返回支持的文件扩展名
func (p *pdfParser) SupportedExtensions() []string {
	return []string{".pdf"}
}

// docxParser DOCX解析器
type docxParser struct{}

// NewDocxParser 创建DOCX解析器
func NewDocxParser() DocumentParser {
	return &docxParser{}
}

// Parse 解析DOCX文档
func (p *docxParser) Parse(ctx context.Context, filePath string) (string, map[string]interface{}, error) {
	// 添加诊断日志
	fmt.Printf("[DEBUG] DOCX解析开始 - 文件路径: %s\n", filePath)

	// 这里应该使用DOCX解析库，如github.com/unidoc/unioffice
	// 为了简化，这里只返回模拟数据
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[DEBUG] DOCX文件读取失败 - 路径: %s, 错误: %v\n", filePath, err)
		return "", nil, fmt.Errorf("failed to read DOCX file: %v", err)
	}

	// DOCX文件是二进制格式，绝对不要直接输出内容到终端
	// 只输出文件大小信息，避免二进制内容导致终端乱码
	fmt.Printf("[DEBUG] DOCX文件读取成功 - 路径: %s, 大小: %d 字节\n", filePath, len(data))

	// 模拟提取文本内容
	content := fmt.Sprintf("DOCX文档内容（大小：%d字节）", len(data))
	metadata := map[string]interface{}{
		"file_size": len(data),
		"type":      "docx",
		"pages":     5, // 模拟页数
	}

	fmt.Printf("[DEBUG] DOCX解析完成 - 路径: %s, 内容长度: %d\n", filePath, len(content))
	return content, metadata, nil
}

// SupportedExtensions 返回支持的文件扩展名
func (p *docxParser) SupportedExtensions() []string {
	return []string{".docx", ".doc"}
}

// swaggerParser Swagger解析器
type swaggerParser struct{}

// NewSwaggerParser 创建Swagger解析器
func NewSwaggerParser() DocumentParser {
	return &swaggerParser{}
}

// Parse 解析Swagger文档
func (p *swaggerParser) Parse(ctx context.Context, filePath string) (string, map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read Swagger file: %v", err)
	}

	content := string(data)
	metadata := extractSwaggerMetadata(content)

	return content, metadata, nil
}

// SupportedExtensions 返回支持的文件扩展名
func (p *swaggerParser) SupportedExtensions() []string {
	return []string{".json", ".yaml", ".yml"}
}

// extractSwaggerMetadata 提取Swagger元数据
func extractSwaggerMetadata(content string) map[string]interface{} {
	metadata := make(map[string]interface{})

	// 简单的元数据提取，实际应用中应该使用Swagger/OpenAPI解析库
	if strings.Contains(content, "\"swagger\"") {
		metadata["spec_version"] = "swagger"
	} else if strings.Contains(content, "\"openapi\"") {
		metadata["spec_version"] = "openapi"
	}

	if strings.Contains(content, "\"info\"") {
		metadata["has_info"] = true
	}

	if strings.Contains(content, "\"paths\"") {
		metadata["has_paths"] = true
	}

	return metadata
}

// openAPIParser OpenAPI解析器
type openAPIParser struct{}

// NewOpenAPIParser 创建OpenAPI解析器
func NewOpenAPIParser() DocumentParser {
	return &openAPIParser{}
}

// Parse 解析OpenAPI文档
func (p *openAPIParser) Parse(ctx context.Context, filePath string) (string, map[string]interface{}, error) {
	// OpenAPI和Swagger格式类似，可以使用相同的解析逻辑
	swaggerParser := NewSwaggerParser()
	return swaggerParser.Parse(ctx, filePath)
}

// SupportedExtensions 返回支持的文件扩展名
func (p *openAPIParser) SupportedExtensions() []string {
	return []string{".json", ".yaml", ".yml"}
}

// javaDocParser JavaDoc解析器
type javaDocParser struct{}

// NewJavaDocParser 创建JavaDoc解析器
func NewJavaDocParser() DocumentParser {
	return &javaDocParser{}
}

// Parse 解析JavaDoc文档
func (p *javaDocParser) Parse(ctx context.Context, filePath string) (string, map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read JavaDoc file: %v", err)
	}

	content := string(data)
	metadata := extractJavaDocMetadata(content)

	return content, metadata, nil
}

// SupportedExtensions 返回支持的文件扩展名
func (p *javaDocParser) SupportedExtensions() []string {
	return []string{".html", ".htm"}
}

// extractJavaDocMetadata 提取JavaDoc元数据
func extractJavaDocMetadata(content string) map[string]interface{} {
	metadata := make(map[string]interface{})

	// 简单的JavaDoc元数据提取
	classCount := strings.Count(content, "class=\"")
	methodCount := strings.Count(content, "method=\"")

	metadata["class_count"] = classCount
	metadata["method_count"] = methodCount

	return metadata
}

// GetParserByExtension 根据文件扩展名获取解析器
func GetParserByExtension(filePath string) DocumentParser {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".md", ".markdown":
		return NewMarkdownParser()
	case ".pdf":
		return NewPDFParser()
	case ".docx", ".doc":
		return NewDocxParser()
	case ".json", ".yaml", ".yml":
		// 简单判断是Swagger还是OpenAPI或普通JSON/YAML
		if data, err := os.ReadFile(filePath); err == nil {
			content := string(data)
			if strings.Contains(content, "\"swagger\"") || strings.Contains(content, "\"openapi\"") {
				return NewSwaggerParser()
			}
		}
		return nil
	case ".html", ".htm":
		return NewJavaDocParser()
	default:
		return nil
	}
}

// pdfGRPCParser PDF gRPC解析器
type pdfGRPCParser struct {
	grpcClient *GRPCClient
}

// NewPDFGRPCParser 创建PDF gRPC解析器
func NewPDFGRPCParser(grpcClient *GRPCClient) DocumentParser {
	return &pdfGRPCParser{
		grpcClient: grpcClient,
	}
}

// Parse 通过gRPC解析PDF文档
func (p *pdfGRPCParser) Parse(ctx context.Context, filePath string) (string, map[string]interface{}, error) {
	fmt.Printf("[DEBUG] 通过gRPC解析PDF文档 - 文件路径: %s\n", filePath)

	content, metadata, err := p.grpcClient.ParsePDFWithGRPC(filePath)
	if err != nil {
		fmt.Printf("[DEBUG] gRPC PDF解析失败 - 路径: %s, 错误: %v\n", filePath, err)
		// 如果gRPC解析失败，回退到本地解析
		localParser := NewPDFParser()
		return localParser.Parse(ctx, filePath)
	}

	fmt.Printf("[DEBUG] gRPC PDF解析成功 - 路径: %s, 内容长度: %d\n", filePath, len(content))
	return content, metadata, nil
}

// SupportedExtensions 返回支持的文件扩展名
func (p *pdfGRPCParser) SupportedExtensions() []string {
	return []string{".pdf"}
}

// docxGRPCParser DOCX gRPC解析器
type docxGRPCParser struct {
	grpcClient *GRPCClient
}

// NewDocxGRPCParser 创建DOCX gRPC解析器
func NewDocxGRPCParser(grpcClient *GRPCClient) DocumentParser {
	return &docxGRPCParser{
		grpcClient: grpcClient,
	}
}

// Parse 通过gRPC解析DOCX文档
func (p *docxGRPCParser) Parse(ctx context.Context, filePath string) (string, map[string]interface{}, error) {
	fmt.Printf("[DEBUG] 通过gRPC解析DOCX文档 - 文件路径: %s\n", filePath)

	content, metadata, err := p.grpcClient.ParseDOCXWithGRPC(filePath)
	if err != nil {
		fmt.Printf("[DEBUG] gRPC DOCX解析失败 - 路径: %s, 错误: %v\n", filePath, err)
		// 如果gRPC解析失败，回退到本地解析
		localParser := NewDocxParser()
		return localParser.Parse(ctx, filePath)
	}

	fmt.Printf("[DEBUG] gRPC DOCX解析成功 - 路径: %s, 内容长度: %d\n", filePath, len(content))
	return content, metadata, nil
}

// SupportedExtensions 返回支持的文件扩展名
func (p *docxGRPCParser) SupportedExtensions() []string {
	return []string{".docx", ".doc"}
}
