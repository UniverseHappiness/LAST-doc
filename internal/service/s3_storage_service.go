package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3StorageService S3存储服务实现
type S3StorageService struct {
	client *s3.S3
	config *StorageConfig
}

// NewS3StorageService 创建S3存储服务实例
func NewS3StorageService(config *StorageConfig) (*S3StorageService, error) {
	awsConfig := &aws.Config{
		Region:           aws.String(config.S3Region),
		DisableSSL:       aws.Bool(config.S3DisableSSL),
		S3ForcePathStyle: aws.Bool(true),
	}

	// 如果提供了自定义endpoint（如MinIO）
	if config.S3Endpoint != "" {
		awsConfig.Endpoint = aws.String(config.S3Endpoint)
	}

	// 设置认证
	if config.S3AccessKey != "" && config.S3SecretKey != "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(
			config.S3AccessKey,
			config.S3SecretKey,
			"",
		)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return &S3StorageService{
		client: s3.New(sess),
		config: config,
	}, nil
}

// SaveFile 保存文件
func (s *S3StorageService) SaveFile(ctx context.Context, file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer src.Close()

	return s.SaveFileStream(ctx, path, src)
}

// SaveFileStream 保存文件流
func (s *S3StorageService) SaveFileStream(ctx context.Context, path string, reader io.Reader) error {
	_, err := s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.S3Bucket),
		Key:    aws.String(path),
		Body:   aws.ReadSeekCloser(reader),
	})

	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return nil
}

// DeleteFile 删除文件
func (s *S3StorageService) DeleteFile(ctx context.Context, path string) error {
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.S3Bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %v", err)
	}

	return nil
}

// GetFile 获取文件内容
func (s *S3StorageService) GetFile(ctx context.Context, path string) ([]byte, error) {
	stream, err := s.GetFileStream(ctx, path)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	return io.ReadAll(stream)
}

// GetFileStream 获取文件流
func (s *S3StorageService) GetFileStream(ctx context.Context, path string) (io.ReadCloser, error) {
	result, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.S3Bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get file from S3: %v", err)
	}

	return result.Body, nil
}

// GenerateFilePath 生成文件存储路径
func (s *S3StorageService) GenerateFilePath(documentID, fileName string) string {
	// 清理文件名，移除特殊字符
	cleanFileName := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, fileName)

	return filepath.Join(documentID, cleanFileName)
}

// FileExists 检查文件是否存在
func (s *S3StorageService) FileExists(ctx context.Context, path string) bool {
	_, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.S3Bucket),
		Key:    aws.String(path),
	})

	return err == nil
}

// CopyFile 复制文件
func (s *S3StorageService) CopyFile(ctx context.Context, srcPath, dstPath string) error {
	_, err := s.client.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.config.S3Bucket),
		CopySource: aws.String(s.config.S3Bucket + "/" + srcPath),
		Key:        aws.String(dstPath),
	})

	if err != nil {
		return fmt.Errorf("failed to copy file in S3: %v", err)
	}

	return nil
}

// MoveFile 移动文件
func (s *S3StorageService) MoveFile(ctx context.Context, srcPath, dstPath string) error {
	// S3中移动文件实际上是复制后删除
	if err := s.CopyFile(ctx, srcPath, dstPath); err != nil {
		return err
	}

	return s.DeleteFile(ctx, srcPath)
}

// GetFileSize 获取文件大小
func (s *S3StorageService) GetFileSize(ctx context.Context, path string) (int64, error) {
	result, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.S3Bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get file size from S3: %v", err)
	}

	return *result.ContentLength, nil
}

// HealthCheck 健康检查
func (s *S3StorageService) HealthCheck(ctx context.Context) error {
	// 通过列出存储桶中的对象来验证连接
	_, err := s.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.config.S3Bucket),
		MaxKeys: aws.Int64(1), // 只获取一个对象，避免过多数据传输
	})

	if err != nil {
		return fmt.Errorf("S3/MinIO connection failed: %v", err)
	}

	return nil
}
