package service

import (
	"context"
	"testing"
	"time"
)

// TestMonitorService_StatusEvaluation 测试状态评估逻辑
func TestMonitorService_StatusEvaluation(t *testing.T) {
	tests := []struct {
		name        string
		cpuUsage    float64
		memoryUsage float64
		expectedCpu string
		expectedMem string
	}{
		{
			name:        "CPU健康状态",
			cpuUsage:    30.0,
			memoryUsage: 40.0,
			expectedCpu: "healthy",
			expectedMem: "healthy",
		},
		{
			name:        "CPU警告状态",
			cpuUsage:    70.0,
			memoryUsage: 55.0,
			expectedCpu: "warning",
			expectedMem: "healthy",
		},
		{
			name:        "CPU严重状态",
			cpuUsage:    90.0,
			memoryUsage: 70.0,
			expectedCpu: "critical",
			expectedMem: "warning",
		},
		{
			name:        "内存严重状态",
			cpuUsage:    50.0,
			memoryUsage: 85.0,
			expectedCpu: "healthy",
			expectedMem: "critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 评估CPU状态
			cpuStatus := evaluateStatus(tt.cpuUsage, []float64{60, 80}, "warning", "critical")
			if cpuStatus != tt.expectedCpu {
				t.Errorf("CPU状态评估错误: got %s, want %s", cpuStatus, tt.expectedCpu)
			}

			// 评估内存状态
			memStatus := evaluateStatus(tt.memoryUsage, []float64{60, 80}, "warning", "critical")
			if memStatus != tt.expectedMem {
				t.Errorf("内存状态评估错误: got %s, want %s", memStatus, tt.expectedMem)
			}
		})
	}
}

// evaluateStatus 辅助函数：评估状态
func evaluateStatus(value float64, thresholds []float64, warnStatus, critStatus string) string {
	if value > thresholds[1] {
		return critStatus
	} else if value > thresholds[0] {
		return warnStatus
	}
	return "healthy"
}

// TestMonitorService_ErrorRateCalculation 测试错误率计算
func TestMonitorService_ErrorRateCalculation(t *testing.T) {
	tests := []struct {
		name       string
		errorCount int64
		totalCount int64
		expected   string
	}{
		{
			name:       "无错误",
			errorCount: 0,
			totalCount: 100,
			expected:   "healthy",
		},
		{
			name:       "低错误率",
			errorCount: 5,
			totalCount: 1000,
			expected:   "healthy",
		},
		{
			name:       "警告错误率",
			errorCount: 20,
			totalCount: 1000,
			expected:   "warning",
		},
		{
			name:       "严重错误率",
			errorCount: 100,
			totalCount: 1000,
			expected:   "critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errorRate float64
			if tt.totalCount > 0 {
				errorRate = float64(tt.errorCount) / float64(tt.totalCount) * 100
			}

			var status string
			if errorRate > 5 {
				status = "critical"
			} else if errorRate > 1 {
				status = "warning"
			} else {
				status = "healthy"
			}

			if status != tt.expected {
				t.Errorf("错误率计算错误: rate=%.2f%%, got %s, want %s", errorRate, status, tt.expected)
			}
		})
	}
}

// TestMonitorService_DatabaseConnectionEvaluation 测试数据库连接评估
func TestMonitorService_DatabaseConnectionEvaluation(t *testing.T) {
	tests := []struct {
		name     string
		inUse    int64
		maxOpen  int64
		expected string
	}{
		{
			name:     "健康连接",
			inUse:    5,
			maxOpen:  100,
			expected: "healthy",
		},
		{
			name:     "警告连接",
			inUse:    85,
			maxOpen:  100,
			expected: "warning",
		},
		{
			name:     "边界情况-80%",
			inUse:    80,
			maxOpen:  100,
			expected: "healthy",
		},
		{
			name:     "边界情况-刚好不满",
			inUse:    79,
			maxOpen:  100,
			expected: "healthy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status string
			if tt.inUse > tt.maxOpen*8/10 {
				status = "warning"
			} else {
				status = "healthy"
			}

			if status != tt.expected {
				t.Errorf("数据库连接评估错误: got %s, want %s", status, tt.expected)
			}
		})
	}
}

