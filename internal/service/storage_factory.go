package service

import (
	"fmt"
	"os"
)

// NewStorageService 根据配置创建存储服务实例
func NewStorageService(config *StorageConfig) (StorageService, error) {
	switch config.Type {
	case StorageTypeLocal:
		return NewLocalStorageService(config.LocalDir), nil
	case StorageTypeS3:
		return NewS3StorageService(config)
	case StorageTypeMinIO:
		return NewS3StorageService(config) // MinIO使用S3兼容模式
	default:
		return NewLocalStorageService(config.LocalDir), nil
	}
}

// NewStorageServiceFromEnv 从环境变量创建存储服务
func NewStorageServiceFromEnv() (StorageService, error) {
	storageTypeStr := getEnv("STORAGE_TYPE", "local")
	storageType := StorageType(storageTypeStr)

	config := &StorageConfig{
		Type:           storageType,
		LocalDir:       getEnv("STORAGE_DIR", "./storage"),
		S3Region:       getEnv("S3_REGION", "us-east-1"),
		S3Bucket:       getEnv("S3_BUCKET", ""),
		S3AccessKey:    getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:    getEnv("S3_SECRET_KEY", ""),
		S3Endpoint:     getEnv("S3_ENDPOINT", ""),
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", ""),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", ""),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", ""),
		MinIOBucket:    getEnv("MINIO_BUCKET", ""),
		MinIOLocation:  getEnv("MINIO_LOCATION", "us-east-1"),
	}

	// 根据存储类型设置SSL配置
	if storageType == StorageTypeS3 || storageType == StorageTypeMinIO {
		config.S3DisableSSL = getEnv("S3_DISABLE_SSL", "false") == "true"
		if storageType == StorageTypeMinIO {
			config.MinIOUseSSL = getEnv("MINIO_USE_SSL", "false") == "true"
			// 对于MinIO，使用S3兼容的配置
			config.S3Endpoint = config.MinIOEndpoint
			config.S3AccessKey = config.MinIOAccessKey
			config.S3SecretKey = config.MinIOSecretKey
			config.S3Bucket = config.MinIOBucket
			config.S3Region = config.MinIOLocation
		}
	}

	return NewStorageService(config)
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// MustNewStorageServiceFromEnv 从环境变量创建存储服务，如果失败则panic
func MustNewStorageServiceFromEnv() StorageService {
	service, err := NewStorageServiceFromEnv()
	if err != nil {
		panic(fmt.Sprintf("Failed to create storage service: %v", err))
	}
	return service
}
