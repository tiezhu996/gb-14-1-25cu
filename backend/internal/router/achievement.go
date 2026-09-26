package router

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/handler"
)

// registerAchievementRoutes 注册成就路由。
func registerAchievementRoutes(api *gin.RouterGroup, h *handler.AchievementHandler, auth gin.HandlerFunc) {
	api.GET("/achievements", auth, h.ListAll)
	api.GET("/achievements/me", auth, h.ListMine)
}
