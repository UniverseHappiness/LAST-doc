package repository

import (
	"testing"
	"time"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB 设置测试数据库
func SetupTestDB(t *testing.T) *gorm.DB {
	// 使用PostgreSQL的内存数据库用于测试
	dsn := "host=localhost user=postgres password=postgres dbname=test_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// 如果连接失败，尝试使用内存中的PostgreSQL
		dsn = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Skipf("无法连接测试数据库，跳过测试: %v", err)
		}
	}

	// 创建测试数据库
	err = db.Exec("CREATE DATABASE IF NOT EXISTS test_db").Error
	if err != nil {
		// 如果创建数据库失败，尝试使用默认数据库
		t.Logf("创建测试数据库失败，使用默认数据库: %v", err)
	} else {
		// 重新连接到新创建的数据库
		dsn = "host=localhost user=postgres password=postgres dbname=test_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Skipf("无法连接到新创建的测试数据库，跳过测试: %v", err)
		}
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.Document{}, &model.DocumentVersion{}, &model.DocumentMetadata{})
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	// 清理测试数据
	db.Exec("DELETE FROM document_versions")
	db.Exec("DELETE FROM document_metadata")
	db.Exec("DELETE FROM documents")

	return db
}

// CreateTestDocument 创建测试文档
func CreateTestDocument() *model.Document {
	return &model.Document{
		Name:        "测试文档",
		Type:        model.DocumentTypeMarkdown,
		Version:     "1.0.0",
		Tags:        []string{"测试", "文档"},
		FilePath:    "/tmp/test.md",
		FileSize:    1024,
		Status:      model.DocumentStatusCompleted,
		Description: "这是一个测试文档",
		Library:     "测试库",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestDocumentMetadata 创建测试文档元数据
func CreateTestDocumentMetadata(documentID string) *model.DocumentMetadata {
	return &model.DocumentMetadata{
		DocumentID: documentID,
		Metadata: map[string]interface{}{
			"author": "测试作者",
			"date":   "2023-01-01",
			"tags":   []string{"测试", "文档"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestDocumentVersion 创建测试文档版本
func CreateTestDocumentVersion(documentID string) *model.DocumentVersion {
	return &model.DocumentVersion{
		DocumentID:  documentID,
		Version:     "1.0.0",
		FilePath:    "/tmp/test_1.0.0.md",
		FileSize:    1024,
		Status:      model.DocumentStatusCompleted,
		Description: "版本1.0.0",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
