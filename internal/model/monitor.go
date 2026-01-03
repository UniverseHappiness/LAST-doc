package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// TimeWithTimeZone 自定义时间类型，确保JSON序列化时使用RFC3339格式
type TimeWithTimeZone time.Time

// MarshalJSON 实现JSON序列化
func (t TimeWithTimeZone) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.RFC3339))
}

// UnmarshalJSON 实现JSON反序列化
func (t *TimeWithTimeZone) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	tm, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	*t = TimeWithTimeZone(tm)
	return nil
}

// Value 实现driver.Valuer接口
func (t TimeWithTimeZone) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 实现sql.Scanner接口
func (t *TimeWithTimeZone) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*t = TimeWithTimeZone(v)
	case []byte:
		tm, err := time.Parse("2006-01-02 15:04:05.999999", string(v))
		if err != nil {
			return err
		}
		*t = TimeWithTimeZone(tm)
	case string:
		tm, err := time.Parse("2006-01-02 15:04:05.999999", v)
		if err != nil {
			return err
		}
		*t = TimeWithTimeZone(tm)
	default:
		return fmt.Errorf("cannot scan %T into TimeWithTimeZone", value)
	}

	return nil
}

// ToTime 转换为标准time.Time
func (t TimeWithTimeZone) ToTime() time.Time {
	return time.Time(t)
}

// SystemMetrics 系统性能指标
type SystemMetrics struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Timestamp time.Time `json:"timestamp" gorm:"not null;index"`

	// CPU指标
	CPUUsage       float64 `json:"cpu_usage"`       // CPU使用率（百分比）
	CPUCores       int     `json:"cpu_cores"`       // CPU核心数
	GoroutineCount int     `json:"goroutine_count"` // Goroutine数量

	// 内存指标
	MemoryAlloc      uint64 `json:"memory_alloc"`       // 已分配内存（字节）
	MemoryTotalAlloc uint64 `json:"memory_total_alloc"` // 累计分配内存（字节）
	MemorySys        uint64 `json:"memory_sys"`         // 从系统获取的内存（字节）
	MemoryHeapAlloc  uint64 `json:"memory_heap_alloc"`  // 堆内存分配（字节）
	MemoryHeapSys    uint64 `json:"memory_heap_sys"`    // 堆内存系统（字节）

	// GC指标
	GCNum        uint32 `json:"gc_num"`         // GC次数
	GCPauseTotal uint64 `json:"gc_pause_total"` // GC总暂停时间（纳秒）
	GCNext       uint64 `json:"gc_next"`        // 下次GC的目标堆大小（字节）

	// 请求指标
	RequestCount   int64 `json:"request_count"`   // 请求总数
	ErrorCount     int64 `json:"error_count"`     // 错误总数
	AverageLatency int64 `json:"average_latency"` // 平均延迟（毫秒）

	// 数据库指标
	DBConnections int `json:"db_connections"` // 数据库连接数
	DBMaxOpen     int `json:"db_max_open"`    // 数据库最大连接数
	DBInUse       int `json:"db_in_use"`      // 数据库使用中的连接数
	DBIdle        int `json:"db_idle"`        // 数据库空闲连接数

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 返回系统指标表名
func (SystemMetrics) TableName() string {
	return "system_metrics"
}

// SystemStatus 系统状态
type SystemStatus struct {
	OverallStatus  string    `json:"overall_status"`  // 总体状态：healthy, warning, critical
	CPUStatus      string    `json:"cpu_status"`      // CPU状态
	MemoryStatus   string    `json:"memory_status"`   // 内存状态
	DatabaseStatus string    `json:"database_status"` // 数据库状态
	ServiceStatus  string    `json:"service_status"`  // 服务状态
	Timestamp      time.Time `json:"timestamp"`
}

// LogEntry 日志条目
type LogEntry struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Timestamp  time.Time `json:"timestamp" gorm:"not null;index"`
	Level      string    `json:"level" gorm:"not null;index"`   // debug, info, warn, error
	Service    string    `json:"service" gorm:"not null;index"` // 服务名称
	Message    string    `json:"message" gorm:"not null"`
	Method     string    `json:"method"`               // HTTP方法
	Path       string    `json:"path"`                 // 请求路径
	UserID     string    `json:"user_id" gorm:"index"` // 用户ID
	Username   string    `json:"username"`             // 用户名
	StatusCode int       `json:"status_code"`          // HTTP状态码
	Latency    int64     `json:"latency"`              // 请求延迟（毫秒）
	IPAddress  string    `json:"ip_address"`           // IP地址
	UserAgent  string    `json:"user_agent"`           // 用户代理
	Error      string    `json:"error"`                // 错误信息
	Stack      string    `json:"stack"`                // 堆栈信息

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 返回日志表名
func (LogEntry) TableName() string {
	return "log_entries"
}

// LogFilter 日志过滤条件
type LogFilter struct {
	Level     string `json:"level"`      // 日志级别
	Service   string `json:"service"`    // 服务名称
	UserID    string `json:"user_id"`    // 用户ID
	StartDate string `json:"start_date"` // 开始日期
	EndDate   string `json:"end_date"`   // 结束日期
	Message   string `json:"message"`    // 消息关键字
	Page      int    `json:"page"`       // 页码
	Size      int    `json:"size"`       // 每页大小
}

// MetricsResponse 指标响应
type MetricsResponse struct {
	Current *SystemMetrics  `json:"current"`
	History []SystemMetrics `json:"history"`
	Average *SystemMetrics  `json:"average"`
	Status  *SystemStatus   `json:"status"`
}

// LogResponse 日志响应
type LogResponse struct {
	Logs       []LogEntry `json:"logs"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Size       int        `json:"size"`
	TotalPages int        `json:"total_pages"`
}

// PerformanceReport 性能报告
type PerformanceReport struct {
	CPUUsage       []float64   `json:"cpu_usage"`       // CPU使用率历史
	MemoryUsage    []uint64    `json:"memory_usage"`    // 内存使用历史
	RequestCount   []int64     `json:"request_count"`   // 请求数量历史
	ErrorCount     []int64     `json:"error_count"`     // 错误数量历史
	AverageLatency []int64     `json:"average_latency"` // 平均延迟历史
	Timestamp      []time.Time `json:"timestamp"`       // 时间戳
}
