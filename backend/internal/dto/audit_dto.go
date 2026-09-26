package dto

import "github.com/blueship581/codelearn/internal/model"

// AuditLogResponse 审计日志响应。
type AuditLogResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	Action     string `json:"action"`
	Entity     string `json:"entity"`
	EntityID   string `json:"entity_id"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	StatusCode int    `json:"status_code"`
	IP         string `json:"ip"`
	RequestID  string `json:"request_id"`
	Detail     string `json:"detail"`
	CreatedAt  string `json:"created_at"`
}

// ToAuditLogResponse 将审计日志转换为响应。
func ToAuditLogResponse(a *model.AuditLog) AuditLogResponse {
	return AuditLogResponse{
		ID:         a.ID.Hex(),
		UserID:     a.UserID.Hex(),
		Username:   a.Username,
		Action:     a.Action,
		Entity:     a.Entity,
		EntityID:   a.EntityID,
		Method:     a.Method,
		Path:       a.Path,
		StatusCode: a.StatusCode,
		IP:         a.IP,
		RequestID:  a.RequestID,
		Detail:     a.Detail,
		CreatedAt:  a.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
