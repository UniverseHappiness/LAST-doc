package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupService 备份服务接口
type BackupService interface {
	BackupDatabase(ctx context.Context) (string, error)
	BackupStorage(ctx context.Context) (string, error)
	CreateFullBackup(ctx context.Context) (string, error)
	RestoreDatabase(ctx context.Context, backupFile string) error
	RestoreStorage(ctx context.Context, backupDir string) error
	GetBackupList(ctx context.Context) ([]BackupInfo, error)
	DeleteBackup(ctx context.Context, backupID string) error
}

// BackupInfo 备份信息
type BackupInfo struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // database, storage, full
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

// BackupServiceImpl 备份服务实现
type BackupServiceImpl struct {
	db             DatabaseBackup
	storageService StorageService
	backupDir      string
}

// DatabaseBackup 数据库备份接口
type DatabaseBackup interface {
	DumpDatabase(ctx context.Context, destFile string) error
	RestoreDatabase(ctx context.Context, sourceFile string) error
}

// NewBackupService 创建备份服务
func NewBackupService(db DatabaseBackup, storageService StorageService, backupDir string) *BackupServiceImpl {
	// 确保备份目录存在
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create backup directory: %v", err))
	}
	return &BackupServiceImpl{
		db:             db,
		storageService: storageService,
		backupDir:      backupDir,
	}
}

// BackupDatabase 备份数据库
func (s *BackupServiceImpl) BackupDatabase(ctx context.Context) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(s.backupDir, fmt.Sprintf("database_%s.sql", timestamp))

	if err := s.db.DumpDatabase(ctx, backupFile); err != nil {
		return "", fmt.Errorf("failed to backup database: %v", err)
	}

	// 验证备份文件
	if _, err := os.Stat(backupFile); err != nil {
		return "", fmt.Errorf("failed to verify backup file: %v", err)
	}

	return backupFile, nil
}

// BackupStorage 备份存储
func (s *BackupServiceImpl) BackupStorage(ctx context.Context) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(s.backupDir, fmt.Sprintf("storage_%s", timestamp))

	// 创建备份目录
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %v", err)
	}

	return backupDir, nil
}

// CreateFullBackup 创建完整备份
func (s *BackupServiceImpl) CreateFullBackup(ctx context.Context) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	fullBackupDir := filepath.Join(s.backupDir, fmt.Sprintf("full_%s", timestamp))

	if err := os.MkdirAll(fullBackupDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create full backup directory: %v", err)
	}

	// 备份数据库
	dbBackupFile := filepath.Join(fullBackupDir, "database.sql")
	if err := s.db.DumpDatabase(ctx, dbBackupFile); err != nil {
		return "", fmt.Errorf("failed to backup database in full backup: %v", err)
	}

	// 创建存储备份标记
	storageBackupFile := filepath.Join(fullBackupDir, "storage_backup_info.txt")
	infoContent := fmt.Sprintf("Backup created at: %s\nStorage type: %s\n", time.Now().Format(time.RFC3339), "local")
	if err := os.WriteFile(storageBackupFile, []byte(infoContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create storage backup info: %v", err)
	}

	return fullBackupDir, nil
}

// RestoreDatabase 恢复数据库
func (s *BackupServiceImpl) RestoreDatabase(ctx context.Context, backupFile string) error {
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupFile)
	}

	if err := s.db.RestoreDatabase(ctx, backupFile); err != nil {
		return fmt.Errorf("failed to restore database: %v", err)
	}

	return nil
}

// RestoreStorage 恢复存储
func (s *BackupServiceImpl) RestoreStorage(ctx context.Context, backupDir string) error {
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return fmt.Errorf("backup directory does not exist: %s", backupDir)
	}

	// 这里可以实现实际的存储恢复逻辑
	// 对于本地存储，可以复制文件回原目录
	// 对于S3/MinIO，可以实现批量上传

	return nil
}

// GetBackupList 获取备份列表
func (s *BackupServiceImpl) GetBackupList(ctx context.Context) ([]BackupInfo, error) {
	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %v", err)
	}

	var backups []BackupInfo
	for _, entry := range entries {
		if entry.IsDir() {
			info, _ := entry.Info()
			backups = append(backups, BackupInfo{
				ID:        entry.Name(),
				Type:      "full",
				Path:      filepath.Join(s.backupDir, entry.Name()),
				Size:      info.Size(),
				CreatedAt: info.ModTime(),
			})
		} else {
			info, _ := entry.Info()
			backupType := "database"
			if filepath.Ext(entry.Name()) != ".sql" {
				backupType = "storage"
			}
			backups = append(backups, BackupInfo{
				ID:        entry.Name(),
				Type:      backupType,
				Path:      filepath.Join(s.backupDir, entry.Name()),
				Size:      info.Size(),
				CreatedAt: info.ModTime(),
			})
		}
	}

	return backups, nil
}

// DeleteBackup 删除备份
func (s *BackupServiceImpl) DeleteBackup(ctx context.Context, backupID string) error {
	backupPath := filepath.Join(s.backupDir, backupID)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup does not exist: %s", backupID)
	}

	// 如果是目录，删除整个目录
	if info, err := os.Stat(backupPath); err == nil && info.IsDir() {
		if err := os.RemoveAll(backupPath); err != nil {
			return fmt.Errorf("failed to delete backup directory: %v", err)
		}
	} else {
		// 删除文件
		if err := os.Remove(backupPath); err != nil {
			return fmt.Errorf("failed to delete backup file: %v", err)
		}
	}

	return nil
}
