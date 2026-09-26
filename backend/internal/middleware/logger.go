package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/constants"
)

// Logger 请求日志中间件：输出 request_id/method/path/status/latency_ms。
func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		logger.Info(constants.LogRequest,
			"request_id", GetRequestID(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", latency,
			"client_ip", c.ClientIP(),
		)
	}
}
