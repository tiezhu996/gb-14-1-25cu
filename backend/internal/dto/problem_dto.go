package dto

import "github.com/blueship581/codelearn/internal/model"

// TestCaseDTO 测试用例 DTO。
type TestCaseDTO struct {
	Input  string `json:"input" binding:"required"`
	Output string `json:"output" binding:"required"`
}

// ProblemRequest 创建/更新题目请求。
type ProblemRequest struct {
	Title       string        `json:"title" binding:"required,max=128"`
	Description string        `json:"description" binding:"required"`
	Difficulty  string        `json:"difficulty" binding:"required,oneof=easy medium hard"`
	Languages   []string      `json:"languages" binding:"required,min=1,dive,oneof=python javascript java"`
	Tags        []string      `json:"tags" binding:"omitempty,dive,max=32"`
	TestCases   []TestCaseDTO `json:"test_cases" binding:"required,min=1,dive"`
	TimeLimit   int           `json:"time_limit" binding:"omitempty,min=1,max=10"`
	Status      string        `json:"status" binding:"omitempty,oneof=draft published archived"`
}

// UpdateProblemStatusRequest 题目状态流转请求。
type UpdateProblemStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=draft published archived"`
}

// ProblemResponse 题目响应（学生视角隐藏期望输出）。
type ProblemResponse struct {
	ID             string        `json:"id"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Difficulty     string        `json:"difficulty"`
	Languages      []string      `json:"languages"`
	Tags           []string      `json:"tags"`
	TestCases      []TestCaseDTO `json:"test_cases,omitempty"`
	TimeLimit      int           `json:"time_limit"`
	Status         string        `json:"status"`
	Points         int           `json:"points"`
	AcceptedCount  int64         `json:"accepted_count"`
	SubmitCount    int64         `json:"submit_count"`
	CreatedByName  string        `json:"created_by_name"`
	CreatedAt      string        `json:"created_at"`
	UpdatedAt      string        `json:"updated_at"`
}

// ToProblemResponse 将模型转换为响应；withOutput=true 时返回期望输出（仅管理员）。
func ToProblemResponse(p *model.Problem, withOutput bool) ProblemResponse {
	tcs := make([]TestCaseDTO, 0, len(p.TestCases))
	for _, tc := range p.TestCases {
		out := ""
		if withOutput {
			out = tc.Output
		}
		tcs = append(tcs, TestCaseDTO{Input: tc.Input, Output: out})
	}
	return ProblemResponse{
		ID:            p.ID.Hex(),
		Title:         p.Title,
		Description:   p.Description,
		Difficulty:    p.Difficulty,
		Languages:     p.Languages,
		Tags:          p.Tags,
		TestCases:     tcs,
		TimeLimit:     p.TimeLimit,
		Status:        p.Status,
		Points:        p.Points,
		AcceptedCount: p.AcceptedCount,
		SubmitCount:   p.SubmitCount,
		CreatedByName: p.CreatedByName,
		CreatedAt:     p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
