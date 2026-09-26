package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/service"
	"github.com/blueship581/codelearn/internal/util"
)

// DashboardHandler 个人学习仪表盘处理器。
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler 构造仪表盘处理器。
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// Get 获取当前用户学习统计。
func (h *DashboardHandler) Get(c *gin.Context) {
	userID := getUserID(c)
	resp, err := h.dashboardService.Get(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}
