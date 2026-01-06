package service

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

// MockHealthCheck 模拟健康检查器，用于测试
type MockHealthCheck struct {
	name    string
	status  HealthStatus
	message string
}

func (m *MockHealthCheck) Check(ctx context.Context) ComponentHealth {
	return ComponentHealth{
		Name:      m.name,
		Status:    m.status,
		Message:   m.message,
		Timestamp: time.Now(),
	}
}

func (m *MockHealthCheck) Name() string {
	return m.name
}

// TestNewHealthService 测试创建健康服务
func TestNewHealthService(t *testing.T) {
	var db *sql.DB // nil DB for testing

	service := NewHealthService(db)

	if service == nil {
		t.Fatal("NewHealthService() 返回 nil")
	}

	// 验证初始化
	if service.startTime.IsZero() {
		t.Error("服务启动时间未设置")
	}

	if service.checks == nil {
		t.Error("健康检查列表未初始化")
	}

	if len(service.checks) != 0 {
		t.Errorf("初始健康检查列表长度 = %d, expected 0", len(service.checks))
	}

	if service.circuitBreaker == nil {
		t.Error("断路器映射未初始化")
	}
}

// TestRegisterCheck 测试注册健康检查
func TestRegisterCheck(t *testing.T) {
	service := NewHealthService(nil)

	mockCheck := &MockHealthCheck{
		name:    "test-check",
		status:  HealthStatusHealthy,
		message: "All good",
	}

	// 注册检查器
	service.RegisterCheck(mockCheck)

	// 验证注册成功
	if len(service.checks) != 1 {
		t.Errorf("健康检查列表长度 = %d, expected 1", len(service.checks))
	}

	if service.checks[0] != mockCheck {
		t.Error("注册的检查器不匹配")
	}
}

// TestRegisterCheck_Multiple 测试注册多个健康检查
func TestRegisterCheck_Multiple(t *testing.T) {
	service := NewHealthService(nil)

	checks := []*MockHealthCheck{
		{name: "check-1", status: HealthStatusHealthy, message: "OK"},
		{name: "check-2", status: HealthStatusHealthy, message: "OK"},
		{name: "check-3", status: HealthStatusDegraded, message: "Warning"},
	}

	// 注册多个检查器
	for _, check := range checks {
		service.RegisterCheck(check)
	}

	// 验证所有检查器都被注册
	if len(service.checks) != len(checks) {
		t.Errorf("健康检查列表长度 = %d, expected %d", len(service.checks), len(checks))
	}
}

// TestCheckHealth_NoChecks 测试无健康检查器时的状态
func TestCheckHealth_NoChecks(t *testing.T) {
	service := NewHealthService(nil)
	ctx := context.Background()

	health := service.CheckHealth(ctx)

	if health == nil {
		t.Fatal("CheckHealth() 返回 nil")
	}

	// 无检查器时应为健康状态
	if health.Status != HealthStatusHealthy {
		t.Errorf("健康状态 = %s, expected %s", health.Status, HealthStatusHealthy)
	}

	if health.Components == nil {
		t.Error("组件列表为nil")
	}

	if len(health.Components) != 0 {
		t.Errorf("组件列表长度 = %d, expected 0", len(health.Components))
	}
}

// TestCheckHealth_AllHealthy 测试所有组件健康
func TestCheckHealth_AllHealthy(t *testing.T) {
	service := NewHealthService(nil)
	ctx := context.Background()

	// 注册健康的检查器
	service.RegisterCheck(&MockHealthCheck{
		name:    "component-1",
		status:  HealthStatusHealthy,
		message: "OK",
	})
	service.RegisterCheck(&MockHealthCheck{
		name:    "component-2",
		status:  HealthStatusHealthy,
		message: "OK",
	})

	health := service.CheckHealth(ctx)

	if health.Status != HealthStatusHealthy {
		t.Errorf("健康状态 = %s, expected %s", health.Status, HealthStatusHealthy)
	}

	if len(health.Components) != 2 {
		t.Errorf("组件数量 = %d, expected 2", len(health.Components))
	}
}

// TestCheckHealth_OneUnhealthy 测试单个组件不健康
func TestCheckHealth_OneUnhealthy(t *testing.T) {
	service := NewHealthService(nil)
	ctx := context.Background()

	// 注册检查器，其中一个不健康
	service.RegisterCheck(&MockHealthCheck{
		name:    "component-1",
		status:  HealthStatusHealthy,
		message: "OK",
	})
	service.RegisterCheck(&MockHealthCheck{
		name:    "component-2",
		status:  HealthStatusUnhealthy,
		message: "Failed",
	})

	health := service.CheckHealth(ctx)

	// 有不健康组件时，整体状态应该是不健康
	if health.Status != HealthStatusUnhealthy {
		t.Errorf("健康状态 = %s, expected %s", health.Status, HealthStatusUnhealthy)
	}
}

// TestCheckHealth_OneDegraded 测试单个组件降级
func TestCheckHealth_OneDegraded(t *testing.T) {
	service := NewHealthService(nil)
	ctx := context.Background()

	// 注册检查器，其中一个降级
	service.RegisterCheck(&MockHealthCheck{
		name:    "component-1",
		status:  HealthStatusHealthy,
		message: "OK",
	})
	service.RegisterCheck(&MockHealthCheck{
		name:    "component-2",
		status:  HealthStatusDegraded,
		message: "Warning",
	})

	health := service.CheckHealth(ctx)

	// 有降级组件时，整体状态应该是降级
	if health.Status != HealthStatusDegraded {
		t.Errorf("健康状态 = %s, expected %s", health.Status, HealthStatusDegraded)
	}
}

