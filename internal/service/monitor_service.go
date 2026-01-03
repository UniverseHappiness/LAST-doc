package service

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"gorm.io/gorm"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/repository"
	"github.com/shirou/gopsutil/v3/cpu"
)

// MonitorService 监控服务接口
type MonitorService interface {
	// CollectMetrics 收集并保存系统指标
	CollectMetrics(ctx context.Context) (*model.SystemMetrics, error)

	// GetCurrentMetrics 获取当前系统指标
	GetCurrentMetrics(ctx context.Context) (*model.SystemMetrics, error)

	// GetMetricsHistory 获取指标历史数据
	GetMetricsHistory(ctx context.Context, startTime, endTime time.Time) ([]model.SystemMetrics, error)

	// GetMetricsReport 获取指标报告
	GetMetricsReport(ctx context.Context, duration time.Duration) (*model.MetricsResponse, error)

	// GetSystemStatus 获取系统状态
	GetSystemStatus(ctx context.Context) (*model.SystemStatus, error)

	// LogRequest 记录请求日志
	LogRequest(ctx context.Context, log *model.LogEntry) error

	// GetLogs 获取日志列表
	GetLogs(ctx context.Context, filter *model.LogFilter) (*model.LogResponse, error)

	// GetLogStats 获取日志统计
	GetLogStats(ctx context.Context, startTime, endTime time.Time) (map[string]int64, error)

	// GetPerformanceReport 获取性能报告
	GetPerformanceReport(ctx context.Context, startTime, endTime time.Time) (*model.PerformanceReport, error)

	// CleanupOldMetrics 清理旧指标
	CleanupOldMetrics(ctx context.Context, retentionDays int) error

	// CleanupOldLogs 清理旧日志
	CleanupOldLogs(ctx context.Context, retentionDays int) error
}

// monitorService 监控服务实现
type monitorService struct {
	metricsRepo repository.MetricsRepository
	logRepo     repository.LogRepository
	db          *gorm.DB
}

// NewMonitorService 创建监控服务实例
func NewMonitorService(
	metricsRepo repository.MetricsRepository,
	logRepo repository.LogRepository,
	db *gorm.DB,
) MonitorService {
	return &monitorService{
		metricsRepo: metricsRepo,
		logRepo:     logRepo,
		db:          db,
	}
}

