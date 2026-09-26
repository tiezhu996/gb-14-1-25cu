package handler

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/service"
	"github.com/blueship581/codelearn/internal/util"
)

// SubmissionHandler 提交评测处理器。
type SubmissionHandler struct {
	submissionService *service.SubmissionService
}

// NewSubmissionHandler 构造提交评测处理器。
func NewSubmissionHandler(submissionService *service.SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{submissionService: submissionService}
}

// Submit 提交代码并评测。
func (h *SubmissionHandler) Submit(c *gin.Context) {
	userID := getUserID(c)
	problemID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的题目 ID"))
		return
	}
	var req dto.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	resp, err := h.submissionService.Submit(c.Request.Context(), userID, problemID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// Get 查询提交详情。
func (h *SubmissionHandler) Get(c *gin.Context) {
	userID := getUserID(c)
	role := getRole(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的提交 ID"))
		return
	}
	resp, err := h.submissionService.Get(c.Request.Context(), id, userID, role)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// List 分页查询提交记录。
func (h *SubmissionHandler) List(c *gin.Context) {
	userID := getUserID(c)
	role := getRole(c)
	page, pageSize := parsePage(c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10"))
	problemID := c.Query("problem_id")
	status := c.Query("status")
	list, total, err := h.submissionService.List(c.Request.Context(), userID, role, page, pageSize, problemID, status)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessPage(c, list, total, page, pageSize)
}
