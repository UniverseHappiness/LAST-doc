package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
)

// MetricsRepository 指标仓库接口
type MetricsRepository interface {
	// CreateMetrics 创建指标记录
	CreateMetrics(ctx context.Context, metrics *model.SystemMetrics) error

	// GetLatestMetrics 获取最新的指标
	GetLatestMetrics(ctx context.Context) (*model.SystemMetrics, error)

	// GetMetricsByTimeRange 获取指定时间范围内的指标
	GetMetricsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]model.SystemMetrics, error)

	// GetAverageMetrics 获取平均指标
	GetAverageMetrics(ctx context.Context, startTime, endTime time.Time) (*model.SystemMetrics, error)

	// DeleteOldMetrics 删除旧指标
	DeleteOldMetrics(ctx context.Context, beforeTime time.Time) error
}

// LogRepository 日志仓库接口
type LogRepository interface {
	// CreateLog 创建日志记录
	CreateLog(ctx context.Context, log *model.LogEntry) error

	// GetLogs 获取日志列表
	GetLogs(ctx context.Context, filter *model.LogFilter) ([]model.LogEntry, int64, error)

	// GetLogByID 根据ID获取日志
	GetLogByID(ctx context.Context, id string) (*model.LogEntry, error)

	// GetLogStats 获取日志统计
	GetLogStats(ctx context.Context, startTime, endTime time.Time) (map[string]int64, error)

	// DeleteOldLogs 删除旧日志
	DeleteOldLogs(ctx context.Context, beforeTime time.Time) error

	// CreateLogs 批量创建日志
	CreateLogs(ctx context.Context, logs []model.LogEntry) error
}

// metricsRepository 指标仓库实现
type metricsRepository struct {
	db *gorm.DB
}

// NewMetricsRepository 创建指标仓库实例
func NewMetricsRepository(db *gorm.DB) MetricsRepository {
	return &metricsRepository{db: db}
}

// CreateMetrics 创建指标记录
func (r *metricsRepository) CreateMetrics(ctx context.Context, metrics *model.SystemMetrics) error {
	return r.db.WithContext(ctx).Create(metrics).Error
}

// GetLatestMetrics 获取最新的指标
func (r *metricsRepository) GetLatestMetrics(ctx context.Context) (*model.SystemMetrics, error) {
	var metrics model.SystemMetrics
	err := r.db.WithContext(ctx).
		Order("timestamp DESC").
		First(&metrics).Error
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

// GetMetricsByTimeRange 获取指定时间范围内的指标
func (r *metricsRepository) GetMetricsByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]model.SystemMetrics, error) {
	var metrics []model.SystemMetrics
	err := r.db.WithContext(ctx).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Order("timestamp ASC").
		Find(&metrics).Error
	return metrics, err
}

// GetAverageMetrics 获取平均指标
func (r *metricsRepository) GetAverageMetrics(ctx context.Context, startTime, endTime time.Time) (*model.SystemMetrics, error) {
	var metrics []model.SystemMetrics
	err := r.db.WithContext(ctx).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Find(&metrics).Error
	if err != nil {
		return nil, err
	}

	if len(metrics) == 0 {
		return nil, nil
	}

	// 计算平均值
	avg := &model.SystemMetrics{
		CPUCores:     metrics[0].CPUCores,
		DBMaxOpen:    metrics[0].DBMaxOpen,
		RequestCount: 0,
		ErrorCount:   0,
	}

	count := float64(len(metrics))
	var sumCPUUsage float64
	var sumMemoryAlloc, sumMemoryTotalAlloc, sumMemorySys, sumMemoryHeapAlloc, sumMemoryHeapSys uint64
	var sumGCNum uint32
	var sumGCPauseTotal uint64
	var sumAverageLatency int64
	var sumDBConnections, sumDBInUse, sumDBIdle int
	var sumGoroutineCount int

	for _, m := range metrics {
		sumCPUUsage += m.CPUUsage
		sumMemoryAlloc += m.MemoryAlloc
		sumMemoryTotalAlloc += m.MemoryTotalAlloc
		sumMemorySys += m.MemorySys
		sumMemoryHeapAlloc += m.MemoryHeapAlloc
		sumMemoryHeapSys += m.MemoryHeapSys
		sumGCNum += m.GCNum
		sumGCPauseTotal += m.GCPauseTotal
		sumAverageLatency += m.AverageLatency
		sumDBConnections += m.DBConnections
		sumDBInUse += m.DBInUse
		sumDBIdle += m.DBIdle
		sumGoroutineCount += m.GoroutineCount
		avg.RequestCount += m.RequestCount
		avg.ErrorCount += m.ErrorCount

		if m.GCNext > avg.GCNext {
			avg.GCNext = m.GCNext
		}
	}

	avg.CPUUsage = sumCPUUsage / count
	avg.MemoryAlloc = uint64(float64(sumMemoryAlloc) / count)
	avg.MemoryTotalAlloc = uint64(float64(sumMemoryTotalAlloc) / count)
	avg.MemorySys = uint64(float64(sumMemorySys) / count)
	avg.MemoryHeapAlloc = uint64(float64(sumMemoryHeapAlloc) / count)
	avg.MemoryHeapSys = uint64(float64(sumMemoryHeapSys) / count)
	avg.GCNum = uint32(float64(sumGCNum) / count)
	avg.GCPauseTotal = uint64(float64(sumGCPauseTotal) / count)
	avg.AverageLatency = sumAverageLatency / int64(len(metrics))
	avg.DBConnections = sumDBConnections / len(metrics)
	avg.DBInUse = sumDBInUse / len(metrics)
	avg.DBIdle = sumDBIdle / len(metrics)
	avg.GoroutineCount = sumGoroutineCount / len(metrics)

	return avg, nil
}

