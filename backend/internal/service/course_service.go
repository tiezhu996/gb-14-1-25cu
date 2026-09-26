package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
)

// CourseService 课程业务：CRUD + 状态机（draft->published->archived）+ 学习记录。
type CourseService struct {
	repo   *repository.CourseRepository
	stat   *repository.UserStatRepository
	logger *slog.Logger
}

// NewCourseService 构造课程服务。
func NewCourseService(repo *repository.CourseRepository, stat *repository.UserStatRepository, logger *slog.Logger) *CourseService {
	return &CourseService{repo: repo, stat: stat, logger: logger}
}

// toModel 将请求转换为课程模型。
func (s *CourseService) toModel(req *dto.CourseRequest, authorID primitive.ObjectID, authorName string) *model.Course {
	chapters := make([]model.Chapter, 0, len(req.Chapters))
	for _, ch := range req.Chapters {
		chapters = append(chapters, model.Chapter{Title: ch.Title, Content: ch.Content, Duration: ch.Duration})
	}
	status := req.Status
	if status == "" {
		status = constants.StatusDraft
	}
	return &model.Course{
		Title:       req.Title,
		Description: req.Description,
		Cover:       req.Cover,
		Markdown:    req.Markdown,
		Chapters:    chapters,
		Difficulty:  req.Difficulty,
		Status:      status,
		AuthorID:    authorID,
		AuthorName:  authorName,
	}
}

// Create 创建课程（管理员）。
func (s *CourseService) Create(ctx context.Context, req *dto.CourseRequest, authorID primitive.ObjectID, authorName string) (*dto.CourseResponse, error) {
	course := s.toModel(req, authorID, authorName)
	if err := s.repo.Create(ctx, course); err != nil {
		if errors.Is(err, repository.ErrCourseExists) {
			return nil, util.WrapAppError(constants.CodeCourseExists, constants.MsgCourseExists, err)
		}
		s.logger.Error(constants.LogCourseCreatedFailed, "title", req.Title, "error", err.Error())
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogCourseCreated, "course_id", course.ID.Hex(), "title", course.Title)
	resp := dto.ToCourseResponse(course)
	return &resp, nil
}

