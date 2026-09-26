package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/service"
)

// Audit 操作审计中间件：对写操作（POST/PUT/PATCH/DELETE）异步落审计日志。
func Audit(auditService *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		method := c.Request.Method
		if method == "GET" || method == "OPTIONS" || method == "HEAD" {
			return
		}
		if c.Writer.Status() >= 400 {
			return
		}
		userID, _ := c.Get("user_id")
		uid, _ := userID.(primitive.ObjectID)
		username, _ := c.Get("username")
		name, _ := username.(string)
		action := strings.ToLower(method)
		entity := guessEntity(c.Request.URL.Path)
		log := &model.AuditLog{
			UserID:     uid,
			Username:   name,
			Action:     action,
			Entity:     entity,
			EntityID:   entityIDFromPath(c.Request.URL.Path),
			Method:     method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			IP:         c.ClientIP(),
			RequestID:  GetRequestID(c),
			Detail:     action + " " + entity,
		}
		auditService.Create(c.Request.Context(), log)
	}
}

// guessEntity 从路径猜测操作实体。
func guessEntity(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, part := range parts {
		if part == "api" && i+2 < len(parts) && parts[i+1] == "v1" {
			return parts[i+2]
		}
	}
	return "unknown"
}

// entityIDFromPath 从路径提取实体 ID。
func entityIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := range parts {
		if i > 0 && parts[i-1] == "v1" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
