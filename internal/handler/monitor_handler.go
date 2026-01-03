package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/service"
)

// MonitorHandler 监控处理器
type MonitorHandler struct {
	monitorService service.MonitorService
}

// NewMonitorHandler 创建监控处理器实例
func NewMonitorHandler(monitorService service.MonitorService) *MonitorHandler {
	return &MonitorHandler{
		monitorService: monitorService,
	}
}

// GetCurrentMetrics 获取当前系统指标
func (h *MonitorHandler) GetCurrentMetrics(c *gin.Context) {
	metrics, err := h.monitorService.GetCurrentMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取系统指标失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}

// GetMetricsHistory 获取指标历史数据
func (h *MonitorHandler) GetMetricsHistory(c *gin.Context) {
	// 解析时间参数
	startTimeStr := c.DefaultQuery("start_time", time.Now().Add(-24*time.Hour).Format(time.RFC3339))
	endTimeStr := c.DefaultQuery("end_time", time.Now().Format(time.RFC3339))

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的开始时间",
			"details": err.Error(),
		})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的结束时间",
			"details": err.Error(),
		})
		return
	}

	// 获取历史数据
	metrics, err := h.monitorService.GetMetricsHistory(c.Request.Context(), startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取指标历史失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
		"count":   len(metrics),
	})
}

// GetMetricsReport 获取指标报告
func (h *MonitorHandler) GetMetricsReport(c *gin.Context) {
	// 解析持续时间参数
	durationStr := c.DefaultQuery("duration", "1h")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的持续时间",
			"details": err.Error(),
		})
		return
	}

	// 获取报告
	report, err := h.monitorService.GetMetricsReport(c.Request.Context(), duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取指标报告失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetSystemStatus 获取系统状态
func (h *MonitorHandler) GetSystemStatus(c *gin.Context) {
	status, err := h.monitorService.GetSystemStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取系统状态失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetLogs 获取日志列表
func (h *MonitorHandler) GetLogs(c *gin.Context) {
	// 解析过滤条件
	filter := &model.LogFilter{
		Level:     c.Query("level"),
		Service:   c.Query("service"),
		UserID:    c.Query("user_id"),
		Message:   c.Query("message"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	filter.Page = page
	filter.Size = size

	// 获取日志
	response, err := h.monitorService.GetLogs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取日志失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetLogStats 获取日志统计
func (h *MonitorHandler) GetLogStats(c *gin.Context) {
	// 解析时间参数
	startTimeStr := c.DefaultQuery("start_time", time.Now().Add(-24*time.Hour).Format(time.RFC3339))
	endTimeStr := c.DefaultQuery("end_time", time.Now().Format(time.RFC3339))

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的开始时间",
			"details": err.Error(),
		})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的结束时间",
			"details": err.Error(),
		})
		return
	}

	// 获取统计
	stats, err := h.monitorService.GetLogStats(c.Request.Context(), startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取日志统计失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// GetPerformanceReport 获取性能报告
func (h *MonitorHandler) GetPerformanceReport(c *gin.Context) {
	// 解析时间参数
	startTimeStr := c.DefaultQuery("start_time", time.Now().Add(-24*time.Hour).Format(time.RFC3339))
	endTimeStr := c.DefaultQuery("end_time", time.Now().Format(time.RFC3339))

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的开始时间",
			"details": err.Error(),
		})
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的结束时间",
			"details": err.Error(),
		})
		return
	}

	// 获取性能报告
	report, err := h.monitorService.GetPerformanceReport(c.Request.Context(), startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取性能报告失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, report)
}

// CleanupOldData 清理旧数据
func (h *MonitorHandler) CleanupOldData(c *gin.Context) {
	// 解析保留天数参数
	retentionDaysStr := c.DefaultQuery("retention_days", "30")
	retentionDays, err := strconv.Atoi(retentionDaysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的保留天数",
			"details": err.Error(),
		})
		return
	}

	// 清理旧指标
	err = h.monitorService.CleanupOldMetrics(c.Request.Context(), retentionDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "清理旧指标失败",
			"details": err.Error(),
		})
		return
	}

	// 清理旧日志
	err = h.monitorService.CleanupOldLogs(c.Request.Context(), retentionDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "清理旧日志失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "旧数据清理成功",
		"retention_days": retentionDays,
	})
}