// TestCheckHealth_MultipleDegraded 测试多个组件降级
func TestCheckHealth_MultipleDegraded(t *testing.T) {
	service := NewHealthService(nil)
	ctx := context.Background()

	// 注册检查器，多个降级
	service.RegisterCheck(&MockHealthCheck{
		name:    "component-1",
		status:  HealthStatusHealthy,
		message: "OK",
	})
	service.RegisterCheck(&MockHealthCheck{
		name:    "component-2",
		status:  HealthStatusDegraded,
		message: "Warning 1",
	})
	service.RegisterCheck(&MockHealthCheck{
		name:    "component-3",
		status:  HealthStatusDegraded,
		message: "Warning 2",
	})

	health := service.CheckHealth(ctx)

	if health.Status != HealthStatusDegraded {
		t.Errorf("健康状态 = %s, expected %s", health.Status, HealthStatusDegraded)
	}
}

// TestCheckHealth_Uptime 测试运行时间计算
func TestCheckHealth_Uptime(t *testing.T) {
	service := NewHealthService(nil)
	ctx := context.Background()

	// 等待一小段时间
	time.Sleep(10 * time.Millisecond)

	health := service.CheckHealth(ctx)

	// 验证运行时间大于0
	if health.Uptime <= 0 {
		t.Errorf("运行时间 = %v, should be > 0", health.Uptime)
	}

	// 验证运行时间合理（不超过1秒）
	if health.Uptime > time.Second {
		t.Errorf("运行时间 = %v, should be < 1s", health.Uptime)
	}
}

// TestCheckHealth_Timestamp 测试时间戳
func TestCheckHealth_Timestamp(t *testing.T) {
	service := NewHealthService(nil)
	ctx := context.Background()

	before := time.Now()
	health := service.CheckHealth(ctx)
	after := time.Now()

	// 验证时间戳在合理范围内
	if health.Timestamp.Before(before) || health.Timestamp.After(after) {
		t.Errorf("健康检查时间戳 = %v, not in expected range [%v, %v]",
			health.Timestamp, before, after)
	}
}

// TestHealthStatus_Constants 测试健康状态常量
func TestHealthStatus_Constants(t *testing.T) {
	tests := []struct {
		name   string
		status HealthStatus
	}{
		{
			name:   "健康状态",
			status: HealthStatusHealthy,
		},
		{
			name:   "不健康状态",
			status: HealthStatusUnhealthy,
		},
		{
			name:   "降级状态",
			status: HealthStatusDegraded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusStr := string(tt.status)
			if statusStr == "" {
				t.Error("健康状态字符串为空")
			}
		})
	}
}

// TestComponentHealth 测试组件健康结构
func TestComponentHealth(t *testing.T) {
	now := time.Now()
	component := ComponentHealth{
		Name:      "test-component",
		Status:    HealthStatusHealthy,
		Message:   "All systems operational",
		Timestamp: now,
	}

	// 验证各字段
	if component.Name != "test-component" {
		t.Errorf("组件名称 = %s, expected test-component", component.Name)
	}

	if component.Status != HealthStatusHealthy {
		t.Errorf("组件状态 = %s, expected %s", component.Status, HealthStatusHealthy)
	}

	if component.Message != "All systems operational" {
		t.Errorf("组件消息 = %s, expected 'All systems operational'", component.Message)
	}

	if component.Timestamp.IsZero() {
		t.Error("组件时间戳为零值")
	}
}

// TestSystemHealth 测试系统健康结构
func TestSystemHealth(t *testing.T) {
	now := time.Now()
	health := SystemHealth{
		Status:    HealthStatusHealthy,
		Timestamp: now,
		Components: []ComponentHealth{
			{
				Name:      "db",
				Status:    HealthStatusHealthy,
				Message:   "OK",
				Timestamp: now,
			},
		},
		Uptime: time.Hour,
	}

	// 验证各字段
	if health.Status != HealthStatusHealthy {
		t.Errorf("系统状态 = %s, expected %s", health.Status, HealthStatusHealthy)
	}

	if health.Timestamp.IsZero() {
		t.Error("系统时间戳为零值")
	}

	if len(health.Components) != 1 {
		t.Errorf("组件数量 = %d, expected 1", len(health.Components))
	}

	if health.Uptime != time.Hour {
		t.Errorf("运行时间 = %v, expected 1h", health.Uptime)
	}
}

// TestNewDatabaseHealthCheck 测试创建数据库健康检查器
func TestNewDatabaseHealthCheck(t *testing.T) {
	var db *sql.DB // nil DB for testing

	check := NewDatabaseHealthCheck(db)

	if check == nil {
		t.Fatal("NewDatabaseHealthCheck() 返回 nil")
	}

	// 验证名称
	if check.Name() != "database" {
		t.Errorf("检查器名称 = %s, expected 'database'", check.Name())
	}
}

// TestDatabaseHealthCheck_Name 测试数据库检查器名称
func TestDatabaseHealthCheck_Name(t *testing.T) {
	check := &DatabaseHealthCheck{}

	name := check.Name()

	if name != "database" {
		t.Errorf("名称 = %s, expected 'database'", name)
	}
}

// TestStorageHealthCheck_Name 测试存储检查器名称
func TestStorageHealthCheck_Name(t *testing.T) {
	check := &StorageHealthCheck{}

	name := check.Name()

	if name != "storage" {
		t.Errorf("名称 = %s, expected 'storage'", name)
	}
}
