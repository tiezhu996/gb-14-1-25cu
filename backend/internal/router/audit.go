package router

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/handler"
)

// registerAuditRoutes 注册审计日志路由（管理员）。
func registerAuditRoutes(api *gin.RouterGroup, h *handler.AuditHandler, auth, admin gin.HandlerFunc) {
	api.GET("/audits", auth, admin, h.List)
}
