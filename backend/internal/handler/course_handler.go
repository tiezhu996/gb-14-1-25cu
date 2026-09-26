package handler

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/service"
	"github.com/blueship581/codelearn/internal/util"
)

// CourseHandler 课程处理器。
type CourseHandler struct {
	courseService *service.CourseService
}

// NewCourseHandler 构造课程处理器。
func NewCourseHandler(courseService *service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

// Create 创建课程（管理员）。
func (h *CourseHandler) Create(c *gin.Context) {
	userID := getUserID(c)
	username := getUsername(c)
	var req dto.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	resp, err := h.courseService.Create(c.Request.Context(), &req, userID, username)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// Update 更新课程。
func (h *CourseHandler) Update(c *gin.Context) {
	userID := getUserID(c)
	role := getRole(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的课程 ID"))
		return
	}
	var req dto.CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	resp, err := h.courseService.Update(c.Request.Context(), id, userID, role, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// UpdateStatus 课程状态流转。
func (h *CourseHandler) UpdateStatus(c *gin.Context) {
	userID := getUserID(c)
	role := getRole(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的课程 ID"))
		return
	}
	var req dto.UpdateCourseStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	if err := h.courseService.UpdateStatus(c.Request.Context(), id, userID, role, req.Status); err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessMessage(c, "课程状态已更新")
}

// List 分页查询课程。
func (h *CourseHandler) List(c *gin.Context) {
	role := getRole(c)
	page, pageSize := parsePage(c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10"))
	difficulty := c.Query("difficulty")
	status := c.Query("status")
	list, total, err := h.courseService.List(c.Request.Context(), page, pageSize, difficulty, status, role)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessPage(c, list, total, page, pageSize)
}

// Get 获取课程详情。
func (h *CourseHandler) Get(c *gin.Context) {
	role := getRole(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的课程 ID"))
		return
	}
	resp, err := h.courseService.Get(c.Request.Context(), id, role)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// Delete 删除课程。
func (h *CourseHandler) Delete(c *gin.Context) {
	userID := getUserID(c)
	role := getRole(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		_ = c.Error(util.NewAppError(constants.CodeBadRequest, "无效的课程 ID"))
		return
	}
	if err := h.courseService.Delete(c.Request.Context(), id, userID, role); err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessMessage(c, "课程已删除")
}

// RecordLearn 记录学习时长。
func (h *CourseHandler) RecordLearn(c *gin.Context) {
	userID := getUserID(c)
	var req dto.LearnRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	if err := h.courseService.RecordLearn(c.Request.Context(), userID, req.DurationMinutes); err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessMessage(c, "学习时长已记录")
}

// CompleteCourse 标记课程完成。
func (h *CourseHandler) CompleteCourse(c *gin.Context) {
	userID := getUserID(c)
	if err := h.courseService.CompleteCourse(c.Request.Context(), userID); err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessMessage(c, "课程已标记完成")
}

// getRole 从上下文读取当前用户角色。
func getRole(c *gin.Context) string {
	v, _ := c.Get("role")
	role, _ := v.(string)
	return role
}

// getUsername 从上下文读取当前用户名。
func getUsername(c *gin.Context) string {
	v, _ := c.Get("username")
	name, _ := v.(string)
	return name
}
