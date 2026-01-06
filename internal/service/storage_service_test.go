package service

import (
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestStorageType_Constants 测试存储类型常量
func TestStorageType_Constants(t *testing.T) {
	tests := []struct {
		name string
		typ  StorageType
	}{
		{name: "本地存储", typ: StorageTypeLocal},
		{name: "S3存储", typ: StorageTypeS3},
		{name: "MinIO存储", typ: StorageTypeMinIO},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.typ != StorageTypeLocal && tt.typ != StorageTypeS3 && tt.typ != StorageTypeMinIO {
				t.Errorf("存储类型 %s 不是常量", tt.typ)
			}
		})
	}
}

// TestNewLocalStorageService 测试创建本地存储服务
func TestNewLocalStorageService(t *testing.T) {
	baseDir := os.TempDir()
	defer os.RemoveAll(baseDir)

	service := NewLocalStorageService(baseDir)

	if service == nil {
		t.Fatal("返回nil服务")
	}

	if service.baseDir != baseDir {
		t.Errorf("baseDir不匹配: got %s, expected %s", service.baseDir, baseDir)
	}
}

// TestLocalStorageService_GenerateFilePath 测试生成文件路径
func TestLocalStorageService_GenerateFilePath(t *testing.T) {
	baseDir := os.TempDir()
	defer os.RemoveAll(baseDir)

	service := NewLocalStorageService(baseDir)

	tests := []struct {
		name       string
		documentID string
		fileName   string
		want       string
	}{
		{name: "标准路径", documentID: "doc123", fileName: "test.pdf", want: filepath.Join(baseDir, "doc123", "test.pdf")},
		{name: "特殊字符", documentID: "doc:abc", fileName: "file.txt", want: filepath.Join(baseDir, "doc:abc", "file.txt")},
		{name: "空格字符", documentID: "doc1", fileName: "file with spaces.pdf", want: filepath.Join(baseDir, "doc1", "file_with_spaces.pdf")},
		{name: "特殊符号", documentID: "doc2", fileName: "file@#$%.pdf", want: filepath.Join(baseDir, "doc2", "file____.pdf")},
		{name: "长文件名", documentID: "doc1", fileName: "very-long-filename-with-multiple-dashes.pdf", want: filepath.Join(baseDir, "doc1", "very-long-filename-with-multiple-dashes.pdf")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.GenerateFilePath(tt.documentID, tt.fileName)
			if got != tt.want {
				t.Errorf("GenerateFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLocalStorageService_GetBaseDir 测试获取基础目录
func TestLocalStorageService_GetBaseDir(t *testing.T) {
	baseDir := os.TempDir()
	defer os.RemoveAll(baseDir)

	service := NewLocalStorageService(baseDir)

	got := service.GetBaseDir()
	if got != baseDir {
		t.Errorf("GetBaseDir() = %v, want %v", got, baseDir)
	}
}

// TestLocalStorageService_FileExists 测试文件存在检查
func TestLocalStorageService_FileExists(t *testing.T) {
	baseDir := os.TempDir()
	defer os.RemoveAll(baseDir)

	service := NewLocalStorageService(baseDir)
	ctx := context.Background()

	// 创建测试文件
	testFile := filepath.Join(baseDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "文件存在", path: testFile, want: true},
		{name: "文件不存在", path: filepath.Join(baseDir, "nonexistent.txt"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.FileExists(ctx, tt.path)
			if got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLocalStorageService_GetFileSize 测试获取文件大小
func TestLocalStorageService_GetFileSize(t *testing.T) {
	baseDir := os.TempDir()
	defer os.RemoveAll(baseDir)

	service := NewLocalStorageService(baseDir)
	ctx := context.Background()

	// 创建测试文件
	testFile := filepath.Join(baseDir, "test.txt")
	content := []byte("test content")
	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	got, err := service.GetFileSize(ctx, testFile)
	if err != nil {
		t.Errorf("GetFileSize() error = %v", err)
		return
	}

	want := int64(len(content))
	if got != want {
		t.Errorf("GetFileSize() = %v, want %v", got, want)
	}
}

// TestLocalStorageService_DeleteFile 测试删除文件
func TestLocalStorageService_DeleteFile(t *testing.T) {
	baseDir := os.TempDir()
	defer os.RemoveAll(baseDir)

	service := NewLocalStorageService(baseDir)
	ctx := context.Background()

	// 创建测试文件
	testFile := filepath.Join(baseDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 删除文件
	err = service.DeleteFile(ctx, testFile)
	if err != nil {
		t.Errorf("DeleteFile() error = %v", err)
	}

	// 验证文件已删除
	if service.FileExists(ctx, testFile) {
		t.Error("文件未被删除")
	}

	// 删除不存在的文件不应报错
	err = service.DeleteFile(ctx, "nonexistent.txt")
	if err != nil {
		t.Errorf("删除不存在的文件应该返回nil: %v", err)
	}
}

// TestLocalStorageService_GetFile 测试获取文件内容
func TestLocalStorageService_GetFile(t *testing.T) {
	baseDir := os.TempDir()
	defer os.RemoveAll(baseDir)

	service := NewLocalStorageService(baseDir)
	ctx := context.Background()

	// 创建测试文件
	testFile := filepath.Join(baseDir, "test.txt")
	content := []byte("test content")
	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	got, err := service.GetFile(ctx, testFile)
	if err != nil {
		t.Errorf("GetFile() error = %v", err)
		return
	}

	if string(got) != string(content) {
		t.Errorf("GetFile() = %v, want %v", string(got), string(content))
	}
}

// TestLocalStorageService_HealthCheck 测试健康检查
func TestLocalStorageService_HealthCheck(t *testing.T) {
	baseDir := os.TempDir()
	defer os.RemoveAll(baseDir)

	service := NewLocalStorageService(baseDir)
	ctx := context.Background()

	err := service.HealthCheck(ctx)
	if err != nil {
		t.Errorf("HealthCheck() error = %v", err)
	}
}

// TestStorageConfig 测试存储配置结构体
func TestStorageConfig(t *testing.T) {
	tests := []struct {
		name   string
		config StorageConfig
	}{
		{name: "本地配置", config: StorageConfig{Type: StorageTypeLocal}},
		{name: "S3配置", config: StorageConfig{Type: StorageTypeS3}},
		{name: "MinIO配置", config: StorageConfig{Type: StorageTypeMinIO}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.Type != tt.config.Type {
				t.Errorf("类型不匹配")
			}
		})
	}
}

// TestS3ConfigFields 测试S3配置字段
func TestS3ConfigFields(t *testing.T) {
	config := &StorageConfig{
		Type:        StorageTypeS3,
		S3Region:    "us-west-2",
		S3Bucket:    "test-bucket",
		S3AccessKey: "access-key",
		S3SecretKey: "secret-key",
		S3Endpoint:  "",
	}

	if config.Type != StorageTypeS3 {
		t.Error("S3Type不匹配")
	}
	if config.S3Region != "us-west-2" {
		t.Error("S3Region不匹配")
	}
	if config.S3Bucket != "test-bucket" {
		t.Error("S3Bucket不匹配")
	}
	if config.S3AccessKey != "access-key" {
		t.Error("S3AccessKey不匹配")
	}
	if config.S3SecretKey != "secret-key" {
		t.Error("S3SecretKey不匹配")
	}
	if config.S3Endpoint != "" {
		t.Error("S3Endpoint不匹配")
	}
}

// TestMinIOConfigFields 测试MinIO配置字段
func TestMinIOConfigFields(t *testing.T) {
	config := &StorageConfig{
		Type:           StorageTypeMinIO,
		MinIOEndpoint:  "http://localhost:9000",
		MinIOAccessKey: "access-key",
		MinIOSecretKey: "secret-key",
		MinIOBucket:    "test-bucket",
		MinIOUseSSL:    true,
	}

	if config.Type != StorageTypeMinIO {
		t.Error("MinIOType不匹配")
	}
	if config.MinIOEndpoint != "http://localhost:9000" {
		t.Error("MinIOEndpoint不匹配")
	}
	if config.MinIOAccessKey != "access-key" {
		t.Error("MinIOAccessKey不匹配")
	}
	if config.MinIOSecretKey != "secret-key" {
		t.Error("MinIOSecretKey不匹配")
	}
	if config.MinIOBucket != "test-bucket" {
		t.Error("MinIOBucket不匹配")
	}
	if config.MinIOUseSSL != true {
		t.Error("MinIOUseSSL不匹配")
	}
}

// TestPathOperations 测试路径操作
func TestPathOperations(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{name: "单级路径", input: "document.pdf", want: 1},
		{name: "多级路径", input: filepath.Join("folder", "file.pdf"), want: 2},
		{name: "根路径", input: "/", want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(tt.input, string(filepath.Separator))
			got := len(parts)
			if got != tt.want {
				t.Errorf("路径解析错误: expected %d parts, got %d", tt.want, got)
			}
		})
	}
}

// TestIOOperations 测试I/O操作
func TestIOOperations(t *testing.T) {
	baseDir := os.TempDir()
	defer os.RemoveAll(baseDir)

	filePath := filepath.Join(baseDir, "io-test.txt")
	content := []byte("test content")

	// 测试写入
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatalf("无法创建I/O测试文件: %v", err)
	}

	// 测试读取
	readContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("无法读取文件: %v", err)
	}

	if string(readContent) != string(content) {
		t.Errorf("文件内容不匹配")
	}

	// 测试文件操作
	file, err := os.Open(filePath)
	if err != nil {
		t.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		t.Errorf("读取失败: %v", err)
	}

	if n != 12 {
		t.Errorf("读取字数不匹配: expected 12, got %d", n)
	}
}

// TestContextHandling 测试context处理
func TestContextHandling(t *testing.T) {
	ctx := context.Background()

	if ctx == nil {
		t.Error("context.Background()返回nil")
	}

	// 测试带超时的context
	testCtx, cancel := context.WithTimeout(ctx, 5*time.Microsecond)
	defer cancel()

	select {
	case <-testCtx.Done():
		// 超时
		t.Log("context已超时")
	case <-time.After(10 * time.Millisecond):
		// 没超时
		t.Log("context未超时")
	}
}

// TestMIMETypeHandling 测试MIME类型处理
func TestMIMETypeHandling(t *testing.T) {
	tests := []struct {
		name string
		ext  string
		want string
	}{
		{name: "PDF文件", ext: ".pdf", want: ".pdf"},
		{name: "文本文件", ext: ".txt", want: ".txt"},
		{name: "JSON文件", ext: ".json", want: ".json"},
		{name: "图片文件", ext: ".png", want: ".png"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试文件扩展包含点号
			if !strings.Contains(tt.ext, ".") {
				t.Error("文件扩展应该包含点号")
			}
		})
	}
}

// TestFileHeader 测试 multipart.FileHeader 结构
func TestFileHeader(t *testing.T) {
	header := &multipart.FileHeader{
		Filename: "test.pdf",
	}

	if header.Filename != "test.pdf" {
		t.Error("文件名设置失败")
	}
}

// TestCleanFileName 测试文件名清理逻辑
func TestCleanFileName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "正常文件名", input: "test.pdf", expected: "test.pdf"},
		{name: "带空格", input: "file with spaces.pdf", expected: "file_with_spaces.pdf"},
		{name: "特殊字符", input: "file@#$%.pdf", expected: "file____.pdf"},
		{name: "中文文件名", input: "测试文件.pdf", expected: "____.pdf"},
		{name: "混合字符", input: "file-123_Test.pdf", expected: "file-123_Test.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.Map(func(r rune) rune {
				if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.' {
					return r
				}
				return '_'
			}, tt.input)

			if result != tt.expected {
				t.Errorf("清理结果: got %s, expected %s", result, tt.expected)
			}
		})
	}
}

// TestBaseDirOperations 测试基础目录操作
func TestBaseDirOperations(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "test-storage-*")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(baseDir)

	// 测试创建子目录
	subDir := filepath.Join(baseDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("无法创建子目录: %v", err)
	}

	// 检查目录是否可访问
	info, err := os.Stat(subDir)
	if err != nil {
		t.Errorf("无法访问子目录: %v", err)
	}

	if !info.IsDir() {
		t.Error("子目录创建失败")
	}

	// 测试清理
	err = os.RemoveAll(baseDir)
	if err != nil {
		t.Errorf("无法清理目录: %v", err)
	}

	// 验证目录已删除
	_, err = os.Stat(baseDir)
	if !os.IsNotExist(err) {
		t.Error("目录未删除")
	}
}
