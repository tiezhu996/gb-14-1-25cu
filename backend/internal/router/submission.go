package router

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/handler"
)

// registerSubmissionRoutes 注册提交评测路由。
func registerSubmissionRoutes(api *gin.RouterGroup, h *handler.SubmissionHandler, auth, admin, audit gin.HandlerFunc) {
	api.GET("/submissions", auth, h.List)
	api.GET("/submissions/:id", auth, h.Get)
}