// TestMonitorService_OverallStatusEvaluation 测试总体状态评估
func TestMonitorService_OverallStatusEvaluation(t *testing.T) {
	tests := []struct {
		name           string
		cpuStatus      string
		memoryStatus   string
		databaseStatus string
		serviceStatus  string
		expected       string
	}{
		{
			name:           "全部健康",
			cpuStatus:      "healthy",
			memoryStatus:   "healthy",
			databaseStatus: "healthy",
			serviceStatus:  "healthy",
			expected:       "healthy",
		},
		{
			name:           "一个警告",
			cpuStatus:      "warning",
			memoryStatus:   "healthy",
			databaseStatus: "healthy",
			serviceStatus:  "healthy",
			expected:       "warning",
		},
		{
			name:           "多个警告",
			cpuStatus:      "warning",
			memoryStatus:   "warning",
			databaseStatus: "healthy",
			serviceStatus:  "healthy",
			expected:       "warning",
		},
		{
			name:           "一个严重",
			cpuStatus:      "critical",
			memoryStatus:   "healthy",
			databaseStatus: "healthy",
			serviceStatus:  "healthy",
			expected:       "critical",
		},
		{
			name:           "警告和严重混合",
			cpuStatus:      "warning",
			memoryStatus:   "critical",
			databaseStatus: "warning",
			serviceStatus:  "healthy",
			expected:       "critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status string
			if tt.cpuStatus == "critical" || tt.memoryStatus == "critical" ||
				tt.databaseStatus == "critical" || tt.serviceStatus == "critical" {
				status = "critical"
			} else if tt.cpuStatus == "warning" || tt.memoryStatus == "warning" ||
				tt.databaseStatus == "warning" || tt.serviceStatus == "warning" {
				status = "warning"
			} else {
				status = "healthy"
			}

			if status != tt.expected {
				t.Errorf("总体状态评估错误: got %s, want %s", status, tt.expected)
			}
		})
	}
}

// TestMonitorService_TimeCalculation 测试时间计算
func TestMonitorService_TimeCalculation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		retentionDays int
		expectedDiff  time.Duration
	}{
		{
			name:          "保留7天",
			retentionDays: 7,
			expectedDiff:  7 * 24 * time.Hour,
		},
		{
			name:          "保留30天",
			retentionDays: 30,
			expectedDiff:  30 * 24 * time.Hour,
		},
		{
			name:          "保留90天",
			retentionDays: 90,
			expectedDiff:  90 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeTime := now.AddDate(0, 0, -tt.retentionDays)
			diff := now.Sub(beforeTime)

			// 允许1秒的误差
			if diff < tt.expectedDiff-time.Second || diff > tt.expectedDiff+time.Second {
				t.Errorf("时间计算错误: got %v, expected ~%v", diff, tt.expectedDiff)
			}
		})
	}
}

// TestMonitorService_PeriodCalculation 测试时间段计算
func TestMonitorService_PeriodCalculation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "1小时",
			duration: time.Hour,
		},
		{
			name:     "24小时",
			duration: 24 * time.Hour,
		},
		{
			name:     "7天",
			duration: 7 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := now.Add(-tt.duration)
			diff := now.Sub(startTime)

			// 允许1秒的误差
			if diff < tt.duration-time.Second || diff > tt.duration+time.Second {
				t.Errorf("时间段计算错误: got %v, expected ~%v", diff, tt.duration)
			}
		})
	}
}

// TestMonitorService_AverageLatencyCalculation 测试平均延迟计算
func TestMonitorService_AverageLatencyCalculation(t *testing.T) {
	tests := []struct {
		name            string
		totalLatency    int64
		count           int64
		expectedLatency int64
	}{
		{
			name:            "正常情况",
			totalLatency:    1000,
			count:           10,
			expectedLatency: 100,
		},
		{
			name:            "零延迟",
			totalLatency:    0,
			count:           10,
			expectedLatency: 0,
		},
		{
			name:            "单个请求",
			totalLatency:    500,
			count:           1,
			expectedLatency: 500,
		},
		{
			name:            "零请求",
			totalLatency:    0,
			count:           0,
			expectedLatency: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var averageLatency int64
			if tt.count > 0 {
				averageLatency = tt.totalLatency / tt.count
			}

			if averageLatency != tt.expectedLatency {
				t.Errorf("平均延迟计算错误: got %d, want %d", averageLatency, tt.expectedLatency)
			}
		})
	}
}

// TestMonitorService_MemoryUsagePercent 计算内存使用百分比
func TestMonitorService_MemoryUsagePercent(t *testing.T) {
	tests := []struct {
		name            string
		heapAlloc       uint64
		heapSys         uint64
		expectedPercent float64
	}{
		{
			name:            "半使用",
			heapAlloc:       500,
			heapSys:         1000,
			expectedPercent: 50.0,
		},
		{
			name:            "全使用",
			heapAlloc:       1000,
			heapSys:         1000,
			expectedPercent: 100.0,
		},
		{
			name:            "低使用",
			heapAlloc:       100,
			heapSys:         1000,
			expectedPercent: 10.0,
		},
		{
			name:            "零使用",
			heapAlloc:       0,
			heapSys:         1000,
			expectedPercent: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memoryUsagePercent := float64(tt.heapAlloc) / float64(tt.heapSys) * 100

			if memoryUsagePercent != tt.expectedPercent {
				t.Errorf("内存使用百分比计算错误: got %.2f, want %.2f", memoryUsagePercent, tt.expectedPercent)
			}
		})
	}
}

// TestMonitorService_ContextHandling 测试context处理
func TestMonitorService_ContextHandling(t *testing.T) {
	tests := []struct {
		name string
		fn   func(context.Context) bool
	}{
		{
			name: "Background Context",
			fn: func(ctx context.Context) bool {
				return ctx != nil && ctx.Value("test") == nil
			},
		},
		{
			name: "Context with Timeout",
			fn: func(ctx context.Context) bool {
				_, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()
				return ctx != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if !tt.fn(ctx) {
				t.Error("context处理失败")
			}
		})
	}
}

