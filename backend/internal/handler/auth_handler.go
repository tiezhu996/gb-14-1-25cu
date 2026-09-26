package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/service"
	"github.com/blueship581/codelearn/internal/util"
)

// AuthHandler 认证处理器：注册/登录。
type AuthHandler struct {
	userService *service.UserService
}

// NewAuthHandler 构造认证处理器。
func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register 注册并返回 JWT。
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	resp, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}

// Login 登录并返回 JWT。
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(util.NewAppError(util.ValidationCode(), util.ValidationMessage()))
		return
	}
	resp, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Success(c, resp)
}
