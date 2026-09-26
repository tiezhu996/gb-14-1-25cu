package dto

import (
	"github.com/blueship581/codelearn/internal/model"
)

// RegisterRequest 注册请求。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Nickname string `json:"nickname" binding:"omitempty,max=32"`
}

// LoginRequest 登录请求。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest 更新个人资料请求。
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,max=32"`
	Avatar   string `json:"avatar" binding:"omitempty,max=512"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// UpdateUserStatusRequest 管理员更新用户状态请求。
type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active banned"`
}

// UserResponse 用户响应。
type UserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	Points      int64  `json:"points"`
	SolvedCount int64  `json:"solved_count"`
	StreakDays  int    `json:"streak_days"`
	CreatedAt   string `json:"created_at"`
}

// LoginResponse 登录响应（含 JWT）。
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// SignInResponse 签到响应。
type SignInResponse struct {
	StreakDays int  `json:"streak_days"`
	FirstToday bool `json:"first_today"`
}

// ToUserResponse 将模型转换为响应。
func ToUserResponse(u *model.User) UserResponse {
	createdAt := ""
	if !u.CreatedAt.IsZero() {
		createdAt = u.CreatedAt.Format("2006-01-02 15:04:05")
	}
	return UserResponse{
		ID:          u.ID.Hex(),
		Username:    u.Username,
		Email:       u.Email,
		Nickname:    u.Nickname,
		Avatar:      u.Avatar,
		Role:        u.Role,
		Status:      u.Status,
		Points:      u.Points,
		SolvedCount: u.SolvedCount,
		StreakDays:  u.StreakDays,
		CreatedAt:   createdAt,
	}
}
