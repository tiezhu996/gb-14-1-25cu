package dto

import "github.com/blueship581/codelearn/internal/model"

// AchievementResponse 成就徽章响应。
type AchievementResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	EarnedAt    string `json:"earned_at,omitempty"`
}

// ToAchievementResponse 将成就定义转换为响应。
func ToAchievementResponse(a *model.Achievement) AchievementResponse {
	return AchievementResponse{
		ID:          a.ID.Hex(),
		Code:        a.Code,
		Name:        a.Name,
		Description: a.Description,
		Icon:        a.Icon,
	}
}

// ToUserAchievementResponse 将用户成就记录转换为响应。
func ToUserAchievementResponse(ua *model.UserAchievement) AchievementResponse {
	return AchievementResponse{
		ID:          ua.ID.Hex(),
		Code:        ua.Code,
		Name:        ua.Name,
		Description: ua.Description,
		Icon:        ua.Icon,
		EarnedAt:    ua.EarnedAt.Format("2006-01-02 15:04:05"),
	}
}
