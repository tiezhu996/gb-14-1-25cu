package router

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/handler"
)

// registerCourseRoutes 注册课程路由。
func registerCourseRoutes(api *gin.RouterGroup, h *handler.CourseHandler, auth, admin, audit gin.HandlerFunc) {
	api.GET("/courses", auth, h.List)
	api.GET("/courses/:id", auth, h.Get)
	api.POST("/courses", auth, admin, audit, h.Create)
	api.PUT("/courses/:id", auth, admin, audit, h.Update)
	api.PUT("/courses/:id/status", auth, admin, audit, h.UpdateStatus)
	api.DELETE("/courses/:id", auth, admin, audit, h.Delete)
	api.POST("/courses/:id/learn", auth, audit, h.RecordLearn)
	api.POST("/courses/:id/complete", auth, audit, h.CompleteCourse)
	_ = constants.RoleAdmin
}
