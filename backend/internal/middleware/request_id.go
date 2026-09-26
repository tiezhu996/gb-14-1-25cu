package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "request_id"

// RequestID 请求追踪中间件：为每个请求生成/透传 request_id。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set(requestIDKey, rid)
		c.Writer.Header().Set("X-Request-ID", rid)
		c.Next()
	}
}

// GetRequestID 从上下文读取 request_id。
func GetRequestID(c *gin.Context) string {
	rid, _ := c.Get(requestIDKey)
	s, _ := rid.(string)
	return s
}
