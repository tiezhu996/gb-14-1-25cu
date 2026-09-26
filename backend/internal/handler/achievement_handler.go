package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/service"
	"github.com/blueship581/codelearn/internal/util"
)

// AchievementHandler 成就处理器。
type AchievementHandler struct {
	achievementService *service.AchievementService
}

// NewAchievementHandler 构造成就处理器。
func NewAchievementHandler(achievementService *service.AchievementService) *AchievementHandler {
	return &AchievementHandler{achievementService: achievementService}
}

// ListAll 查询全部成就定义。
func (h *AchievementHandler) ListAll(c *gin.Context) {
	list, err := h.achievementService.ListAll(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, list)
}

// ListMine 查询我的成就。
func (h *AchievementHandler) ListMine(c *gin.Context) {
	userID := getUserID(c)
	list, err := h.achievementService.ListMine(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, list)
}