// DeleteOldMetrics 删除旧指标
func (r *metricsRepository) DeleteOldMetrics(ctx context.Context, beforeTime time.Time) error {
	return r.db.WithContext(ctx).
		Where("timestamp < ?", beforeTime).
		Delete(&model.SystemMetrics{}).Error
}

// logRepository 日志仓库实现
type logRepository struct {
	db *gorm.DB
}

// NewLogRepository 创建日志仓库实例
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

// CreateLog 创建日志记录
func (r *logRepository) CreateLog(ctx context.Context, log *model.LogEntry) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetLogs 获取日志列表
func (r *logRepository) GetLogs(ctx context.Context, filter *model.LogFilter) ([]model.LogEntry, int64, error) {
	var logs []model.LogEntry
	var total int64

	// 设置默认值
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Size <= 0 {
		filter.Size = 20
	}
	if filter.Size > 100 {
		filter.Size = 100
	}

	query := r.db.WithContext(ctx).Model(&model.LogEntry{})

	// 应用过滤条件
	if filter.Level != "" {
		query = query.Where("level = ?", filter.Level)
	}
	if filter.Service != "" {
		query = query.Where("service = ?", filter.Service)
	}
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Message != "" {
		query = query.Where("message LIKE ?", "%"+filter.Message+"%")
	}
	if filter.StartDate != "" {
		startTime, err := time.Parse(time.RFC3339, filter.StartDate)
		if err == nil {
			query = query.Where("timestamp >= ?", startTime)
		}
	}
	if filter.EndDate != "" {
		endTime, err := time.Parse(time.RFC3339, filter.EndDate)
		if err == nil {
			query = query.Where("timestamp <= ?", endTime)
		}
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (filter.Page - 1) * filter.Size
	err := query.Order("timestamp DESC").
		Offset(offset).
		Limit(filter.Size).
		Find(&logs).Error

	return logs, total, err
}

// GetLogByID 根据ID获取日志
func (r *logRepository) GetLogByID(ctx context.Context, id string) (*model.LogEntry, error) {
	var log model.LogEntry
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetLogStats 获取日志统计
func (r *logRepository) GetLogStats(ctx context.Context, startTime, endTime time.Time) (map[string]int64, error) {
	stats := make(map[string]int64)

	var results []struct {
		Level string
		Count int64
	}

	err := r.db.WithContext(ctx).
		Model(&model.LogEntry{}).
		Select("level, count(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("level").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		stats[result.Level] = result.Count
	}

	return stats, nil
}

// DeleteOldLogs 删除旧日志
func (r *logRepository) DeleteOldLogs(ctx context.Context, beforeTime time.Time) error {
	return r.db.WithContext(ctx).
		Where("timestamp < ?", beforeTime).
		Delete(&model.LogEntry{}).Error
}

// CreateLogs 批量创建日志
func (r *logRepository) CreateLogs(ctx context.Context, logs []model.LogEntry) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&logs).Error
}
