package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// StorageService 存储服务接口
type StorageService interface {
	SaveFile(ctx context.Context, file *multipart.FileHeader, path string) error
	DeleteFile(ctx context.Context, path string) error
	GetFile(ctx context.Context, path string) ([]byte, error)
	GenerateFilePath(documentID, fileName string) string
	FileExists(ctx context.Context, path string) bool
	CopyFile(ctx context.Context, srcPath, dstPath string) error
	MoveFile(ctx context.Context, srcPath, dstPath string) error
	GetFileSize(ctx context.Context, path string) (int64, error)
	GetFileStream(ctx context.Context, path string) (io.ReadCloser, error)
	SaveFileStream(ctx context.Context, path string, reader io.Reader) error
	HealthCheck(ctx context.Context) error
}

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeS3    StorageType = "s3"
	StorageTypeMinIO StorageType = "minio"
)

// StorageConfig 存储配置
type StorageConfig struct {
	Type StorageType

	// 本地存储配置
	LocalDir string

	// S3配置
	S3Region     string
	S3Bucket     string
	S3AccessKey  string
	S3SecretKey  string
	S3Endpoint   string
	S3DisableSSL bool

	// MinIO配置
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIOLocation  string
	MinIOUseSSL    bool
}

// LocalStorageService 本地存储服务实现
type LocalStorageService struct {
	baseDir string
}

// NewLocalStorageService 创建本地存储服务实例
func NewLocalStorageService(baseDir string) *LocalStorageService {
	// 确保基础目录存在
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create base storage directory: %v", err))
	}
	return &LocalStorageService{
		baseDir: baseDir,
	}
}

// SaveFile 保存文件
func (s *LocalStorageService) SaveFile(ctx context.Context, file *multipart.FileHeader, path string) error {
	// 确保目标目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// 打开源文件
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	return nil
}

// DeleteFile 删除文件
func (s *LocalStorageService) DeleteFile(ctx context.Context, path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

// GetFile 获取文件内容
func (s *LocalStorageService) GetFile(ctx context.Context, path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return data, nil
}

// GenerateFilePath 生成文件存储路径
func (s *LocalStorageService) GenerateFilePath(documentID, fileName string) string {
	// 清理文件名，移除特殊字符
	cleanFileName := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, fileName)

	return filepath.Join(s.baseDir, documentID, cleanFileName)
}

// GetBaseDir 获取基础存储目录
func (s *LocalStorageService) GetBaseDir() string {
	return s.baseDir
}

// FileExists 检查文件是否存在
func (s *LocalStorageService) FileExists(ctx context.Context, path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CopyFile 复制文件
func (s *LocalStorageService) CopyFile(ctx context.Context, srcPath, dstPath string) error {
	// 确保目标目录存在
	dir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// 打开源文件
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	return nil
}

// MoveFile 移动文件
func (s *LocalStorageService) MoveFile(ctx context.Context, srcPath, dstPath string) error {
	// 确保目标目录存在
	dir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// 移动文件
	if err := os.Rename(srcPath, dstPath); err != nil {
		return fmt.Errorf("failed to move file: %v", err)
	}

	return nil
}

// GetFileSize 获取文件大小
func (s *LocalStorageService) GetFileSize(ctx context.Context, path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %v", err)
	}
	return info.Size(), nil
}

// GetFileStream 获取文件流
func (s *LocalStorageService) GetFileStream(ctx context.Context, path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	return file, nil
}

// SaveFileStream 保存文件流
func (s *LocalStorageService) SaveFileStream(ctx context.Context, path string, reader io.Reader) error {
	// 确保目标目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// 创建目标文件
	dst, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// 复制流内容
	if _, err = io.Copy(dst, reader); err != nil {
		return fmt.Errorf("failed to copy stream content: %v", err)
	}

	return nil
}

// HealthCheck 健康检查
func (s *LocalStorageService) HealthCheck(ctx context.Context) error {
	// 检查基础目录是否可访问
	if _, err := os.Stat(s.baseDir); err != nil {
		return fmt.Errorf("base directory not accessible: %v", err)
	}

	// 尝试创建临时文件进行写入测试
	testFile := filepath.Join(s.baseDir, ".health_check_test")
	if err := os.WriteFile(testFile, []byte("health check"), 0644); err != nil {
		return fmt.Errorf("failed to write test file: %v", err)
	}
	defer os.Remove(testFile)

	return nil
}
