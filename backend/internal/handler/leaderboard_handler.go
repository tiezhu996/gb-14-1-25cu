package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/service"
	"github.com/blueship581/codelearn/internal/util"
)

// LeaderboardHandler 排行榜处理器。
type LeaderboardHandler struct {
	leaderboardService *service.LeaderboardService
}

// NewLeaderboardHandler 构造排行榜处理器。
func NewLeaderboardHandler(leaderboardService *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{leaderboardService: leaderboardService}
}

// Get 查询排行榜（daily/weekly/total）。
func (h *LeaderboardHandler) Get(c *gin.Context) {
	period := c.DefaultQuery("period", "total")
	limit := int64(50)
	if v := c.Query("limit"); v != "" {
		if n := parseInt64(v); n > 0 && n <= 100 {
			limit = n
		}
	}
	list, err := h.leaderboardService.Get(c.Request.Context(), period, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, list)
}