// Update 更新课程（作者或管理员）。
func (s *CourseService) Update(ctx context.Context, courseID, operatorID primitive.ObjectID, operatorRole string, req *dto.CourseRequest) (*dto.CourseResponse, error) {
	course, err := s.repo.FindByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, repository.ErrCourseNotFound) {
			return nil, util.WrapAppError(constants.CodeCourseNotFound, constants.MsgCourseNotFound, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if course.AuthorID != operatorID && operatorRole != constants.RoleAdmin {
		return nil, util.WrapAppError(constants.CodeCourseOwner, constants.MsgCourseOwner, nil)
	}
	chapters := make([]model.Chapter, 0, len(req.Chapters))
	for _, ch := range req.Chapters {
		chapters = append(chapters, model.Chapter{Title: ch.Title, Content: ch.Content, Duration: ch.Duration})
	}
	update := bson.M{
		"title":       req.Title,
		"description": req.Description,
		"cover":       req.Cover,
		"markdown":    req.Markdown,
		"chapters":    chapters,
		"difficulty":  req.Difficulty,
	}
	if err := s.repo.Update(ctx, courseID, update); err != nil {
		if errors.Is(err, repository.ErrCourseExists) {
			return nil, util.WrapAppError(constants.CodeCourseExists, constants.MsgCourseExists, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogCourseUpdated, "course_id", courseID.Hex(), "title", req.Title)
	return s.Get(ctx, courseID, operatorRole)
}

// UpdateStatus 课程状态机流转：draft->published->archived（作者或管理员）。
func (s *CourseService) UpdateStatus(ctx context.Context, courseID, operatorID primitive.ObjectID, operatorRole, status string) error {
	course, err := s.repo.FindByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, repository.ErrCourseNotFound) {
			return util.WrapAppError(constants.CodeCourseNotFound, constants.MsgCourseNotFound, err)
		}
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if course.AuthorID != operatorID && operatorRole != constants.RoleAdmin {
		return util.WrapAppError(constants.CodeCourseOwner, constants.MsgCourseOwner, nil)
	}
	if !canTransition(course.Status, status) {
		return util.WrapAppError(constants.CodeCourseLocked, constants.MsgCourseLocked, nil)
	}
	if err := s.repo.UpdateStatus(ctx, courseID, status); err != nil {
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if status == constants.StatusPublished {
		s.logger.Info(constants.LogCoursePublish, "course_id", courseID.Hex())
	} else if status == constants.StatusArchived {
		s.logger.Info(constants.LogCourseArchived, "course_id", courseID.Hex())
	}
	return nil
}

// List 分页查询课程（学生只看已发布）。
func (s *CourseService) List(ctx context.Context, page, pageSize int64, difficulty, status, role string) ([]dto.CourseResponse, int64, error) {
	filter := bson.M{}
	if difficulty != "" {
		filter["difficulty"] = difficulty
	}
	if status != "" {
		filter["status"] = status
	} else if role != constants.RoleAdmin {
		filter["status"] = constants.StatusPublished
	}
	courses, total, err := s.repo.List(ctx, filter, (page-1)*pageSize, pageSize)
	if err != nil {
		s.logger.Error(constants.LogCourseListFailed, "error", err.Error())
		return nil, 0, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	out := make([]dto.CourseResponse, 0, len(courses))
	for _, c := range courses {
		out = append(out, dto.ToCourseResponse(c))
	}
	return out, total, nil
}

// Get 获取课程详情。
func (s *CourseService) Get(ctx context.Context, courseID primitive.ObjectID, role string) (*dto.CourseResponse, error) {
	course, err := s.repo.FindByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, repository.ErrCourseNotFound) {
			return nil, util.WrapAppError(constants.CodeCourseNotFound, constants.MsgCourseNotFound, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if course.Status != constants.StatusPublished && role != constants.RoleAdmin && course.Status != constants.StatusDraft {
		return nil, util.WrapAppError(constants.CodeCourseLocked, constants.MsgCourseLocked, nil)
	}
	resp := dto.ToCourseResponse(course)
	return &resp, nil
}

// Delete 删除课程（作者或管理员）。
func (s *CourseService) Delete(ctx context.Context, courseID, operatorID primitive.ObjectID, operatorRole string) error {
	course, err := s.repo.FindByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, repository.ErrCourseNotFound) {
			return util.WrapAppError(constants.CodeCourseNotFound, constants.MsgCourseNotFound, err)
		}
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if course.AuthorID != operatorID && operatorRole != constants.RoleAdmin {
		return util.WrapAppError(constants.CodeCourseOwner, constants.MsgCourseOwner, nil)
	}
	if err := s.repo.Delete(ctx, courseID); err != nil {
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	return nil
}

// RecordLearn 记录学习时长（原子累加 + 日活跃热力图）。
func (s *CourseService) RecordLearn(ctx context.Context, userID primitive.ObjectID, minutes int) error {
	dayKey := util.SignInDailyKey(time.Now())
	if err := s.stat.AddLearning(ctx, userID, int64(minutes), dayKey); err != nil {
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogCourseLearnRecorded, "user_id", userID.Hex(), "minutes", minutes)
	return nil
}

// CompleteCourse 标记课程完成（原子累加完成数）。
func (s *CourseService) CompleteCourse(ctx context.Context, userID primitive.ObjectID) error {
	dayKey := util.SignInDailyKey(time.Now())
	if err := s.stat.CompleteCourse(ctx, userID, dayKey); err != nil {
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogCourseCompleted, "user_id", userID.Hex())
	return nil
}

// canTransition 内容状态机：草稿->已发布->已归档；已发布可回草稿；已归档终态。
func canTransition(from, to string) bool {
	if from == to {
		return true
	}
	switch from {
	case constants.StatusDraft:
		return to == constants.StatusPublished
	case constants.StatusPublished:
		return to == constants.StatusArchived || to == constants.StatusDraft
	case constants.StatusArchived:
		return false
	}
	return false
}
