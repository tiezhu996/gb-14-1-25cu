package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/blueship581/codelearn/internal/service"
	"github.com/blueship581/codelearn/internal/util"
)

// AuditHandler 审计日志处理器。
type AuditHandler struct {
	auditService *service.AuditService
}

// NewAuditHandler 构造审计处理器。
func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// List 分页查询审计日志（管理员）。
func (h *AuditHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c.DefaultQuery("page", "1"), c.DefaultQuery("page_size", "10"))
	username := c.Query("username")
	entity := c.Query("entity")
	list, total, err := h.auditService.List(c.Request.Context(), page, pageSize, username, entity)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.SuccessPage(c, list, total, page, pageSize)
}
