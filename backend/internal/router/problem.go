package router

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/handler"
)

// registerProblemRoutes 注册题目路由。
func registerProblemRoutes(api *gin.RouterGroup, ph *handler.ProblemHandler, sh *handler.SubmissionHandler,
	dh *handler.DiscussionHandler, auth, admin, audit gin.HandlerFunc) {
	api.GET("/problems", auth, ph.List)
	api.GET("/problems/:id", auth, ph.Get)
	api.POST("/problems", auth, admin, audit, ph.Create)
	api.PUT("/problems/:id", auth, admin, audit, ph.Update)
	api.PUT("/problems/:id/status", auth, admin, audit, ph.UpdateStatus)
	api.DELETE("/problems/:id", auth, admin, audit, ph.Delete)
	// 提交评测
	api.POST("/problems/:id/submit", auth, audit, sh.Submit)
	// 讨论区
	api.GET("/problems/:id/discussions", auth, dh.ListByProblem)
	api.POST("/problems/:id/discussions", auth, audit, dh.Create)
	api.PUT("/problems/:id/discussions/:discussion_id/best", auth, audit, dh.MarkBest)
}
