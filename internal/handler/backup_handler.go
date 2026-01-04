package handler

import (
	"net/http"
	"time"

	"github.com/UniverseHappiness/LAST-doc/internal/service"
	"github.com/gin-gonic/gin"
)

// BackupHandler 备份管理处理器
type BackupHandler struct {
	backupService service.BackupService
}

// NewBackupHandler 创建备份管理处理器
func NewBackupHandler(backupService service.BackupService) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
	}
}

// CreateBackup 创建备份
// @Summary 创建系统备份
// @Description 创建完整的系统备份（数据库+存储）
// @Tags Backup
// @Accept json
// @Produce json
// @Param type query string false "备份类型 (database, storage, full)" Enums(database, storage, full)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/create [post]
func (h *BackupHandler) CreateBackup(c *gin.Context) {
	ctx := c.Request.Context()

	backupType := c.DefaultQuery("type", "full")
	var backupPath string
	var err error

	switch backupType {
	case "database":
		backupPath, err = h.backupService.BackupDatabase(ctx)
	case "storage":
		backupPath, err = h.backupService.BackupStorage(ctx)
	case "full":
		backupPath, err = h.backupService.CreateFullBackup(ctx)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid backup type. Must be 'database', 'storage', or 'full'",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Backup created successfully",
		"type":       backupType,
		"backupPath": backupPath,
		"createdAt":  time.Now().Format(time.RFC3339),
	})
}

// ListBackups 获取备份列表
// @Summary 获取备份列表
// @Description 获取所有系统备份的列表
// @Tags Backup
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/list [get]
func (h *BackupHandler) ListBackups(c *gin.Context) {
	ctx := c.Request.Context()

	backups, err := h.backupService.GetBackupList(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   len(backups),
		"backups": backups,
	})
}

// RestoreBackup 恢复备份
// @Summary 恢复系统备份
// @Description 从指定备份恢复系统
// @Tags Backup
// @Accept json
// @Produce json
// @Param backupId path string true "备份ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/restore/{backupId} [post]
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	ctx := c.Request.Context()

	backupID := c.Param("backupId")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "backupId is required",
		})
		return
	}

	// 查找备份
	backups, err := h.backupService.GetBackupList(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var targetBackup *service.BackupInfo
	for i := range backups {
		if backups[i].ID == backupID {
			targetBackup = &backups[i]
			break
		}
	}

	if targetBackup == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Backup not found",
		})
		return
	}

	// 根据备份类型执行恢复
	var restoreErr error
	switch targetBackup.Type {
	case "database":
		restoreErr = h.backupService.RestoreDatabase(ctx, targetBackup.Path)
	case "storage", "full":
		restoreErr = h.backupService.RestoreStorage(ctx, targetBackup.Path)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unknown backup type",
		})
		return
	}

	if restoreErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": restoreErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Backup restored successfully",
		"backupId":   backupID,
		"restoredAt": time.Now().Format(time.RFC3339),
	})
}

// DeleteBackup 删除备份
// @Summary 删除备份
// @Description 删除指定的系统备份
// @Tags Backup
// @Accept json
// @Produce json
// @Param backupId path string true "备份ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/backup/{backupId} [delete]
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	ctx := c.Request.Context()

	backupID := c.Param("backupId")
	if backupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "backupId is required",
		})
		return
	}

	if err := h.backupService.DeleteBackup(ctx, backupID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Backup deleted successfully",
		"backupId": backupID,
	})
}
