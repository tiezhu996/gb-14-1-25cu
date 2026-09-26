package router

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/handler"
)

// registerLeaderboardRoutes 注册排行榜路由。
func registerLeaderboardRoutes(api *gin.RouterGroup, h *handler.LeaderboardHandler, auth gin.HandlerFunc) {
	api.GET("/leaderboard", auth, h.Get)
}
