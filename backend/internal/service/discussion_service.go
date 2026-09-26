package service

import (
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
)

// DiscussionService 讨论社区业务：发布/排序/投票/最佳答案。
type DiscussionService struct {
	repo      *repository.DiscussionRepository
	userRepo  *repository.UserRepository
	problemRepo *repository.ProblemRepository
	logger    *slog.Logger
}

// NewDiscussionService 构造讨论服务。
func NewDiscussionService(repo *repository.DiscussionRepository, userRepo *repository.UserRepository,
	problemRepo *repository.ProblemRepository, logger *slog.Logger) *DiscussionService {
	return &DiscussionService{repo: repo, userRepo: userRepo, problemRepo: problemRepo, logger: logger}
}

// Create 发布讨论帖（绑定题目）。
func (s *DiscussionService) Create(ctx context.Context, userID primitive.ObjectID, problemID primitive.ObjectID, req *dto.CreateDiscussionRequest) (*dto.DiscussionResponse, error) {
	if _, err := s.problemRepo.FindByID(ctx, problemID); err != nil {
		return nil, util.WrapAppError(constants.CodeProblemNotFound, constants.MsgProblemNotFound, err)
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeUserNotFound, constants.MsgUserNotFound, err)
	}
	d := &model.Discussion{
		ProblemID: problemID,
		UserID:    userID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Title:     req.Title,
		Content:   req.Content,
		Code:      req.Code,
		Language:  req.Language,
		Status:    constants.DiscussionActive,
		VoterIDs:  []primitive.ObjectID{},
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogDiscussionCreated, "discussion_id", d.ID.Hex(), "problem_id", problemID.Hex(), "username", user.Username)
	resp := dto.ToDiscussionResponse(d)
	return &resp, nil
}

// ListByProblem 查询题目讨论区（排序：best/new/votes）。
func (s *DiscussionService) ListByProblem(ctx context.Context, problemID primitive.ObjectID, sort string) ([]dto.DiscussionResponse, error) {
	list, err := s.repo.ListByProblem(ctx, problemID, sort)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	out := make([]dto.DiscussionResponse, 0, len(list))
	for _, d := range list {
		out = append(out, dto.ToDiscussionResponse(d))
	}
	return out, nil
}

// Vote 点赞/取消点赞（幂等：voter_ids 数组）。
func (s *DiscussionService) Vote(ctx context.Context, userID primitive.ObjectID, discussionID primitive.ObjectID) (*dto.DiscussionResponse, error) {
	d, err := s.repo.FindByID(ctx, discussionID)
	if err != nil {
		if errors.Is(err, repository.ErrDiscussionNotFound) {
			return nil, util.WrapAppError(constants.CodeDiscussionNotFound, constants.MsgDiscussionNotFound, err)
		}
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if d.Status != constants.DiscussionActive {
		return nil, util.WrapAppError(constants.CodeDiscussionLocked, constants.MsgDiscussionLocked, nil)
	}
	voted, err := s.repo.HasVoted(ctx, discussionID, userID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	delta := 1
	if voted {
		delta = -1
	}
	if err := s.repo.Vote(ctx, discussionID, userID, delta); err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogDiscussionVoted, "discussion_id", discussionID.Hex(), "delta", delta)
	d.VoteCount += delta
	resp := dto.ToDiscussionResponse(d)
	return &resp, nil
}

// MarkBest 标记最佳答案（题目创建者或管理员；清除同题旧最佳）。
func (s *DiscussionService) MarkBest(ctx context.Context, operatorID primitive.ObjectID, operatorRole string, problemID, discussionID primitive.ObjectID) error {
	d, err := s.repo.FindByID(ctx, discussionID)
	if err != nil {
		if errors.Is(err, repository.ErrDiscussionNotFound) {
			return util.WrapAppError(constants.CodeDiscussionNotFound, constants.MsgDiscussionNotFound, err)
		}
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	problem, err := s.problemRepo.FindByID(ctx, problemID)
	if err != nil {
		return util.WrapAppError(constants.CodeProblemNotFound, constants.MsgProblemNotFound, err)
	}
	if problem.CreatedBy != operatorID && operatorRole != constants.RoleAdmin {
		return util.WrapAppError(constants.CodeDiscussionDenied, constants.MsgDiscussionDenied, nil)
	}
	if d.ProblemID != problemID {
		return util.WrapAppError(constants.CodeDiscussionDenied, constants.MsgDiscussionDenied, nil)
	}
	if err := s.repo.MarkBest(ctx, problemID, discussionID); err != nil {
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogDiscussionBestMarked, "discussion_id", discussionID.Hex())
	return nil
}

// Hide 隐藏讨论帖（管理员或本人）。
func (s *DiscussionService) Hide(ctx context.Context, operatorID primitive.ObjectID, operatorRole string, discussionID primitive.ObjectID) error {
	d, err := s.repo.FindByID(ctx, discussionID)
	if err != nil {
		if errors.Is(err, repository.ErrDiscussionNotFound) {
			return util.WrapAppError(constants.CodeDiscussionNotFound, constants.MsgDiscussionNotFound, err)
		}
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	if d.UserID != operatorID && operatorRole != constants.RoleAdmin {
		return util.WrapAppError(constants.CodeDiscussionDenied, constants.MsgDiscussionDenied, nil)
	}
	if err := s.repo.UpdateStatus(ctx, discussionID, constants.DiscussionHidden); err != nil {
		return util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	s.logger.Info(constants.LogDiscussionHidden, "discussion_id", discussionID.Hex())
	return nil
}
