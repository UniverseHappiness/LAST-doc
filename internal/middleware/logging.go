package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/UniverseHappiness/LAST-doc/internal/model"
	"github.com/UniverseHappiness/LAST-doc/internal/service"
	"github.com/gin-gonic/gin"
)

// LoggingMiddleware 日志记录中间件
type LoggingMiddleware struct {
	monitorService service.MonitorService
}

// NewLoggingMiddleware 创建日志记录中间件实例
func NewLoggingMiddleware(monitorService service.MonitorService) *LoggingMiddleware {
	return &LoggingMiddleware{
		monitorService: monitorService,
	}
}

// LogRequest 记录HTTP请求
func (m *LoggingMiddleware) LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算请求耗时
		latency := time.Since(start).Milliseconds()

		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		// 获取用户ID（从context中）
		userID, _ := c.Get("user_id")

		// 确定日志级别
		level := "info"
		if statusCode >= 500 {
			level = "error"
		} else if statusCode >= 400 {
			level = "warn"
		}

		// 获取用户名（从context中）
		username, _ := c.Get("username")
		var usernameStr string
		if username != nil {
			usernameStr = fmt.Sprintf("%v", username)
		}

		// 创建日志条目 - 使用本地时间（Asia/Shanghai）
		// 注意：Go的time.Now()返回UTC时间，需要转换为本地时间
		now := time.Now()
		// 使用In方法指定时区为Asia/Shanghai
		shanghaiZone, err := time.LoadLocation("Asia/Shanghai")
		if err == nil {
			now = now.In(shanghaiZone)
		}

		logEntry := &model.LogEntry{
			Timestamp:  now,
			Level:      level,
			Service:    "api",
			Message:    fmt.Sprintf("%s %s - %d", method, path, statusCode),
			Method:     method,
			Path:       path,
			UserID:     fmt.Sprintf("%v", userID),
			Username:   usernameStr,
			StatusCode: statusCode,
			Latency:    latency,
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}

		// 如果响应中有错误，记录错误信息
		if len(c.Errors) > 0 {
			errMsg := c.Errors.String()
			logEntry.Error = errMsg
			logEntry.Level = "error"
		}

		// 异步记录日志（不阻塞请求）
		go func() {
			// 使用新的context，避免请求context被取消
			ctx := context.Background()
			if err := m.monitorService.LogRequest(ctx, logEntry); err != nil {
				// 记录失败，打印到控制台
				fmt.Printf("Failed to log request: %v\n", err)
			}
		}()
	}
}