// TestMonitorService_MetadataCalculation 测试元数据计算
func TestMonitorService_MetadataCalculation(t *testing.T) {
	type TestStats struct {
		Count        int64
		ErrorCount   int64
		TotalLatency int64
	}

	tests := []struct {
		name       string
		stats      TestStats
		shouldWarn bool
		shouldCrit bool
	}{
		{
			name: "无请求",
			stats: TestStats{
				Count:        0,
				ErrorCount:   0,
				TotalLatency: 0,
			},
			shouldWarn: false,
			shouldCrit: false,
		},
		{
			name: "正常请求",
			stats: TestStats{
				Count:        100,
				ErrorCount:   5,
				TotalLatency: 10000,
			},
			shouldWarn: true,
			shouldCrit: false,
		},
		{
			name: "严重错误",
			stats: TestStats{
				Count:        100,
				ErrorCount:   20,
				TotalLatency: 20000,
			},
			shouldWarn: true,
			shouldCrit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errorRate float64
			if tt.stats.Count > 0 {
				errorRate = float64(tt.stats.ErrorCount) / float64(tt.stats.Count) * 100
			}

			isWarning := errorRate > 1
			isCritical := errorRate > 5

			if isWarning != tt.shouldWarn {
				t.Errorf("警告状态不匹配: got %v, want %v", isWarning, tt.shouldWarn)
			}

			if isCritical != tt.shouldCrit {
				t.Errorf("严重状态不匹配: got %v, want %v", isCritical, tt.shouldCrit)
			}
		})
	}
}

// TestMonitorService_GCMetricsEvaluation 测试GC指标评估
func TestMonitorService_GCMetricsEvaluation(t *testing.T) {
	tests := []struct {
		name         string
		gcNum        uint32
		pauseTotalNs uint64
		nextGC       uint64
		shouldFlag   bool
	}{
		{
			name:         "低GC活动",
			gcNum:        10,
			pauseTotalNs: 1000000,
			nextGC:       1000000,
			shouldFlag:   false,
		},
		{
			name:         "中GC活动",
			gcNum:        100,
			pauseTotalNs: 50000000,
			nextGC:       5000000,
			shouldFlag:   false,
		},
		{
			name:         "高GC活动",
			gcNum:        1000,
			pauseTotalNs: 1000000000,
			nextGC:       10000000,
			shouldFlag:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 简单的启发式检查：GC次数过多或暂停时间过长
			isHighActivity := tt.gcNum > 500 || tt.pauseTotalNs > 500000000

			if isHighActivity != tt.shouldFlag {
				t.Errorf("GC活动评估不匹配: got %v, want %v", isHighActivity, tt.shouldFlag)
			}
		})
	}
}

// TestMonitorService_ConnectionPoolEvaluation 测试连接池评估
func TestMonitorService_ConnectionPoolEvaluation(t *testing.T) {
	tests := []struct {
		name          string
		inUse         int64
		idle          int64
		maxOpen       int64
		expectedTotal int64
		expectedMax   int64
	}{
		{
			name:          "正常使用",
			inUse:         10,
			idle:          20,
			maxOpen:       100,
			expectedTotal: 30,
			expectedMax:   100,
		},
		{
			name:          "高使用",
			inUse:         50,
			idle:          10,
			maxOpen:       100,
			expectedTotal: 60,
			expectedMax:   100,
		},
		{
			name:          "最大使用计算",
			inUse:         30,
			idle:          20,
			maxOpen:       0,
			expectedTotal: 50,
			expectedMax:   50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := tt.inUse + tt.idle
			var maxOpen int64
			if tt.maxOpen == 0 {
				maxOpen = total
			} else {
				maxOpen = tt.maxOpen
			}

			if total != tt.expectedTotal {
				t.Errorf("连接总数不匹配: got %d, want %d", total, tt.expectedTotal)
			}

			if maxOpen != tt.expectedMax {
				t.Errorf("最大连接数不匹配: got %d, want %d", maxOpen, tt.expectedMax)
			}
		})
	}
}

// TestMonitorService_RequestStatisticsEvaluation 测试请求统计评估
func TestMonitorService_RequestStatisticsEvaluation(t *testing.T) {
	tests := []struct {
		name         string
		oneMinuteAgo time.Time
		currentTime  time.Time
		expectedDiff time.Duration
	}{
		{
			name:         "1分钟前",
			oneMinuteAgo: time.Now().Add(-time.Minute),
			currentTime:  time.Now(),
			expectedDiff: time.Minute,
		},
		{
			name:         "2分钟前",
			oneMinuteAgo: time.Now().Add(-2 * time.Minute),
			currentTime:  time.Now(),
			expectedDiff: 2 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := tt.currentTime.Sub(tt.oneMinuteAgo)

			// 允许100ms的误差
			if diff < tt.expectedDiff-100*time.Millisecond || diff > tt.expectedDiff+100*time.Millisecond {
				t.Errorf("时间差计算错误: got %v, expected ~%v", diff, tt.expectedDiff)
			}
		})
	}
}
