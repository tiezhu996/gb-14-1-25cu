package dto

import "github.com/blueship581/codelearn/internal/model"

// ChapterDTO 章节 DTO。
type ChapterDTO struct {
	Title    string `json:"title" binding:"required,max=128"`
	Content  string `json:"content" binding:"required"`
	Duration int    `json:"duration" binding:"required,min=1,max=600"`
}

// CourseRequest 创建/更新课程请求。
type CourseRequest struct {
	Title       string       `json:"title" binding:"required,max=128"`
	Description string       `json:"description" binding:"required,max=512"`
	Cover       string       `json:"cover" binding:"omitempty,max=512"`
	Markdown    string       `json:"markdown" binding:"omitempty"`
	Chapters    []ChapterDTO `json:"chapters" binding:"omitempty,dive"`
	Difficulty  string       `json:"difficulty" binding:"required,oneof=easy medium hard"`
	Status      string       `json:"status" binding:"omitempty,oneof=draft published archived"`
}

// UpdateCourseStatusRequest 课程状态流转请求。
type UpdateCourseStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=draft published archived"`
}

// LearnRecordRequest 学习时长记录请求。
type LearnRecordRequest struct {
	DurationMinutes int `json:"duration_minutes" binding:"required,min=1,max=600"`
}

// CourseResponse 课程响应。
type CourseResponse struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Cover       string        `json:"cover"`
	Markdown    string        `json:"markdown"`
	Chapters    []ChapterDTO  `json:"chapters"`
	Difficulty  string        `json:"difficulty"`
	Status      string        `json:"status"`
	AuthorID    string        `json:"author_id"`
	AuthorName  string        `json:"author_name"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

// ToCourseResponse 将模型转换为响应。
func ToCourseResponse(c *model.Course) CourseResponse {
	chapters := make([]ChapterDTO, 0, len(c.Chapters))
	for _, ch := range c.Chapters {
		chapters = append(chapters, ChapterDTO{Title: ch.Title, Content: ch.Content, Duration: ch.Duration})
	}
	return CourseResponse{
		ID:          c.ID.Hex(),
		Title:       c.Title,
		Description: c.Description,
		Cover:       c.Cover,
		Markdown:    c.Markdown,
		Chapters:    chapters,
		Difficulty:  c.Difficulty,
		Status:      c.Status,
		AuthorID:    c.AuthorID.Hex(),
		AuthorName:  c.AuthorName,
		CreatedAt:   c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
