package service

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// HealthStatus 表示健康状态
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
)

// ComponentHealth 组件健康状态
type ComponentHealth struct {
	Name      string       `json:"name"`
	Status    HealthStatus `json:"status"`
	Message   string       `json:"message"`
	Timestamp time.Time    `json:"timestamp"`
}

// SystemHealth 系统健康状态
type SystemHealth struct {
	Status     HealthStatus      `json:"status"`
	Timestamp  time.Time         `json:"timestamp"`
	Components []ComponentHealth `json:"components"`
	Uptime     time.Duration     `json:"uptime"`
}

// HealthCheck 健康检查接口
type HealthCheck interface {
	Check(ctx context.Context) ComponentHealth
	Name() string
}

// HealthService 健康检查服务
type HealthService struct {
	db             *sql.DB
	checks         []HealthCheck
	startTime      time.Time
	mu             sync.RWMutex
	circuitBreaker map[string]*CircuitBreaker
}

// NewHealthService 创建健康检查服务
func NewHealthService(db *sql.DB) *HealthService {
	return &HealthService{
		db:             db,
		checks:         make([]HealthCheck, 0),
		startTime:      time.Now(),
		circuitBreaker: make(map[string]*CircuitBreaker),
	}
}

// RegisterCheck 注册健康检查
func (s *HealthService) RegisterCheck(check HealthCheck) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checks = append(s.checks, check)
}

// GetCircuitBreaker 获取断路器（如不存在则创建）
func (s *HealthService) GetCircuitBreaker(name string, threshold int, timeout time.Duration) *CircuitBreaker {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cb, exists := s.circuitBreaker[name]; exists {
		return cb
	}

	cb := NewCircuitBreaker(name, threshold, timeout)
	s.circuitBreaker[name] = cb
	return cb
}

// CheckHealth 执行健康检查
func (s *HealthService) CheckHealth(ctx context.Context) *SystemHealth {
	s.mu.RLock()
	defer s.mu.RUnlock()

	health := &SystemHealth{
		Timestamp:  time.Now(),
		Uptime:     time.Since(s.startTime),
		Components: make([]ComponentHealth, 0, len(s.checks)),
	}

	// 检查每个组件
	overallStatus := HealthStatusHealthy
	for _, check := range s.checks {
		componentHealth := check.Check(ctx)
		health.Components = append(health.Components, componentHealth)

		// 更新整体状态
		if componentHealth.Status == HealthStatusUnhealthy {
			overallStatus = HealthStatusUnhealthy
		} else if componentHealth.Status == HealthStatusDegraded && overallStatus == HealthStatusHealthy {
			overallStatus = HealthStatusDegraded
		}
	}

	health.Status = overallStatus
	return health
}

// DatabaseHealthCheck 数据库健康检查
type DatabaseHealthCheck struct {
	db *sql.DB
}

// NewDatabaseHealthCheck 创建数据库健康检查
func NewDatabaseHealthCheck(db *sql.DB) *DatabaseHealthCheck {
	return &DatabaseHealthCheck{db: db}
}

// Check 执行检查
func (c *DatabaseHealthCheck) Check(ctx context.Context) ComponentHealth {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := c.db.PingContext(ctx)
	if err != nil {
		return ComponentHealth{
			Name:      c.Name(),
			Status:    HealthStatusUnhealthy,
			Message:   fmt.Sprintf("数据库连接失败: %v", err),
			Timestamp: time.Now(),
		}
	}

	return ComponentHealth{
		Name:      c.Name(),
		Status:    HealthStatusHealthy,
		Message:   "数据库连接正常",
		Timestamp: time.Now(),
	}
}

// Name 返回组件名称
func (c *DatabaseHealthCheck) Name() string {
	return "database"
}

// StorageHealthCheck 存储健康检查
type StorageHealthCheck struct {
	storageService StorageService
}

// NewStorageHealthCheck 创建存储健康检查
func NewStorageHealthCheck(storageService StorageService) *StorageHealthCheck {
	return &StorageHealthCheck{storageService: storageService}
}

// Check 执行检查
func (c *StorageHealthCheck) Check(ctx context.Context) ComponentHealth {
	// 检查存储是否可访问
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := c.storageService.HealthCheck(ctx); err != nil {
		return ComponentHealth{
			Name:      c.Name(),
			Status:    HealthStatusUnhealthy,
			Message:   fmt.Sprintf("存储服务不可用: %v", err),
			Timestamp: time.Now(),
		}
	}

	return ComponentHealth{
		Name:      c.Name(),
		Status:    HealthStatusHealthy,
		Message:   "存储服务正常",
		Timestamp: time.Now(),
	}
}

// Name 返回组件名称
func (c *StorageHealthCheck) Name() string {
	return "storage"
}
