package handler

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/service"
	"github.com/blueship581/codelearn/internal/util"
)

// DiscussionHandler 讨论社区处理器。
type DiscussionHandler struct {
	discussionService *service.DiscussionService
}

// NewDiscussionHandler 构造讨论处理器。
func NewDiscussionHandler(discussionService *service.DiscussionService) *DiscussionHandler {
	return &DiscussionHandler{discussionService: discussionService}
}

// Create 发布讨论帖。
func (h *DiscussionHandler) Create(c *gin.Context) {
	userID := getUserID(c)
	problemID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的题目 ID"))
		return
	}
	var req dto.CreateDiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	resp, err := h.discussionService.Create(c.Request.Context(), userID, problemID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// ListByProblem 查询题目讨论区。
func (h *DiscussionHandler) ListByProblem(c *gin.Context) {
	problemID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的题目 ID"))
		return
	}
	sort := c.DefaultQuery("sort", "best")
	list, err := h.discussionService.ListByProblem(c.Request.Context(), problemID, sort)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, list)
}

// Vote 点赞/取消点赞。
func (h *DiscussionHandler) Vote(c *gin.Context) {
	userID := getUserID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的讨论 ID"))
		return
	}
	resp, err := h.discussionService.Vote(c.Request.Context(), userID, id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// MarkBest 标记最佳答案（题目作者或管理员）。
func (h *DiscussionHandler) MarkBest(c *gin.Context) {
	userID := getUserID(c)
	role := getRole(c)
	problemID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的题目 ID"))
		return
	}
	discussionID, err := primitive.ObjectIDFromHex(c.Param("discussion_id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的讨论 ID"))
		return
	}
	if err := h.discussionService.MarkBest(c.Request.Context(), userID, role, problemID, discussionID); err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessMessage(c, "已设为最佳答案")
}

// Hide 隐藏讨论帖。
func (h *DiscussionHandler) Hide(c *gin.Context) {
	userID := getUserID(c)
	role := getRole(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的讨论 ID"))
		return
	}
	if err := h.discussionService.Hide(c.Request.Context(), userID, role, id); err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessMessage(c, "讨论帖已隐藏")
}
