package service

import (
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
)

// ProblemService 题目业务：CRUD + 状态机 + 难度积分。
type ProblemService struct {
	repo   *repository.ProblemRepository
	logger *slog.Logger
}

// NewProblemService 构造题目服务。
func NewProblemService(repo *repository.ProblemRepository, logger *slog.Logger) *ProblemService {
	return &ProblemService{repo: repo, logger: logger}
}

// toModel 将请求转换为题目模型。
func (s *ProblemService) toModel(req *dto.ProblemRequest, creatorID primitive.ObjectID, creatorName string) *model.Problem {
	tcs := make([]model.TestCase, 0, len(req.TestCases))
	for _, tc := range req.TestCases {
		tcs = append(tcs, model.TestCase{Input: tc.Input, Output: tc.Output})
	}
	timeLimit := req.TimeLimit
	if timeLimit <= 0 {
		timeLimit = 10
	}
	status := req.Status
	if status == "" {
		status = constants.StatusDraft
	}
	return &model.Problem{
		Title:        req.Title,
		Description:  req.Description,
		Difficulty:   req.Difficulty,
		Languages:    req.Languages,
		Tags:         req.Tags,
		TestCases:    tcs,
		TimeLimit:    timeLimit,
		Status:       status,
		Points:       constants.DifficultyPoints[req.Difficulty],
		CreatedBy:    creatorID,
		CreatedByName: creatorName,
	}
}

// Create 创建题目（管理员）。
func (s *ProblemService) Create(ctx context.Context, req *dto.ProblemRequest, creatorID primitive.ObjectID, creatorName string) (*dto.ProblemResponse, error) {
	if len(req.TestCases) == 0 {
		return nil, util.WrapAppError(constants.CodeProblemNoCases, constants.MsgProblemNoCases, nil)
	}
	problem := s.toModel(req, creatorID, creatorName)
	if err := s.repo.Create(ctx, problem); err != nil {
		if errors.Is(err, repository.ErrProblemExists) {
			return nil, util.WrapAppError(constants.CodeProblemExists, constants.MsgProblemExists, err)
		}
		s.logger.Error(constants.LogProblemCreatedFailed, "title", req.Title, "error", err.Error())
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogProblemCreated, "problem_id", problem.ID.Hex(), "title", problem.Title)
	resp := dto.ToProblemResponse(problem, true)
	return &resp, nil
}

