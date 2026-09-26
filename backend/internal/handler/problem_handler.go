package handler

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/service"
	"github.com/blueship581/codelearn/internal/util"
)

// ProblemHandler 题目处理器。
type ProblemHandler struct {
	problemService *service.ProblemService
}

// NewProblemHandler 构造题目处理器。
func NewProblemHandler(problemService *service.ProblemService) *ProblemHandler {
	return &ProblemHandler{problemService: problemService}
}

// Create 创建题目（管理员）。
func (h *ProblemHandler) Create(c *gin.Context) {
	userID := getUserID(c)
	username := getUsername(c)
	var req dto.ProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	resp, err := h.problemService.Create(c.Request.Context(), &req, userID, username)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// Update 更新题目。
func (h *ProblemHandler) Update(c *gin.Context) {
	userID := getUserID(c)
	role := getRole(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的题目 ID"))
		return
	}
	var req dto.ProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	resp, err := h.problemService.Update(c.Request.Context(), id, userID, role, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// UpdateStatus 题目状态流转。
func (h *ProblemHandler) UpdateStatus(c *gin.Context) {
	userID := getUserID(c)
	role := getRole(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的题目 ID"))
		return
	}
	var req dto.UpdateProblemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	if err := h.problemService.UpdateStatus(c.Request.Context(), id, userID, role, req.Status); err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessMessage(c, "题目状态已更新")
}

// List 分页查询题目。
func (h *ProblemHandler) List(c *gin.Context) {
	role := getRole(c)
	page, pageSize := parsePage(c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10"))
	difficulty := c.Query("difficulty")
	status := c.Query("status")
	tag := c.Query("tag")
	list, total, err := h.problemService.List(c.Request.Context(), page, pageSize, difficulty, status, tag, role)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessPage(c, list, total, page, pageSize)
}

// Get 获取题目详情。
func (h *ProblemHandler) Get(c *gin.Context) {
	role := getRole(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的题目 ID"))
		return
	}
	resp, err := h.problemService.Get(c.Request.Context(), id, role)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// Delete 删除题目。
func (h *ProblemHandler) Delete(c *gin.Context) {
	userID := getUserID(c)
	role := getRole(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的题目 ID"))
		return
	}
	if err := h.problemService.Delete(c.Request.Context(), id, userID, role); err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessMessage(c, "题目已删除")
}
