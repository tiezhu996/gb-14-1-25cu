package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/service"
	"github.com/blueship581/codelearn/internal/util"
)

// UserHandler 用户处理器。
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 构造用户处理器。
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMe 获取当前登录用户资料。
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := getUserID(c)
	resp, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// UpdateProfile 更新个人资料。
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := getUserID(c)
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	resp, err := h.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// SignIn 每日签到。
func (h *UserHandler) SignIn(c *gin.Context) {
	userID := getUserID(c)
	resp, err := h.userService.SignIn(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	msg := "签到成功"
	if resp.StreakDays > 0 {
		msg = "签到成功，连续签到 " + strconv.Itoa(resp.StreakDays) + " 天"
	}
	util.SuccessMessage(c, msg)
}

// List 管理员分页查询用户。
func (h *UserHandler) List(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	p, ps := parsePage(page, pageSize)
	role := c.Query("role")
	status := c.Query("status")
	list, total, err := h.userService.List(c.Request.Context(), p, ps, role, status)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessPage(c, list, total, p, ps)
}

// UpdateStatus 管理员更新用户状态。
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	adminID := getUserID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的用户 ID"))
		return
	}
	var req dto.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	if err := h.userService.UpdateStatus(c.Request.Context(), adminID, id, req.Status); err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessMessage(c, constants.MsgUserStatusChanged)
}

// UpdateRole 管理员更新用户角色。
func (h *UserHandler) UpdateRole(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的用户 ID"))
		return
	}
	var req struct {
		Role string `json:"role" binding:"required,oneof=student admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	if err := h.userService.UpdateRole(c.Request.Context(), id, req.Role); err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessMessage(c, constants.MsgUserRoleChanged)
}

// getUserID 从上下文读取当前用户 ID。
func getUserID(c *gin.Context) primitive.ObjectID {
	v, _ := c.Get("user_id")
	id, _ := v.(primitive.ObjectID)
	return id
}

// parsePage 解析分页参数。
func parsePage(page, pageSize string) (int64, int64) {
	p := int64(1)
	ps := int64(constants.DefaultPageSize)
	if v := parseInt64(page); v > 0 {
		p = v
	}
	if v := parseInt64(pageSize); v > 0 {
		ps = v
	}
	if ps > constants.MaxPageSize {
		ps = constants.MaxPageSize
	}
	return p, ps
}

func parseInt64(s string) int64 {
	var n int64
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0
		}
		n = n*10 + int64(ch-'0')
	}
	return n
}
