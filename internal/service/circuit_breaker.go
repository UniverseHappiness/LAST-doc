package service

import (
	"sync"
	"time"
)

// State 断路器状态
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker 断路器
type CircuitBreaker struct {
	name         string
	maxFailures  int
	timeout      time.Duration
	mu           sync.RWMutex
	state        State
	failures     int
	lastFailTime time.Time
}

// NewCircuitBreaker 创建断路器
func NewCircuitBreaker(name string, maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:        name,
		maxFailures: maxFailures,
		timeout:     timeout,
		state:       StateClosed,
		failures:    0,
	}
}

// Execute 执行操作，带断路器保护
func (cb *CircuitBreaker) Execute(fn func() error) error {
	// 检查断路器状态
	if cb.allow() {
		err := fn()
		cb.recordResult(err)
		return err
	}

	return &CircuitBreakerOpenError{Breaker: cb.name}
}

// allow 判断是否允许执行
func (cb *CircuitBreaker) allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// 如果处于打开状态，检查是否可以转为半开
	if cb.state == StateOpen {
		if time.Since(cb.lastFailTime) > cb.timeout {
			cb.state = StateHalfOpen
			cb.failures = 0
			return true
		}
		return false
	}

	// 关闭或半开状态都允许执行
	return true
}

// recordResult 记录执行结果
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()

		// 如果失败次数超过阈值，打开断路器
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		}
	} else {
		// 成功则重置失败计数
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
		}
		cb.failures = 0
	}
}

// State 获取当前状态
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Failures 获取失败次数
func (cb *CircuitBreaker) Failures() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

// CircuitBreakerOpenError 断路器打开错误
type CircuitBreakerOpenError struct {
	Breaker string
}

func (e *CircuitBreakerOpenError) Error() string {
	return "circuit breaker is open for " + e.Breaker
}
