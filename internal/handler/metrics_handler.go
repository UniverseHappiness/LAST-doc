package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP 请求总数
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP 请求持续时间
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	// 注册 metrics
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// MetricsHandler metrics处理器
type MetricsHandler struct{}

// NewMetricsHandler 创建metrics处理器
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// ServeMetrics 暴露 metrics 端点
func (h *MetricsHandler) ServeMetrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

// RecordRequest 记录请求指标
func RecordRequest(method, path string, status int, duration float64) {
	httpRequestsTotal.WithLabelValues(method, path, string(rune(status))).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration)
}
