package main

import (
	// "context"
	"fmt"
	"log"
	"os"

	// "path/filepath"

	// "github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/UniverseHappiness/LAST-doc/internal/handler"
	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/repository"
	"github.com/UniverseHappiness/LAST-doc/internal/router"
	"github.com/UniverseHappiness/LAST-doc/internal/service"
)

func main() {
	// 从环境变量获取配置
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "ai_doc_library")
	serverPort := getEnv("SERVER_PORT", "8080")
	baseStorageDir := getEnv("STORAGE_DIR", "./storage")

	// 构建数据库连接字符串
	dsn := buildDSN(dbHost, dbPort, dbUser, dbPassword, dbName)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	err = autoMigrate(db)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建存储目录
	if err := os.MkdirAll(baseStorageDir, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	// 初始化仓库
	documentRepo := repository.NewDocumentRepository(db)
	versionRepo := repository.NewDocumentVersionRepository(db)
	metadataRepo := repository.NewDocumentMetadataRepository(db)

	// 初始化服务
	storageService := service.NewLocalStorageService(baseStorageDir)
	parserService := service.NewParserService()
	documentService := service.NewDocumentService(
		documentRepo,
		versionRepo,
		metadataRepo,
		storageService,
		parserService,
		baseStorageDir,
	)

	// 初始化处理器
	documentHandler := handler.NewDocumentHandler(documentService)

	// 初始化路由器
	router := router.NewRouter(documentHandler)
	r := router.SetupRoutes()

	// 启动服务器
	log.Printf("Server starting on port %s", serverPort)
	log.Printf("Storage directory: %s", baseStorageDir)
	if err := r.Run(":" + serverPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// buildDSN 构建数据库连接字符串
func buildDSN(host, port, user, password, dbname string) string {
	return "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=Asia/Shanghai"
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	// 首先执行自动迁移
	err := db.AutoMigrate(
		&model.Document{},
		&model.DocumentVersion{},
		&model.DocumentMetadata{},
	)
	if err != nil {
		return err
	}

	// 为document_versions表添加复合唯一约束
	// 确保同一文档的版本号唯一
	if err := addUniqueConstraint(db); err != nil {
		return err
	}

	return nil
}

// addUniqueConstraint 添加数据库唯一约束
func addUniqueConstraint(db *gorm.DB) error {
	// 检查约束是否已存在
	var count int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.table_constraints
		WHERE table_name = 'document_versions'
		AND constraint_name = 'idx_document_version_unique'
	`).Scan(&count).Error; err != nil {
		return err
	}

	// 如果约束不存在，则创建
	if count == 0 {
		log.Println("准备添加document_versions表的复合唯一约束...")

		// 首先清理重复数据
		if err := cleanDuplicateVersions(db); err != nil {
			return fmt.Errorf("failed to clean duplicate versions: %v", err)
		}

		// 然后添加唯一约束
		if err := db.Exec(`
			ALTER TABLE document_versions
			ADD CONSTRAINT idx_document_version_unique
			UNIQUE (document_id, version)
		`).Error; err != nil {
			return fmt.Errorf("failed to add unique constraint: %v", err)
		}
		log.Println("成功添加document_versions表的复合唯一约束 (document_id, version)")
	}

	return nil
}

// cleanDuplicateVersions 清理重复的文档版本数据
func cleanDuplicateVersions(db *gorm.DB) error {
	log.Println("开始清理重复的文档版本数据...")

	// 查看重复数据
	var duplicates []struct {
		DocumentID     string
		Version        string
		DuplicateCount int64
	}

	if err := db.Raw(`
		SELECT document_id, version, COUNT(*) as duplicate_count
		FROM document_versions
		GROUP BY document_id, version
		HAVING COUNT(*) > 1
	`).Scan(&duplicates).Error; err != nil {
		return err
	}

	if len(duplicates) == 0 {
		log.Println("没有发现重复的版本数据")
		return nil
	}

	log.Printf("发现 %d 组重复的版本数据", len(duplicates))

	// 创建临时表存储需要保留的版本ID（每个document_id+version组合保留最新的一个）
	if err := db.Exec(`
		CREATE TEMPORARY TABLE versions_to_keep AS
		WITH ranked_versions AS (
			SELECT
				id,
				document_id,
				version,
				ROW_NUMBER() OVER (PARTITION BY document_id, version ORDER BY created_at DESC) as rn
			FROM document_versions
		)
		SELECT id FROM ranked_versions WHERE rn = 1
	`).Error; err != nil {
		return err
	}

	// 删除重复的版本（保留最新的）
	result := db.Exec(`
		DELETE FROM document_versions
		WHERE id NOT IN (SELECT id FROM versions_to_keep)
	`)

	if result.Error != nil {
		// 删除临时表
		db.Exec("DROP TABLE IF EXISTS versions_to_keep")
		return result.Error
	}

	deletedCount := result.RowsAffected
	log.Printf("删除了 %d 个重复的版本记录", deletedCount)

	// 删除临时表
	if err := db.Exec("DROP TABLE versions_to_keep").Error; err != nil {
		return err
	}

	// 验证清理结果
	var remainingDuplicates int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM (
			SELECT document_id, version
			FROM document_versions
			GROUP BY document_id, version
			HAVING COUNT(*) > 1
		) as remaining
	`).Scan(&remainingDuplicates).Error; err != nil {
		return err
	}

	if remainingDuplicates > 0 {
		return fmt.Errorf("清理后仍有 %d 组重复数据", remainingDuplicates)
	}

	log.Println("重复数据清理完成")
	return nil
}
