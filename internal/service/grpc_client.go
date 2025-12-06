package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/UniverseHappiness/LAST-doc/proto"
)

// GRPCClient gRPC客户端
type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.DocumentParserServiceClient
}

// NewGRPCClient 创建gRPC客户端
func NewGRPCClient() *GRPCClient {
	return &GRPCClient{}
}

// Connect 连接到gRPC服务器
func (c *GRPCClient) Connect(serverAddr string) error {
	var err error

	// 连接到解析服务
	c.conn, err = grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("DEBUG: gRPC连接失败 - 服务地址: %s, 错误: %v", serverAddr, err)
		return fmt.Errorf("连接解析服务失败: %v", err)
	}

	// 创建客户端实例
	c.client = pb.NewDocumentParserServiceClient(c.conn)

	log.Printf("DEBUG: gRPC客户端连接成功 - 服务地址: %s", serverAddr)
	return nil
}

// Close 关闭连接
func (c *GRPCClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// ParsePDFWithGRPC 通过gRPC调用Python服务解析PDF
func (c *GRPCClient) ParsePDFWithGRPC(filePath string) (string, map[string]interface{}, error) {
	// 获取当前工作目录和绝对路径用于诊断
	currentDir, _ := os.Getwd()
	absPath, _ := filepath.Abs(filePath)
	log.Printf("DEBUG: 通过gRPC调用Python服务解析PDF - Go工作目录: %s", currentDir)
	log.Printf("DEBUG: 通过gRPC调用Python服务解析PDF - 原始路径: %s", filePath)
	log.Printf("DEBUG: 通过gRPC调用Python服务解析PDF - 绝对路径: %s", absPath)

	if c.client == nil {
		return "", nil, fmt.Errorf("gRPC客户端未连接")
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 调用gRPC服务，使用绝对路径
	req := &pb.ParsePDFRequest{FilePath: absPath}
	resp, err := c.client.ParsePDF(ctx, req)
	if err != nil {
		log.Printf("DEBUG: gRPC PDF解析失败 - 原始路径: %s, 绝对路径: %s, 错误: %v", filePath, absPath, err)
		return "", nil, fmt.Errorf("gRPC PDF解析失败: %v", err)
	}

	if !resp.Success {
		return "", nil, fmt.Errorf("PDF解析服务返回错误: %s", resp.ErrorMessage)
	}

	// 转换元数据格式
	metadata := make(map[string]interface{})
	for k, v := range resp.Metadata {
		metadata[k] = v
	}

	// 检查返回的内容是否包含非UTF-8字符
	content := resp.Content
	isValidUTF8 := true
	for i, r := range content {
		if r == utf8.RuneError {
			// 检查是否真的是UTF-8错误
			_, size := utf8.DecodeRuneInString(content[i:])
			if size == 1 {
				isValidUTF8 = false
				break
			}
		}
	}

	// 安全地输出内容预览，避免乱码
	if isValidUTF8 {
		// 如果内容不太长，显示预览
		preview := content
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		log.Printf("DEBUG: gRPC PDF解析完成 - 文件路径: %s, 内容长度: %d", filePath, len(content))
		// 单独输出预览，确保即使预览中有问题也不会影响主要日志
		log.Printf("DEBUG: PDF内容预览: %q", preview)
	} else {
		log.Printf("DEBUG: gRPC PDF解析完成 - 文件路径: %s, 内容长度: %d, 内容包含非UTF-8字符", filePath, len(content))
		// 不输出内容预览，避免乱码
	}
	return content, metadata, nil
}

// ParseDOCXWithGRPC 通过gRPC调用Python服务解析DOCX
func (c *GRPCClient) ParseDOCXWithGRPC(filePath string) (string, map[string]interface{}, error) {
	// 获取当前工作目录和绝对路径用于诊断
	currentDir, _ := os.Getwd()
	absPath, _ := filepath.Abs(filePath)
	log.Printf("DEBUG: 通过gRPC调用Python服务解析DOCX - Go工作目录: %s", currentDir)
	log.Printf("DEBUG: 通过gRPC调用Python服务解析DOCX - 原始路径: %s", filePath)
	log.Printf("DEBUG: 通过gRPC调用Python服务解析DOCX - 绝对路径: %s", absPath)

	if c.client == nil {
		return "", nil, fmt.Errorf("gRPC客户端未连接")
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 调用gRPC服务，使用绝对路径
	req := &pb.ParseDOCXRequest{FilePath: absPath}
	resp, err := c.client.ParseDOCX(ctx, req)
	if err != nil {
		log.Printf("DEBUG: gRPC DOCX解析失败 - 原始路径: %s, 绝对路径: %s, 错误: %v", filePath, absPath, err)
		return "", nil, fmt.Errorf("gRPC DOCX解析失败: %v", err)
	}

	if !resp.Success {
		return "", nil, fmt.Errorf("DOCX解析服务返回错误: %s", resp.ErrorMessage)
	}

	// 转换元数据格式
	metadata := make(map[string]interface{})
	for k, v := range resp.Metadata {
		metadata[k] = v
	}

	// 检查返回的内容是否包含非UTF-8字符
	content := resp.Content
	isValidUTF8 := true
	for i, r := range content {
		if r == utf8.RuneError {
			// 检查是否真的是UTF-8错误
			_, size := utf8.DecodeRuneInString(content[i:])
			if size == 1 {
				isValidUTF8 = false
				break
			}
		}
	}

	// 安全地输出内容预览，避免乱码
	if isValidUTF8 {
		// 如果内容不太长，显示预览
		preview := content
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		log.Printf("DEBUG: gRPC DOCX解析完成 - 文件路径: %s, 内容长度: %d", filePath, len(content))
		// 单独输出预览，确保即使预览中有问题也不会影响主要日志
		log.Printf("DEBUG: DOCX内容预览: %q", preview)
	} else {
		log.Printf("DEBUG: gRPC DOCX解析完成 - 文件路径: %s, 内容长度: %d, 内容包含非UTF-8字符", filePath, len(content))
		// 不输出内容预览，避免乱码
	}
	return content, metadata, nil
}

// HealthCheck 健康检查
func (c *GRPCClient) HealthCheck(service string) (bool, string, error) {
	log.Printf("DEBUG: 执行健康检查 - 服务: %s", service)

	if c.client == nil {
		return false, "gRPC客户端未连接", nil
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 调用健康检查
	req := &pb.HealthCheckRequest{Service: service}
	resp, err := c.client.HealthCheck(ctx, req)
	if err != nil {
		log.Printf("DEBUG: 健康检查失败 - 服务: %s, 错误: %v", service, err)
		return false, fmt.Sprintf("健康检查失败: %v", err), nil
	}

	return resp.Healthy, resp.Message, nil
}
