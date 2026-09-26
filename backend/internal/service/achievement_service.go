package service

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/blueship581/codelearn/internal/constants"
	"github.com/blueship581/codelearn/internal/dto"
	"github.com/blueship581/codelearn/internal/model"
	"github.com/blueship581/codelearn/internal/repository"
	"github.com/blueship581/codelearn/internal/util"
)

// AchievementService 成就系统：根据签到/解题事件检查并授予徽章（幂等）。
type AchievementService struct {
	achRepo *repository.AchievementRepository
	subRepo *repository.SubmissionRepository
	userRepo *repository.UserRepository
	logger  *slog.Logger
}

// NewAchievementService 构造成就服务。
func NewAchievementService(achRepo *repository.AchievementRepository, subRepo *repository.SubmissionRepository,
	userRepo *repository.UserRepository, logger *slog.Logger) *AchievementService {
	return &AchievementService{achRepo: achRepo, subRepo: subRepo, userRepo: userRepo, logger: logger}
}

// ListAll 查询全部成就定义。
func (s *AchievementService) ListAll(ctx context.Context) ([]dto.AchievementResponse, error) {
	list, err := s.achRepo.ListAll(ctx)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	out := make([]dto.AchievementResponse, 0, len(list))
	for _, a := range list {
		out = append(out, dto.ToAchievementResponse(a))
	}
	return out, nil
}

// ListMine 查询我的成就。
func (s *AchievementService) ListMine(ctx context.Context, userID primitive.ObjectID) ([]dto.AchievementResponse, error) {
	list, err := s.achRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternal, constants.MsgInternalError, err)
	}
	out := make([]dto.AchievementResponse, 0, len(list))
	for _, ua := range list {
		out = append(out, dto.ToUserAchievementResponse(ua))
	}
	return out, nil
}

// CheckAfterSignIn 签到后检查连续签到成就。
func (s *AchievementService) CheckAfterSignIn(ctx context.Context, userID primitive.ObjectID, streak int) {
	if streak >= 7 {
		s.grant(ctx, userID, constants.AchievementStreak7)
	}
}

// CheckAfterSubmission 提交通过后检查解题类成就。
func (s *AchievementService) CheckAfterSubmission(ctx context.Context, userID primitive.ObjectID, status string, problem *model.Problem) {
	if status != constants.SubmissionAccepted {
		return
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return
	}
	s.grant(ctx, userID, constants.AchievementFirstAC)
	if problem.Difficulty == constants.DifficultyHard {
		s.grant(ctx, userID, constants.AchievementFirstHard)
	}
	if user.SolvedCount >= 10 {
		s.grant(ctx, userID, constants.AchievementSolved10)
	}
	if user.SolvedCount >= 100 {
		s.grant(ctx, userID, constants.AchievementSolved100)
	}
}

// grant 授予成就（唯一索引幂等）。
func (s *AchievementService) grant(ctx context.Context, userID primitive.ObjectID, code string) {
	def, err := s.achRepo.FindByCode(ctx, code)
	if err != nil {
		return
	}
	has, err := s.achRepo.HasCode(ctx, userID, code)
	if err != nil || has {
		return
	}
	ua, err := s.achRepo.Grant(ctx, userID, def)
	if err != nil {
		return
	}
	s.logger.Info(constants.LogAchievementGranted, "user_id", userID.Hex(), "code", code, "earned_at", ua.EarnedAt.Format(time.RFC3339))
}
