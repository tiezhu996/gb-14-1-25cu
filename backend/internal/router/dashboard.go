package router

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/handler"
)

// registerDashboardRoutes 注册仪表盘路由。
func registerDashboardRoutes(api *gin.RouterGroup, h *handler.DashboardHandler, auth gin.HandlerFunc) {
	api.GET("/dashboard/me", auth, h.Get)
}
