package router

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/handler"
)

// registerDiscussionRoutes 注册讨论路由。
func registerDiscussionRoutes(api *gin.RouterGroup, h *handler.DiscussionHandler, auth, admin, audit gin.HandlerFunc) {
	api.POST("/discussions/:id/vote", auth, audit, h.Vote)
	api.PUT("/discussions/:id/hide", auth, audit, h.Hide)
}