// CollectMetrics 收集并保存系统指标
func (s *monitorService) CollectMetrics(ctx context.Context) (*model.SystemMetrics, error) {
	// 收集内存统计
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 获取数据库状态
	sqlDB, err := s.db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	dbStats := sqlDB.Stats()

	// 使用gopsutil获取真实的CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Printf("获取CPU使用率失败: %v\n", err)
		cpuPercent = []float64{0.0}
	}

	// 获取数据库配置的最大连接数
	maxOpenConns := dbStats.MaxOpenConnections
	if maxOpenConns == 0 {
		// 如果没有设置，使用GORM默认值或显示实际使用数
		// 显示当前使用的连接数+空闲连接数作为"最大"参考值
		maxOpenConns = dbStats.Idle + dbStats.InUse
	}

	// 计算请求指标（从日志表统计最近1分钟的请求）
	oneMinuteAgo := time.Now().Add(-time.Minute)

	// 定义请求统计结构体
	requestStats := struct {
		Count        int64 `db:"count"`
		ErrorCount   int64 `db:"error_count"`
		TotalLatency int64 `db:"total_latency"`
	}{
		Count:        0,
		ErrorCount:   0,
		TotalLatency: 0,
	}

	// 查询最近1分钟的日志统计（使用原始SQL）
	err = s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) as count,
			COUNT(CASE WHEN level = 'error' OR status_code >= 400 THEN 1 END) as error_count,
			COALESCE(SUM(latency), 0) as total_latency
		FROM log_entries
		WHERE timestamp >= ?
	`, oneMinuteAgo).Scan(&requestStats).Error

	if err != nil {
		// 查询失败，使用默认值（已经初始化为0）
		fmt.Printf("统计请求指标失败: %v\n", err)
	} else {
		fmt.Printf("SQL查询成功: Count=%d, Errors=%d, LatencySum=%d\n",
			requestStats.Count, requestStats.ErrorCount, requestStats.TotalLatency)
	}

	// 计算平均延迟
	var averageLatency int64
	if requestStats.Count > 0 {
		averageLatency = requestStats.TotalLatency / requestStats.Count
	}

	fmt.Printf("请求统计: 总数=%d, 错误=%d, 总延迟=%dms, 平均=%dms\n",
		requestStats.Count, requestStats.ErrorCount, requestStats.TotalLatency, averageLatency)

	// 创建指标记录
	metrics := &model.SystemMetrics{
		Timestamp:      time.Now(),
		CPUCores:       runtime.NumCPU(),
		CPUUsage:       cpuPercent[0], // 使用gopsutil获取的真实CPU使用率
		GoroutineCount: runtime.NumGoroutine(),

		// 内存指标
		MemoryAlloc:      memStats.Alloc,
		MemoryTotalAlloc: memStats.TotalAlloc,
		MemorySys:        memStats.Sys,
		MemoryHeapAlloc:  memStats.HeapAlloc,
		MemoryHeapSys:    memStats.HeapSys,

		// GC指标
		GCNum:        memStats.NumGC,
		GCPauseTotal: memStats.PauseTotalNs,
		GCNext:       memStats.NextGC,

		// 请求指标（从日志表统计）
		RequestCount:   requestStats.Count,
		ErrorCount:     requestStats.ErrorCount,
		AverageLatency: averageLatency,

		// 数据库指标
		DBConnections: dbStats.Idle + dbStats.InUse,
		DBMaxOpen:     maxOpenConns, // 使用计算的最大连接数
		DBInUse:       dbStats.InUse,
		DBIdle:        dbStats.Idle,
	}

	// 保存指标
	err = s.metricsRepo.CreateMetrics(ctx, metrics)
	if err != nil {
		return nil, fmt.Errorf("保存指标失败: %v", err)
	}

	return metrics, nil
}

// GetCurrentMetrics 获取当前系统指标
func (s *monitorService) GetCurrentMetrics(ctx context.Context) (*model.SystemMetrics, error) {
	// 先尝试获取最新的指标
	latest, err := s.metricsRepo.GetLatestMetrics(ctx)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 如果没有指标或指标过期（超过10秒），则收集新指标
	if latest == nil || time.Since(latest.Timestamp) > 10*time.Second {
		return s.CollectMetrics(ctx)
	}

	return latest, nil
}

// GetMetricsHistory 获取指标历史数据
func (s *monitorService) GetMetricsHistory(ctx context.Context, startTime, endTime time.Time) ([]model.SystemMetrics, error) {
	return s.metricsRepo.GetMetricsByTimeRange(ctx, startTime, endTime)
}

// GetMetricsReport 获取指标报告
func (s *monitorService) GetMetricsReport(ctx context.Context, duration time.Duration) (*model.MetricsResponse, error) {
	now := time.Now()
	startTime := now.Add(-duration)

	// 获取当前指标
	current, err := s.GetCurrentMetrics(ctx)
	if err != nil {
		return nil, err
	}

	// 获取历史指标
	history, err := s.metricsRepo.GetMetricsByTimeRange(ctx, startTime, now)
	if err != nil {
		return nil, err
	}

	// 获取平均指标
	average, err := s.metricsRepo.GetAverageMetrics(ctx, startTime, now)
	if err != nil {
		return nil, err
	}

	// 获取系统状态
	status, err := s.GetSystemStatus(ctx)
	if err != nil {
		return nil, err
	}

	return &model.MetricsResponse{
		Current: current,
		History: history,
		Average: average,
		Status:  status,
	}, nil
}

// GetSystemStatus 获取系统状态
func (s *monitorService) GetSystemStatus(ctx context.Context) (*model.SystemStatus, error) {
	metrics, err := s.GetCurrentMetrics(ctx)
	if err != nil {
		return nil, err
	}

	status := &model.SystemStatus{
		Timestamp: time.Now(),
	}

	// 评估CPU状态
	if metrics.CPUUsage > 80 {
		status.CPUStatus = "critical"
	} else if metrics.CPUUsage > 60 {
		status.CPUStatus = "warning"
	} else {
		status.CPUStatus = "healthy"
	}

	// 评估内存状态
	memoryUsagePercent := float64(metrics.MemoryHeapAlloc) / float64(metrics.MemoryHeapSys) * 100
	if memoryUsagePercent > 80 {
		status.MemoryStatus = "critical"
	} else if memoryUsagePercent > 60 {
		status.MemoryStatus = "warning"
	} else {
		status.MemoryStatus = "healthy"
	}

	// 评估数据库状态
	if metrics.DBInUse > metrics.DBMaxOpen*8/10 {
		status.DatabaseStatus = "warning"
	} else {
		status.DatabaseStatus = "healthy"
	}

	// 评估请求状态
	if metrics.ErrorCount > 0 {
		errorRate := float64(metrics.ErrorCount) / float64(metrics.RequestCount) * 100
		if errorRate > 5 {
			status.ServiceStatus = "critical"
		} else if errorRate > 1 {
			status.ServiceStatus = "warning"
		} else {
			status.ServiceStatus = "healthy"
		}
	} else {
		status.ServiceStatus = "healthy"
	}

	// 评估总体状态
	if status.CPUStatus == "critical" || status.MemoryStatus == "critical" ||
		status.DatabaseStatus == "critical" || status.ServiceStatus == "critical" {
		status.OverallStatus = "critical"
	} else if status.CPUStatus == "warning" || status.MemoryStatus == "warning" ||
		status.DatabaseStatus == "warning" || status.ServiceStatus == "warning" {
		status.OverallStatus = "warning"
	} else {
		status.OverallStatus = "healthy"
	}

	return status, nil
}

// LogRequest 记录请求日志
func (s *monitorService) LogRequest(ctx context.Context, log *model.LogEntry) error {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}
	return s.logRepo.CreateLog(ctx, log)
}

// GetLogs 获取日志列表
func (s *monitorService) GetLogs(ctx context.Context, filter *model.LogFilter) (*model.LogResponse, error) {
	logs, total, err := s.logRepo.GetLogs(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(filter.Size) - 1) / int64(filter.Size))

	return &model.LogResponse{
		Logs:       logs,
		Total:      total,
		Page:       filter.Page,
		Size:       filter.Size,
		TotalPages: totalPages,
	}, nil
}

// GetLogStats 获取日志统计
func (s *monitorService) GetLogStats(ctx context.Context, startTime, endTime time.Time) (map[string]int64, error) {
	return s.logRepo.GetLogStats(ctx, startTime, endTime)
}

// GetPerformanceReport 获取性能报告
func (s *monitorService) GetPerformanceReport(ctx context.Context, startTime, endTime time.Time) (*model.PerformanceReport, error) {
	metrics, err := s.metricsRepo.GetMetricsByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	report := &model.PerformanceReport{
		CPUUsage:       make([]float64, len(metrics)),
		MemoryUsage:    make([]uint64, len(metrics)),
		RequestCount:   make([]int64, len(metrics)),
		ErrorCount:     make([]int64, len(metrics)),
		AverageLatency: make([]int64, len(metrics)),
		Timestamp:      make([]time.Time, len(metrics)),
	}

	for i, m := range metrics {
		report.CPUUsage[i] = m.CPUUsage
		report.MemoryUsage[i] = m.MemoryHeapAlloc
		report.RequestCount[i] = m.RequestCount
		report.ErrorCount[i] = m.ErrorCount
		report.AverageLatency[i] = m.AverageLatency
		report.Timestamp[i] = m.Timestamp
	}

	return report, nil
}

// CleanupOldMetrics 清理旧指标
func (s *monitorService) CleanupOldMetrics(ctx context.Context, retentionDays int) error {
	beforeTime := time.Now().AddDate(0, 0, -retentionDays)
	return s.metricsRepo.DeleteOldMetrics(ctx, beforeTime)
}

// CleanupOldLogs 清理旧日志
func (s *monitorService) CleanupOldLogs(ctx context.Context, retentionDays int) error {
	beforeTime := time.Now().AddDate(0, 0, -retentionDays)
	return s.logRepo.DeleteOldLogs(ctx, beforeTime)
}