// Update 更新题目（创建者或管理员）。
func (s *ProblemService) Update(ctx context.Context, problemID, operatorID primitive.ObjectID, operatorRole string, req *dto.ProblemRequest) (*dto.ProblemResponse, error) {
	problem, err := s.repo.FindByID(ctx, problemID)
	if err != nil {
		if errors.Is(err, repository.ErrProblemNotFound) {
			return nil, util.WrapAppError(constants.CodeProblemNotFound, constants.MsgProblemNotFound, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if problem.CreatedBy != operatorID && operatorRole != constants.RoleAdmin {
		return nil, util.WrapAppError(constants.CodeProblemOwner, constants.MsgProblemOwner, nil)
	}
	if len(req.TestCases) == 0 {
		return nil, util.WrapAppError(constants.CodeProblemNoCases, constants.MsgProblemNoCases, nil)
	}
	tcs := make([]model.TestCase, 0, len(req.TestCases))
	for _, tc := range req.TestCases {
		tcs = append(tcs, model.TestCase{Input: tc.Input, Output: tc.Output})
	}
	timeLimit := req.TimeLimit
	if timeLimit <= 0 {
		timeLimit = problem.TimeLimit
	}
	update := bson.M{
		"title":       req.Title,
		"description": req.Description,
		"difficulty":  req.Difficulty,
		"languages":   req.Languages,
		"tags":        req.Tags,
		"test_cases":  tcs,
		"time_limit":  timeLimit,
		"points":      constants.DifficultyPoints[req.Difficulty],
	}
	if err := s.repo.Update(ctx, problemID, update); err != nil {
		if errors.Is(err, repository.ErrProblemExists) {
			return nil, util.WrapAppError(constants.CodeProblemExists, constants.MsgProblemExists, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogProblemUpdated, "problem_id", problemID.Hex(), "title", req.Title)
	return s.Get(ctx, problemID, operatorRole)
}

// UpdateStatus 题目状态机流转（创建者或管理员）。
func (s *ProblemService) UpdateStatus(ctx context.Context, problemID, operatorID primitive.ObjectID, operatorRole, status string) error {
	problem, err := s.repo.FindByID(ctx, problemID)
	if err != nil {
		if errors.Is(err, repository.ErrProblemNotFound) {
			return util.WrapAppError(constants.CodeProblemNotFound, constants.MsgProblemNotFound, err)
		}
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if problem.CreatedBy != operatorID && operatorRole != constants.RoleAdmin {
		return util.WrapAppError(constants.CodeProblemOwner, constants.MsgProblemOwner, nil)
	}
	if !canTransition(problem.Status, status) {
		return util.WrapAppError(constants.CodeProblemLocked, constants.MsgProblemLocked, nil)
	}
	if err := s.repo.UpdateStatus(ctx, problemID, status); err != nil {
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogProblemPublish, "problem_id", problemID.Hex(), "status", status)
	return nil
}

// List 分页查询题目（学生只看已发布）。
func (s *ProblemService) List(ctx context.Context, page, pageSize int64, difficulty, status, tag, role string) ([]dto.ProblemResponse, int64, error) {
	filter := bson.M{}
	if difficulty != "" {
		filter["difficulty"] = difficulty
	}
	if status != "" {
		filter["status"] = status
	} else if role != constants.RoleAdmin {
		filter["status"] = constants.StatusPublished
	}
	if tag != "" {
		filter["tags"] = tag
	}
	problems, total, err := s.repo.List(ctx, filter, (page-1)*pageSize, pageSize)
	if err != nil {
		s.logger.Error(constants.LogProblemListFailed, "error", err.Error())
		return nil, 0, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	withOutput := role == constants.RoleAdmin
	out := make([]dto.ProblemResponse, 0, len(problems))
	for _, p := range problems {
		out = append(out, dto.ToProblemResponse(p, withOutput))
	}
	return out, total, nil
}

// Get 获取题目详情（期望输出仅管理员可见）。
func (s *ProblemService) Get(ctx context.Context, problemID primitive.ObjectID, role string) (*dto.ProblemResponse, error) {
	problem, err := s.repo.FindByID(ctx, problemID)
	if err != nil {
		if errors.Is(err, repository.ErrProblemNotFound) {
			return nil, util.WrapAppError(constants.CodeProblemNotFound, constants.MsgProblemNotFound, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if problem.Status != constants.StatusPublished && role != constants.RoleAdmin {
		return nil, util.WrapAppError(constants.CodeProblemLocked, constants.MsgProblemLocked, nil)
	}
	resp := dto.ToProblemResponse(problem, role == constants.RoleAdmin)
	return &resp, nil
}

// Delete 删除题目（创建者或管理员）。
func (s *ProblemService) Delete(ctx context.Context, problemID, operatorID primitive.ObjectID, operatorRole string) error {
	problem, err := s.repo.FindByID(ctx, problemID)
	if err != nil {
		if errors.Is(err, repository.ErrProblemNotFound) {
			return util.WrapAppError(constants.CodeProblemNotFound, constants.MsgProblemNotFound, err)
		}
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if problem.CreatedBy != operatorID && operatorRole != constants.RoleAdmin {
		return util.WrapAppError(constants.CodeProblemOwner, constants.MsgProblemOwner, nil)
	}
	if err := s.repo.Delete(ctx, problemID); err != nil {
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	return nil
}
