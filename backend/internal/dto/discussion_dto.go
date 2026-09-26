package dto

import "github.com/blueship581/codelearn/internal/model"

// CreateDiscussionRequest 发布讨论帖请求。
type CreateDiscussionRequest struct {
	Title    string `json:"title" binding:"required,max=128"`
	Content  string `json:"content" binding:"required,max=20000"`
	Code     string `json:"code" binding:"omitempty,max=20000"`
	Language string `json:"language" binding:"omitempty,oneof=python javascript java"`
}

// DiscussionResponse 讨论帖响应。
type DiscussionResponse struct {
	ID        string `json:"id"`
	ProblemID string `json:"problem_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Code      string `json:"code"`
	Language  string `json:"language"`
	IsBest    bool   `json:"is_best"`
	VoteCount int    `json:"vote_count"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// ToDiscussionResponse 将模型转换为响应。
func ToDiscussionResponse(d *model.Discussion) DiscussionResponse {
	return DiscussionResponse{
		ID:        d.ID.Hex(),
		ProblemID: d.ProblemID.Hex(),
		UserID:    d.UserID.Hex(),
		Username:  d.Username,
		Nickname:  d.Nickname,
		Title:     d.Title,
		Content:   d.Content,
		Code:      d.Code,
		Language:  d.Language,
		IsBest:    d.IsBest,
		VoteCount: d.VoteCount,
		Status:    d.Status,
		CreatedAt: d.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
