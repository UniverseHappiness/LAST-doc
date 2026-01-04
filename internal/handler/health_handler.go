package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/UniverseHappiness/LAST-doc/internal/service"
	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	healthService *service.HealthService
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// CheckHealth 健康检查端点
// @Summary 系统健康检查
// @Description 检查系统各组件的健康状态
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) CheckHealth(c *gin.Context) {
	ctx, cancel := contextWithTimeout(c, 5*time.Second)
	defer cancel()

	health := h.healthService.CheckHealth(ctx)

	c.JSON(http.StatusOK, gin.H{
		"status":     health.Status,
		"timestamp":  health.Timestamp,
		"uptime":     health.Uptime.String(),
		"components": health.Components,
	})
}

// LivenessProbe 存活探针
// @Summary Kubernetes存活探针
// @Description 用于Kubernetes的存活探针检查
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/live [get]
func (h *HealthHandler) LivenessProbe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// ReadinessProbe 就绪探针
// @Summary Kubernetes就绪探针
// @Description 用于Kubernetes的就绪探针检查
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]string
// @Router /health/ready [get]
func (h *HealthHandler) ReadinessProbe(c *gin.Context) {
	ctx, cancel := contextWithTimeout(c, 3*time.Second)
	defer cancel()

	health := h.healthService.CheckHealth(ctx)

	if health.Status == service.HealthStatusHealthy {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": health.Timestamp,
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not_ready",
			"message": "Service is not ready to accept traffic",
		})
	}
}

// CircuitBreakers 获取断路器状态
// @Summary 获取断路器状态
// @Description 获取所有断路器的当前状态
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health/circuit-breakers [get]
func (h *HealthHandler) CircuitBreakers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Circuit breaker status",
	})
}

func contextWithTimeout(c *gin.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), timeout)
}
